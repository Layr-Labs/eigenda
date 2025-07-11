//! Following eigenda-contracts/lib/eigenlayer-middleware/src/BLSSignatureChecker.sol
//! from 6797f3821db92c2214aaa6f137d94c603011ac2a lib/eigenlayer-middleware (v0.5.4-mainnet-rewards-v2-1-g6797f38)

#![no_std]
extern crate alloc;

mod bitmap_utils;
mod error;
mod hashing;
mod types;
mod validation;

use alloc::vec::Vec;

use ark_bn254::{G1Affine, G1Projective};
use ark_ec::{AffineRepr, CurveGroup, PrimeGroup};
use ark_ff::BigInteger256;
use bitmap_utils::{Bitmap, bit_indices_to_bitmap};
use error::SignaturesVerificationError;
use hashbrown::HashMap;
use hashing::Hash;
use types::{NonSignerInfo, NonSignerStakesAndSignature, QuorumStakeTotals, ReferenceBlock};

#[derive(Default, Debug)]
pub struct SignaturesVerification {
    pub quorum_stake_totals: QuorumStakeTotals,
    pub signatory_record_hash: Hash,
}

pub trait SignatureVerifier {
    /// In the context of EigenDA blob certificates this is the hash of the batchRoot
    ///
    /// Quorums for which there were signatures i.e. if 3 quorums exist and
    /// `quorum_numbers = [0, 2]` it means Quorum 1 has not been signed by any signer
    /// whereas Quorums 0 and 2 were signed by at least one signer
    ///
    /// The block at which stake information is queried
    fn verify_signatures<'a>(
        &'a self,
        msg_hash: Hash,
        quorum_numbers: &'a [u8],
        reference_block_number: u32,
        current_block_number: u32,
        params: &'a NonSignerStakesAndSignature,
        reference_block: &'a ReferenceBlock,
    ) -> Result<SignaturesVerification, SignaturesVerificationError<'a>>;
}

#[derive(Default, Debug)]
pub struct BlsSignaturesVerifier;

impl SignatureVerifier for BlsSignaturesVerifier {
    fn verify_signatures<'a>(
        &'a self,
        _msg_hash: Hash,
        signing_quorum_numbers: &'a [u8],
        reference_block_number: u32,
        current_block_number: u32,
        params: &'a NonSignerStakesAndSignature,
        reference_block: &'a ReferenceBlock,
    ) -> Result<SignaturesVerification, SignaturesVerificationError<'a>> {
        validation::validate_inputs(
            signing_quorum_numbers,
            reference_block_number,
            current_block_number,
            params,
        )?;

        let _signers_aggregate_pubkey = compute_signers_apk(
            signing_quorum_numbers,
            u8::MAX, // todo
            &params.non_signer_pubkeys,
            &reference_block.hash_to_bitmap,
            &params.quorum_apks,
        )?;

        let _stake_totals = QuorumStakeTotals::default();
        let _non_signer_info = NonSignerInfo::default();

        Ok(SignaturesVerification::default())
    }
}

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
fn compute_signers_apk<'a>(
    signing_quorum_numbers: &'a [u8],
    upper_bound_bit_index: u8,
    non_signer_pks: &'a [G1Affine],
    hash_to_bitmap: &'a HashMap<Hash, Bitmap>,
    quorum_apks: &'a [G1Affine],
) -> Result<G1Affine, SignaturesVerificationError<'a>> {
    let signing_quorum_bitmap =
        bit_indices_to_bitmap(signing_quorum_numbers, upper_bound_bit_index)?;

    let non_signer_bitmaps = non_signer_pks
        .iter()
        .map(|non_signer_pk| {
            let hash = hashing::hash_g1_point(non_signer_pk);
            hash_to_bitmap.get(&hash)
        })
        .collect::<Option<Vec<_>>>()
        .ok_or(SignaturesVerificationError::SignerBitmapNotFound)?;

    let non_signers_apk = non_signer_pks
        .iter()
        .zip(non_signer_bitmaps.iter())
        .map(|(non_signer_pk, non_signer_bitmap)| {
            let weight = (**non_signer_bitmap & signing_quorum_bitmap).count_ones();
            let weight = BigInteger256::from(weight as u64);

            non_signer_pk.into_group().mul_bigint(weight)
        })
        .sum::<G1Projective>();

    debug_assert!(quorum_apks.len() == signing_quorum_numbers.len());
    let total_apk = quorum_apks
        .iter()
        .map(|quorum_apk| quorum_apk.into_group())
        .sum::<G1Projective>();

    let signers_apk = total_apk - non_signers_apk;

    Ok(signers_apk.into_affine())
}

fn _compute_stakes() {}

#[cfg(test)]
mod tests {
    use alloc::{vec, vec::Vec};

    use ark_bn254::G1Affine;
    use ark_ec::{AffineRepr, CurveGroup, PrimeGroup};
    use ark_ff::BigInteger256;
    use bitvec::array::BitArray;
    use hashbrown::HashMap;

    use crate::{compute_signers_apk, hashing::hash_g1_point};

    #[test]
    fn test_compute_signers_apk_for_3_quorums_and_6_signers() {
        let signing_quorum_numbers = vec![0, 2];
        let upper_bound_bit_index = u8::MAX;

        let generator = G1Affine::generator();
        let ppk = |n: u64| {
            generator
                .into_group()
                .mul_bigint(BigInteger256::from(n + 1))
        };
        let pk = |n: u64| ppk(n).into_affine();

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
                let hash = hash_g1_point(pk);
                (hash, bitmap)
            })
            .collect::<HashMap<_, _>>();

        // since quorum_apks.len() == quorum_numbers.len() it means Quorum 1 is not even part of  quorum_apks
        let quorum_apks = vec![
            ppk(0) + ppk(2) + ppk(4), // Quorum 0
            ppk(0) + ppk(1) + ppk(2) + ppk(3) // Quorum 2
        ]
            .into_iter()
            .map(CurveGroup::into_affine)
            .collect::<Vec<_>>();

        let actual = compute_signers_apk(
            &signing_quorum_numbers,
            upper_bound_bit_index,
            &non_signer_pks,
            &hash_to_bitmap,
            &quorum_apks,
        )
        .unwrap();

        let expected = (ppk(3) + ppk(4)).into_affine();

        assert_eq!(actual, expected);
    }
}
