//! High-level contract interfaces for EigenDA data extraction
//!
//! This module provides convenient interfaces for each EigenDA smart contract,
//! aggregating the storage keys needed for certificate verification.

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
    /// - Quorum update block numbers (if stale-stakes-forbidden feature is enabled)
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
/// enabling efficient signature verification.
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

#[cfg(feature = "stale-stakes-forbidden")]
mod stale_stakes_forbidden {
    //! Additional contract interfaces for guarding against stale stakes
    //!
    //! These interfaces are only available when the `stale-stakes-forbidden` feature is enabled.
    //! They provide access to parameters that control whether "stale" operator stakes can be
    //! used for verification, adding extra security against attacks using outdated stake information.
    //!
    //! > **Note:** This functionality has been gated behind this flag because relevant
    //! > EigenDA contracts are currently deployed having `staleStakesForbidden` set
    //! > to `false`. Should that change in the future, the functionality can be
    //! > activated by simply enabling this feature.

    use super::*;

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
