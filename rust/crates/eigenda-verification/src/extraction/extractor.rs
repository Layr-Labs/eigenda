use alloy_primitives::aliases::{U96, U192};
use alloy_primitives::{B256, Bytes, StorageKey, U256};
use hashbrown::HashMap;
use reth_trie_common::StorageProof;
pub use stale_stakes_forbidden::*;
use tracing::instrument;

use crate::cert::StandardCommitment;
use crate::cert::solidity::{SecurityThresholds, StakeUpdate, VersionedBlobParams};
use crate::extraction::{CertExtractionError, decode_helpers, storage_key_helpers};
use crate::verification::cert::bitmap::Bitmap;
use crate::verification::cert::hash::TruncHash;
use crate::verification::cert::types::history::History;
use crate::verification::cert::types::{QuorumNumber, Stake, Version};

// Storage slot constants for EigenDA contract variables
// These correspond to specific storage slots in the deployed contracts
// These can be verified by running for example `forge inspect RegistryCoordinator storageLayout`
// from the contracts subdir.
// TODO(samlaf): we need to make sure these are kept in sync with the deployed contracts!
// Prob want to automate this in CI somehow.

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

/// Storage slot for certificate verifiers address mapping in EigenDaCertVerifierRouter
const CERT_VERIFIERS_ADDRESS_MAPPING_SLOT: u64 = 101;

/// Storage slot for certificate verifiers ABNs array in EigenDaCertVerifierRouter
pub const CERT_VERIFIER_ABNS_ARRAY_SLOT: u64 = 102;

/// Storage slot for security thresholds V2 in EigenDaCertVerifier
const SECURITY_THRESHOLDS_V2_VARIABLE_SLOT: u64 = 0;

/// Storage slot for required quorum numbers V2 in EigenDaCertVerifier
const QUORUM_NUMBERS_REQUIRED_V2_VARIABLE_SLOT: u64 = 1;

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
#[derive(Default)]
pub struct QuorumCountExtractor;

impl QuorumCountExtractor {
    /// Create a new quorum count extractor
    pub fn new() -> Self {
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
#[derive(Default)]
pub struct NextBlobVersionExtractor;

impl NextBlobVersionExtractor {
    /// Create a new next blob version extractor
    pub fn new() -> Self {
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
                storage_key_helpers::mapping_to_dynamic_array_key(
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
                storage_key_helpers::mapping_to_dynamic_array_key(
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
                storage_key_helpers::mapping_to_dynamic_array_key(
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
                    let storage_key = storage_key_helpers::nested_mapping_to_dynamic_array_key(
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
                    let storage_key = storage_key_helpers::nested_mapping_to_dynamic_array_key(
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

/// Extractor for the length of the certificate verifiers ABNs array.
///
/// This extractor is used to determine how many certificate verifiers are registered.
/// It is needed to prove an ABN is currently active in case that ABN is the last
/// registered in the contract.
pub struct CertVerifierABNsLenExtractor;

impl CertVerifierABNsLenExtractor {
    /// Create a new certificate verifier ABNs length extractor
    pub fn new() -> Self {
        Self {}
    }
}

impl Default for CertVerifierABNsLenExtractor {
    /// Create a default instance of the extractor
    fn default() -> Self {
        Self::new()
    }
}

impl StorageKeyProvider for CertVerifierABNsLenExtractor {
    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn storage_keys(&self) -> Vec<StorageKey> {
        vec![storage_key_helpers::simple_slot_key(
            CERT_VERIFIER_ABNS_ARRAY_SLOT,
        )]
    }
}

/// Extractor for the certificate verifiers ABNs array.
/// This struct is used to extract a specified number of certificate verifier ABNs from storage.
pub struct CertVerifierABNsExtractor {
    /// The number of certificate verifier ABNs to extract.
    /// Should be fetched using CertVerifierABNsLenExtractor beforehand,
    /// to make sure all ABNs are retrieved.
    pub num_abns: usize,
}

impl CertVerifierABNsExtractor {
    /// Create a new certificate verifier ABNs extractor
    ///
    /// # Arguments
    /// * `num_abns` - Number of ABNs to extract from storage
    pub fn new(num_abns: usize) -> Self {
        Self { num_abns }
    }
}

impl StorageKeyProvider for CertVerifierABNsExtractor {
    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn storage_keys(&self) -> Vec<StorageKey> {
        let keys: Vec<_> = (0..self.num_abns)
            .map(|i| {
                storage_key_helpers::dynamic_array_key(CERT_VERIFIER_ABNS_ARRAY_SLOT, i as u32)
            })
            .collect();
        keys
    }
}

/// Extracts cert verifier information for a given set of ABNs (activation block numbers).
///
/// This struct is used to retrieve the cert verifiers associated with the provided ABNs.
/// ABNs are required to identify which cert verifiers' data should be extracted from storage,
/// as each ABN corresponds to a specific cert verifier in the mapping.
pub struct CertVerifiersExtractor<'a> {
    /// The list of ABNs (activation block numbers) for which cert verifiers are to be extracted.
    ///
    /// Each ABN uniquely identifies a cert verifier in the contract's storage mapping.
    pub abns: &'a [u32],
}

impl<'a> CertVerifiersExtractor<'a> {
    /// Create a new cert verifiers extractor
    ///
    /// # Arguments
    /// * `abns` - Slice of activation block numbers identifying cert verifiers
    pub fn new(abns: &'a [u32]) -> Self {
        Self { abns }
    }
}

impl<'a> StorageKeyProvider for CertVerifiersExtractor<'a> {
    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn storage_keys(&self) -> Vec<StorageKey> {
        let keys: Vec<_> = self
            .abns
            .iter()
            .map(|abn| {
                storage_key_helpers::mapping_key(
                    U256::from(*abn),
                    CERT_VERIFIERS_ADDRESS_MAPPING_SLOT,
                )
            })
            .collect();
        keys
    }
}

/// Extractor for security thresholds from the certificate verifier.
///
/// This extractor fetches the security threshold parameters that define the minimum
/// requirements for certificate validation, including confirmation and adversary thresholds
/// that determine the minimum stake percentages needed for valid signatures.
#[derive(Default)]
pub struct SecurityThresholdsV2Extractor;

impl SecurityThresholdsV2Extractor {
    /// Create a new security thresholds extractor
    pub fn new() -> Self {
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
#[derive(Default)]
pub struct QuorumNumbersRequiredV2Extractor;

impl QuorumNumbersRequiredV2Extractor {
    /// Create a new required quorum numbers extractor
    pub fn new() -> Self {
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
/// Example on Holesky: 0x0001 (indicating quorum 0 is required).
impl DataDecoder for QuorumNumbersRequiredV2Extractor {
    type Output = Bytes;

    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError> {
        use CertExtractionError::*;

        let storage_key = &self.storage_keys()[0];
        let proof = decode_helpers::find_required_proof(
            storage_proofs,
            storage_key,
            "quorumNumbersRequiredV2",
        )?;

        // By design there can be at most 256 quorums (meaning this value occupies only 8 bytes)
        // Quorum numbers are stored as the Ethereum "bytes" type
        // Ethereum encodes bytes (like strings) of length < 32 (called "short form", our case) as follows:
        //   The actual bytes are stored left-aligned (i.e., starting at the most significant byte).
        //   The LSB either stores length * 2 (short form) or length * 2 + 1 (long form)
        //     So the LSB serves a dual purpose:
        //       - its parity indicates whether we're dealing with short or long form
        //       - it also stores the length of the payload
        let be = proof.value.to_be_bytes::<32>();

        let is_long_form = (be[31] & 1) == 1;
        if is_long_form {
            return Err(UnexpectedEthereumBytesLongForm);
        }

        // Since the LSB stores (len << 1) we can recover the length with just LSB >> 1
        let len = (be[31] >> 1) as usize;
        Ok(be[..len].to_vec().into())
    }
}

#[cfg(test)]
mod tests {
    use alloy_primitives::U256;
    use reth_trie_common::StorageProof;

    use super::*;

    fn create_mock_certificate() -> StandardCommitment {
        let commitment_hex = "02f90389e5a0c769488dd5264b3ef21dce7ee2d42fba43e1f83ff228f501223e38818cb14492833f44fcf901eff901caf9018180820001f90159f842a0012e810ffc0a83074b3d14db9e78bbae623f7770cac248df9e73fac6b9d59d17a02a916ffbbf9dde4b7ebe94191a29ff686422d7dcb3b47ecb03c6ada75a9c15c8f888f842a01811c8b4152fce9b8c4bae61a3d097e61dfc43dc7d45363d19e7c7f1374034ffa001edc62174217cdce60a4b52fa234ac0d96db4307dac9150e152ba82cbb4d2f1f842a00f423b0dbc1fe95d2e3f7dbac6c099e51dbf73400a4b3f26b9a29665b4ac58a8a01855a2bd56c0e8f4cc85ac149cf9a531673d0e89e22f0d6c4ae419ed7c5d2940f888f842a02667cbb99d60fa0d7f3544141d3d531dceeeb50b06e5a0cdc42338a359138ae4a00dff4c929d8f8a307c19bba6e8006fe6700f6554cef9eb3797944f89472ffb30f842a004c17a6225acd5b4e7d672a1eb298c5358f4f6f17d04fd1ee295d0c0d372fa84a024bc3ad4d5e54f54f71db382ce276f37ac3c260cc74306b832e8a3c93c7951d302a0e43e11e2405c2fd1d880af8612d969b654827e0ba23d9feb3722ccce6226fce7b8411ddf4553c79c0515516fd3c8b3ae6a756b05723f4d0ebe98a450c8bcc96cbb355ef07a44eeb56f831be73647e4da20e22fa859f984ee41d6efcd3692063b0b0601c2800101a0a69e552a6fc2ff75d32edaf5313642ddeebe60d2069435d12e266ce800e9e96bf9016bc0c0f888f842a00d45727a99053af8d38d4716ab83ace676096e7506b6b7aa6953e87bc04a023ca016c030c31dd1c94062948ecdce2e67c4e6626c16af0033dcdb7a96362c937d48f842a00a95fac74aba7e3fbd24bc62457ce6981803d8f5fef28871d3d5e2af05d50cd4a0117400693917cd50d9bc28d4ab4fadf93a23e771f303637f8d1f83cd0632c3fcf888f842a0301bfced3253e99e8d50f2fed62313a16d714013d022a4dc4294656276f10d1ba0152e047a83c326a9d81dac502ec429b662b58ee119ca4c8748a355b539c24131f842a01944b5b4a3e93d46b0fe4370128c6cdcd066ae6b036b019a20f8d22fe9a10d67a00ddf3421722967c0bd965b9fc9e004bf01183b6206fec8de65e40331d185372ef842a02db8fb278708abf8878ebf578872ab35ee914ad8196b78de16b34498222ac1c2a02ff9d9a5184684f4e14530bde3a61a2f9adaa74734dff104b61ba3d963a644dac68207388208b7c68209998209c5c2c0c0820001";
        let raw_commitment = hex::decode(commitment_hex).unwrap();
        StandardCommitment::from_rlp_bytes(raw_commitment.as_slice()).unwrap()
    }

    fn create_storage_proof(key: StorageKey, value: U256) -> StorageProof {
        StorageProof {
            key,
            value,
            ..Default::default()
        }
    }

    #[test]
    fn quorum_count_extractor() {
        let extractor = QuorumCountExtractor::new();

        let keys = extractor.storage_keys();
        assert_eq!(
            keys[0],
            storage_key_helpers::simple_slot_key(QUORUM_COUNT_VARIABLE_SLOT)
        );

        let storage_key = keys[0];
        let proof = create_storage_proof(storage_key, U256::from(5u8));
        let proofs = vec![proof];
        let result = extractor.decode_data(&proofs).unwrap();
        assert_eq!(result, 5u8);

        let empty_proofs = vec![];
        let err = extractor.decode_data(&empty_proofs).unwrap_err();
        assert!(matches!(err, CertExtractionError::MissingStorageProof(_)));
    }

    #[test]
    fn versioned_blob_params_extractor() {
        let cert = create_mock_certificate();
        let extractor = VersionedBlobParamsExtractor::new(&cert);

        let keys = extractor.storage_keys();
        assert_eq!(keys.len(), 1);
        let expected_key = storage_key_helpers::mapping_key(
            U256::from(cert.version()),
            VERSIONED_BLOB_PARAMS_MAPPING_SLOT,
        );
        assert_eq!(keys[0], expected_key);

        let storage_key = keys[0];
        let mut value_bytes = [0u8; 32];
        value_bytes[0..4].copy_from_slice(&100u32.to_le_bytes());
        value_bytes[4..8].copy_from_slice(&50u32.to_le_bytes());
        value_bytes[8] = 80u8;
        let value = U256::from_le_bytes(value_bytes);

        let proof = create_storage_proof(storage_key, value);
        let proofs = vec![proof];
        let result = extractor.decode_data(&proofs).unwrap();
        let version = cert.version();
        let params = result.get(&version).unwrap();

        assert_eq!(params.maxNumOperators, 100);
        assert_eq!(params.numChunks, 50);
        assert_eq!(params.codingRate, 80);
    }

    #[test]
    fn next_blob_version_extractor() {
        let extractor = NextBlobVersionExtractor::new();

        let keys = extractor.storage_keys();
        assert_eq!(keys.len(), 1);
        assert_eq!(
            keys[0],
            storage_key_helpers::simple_slot_key(NEXT_BLOB_VERSION_SLOT)
        );

        let storage_key = keys[0];
        let proof = create_storage_proof(storage_key, U256::from(42u16));
        let proofs = vec![proof];
        let result = extractor.decode_data(&proofs).unwrap();
        assert_eq!(result, 42u16);
    }

    #[test]
    fn security_thresholds_v2_extractor() {
        let extractor = SecurityThresholdsV2Extractor::new();

        let keys = extractor.storage_keys();
        assert_eq!(keys.len(), 1);
        assert_eq!(
            keys[0],
            storage_key_helpers::simple_slot_key(SECURITY_THRESHOLDS_V2_VARIABLE_SLOT)
        );

        let storage_key = keys[0];
        let mut value_bytes = [0u8; 32];
        value_bytes[0] = 55u8;
        value_bytes[1] = 33u8;
        let value = U256::from_le_bytes(value_bytes);

        let proof = create_storage_proof(storage_key, value);
        let proofs = vec![proof];
        let result = extractor.decode_data(&proofs).unwrap();
        assert_eq!(result.confirmationThreshold, 55);
        assert_eq!(result.adversaryThreshold, 33);
    }

    #[test]
    fn quorum_numbers_required_v2_extractor() {
        let extractor = QuorumNumbersRequiredV2Extractor::new();

        let keys = extractor.storage_keys();
        assert_eq!(keys.len(), 1);
        assert_eq!(
            keys[0],
            storage_key_helpers::simple_slot_key(QUORUM_NUMBERS_REQUIRED_V2_VARIABLE_SLOT)
        );

        let storage_key = keys[0];
        let mut value_bytes = [0u8; 32];
        value_bytes[0] = 0u8; // quorum 0
        value_bytes[1] = 1u8; // quorum 1  
        value_bytes[31] = 4u8; // length = 2, encoded as (length * 2)
        let value = U256::from_be_bytes(value_bytes);

        let proof = create_storage_proof(storage_key, value);
        let proofs = vec![proof];
        let result = extractor.decode_data(&proofs).unwrap();
        assert_eq!(result.len(), 2);
        assert_eq!(result[0], 0u8);
        assert_eq!(result[1], 1u8);
    }

    #[test]
    fn operator_bitmap_history_extractor() {
        let cert = create_mock_certificate();
        let extractor = OperatorBitmapHistoryExtractor::new(&cert);

        let keys = extractor.storage_keys();
        assert_eq!(keys.len(), cert.non_signers_pk_hashes().len());

        let proofs = vec![];
        let result = extractor.decode_data(&proofs).unwrap();
        assert!(result.is_empty());
    }

    #[test]
    fn apk_history_extractor() {
        let cert = create_mock_certificate();
        let extractor = ApkHistoryExtractor::new(&cert);

        let keys = extractor.storage_keys();
        assert_eq!(keys.len(), cert.signed_quorum_numbers().len());

        let proofs = vec![];
        let err = extractor.decode_data(&proofs).unwrap_err();
        assert!(matches!(err, CertExtractionError::MissingStorageProof(_)));
    }

    #[test]
    fn total_stake_history_extractor() {
        let cert = create_mock_certificate();
        let extractor = TotalStakeHistoryExtractor::new(&cert);

        let keys = extractor.storage_keys();
        assert_eq!(keys.len(), cert.signed_quorum_numbers().len());

        let proofs = vec![];
        let err = extractor.decode_data(&proofs).unwrap_err();
        assert!(matches!(err, CertExtractionError::MissingStorageProof(_)));
    }

    #[test]
    fn operator_stake_history_extractor() {
        let cert = create_mock_certificate();
        let extractor = OperatorStakeHistoryExtractor::new(&cert);

        let keys = extractor.storage_keys();
        let expected_len = cert.signed_quorum_numbers().len()
            * cert.non_signers_pk_hashes().len()
            * cert
                .non_signer_stake_indices()
                .iter()
                .map(|v| v.len())
                .sum::<usize>();
        assert_eq!(keys.len(), expected_len);

        let proofs = vec![];
        let result = extractor.decode_data(&proofs).unwrap();
        assert!(result.is_empty());
    }

    #[test]
    fn cert_verifier_abns_len_extractor() {
        let extractor = CertVerifierABNsLenExtractor::new();

        let keys = extractor.storage_keys();
        assert_eq!(keys.len(), 1);
    }

    #[test]
    fn cert_verifier_abns_extractor() {
        let extractor = CertVerifierABNsExtractor::new(3);

        let keys = extractor.storage_keys();
        assert_eq!(keys.len(), 3);
    }

    #[test]
    fn cert_verifiers_extractor() {
        let abns = vec![1u32, 2u32, 3u32];
        let extractor = CertVerifiersExtractor::new(&abns);

        let keys = extractor.storage_keys();
        assert_eq!(keys.len(), abns.len());
    }
}

mod stale_stakes_forbidden {
    use tracing::instrument;

    use super::*;
    use crate::verification::cert::types::BlockNumber;

    const QUORUM_UPDATE_BLOCK_NUMBER_MAPPING_SLOT: u64 = 155;
    const STALE_STAKES_FORBIDDEN_VARIABLE_SLOT: u64 = 201;
    const MIN_WITHDRAWAL_DELAY_BLOCKS_VARIABLE_SLOT: u64 = 157;

    /// Extractor for the stale stakes forbidden flag from the service manager.
    ///
    /// This extractor determines whether stale stakes are forbidden in the current
    /// configuration. When enabled, this prevents operators from using outdated
    /// stake information for validation.
    #[derive(Default)]
    pub struct StaleStakesForbiddenExtractor;

    impl StaleStakesForbiddenExtractor {
        /// Create a new stale stakes forbidden extractor
        pub fn new() -> Self {
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
    #[derive(Default)]
    pub struct MinWithdrawalDelayBlocksExtractor;

    impl MinWithdrawalDelayBlocksExtractor {
        /// Create a new minimum withdrawal delay blocks extractor
        pub fn new() -> Self {
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

    #[cfg(test)]
    mod tests {
        use alloy_primitives::U256;
        use reth_trie_common::StorageProof;

        use super::*;

        fn create_mock_certificate() -> StandardCommitment {
            let commitment_hex = "02f90389e5a0c769488dd5264b3ef21dce7ee2d42fba43e1f83ff228f501223e38818cb14492833f44fcf901eff901caf9018180820001f90159f842a0012e810ffc0a83074b3d14db9e78bbae623f7770cac248df9e73fac6b9d59d17a02a916ffbbf9dde4b7ebe94191a29ff686422d7dcb3b47ecb03c6ada75a9c15c8f888f842a01811c8b4152fce9b8c4bae61a3d097e61dfc43dc7d45363d19e7c7f1374034ffa001edc62174217cdce60a4b52fa234ac0d96db4307dac9150e152ba82cbb4d2f1f842a00f423b0dbc1fe95d2e3f7dbac6c099e51dbf73400a4b3f26b9a29665b4ac58a8a01855a2bd56c0e8f4cc85ac149cf9a531673d0e89e22f0d6c4ae419ed7c5d2940f888f842a02667cbb99d60fa0d7f3544141d3d531dceeeb50b06e5a0cdc42338a359138ae4a00dff4c929d8f8a307c19bba6e8006fe6700f6554cef9eb3797944f89472ffb30f842a004c17a6225acd5b4e7d672a1eb298c5358f4f6f17d04fd1ee295d0c0d372fa84a024bc3ad4d5e54f54f71db382ce276f37ac3c260cc74306b832e8a3c93c7951d302a0e43e11e2405c2fd1d880af8612d969b654827e0ba23d9feb3722ccce6226fce7b8411ddf4553c79c0515516fd3c8b3ae6a756b05723f4d0ebe98a450c8bcc96cbb355ef07a44eeb56f831be73647e4da20e22fa859f984ee41d6efcd3692063b0b0601c2800101a0a69e552a6fc2ff75d32edaf5313642ddeebe60d2069435d12e266ce800e9e96bf9016bc0c0f888f842a00d45727a99053af8d38d4716ab83ace676096e7506b6b7aa6953e87bc04a023ca016c030c31dd1c94062948ecdce2e67c4e6626c16af0033dcdb7a96362c937d48f842a00a95fac74aba7e3fbd24bc62457ce6981803d8f5fef28871d3d5e2af05d50cd4a0117400693917cd50d9bc28d4ab4fadf93a23e771f303637f8d1f83cd0632c3fcf888f842a0301bfced3253e99e8d50f2fed62313a16d714013d022a4dc4294656276f10d1ba0152e047a83c326a9d81dac502ec429b662b58ee119ca4c8748a355b539c24131f842a01944b5b4a3e93d46b0fe4370128c6cdcd066ae6b036b019a20f8d22fe9a10d67a00ddf3421722967c0bd965b9fc9e004bf01183b6206fec8de65e40331d185372ef842a02db8fb278708abf8878ebf578872ab35ee914ad8196b78de16b34498222ac1c2a02ff9d9a5184684f4e14530bde3a61a2f9adaa74734dff104b61ba3d963a644dac68207388208b7c68209998209c5c2c0c0820001";
            let raw_commitment = hex::decode(commitment_hex).expect("Invalid hex in test data");
            StandardCommitment::from_rlp_bytes(raw_commitment.as_slice())
                .expect("Failed to parse test certificate")
        }

        fn create_storage_proof(key: StorageKey, value: U256) -> StorageProof {
            StorageProof {
                key,
                value,
                ..Default::default()
            }
        }

        #[cfg(test)]
        mod stale_stakes_extractors {
            use super::*;

            #[test]
            fn stale_stakes_forbidden_extractor() {
                let extractor = StaleStakesForbiddenExtractor::new();

                let keys = extractor.storage_keys();
                assert_eq!(
                    keys[0],
                    storage_key_helpers::simple_slot_key(STALE_STAKES_FORBIDDEN_VARIABLE_SLOT)
                );

                let storage_key = keys[0];
                let proof_true = create_storage_proof(storage_key, U256::from(1u8));
                let proofs_true = vec![proof_true];
                let result = extractor.decode_data(&proofs_true).unwrap();
                assert!(result);

                let proof_false = create_storage_proof(storage_key, U256::ZERO);
                let proofs_false = vec![proof_false];
                let result = extractor.decode_data(&proofs_false).unwrap();
                assert!(!result);
            }

            #[test]
            fn min_withdrawal_delay_blocks_extractor() {
                let extractor = MinWithdrawalDelayBlocksExtractor::new();

                let keys = extractor.storage_keys();
                assert_eq!(
                    keys[0],
                    storage_key_helpers::simple_slot_key(MIN_WITHDRAWAL_DELAY_BLOCKS_VARIABLE_SLOT)
                );

                let storage_key = keys[0];
                let proof = create_storage_proof(storage_key, U256::from(7200u32));
                let proofs = vec![proof];
                let result = extractor.decode_data(&proofs).unwrap();
                assert_eq!(result, 7200u32);
            }

            #[test]
            fn quorum_update_block_number_extractor() {
                let cert = create_mock_certificate();
                let extractor = QuorumUpdateBlockNumberExtractor::new(&cert);

                let keys = extractor.storage_keys();
                assert_eq!(keys.len(), cert.signed_quorum_numbers().len());

                let proofs = vec![];
                let err = extractor.decode_data(&proofs).unwrap_err();
                assert!(matches!(err, CertExtractionError::MissingStorageProof(_)));
            }
        }
    }
}
