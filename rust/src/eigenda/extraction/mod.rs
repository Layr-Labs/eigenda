#[cfg(feature = "native")]
pub mod contract;
pub mod decode_helpers;
pub mod storage_key_helpers;
#[cfg(feature = "stale-stakes-forbidden")]
pub use stale_stakes_forbidden::*;

use alloy_primitives::{
    Address, B256, Bytes, StorageKey, U256,
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
            QuorumNumber, RelayKey, Stake, Version,
            history::{History, HistoryError},
            solidity::{SecurityThresholds, StakeUpdate, VersionedBlobParams},
        },
    },
};

const RELAY_KEY_TO_RELAY_INFO_MAPPING_SLOT: u64 = 101u64;
const VERSIONED_BLOB_PARAMS_MAPPING_SLOT: u64 = 4;
const QUORUM_COUNT_VARIABLE_SLOT: u64 = 150;
const OPERATOR_BITMAP_HISTORY_MAPPING_SLOT: u64 = 152;
const APK_HISTORY_MAPPING_SLOT: u64 = 4;
const TOTAL_STAKE_HISTORY_MAPPING_SLOT: u64 = 1;
const OPERATOR_STAKE_HISTORY_MAPPING_SLOT: u64 = 2;
const SECURITY_THRESHOLDS_V2_VARIABLE_SLOT: u64 = 0;
const QUORUM_NUMBERS_REQUIRED_V2_VARIABLE_SLOT: u64 = 1;

#[derive(Debug, Error, PartialEq)]
pub enum CertExtractionError {
    #[error("Failed to extract StorageProof for {0}")]
    MissingStorageProof(String),

    #[error(transparent)]
    WrapHistoryError(#[from] HistoryError),

    #[error(transparent)]
    WrapAlloySolTypesError(#[from] alloy_sol_types::Error),
}

pub trait StorageKeyProvider {
    fn storage_keys(&self) -> Vec<StorageKey>;
}

pub trait DataDecoder: StorageKeyProvider {
    type Output;

    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError>;
}

pub struct QuorumCountExtractor;

impl QuorumCountExtractor {
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

// REGISTRY_COORDINATOR::quorumCount (3 on holesky)
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

pub struct RelayKeyToRelayInfoExtractor {
    pub relay_keys: Vec<RelayKey>,
}

impl RelayKeyToRelayInfoExtractor {
    pub fn new(certificate: &StandardCommitment) -> Self {
        Self {
            relay_keys: certificate.relay_keys().to_vec(),
        }
    }
}

impl StorageKeyProvider for RelayKeyToRelayInfoExtractor {
    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn storage_keys(&self) -> Vec<StorageKey> {
        self.relay_keys
            .iter()
            .map(|&relay_key| {
                storage_key_helpers::mapping_key(
                    U256::from(relay_key),
                    RELAY_KEY_TO_RELAY_INFO_MAPPING_SLOT,
                )
            })
            .collect()
    }
}

// RELAY_REGISTRY::relayKeyToAddress
// TODO: been using relayKeyToRelayInfo but there's no need because relayKeyToAddress exists
impl DataDecoder for RelayKeyToRelayInfoExtractor {
    type Output = HashMap<RelayKey, Address>;

    #[instrument(skip_all, fields(component = std::any::type_name::<Self>()))]
    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError> {
        self.storage_keys()
            .iter()
            .zip(self.relay_keys.iter())
            .map(|(storage_key, &relay_key)| {
                decode_helpers::find_required_proof(storage_proofs, storage_key, "relayKeyToInfo")
                    .map(|proof| (relay_key, Address::from_word(proof.value.into())))
            })
            .collect()
    }
}

pub struct VersionedBlobParamsExtractor {
    pub version: u16,
}

impl VersionedBlobParamsExtractor {
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

pub struct OperatorBitmapHistoryExtractor {
    pub non_signers_pk_hashes: Vec<B256>,
    pub non_signer_quorum_bitmap_indices: Vec<u32>,
}

impl OperatorBitmapHistoryExtractor {
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

// REGISTRY_COORDINATOR::getQuorumBitmapAtBlockNumberByIndex (accesses _operatorBitmapHistory)
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

pub struct ApkHistoryExtractor {
    pub signed_quorum_numbers: Bytes,
    pub quorum_apk_indices: Vec<u32>,
}

impl ApkHistoryExtractor {
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

// BLS_APK_REGISTRY::apkHistory
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

pub struct TotalStakeHistoryExtractor {
    pub signed_quorum_numbers: Bytes,
    pub non_signer_total_stake_indices: Vec<u32>,
}

impl TotalStakeHistoryExtractor {
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

// STAKE_REGISTRY::getTotalStakeAtBlockNumberFromIndex (accesses _totalStakeHistory)
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

pub struct OperatorStakeHistoryExtractor {
    pub signed_quorum_numbers: Bytes,
    pub non_signers_pk_hashes: Vec<B256>,
    pub non_signer_stake_indices: Vec<Vec<u32>>,
}

impl OperatorStakeHistoryExtractor {
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
                // that map to non-existent data will return empty but won't fail
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

// STAKE_REGISTRY::getStakeAtBlockNumberAndIndex (accesses operatorStakeHistory)
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

pub struct SecurityThresholdsV2Extractor;

impl SecurityThresholdsV2Extractor {
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

// _CERT_VERIFIER_V2::securityThresholdsV2
// (confirmationThreshold: 55u8, adversaryThreshold: 33u8 on holesky)
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

pub struct QuorumNumbersRequiredV2Extractor;

impl QuorumNumbersRequiredV2Extractor {
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

// _CERT_VERIFIER_V2::quorumNumbersRequiredV2 (0x0001 on holesky)
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

    pub struct StaleStakesForbiddenExtractor;

    impl StaleStakesForbiddenExtractor {
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

    // _SERVICE_MANAGER::staleStakesForbidden (false on holesky)
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

    pub struct MinWithdrawalDelayBlocksExtractor;

    impl MinWithdrawalDelayBlocksExtractor {
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

    // DELEGATION_MANAGER::minWithdrawalDelayBlocks
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
            Ok(proof.value.to::<u32>())
        }
    }

    pub struct QuorumUpdateBlockNumberExtractor {
        pub signed_quorum_numbers: Bytes,
    }

    impl QuorumUpdateBlockNumberExtractor {
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

    // REGISTRY_COORDINATOR::quorumUpdateBlockNumber
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
