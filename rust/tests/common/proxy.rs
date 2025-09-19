use std::borrow::Cow;

use testcontainers::core::{ContainerPort, WaitFor};
use testcontainers::runners::AsyncRunner;
use testcontainers::{ContainerAsync, Image};

use crate::common::SEQUENCER_SIGNER;

const NAME: &str = "ghcr.io/layr-labs/eigenda-proxy";
const TAG: &str = "2.2.1";
const READY_MSG: &str = "Started EigenDA proxy server";
const PORT: ContainerPort = ContainerPort::Tcp(3100);

#[allow(dead_code)]
#[derive(Debug)]
pub enum ProxyNetwork {
    /// Run the proxy with the in-memory backend that mocks a real backing EigenDA network.
    InMemory,
    /// Run the proxy against the Holesky network
    Holesky,
    /// Run the proxy against the Sepolia network
    Sepolia,
}

/// Start the proxy server.
pub async fn start_proxy(
    mode: ProxyNetwork,
) -> Result<(String, ContainerAsync<EigenDaProxy>), anyhow::Error> {
    let container = EigenDaProxy::new(mode).await.start().await?;
    let host_port = container.get_host_port_ipv4(PORT).await?;
    let url = format!("http://127.0.0.1:{host_port}");

    Ok((url, container))
}

/// EigenDAProxy image for testcontainers
#[derive(Debug)]
pub struct EigenDaProxy {
    cmd_args: Vec<String>,
}

impl EigenDaProxy {
    pub async fn new(mode: ProxyNetwork) -> Self {
        let mut cmd_args = vec![
            "--port".to_string(),
            PORT.as_u16().to_string(),
            "--storage.dispersal-backend".to_string(),
            "v2".to_string(),
            "--storage.backends-to-enable".to_string(),
            "v2".to_string(),
            "--eigenda.v2.signer-payment-key-hex".to_string(),
            SEQUENCER_SIGNER.to_string(),
        ];

        match mode {
            ProxyNetwork::InMemory => {
                cmd_args.push("--memstore.enabled".to_string());
                cmd_args.push("true".to_string());
            }
            ProxyNetwork::Holesky => {
                cmd_args.push("--eigenda.v2.network".to_string());
                cmd_args.push("holesky_testnet".to_string());
                cmd_args.push(
                    "--eigenda.v2.cert-verifier-router-or-immutable-verifier-addr".to_string(),
                );
                cmd_args.push("0x036bB27A1F03350bDcccF344b497Ef22604006a3".to_string());
                cmd_args.push("--eigenda.v2.eth-rpc".to_string());
                cmd_args.push("wss://ethereum-holesky-rpc.publicnode.com".to_string());
            }
            ProxyNetwork::Sepolia => {
                cmd_args.push("--eigenda.v2.network".to_string());
                cmd_args.push("sepolia_testnet".to_string());
                cmd_args.push(
                    "--eigenda.v2.cert-verifier-router-or-immutable-verifier-addr".to_string(),
                );
                cmd_args.push("0x58D2B844a894f00b7E6F9F492b9F43aD54Cd4429".to_string());
                cmd_args.push("--eigenda.v2.eth-rpc".to_string());
                cmd_args.push("wss://ethereum-sepolia-rpc.publicnode.com".to_string());
            }
        };

        Self { cmd_args }
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

    fn cmd(&self) -> impl IntoIterator<Item = impl Into<Cow<'_, str>>> {
        &self.cmd_args
    }

    fn expose_ports(&self) -> &[ContainerPort] {
        &[PORT]
    }
}
