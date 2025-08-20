use ark_bn254::G1Affine;
use ark_serialize::{CanonicalDeserialize, CanonicalSerialize};
use rust_kzg_bn254_prover::srs::SRS;

// used in build.rs and tests
#[allow(dead_code)]
pub const POINTS_TO_LOAD: u32 = 1024;

#[derive(CanonicalSerialize, CanonicalDeserialize)]
pub struct SerializableSRS {
    pub g1: Vec<G1Affine>,
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
