//! Implementation of EigenDA blob verification as defined in
//! [EigenDA specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation)!

pub mod error;
pub mod srs;

use std::sync::LazyLock;

use crate::eigenda::cert::{BlobCommitment, G1Point};
use ark_bn254::G1Affine;
use ark_serialize::CanonicalDeserialize;
use rust_kzg_bn254_primitives::blob::Blob;
use rust_kzg_bn254_prover::{kzg::KZG, srs::SRS};

use crate::eigenda::verification::{
    blob::{error::BlobVerificationError, srs::SerializableSRS},
    cert::types::conversions::IntoExt,
};

const FIELD_ELEMENT_GUARD_BYTE: u8 = 0;
const BYTES_PER_SYMBOL: usize = 32;
const BYTES_PER_CHUNK: usize = BYTES_PER_SYMBOL - 1;
const HEADER_SYMBOLS_LEN: usize = 1;
const HEADER_BYTES_LEN: usize = HEADER_SYMBOLS_LEN * BYTES_PER_SYMBOL;
const VERSION: u8 = 0;
const SRS_BYTES: &[u8] = include_bytes!(concat!(env!("OUT_DIR"), "/srs.bin"));

pub static SRS: LazyLock<SRS> = LazyLock::new(|| {
    SerializableSRS::deserialize_compressed(SRS_BYTES)
        .expect("Failed to deserialize precomputed SRS")
        .into()
});

/// Verifies that `blob` passes all the checks defined in
/// [EigenDA specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation)!
pub fn verify(
    blob_commitment: &BlobCommitment,
    payload: &[u8],
) -> Result<(), BlobVerificationError> {
    let blob = payload_into_blob(payload)?;
    let blob_symbols_len = blob.len() / BYTES_PER_SYMBOL;

    let BlobCommitment {
        commitment, length, ..
    } = blob_commitment;

    verify_blob_symbols_len_against_commitment(blob_symbols_len, *length as usize)?;
    verify_commitment_len_is_power_of_two(*length)?;
    verify_kzg_commitment(&blob, *commitment)?;

    Ok(())
}

/// [EigenDA specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation)!
///
/// 1. Verify that received blob length is <= the length in the cert's BlobCommitment
///
/// We don't check for equality (blob_len == commitment_len) because trailing 0x00s
/// may have been removed in transmission and that's acceptable
fn verify_blob_symbols_len_against_commitment(
    blob_symbols_len: usize,
    commitment_symbols_len: usize,
) -> Result<(), BlobVerificationError> {
    use BlobVerificationError::*;

    (blob_symbols_len <= commitment_symbols_len)
        .then_some(())
        .ok_or(BlobLargerThanCommitmentLength(
            blob_symbols_len,
            commitment_symbols_len,
        ))
}

/// [EigenDA specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation)!
///
/// 2. Verify that the blob length claimed in the BlobCommitment is greater than 0
/// 3. Verify that the blob length claimed in the BlobCommitment is a power of two
///
/// Since 0 is not a power of two, verification 3. subsumes 2.
#[inline]
fn verify_commitment_len_is_power_of_two(
    commitment_symbols_len: u32,
) -> Result<(), BlobVerificationError> {
    use BlobVerificationError::*;

    commitment_symbols_len
        .is_power_of_two()
        .then_some(())
        .ok_or(CommitmentLengthNotPowerOfTwo(commitment_symbols_len))
}

/// [EigenDA specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation)!
/// 4. Verify that the payload length claimed in the encoded payload header is <= the
///    maximum permissible payload length, as calculated from the length in the
///    BlobCommitment
///
///   1. The maximum permissible payload length is computed by looking at the
///      claimed blob length, and determining how many bytes would remain if you were
///      to remove the encoding which is performed when converting a payload into a
///      encodedPayload. This presents an upper bound for payload length: e.g.
///      "If the payload were any bigger than x, then the process of converting it
///      to an encodedPayload would have yielded a blob of larger size than claimed"
/// 5. If the bytes received for the blob are longer than necessary to convey
///    the payload, as determined by the claimed payload length, then verify that
///    all extra bytes are 0x0
///
///   1. Due to how padding of a blob works, it's possible that there may be
///      trailing 0x0 bytes, but there shouldn't be any trailing bytes that aren't
///      equal to 0x0
///
/// Requirements 4. and 5. are satisfied by construction
fn payload_into_blob(payload: &[u8]) -> Result<Blob, BlobVerificationError> {
    let header = construct_header(payload)?;
    let header_bytes_len = header.len();

    let payload = pad_payload(payload);
    let payload_bytes_len = payload.len();

    let blob_symbols_len =
        ((payload_bytes_len + header_bytes_len) / BYTES_PER_SYMBOL).next_power_of_two();
    let blob_bytes_len = blob_symbols_len * BYTES_PER_SYMBOL;

    let mut blob = vec![0; blob_bytes_len];
    blob[..header_bytes_len].copy_from_slice(&header);
    let from = header_bytes_len;
    let to = header_bytes_len + payload_bytes_len;
    blob[from..to].copy_from_slice(&payload);

    Ok(Blob::new(&blob))
}

/// Constructs the blob header according to EigenDA specification.
///
/// The header is a 32-byte structure with the following format:
/// - Byte 0: Field element guard byte (0x00)
/// - Byte 1: Version byte (0x00)
/// - Bytes 2-5: Payload length as big-endian u32
/// - Bytes 6-31: Zero padding
///
/// # Arguments
/// * `payload` - The payload data to encode in the header
///
/// # Returns
/// * `Result<[u8; HEADER_BYTES_LEN], BlobVerificationError>` - The constructed header or an error if payload is too large
fn construct_header(payload: &[u8]) -> Result<[u8; HEADER_BYTES_LEN], BlobVerificationError> {
    let mut header = [0; HEADER_BYTES_LEN];
    header[0] = FIELD_ELEMENT_GUARD_BYTE;
    header[1] = VERSION;
    let payload_len: u32 = payload.len().try_into()?;
    header[2..6].copy_from_slice(&payload_len.to_be_bytes());
    Ok(header)
}

/// Pads and encodes payload data into symbols for blob creation.
///
/// This function transforms raw payload data into a format suitable for EigenDA blob encoding
/// by splitting it into chunks and adding field element guard bytes. Each 31-byte chunk of
/// payload data becomes a 32-byte symbol with a guard byte prefix.
///
/// # Process
/// 1. Divides payload into 31-byte chunks (BYTES_PER_CHUNK)
/// 2. Pads the last chunk with zeros if needed
/// 3. Converts each chunk into a 32-byte symbol by prepending a field element guard byte
///
/// # Arguments
/// * `payload` - The raw payload data to pad and encode
///
/// # Returns
/// * `Vec<u8>` - The padded and encoded payload as a vector of symbols
fn pad_payload(payload: &[u8]) -> Vec<u8> {
    let chunks = payload.len().div_ceil(BYTES_PER_CHUNK);

    let chunk_bytes_len = chunks * BYTES_PER_CHUNK;
    let mut src = Vec::with_capacity(chunk_bytes_len);
    src.extend_from_slice(payload);
    src.resize(chunk_bytes_len, 0u8);

    let symbol_bytes_len = chunks * BYTES_PER_SYMBOL;
    let mut dst = vec![0; symbol_bytes_len];

    for (src, dst) in src
        .chunks_exact(BYTES_PER_CHUNK)
        .zip(dst.chunks_exact_mut(BYTES_PER_SYMBOL))
    {
        dst[0] = FIELD_ELEMENT_GUARD_BYTE;
        dst[1..].copy_from_slice(src);
    }

    dst
}

/// [EigenDA specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation)!
///
/// 6. Verify the KZG commitment. This can either be done:
///
///   1. directly: recomputing the commitment using SRS points and checking
///      that the two commitments match (this is the current implemented way)
///   2. indirectly: verifying a point opening using Fiat-Shamir (see this [issue](https://github.com/Layr-Labs/eigenda/issues/1037))
///
/// > the PR is still open so we don't have the data for option 2.
fn verify_kzg_commitment(
    blob: &Blob,
    claimed_commitment: G1Point,
) -> Result<(), BlobVerificationError> {
    use BlobVerificationError::*;

    let computed_commitment = KZG::new().commit_blob(blob, &SRS)?;

    let claimed_commitment: G1Affine = claimed_commitment.into_ext();

    (computed_commitment == claimed_commitment)
        .then_some(())
        .ok_or(InvalidKzgCommitment)
}

#[cfg(test)]
mod test {
    use std::str::FromStr;

    use ark_bn254::{Fq, G1Affine, G2Affine};
    use rust_kzg_bn254_primitives::errors::KzgError;

    use crate::eigenda::{
        cert::{BlobCommitment, G1Point},
        verification::{
            blob::{
                error::BlobVerificationError::*, payload_into_blob, srs::POINTS_TO_LOAD, verify,
                verify_blob_symbols_len_against_commitment, verify_commitment_len_is_power_of_two,
                verify_kzg_commitment,
            },
            cert::types::conversions::IntoExt,
        },
    };

    #[test]
    fn verify_succeeds_with_known_commitment() {
        let payload = [123; 512];

        let known_commitment = G1Affine::new_unchecked(
            Fq::from_str(
                "14744258532267160547483505594354502788777214273862365248297251133183543768320",
            )
            .unwrap(),
            Fq::from_str(
                "14747463945321045950305747275042450369190644326153769248149140572576072465547",
            )
            .unwrap(),
        )
        .into_ext();

        let blob_commitment = BlobCommitment {
            commitment: known_commitment,
            length_commitment: G2Affine::default().into_ext(),
            length_proof: G2Affine::default().into_ext(),
            length: 32 + 32,
        };

        assert_eq!(verify(&blob_commitment, &payload), Ok(()));
    }

    #[test]
    fn verify_fails_when_payload_is_too_large() {
        let payload = vec![42u8; u32::MAX as usize + 1]; // 4GB vec

        let blob_commitment = BlobCommitment {
            commitment: G1Affine::default().into_ext(),
            length_commitment: G2Affine::default().into_ext(),
            length_proof: G2Affine::default().into_ext(),
            length: Default::default(),
        };

        let result = verify(&blob_commitment, &payload);
        assert!(matches!(result, Err(BlobTooLarge(_))));
    }

    #[test]
    fn test_verify_blob_symbols_len_against_commitment() {
        assert_eq!(verify_blob_symbols_len_against_commitment(42, 43), Ok(()));
        assert_eq!(verify_blob_symbols_len_against_commitment(42, 42), Ok(()));
        assert_eq!(
            verify_blob_symbols_len_against_commitment(42, 41),
            Err(BlobLargerThanCommitmentLength(42, 41))
        );
    }

    #[test]
    fn test_verify_commitment_symbols_len_is_power_of_two() {
        assert_eq!(verify_commitment_len_is_power_of_two(0b1000), Ok(()));
        assert_eq!(
            verify_commitment_len_is_power_of_two(0b0111),
            Err(CommitmentLengthNotPowerOfTwo(0b0111))
        );
    }

    #[test]
    fn verify_kzg_commitment_fails_when_payload_is_too_big_for_srs() {
        const LEN: usize = (POINTS_TO_LOAD as usize + 1) * 32;
        let payload = [42u8; LEN];
        let blob = payload_into_blob(&payload).unwrap();
        let claimed_commitment: G1Point = G1Affine::default().into_ext();

        assert_eq!(
            verify_kzg_commitment(&blob, claimed_commitment),
            Err(WrapKzgError(KzgError::SerializationError(
                "polynomial length is not correct".into()
            )))
        );
    }
}
