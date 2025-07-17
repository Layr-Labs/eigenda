// todo: consider using the existing crate `rust-eigenda-v2-common`
use alloc::vec::Vec;

use ark_bn254::{G1Affine, G2Affine};
use hashbrown::HashMap;

use crate::{bitmap_utils::Bitmap, hash::BeHash};

#[derive(Default, Debug)]
pub struct NonSignerStakesAndSignature {
    /// Quorum aggregated pubkeys
    pub quorum_apks: Vec<G1Affine>,

    /// Quorum aggregated pubkey indices
    pub quorum_apk_indices: Vec<u8>,
    pub total_stake_indices: Vec<u8>,
    pub non_signer_stake_indices: Vec<u8>,
    pub non_signer_pubkeys: Vec<G1Affine>,
    pub non_signer_quorum_bitmap_indices: Vec<u8>,
    pub apk_g2: G2Affine,
    pub sigma: G1Affine,
}

#[derive(Default, Debug)]
pub struct QuorumStakeTotals {
    pub total_stake_for_quorum: Vec<u128>,  // u96 in sol
    pub signed_stake_for_quorum: Vec<u128>, // u96 in sol
}

#[derive(Default, Debug)]
pub struct NonSignerInfo {
    pub _quorum_bitmaps: Vec<Bitmap>,
    pub _non_signer_pubkey_hashes: Vec<BeHash>,
}

#[derive(Default, Debug)]
pub struct ReferenceBlock {
    pub number: u32,
    pub hash_to_bitmap: HashMap<BeHash, Bitmap>,
}
