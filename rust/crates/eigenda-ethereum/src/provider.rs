use alloy_consensus::Header;
use alloy_primitives::{Address, U256};
use alloy_provider::network::Ethereum;
use alloy_provider::{DynProvider, PendingTransactionBuilder, Provider, ProviderBuilder};
use alloy_rpc_client::RpcClient;
use alloy_rpc_types_eth::{Block, BlockId, BlockNumberOrTag, TransactionRequest};
use alloy_signer_local::PrivateKeySigner;
use alloy_sol_types::sol;
use alloy_transport::layers::RetryBackoffLayer;
use alloy_transport::{RpcError, TransportErrorKind};
use eigenda_verification::cert::StandardCommitment;
use eigenda_verification::extraction::extractor::CERT_VERIFIER_ABNS_ARRAY_SLOT;
use eigenda_verification::extraction::{CertStateData, contract};
use futures::future::try_join_all;
use futures::{TryFutureExt, try_join};
use reth_trie_common::AccountProof;
use rustls::crypto::{CryptoProvider, aws_lc_rs};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
use tracing::instrument;

use crate::address::EthereumAddress;
use crate::contracts::{
    EIGENDA_DIRECTORY_HOODI, EIGENDA_DIRECTORY_INABOX, EIGENDA_DIRECTORY_MAINNET,
    EIGENDA_DIRECTORY_SEPOLIA, EigenDaContracts,
};

sol! {

    #[sol(rpc)]
    contract EigenDACertVerifierRouter {
        function getCertVerifierAt(uint32 referenceBlockNumber) external view returns (address);
        function certVerifierABNs(uint256 index) external view returns (uint32);
    }
}

/// Default maximal number of times we retry requests.
const DEFAULT_MAX_RETRY_TIMES: u32 = 10;

/// Default starting delay at which requests will be retried. In milliseconds.
const DEFAULT_INITIAL_BACKOFF: u64 = 1000;

/// Default compute units per second.
const DEFAULT_COMPUTE_UNITS: u64 = u64::MAX;

/// Network the adapter is running against.
#[derive(Debug, Clone, Copy, JsonSchema, PartialEq, Serialize, Deserialize)]
pub enum Network {
    /// Ethereum mainnet.
    Mainnet,
    /// Hoodi testnet.
    Hoodi,
    /// Sepolia testnet.
    Sepolia,
    /// Inabox local devnet.
    Inabox,
}

/// Configuration for the EigenDA Ethereum provider
///
/// # Required Traits
///
/// This type **must** implement [`JsonSchema`](schemars::JsonSchema) because it's used
/// in the Sovereign SDK's DA service configuration:
/// <https://github.com/Sovereign-Labs/sovereign-sdk/blob/e099285e0bae55812f35af3446240daca4470bf9/crates/rollup-interface/src/node/da.rs#L118>
#[derive(Debug, Clone, JsonSchema, PartialEq, Serialize, Deserialize)]
pub struct EigenDaProviderConfig {
    /// Network the adapter is running against.
    pub network: Network,

    /// URL of the Ethereum RPC node.
    pub rpc_url: String,

    /// Optional address of an EigenDACertVerifierRouter contract. See
    /// <https://layr-labs.github.io/eigenda/integration/spec/4-contracts.html#eigendacertverifierrouter>
    /// If None, the default EigenDA maintained Router for the selected network will be used.
    /// For a trustless integration, we strongly recommend that teams deploy and use their own Router contract. See
    /// <https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#upgradable-quorums-and-thresholds-for-optimistic-verification>
    /// for more details.
    pub cert_verifier_router_address: Option<EthereumAddress>,

    /// The number of compute units per second for the provider. Used in cases
    /// when the Ethereum node is hosted at providers like Alchemy that track
    /// compute units used when making a requests. If None, it means the node is
    /// not tracking compute units.
    pub compute_units: Option<u64>,

    /// The maximal number of times we retry requests to the node before
    /// returning the error.
    pub max_retry_times: Option<u32>,

    /// The initial backoff in milliseconds used when retrying Ethereum
    /// requests. It is increased on each subsequent retry.
    pub initial_backoff: Option<u64>,
}

/// Thin wrapper around the Alloy Ethereum provider with EigenDA-specific helpers.
#[derive(Debug, Clone)]
pub struct EigenDaProvider {
    /// Shared Alloy provider used for all Ethereum RPC calls.
    pub ethereum: DynProvider,

    /// EigenDA relevant contracts
    contracts: EigenDaContracts,
}

impl EigenDaProvider {
    /// Initialize the EigenDA Ethereum provider
    pub async fn new(
        config: &EigenDaProviderConfig,
        signer: PrivateKeySigner,
    ) -> Result<Self, RpcError<TransportErrorKind>> {
        let _ = CryptoProvider::install_default(aws_lc_rs::default_provider());

        let max_retry_times = config.max_retry_times.unwrap_or(DEFAULT_MAX_RETRY_TIMES);

        let backoff = config.initial_backoff.unwrap_or(DEFAULT_INITIAL_BACKOFF);

        let compute_units_per_second = config.compute_units.unwrap_or(DEFAULT_COMPUTE_UNITS);

        let retry_layer =
            RetryBackoffLayer::new(max_retry_times, backoff, compute_units_per_second);

        let client = RpcClient::builder()
            .layer(retry_layer)
            .connect(&config.rpc_url)
            .await?;

        let ethereum = ProviderBuilder::new()
            .wallet(signer)
            .connect_client(client)
            .erased();

        let directory_address = match config.network {
            Network::Mainnet => EIGENDA_DIRECTORY_MAINNET,
            Network::Hoodi => EIGENDA_DIRECTORY_HOODI,
            Network::Sepolia => EIGENDA_DIRECTORY_SEPOLIA,
            Network::Inabox => EIGENDA_DIRECTORY_INABOX,
        };

        let contracts = EigenDaContracts::new(
            &ethereum,
            directory_address,
            config.cert_verifier_router_address.map(|a| a.into()),
        )
        .await?;

        Ok(Self {
            ethereum,
            contracts,
        })
    }

    /// Broadcasts a transaction via the underlying Ethereum provider.
    pub async fn send_transaction(
        &self,
        tx: TransactionRequest,
    ) -> Result<PendingTransactionBuilder<Ethereum>, RpcError<TransportErrorKind>> {
        self.ethereum.send_transaction(tx).await
    }

    /// Fetches the block header for the given height if it exists.
    pub async fn fetch_ancestor(
        &self,
        block_height: u64,
    ) -> Result<Option<Header>, RpcError<TransportErrorKind>> {
        let block = self
            .ethereum
            .get_block_by_number(block_height.into())
            .await?;

        let header = block.map(|block| block.header.into_consensus());
        Ok(header)
    }

    /// Fetches a block by its number, including full transactions.
    pub async fn get_block_by_number(
        &self,
        number: BlockNumberOrTag,
    ) -> Result<Option<Block>, RpcError<TransportErrorKind>> {
        self.ethereum.get_block_by_number(number).full().await
    }

    /// Fetches a block by a [`BlockId`], returning full transaction data when available.
    pub async fn get_block(
        &self,
        block: BlockId,
    ) -> Result<Option<Block>, RpcError<TransportErrorKind>> {
        self.ethereum.get_block(block).await
    }

    /// Fetches all ABNs registered in the cert verifier router stored in self.contracts.
    // TODO(samlaf): we should add a function in the Router contract to fetch all abns at once.
    async fn get_router_abns(&self) -> Result<Vec<u32>, alloy_contract::Error> {
        let router =
            EigenDACertVerifierRouter::new(self.contracts.cert_verifier_router, &self.ethereum);

        let num_abns = self
            .ethereum
            .get_storage_at(
                self.contracts.cert_verifier_router,
                U256::from(CERT_VERIFIER_ABNS_ARRAY_SLOT),
            )
            .await?;

        let abn_futs = (0..num_abns.to::<u64>()).map(|i| {
            let router = router.clone();
            async move { router.certVerifierABNs(U256::from(i)).call().await }
        });

        let abns: Vec<u32> = try_join_all(abn_futs).await?;
        Ok(abns)
    }

    /// Fetches the address of the cert verifier active at a given reference block number
    /// according to the cert verifier router stored in self.contracts.
    async fn get_cert_verifier_at_rbn(
        &self,
        reference_block_number: u32,
    ) -> Result<Address, alloy_contract::Error> {
        let router =
            EigenDACertVerifierRouter::new(self.contracts.cert_verifier_router, &self.ethereum);

        let addr: Address = router
            .getCertVerifierAt(reference_block_number)
            .call()
            .await?;

        Ok(addr)
    }

    /// Fetches the relevant state used to validate the EigenDA certificate.
    ///
    /// See the contracts storage diagram in the [crate documentation](crate#contracts-storage-diagram)
    /// to get a visual understanding of the different pieces of state being fetched here.
    #[instrument(skip_all)]
    pub async fn fetch_cert_state(
        &self,
        block_height: u64,
        cert: &StandardCommitment,
    ) -> Result<CertStateData, alloy_contract::Error> {
        // First we extract all the cert-dependent storage slots from the registry contracts.
        let keys = contract::RegistryCoordinator::storage_keys(cert);
        let registry_coordinator_fut = self
            .ethereum
            .get_proof(self.contracts.registry_coordinator, keys)
            .number(block_height)
            .into_future()
            .map_err(alloy_contract::Error::TransportError);

        let keys = contract::EigenDaThresholdRegistry::storage_keys(cert);
        let threshold_registry_fut = self
            .ethereum
            .get_proof(self.contracts.threshold_registry, keys)
            .number(block_height)
            .into_future()
            .map_err(alloy_contract::Error::TransportError);

        let keys = contract::BlsApkRegistry::storage_keys(cert);
        let bls_apk_registry_fut = self
            .ethereum
            .get_proof(self.contracts.bls_apk_registry, keys)
            .number(block_height)
            .into_future()
            .map_err(alloy_contract::Error::TransportError);

        let keys = contract::StakeRegistry::storage_keys(cert);
        let stake_registry_fut = self
            .ethereum
            .get_proof(self.contracts.stake_registry, keys)
            .number(block_height)
            .into_future()
            .map_err(alloy_contract::Error::TransportError);

        let keys = contract::ServiceManager::storage_keys();
        let service_manager_fut = self
            .ethereum
            .get_proof(self.contracts.service_manager, keys)
            .number(block_height)
            .into_future()
            .map_err(alloy_contract::Error::TransportError);

        let keys = contract::DelegationManager::storage_keys();
        let delegation_manager_fut = self
            .ethereum
            .get_proof(self.contracts.delegation_manager, keys)
            .number(block_height)
            .into_future()
            .map_err(alloy_contract::Error::TransportError);

        let cert_verifier_router_fut = async {
            let abns = self.get_router_abns().await?;
            let keys = contract::EigenDaCertVerifierRouter::storage_keys(&abns);
            self.ethereum
                .get_proof(self.contracts.cert_verifier_router, keys)
                .number(block_height)
                .await
                .map_err(alloy_contract::Error::TransportError)
        };

        let cert_verifier_fut = async {
            let cert_verifier_addr = self
                // rbn is u32 but reference_block casts it to u64, so its safe to cast it back to u32 here.
                .get_cert_verifier_at_rbn(cert.reference_block() as u32)
                .await?;

            let keys = contract::EigenDaCertVerifier::storage_keys();
            self.ethereum
                .get_proof(cert_verifier_addr, keys)
                .number(block_height)
                .await
                .map_err(alloy_contract::Error::TransportError)
        };

        let (
            threshold_registry,
            registry_coordinator,
            service_manager,
            bls_apk_registry,
            stake_registry,
            delegation_manager,
            cert_verifier_router,
            cert_verifier,
        ) = try_join!(
            threshold_registry_fut,
            registry_coordinator_fut,
            service_manager_fut,
            bls_apk_registry_fut,
            stake_registry_fut,
            delegation_manager_fut,
            cert_verifier_router_fut,
            cert_verifier_fut,
        )?;

        Ok(CertStateData {
            threshold_registry: AccountProof::from(threshold_registry),
            registry_coordinator: AccountProof::from(registry_coordinator),
            service_manager: AccountProof::from(service_manager),
            bls_apk_registry: AccountProof::from(bls_apk_registry),
            stake_registry: AccountProof::from(stake_registry),
            delegation_manager: AccountProof::from(delegation_manager),
            cert_verifier_router: AccountProof::from(cert_verifier_router),
            cert_verifier: AccountProof::from(cert_verifier),
        })
    }
}

#[cfg(test)]
/// Testing utilities for Ethereum provider functionality.
pub mod tests {
    use std::borrow::Cow;

    use alloy_provider::RootProvider;
    use alloy_provider::ext::AnvilApi;
    use alloy_rpc_types::anvil::MineOptions;
    use testcontainers::core::{ContainerPort, WaitFor};
    use testcontainers::runners::AsyncRunner;
    use testcontainers::{ContainerAsync, Image};

    /// Start local ethereum development node.
    #[allow(dead_code)]
    pub async fn start_ethereum_dev_node(
        mining: MiningKind,
    ) -> Result<(String, ContainerAsync<AnvilNode>), anyhow::Error> {
        let container = AnvilNode::new(mining).start().await?;
        let host_port = container.get_host_port_ipv4(PORT).await?;
        let url = format!("http://127.0.0.1:{host_port}");

        Ok((url, container))
    }

    const NAME: &str = "ghcr.io/foundry-rs/foundry";
    const TAG: &str = "stable";
    const READY_MSG: &str = "Listening on";
    const PORT: ContainerPort = ContainerPort::Tcp(8548);

    /// Defines different mining modes for the Anvil test node.
    #[derive(Debug, Default, Clone, Copy)]
    pub enum MiningKind {
        /// Mining interval in seconds.
        #[allow(dead_code)]
        Interval(u64),
        /// Mine the block after each submitted transaction.
        #[default]
        EachTransaction,
        /// The blocks should be mined manually by the user.
        #[allow(dead_code)]
        Manual,
    }

    /// If node is started with [`MiningKind::Manual`]. We should use this
    /// function to advance the chain.
    #[allow(dead_code)]
    pub async fn mine_block(ethereum_rpc_url: &str, n_blocks: u64) -> Result<(), anyhow::Error> {
        let ethereum: RootProvider = RootProvider::connect(ethereum_rpc_url).await?;
        ethereum
            .evm_mine(Some(MineOptions::Options {
                timestamp: None,
                blocks: Some(n_blocks),
            }))
            .await?;

        Ok(())
    }

    /// AnvilNode image for testcontainers
    #[derive(Debug, Default)]
    pub struct AnvilNode {
        mining: MiningKind,
    }

    impl AnvilNode {
        /// Create a new AnvilNode with the specified mining configuration.
        pub fn new(mining: MiningKind) -> Self {
            Self { mining }
        }
    }

    impl Image for AnvilNode {
        fn name(&self) -> &str {
            NAME
        }

        fn tag(&self) -> &str {
            TAG
        }

        fn ready_conditions(&self) -> Vec<testcontainers::core::WaitFor> {
            vec![WaitFor::message_on_stdout(READY_MSG)]
        }

        fn expose_ports(&self) -> &[ContainerPort] {
            &[PORT]
        }

        fn cmd(&self) -> impl IntoIterator<Item = impl Into<Cow<'_, str>>> {
            let mining = match self.mining {
                MiningKind::Interval(interval) => format!("--block-time {interval}"),
                MiningKind::EachTransaction => "".to_string(), // This is set by default if no flag passed
                MiningKind::Manual => "--no-mining".to_string(),
            };

            let command = format!("anvil --host 0.0.0.0 --port {} {mining}", PORT.as_u16());
            std::iter::once(command)
        }
    }
}
