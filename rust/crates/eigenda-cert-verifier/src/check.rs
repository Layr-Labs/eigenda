use alloc::vec::Vec;
use alloy_primitives::{Address, B256, Bytes};
use eigenda_cert::G1Point;
use hashbrown::HashMap;

use crate::{
    bitmap::{Bitmap, bit_indices_to_bitmap},
    convert,
    error::CertVerificationError::{self, *},
    hash::{TruncatedB256, keccak_v256},
    types::{
        BlockNumber, NonSigner, Quorum, QuorumNumber, RelayKey, Version,
        history::History,
        solidity::{RelayInfo, SecurityThresholds, VersionedBlobParams},
    },
};

const THRESHOLD_DENOMINATOR: u128 = 100; // uint256 in sol

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
    signed_quorums: &[QuorumNumber],
    reference_block: BlockNumber,
    last_updated_at_block_by_quorum: HashMap<u8, BlockNumber>,
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

pub fn cert_apks_equal_storage_apks(
    signed_quorums: &[QuorumNumber],
    reference_block: BlockNumber,
    apk_for_each_quorum: &[G1Point],
    apk_index_for_each_quorum: Vec<BlockNumber>,
    apk_trunc_hash_history_by_quorum: HashMap<QuorumNumber, History<TruncatedB256>>,
) -> Result<(), CertVerificationError> {
    signed_quorums
        .iter()
        .zip(apk_for_each_quorum.iter())
        .zip(apk_index_for_each_quorum.into_iter())
        .try_for_each(|((signed_quorum, cert_apk), apk_index)| {
            let cert_apk_hash = convert::point_to_hash(cert_apk);
            let cert_apk_trunc_hash = &cert_apk_hash[..24];

            let storage_apk_trunc_hash = apk_trunc_hash_history_by_quorum
                .get(signed_quorum)
                .ok_or(MissingQuorumEntry)?
                .try_get_at(apk_index)?
                .try_get_against(reference_block)?;

            (cert_apk_trunc_hash == storage_apk_trunc_hash)
                .then_some(())
                .ok_or(CertApkDoesNotEqualStorageApk)
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

        (relay_info.relayAddress != Address::default())
            .then_some(())
            .ok_or(RelayKeyNotSet)
    })
}

pub fn security_assumptions_are_met(
    version: Version,
    version_to_versioned_blob_params: &HashMap<Version, VersionedBlobParams>,
    security_thresholds: &SecurityThresholds,
) -> Result<(), CertVerificationError> {
    let SecurityThresholds {
        confirmationThreshold,
        adversaryThreshold,
    } = security_thresholds;

    let VersionedBlobParams {
        maxNumOperators,
        numChunks,
        codingRate,
    } = version_to_versioned_blob_params
        .get(&version)
        .ok_or(MissingVersionEntry)?;

    if (confirmationThreshold > adversaryThreshold) == false {
        return Err(ConfirmationThresholdNotGreaterThanAdversaryThreshold);
    }

    let confirmation_threshold = *confirmationThreshold as u64;
    let adversary_threshold = *adversaryThreshold as u64;
    let coding_rate = *codingRate as u64;
    let num_chunks = *numChunks as u64;
    let max_num_operators = *maxNumOperators as u64;

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
    //     which means gamma*codingRate >= 100
    // Conclusion: underflow will happen whenever gamma*codingRate < 100
    //
    // Another consideration: n * numChunks ∈ [0, 10_000] * [0, 2^32]
    //     where the upper bound can overflow if represented as u32 hence the casts to u64
    //     same for maxNumOperators * 10_000

    if n < max_num_operators * 10_000 {
        return Err(UnmetSecurityAssumptions);
    }

    Ok(())
}

pub fn confirmed_quorums_contain_blob_quorums(
    confirmation_threshold: u8,
    quorums: &[Quorum],
    blob_quorums: &Bytes,
) -> Result<(), CertVerificationError> {
    let blob_quorums = bit_indices_to_bitmap(blob_quorums, None)?;

    let mut confirmed_quorums = Bitmap::default();

    quorums.iter().try_for_each(|quorum| {
        let Quorum {
            number,
            total_stake,
            signed_stake,
            ..
        } = *quorum;

        let left = signed_stake
            .checked_mul(THRESHOLD_DENOMINATOR)
            .ok_or(Overflow)?;

        let right = total_stake
            .checked_mul(confirmation_threshold as u128)
            .ok_or(Overflow)?;

        confirmed_quorums.set(number as usize, left >= right);

        Ok(())
    })?;

    contains(confirmed_quorums, blob_quorums)
        .then_some(())
        .ok_or(ConfirmedQuorumsDoNotContainBlobQuorums)
}

pub fn blob_quorums_contain_required_quorums(
    blob_quorums: &Bytes,
    required_quorums: &Bytes,
) -> Result<(), CertVerificationError> {
    let required_quorums = bit_indices_to_bitmap(required_quorums, None)?;
    let blob_quorums = bit_indices_to_bitmap(blob_quorums, None)?;
    contains(blob_quorums, required_quorums)
        .then_some(())
        .ok_or(BlobQuorumsDoNotContainRequiredQuorums)
}

/// Returns true if `container` contains all bits set in `contained`
#[inline]
fn contains(container: Bitmap, contained: Bitmap) -> bool {
    container & contained == contained
}

pub fn leaf_node_belongs_to_merkle_tree(
    leaf_node: B256,
    expected_root: B256,
    proof: Bytes,
    sibling_path: u32,
) -> Result<(), CertVerificationError> {
    if proof.len() % 32 != 0 {
        return Err(MerkleProofLengthNotMultipleOf32Bytes);
    }

    let sibling_path = Bitmap::new([sibling_path as u64, 0, 0, 0]);

    let proof_depth = proof.len() / 32;
    if sibling_path.len() < proof_depth {
        return Err(MerkleProofPathTooShort);
    }

    let mut current_node = leaf_node;
    for (i, sibling_node) in proof.chunks(32).enumerate() {
        // safety: the above `proof.len() % 32 != 0` guarantees proof is a multiple of 32
        let sibling_node = sibling_node.try_into().unwrap();
        let is_sibling_node_on_the_left = sibling_path[i];
        let (left_node, right_node) = match is_sibling_node_on_the_left {
            true => (sibling_node, current_node),
            false => (current_node, sibling_node),
        };
        let parent_node = keccak_v256([left_node, right_node].into_iter());
        current_node = parent_node;
    }

    let actual_root = current_node;
    (actual_root == expected_root)
        .then_some(())
        .ok_or(LeafNodeDoesNotBelongToMerkleTree)
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
            pk_hash: pk_hash.into(),
            ..Default::default()
        });
        let result = check::non_signers_strictly_sorted_by_hash(non_signers);
        assert_eq!(result, Ok(()));
    }

    #[test]
    fn sorted_by_hash_but_not_strictly() {
        let non_signers = &[[42u8; 32], [43u8; 32], [43u8; 32]].map(|pk_hash| NonSigner {
            pk_hash: pk_hash.into(),
            ..Default::default()
        });
        let err = check::non_signers_strictly_sorted_by_hash(non_signers).unwrap_err();
        assert_eq!(err, NotStrictlySortedByHash);
    }

    #[test]
    fn not_sorted_by_hash() {
        let non_signers = &[[44u8; 32], [43u8; 32], [42u8; 32]].map(|pk_hash| NonSigner {
            pk_hash: pk_hash.into(),
            ..Default::default()
        });
        let err = check::non_signers_strictly_sorted_by_hash(non_signers).unwrap_err();
        assert_eq!(err, NotStrictlySortedByHash);
    }

    #[test]
    fn empty_vec() {
        let result = check::non_signers_strictly_sorted_by_hash(&[]);
        assert_eq!(result, Ok(()));
    }

    #[test]
    fn just_one_signer() {
        let non_signers = &[[42u8; 32]].map(|pk_hash| NonSigner {
            pk_hash: pk_hash.into(),
            ..Default::default()
        });
        let result = check::non_signers_strictly_sorted_by_hash(non_signers);
        assert_eq!(result, Ok(()));
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

        assert_eq!(result, Ok(()));
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
mod test_cert_apks_equal_storage_apks {
    use alloc::vec;
    use ark_bn254::{Fr, G1Projective};
    use ark_ec::{CurveGroup, PrimeGroup};
    use hashbrown::HashMap;

    use crate::{
        check, convert,
        error::CertVerificationError::*,
        hash::TruncatedB256,
        types::{
            BlockNumber,
            conversions::IntoExt,
            history::{History, Update},
        },
    };

    #[test]
    fn cert_apk_equal_storage_apk() {
        let apk = (G1Projective::generator() * Fr::from(42)).into_affine();
        let apk_hash = convert::point_to_hash(&apk.into_ext());
        let apk_trunc_hash: TruncatedB256 = apk_hash[..24].try_into().unwrap();

        let signed_quorums = [0];
        let reference_block = 42;
        let apk_for_each_quorum = [apk.into_ext()];
        let apk_index_for_each_quorum = vec![0];

        let update = Update::new(42, 43, apk_trunc_hash.clone()).unwrap();
        let history = HashMap::from([(0, update)]);
        let apk_trunc_hash_history = History(history);
        let apk_trunc_hash_history_by_quorum = HashMap::from([(0, apk_trunc_hash_history)]);

        let result = check::cert_apks_equal_storage_apks(
            &signed_quorums,
            reference_block,
            &apk_for_each_quorum,
            apk_index_for_each_quorum,
            apk_trunc_hash_history_by_quorum,
        );

        assert_eq!(result, Ok(()));
    }

    #[test]
    fn cert_apk_does_not_equal_storage_apk() {
        let cert_apk = (G1Projective::generator() * Fr::from(42)).into_affine();
        let storage_apk = (G1Projective::generator() * Fr::from(43)).into_affine();
        let storage_apk_hash = convert::point_to_hash(&storage_apk.into_ext());
        let storage_apk_trunc_hash: TruncatedB256 = storage_apk_hash[..24].try_into().unwrap();

        let signed_quorums = [0];
        let reference_block = 42;
        let apk_for_each_quorum = [cert_apk.into_ext()];
        let apk_index_for_each_quorum = vec![0];

        let update = Update::new(42, 43, storage_apk_trunc_hash.clone()).unwrap();
        let history = HashMap::from([(0, update)]);
        let apk_trunc_hash_history = History(history);
        let apk_trunc_hash_history_by_quorum = HashMap::from([(0, apk_trunc_hash_history)]);

        let err = check::cert_apks_equal_storage_apks(
            &signed_quorums,
            reference_block,
            &apk_for_each_quorum,
            apk_index_for_each_quorum,
            apk_trunc_hash_history_by_quorum,
        )
        .unwrap_err();

        assert_eq!(err, CertApkDoesNotEqualStorageApk);
    }

    #[test]
    fn missing_quorum_entry() {
        let apk = (G1Projective::generator() * Fr::from(42)).into_affine();

        let signed_quorums = [0];
        let reference_block = 42;
        let apk_for_each_quorum = [apk.into_ext()];

        let apk_index_for_each_quorum = vec![0];

        let err = check::cert_apks_equal_storage_apks(
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
    fn missing_history_entry() {
        let apk = (G1Projective::generator() * Fr::from(42)).into_affine();

        let signed_quorums = [0];
        let reference_block = 42;
        let apk_for_each_quorum = [apk.into_ext()];
        let apk_index_for_each_quorum = vec![0];

        let apk_trunc_hash_history = History(Default::default());
        let apk_trunc_hash_history_by_quorum = HashMap::from([(0, apk_trunc_hash_history)]);

        let err = check::cert_apks_equal_storage_apks(
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
        let apk_for_each_quorum = [apk.into_ext()];
        let apk_index_for_each_quorum = vec![0];

        let update = Update::new(42, 43, Default::default()).unwrap();
        let history = HashMap::from([(0, update)]);
        let apk_trunc_hash_history = History(history);
        let apk_trunc_hash_history_by_quorum = HashMap::from([(0, apk_trunc_hash_history)]);

        let err = check::cert_apks_equal_storage_apks(
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
    use alloy_primitives::Address;
    use hashbrown::HashMap;

    use crate::{check, error::CertVerificationError::*, types::solidity::RelayInfo};

    #[test]
    fn success_when_all_relay_keys_are_set() {
        let relay_keys = vec![0];

        let relay_key_to_relay_info = HashMap::from([(
            0,
            RelayInfo {
                relayAddress: [42u8; 20].into(),
                relayURL: Default::default(),
            },
        )]);

        let result = check::relay_keys_are_set(&relay_keys, &relay_key_to_relay_info);

        assert_eq!(result, Ok(()));
    }

    #[test]
    fn relay_keys_are_set_fails_with_missing_relay_key() {
        let relay_keys = vec![99]; // 99 not found on storage

        let relay_key_to_relay_info = HashMap::from([(
            42,
            RelayInfo {
                relayAddress: [42u8; 20].into(),
                relayURL: Default::default(),
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
                relayAddress: Address::default().into(),
                relayURL: Default::default(),
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
        types::{
            Version,
            solidity::{SecurityThresholds, VersionedBlobParams},
        },
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
            Version::MAX,
            &version_to_versioned_blob_params,
            &security_thresholds,
        )
        .unwrap_err();

        assert_eq!(err, MissingVersionEntry);
    }

    #[test]
    fn security_assumptions_are_met_fails_when_confirmation_threshold_equals_adversary_threshold() {
        let (version, version_to_versioned_blob_params, mut security_thresholds) = success_inputs();

        security_thresholds.confirmationThreshold = security_thresholds.adversaryThreshold;

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

        security_thresholds.confirmationThreshold = security_thresholds.adversaryThreshold - 1;

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

        // to trigger overflow (gamma * codingRate) < 100
        // where gamma = confirmation_threshold - adversary_threshold
        security_thresholds.confirmationThreshold = 101;
        security_thresholds.adversaryThreshold = 100;
        // gamma = 101 - 100 = 1
        let params = version_to_versioned_blob_params.get_mut(&version).unwrap();
        params.codingRate = 99;

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
        security_thresholds.adversaryThreshold = 2; // instead of 1, resulting in gamma = 99

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
                maxNumOperators: 99,
                numChunks: 100,
                codingRate: 100,
            },
        )]);
        let security_thresholds = SecurityThresholds {
            confirmationThreshold: 101,
            adversaryThreshold: 1,
        };

        // gamma = confirmation_threshold - adversary_threshold = 101 - 1 = 100
        // inverse = 1_000_000 / (gamma * codingRate) = 1_000_000 / (100 * 100) = 100
        // n = (10_000 - inverse) * numChunks = (10_000 - 100) * 100 = 990_000
        // maxNumOperators * 10_000 = 99 * 10_000 = 990_000
        // 990_000 >= 990_000

        (
            version,
            version_to_versioned_blob_params,
            security_thresholds,
        )
    }
}

#[cfg(test)]
mod test_confirmed_quorums_contains_blob_quorums {
    use crate::{check, error::CertVerificationError::*, types::Quorum};
    use ark_bn254::G1Affine;

    #[test]
    fn success_when_confirmed_quorums_contain_blob_quorums() {
        let confirmation_threshold = 100;

        // in this example:
        //     quorum is confirmed if signed_stake * 100 > total_stake * 100
        //     quorum is confirmed if signed_stake * THRESHOLD_DENOMINATOR >= total_skate * confirmation_threshold
        let quorums = [
            Quorum {
                number: 0,
                total_stake: 42,
                signed_stake: 43,
                ..Default::default()
            },
            Quorum {
                number: 1,
                apk: G1Affine::default(),
                total_stake: 42,
                signed_stake: 42,
                ..Default::default()
            },
            Quorum {
                number: 2,
                total_stake: 42,
                signed_stake: 41,
                ..Default::default()
            },
        ];

        // in this example blob_quorums contains only confirmed quorums (0, 1 and 2)
        let blob_quorums = [0, 1].into();

        let result = check::confirmed_quorums_contain_blob_quorums(
            confirmation_threshold,
            &quorums,
            &blob_quorums,
        );

        assert_eq!(result, Ok(()));
    }

    #[test]
    fn confirmed_quorums_do_not_contain_blob_quorums() {
        let confirmation_threshold = 100;

        let quorums = [
            Quorum {
                number: 0,
                total_stake: 42,
                signed_stake: 43,
                ..Default::default()
            },
            Quorum {
                number: 1,
                apk: G1Affine::default(),
                total_stake: 42,
                signed_stake: 42,
                ..Default::default()
            },
            Quorum {
                number: 2,
                total_stake: 42,
                signed_stake: 41,
                ..Default::default()
            },
        ];

        // blob_quorums contains unconfirmed quorum 1
        let blob_quorums = [1, 2].into();

        let err = check::confirmed_quorums_contain_blob_quorums(
            confirmation_threshold,
            &quorums,
            &blob_quorums,
        )
        .unwrap_err();

        assert_eq!(err, ConfirmedQuorumsDoNotContainBlobQuorums);
    }

    #[test]
    fn overflow_in_signed_stake_multiplication() {
        let confirmation_threshold = 100;

        let quorums = [Quorum {
            number: 0,
            total_stake: 42,
            signed_stake: u128::MAX, // Will overflow when multiplied by THRESHOLD_DENOMINATOR
            ..Default::default()
        }];

        let blob_quorums = [0].into();

        let err = check::confirmed_quorums_contain_blob_quorums(
            confirmation_threshold,
            &quorums,
            &blob_quorums,
        )
        .unwrap_err();

        assert_eq!(err, Overflow);
    }

    #[test]
    fn overflow_in_total_stake_multiplication() {
        let confirmation_threshold = u8::MAX; // Will cause overflow when cast to u128 and multiplied

        let quorums = [Quorum {
            number: 0,
            total_stake: u128::MAX,
            signed_stake: 43,
            ..Default::default()
        }];

        let blob_quorums = [0].into();

        let err = check::confirmed_quorums_contain_blob_quorums(
            confirmation_threshold,
            &quorums,
            &blob_quorums,
        )
        .unwrap_err();

        assert_eq!(err, Overflow);
    }

    #[test]
    fn blob_quorums_bit_indices_not_sorted() {
        let confirmation_threshold = 100;
        let quorums = [Quorum {
            number: 0,
            total_stake: 42,
            signed_stake: 43,
            ..Default::default()
        }];

        let blob_quorums = [1, 0].into(); // Not sorted

        let err = check::confirmed_quorums_contain_blob_quorums(
            confirmation_threshold,
            &quorums,
            &blob_quorums,
        )
        .unwrap_err();

        assert_eq!(err, BitIndicesNotSorted);
    }
}

#[cfg(test)]
mod test_blob_quorums_contains_required_quorums {
    use crate::{check, error::CertVerificationError::*};

    #[test]
    fn success_when_blob_quorums_contain_required_quorums() {
        let blob_quorums = [0, 1, 2, 3].into();
        let required_quorums = [1, 2].into();

        let result = check::blob_quorums_contain_required_quorums(&blob_quorums, &required_quorums);

        assert_eq!(result, Ok(()));
    }

    #[test]
    fn blob_quorums_do_not_contain_required_quorums() {
        let blob_quorums = [0, 1].into();
        let required_quorums = [1, 2, 3].into(); // 2 and 3 are not in blob_quorums

        let err = check::blob_quorums_contain_required_quorums(&blob_quorums, &required_quorums)
            .unwrap_err();

        assert_eq!(err, BlobQuorumsDoNotContainRequiredQuorums);
    }

    #[test]
    fn required_quorums_bit_indices_not_sorted() {
        let blob_quorums = [0, 1].into();
        let required_quorums = [2, 1].into(); // Not sorted

        let err = check::blob_quorums_contain_required_quorums(&blob_quorums, &required_quorums)
            .unwrap_err();

        assert_eq!(err, BitIndicesNotSorted);
    }

    #[test]
    fn blob_quorums_bit_indices_not_sorted() {
        let blob_quorums = [1, 0].into(); // Not sorted
        let required_quorums = [0].into();

        let err = check::blob_quorums_contain_required_quorums(&blob_quorums, &required_quorums)
            .unwrap_err();

        assert_eq!(err, BitIndicesNotSorted);
    }
}

#[cfg(test)]
mod test_leaf_node_belongs_to_merkle_tree {
    use alloc::vec::Vec;
    use alloy_primitives::FixedBytes;

    use crate::{check, error::CertVerificationError::*, hash::keccak_v256};

    #[test]
    fn single_level_tree_left_child() {
        //   1||2
        //  /    \
        // 1      2

        let left_child: FixedBytes<32> = [1; 32].into();
        let right_sibling: FixedBytes<32> = [2; 32].into();
        let expected_root: FixedBytes<32> = keccak_v256([left_child, right_sibling].into_iter());

        let proof = right_sibling.into();

        // path: ... 0000 0000
        let path = 0;

        let result =
            check::leaf_node_belongs_to_merkle_tree(left_child, expected_root, proof, path);

        assert_eq!(result, Ok(()));
    }

    #[test]
    fn single_level_tree_right_child() {
        //   1||2
        //  /    \
        // 1      2

        let right_child: FixedBytes<32> = [2; 32].into();
        let left_sibling: FixedBytes<32> = [1; 32].into();
        let expected_root: FixedBytes<32> = keccak_v256([left_sibling, right_child].into_iter());

        let proof = left_sibling.into();

        // path: ... 0000 0001
        let path = 1;

        let result =
            check::leaf_node_belongs_to_merkle_tree(right_child, expected_root, proof, path);

        assert_eq!(result, Ok(()));
    }

    #[test]
    fn two_level_left_leaning_tree_left_child_inclusion() {
        //      (1||2)||3
        //        /    \
        //    1||2      3
        //   /    \
        // *1*     2

        let left_child: FixedBytes<32> = [1; 32].into();
        let right_sibling: FixedBytes<32> = [2; 32].into();
        let right_pibling: FixedBytes<32> = [3; 32].into();

        let parent = keccak_v256([left_child, right_sibling].into_iter());
        let expected_root = keccak_v256([parent, right_pibling].into_iter());

        let proof = [&right_sibling[..], &right_pibling[..]].concat().into();

        let path = 0;

        let result =
            check::leaf_node_belongs_to_merkle_tree(left_child, expected_root, proof, path);

        assert_eq!(result, Ok(()));
    }

    #[test]
    fn two_level_left_leaning_tree_right_child_inclusion() {
        //     (1||2)||3
        //       /    \
        //   1||2      3
        //  /    \
        // 1     *2*

        let right_child: FixedBytes<32> = [2; 32].into();
        let left_sibling: FixedBytes<32> = [1; 32].into();
        let right_pibling: FixedBytes<32> = [3; 32].into();

        let parent = keccak_v256([right_child, left_sibling].into_iter());
        let expected_root = keccak_v256([parent, right_pibling].into_iter());

        let proof = [&left_sibling[..], &right_pibling[..]].concat().into();

        // path: ... 0000 0000
        let path = 0;

        let result =
            check::leaf_node_belongs_to_merkle_tree(right_child, expected_root, proof, path);

        assert_eq!(result, Ok(()));
    }

    #[test]
    fn two_level_right_leaning_tree_left_child_inclusion() {
        // (1||2)||3
        //   /    \
        //  3    1||2
        //      /    \
        //    *1*     2

        let left_child: FixedBytes<32> = [1; 32].into();
        let right_sibling: FixedBytes<32> = [2; 32].into();
        let left_pibling: FixedBytes<32> = [3; 32].into();

        let parent = keccak_v256([left_child, right_sibling].into_iter());
        let expected_root = keccak_v256([left_pibling, parent].into_iter());

        let proof = [&right_sibling[..], &left_pibling[..]].concat().into();

        // path: ... 0000 0010
        let path = 2;

        let result =
            check::leaf_node_belongs_to_merkle_tree(left_child, expected_root, proof, path);

        assert_eq!(result, Ok(()));
    }

    #[test]
    fn two_level_right_leaning_tree_right_child_inclusion() {
        // (1||2)||3
        //   /    \
        //  3    1||2
        //      /    \
        //     1     *2*

        let right_child: FixedBytes<32> = [2; 32].into();
        let left_sibling: FixedBytes<32> = [1; 32].into();
        let left_pibling: FixedBytes<32> = [3; 32].into();

        let parent = keccak_v256([left_sibling, right_child].into_iter());
        let expected_root = keccak_v256([left_pibling, parent].into_iter());

        let proof = [&left_sibling[..], &left_pibling[..]].concat().into();

        // path: ... 0000 0011
        let path = 3;

        let result =
            check::leaf_node_belongs_to_merkle_tree(right_child, expected_root, proof, path);

        assert_eq!(result, Ok(()));
    }

    #[test]
    fn three_level_tree_complex_path() {
        //   ((3||(1||2))||4)
        //        /    \
        //   3||(1||2)  4
        //  /      \
        // 3      1||2
        //       /    \
        //     *1*     2

        let left_child: FixedBytes<32> = [1; 32].into();
        let right_sibling: FixedBytes<32> = [2; 32].into();
        let left_pibling: FixedBytes<32> = [3; 32].into();
        let right_grandparent: FixedBytes<32> = [4; 32].into();

        let right_parent = keccak_v256([left_child, right_sibling].into_iter());
        let left_grandparent = keccak_v256([left_pibling, right_parent].into_iter());
        let expected_root = keccak_v256([left_grandparent, right_grandparent].into_iter());

        let proof = [
            &right_sibling[..],
            &left_pibling[..],
            &right_grandparent[..],
        ]
        .concat()
        .into();

        // path: ... 0000 0010
        let path = 2;

        let result =
            check::leaf_node_belongs_to_merkle_tree(left_child, expected_root, proof, path);

        assert_eq!(result, Ok(()));
    }

    #[test]
    fn empty_proof_leaf_is_root() {
        let leaf: FixedBytes<32> = [1; 32].into();
        let expected_root = leaf;

        let proof = [].into();
        // path: ... 0000 0000
        let path = 0;

        let result = check::leaf_node_belongs_to_merkle_tree(leaf, expected_root, proof, path);

        assert_eq!(result, Ok(()));
    }

    #[test]
    fn proof_length_not_multiple_of_32() {
        let leaf: FixedBytes<32> = [1; 32].into();
        let expected_root: FixedBytes<32> = [2; 32].into();

        let proof = [1; 31].into(); // 31 bytes, not 32
        // path: ... 0000 0000
        let path = 0;

        let err =
            check::leaf_node_belongs_to_merkle_tree(leaf, expected_root, proof, path).unwrap_err();

        assert_eq!(err, MerkleProofLengthNotMultipleOf32Bytes);
    }

    #[test]
    fn path_too_short() {
        let leaf: FixedBytes<32> = [0; 32].into();
        let expected_root: FixedBytes<32> = [0; 32].into();

        let proof = [0; 257 * 32].into(); // path.len() == 256
        // path: ... 0000 0000
        let path = 0;

        let err =
            check::leaf_node_belongs_to_merkle_tree(leaf, expected_root, proof, path).unwrap_err();

        assert_eq!(err, MerkleProofPathTooShort);
    }

    #[test]
    fn invalid_proof_wrong_sibling() {
        //    1||2
        //   /    \
        // *1*     2

        let left_child: FixedBytes<32> = [1; 32].into();
        let correct_right_sibling: FixedBytes<32> = [2; 32].into();
        let wrong_right_sibling: FixedBytes<32> = [3; 32].into();
        let expected_root = keccak_v256([left_child, correct_right_sibling].into_iter());

        let proof = wrong_right_sibling.into();
        // path: ... 0000 0000
        let path = 0;

        let err = check::leaf_node_belongs_to_merkle_tree(left_child, expected_root, proof, path)
            .unwrap_err();

        assert_eq!(err, LeafNodeDoesNotBelongToMerkleTree);
    }

    #[test]
    fn invalid_proof_wrong_position() {
        //    1||2
        //   /    \
        // *1*     2

        let left_child: FixedBytes<32> = [1; 32].into();
        let right_sibling: FixedBytes<32> = [2; 32].into();
        let expected_root = keccak_v256([left_child, right_sibling].into_iter());

        let proof = right_sibling.into();
        // path: ... 0000 0001 (should be 0000 0000)
        let path = 1;

        let err = check::leaf_node_belongs_to_merkle_tree(left_child, expected_root, proof, path)
            .unwrap_err();

        assert_eq!(err, LeafNodeDoesNotBelongToMerkleTree);
    }

    #[test]
    fn max_depth_proof() {
        //      ...
        //    255||0
        //    /     \
        // *255*     0
        let mut left_current_node = [255; 32].into();
        let mut proof = Vec::new();

        for i in 0..=255u8 {
            let right_sibling_node = [i; 32].into();
            left_current_node = keccak_v256([left_current_node, right_sibling_node].into_iter());
            proof.extend_from_slice(right_sibling_node.as_ref());
        }

        let proof = proof.into();

        let leaf = [255; 32].into();
        let expected_root = left_current_node;

        // path: ... 0000 0000
        let path = 0;

        let result = check::leaf_node_belongs_to_merkle_tree(leaf, expected_root, proof, path);

        assert_eq!(result, Ok(()));
    }
}
