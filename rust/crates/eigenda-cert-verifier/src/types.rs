use alloy_primitives::B256;
use ark_bn254::G1Affine;
use hashbrown::HashMap;

use crate::{
    bitmap::Bitmap,
    hash::TruncatedB256,
    types::{
        history::History,
        solidity::{RelayInfo, VersionedBlobParams},
    },
};

pub mod conversions;
pub mod history;
pub mod solidity;

pub type QuorumNumber = u8;
pub type Stake = u128; // u96 in sol
pub type BlockNumber = u32;
pub type RelayKey = u32;
pub type Version = u16;

#[derive(Default, Debug, Clone)]
pub struct Storage {
    // number of quorums initialized in the RegistryCoordinator
    pub initialized_quorums_count: u8,
    pub current_block: BlockNumber,
    pub reject_staleness: bool,
    pub min_withdrawal_delay_blocks: BlockNumber,

    pub quorum_membership_history_by_signer: HashMap<B256, History<Bitmap>>,
    pub stake_history_by_signer_and_quorum: HashMap<B256, HashMap<QuorumNumber, History<Stake>>>,
    pub total_stake_history_by_quorum: HashMap<QuorumNumber, History<Stake>>,
    pub apk_trunc_hash_history_by_quorum: HashMap<QuorumNumber, History<TruncatedB256>>,

    pub last_updated_at_block_by_quorum: HashMap<QuorumNumber, BlockNumber>,
    pub relay_key_to_relay_info: HashMap<RelayKey, RelayInfo>,
    pub version_to_versioned_blob_params: HashMap<Version, VersionedBlobParams>,
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
    pub quorum_membership: Bitmap,
}
