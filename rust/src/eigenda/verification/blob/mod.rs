//! Implementation of EigenDA blob verification as defined in
//! [EigenDA specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation)!
//!
//! # Header layout:
//!
//! offset:    0        1        2  3  4  5             6 .. 31
//!         [0x00][version][  payload_len_be (u32)  ][   zeros…   ]  => 32 bytes total
//!
//! # Blob layout:
//!
//! [header: 32 bytes][payload: n bytes][padding: m bytes]

// TODO: what is the maximum blob size allowed? That is how many SRS points should be read?

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

const HEADER_LEN: u32 = 32;
const SRS_BYTES: &[u8] = include_bytes!(concat!(env!("OUT_DIR"), "/srs.bin"));

pub static SRS: LazyLock<SRS> = LazyLock::new(|| {
    SerializableSRS::deserialize_compressed(SRS_BYTES)
        .expect("Failed to deserialize precomputed SRS")
        .into()
});

/// Verifies that `blob` passes all the checks defined in
/// [EigenDA specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation)!
pub fn verify(blob_commitment: &BlobCommitment, blob: &[u8]) -> Result<(), BlobVerificationError> {
    use BlobVerificationError::*;

    // TODO: check error shows the blob len
    let blob_len = blob.len().try_into().map_err(BlobTooLarge)?;

    let BlobCommitment {
        commitment, length, ..
    } = blob_commitment;

    verify_blob_len_against_commitment_len(blob_len, *length)?;
    verify_commitment_len_is_power_of_two(*length)?;
    let payload_len = verify_payload_not_greater_than_upper_bound(blob, blob_len)?;
    verify_trailings_bytes_are_all_zero(blob, payload_len)?;
    verify_kzg_commitment(blob, *commitment)?;

    Ok(())
}

/// [EigenDA specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation)!
///
/// 1. Verify that received blob length is <= the length in the cert's BlobCommitment
///
/// We don't check for equality (blob_len == commitment_len) because trailing 0x00s
/// may have been removed in transmission and that's acceptable
#[inline]
fn verify_blob_len_against_commitment_len(
    blob_len: u32,
    commitment_len: u32,
) -> Result<(), BlobVerificationError> {
    use BlobVerificationError::*;

    (blob_len <= commitment_len)
        .then_some(())
        .ok_or(BlobLargerThanCommitmentLength(blob_len, commitment_len))
}

/// [EigenDA specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation)!
///
/// 2. Verify that the blob length claimed in the BlobCommitment is greater than 0
/// 3. Verify that the blob length claimed in the BlobCommitment is a power of two
///
/// Since 0 is not a power of two, verification 3. subsumes 2.
#[inline]
fn verify_commitment_len_is_power_of_two(commitment_len: u32) -> Result<(), BlobVerificationError> {
    use BlobVerificationError::*;

    commitment_len
        .is_power_of_two()
        .then_some(())
        .ok_or(CommitmentLengthNotPowerOfTwo(commitment_len))
}

/// [EigenDA specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation)!
///
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
fn verify_payload_not_greater_than_upper_bound(
    blob: &[u8],
    blob_len: u32,
) -> Result<u32, BlobVerificationError> {
    use BlobVerificationError::*;

    const NON_EMPTY_HEADER_LEN: usize = 6;
    let first_chunk: &[u8; NON_EMPTY_HEADER_LEN] = blob
        .first_chunk::<NON_EMPTY_HEADER_LEN>()
        .ok_or(BlobTooSmallForHeader(blob_len))?;

    const PAYLOAD_BYTE_LEN: usize = 4;
    let be_bytes: [u8; PAYLOAD_BYTE_LEN] = [
        first_chunk[2],
        first_chunk[3],
        first_chunk[4],
        first_chunk[5],
    ];
    let payload_len = u32::from_be_bytes(be_bytes);

    let payload_len_upper_bound = blob_len
        .checked_sub(HEADER_LEN)
        .ok_or(BlobTooSmallForHeader(blob_len))?;

    (payload_len <= payload_len_upper_bound)
        .then_some(payload_len)
        .ok_or(PayloadLengthLargerThanUpperBound(
            payload_len,
            payload_len_upper_bound,
        ))
}

/// [EigenDA specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation)!
///
/// 5. If the bytes received for the blob are longer than necessary to convey
///    the payload, as determined by the claimed payload length, then verify that
///    all extra bytes are 0x0
///
///   1. Due to how padding of a blob works, it's possible that there may be
///      trailing 0x0 bytes, but there shouldn't be any trailing bytes that aren't
///      equal to 0x0
fn verify_trailings_bytes_are_all_zero(
    blob: &[u8],
    payload_len: u32,
) -> Result<(), BlobVerificationError> {
    use BlobVerificationError::*;

    let offset = (HEADER_LEN + payload_len) as usize;
    blob.get(offset..)
        .ok_or(BlobTooSmallForHeaderAndPayload(payload_len))?
        .iter()
        .all(|&byte| byte == 0)
        .then_some(())
        .ok_or(NonZeroTrailingBytes)
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
    blob: &[u8],
    claimed_commitment: G1Point,
) -> Result<(), BlobVerificationError> {
    use BlobVerificationError::*;

    let polynomial = Blob::new(blob).to_polynomial_coeff_form();
    let computed_commitment = KZG::new().commit_coeff_form(&polynomial, &SRS)?;

    let claimed_commitment: G1Affine = claimed_commitment.into_ext();

    (computed_commitment == claimed_commitment)
        .then_some(())
        .ok_or(InvalidKzgCommitment)
}

#[cfg(test)]
mod test {
    use crate::eigenda::cert::{BlobCommitment, G1Point};
    use ark_bn254::{G1Affine, G2Affine};
    use rust_kzg_bn254_primitives::{blob::Blob, errors::KzgError};
    use rust_kzg_bn254_prover::kzg::KZG;

    use crate::eigenda::verification::{
        blob::{
            SRS, error::BlobVerificationError::*, srs::POINTS_TO_LOAD, verify,
            verify_blob_len_against_commitment_len, verify_commitment_len_is_power_of_two,
            verify_kzg_commitment, verify_payload_not_greater_than_upper_bound,
            verify_trailings_bytes_are_all_zero,
        },
        cert::types::conversions::IntoExt,
    };

    #[test]
    fn verify_fails_when_blob_is_too_large() {
        let mut blob = vec![42u8; u32::MAX as usize + 1]; // 4GB vec
        let payload_len = 32u32.to_be_bytes();
        blob[2..6].copy_from_slice(&payload_len);

        let blob_commitment = BlobCommitment {
            commitment: G1Affine::default().into_ext(),
            length_commitment: G2Affine::default().into_ext(),
            length_proof: G2Affine::default().into_ext(),
            length: Default::default(),
        };

        let result = verify(&blob_commitment, &blob);
        assert!(matches!(result, Err(BlobTooLarge(_))));
    }

    #[test]
    fn verify_succeeds() {
        let mut blob = [42u8; 64];
        let payload_len = 32u32.to_be_bytes();
        blob[2..6].copy_from_slice(&payload_len);

        // reproduce the calculation, more flexible than hardcoding the commitment
        let polynomial = Blob::new(&blob).to_polynomial_coeff_form();
        let commitment = KZG::new().commit_coeff_form(&polynomial, &SRS).unwrap();

        let blob_commitment = BlobCommitment {
            commitment: commitment.into_ext(),
            length_commitment: G2Affine::default().into_ext(),
            length_proof: G2Affine::default().into_ext(),
            length: 32 + 32,
        };

        assert_eq!(verify(&blob_commitment, &blob), Ok(()));
    }

    #[test]
    fn test_verify_blob_len_against_commitment_len() {
        assert_eq!(verify_blob_len_against_commitment_len(42, 43), Ok(()));
        assert_eq!(verify_blob_len_against_commitment_len(42, 42), Ok(()));
        assert_eq!(
            verify_blob_len_against_commitment_len(42, 41),
            Err(BlobLargerThanCommitmentLength(42, 41))
        );
    }

    #[test]
    fn test_verify_commitment_len_is_power_of_two() {
        assert_eq!(verify_commitment_len_is_power_of_two(0b1000), Ok(()));
        assert_eq!(
            verify_commitment_len_is_power_of_two(0b0111),
            Err(CommitmentLengthNotPowerOfTwo(0b0111))
        );
    }

    #[test]
    fn verify_payload_not_greater_than_upper_bound_when_payload_is_too_small() {
        assert!(verify_payload_not_greater_than_upper_bound(&[0, 1, 2, 3, 4], 5).is_err());
        assert!(verify_payload_not_greater_than_upper_bound(&[0, 1, 2, 3, 4, 5], 6).is_err());
    }

    #[test]
    fn verify_payload_greater_than_upper_bound_fails_when_blob_shorter_than_header() {
        let blob = &[42u8; 5];
        let blob_len = blob.len() as u32;
        assert_eq!(
            verify_payload_not_greater_than_upper_bound(blob, blob_len),
            Err(BlobTooSmallForHeader(blob_len))
        );

        let blob = &[42u8; 7];
        let blob_len = blob.len() as u32;
        assert_eq!(
            verify_payload_not_greater_than_upper_bound(blob, blob_len),
            Err(BlobTooSmallForHeader(blob_len))
        );
    }

    #[test]
    fn verify_payload_greater_than_upper_bound_fails_when_payload_len_larger_than_upper_bound() {
        // We are claiming a payload length of 33 bytes
        // The header alone occupies 32 bytes
        // So a [u8; 64] blob can't possibly fit header + payload
        let payload_len = 33u32.to_be_bytes();
        let mut blob = [0u8; 64];
        blob[2..6].copy_from_slice(&payload_len);
        let blob_len = blob.len() as u32;
        assert_eq!(
            verify_payload_not_greater_than_upper_bound(&blob, blob_len),
            Err(PayloadLengthLargerThanUpperBound(33, 32))
        );
    }

    #[test]
    fn verify_payload_greater_than_upper_bound_succeeds() {
        // We are claiming a payload length of 32 bytes
        // The header alone occupies 32 bytes
        // So a [u8; 64] blob can exactly fit header + payload
        let payload_len = 32u32.to_be_bytes();
        let mut blob = [0u8; 64];
        blob[2..6].copy_from_slice(&payload_len);
        let blob_len = blob.len() as u32;
        assert_eq!(
            verify_payload_not_greater_than_upper_bound(&blob, blob_len),
            Ok(32)
        );
    }

    #[test]
    fn verify_trailings_bytes_are_all_zero_fails_when_blob_smaller_than_payload() {
        // We are claiming a payload length of 33 bytes
        // The header alone occupies 32 bytes
        // So a [u8; 64] blob can't possibly fit header + payload
        let payload_len = 33u32;
        let payload_len_bytes = payload_len.to_be_bytes();
        let mut blob = [0u8; 64];
        blob[2..6].copy_from_slice(&payload_len_bytes);
        assert_eq!(
            verify_trailings_bytes_are_all_zero(&blob, payload_len),
            Err(BlobTooSmallForHeaderAndPayload(33))
        );
    }

    #[test]
    fn verify_trailings_bytes_are_all_zero_fails_when_trailing_bytes_not_all_zero() {
        let payload_len = 32u32;
        let payload_len_bytes = payload_len.to_be_bytes();
        let mut blob = [1u8; 65]; // 1 trailing byte not equal to 1
        blob[2..6].copy_from_slice(&payload_len_bytes);
        assert_eq!(
            verify_trailings_bytes_are_all_zero(&blob, payload_len),
            Err(NonZeroTrailingBytes)
        );
    }

    #[test]
    fn verify_trailings_bytes_are_all_zero_succeeds() {
        let payload_len = 32u32;
        let payload_len_bytes = payload_len.to_be_bytes();
        let mut blob = [0u8; 65]; // 1 trailing byte equal to 0
        blob[2..6].copy_from_slice(&payload_len_bytes);
        assert_eq!(
            verify_trailings_bytes_are_all_zero(&blob, payload_len),
            Ok(())
        );
    }

    #[test]
    fn verify_kzg_commitment_fails_when_blob_is_too_big_for_srs() {
        const LEN: usize = (POINTS_TO_LOAD as usize + 1) * 32;
        let mut blob = [42u8; LEN];
        let payload_len = (LEN as u32 - 32u32).to_be_bytes();
        blob[2..6].copy_from_slice(&payload_len);
        let claimed_commitment: G1Point = G1Affine::default().into_ext();
        assert_eq!(
            verify_kzg_commitment(&blob, claimed_commitment),
            Err(WrapKzgError(KzgError::SerializationError(
                "polynomial length is not correct".into()
            )))
        );
    }

    #[test]
    fn verify_kzg_commitment_succeeds() {
        let mut blob = [42u8; 42];
        let payload_len = 10u32.to_be_bytes();
        blob[2..6].copy_from_slice(&payload_len);

        // reproduce the calculation, more flexible than hardcoding the commitment
        let polynomial = Blob::new(&blob).to_polynomial_coeff_form();
        let claimed_commitment = KZG::new().commit_coeff_form(&polynomial, &SRS).unwrap();
        let claimed_commitment: G1Point = claimed_commitment.into_ext();

        assert_eq!(verify_kzg_commitment(&blob, claimed_commitment), Ok(()));
    }
}
