#[cfg(feature = "native")]
pub mod contract;
pub mod decode_helpers;
pub mod storage_key_helpers;

use alloy_primitives::{
    Address, B256, Bytes, StorageKey, U256,
    aliases::{U96, U192},
};
use hashbrown::HashMap;
use reth_trie_common::StorageProof;
use thiserror::Error;

use crate::eigenda::{
    types::StandardCommitment,
    verification::cert::{
        bitmap::Bitmap,
        hash::TruncatedB256,
        types::{
            BlockNumber, QuorumNumber, RelayKey, Stake, Version,
            history::{History, HistoryError},
            solidity::{SecurityThresholds, StakeUpdate, VersionedBlobParams},
        },
    },
};

const RELAY_KEY_TO_RELAY_INFO_MAPPING_SLOT: u64 = 101u64;
const VERSIONED_BLOB_PARAMS_MAPPING_SLOT: u64 = 4;
const QUORUM_COUNT_VARIABLE_SLOT: u64 = 150;
const OPERATOR_BITMAP_HISTORY_MAPPING_SLOT: u64 = 152;
const QUORUM_UPDATE_BLOCK_NUMBER_MAPPING_SLOT: u64 = 155;
const STALE_STAKES_FORBIDDEN_VARIABLE_SLOT: u64 = 0;
const MIN_WITHDRAWAL_DELAY_BLOCKS_VARIABLE_SLOT: u64 = 157;
const APK_HISTORY_MAPPING_SLOT: u64 = 4;
const TOTAL_STAKE_HISTORY_MAPPING_SLOT: u64 = 1;
const OPERATOR_STAKE_HISTORY_MAPPING_SLOT: u64 = 2;
const SECURITY_THRESHOLDS_V2_VARIABLE_SLOT: u64 = 0;
const QUORUM_NUMBERS_REQUIRED_V2_VARIABLE_SLOT: u64 = 1;

#[derive(Debug, Error, PartialEq)]
pub enum CertExtractionError {
    #[error("Failed to extract StorageProof for {0}")]
    MissingStorageProof(&'static str),

    #[error(transparent)]
    WrapHistoryError(#[from] HistoryError),
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
    fn storage_keys(&self) -> Vec<StorageKey> {
        vec![storage_key_helpers::simple_slot_key(
            QUORUM_COUNT_VARIABLE_SLOT,
        )]
    }
}

impl DataDecoder for QuorumCountExtractor {
    type Output = u8;

    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError> {
        let storage_key = &self.storage_keys()[0];
        let proof =
            decode_helpers::find_required_proof(storage_proofs, storage_key, "quorumCount")?;
        let quorum_count = proof.value.to::<u8>();
        Ok(quorum_count)
    }
}

pub struct StaleStakesForbiddenExtractor;

impl StaleStakesForbiddenExtractor {
    pub fn new(_certificate: &StandardCommitment) -> Self {
        Self {}
    }
}

impl StorageKeyProvider for StaleStakesForbiddenExtractor {
    fn storage_keys(&self) -> Vec<StorageKey> {
        vec![storage_key_helpers::simple_slot_key(
            STALE_STAKES_FORBIDDEN_VARIABLE_SLOT,
        )]
    }
}

impl DataDecoder for StaleStakesForbiddenExtractor {
    type Output = bool;

    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError> {
        let storage_key = &self.storage_keys()[0];
        let proof =
            decode_helpers::find_required_proof(storage_proofs, storage_key, "quorumCount")?;
        Ok(proof.value.is_zero())
    }
}

pub struct MinWithdrawalDelayBlocksExtractor;

impl MinWithdrawalDelayBlocksExtractor {
    pub fn new(_certificate: &StandardCommitment) -> Self {
        Self {}
    }
}

impl StorageKeyProvider for MinWithdrawalDelayBlocksExtractor {
    fn storage_keys(&self) -> Vec<StorageKey> {
        vec![storage_key_helpers::simple_slot_key(
            MIN_WITHDRAWAL_DELAY_BLOCKS_VARIABLE_SLOT,
        )]
    }
}

impl DataDecoder for MinWithdrawalDelayBlocksExtractor {
    type Output = u32;

    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError> {
        let storage_key = &self.storage_keys()[0];
        let proof =
            decode_helpers::find_required_proof(storage_proofs, storage_key, "quorumCount")?;
        let min_withdrawal_delay_blocks = proof.value.to::<u32>();
        Ok(min_withdrawal_delay_blocks)
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

impl DataDecoder for QuorumUpdateBlockNumberExtractor {
    type Output = HashMap<QuorumNumber, BlockNumber>;

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
                .map(|proof| {
                    let block_number = proof.value.to::<BlockNumber>();
                    (quorum_number, block_number)
                })
            })
            .collect()
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

impl DataDecoder for RelayKeyToRelayInfoExtractor {
    type Output = HashMap<RelayKey, Address>;

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
        let le = proof.value.to_le_bytes::<9>();

        let key = self.version;
        let value = VersionedBlobParams {
            maxNumOperators: u32::from_le_bytes(le[0..4].try_into().unwrap()),
            numChunks: u32::from_le_bytes(le[4..8].try_into().unwrap()),
            codingRate: le[8],
        };
        let versioned_blob_params = HashMap::from([(key, value)]);
        Ok(versioned_blob_params)
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

impl DataDecoder for OperatorBitmapHistoryExtractor {
    type Output = HashMap<B256, History<Bitmap>>;

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
                let bitmap = Bitmap::new([lo, mid, hi, 0]);

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

// TODO: make structs generic over lifetime
impl ApkHistoryExtractor {
    pub fn new(certificate: &StandardCommitment) -> Self {
        Self {
            signed_quorum_numbers: certificate.signed_quorum_numbers().clone(),
            quorum_apk_indices: certificate.quorum_apk_indices().to_vec(),
        }
    }
}

impl StorageKeyProvider for ApkHistoryExtractor {
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

impl DataDecoder for ApkHistoryExtractor {
    type Output = HashMap<QuorumNumber, History<TruncatedB256>>;

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

                let apk_hash_bytes: [u8; 24] = le[..24].try_into().unwrap();
                let apk_hash: TruncatedB256 = apk_hash_bytes;
                let update_block = u32::from_le_bytes(le[24..28].try_into().unwrap());
                let next_update_block = u32::from_le_bytes(le[28..32].try_into().unwrap());

                let update =
                    decode_helpers::create_update(update_block, next_update_block, apk_hash)?;
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

impl DataDecoder for TotalStakeHistoryExtractor {
    type Output = HashMap<QuorumNumber, History<Stake>>;

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

impl DataDecoder for OperatorStakeHistoryExtractor {
    type Output = HashMap<B256, HashMap<QuorumNumber, History<Stake>>>;

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
                    let le = proof.value.to_le_bytes::<20>();
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
    fn storage_keys(&self) -> Vec<StorageKey> {
        vec![storage_key_helpers::simple_slot_key(
            SECURITY_THRESHOLDS_V2_VARIABLE_SLOT,
        )]
    }
}

impl DataDecoder for SecurityThresholdsV2Extractor {
    type Output = SecurityThresholds;

    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError> {
        let storage_key = &self.storage_keys()[0];
        let proof =
            decode_helpers::find_required_proof(storage_proofs, storage_key, "quorumCount")?;

        let [confirmation_threshold, adversary_threshold] = proof.value.to_le_bytes::<2>();

        Ok(SecurityThresholds {
            confirmationThreshold: confirmation_threshold,
            adversaryThreshold: adversary_threshold,
        })
    }
}

pub struct QuorumNumbersRequiredV2Extractor;

impl QuorumNumbersRequiredV2Extractor {
    pub fn new(_certificate: &StandardCommitment) -> Self {
        Self {}
    }
}

impl StorageKeyProvider for QuorumNumbersRequiredV2Extractor {
    fn storage_keys(&self) -> Vec<StorageKey> {
        vec![storage_key_helpers::simple_slot_key(
            QUORUM_NUMBERS_REQUIRED_V2_VARIABLE_SLOT,
        )]
    }
}

impl DataDecoder for QuorumNumbersRequiredV2Extractor {
    type Output = Bytes;

    fn decode_data(
        &self,
        storage_proofs: &[StorageProof],
    ) -> Result<Self::Output, CertExtractionError> {
        let storage_key = &self.storage_keys()[0];
        let proof =
            decode_helpers::find_required_proof(storage_proofs, storage_key, "quorumCount")?;

        // there can be at most 256 quorums
        let bytes = proof.value.to_le_bytes::<32>();

        // quorum numbers are ordered so it's safe (and necessary) to trim
        let bytes = decode_helpers::trim_trailing_zeros(&bytes);

        Ok(bytes.to_vec().into())
    }
}
