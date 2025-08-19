use std::{hash::Hash, str::FromStr};

use alloy_consensus::{
    EthereumTxEnvelope, Header, Transaction, TxEip4844,
    serde_bincode_compat::{self},
};
use alloy_eips::Typed2718;
use alloy_primitives::{Address, AddressError, B256, FixedBytes, wrap_fixed_bytes};
use borsh::{BorshDeserialize, BorshSerialize};
use bytes::Bytes;
use reth_trie_common::{AccountProof, proof::ProofVerificationError};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
use serde_with::serde_as;
use sov_rollup_interface::{
    BasicAddress,
    da::{BlobReaderTrait, BlockHashTrait, BlockHeaderTrait, CountedBufReader, DaSpec, Time},
    sov_universal_wallet::UniversalWallet,
};

use crate::{
    eigenda::{
        extraction::{
            ApkHistoryExtractor, DataDecoder, MinWithdrawalDelayBlocksExtractor,
            OperatorBitmapHistoryExtractor, OperatorStakeHistoryExtractor, QuorumCountExtractor,
            QuorumNumbersRequiredV2Extractor, QuorumUpdateBlockNumberExtractor,
            RelayKeyToRelayInfoExtractor, SecurityThresholdsV2Extractor,
            StaleStakesForbiddenExtractor, TotalStakeHistoryExtractor,
            VersionedBlobParamsExtractor,
        },
        types::StandardCommitment,
        verification::cert::{
            CertVerificationInputs, error::CertVerificationError, types::Storage,
        },
    },
    verifier::{EigenDaCompletenessProof, EigenDaInclusionProof},
};

/// A specification for the types used by a DA layer.
#[derive(Clone, Debug, Default, PartialEq, Eq, BorshDeserialize, BorshSerialize)]
pub struct EigenDaSpec;

impl DaSpec for EigenDaSpec {
    /// The hash of a DA layer block
    type SlotHash = EthereumHash;

    /// The block header type used by the DA layer
    type BlockHeader = EthereumBlockHeader;

    /// The transaction type used by the DA layer.
    type BlobTransaction = BlobWithSender;

    /// How transactions can be identified on the DA layer.
    type TransactionId = EthereumHash;

    /// The type used to represent addresses on the DA layer.
    type Address = EthereumAddress;

    /// A proof that each tx in a set of blob transactions is included in a given block.
    type InclusionMultiProof = EigenDaInclusionProof;

    /// A proof that a claimed set of transactions is complete.
    type CompletenessProof = EigenDaCompletenessProof;

    /// The parameters of the rollup which are baked into the state-transition function.
    type ChainParams = RollupParams;
}

#[derive(Debug, Copy, Clone, PartialEq, Eq, Serialize, Deserialize, Hash)]
pub struct RollupParams {
    /// The account to which we are storing the certificates of the batch blobs
    pub rollup_batch_namespace: NamespaceId,
    /// The account to which we are storing the certificates of the proof blobs
    pub rollup_proof_namespace: NamespaceId,
    /// A cert is considered valid when it is included onchain before the cert's ReferenceBlockNumber (RBN) + the cert's CPW (Cert punctuality window).
    ///
    /// https://docs.eigencloud.xyz/products/eigenda/integrations-guides/rollup-guides/glossary#cert-punctuality-window
    pub cert_recency_window: u64,
}

/// A namespace id used to identify transactions of the sequencer. The namespace
/// is a regular [`EthereumAddress`]. We say that the specific transaction is
/// part of a namespace if the receiver equals the [`EthereumAddress`] used as a namespace.
#[derive(Debug, Copy, Clone, PartialEq, Eq, Serialize, Deserialize, Hash)]
#[serde(transparent)]
pub struct NamespaceId(EthereumAddress);

impl NamespaceId {
    /// Check if namespace contains this transaction. The namespace contains
    /// transaction if the receiver of transaction is the address used as a namespace.
    pub fn contains<T>(&self, tx: &T) -> bool
    where
        T: Typed2718 + Transaction,
    {
        tx.is_eip1559() && tx.to().is_some_and(|to| to == self.0.0)
    }
}

impl From<NamespaceId> for Address {
    fn from(value: NamespaceId) -> Self {
        value.0.0
    }
}

impl FromStr for NamespaceId {
    type Err = AddressError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        Ok(Self(EthereumAddress::from_str(s)?))
    }
}

/// An Ethereum block header containing only relevant information.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize)]
pub struct EthereumBlockHeader(Header);

impl EthereumBlockHeader {
    /// The function checks if **this** [`EthereumBlockHeader`] is a direct
    /// parent to the [`EthereumBlockHeader`] passed as an argument.
    ///
    /// Note: The relationship is checked by hashing this header and checking it
    /// against the parent_hash of `maybe_child`.
    pub fn is_parent(&self, maybe_child: &EthereumBlockHeader) -> bool {
        let our_hash = self.0.hash_slow();
        maybe_child.0.parent_hash == our_hash
    }
}

impl From<Header> for EthereumBlockHeader {
    fn from(header: Header) -> Self {
        Self(header)
    }
}

impl AsRef<Header> for EthereumBlockHeader {
    fn as_ref(&self) -> &Header {
        &self.0
    }
}

impl BlockHeaderTrait for EthereumBlockHeader {
    type Hash = EthereumHash;

    fn prev_hash(&self) -> Self::Hash {
        self.0.parent_hash.into()
    }

    fn hash(&self) -> Self::Hash {
        self.0.hash_slow().into()
    }

    fn height(&self) -> u64 {
        self.0.number
    }

    fn time(&self) -> Time {
        let timestamp = self
            .0
            .timestamp
            .try_into()
            .expect("is able to convert to i64");
        Time::from_secs(timestamp)
    }
}

#[derive(
    Debug,
    derive_more::Display,
    Clone,
    Copy,
    Serialize,
    Deserialize,
    Hash,
    PartialEq,
    Eq,
    PartialOrd,
    Ord,
    UniversalWallet,
)]
pub struct EthereumAddress(#[sov_wallet(as_ty = "[u8; 20]", display = "hex")] Address);

impl BasicAddress for EthereumAddress {}

impl BorshSerialize for EthereumAddress {
    fn serialize<W: std::io::Write>(&self, writer: &mut W) -> std::io::Result<()> {
        writer.write_all(&self.0.0.0)
    }
}

impl BorshDeserialize for EthereumAddress {
    fn deserialize_reader<R: std::io::Read>(reader: &mut R) -> std::io::Result<Self> {
        let bytes = <[u8; 20]>::deserialize_reader(reader)?;
        Ok(Self(bytes.into()))
    }
}

impl JsonSchema for EthereumAddress {
    fn schema_name() -> String {
        "EthereumAddress".to_string()
    }

    fn json_schema(_generator: &mut schemars::r#gen::SchemaGenerator) -> schemars::schema::Schema {
        serde_json::from_value(serde_json::json!({
            "type": "string",
            "pattern": "^0x[a-fA-F0-9]{40}$",
            "description": "An Ethereum address",
        }))
        .expect("valid schema")
    }
}

impl AsRef<[u8]> for EthereumAddress {
    fn as_ref(&self) -> &[u8] {
        self.0.as_ref()
    }
}

impl FromStr for EthereumAddress {
    type Err = AddressError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        Ok(EthereumAddress(Address::parse_checksummed(s, None)?))
    }
}

impl TryFrom<&[u8]> for EthereumAddress {
    type Error = anyhow::Error;

    fn try_from(value: &[u8]) -> Result<Self, Self::Error> {
        Ok(EthereumAddress(Address::try_from(value)?))
    }
}

impl From<Address> for EthereumAddress {
    fn from(value: Address) -> Self {
        Self(value)
    }
}

impl From<EthereumAddress> for Address {
    fn from(value: EthereumAddress) -> Self {
        value.0
    }
}

wrap_fixed_bytes!(pub struct EthereumHash<32>;);

impl BlockHashTrait for EthereumHash {}

impl BorshSerialize for EthereumHash {
    fn serialize<W: std::io::Write>(&self, writer: &mut W) -> Result<(), std::io::Error> {
        BorshSerialize::serialize(&self.0.0, writer)
    }
}

impl BorshDeserialize for EthereumHash {
    fn deserialize_reader<R: std::io::Read>(reader: &mut R) -> std::io::Result<Self> {
        let bytes = <[u8; 32]>::deserialize_reader(reader)?;
        Ok(Self(FixedBytes::from(bytes)))
    }
}

/// A blob containing the sender address.
#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct BlobWithSender {
    /// The address that submitted blob to the chain.
    pub sender: EthereumAddress,
    /// The ethereum transaction hash in which the blob was included.
    pub tx_hash: EthereumHash,
    /// The actual blob of bytes
    pub blob: CountedBufReader<Bytes>,
}

impl BlobWithSender {
    pub fn new<Address, Hash>(sender: Address, tx_hash: Hash, blob: Bytes) -> Self
    where
        Address: Into<EthereumAddress>,
        Hash: Into<EthereumHash>,
    {
        Self {
            sender: sender.into(),
            tx_hash: tx_hash.into(),
            blob: CountedBufReader::new(blob),
        }
    }
}

impl BlobReaderTrait for BlobWithSender {
    type Address = EthereumAddress;

    type BlobHash = EthereumHash;

    fn sender(&self) -> Self::Address {
        self.sender
    }

    fn hash(&self) -> Self::BlobHash {
        self.tx_hash
    }

    fn verified_data(&self) -> &[u8] {
        self.blob.accumulator()
    }

    fn total_len(&self) -> usize {
        self.blob.total_len()
    }

    #[cfg(feature = "native")]
    fn advance(&mut self, num_bytes: usize) -> &[u8] {
        self.blob.advance(num_bytes);
        self.verified_data()
    }
}

/// Struct that holds an Ethereum transaction with an actual blob
/// persisted to EigenDA
#[serde_as]
#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct TransactionWithBlob {
    /// The transaction that holds a certificate.
    #[serde_as(as = "serde_bincode_compat::EthereumTxEnvelope<'_>")]
    pub tx: EthereumTxEnvelope<TxEip4844>,
    /// The blob persisted to the EigenDA.
    pub blob: Option<Bytes>,
}

/// Data tracked for the specific ancestor.
#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct AncestorMetadata {
    // Header for the ancestor block.
    pub header: EthereumBlockHeader,
    // The data needed to validate the certificate referencing this ancestor.
    // It's `Some` only in cases when we have a certificate that references this
    // ancestor.
    pub data: Option<AncestorStateData>,
}

/// Contains data needed to validate the certificate using the ancestor as the
/// reference block. It also contains proofs used to verify the data.
#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct AncestorStateData {
    eigen_da_relay_registry: AccountProof,
    eigen_da_threshold_registry: AccountProof,
    registry_coordinator: AccountProof,
    bls_signature_checker: AccountProof,
    delegation_manager: AccountProof,
    bls_apk_registry: AccountProof,
    stake_registry: AccountProof,
    eigen_da_cert_verifier: AccountProof,
}

impl AncestorStateData {
    #[allow(clippy::too_many_arguments)]
    pub fn new(
        eigen_da_relay_registry: AccountProof,
        eigen_da_threshold_registry: AccountProof,
        registry_coordinator: AccountProof,
        bls_signature_checker: AccountProof,
        delegation_manager: AccountProof,
        bls_apk_registry: AccountProof,
        stake_registry: AccountProof,
        eigen_da_cert_verifier: AccountProof,
    ) -> Self {
        Self {
            eigen_da_relay_registry,
            eigen_da_threshold_registry,
            registry_coordinator,
            bls_signature_checker,
            delegation_manager,
            bls_apk_registry,
            stake_registry,
            eigen_da_cert_verifier,
        }
    }

    pub fn verify(&self, state_root: B256) -> Result<(), ProofVerificationError> {
        self.eigen_da_relay_registry.verify(state_root)?;
        self.eigen_da_threshold_registry.verify(state_root)?;
        self.registry_coordinator.verify(state_root)?;
        self.bls_signature_checker.verify(state_root)?;
        self.delegation_manager.verify(state_root)?;
        self.bls_apk_registry.verify(state_root)?;
        self.stake_registry.verify(state_root)?;
        self.eigen_da_cert_verifier.verify(state_root)?;

        Ok(())
    }

    /// Extract the data that this ancestor contains.
    ///
    /// NOTE: The data extracted is not verified. To verify the data, ensure
    /// that the [`AncestorStateData::verify`] is called.
    pub fn extract(
        &self,
        cert: &StandardCommitment,
        current_block: u32,
    ) -> Result<CertVerificationInputs, CertVerificationError> {
        // TODO: can we make the association (to contract) type safe?

        let quorum_count = QuorumCountExtractor::new(cert)
            .decode_data(&self.registry_coordinator.storage_proofs)?;

        let stale_stakes_forbidden = StaleStakesForbiddenExtractor::new(cert)
            .decode_data(&self.bls_signature_checker.storage_proofs)?;

        let min_withdrawal_delay_blocks = MinWithdrawalDelayBlocksExtractor::new(cert)
            .decode_data(&self.delegation_manager.storage_proofs)?;

        let quorum_bitmap_history = OperatorBitmapHistoryExtractor::new(cert)
            .decode_data(&self.registry_coordinator.storage_proofs)?;

        let operator_stake_history = OperatorStakeHistoryExtractor::new(cert)
            .decode_data(&self.stake_registry.storage_proofs)?;

        let total_stake_history = TotalStakeHistoryExtractor::new(cert)
            .decode_data(&self.stake_registry.storage_proofs)?;

        let apk_history =
            ApkHistoryExtractor::new(cert).decode_data(&self.bls_apk_registry.storage_proofs)?;

        let quorum_update_block_number = QuorumUpdateBlockNumberExtractor::new(cert)
            .decode_data(&self.registry_coordinator.storage_proofs)?;

        let relay_key_to_relay_address = RelayKeyToRelayInfoExtractor::new(cert)
            .decode_data(&self.eigen_da_relay_registry.storage_proofs)?;

        let versioned_blob_params = VersionedBlobParamsExtractor::new(cert)
            .decode_data(&self.eigen_da_threshold_registry.storage_proofs)?;

        let storage = Storage {
            quorum_count,
            current_block,
            stale_stakes_forbidden,
            min_withdrawal_delay_blocks,
            quorum_bitmap_history,
            operator_stake_history,
            total_stake_history,
            apk_history,
            quorum_update_block_number,
            relay_key_to_relay_address,
            versioned_blob_params,
        };

        let security_thresholds = SecurityThresholdsV2Extractor::new(cert)
            .decode_data(&self.eigen_da_cert_verifier.storage_proofs)?;

        let required_quorum_numbers = QuorumNumbersRequiredV2Extractor::new(cert)
            .decode_data(&self.eigen_da_cert_verifier.storage_proofs)?;

        let inputs = CertVerificationInputs {
            batch_header: cert.batch_header_v2().clone(),
            blob_inclusion_info: cert.blob_inclusion_info().clone(),
            non_signer_stakes_and_signature: cert.nonsigner_stake_and_signature().clone(),
            security_thresholds,
            required_quorum_numbers,
            signed_quorum_numbers: cert.signed_quorum_numbers().clone(),
            storage,
        };

        Ok(inputs)
    }
}

#[cfg(test)]
mod tests {
    use std::str::FromStr;

    use sov_rollup_interface::sov_universal_wallet::schema::Schema;

    use crate::spec::{EthereumAddress, EthereumHash};

    const ADDRESS: &str = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266";
    const HASH: &str = "0x1234567890123456789012345678901234567890123456789012345678901234";

    // TODO: Add more tests and cleanup. Maybe even use proptest.

    #[test]
    fn test_ethereum_address_schema() {
        let address = EthereumAddress::from_str(ADDRESS).unwrap();
        let schema = Schema::of_single_type::<EthereumAddress>().unwrap();

        let borsh_bytes = borsh::to_vec(&address).unwrap();
        let deserialized: EthereumAddress = borsh::from_slice(&borsh_bytes).unwrap();
        assert_eq!(deserialized, address);

        let displayed_from_schema = schema.display(0, &borsh_bytes).unwrap();

        let lowercase_address = ADDRESS.to_lowercase();
        assert_eq!(&displayed_from_schema, &lowercase_address);
    }

    #[test]
    fn test_address_display_from_string() {
        let address = EthereumAddress::from_str(ADDRESS).unwrap();
        let output = format!("{}", address);
        assert_eq!(ADDRESS, output);
    }

    #[test]
    fn test_ethereum_address_try_from() {
        let bytes = hex::decode(&ADDRESS[2..]).unwrap();
        let address = EthereumAddress::try_from(bytes.as_slice()).unwrap();
        assert_eq!(address.to_string(), ADDRESS);
    }

    #[test]
    fn test_ethereum_address_as_ref() {
        let address = EthereumAddress::from_str(ADDRESS).unwrap();
        let bytes: &[u8] = address.as_ref();
        assert_eq!(bytes.len(), 20); // Ethereum addresses are 20 bytes
        assert_eq!(hex::encode(bytes), ADDRESS[2..].to_lowercase());
    }

    #[test]
    fn test_ethereum_address_invalid_input() {
        // Test invalid length
        let result = EthereumAddress::from_str("0x1234");
        assert!(result.is_err());

        // Test invalid hex
        let result = EthereumAddress::from_str("0xg39Fd6e51aad88F6F4ce6aB8827279cffFb92266");
        assert!(result.is_err());
    }

    #[test]
    fn test_ethereum_hash_serialization() {
        let ethereum_hash = EthereumHash::from_str(HASH).unwrap();
        let serde_serialized = serde_json::to_string(&ethereum_hash).unwrap();
        let serde_deserialized: EthereumHash = serde_json::from_str(&serde_serialized).unwrap();

        assert_eq!(ethereum_hash, serde_deserialized);

        let borsh_serialized = borsh::to_vec(&ethereum_hash).unwrap();
        let borsh_deserialized: EthereumHash = borsh::from_slice(&borsh_serialized).unwrap();

        assert_eq!(ethereum_hash, borsh_deserialized);
    }
}
