use alloy_consensus::{EthereumTxEnvelope, Transaction, TxEip4844};
use alloy_primitives::B256;
use eigenda_verification::cert::StandardCommitment;
use eigenda_verification::verification::cert::CertVerificationInputs;
use eigenda_verification::verification::cert::types::Storage;
use eigenda_verification::verification::cert::types::history::HistoryError;
use reth_trie_common::AccountProof;
use reth_trie_common::proof::ProofVerificationError;
use serde::{Deserialize, Serialize};
use thiserror::Error;
use tracing::instrument;

use crate::extraction::extractor::{
    ApkHistoryExtractor, DataDecoder, NextBlobVersionExtractor, OperatorBitmapHistoryExtractor,
    OperatorStakeHistoryExtractor, QuorumCountExtractor, QuorumNumbersRequiredV2Extractor,
    SecurityThresholdsV2Extractor, TotalStakeHistoryExtractor, VersionedBlobParamsExtractor,
};
#[cfg(feature = "stale-stakes-forbidden")]
use crate::extraction::extractor::{
    MinWithdrawalDelayBlocksExtractor, QuorumUpdateBlockNumberExtractor,
    StaleStakesForbiddenExtractor,
};

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
    WrapHistoryError(#[from] HistoryError),

    /// Error from Alloy Solidity types decoding
    #[error(transparent)]
    WrapAlloySolTypesError(#[from] alloy_sol_types::Error),

    /// Error for when Ethereum Bytes are expected to be encoded in short form but long form is found instead
    #[error("Unexpected ethereum bytes long form")]
    UnexpectedEthereumBytesLongForm,
}

/// Contains data needed to validate the certificate. It also contains proofs
/// used to verify the data.
#[derive(Clone, Debug, PartialEq, Serialize, Deserialize, Default)]
pub struct CertStateData {
    /// Proof for threshold registry contract state.
    pub threshold_registry: AccountProof,
    /// Proof for registry coordinator contract state.
    pub registry_coordinator: AccountProof,
    #[cfg(feature = "stale-stakes-forbidden")]
    /// Proof for service manager contract state.
    pub service_manager: AccountProof,
    /// Proof for BLS aggregate public key registry contract state.
    pub bls_apk_registry: AccountProof,
    /// Proof for stake registry contract state.
    pub stake_registry: AccountProof,
    /// Proof for certificate verifier contract state.
    pub cert_verifier: AccountProof,
    #[cfg(feature = "stale-stakes-forbidden")]
    /// Proof for delegation manager contract state.
    pub delegation_manager: AccountProof,
}

impl CertStateData {
    #![allow(clippy::result_large_err)]
    /// Verify all contract state proofs against the given state root.
    pub fn verify(&self, state_root: B256) -> Result<(), ProofVerificationError> {
        self.threshold_registry.verify(state_root)?;
        self.registry_coordinator.verify(state_root)?;
        #[cfg(feature = "stale-stakes-forbidden")]
        self.service_manager.verify(state_root)?;
        self.bls_apk_registry.verify(state_root)?;
        self.stake_registry.verify(state_root)?;
        self.cert_verifier.verify(state_root)?;
        #[cfg(feature = "stale-stakes-forbidden")]
        self.delegation_manager.verify(state_root)?;

        Ok(())
    }

    ///
    /// NOTE: The data extracted is not verified. To verify the data, ensure
    /// that the [`CertStateData::verify`] is called.
    #[instrument(skip_all)]
    pub fn extract(
        &self,
        cert: &StandardCommitment,
        current_block: u32,
    ) -> Result<CertVerificationInputs, CertExtractionError> {
        let quorum_count = QuorumCountExtractor::new(cert)
            .decode_data(&self.registry_coordinator.storage_proofs)?;

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

        let next_blob_version = NextBlobVersionExtractor::new(cert)
            .decode_data(&self.threshold_registry.storage_proofs)?;

        #[cfg(feature = "stale-stakes-forbidden")]
        let staleness = {
            use eigenda_verification::verification::cert::types::Staleness;

            let stale_stakes_forbidden = StaleStakesForbiddenExtractor::new(cert)
                .decode_data(&self.service_manager.storage_proofs)?;

            let min_withdrawal_delay_blocks = MinWithdrawalDelayBlocksExtractor::new(cert)
                .decode_data(&self.delegation_manager.storage_proofs)?;

            let quorum_update_block_number = QuorumUpdateBlockNumberExtractor::new(cert)
                .decode_data(&self.registry_coordinator.storage_proofs)?;

            Staleness {
                stale_stakes_forbidden,
                min_withdrawal_delay_blocks,
                quorum_update_block_number,
            }
        };

        let storage = Storage {
            quorum_count,
            current_block,
            quorum_bitmap_history,
            operator_stake_history,
            total_stake_history,
            apk_history,
            versioned_blob_params,
            next_blob_version,
            #[cfg(feature = "stale-stakes-forbidden")]
            staleness,
        };

        let security_thresholds = SecurityThresholdsV2Extractor::new(cert)
            .decode_data(&self.cert_verifier.storage_proofs)?;

        let required_quorum_numbers = QuorumNumbersRequiredV2Extractor::new(cert)
            .decode_data(&self.cert_verifier.storage_proofs)?;

        let inputs = CertVerificationInputs {
            batch_header: cert.batch_header_v2().clone(),
            blob_inclusion_info: cert.blob_inclusion_info().clone(),
            non_signer_stakes_and_signature: cert.nonsigner_stake_and_signature().clone(),
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers: cert.signed_quorum_numbers().clone(),
            storage,
        };

        Ok(inputs)
    }
}

/// Extract certificate from the transaction. Return None if no parsable
/// certificate exists.
#[instrument(skip_all)]
pub fn extract_certificate(tx: &EthereumTxEnvelope<TxEip4844>) -> Option<StandardCommitment> {
    let raw_cert = tx.as_eip1559()?.input();
    StandardCommitment::from_rlp_bytes(raw_cert).ok()
}
