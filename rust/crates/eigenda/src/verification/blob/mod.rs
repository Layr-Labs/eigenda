//! EigenDA blob verification using KZG polynomial commitments
//!
//! This module implements the blob validation stage of EigenDA verification,
//! ensuring that blob data matches its cryptographic commitment using KZG proofs
//! over the BN254 curve.
//!
//! ## Overview
//!
//! Blob verification validates that received data matches the commitment specified
//! in an EigenDA certificate. This prevents data tampering and ensures integrity
//! of the data availability guarantees.
//!
//! ## Verification Process
//!
//! The verification follows the [EigenDA specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation):
//!
//! 1. **Length Validation**: Ensure received blob length ≤ committed length
//! 2. **Power-of-two Check**: Verify commitment length is a power of two
//! 3. **Payload Encoding**: Transform payload into proper blob format
//! 4. **Header Validation**: Verify encoded payload header constraints
//! 5. **Padding Verification**: Ensure all extra bytes are zero
//! 6. **KZG Commitment**: Verify the cryptographic commitment matches
//!
//! ## Blob Encoding Format
//!
//! EigenDA uses a specific encoding format for blobs:
//!
//! ```text
//! [32-byte header][padded payload symbols...]
//!
//! Header format:
//! - Byte 0: Field element guard (0x00)
//! - Byte 1: Version (0x00)  
//! - Bytes 2-5: Payload length (big-endian u32)
//! - Bytes 6-31: Zero padding
//!
//! Payload symbols:
//! - Each 31-byte payload chunk becomes a 32-byte symbol
//! - Symbols are prefixed with field element guard byte (0x00)
//! - Final chunk padded with zeros if needed
//! ```
//!
//! ## KZG Verification
//!
//! The module uses KZG polynomial commitments over BN254 for cryptographic verification:
//! - Recomputes the commitment from blob data using SRS points
//! - Compares computed commitment with claimed commitment
//! - Uses precomputed SRS (Structured Reference String)

pub mod codec;
pub mod error;
/// Structured Reference String (SRS) for KZG commitments.
///
/// This module manages the trusted setup parameters required for KZG polynomial
/// commitments used in EigenDA blob verification.
pub mod srs;

use std::sync::LazyLock;

use ark_bn254::G1Affine;
use ark_serialize::CanonicalDeserialize;
use rust_kzg_bn254_primitives::blob::Blob;
use rust_kzg_bn254_prover::kzg::KZG;
use rust_kzg_bn254_prover::srs::SRS;

use crate::cert::{BlobCommitment, G1Point};
use crate::verification::blob::codec::BYTES_PER_SYMBOL;
use crate::verification::blob::error::BlobVerificationError;
use crate::verification::blob::srs::SerializableSRS;

const SRS_BYTES: &[u8] = include_bytes!(concat!(env!("OUT_DIR"), "/srs.bin"));

/// Global SRS instance loaded lazily for KZG operations
///
/// The SRS is precomputed at build time and embedded in the binary for
/// efficient KZG commitment verification without runtime generation.
pub static SRS: LazyLock<SRS> = LazyLock::new(|| {
    SerializableSRS::deserialize_compressed(SRS_BYTES)
        .expect("Failed to deserialize precomputed SRS")
        .into()
});

/// Verifies that `blob` passes all the checks defined in
/// [EigenDA specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation)!
/// Verify blob data against its KZG commitment
///
/// Performs comprehensive validation of blob data according to the EigenDA
/// specification, including length checks, encoding validation, and KZG
/// commitment verification.
///
/// # Arguments
/// * `blob_commitment` - The commitment from the EigenDA certificate
/// * `payload` - Raw blob data to verify
///
/// # Returns
/// `Ok(())` if the blob is valid and matches the commitment
///
/// # Errors
/// Returns [`BlobVerificationError`] for various validation failures:
/// - Blob larger than committed length
/// - Invalid commitment length (not power of two)
/// - Payload too large for encoding
/// - KZG commitment mismatch
///
/// # Reference
/// [EigenDA Specification - Blob Validation](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation)
pub fn verify(
    blob_commitment: &BlobCommitment,
    encoded_payload: &[u8],
) -> Result<(), BlobVerificationError> {
    let blob = Blob::new(encoded_payload);
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
///
/// 6. Verify the KZG commitment. This can either be done:
///
///   1. directly: recomputing the commitment using SRS points and checking
///      that the two commitments match (this is the current implemented way)
///   2. indirectly: verifying a point opening using Fiat-Shamir (see this [issue](https://github.com/Layr-Labs/eigenda/issues/1037))
///
/// > the referenced PR is still open so we don't have the means to implement option 2.
fn verify_kzg_commitment(
    blob: &Blob,
    claimed_commitment: G1Point,
) -> Result<(), BlobVerificationError> {
    use BlobVerificationError::*;

    // for a large number of SRS points this is slow: ~40s in debug (~3s in release) on an M2 due to the 16MiB SRS one-time deserialization
    // that is first materialized here when the LazyLock is first accessed
    let computed_commitment = KZG::new().commit_blob(blob, &SRS)?;

    let claimed_commitment: G1Affine = claimed_commitment.into();

    (computed_commitment == claimed_commitment)
        .then_some(())
        .ok_or(InvalidKzgCommitment)
}

#[cfg(test)]
mod test {
    use crate::verification::blob::error::BlobVerificationError::*;
    use crate::verification::blob::{
        verify_blob_symbols_len_against_commitment, verify_commitment_len_is_power_of_two,
    };

    // This test takes ~40s in debug (~3s in release) on an M2 due to 16MiB SRS one-time deserialization
    // Using LazyLock is very advantageous for testing since many tests don't actually ever access
    // the expensive SRS resource which means it doesn't ever get deserialized in tests that don't
    // use it
    #[test]
    #[cfg(not(debug_assertions))]
    fn verify_succeeds_with_known_commitment() {
        use std::str::FromStr;

        use ark_bn254::{Fq, G1Affine, G2Affine};

        use crate::cert::BlobCommitment;
        use crate::verification::blob::codec::tests::encode_raw_payload;
        use crate::verification::blob::verify;

        let raw_payload = [123; 512];
        let encoded_payload = encode_raw_payload(&raw_payload).unwrap();

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
        .into();

        let blob_commitment = BlobCommitment {
            commitment: known_commitment,
            length_commitment: G2Affine::default().into(),
            length_proof: G2Affine::default().into(),
            length: 32 + 32,
        };

        assert_eq!(verify(&blob_commitment, &encoded_payload), Ok(()));
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
}
