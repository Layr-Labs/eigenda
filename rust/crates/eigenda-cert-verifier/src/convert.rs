use ark_bn254::{Fq, G1Affine};
use ark_ec::AffineRepr;
use ark_ff::{BigInt, Field, MontFp, PrimeField};

use crate::hash::{self, BeHash};

const ONE: Fq = MontFp!("1");
const THREE: Fq = MontFp!("3");

pub fn point_to_hash(point: &G1Affine) -> Option<BeHash> {
    point.xy().map(|(x, y)| {
        let x_bytes = fq_to_bytes_be(&x);
        let y_bytes = fq_to_bytes_be(&y);
        hash::keccak256(&[&x_bytes, &y_bytes])
    })
}

pub fn hash_to_point<'a>(hash: &BeHash) -> G1Affine {
    let x = hash_to_big_int(hash);
    let mut x = Fq::from(x);
    // Won't overflow the stack because:
    // - exactly half of non-zero field elements satisfy y^2 = x^3 + 3
    // - it's a finite field
    // So if x does not satisfy the equation, x + 1 probably will
    // In practice it'll take at most a few trials to succeed (90% of the time less than 3 tries are required)
    // It is also deterministic
    loop {
        let y = (x * x * x + THREE).sqrt();
        if let Some(y) = y {
            // `new_unchecked`: we've manually validated that (x, y) belongs to the curve
            return G1Affine::new_unchecked(x, y);
        }
        x += ONE;
    }
}

pub fn hash_to_big_int(hash: &BeHash) -> BigInt<4> {
    let mut limbs = [0u64; 4];

    for (i, chunk) in hash.chunks_exact(8).enumerate() {
        // ark-ff expects little-endian limbs so we reverse limb order ([3-i])
        // safe to unwrap because `chunk` is guaranteed to be convertible to [u8; 8] given `hash` is [u8; 32]
        limbs[3 - i] = u64::from_be_bytes(chunk.try_into().unwrap());
    }

    BigInt::new(limbs)
}

pub fn fq_to_bytes_be(fq: &Fq) -> [u8; 32] {
    let fq = fq.into_bigint();
    let mut bytes_be = [0u8; 32];
    // .rev() to make it big-endian
    for (i, limb) in fq.0.iter().rev().enumerate() {
        bytes_be[i * 8..(i + 1) * 8].copy_from_slice(&limb.to_be_bytes());
    }
    bytes_be
}

#[cfg(test)]
mod tests {
    use ark_bn254::G1Affine;
    use ark_ec::AffineRepr;

    use crate::convert;

    #[test]
    fn convert_point_to_hash() {
        let point = G1Affine::generator();
        let actual = convert::point_to_hash(&point).unwrap();
        let actual = hex::encode(actual);
        let generated_from_js = "e90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0";

        assert_eq!(actual, generated_from_js);
    }
}
