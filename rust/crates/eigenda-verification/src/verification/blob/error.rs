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
    /// encoded payload decoding error
    #[error("cannot decode an encoded payload")]
    DecodingError(#[from] EncodedPayloadDecodingError),

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

/// List of error can happen during decoding an encoded payload
#[derive(Debug, thiserror::Error, PartialEq)]
pub enum EncodedPayloadDecodingError {
    /// the input encoded payload has wrong size
    #[error(
        "invalid number of bytes in the encoded payload {0}, that is not multiple of bytes per field element"
    )]
    InvalidLengthEncodedPayload(u64),
    /// encoded payload must contain a power of 2 number of field elements
    #[error(
        "encoded payload must be a power of 2 field elements (32 bytes chunks), but got {0} field elements"
    )]
    InvalidPowerOfTwoLength(usize),
    /// encoded payload header validation error
    #[error("encoded payload header first byte must be 0x00, but got {0:#04x}")]
    InvalidHeaderFirstByte(u8),
    /// encoded payload too short for header
    #[error(
        "encoded payload is too small ({0} bytes), it is shorter than the 32 byte header required"
    )]
    EncodedPayloadTooShortForHeader(
        /// Actual payload length
        usize,
    ),
    /// unknown encoded payload header version
    #[error("unknown encoded payload header version: {0}")]
    UnknownEncodingVersion(u8),
    /// length of unpadded data is less than claimed in header
    #[error(
        "length of unpadded data {actual} is less than length claimed in encoded payload header {claimed}"
    )]
    DecodedPayloadBodyTooShort {
        /// Actual decoded body length that potentially has padding
        actual: usize,
        /// Claimed length from header
        claimed: u32,
    },
    /// every multiple 32 bytes for storing a field element requires the first byte to be zero
    #[error("non-zero byte encountered in the first byte of multiples of 32 bytes: {0}")]
    InvalidFirstByteFieldElementPadding(u8),
    /// padding are applied to the encoded payload body to ensure encoded length is power of 2, padding must be 0
    #[error("non-zero padding byte encountered in the encoded payload body: {0}")]
    InvalidEncodedPayloadBodyPadding(u8),
    /// padding are applied to the encoded payload header to ensure the header takes 32 bytes, padding must be 0
    #[error("non-zero padding byte encountered in the encoded payload header: {0}")]
    InvalidEncodedPayloadHeaderPadding(u8),
}
