use alloy_primitives::StorageKey;

#[cfg(feature = "stale-stakes-forbidden")]
pub use stale_stakes_forbidden::*;

#[cfg(feature = "stale-stakes-forbidden")]
use crate::eigenda::extraction::{
    MinWithdrawalDelayBlocksExtractor, QuorumUpdateBlockNumberExtractor,
    StaleStakesForbiddenExtractor,
};

use crate::eigenda::{
    cert::StandardCommitment,
    extraction::{
        ApkHistoryExtractor, NextBlobVersionExtractor, OperatorBitmapHistoryExtractor,
        OperatorStakeHistoryExtractor, QuorumCountExtractor, QuorumNumbersRequiredV2Extractor,
        SecurityThresholdsV2Extractor, StorageKeyProvider, TotalStakeHistoryExtractor,
        VersionedBlobParamsExtractor,
    },
};

pub struct RegistryCoordinator;

impl RegistryCoordinator {
    pub fn storage_keys(certificate: &StandardCommitment) -> Vec<StorageKey> {
        let quorum_count = QuorumCountExtractor::new(certificate).storage_keys();
        let quorum_bitmap_history = OperatorBitmapHistoryExtractor::new(certificate).storage_keys();
        #[cfg(feature = "stale-stakes-forbidden")]
        let quorum_update_block_number =
            QuorumUpdateBlockNumberExtractor::new(certificate).storage_keys();

        [
            quorum_count,
            quorum_bitmap_history,
            #[cfg(feature = "stale-stakes-forbidden")]
            quorum_update_block_number,
        ]
        .into_iter()
        .flatten()
        .collect()
    }
}

pub struct StakeRegistry;

impl StakeRegistry {
    pub fn storage_keys(certificate: &StandardCommitment) -> Vec<StorageKey> {
        let operator_stake_history = OperatorStakeHistoryExtractor::new(certificate).storage_keys();
        let total_stake_history = TotalStakeHistoryExtractor::new(certificate).storage_keys();

        [operator_stake_history, total_stake_history]
            .into_iter()
            .flatten()
            .collect()
    }
}

pub struct BlsApkRegistry;

impl BlsApkRegistry {
    pub fn storage_keys(certificate: &StandardCommitment) -> Vec<StorageKey> {
        ApkHistoryExtractor::new(certificate).storage_keys()
    }
}

pub struct EigenDaThresholdRegistry;

impl EigenDaThresholdRegistry {
    pub fn storage_keys(certificate: &StandardCommitment) -> Vec<StorageKey> {
        let versioned_blob_params = VersionedBlobParamsExtractor::new(certificate).storage_keys();
        let next_blob_version = NextBlobVersionExtractor::new(certificate).storage_keys();

        [versioned_blob_params, next_blob_version]
            .into_iter()
            .flatten()
            .collect()
    }
}

pub struct EigenDaCertVerifier;

impl EigenDaCertVerifier {
    pub fn storage_keys(certificate: &StandardCommitment) -> Vec<StorageKey> {
        let security_thresholds = SecurityThresholdsV2Extractor::new(certificate).storage_keys();
        let required_quorum_numbers =
            QuorumNumbersRequiredV2Extractor::new(certificate).storage_keys();

        [security_thresholds, required_quorum_numbers]
            .into_iter()
            .flatten()
            .collect()
    }
}

#[cfg(feature = "stale-stakes-forbidden")]
mod stale_stakes_forbidden {
    use super::*;

    pub struct ServiceManager;

    impl ServiceManager {
        pub fn storage_keys(certificate: &StandardCommitment) -> Vec<StorageKey> {
            StaleStakesForbiddenExtractor::new(certificate).storage_keys()
        }
    }

    pub struct DelegationManager;

    impl DelegationManager {
        pub fn storage_keys(certificate: &StandardCommitment) -> Vec<StorageKey> {
            MinWithdrawalDelayBlocksExtractor::new(certificate).storage_keys()
        }
    }
}
