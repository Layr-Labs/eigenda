//! EigenDA data extraction from Ethereum contract storage
//!
//! This module provides utilities for extracting and decoding data from EigenDA
//! protocol smart contracts deployed on Ethereum. It enables verification of
//! blob certificates by fetching the necessary on-chain state data.
//!
//! ## Architecture
//!
//! The extraction system follows a trait-based approach:
//! - [`StorageKeyProvider`]: Generates storage keys for contract data
//! - [`DataDecoder`]: Decodes storage proofs into typed data structures
//!
//! ## Key Components
//!
//! - **Extractors**: Specialized types for extracting specific data from contracts
//! - **Contract Interfaces**: High-level interfaces for each EigenDA contract
//! - **Storage Helpers**: Utilities for generating Ethereum storage keys
//! - **Decode Helpers**: Utilities for parsing storage proofs
//!
//! ## Contract Data Extracted
//!
//! - Quorum configurations and counts
//! - Operator stake histories and bitmap histories  
//! - Aggregated public key (APK) histories
//! - Blob versioning parameters
//! - Security thresholds and required quorum numbers
//! - Stale stake prevention settings (feature-gated)

#[cfg(feature = "native")]
pub mod contract;
pub mod decode_helpers;
pub mod storage_key_helpers;
#[cfg(feature = "stale-stakes-forbidden")]
pub use stale_stakes_forbidden::*;

use alloy_primitives::{
    B256, Bytes, StorageKey, U256,
    aliases::{U96, U192},
};
use hashbrown::HashMap;
use reth_trie_common::StorageProof;
use thiserror::Error;
use tracing::instrument;

use crate::eigenda::{
    cert::StandardCommitment,
    verification::cert::{
        bitmap::Bitmap,
        hash::TruncHash,
        types::{
            QuorumNumber, Stake, Version,
            history::{History, HistoryError},
            solidity::{SecurityThresholds, StakeUpdate, VersionedBlobParams},
        },
    },
};

// Storage slot constants for EigenDA contract variables
// These correspond to specific storage slots in the deployed contracts

/// Storage slot for versioned blob parameters mapping in EigenDaThresholdRegistry
const VERSIONED_BLOB_PARAMS_MAPPING_SLOT: u64 = 4;
/// Storage slot for next blob version in EigenDaThresholdRegistry
const NEXT_BLOB_VERSION_SLOT: u64 = 3;
/// Storage slot for quorum count in RegistryCoordinator
const QUORUM_COUNT_VARIABLE_SLOT: u64 = 150;
/// Storage slot for operator bitmap history mapping in RegistryCoordinator
const OPERATOR_BITMAP_HISTORY_MAPPING_SLOT: u64 = 152;
/// Storage slot for APK history mapping in BlsApkRegistry
const APK_HISTORY_MAPPING_SLOT: u64 = 4;
/// Storage slot for total stake history mapping in StakeRegistry
const TOTAL_STAKE_HISTORY_MAPPING_SLOT: u64 = 1;
/// Storage slot for operator stake history mapping in StakeRegistry
const OPERATOR_STAKE_HISTORY_MAPPING_SLOT: u64 = 2;
/// Storage slot for security thresholds V2 in EigenDaCertVerifier
const SECURITY_THRESHOLDS_V2_VARIABLE_SLOT: u64 = 0;
/// Storage slot for required quorum numbers V2 in EigenDaCertVerifier
const QUORUM_NUMBERS_REQUIRED_V2_VARIABLE_SLOT: u64 = 1;

/// Errors that can occur during certificate data extraction
#[derive(Debug, Error, PartialEq)]
pub enum CertExtractionError {
    /// Storage proof was not found for the requested variable
    #[error("Failed to extract StorageProof for {0}")]
    MissingStorageProof(String),

    /// Error from history data processing
    #[error(transparent)]
    WrapHistoryError(#[from] HistoryError),

    /// Error from Alloy Solidity types decoding
    #[error(transparent)]
    WrapAlloySolTypesError(#[from] alloy_sol_types::Error),
}

/// Trait for types that can provide storage keys for data extraction
///
/// This trait is implemented by extractors to specify which storage locations
/// they need to read from Ethereum contracts.
pub trait StorageKeyProvider {
    /// Returns the storage keys that need to be fetched from the blockchain
    fn storage_keys(&self) -> Vec<StorageKey>;
}

/// Trait for types that can decode storage proofs into typed data
///
/// This trait extends [`StorageKeyProvider`] to also handle the decoding of
/// the fetched storage data into application-specific types.
pub trait DataDecoder: StorageKeyProvider {
    /// The type of data this decoder produces
    type Output;

    /// Decode storage proofs into the output type
    ///
    /// # Arguments
    /// * `storage_proofs` - Array of storage proofs from the blockchain
    ///
    /// # Returns
    /// The decoded data of type `Self::Output`
    ///
    /// # Errors
    /// Returns [`CertExtractionError`] if required proofs are missing or decoding fails
    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError>;
}

/// Extractor for the total number of quorums in the registry
///
/// Reads the `quorumCount` variable from the RegistryCoordinator contract.
pub struct QuorumCountExtractor;

impl QuorumCountExtractor {
    /// Create a new quorum count extractor
    ///
    /// # Arguments
    /// * `_certificate` - Certificate (not used but kept for consistency)
    pub fn new(_certificate: &StandardCommitment) -> Self {
        Self {}
    }
}

impl StorageKeyProvider for QuorumCountExtractor {
    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn storage_keys(&self) -> Vec<StorageKey> {
        vec![storage_key_helpers::simple_slot_key(
            QUORUM_COUNT_VARIABLE_SLOT,
        )]
    }
}

impl DataDecoder for QuorumCountExtractor {
    type Output = u8;

    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError> {
        let storage_key = &self.storage_keys()[0];
        let proof =
            decode_helpers::find_required_proof(storage_proofs, storage_key, "quorumCount")?;
        Ok(proof.value.to::<u8>())
    }
}

/// Extractor for versioned blob parameters
///
/// Reads blob configuration parameters for a specific version from the
/// EigenDaThresholdRegistry contract
pub struct VersionedBlobParamsExtractor {
    /// The blob version to extract parameters for
    pub version: u16,
}

impl VersionedBlobParamsExtractor {
    /// Create a new versioned blob parameters extractor
    ///
    /// # Arguments
    /// * `certificate` - Certificate containing the blob version
    pub fn new(certificate: &StandardCommitment) -> Self {
        Self {
            version: certificate.version(),
        }
    }
}

impl StorageKeyProvider for VersionedBlobParamsExtractor {
    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn storage_keys(&self) -> Vec<StorageKey> {
        let version = U256::from(self.version);
        vec![storage_key_helpers::mapping_key(
            version,
            VERSIONED_BLOB_PARAMS_MAPPING_SLOT,
        )]
    }
}

impl DataDecoder for VersionedBlobParamsExtractor {
    type Output = HashMap<Version, VersionedBlobParams>;

    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError> {
        let storage_key = &self.storage_keys()[0];
        let proof = decode_helpers::find_required_proof(
            storage_proofs,
            storage_key,
            "versionedBlobParams",
        )?;
        let key = self.version;
        let le = proof.value.to_le_bytes::<32>();
        let value = VersionedBlobParams {
            maxNumOperators: u32::from_le_bytes(le[0..4].try_into().unwrap()),
            numChunks: u32::from_le_bytes(le[4..8].try_into().unwrap()),
            codingRate: le[8],
        };
        Ok(HashMap::from([(key, value)]))
    }
}

/// Extractor for the next blob version from the threshold registry.
///
/// Reads the `nextBlobVersion` variable from the EigenDaThresholdRegistry contract.
/// This indicates the next version number that will be assigned to blob parameters.
pub struct NextBlobVersionExtractor;

impl NextBlobVersionExtractor {
    /// Create a new next blob version extractor
    ///
    /// # Arguments
    /// * `_certificate` - Certificate (not used but kept for consistency)
    pub fn new(_certificate: &StandardCommitment) -> Self {
        Self
    }
}

impl StorageKeyProvider for NextBlobVersionExtractor {
    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn storage_keys(&self) -> Vec<StorageKey> {
        vec![storage_key_helpers::simple_slot_key(NEXT_BLOB_VERSION_SLOT)]
    }
}

impl DataDecoder for NextBlobVersionExtractor {
    type Output = Version;

    /// Decode the next blob version from storage proofs
    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError> {
        let storage_key = &self.storage_keys()[0];
        let proof =
            decode_helpers::find_required_proof(storage_proofs, storage_key, "nextBlobVersion")?;
        let next_blob_version = proof.value.to::<Self::Output>();

        tracing::info!(?next_blob_version);

        Ok(next_blob_version)
    }
}

/// Extractor for operator bitmap history from the registry coordinator.
///
/// This extractor fetches historical quorum membership data for non-signing operators.
/// The bitmap indicates which quorums each operator was a member of at specific block heights.
/// This information is needed to verify that non-signers were actually part of the required
/// quorums at the time the certificate was created.
pub struct OperatorBitmapHistoryExtractor {
    /// Public key hashes of operators that did not sign the certificate
    pub non_signers_pk_hashes: Vec<B256>,
    /// Indices for looking up bitmap history entries for each non-signer
    pub non_signer_quorum_bitmap_indices: Vec<u32>,
}

impl OperatorBitmapHistoryExtractor {
    /// Create a new operator bitmap history extractor
    ///
    /// # Arguments
    /// * `certificate` - Certificate containing non-signer information
    pub fn new(certificate: &StandardCommitment) -> Self {
        Self {
            non_signers_pk_hashes: certificate.non_signers_pk_hashes(),
            non_signer_quorum_bitmap_indices: certificate
                .non_signer_quorum_bitmap_indices()
                .to_vec(),
        }
    }
}

impl StorageKeyProvider for OperatorBitmapHistoryExtractor {
    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn storage_keys(&self) -> Vec<StorageKey> {
        self.non_signers_pk_hashes
            .iter()
            .zip(self.non_signer_quorum_bitmap_indices.iter())
            .map(|(&operator_id, &index)| {
                storage_key_helpers::dynamic_array_key(
                    operator_id.into(),
                    OPERATOR_BITMAP_HISTORY_MAPPING_SLOT,
                    index,
                )
            })
            .collect()
    }
}

/// Extracts operator bitmap history from RegistryCoordinator::_operatorBitmapHistory.
impl DataDecoder for OperatorBitmapHistoryExtractor {
    type Output = HashMap<B256, History<Bitmap>>;

    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError> {
        self.storage_keys()
            .iter()
            .zip(self.non_signers_pk_hashes.iter())
            .zip(self.non_signer_quorum_bitmap_indices.iter())
            .map(|((&storage_key, &operator_id), &index)| {
                let proof = decode_helpers::find_required_proof(
                    storage_proofs,
                    &storage_key,
                    "_operatorBitmapHistory",
                )?;
                let le = proof.value.to_le_bytes::<32>();
                let update_block = u32::from_le_bytes(le[0..4].try_into().unwrap());
                let next_update_block = u32::from_le_bytes(le[4..8].try_into().unwrap());

                let quorum_bitmap = U192::from_le_bytes::<24>(le[8..32].try_into().unwrap());
                let [lo, mid, hi] = quorum_bitmap.into_limbs();
                let bitmap = Bitmap::new([lo as usize, mid as usize, hi as usize, 0]);

                let update =
                    decode_helpers::create_update(update_block, next_update_block, bitmap)?;
                let history = HashMap::from([(index, update)]);

                Ok((operator_id, History(history)))
            })
            .collect()
    }
}

/// Extractor for aggregate public key (APK) history from the BLS APK registry.
///
/// This extractor fetches the historical aggregate public keys for each quorum that signed
/// the certificate. The APK represents the combined public key of all operators in a quorum
/// at a specific block height, which is essential for verifying BLS aggregate signatures.
pub struct ApkHistoryExtractor {
    /// Numbers of quorums that signed the certificate
    pub signed_quorum_numbers: Bytes,
    /// Indices for looking up APK history entries for each quorum
    pub quorum_apk_indices: Vec<u32>,
}

impl ApkHistoryExtractor {
    /// Create a new APK history extractor
    ///
    /// # Arguments
    /// * `certificate` - Certificate containing signed quorum information
    pub fn new(certificate: &StandardCommitment) -> Self {
        Self {
            signed_quorum_numbers: certificate.signed_quorum_numbers().clone(),
            quorum_apk_indices: certificate.quorum_apk_indices().to_vec(),
        }
    }
}

impl StorageKeyProvider for ApkHistoryExtractor {
    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn storage_keys(&self) -> Vec<StorageKey> {
        self.signed_quorum_numbers
            .iter()
            .zip(self.quorum_apk_indices.iter())
            .map(|(&signed_quorum_number, &index)| {
                storage_key_helpers::dynamic_array_key(
                    U256::from(signed_quorum_number),
                    APK_HISTORY_MAPPING_SLOT,
                    index,
                )
            })
            .collect()
    }
}

/// Extracts APK history from BlsApkRegistry::apkHistory.
/// Contains the aggregate public keys for each quorum at different block heights.
impl DataDecoder for ApkHistoryExtractor {
    type Output = HashMap<QuorumNumber, History<TruncHash>>;

    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError> {
        self.storage_keys()
            .iter()
            .zip(self.signed_quorum_numbers.iter())
            .zip(self.quorum_apk_indices.iter())
            .map(|((&storage_key, &signed_quorum_number), &index)| {
                let proof = decode_helpers::find_required_proof(
                    storage_proofs,
                    &storage_key,
                    "apkHistory",
                )?;
                let le = proof.value.to_le_bytes::<32>();

                let mut apk_trunc_hash_bytes: [u8; 24] = le[..24].try_into().unwrap();
                apk_trunc_hash_bytes.reverse();

                let apk_trunc_hash: TruncHash = apk_trunc_hash_bytes.into();
                let update_block = u32::from_le_bytes(le[24..28].try_into().unwrap());
                let next_update_block = u32::from_le_bytes(le[28..32].try_into().unwrap());

                let update =
                    decode_helpers::create_update(update_block, next_update_block, apk_trunc_hash)?;
                let history = HashMap::from([(index, update)]);

                Ok((signed_quorum_number, History(history)))
            })
            .collect()
    }
}

/// Extractor for total stake history from the stake registry.
///
/// This extractor fetches the historical total stake amounts for each quorum at specific
/// block heights. The total stake is used to calculate voting thresholds and determine
/// whether sufficient stake participated in signing the certificate.
pub struct TotalStakeHistoryExtractor {
    /// Numbers of quorums that signed the certificate
    pub signed_quorum_numbers: Bytes,
    /// Indices for looking up total stake history entries
    pub non_signer_total_stake_indices: Vec<u32>,
}

impl TotalStakeHistoryExtractor {
    /// Create a new total stake history extractor
    ///
    /// # Arguments
    /// * `certificate` - Certificate containing quorum and stake index information
    pub fn new(certificate: &StandardCommitment) -> Self {
        Self {
            signed_quorum_numbers: certificate.signed_quorum_numbers().clone(),
            non_signer_total_stake_indices: certificate.non_signer_total_stake_indices().to_vec(),
        }
    }
}

impl StorageKeyProvider for TotalStakeHistoryExtractor {
    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn storage_keys(&self) -> Vec<StorageKey> {
        self.signed_quorum_numbers
            .iter()
            .zip(self.non_signer_total_stake_indices.iter())
            .map(|(&signed_quorum_number, &index)| {
                storage_key_helpers::dynamic_array_key(
                    U256::from(signed_quorum_number),
                    TOTAL_STAKE_HISTORY_MAPPING_SLOT,
                    index,
                )
            })
            .collect()
    }
}

/// Extracts total stake history from StakeRegistry::_totalStakeHistory.
/// This is used by getTotalStakeAtBlockNumberFromIndex for stake calculations.
impl DataDecoder for TotalStakeHistoryExtractor {
    type Output = HashMap<QuorumNumber, History<Stake>>;

    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError> {
        self.storage_keys()
            .iter()
            .zip(self.signed_quorum_numbers.iter())
            .zip(self.non_signer_total_stake_indices.iter())
            .map(|((&storage_key, &signed_quorum_number), &index)| {
                let proof = decode_helpers::find_required_proof(
                    storage_proofs,
                    &storage_key,
                    "_totalStakeHistory",
                )?;
                let le = proof.value.to_le_bytes::<32>();
                let stake_update = StakeUpdate {
                    updateBlockNumber: u32::from_le_bytes(le[0..4].try_into().unwrap()),
                    nextUpdateBlockNumber: u32::from_le_bytes(le[4..8].try_into().unwrap()),
                    stake: U96::from_le_bytes::<12>(le[8..20].try_into().unwrap()),
                };

                let stake = stake_update.stake.to::<U96>();
                let update = decode_helpers::create_update(
                    stake_update.updateBlockNumber,
                    stake_update.nextUpdateBlockNumber,
                    stake,
                )?;

                let history = HashMap::from([(index, update)]);
                Ok((signed_quorum_number, History(history)))
            })
            .collect()
    }
}

/// Extractor for individual operator stake history from the stake registry.
///
/// This extractor fetches the historical stake amounts for individual operators
/// who did not sign the certificate. This data is needed to calculate the exact
/// stake distribution and verify that non-signers' stakes are properly accounted
/// for in the threshold calculations.
pub struct OperatorStakeHistoryExtractor {
    /// Numbers of quorums that signed the certificate
    pub signed_quorum_numbers: Bytes,
    /// Public key hashes of operators that did not sign
    pub non_signers_pk_hashes: Vec<B256>,
    /// Nested indices for looking up stake history (per quorum, per operator)
    pub non_signer_stake_indices: Vec<Vec<u32>>,
}

impl OperatorStakeHistoryExtractor {
    /// Create a new operator stake history extractor
    ///
    /// # Arguments
    /// * `certificate` - Certificate containing non-signer and stake index information
    pub fn new(certificate: &StandardCommitment) -> Self {
        Self {
            signed_quorum_numbers: certificate.signed_quorum_numbers().clone(),
            non_signers_pk_hashes: certificate.non_signers_pk_hashes(),
            non_signer_stake_indices: certificate.non_signer_stake_indices().to_vec(),
        }
    }
}

impl StorageKeyProvider for OperatorStakeHistoryExtractor {
    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn storage_keys(&self) -> Vec<StorageKey> {
        let mut storage_keys = vec![];

        for (&signed_quorum_number, stake_index_for_each_required_non_signer) in self
            .signed_quorum_numbers
            .iter()
            .zip(&self.non_signer_stake_indices)
        {
            for &operator_id in &self.non_signers_pk_hashes {
                // without peeking at the actual data it's impossible to associate indices with
                // any one non_signer so it's necessary to do this cartesian product. Storage keys
                // that map to non-existent data will return empty but won't fail. When retrieved
                // an empty value will be returned for inexisting storage keys
                for &stake_index in stake_index_for_each_required_non_signer {
                    let storage_key = storage_key_helpers::nested_dynamic_array_key(
                        operator_id.into(),
                        OPERATOR_STAKE_HISTORY_MAPPING_SLOT,
                        U256::from(signed_quorum_number),
                        stake_index,
                    );
                    storage_keys.push(storage_key);
                }
            }
        }

        storage_keys
    }
}

/// Extracts operator stake history from StakeRegistry::operatorStakeHistory.
impl DataDecoder for OperatorStakeHistoryExtractor {
    type Output = HashMap<B256, HashMap<QuorumNumber, History<Stake>>>;

    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError> {
        let mut out: HashMap<B256, HashMap<QuorumNumber, History<Stake>>> = HashMap::new();

        for (&signed_quorum_number, stake_index_for_each_required_non_signer) in self
            .signed_quorum_numbers
            .iter()
            .zip(&self.non_signer_stake_indices)
        {
            for &operator_id in &self.non_signers_pk_hashes {
                // Same cartesian product is necessary as for the StorageKeyProvider impl
                for &stake_index in stake_index_for_each_required_non_signer {
                    let storage_key = storage_key_helpers::nested_dynamic_array_key(
                        operator_id.into(),
                        OPERATOR_STAKE_HISTORY_MAPPING_SLOT,
                        U256::from(signed_quorum_number),
                        stake_index,
                    );

                    let proof = decode_helpers::find_required_proof(
                        storage_proofs,
                        &storage_key,
                        "operatorStakeHistory",
                    )?;
                    let le = proof.value.to_le_bytes::<32>();
                    let stake_update = StakeUpdate {
                        updateBlockNumber: u32::from_le_bytes(le[0..4].try_into().unwrap()),
                        nextUpdateBlockNumber: u32::from_le_bytes(le[4..8].try_into().unwrap()),
                        stake: U96::from_le_bytes::<12>(le[8..20].try_into().unwrap()),
                    };

                    let stake = stake_update.stake.to::<U96>();
                    let update = decode_helpers::create_update(
                        stake_update.updateBlockNumber,
                        stake_update.nextUpdateBlockNumber,
                        stake,
                    )?;

                    let operator_id: B256 = operator_id;

                    out.entry(operator_id)
                        .or_default()
                        .entry(signed_quorum_number)
                        .or_insert_with(|| History(HashMap::new()))
                        .0
                        .insert(stake_index, update);
                }
            }
        }

        Ok(out)
    }
}

/// Extractor for security thresholds from the certificate verifier.
///
/// This extractor fetches the security threshold parameters that define the minimum
/// requirements for certificate validation, including confirmation and adversary thresholds
/// that determine the minimum stake percentages needed for valid signatures.
pub struct SecurityThresholdsV2Extractor;

impl SecurityThresholdsV2Extractor {
    /// Create a new security thresholds extractor
    ///
    /// # Arguments
    /// * `_certificate` - Certificate (not used but kept for consistency)
    pub fn new(_certificate: &StandardCommitment) -> Self {
        Self {}
    }
}

impl StorageKeyProvider for SecurityThresholdsV2Extractor {
    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn storage_keys(&self) -> Vec<StorageKey> {
        vec![storage_key_helpers::simple_slot_key(
            SECURITY_THRESHOLDS_V2_VARIABLE_SLOT,
        )]
    }
}

/// Extracts security thresholds from EigenDaCertVerifier::securityThresholdsV2.
/// Example on Holesky: confirmationThreshold=55%, adversaryThreshold=33%.
impl DataDecoder for SecurityThresholdsV2Extractor {
    type Output = SecurityThresholds;

    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError> {
        let storage_key = &self.storage_keys()[0];
        let proof = decode_helpers::find_required_proof(
            storage_proofs,
            storage_key,
            "securityThresholdsV2",
        )?;

        let le = proof.value.to_le_bytes::<32>();

        let security_thresholds = SecurityThresholds {
            confirmationThreshold: le[0],
            adversaryThreshold: le[1],
        };

        Ok(security_thresholds)
    }
}

/// Extractor for required quorum numbers from the certificate verifier.
///
/// This extractor fetches the list of quorum numbers that are required to participate
/// in certificate signing for the certificate to be considered valid. This defines
/// which quorums must have sufficient stake participation.
pub struct QuorumNumbersRequiredV2Extractor;

impl QuorumNumbersRequiredV2Extractor {
    /// Create a new required quorum numbers extractor
    ///
    /// # Arguments
    /// * `_certificate` - Certificate (not used but kept for consistency)
    pub fn new(_certificate: &StandardCommitment) -> Self {
        Self {}
    }
}

impl StorageKeyProvider for QuorumNumbersRequiredV2Extractor {
    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn storage_keys(&self) -> Vec<StorageKey> {
        vec![storage_key_helpers::simple_slot_key(
            QUORUM_NUMBERS_REQUIRED_V2_VARIABLE_SLOT,
        )]
    }
}

/// Extracts required quorum numbers from EigenDaCertVerifier::quorumNumbersRequiredV2.
/// Example on Holesky: 0x0001 (indicating quorum 0 and 1 are required).
impl DataDecoder for QuorumNumbersRequiredV2Extractor {
    type Output = Bytes;

    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError> {
        let storage_key = &self.storage_keys()[0];
        let proof = decode_helpers::find_required_proof(
            storage_proofs,
            storage_key,
            "quorumNumbersRequiredV2",
        )?;

        // there can be at most 256 quorums
        let be = proof.value.to_be_bytes::<32>();
        let len = (be[31] / 2) as usize;
        Ok(be[..len].to_vec().into())
    }
}

#[cfg(feature = "stale-stakes-forbidden")]
mod stale_stakes_forbidden {
    use tracing::instrument;

    use super::*;
    use crate::eigenda::verification::cert::types::BlockNumber;

    const QUORUM_UPDATE_BLOCK_NUMBER_MAPPING_SLOT: u64 = 155;
    const STALE_STAKES_FORBIDDEN_VARIABLE_SLOT: u64 = 201;
    const MIN_WITHDRAWAL_DELAY_BLOCKS_VARIABLE_SLOT: u64 = 157;

    /// Extractor for the stale stakes forbidden flag from the service manager.
    ///
    /// This extractor determines whether stale stakes are forbidden in the current
    /// configuration. When enabled, this prevents operators from using outdated
    /// stake information for validation.
    pub struct StaleStakesForbiddenExtractor;

    impl StaleStakesForbiddenExtractor {
        /// Create a new stale stakes forbidden extractor
        ///
        /// # Arguments
        /// * `_certificate` - Certificate (not used but kept for consistency)
        pub fn new(_certificate: &StandardCommitment) -> Self {
            Self {}
        }
    }

    impl StorageKeyProvider for StaleStakesForbiddenExtractor {
        #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
        fn storage_keys(&self) -> Vec<StorageKey> {
            vec![storage_key_helpers::simple_slot_key(
                STALE_STAKES_FORBIDDEN_VARIABLE_SLOT,
            )]
        }
    }

    /// Extracts stale stakes flag from EigenDAServiceManager::staleStakesForbidden.
    /// Example on Holesky: false (stale stakes are allowed).
    impl DataDecoder for StaleStakesForbiddenExtractor {
        type Output = bool;

        #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
        fn decode_data(
            &self,
            storage_proofs: &[StorageProof],
        ) -> Result<Self::Output, CertExtractionError> {
            let storage_key = &self.storage_keys()[0];
            let proof = decode_helpers::find_required_proof(
                storage_proofs,
                storage_key,
                "staleStakesForbidden",
            )?;
            Ok(!proof.value.is_zero())
        }
    }

    /// Extractor for minimum withdrawal delay blocks from the delegation manager.
    ///
    /// This extractor fetches the minimum number of blocks that must pass before
    /// stake withdrawals can be completed. This delay is a security mechanism
    /// to prevent rapid stake changes that could affect validation integrity.
    pub struct MinWithdrawalDelayBlocksExtractor;

    impl MinWithdrawalDelayBlocksExtractor {
        /// Create a new minimum withdrawal delay blocks extractor
        ///
        /// # Arguments
        /// * `_certificate` - Certificate (not used but kept for consistency)
        pub fn new(_certificate: &StandardCommitment) -> Self {
            Self {}
        }
    }

    impl StorageKeyProvider for MinWithdrawalDelayBlocksExtractor {
        #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
        fn storage_keys(&self) -> Vec<StorageKey> {
            vec![storage_key_helpers::simple_slot_key(
                MIN_WITHDRAWAL_DELAY_BLOCKS_VARIABLE_SLOT,
            )]
        }
    }

    /// Extracts minimum withdrawal delay from DelegationManager::minWithdrawalDelayBlocks.
    /// Defines the security delay period for stake withdrawals.
    impl DataDecoder for MinWithdrawalDelayBlocksExtractor {
        type Output = u32;

        #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
        fn decode_data(
            &self,
            storage_proofs: &[StorageProof],
        ) -> Result<Self::Output, CertExtractionError> {
            let storage_key = &self.storage_keys()[0];
            let proof = decode_helpers::find_required_proof(
                storage_proofs,
                storage_key,
                "minWithdrawalDelayBlocks",
            )?;
            Ok(proof.value.to::<Self::Output>())
        }
    }

    /// Extractor for quorum update block numbers from the registry coordinator.
    ///
    /// This extractor fetches the block numbers when each quorum was last updated.
    /// This information is used in conjunction with stale stakes prevention to ensure
    /// that stake information is sufficiently recent for validation purposes.
    pub struct QuorumUpdateBlockNumberExtractor {
        /// Numbers of quorums that signed the certificate
        pub signed_quorum_numbers: Bytes,
    }

    impl QuorumUpdateBlockNumberExtractor {
        /// Create a new quorum update block number extractor
        ///
        /// # Arguments
        /// * `certificate` - Certificate containing signed quorum information
        pub fn new(certificate: &StandardCommitment) -> Self {
            Self {
                signed_quorum_numbers: certificate.signed_quorum_numbers().clone(),
            }
        }
    }

    impl StorageKeyProvider for QuorumUpdateBlockNumberExtractor {
        #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
        fn storage_keys(&self) -> Vec<StorageKey> {
            self.signed_quorum_numbers
                .iter()
                .map(|&quorum_number| {
                    storage_key_helpers::mapping_key(
                        U256::from(quorum_number),
                        QUORUM_UPDATE_BLOCK_NUMBER_MAPPING_SLOT,
                    )
                })
                .collect()
        }
    }

    /// Extracts quorum update blocks from RegistryCoordinator::quorumUpdateBlockNumber.
    /// Tracks when each quorum configuration was last modified.
    impl DataDecoder for QuorumUpdateBlockNumberExtractor {
        type Output = HashMap<QuorumNumber, BlockNumber>;

        #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
        fn decode_data(
            &self,
            storage_proofs: &[StorageProof],
        ) -> Result<Self::Output, CertExtractionError> {
            self.storage_keys()
                .iter()
                .zip(self.signed_quorum_numbers.iter())
                .map(|(storage_key, &quorum_number)| {
                    decode_helpers::find_required_proof(
                        storage_proofs,
                        storage_key,
                        "quorumUpdateBlockNumber",
                    )
                    .map(|proof| (quorum_number, proof.value.to::<BlockNumber>()))
                })
                .collect()
        }
    }
}
