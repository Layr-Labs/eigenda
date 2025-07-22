use alloc::vec::Vec;

use crate::{
    check, convert,
    error::CertVerificationError::{self, *},
    hash, signature,
    types::{Cert, Chain, NonSigner, NonSignerStakesAndSignature, Quorum, Stake},
};

pub fn verify(cert: Cert, chain: Chain) -> Result<(), CertVerificationError> {
    let Cert {
        msg_hash,
        reference_block,
        signed_quorums,
        params,
    } = cert;

    let NonSignerStakesAndSignature {
        apk_for_each_quorum,
        apk_index_for_each_quorum,
        total_stake_index_for_each_quorum,
        stake_index_for_each_quorum_and_required_non_signer,
        pk_for_each_non_signer,
        quorum_membership_index_for_each_non_signer,
        apk_g2,
        sigma,
    } = params;

    let Chain {
        initialized_quorums_count,
        current_block,
        reject_staleness,
        min_withdrawal_delay_blocks,
        quorum_membership_history_by_signer,
        stake_history_by_signer_and_quorum,
        total_stake_history_by_quorum,
        apk_trunc_hash_history_by_quorum,
        last_updated_at_block_by_quorum,
    } = chain;

    if reference_block >= current_block {
        return Err(ReferenceBlockDoesNotPrecedeCurrentBlock);
    }

    let lengths = [
        pk_for_each_non_signer.len(),
        quorum_membership_index_for_each_non_signer.len(),
    ];
    check::non_zero_equal_lengths(&lengths)?;

    let lengths = [
        signed_quorums.len(),
        apk_for_each_quorum.len(),
        apk_index_for_each_quorum.len(),
        total_stake_index_for_each_quorum.len(),
        stake_index_for_each_quorum_and_required_non_signer.len(),
    ];
    check::non_zero_equal_lengths(&lengths)?;

    if reject_staleness {
        check::quorums_last_updated_after_most_recent_stale_block(
            &signed_quorums,
            reference_block,
            last_updated_at_block_by_quorum,
            min_withdrawal_delay_blocks,
        )?;
    }

    check::cert_apks_equal_chain_apks(
        &signed_quorums,
        reference_block,
        &apk_for_each_quorum,
        apk_index_for_each_quorum,
        apk_trunc_hash_history_by_quorum,
    )?;

    // assumption: collection_a[i] corresponds to collection_b[i] for all i
    let non_signers = pk_for_each_non_signer
        .into_iter()
        .zip(quorum_membership_index_for_each_non_signer.into_iter())
        .map(|(pk, quorum_membership_index)| {
            let pk_hash = convert::point_to_hash(pk)?;

            let quorum_membership = quorum_membership_history_by_signer
                .get(&pk_hash)
                .ok_or(MissingSignerEntry)?
                .try_get_at(quorum_membership_index)?
                .try_get_against(reference_block)?;

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
    let quorums = signed_quorums
        .into_iter()
        .zip(apk_for_each_quorum.into_iter())
        .zip(total_stake_index_for_each_quorum.into_iter())
        .zip(stake_index_for_each_quorum_and_required_non_signer.into_iter())
        .map(
            |(
                ((signed_quorum, apk), total_stake_index),
                stake_index_for_each_required_non_signer,
            )| {
                let total_stake = total_stake_history_by_quorum
                    .get(&signed_quorum)
                    .ok_or(MissingQuorumEntry)?
                    .try_get_at(total_stake_index)?
                    .try_get_against(reference_block)?;

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
                            .try_get_against(reference_block)
                    })
                    .sum::<Result<Stake, _>>()?;

                let signed_stake = total_stake.checked_sub(unsigned_stake).ok_or(Underflow)?;

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

    if signature::verification::verify(msg_hash, signers_apk, apk_g2, sigma) == false {
        return Err(SignatureVerificationFailed);
    }

    let _signatory_record_hash = hash::signature_record(reference_block, &non_signers);

    Ok(())
}

#[cfg(test)]
mod tests {
    use alloc::vec;
    use ark_bn254::{Fr, G1Affine, G1Projective, G2Projective};
    use ark_ec::{CurveGroup, PrimeGroup};
    use hashbrown::HashMap;

    use crate::{
        bitmap::Bitmap,
        cert, convert,
        error::CertVerificationError::*,
        hash::BeHash,
        types::{
            Cert, Chain, NonSignerStakesAndSignature,
            history::{History, Update},
        },
    };

    #[test]
    fn success() {
        let (cert, chain) = success_inputs();

        let result = cert::verify(cert, chain);
        assert_eq!(result, Ok(()));
    }

    #[test]
    fn reference_block_past_current_block() {
        let (mut cert, mut chain) = success_inputs();
        cert.reference_block = 43;
        chain.current_block = 42;

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, ReferenceBlockDoesNotPrecedeCurrentBlock);
    }

    #[test]
    fn reference_block_at_current_block() {
        let (mut cert, mut chain) = success_inputs();
        cert.reference_block = 42;
        chain.current_block = 42;

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, ReferenceBlockDoesNotPrecedeCurrentBlock);
    }

    #[test]
    fn empty_non_signer_vecs() {
        let (mut cert, chain) = success_inputs();
        cert.params.pk_for_each_non_signer.clear();

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, EmptyVec);
    }

    #[test]
    fn empty_quorum_vecs() {
        let (mut cert, chain) = success_inputs();
        cert.signed_quorums.clear();

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, EmptyVec);
    }

    #[test]
    fn reject_staleness() {
        let (cert, mut chain) = success_inputs();
        chain.reject_staleness = true;
        chain.last_updated_at_block_by_quorum.insert(0, 41);

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, StaleQuorum);
    }

    #[test]
    fn cert_apk_not_equal_chain_apk() {
        let (mut cert, chain) = success_inputs();
        cert.params.apk_for_each_quorum[0] = G1Affine::identity();

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, PointAtInfinity);
    }

    #[test]
    fn non_signers_point_at_infinity() {
        let (mut cert, chain) = success_inputs();
        cert.params.pk_for_each_non_signer[0] = G1Affine::identity();

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, PointAtInfinity);
    }

    #[test]
    fn quorum_membership_history_missing_signer_entry() {
        let (cert, mut chain) = success_inputs();
        chain.quorum_membership_history_by_signer.clear();

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, MissingSignerEntry);
    }

    #[test]
    fn quorum_membership_history_missing_history_entry() {
        let (mut cert, chain) = success_inputs();
        cert.params.quorum_membership_index_for_each_non_signer[0] = 42;

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, MissingHistoryEntry);
    }

    #[test]
    fn quorum_membership_history_reference_block_not_in_interval() {
        let (cert, mut chain) = success_inputs();
        chain
            .quorum_membership_history_by_signer
            .iter_mut()
            .for_each(|(_, v)| {
                v.0.insert(0, Update::new(141, 143, Default::default()).unwrap());
            });

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, ElementNotInInterval);
    }

    #[test]
    fn non_signers_not_strictly_sorted_by_hash() {
        let (mut cert, chain) = success_inputs();
        cert.params.pk_for_each_non_signer.reverse();

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, NotStrictlySortedByHash);
    }

    #[test]
    fn total_stake_history_missing_quorum_entry() {
        let (cert, mut chain) = success_inputs();
        chain.total_stake_history_by_quorum.clear();

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, MissingQuorumEntry);
    }

    #[test]
    fn total_stake_history_missing_history_entry() {
        let (cert, mut chain) = success_inputs();
        chain
            .total_stake_history_by_quorum
            .insert(0, Default::default());

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, MissingHistoryEntry);
    }

    #[test]
    fn total_stake_history_reference_block_not_in_interval() {
        let (cert, mut chain) = success_inputs();
        chain
            .total_stake_history_by_quorum
            .iter_mut()
            .for_each(|(_, v)| {
                v.0.insert(0, Update::new(141, 143, Default::default()).unwrap());
            });

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, ElementNotInInterval);
    }

    #[test]
    fn stake_history_missing_signer_entry() {
        let (cert, mut chain) = success_inputs();
        chain.stake_history_by_signer_and_quorum.clear();

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, MissingSignerEntry);
    }

    #[test]
    fn stake_history_missing_quorum_entry() {
        let (cert, mut chain) = success_inputs();
        chain
            .stake_history_by_signer_and_quorum
            .iter_mut()
            .for_each(|(_, stake_history_by_quorum)| {
                stake_history_by_quorum.clear();
            });

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, MissingQuorumEntry);
    }

    #[test]
    fn stake_history_missing_history_entry() {
        let (cert, mut chain) = success_inputs();
        chain
            .stake_history_by_signer_and_quorum
            .iter_mut()
            .for_each(|(_, stake_history_by_quorum)| {
                stake_history_by_quorum.insert(0, Default::default());
            });

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, MissingHistoryEntry);
    }

    #[test]
    fn stake_history_reference_block_not_in_interval() {
        let (cert, mut chain) = success_inputs();
        chain
            .stake_history_by_signer_and_quorum
            .iter_mut()
            .for_each(|(_, stake_history_by_quorum)| {
                stake_history_by_quorum.iter_mut().for_each(|(_, v)| {
                    v.0.insert(0, Update::new(141, 143, Default::default()).unwrap());
                })
            });

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, ElementNotInInterval);
    }

    #[test]
    fn stake_underflow() {
        let (cert, mut chain) = success_inputs();

        chain
            .total_stake_history_by_quorum
            .iter_mut()
            .for_each(|(_, v)| {
                v.0.insert(0, Update::new(41, 43, 29).unwrap());
            });

        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, Underflow);
    }

    #[test]
    fn aggregation_failure() {
        let (cert, mut chain) = success_inputs();
        chain.initialized_quorums_count = 1;
        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, BitIndexNotLessThanUpperBound);
    }

    #[test]
    fn signature_verification_failure() {
        let (mut cert, chain) = success_inputs();
        cert.msg_hash = BeHash::default();
        let err = cert::verify(cert, chain).unwrap_err();
        assert_eq!(err, SignatureVerificationFailed);
    }

    fn success_inputs() -> (Cert, Chain) {
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

        let msg_hash = [42u8; 32];
        let msg_point = convert::hash_to_point(msg_hash);

        let sig_at_quorum_2_by_signer_3 = (msg_point * signer3_sk).into_affine();
        let sig_at_quorum_0_by_signer_4 = (msg_point * signer4_sk).into_affine();
        let sigma = (sig_at_quorum_2_by_signer_3 + sig_at_quorum_0_by_signer_4).into_affine();

        let apk_for_each_quorum = [
            (non_signer0_g1_pk + non_signer2_g1_pk + signer4_g1_pk).into_affine(),
            (non_signer0_g1_pk + non_signer1_g1_pk + non_signer2_g1_pk + signer3_g1_pk)
                .into_affine(),
        ];
        let params = NonSignerStakesAndSignature {
            apk_for_each_quorum: apk_for_each_quorum.to_vec(),
            apk_index_for_each_quorum: vec![0, 0],
            total_stake_index_for_each_quorum: vec![0, 0],
            stake_index_for_each_quorum_and_required_non_signer: vec![vec![0, 0, 0], vec![0, 0, 0]],
            pk_for_each_non_signer: vec![non_signer0_g1_pk, non_signer1_g1_pk, non_signer2_g1_pk],
            quorum_membership_index_for_each_non_signer: vec![0, 0, 0],
            apk_g2,
            sigma,
        };

        let signed_quorums = [0, 2];

        let cert = Cert {
            msg_hash,
            reference_block: 42,

            // quorum 1 had no signatures
            // quorums 0 and 2 had at least one signature (exactly one in this example)
            signed_quorums: signed_quorums.to_vec(),

            params,
        };

        let non_signer0_pk_hash = convert::point_to_hash(non_signer0_g1_pk).unwrap();
        let non_signer1_pk_hash = convert::point_to_hash(non_signer1_g1_pk).unwrap();
        let non_signer2_pk_hash = convert::point_to_hash(non_signer2_g1_pk).unwrap();
        let signer3_pk_hash = convert::point_to_hash(signer3_g1_pk).unwrap();
        let signer4_pk_hash = convert::point_to_hash(signer4_g1_pk).unwrap();
        let optional_non_signer5_pk_hash =
            convert::point_to_hash(optional_non_signer5_g1_pk).unwrap();

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
                let stake_history_by_quorum = signed_quorums
                    .into_iter()
                    .map(|quorum| {
                        let update = Update::new(41, 43, 10).unwrap();
                        let history = HashMap::from([(0, update)]);
                        (quorum, History(history))
                    })
                    .collect();
                (pk_hash, stake_history_by_quorum)
            })
            .collect::<HashMap<BeHash, _>>();

        let total_stake_history_by_quorum = signed_quorums
            .into_iter()
            .map(|quorum| {
                let update = Update::new(41, 43, 30).unwrap();
                let history = HashMap::from([(0, update)]);
                (quorum, History(history))
            })
            .collect();

        let apk_trunc_hash_history_by_quorum = signed_quorums
            .into_iter()
            .zip(apk_for_each_quorum)
            .map(|(quorum, apk)| {
                let apk_hash = convert::point_to_hash(apk).unwrap();
                let apk_trunch_hash: [u8; 24] = apk_hash[..24].try_into().unwrap();
                let update = Update::new(41, 43, apk_trunch_hash).unwrap();
                let history = HashMap::from([(0, update)]);
                (quorum, History(history))
            })
            .collect();

        let last_updated_at_block_by_quorum = signed_quorums
            .into_iter()
            .map(|quorum| (quorum, 42))
            .collect();

        let chain = Chain {
            initialized_quorums_count: u8::MAX,
            current_block: 43,
            reject_staleness: false,
            min_withdrawal_delay_blocks: 1,
            quorum_membership_history_by_signer,
            stake_history_by_signer_and_quorum,
            total_stake_history_by_quorum,
            apk_trunc_hash_history_by_quorum,
            last_updated_at_block_by_quorum,
        };

        (cert, chain)
    }
}
