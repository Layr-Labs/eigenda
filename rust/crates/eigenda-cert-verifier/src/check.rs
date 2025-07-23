use alloc::vec::Vec;
use ark_bn254::G1Affine;
use hashbrown::HashMap;

use crate::{
    convert,
    error::CertVerificationError::{self, *},
    types::{
        Address, NonSigner, RelayInfo, RelayKey, SecurityThresholds, Version, VersionedBlobParams,
        history::History,
    },
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

pub fn relay_keys_are_set(
    relay_keys: &[RelayKey],
    relay_key_to_relay_info: &HashMap<RelayKey, RelayInfo>,
) -> Result<(), CertVerificationError> {
    relay_keys.iter().try_for_each(|relay_key| {
        let relay_info = relay_key_to_relay_info
            .get(relay_key)
            .ok_or(MissingRelayKeyEntry)?;
        match relay_info.address == Address::default() {
            true => Err(RelayKeyNotSet),
            false => Ok(()),
        }
    })
}

pub fn security_assumptions_are_met(
    version: Version,
    version_to_versioned_blob_params: &HashMap<Version, VersionedBlobParams>,
    security_thresholds: &SecurityThresholds,
) -> Result<(), CertVerificationError> {
    let SecurityThresholds {
        confirmation_threshold,
        adversary_threshold,
    } = security_thresholds;

    let VersionedBlobParams {
        max_num_operators,
        num_chunks,
        coding_rate,
    } = version_to_versioned_blob_params
        .get(&version)
        .ok_or(MissingVersionEntry)?;

    if (confirmation_threshold > adversary_threshold) == false {
        return Err(ConfirmationThresholdNotGreaterThanAdversaryThreshold);
    }

    let confirmation_threshold = *confirmation_threshold as u64;
    let adversary_threshold = *adversary_threshold as u64;
    let coding_rate = *coding_rate as u64;
    let num_chunks = *num_chunks as u64;
    let max_num_operators = *max_num_operators as u64;

    // safety: cannot underflow due to the `confirmation_threshold > adversary_threshold` check
    let gamma = confirmation_threshold - adversary_threshold;

    let denominator = gamma * coding_rate;

    // safety: cannot be 0 due to the `confirmation_threshold > adversary_threshold` check
    let inverse = 1_000_000 / denominator;

    let n = 10_000u64.checked_sub(inverse).ok_or(Underflow)? * num_chunks;

    // Overflow analysis:
    //
    // confirmation_threshold ∈ [0, 255]
    // adversary_threshold ∈ [0, 255]
    // gamma ∈ [1, 255] (not [0, 255] due to the `confirmation_threshold > adversary_threshold` check)
    // denominator ∈ [1*1, 255*255]
    // inverse ∈ [1_000_000 / (255*255), 1_000_000 / (1*1)]
    //     in the calculation of n that follows, inverse cannot exceed 10_000
    //     so inverse must instead ∈ [1_000_000 / (255*255), 1_000_000 / 100]
    //     which means gamma*coding_rate >= 100
    // Conclusion: underflow will happen whenever gamma*coding_rate < 100
    //
    // Another consideration: n * num_chunks ∈ [0, 10_000] * [0, 2^32]
    //     where the upper bound can overflow if represented as u32 hence the casts to u64
    //     same for max_num_operators * 10_000

    if n < max_num_operators * 10_000 {
        return Err(UnmetSecurityAssumptions);
    }

    Ok(())
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
    fn all_lengths_zero() {
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

#[cfg(test)]
mod test_relay_keys_are_set {
    use alloc::vec;
    use hashbrown::HashMap;

    use crate::{
        check,
        error::CertVerificationError::*,
        types::{Address, RelayInfo},
    };

    #[test]
    fn success_when_all_relay_keys_are_set() {
        let relay_keys = vec![0];

        let relay_key_to_relay_info = HashMap::from([(
            0,
            RelayInfo {
                address: [42u8; 20],
                url: Default::default(),
            },
        )]);

        let result = check::relay_keys_are_set(&relay_keys, &relay_key_to_relay_info);

        assert_eq!(result, Ok(()));
    }

    #[test]
    fn relay_keys_are_set_fails_with_missing_relay_key() {
        let relay_keys = vec![99]; // 99 not found on chain

        let relay_key_to_relay_info = HashMap::from([(
            42,
            RelayInfo {
                address: [42u8; 20],
                url: Default::default(),
            },
        )]);

        let err = check::relay_keys_are_set(&relay_keys, &relay_key_to_relay_info).unwrap_err();

        assert_eq!(err, MissingRelayKeyEntry);
    }

    #[test]
    fn relay_keys_are_set_fails_when_corresponding_address_is_not_set() {
        let relay_keys = vec![42];

        let relay_key_to_relay_info = HashMap::from([(
            42,
            RelayInfo {
                address: Address::default(),
                url: Default::default(),
            },
        )]);

        let err = check::relay_keys_are_set(&relay_keys, &relay_key_to_relay_info).unwrap_err();

        assert_eq!(err, RelayKeyNotSet);
    }
}

#[cfg(test)]
mod test_security_assumptions_are_met {
    use hashbrown::HashMap;

    use crate::{
        check,
        error::CertVerificationError::*,
        types::{SecurityThresholds, Version, VersionedBlobParams},
    };

    #[test]
    fn success_when_security_assumptions_are_met() {
        let (version, version_to_versioned_blob_params, security_thresholds) = success_inputs();

        let result = check::security_assumptions_are_met(
            version,
            &version_to_versioned_blob_params,
            &security_thresholds,
        );

        assert_eq!(result, Ok(()));
    }

    #[test]
    fn security_assumptions_are_met_fails_with_missing_version_entry() {
        let (_version, version_to_versioned_blob_params, security_thresholds) = success_inputs();

        let err = check::security_assumptions_are_met(
            u16::MAX,
            &version_to_versioned_blob_params,
            &security_thresholds,
        )
        .unwrap_err();

        assert_eq!(err, MissingVersionEntry);
    }

    #[test]
    fn security_assumptions_are_met_fails_when_confirmation_threshold_equals_adversary_threshold() {
        let (version, version_to_versioned_blob_params, mut security_thresholds) = success_inputs();

        security_thresholds.confirmation_threshold = security_thresholds.adversary_threshold;

        let err = check::security_assumptions_are_met(
            version,
            &version_to_versioned_blob_params,
            &security_thresholds,
        )
        .unwrap_err();

        assert_eq!(err, ConfirmationThresholdNotGreaterThanAdversaryThreshold);
    }

    #[test]
    fn security_assumptions_are_met_fails_when_confirmation_threshold_less_than_adversary_threshold()
     {
        let (version, version_to_versioned_blob_params, mut security_thresholds) = success_inputs();

        security_thresholds.confirmation_threshold = security_thresholds.adversary_threshold - 1;

        let err = check::security_assumptions_are_met(
            version,
            &version_to_versioned_blob_params,
            &security_thresholds,
        )
        .unwrap_err();

        assert_eq!(err, ConfirmationThresholdNotGreaterThanAdversaryThreshold);
    }

    #[test]
    fn security_assumptions_are_met_fails_with_underflow() {
        let (version, mut version_to_versioned_blob_params, mut security_thresholds) =
            success_inputs();

        // to trigger overflow (gamma * coding_rate) < 100
        // where gamma = confirmation_threshold - adversary_threshold
        security_thresholds.confirmation_threshold = 101;
        security_thresholds.adversary_threshold = 100;
        // gamma = 101 - 100 = 1
        let params = version_to_versioned_blob_params.get_mut(&version).unwrap();
        params.coding_rate = 99;

        let err = check::security_assumptions_are_met(
            version,
            &version_to_versioned_blob_params,
            &security_thresholds,
        )
        .unwrap_err();

        assert_eq!(err, Underflow);
    }

    #[test]
    fn security_assumptions_are_met_fails_with_unmet_security_assumptions() {
        let (version, version_to_versioned_blob_params, mut security_thresholds) = success_inputs();

        // from success_inputs:
        // gamma = confirmation_threshold - adversary_threshold = 101 - 1 = 100
        // since the success_inputs are at the limit
        // any disturbance will cause UnmetSecurityAssumptions so
        security_thresholds.adversary_threshold = 2; // instead of 1, resulting in gamma = 99

        let err = check::security_assumptions_are_met(
            version,
            &version_to_versioned_blob_params,
            &security_thresholds,
        )
        .unwrap_err();

        assert_eq!(err, UnmetSecurityAssumptions);
    }

    fn success_inputs() -> (
        Version,
        HashMap<Version, VersionedBlobParams>,
        SecurityThresholds,
    ) {
        let version = 42u16;
        let version_to_versioned_blob_params = HashMap::from([(
            version,
            VersionedBlobParams {
                max_num_operators: 99,
                num_chunks: 100,
                coding_rate: 100,
            },
        )]);
        let security_thresholds = SecurityThresholds {
            confirmation_threshold: 101,
            adversary_threshold: 1,
        };

        // gamma = confirmation_threshold - adversary_threshold = 101 - 1 = 100
        // inverse = 1_000_000 / (gamma * coding_rate) = 1_000_000 / (100 * 100) = 100
        // n = (10_000 - inverse) * num_chunks = (10_000 - 100) * 100 = 990_000
        // max_num_operators * 10_000 = 99 * 10_000 = 990_000
        // 990_000 >= 990_000

        (
            version,
            version_to_versioned_blob_params,
            security_thresholds,
        )
    }
}
