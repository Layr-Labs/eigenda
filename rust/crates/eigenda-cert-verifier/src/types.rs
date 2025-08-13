use alloy_primitives::{Address, B256, aliases::U96};
use ark_bn254::G1Affine;
use hashbrown::HashMap;

use crate::{
    bitmap::Bitmap,
    hash::TruncatedB256,
    types::{history::History, solidity::VersionedBlobParams},
};

pub mod conversions;
pub mod history;
pub mod solidity;

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
    pub stale_stakes_forbidden: bool,
    pub min_withdrawal_delay_blocks: BlockNumber,

    pub quorum_update_block_number: HashMap<QuorumNumber, BlockNumber>,
    pub relay_key_to_relay_address: HashMap<RelayKey, Address>,
    pub versioned_blob_params: HashMap<Version, VersionedBlobParams>,

    pub quorum_bitmap_history: HashMap<B256, History<Bitmap>>,
    pub apk_history: HashMap<QuorumNumber, History<TruncatedB256>>,
    pub total_stake_history: HashMap<QuorumNumber, History<Stake>>,
    pub operator_stake_history: HashMap<B256, HashMap<QuorumNumber, History<Stake>>>,
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
