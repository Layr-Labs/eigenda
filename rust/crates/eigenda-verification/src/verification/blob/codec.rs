//! # EigenDA Payload Encoding/Decoding
//!
//! This module implements the EigenDA payload encoding and decoding functionality according to the
//! [EigenDA specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#decoding-an-encoded-payload).
//!
//! ## Overview
//!
//! EigenDA stores arbitrary data as encoded payloads that undergo a specific [encoding process](https://layr-labs.github.io/eigenda/assets/integration/payload-to-blob-encoding.png):
//! 1. Raw payload data is prefixed with a header containing metadata
//! 2. The data is split into 31-byte chunks and each chunk is prefixed with a guard byte
//! 3. The resulting encoded payload is padded to a power-of-two length for cryptographic operations
//!
//! ## Encoded Payload Structure
//!
//! | Header | Encoded Payload |
//! |--------|-----------------|
//! | Header (32 bytes) | Symbol 1 (32 bytes) |
//! |                   | Symbol 2 (32 bytes) |
//! |                   | ... |
//! |                   | Symbol N (32 bytes) |
//! |                   | 0-Padding |
//!
//! ### Header Format (32 bytes)
//!
//! | Byte | 0 | 1 | 2-5 | 6-31 |
//! |------|---|---|-----|------|
//! | Field | Guard Byte | Version | Payload Length | Zero Padding |
//! | Value | 0 | 0 | Big-endian u32 | 0x00... |
//!
//! ### Symbol Format (32 bytes each)
//!
//! | Byte | 0 | 1-31 |
//! |------|---|------|
//! | Field | Guard Byte | Payload Data (31 bytes max) |
//! | Value | 0 | raw payload chunk + 0-padding |
//!
//! ## Notes
//!
//! - All symbols are guaranteed to be valid BN254 field elements
//! - **Version 0**: Current specification (only version supported)
//! - **Endianness**: Big-endian encoding
//! - **Field**: BN254 elliptic curve field (order ≈ 2^254)

use crate::verification::blob::{BlobVerificationError, EncodedPayloadDecodingError};

/// Size of each symbol in bytes.
///
/// EigenDA organizes data into 32-byte symbols that are compatible with BN254
/// field elements. Each symbol contains 1 guard byte + 31 bytes of payload data.
pub const BYTES_PER_SYMBOL: usize = 32;

/// Size of the payload data portion within each symbol.
///
/// Since each symbol is 32 bytes total and requires 1 guard byte, the remaining
/// 31 bytes are available for actual payload data.
pub const BYTES_PER_CHUNK: usize = BYTES_PER_SYMBOL - 1;

/// Number of symbols used for the encoded payload header.
///
/// The header is exactly one symbol (32 bytes) containing metadata about the encoded payload.
pub const HEADER_SYMBOLS_LEN: usize = 1;

/// Size of the encoded payload header in bytes.
///
/// The header is always exactly 32 bytes, containing the guard byte, version,
/// payload length, and zero padding.
pub const HEADER_BYTES_LEN: usize = HEADER_SYMBOLS_LEN * BYTES_PER_SYMBOL;

/// The PAYLOAD_ENCODING_VERSION_0 requires payload to be encoded as follows
/// - begin with 32 byte header = [0x00, version byte 0, uint32 len of data, 0x00, 0x00,..., 0x00]
/// - followed by the encoded data [0x00, 31 bytes of data, 0x00, 31 bytes of data,...]
pub const PAYLOAD_ENCODING_VERSION_0: u8 = 0x0;

/// Extracts the raw payload from an EigenDA encoded payload.
///
/// This function reverses the encoding process performed by [`encode_raw_payload`], parsing
/// the encoded payload to recover the raw payload data. It performs strict validation
/// of the encoded payload format according to the EigenDA specification.
///
/// # Arguments
///
/// * `encoded payload` - A slice containing the complete encoded data
///
/// # Returns
///
/// * `Ok(Vec<u8>)` - The raw payload data
/// * Err([EncodedPayloadDecodingError]) - if some encoding invariants are violated
pub fn decode_encoded_payload(encoded_payload: &[u8]) -> Result<Vec<u8>, BlobVerificationError> {
    // Check length invariant
    check_len_invariant(encoded_payload)?;

    // Decode header to get claimed payload length
    let payload_len_in_header = decode_header(encoded_payload)?;

    // Decode payload using the helper method
    decode_payload(encoded_payload, payload_len_in_header)
}

/// Checks whether the encoded payload satisfies its length invariant.
/// EncodedPayloads must contain a power of 2 number of Field Elements, each of length 32.
/// This means the only valid encoded payloads have byte lengths of 32, 64, 128, 256, etc.
///
/// Note that this function only checks the length invariant, meaning that it doesn't check that
/// the 32 byte chunks are valid bn254 elements.
fn check_len_invariant(encoded_payload: &[u8]) -> Result<(), BlobVerificationError> {
    // this check is redundant since 0 is not a valid power of 32, but we keep it for clarity.
    if encoded_payload.len() < HEADER_BYTES_LEN {
        return Err(
            EncodedPayloadDecodingError::EncodedPayloadTooShortForHeader(encoded_payload.len())
                .into(),
        );
    }

    if encoded_payload.len() % BYTES_PER_SYMBOL != 0 {
        return Err(EncodedPayloadDecodingError::InvalidLengthEncodedPayload(
            encoded_payload.len() as u64,
        )
        .into());
    }

    // Check encoded payload has a power of two number of field elements
    let num_field_elements = encoded_payload.len() / BYTES_PER_SYMBOL;
    if !num_field_elements.is_power_of_two() {
        return Err(
            EncodedPayloadDecodingError::InvalidPowerOfTwoLength(num_field_elements).into(),
        );
    }
    Ok(())
}

/// Validates the header (first field element = 32 bytes) of the encoded payload,
/// and returns the claimed length of the payload if the header is valid.
fn decode_header(encoded_payload: &[u8]) -> Result<u32, BlobVerificationError> {
    if encoded_payload.len() < HEADER_BYTES_LEN {
        return Err(
            EncodedPayloadDecodingError::EncodedPayloadTooShortForHeader(encoded_payload.len())
                .into(),
        );
    }
    // this ensures the header 32 bytes is a valid field element
    if encoded_payload[0] != 0x00 {
        return Err(EncodedPayloadDecodingError::InvalidHeaderFirstByte(encoded_payload[0]).into());
    }
    let payload_length = match encoded_payload[1] {
        version if version == PAYLOAD_ENCODING_VERSION_0 => u32::from_be_bytes([
            encoded_payload[2],
            encoded_payload[3],
            encoded_payload[4],
            encoded_payload[5],
        ]),
        version => {
            return Err(EncodedPayloadDecodingError::UnknownEncodingVersion(version).into());
        }
    };

    // all the remaining bytes in the payload header must be zero
    for b in &encoded_payload[6..HEADER_BYTES_LEN] {
        if *b != 0x00 {
            return Err(EncodedPayloadDecodingError::InvalidEncodedPayloadHeaderPadding(*b).into());
        }
    }

    Ok(payload_length)
}

/// Decodes the payload from the encoded payload bytes.
/// Removes internal padding and extracts the payload data based on the claimed length.
fn decode_payload(
    encoded_payload: &[u8],
    payload_len: u32,
) -> Result<Vec<u8>, BlobVerificationError> {
    let body = &encoded_payload[HEADER_BYTES_LEN..];

    // Decode the body by removing internal 0 byte padding (0x00 initial byte for every 32 byte chunk)
    // this ensures every 32 bytes is a valid field element
    let mut decoded_body = check_and_remove_zero_padding_for_field_elements(body)?;

    // data length is checked when constructing an encoded payload. If this error is encountered, that means there
    // must be a flaw in the logic at construction time (or someone was bad and didn't use the proper construction methods)
    if decoded_body.len() < payload_len as usize {
        return Err(EncodedPayloadDecodingError::DecodedPayloadBodyTooShort {
            actual: decoded_body.len(),
            claimed: payload_len,
        }
        .into());
    }

    for b in &decoded_body[payload_len as usize..] {
        if *b != 0x00 {
            return Err(EncodedPayloadDecodingError::InvalidEncodedPayloadBodyPadding(*b).into());
        }
    }

    decoded_body.truncate(payload_len as usize);
    Ok(decoded_body)
}

/// check_and_remove_zero_padding_for_field_elements checks if the first byte of every mulitple of 32 bytes is 0x00,
/// it enforces the spec in <https://layr-labs.github.io/eigenda/integration/spec/3-data-structs.html#encoding-payload-version-0x0>
/// then the function returns bytes with the zero-padding bytes removed.
/// this ensures every multiple of 32 bytes is a valid field element
fn check_and_remove_zero_padding_for_field_elements(
    encoded_body: &[u8],
) -> Result<Vec<u8>, BlobVerificationError> {
    if encoded_body.len() % BYTES_PER_SYMBOL != 0 {
        return Err(EncodedPayloadDecodingError::InvalidLengthEncodedPayload(
            encoded_body.len() as u64
        )
        .into());
    }

    let num_field_elements = encoded_body.len() / BYTES_PER_SYMBOL;
    let mut decoded_body = Vec::with_capacity(num_field_elements * 31);
    for chunk in encoded_body.chunks_exact(BYTES_PER_SYMBOL) {
        if chunk[0] != 0x00 {
            return Err(
                EncodedPayloadDecodingError::InvalidFirstByteFieldElementPadding(chunk[0]).into(),
            );
        }
        decoded_body.extend_from_slice(&chunk[1..32]);
    }
    Ok(decoded_body)
}

#[cfg(any(test, feature = "test-utils"))]
/// Test utilities for blob codec operations
///
/// This module provides helper functions for encoding raw payloads into the
/// EigenDA blob format for use in tests and benchmarks. These utilities are
/// only available when the `test-utils` feature is enabled or during testing.
pub mod tests_utils {
    use crate::verification::blob::BlobVerificationError::{self, *};
    use crate::verification::blob::codec::{BYTES_PER_CHUNK, BYTES_PER_SYMBOL, HEADER_BYTES_LEN};

    /// Guard byte value used to prefix field elements in the EigenDA encoding.
    ///
    /// This byte is prepended to each 31-byte chunk to create 32-byte symbols that
    /// are compatible with the BN254 field arithmetic used in EigenDA's cryptographic
    /// operations. The value 0 ensures that the resulting 32-byte value is always
    /// less than the BN254 field modulus.
    pub const FIELD_ELEMENT_GUARD_BYTE: u8 = 0;

    /// Version byte used in the encoded payload header.
    ///
    /// This indicates the encoding format version. Currently, only version 0 is defined.
    pub const VERSION: u8 = 0;

    /// Encodes a raw payload into an EigenDA-compatible encoded payload format.
    ///
    /// This function transforms arbitrary raw payload data into the standardized EigenDA encoded payload
    /// format, which is designed for efficient storage and cryptographic operations on
    /// the EigenDA network. The resulting encoded payload can be decoded back to the raw
    /// payload using [`decode_encoded_payload`].
    ///
    /// # Process
    ///
    /// 1. **Header Construction**: Creates a 32-byte header containing metadata
    /// 2. **Payload Chunking**: Splits the payload into 31-byte chunks
    /// 3. **Symbol Creation**: Prefixes each chunk with a guard byte to form 32-byte symbols
    /// 4. **Power-of-Two Padding**: Expands the encoded payload to the next power-of-two size
    /// 5. **Zero Padding**: Fills unused space with zero bytes
    ///
    /// # Arguments
    ///
    /// * `raw_payload` - A slice containing the raw data to encode
    ///
    /// # Returns
    ///
    /// * `Ok(Vec<u8>)` - The encoded payload data with power-of-two size
    /// * `Err(BlobVerificationError)` - Error conditions:
    ///   - [`BlobTooLarge`](BlobVerificationError::BlobTooLarge) if payload exceeds `u32::MAX` bytes
    ///
    /// # Encoded payload Structure
    ///
    /// The resulting encoded payload has this structure:
    /// ```text
    /// [Header: 32 bytes][Encoded Payload: variable][Zero Padding: to power of 2]
    /// ```
    ///
    /// Where the encoded payload consists of symbols:
    /// ```text
    /// [Guard:1][Data:31][Guard:1][Data:31]...[Guard:1][Data+Pad:31]
    /// ```
    ///
    /// # Notes
    ///
    /// This function satisfies requirements 4 and 5 from the
    /// [EigenDA specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation)
    /// by construction:
    /// - The payload length in the header provides an upper bound for payload size validation
    /// - All padding bytes are guaranteed to be zero
    #[cfg(any(test, feature = "test-utils"))]
    pub fn encode_raw_payload(raw_payload: &[u8]) -> Result<Vec<u8>, BlobVerificationError> {
        let header = construct_header(raw_payload)?;

        let padded_payload = pad_raw_payload(raw_payload)?;
        let padded_payload_bytes_len = padded_payload.len();

        let encoded_payload_len = HEADER_BYTES_LEN
            .checked_add(padded_payload_bytes_len)
            .ok_or(Overflow)?;

        let encoded_payload_symbols_len = encoded_payload_len
            .div_ceil(BYTES_PER_SYMBOL)
            .checked_next_power_of_two()
            .ok_or(Overflow)?;

        let encoded_payload_bytes_len = encoded_payload_symbols_len
            .checked_mul(BYTES_PER_SYMBOL)
            .ok_or(Overflow)?;

        let mut encoded_payload = vec![0; encoded_payload_bytes_len];
        encoded_payload[..HEADER_BYTES_LEN].copy_from_slice(&header);
        encoded_payload[HEADER_BYTES_LEN..encoded_payload_len].copy_from_slice(&padded_payload);

        Ok(encoded_payload)
    }

    /// Constructs the 32-byte blob header according to EigenDA specification.
    ///
    /// The header contains essential metadata about the blob and follows a strict format
    /// to ensure compatibility with EigenDA's cryptographic operations and verification
    /// processes.
    ///
    /// # Header Layout
    ///
    /// | Offset | Size | Field | Description |
    /// |--------|------|-------|-------------|
    /// | 0 | 1 | Guard Byte | 0x00 (field element guard) |
    /// | 1 | 1 | Version | 0x00 (format version) |
    /// | 2-5 | 4 | Payload Length | Big-endian u32 (raw payload size) |
    /// | 6-31 | 26 | Padding | 0x00... (zero padding) |
    ///
    /// # Implementation Details
    ///
    /// - **Guard Byte**: Ensures the header forms a valid BN254 field element
    /// - **Version**: Future-proofs the format (currently only version 0 exists)
    /// - **Length Encoding**: Big-endian u32 supports payloads up to 4GB
    /// - **Zero Padding**: Guarantees the header is exactly 32 bytes
    ///
    /// # Arguments
    ///
    /// * `raw_payload` - Slice containing the raw payload data to encode metadata for
    ///
    /// # Returns
    ///
    /// * `Ok([u8; 32])` - The constructed header bytes
    /// * `Err(BlobVerificationError::BlobTooLarge)` - If payload length exceeds `u32::MAX`
    pub fn construct_header(
        raw_payload: &[u8],
    ) -> Result<[u8; HEADER_BYTES_LEN], BlobVerificationError> {
        let mut header = [0; HEADER_BYTES_LEN];
        header[0] = FIELD_ELEMENT_GUARD_BYTE;
        header[1] = VERSION;
        let raw_payload_len: u32 = raw_payload.len().try_into()?;
        header[2..6].copy_from_slice(&raw_payload_len.to_be_bytes());
        Ok(header)
    }

    /// Transforms raw payload data into field element symbols for cryptographic operations.
    ///
    /// This function is a critical component of the EigenDA encoding process that converts
    /// arbitrary payload data into symbols compatible with BN254 field arithmetic. Each
    /// symbol is exactly 32 bytes and forms a valid field element.
    ///
    /// # Transformation Process
    ///
    /// 1. **Chunking**: Divides payload into 31-byte chunks (maximum data per symbol)
    /// 2. **Padding**: Extends the last chunk to 31 bytes with zero bytes if needed
    /// 3. **Symbol Creation**: Prepends each chunk with a guard byte (0x00) to form 32-byte symbols
    /// 4. **Field Element Guarantee**: Each symbol is guaranteed to be < BN254 field modulus
    ///
    /// # Symbol Structure
    ///
    /// | Byte | Content |
    /// |------|---------|
    /// | 0 | Guard (0x00) |
    /// | 1-31 | Payload Data (padded with zeros if needed) |
    ///
    /// # Mathematical Properties
    ///
    /// - Each 32-byte symbol represents a value < 2^255 (BN254 field modulus ≈ 2^254)
    /// - Guard byte ensures 0 ≤ symbol_value < BN254_MODULUS
    /// - Enables efficient polynomial operations in cryptographic proofs
    ///
    /// # Arguments
    ///
    /// * `raw_payload` - Slice containing the raw data to transform into symbols
    ///
    /// # Returns
    ///
    /// * `Ok(Vec<u8>)` - Encoded symbols as a flat byte vector
    ///   - Length: `ceil(payload.len() / 31) * 32` bytes
    ///   - Empty payload returns empty vector (0 symbols)
    /// * `Err(BlobVerificationError::Overflow)` - If arithmetic operations overflow
    ///
    /// # Notes
    ///
    /// The function uses a two-stage approach:
    /// 1. Expand payload to chunk-aligned size with zero padding
    /// 2. Transform chunks into symbols by interleaving guard bytes
    pub fn pad_raw_payload(raw_payload: &[u8]) -> Result<Vec<u8>, BlobVerificationError> {
        let chunks = raw_payload.len().div_ceil(BYTES_PER_CHUNK);

        let chunk_bytes_len = chunks.checked_mul(BYTES_PER_CHUNK).ok_or(Overflow)?;
        let mut src = Vec::with_capacity(chunk_bytes_len);
        src.extend_from_slice(raw_payload);
        src.resize(chunk_bytes_len, 0u8);

        let symbol_bytes_len = chunks.checked_mul(BYTES_PER_SYMBOL).ok_or(Overflow)?;
        let mut dst = vec![0; symbol_bytes_len];

        for (src, dst) in src
            .chunks_exact(BYTES_PER_CHUNK)
            .zip(dst.chunks_exact_mut(BYTES_PER_SYMBOL))
        {
            dst[0] = FIELD_ELEMENT_GUARD_BYTE;
            dst[1..].copy_from_slice(src);
        }

        Ok(dst)
    }
}

#[cfg(test)]
mod tests {
    use crate::verification::blob::codec::tests_utils::{
        FIELD_ELEMENT_GUARD_BYTE, VERSION, construct_header, encode_raw_payload, pad_raw_payload,
    };
    use crate::verification::blob::codec::{
        BYTES_PER_CHUNK, BYTES_PER_SYMBOL, HEADER_BYTES_LEN, decode_encoded_payload,
    };

    #[test]
    fn roundtrip_boundary_cases() {
        // Test critical boundary cases around chunk/symbol boundaries
        for size in [0, 1, 30, 31, 32, 61, 62, 63, 100, 512, 1000, 2048] {
            let raw_payload: Vec<u8> = (0..size).map(|i| (i % 256) as u8).collect();

            let encoded_payload = encode_raw_payload(&raw_payload).unwrap();
            let recovered_payload = decode_encoded_payload(&encoded_payload).unwrap();

            assert_eq!(
                raw_payload, recovered_payload,
                "Failed roundtrip for size {size}",
            );
        }
    }

    #[test]
    fn encoded_payload_structure_properties() {
        let payload = vec![1, 2, 3, 4, 5];
        let encoded_payload = encode_raw_payload(&payload).unwrap();

        assert!(encoded_payload.len().is_power_of_two());

        assert!(encoded_payload.len() >= HEADER_BYTES_LEN + BYTES_PER_SYMBOL);

        assert_eq!(encoded_payload[0], FIELD_ELEMENT_GUARD_BYTE);
        assert_eq!(encoded_payload[1], VERSION);

        let claimed_len = u32::from_be_bytes([
            encoded_payload[2],
            encoded_payload[3],
            encoded_payload[4],
            encoded_payload[5],
        ]);
        assert_eq!(claimed_len, payload.len() as u32);

        for &byte in &encoded_payload[6..HEADER_BYTES_LEN] {
            assert_eq!(byte, 0);
        }
    }

    #[test]
    fn encoded_payload_too_small_for_header() {
        use crate::verification::blob::error::{
            BlobVerificationError, EncodedPayloadDecodingError,
        };

        for size in [1, 2, 5, 9, 16] {
            let small_encoded_payload = vec![0u8; size];
            assert!(matches!(
                decode_encoded_payload(&small_encoded_payload),
                Err(BlobVerificationError::DecodingError(
                    EncodedPayloadDecodingError::EncodedPayloadTooShortForHeader(len)
                )) if len == size
            ),);
        }
    }

    #[test]
    fn claimed_encoded_payload_len_too_small_for_actual_payload() {
        use crate::verification::blob::error::{
            BlobVerificationError, EncodedPayloadDecodingError,
        };

        // Create a properly sized (power of two) encoded payload, but with a header claiming
        // more payload data than is actually present
        let mut encoded_payload = vec![0u8; 128]; // Power of two
        encoded_payload[0] = FIELD_ELEMENT_GUARD_BYTE;
        encoded_payload[1] = VERSION;

        // Claim we have 100 bytes of payload, which would require:
        // - 100 bytes / 31 bytes per chunk = 4 chunks (rounded up)
        // - 4 chunks * 32 bytes per symbol = 128 bytes
        // - Total: 32 (header) + 128 (payload) = 160 bytes
        // But we only have 128 bytes total, so this should fail
        let payload_len: u32 = 100;
        encoded_payload[2..6].copy_from_slice(&payload_len.to_be_bytes());

        assert!(matches!(
            decode_encoded_payload(&encoded_payload),
            Err(BlobVerificationError::DecodingError(
                EncodedPayloadDecodingError::DecodedPayloadBodyTooShort { .. }
            ))
        ));
    }

    #[test]
    fn empty_payload_roundtrip() {
        let encoded_payload = encode_raw_payload(&[]).unwrap();
        let recovered = decode_encoded_payload(&encoded_payload).unwrap();
        assert!(recovered.is_empty());
    }

    #[test]
    fn construct_header_format() {
        for (payload, expected_len) in [
            (vec![], 0u32),
            (vec![1, 2, 3, 4, 5], 5u32),
            (vec![0u8; 1000], 1000u32),
        ] {
            let header = construct_header(&payload).unwrap();

            assert_eq!(header[0], FIELD_ELEMENT_GUARD_BYTE);
            assert_eq!(header[1], VERSION);
            assert_eq!(
                u32::from_be_bytes([header[2], header[3], header[4], header[5]]),
                expected_len
            );

            for &byte in &header[6..] {
                assert_eq!(byte, 0);
            }
        }
    }

    #[test]
    fn pad_empty_payload() {
        let result = pad_raw_payload(&[]).unwrap();
        assert_eq!(result.len(), 0);
    }

    #[test]
    fn pad_single_byte() {
        let payload = vec![42];
        let result = pad_raw_payload(&payload).unwrap();

        assert_eq!(result.len(), BYTES_PER_SYMBOL);
        assert_eq!(result[0], FIELD_ELEMENT_GUARD_BYTE);
        assert_eq!(result[1], 42);

        for &byte in &result[2..] {
            assert_eq!(byte, 0);
        }
    }

    #[test]
    fn pad_exact_chunk_size() {
        let payload = vec![0u8; BYTES_PER_CHUNK];
        let result = pad_raw_payload(&payload).unwrap();

        assert_eq!(result.len(), BYTES_PER_SYMBOL);
        assert_eq!(result[0], FIELD_ELEMENT_GUARD_BYTE);

        assert_eq!(&result[1..], &payload);
    }

    #[test]
    fn pad_multiple_exact_chunks() {
        let payload = vec![0u8; BYTES_PER_CHUNK * 2];
        let result = pad_raw_payload(&payload).unwrap();

        assert_eq!(result.len(), BYTES_PER_SYMBOL * 2);

        assert_eq!(result[0], FIELD_ELEMENT_GUARD_BYTE);
        assert_eq!(result[BYTES_PER_SYMBOL], FIELD_ELEMENT_GUARD_BYTE);

        for (i, &expected_byte) in payload.iter().enumerate() {
            let symbol_idx = i / BYTES_PER_CHUNK;
            let byte_idx = i % BYTES_PER_CHUNK;
            let result_idx = symbol_idx * BYTES_PER_SYMBOL + byte_idx + 1;
            assert_eq!(result[result_idx], expected_byte);
        }
    }

    #[test]
    fn pad_with_partial_chunk() {
        let payload = vec![0u8; BYTES_PER_CHUNK * 2 + 5];
        let result = pad_raw_payload(&payload).unwrap();

        assert_eq!(result.len(), BYTES_PER_SYMBOL * 3);

        for symbol in 0..3 {
            assert_eq!(result[symbol * BYTES_PER_SYMBOL], FIELD_ELEMENT_GUARD_BYTE);
        }

        for (i, &expected_byte) in payload.iter().enumerate() {
            let symbol_idx = i / BYTES_PER_CHUNK;
            let byte_idx = i % BYTES_PER_CHUNK;
            let result_idx = symbol_idx * BYTES_PER_SYMBOL + byte_idx + 1;
            assert_eq!(result[result_idx], expected_byte);
        }

        let last_symbol_start = 2 * BYTES_PER_SYMBOL;
        for i in 6..BYTES_PER_CHUNK {
            assert_eq!(result[last_symbol_start + i + 1], 0);
        }
    }

    #[test]
    fn reject_non_power_of_two_length() {
        use crate::verification::blob::error::{
            BlobVerificationError, EncodedPayloadDecodingError,
        };

        // Test various lengths that aren't multiples of 32 (but are >= 32 bytes)
        for size in [33, 50, 63, 100, 129, 200, 1000] {
            let small_encoded_payload = vec![0u8; size];
            assert!(matches!(
                decode_encoded_payload(&small_encoded_payload),
                Err(BlobVerificationError::DecodingError(
                    EncodedPayloadDecodingError::InvalidLengthEncodedPayload(len)
                )) if len == size as u64
            ));
        }

        // Then multiples of 32 that aren't powers of 2
        for size in [96, 224] {
            let encoded_payload = vec![0u8; size];
            assert!(
                matches!(
                    decode_encoded_payload(&encoded_payload),
                    Err(BlobVerificationError::DecodingError(
                        EncodedPayloadDecodingError::InvalidPowerOfTwoLength(len)
                    )) if len == size / 32
                ),
                "Failed to reject size {size}"
            );
        }
    }

    #[test]
    fn reject_invalid_header_guard_byte() {
        use crate::verification::blob::error::{
            BlobVerificationError, EncodedPayloadDecodingError,
        };

        let mut encoded_payload = vec![0u8; 32]; // Power of two
        encoded_payload[0] = 0x01; // Invalid guard byte
        encoded_payload[1] = VERSION;

        assert!(matches!(
            decode_encoded_payload(&encoded_payload),
            Err(BlobVerificationError::DecodingError(
                EncodedPayloadDecodingError::InvalidHeaderFirstByte(0x01)
            ))
        ));
    }

    #[test]
    fn reject_invalid_version() {
        use crate::verification::blob::error::{
            BlobVerificationError, EncodedPayloadDecodingError,
        };

        let mut encoded_payload = vec![0u8; 32]; // Power of two
        encoded_payload[0] = FIELD_ELEMENT_GUARD_BYTE;
        encoded_payload[1] = 0x01; // Invalid version

        assert!(matches!(
            decode_encoded_payload(&encoded_payload),
            Err(BlobVerificationError::DecodingError(
                EncodedPayloadDecodingError::UnknownEncodingVersion(0x01)
            ))
        ));
    }

    #[test]
    fn reject_invalid_header_padding() {
        use crate::verification::blob::error::{
            BlobVerificationError, EncodedPayloadDecodingError,
        };

        // Test padding at various offsets in the header
        for test_offset in 6..32 {
            let mut encoded_payload = vec![0u8; 32]; // Power of two
            encoded_payload[0] = FIELD_ELEMENT_GUARD_BYTE;
            encoded_payload[1] = VERSION;
            encoded_payload[test_offset] = 0x42; // Invalid padding byte

            assert!(
                matches!(
                    decode_encoded_payload(&encoded_payload),
                    Err(BlobVerificationError::DecodingError(
                        EncodedPayloadDecodingError::InvalidEncodedPayloadHeaderPadding(0x42)
                    ))
                ),
                "Failed to reject invalid header padding at offset {test_offset}"
            );
        }
    }

    #[test]
    fn reject_invalid_symbol_guard_byte() {
        use crate::verification::blob::error::{
            BlobVerificationError, EncodedPayloadDecodingError,
        };

        let payload = vec![1u8; 5];
        let mut encoded_payload = encode_raw_payload(&payload).unwrap();

        // Corrupt the guard byte of the first payload symbol (after the header)
        encoded_payload[HEADER_BYTES_LEN] = 0x99;

        assert!(matches!(
            decode_encoded_payload(&encoded_payload),
            Err(BlobVerificationError::DecodingError(
                EncodedPayloadDecodingError::InvalidFirstByteFieldElementPadding(0x99)
            ))
        ));
    }

    #[test]
    fn reject_invalid_padding_in_last_chunk() {
        use crate::verification::blob::error::{
            BlobVerificationError, EncodedPayloadDecodingError,
        };

        // Create a payload that doesn't fill the last chunk completely
        let payload = vec![1u8; 5]; // 5 bytes will leave padding in the last chunk
        let mut encoded_payload = encode_raw_payload(&payload).unwrap();

        // The encoded structure is:
        // [Header: 32 bytes][Symbol: 32 bytes (guard + 31 data/padding)]
        // The 5 bytes of payload go into the symbol, leaving bytes 6-31 as padding
        // Corrupt one of the padding bytes in the symbol
        encoded_payload[HEADER_BYTES_LEN + 1 + 10] = 0xAB; // Offset: header + guard + 10th byte

        assert!(matches!(
            decode_encoded_payload(&encoded_payload),
            Err(BlobVerificationError::DecodingError(
                EncodedPayloadDecodingError::InvalidEncodedPayloadBodyPadding(0xAB)
            ))
        ));
    }

    #[test]
    fn reject_invalid_trailing_padding() {
        use crate::verification::blob::error::{
            BlobVerificationError, EncodedPayloadDecodingError,
        };

        let payload = vec![1u8; 5];
        let mut encoded_payload = encode_raw_payload(&payload).unwrap();

        // The encoded_payload is padded to a power of two
        // Corrupt a trailing padding byte (after the claimed encoded payload length)
        let last_idx = encoded_payload.len() - 1;
        encoded_payload[last_idx] = 0xCD;

        assert!(
            matches!(
                decode_encoded_payload(&encoded_payload),
                Err(BlobVerificationError::DecodingError(
                    EncodedPayloadDecodingError::InvalidEncodedPayloadBodyPadding(0xCD)
                ))
            ),
            "Should reject non-zero trailing padding"
        );
    }

    #[test]
    fn accept_valid_encoded_payload_with_various_padding() {
        // Test that valid encoded payloads with different amounts of padding work correctly
        for payload_size in [1, 5, 31, 32, 62, 100] {
            let payload = vec![0xFFu8; payload_size];
            let encoded_payload = encode_raw_payload(&payload).unwrap();
            let decoded = decode_encoded_payload(&encoded_payload).unwrap();
            assert_eq!(payload, decoded, "Failed for payload size {payload_size}");
        }
    }
}

#[cfg(all(test, feature = "arbitrary"))]
mod proptests {
    use proptest::prelude::*;

    use crate::verification::blob::codec::decode_encoded_payload;
    use crate::verification::blob::codec::tests_utils::encode_raw_payload;

    proptest! {
        #[test]
        fn prop_roundtrip_encode_decode_random_payloads(
            payload in prop::collection::vec(any::<u8>(), 0..=8192)
        ) {
            let encoded_payload = encode_raw_payload(&payload)?;
            let recovered_payload = decode_encoded_payload(&encoded_payload)?;
            prop_assert_eq!(payload, recovered_payload);
        }
    }
}
