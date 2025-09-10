use ark_bn254::G1Affine;
use ark_serialize::{CanonicalDeserialize, CanonicalSerialize};
use rust_kzg_bn254_prover::srs::SRS;

/// Maximum number of SRS points to load for KZG operations.
///
/// This determines the maximum polynomial degree that can be committed to,
/// effectively limiting the maximum blob size that can be verified.
/// With 4096 points, blobs up to ~127KB can be handled (assuming ~32 bytes per coefficient).
///
/// Used in build.rs to determine how many points to load from the trusted setup,
/// and in tests to validate the SRS configuration.
// used in build.rs and tests
#[allow(dead_code)]
pub const POINTS_TO_LOAD: u32 = 16 * 1024 * 1024 / 32; // 16 MiB worth of points

/// Serializable wrapper for the Structured Reference String (SRS).
///
/// The SRS contains the cryptographic parameters needed for KZG polynomial
/// commitments and proof verification. This wrapper allows the SRS to be
/// serialized/deserialized using arkworks' canonical format for storage
/// and transport.
///
/// KZG commitments require a trusted setup ceremony to generate these points,
/// which are then used for polynomial commitment and verification operations.
#[derive(CanonicalSerialize, CanonicalDeserialize)]
pub struct SerializableSRS {
    /// G1 curve points from the trusted setup, used for polynomial commitments
    /// The number of points determines the maximum polynomial degree supported
    pub g1: Vec<G1Affine>,
    /// Number of valid points in the SRS
    pub order: u32,
}

impl From<SRS> for SerializableSRS {
    fn from(srs: SRS) -> Self {
        Self {
            g1: srs.g1,
            order: srs.order,
        }
    }
}

impl From<SerializableSRS> for SRS {
    fn from(srs: SerializableSRS) -> Self {
        Self {
            g1: srs.g1,
            order: srs.order,
        }
    }
}
