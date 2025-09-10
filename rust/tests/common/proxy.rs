use std::{borrow::Cow, collections::HashMap};

use testcontainers::{
    ContainerAsync, Image,
    core::{ContainerPort, WaitFor},
    runners::AsyncRunner,
};

use crate::common::SEQUENCER_SIGNER;

const NAME: &str = "ghcr.io/layr-labs/eigenda-proxy";
const TAG: &str = "2.2.1";
const READY_MSG: &str = "Started EigenDA proxy server";
const PORT: ContainerPort = ContainerPort::Tcp(3100);

#[derive(Debug)]
pub enum ProxyMode {
    /// Run the proxy with the in-memory backend that mocks a real backing EigenDA network.
    #[allow(dead_code)]
    InMemory,
    /// Run the proxy against the holesky network
    Holesky,
}

/// Start the proxy server.
pub async fn start_proxy() -> Result<(String, ContainerAsync<EigenDaProxy>), anyhow::Error> {
    let mode = ProxyMode::Holesky;
    let container = EigenDaProxy::new(mode).await.start().await?;
    let host_port = container.get_host_port_ipv4(PORT).await?;
    let url = format!("http://127.0.0.1:{host_port}");

    Ok((url, container))
}

/// EigenDAProxy image for testcontainers
#[derive(Debug)]
pub struct EigenDaProxy {
    env_vars: HashMap<String, String>,
}

impl EigenDaProxy {
    pub async fn new(mode: ProxyMode) -> Self {
        let env_vars = match mode {
            ProxyMode::InMemory => in_memory_env_vars(),
            ProxyMode::Holesky => holisky_env_vars().await,
        };

        Self { env_vars }
    }
}

impl Image for EigenDaProxy {
    fn name(&self) -> &str {
        NAME
    }

    fn tag(&self) -> &str {
        TAG
    }

    fn ready_conditions(&self) -> Vec<WaitFor> {
        vec![WaitFor::message_on_stdout(READY_MSG)]
    }

    fn env_vars(
        &self,
    ) -> impl IntoIterator<Item = (impl Into<Cow<'_, str>>, impl Into<Cow<'_, str>>)> {
        &self.env_vars
    }

    fn expose_ports(&self) -> &[ContainerPort] {
        &[PORT]
    }
}

// Config from https://github.com/Layr-Labs/eigenda/blob/e4bfcb45e3f9504f38b911ea34dc5d8dcb18cb99/api/proxy/.env.exampleV2.holesky
async fn holisky_env_vars() -> HashMap<String, String> {
    HashMap::from([
        ("EIGENDA_PROXY_PORT".to_string(), PORT.as_u16().to_string()),
        (
            "EIGENDA_PROXY_EIGENDA_V2_SIGNER_PRIVATE_KEY_HEX".to_string(),
            SEQUENCER_SIGNER.to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_V2_ETH_RPC".to_string(),
            "https://ethereum-holesky-rpc.publicnode.com".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_V2_GRPC_DISABLE_TLS".to_string(),
            "false".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_V2_DISABLE_POINT_EVALUATION".to_string(),
            "false".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_V2_PUT_RETRIES".to_string(),
            "3".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_V2_DISPERSE_BLOB_TIMEOUT".to_string(),
            "2m".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_V2_CERTIFY_BLOB_TIMEOUT".to_string(),
            "2m".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_V2_BLOB_STATUS_POLL_INTERVAL".to_string(),
            "1s".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_V2_CONTRACT_CALL_TIMEOUT".to_string(),
            "5s".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_V2_RELAY_TIMEOUT".to_string(),
            "5s".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_V2_VALIDATOR_TIMEOUT".to_string(),
            "2m".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_V2_BLOB_PARAMS_VERSION".to_string(),
            "0".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_V2_MAX_BLOB_LENGTH".to_string(),
            "1MiB".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_V2_NETWORK".to_string(),
            "holesky_testnet".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_V2_CERT_VERIFIER_ROUTER_OR_IMMUTABLE_VERIFIER_ADDR".to_string(),
            "0xd305aeBcdEc21D00fDF8796CE37d0e74836a6B6e".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_V2_RBN_RECENCY_WINDOW_SIZE".to_string(),
            "0".to_string(),
        ),
        (
            "EIGENDA_PROXY_STORAGE_BACKENDS_TO_ENABLE".to_string(),
            "V2".to_string(),
        ),
        (
            "EIGENDA_PROXY_STORAGE_DISPERSAL_BACKEND".to_string(),
            "V2".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_TARGET_CACHE_PATH".to_string(),
            "resources/SRSTables".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_TARGET_KZG_G1_PATH".to_string(),
            "resources/g1.point".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_TARGET_KZG_G2_PATH".to_string(),
            "resources/g2.point".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_TARGET_KZG_G2_TRAILING_PATH".to_string(),
            "resources/g2.trailing.point".to_string(),
        ),
        (
            "EIGENDA_PROXY_EIGENDA_CERT_VERIFICATION_DISABLED".to_string(),
            "false".to_string(),
        ),
        ("EIGENDA_PROXY_LOG_FORMAT".to_string(), "text".to_string()),
        ("EIGENDA_PROXY_LOG_LEVEL".to_string(), "INFO".to_string()),
        (
            "EIGENDA_PROXY_MEMSTORE_ENABLED".to_string(),
            "false".to_string(),
        ),
        (
            "EIGENDA_PROXY_MEMSTORE_EXPIRATION".to_string(),
            "25m0s".to_string(),
        ),
    ])
}

fn in_memory_env_vars() -> HashMap<String, String> {
    HashMap::from([
        ("EIGENDA_PROXY_PORT".to_string(), PORT.as_u16().to_string()),
        (
            "EIGENDA_PROXY_MEMSTORE_ENABLED".to_owned(),
            "true".to_string(),
        ),
    ])
}
