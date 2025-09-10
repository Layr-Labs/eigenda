use std::{hash::Hash, str::FromStr};

use alloy_consensus::{EthereumTxEnvelope, Header, Transaction, TxEip4844};
use alloy_eips::Typed2718;
use alloy_primitives::{Address, AddressError, B256, FixedBytes, TxHash, wrap_fixed_bytes};
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
use tracing::instrument;

#[cfg(feature = "stale-stakes-forbidden")]
use crate::eigenda::extraction::{
    MinWithdrawalDelayBlocksExtractor, QuorumUpdateBlockNumberExtractor,
    StaleStakesForbiddenExtractor,
};

use crate::{
    eigenda::{
        cert::StandardCommitment,
        extraction::{
            ApkHistoryExtractor, DataDecoder, NextBlobVersionExtractor,
            OperatorBitmapHistoryExtractor, OperatorStakeHistoryExtractor, QuorumCountExtractor,
            QuorumNumbersRequiredV2Extractor, SecurityThresholdsV2Extractor,
            TotalStakeHistoryExtractor, VersionedBlobParamsExtractor,
        },
        verification::cert::{
            CertVerificationInputs, error::CertVerificationError, types::Storage,
        },
    },
    ethereum::tx::serde_bincode_compat::{self},
    verifier::{EigenDaCompletenessProof, EigenDaInclusionProof},
};

/// A specification for the types used by a DA layer.
#[derive(
    Clone, Debug, Default, PartialEq, Eq, Serialize, Deserialize, BorshDeserialize, BorshSerialize,
)]
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
    /// <https://docs.eigencloud.xyz/products/eigenda/integrations-guides/rollup-guides/glossary#cert-punctuality-window>
    pub cert_recency_window: u64,
}

/// A namespace id used to identify transactions of the sequencer. We say that
/// the specific transaction is part of a namespace if the receiver equals the
/// [`Address`] used as a namespace.
#[derive(Debug, Copy, Clone, PartialEq, Eq, Serialize, Deserialize, Hash)]
#[serde(transparent)]
pub struct NamespaceId(Address);

impl NamespaceId {
    pub const fn from_bytes(bytes: [u8; 20]) -> Self {
        Self(Address(FixedBytes(bytes)))
    }

    /// Check if namespace contains this transaction. The namespace contains
    /// transaction if the receiver of transaction is the address used as a namespace.
    pub fn contains<T>(&self, tx: &T) -> bool
    where
        T: Typed2718 + Transaction,
    {
        tx.is_eip1559() && tx.to().is_some_and(|to| to == self.0)
    }
}

impl From<[u8; 20]> for NamespaceId {
    fn from(bytes: [u8; 20]) -> Self {
        Self::from_bytes(bytes)
    }
}

impl From<Address> for NamespaceId {
    fn from(address: Address) -> Self {
        Self(address)
    }
}

impl From<NamespaceId> for Address {
    fn from(value: NamespaceId) -> Self {
        value.0
    }
}

impl FromStr for NamespaceId {
    type Err = AddressError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        Ok(Self(Address::from_str(s)?))
    }
}

/// An Ethereum block header containing only relevant information.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize)]
pub struct EthereumBlockHeader(Header);

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

#[cfg(feature = "arbitrary")]
mod arbitrary_impl {
    use prop::arbitrary::any;
    use prop::strategy::Strategy;
    use proptest::prelude::prop;
    use proptest::strategy::BoxedStrategy;

    use super::*;

    fn new(bytes: [u8; 20]) -> EthereumAddress {
        EthereumAddress(Address(FixedBytes::from(bytes)))
    }

    impl<'a> ::arbitrary::Arbitrary<'a> for EthereumAddress {
        fn arbitrary(u: &mut arbitrary::Unstructured<'a>) -> arbitrary::Result<Self> {
            let bytes = <[u8; 20]>::arbitrary(u)?;
            Ok(new(bytes))
        }
    }

    impl proptest::arbitrary::Arbitrary for EthereumAddress {
        type Parameters = ();

        fn arbitrary_with((): Self::Parameters) -> Self::Strategy {
            any::<[u8; 20]>().prop_map(new).boxed()
        }

        type Strategy = BoxedStrategy<Self>;
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
    pub sender: Address,
    /// The Ethereum transaction hash in which the blob was included.
    pub tx_hash: TxHash,
    /// The actual blob of bytes
    pub blob: CountedBufReader<Bytes>,
}

impl BlobWithSender {
    pub fn new(sender: Address, tx_hash: TxHash, blob: Bytes) -> Self {
        Self {
            sender,
            tx_hash,
            blob: CountedBufReader::new(blob),
        }
    }
}

impl BlobReaderTrait for BlobWithSender {
    type Address = EthereumAddress;
    type BlobHash = EthereumHash;

    fn sender(&self) -> Self::Address {
        EthereumAddress(self.sender)
    }

    fn hash(&self) -> Self::BlobHash {
        EthereumHash(self.tx_hash)
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

/// Struct that holds an Ethereum transaction. If the transaction contains a
/// certificate. Then the [`TransactionWithBlob`] also contains a
/// [`CertificateStateData`] used to verify the included certificate. If the
/// certificate is valid, then it should also contain a data blob that was
/// persisted to the EigenDA.
#[serde_as]
#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct TransactionWithBlob {
    /// The transaction that holds a certificate.
    #[serde_as(as = "serde_bincode_compat::EthereumTxEnvelope<'_>")]
    pub tx: EthereumTxEnvelope<TxEip4844>,
    /// Data used to verify the certificate.
    pub cert_state: Option<CertificateStateData>,
    /// The blob persisted to the EigenDA.
    pub blob: Option<Bytes>,
}

/// Contains data needed to validate the certificate. It also contains proofs
/// used to verify the data.
#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct CertificateStateData {
    pub threshold_registry: AccountProof,
    pub registry_coordinator: AccountProof,
    #[cfg(feature = "stale-stakes-forbidden")]
    pub service_manager: AccountProof,
    pub bls_apk_registry: AccountProof,
    pub stake_registry: AccountProof,
    pub cert_verifier: AccountProof,
    #[cfg(feature = "stale-stakes-forbidden")]
    pub delegation_manager: AccountProof,
}

impl CertificateStateData {
    #![allow(clippy::result_large_err)]
    pub fn verify(&self, state_root: B256) -> Result<(), ProofVerificationError> {
        self.threshold_registry.verify(state_root)?;
        self.registry_coordinator.verify(state_root)?;
        #[cfg(feature = "stale-stakes-forbidden")]
        self.service_manager.verify(state_root)?;
        self.bls_apk_registry.verify(state_root)?;
        self.stake_registry.verify(state_root)?;
        self.cert_verifier.verify(state_root)?;
        #[cfg(feature = "stale-stakes-forbidden")]
        self.delegation_manager.verify(state_root)?;

        Ok(())
    }

    /// Extract the data.
    ///
    /// NOTE: The data extracted is not verified. To verify the data, ensure
    /// that the [`CertificateStateData::verify`] is called.
    #[instrument(skip_all)]
    pub fn extract(
        &self,
        cert: &StandardCommitment,
        current_block: u32,
    ) -> Result<CertVerificationInputs, CertVerificationError> {
        let quorum_count = QuorumCountExtractor::new(cert)
            .decode_data(&self.registry_coordinator.storage_proofs)?;

        let quorum_bitmap_history = OperatorBitmapHistoryExtractor::new(cert)
            .decode_data(&self.registry_coordinator.storage_proofs)?;

        let operator_stake_history = OperatorStakeHistoryExtractor::new(cert)
            .decode_data(&self.stake_registry.storage_proofs)?;

        let total_stake_history = TotalStakeHistoryExtractor::new(cert)
            .decode_data(&self.stake_registry.storage_proofs)?;

        let apk_history =
            ApkHistoryExtractor::new(cert).decode_data(&self.bls_apk_registry.storage_proofs)?;

        let versioned_blob_params = VersionedBlobParamsExtractor::new(cert)
            .decode_data(&self.threshold_registry.storage_proofs)?;

        let next_blob_version = NextBlobVersionExtractor::new(cert)
            .decode_data(&self.threshold_registry.storage_proofs)?;

        #[cfg(feature = "stale-stakes-forbidden")]
        let staleness = {
            use crate::eigenda::verification::cert::types::Staleness;

            let stale_stakes_forbidden = StaleStakesForbiddenExtractor::new(cert)
                .decode_data(&self.service_manager.storage_proofs)?;

            let min_withdrawal_delay_blocks = MinWithdrawalDelayBlocksExtractor::new(cert)
                .decode_data(&self.delegation_manager.storage_proofs)?;

            let quorum_update_block_number = QuorumUpdateBlockNumberExtractor::new(cert)
                .decode_data(&self.registry_coordinator.storage_proofs)?;

            Staleness {
                stale_stakes_forbidden,
                min_withdrawal_delay_blocks,
                quorum_update_block_number,
            }
        };

        let storage = Storage {
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
        };

        let security_thresholds = SecurityThresholdsV2Extractor::new(cert)
            .decode_data(&self.cert_verifier.storage_proofs)?;

        let required_quorum_numbers = QuorumNumbersRequiredV2Extractor::new(cert)
            .decode_data(&self.cert_verifier.storage_proofs)?;

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

    use crate::spec::EthereumAddress;

    const ADDR_1: &str = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266";
    const ADDR_2: &str = "0x1234567890123456789012345678901234567890";

    #[test]
    fn test_ethereum_address_schema() {
        let raw_address_str = ADDR_1;
        let address = EthereumAddress::from_str(raw_address_str).unwrap();

        let schema = Schema::of_single_type::<EthereumAddress>().unwrap();

        let borsh_bytes = borsh::to_vec(&address).unwrap();
        let deserialized: EthereumAddress = borsh::from_slice(&borsh_bytes).unwrap();
        assert_eq!(deserialized, address);

        let displayed_from_schema = schema.display(0, &borsh_bytes).unwrap();
        assert_eq!(&displayed_from_schema, &raw_address_str.to_lowercase());
    }

    #[test]
    fn test_address_display_from_string() {
        let raw_address_str = ADDR_1;
        let address = EthereumAddress::from_str(raw_address_str).unwrap();
        let output = format!("{address}");
        assert_eq!(raw_address_str, output);
    }

    #[test]
    fn test_from_string_for_registering() {
        let raw_address_str = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266";
        let address = EthereumAddress::from_str(raw_address_str).unwrap();
        let raw_bytes = address.as_ref().to_vec();
        let expected_bytes = vec![
            0xf3, 0x9f, 0xd6, 0xe5, 0x1a, 0xad, 0x88, 0xf6, 0xf4, 0xce, 0x6a, 0xb8, 0x82, 0x72,
            0x79, 0xcf, 0xff, 0xb9, 0x22, 0x66,
        ];

        assert_eq!(expected_bytes, raw_bytes);
    }

    #[test]
    fn test_address_display_try_vec() {
        let raw_address_str = ADDR_2;
        let raw_address: Vec<u8> = hex::decode(&raw_address_str[2..]).unwrap();
        let address = EthereumAddress::try_from(&raw_address[..]).unwrap();
        let output = format!("{address}");
        assert_eq!(raw_address_str, output);
    }

    #[test]
    fn test_from_str_and_from_slice_same() {
        let raw_address_str = ADDR_2;
        let raw_address_bytes = hex::decode(&raw_address_str[2..]).unwrap();

        let address_from_str = EthereumAddress::from_str(raw_address_str).unwrap();
        let address_from_slice = EthereumAddress::try_from(raw_address_bytes.as_slice()).unwrap();

        assert_eq!(address_from_str, address_from_slice);
    }

    #[cfg(feature = "arbitrary")]
    mod proptest_tests {
        use super::*;
        use test_strategy::proptest;

        #[proptest]
        fn validate_json_schema(input: EthereumAddress) {
            let schema = serde_json::to_value(schemars::schema_for!(EthereumAddress)).unwrap();
            let json = serde_json::to_value(input).unwrap();

            jsonschema::validate(&schema, &json)
                .map_err(|e| e.kind)
                .unwrap()
        }

        #[proptest]
        fn ord_invariants(values: [EthereumAddress; 3]) {
            reltester::ord(&values[0], &values[1], &values[2]).unwrap();
        }

        #[proptest]
        fn hash_invariants(values: [EthereumAddress; 2]) {
            reltester::hash(&values[0], &values[1]).unwrap();
        }

        #[proptest]
        fn test_try_from_any_slice(
            #[any(proptest::sample::size_range(0..100).lift())] input: Vec<u8>,
        ) {
            let _ = EthereumAddress::try_from(&input[..]);
        }

        #[proptest]
        fn test_from_str_anything(#[strategy("\\PC*")] input: String) {
            let _ = EthereumAddress::from_str(&input);
        }

        #[proptest]
        fn test_from_str_hex_addresses(#[strategy("0x[a-fA-F0-9]{40}")] input: String) {
            let result = EthereumAddress::from_str(&input);
            if let Ok(address) = result {
                let output = format!("{address}");
                assert_eq!(input, output);
            }
        }

        #[proptest]
        fn test_try_from_bytes(input: [u8; 20]) {
            let hex_str = format!("0x{}", hex::encode(input));
            let address_from_bytes = EthereumAddress::try_from(input.as_slice()).unwrap();
            let output = format!("{address_from_bytes}");

            assert_eq!(hex_str, output.to_lowercase());
        }
    }
}
