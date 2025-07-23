// TODO: consider using the existing crate `rust-eigenda-v2-common`
use alloc::{string::String, vec::Vec};

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
pub type RelayKey = u32;
pub type Version = u16;

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
    pub relay_key_to_relay_info: HashMap<RelayKey, RelayInfo>,
    pub version_to_versioned_blob_params: HashMap<Version, VersionedBlobParams>,
}

// todo: match the data schema of IEigenDAStructs.sol
#[derive(Default, Debug, Clone)]
pub struct Cert {
    pub msg_hash: BeHash,
    pub reference_block: BlockNumber,
    pub signed_quorums: Vec<QuorumNumber>,
    pub params: NonSignerStakesAndSignature,
    pub blob_inclusion_info: BlobInclusionInfo,
    pub security_thresholds: SecurityThresholds,
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

pub type Address = [u8; 20];

#[derive(Default, Debug, Clone)]
pub struct RelayInfo {
    pub address: Address,
    pub url: String,
}

#[derive(Default, Debug, Clone)]
pub struct BlobInclusionInfo {
    pub blob_certificate: BlobCertificate,
    pub blob_index: Vec<u32>,
    pub inclusion_proof: Vec<u8>,
}

#[derive(Default, Debug, Clone)]
pub struct BlobCertificate {
    pub blob_header: BlobHeaderV2,
    pub signature: Vec<u8>,
    pub relay_keys: Vec<RelayKey>,
}

#[derive(Default, Debug, Clone)]
pub struct BlobHeaderV2 {
    pub version: u16,
    pub quorum_numbers: Vec<u8>,
    pub commitment: BlobCommitment,
    pub payment_header_hash: [u8; 32],
}

#[derive(Default, Debug, Clone)]
pub struct BlobCommitment {
    pub commitment: G1Affine,
    pub length_commitment: G2Affine,
    pub length_proof: G2Affine,
    pub length: u32,
}

#[derive(Default, Debug, Clone)]
pub struct SecurityThresholds {
    pub confirmation_threshold: u8,
    pub adversary_threshold: u8,
}

#[derive(Default, Debug, Clone)]
pub struct VersionedBlobParams {
    pub max_num_operators: u32,
    pub num_chunks: u32,
    pub coding_rate: u8,
}
