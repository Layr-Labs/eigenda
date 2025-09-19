#[path = "src/eigenda/verification/blob/srs.rs"]
mod srs;

use std::env;
use std::path::Path;

use ark_serialize::CanonicalSerialize;
use rust_kzg_bn254_prover::srs::SRS;
use srs::SerializableSRS;

use crate::srs::POINTS_TO_LOAD;

fn main() {
    println!("cargo:rerun-if-changed=resources/g1.point");

    let path = "resources/g1.point";

    if !Path::new(path).exists() {
        panic!("g1.point file not found at {path}");
    }

    let srs = SRS::new(path, 268435456, POINTS_TO_LOAD).expect("Failed to create SRS");

    let wrapper: SerializableSRS = srs.into();
    let mut serialized = Vec::new();
    wrapper
        .serialize_compressed(&mut serialized)
        .expect("Failed to serialize SRS");

    let out_dir = env::var("OUT_DIR").unwrap();
    let path = Path::new(&out_dir).join("srs.bin");
    std::fs::write(&path, &serialized).expect("Failed to write serialized SRS");
}
