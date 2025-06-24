use crate::spec::EigenDaSpec;
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

#[derive(Debug, Serialize, Deserialize)]
pub struct EigenDaInclusionProof;
