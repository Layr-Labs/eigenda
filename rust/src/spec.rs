use std::{hash::Hash, num::TryFromIntError, str::FromStr};

use alloy::{
    consensus::Transaction as TTransaction,
    eips::Typed2718,
    primitives::{Address, Bytes, FixedBytes, wrap_fixed_bytes},
    rpc::types::{Header, Transaction},
};
use borsh::{BorshDeserialize, BorshSerialize};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
use sov_rollup_interface::{
    BasicAddress,
    da::{BlobReaderTrait, BlockHashTrait, BlockHeaderTrait, CountedBufReader, DaSpec, Time},
    sov_universal_wallet::UniversalWallet,
};

use crate::verifier::{EigenDaCompletenessProof, EigenDaInclusionProof};

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
    /// For example, this could be a range proof demonstrating that
    /// the provided BlobTransactions represent the entire contents
    /// of Celestia namespace in a given block
    type CompletenessProof = EigenDaCompletenessProof;

    /// The parameters of the rollup which are baked into the state-transition function.
    /// For example, this could include the namespace of the rollup on Celestia.
    type ChainParams = RollupParams;
}

#[derive(Debug, Copy, Clone, PartialEq, Eq, Serialize, Deserialize, Hash)]
pub struct RollupParams {
    /// The account to which we are storing the certificates of the batch blobs
    pub rollup_batch_account: NamespaceId,
    /// The account to which we are storing the certificates of the proof blobs
    pub rollup_proof_account: NamespaceId,
}

/// A namespace id used to identify transactions of the sequencer.
#[derive(Debug, Copy, Clone, PartialEq, Eq, Serialize, Deserialize, Hash)]
#[serde(transparent)]
pub struct NamespaceId(EthereumAddress);

impl NamespaceId {
    /// Check if namespace contains this transaction. The namespace contains
    /// transaction if the receiver of transaction is the address used as a namespace.
    pub fn contains(&self, tx: &Transaction) -> bool {
        tx.is_eip4844() && tx.to().is_some_and(|to| to == self.0.0)
    }
}

impl From<NamespaceId> for Address {
    fn from(value: NamespaceId) -> Self {
        value.0.0
    }
}

impl FromStr for NamespaceId {
    type Err = alloy::primitives::AddressError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        Ok(Self(EthereumAddress::from_str(s)?))
    }
}

/// An Ethereum block header containing only relevant information.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize)]
pub struct EthereumBlockHeader {
    parent_hash: EthereumHash,
    hash: EthereumHash,
    timestamp: Time,
    height: u64,
}

impl TryFrom<Header> for EthereumBlockHeader {
    type Error = TryFromIntError;

    fn try_from(header: Header) -> Result<Self, Self::Error> {
        let timestamp = header.timestamp.try_into()?;

        Ok(Self {
            timestamp: Time::from_secs(timestamp),
            parent_hash: EthereumHash::from(header.parent_hash),
            hash: EthereumHash::from(header.hash_slow()),
            height: header.number,
        })
    }
}

impl BlockHeaderTrait for EthereumBlockHeader {
    type Hash = EthereumHash;

    fn prev_hash(&self) -> Self::Hash {
        self.parent_hash
    }

    fn hash(&self) -> Self::Hash {
        self.hash
    }

    fn height(&self) -> u64 {
        self.height
    }

    fn time(&self) -> Time {
        self.timestamp.clone()
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
pub struct EthereumAddress(
    #[sov_wallet(as_ty = "[u8; 20]", display = "hex")] alloy::primitives::Address,
);

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
    type Err = alloy::primitives::AddressError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        Ok(EthereumAddress(
            alloy::primitives::Address::parse_checksummed(s, None)?,
        ))
    }
}

impl TryFrom<&[u8]> for EthereumAddress {
    type Error = anyhow::Error;

    fn try_from(value: &[u8]) -> Result<Self, Self::Error> {
        Ok(EthereumAddress(alloy::primitives::Address::try_from(
            value,
        )?))
    }
}

impl From<alloy::primitives::Address> for EthereumAddress {
    fn from(value: alloy::primitives::Address) -> Self {
        Self(value)
    }
}

impl From<EthereumAddress> for alloy::primitives::Address {
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
    pub fn new<Address, Hash>(sender: Address, tx_hash: Hash, blob: Vec<u8>) -> Self
    where
        Address: Into<EthereumAddress>,
        Hash: Into<EthereumHash>,
    {
        Self {
            sender: sender.into(),
            tx_hash: tx_hash.into(),
            blob: CountedBufReader::new(blob.into()),
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

/// Struct that holds an Ethereum transaction with certificate and an actual blob
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize)]
pub struct TransactionWithBlob {
    /// The transaction that holds a certificate
    pub transaction: Transaction,
    /// The actual blob persisted to the EigenDa
    pub blob: Option<Vec<u8>>,
}

#[cfg(test)]
mod tests {
    use std::str::FromStr;

    use alloy::hex::{self};
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
