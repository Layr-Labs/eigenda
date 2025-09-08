use crate::eigenda::{cert::G1Point, verification::cert::hash::streaming_keccak256};
use alloy_primitives::{B256, Uint};
use ark_bn254::{Fq, G1Affine};
use ark_ff::{BigInt, BigInteger, Field, MontFp, PrimeField};

const ONE: Fq = MontFp!("1");
const THREE: Fq = MontFp!("3");

pub fn point_to_hash(point: &G1Point) -> B256 {
    let x_bytes: [u8; 32] = point.x.to_be_bytes();
    let y_bytes: [u8; 32] = point.y.to_be_bytes();
    streaming_keccak256(&[&x_bytes, &y_bytes])
}

pub(crate) fn hash_to_point(hash: B256) -> G1Affine {
    let x = hash_to_big_int(hash);
    let mut x = Fq::new(x);
    // safety: won't overflow the stack because:
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

#[inline]
fn hash_to_big_int(hash: B256) -> BigInt<4> {
    let mut limbs = [0u64; 4];

    for (i, chunk) in hash.chunks_exact(8).enumerate() {
        // ark-ff expects little-endian limbs so we reverse limb order ([3-i])
        // safe to unwrap because `chunk` is guaranteed to be convertible to [u8; 8] given `hash` is [u8; 32]
        limbs[3 - i] = u64::from_be_bytes(chunk.try_into().unwrap());
    }

    BigInt::new(limbs)
}

#[inline]
pub(crate) fn fq_to_bytes_be(fq: Fq) -> [u8; 32] {
    // safety: Fq is 256 bits
    fq.into_bigint().to_bytes_be().try_into().unwrap()
}

#[inline]
pub(crate) fn fq_to_uint(fq: Fq) -> Uint<256, 4> {
    Uint::from_limbs(fq.into_bigint().0)
}

#[cfg(test)]
mod tests {
    use crate::eigenda::{
        cert::G1Point,
        verification::cert::{convert, types::conversions::IntoExt},
    };
    use alloy_primitives::Uint;
    use ark_bn254::G1Affine;
    use ark_ec::AffineRepr;

    #[test]
    fn convert_point_to_hash() {
        let point = G1Affine::generator();
        let actual = convert::point_to_hash(&point.into_ext());
        let actual = hex::encode(actual);
        let expected = "e90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0";

        assert_eq!(actual, expected);
    }

    #[test]
    fn convert_infinity_to_hash() {
        let point = G1Affine::identity();
        let actual = convert::point_to_hash(&point.into_ext());
        let point = G1Point {
            x: Uint::from_be_bytes([0u8; 32]),
            y: Uint::from_be_bytes([0u8; 32]),
        };
        let expected = convert::point_to_hash(&point);
        assert_eq!(actual, expected);
    }
}
