//! High-level contract interfaces for EigenDA data extraction
//!
//! This module provides convenient interfaces for each EigenDA smart contract,
//! aggregating the storage keys needed for certificate verification.

use alloy_primitives::StorageKey;
pub use stale_stakes_forbidden::*;

use crate::cert::StandardCommitment;
use crate::extraction::extractor::{
    ApkHistoryExtractor, NextBlobVersionExtractor, OperatorBitmapHistoryExtractor,
    OperatorStakeHistoryExtractor, QuorumCountExtractor, QuorumNumbersRequiredV2Extractor,
    QuorumUpdateBlockNumberExtractor, SecurityThresholdsV2Extractor, StorageKeyProvider,
    TotalStakeHistoryExtractor, VersionedBlobParamsExtractor,
};

/// Interface for the RegistryCoordinator contract
///
/// Manages operator registration, quorum membership, and coordination
/// between different EigenDA registry components.
pub struct RegistryCoordinator;

impl RegistryCoordinator {
    /// Get all storage keys needed for subsequent data extraction
    ///
    /// # Arguments
    /// * `certificate` - The certificate being verified
    ///
    /// # Returns
    /// Vector of storage keys for:
    /// - Quorum count
    /// - Operator bitmap histories  
    /// - Quorum update block numbers
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

/// Interface for the StakeRegistry contract
///
/// Tracks operator stakes across different quorums maintaining
/// historical stake information
pub struct StakeRegistry;

impl StakeRegistry {
    /// Get all storage keys needed for subsequent data extraction
    ///
    /// # Arguments
    /// * `certificate` - The certificate being verified
    ///
    /// # Returns  
    /// Vector of storage keys for:
    /// - Individual operator stake histories
    /// - Total stake histories per quorum
    pub fn storage_keys(certificate: &StandardCommitment) -> Vec<StorageKey> {
        let operator_stake_history = OperatorStakeHistoryExtractor::new(certificate).storage_keys();
        let total_stake_history = TotalStakeHistoryExtractor::new(certificate).storage_keys();

        [operator_stake_history, total_stake_history]
            .into_iter()
            .flatten()
            .collect()
    }
}

/// Interface for the BlsApkRegistry contract
///
/// Manages BLS aggregate public keys (APKs) for each quorum,
/// enabling signature verification.
pub struct BlsApkRegistry;

impl BlsApkRegistry {
    /// Get all storage keys needed for subsequent data extraction
    ///
    /// # Arguments
    /// * `certificate` - The certificate being verified
    ///
    /// # Returns
    /// Vector of storage keys for APK histories
    pub fn storage_keys(certificate: &StandardCommitment) -> Vec<StorageKey> {
        ApkHistoryExtractor::new(certificate).storage_keys()
    }
}

/// Interface for the EigenDaThresholdRegistry contract
///
/// Manages blob versioning parameters and thresholds for
/// data availability requirements.
pub struct EigenDaThresholdRegistry;

impl EigenDaThresholdRegistry {
    /// Get all storage keys needed for subsequent data extraction
    ///
    /// # Arguments
    /// * `certificate` - The certificate being verified
    ///
    /// # Returns
    /// Vector of storage keys for:
    /// - Versioned blob parameters
    /// - Next blob version
    pub fn storage_keys(certificate: &StandardCommitment) -> Vec<StorageKey> {
        let versioned_blob_params = VersionedBlobParamsExtractor::new(certificate).storage_keys();
        let next_blob_version = NextBlobVersionExtractor::new(certificate).storage_keys();

        [versioned_blob_params, next_blob_version]
            .into_iter()
            .flatten()
            .collect()
    }
}

/// Interface for the EigenDaCertVerifier contract
///
/// Contains security parameters and requirements for certificate
/// verification, including thresholds and required quorum numbers.
pub struct EigenDaCertVerifier;

impl EigenDaCertVerifier {
    /// Get all storage keys needed for subsequent data extraction
    ///
    /// # Arguments
    /// * `certificate` - The certificate being verified
    ///
    /// # Returns
    /// Vector of storage keys for:
    /// - Security thresholds (confirmation and adversary thresholds)
    /// - Required quorum numbers
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

mod stale_stakes_forbidden {
    //! Additional contract interfaces for guarding against stale stakes
    //!
    //! These interfaces expose EigenDA contract storage required for stale stake prevention.

    use alloy_primitives::StorageKey;

    use crate::cert::StandardCommitment;
    use crate::extraction::extractor::{
        MinWithdrawalDelayBlocksExtractor, StaleStakesForbiddenExtractor, StorageKeyProvider,
    };

    /// Interface for the EigenDA Service Manager contract (stale stakes functionality)
    ///
    /// Provides access to the `staleStakesForbidden` flag which controls whether
    /// the system accepts operator stakes that may be considered "stale" during
    /// certificate verification.
    pub struct ServiceManager;

    impl ServiceManager {
        /// Get all storage keys needed for subsequent data extraction
        ///
        /// # Arguments
        /// * `certificate` - The certificate being verified
        ///
        /// # Returns
        /// Vector containing the storage key for the `staleStakesForbidden` boolean flag.
        /// When this flag is `true`, additional staleness checks are performed during
        /// verification to ensure operator stakes were updated recently enough relative
        /// to the reference block number.
        ///
        /// # Security Context
        /// When `staleStakesForbidden` is enabled, the system prevents potential attacks
        /// where an adversary could exploit the delay between stake updates and verification
        /// by using operator stake information that is too old to be trustworthy.
        pub fn storage_keys(certificate: &StandardCommitment) -> Vec<StorageKey> {
            StaleStakesForbiddenExtractor::new(certificate).storage_keys()
        }
    }

    /// Interface for the EigenLayer Delegation Manager contract (stale stakes functionality)
    ///
    /// Provides access to withdrawal delay parameters that define the time window
    /// for determining stake staleness.
    pub struct DelegationManager;

    impl DelegationManager {
        /// Get all storage keys needed for subsequent data extraction
        ///
        /// # Arguments
        /// * `certificate` - The certificate being verified
        ///
        /// # Returns
        /// Vector containing the storage key for `minWithdrawalDelayBlocks`.
        /// This value defines the minimum number of blocks that must pass before
        /// a withdrawal can be completed, and is used as the threshold for determining
        /// whether operator stakes are "stale" when `staleStakesForbidden` is enabled.
        ///
        /// # Staleness Logic
        /// Stakes are considered stale if the last quorum update occurred more than
        /// `minWithdrawalDelayBlocks` blocks before the `referenceBlockNumber`.
        /// This ensures that operator stakes reflect a recent enough view of the
        /// network state to be trusted for verification.
        pub fn storage_keys(certificate: &StandardCommitment) -> Vec<StorageKey> {
            MinWithdrawalDelayBlocksExtractor::new(certificate).storage_keys()
        }
    }
}

#[cfg(test)]
mod tests {
    use std::collections::HashSet;

    use crate::cert::StandardCommitment;
    use crate::extraction::contract::{
        BlsApkRegistry, DelegationManager, EigenDaCertVerifier, EigenDaThresholdRegistry,
        RegistryCoordinator, ServiceManager, StakeRegistry,
    };
    use crate::extraction::extractor::{
        ApkHistoryExtractor, MinWithdrawalDelayBlocksExtractor, NextBlobVersionExtractor,
        OperatorBitmapHistoryExtractor, OperatorStakeHistoryExtractor, QuorumCountExtractor,
        QuorumNumbersRequiredV2Extractor, QuorumUpdateBlockNumberExtractor,
        SecurityThresholdsV2Extractor, StaleStakesForbiddenExtractor, StorageKeyProvider,
        TotalStakeHistoryExtractor, VersionedBlobParamsExtractor,
    };

    fn create_test_commitment() -> StandardCommitment {
        let commitment_hex = "02f90389e5a0c769488dd5264b3ef21dce7ee2d42fba43e1f83ff228f501223e38818cb14492833f44fcf901eff901caf9018180820001f90159f842a0012e810ffc0a83074b3d14db9e78bbae623f7770cac248df9e73fac6b9d59d17a02a916ffbbf9dde4b7ebe94191a29ff686422d7dcb3b47ecb03c6ada75a9c15c8f888f842a01811c8b4152fce9b8c4bae61a3d097e61dfc43dc7d45363d19e7c7f1374034ffa001edc62174217cdce60a4b52fa234ac0d96db4307dac9150e152ba82cbb4d2f1f842a00f423b0dbc1fe95d2e3f7dbac6c099e51dbf73400a4b3f26b9a29665b4ac58a8a01855a2bd56c0e8f4cc85ac149cf9a531673d0e89e22f0d6c4ae419ed7c5d2940f888f842a02667cbb99d60fa0d7f3544141d3d531dceeeb50b06e5a0cdc42338a359138ae4a00dff4c929d8f8a307c19bba6e8006fe6700f6554cef9eb3797944f89472ffb30f842a004c17a6225acd5b4e7d672a1eb298c5358f4f6f17d04fd1ee295d0c0d372fa84a024bc3ad4d5e54f54f71db382ce276f37ac3c260cc74306b832e8a3c93c7951d302a0e43e11e2405c2fd1d880af8612d969b654827e0ba23d9feb3722ccce6226fce7b8411ddf4553c79c0515516fd3c8b3ae6a756b05723f4d0ebe98a450c8bcc96cbb355ef07a44eeb56f831be73647e4da20e22fa859f984ee41d6efcd3692063b0b0601c2800101a0a69e552a6fc2ff75d32edaf5313642ddeebe60d2069435d12e266ce800e9e96bf9016bc0c0f888f842a00d45727a99053af8d38d4716ab83ace676096e7506b6b7aa6953e87bc04a023ca016c030c31dd1c94062948ecdce2e67c4e6626c16af0033dcdb7a96362c937d48f842a00a95fac74aba7e3fbd24bc62457ce6981803d8f5fef28871d3d5e2af05d50cd4a0117400693917cd50d9bc28d4ab4fadf93a23e771f303637f8d1f83cd0632c3fcf888f842a0301bfced3253e99e8d50f2fed62313a16d714013d022a4dc4294656276f10d1ba0152e047a83c326a9d81dac502ec429b662b58ee119ca4c8748a355b539c24131f842a01944b5b4a3e93d46b0fe4370128c6cdcd066ae6b036b019a20f8d22fe9a10d67a00ddf3421722967c0bd965b9fc9e004bf01183b6206fec8de65e40331d185372ef842a02db8fb278708abf8878ebf578872ab35ee914ad8196b78de16b34498222ac1c2a02ff9d9a5184684f4e14530bde3a61a2f9adaa74734dff104b61ba3d963a644dac68207388208b7c68209998209c5c2c0c0820001";
        let raw_commitment = hex::decode(commitment_hex).unwrap();
        StandardCommitment::from_rlp_bytes(raw_commitment.as_slice()).unwrap()
    }

    #[test]
    fn registry_coordinator_storage_keys() {
        let certificate = create_test_commitment();
        let keys = RegistryCoordinator::storage_keys(&certificate);

        assert!(!keys.is_empty(), "Should generate required data");

        let keys_set: HashSet<_> = keys.iter().collect();
        assert_eq!(
            keys_set.len(),
            keys.len(),
            "All generated items should be unique"
        );

        // Verify expected item count based on feature flags
        let quorum_count_keys = QuorumCountExtractor::new(&certificate).storage_keys();
        let quorum_bitmap_keys = OperatorBitmapHistoryExtractor::new(&certificate).storage_keys();
        let quorum_update_keys = QuorumUpdateBlockNumberExtractor::new(&certificate).storage_keys();
        let expected_total =
            quorum_count_keys.len() + quorum_bitmap_keys.len() + quorum_update_keys.len();
        assert_eq!(
            keys.len(),
            expected_total,
            "Should include all required data"
        );
    }

    #[test]
    fn stake_registry_storage_keys() {
        let certificate = create_test_commitment();
        let keys = StakeRegistry::storage_keys(&certificate);

        assert!(!keys.is_empty(), "Should generate required data");

        let operator_stake_keys = OperatorStakeHistoryExtractor::new(&certificate).storage_keys();
        let total_stake_keys = TotalStakeHistoryExtractor::new(&certificate).storage_keys();
        let expected_total = operator_stake_keys.len() + total_stake_keys.len();

        assert_eq!(
            keys.len(),
            expected_total,
            "Should include all expected data"
        );

        let keys_set: HashSet<_> = keys.iter().collect();
        assert_eq!(
            keys_set.len(),
            keys.len(),
            "All generated items should be unique"
        );
    }

    #[test]
    fn bls_apk_registry_storage_keys() {
        let certificate = create_test_commitment();
        let keys = BlsApkRegistry::storage_keys(&certificate);

        assert!(!keys.is_empty(), "Should generate required data");

        let apk_history_keys = ApkHistoryExtractor::new(&certificate).storage_keys();
        assert_eq!(
            keys.len(),
            apk_history_keys.len(),
            "Should match expected data size"
        );
        assert_eq!(
            keys, apk_history_keys,
            "Should return exactly the required data"
        );
    }

    #[test]
    fn eigen_da_threshold_registry_storage_keys() {
        let certificate = create_test_commitment();
        let keys = EigenDaThresholdRegistry::storage_keys(&certificate);

        assert!(!keys.is_empty(), "Should generate required data");

        let versioned_blob_keys = VersionedBlobParamsExtractor::new(&certificate).storage_keys();
        let next_blob_keys = NextBlobVersionExtractor::new(&certificate).storage_keys();
        let expected_total = versioned_blob_keys.len() + next_blob_keys.len();

        assert_eq!(
            keys.len(),
            expected_total,
            "Should include all expected data"
        );

        let keys_set: HashSet<_> = keys.iter().collect();
        assert_eq!(
            keys_set.len(),
            keys.len(),
            "All generated items should be unique"
        );
    }

    #[test]
    fn eigen_da_cert_verifier_storage_keys() {
        let certificate = create_test_commitment();
        let keys = EigenDaCertVerifier::storage_keys(&certificate);

        assert!(!keys.is_empty(), "Should generate required data");

        let security_threshold_keys =
            SecurityThresholdsV2Extractor::new(&certificate).storage_keys();
        let quorum_numbers_keys =
            QuorumNumbersRequiredV2Extractor::new(&certificate).storage_keys();
        let expected_total = security_threshold_keys.len() + quorum_numbers_keys.len();

        assert_eq!(
            keys.len(),
            expected_total,
            "Should include all expected data"
        );

        let keys_set: HashSet<_> = keys.iter().collect();
        assert_eq!(
            keys_set.len(),
            keys.len(),
            "All generated items should be unique"
        );
    }

    #[test]
    fn service_manager_storage_keys() {
        let certificate = create_test_commitment();
        let keys = ServiceManager::storage_keys(&certificate);

        assert!(!keys.is_empty(), "Should generate required data");

        let stale_stakes_keys = StaleStakesForbiddenExtractor::new(&certificate).storage_keys();
        assert_eq!(
            keys.len(),
            stale_stakes_keys.len(),
            "Should match expected data size"
        );
        assert_eq!(
            keys, stale_stakes_keys,
            "Should return exactly the required data"
        );
    }

    #[test]
    fn delegation_manager_storage_keys() {
        let certificate = create_test_commitment();
        let keys = DelegationManager::storage_keys(&certificate);

        assert!(!keys.is_empty(), "Should generate required data");

        let min_withdrawal_keys =
            MinWithdrawalDelayBlocksExtractor::new(&certificate).storage_keys();
        assert_eq!(
            keys.len(),
            min_withdrawal_keys.len(),
            "Should match expected data size"
        );
        assert_eq!(
            keys, min_withdrawal_keys,
            "Should return exactly the required data"
        );
    }
}
