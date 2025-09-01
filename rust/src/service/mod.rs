pub mod config;

use std::{future::ready, ops::Not, str::FromStr, time::Duration};

use alloy_eips::{BlockId, BlockNumberOrTag};
use alloy_network::TransactionBuilder;
use alloy_primitives::TxHash;
use alloy_provider::{DynProvider, Provider};
use alloy_rpc_types_eth::{EIP1186AccountProofResponse, Transaction, TransactionRequest};
use alloy_signer_local::{LocalSigner, PrivateKeySigner};
use alloy_transport::{RpcError, TransportErrorKind};
use async_trait::async_trait;
use futures::{
    StreamExt, TryStreamExt,
    future::{Either, try_join_all},
    stream,
};
use reth_trie_common::AccountProof;
use serde::{Deserialize, Serialize};
use sov_rollup_interface::{
    common::HexHash,
    da::{BlobReaderTrait, BlockHeaderTrait, DaProof, DaSpec, RelevantBlobs, RelevantProofs, Time},
    node::da::{DaService, SlotData, SubmitBlobReceipt},
};
use thiserror::Error;
use tokio::{sync::oneshot, time::sleep, try_join};
use tracing::{debug, error, instrument, warn};

use crate::{
    eigenda::{
        cert::StandardCommitment,
        extraction::contract,
        proxy::{ProxyClient, ProxyError},
        verification::{verify_cert, verify_cert_recency},
    },
    ethereum::{extract_certificate, provider::init_ethereum_provider, tx::map_eip4844},
    service::config::{EigenDaConfig, EigenDaContracts, Network},
    spec::{
        AncestorMetadata, AncestorStateData, BlobWithSender, EigenDaSpec, EthereumAddress,
        EthereumBlockHeader, EthereumHash, NamespaceId, RollupParams, TransactionWithBlob,
    },
    verifier::{EigenDaCompletenessProof, EigenDaInclusionProof, EigenDaVerifier},
};

/// Possible errors that can happen when using [`EigenDaService`].
#[derive(Debug, Error)]
pub enum EigenDaServiceError {
    #[error("Configuration error: {0}")]
    Configuration(String),

    #[error("Error received from the EigenDA proxy: {0}")]
    ProxyError(#[from] ProxyError),

    #[error("Error received from the Ethereum node: {0}")]
    EthereumRpcError(#[from] RpcError<TransportErrorKind>),

    #[error("Ancestor at height ({0}) is missing")]
    AncestorMissing(u64),
}

/// EigenDaService is responsible for interacting with the EigenDA data availability layer.
/// It provides functionality to submit blobs (data) to EigenDA and interfaces with Ethereum
/// for block information and finality status.
#[derive(Debug, Clone)]
pub struct EigenDaService {
    /// Client for interacting with the EigenDA proxy node
    proxy: ProxyClient,
    /// Provider for interacting with an Ethereum node
    ethereum: DynProvider,
    /// The account to which we are storing the certificates of the batch blobs
    rollup_batch_namespace: NamespaceId,
    /// The account to which we are storing the certificates of the proof blobs
    rollup_proof_namespace: NamespaceId,
    /// Cert recency window
    cert_recency_window: u64,
    /// Private key of the sequencer. It is used to sign the transactions
    /// persisting the certificates to Ethereum
    sequencer_signer: PrivateKeySigner,
    /// EigenDA relevant contracts
    contracts: EigenDaContracts,
}

impl EigenDaService {
    /// Initialize new [`EigenDaService`] with provided [`EigenDaConfig`] and [`RollupParams`].
    pub async fn new(
        config: EigenDaConfig,
        params: RollupParams,
    ) -> Result<Self, EigenDaServiceError> {
        if params.rollup_batch_namespace == params.rollup_proof_namespace {
            return Err(EigenDaServiceError::Configuration(
                "Namespaces should not be equal".to_string(),
            ));
        }

        // Setup ethereum client
        let sequencer_signer = LocalSigner::from_str(&config.sequencer_signer)
            .map_err(|err| EigenDaServiceError::Configuration(err.to_string()))?;
        let ethereum = init_ethereum_provider(&config, sequencer_signer.clone()).await?;

        // Setup proxy client
        let proxy = ProxyClient::new(&config)?;

        // Set contracts
        let contracts = match config.network {
            Network::Mainnet => EigenDaContracts::mainnet(),
            Network::Holesky => EigenDaContracts::holesky(),
        };

        Ok(Self {
            proxy,
            ethereum,
            rollup_batch_namespace: params.rollup_batch_namespace,
            rollup_proof_namespace: params.rollup_proof_namespace,
            cert_recency_window: params.cert_recency_window,
            sequencer_signer,
            contracts,
        })
    }

    /// Submit a blob to the EigenDA.
    #[instrument(skip_all)]
    async fn submit_blob_to_namespace(
        &self,
        blob: &[u8],
        namespace: NamespaceId,
    ) -> Result<SubmitBlobReceipt<EthereumHash>, anyhow::Error> {
        // Submit blob to the EigenDA
        let certificate = self.proxy.store_blob(blob).await?;
        debug!(?certificate, "Certificate was received by EigenDa");

        // Persist certificate to the ethereum
        let da_transaction_id = self.submit_certificate(&certificate, namespace).await?;
        debug!(
            th_hash = %da_transaction_id,
            "Certificate was submitted to Ethereum"
        );
        let blob_hash = HexHash::new(da_transaction_id.into());

        Ok(SubmitBlobReceipt {
            blob_hash,
            da_transaction_id: EthereumHash::from(da_transaction_id),
        })
    }

    /// Submit the certificate to the ethereum
    async fn submit_certificate(
        &self,
        certificate: &StandardCommitment,
        namespace: NamespaceId,
    ) -> Result<TxHash, EigenDaServiceError> {
        let bytes = certificate.to_rlp_bytes();

        let tx = TransactionRequest::default()
            // This is technically not needed. Because the `from` field is
            // automatically filled with the signer's address. We are specifying
            // it explicitly for the readability.
            .with_from(self.sequencer_signer.address())
            .with_to(namespace.into())
            .with_input(bytes);

        let transaction = self.ethereum.send_transaction(tx).await?;

        Ok(*transaction.tx_hash())
    }

    async fn process_transactions_with_metadata(
        &self,
        header: &EthereumBlockHeader,
        transactions: Vec<Transaction>,
    ) -> Result<Vec<(TransactionWithBlob, Option<AncestorMetadata>)>, EigenDaServiceError> {
        let mut block_transactions = Vec::with_capacity(transactions.len());

        for transaction in transactions {
            let tx = map_eip4844(transaction.into_inner());

            // Transaction is not relevant for the rollup. We still need it for
            // later when proving the completeness
            if self.rollup_batch_namespace.contains(&tx).not()
                && self.rollup_proof_namespace.contains(&tx).not()
            {
                block_transactions.push((TransactionWithBlob { tx, blob: None }, None));
                continue;
            }

            // Certificate is malformed
            let Some(cert) = extract_certificate(&tx) else {
                block_transactions.push((TransactionWithBlob { tx, blob: None }, None));
                continue;
            };

            // Verify certificate recency
            if let Err(err) =
                verify_cert_recency(header, cert.reference_block(), self.cert_recency_window)
            {
                warn!(
                    ?header,
                    ?cert,
                    ?err,
                    "Certificate recency verification failed. Ignoring."
                );
                block_transactions.push((
                    TransactionWithBlob { tx, blob: None },
                    // We don't need to store an ancestor to prove that the
                    // certificate recency is invalid.
                    None,
                ));
                continue;
            };

            // Verify the certificate against the ancestor referenced
            let ancestor = self.fetch_referenced_ancestor(&cert).await?;
            if let Err(err) = verify_cert(header, &ancestor, &cert) {
                warn!(
                    ?header,
                    ?cert,
                    ?err,
                    "Certificate verification failed. Ignoring."
                );
                block_transactions.push((
                    TransactionWithBlob { tx, blob: None },
                    // We need the ancestor for the invalid certificate so that
                    // we can prove that the certificate is really invalid and
                    // that was the reason we skipped it.
                    Some(ancestor),
                ));
                continue;
            };

            // The blob should always be available for the valid certificate
            let blob = self.proxy.get_blob(&cert).await?;

            let transaction = TransactionWithBlob {
                tx,
                blob: Some(blob),
            };
            block_transactions.push((transaction, Some(ancestor)));
        }

        Ok(block_transactions)
    }

    /// Fetch [`AncestorMetadata`] with data needed to verify the certificate.
    async fn fetch_referenced_ancestor(
        &self,
        certificate: &StandardCommitment,
    ) -> Result<AncestorMetadata, EigenDaServiceError> {
        let block_height = certificate.reference_block();

        let (mut ancestor, data) = try_join!(
            self.fetch_ancestor(block_height),
            self.fetch_ancestor_state(block_height, certificate)
        )?;
        ancestor.data.replace(data);

        Ok(ancestor)
    }

    /// Fetch [`AncestorMetadata`] only with the header set. The
    /// [`AncestorMetadata::data`] is set as None.
    async fn fetch_ancestor(
        &self,
        block_height: u64,
    ) -> Result<AncestorMetadata, EigenDaServiceError> {
        let block = self
            .ethereum
            .get_block_by_number(block_height.into())
            .await?
            .ok_or_else(|| EigenDaServiceError::AncestorMissing(block_height))?;
        let header = block.header.into_consensus();

        Ok(AncestorMetadata {
            header: EthereumBlockHeader::from(header),
            data: None,
        })
    }

    /// Prepare a valid ancestor chain. Based on the `referenced_ancestors` we
    /// fetch the remaining ancestors needed to have a contiguous chain. The
    /// first ancestor is the earliest referenced ancestor, the last ancestor is
    /// the parent of the block on the `current_height`.
    async fn prepare_ancestor_chain(
        &self,
        current_height: u64,
        referenced_ancestors: Vec<AncestorMetadata>,
    ) -> Result<Vec<AncestorMetadata>, EigenDaServiceError> {
        // If no ancestors, we return an empty chain
        let Some(oldest_ancestor_height) = referenced_ancestors
            .iter()
            .map(|ancestor| ancestor.header.height())
            .min()
        else {
            return Ok(vec![]);
        };

        if oldest_ancestor_height >= current_height {
            warn!(
                ?oldest_ancestor_height,
                %current_height,
                "oldest_ancestor_height is >= current_height"
            );
            return Ok(vec![]);
        }

        let ancestors_fut = (oldest_ancestor_height..current_height).map(|height| {
            if let Some(referenced) = referenced_ancestors
                .iter()
                .find(|ancestor| ancestor.header.height() == height)
            {
                Either::Left(ready(Ok(referenced.clone())))
            } else {
                Either::Right(self.fetch_ancestor(height))
            }
        });

        let ancestors = stream::iter(ancestors_fut)
            .buffered(10)
            .try_collect()
            .await?;

        Ok(ancestors)
    }

    /// Fetches the relevant state used at certificate creation. This state is
    /// later used to verify the EigenDA certificate construction.
    async fn fetch_ancestor_state(
        &self,
        block_height: u64,
        cert: &StandardCommitment,
    ) -> Result<AncestorStateData, EigenDaServiceError> {
        let keys = contract::EigenDaRelayRegistry::storage_keys(cert);
        let eigen_da_relay_registry_fut = self
            .ethereum
            .get_proof(self.contracts.eigen_da_relay_registry, keys)
            .number(block_height)
            .into_future();

        let keys = contract::EigenDaThresholdRegistry::storage_keys(cert);
        let eigen_da_threshold_registry_fut = self
            .ethereum
            .get_proof(self.contracts.eigen_da_threshold_registry, keys)
            .number(block_height)
            .into_future();

        let keys = contract::RegistryCoordinator::storage_keys(cert);
        let registry_coordinator_fut = self
            .ethereum
            .get_proof(self.contracts.registry_coordinator, keys)
            .number(block_height)
            .into_future();

        let keys = contract::BlsSignatureChecker::storage_keys(cert);
        let bls_signature_checker_fut = self
            .ethereum
            .get_proof(self.contracts.bls_signature_checker, keys)
            .number(block_height)
            .into_future();

        let keys = contract::DelegationManager::storage_keys(cert);
        let delegation_manager_fut = self
            .ethereum
            .get_proof(self.contracts.delegation_manager, keys)
            .number(block_height)
            .into_future();

        let keys = contract::BlsApkRegistry::storage_keys(cert);
        let bls_apk_registry_fut = self
            .ethereum
            .get_proof(self.contracts.bls_apk_registry, keys)
            .number(block_height)
            .into_future();

        let keys = contract::StakeRegistry::storage_keys(cert);
        let stake_registry_fut = self
            .ethereum
            .get_proof(self.contracts.stake_registry, keys)
            .number(block_height)
            .into_future();

        let keys = contract::EigenDaCertVerifier::storage_keys(cert);
        let eigen_da_cert_verifier_fut = self
            .ethereum
            .get_proof(self.contracts.eigen_da_cert_verifier, keys)
            .number(block_height)
            .into_future();

        let responses = try_join_all([
            eigen_da_relay_registry_fut,
            eigen_da_threshold_registry_fut,
            registry_coordinator_fut,
            bls_signature_checker_fut,
            delegation_manager_fut,
            bls_apk_registry_fut,
            stake_registry_fut,
            eigen_da_cert_verifier_fut,
        ])
        .await?;

        let [
            eigen_da_relay_registry,
            eigen_da_threshold_registry,
            registry_coordinator,
            bls_signature_checker,
            delegation_manager,
            bls_apk_registry,
            stake_registry,
            eigen_da_cert_verifier,
        ]: [EIP1186AccountProofResponse; 8] = responses.try_into().expect("Expected 8 elements");

        Ok(AncestorStateData::new(
            AccountProof::from(eigen_da_relay_registry),
            AccountProof::from(eigen_da_threshold_registry),
            AccountProof::from(registry_coordinator),
            AccountProof::from(bls_signature_checker),
            AccountProof::from(delegation_manager),
            AccountProof::from(bls_apk_registry),
            AccountProof::from(stake_registry),
            AccountProof::from(eigen_da_cert_verifier),
        ))
    }
}

#[async_trait]
impl DaService for EigenDaService {
    /// A handle to the types used by the DA layer.
    type Spec = EigenDaSpec;

    /// [`serde`]-compatible configuration data for this [`DaService`]. Parsed
    /// from TOML.
    type Config = EigenDaConfig;

    /// The verifier for this DA layer.
    type Verifier = EigenDaVerifier;

    /// A DA layer block, possibly excluding some irrelevant information.
    type FilteredBlock = EthereumBlock;

    /// The error type for fallible methods.
    type Error = anyhow::Error;

    /// Fetch the block at the given height, waiting for one to be mined if necessary.
    ///
    /// The returned block may not be final, and can be reverted without a consensus violation.
    /// Calls to this method for the same height are allowed to return different results.
    /// Should always return the block at that height on the best fork.
    async fn get_block_at(&self, height: u64) -> Result<Self::FilteredBlock, Self::Error> {
        let number = BlockNumberOrTag::Number(height);

        // Poll until the requested block is mined
        let poll_interval = Duration::from_secs(10);
        let block = loop {
            match self.ethereum.get_block_by_number(number).full().await {
                Ok(Some(block)) => break block,
                Ok(None) => {
                    sleep(poll_interval).await;
                    continue;
                }
                Err(err) => {
                    warn!(?err, "error occurred while getting the block");
                    continue;
                }
            }
        };

        let header = EthereumBlockHeader::from(block.header.clone().into_consensus());

        // Iterate over transactions in the block and fetch sequencer relevant data
        let transactions_with_ancestors = self
            .process_transactions_with_metadata(&header, block.into_transactions_vec())
            .await?;
        let (transactions, ancestors): (_, Vec<_>) =
            transactions_with_ancestors.into_iter().unzip();

        let ancestors = ancestors.into_iter().flatten().collect();
        let ancestors = self.prepare_ancestor_chain(height, ancestors).await?;

        Ok(EthereumBlock {
            ancestors,
            header,
            transactions,
        })
    }

    /// Fetch the [`DaSpec::BlockHeader`] of the last finalized block.
    /// If there's no finalized block yet, it should return an error.
    async fn get_last_finalized_block_header(
        &self,
    ) -> Result<<Self::Spec as DaSpec>::BlockHeader, Self::Error> {
        let block = BlockId::finalized();
        let block = self
            .ethereum
            .get_block(block)
            .await?
            .ok_or_else(|| anyhow::anyhow!("No finalized block"))?;
        let header = block.header.into_consensus();

        Ok(EthereumBlockHeader::from(header))
    }

    /// Fetch the head block of the most popular fork.
    ///
    /// More like utility method, to provide better user experience
    async fn get_head_block_header(
        &self,
    ) -> Result<<Self::Spec as DaSpec>::BlockHeader, Self::Error> {
        let block = BlockId::latest();
        let block = self
            .ethereum
            .get_block(block)
            .await?
            .ok_or_else(|| anyhow::anyhow!("No finalized block"))?;
        let header = block.header.into_consensus();

        Ok(EthereumBlockHeader::from(header))
    }

    /// Extract the relevant transactions from a block.
    fn extract_relevant_blobs(
        &self,
        block: &Self::FilteredBlock,
    ) -> RelevantBlobs<<Self::Spec as DaSpec>::BlobTransaction> {
        let proof_blobs = block.extract_relevant_blobs(self.rollup_proof_namespace);
        let batch_blobs = block.extract_relevant_blobs(self.rollup_batch_namespace);

        RelevantBlobs {
            proof_blobs,
            batch_blobs,
        }
    }

    /// Generate a proof that the relevant blob transactions have been extracted correctly from the DA layer
    /// block.
    async fn get_extraction_proof(
        &self,
        block: &Self::FilteredBlock,
        _blobs: &RelevantBlobs<<Self::Spec as DaSpec>::BlobTransaction>,
    ) -> RelevantProofs<
        <Self::Spec as DaSpec>::InclusionMultiProof,
        <Self::Spec as DaSpec>::CompletenessProof,
    > {
        let proof = block.get_extraction_proof(self.rollup_proof_namespace);
        let batch = block.get_extraction_proof(self.rollup_batch_namespace);

        RelevantProofs { proof, batch }
    }

    /// Send a transaction directly to the DA layer.
    /// This method is infallible: the SubmitBlobReceipt is returned via the `oneshot::Receiver` after the blob is posted to the DA.
    async fn send_transaction(
        &self,
        blob: &[u8],
    ) -> oneshot::Receiver<
        Result<SubmitBlobReceipt<<Self::Spec as DaSpec>::TransactionId>, Self::Error>,
    > {
        let (tx, rx) = oneshot::channel();
        let result = self
            .submit_blob_to_namespace(blob, self.rollup_batch_namespace)
            .await;
        tx.send(result).expect("receiver exists");

        rx
    }

    /// Sends a proof to the DA layer.
    /// This method is infallible: the SubmitBlobReceipt is returned via the `oneshot::Receiver` after the blob is posted to the DA.
    async fn send_proof(
        &self,
        aggregated_proof_data: &[u8],
    ) -> oneshot::Receiver<
        Result<SubmitBlobReceipt<<Self::Spec as DaSpec>::TransactionId>, Self::Error>,
    > {
        let (tx, rx) = oneshot::channel();
        let result = self
            .submit_blob_to_namespace(aggregated_proof_data, self.rollup_proof_namespace)
            .await;
        tx.send(result).expect("receiver exists");

        rx
    }

    /// Fetches all proofs at a specified block height.
    async fn get_proofs_at(&self, height: u64) -> Result<Vec<Vec<u8>>, Self::Error> {
        let block = self.get_block_at(height).await?;
        let blobs = block.extract_relevant_blobs(self.rollup_proof_namespace);
        let proofs = blobs
            .into_iter()
            .map(|mut b| b.full_data().to_vec())
            .collect::<Vec<Vec<u8>>>();

        Ok(proofs)
    }

    /// Returns a [`DaSpec::Address`] that signs blobs submitted by this instance of [`DaService`]
    async fn get_signer(&self) -> <Self::Spec as DaSpec>::Address {
        EthereumAddress::from(self.sequencer_signer.address())
    }
}

/// An Ethereum block containing relevant information.
#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct EthereumBlock {
    /// List of ancestor blocks. The first ancestor in a list is the earliest
    /// reference block from which the values for certificate creation were
    /// sourced. The last header in the list is a parent of this block.
    pub ancestors: Vec<AncestorMetadata>,
    /// The current block header.
    pub header: EthereumBlockHeader,
    /// Transactions included in this block.
    pub transactions: Vec<TransactionWithBlob>,
}

impl EthereumBlock {
    /// Extract all rollup's blobs (indicated by namespace) from this block.
    pub fn extract_relevant_blobs(&self, namespace: NamespaceId) -> Vec<BlobWithSender> {
        self.transactions
            .iter()
            .filter_map(|tx| {
                namespace
                    .contains(&tx.tx)
                    .then(|| tx.blob.clone().map(|blob| (&tx.tx, blob)))
                    .flatten()
            })
            .filter_map(|(tx, blob)| {
                let sender = tx.recover_signer().ok()?;
                let tx_hash = tx.hash().to_owned();

                Some(BlobWithSender::new(sender, tx_hash, blob))
            })
            .collect::<Vec<_>>()
    }

    /// Get the inclusion and completeness proofs pair for rollup's blobs (indicated by namespace)
    /// contained in this block.
    pub fn get_extraction_proof(
        &self,
        namespace: NamespaceId,
    ) -> DaProof<EigenDaInclusionProof, EigenDaCompletenessProof> {
        let maybe_relevant_txs = self
            .transactions
            .iter()
            .filter(|tx| namespace.contains(&tx.tx))
            .cloned()
            .collect();

        DaProof {
            inclusion_proof: EigenDaInclusionProof::new(maybe_relevant_txs),
            completeness_proof: EigenDaCompletenessProof::new(
                self.ancestors.clone(),
                self.transactions.clone(),
            ),
        }
    }
}

impl SlotData for EthereumBlock {
    type BlockHeader = EthereumBlockHeader;

    fn hash(&self) -> [u8; 32] {
        self.header.hash().into()
    }

    fn header(&self) -> &Self::BlockHeader {
        &self.header
    }

    fn timestamp(&self) -> Time {
        self.header.time()
    }
}
