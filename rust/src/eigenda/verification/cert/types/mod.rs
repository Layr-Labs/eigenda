pub mod conversions;
pub mod history;
pub mod solidity;

use alloy_primitives::{B256, aliases::U96};
use ark_bn254::G1Affine;
use hashbrown::HashMap;

use crate::eigenda::verification::cert::{
    bitmap::Bitmap,
    hash::TruncHash,
    types::{history::History, solidity::VersionedBlobParams},
};

pub type QuorumNumber = u8;
pub type Stake = U96;
pub type BlockNumber = u32;
pub type RelayKey = u32;
pub type Version = u16;

#[derive(Default, Debug, Clone)]
pub struct Storage {
    // number of quorums initialized in the RegistryCoordinator
    pub quorum_count: u8,
    pub current_block: BlockNumber,
    pub versioned_blob_params: HashMap<Version, VersionedBlobParams>,
    pub next_blob_version: Version,
    pub quorum_bitmap_history: HashMap<B256, History<Bitmap>>,
    pub apk_history: HashMap<QuorumNumber, History<TruncHash>>,
    pub total_stake_history: HashMap<QuorumNumber, History<Stake>>,
    pub operator_stake_history: HashMap<B256, HashMap<QuorumNumber, History<Stake>>>,
    #[cfg(feature = "stale-stakes-forbidden")]
    pub staleness: Staleness,
}

#[cfg(feature = "stale-stakes-forbidden")]
#[derive(Default, Debug, Clone)]
pub struct Staleness {
    pub stale_stakes_forbidden: bool,
    pub min_withdrawal_delay_blocks: BlockNumber,
    pub quorum_update_block_number: HashMap<QuorumNumber, BlockNumber>,
}

#[derive(Default, Debug, Clone)]
pub(crate) struct Quorum {
    pub number: QuorumNumber,
    pub apk: G1Affine,
    pub total_stake: Stake,
    pub signed_stake: Stake,
}

#[derive(Default, Debug, Clone)]
pub(crate) struct NonSigner {
    pub pk: G1Affine,
    pub pk_hash: B256,
    pub quorum_bitmap_history: Bitmap,
}
