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

    #[error("Certificate quorum apk not equal to storage quorum apk")]
    CertApkDoesNotEqualStorageApk,

    #[error("Unexpected unequal lengths")]
    UnequalLengths,

    #[error("Empty vec")]
    EmptyVec,

    #[error("Overflow")]
    Overflow,

    #[error("Underflow")]
    Underflow,

    #[error("Missing relay key entry")]
    MissingRelayKeyEntry,

    #[error("Relay key not set")]
    RelayKeyNotSet,

    #[error("Missing version entry")]
    MissingVersionEntry,

    #[error("Confirmation threshold not greater than adversary threshold")]
    ConfirmationThresholdNotGreaterThanAdversaryThreshold,

    #[error("Unmet security assumptions")]
    UnmetSecurityAssumptions,

    #[error("Required quorums not subset of blob quorums")]
    BlobQuorumsDoNotContainRequiredQuorums,

    #[error("Blob quorums not subset of confirmed quorums")]
    ConfirmedQuorumsDoNotContainBlobQuorums,

    #[error("Merkle proof length not multiple of 32 bytes")]
    MerkleProofLengthNotMultipleOf32Bytes,

    #[error("Leaf node does not belong to merkle tree")]
    LeafNodeDoesNotBelongToMerkleTree,

    #[error("Merkle proof path too short")]
    MerkleProofPathTooShort,
}
