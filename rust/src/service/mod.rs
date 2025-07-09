mod ethereum;
mod proxy;

use std::str::FromStr;
use std::time::Duration;

use alloy::network::TransactionBuilder4844;
use alloy::providers::{DynProvider, ProviderBuilder};
use alloy::rpc::types::Transaction;
use alloy::signers::local::{LocalSigner, LocalSignerError, PrivateKeySigner};
use alloy::{
    consensus::{SidecarBuilder, SimpleCoder},
    eips::{BlockId, BlockNumberOrTag},
    network::TransactionBuilder,
    providers::Provider,
    rpc::types::TransactionRequest,
    transports::{RpcError, TransportErrorKind},
};
use async_trait::async_trait;
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
use sov_rollup_interface::{
    common::HexHash,
    da::{BlockHeaderTrait, DaSpec, RelevantBlobs, RelevantProofs, Time},
    node::da::{DaService, SlotData, SubmitBlobReceipt},
};
use thiserror::Error;
use tokio::sync::oneshot;
use tokio::time::sleep;
use tracing::{debug, instrument, warn};

use crate::spec::{BlobWithSender, NamespaceId, RollupParams};
use crate::{
    service::proxy::{ProxyClient, ProxyError},
    spec::{EigenDaSpec, EthereumBlockHeader, EthereumHash},
    verifier::EigenDaVerifier,
};

/// Configuration for the [`EigenDaService`].
#[derive(Debug, JsonSchema, PartialEq)]
pub struct EigenDaConfig {
    /// URL of the Ethereum RPC node
    pub ethereum_rpc_url: String,
    /// URL of the EigenDA proxy node
    pub proxy_url: String,
    /// Private key of the sequencer. The account with corresponding private key
    /// is used by the sequencer to persist the certificates to Ethereum.
    /// Expected private key in the HEX format.
    pub sequencer_signer: String,
}

/// Possible errors that can happen when using [`EigenDaService`].
#[derive(Debug, Error)]
pub enum EigenDaServiceError {
    #[error("ProxyError: {0}")]
    ProxyError(#[from] ProxyError),

    #[error("EthereumRpcError: {0}")]
    EthereumRpcError(#[from] RpcError<TransportErrorKind>),

    #[error("LocalSignerError: {0}")]
    LocalSignerError(#[from] LocalSignerError),
}

/// EigenDaService is responsible for interacting with the EigenDA data availability layer.
/// It provides functionality to submit blobs (data) to EigenDA and interfaces with Ethereum
/// for block information and finality status.
#[derive(Clone)]
pub struct EigenDaService {
    /// Client for interacting with the EigenDA proxy node
    proxy: ProxyClient,
    /// Provider for interacting with an Ethereum node
    ethereum: DynProvider,
    /// The account to which we are storing the certificates of the batch blobs
    rollup_batch_namespace: NamespaceId,
    /// The account to which we are storing the certificates of the proof blobs
    rollup_proof_namespace: NamespaceId,
    /// Private key of the sequencer. It is used to sign the transactions
    /// persisting the certificates to Ethereum
    sequencer_signer: PrivateKeySigner,
}

impl EigenDaService {
    /// Initialize new [`EigenDaService`] with provided [`EigenDaConfig`] and [`RollupParams`].
    pub async fn new(
        config: EigenDaConfig,
        params: RollupParams,
    ) -> Result<Self, EigenDaServiceError> {
        let client = ProxyClient::new(config.proxy_url)?;

        let sequencer_signer = LocalSigner::from_str(&config.sequencer_signer)?;
        let ethereum = ProviderBuilder::new()
            .wallet(sequencer_signer.clone())
            .connect(&config.ethereum_rpc_url)
            .await?
            .erased();

        Ok(Self {
            proxy: client,
            ethereum,
            rollup_batch_namespace: params.rollup_batch_account,
            rollup_proof_namespace: params.rollup_proof_account,
            sequencer_signer,
        })
    }
}

impl EigenDaService {
    /// Submit a blob to the EigenDA.
    #[instrument(skip_all)]
    async fn submit_blob_to_namespace(
        &self,
        blob: &[u8],
        namespace: NamespaceId,
    ) -> Result<SubmitBlobReceipt<EthereumHash>, anyhow::Error> {
        // Submit blob to the EigenDA
        let certificate = self.proxy.store_blob(blob).await?;
        debug!(?certificate, "Certificate was received");

        // Persist certificate to the ethereum
        let da_transaction_id = self.submit_certificate(&certificate, namespace).await?;
        // TODO: Check how should the blob_hash be actually computed. Is it ok
        // to use the transaction id?
        let blob_hash = HexHash::new(da_transaction_id.into());

        Ok(SubmitBlobReceipt {
            blob_hash,
            da_transaction_id,
        })
    }

    /// Submit the certificate to the ethereum
    async fn submit_certificate(
        &self,
        certificate: &[u8],
        namespace: NamespaceId,
    ) -> Result<EthereumHash, EigenDaServiceError> {
        let sidecar = SidecarBuilder::<SimpleCoder>::from_slice(certificate);
        let sidecar = sidecar.build().unwrap();

        let tx = TransactionRequest::default()
            // This is technically not needed. Because the `from` field is
            // automatically filled with the signer's address. We are specifying
            // it explicitly for the readability.
            .with_from(self.sequencer_signer.address())
            .with_to(namespace.into())
            .with_blob_sidecar(sidecar);

        let transaction = self.ethereum.send_transaction(tx).await?;
        Ok(transaction.tx_hash().to_owned().into())
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
    /// Should always returns the block at that height on the best fork.
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

        // TODO: Retrieve blobs from the beacon chain and append them to
        // transactions.
        let transactions = block
            .transactions
            .into_transactions()
            .map(|transaction| EthereumTransactionWithBlob {
                transaction,
                blob: todo!(),
            })
            .collect();

        Ok(EthereumBlock {
            header: EthereumBlockHeader::try_from(block.header)?,
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

        Ok(EthereumBlockHeader::try_from(block.header)?)
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

        Ok(EthereumBlockHeader::try_from(block.header)?)
    }

    /// Extract the relevant transactions from a block. For example, this method might return
    /// all of the blob transactions from a set rollup namespaces on Celestia.
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
        blobs: &RelevantBlobs<<Self::Spec as DaSpec>::BlobTransaction>,
    ) -> RelevantProofs<
        <Self::Spec as DaSpec>::InclusionMultiProof,
        <Self::Spec as DaSpec>::CompletenessProof,
    > {
        todo!()
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
        todo!()
    }
}

/// An Ethereum block containing only relevant information.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize)]
pub struct EthereumBlock {
    header: EthereumBlockHeader,
    transactions: Vec<EthereumTransactionWithBlob>,
}

impl EthereumBlock {
    /// Extract all rollup's blobs (indicated by namespace) from this block.
    pub fn extract_relevant_blobs(&self, namespace: NamespaceId) -> Vec<BlobWithSender> {
        self.transactions
            .iter()
            .filter_map(|tx: &EthereumTransactionWithBlob| {
                let recovered = tx.transaction.as_recovered();
                let address = recovered.signer();
                let tx_hash = recovered.hash().to_owned();
                let blob = tx.blob.clone();

                namespace
                    .contains(&tx.transaction)
                    .then_some(BlobWithSender::new(address, tx_hash, blob))
            })
            .collect::<Vec<_>>()
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

#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize)]
struct EthereumTransactionWithBlob {
    pub transaction: Transaction,
    pub blob: Vec<u8>,
}

#[cfg(test)]
mod tests {
    use std::str::FromStr;

    use sov_rollup_interface::{da::DaVerifier, node::da::DaService};

    use crate::{
        service::{
            EigenDaConfig, EigenDaService, EigenDaServiceError,
            ethereum::tests::{MiningKind, mine_block, start_ethereum_dev_node},
            proxy::tests::start_proxy,
        },
        spec::{NamespaceId, RollupParams},
        verifier::EigenDaVerifier,
    };

    #[tokio::test]
    async fn submit_extract_verify_e2e() {
        let (proxy_url, _proxy_container) = start_proxy().await.unwrap();
        let (ethereum_rpc_url, _anvil_container) =
            start_ethereum_dev_node(MiningKind::Manual).await.unwrap();

        let service = setup_service(proxy_url, ethereum_rpc_url.clone())
            .await
            .unwrap();
        let verifier = EigenDaVerifier::new(RollupParams {
            rollup_batch_account: service.rollup_batch_namespace,
            rollup_proof_account: service.rollup_proof_namespace,
        });

        let blobs = [vec![123; 123], vec![15; 45], vec![8; 1234], vec![2; 1]];
        let proofs = [vec![43; 87], vec![112; 135], vec![2; 994], vec![1; 1]];

        // Post the rollup data to the network
        for blob in blobs {
            service
                .send_transaction(&blob)
                .await
                .await
                .unwrap()
                .unwrap();
        }
        for proof in proofs {
            service.send_proof(&proof).await.await.unwrap().unwrap();
        }

        // Mine the block
        mine_block(&ethereum_rpc_url, 1).await.unwrap();

        // Extract rollup data from the block
        let block_height = 1;
        let block = service.get_block_at(block_height).await.unwrap();
        let blobs = service.extract_relevant_blobs(&block);
        let proofs = service.get_extraction_proof(&block, &blobs).await;

        // Simulate we're sending this to zkvm
        let header = risc0_zkvm::serde::to_vec(&block.header).unwrap();
        let blobs = risc0_zkvm::serde::to_vec(&blobs).unwrap();
        let proofs = risc0_zkvm::serde::to_vec(&proofs).unwrap();

        // Receive on zkvm side
        let header = risc0_zkvm::serde::from_slice(&header).unwrap();
        let blobs = risc0_zkvm::serde::from_slice(&blobs).unwrap();
        let proofs = risc0_zkvm::serde::from_slice(&proofs).unwrap();

        // Verify
        verifier
            .verify_relevant_tx_list(&header, &blobs, proofs)
            .unwrap();
    }

    async fn setup_service(
        proxy_url: String,
        ethereum_rpc_url: String,
    ) -> Result<EigenDaService, EigenDaServiceError> {
        // ! These keys should not be used in production. They are the keys used
        // by Anvil node for the predefined accounts. !
        let sequencer_signer =
            "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80".to_string();
        let rollup_batch_account =
            NamespaceId::from_str("0x70997970C51812dc3A010C7d01b50e0d17dc79C8").unwrap();
        let rollup_proof_account =
            NamespaceId::from_str("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC").unwrap();

        let config = EigenDaConfig {
            ethereum_rpc_url,
            proxy_url,
            sequencer_signer,
        };
        let params = RollupParams {
            rollup_batch_account,
            rollup_proof_account,
        };

        EigenDaService::new(config, params).await
    }
}
