//! BLS signature verification using bilinear pairings
//!
//! This module implements BLS signature verification for EigenDA certificates using
//! the BN254 pairing-friendly elliptic curve. It verifies that aggregate signatures
//! were indeed created by the claimed aggregate public keys.
//!
//! ## BLS Signature Verification
//!
//! The verification process uses bilinear pairings to check the equation:
//! `e(σ + γG₁, -G₂) · e(H(m) + γG₁, APK_G₂) = 1`
//!
//! Where:
//! - `σ`: The aggregate signature (G₁ point)
//! - `APK_G₂`: Aggregate public key on G₂
//! - `H(m)`: Hash-to-curve of the message  
//! - `γ`: Challenge derived from Fiat-Shamir heuristic
//! - `G₁`, `G₂`: Curve generators
//!
//! ## Security
//!
//! The challenge `γ` is computed as `keccak256(H(m) || APK_G₁ || APK_G₂ || σ)`
//! to prevent rogue public key attacks in the aggregate setting.

use std::sync::LazyLock;

use alloy_primitives::B256;
use ark_bn254::{Bn254, Fr, G1Affine, G2Affine};
use ark_ec::bn::G2Prepared;
use ark_ec::pairing::{Pairing, PairingOutput};
use ark_ec::{AffineRepr, CurveGroup};
use ark_ff::{AdditiveGroup, PrimeField};

use crate::verification::cert::convert;
use crate::verification::cert::hash::streaming_keccak256;

static PRECOMPUTED_NEG_G2: LazyLock<G2Prepared<ark_bn254::Config>> =
    LazyLock::new(|| G2Prepared::from(-G2Affine::generator()));

/// Verify a BLS signature using bilinear pairings.
///
/// Checks if the signature `sigma` was created by holders of the aggregate public key
/// (`apk_g1`, `apk_g2`) over the message `msg_hash`.
///
/// # Arguments
/// * `msg_hash` - 32-byte hash of the message that was signed
/// * `apk_g1` - Aggregate public key on G₁ (used for challenge computation)
/// * `apk_g2` - Aggregate public key on G₂ (used in pairing verification)  
/// * `sigma` - Aggregate signature to verify (G₁ point)
///
/// # Returns
/// `true` if the signature is valid, `false` otherwise
///
/// # Algorithm
/// Verifies the equation: `e(σ + γG₁, -G₂) · e(H(m) + γG₁, APK_G₂) = 1`
///
/// Where `γ = keccak256(msg_hash || apk_g1 || apk_g2 || sigma)` is a Fiat-Shamir challenge
/// that prevents rogue public key attacks in the aggregate signature setting.
pub fn verify(msg_hash: B256, apk_g1: G1Affine, apk_g2: G2Affine, sigma: G1Affine) -> bool {
    let Some(gamma) = compute_gamma(msg_hash, apk_g1, apk_g2, sigma) else {
        return false;
    };
    let msg_point = convert::hash_to_point(msg_hash);

    let a1 = (sigma + apk_g1 * gamma).into_affine();
    let a2 = PRECOMPUTED_NEG_G2.clone();
    let b1 = (msg_point + G1Affine::generator() * gamma).into_affine();
    let b2 = G2Prepared::from(apk_g2);

    let miller_result = Bn254::multi_miller_loop([a1, b1], [a2, b2]);
    let pairing_result = Bn254::final_exponentiation(miller_result);
    // `pairing_result` could be None if one of `a1`, `b1`, `a2`, `b2` is at infinity
    // a PairingOutput::zero() has an underlying TargetField::one()
    // which is the RHS of e(sigma + apk_g1 * gamma, -G2) * e(msg_hash + G1 * gamma, apk_g2) == 1
    pairing_result == Some(PairingOutput::ZERO)
}

/// Compute the Fiat-Shamir challenge for BLS signature verification.
///
/// Creates a cryptographic challenge by hashing all public parameters
///
/// # Arguments
/// * `msg_hash` - Hash of the signed message
/// * `apk_g1` - Aggregate public key on G₁
/// * `apk_g2` - Aggregate public key on G₂  
/// * `sigma` - Signature being verified
///
/// # Returns
/// * `Some(Fr)` - Challenge scalar if all points are valid (not at infinity)
/// * `None` - If any input point is at infinity (invalid)
fn compute_gamma(
    msg_hash: B256,
    apk_g1: G1Affine,
    apk_g2: G2Affine,
    sigma: G1Affine,
) -> Option<Fr> {
    // returns None if any point is at infinity
    let (apk_g1_x, apk_g1_y) = apk_g1.xy()?;
    let (apk_g2_x, apk_g2_y) = apk_g2.xy()?;
    let (sigma_x, sigma_y) = sigma.xy()?;

    let gamma = streaming_keccak256(&[
        msg_hash.as_slice(),
        &convert::fq_to_bytes_be(apk_g1_x),
        &convert::fq_to_bytes_be(apk_g1_y),
        &convert::fq_to_bytes_be(apk_g2_x.c0),
        &convert::fq_to_bytes_be(apk_g2_x.c1),
        &convert::fq_to_bytes_be(apk_g2_y.c0),
        &convert::fq_to_bytes_be(apk_g2_y.c1),
        &convert::fq_to_bytes_be(sigma_x),
        &convert::fq_to_bytes_be(sigma_y),
    ]);

    let gamma = Fr::from_be_bytes_mod_order(&*gamma);
    Some(gamma)
}

#[cfg(test)]
mod tests {
    use ark_bn254::{Fr, G1Affine, G1Projective, G2Affine, G2Projective};
    use ark_ec::{AffineRepr, CurveGroup, PrimeGroup};

    use crate::verification::cert::convert;
    use crate::verification::cert::signature::verification::{compute_gamma, verify};

    #[test]
    fn signature_roundtrip() {
        let sk = Fr::from(42);
        let apk_g1 = (G1Projective::generator() * sk).into_affine();
        let apk_g2 = (G2Projective::generator() * sk).into_affine();
        let msg_hash = [42u8; 32].into();
        let msg_point = convert::hash_to_point(msg_hash);
        let sigma = (msg_point * sk).into_affine();
        let result = verify(msg_hash, apk_g1, apk_g2, sigma);
        assert!(result);
    }

    #[test]
    fn signature_not_signed_by_expected_signer() {
        let expected_signer_sk = Fr::from(42);
        let apk_g1 = (G1Projective::generator() * expected_signer_sk).into_affine();
        let apk_g2 = (G2Projective::generator() * expected_signer_sk).into_affine();
        let msg_hash = [42u8; 32].into();
        let msg_point = convert::hash_to_point(msg_hash);

        let actual_signer_sk = Fr::from(43);
        let sigma = (msg_point * actual_signer_sk).into_affine();
        let result = verify(msg_hash, apk_g1, apk_g2, sigma);
        assert!(!result);
    }

    #[test]
    fn inputs_at_infinity() {
        let msg_hash = [42u8; 32].into();

        let sk = Fr::from(42);
        let apk_g1 = (G1Projective::generator() * sk).into_affine();
        let apk_g2 = (G2Projective::generator() * sk).into_affine();
        let sigma = G1Affine::generator();

        let result = verify(msg_hash, G1Affine::identity(), apk_g2, sigma);
        assert!(!result);

        let result = verify(msg_hash, apk_g1, G2Affine::identity(), sigma);
        assert!(!result);

        let result = verify(msg_hash, apk_g1, apk_g2, G1Affine::identity());
        assert!(!result);
    }

    #[test]
    fn compute_gamma_baseline() {
        use ark_ff::{BigInteger, PrimeField};

        let msg_hash = [42u8; 32].into();
        let sk = Fr::from(12345);
        let apk_g1 = (G1Projective::generator() * sk).into_affine();
        let apk_g2 = (G2Projective::generator() * sk).into_affine();
        let sigma = (G1Projective::generator() * Fr::from(67890)).into_affine();

        let gamma = compute_gamma(msg_hash, apk_g1, apk_g2, sigma).unwrap();
        let actual = hex::encode(gamma.into_bigint().to_bytes_be());
        let expected = "1866953a8361306ca9a0b59082525a8e917e686c9cf66fa00cb3bcf3ecae6164";

        assert_eq!(actual, expected);
    }
}
