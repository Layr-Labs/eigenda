use alloy_primitives::{Address, address};
use alloy_provider::DynProvider;
use alloy_provider::Provider;
use alloy_rpc_types_eth::TransactionRequest;
use alloy_transport::{RpcError, TransportErrorKind};
use core::fmt::Debug;
use serde::{Deserialize, Serialize};

use crate::contracts::IEigenDADirectory::getAddressCall;
use alloy_sol_types::{SolCall, sol};

sol! {
    interface IEigenDADirectory {
        function getAddress(string memory name) external view returns (address);
    }
}

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
    /// A CertVerifierRouter is an upgradable contract that routes cert verification requests to the appropriate CertVerifier contract.
    /// This allows for dynamic updates to the cert verification logic without changing the address that consumers interact with.
    /// For trustless integrations, it is recommended to deploy and use a dedicated CertVerifierRouter contract.
    /// See https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#upgradable-quorums-and-thresholds-for-optimistic-verification for more details.
    ///
    /// # Details
    ///
    /// The cert_verifier contract address at a specific (reference) block number is read from it
    pub cert_verifier_router: Address,

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

impl EigenDaContracts {
    /// Query the EigenDADirectory contract to fetch all required contract addresses
    pub async fn new(
        ethereum: &DynProvider,
        directory_address: Address,
        cert_verifier_router_address: Option<Address>,
    ) -> Result<EigenDaContracts, RpcError<TransportErrorKind>> {
        let eigen_da_contracts = EigenDaContracts {
            threshold_registry: get_address(ethereum, "THRESHOLD_REGISTRY", directory_address)
                .await?,
            registry_coordinator: get_address(ethereum, "REGISTRY_COORDINATOR", directory_address)
                .await?,
            service_manager: get_address(ethereum, "SERVICE_MANAGER", directory_address).await?,
            bls_apk_registry: get_address(ethereum, "BLS_APK_REGISTRY", directory_address).await?,
            stake_registry: get_address(ethereum, "STAKE_REGISTRY", directory_address).await?,
            cert_verifier_router: match cert_verifier_router_address {
                Some(addr) => addr,
                None => get_address(ethereum, "CERT_VERIFIER_ROUTER", directory_address).await?,
            },
            delegation_manager: get_address(ethereum, "DELEGATION_MANAGER", directory_address)
                .await?,
        };

        Ok(eigen_da_contracts)
    }
}

/// The function performs a contract call to the EigenDA contract directory
/// to look up an address associated with a given contract name. It uses the
/// `getAddress` function from the directory contract.
async fn get_address(
    ethereum: &DynProvider,
    name: &'static str,
    directory_address: Address,
) -> Result<Address, RpcError<TransportErrorKind>> {
    let input = getAddressCall {
        name: name.to_string(),
    };

    let tx = TransactionRequest::default()
        .to(directory_address)
        .input(input.abi_encode().into());

    let src = ethereum.call(tx).await?;

    Ok(Address::from_slice(&src[12..32]))
}
