use core::fmt::{Debug, Formatter, Result};

use alloy_primitives::{Address, address};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

/// Configuration for the [`crate::service::EigenDaService`].
#[derive(Clone, JsonSchema, PartialEq, Serialize, Deserialize)]
pub struct EigenDaConfig {
    /// Network the adapter is running against.
    pub network: Network,
    /// URL of the Ethereum RPC node.
    pub ethereum_rpc_url: String,
    /// The number of compute units per second for the provider. Used in cases
    /// when the Ethereum node is hosted at providers like Alchemy that track
    /// compute units used when making a requests. If None, it means the node is
    /// not tracking compute units.
    pub ethereum_compute_units: Option<u64>,
    /// The maximal number of times we retry requests to the node before
    /// returning the error.
    pub ethereum_max_retry_times: Option<u32>,
    /// The initial backoff in milliseconds used when retrying Ethereum
    /// requests. It is increased on each subsequent retry.
    pub ethereum_initial_backoff: Option<u64>,
    /// URL of the EigenDA proxy node.
    pub proxy_url: String,
    /// The initial backoff in milliseconds used when retrying EigenDA proxy
    /// requests. It is increased on each subsequent retry.
    pub proxy_min_retry_delay: Option<u64>,
    /// The maximal backoff in milliseconds used when retrying EigenDA proxy requests.
    pub proxy_max_retry_delay: Option<u64>,
    /// The maximal number of times we retry requests to the EigenDA proxy
    /// before returning the error.
    pub proxy_max_retry_times: Option<u64>,
    /// Private key of the sequencer. The account with corresponding private key
    /// is used by the sequencer to persist the certificates to Ethereum.
    /// Expected private key in the HEX format.
    pub sequencer_signer: String,
}

impl Debug for EigenDaConfig {
    fn fmt(&self, f: &mut Formatter<'_>) -> Result {
        f.debug_struct("EigenDaConfig")
            .field("network", &self.network)
            .field("ethereum_rpc_url", &self.ethereum_rpc_url)
            .field("ethereum_compute_units", &self.ethereum_compute_units)
            .field("ethereum_max_retry_times", &self.ethereum_max_retry_times)
            .field("ethereum_initial_backoff", &self.ethereum_initial_backoff)
            .field("proxy_url", &self.proxy_url)
            .field("proxy_min_retry_delay", &self.proxy_min_retry_delay)
            .field("proxy_max_retry_delay", &self.proxy_max_retry_delay)
            .field("proxy_max_retry_times", &self.proxy_max_retry_times)
            .field("sequencer_signer", &"[REDACTED]")
            .finish()
    }
}

/// Network the adapter is running against.
#[derive(Debug, Clone, JsonSchema, PartialEq, Serialize, Deserialize)]
pub enum Network {
    Mainnet,
    Holesky,
}

/// EigenDA relevant contracts.
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
    /// The quorumCount variable is read from it. See [eigenda_cert_verifier::types::Storage]'s `quorum_count`
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
    #[cfg(feature = "stale-stakes-forbidden")]
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
    #[cfg(feature = "stale-stakes-forbidden")]
    pub delegation_manager: Address,
}

impl EigenDaContracts {
    /// Initialize contracts used by the Mainnet.
    ///
    /// Instructions on how they were retrieved:
    /// * https://docs.eigencloud.xyz/products/eigenda/networks/mainnet#contract-addresses
    ///
    /// Except for `delegation_manager` taken from:
    /// * https://github.com/Layr-Labs/eigenlayer-contracts/blob/cd5612ec76e31b4f7768f3a2308f658e476d94ea/script/configs/mainnet.json
    pub fn mainnet() -> Self {
        Self {
            threshold_registry: address!("0xdb4c89956eEa6F606135E7d366322F2bDE609F15"),
            registry_coordinator: address!("0x0BAAc79acD45A023E19345c352d8a7a83C4e5656"),
            #[cfg(feature = "stale-stakes-forbidden")]
            service_manager: address!("0x870679E138bCdf293b7Ff14dD44b70FC97e12fc0"),
            bls_apk_registry: address!("0x00A5Fd09F6CeE6AE9C8b0E5e33287F7c82880505"),
            stake_registry: address!("0x006124Ae7976137266feeBFb3F4D2BE4C073139D"),
            cert_verifier: address!("0x61692e93b6B045c444e942A91EcD1527F23A3FB7"),
            #[cfg(feature = "stale-stakes-forbidden")]
            delegation_manager: address!("0x870679E138bCdf293b7Ff14dD44b70FC97e12fc0"),
        }
    }

    /// Initialize contracts used by the Holesky.
    ///
    /// Instructions on how they were retrieved:
    /// * https://docs.eigencloud.xyz/products/eigenda/networks/holesky#contract-addresses
    ///
    /// Except for `delegation_manager` taken from:
    /// * https://github.com/Layr-Labs/eigenlayer-contracts/blob/cd5612ec76e31b4f7768f3a2308f658e476d94ea/script/configs/holesky.json
    pub fn holesky() -> Self {
        Self {
            threshold_registry: address!("0x76d131CFBD900dA12f859a363Fb952eEDD1d1Ec1"),
            registry_coordinator: address!("0x53012C69A189cfA2D9d29eb6F19B32e0A2EA3490"),
            #[cfg(feature = "stale-stakes-forbidden")]
            service_manager: address!("0xD4A7E1Bd8015057293f0D0A557088c286942e84b"),
            bls_apk_registry: address!("0x066cF95c1bf0927124DFB8B02B401bc23A79730D"),
            stake_registry: address!("0xBDACD5998989Eec814ac7A0f0f6596088AA2a270"),
            // CERT_VERIFIER (also available: CERT_VERIFIER_V1, CERT_VERIFIER_V2)
            cert_verifier: address!("0x036bB27A1F03350bDcccF344b497Ef22604006a3"),
            #[cfg(feature = "stale-stakes-forbidden")]
            delegation_manager: address!("0xA44151489861Fe9e3055d95adC98FbD462B948e7"),
        }
    }
}
