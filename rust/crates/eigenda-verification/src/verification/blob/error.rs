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
    /// Encoded payload length is not a power of two
    #[error("Encoded payload length ({0} bytes) is not a power of two")]
    EncodedPayloadLengthNotPowerOfTwo(usize),

    /// Invalid guard byte in header or symbol
    #[error("Invalid guard byte at offset {offset}: expected 0x00, found 0x{found:02x}")]
    EncodedPayloadInvalidGuardByte {
        /// Offset where the invalid guard byte was found
        offset: usize,
        /// The actual byte value found
        found: u8,
    },

    /// Invalid version byte in header
    #[error("Invalid version byte: expected 0x00, found 0x{0:02x}")]
    EncodedPayloadHeaderInvalidVersion(u8),

    /// Invalid padding in header (bytes 6-31 must be zeros)
    #[error("Invalid header padding at offset {offset}: expected 0x00, found 0x{found:02x}")]
    EncodedPayloadInvalidHeaderPadding {
        /// Offset where the invalid padding was found
        offset: usize,
        /// The actual byte value found
        found: u8,
    },

    /// Invalid padding bytes (all padding must be zeros)
    #[error(
        "Invalid encoded payload padding at offset {offset}: expected 0x00, found 0x{found:02x}"
    )]
    EncodedPayloadInvalidPadding {
        /// Offset where the invalid padding was found
        offset: usize,
        /// The actual byte value found
        found: u8,
    },

    /// EncodedPayload is too small to contain the required 32-byte header
    #[error("EncodedPayload is too small ({0} bytes), it is shorter than the 32 byte header")]
    EncodedPayloadTooSmallForHeader(usize),

    /// EncodedPayload is too small to contain the claimed payload length
    #[error(
        "EncodedPayload is too small ({encoded_payload_bytes_len} bytes), it can't hold claimed encoded payload length ({claimed_encoded_payload_bytes_len} bytes)"
    )]
    EncodedPayloadTooSmallForHeaderAndPayload {
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
