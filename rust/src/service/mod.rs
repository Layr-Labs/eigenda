mod proxy;

use alloy::{
    eips::{BlockId, BlockNumberOrTag},
    providers::{Provider, RootProvider},
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
use tracing::{debug, instrument};

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
}

/// Possible errors that can happen when using [`EigenDaService`].
#[derive(Debug, Error)]
pub enum EigenDaServiceError {
    #[error("ProxyError: {0}")]
    ProxyError(#[from] ProxyError),

    #[error("EthereumRpcError: {0}")]
    EthereumRpcError(#[from] RpcError<TransportErrorKind>),
}

/// EigenDaService is responsible for interacting with the EigenDA data availability layer.
/// It provides functionality to submit blobs (data) to EigenDA and interfaces with Ethereum
/// for block information and finality status.
#[derive(Clone)]
pub struct EigenDaService {
    /// Client for interacting with the EigenDA proxy node
    proxy: ProxyClient,
    /// Provider for interacting with an Ethereum node
    ethereum: RootProvider,
}

impl EigenDaService {
    /// Initialize new [`EigenDaService`] with provided [`EigenDaConfig`].
    pub async fn new(config: EigenDaConfig) -> Result<Self, EigenDaServiceError> {
        let client = ProxyClient::new(config.proxy_url)?;
        let ethereum = RootProvider::connect(&config.ethereum_rpc_url).await?;

        Ok(Self {
            proxy: client,
            ethereum,
        })
    }
}

impl EigenDaService {
    /// Submit a blob to the EigenDA.
    #[instrument(skip_all)]
    async fn submit_blob(
        &self,
        blob: &[u8],
        // TODO: Namespace the blob being submitted
    ) -> Result<SubmitBlobReceipt<EthereumHash>, anyhow::Error> {
        // Submit blob to the EigenDA
        let certificate = self.proxy.store_blob(blob).await?;
        debug!(?certificate, "Certificate was received");

        // TODO: What should be used as a blob_hash?
        let blob_hash = HexHash::new(todo!());

        // TODO: The da_transaction_id is the transaction id on the L1 in which
        // the blob was submitted in. In our case that should be a transaction
        // in which the certificate was persisted on chain.
        let da_transaction_id = todo!();

        Ok(SubmitBlobReceipt {
            blob_hash,
            da_transaction_id,
        })
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
        let block = self.ethereum.get_block_by_number(number).await?;

        // TODO: If we receive an Option::None. We should wait for the block to
        // be mined.
        let block = block.unwrap();

        Ok(EthereumBlock {
            header: EthereumBlockHeader::try_from(block.header)?,
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
        todo!()
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
        let result = self.submit_blob(blob).await;
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
        let result = self.submit_blob(aggregated_proof_data).await;
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
