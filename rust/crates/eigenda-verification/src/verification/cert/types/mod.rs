//! Type definitions for EigenDA certificate verification
//!
//! This module defines the core data structures and type aliases used throughout
//! the EigenDA certificate verification process, including on-chain state
//! representations and verification context.

pub mod conversions;
/// Historical data tracking for operator state changes.
///
/// This module provides utilities for tracking temporal data about operator
/// states, stakes, and quorum memberships across blockchain history.
pub mod history;

use alloy_primitives::B256;
use alloy_primitives::aliases::U96;
use ark_bn254::G1Affine;
use hashbrown::HashMap;

use crate::cert::solidity::{SecurityThresholds, VersionedBlobParams};
use crate::verification::cert::bitmap::Bitmap;
use crate::verification::cert::hash::TruncHash;
use crate::verification::cert::types::history::History;

/// Identifier for a quorum (0-255)
pub type QuorumNumber = u8;

/// Stake amount using 96-bit precision to match Ethereum's uint96
pub type Stake = U96;

/// Ethereum block number  
pub type BlockNumber = u32;

/// Key identifier for data relays
pub type RelayKey = u32;

/// Version number for blob parameters and configurations
pub type Version = u16;

/// Complete on-chain state data required for certificate verification.
///
/// This structure aggregates all the historical and current state information
/// needed to verify an EigenDA certificate, including operator stakes, quorum
/// configurations, and cryptographic commitments.
#[derive(Default, Debug, Clone)]
pub struct Storage {
    /// Total number of quorums initialized in the RegistryCoordinator
    pub quorum_count: u8,
    /// Current block number
    pub current_block: BlockNumber,
    /// Blob configuration parameters by version
    pub versioned_blob_params: HashMap<Version, VersionedBlobParams>,
    /// Next blob version
    pub next_blob_version: Version,
    /// Historical quorum membership bitmaps for each operator
    pub quorum_bitmap_history: HashMap<B256, History<Bitmap>>,
    /// Historical aggregate public key hashes for each quorum
    pub apk_history: HashMap<QuorumNumber, History<TruncHash>>,
    /// Historical total stake amounts for each quorum
    pub total_stake_history: HashMap<QuorumNumber, History<Stake>>,
    /// Historical individual operator stakes per quorum
    pub operator_stake_history: HashMap<B256, HashMap<QuorumNumber, History<Stake>>>,
    /// Security thresholds for confirmation and adversary limits
    pub security_thresholds: SecurityThresholds,
    /// Quorum numbers required to sign certificates
    pub required_quorum_numbers: alloy_primitives::Bytes,
    /// Historical on-chain storage data for verification
    /// Stale stake prevention data (feature-gated)
    #[cfg(feature = "stale-stakes-forbidden")]
    pub staleness: Staleness,
}

/// Stale stake prevention configuration and tracking data.
///
/// This structure contains information used to prevent the use of outdated
/// stake information in certificate verification, enhancing security by
/// ensuring operators can't use stale state to their advantage.
#[cfg(feature = "stale-stakes-forbidden")]
#[derive(Default, Debug, Clone)]
pub struct Staleness {
    /// Whether stale stakes are forbidden in the current configuration
    pub stale_stakes_forbidden: bool,
    /// Minimum number of blocks that must pass before stake can be withdrawn
    pub min_withdrawal_delay_blocks: BlockNumber,
    /// Block number when each quorum was last updated
    pub quorum_update_block_number: HashMap<QuorumNumber, BlockNumber>,
}

/// Quorum state during certificate verification.
///
/// Represents the computed state of a quorum at the time of certificate
/// verification, including stake calculations and aggregate public key.
#[derive(Default, Debug, Clone)]
pub(crate) struct Quorum {
    /// Quorum identifier number
    pub number: QuorumNumber,
    /// Aggregate public key for this quorum (G1 point)
    pub apk: G1Affine,
    /// Total stake registered in this quorum
    pub total_stake: Stake,
    /// Stake that participated in signing (total_stake - non_signer_stake)
    pub signed_stake: Stake,
}

/// Non-signing operator information during certificate verification.
///
/// Represents an operator that did not participate in signing the certificate,
/// along with their public key and quorum membership information.
#[derive(Default, Debug, Clone)]
pub(crate) struct NonSigner {
    /// Operator's public key (G1 point)
    pub pk: G1Affine,
    /// Hash of the operator's public key (used as operator ID)
    pub pk_hash: B256,
    /// Bitmap indicating which quorums this operator belonged to
    pub quorum_bitmap_history: Bitmap,
}
