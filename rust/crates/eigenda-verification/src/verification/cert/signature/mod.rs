//! BLS signature operations for EigenDA certificate verification
//!
//! This module provides BLS signature aggregation and verification functionality
//! specifically tailored for EigenDA's operator signature scheme. It handles
//! the logic of aggregating operator public keys while accounting for
//! non-signing operators and verifying the resulting signatures against batch commitments.
//!
//! ## Key Components
//!
//! - [`aggregation`]: Computes aggregate public keys from quorum operators, handling non-signers
//! - [`verification`]: Verifies BLS signatures using bilinear pairings
//!
//! ## BLS Signature Scheme
//!
//! EigenDA uses BLS signatures on the BN254 curve to enable efficient signature aggregation.
//! Multiple operators can sign the same message, and their signatures can be combined into
//! a single aggregate signature that can be verified against an aggregate public key.

/// BLS public key aggregation logic for combining operator keys.
pub mod aggregation;
/// BLS signature verification using bilinear pairings.
pub mod verification;

#[cfg(test)]
mod tests {
    use std::str::FromStr;

    use alloy_primitives::{B256, U256};

    use crate::cert::{G1Point, G2Point, NonSignerStakesAndSignature};
    use crate::verification::cert::signature::aggregation::aggregate;
    use crate::verification::cert::signature::verification::verify;
    use crate::verification::cert::types::{NonSigner, Quorum, Stake};

    #[test]
    fn signature_verification_without_non_signers() {
        let msg_hash =
            B256::from_str("0xc11f0d6546b185e583cb7d31824c0fdf4af1dc04579fcbb5538ff6c205f6ecc4")
                .unwrap();

        let params = NonSignerStakesAndSignature {
            non_signer_quorum_bitmap_indices: vec![],
            non_signer_pubkeys: vec![],
            quorum_apks: vec![
                G1Point {
                    x: U256::from_str(
                        "647887176094346434688797418329165908112788375706471933112226398612018692311",
                    )
                    .unwrap(),
                    y: U256::from_str(
                        "14219015594739757037737335153756242541699018088640667335296076363950011933479",
                    )
                    .unwrap(),
                },
                G1Point {
                    x: U256::from_str(
                        "6182682689227032767282175811228041488012494622337860227375748139742433007060",
                    )
                    .unwrap(),
                    y: U256::from_str(
                        "3937555473299642407446407290166920042709516259189610965714253279007332654630",
                    )
                    .unwrap(),
                },
            ],
            apk_g2: G2Point {
                x: vec![
                    U256::from_str(
                        "2971582905681448632396838815389593577218918217682961002224335998108796877821",
                    )
                    .unwrap(),
                    U256::from_str(
                        "20493015775924070127190293208207752271841430906645021627145870133490690913120",
                    )
                    .unwrap(),
                ],
                y: vec![
                    U256::from_str(
                        "1352394632334497324545086446186502637904528128084134970457703718550262010278",
                    )
                    .unwrap(),
                    U256::from_str(
                        "2360571446350899391547904541365466568108120225676871506677828765446847764586",
                    )
                    .unwrap(),
                ],
            },
            sigma: G1Point {
                x: U256::from_str(
                    "7229513079519707806356434796736516602069750608278578152681096587215959229139",
                )
                .unwrap(),
                y: U256::from_str(
                    "11534913467352427575310279662799880782898289594350659580468941325380622942260",
                )
                .unwrap(),
            },
            quorum_apk_indices: vec![1873, 2247],
            total_stake_indices: vec![2500, 2541],
            non_signer_stake_indices: vec![vec![], vec![]],
        };

        let signed_quorums_numbers = [0u8, 1u8];

        let quorums = signed_quorums_numbers
            .iter()
            .zip(params.quorum_apks.iter())
            .map(|(number, apk)| Quorum {
                number: *number,
                apk: (*apk).into(),
                total_stake: Stake::default(),
                signed_stake: Stake::default(),
            })
            .collect::<Vec<_>>();

        let non_signers: Vec<NonSigner> = vec![];

        let apk_g1 = aggregate(u8::MAX, &non_signers, &quorums).unwrap();

        let is_signature_valid =
            verify(msg_hash, apk_g1, params.apk_g2.into(), params.sigma.into());

        assert!(is_signature_valid);
    }
}
