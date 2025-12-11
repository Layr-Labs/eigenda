use alloy_primitives::B256;
use reth_trie_common::AccountProof;
use reth_trie_common::proof::ProofVerificationError;
use serde::{Deserialize, Serialize};
use thiserror::Error;
use tracing::instrument;

use crate::cert::StandardCommitment;
use crate::extraction::extractor::{
    ApkHistoryExtractor, DataDecoder, MinWithdrawalDelayBlocksExtractor, NextBlobVersionExtractor,
    OperatorBitmapHistoryExtractor, OperatorStakeHistoryExtractor, QuorumCountExtractor,
    QuorumNumbersRequiredV2Extractor, QuorumUpdateBlockNumberExtractor,
    SecurityThresholdsV2Extractor, StaleStakesForbiddenExtractor, TotalStakeHistoryExtractor,
    VersionedBlobParamsExtractor,
};
use crate::verification::cert::types::history::HistoryError;
use crate::verification::cert::types::{Staleness, Storage};
use crate::verification::cert::{Cert, CertVerificationInputs};

/// Contract-specific extraction logic and storage key generators.
pub mod contract;

/// Helper functions for decoding contract storage data.
pub mod decode_helpers;

/// Core extraction traits and implementations for certificate data.
pub mod extractor;

/// Utilities for generating Ethereum contract storage keys.
pub mod storage_key_helpers;

/// Errors that can occur during certificate data extraction
#[derive(Debug, Error, PartialEq)]
pub enum CertExtractionError {
    /// Storage proof was not found for the requested variable
    #[error("Failed to extract StorageProof for {0}")]
    MissingStorageProof(String),

    /// Error from history data processing
    #[error(transparent)]
    HistoryError(#[from] HistoryError),

    /// Error from Alloy Solidity types decoding
    #[error(transparent)]
    AlloySolTypesError(#[from] alloy_sol_types::Error),

    /// Error for when Ethereum Bytes are expected to be encoded in short form but long form is found instead
    #[error("Unexpected ethereum bytes long form")]
    UnexpectedEthereumBytesLongForm,
}

/// Contains data needed to validate the certificate. It also contains proofs
/// used to verify the data.
///
/// AccountProof values both verify storage proofs and carry the raw slots we later decode.
/// Verification and data extraction happen on separate call paths, so we keep this struct as a
/// standalone carrier instead of hiding it inside one helper function.
/// Parsing up-front may be wasteful since proving does not need the data and failure would
/// mean we parsed prematurely.
#[derive(Clone, Debug, PartialEq, Serialize, Deserialize, Default)]
pub struct CertStateData {
    /// Proof for threshold registry contract state.
    pub threshold_registry: AccountProof,
    /// Proof for registry coordinator contract state.
    pub registry_coordinator: AccountProof,
    /// Proof for service manager contract state.
    pub service_manager: AccountProof,
    /// Proof for BLS aggregate public key registry contract state.
    pub bls_apk_registry: AccountProof,
    /// Proof for stake registry contract state.
    pub stake_registry: AccountProof,
    /// Proof for delegation manager contract state.
    pub delegation_manager: AccountProof,
    /// Proof for cert verifier router contract state.
    pub cert_verifier_router: AccountProof,
    /// Proof for certificate verifier contract state.
    pub cert_verifier: AccountProof,
}

impl CertStateData {
    #![allow(clippy::result_large_err)]
    /// Verify all contract state proofs against the given state root.
    pub fn verify(&self, state_root: B256) -> Result<(), ProofVerificationError> {
        self.threshold_registry.verify(state_root)?;
        self.registry_coordinator.verify(state_root)?;
        self.service_manager.verify(state_root)?;
        self.bls_apk_registry.verify(state_root)?;
        self.stake_registry.verify(state_root)?;
        self.delegation_manager.verify(state_root)?;

        self.cert_verifier_router.verify(state_root)?;
        self.cert_verifier.verify(state_root)?;
        // TODO(samlaf): verify that the cert_verifier matches the expected ABN from the router
        Ok(())
    }

    /// Extract certificate verification inputs from contract state data.
    ///
    /// Decodes all required contract storage data from the proofs to construct
    /// verification inputs for certificate validation.
    ///
    /// # Arguments
    /// * `cert` - The certificate to extract data for
    /// * `current_block` - Current block height for verification context
    ///
    /// # Returns
    /// [`CertVerificationInputs`] containing all data needed for certificate verification
    ///
    /// # Errors
    /// Returns [`CertExtractionError`] if:
    /// - Storage proofs are missing for required contract variables
    /// - Data decoding fails
    /// - Historical data is inconsistent
    ///
    /// # Safety
    /// The data extracted is not cryptographically verified. To verify the data,
    /// ensure that [`CertStateData::verify`] is called before extraction.
    #[instrument(skip_all)]
    pub fn extract(
        &self,
        cert: &StandardCommitment,
        current_block: u32,
    ) -> Result<CertVerificationInputs, CertExtractionError> {
        let quorum_count =
            QuorumCountExtractor::new().decode_data(&self.registry_coordinator.storage_proofs)?;

        let quorum_bitmap_history = OperatorBitmapHistoryExtractor::new(cert)
            .decode_data(&self.registry_coordinator.storage_proofs)?;

        let operator_stake_history = OperatorStakeHistoryExtractor::new(cert)
            .decode_data(&self.stake_registry.storage_proofs)?;

        let total_stake_history = TotalStakeHistoryExtractor::new(cert)
            .decode_data(&self.stake_registry.storage_proofs)?;

        let apk_history =
            ApkHistoryExtractor::new(cert).decode_data(&self.bls_apk_registry.storage_proofs)?;

        let versioned_blob_params = VersionedBlobParamsExtractor::new(cert)
            .decode_data(&self.threshold_registry.storage_proofs)?;

        let next_blob_version =
            NextBlobVersionExtractor::new().decode_data(&self.threshold_registry.storage_proofs)?;

        let staleness = {
            let stale_stakes_forbidden = StaleStakesForbiddenExtractor::new()
                .decode_data(&self.service_manager.storage_proofs)?;

            let min_withdrawal_delay_blocks = MinWithdrawalDelayBlocksExtractor::new()
                .decode_data(&self.delegation_manager.storage_proofs)?;

            let quorum_update_block_number = QuorumUpdateBlockNumberExtractor::new(cert)
                .decode_data(&self.registry_coordinator.storage_proofs)?;

            Staleness {
                stale_stakes_forbidden,
                min_withdrawal_delay_blocks,
                quorum_update_block_number,
            }
        };

        let security_thresholds =
            SecurityThresholdsV2Extractor::new().decode_data(&self.cert_verifier.storage_proofs)?;

        let required_quorum_numbers = QuorumNumbersRequiredV2Extractor::new()
            .decode_data(&self.cert_verifier.storage_proofs)?;

        let storage = Storage {
            quorum_count,
            current_block,
            quorum_bitmap_history,
            operator_stake_history,
            total_stake_history,
            apk_history,
            versioned_blob_params,
            next_blob_version,
            security_thresholds,
            required_quorum_numbers,
            staleness,
        };

        let cert = Cert {
            batch_header: cert.batch_header_v2().clone(),
            blob_inclusion_info: cert.blob_inclusion_info().clone(),
            non_signer_stakes_and_signature: cert.nonsigner_stake_and_signature().clone(),
            signed_quorum_numbers: cert.signed_quorum_numbers().clone(),
        };

        let inputs = CertVerificationInputs { cert, storage };

        Ok(inputs)
    }
}
