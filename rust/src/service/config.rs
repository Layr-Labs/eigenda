use core::fmt::{Debug, Formatter};

use alloy_primitives::{Address, address};
use alloy_provider::{DynProvider, Provider};
use alloy_rpc_types_eth::TransactionRequest;
use alloy_sol_types::{SolCall, sol};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use crate::service::EigenDaServiceError;
use crate::service::config::IEigenDADirectory::getAddressCall;

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

sol! {
    interface IEigenDADirectory {
        function getAddress(string memory name) external view returns (address);
    }
}

impl Debug for EigenDaConfig {
    fn fmt(&self, f: &mut Formatter<'_>) -> core::fmt::Result {
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
    Sepolia,
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
    /// The EigenDA Directory contract address on mainnet is
    /// `0x64AB2e9A86FA2E183CB6f01B2D4050c1c2dFAad4`. This address serves as the
    /// central registry for all EigenDA contract addresses. The method
    /// dynamically queries the EigenDADirectory contract to retrieve [`EigenDaContracts`].
    ///
    /// <https://docs.eigencloud.xyz/products/eigenda/networks/mainnet#contract-addresses>
    pub async fn mainnet(provider: &DynProvider) -> Result<Self, EigenDaServiceError> {
        let directory_address = address!("0x64AB2e9A86FA2E183CB6f01B2D4050c1c2dFAad4");
        let eigen_da_contracts = Self::from_directory(provider, directory_address).await?;
        Ok(eigen_da_contracts)
    }

    /// Initialize contracts used by the Holesky.
    ///
    /// The EigenDA Directory contract address on Holesky is
    /// `0x90776Ea0E99E4c38aA1Efe575a61B3E40160A2FE`. This address serves as the
    /// central registry for all EigenDA contract addresses. The method
    /// dynamically queries the EigenDADirectory contract to retrieve [`EigenDaContracts`].
    ///
    /// <https://docs.eigencloud.xyz/products/eigenda/networks/holesky#contract-addresses>
    pub async fn holesky(provider: &DynProvider) -> Result<Self, EigenDaServiceError> {
        let directory_address = address!("0x90776Ea0E99E4c38aA1Efe575a61B3E40160A2FE");
        let eigen_da_contracts = Self::from_directory(provider, directory_address).await?;
        Ok(eigen_da_contracts)
    }

    /// Initialize contracts used by the Sepolia.
    ///
    /// The EigenDA Directory contract address on Sepolia is
    /// `0x9620dC4B3564198554e4D2b06dEFB7A369D90257`. This address serves as the
    /// central registry for all EigenDA contract addresses. The method
    /// dynamically queries the EigenDADirectory contract to retrieve [`EigenDaContracts`].
    ///
    /// <https://docs.eigencloud.xyz/products/eigenda/networks/sepolia#contract-addresses>
    pub async fn sepolia(provider: &DynProvider) -> Result<Self, EigenDaServiceError> {
        let directory_address = address!("0x9620dC4B3564198554e4D2b06dEFB7A369D90257");
        let eigen_da_contracts = Self::from_directory(provider, directory_address).await?;
        Ok(eigen_da_contracts)
    }

    /// Query the EigenDADirectory contract to fetch all required contract addresses
    async fn from_directory(
        provider: &DynProvider,
        directory_address: Address,
    ) -> Result<Self, EigenDaServiceError> {
        let eigen_da_contracts = Self {
            threshold_registry: get_address("THRESHOLD_REGISTRY", provider, directory_address)
                .await?,
            registry_coordinator: get_address("REGISTRY_COORDINATOR", provider, directory_address)
                .await?,
            #[cfg(feature = "stale-stakes-forbidden")]
            service_manager: get_address("SERVICE_MANAGER", provider, directory_address).await?,
            bls_apk_registry: get_address("BLS_APK_REGISTRY", provider, directory_address).await?,
            stake_registry: get_address("STAKE_REGISTRY", provider, directory_address).await?,
            cert_verifier: get_address("CERT_VERIFIER", provider, directory_address).await?,
            #[cfg(feature = "stale-stakes-forbidden")]
            delegation_manager: get_address("DELEGATION_MANAGER", provider, directory_address)
                .await?,
        };

        Ok(eigen_da_contracts)
    }
}

/// The function performs a contract call to the EigenDA contract directory
/// to look up an address associated with a given contract name. It uses the
/// `getAddress` function from the directory contract.
async fn get_address(
    name: &'static str,
    provider: &DynProvider,
    directory_address: Address,
) -> Result<Address, EigenDaServiceError> {
    let input = getAddressCall {
        name: name.to_string(),
    };

    let tx = TransactionRequest::default()
        .to(directory_address)
        .input(input.abi_encode().into());

    let src = provider.call(tx).await?;

    Ok(Address::from_slice(&src[12..32]))
}
