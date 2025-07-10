use ark_bn254::G1Affine;
use ark_ff::{BigInteger, PrimeField};
use tiny_keccak::{Hasher, Keccak};

pub fn hash_g1_point(point: &G1Affine) -> [u8; 32] {
    let x_bytes = point.x.into_bigint().to_bytes_be();
    let y_bytes = point.y.into_bigint().to_bytes_be();

    // pad
    let mut packed = [0u8; 64];
    packed[32 - x_bytes.len()..32].copy_from_slice(&x_bytes);
    packed[64 - y_bytes.len()..64].copy_from_slice(&y_bytes);

    let mut hasher = Keccak::v256();
    hasher.update(&packed[..]);

    let mut output = [0u8; 32];
    hasher.finalize(&mut output[..]);
    output
}

#[cfg(test)]
mod tests {
    use super::*;
    use ark_ec::AffineRepr;

    #[test]
    fn test_hash_g1_point() {
        let point = G1Affine::generator();
        let actual = hash_g1_point(&point);
        let actual = hex::encode(actual);
        let generated_from_js = "e90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0";

        assert_eq!(actual, generated_from_js);
    }
}
