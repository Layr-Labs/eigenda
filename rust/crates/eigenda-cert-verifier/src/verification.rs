use ark_bn254::{Bn254, Fr, G1Affine, G2Affine};
use ark_ec::{
    AffineRepr, CurveGroup,
    pairing::{Pairing, PairingOutput},
};
use ark_ff::{PrimeField, Zero};

use crate::{
    convert,
    hash::{self, BeHash},
};

/// Verifies `sigma` signature over `msg_hash` by the (`apk_g1`, `apk_g2`) pubkey
/// by checking e(sigma + apk_g1 * gamma, -G2) * e(msg_hash + G1 * gamma, apk_g2) == 1
pub fn verify(msg_hash: &BeHash, apk_g1: G1Affine, apk_g2: G2Affine, sigma: G1Affine) -> bool {
    let Some(gamma) = compute_gamma(msg_hash, &apk_g1, &apk_g2, &sigma) else {
        return false;
    };
    let msg_point = convert::hash_to_point(msg_hash);

    let a1 = (sigma + apk_g1 * gamma).into_affine();
    let a2 = -G2Affine::generator();
    let b1 = (msg_point + G1Affine::generator() * gamma).into_affine();
    let b2 = apk_g2;

    let g1_iter = [a1, b1].into_iter();
    let g2_iter = [a2, b2].into_iter();

    let miller_result = Bn254::multi_miller_loop(g1_iter, g2_iter);
    let pairing_result = Bn254::final_exponentiation(miller_result);
    // a PairingOutput::zero() has an underlying TargetField::one()
    // which is the RHS of e(sigma + apk_g1 * gamma, -G2) * e(msg_hash + G1 * gamma, apk_g2) == 1
    pairing_result == Some(PairingOutput::zero())
}

fn compute_gamma(
    msg_hash: &BeHash,
    apk_g1: &G1Affine,
    apk_g2: &G2Affine,
    sigma: &G1Affine,
) -> Option<Fr> {
    // returns None if any point is at infinity
    let (apk_g1_x, apk_g1_y) = apk_g1.xy()?;
    let (apk_g2_x, apk_g2_y) = apk_g2.xy()?;
    let (sigma_x, sigma_y) = sigma.xy()?;

    let gamma = hash::keccak256(&[
        msg_hash,
        &convert::fq_to_bytes_be(&apk_g1_x),
        &convert::fq_to_bytes_be(&apk_g1_y),
        &convert::fq_to_bytes_be(&apk_g2_x.c0),
        &convert::fq_to_bytes_be(&apk_g2_x.c1),
        &convert::fq_to_bytes_be(&apk_g2_y.c0),
        &convert::fq_to_bytes_be(&apk_g2_y.c1),
        &convert::fq_to_bytes_be(&sigma_x),
        &convert::fq_to_bytes_be(&sigma_y),
    ]);

    let gamma = Fr::from_be_bytes_mod_order(&gamma);
    Some(gamma)
}

#[cfg(test)]
mod tests {
    use ark_bn254::{Fr, G1Affine, G1Projective, G2Affine, G2Projective};
    use ark_ec::{AffineRepr, CurveGroup, PrimeGroup};
    use ark_ff::{BigInteger, PrimeField};

    use crate::convert;
    use crate::verification::verify;

    #[test]
    fn signature_roundtrip() {
        let sk = Fr::from(42u64);
        let apk_g1 = (G1Projective::generator() * sk).into_affine();
        let apk_g2 = (G2Projective::generator() * sk).into_affine();
        let msg_hash = [42u8; 32];
        let msg_point = convert::hash_to_point(&msg_hash);
        let sigma = (msg_point * sk).into_affine();
        let result = verify(&msg_hash, apk_g1, apk_g2, sigma);
        assert_eq!(result, true);
    }

    #[test]
    fn signature_not_signed_by_expected_signer() {
        let expected_signer_sk = Fr::from(42u64);
        let apk_g1 = (G1Projective::generator() * expected_signer_sk).into_affine();
        let apk_g2 = (G2Projective::generator() * expected_signer_sk).into_affine();
        let msg_hash = [42u8; 32];
        let msg_point = convert::hash_to_point(&msg_hash);

        let actual_signer_sk = Fr::from(43u64);
        let sigma = (msg_point * actual_signer_sk).into_affine();
        let result = verify(&msg_hash, apk_g1, apk_g2, sigma);
        assert_eq!(result, false);
    }

    #[test]
    fn inputs_at_infinity() {
        let msg_hash = [42u8; 32];

        let sk = Fr::from(42u64);
        let apk_g1 = (G1Projective::generator() * sk).into_affine();
        let apk_g2 = (G2Projective::generator() * sk).into_affine();
        let sigma = G1Affine::generator();

        let result = verify(&msg_hash, G1Affine::identity(), apk_g2, sigma);
        assert_eq!(result, false);

        let result = verify(&msg_hash, apk_g1, G2Affine::identity(), sigma);
        assert_eq!(result, false);

        let result = verify(&msg_hash, apk_g1, apk_g2, G1Affine::identity());
        assert_eq!(result, false);
    }

    #[test]
    fn behaviour_of_to_bytes_be() {
        let x = G1Affine::generator().x().unwrap().into_bigint();
        let expected = [
            x.0[3].to_be_bytes(),
            x.0[2].to_be_bytes(),
            x.0[1].to_be_bytes(),
            x.0[0].to_be_bytes(),
        ]
        .concat();
        let actual = x.to_bytes_be();
        assert_eq!(actual, expected);
    }
}
