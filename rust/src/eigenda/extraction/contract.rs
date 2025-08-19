use alloy_primitives::StorageKey;

use crate::eigenda::{
    extraction::{
        ApkHistoryExtractor, MinWithdrawalDelayBlocksExtractor, OperatorBitmapHistoryExtractor,
        OperatorStakeHistoryExtractor, QuorumCountExtractor, QuorumNumbersRequiredV2Extractor,
        QuorumUpdateBlockNumberExtractor, RelayKeyToRelayInfoExtractor,
        SecurityThresholdsV2Extractor, StaleStakesForbiddenExtractor, StorageKeyProvider,
        TotalStakeHistoryExtractor, VersionedBlobParamsExtractor,
    },
    types::StandardCommitment,
};

pub struct EigenDaRelayRegistry;

impl EigenDaRelayRegistry {
    pub fn storage_keys(certificate: &StandardCommitment) -> Vec<StorageKey> {
        RelayKeyToRelayInfoExtractor::new(certificate).storage_keys()
    }
}

pub struct BlsSignatureChecker;

impl BlsSignatureChecker {
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

pub struct RegistryCoordinator;

impl RegistryCoordinator {
    pub fn storage_keys(certificate: &StandardCommitment) -> Vec<StorageKey> {
        let quorum_count = QuorumCountExtractor::new(certificate).storage_keys();
        let quorum_bitmap_history = OperatorBitmapHistoryExtractor::new(certificate).storage_keys();
        let quorum_update_block_number =
            QuorumUpdateBlockNumberExtractor::new(certificate).storage_keys();

        [
            quorum_count,
            quorum_bitmap_history,
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
        VersionedBlobParamsExtractor::new(certificate).storage_keys()
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
