use crate::spec::{EigenDaSpec, NamespaceId, TransactionWithBlob};
use serde::{Deserialize, Serialize};
use sov_rollup_interface::da::{DaSpec, DaVerifier, RelevantBlobs, RelevantProofs};
use thiserror::Error;

/// Errors that may occur when verifying with the [`EigenDaVerifier`].
#[derive(Debug, Error)]
pub enum VerifierError {}

#[derive(Debug, Clone)]
pub struct EigenDaVerifier;

impl DaVerifier for EigenDaVerifier {
    type Spec = EigenDaSpec;

    type Error = VerifierError;

    fn new(params: <Self::Spec as DaSpec>::ChainParams) -> Self {
        EigenDaVerifier
    }

    fn verify_relevant_tx_list(
        &self,
        block_header: &<Self::Spec as DaSpec>::BlockHeader,
        relevant_blobs: &RelevantBlobs<<Self::Spec as DaSpec>::BlobTransaction>,
        relevant_proofs: RelevantProofs<
            <Self::Spec as DaSpec>::InclusionMultiProof,
            <Self::Spec as DaSpec>::CompletenessProof,
        >,
    ) -> Result<(), Self::Error> {
        todo!()
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct EigenDaCompletenessProof;

impl EigenDaCompletenessProof {
    /// Create a new completeness proof from a complete and ordered list of
    /// transactions that can be relevant for the rollup.
    pub fn new(maybe_relevant_txs: Vec<TransactionWithBlob>) -> Self {
        todo!()
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct EigenDaInclusionProof;

impl EigenDaInclusionProof {
    /// Create a new inclusion proof from a namespace and complete and ordered
    /// set of transactions in the block.
    pub fn new(namespace: NamespaceId, transactions: &[TransactionWithBlob]) -> Self {
        todo!()
    }
}
