//! EigenDA certificate verification using BLS signature aggregation
//!
//! This module implements comprehensive verification of EigenDA certificates,
//! validating the cryptographic integrity and security properties of data
//! availability certificates.
//!
//! ## Overview
//!
//! Certificate verification ensures that:
//! - The certificate was signed by a sufficient stake-weighted quorum
//! - All cryptographic signatures are valid (BLS signature aggregation)
//! - Security thresholds are met for data availability guarantees
//! - Historical operator state is consistent at the reference block
//!
//! ## Verification Process
//!
//! The verification follows a multi-stage approach:
//!
//! 1. **Signature Verification**: Validate BLS aggregate signatures
//! 2. **Stake Validation**: Ensure sufficient stake signed the certificate  
//! 3. **Quorum Checks**: Verify required quorums participated
//! 4. **Security Thresholds**: Enforce minimum security requirements
//! 5. **Historical Consistency**: Validate operator state at reference block
//!
//! ## BLS Signature Aggregation
//!
//! EigenDA uses BLS signatures over the BN254 curve for efficient aggregation:
//! - Individual operator signatures are aggregated into a single signature
//! - Public keys are aggregated using elliptic curve operations
//! - Verification checks the aggregate signature against the aggregate public key
//!
//! ## Security Model
//!
//! The verification enforces EigenDA's security model:
//! - **Confirmation Threshold**: Minimum percentage of honest stake required
//! - **Adversary Threshold**: Maximum percentage of adversarial stake tolerated
//! - **Quorum Requirements**: Specific quorums that must participate
//!
//! ## Reference Implementation
//!
//! Based on the [EigenDA Solidity implementation](https://github.com/Layr-Labs/eigenda/blob/60d438705b30e899777736cdffcc478ded08cc76/contracts/src/integrations/cert/libraries/EigenDACertVerificationLib.sol#L125)

pub mod bitmap;
mod check;
pub mod convert;
pub mod error;
pub mod hash;
mod signature;
pub mod types;

use alloy_primitives::{B256, Bytes};
use ark_bn254::{G1Affine, G2Affine};
use hashbrown::HashMap;
use tracing::instrument;

use crate::cert::solidity::SecurityThresholds;
use crate::cert::{BatchHeaderV2, BlobInclusionInfo, G1Point, NonSignerStakesAndSignature};
use crate::verification::cert::error::CertVerificationError::{self, *};
use crate::verification::cert::hash::HashExt;
use crate::verification::cert::types::history::History;
use crate::verification::cert::types::{
    BlockNumber, NonSigner, Quorum, QuorumNumber, Stake, Storage,
};

/// Input parameters for certificate verification
///
/// Contains all the data needed to perform comprehensive certificate validation,
/// including on-chain state data, signature information, and security parameters.
#[derive(Clone, Debug)]
pub struct CertVerificationInputs {
    /// Batch header containing the merkle root and reference block number
    pub batch_header: BatchHeaderV2,
    /// Blob inclusion proof and certificate information
    pub blob_inclusion_info: BlobInclusionInfo,
    /// Non-signer information and aggregated signatures
    pub non_signer_stakes_and_signature: NonSignerStakesAndSignature,
    /// Security thresholds for confirmation and adversary limits
    pub security_thresholds: SecurityThresholds,
    /// Quorum numbers required to sign certificates
    pub required_quorum_numbers: alloy_primitives::Bytes,
    /// Quorum numbers that actually signed this certificate
    pub signed_quorum_numbers: alloy_primitives::Bytes,
    /// Historical on-chain storage data for verification
    pub storage: Storage,
}

/// Performs comprehensive EigenDA certificate verification.
///
/// This is the main entry point for validating data availability certificates in the EigenDA
/// system. It implements a multi-stage verification process that ensures cryptographic integrity,
/// sufficient stake participation, and compliance with security parameters.
///
/// # Verification Process
///
/// The function executes the following verification stages in order:
///
/// ## 1. Blob Inclusion Verification
/// - Validates the blob certificate is included in the batch using Merkle proofs
/// - Ensures the blob index corresponds to the correct position in the batch
///
/// ## 2. Version and Security Validation  
/// - Checks blob version compatibility against available versions
/// - Enforces security assumptions are met for the blob's coding parameters
/// - Validates confirmation thresholds and adversarial assumptions
///
/// ## 3. Input Validation
/// - Ensures signed quorum numbers are not empty
/// - Verifies corresponding array lengths match across all input collections
/// - Validates reference block precedes current block
///
/// ## 4. Non-Signer Processing
/// - Reconstructs non-signer data from public keys and bitmap indices
/// - Validates non-signers are sorted by hash (required for verification)
/// - Retrieves historical quorum participation bitmaps at reference block
///
/// ## 5. Quorum Stake Calculation
/// - Processes each signing quorum to compute stake distributions
/// - Calculates signed stake by subtracting non-signer stakes from totals
/// - Validates sufficient stake participated in each quorum
///
/// ## 6. Signature Aggregation and Verification
/// - Aggregates public keys of signing operators across all quorums
/// - Computes expected aggregate public key excluding non-signers
/// - Verifies BLS signature against batch header hash using aggregated keys
///
/// ## 7. Security Threshold Enforcement
/// - Validates quorums meeting confirmation threshold include blob quorums
/// - Ensures blob quorums contain all required quorum numbers
/// - Enforces minimum security guarantees for data availability
///
/// # Arguments
///
/// * `inputs` - Complete verification input containing:
///   - `batch_header` - Batch metadata with reference block and root hash
///   - `blob_inclusion_info` - Certificate and Merkle inclusion proof  
///   - `non_signer_stakes_and_signature` - BLS signature data and non-signer info
///   - `security_thresholds` - Required confirmation and adversarial thresholds
///   - `required_quorum_numbers` - Quorums mandated for this certificate type
///   - `signed_quorum_numbers` - Quorums that actually signed the certificate
///   - `storage` - Historical on-chain state data for validation
///
/// # Returns
///
/// * `Ok(())` - Certificate passes all verification checks and is valid
/// * `Err(CertVerificationError)` - Verification failed with specific error details
///
/// # Errors
///
/// Returns [`CertVerificationError`] for various validation failures:
///
/// ## Cryptographic Failures
/// - `SignatureVerificationFailed` - BLS signature validation failed
/// - `LeafNodeDoesNotBelongToMerkleTree` - Invalid inclusion proof
///
/// ## Stake and Quorum Failures  
/// - `InsufficientStake` - Not enough stake signed the certificate
/// - `EmptyBlobQuorums` - No quorums specified for the blob
/// - `MissingQuorumEntry` - Referenced quorum not found in historical data
///
/// ## Parameter Validation Failures
/// - `UnsupportedVersion` - Blob version not supported
/// - `SecurityAssumptionsNotMet` - Coding parameters violate security model
/// - `ReferenceBlockDoesNotPrecedeCurrentBlock` - Invalid block ordering
///
/// ## Data Consistency Failures
/// - `MissingSignerEntry` - Operator not found in historical data
/// - `ArrayLengthMismatch` - Input array lengths don't correspond
/// - `NonSignersNotSorted` - Non-signers not properly ordered
///
/// # Security Considerations
///
/// This function is critical for EigenDA's security model. It ensures:
/// - Only certificates with sufficient economic backing are accepted
/// - Historical operator state is accurately reflected at reference blocks
/// - BLS signature aggregation is performed correctly to prevent forgeries
/// - Security parameters enforce adequate redundancy for data recovery
#[instrument]
pub fn verify(inputs: CertVerificationInputs) -> Result<(), CertVerificationError> {
    let CertVerificationInputs {
        batch_header,
        blob_inclusion_info,
        non_signer_stakes_and_signature,
        security_thresholds,
        required_quorum_numbers,
        signed_quorum_numbers,
        storage,
    } = inputs;

    let NonSignerStakesAndSignature {
        non_signer_quorum_bitmap_indices,
        non_signer_pubkeys,
        quorum_apks,
        apk_g2,
        sigma,
        quorum_apk_indices,
        total_stake_indices,
        non_signer_stake_indices,
    } = non_signer_stakes_and_signature;

    let Storage {
        quorum_count,
        current_block,
        quorum_bitmap_history,
        operator_stake_history,
        total_stake_history,
        apk_history,
        versioned_blob_params,
        next_blob_version,
        #[cfg(feature = "stale-stakes-forbidden")]
        staleness,
    } = storage;

    check::blob_inclusion(
        &blob_inclusion_info.blob_certificate,
        batch_header.batch_root.into(),
        blob_inclusion_info.inclusion_proof,
        blob_inclusion_info.blob_index,
    )?;

    let cert_blob_version = blob_inclusion_info.blob_certificate.blob_header.version;
    check::blob_version(cert_blob_version, next_blob_version)?;

    check::security_assumptions_are_met(
        cert_blob_version,
        &versioned_blob_params,
        &security_thresholds,
    )?;

    check::not_empty(&signed_quorum_numbers)?;

    let lengths = [
        signed_quorum_numbers.len(),
        quorum_apks.len(),
        quorum_apk_indices.len(),
        total_stake_indices.len(),
        non_signer_stake_indices.len(),
    ];

    check::equal_lengths(&lengths).unwrap();

    let lengths = [
        non_signer_pubkeys.len(),
        non_signer_quorum_bitmap_indices.len(),
    ];

    check::equal_lengths(&lengths).unwrap();

    if batch_header.reference_block_number >= current_block {
        return Err(ReferenceBlockDoesNotPrecedeCurrentBlock(
            batch_header.reference_block_number,
            current_block,
        ));
    }

    // assumption: collection_a[i] corresponds to collection_b[i] for all i
    let non_signers = non_signer_pubkeys
        .into_iter()
        .zip(non_signer_quorum_bitmap_indices.into_iter())
        .map(|(pk, quorum_bitmap_history_index)| {
            let pk_hash = convert::point_to_hash(&pk);

            let quorum_bitmap_history = quorum_bitmap_history
                .get(&pk_hash)
                .ok_or(MissingSignerEntry)?
                .try_get_at(quorum_bitmap_history_index)?
                .try_get_against(batch_header.reference_block_number)?;

            let non_signer = NonSigner {
                pk: pk.into(),
                pk_hash,
                quorum_bitmap_history,
            };
            Ok::<_, CertVerificationError>(non_signer)
        })
        .collect::<Result<Vec<_>, _>>()?;

    check::non_signers_strictly_sorted_by_hash(&non_signers)?;

    let quorums = process_quorums(
        &signed_quorum_numbers,
        &quorum_apks,
        &total_stake_indices,
        &non_signer_stake_indices,
        &total_stake_history,
        batch_header.reference_block_number,
        &operator_stake_history,
        &non_signers,
    )?;

    let signers_apk = signature::aggregation::aggregate(quorum_count, &non_signers, &quorums)?;

    #[cfg(feature = "stale-stakes-forbidden")]
    if staleness.stale_stakes_forbidden {
        check::quorums_last_updated_after_most_recent_stale_block(
            &signed_quorum_numbers,
            batch_header.reference_block_number,
            staleness.quorum_update_block_number,
            staleness.min_withdrawal_delay_blocks,
        )?;
    }

    check::cert_apks_equal_storage_apks(
        &signed_quorum_numbers,
        batch_header.reference_block_number,
        &quorum_apks,
        quorum_apk_indices,
        apk_history,
    )?;

    let msg_hash = batch_header.hash_ext();
    let apk_g2: G2Affine = apk_g2.into();
    let sigma: G1Affine = sigma.into();

    if !signature::verification::verify(msg_hash, signers_apk, apk_g2, sigma) {
        return Err(SignatureVerificationFailed);
    }

    let blob_quorums = blob_inclusion_info
        .blob_certificate
        .blob_header
        .quorum_numbers;

    if blob_quorums.is_empty() {
        return Err(EmptyBlobQuorums);
    }

    check::confirmed_quorums_contain_blob_quorums(
        security_thresholds.confirmationThreshold,
        &quorums,
        &blob_quorums,
    )?;

    check::blob_quorums_contain_required_quorums(&blob_quorums, &required_quorum_numbers)?;

    Ok(())
}

/// Processes and validates quorum data for certificate verification.
///
/// This function computes the stake distribution for each quorum involved in signing
/// a certificate, calculating both total stake and signed stake by accounting for
/// non-signing operators. It constructs validated `Quorum` objects containing the
/// aggregate public key and stake information needed for BLS signature verification.
///
/// # Returns
///
/// * `Ok(Vec<Quorum>)` - Vector of processed quorums with computed stake distributions
/// * `Err(CertVerificationError)` - If stake calculation fails due to:
///   - Missing quorum or signer entries in historical data
///   - Invalid stake indices or block number references  
///   - Arithmetic underflow when computing signed stake
///
/// # Algorithm
///
/// For each quorum:
/// 1. **Total Stake Lookup**: Retrieves total stake at the reference block using the provided index
/// 2. **Non-Signer Filtering**: Identifies non-signers required to participate in this quorum
/// 3. **Unsigned Stake Calculation**: Sums stake of all filtered non-signers at reference block
/// 4. **Signed Stake Computation**: Subtracts unsigned stake from total stake
/// 5. **Quorum Construction**: Creates validated quorum with APK and computed stakes
///
/// # Invariants
///
/// - All input collections must have corresponding elements at the same indices
/// - `signed_stake = total_stake - unsigned_stake` must not underflow
/// - Historical data must exist for all referenced quorums and operators
/// - Non-signer quorum bitmaps must accurately reflect participation requirements
#[allow(clippy::too_many_arguments)]
fn process_quorums(
    signed_quorum_numbers: &Bytes,
    quorum_apks: &[G1Point],
    total_stake_indices: &[u32],
    non_signer_stake_indices: &[Vec<u32>],
    total_stake_history: &HashMap<QuorumNumber, History<Stake>>,
    reference_block_number: BlockNumber,
    operator_stake_history: &HashMap<B256, HashMap<QuorumNumber, History<Stake>>>,
    non_signers: &[NonSigner],
) -> Result<Vec<Quorum>, CertVerificationError> {
    // assumption: collection_a[i] corresponds to collection_b[i] for all i, for all (a, b)
    signed_quorum_numbers
        .iter()
        .zip(quorum_apks.iter())
        .zip(total_stake_indices.iter())
        .zip(non_signer_stake_indices.iter())
        .map(
            |(
                ((signed_quorum, apk), total_stake_index),
                stake_index_for_each_required_non_signer,
            )| {
                let total_stake = total_stake_history
                    .get(signed_quorum)
                    .ok_or(MissingQuorumEntry)?
                    .try_get_at(*total_stake_index)?
                    .try_get_against(reference_block_number)?;

                let bit = *signed_quorum as usize;
                let unsigned_stake = non_signers
                    .iter()
                    .filter(|non_signer| {
                        // whether signer was required to sign this quorum
                        non_signer.quorum_bitmap_history[bit]
                    })
                    // assumption: collection_a[i] corresponds to collection_b[i] for all i
                    .zip(stake_index_for_each_required_non_signer.iter())
                    .map(|(required_non_signer, stake_index)| {
                        let stake = operator_stake_history
                            .get(&required_non_signer.pk_hash)
                            .ok_or(MissingSignerEntry)?
                            .get(signed_quorum)
                            .ok_or(MissingQuorumEntry)?
                            .try_get_at(*stake_index)?
                            .try_get_against(reference_block_number)?;
                        Ok(stake)
                    })
                    .sum::<Result<_, CertVerificationError>>()?;

                let signed_stake = total_stake.checked_sub(unsigned_stake).ok_or(Underflow)?;

                let apk: G1Affine = (*apk).into();
                let quorum = Quorum {
                    number: *signed_quorum,
                    apk,
                    total_stake,
                    signed_stake,
                };

                Ok::<_, CertVerificationError>(quorum)
            },
        )
        .collect()
}

#[cfg(test)]
mod tests {
    use alloy_primitives::aliases::U96;
    use alloy_primitives::{B256, Bytes, keccak256};
    use alloy_sol_types::SolValue;
    use ark_bn254::{Fr, G1Affine, G1Projective, G2Projective};
    use ark_ec::{CurveGroup, PrimeGroup};
    use hashbrown::HashMap;

    use crate::cert::solidity::{SecurityThresholds, VersionedBlobParams};
    use crate::cert::{
        BatchHeaderV2, BlobCertificate, BlobCommitment, BlobHeaderV2, BlobInclusionInfo, G1Point,
        NonSignerStakesAndSignature,
    };
    use crate::verification::cert::bitmap::Bitmap;
    use crate::verification::cert::bitmap::BitmapError::*;
    use crate::verification::cert::error::CertVerificationError::*;
    use crate::verification::cert::hash::{HashExt, TruncHash, streaming_keccak256};
    #[cfg(feature = "stale-stakes-forbidden")]
    use crate::verification::cert::types::Staleness;
    use crate::verification::cert::types::Storage;
    use crate::verification::cert::types::history::HistoryError::*;
    use crate::verification::cert::types::history::{History, Update};
    use crate::verification::cert::{CertVerificationInputs, convert, verify};

    #[test]
    fn success() {
        let inputs = success_inputs();

        let result = verify(inputs);
        assert_eq!(result, Ok(()));
    }

    #[test]
    fn leaf_node_does_not_belong_to_merkle_tree() {
        let mut inputs = success_inputs();

        // any change to blobCertificate causes the leaf node hash to differ
        inputs.blob_inclusion_info.blob_certificate.signature = [0u8; 32].into();

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, LeafNodeDoesNotBelongToMerkleTree);
    }

    #[test]
    fn reference_block_past_current_block() {
        let mut inputs = success_inputs();

        inputs.batch_header.reference_block_number = 43;
        inputs.storage.current_block = 42;

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, ReferenceBlockDoesNotPrecedeCurrentBlock(43, 42));
    }

    #[test]
    fn reference_block_at_current_block() {
        let mut inputs = success_inputs();

        inputs.batch_header.reference_block_number = 42;
        inputs.storage.current_block = 42;

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, ReferenceBlockDoesNotPrecedeCurrentBlock(42, 42));
    }

    #[test]
    fn empty_non_signer_vecs() {
        let mut inputs = success_inputs();

        inputs
            .non_signer_stakes_and_signature
            .non_signer_pubkeys
            .clear();

        inputs
            .non_signer_stakes_and_signature
            .non_signer_quorum_bitmap_indices
            .clear();

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, SignatureVerificationFailed);
    }

    #[test]
    fn empty_quorum_vecs() {
        let mut inputs = success_inputs();

        inputs.signed_quorum_numbers = [].into();

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, EmptyVec);
    }

    #[cfg(feature = "stale-stakes-forbidden")]
    #[test]
    fn stale_stakes_forbidden() {
        let mut inputs = success_inputs();

        inputs.storage.staleness.stale_stakes_forbidden = true;
        inputs
            .storage
            .staleness
            .quorum_update_block_number
            .insert(0, 41);

        inputs.storage.staleness.min_withdrawal_delay_blocks = 1;

        let err = verify(inputs).unwrap_err();

        assert_eq!(
            err,
            StaleQuorum {
                last_updated_at_block: 41,
                most_recent_stale_block: 41,
                window: 1,
            }
        );
    }

    #[test]
    fn quorum_bitmap_history_history_missing_signer_entry() {
        let mut inputs = success_inputs();

        inputs.storage.quorum_bitmap_history.clear();

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, MissingSignerEntry);
    }

    #[test]
    fn quorum_bitmap_history_history_missing_history_entry() {
        let mut inputs = success_inputs();

        inputs
            .non_signer_stakes_and_signature
            .non_signer_quorum_bitmap_indices[0] = 42;

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, WrapHistoryError(MissingHistoryEntry(42)));
    }

    #[test]
    fn quorum_bitmap_history_history_reference_block_not_in_interval() {
        let mut inputs = success_inputs();

        inputs
            .storage
            .quorum_bitmap_history
            .iter_mut()
            .for_each(|(_, v)| {
                v.0.insert(0, Update::new(141, 143, Default::default()).unwrap());
            });

        let err = verify(inputs).unwrap_err();

        assert_eq!(
            err,
            WrapHistoryError(ElementNotInInterval("42".into(), "[141, 143)".into()))
        );
    }

    #[test]
    fn non_signers_not_strictly_sorted_by_hash() {
        let mut inputs = success_inputs();

        inputs
            .non_signer_stakes_and_signature
            .non_signer_pubkeys
            .reverse();

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, NotStrictlySortedByHash);
    }

    #[test]
    fn total_stake_history_missing_quorum_entry() {
        let mut inputs = success_inputs();

        inputs.storage.total_stake_history.clear();

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, MissingQuorumEntry);
    }

    #[test]
    fn total_stake_history_missing_history_entry() {
        let mut inputs = success_inputs();

        inputs
            .storage
            .total_stake_history
            .insert(0, Default::default());

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, WrapHistoryError(MissingHistoryEntry(0)));
    }

    #[test]
    fn total_stake_history_reference_block_not_in_interval() {
        let mut inputs = success_inputs();

        inputs
            .storage
            .total_stake_history
            .iter_mut()
            .for_each(|(_, v)| {
                v.0.insert(0, Update::new(141, 143, Default::default()).unwrap());
            });

        let err = verify(inputs).unwrap_err();

        assert_eq!(
            err,
            WrapHistoryError(ElementNotInInterval("42".into(), "[141, 143)".into()))
        );
    }

    #[test]
    fn stake_history_missing_signer_entry() {
        let mut inputs = success_inputs();

        inputs.storage.operator_stake_history.clear();

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, MissingSignerEntry);
    }

    #[test]
    fn stake_history_missing_quorum_entry() {
        let mut inputs = success_inputs();

        inputs.storage.operator_stake_history.iter_mut().for_each(
            |(_, stake_history_by_quorum)| {
                stake_history_by_quorum.clear();
            },
        );

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, MissingQuorumEntry);
    }

    #[test]
    fn stake_history_missing_history_entry() {
        let mut inputs = success_inputs();

        inputs.storage.operator_stake_history.iter_mut().for_each(
            |(_, stake_history_by_quorum)| {
                stake_history_by_quorum.insert(0, Default::default());
            },
        );

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, WrapHistoryError(MissingHistoryEntry(0)));
    }

    #[test]
    fn stake_history_reference_block_not_in_interval() {
        let mut inputs = success_inputs();

        inputs.storage.operator_stake_history.iter_mut().for_each(
            |(_, stake_history_by_quorum)| {
                stake_history_by_quorum.iter_mut().for_each(|(_, v)| {
                    v.0.insert(0, Update::new(141, 143, Default::default()).unwrap());
                })
            },
        );

        let err = verify(inputs).unwrap_err();

        assert_eq!(
            err,
            WrapHistoryError(ElementNotInInterval("42".into(), "[141, 143)".into()))
        );
    }

    #[test]
    fn stake_underflow() {
        let mut inputs = success_inputs();

        inputs
            .storage
            .total_stake_history
            .iter_mut()
            .for_each(|(_, v)| {
                v.0.insert(0, Update::new(41, 43, U96::from(29)).unwrap());
            });

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, Underflow);
    }

    #[test]
    fn aggregation_failure() {
        let mut inputs = success_inputs();

        inputs.storage.quorum_count = 1;

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, WrapBitmapError(IndexThanOrEqualToUpperBound));
    }

    #[test]
    fn signature_verification_failure() {
        let mut inputs = success_inputs();

        inputs.non_signer_stakes_and_signature.sigma = G1Point::default();

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, SignatureVerificationFailed);
    }

    #[test]
    fn security_assumptions_not_met() {
        let mut inputs = success_inputs();

        let params = inputs.storage.versioned_blob_params.get_mut(&42).unwrap();
        params.numChunks = 43;

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, UnmetSecurityAssumptions);
    }

    #[test]
    fn confirmed_quorums_do_not_contain_blob_quorums() {
        let mut inputs = success_inputs();

        inputs
            .storage
            .versioned_blob_params
            .iter_mut()
            .for_each(|(_, versioned_blob_params)| {
                versioned_blob_params.maxNumOperators = 0;
            });

        inputs
            .blob_inclusion_info
            .blob_certificate
            .blob_header
            .quorum_numbers = [0, 1, 2].into(); // while confirmed_quorums: [0, 2]

        // any change to blobCertificate requires recomputing...
        let secret_keys = vec![Fr::from(43u64), Fr::from(44u64)];
        let (batch_header, sigma) =
            compute_batch_header_and_sigma(&inputs.blob_inclusion_info, secret_keys);

        inputs.batch_header = batch_header;

        inputs.non_signer_stakes_and_signature.sigma = sigma.into();

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, ConfirmedQuorumsDoNotContainBlobQuorums);
    }

    #[test]
    fn blob_quorums_do_not_contain_required_quorums() {
        let mut inputs = success_inputs();
        inputs.required_quorum_numbers = [1].into(); // 3 is not in blob_quorums: [0, 2]

        let err = verify(inputs).unwrap_err();

        assert_eq!(err, BlobQuorumsDoNotContainRequiredQuorums);
    }

    fn success_inputs() -> CertVerificationInputs {
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
                    commitment: BlobCommitment::default(),
                    payment_header_hash: [42; 32],
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
                non_signer0_g1_pk.into(),
                non_signer1_g1_pk.into(),
                non_signer2_g1_pk.into(),
            ],
            quorum_apks: vec![apk_for_each_quorum[0].into(), apk_for_each_quorum[1].into()],
            apk_g2: apk_g2.into(),
            sigma: sigma.into(),
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

        let non_signer0_pk_hash = convert::point_to_hash(&non_signer0_g1_pk.into());
        let non_signer1_pk_hash = convert::point_to_hash(&non_signer1_g1_pk.into());
        let non_signer2_pk_hash = convert::point_to_hash(&non_signer2_g1_pk.into());
        let signer3_pk_hash = convert::point_to_hash(&signer3_g1_pk.into());
        let signer4_pk_hash = convert::point_to_hash(&signer4_g1_pk.into());
        let optional_non_signer5_pk_hash =
            convert::point_to_hash(&optional_non_signer5_g1_pk.into());

        // by sheer coincidence the first 3 hashes are already sorted
        let pk_hashes = [
            non_signer0_pk_hash,
            non_signer1_pk_hash,
            non_signer2_pk_hash,
            signer3_pk_hash,
            signer4_pk_hash,
            optional_non_signer5_pk_hash,
        ];

        let quorum_bitmap_history = {
            let quorum_bitmap_histories = vec![
                Bitmap::new([5, 0, 0, 0]), // 1 0 1
                Bitmap::new([6, 0, 0, 0]), // 1 1 0
                Bitmap::new([7, 0, 0, 0]), // 1 1 1
                Bitmap::new([4, 0, 0, 0]), // 1 0 0
                Bitmap::new([1, 0, 0, 0]), // 0 0 1
                Bitmap::new([0, 0, 0, 0]), // 0 0 0
            ];

            pk_hashes
                .into_iter()
                .zip(quorum_bitmap_histories)
                .map(|(pk_hash, quorum_bitmap_history)| {
                    let update = Update::new(41, 43, quorum_bitmap_history).unwrap();
                    let history = HashMap::from([(0, update)]);
                    (pk_hash, History(history))
                })
                .collect()
        };

        let operator_stake_history = pk_hashes
            .into_iter()
            .map(|pk_hash| {
                let stake_history_by_quorum = signed_quorum_numbers
                    .clone()
                    .into_iter()
                    .map(|quorum| {
                        let update = Update::new(41, 43, U96::from(10)).unwrap();
                        let history = HashMap::from([(0, update)]);
                        (quorum, History(history))
                    })
                    .collect();
                (pk_hash, stake_history_by_quorum)
            })
            .collect::<HashMap<B256, _>>();

        let total_stake_history = signed_quorum_numbers
            .clone()
            .into_iter()
            .map(|quorum| {
                let update = Update::new(41, 43, U96::from(100)).unwrap();
                let history = HashMap::from([(0, update)]);
                (quorum, History(history))
            })
            .collect();

        let apk_history = signed_quorum_numbers
            .clone()
            .into_iter()
            .zip(apk_for_each_quorum)
            .map(|(quorum, apk)| {
                let apk_hash = convert::point_to_hash(&apk.into());
                let apk_trunc_hash: [u8; 24] = apk_hash[..24].try_into().unwrap();
                let apk_trunc_hash: TruncHash = apk_trunc_hash.into();
                let update = Update::new(41, 43, apk_trunc_hash).unwrap();
                let history = HashMap::from([(0, update)]);
                (quorum, History(history))
            })
            .collect();

        let versioned_blob_params = HashMap::from([(
            42,
            VersionedBlobParams {
                maxNumOperators: 42,
                numChunks: 44,
                codingRate: 42,
            },
        )]);

        let next_blob_version = 43;

        #[cfg(feature = "stale-stakes-forbidden")]
        let staleness = {
            let quorum_update_block_number = signed_quorum_numbers
                .clone()
                .into_iter()
                .map(|quorum| (quorum, 42))
                .collect();

            Staleness {
                stale_stakes_forbidden: true,
                min_withdrawal_delay_blocks: 10,
                quorum_update_block_number,
            }
        };

        let storage = Storage {
            quorum_count: u8::MAX,
            current_block: 43,
            quorum_bitmap_history,
            operator_stake_history,
            total_stake_history,
            apk_history,
            versioned_blob_params,
            next_blob_version,
            #[cfg(feature = "stale-stakes-forbidden")]
            staleness,
        };

        let required_quorum_numbers: Bytes = [0, 2].into();

        CertVerificationInputs {
            batch_header,
            blob_inclusion_info,
            non_signer_stakes_and_signature,
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers,
            storage,
        }
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
        let left_child = keccak256(&encoded);

        let right_sibling = [42u8; 32].into();
        let batch_root = streaming_keccak256(&[left_child, right_sibling]);

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
