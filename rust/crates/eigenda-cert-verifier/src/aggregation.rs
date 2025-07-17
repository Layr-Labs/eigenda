use alloc::vec::Vec;

use ark_bn254::{G1Affine, G1Projective};
use ark_ec::{AffineRepr, CurveGroup, PrimeGroup};
use ark_ff::BigInteger256;
use hashbrown::HashMap;

use crate::{
    bitmap_utils::{Bitmap, bit_indices_to_bitmap},
    convert::point_to_hash,
    error::SignaturesVerificationError,
    hash::BeHash,
};

// Example:
//
// quorum_numbers: [0, 2] translate to this bitmap:
//         +-----+-----+-----+
// index:  |  2  |  1  |  0  |
//         +-----+-----+-----+
// bitmap: |  1  |  0  |  1  |
//         +-----+-----+-----+
//
// Quorum 1 being 0 means no signers that were supposed to sign actually did
// Quorums 0 and 2 being 1 means at least one signer that was supposed to sign
// actually did (todo: confirm this assumption, it may differ slightly)
//
// Let's assume there exist 6 signers, the first 3 being non-signers
// Let's list their quorum bitmaps. For each signer a quorum bitmap says whether the
// signer was supposed to sign at each quorum (1) or not supposed to sign (0)
//
// In the solidity implementation each of these bitmaps is queried from the reference
// block from the non-signer pubkey, that is, these bitmaps are not passed in
// directly as input to this function
//
// Signer 0 was supposed to sign at quorums 0 and 2 (but assume signed neither)
//         +-----+-----+-----+
// index:  |  2  |  1  |  0  |
//         +-----+-----+-----+
// bitmap: |  1  |  0  |  1  |
//         +-----+-----+-----+
//
// Signer 1 was supposed to sign at quorums 1 and 2 (but assume signed neither)
//         +-----+-----+-----+
// index:  |  2  |  1  |  0  |
//         +-----+-----+-----+
// bitmap: |  1  |  1  |  0  |
//         +-----+-----+-----+
//
// Signer 2 was supposed to sign at all quorums (but assume signed none)
//         +-----+-----+-----+
// index:  |  2  |  1  |  0  |
//         +-----+-----+-----+
// bitmap: |  1  |  1  |  1  |
//         +-----+-----+-----+
//
// Signer 3 was supposed to sign at quorum 2 (assume it did sign it)
//         +-----+-----+-----+
// index:  |  2  |  1  |  0  |
//         +-----+-----+-----+
// bitmap: |  1  |  0  |  0  |
//         +-----+-----+-----+
//
// Signer 4 was supposed to sign at quorum 0 (assume it did sign it)
//         +-----+-----+-----+
// index:  |  2  |  1  |  0  |
//         +-----+-----+-----+
// bitmap: |  0  |  0  |  1  |
//         +-----+-----+-----+
//
// Signer 5 was not supposed to sign at any quorum (assume it did not sign any)
//         +-----+-----+-----+
// index:  |  2  |  1  |  0  |
//         +-----+-----+-----+
// bitmap: |  0  |  0  |  0  |
//         +-----+-----+-----+
//
// The above example bitmaps specify only whether each signer was supposed to sign,
// they say nothing about whether they did or did not sign. So every statement above
// about them signing or not is for the sake of example. The actual non-signers will
// be given through `params.non_signer_pubkeys`, which for this example would be the
// vec [PK0, PK1, PK2]
//
// Since the signature is over the batch root from a tree of all blob certificates
// it means that a signer either signs all quorums it was assigned to or signs none,
// that is, it cannot sign some quorums but not others. This is important for the
// correctness of this implementation (todo: confirm this assumption and its
// conclusion)
//
// The calculation starts by iterating over non-signer quorum bitmaps (in solidity
// each bitmap is queried from the reference block). Each non-signer quorum bitmap
// (the ones shown above) is ANDed with the quorum bitmap:
//
//              signer bitmap      &  quorum numbers bitmap =     result bitmap
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
// This effectively ignores extra signatures that were not called for, that is, if
// Signer 5 was not supposed to sign in Quorum 0 but did anyway (unlike in the
// example above) this would just be ignored
//
// Note that consideration of stakes is a separate matter that is thus calculated
// separately
//
// Signers 0, 1 and 2 are non-signers and the result bitmaps encode how many
// signatures were expected from them:
//
// 2 signatures were expected from Signer 0 (at quorums 0 and 2) but neither was provided
// 1 signature was expected from Signer 1 (at quorum 2) but it was not provided
// 2 signatures were expected from Signer 2 (at quorums 0 and 2) but neither was provided
//
// Note that quorum 1 was not signed at all by any signer so it's excluded from
// further consideration, that is, in what follows the expected (but not provided)
// signatures of Signers 1 and 2 for Quorum 1 will be simply ignored
//
// All of the above considerations translate to an initial aggregate pubkey of:
// APK = -(2*PK0 + 1*PK1 + 2*PK2)
//
// That is, the calculation starts by subtracting pubkeys corresponding to how
// many signatures are missing from them
//
// Each quorum has an associated aggregate pubkey that corresponds to the sum of the
// pubkeys that were supposed to sign:
//
// For Quorum 0, Signers 0, 2 and 4 were supposed to sign
// For Quorum 1, Signers 1 and 2 were supposed to sign
// For Quorum 2, Signers 0, 1 and 2 were supposed to sign
//
// So the aggregate pubkeys of each quorum are:
//
// APK of Quorum 0: PK0 + PK2 + PK4
// APK of Quorum 1: PK1 + PK2 (which will be ignored because there were no signers)
// APK of Quorum 2: PK0 + PK1 + PK2 + PK3
//
// These are provided by `params.quorum_apks`
//
// The resulting aggregate pubkey is the sum of all quorums' aggregate pubkeys and
// the negated aggregate pubkey calculated earlier:
//
//       -    non-signers APK     +   Quorum 0 APK     + Quorum 1 APK +      Quorum 2 APK
// APK = -(2*PK0 + 1*PK1 + 2*PK2) + (PK0 + PK2 + PK4)  +      0       + (PK0 + PK1 + PK2 + PK3)
//
// After cancelling out terms, the resulting APK is PK3 + PK4 as expected since
// those were the only signers that were both expected to sign and did sign
pub fn compute_signers_apk<'a>(
    signed_quorum_numbers: &'a [u8],
    upper_bound_bit_index: u8,
    non_signer_pks: &'a [G1Affine],
    hash_to_bitmap: &'a HashMap<BeHash, Bitmap>,
    quorum_apks: &'a [G1Affine],
) -> Result<G1Affine, SignaturesVerificationError<'a>> {
    let signed_quorum_bitmap = bit_indices_to_bitmap(signed_quorum_numbers, upper_bound_bit_index)?;

    let non_signer_bitmaps = non_signer_pks
        .iter()
        .map(|non_signer_pk| {
            // SignerBitmapNotFound: if the hashed `non_signer_pk` cannot be found in `hash_to_bitmap`
            // OrPubkeyAtInfinity: if `non_signer_pk` is the point at infinity
            point_to_hash(non_signer_pk).and_then(|hash| hash_to_bitmap.get(&hash))
        })
        .collect::<Option<Vec<_>>>()
        .ok_or(SignaturesVerificationError::SignerBitmapNotFoundOrPubkeyAtInfinity)?;

    let non_signers_apk = non_signer_pks
        .iter()
        .zip(non_signer_bitmaps.iter())
        .map(|(non_signer_pk, non_signer_bitmap)| {
            let weight = (**non_signer_bitmap & signed_quorum_bitmap).count_ones();
            let weight = BigInteger256::from(weight as u64);

            non_signer_pk.into_group().mul_bigint(weight)
        })
        .sum::<G1Projective>();

    let total_apk = quorum_apks
        .iter()
        .map(|quorum_apk| quorum_apk.into_group())
        .sum::<G1Projective>();

    let signers_apk = total_apk - non_signers_apk;

    Ok(signers_apk.into_affine())
}

#[cfg(test)]
mod tests {
    use alloc::{vec, vec::Vec};

    use ark_bn254::{G1Affine, G1Projective};
    use ark_ec::{AffineRepr, CurveGroup, PrimeGroup};
    use ark_ff::BigInteger256;
    use bitvec::array::BitArray;
    use hashbrown::HashMap;

    use crate::{
        aggregation::compute_signers_apk, bitmap_utils::Bitmap, convert::point_to_hash,
        error::SignaturesVerificationError, hash::BeHash,
    };

    #[test]
    fn compute_signers_apk_for_3_quorums_and_6_signers() {
        let (non_signer_pks, hash_to_bitmap, quorum_apks) = inputs_for_3_quorums_and_6_signers();

        let signed_quorum_numbers = vec![0, 2];
        let upper_bound_bit_index = u8::MAX;

        let actual = compute_signers_apk(
            &signed_quorum_numbers,
            upper_bound_bit_index,
            &non_signer_pks,
            &hash_to_bitmap,
            &quorum_apks,
        )
        .unwrap();

        let expected = (ppk(3) + ppk(4)).into_affine();

        assert_eq!(actual, expected);
    }

    #[test]
    fn compute_signers_apk_fails_given_too_long_signed_quorum_numbers() {
        let (non_signer_pks, hash_to_bitmap, quorum_apks) = inputs_for_3_quorums_and_6_signers();

        let signed_quorum_numbers = vec![1u8; 256 + 1];
        let upper_bound_bit_index = u8::MAX;

        let result = compute_signers_apk(
            &signed_quorum_numbers,
            upper_bound_bit_index,
            &non_signer_pks,
            &hash_to_bitmap,
            &quorum_apks,
        );

        assert!(result.is_err());
    }

    #[test]
    fn compute_signers_apk_fails_given_empty_hash_to_bitmap() {
        let (non_signer_pks, mut hash_to_bitmap, quorum_apks) =
            inputs_for_3_quorums_and_6_signers();

        hash_to_bitmap.clear();

        let signed_quorum_numbers = vec![0, 2];
        let upper_bound_bit_index = u8::MAX;

        let result = compute_signers_apk(
            &signed_quorum_numbers,
            upper_bound_bit_index,
            &non_signer_pks,
            &hash_to_bitmap,
            &quorum_apks,
        )
        .unwrap_err();

        assert_eq!(
            result,
            SignaturesVerificationError::SignerBitmapNotFoundOrPubkeyAtInfinity
        );
    }

    fn inputs_for_3_quorums_and_6_signers()
    -> (Vec<G1Affine>, HashMap<BeHash, Bitmap>, Vec<G1Affine>) {
        let non_signer_pks = vec![pk(0), pk(1), pk(2)];
        let signer_pks = vec![pk(0), pk(1), pk(2), pk(3), pk(4), pk(5)];

        let signer_bitmaps = vec![
            BitArray::new([5, 0, 0, 0]),
            BitArray::new([6, 0, 0, 0]),
            BitArray::new([7, 0, 0, 0]),
            BitArray::new([4, 0, 0, 0]),
            BitArray::new([1, 0, 0, 0]),
            BitArray::new([0, 0, 0, 0]),
        ];

        let hash_to_bitmap = signer_pks
            .iter()
            .zip(signer_bitmaps.into_iter())
            .map(|(pk, bitmap)| {
                let hash = point_to_hash(pk).unwrap();
                (hash, bitmap)
            })
            .collect::<HashMap<_, _>>();

        // since quorum_apks.len() == quorum_numbers.len() it means Quorum 1 is not even part of  quorum_apks
        let quorum_apks = vec![
            ppk(0) + ppk(2) + ppk(4),          // Quorum 0
            ppk(0) + ppk(1) + ppk(2) + ppk(3), // Quorum 2
        ]
        .into_iter()
        .map(CurveGroup::into_affine)
        .collect::<Vec<_>>();

        (non_signer_pks, hash_to_bitmap, quorum_apks)
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
