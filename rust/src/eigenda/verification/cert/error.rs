use thiserror::Error;

use crate::eigenda::{
    extraction::CertExtractionError,
    verification::cert::{
        bitmap::BitmapError,
        hash::TruncHash,
        types::{RelayKey, history::HistoryError},
    },
};

#[derive(Error, Debug, PartialEq)]
pub enum CertVerificationError {
    #[error("The recency window was missed, inclusion_height ({0}), recency height ({1})")]
    RecencyWindowMissed(u64, u64),

    #[error("Ancestor data missing")]
    AncestorDataMissing,

    #[error("Reference block {0} must precede current block {1}")]
    ReferenceBlockDoesNotPrecedeCurrentBlock(u32, u32),

    #[error("Expected pubkeys to be sorted by their hashes")]
    NotStrictlySortedByHash,

    #[cfg(feature = "stale-stakes-forbidden")]
    #[error(
        "Stale quorum, last updated at block {last_updated_at_block} should be greater than most recent stale block {most_recent_stale_block}"
    )]
    StaleQuorum {
        last_updated_at_block: u32,
        most_recent_stale_block: u32,
        window: u32,
    },

    #[error("Signature verification failed")]
    SignatureVerificationFailed,

    #[error("Missing quorum entry")]
    MissingQuorumEntry,

    #[error("Missing signer entry")]
    MissingSignerEntry,

    #[error(
        "Certificate apk truncated hash {cert_apk_trunc_hash} not equal to storage apk truncated hash {storage_apk_trunc_hash}"
    )]
    CertApkDoesNotEqualStorageApk {
        cert_apk_trunc_hash: TruncHash,
        storage_apk_trunc_hash: TruncHash,
    },

    #[error("Unexpected unequal lengths")]
    UnequalLengths,

    #[error("Empty vec")]
    EmptyVec,

    #[error("Overflow")]
    Overflow,

    #[error("Underflow")]
    Underflow,

    #[error("Missing relay key entry {0}")]
    MissingRelayKeyEntry(RelayKey),

    #[error("Relay key not set")]
    RelayKeyNotSet,

    #[error("Missing version entry {0}")]
    MissingVersionEntry(u16),

    #[error("Confirmation threshold  {0} less than or equal to adversary threshold {1}")]
    ConfirmationThresholdLessThanOrEqualToAdversaryThreshold(u8, u8),

    #[error("Unmet security assumptions")]
    UnmetSecurityAssumptions,

    #[error("Required quorums not subset of blob quorums")]
    BlobQuorumsDoNotContainRequiredQuorums,

    #[error("Blob quorums not subset of confirmed quorums")]
    ConfirmedQuorumsDoNotContainBlobQuorums,

    #[error("Merkle proof length ({0}) not multiple of 32 bytes")]
    MerkleProofLengthNotMultipleOf32Bytes(usize),

    #[error("Leaf node does not belong to merkle tree")]
    LeafNodeDoesNotBelongToMerkleTree,

    #[error("Merkle proof path too short, expected {proof_depth}, found {sibling_path_len}")]
    MerkleProofPathTooShort {
        sibling_path_len: usize,
        proof_depth: usize,
    },

    #[error(transparent)]
    WrapCertExtractionError(#[from] CertExtractionError),

    #[error(transparent)]
    WrapHistoryError(#[from] HistoryError),

    #[error(transparent)]
    WrapBitmapError(#[from] BitmapError),
}
