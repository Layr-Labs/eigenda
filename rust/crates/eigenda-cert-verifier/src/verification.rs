use core::iter::once;

use alloc::vec::Vec;
use alloy_primitives::Bytes;
use alloy_sol_types::SolValue;
use ark_bn254::{G1Affine, G2Affine};
use eigenda_cert::{BatchHeaderV2, BlobInclusionInfo, NonSignerStakesAndSignature};

use crate::{
    check, convert,
    error::CertVerificationError::{self, *},
    hash::{self, HashExt, keccak_v256},
    signature,
    types::{
        NonSigner, Quorum, Stake, Storage, conversions::IntoExt, solidity::SecurityThresholds,
    },
};

pub fn verify(
    batch_header: BatchHeaderV2,
    blob_inclusion_info: BlobInclusionInfo,
    non_signer_stakes_and_signature: NonSignerStakesAndSignature,
    security_thresholds: SecurityThresholds,
    required_quorum_numbers: Bytes,
    signed_quorum_numbers: Bytes,
    storage: Storage,
) -> Result<(), CertVerificationError> {
    let Storage {
        initialized_quorums_count,
        current_block,
        reject_staleness,
        min_withdrawal_delay_blocks,
        quorum_membership_history_by_signer,
        stake_history_by_signer_and_quorum,
        total_stake_history_by_quorum,
        apk_trunc_hash_history_by_quorum,
        last_updated_at_block_by_quorum,
        relay_key_to_relay_info,
        version_to_versioned_blob_params,
    } = storage;

    let blob_certificate = blob_inclusion_info.blob_certificate.hash_ext();
    let encoded = blob_certificate.abi_encode_packed();
    let leaf_node = keccak_v256(once(encoded));
    check::leaf_node_belongs_to_merkle_tree(
        leaf_node,
        batch_header.batch_root.into(),
        blob_inclusion_info.inclusion_proof,
        blob_inclusion_info.blob_index,
    )?;

    if batch_header.reference_block_number >= current_block {
        return Err(ReferenceBlockDoesNotPrecedeCurrentBlock);
    }

    let lengths = [
        non_signer_stakes_and_signature.non_signer_pubkeys.len(),
        non_signer_stakes_and_signature
            .non_signer_quorum_bitmap_indices
            .len(),
    ];
    check::non_zero_equal_lengths(&lengths)?;

    let lengths = [
        signed_quorum_numbers.len(),
        non_signer_stakes_and_signature.quorum_apks.len(),
        non_signer_stakes_and_signature.quorum_apk_indices.len(),
        non_signer_stakes_and_signature.total_stake_indices.len(),
        non_signer_stakes_and_signature
            .non_signer_stake_indices
            .len(),
    ];
    check::non_zero_equal_lengths(&lengths)?;

    if reject_staleness {
        check::quorums_last_updated_after_most_recent_stale_block(
            &signed_quorum_numbers,
            batch_header.reference_block_number,
            last_updated_at_block_by_quorum,
            min_withdrawal_delay_blocks,
        )?;
    }

    check::cert_apks_equal_storage_apks(
        &signed_quorum_numbers,
        batch_header.reference_block_number,
        &non_signer_stakes_and_signature.quorum_apks,
        non_signer_stakes_and_signature.quorum_apk_indices,
        apk_trunc_hash_history_by_quorum,
    )?;

    // assumption: collection_a[i] corresponds to collection_b[i] for all i
    let non_signers = non_signer_stakes_and_signature
        .non_signer_pubkeys
        .into_iter()
        .zip(
            non_signer_stakes_and_signature
                .non_signer_quorum_bitmap_indices
                .into_iter(),
        )
        .map(|(pk, quorum_membership_index)| {
            let pk_hash = convert::point_to_hash(&pk);

            let quorum_membership = quorum_membership_history_by_signer
                .get(&pk_hash)
                .ok_or(MissingSignerEntry)?
                .try_get_at(quorum_membership_index)?
                .try_get_against(batch_header.reference_block_number)?;

            let pk: G1Affine = pk.into_ext();
            let non_signer = NonSigner {
                pk,
                pk_hash,
                quorum_membership,
            };
            Ok(non_signer)
        })
        .collect::<Result<Vec<_>, _>>()?;

    check::non_signers_strictly_sorted_by_hash(&non_signers)?;

    // assumption: collection_a[i] corresponds to collection_b[i] for all i, for all (a, b)
    let quorums = signed_quorum_numbers
        .into_iter()
        .zip(non_signer_stakes_and_signature.quorum_apks.into_iter())
        .zip(
            non_signer_stakes_and_signature
                .total_stake_indices
                .into_iter(),
        )
        .zip(
            non_signer_stakes_and_signature
                .non_signer_stake_indices
                .into_iter(),
        )
        .map(
            |(
                ((signed_quorum, apk), total_stake_index),
                stake_index_for_each_required_non_signer,
            )| {
                let total_stake = total_stake_history_by_quorum
                    .get(&signed_quorum)
                    .ok_or(MissingQuorumEntry)?
                    .try_get_at(total_stake_index)?
                    .try_get_against(batch_header.reference_block_number)?;

                let bit = signed_quorum as usize;
                let unsigned_stake = non_signers
                    .iter()
                    .filter(|non_signer| {
                        let was_required_to_sign_this_quorum = non_signer.quorum_membership[bit];
                        was_required_to_sign_this_quorum
                    })
                    // assumption: collection_a[i] corresponds to collection_b[i] for all i
                    .zip(stake_index_for_each_required_non_signer.into_iter())
                    .map(|(required_non_signer, stake_index)| {
                        stake_history_by_signer_and_quorum
                            .get(&required_non_signer.pk_hash)
                            .ok_or(MissingSignerEntry)?
                            .get(&signed_quorum)
                            .ok_or(MissingQuorumEntry)?
                            .try_get_at(stake_index)?
                            .try_get_against(batch_header.reference_block_number)
                    })
                    .sum::<Result<Stake, _>>()?;

                let signed_stake = total_stake.checked_sub(unsigned_stake).ok_or(Underflow)?;

                let apk: G1Affine = apk.into_ext();
                let quorum = Quorum {
                    number: signed_quorum,
                    apk,
                    total_stake,
                    signed_stake,
                };

                Ok(quorum)
            },
        )
        .collect::<Result<Vec<_>, _>>()?;

    let signers_apk =
        signature::aggregation::aggregate(initialized_quorums_count, &non_signers, &quorums)?;

    let msg_hash = batch_header.hash_ext();
    let apk_g2: G2Affine = non_signer_stakes_and_signature.apk_g2.into_ext();
    let sigma: G1Affine = non_signer_stakes_and_signature.sigma.into_ext();

    if signature::verification::verify(msg_hash, signers_apk, apk_g2, sigma) == false {
        return Err(SignatureVerificationFailed);
    }

    let pk_hashes = non_signers
        .iter()
        .map(|non_signer| non_signer.pk_hash)
        .collect::<Vec<_>>();
    let _signatory_record_hash =
        hash::signature_record(batch_header.reference_block_number, &pk_hashes);

    check::relay_keys_are_set(
        &blob_inclusion_info.blob_certificate.relay_keys,
        &relay_key_to_relay_info,
    )?;

    let version = blob_inclusion_info.blob_certificate.blob_header.version;
    check::security_assumptions_are_met(
        version,
        &version_to_versioned_blob_params,
        &security_thresholds,
    )?;

    let blob_quorums = blob_inclusion_info
        .blob_certificate
        .blob_header
        .quorum_numbers;

    check::confirmed_quorums_contain_blob_quorums(
        security_thresholds.confirmationThreshold,
        &quorums,
        &blob_quorums,
    )?;

    check::blob_quorums_contain_required_quorums(&blob_quorums, &required_quorum_numbers)?;

    Ok(())
}

#[cfg(test)]
mod tests {
    use core::iter::once;

    use alloc::vec;
    use alloc::vec::Vec;
    use alloy_primitives::{B256, Bytes};
    use alloy_sol_types::SolValue;
    use ark_bn254::{Fr, G1Affine, G1Projective, G2Projective};
    use ark_ec::{CurveGroup, PrimeGroup};
    use eigenda_cert::{
        BatchHeaderV2, BlobCertificate, BlobCommitment, BlobHeaderV2, BlobInclusionInfo, G1Point,
        NonSignerStakesAndSignature,
    };
    use hashbrown::HashMap;

    use crate::{
        bitmap::Bitmap,
        convert,
        error::CertVerificationError::*,
        hash::{HashExt, keccak_v256},
        types::{
            Storage,
            conversions::{DefaultExt, IntoExt},
            history::{History, Update},
            solidity::{RelayInfo, SecurityThresholds, VersionedBlobParams},
        },
        verification,
    };

    #[test]
    fn success() {
        let (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        ) = success_inputs();

        let result = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        );
        assert_eq!(result, Ok(()));
    }

    #[test]
    fn leaf_node_does_not_belong_to_merkle_tree() {
        let (
            batch_header,
            mut blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        ) = success_inputs();

        // any change to blobCertificate causes the leaf node hash to differ
        blob_inclusion_info.blob_certificate.signature = [0u8; 32].into();

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, LeafNodeDoesNotBelongToMerkleTree);
    }

    #[test]
    fn reference_block_past_current_block() {
        let (
            mut batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            mut storage,
        ) = success_inputs();

        batch_header.reference_block_number = 43;
        storage.current_block = 42;

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, ReferenceBlockDoesNotPrecedeCurrentBlock);
    }

    #[test]
    fn reference_block_at_current_block() {
        let (
            mut batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            mut storage,
        ) = success_inputs();

        batch_header.reference_block_number = 42;
        storage.current_block = 42;

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, ReferenceBlockDoesNotPrecedeCurrentBlock);
    }

    #[test]
    fn empty_non_signer_vecs() {
        let (
            batch_header,
            blob_inclusion_info,
            mut non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        ) = success_inputs();

        non_signer_stakes_and_signature.non_signer_pubkeys.clear();

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, EmptyVec);
    }

    #[test]
    fn empty_quorum_vecs() {
        let (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            _signed_quorum_numbers,
            storage,
        ) = success_inputs();

        let signed_quorum_numbers: Bytes = [].into();

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, EmptyVec);
    }

    #[test]
    fn reject_staleness() {
        let (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            mut storage,
        ) = success_inputs();

        storage.reject_staleness = true;
        storage.last_updated_at_block_by_quorum.insert(0, 41);

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, StaleQuorum);
    }

    #[test]
    fn cert_apk_not_equal_storage_apk() {
        let (
            batch_header,
            blob_inclusion_info,
            mut non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        ) = success_inputs();

        non_signer_stakes_and_signature.quorum_apks[0] = G1Point {
            x: alloy_primitives::Uint::ONE,
            y: alloy_primitives::Uint::ONE,
        };

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, CertApkDoesNotEqualStorageApk);
    }

    #[test]
    fn quorum_membership_history_missing_signer_entry() {
        let (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            mut storage,
        ) = success_inputs();

        storage.quorum_membership_history_by_signer.clear();

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, MissingSignerEntry);
    }

    #[test]
    fn quorum_membership_history_missing_history_entry() {
        let (
            batch_header,
            blob_inclusion_info,
            mut non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        ) = success_inputs();

        non_signer_stakes_and_signature.non_signer_quorum_bitmap_indices[0] = 42;

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, MissingHistoryEntry);
    }

    #[test]
    fn quorum_membership_history_reference_block_not_in_interval() {
        let (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            mut storage,
        ) = success_inputs();

        storage
            .quorum_membership_history_by_signer
            .iter_mut()
            .for_each(|(_, v)| {
                v.0.insert(0, Update::new(141, 143, Default::default()).unwrap());
            });

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, ElementNotInInterval);
    }

    #[test]
    fn non_signers_not_strictly_sorted_by_hash() {
        let (
            batch_header,
            blob_inclusion_info,
            mut non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        ) = success_inputs();

        non_signer_stakes_and_signature.non_signer_pubkeys.reverse();

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, NotStrictlySortedByHash);
    }

    #[test]
    fn total_stake_history_missing_quorum_entry() {
        let (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            mut storage,
        ) = success_inputs();

        storage.total_stake_history_by_quorum.clear();

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, MissingQuorumEntry);
    }

    #[test]
    fn total_stake_history_missing_history_entry() {
        let (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            mut storage,
        ) = success_inputs();

        storage
            .total_stake_history_by_quorum
            .insert(0, Default::default());

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, MissingHistoryEntry);
    }

    #[test]
    fn total_stake_history_reference_block_not_in_interval() {
        let (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            mut storage,
        ) = success_inputs();

        storage
            .total_stake_history_by_quorum
            .iter_mut()
            .for_each(|(_, v)| {
                v.0.insert(0, Update::new(141, 143, Default::default()).unwrap());
            });

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, ElementNotInInterval);
    }

    #[test]
    fn stake_history_missing_signer_entry() {
        let (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            mut storage,
        ) = success_inputs();

        storage.stake_history_by_signer_and_quorum.clear();

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, MissingSignerEntry);
    }

    #[test]
    fn stake_history_missing_quorum_entry() {
        let (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            mut storage,
        ) = success_inputs();

        storage
            .stake_history_by_signer_and_quorum
            .iter_mut()
            .for_each(|(_, stake_history_by_quorum)| {
                stake_history_by_quorum.clear();
            });

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, MissingQuorumEntry);
    }

    #[test]
    fn stake_history_missing_history_entry() {
        let (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            mut storage,
        ) = success_inputs();

        storage
            .stake_history_by_signer_and_quorum
            .iter_mut()
            .for_each(|(_, stake_history_by_quorum)| {
                stake_history_by_quorum.insert(0, Default::default());
            });

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, MissingHistoryEntry);
    }

    #[test]
    fn stake_history_reference_block_not_in_interval() {
        let (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            mut storage,
        ) = success_inputs();

        storage
            .stake_history_by_signer_and_quorum
            .iter_mut()
            .for_each(|(_, stake_history_by_quorum)| {
                stake_history_by_quorum.iter_mut().for_each(|(_, v)| {
                    v.0.insert(0, Update::new(141, 143, Default::default()).unwrap());
                })
            });

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, ElementNotInInterval);
    }

    #[test]
    fn stake_underflow() {
        let (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            mut storage,
        ) = success_inputs();

        storage
            .total_stake_history_by_quorum
            .iter_mut()
            .for_each(|(_, v)| {
                v.0.insert(0, Update::new(41, 43, 29).unwrap());
            });

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, Underflow);
    }

    #[test]
    fn aggregation_failure() {
        let (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            mut storage,
        ) = success_inputs();

        storage.initialized_quorums_count = 1;

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, BitIndexNotLessThanUpperBound);
    }

    #[test]
    fn signature_verification_failure() {
        let (
            batch_header,
            blob_inclusion_info,
            mut non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        ) = success_inputs();

        non_signer_stakes_and_signature.sigma = G1Point::default_ext();

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, SignatureVerificationFailed);
    }

    #[test]
    fn relay_keys_not_set() {
        let (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            mut storage,
        ) = success_inputs();

        let relay_info = storage.relay_key_to_relay_info.get_mut(&42).unwrap();
        relay_info.relayAddress = Default::default();

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, RelayKeyNotSet);
    }

    #[test]
    fn security_assumptions_not_met() {
        let (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            mut storage,
        ) = success_inputs();

        let params = storage
            .version_to_versioned_blob_params
            .get_mut(&42)
            .unwrap();
        params.numChunks = 43;

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, UnmetSecurityAssumptions);
    }

    #[test]
    fn confirmed_quorums_do_not_contain_blob_quorums() {
        let (
            _batch_header,
            mut blob_inclusion_info,
            mut non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            mut storage,
        ) = success_inputs();

        storage
            .version_to_versioned_blob_params
            .iter_mut()
            .for_each(|(_, versioned_blob_params)| {
                versioned_blob_params.maxNumOperators = 0;
            });

        blob_inclusion_info
            .blob_certificate
            .blob_header
            .quorum_numbers = [0, 1, 2].into(); // while confirmed_quorums: [0, 2]

        // any change to blobCertificate requires recomputing...
        let secret_keys = vec![Fr::from(43u64), Fr::from(44u64)];
        let (batch_header, sigma) =
            compute_batch_header_and_sigma(&blob_inclusion_info, secret_keys);

        non_signer_stakes_and_signature.sigma = sigma.into_ext();

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, ConfirmedQuorumsDoNotContainBlobQuorums);
    }

    #[test]
    fn blob_quorums_do_not_contain_required_quorums() {
        let (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            _required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        ) = success_inputs();

        let required_quorum_numbers: Bytes = [1].into(); // 3 is not in blob_quorums: [0, 2]

        let err = verification::verify(
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
        .unwrap_err();

        assert_eq!(err, BlobQuorumsDoNotContainRequiredQuorums);
    }

    fn success_inputs() -> (
        BatchHeaderV2,
        BlobInclusionInfo,
        NonSignerStakesAndSignature,
        SecurityThresholds,
        Bytes,
        Bytes,
        Storage,
    ) {
        let g1 = G1Projective::generator();
        let g2 = G2Projective::generator();

        let non_signer0_sk = Fr::from(40u64);
        let non_signer0_g1_pk = (g1 * non_signer0_sk).into_affine();

        let non_signer1_sk = Fr::from(41u64);
        let non_signer1_g1_pk = (g1 * non_signer1_sk).into_affine();

        let non_signer2_sk = Fr::from(42u64);
        let non_signer2_g1_pk = (g1 * non_signer2_sk).into_affine();

        let signer3_sk = Fr::from(43u64);
        let signer3_g1_pk = (g1 * signer3_sk).into_affine();
        let signer3_g2_pk = (g2 * signer3_sk).into_affine();

        let signer4_sk = Fr::from(44u64);
        let signer4_g1_pk = (g1 * signer4_sk).into_affine();
        let signer4_g2_pk = (g2 * signer4_sk).into_affine();

        let optional_non_signer5_sk = Fr::from(45u64);
        let optional_non_signer5_g1_pk = (g1 * optional_non_signer5_sk).into_affine();

        let _apk_g1 = (signer3_g1_pk + signer4_g1_pk).into_affine();
        let apk_g2 = (signer3_g2_pk + signer4_g2_pk).into_affine();

        let blob_inclusion_info = BlobInclusionInfo {
            blob_certificate: BlobCertificate {
                blob_header: BlobHeaderV2 {
                    version: 42,
                    quorum_numbers: [0, 2].into(),
                    commitment: BlobCommitment::default_ext(),
                    payment_header_hash: [42; 32].into(),
                },
                signature: [].into(),
                relay_keys: vec![42],
            },
            blob_index: 0,
            inclusion_proof: [42u8; 32].into(),
        };

        let (batch_header, sigma) =
            compute_batch_header_and_sigma(&blob_inclusion_info, vec![signer3_sk, signer4_sk]);

        // let sig_at_quorum_2_by_signer_3 = (msg_point * signer3_sk).into_affine();
        // let sig_at_quorum_0_by_signer_4 = (msg_point * signer4_sk).into_affine();
        // let sigma = (sig_at_quorum_2_by_signer_3 + sig_at_quorum_0_by_signer_4).into_affine();

        let apk_for_each_quorum = [
            (non_signer0_g1_pk + non_signer2_g1_pk + signer4_g1_pk).into_affine(),
            (non_signer0_g1_pk + non_signer1_g1_pk + non_signer2_g1_pk + signer3_g1_pk)
                .into_affine(),
        ];

        let non_signer_stakes_and_signature = NonSignerStakesAndSignature {
            non_signer_quorum_bitmap_indices: vec![0, 0, 0],
            non_signer_pubkeys: vec![
                non_signer0_g1_pk.into_ext(),
                non_signer1_g1_pk.into_ext(),
                non_signer2_g1_pk.into_ext(),
            ],
            quorum_apks: vec![
                apk_for_each_quorum[0].into_ext(),
                apk_for_each_quorum[1].into_ext(),
            ],
            apk_g2: apk_g2.into_ext(),
            sigma: sigma.into_ext(),
            quorum_apk_indices: vec![0, 0],
            total_stake_indices: vec![0, 0],
            non_signer_stake_indices: vec![vec![0, 0, 0], vec![0, 0, 0]],
        };
        // quorum 1 had no signatures
        // quorums 0 and 2 had at least one signature (exactly one in this example)
        let signed_quorum_numbers: Bytes = [0, 2].into();

        let security_thresholds = SecurityThresholds {
            // further down I set codingRate = 42
            // since (confirmation_threshold - adversary_threshold) * codingRate >= 100
            // and confirmation_threshold > adversary_threshold
            // I set the following:
            // the above condition would be met with confirmation_threshold: 100
            // but would result in n = 0 in `n < maxNumOperators` thus not meeting security assumptions
            confirmationThreshold: 66,
            adversaryThreshold: 0,
        };

        let non_signer0_pk_hash = convert::point_to_hash(&non_signer0_g1_pk.into_ext());
        let non_signer1_pk_hash = convert::point_to_hash(&non_signer1_g1_pk.into_ext());
        let non_signer2_pk_hash = convert::point_to_hash(&non_signer2_g1_pk.into_ext());
        let signer3_pk_hash = convert::point_to_hash(&signer3_g1_pk.into_ext());
        let signer4_pk_hash = convert::point_to_hash(&signer4_g1_pk.into_ext());
        let optional_non_signer5_pk_hash =
            convert::point_to_hash(&optional_non_signer5_g1_pk.into_ext());

        // by sheer coincidence the first 3 hashes are already sorted
        let pk_hashes = [
            non_signer0_pk_hash,
            non_signer1_pk_hash,
            non_signer2_pk_hash,
            signer3_pk_hash,
            signer4_pk_hash,
            optional_non_signer5_pk_hash,
        ];

        let quorum_membership_history_by_signer = {
            let quorum_memberships = vec![
                Bitmap::new([5, 0, 0, 0]), // 1 0 1
                Bitmap::new([6, 0, 0, 0]), // 1 1 0
                Bitmap::new([7, 0, 0, 0]), // 1 1 1
                Bitmap::new([4, 0, 0, 0]), // 1 0 0
                Bitmap::new([1, 0, 0, 0]), // 0 0 1
                Bitmap::new([0, 0, 0, 0]), // 0 0 0
            ];

            pk_hashes
                .into_iter()
                .zip(quorum_memberships.into_iter())
                .map(|(pk_hash, quorum_membership)| {
                    let update = Update::new(41, 43, quorum_membership).unwrap();
                    let history = HashMap::from([(0, update)]);
                    (pk_hash, History(history))
                })
                .collect()
        };

        let stake_history_by_signer_and_quorum = pk_hashes
            .into_iter()
            .map(|pk_hash| {
                let stake_history_by_quorum = signed_quorum_numbers
                    .clone()
                    .into_iter()
                    .map(|quorum| {
                        let update = Update::new(41, 43, 10).unwrap();
                        let history = HashMap::from([(0, update)]);
                        (quorum, History(history))
                    })
                    .collect();
                (pk_hash, stake_history_by_quorum)
            })
            .collect::<HashMap<B256, _>>();

        let total_stake_history_by_quorum = signed_quorum_numbers
            .clone()
            .into_iter()
            .map(|quorum| {
                let update = Update::new(41, 43, 100).unwrap();
                let history = HashMap::from([(0, update)]);
                (quorum, History(history))
            })
            .collect();

        let apk_trunc_hash_history_by_quorum = signed_quorum_numbers
            .clone()
            .into_iter()
            .zip(apk_for_each_quorum)
            .map(|(quorum, apk)| {
                let apk_hash = convert::point_to_hash(&apk.into_ext());
                let apk_trunc_hash: [u8; 24] = apk_hash[..24].try_into().unwrap();
                let update = Update::new(41, 43, apk_trunc_hash).unwrap();
                let history = HashMap::from([(0, update)]);
                (quorum, History(history))
            })
            .collect();

        let last_updated_at_block_by_quorum = signed_quorum_numbers
            .clone()
            .into_iter()
            .map(|quorum| (quorum, 42))
            .collect();

        let relay_key_to_relay_info = HashMap::from([(
            42,
            RelayInfo {
                relayAddress: [42u8; 20].into(),
                relayURL: Default::default(),
            },
        )]);

        let version_to_versioned_blob_params = HashMap::from([(
            42,
            VersionedBlobParams {
                maxNumOperators: 42,
                numChunks: 44,
                codingRate: 42,
            },
        )]);

        let storage = Storage {
            initialized_quorums_count: u8::MAX,
            current_block: 43,
            reject_staleness: false,
            min_withdrawal_delay_blocks: 1,
            quorum_membership_history_by_signer,
            stake_history_by_signer_and_quorum,
            total_stake_history_by_quorum,
            apk_trunc_hash_history_by_quorum,
            last_updated_at_block_by_quorum,
            relay_key_to_relay_info,
            version_to_versioned_blob_params,
        };

        let required_quorum_numbers: Bytes = [0, 2].into();

        (
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        )
    }

    fn compute_batch_header_and_sigma(
        blob_inclusion_info: &BlobInclusionInfo,
        secret_keys: Vec<Fr>,
    ) -> (BatchHeaderV2, G1Affine) {
        //   C || 42
        //  /      \
        // C        42

        let encoded = blob_inclusion_info
            .blob_certificate
            .hash_ext()
            .abi_encode_packed();
        let left_child = keccak_v256(once(encoded));

        let right_sibling = [42u8; 32].into();
        let batch_root = keccak_v256([left_child, right_sibling].into_iter());

        let batch_header = BatchHeaderV2 {
            batch_root: batch_root.into(),
            reference_block_number: 42,
        };

        let msg_hash = batch_header.hash_ext();
        let msg_point = convert::hash_to_point(msg_hash);

        let sigma = secret_keys
            .iter()
            .map(|secret_key| msg_point * secret_key)
            .sum::<G1Projective>()
            .into_affine();

        (batch_header, sigma)
    }
}
