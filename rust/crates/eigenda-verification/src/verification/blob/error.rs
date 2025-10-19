//! Error types for EigenDA blob verification
//!
//! This module defines all possible errors that can occur during blob
//! verification against KZG commitments.

use std::num::TryFromIntError;

use rust_kzg_bn254_primitives::errors::KzgError;
use thiserror::Error;

/// Errors that can occur during blob verification
#[derive(Debug, Error, PartialEq)]
pub enum BlobVerificationError {
    /// Blob is too small to contain the required 32-byte header
    #[error("Blob is too small ({0} bytes), it is shorter than the 32 byte header")]
    BlobTooSmallForHeader(usize),

    /// Blob is too small to contain the claimed payload length
    #[error(
        "Blob is too small ({encoded_payload_bytes_len} bytes), it can't hold claimed encoded payload length ({claimed_encoded_payload_bytes_len} bytes)"
    )]
    BlobTooSmallForHeaderAndPayload {
        /// Actual size of the encoded payload in bytes
        encoded_payload_bytes_len: usize,
        /// Claimed size of the encoded payload in bytes
        claimed_encoded_payload_bytes_len: usize,
    },

    /// Blob length exceeds the maximum representable size (u32::MAX)
    #[error("Blob length does not fit into a u32 variable: {0}")]
    BlobTooLarge(#[from] TryFromIntError),

    /// Received blob is larger than the length specified in the certificate
    #[error("Blob with length {0} exceeds the certificate's commitment length of {1}")]
    BlobLargerThanCommitmentLength(usize, usize),

    /// Commitment length is not a power of two (required for KZG)
    #[error("Commitment length ({0}) not power of two")]
    CommitmentLengthNotPowerOfTwo(u32),

    /// KZG commitment verification failed (computed ≠ claimed commitment)
    #[error("Invalid kzg commitment")]
    InvalidKzgCommitment,

    /// Underlying KZG cryptographic library error
    #[error("Kzg error: {0}")]
    KzgError(#[from] KzgError),

    /// Arithmetic overflow occurred during payload processing
    #[error("Arithmetic overflow during payload processing")]
    Overflow,
}
