//! Error types for EigenDA certificate verification
//!
//! This module defines all possible errors that can occur during certificate
//! verification, covering cryptographic validation, stake verification, and
//! on-chain state consistency checks.

use thiserror::Error;

use crate::eigenda::{
    extraction::CertExtractionError,
    verification::cert::{
        bitmap::BitmapError,
        hash::TruncHash,
        types::{Version, history::HistoryError},
    },
};

/// Errors that can occur during certificate verification
#[derive(Error, Debug, PartialEq)]
pub enum CertVerificationError {
    /// Certificate is too old relative to the current block (recency validation failure)
    #[error("The recency window was missed, inclusion_height ({0}), recency height ({1})")]
    RecencyWindowMissed(u64, u64),

    /// Certificate's reference block is not before the current block (temporal ordering violation)
    #[error("Reference block {0} must precede current block {1}")]
    ReferenceBlockDoesNotPrecedeCurrentBlock(u32, u32),

    /// Operator public keys are not properly sorted by their hash values
    #[error("Expected pubkeys to be sorted by their hashes")]
    NotStrictlySortedByHash,

    /// Quorum state is stale and cannot be used for verification (feature-gated)
    #[cfg(feature = "stale-stakes-forbidden")]
    #[error(
        "Stale quorum, last updated at block {last_updated_at_block} should be greater than most recent stale block {most_recent_stale_block}"
    )]
    StaleQuorum {
        last_updated_at_block: u32,
        most_recent_stale_block: u32,
        window: u32,
    },

    /// BLS signature verification failed (cryptographic validation failure)
    #[error("Signature verification failed")]
    SignatureVerificationFailed,

    /// Required quorum data is missing from on-chain storage
    #[error("Missing quorum entry")]
    MissingQuorumEntry,

    /// Required signer data is missing from on-chain storage
    #[error("Missing signer entry")]
    MissingSignerEntry,

    /// Aggregate public key hash in certificate doesn't match on-chain value
    #[error(
        "Certificate apk truncated hash {cert_apk_trunc_hash} not equal to storage apk truncated hash {storage_apk_trunc_hash}"
    )]
    CertApkDoesNotEqualStorageApk {
        cert_apk_trunc_hash: TruncHash,
        storage_apk_trunc_hash: TruncHash,
    },

    /// Array or vector lengths don't match when they should be equal
    #[error("Unexpected unequal lengths")]
    UnequalLengths,

    /// Required data structure is empty when it shouldn't be
    #[error("Empty vec")]
    EmptyVec,

    /// Arithmetic overflow occurred during stake or threshold calculations
    #[error("Overflow")]
    Overflow,

    /// Arithmetic underflow occurred during stake or threshold calculations  
    #[error("Underflow")]
    Underflow,

    /// Required blob version configuration not found in threshold registry
    #[error("Missing version entry {0}")]
    MissingVersionEntry(u16),

    /// Security thresholds are incorrectly configured (confirmation must be > adversary)
    #[error("Confirmation threshold  {0} less than or equal to adversary threshold {1}")]
    ConfirmationThresholdLessThanOrEqualToAdversaryThreshold(u8, u8),

    /// Certificate fails to meet the required security assumptions for validity
    #[error("Unmet security assumptions")]
    UnmetSecurityAssumptions,

    /// Not all required quorums are present in the blob's quorum
    #[error("Required quorums not subset of blob quorums")]
    BlobQuorumsDoNotContainRequiredQuorums,

    /// Some blob quorums didn't meet confirmation thresholds
    #[error("Blob quorums not subset of confirmed quorums")]
    ConfirmedQuorumsDoNotContainBlobQuorums,

    /// Merkle inclusion proof has invalid format (must be multiple of 32 bytes)
    #[error("Merkle proof length ({0}) not multiple of 32 bytes")]
    MerkleProofLengthNotMultipleOf32Bytes(usize),

    /// Merkle proof verification failed - leaf doesn't belong to claimed tree
    #[error("Leaf node does not belong to merkle tree")]
    LeafNodeDoesNotBelongToMerkleTree,

    /// Merkle proof path is incomplete for the claimed tree depth
    #[error("Merkle proof path too short, expected {proof_depth}, found {sibling_path_len}")]
    MerkleProofPathTooShort {
        sibling_path_len: usize,
        proof_depth: usize,
    },

    /// Error occurred during on-chain data extraction (storage proofs, contract data)
    #[error(transparent)]
    WrapCertExtractionError(#[from] CertExtractionError),

    /// Error occurred during historical data processing (invalid block ranges, etc.)
    #[error(transparent)]
    WrapHistoryError(#[from] HistoryError),

    /// Error occurred during quorum bitmap operations (invalid bitmap format)
    #[error(transparent)]
    WrapBitmapError(#[from] BitmapError),

    #[error(
        "Certificate blob version ({0}) must be less than Threshold Registry's next blob version ({1})"
    )]
    InvalidBlobVersion(Version, Version),

    #[error("A blob certificate containing no quorum numbers is invalid")]
    EmptyBlobQuorums,
}
