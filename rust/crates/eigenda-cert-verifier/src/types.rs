// todo: consider using the existing crate `rust-eigenda-v2-common`
use alloc::vec::Vec;

use ark_bn254::{G1Affine, G2Affine};
use hashbrown::HashMap;

use crate::{
    bitmap::Bitmap,
    hash::{BeHash, TruncatedBeHash},
    types::history::History,
};

pub mod history;

pub type QuorumNumber = u8;
pub type Stake = u128; // u96 in sol
pub type BlockNumber = u32;
pub type SignerId = BeHash;

#[derive(Default, Debug, Clone)]
pub struct Chain {
    // number of quorums initialized in the RegistryCoordinator
    pub initialized_quorums_count: u8,
    pub current_block: BlockNumber,
    pub reject_staleness: bool,
    pub min_withdrawal_delay_blocks: BlockNumber,

    pub quorum_membership_history_by_signer: HashMap<SignerId, History<Bitmap>>,
    pub stake_history_by_signer_and_quorum:
        HashMap<SignerId, HashMap<QuorumNumber, History<Stake>>>,
    pub total_stake_history_by_quorum: HashMap<QuorumNumber, History<Stake>>,
    pub apk_trunc_hash_history_by_quorum: HashMap<QuorumNumber, History<TruncatedBeHash>>,
    pub last_updated_at_block_by_quorum: HashMap<QuorumNumber, BlockNumber>,
}

#[derive(Default, Debug, Clone)]
pub struct Cert {
    pub msg_hash: BeHash,
    pub reference_block: BlockNumber,
    pub signed_quorums: Vec<QuorumNumber>,
    pub params: NonSignerStakesAndSignature,
}

#[derive(Default, Debug, Clone)]
pub struct NonSignerStakesAndSignature {
    pub apk_for_each_quorum: Vec<G1Affine>,
    pub apk_index_for_each_quorum: Vec<u32>,
    pub total_stake_index_for_each_quorum: Vec<u32>,
    pub stake_index_for_each_quorum_and_required_non_signer: Vec<Vec<u32>>,
    pub pk_for_each_non_signer: Vec<G1Affine>,
    pub quorum_membership_index_for_each_non_signer: Vec<u32>,
    pub apk_g2: G2Affine,
    pub sigma: G1Affine,
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
    pub pk_hash: BeHash,
    pub quorum_membership: Bitmap,
}
