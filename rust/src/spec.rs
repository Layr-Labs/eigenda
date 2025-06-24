use core::fmt;
use std::{fmt::Display, str::FromStr};

use borsh::{BorshDeserialize, BorshSerialize};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
use sov_rollup_interface::{
    BasicAddress,
    da::{BlobReaderTrait, BlockHashTrait, BlockHeaderTrait, DaSpec, Time},
    sov_universal_wallet::UniversalWallet,
};

use crate::verifier::{EigenDaCompletenessProof, EigenDaInclusionProof};

/// A specification for the types used by a DA layer.
#[derive(Clone, Debug, Default, PartialEq, Eq, BorshDeserialize, BorshSerialize)]
pub struct EigenDaSpec;

impl DaSpec for EigenDaSpec {
    /// The hash of a DA layer block
    type SlotHash = EthereumBlockHash;

    /// The block header type used by the DA layer
    type BlockHeader = EthereumBlockHeader;

    /// The transaction type used by the DA layer.
    type BlobTransaction = BlobWithSender;

    /// How transactions can be identified on the DA layer.
    type TransactionId = EthereumBlockHash;

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
pub struct RollupParams;

/// An Ethereum block header containing only relevant information.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize)]
pub struct EthereumBlockHeader;

impl BlockHeaderTrait for EthereumBlockHeader {
    type Hash = EthereumBlockHash;

    fn prev_hash(&self) -> Self::Hash {
        todo!()
    }

    fn hash(&self) -> Self::Hash {
        todo!()
    }

    fn height(&self) -> u64 {
        todo!()
    }

    fn time(&self) -> Time {
        todo!()
    }
}

/// A wrapper for the Ethereum block hash.
#[derive(
    Clone,
    Copy,
    Debug,
    PartialEq,
    Eq,
    Serialize,
    Deserialize,
    Hash,
    BorshDeserialize,
    BorshSerialize,
)]
pub struct EthereumBlockHash;

impl BlockHashTrait for EthereumBlockHash {}

impl FromStr for EthereumBlockHash {
    type Err = ();

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        todo!()
    }
}

impl AsRef<[u8]> for EthereumBlockHash {
    fn as_ref(&self) -> &[u8] {
        todo!()
    }
}

impl From<[u8; 32]> for EthereumBlockHash {
    fn from(value: [u8; 32]) -> Self {
        todo!()
    }
}

impl From<EthereumBlockHash> for [u8; 32] {
    fn from(value: EthereumBlockHash) -> Self {
        todo!()
    }
}

impl fmt::Display for EthereumBlockHash {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        todo!()
    }
}

#[derive(
    Debug,
    Clone,
    JsonSchema,
    BorshSerialize,
    BorshDeserialize,
    Serialize,
    Deserialize,
    Hash,
    PartialEq,
    Eq,
    PartialOrd,
    Ord,
    UniversalWallet,
)]
pub struct EthereumAddress;

impl BasicAddress for EthereumAddress {}

impl AsRef<[u8]> for EthereumAddress {
    fn as_ref(&self) -> &[u8] {
        todo!()
    }
}

impl FromStr for EthereumAddress {
    type Err = anyhow::Error;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        todo!()
    }
}

impl<'a> TryFrom<&'a [u8]> for EthereumAddress {
    type Error = anyhow::Error;

    fn try_from(value: &'a [u8]) -> Result<Self, Self::Error> {
        todo!()
    }
}

impl Display for EthereumAddress {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        todo!()
    }
}

/// A blob containing the sender address.
#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct BlobWithSender;

impl BlobReaderTrait for BlobWithSender {
    type Address = EthereumAddress;

    type BlobHash = EthereumBlockHash;

    fn sender(&self) -> Self::Address {
        todo!()
    }

    fn hash(&self) -> Self::BlobHash {
        todo!()
    }

    fn verified_data(&self) -> &[u8] {
        todo!()
    }

    fn total_len(&self) -> usize {
        todo!()
    }

    #[cfg(feature = "native")]
    fn advance(&mut self, num_bytes: usize) -> &[u8] {
        todo!()
    }
}
