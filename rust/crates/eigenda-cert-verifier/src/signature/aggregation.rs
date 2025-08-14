use alloc::vec::Vec;
use ark_bn254::{Fr, G1Affine, G1Projective};

use crate::{
    bitmap::bit_indices_to_bitmap,
    error::CertVerificationError,
    types::{NonSigner, Quorum},
};

// Example:
//
// signed_quorums: [0, 2] translate to this bitmap:
//         +-----+-----+-----+
// index:  |  2  |  1  |  0  |
//         +-----+-----+-----+
// bitmap: |  1  |  0  |  1  |
//         +-----+-----+-----+
//
// Quorum 1 being 0 means no signers that were required to sign actually did
// Quorums 0 and 2 being 1 means at least one signer that was required to sign
// actually did
//
// Let's assume there exist 6 signers, the first 3 being non-signers
// For each non-signer a quorum membership bitmap says whether they
// were required to sign at each quorum (1) or not (0)
//
// Signer 0 was required to sign at quorums 0 and 2 (but assume signed neither)
//         +-----+-----+-----+
// index:  |  2  |  1  |  0  |
//         +-----+-----+-----+
// bitmap: |  1  |  0  |  1  |
//         +-----+-----+-----+
//
// Signer 1 was required to sign at quorums 1 and 2 (but assume signed neither)
//         +-----+-----+-----+
// index:  |  2  |  1  |  0  |
//         +-----+-----+-----+
// bitmap: |  1  |  1  |  0  |
//         +-----+-----+-----+
//
// Signer 2 was required to sign at all quorums (but assume signed none)
//         +-----+-----+-----+
// index:  |  2  |  1  |  0  |
//         +-----+-----+-----+
// bitmap: |  1  |  1  |  1  |
//         +-----+-----+-----+
//
// Signer 3 was required to sign at quorum 2 (assume it did sign it)
//         +-----+-----+-----+
// index:  |  2  |  1  |  0  |
//         +-----+-----+-----+
// bitmap: |  1  |  0  |  0  |
//         +-----+-----+-----+
//
// Signer 4 was required to sign at quorum 0 (assume it did sign it)
//         +-----+-----+-----+
// index:  |  2  |  1  |  0  |
//         +-----+-----+-----+
// bitmap: |  0  |  0  |  1  |
//         +-----+-----+-----+
//
// Signer 5 was not required to sign at any quorum (assume it did not sign any since it was not required)
//         +-----+-----+-----+
// index:  |  2  |  1  |  0  |
//         +-----+-----+-----+
// bitmap: |  0  |  0  |  0  |
//         +-----+-----+-----+
//
// The above example quorum membership bitmaps specify only whether each signer was
// required to sign, they say nothing about whether they actually did or did not sign.
// So every statement above about them signing or not is for the sake of example.
// Following the example then non-signers have pubkeys [PK0, PK1, PK2]
// while signers have pubkeys [PK3, PK4]. PK5 belongs to neither set
//
// Since the signature is over the batch root from a tree of all blob certificates
// it means that a signer either signs all quorums it was assigned to
// (because the batch root represents all) or signs none at all,
// that is, they cannot sign some quorums but not others. This is important for the
// correctness of this implementation
//
// At its core the calculation iterates over non-signer quorum membership bitmaps
// ANDing each against `signed_quorums` to get as result the number of `required_non_signers`
// In other words, given `non_signers` = `required_non_signers` + `optional_non_signers`,
// the calculation filters out the optional_non_signers leaving only required_non_signers.
//
//            signer membership    &     signed quorums     =    required_signers
// Quorum:      2     1     0             2     1     0           2     1     0
//           +-----+-----+-----+       +-----+-----+-----+     +-----+-----+-----+
// Signer 0: |  1  |  0  |  1  |   &   |  1  |  0  |  1  |  =  |  1  |  0  |  1  |
//           +-----+-----+-----+       +-----+-----+-----+     +-----+-----+-----+
//
//           +-----+-----+-----+       +-----+-----+-----+     +-----+-----+-----+
// Signer 1: |  1  |  1  |  0  |   &   |  1  |  0  |  1  |  =  |  1  |  0  |  0  |
//           +-----+-----+-----+       +-----+-----+-----+     +-----+-----+-----+
//
//           +-----+-----+-----+       +-----+-----+-----+     +-----+-----+-----+
// Signer 2: |  1  |  1  |  1  |   &   |  1  |  0  |  1  |  =  |  1  |  0  |  1  |
//           +-----+-----+-----+       +-----+-----+-----+     +-----+-----+-----+
//
//           +-----+-----+-----+       +-----+-----+-----+     +-----+-----+-----+
// Signer 3: |  1  |  0  |  0  |   &   |  1  |  0  |  1  |  =  |  1  |  0  |  0  |
//           +-----+-----+-----+       +-----+-----+-----+     +-----+-----+-----+
//
//           +-----+-----+-----+       +-----+-----+-----+     +-----+-----+-----+
// Signer 4: |  0  |  0  |  1  |   &   |  1  |  0  |  1  |  =  |  0  |  0  |  1  |
//           +-----+-----+-----+       +-----+-----+-----+     +-----+-----+-----+
//
//           +-----+-----+-----+       +-----+-----+-----+     +-----+-----+-----+
// Signer 5: |  0  |  0  |  0  |   &   |  1  |  0  |  1  |  =  |  0  |  0  |  0  |
//           +-----+-----+-----+       +-----+-----+-----+     +-----+-----+-----+
//
// In the above, Signers 3, 4 and 5 are not `non_signers` so they have never been
// included in the calculation.
//
// Signers 0, 1 and 2 are non-signers and the resulting `missing_signatures` bitmap
// encode how many signatures were expected but not provided:
//
// 2 signatures were expected from Signer 0 (at quorums 0 and 2) but neither was provided
// 1 signature was expected from Signer 1 (at quorum 2) but it was not provided
// 2 signatures were expected from Signer 2 (at quorums 0 and 2) but neither was provided
//
// Note that quorum 1 was not signed at all by any signer so it's excluded from
// further consideration, that is, in what follows the expected (but not provided)
// signatures of Signers 1 and 2 for Quorum 1 will be simply ignored
// This is also why the example `signed_quorums` is [0, 2] instead of [0, 1, 2].
//
// All of the above considerations translate to an initial aggregate pubkey of:
// APK = -(2*PK0 + 1*PK1 + 2*PK2)
//
// That is, the calculation starts by subtracting pubkeys of non-signers proportional
// to how many signatures are missing from each
//
// Each quorum has an associated aggregate pubkey that corresponds to the sum of the
// pubkeys that were required to sign:
//
// For Quorum 0, Signers 0, 2 and 4 were required to sign
// For Quorum 1, Signers 1 and 2 were required to sign
// For Quorum 2, Signers 0, 1 and 2 were required to sign
//
// So the aggregate pubkeys of each quorum are:
//
// APK of Quorum 0: PK0 + PK2 + PK4
// APK of Quorum 1: PK1 + PK2 (which is ignored because there were no signers)
// APK of Quorum 2: PK0 + PK1 + PK2 + PK3
//
// The resulting aggregate pubkey is the sum of all quorums' aggregate pubkeys and
// the negated aggregate pubkey calculated earlier:
//
//       -    non-signers APK     +   Quorum 0 APK     + Quorum 1 APK +      Quorum 2 APK
// APK = -(2*PK0 + 1*PK1 + 2*PK2) + (PK0 + PK2 + PK4)  +   IDENTITY   + (PK0 + PK1 + PK2 + PK3)
//
// After cancelling out terms, the resulting `signers` APK is PK3 + PK4 as expected
// since those were the only signers that were both expected to sign and did sign
pub fn aggregate(
    quorum_count: u8,
    non_signers: &[NonSigner],
    quorums: &[Quorum],
) -> Result<G1Affine, CertVerificationError> {
    let total_apk = quorums
        .iter()
        .map(|quorum| quorum.apk)
        .sum::<G1Projective>();

    let bit_indices = quorums
        .iter()
        .map(|quorum| quorum.number)
        .collect::<Vec<_>>();

    let signed_quorums = bit_indices_to_bitmap(&bit_indices.into(), Some(quorum_count))?;

    let non_signers_apk = non_signers
        .iter()
        .map(|non_signer| {
            let missing_signatures = non_signer.quorum_bitmap_history & signed_quorums;
            let missing_signatures = missing_signatures.count_ones();
            let missing_signatures = Fr::from(missing_signatures as u64);
            // TODO: benchmark then consider implementing `scalar_mul_tiny`
            non_signer.pk * missing_signatures
        })
        .sum::<G1Projective>();

    let signers_apk = total_apk - non_signers_apk;

    Ok(signers_apk.into())
}

#[cfg(test)]
mod tests {
    use alloc::{vec, vec::Vec};
    use ark_bn254::{G1Affine, G1Projective};
    use ark_ec::{AffineRepr, CurveGroup, PrimeGroup};
    use ark_ff::BigInteger256;
    use bitvec::array::BitArray;

    use crate::{
        convert,
        error::CertVerificationError::*,
        signature::aggregation::aggregate,
        types::{NonSigner, Quorum, conversions::IntoExt},
    };

    #[test]
    fn compute_signers_apk_fails_with_too_many_quorums() {
        let quorums = vec![Default::default(); 256 + 1];
        let err = aggregate(Default::default(), Default::default(), &quorums).unwrap_err();
        assert_eq!(err, BitIndicesGreaterThanMaxLength);
    }

    #[test]
    fn compute_signers_apk_for_3_quorums_and_6_signers() {
        let (quorum_count, non_signers, quorums) = inputs_for_3_quorums_and_6_signers();

        let actual = aggregate(quorum_count, &non_signers, &quorums).unwrap();

        let expected = (ppk(3) + ppk(4)).into_affine();

        assert_eq!(actual, expected);
    }

    fn inputs_for_3_quorums_and_6_signers() -> (u8, Vec<NonSigner>, Vec<Quorum>) {
        let signed_quorums = vec![0, 2];
        let quorum_count = u8::MAX;

        let non_signer_pks = vec![pk(0), pk(1), pk(2)];

        let non_signer_quorum_bitmap_history = vec![
            BitArray::new([5, 0, 0, 0]), // 1 0 1
            BitArray::new([6, 0, 0, 0]), // 1 1 0
            BitArray::new([7, 0, 0, 0]), // 1 1 1
                                         // BitArray::new([4, 0, 0, 0]), // 1 0 0
                                         // BitArray::new([1, 0, 0, 0]), // 0 0 1
                                         // BitArray::new([0, 0, 0, 0]), // 0 0 0
        ];

        let non_signers = non_signer_pks
            .into_iter()
            .zip(non_signer_quorum_bitmap_history.into_iter())
            .map(|(pk, quorum_bitmap_history)| NonSigner {
                pk,
                pk_hash: convert::point_to_hash(&pk.into_ext()),
                quorum_bitmap_history,
            })
            .collect();

        let apks = vec![
            ppk(0) + ppk(2) + ppk(4),          // Quorum 0
            ppk(0) + ppk(1) + ppk(2) + ppk(3), // Quorum 2
        ];

        let quorums = signed_quorums
            .iter()
            .zip(apks.into_iter())
            .map(|(signed_quorum_number, apk)| Quorum {
                number: *signed_quorum_number,
                apk: apk.into_affine(),
                ..Default::default()
            })
            .collect();

        (quorum_count, non_signers, quorums)
    }

    fn pk(n: u64) -> G1Affine {
        ppk(n).into_affine()
    }

    fn ppk(n: u64) -> G1Projective {
        let generator = G1Affine::generator();
        generator
            .into_group()
            .mul_bigint(BigInteger256::from(n + 1))
    }
}
