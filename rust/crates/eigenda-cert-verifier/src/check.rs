use alloc::vec::Vec;
use ark_bn254::G1Affine;
use hashbrown::HashMap;

use crate::{
    convert,
    error::CertVerificationError::{self, *},
    types::{NonSigner, history::History},
};

pub fn non_zero_equal_lengths(lengths: &[usize]) -> Result<(), CertVerificationError> {
    match lengths.first() {
        None | Some(0) => Err(EmptyVec),
        Some(first) => lengths
            .iter()
            .all(|length| length == first)
            .then_some(())
            .ok_or(UnequalLengths),
    }
}

pub fn non_signers_strictly_sorted_by_hash(
    non_signers: &[NonSigner],
) -> Result<(), CertVerificationError> {
    non_signers
        // if `non_signers.len() < 2` windows yields no elements
        .windows(2)
        .all(|window| matches!(window, [prev, curr] if prev.pk_hash < curr.pk_hash))
        .then_some(())
        .ok_or(NotStrictlySortedByHash)
}

pub fn quorums_last_updated_after_most_recent_stale_block(
    signed_quorums: &[u8],
    reference_block: u32,
    last_updated_at_block_by_quorum: HashMap<u8, u32>,
    window: u32,
) -> Result<(), CertVerificationError> {
    signed_quorums.iter().try_for_each(|signed_quorum| {
        let last_updated_at_block = *last_updated_at_block_by_quorum
            .get(signed_quorum)
            .ok_or(MissingQuorumEntry)?;

        let most_recent_stale_block = reference_block.checked_sub(window).ok_or(Underflow)?;
        let is_recent = last_updated_at_block > most_recent_stale_block;
        is_recent.then_some(()).ok_or(StaleQuorum)
    })
}

pub fn cert_apks_equal_chain_apks(
    signed_quorums: &[u8],
    reference_block: u32,
    apk_for_each_quorum: &[G1Affine],
    apk_index_for_each_quorum: Vec<u32>,
    apk_trunc_hash_history_by_quorum: HashMap<u8, History<[u8; 24]>>,
) -> Result<(), CertVerificationError> {
    signed_quorums
        .iter()
        .zip(apk_for_each_quorum.into_iter())
        .zip(apk_index_for_each_quorum.into_iter())
        .try_for_each(|((signed_quorum, &cert_apk), apk_index)| {
            let cert_apk_hash = convert::point_to_hash(cert_apk)?;
            let cert_apk_trunc_hash = &cert_apk_hash[..24];

            let chain_apk_trunc_hash = apk_trunc_hash_history_by_quorum
                .get(signed_quorum)
                .ok_or(MissingQuorumEntry)?
                .try_get_at(apk_index)?
                .try_get_against(reference_block)?;

            (cert_apk_trunc_hash == chain_apk_trunc_hash)
                .then_some(())
                .ok_or(CertApkDoesNotEqualChainApk)
        })
}

#[cfg(test)]
mod test_non_zero_equal_lengths {
    use crate::{check, error::CertVerificationError::*};

    #[test]
    fn non_zero_equal_lengths_success() {
        assert!(check::non_zero_equal_lengths(&[42, 42, 42, 42]).is_ok());
    }

    #[test]
    fn different_lengths_where_none_is_zero() {
        let err = check::non_zero_equal_lengths(&[42, 43, 44, 45]).unwrap_err();
        assert_eq!(err, UnequalLengths);
    }

    #[test]
    fn first_length_zero_but_otherwise_equal_lengths() {
        let err = check::non_zero_equal_lengths(&[0, 42, 42, 42]).unwrap_err();
        assert_eq!(err, EmptyVec);
    }

    #[test]
    fn all_lenghts_zero() {
        let err = check::non_zero_equal_lengths(&[0, 0, 0, 0]).unwrap_err();
        assert_eq!(err, EmptyVec);
    }

    #[test]
    fn some_length_zero_but_otherwise_equal_lengths() {
        let err = check::non_zero_equal_lengths(&[42, 42, 0, 42]).unwrap_err();
        assert_eq!(err, UnequalLengths);
    }
}

#[cfg(test)]
mod test_non_signers_strictly_sorted_by_hash {
    use crate::{check, error::CertVerificationError::*, types::NonSigner};

    #[test]
    fn strictly_sorted_by_hash() {
        let non_signers = &[[42u8; 32], [43u8; 32], [44u8; 32]].map(|pk_hash| NonSigner {
            pk_hash,
            ..Default::default()
        });
        let result = check::non_signers_strictly_sorted_by_hash(non_signers);
        assert!(result.is_ok());
    }

    #[test]
    fn sorted_by_hash_but_not_strictly() {
        let non_signers = &[[42u8; 32], [43u8; 32], [43u8; 32]].map(|pk_hash| NonSigner {
            pk_hash,
            ..Default::default()
        });
        let err = check::non_signers_strictly_sorted_by_hash(non_signers).unwrap_err();
        assert_eq!(err, NotStrictlySortedByHash);
    }

    #[test]
    fn not_sorted_by_hash() {
        let non_signers = &[[44u8; 32], [43u8; 32], [42u8; 32]].map(|pk_hash| NonSigner {
            pk_hash,
            ..Default::default()
        });
        let err = check::non_signers_strictly_sorted_by_hash(non_signers).unwrap_err();
        assert_eq!(err, NotStrictlySortedByHash);
    }

    #[test]
    fn empty_vec() {
        let result = check::non_signers_strictly_sorted_by_hash(&[]);
        assert!(result.is_ok());
    }

    #[test]
    fn just_one_signer() {
        let non_signers = &[[42u8; 32]].map(|pk_hash| NonSigner {
            pk_hash,
            ..Default::default()
        });
        let result = check::non_signers_strictly_sorted_by_hash(non_signers);
        assert!(result.is_ok());
    }
}

#[cfg(test)]
mod test_quorums_last_updated_after_most_recent_stale_block {
    use crate::{check, error::CertVerificationError::*};

    #[test]
    fn quorum_last_updated_before_most_recent_stale_block() {
        let reference_block = 42;
        let window = 1;
        let most_recent_stale_block = reference_block - window;

        let signed_quorums = [0];
        let last_updated_at_block_by_quorum = signed_quorums
            .into_iter()
            .map(|signed_quorum| (signed_quorum, most_recent_stale_block - 1))
            .collect();

        let err = check::quorums_last_updated_after_most_recent_stale_block(
            &signed_quorums,
            reference_block,
            last_updated_at_block_by_quorum,
            window,
        )
        .unwrap_err();

        assert_eq!(err, StaleQuorum);
    }

    #[test]
    fn quorum_last_updated_at_most_recent_stale_block() {
        let reference_block = 42;
        let window = 1;
        let most_recent_stale_block = reference_block - window;

        let signed_quorums = [0];
        let last_updated_at_block_by_quorum = signed_quorums
            .into_iter()
            .map(|signed_quorum| (signed_quorum, most_recent_stale_block))
            .collect();

        let err = check::quorums_last_updated_after_most_recent_stale_block(
            &signed_quorums,
            reference_block,
            last_updated_at_block_by_quorum,
            window,
        )
        .unwrap_err();

        assert_eq!(err, StaleQuorum);
    }

    #[test]
    fn quorums_last_updated_after_most_recent_stale_block() {
        let reference_block = 42;
        let window = 1;
        let most_recent_stale_block = reference_block - window;

        let signed_quorums = [0];
        let last_updated_at_block_by_quorum = signed_quorums
            .into_iter()
            .map(|signed_quorum| (signed_quorum, most_recent_stale_block + 1))
            .collect();

        let result = check::quorums_last_updated_after_most_recent_stale_block(
            &signed_quorums,
            reference_block,
            last_updated_at_block_by_quorum,
            window,
        );

        assert!(result.is_ok());
    }

    #[test]
    fn missing_quorum_entry() {
        let reference_block = 42;
        let window = 1;

        let signed_quorums = [0];
        let err = check::quorums_last_updated_after_most_recent_stale_block(
            &signed_quorums,
            reference_block,
            Default::default(),
            window,
        )
        .unwrap_err();

        assert_eq!(err, MissingQuorumEntry);
    }

    #[test]
    fn underflow() {
        let reference_block = 42;
        let window = 43;
        let signed_quorums = [0];
        let last_updated_at_block_by_quorum = signed_quorums
            .into_iter()
            .map(|signed_quorum| (signed_quorum, Default::default()))
            .collect();

        let err = check::quorums_last_updated_after_most_recent_stale_block(
            &signed_quorums,
            reference_block,
            last_updated_at_block_by_quorum,
            window,
        )
        .unwrap_err();

        assert_eq!(err, Underflow);
    }
}

#[cfg(test)]
mod test_cert_apks_equal_chain_apks {
    use alloc::vec;
    use ark_bn254::{Fr, G1Affine, G1Projective};
    use ark_ec::{CurveGroup, PrimeGroup};
    use hashbrown::HashMap;

    use crate::{
        check, convert,
        error::CertVerificationError::*,
        hash::TruncatedBeHash,
        types::{
            BlockNumber,
            history::{History, Update},
        },
    };

    #[test]
    fn cert_apk_equal_chain_apk() {
        let apk = (G1Projective::generator() * Fr::from(42)).into_affine();
        let apk_hash = convert::point_to_hash(apk).unwrap();
        let apk_trunc_hash: TruncatedBeHash = apk_hash[..24].try_into().unwrap();

        let signed_quorums = [0];
        let reference_block = 42;
        let apk_for_each_quorum = [apk];
        let apk_index_for_each_quorum = vec![0];

        let update = Update::new(42, 43, apk_trunc_hash.clone()).unwrap();
        let history = HashMap::from([(0, update)]);
        let apk_trunc_hash_history = History(history);
        let apk_trunc_hash_history_by_quorum = HashMap::from([(0, apk_trunc_hash_history)]);

        let result = check::cert_apks_equal_chain_apks(
            &signed_quorums,
            reference_block,
            &apk_for_each_quorum,
            apk_index_for_each_quorum,
            apk_trunc_hash_history_by_quorum,
        );

        assert!(result.is_ok());
    }

    #[test]
    fn cert_apk_does_not_equal_chain_apk() {
        let cert_apk = (G1Projective::generator() * Fr::from(42)).into_affine();
        let chain_apk = (G1Projective::generator() * Fr::from(43)).into_affine();
        let chain_apk_hash = convert::point_to_hash(chain_apk).unwrap();
        let chain_apk_trunc_hash: TruncatedBeHash = chain_apk_hash[..24].try_into().unwrap();

        let signed_quorums = [0];
        let reference_block = 42;
        let apk_for_each_quorum = [cert_apk];
        let apk_index_for_each_quorum = vec![0];

        let update = Update::new(42, 43, chain_apk_trunc_hash.clone()).unwrap();
        let history = HashMap::from([(0, update)]);
        let apk_trunc_hash_history = History(history);
        let apk_trunc_hash_history_by_quorum = HashMap::from([(0, apk_trunc_hash_history)]);

        let err = check::cert_apks_equal_chain_apks(
            &signed_quorums,
            reference_block,
            &apk_for_each_quorum,
            apk_index_for_each_quorum,
            apk_trunc_hash_history_by_quorum,
        )
        .unwrap_err();

        assert_eq!(err, CertApkDoesNotEqualChainApk);
    }

    #[test]
    fn missing_quorum_entry() {
        let apk = (G1Projective::generator() * Fr::from(42)).into_affine();

        let signed_quorums = [0];
        let reference_block = 42;
        let apk_for_each_quorum = [apk];

        let apk_index_for_each_quorum = vec![0];

        let err = check::cert_apks_equal_chain_apks(
            &signed_quorums,
            reference_block,
            &apk_for_each_quorum,
            apk_index_for_each_quorum,
            Default::default(),
        )
        .unwrap_err();

        assert_eq!(err, MissingQuorumEntry);
    }

    #[test]
    fn point_at_infinity() {
        let apk = G1Affine::identity();

        let signed_quorums = [0];
        let reference_block = 42;
        let apk_for_each_quorum = [apk];

        let apk_index_for_each_quorum = vec![0];

        let err = check::cert_apks_equal_chain_apks(
            &signed_quorums,
            reference_block,
            &apk_for_each_quorum,
            apk_index_for_each_quorum,
            Default::default(),
        )
        .unwrap_err();

        assert_eq!(err, PointAtInfinity);
    }

    #[test]
    fn missing_history_entry() {
        let apk = (G1Projective::generator() * Fr::from(42)).into_affine();

        let signed_quorums = [0];
        let reference_block = 42;
        let apk_for_each_quorum = [apk];
        let apk_index_for_each_quorum = vec![0];

        let apk_trunc_hash_history = History(Default::default());
        let apk_trunc_hash_history_by_quorum = HashMap::from([(0, apk_trunc_hash_history)]);

        let err = check::cert_apks_equal_chain_apks(
            &signed_quorums,
            reference_block,
            &apk_for_each_quorum,
            apk_index_for_each_quorum,
            apk_trunc_hash_history_by_quorum,
        )
        .unwrap_err();

        assert_eq!(err, MissingHistoryEntry);
    }

    #[test]
    fn stale_reference_block() {
        let apk = (G1Projective::generator() * Fr::from(42)).into_affine();

        let signed_quorums = [0];
        const STALE_REFERENCE_BLOCK: BlockNumber = 41;
        let apk_for_each_quorum = [apk];
        let apk_index_for_each_quorum = vec![0];

        let update = Update::new(42, 43, Default::default()).unwrap();
        let history = HashMap::from([(0, update)]);
        let apk_trunc_hash_history = History(history);
        let apk_trunc_hash_history_by_quorum = HashMap::from([(0, apk_trunc_hash_history)]);

        let err = check::cert_apks_equal_chain_apks(
            &signed_quorums,
            STALE_REFERENCE_BLOCK,
            &apk_for_each_quorum,
            apk_index_for_each_quorum,
            apk_trunc_hash_history_by_quorum,
        )
        .unwrap_err();

        assert_eq!(err, ElementNotInInterval);
    }
}
