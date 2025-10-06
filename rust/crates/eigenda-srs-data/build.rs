//! Build script for srs-data crate.
//!
//! This script generates compile-time Rust code for the SRS (Structured Reference String)
//! by reading the g1.point file and creating a static G1Affine point array that can be
//! embedded directly in the binary at compile time
use std::path::Path;
use std::{env, fs, mem};

use ark_bn254::G1Affine;
use rust_kzg_bn254_prover::srs::SRS;

const POINTS_TO_LOAD: u32 = 16 * 1024 * 1024 / 32;

fn main() {
    println!("cargo:rerun-if-changed=resources/g1.point");

    let path = "resources/g1.point";

    let order = POINTS_TO_LOAD * 32;
    let srs = SRS::new(path, order, POINTS_TO_LOAD).expect("Failed to create SRS");
    assert_eq!(srs.g1.len(), POINTS_TO_LOAD as usize);

    let out_dir = env::var("OUT_DIR").unwrap();
    let out_path = Path::new(&out_dir);

    let g1_slice = &srs.g1[..];
    // SAFETY: Converting G1Affine slice to byte slice for serialization.
    // - g1_slice is a valid reference to G1Affine elements with known lifetime
    // - G1Affine has a well-defined memory layout from ark-bn254
    // - size_of_val() ensures the byte slice doesn't exceed source data bounds
    // - The resulting byte slice lifetime is bounded by the original slice
    let g1_bytes = unsafe {
        std::slice::from_raw_parts(g1_slice.as_ptr() as *const u8, size_of_val(g1_slice))
    };

    let g1_path = out_path.join("srs_points.bin");
    fs::write(&g1_path, g1_bytes).expect("Failed to write G1 points");

    let byte_size = POINTS_TO_LOAD as usize * mem::size_of::<G1Affine>();

    macro_rules! generate_constants {
        ($points:expr, $byte_size:expr) => {
            format!(
                r#"// Auto-generated constants - DO NOT EDIT

/// Number of G1 points to load from the SRS data.
/// This represents the maximum degree of polynomials that can be committed.
pub const POINTS_TO_LOAD: usize = {};

/// Total byte size of the embedded SRS point data.
/// This is calculated as POINTS_TO_LOAD * size_of::<G1Affine>().
pub const BYTE_SIZE: usize = {};
"#,
                $points, $byte_size
            )
        };
    }

    let constants_content = generate_constants!(POINTS_TO_LOAD, byte_size);
    let constants_path = out_path.join("constants.rs");
    fs::write(&constants_path, constants_content).expect("Failed to write constants");
}
