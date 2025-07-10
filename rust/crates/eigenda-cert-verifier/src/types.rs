use alloc::vec::Vec;
use ark_bn254::G1Affine;
use primitive_types::U256;

#[derive(Default, Debug)]
pub struct NonSignerStakesAndSignature {
    /// Quorum aggregated pubkeys
    pub quorum_apks: Vec<u8>,

    /// Quorum aggregated pubkey indices
    pub quorum_apk_indices: Vec<u8>,
    pub total_stake_indices: Vec<u8>,
    pub non_signer_stake_indices: Vec<u8>,
    pub non_signer_pubkeys: Vec<G1Affine>,
    pub non_signer_quorum_bitmap_indices: Vec<u8>,
}

#[derive(Default, Debug)]
pub struct QuorumStakeTotals {
    pub total_stake_for_quorum: Vec<u128>,  // u96 in sol
    pub signed_stake_for_quorum: Vec<u128>, // u96 in sol
}

#[derive(Default, Debug)]
pub struct NonSignerInfo {
    pub _quorum_bitmaps: Vec<U256>,
    pub _non_signer_pubkey_hashes: Vec<[u8; 32]>,
}
