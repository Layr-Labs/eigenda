use thiserror::Error;

#[derive(Error, Debug, PartialEq)]
pub enum CertVerificationError {
    #[error("Reference block must precede current block")]
    ReferenceBlockDoesNotPrecedeCurrentBlock,

    #[error("Bit indices length exceeds max byte slice length")]
    BitIndicesGreaterThanMaxLength,

    #[error("Bit indices not unique")]
    BitIndicesNotUnique,

    #[error("Bit indices not ordered")]
    BitIndicesNotSorted,

    #[error("One or more bit indices are greater than or equal to the provided upper bound")]
    BitIndexNotLessThanUpperBound,

    #[error("Unexpected identity point in operation")]
    PointAtInfinity,

    #[error("Expected pubkeys to be sorted by their hashes")]
    NotStrictlySortedByHash,

    #[error("Stale quorum")]
    StaleQuorum,

    #[error("Signature verification failed")]
    SignatureVerificationFailed,

    #[error("Element not in interval")]
    ElementNotInInterval,

    #[error("Degenerate interval")]
    DegenerateInterval,

    #[error("Missing quorum entry")]
    MissingQuorumEntry,

    #[error("Missing signer entry")]
    MissingSignerEntry,

    #[error("Missing history entry")]
    MissingHistoryEntry,

    #[error("Certificate quorum apk not equal to chain quorum apk")]
    CertApkDoesNotEqualChainApk,

    #[error("Unexpected unequal lengths")]
    UnequalLengths,

    #[error("Empty vec")]
    EmptyVec,

    #[error("Underflow")]
    Underflow,
}
