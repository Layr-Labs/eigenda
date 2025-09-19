//! Conversion utilities between different cryptographic representations
//!
//! This module provides functions for converting between EigenDA's G1Point
//! representation and arkworks' G1Affine, as well as utilities for
//! deterministic hash-to-curve operations.

use crate::eigenda::{cert::G1Point, verification::cert::hash::streaming_keccak256};
use alloy_primitives::{B256, Uint};
use ark_bn254::{Fq, G1Affine};
use ark_ff::{BigInt, BigInteger, Field, MontFp, PrimeField};

/// Field element one in Montgomery form
const ONE: Fq = MontFp!("1");

/// Field element three in Montgomery form (used in BN254 curve equation y² = x³ + 3)
const THREE: Fq = MontFp!("3");

/// Convert a G1 point to its hash representation.
///
/// Computes keccak256(x_bytes || y_bytes) where coordinates are encoded
/// as big-endian 32-byte arrays. This matches EigenDA's operator ID
/// generation from public keys.
///
/// # Arguments
/// * `point` - G1 point to hash
///
/// # Returns
/// 32-byte hash that uniquely identifies the point
pub fn point_to_hash(point: &G1Point) -> B256 {
    let x_bytes: [u8; 32] = point.x.to_be_bytes();
    let y_bytes: [u8; 32] = point.y.to_be_bytes();
    streaming_keccak256(&[&x_bytes, &y_bytes])
}

/// Convert a hash to a deterministic point on the BN254 curve.
///
/// Uses a simple try-and-increment method: treats the hash as an x-coordinate
/// and checks if it yields a valid point. If not, increments x and tries again.
/// This is deterministic and will always find a valid point.
///
/// # Arguments  
/// * `hash` - 32-byte hash to convert to a curve point
///
/// # Returns
/// A valid G1 point derived deterministically from the hash
pub(crate) fn hash_to_point(hash: B256) -> G1Affine {
    let x = hash_to_big_int(hash);
    let mut x = Fq::new(x);
    // safety: won't overflow the stack because:
    // - exactly half of non-zero field elements satisfy y^2 = x^3 + 3
    // - it's a finite field
    // So if x does not satisfy the equation, x + 1 probably will
    // In practice it'll take at most a few trials to succeed (90% of the time less than 3 tries are required)
    // Thus it is deterministic
    loop {
        let y = (x * x * x + THREE).sqrt();
        if let Some(y) = y {
            // `new_unchecked`: we've manually validated that (x, y) belongs to the curve
            return G1Affine::new_unchecked(x, y);
        }
        x += ONE;
    }
}

/// Convert a 32-byte B256 to arkworks BigInt representation.
///
/// Converts from big-endian byte representation to the little-endian
/// limb format expected by arkworks.
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

/// Convert field element to big-endian byte representation.
///
/// # Arguments
/// * `fq` - Field element to convert
///
/// # Returns
/// 32-byte big-endian representation
#[inline]
pub(crate) fn fq_to_bytes_be(fq: Fq) -> [u8; 32] {
    // safety: Fq is 256 bits
    fq.into_bigint().to_bytes_be().try_into().unwrap()
}

/// Convert field element to Uint representation.
///
/// # Arguments
/// * `fq` - Field element to convert
///
/// # Returns
/// 256-bit unsigned integer with 4 limbs
#[inline]
pub(crate) fn fq_to_uint(fq: Fq) -> Uint<256, 4> {
    Uint::from_limbs(fq.into_bigint().0)
}

#[cfg(test)]
mod tests {
    use alloy_primitives::{Uint, hex};
    use ark_bn254::{Fq, G1Affine};
    use ark_ec::AffineRepr;

    use crate::eigenda::{
        cert::G1Point,
        verification::cert::convert::{
            self, fq_to_bytes_be, fq_to_uint, hash_to_big_int, hash_to_point,
        },
    };

    #[test]
    fn convert_point_to_hash() {
        let point = G1Affine::generator();
        let actual = convert::point_to_hash(&point.into());
        let actual = hex::encode(actual);
        let expected = "e90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0";

        assert_eq!(actual, expected);
    }

    #[test]
    fn convert_infinity_to_hash() {
        let point = G1Affine::identity();
        let actual = convert::point_to_hash(&point.into());
        let point = G1Point {
            x: Uint::from_be_bytes([0u8; 32]),
            y: Uint::from_be_bytes([0u8; 32]),
        };
        let expected = convert::point_to_hash(&point);
        assert_eq!(actual, expected);
    }

    #[test]
    fn hash_to_point_test() {
        let hash = hex!("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef");
        let point = hash_to_point(hash.into());
        assert!(point.is_on_curve());
        assert!(!point.is_zero());
    }

    #[test]
    fn hash_to_big_int_test() {
        let hash = hex!("0000000000000000000000000000000000000000000000000000000000000001");
        let actual = hash_to_big_int(hash.into()).0;
        let expected = [1, 0, 0, 0];
        assert_eq!(actual, expected);
    }

    #[test]
    fn fq_to_bytes_be_test() {
        let fq = Fq::from(42u64);
        let actual = fq_to_bytes_be(fq);
        let expected = hex!("000000000000000000000000000000000000000000000000000000000000002a");
        assert_eq!(actual, expected);
    }

    #[test]
    fn fq_to_uint_test() {
        let fq = Fq::from(123u64);
        let actual = fq_to_uint(fq);
        let expected = Uint::from(123u64);
        assert_eq!(actual, expected);
    }
}
