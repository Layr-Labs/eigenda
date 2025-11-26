use core::fmt::Debug;

use alloy_primitives::{Address, address};
use serde::{Deserialize, Serialize};

/// EigenDA directory address on Ethereum mainnet.
pub const EIGENDA_DIRECTORY_MAINNET: Address =
    address!("0x64AB2e9A86FA2E183CB6f01B2D4050c1c2dFAad4");
/// EigenDA directory address on the Hoodi test network.
pub const EIGENDA_DIRECTORY_HOODI: Address = address!("0x5a44e56e88abcf610c68340c6814ae7f5c4369fd");
/// EigenDA directory address on the Sepolia test network.
pub const EIGENDA_DIRECTORY_SEPOLIA: Address =
    address!("0x9620dC4B3564198554e4D2b06dEFB7A369D90257");
/// EigenDA directory address on the Inabox local devnet.
/// This address could get outdated if contract deployment script changes...
/// run `make start-inabox` and get the EIGENDA_DIRECTORY_ADDR printed to stdout.
pub const EIGENDA_DIRECTORY_INABOX: Address =
    address!("0x1613beB3B2C4f22Ee086B2b38C1476A3cE7f78E8");

/// EigenDA CertVerifier v3.1.0 address on Ethereum mainnet.
pub const STATIC_CERT_VERIFIER_MAINNET: Address =
    address!("0x46766C6426eF4D3092f73F72660A8b7B510E6846");
/// EigenDA CertVerifier v3.1.0 address on the Hoodi test network.
pub const STATIC_CERT_VERIFIER_HOODI: Address =
    address!("0xe0F78542A950A8695f43B19Ad1Db654249e12643");
/// EigenDA CertVerifier v3.1.0 address on the Sepolia test network.
pub const STATIC_CERT_VERIFIER_SEPOLIA: Address =
    address!("0x19a469Ddb7199c7EB9E40455978b39894BB90974");
/// EigenDA CertVerifier v3.1.0 address on the Inabox local devnet.
/// To fetch this address, run `make start-inabox` and run
/// ```bash
/// export EIGENDA_CERT_VERIFIER_ROUTER_ADDR=$(cast call $EIGENDA_DIRECTORY_ADDR "getAddress(string)(address)" "CERT_VERIFIER_ROUTER")
/// cast call $EIGENDA_CERT_VERIFIER_ROUTER_ADDR "getCertVerifierAt(uint32)(address)" 0
/// ```
pub const STATIC_CERT_VERIFIER_INABOX: Address =
    address!("0x172076E0166D1F9Cc711C77Adf8488051744980C");

/// EigenDA relevant contracts. Addresses are retrieved from the the EigenDADirectory contract for
/// the respective network (i.e. Mainnet, Hoodi)
#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct EigenDaContracts {
    /// # Ethereum description
    ///
    /// The `EigenDAThresholdRegistry` contract.
    ///
    /// # Details
    ///
    /// The `versionedBlobParams` mapping is read from it
    pub threshold_registry: Address,

    /// # Ethereum description
    ///
    /// A `RegistryCoordinator` that has three registries:
    ///   1. a `StakeRegistry` that keeps track of operators' stakes
    ///   2. a `BLSApkRegistry` that keeps track of operators' BLS public keys and aggregate BLS public keys for each quorum
    ///   3. an `IndexRegistry` that keeps track of an ordered list of operators for each quorum
    ///
    /// # Details
    ///
    /// The quorumCount variable is read from it
    /// The _operatorBitmapHistory mapping is read from it
    /// The quorumUpdateBlockNumber mapping is read from it
    pub registry_coordinator: Address,

    /// # Ethereum description
    ///
    /// Primary entrypoint for procuring services from EigenDA.
    /// This contract is used for:
    /// - initializing the data store by the disperser
    /// - confirming the data store by the disperser with inferred aggregated signatures of the quorum
    /// - freezing operators as the result of various "challenges"
    ///
    /// # Details
    ///
    /// The staleStakesForbidden variable is read from it
    pub service_manager: Address,

    /// # Ethereum description
    ///
    /// The `BlsApkRegistry` contract.
    ///
    /// # Details
    ///
    /// The apkHistory mapping is read from it
    pub bls_apk_registry: Address,

    /// # Ethereum description
    ///
    /// A `Registry` that keeps track of stakes of operators for up to 256 quorums.
    /// Specifically, it keeps track of
    ///   1. The stake of each operator in all the quorums they are a part of for block ranges
    ///   2. The total stake of all operators in each quorum for block ranges
    ///   3. The minimum stake required to register for each quorum
    ///
    /// It allows an additional functionality (in addition to registering and deregistering) to update the stake of an operator.
    ///
    /// # Details
    ///
    /// The _totalStakeHistory mapping is read from it
    /// The operatorStakeHistory mapping is read from it
    pub stake_registry: Address,

    /// # Ethereum description
    ///
    /// A CertVerifier is an immutable contract that is used by a consumer to verify EigenDA blob certificates
    /// For V2 verification this contract is deployed with immutable security thresholds and required quorum numbers,
    /// to change these values or verification behavior a new CertVerifier must be deployed
    ///
    /// # Details
    ///
    /// The quorumNumbersRequiredV2 variable is read from it
    /// The securityThresholdsV2 variable is read from it
    pub cert_verifier: Address,

    /// # Ethereum description
    ///
    /// This is the contract for delegation in EigenLayer. The main functionalities of this contract are
    /// - enabling anyone to register as an operator in EigenLayer
    /// - allowing operators to specify parameters related to stakers who delegate to them
    /// - enabling any staker to delegate its stake to the operator of its choice (a given staker can only delegate to a single operator at a time)
    /// - enabling a staker to undelegate its assets from the operator it is delegated to (performed as part of the withdrawal process, initiated through the StrategyManager)
    ///
    /// # Details
    ///
    /// The minWithdrawalDelayBlocks variable is read from it
    pub delegation_manager: Address,
}
