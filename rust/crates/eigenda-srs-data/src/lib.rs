//! Generated SRS data for EigenDA blob verification
//!
//! This crate contains compile-time embedded Structured Reference String (SRS) data
//! as raw bytes that are transmuted into a G1Affine array.

use std::borrow::Cow;
use std::sync::LazyLock;

use ark_bn254::G1Affine;
use rust_kzg_bn254_prover::srs::SRS;

include!(concat!(env!("OUT_DIR"), "/constants.rs"));

// SAFETY: Transmuting compile-time embedded binary data to typed G1Affine array.
// - Binary data originates from the same G1Affine structures in build.rs
// - BYTE_SIZE constant ensures exact size match: POINTS_TO_LOAD * size_of::<G1Affine>()
// - G1Affine has stable, well-defined memory representation from ark-bn254
// - Both source and target arrays have identical size and alignment requirements
// - Static lifetime is appropriate for compile-time embedded data
static SRS_POINTS: &[G1Affine; POINTS_TO_LOAD] = unsafe {
    &core::mem::transmute::<[u8; BYTE_SIZE], [G1Affine; POINTS_TO_LOAD]>(*include_bytes!(concat!(
        env!("OUT_DIR"),
        "/srs_points.bin"
    )))
};

/// Globally accessible SRS (Structured Reference String) for KZG operations.
///
/// This static contains precomputed G1 curve points loaded from embedded binary data.
/// The SRS is lazily initialized on first access and provides the cryptographic
/// parameters needed for KZG polynomial commitments and proofs.
pub static SRS: LazyLock<SRS<'static>> = LazyLock::new(|| SRS {
    g1: Cow::Borrowed(SRS_POINTS),
    order: (POINTS_TO_LOAD * 32) as u32,
});
