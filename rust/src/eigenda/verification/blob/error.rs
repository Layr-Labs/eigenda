use std::num::TryFromIntError;

use rust_kzg_bn254_primitives::errors::KzgError;
use thiserror::Error;

#[derive(Debug, Error, PartialEq)]
pub enum BlobVerificationError {
    #[error("Blob is too small ({0} bytes), it is shorter than the 32 byte header")]
    BlobTooSmallForHeader(u32),

    #[error("Blob is too small ({0} bytes), it can't hold header (32 bytes) + payload ({0} bytes)")]
    BlobTooSmallForHeaderAndPayload(usize),

    #[error("Blob length does not fit into a u32 variable: {0}")]
    BlobTooLarge(#[from] TryFromIntError),

    #[error("Blob with length {0} exceeds the certificate's commitment length of {1}")]
    BlobLargerThanCommitmentLength(usize, usize),

    #[error("Commitment length ({0}) not power of two")]
    CommitmentLengthNotPowerOfTwo(u32),

    #[error("Payload length ({0}) larger than upper bound ({1})")]
    PayloadLengthLargerThanUpperBound(u32, u32),

    #[error("Blob's trailing bytes should all be 0x0")]
    NonZeroTrailingBytes,

    #[error("Invalid kzg commitment")]
    InvalidKzgCommitment,

    #[error("Kzg error: {0}")]
    WrapKzgError(#[from] KzgError),
}
