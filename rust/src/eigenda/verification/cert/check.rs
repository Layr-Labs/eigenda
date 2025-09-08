use crate::eigenda::{
    cert::{BlobCertificate, G1Point},
    verification::cert::hash::{HashExt, streaming_keccak256},
};
use alloy_primitives::{B256, Bytes, aliases::U96, keccak256};
use alloy_sol_types::SolValue;
use hashbrown::HashMap;
use tracing::{Level, instrument};

use crate::eigenda::verification::cert::{
    bitmap::{Bitmap, bit_indices_to_bitmap},
    convert,
    error::CertVerificationError::{self, *},
    hash::TruncHash,
    types::{
        BlockNumber, NonSigner, Quorum, QuorumNumber, Version,
        history::History,
        solidity::{SecurityThresholds, VersionedBlobParams},
    },
};

const THRESHOLD_DENOMINATOR: u128 = 100; // uint256 in sol

/// Validate that the certificate blob's version is valid. Otherwise it'll result a `coding_rate = 0`
/// which in turn will lead to division by zero at the subsequent `check::security_assumptions_are_met`
pub fn blob_version(
    cert_blob_version: Version,
    next_blob_version: Version,
) -> Result<(), CertVerificationError> {
    (cert_blob_version < next_blob_version)
        .then_some(())
        .ok_or(InvalidBlobVersion(cert_blob_version, next_blob_version))
}

#[instrument(level = Level::DEBUG, skip_all)]
pub fn equal_lengths(lengths: &[usize]) -> Result<(), CertVerificationError> {
    let Some(first) = lengths.first() else {
        return Err(EmptyVec);
    };

    lengths
        .iter()
        .all(|length| length == first)
        .then_some(())
        .ok_or(UnequalLengths)
}

#[instrument(level = Level::DEBUG, skip_all)]
pub fn not_empty<T>(slice: &[T]) -> Result<(), CertVerificationError> {
    (!slice.is_empty()).then_some(()).ok_or(EmptyVec)
}

#[instrument(level = Level::DEBUG, skip_all)]
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

#[cfg(feature = "stale-stakes-forbidden")]
#[instrument(level = Level::DEBUG, skip_all)]
pub fn quorums_last_updated_after_most_recent_stale_block(
    signed_quorums: &[QuorumNumber],
    reference_block: BlockNumber,
    quorum_update_block_number: HashMap<u8, BlockNumber>,
    window: u32,
) -> Result<(), CertVerificationError> {
    signed_quorums.iter().try_for_each(|signed_quorum| {
        let last_updated_at_block = *quorum_update_block_number
            .get(signed_quorum)
            .ok_or(MissingQuorumEntry)?;

        let most_recent_stale_block = reference_block.checked_sub(window).ok_or(Underflow)?;
        let is_recent = last_updated_at_block > most_recent_stale_block;
        is_recent.then_some(()).ok_or(StaleQuorum {
            last_updated_at_block,
            most_recent_stale_block,
            window,
        })
    })
}

#[instrument(level = Level::DEBUG, skip_all)]
pub fn cert_apks_equal_storage_apks(
    signed_quorums: &[QuorumNumber],
    reference_block: BlockNumber,
    apk_for_each_quorum: &[G1Point],
    apk_index_for_each_quorum: Vec<BlockNumber>,
    apk_history: HashMap<QuorumNumber, History<TruncHash>>,
) -> Result<(), CertVerificationError> {
    signed_quorums
        .iter()
        .zip(apk_for_each_quorum.iter())
        .zip(apk_index_for_each_quorum)
        .try_for_each(|((signed_quorum, cert_apk), apk_index)| {
            let cert_apk_hash = convert::point_to_hash(cert_apk);
            let cert_apk_trunc_hash: [u8; 24] = cert_apk_hash[..24].try_into().unwrap();
            let cert_apk_trunc_hash: TruncHash = cert_apk_trunc_hash.into();

            let storage_apk_trunc_hash = apk_history
                .get(signed_quorum)
                .ok_or(MissingQuorumEntry)?
                .try_get_at(apk_index)?
                .try_get_against(reference_block)?;

            (cert_apk_trunc_hash == storage_apk_trunc_hash)
                .then_some(())
                .ok_or(CertApkDoesNotEqualStorageApk {
                    cert_apk_trunc_hash,
                    storage_apk_trunc_hash,
                })
        })
}

#[instrument(level = Level::DEBUG, skip_all)]
pub fn security_assumptions_are_met(
    cert_blob_version: Version,
    versioned_blob_params: &HashMap<Version, VersionedBlobParams>,
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
    } = versioned_blob_params
        .get(&cert_blob_version)
        .ok_or(MissingVersionEntry(cert_blob_version))?;

    if confirmationThreshold <= adversaryThreshold {
        return Err(ConfirmationThresholdLessThanOrEqualToAdversaryThreshold(
            *confirmationThreshold,
            *adversaryThreshold,
        ));
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

    (n >= max_num_operators * 10_000)
        .then_some(())
        .ok_or(UnmetSecurityAssumptions)
}

#[instrument(level = Level::DEBUG, skip_all)]
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
            .checked_mul(U96::from(THRESHOLD_DENOMINATOR))
            .ok_or(Overflow)?;

        let right = total_stake
            .checked_mul(U96::from(confirmation_threshold))
            .ok_or(Overflow)?;

        confirmed_quorums.set(number as usize, left >= right);

        Ok::<_, CertVerificationError>(())
    })?;

    contains(confirmed_quorums, blob_quorums)
        .then_some(())
        .ok_or(ConfirmedQuorumsDoNotContainBlobQuorums)
}

#[instrument(level = Level::DEBUG, skip_all)]
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

#[instrument(level = Level::DEBUG, skip_all)]
pub fn blob_inclusion(
    blob_certificate: &BlobCertificate,
    expected_root: B256,
    proof: Bytes,
    sibling_path: u32,
) -> Result<(), CertVerificationError> {
    let blob_certificate = blob_certificate.hash_ext();
    let encoded = blob_certificate.abi_encode_packed();
    let leaf_node = keccak256(&encoded);
    leaf_node_belongs_to_merkle_tree(leaf_node, expected_root, proof, sibling_path)
}

#[instrument(level = Level::DEBUG, skip_all)]
fn leaf_node_belongs_to_merkle_tree(
    leaf_node: B256,
    expected_root: B256,
    proof: Bytes,
    sibling_path: u32,
) -> Result<(), CertVerificationError> {
    let proof_len = proof.len();
    if proof_len % 32 != 0 {
        return Err(MerkleProofLengthNotMultipleOf32Bytes(proof_len));
    }

    // will only fail when proof_depth exceeds u32::MAX
    let sibling_path = Bitmap::new([sibling_path as usize, 0, 0, 0]);

    let proof_depth = proof.len() / 32;
    let sibling_path_len = sibling_path.len();
    if sibling_path_len < proof_depth {
        return Err(MerkleProofPathTooShort {
            sibling_path_len,
            proof_depth,
        });
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
        let parent_node = streaming_keccak256(&[left_node, right_node]);
        current_node = parent_node;
    }

    let actual_root = current_node;
    (actual_root == expected_root)
        .then_some(())
        .ok_or(LeafNodeDoesNotBelongToMerkleTree)
}

#[cfg(test)]
mod test_blob_version {
    use crate::eigenda::verification::cert::{check, error::CertVerificationError::*};

    #[test]
    fn success_when_cert_version_less_than_next_version() {
        let result = check::blob_version(42, 43);
        assert_eq!(result, Ok(()));
    }

    #[test]
    fn invalid_blob_version_when_cert_version_equals_next_version() {
        let err = check::blob_version(42, 42).unwrap_err();
        assert_eq!(err, InvalidBlobVersion(42, 42));
    }

    #[test]
    fn invalid_blob_version_when_cert_version_greater_than_next_version() {
        let err = check::blob_version(43, 42).unwrap_err();
        assert_eq!(err, InvalidBlobVersion(43, 42));
    }
}

#[cfg(test)]
mod test_equal_lengths_and_not_empty {
    use crate::eigenda::verification::cert::{check, error::CertVerificationError::*};

    #[test]
    fn equal_lengths_success() {
        let result = check::equal_lengths(&[42, 42, 42, 42]);
        assert_eq!(result, Ok(()));
    }

    #[test]
    fn different_lengths_where_none_is_zero() {
        let err = check::equal_lengths(&[42, 43, 44, 45]).unwrap_err();
        assert_eq!(err, UnequalLengths);
    }

    #[test]
    fn first_length_zero_but_otherwise_equal_lengths() {
        let err = check::equal_lengths(&[0, 42, 42, 42]).unwrap_err();
        assert_eq!(err, UnequalLengths);
    }

    #[test]
    fn all_lengths_zero() {
        let result = check::equal_lengths(&[0, 0, 0, 0]);
        assert_eq!(result, Ok(()));
    }

    #[test]
    fn some_length_zero_but_otherwise_equal_lengths() {
        let err = check::equal_lengths(&[42, 42, 0, 42]).unwrap_err();
        assert_eq!(err, UnequalLengths);
    }

    #[test]
    fn not_empty_failure() {
        let err = check::not_empty::<u8>(&[]).unwrap_err();
        assert_eq!(err, EmptyVec);
    }

    #[test]
    fn not_empty_success() {
        let result = check::not_empty(&[42]);
        assert_eq!(result, Ok(()));
    }
}

#[cfg(feature = "stale-stakes-forbidden")]
#[cfg(test)]
mod test_non_signers_strictly_sorted_by_hash {
    use crate::eigenda::verification::cert::{
        check, error::CertVerificationError::*, types::NonSigner,
    };

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

#[cfg(feature = "stale-stakes-forbidden")]
#[cfg(test)]
mod test_quorums_last_updated_after_most_recent_stale_block {
    use crate::eigenda::verification::cert::{check, error::CertVerificationError::*};

    #[test]
    fn quorums_last_updated_after_most_recent_stale_block() {
        let reference_block = 42;
        let window = 1;
        let most_recent_stale_block = reference_block - window;

        let signed_quorums = [0];
        let quorum_update_block_number = signed_quorums
            .into_iter()
            .map(|signed_quorum| (signed_quorum, most_recent_stale_block + 1))
            .collect();

        let result = check::quorums_last_updated_after_most_recent_stale_block(
            &signed_quorums,
            reference_block,
            quorum_update_block_number,
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
        let quorum_update_block_number = signed_quorums
            .into_iter()
            .map(|signed_quorum| (signed_quorum, most_recent_stale_block - 1))
            .collect();

        let err = check::quorums_last_updated_after_most_recent_stale_block(
            &signed_quorums,
            reference_block,
            quorum_update_block_number,
            window,
        )
        .unwrap_err();

        assert_eq!(
            err,
            StaleQuorum {
                last_updated_at_block: 40,
                most_recent_stale_block: 41,
                window,
            }
        );
    }

    #[test]
    fn quorum_last_updated_at_most_recent_stale_block() {
        let reference_block = 42;
        let window = 1;
        let most_recent_stale_block = reference_block - window;

        let signed_quorums = [0];
        let quorum_update_block_number = signed_quorums
            .into_iter()
            .map(|signed_quorum| (signed_quorum, most_recent_stale_block))
            .collect();

        let err = check::quorums_last_updated_after_most_recent_stale_block(
            &signed_quorums,
            reference_block,
            quorum_update_block_number,
            window,
        )
        .unwrap_err();

        assert_eq!(
            err,
            StaleQuorum {
                last_updated_at_block: 41,
                most_recent_stale_block: 41,
                window,
            }
        );
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
        let quorum_update_block_number = signed_quorums
            .into_iter()
            .map(|signed_quorum| (signed_quorum, Default::default()))
            .collect();

        let err = check::quorums_last_updated_after_most_recent_stale_block(
            &signed_quorums,
            reference_block,
            quorum_update_block_number,
            window,
        )
        .unwrap_err();

        assert_eq!(err, Underflow);
    }
}

#[cfg(test)]
mod test_cert_apks_equal_storage_apks {
    use ark_bn254::{Fr, G1Projective};
    use ark_ec::{CurveGroup, PrimeGroup};
    use hashbrown::HashMap;

    use crate::eigenda::verification::cert::{
        check, convert,
        error::CertVerificationError::*,
        hash::TruncHash,
        types::{
            BlockNumber,
            conversions::IntoExt,
            history::{History, HistoryError::*, Update},
        },
    };

    #[test]
    fn cert_apk_equal_storage_apk() {
        let apk = (G1Projective::generator() * Fr::from(42)).into_affine();
        let apk_hash = convert::point_to_hash(&apk.into_ext());
        let apk_trunc_hash: [u8; 24] = apk_hash[..24].try_into().unwrap();
        let apk_trunc_hash: TruncHash = apk_trunc_hash.into();

        let signed_quorums = [0];
        let reference_block = 42;
        let apk_for_each_quorum = [apk.into_ext()];
        let apk_index_for_each_quorum = vec![0];

        let update = Update::new(42, 43, apk_trunc_hash).unwrap();
        let history = HashMap::from([(0, update)]);
        let apk_trunc_hash_history = History(history);
        let apk_history = HashMap::from([(0, apk_trunc_hash_history)]);

        let result = check::cert_apks_equal_storage_apks(
            &signed_quorums,
            reference_block,
            &apk_for_each_quorum,
            apk_index_for_each_quorum,
            apk_history,
        );

        assert_eq!(result, Ok(()));
    }

    #[test]
    fn cert_apk_does_not_equal_storage_apk() {
        let cert_apk = (G1Projective::generator() * Fr::from(42)).into_affine();
        let storage_apk = (G1Projective::generator() * Fr::from(43)).into_affine();
        let storage_apk_hash = convert::point_to_hash(&storage_apk.into_ext());
        let storage_apk_trunc_hash: [u8; 24] = storage_apk_hash[..24].try_into().unwrap();
        let storage_apk_trunc_hash: TruncHash = storage_apk_trunc_hash.into();

        let signed_quorums = [0];
        let reference_block = 42;
        let apk_for_each_quorum = [cert_apk.into_ext()];
        let apk_index_for_each_quorum = vec![0];

        let update = Update::new(42, 43, storage_apk_trunc_hash).unwrap();
        let history = HashMap::from([(0, update)]);
        let apk_trunc_hash_history = History(history);
        let apk_history = HashMap::from([(0, apk_trunc_hash_history)]);

        let err = check::cert_apks_equal_storage_apks(
            &signed_quorums,
            reference_block,
            &apk_for_each_quorum,
            apk_index_for_each_quorum,
            apk_history,
        )
        .unwrap_err();

        let cert_apk_hash = convert::point_to_hash(&cert_apk.into_ext());
        let cert_apk_trunc_hash: [u8; 24] = cert_apk_hash[..24].try_into().unwrap();
        let cert_apk_trunc_hash = cert_apk_trunc_hash.into();

        assert_eq!(
            err,
            CertApkDoesNotEqualStorageApk {
                cert_apk_trunc_hash,
                storage_apk_trunc_hash,
            }
        );
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
        let apk_history = HashMap::from([(0, apk_trunc_hash_history)]);

        let err = check::cert_apks_equal_storage_apks(
            &signed_quorums,
            reference_block,
            &apk_for_each_quorum,
            apk_index_for_each_quorum,
            apk_history,
        )
        .unwrap_err();

        assert_eq!(err, WrapHistoryError(MissingHistoryEntry(0)));
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
        let apk_history = HashMap::from([(0, apk_trunc_hash_history)]);

        let err = check::cert_apks_equal_storage_apks(
            &signed_quorums,
            STALE_REFERENCE_BLOCK,
            &apk_for_each_quorum,
            apk_index_for_each_quorum,
            apk_history,
        )
        .unwrap_err();

        assert_eq!(
            err,
            WrapHistoryError(ElementNotInInterval("41".into(), "[42, 43)".into()))
        );
    }
}

#[cfg(test)]
mod test_security_assumptions_are_met {
    use hashbrown::HashMap;

    use crate::eigenda::verification::cert::{
        check,
        error::CertVerificationError::*,
        types::{
            Version,
            solidity::{SecurityThresholds, VersionedBlobParams},
        },
    };

    #[test]
    fn success_when_security_assumptions_are_met() {
        let (version, versioned_blob_params, security_thresholds) = success_inputs();

        let result = check::security_assumptions_are_met(
            version,
            &versioned_blob_params,
            &security_thresholds,
        );

        assert_eq!(result, Ok(()));
    }

    #[test]
    fn security_assumptions_are_met_fails_with_missing_version_entry() {
        let (_version, versioned_blob_params, security_thresholds) = success_inputs();

        let err = check::security_assumptions_are_met(
            Version::MAX,
            &versioned_blob_params,
            &security_thresholds,
        )
        .unwrap_err();

        assert_eq!(err, MissingVersionEntry(Version::MAX));
    }

    #[test]
    fn security_assumptions_are_met_fails_when_confirmation_threshold_equals_adversary_threshold() {
        let (version, versioned_blob_params, mut security_thresholds) = success_inputs();

        security_thresholds.confirmationThreshold = security_thresholds.adversaryThreshold;

        let err = check::security_assumptions_are_met(
            version,
            &versioned_blob_params,
            &security_thresholds,
        )
        .unwrap_err();

        assert_eq!(
            err,
            ConfirmationThresholdLessThanOrEqualToAdversaryThreshold(1, 1)
        );
    }

    #[test]
    fn security_assumptions_are_met_fails_when_confirmation_threshold_less_than_adversary_threshold()
     {
        let (version, versioned_blob_params, mut security_thresholds) = success_inputs();

        security_thresholds.confirmationThreshold = security_thresholds.adversaryThreshold - 1;

        let err = check::security_assumptions_are_met(
            version,
            &versioned_blob_params,
            &security_thresholds,
        )
        .unwrap_err();

        assert_eq!(
            err,
            ConfirmationThresholdLessThanOrEqualToAdversaryThreshold(0, 1)
        );
    }

    #[test]
    fn security_assumptions_are_met_fails_with_underflow() {
        let (version, mut versioned_blob_params, mut security_thresholds) = success_inputs();

        // to trigger overflow (gamma * codingRate) < 100
        // where gamma = confirmation_threshold - adversary_threshold
        security_thresholds.confirmationThreshold = 101;
        security_thresholds.adversaryThreshold = 100;
        // gamma = 101 - 100 = 1
        let params = versioned_blob_params.get_mut(&version).unwrap();
        params.codingRate = 99;

        let err = check::security_assumptions_are_met(
            version,
            &versioned_blob_params,
            &security_thresholds,
        )
        .unwrap_err();

        assert_eq!(err, Underflow);
    }

    #[test]
    fn security_assumptions_are_met_fails_with_unmet_security_assumptions() {
        let (version, versioned_blob_params, mut security_thresholds) = success_inputs();

        // from success_inputs:
        // gamma = confirmation_threshold - adversary_threshold = 101 - 1 = 100
        // since the success_inputs are at the limit
        // any disturbance will cause UnmetSecurityAssumptions so
        security_thresholds.adversaryThreshold = 2; // instead of 1, resulting in gamma = 99

        let err = check::security_assumptions_are_met(
            version,
            &versioned_blob_params,
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
        let versioned_blob_params = HashMap::from([(
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

        (version, versioned_blob_params, security_thresholds)
    }
}

#[cfg(test)]
mod test_confirmed_quorums_contains_blob_quorums {
    use alloy_primitives::aliases::U96;
    use ark_bn254::G1Affine;

    use crate::eigenda::verification::cert::{
        bitmap::BitmapError::*, check, error::CertVerificationError::*, types::Quorum,
    };

    #[test]
    fn success_when_confirmed_quorums_contain_blob_quorums() {
        let confirmation_threshold = 100;

        // in this example:
        //     quorum is confirmed if signed_stake * 100 > total_stake * 100
        //     quorum is confirmed if signed_stake * THRESHOLD_DENOMINATOR >= total_skate * confirmation_threshold
        let quorums = [
            Quorum {
                number: 0,
                total_stake: U96::from(42),
                signed_stake: U96::from(43),
                ..Default::default()
            },
            Quorum {
                number: 1,
                apk: G1Affine::default(),
                total_stake: U96::from(42),
                signed_stake: U96::from(42),
            },
            Quorum {
                number: 2,
                total_stake: U96::from(42),
                signed_stake: U96::from(41),
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
                total_stake: U96::from(42),
                signed_stake: U96::from(43),
                ..Default::default()
            },
            Quorum {
                number: 1,
                apk: G1Affine::default(),
                total_stake: U96::from(42),
                signed_stake: U96::from(42),
            },
            Quorum {
                number: 2,
                total_stake: U96::from(42),
                signed_stake: U96::from(41),
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
            total_stake: U96::from(42),
            signed_stake: U96::MAX, // Will overflow when multiplied by THRESHOLD_DENOMINATOR
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
            total_stake: U96::MAX,
            signed_stake: U96::from(43),
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
            total_stake: U96::from(42),
            signed_stake: U96::from(43),
            ..Default::default()
        }];

        let blob_quorums = [1, 0].into(); // Not sorted

        let err = check::confirmed_quorums_contain_blob_quorums(
            confirmation_threshold,
            &quorums,
            &blob_quorums,
        )
        .unwrap_err();

        assert_eq!(err, WrapBitmapError(IndicesNotSorted));
    }
}

#[cfg(test)]
mod test_blob_quorums_contains_required_quorums {
    use crate::eigenda::verification::cert::{
        bitmap::BitmapError::*, check, error::CertVerificationError::*,
    };

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

        assert_eq!(err, WrapBitmapError(IndicesNotSorted));
    }

    #[test]
    fn blob_quorums_bit_indices_not_sorted() {
        let blob_quorums = [1, 0].into(); // Not sorted
        let required_quorums = [0].into();

        let err = check::blob_quorums_contain_required_quorums(&blob_quorums, &required_quorums)
            .unwrap_err();

        assert_eq!(err, WrapBitmapError(IndicesNotSorted));
    }
}

#[cfg(test)]
mod test_leaf_node_belongs_to_merkle_tree {
    use alloy_primitives::FixedBytes;

    use crate::eigenda::verification::cert::{
        check, error::CertVerificationError::*, hash::streaming_keccak256,
    };

    #[test]
    fn single_level_tree_left_child() {
        //   1||2
        //  /    \
        // 1      2

        let left_child: FixedBytes<32> = [1; 32].into();
        let right_sibling: FixedBytes<32> = [2; 32].into();
        let expected_root: FixedBytes<32> = streaming_keccak256(&[left_child, right_sibling]);

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
        let expected_root: FixedBytes<32> = streaming_keccak256(&[left_sibling, right_child]);

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

        let parent = streaming_keccak256(&[left_child, right_sibling]);
        let expected_root = streaming_keccak256(&[parent, right_pibling]);

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

        let parent = streaming_keccak256(&[right_child, left_sibling]);
        let expected_root = streaming_keccak256(&[parent, right_pibling]);

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

        let parent = streaming_keccak256(&[left_child, right_sibling]);
        let expected_root = streaming_keccak256(&[left_pibling, parent]);

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

        let parent = streaming_keccak256(&[left_sibling, right_child]);
        let expected_root = streaming_keccak256(&[left_pibling, parent]);

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

        let right_parent = streaming_keccak256(&[left_child, right_sibling]);
        let left_grandparent = streaming_keccak256(&[left_pibling, right_parent]);
        let expected_root = streaming_keccak256(&[left_grandparent, right_grandparent]);

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

        assert_eq!(err, MerkleProofLengthNotMultipleOf32Bytes(31));
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

        assert_eq!(
            err,
            MerkleProofPathTooShort {
                sibling_path_len: 256,
                proof_depth: 257,
            }
        );
    }

    #[test]
    fn invalid_proof_wrong_sibling() {
        //    1||2
        //   /    \
        // *1*     2

        let left_child: FixedBytes<32> = [1; 32].into();
        let correct_right_sibling: FixedBytes<32> = [2; 32].into();
        let wrong_right_sibling: FixedBytes<32> = [3; 32].into();
        let expected_root = streaming_keccak256(&[left_child, correct_right_sibling]);

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
        let expected_root = streaming_keccak256(&[left_child, right_sibling]);

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
            left_current_node = streaming_keccak256(&[left_current_node, right_sibling_node]);
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
