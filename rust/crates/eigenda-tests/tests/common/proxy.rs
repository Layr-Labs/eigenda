use std::borrow::Cow;
use std::time::Duration;

use testcontainers::core::{ContainerPort, WaitFor};
use testcontainers::runners::AsyncRunner;
use testcontainers::{ContainerAsync, Image, ImageExt};

const NAME: &str = "ghcr.io/layr-labs/eigenda-proxy";
const TAG: &str = "2.4.1";
const PORT: ContainerPort = ContainerPort::Tcp(3100);
const READY_MSG: &str = "Started EigenDA Proxy REST ALT DA server";

// TODO(samlaf): add support for inabox
#[allow(dead_code)]
#[derive(Debug)]
pub enum ProxyNetwork {
    /// Run the proxy against the Sepolia network
    Sepolia,
}

/// Start the proxy server.
pub async fn start_proxy(
    mode: ProxyNetwork,
    signer_sk_hex: &str,
) -> Result<(String, ContainerAsync<EigenDaProxy>), anyhow::Error> {
    let container = EigenDaProxy::new(mode, signer_sk_hex)
        .with_startup_timeout(Duration::from_secs(30))
        .start()
        .await?;
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
    pub fn new(mode: ProxyNetwork, signer_sk_hex: &str) -> Self {
        let mut cmd_args = vec![
            "--port".to_string(),
            PORT.as_u16().to_string(),
            "--apis.enabled".to_string(),
            "standard".to_string(),
            "--storage.backends-to-enable".to_string(),
            "v2".to_string(),
            "--storage.dispersal-backend".to_string(),
            "v2".to_string(),
            "--eigenda.v2.signer-payment-key-hex".to_string(),
            signer_sk_hex.to_string(),
        ];

        match mode {
            ProxyNetwork::Sepolia => {
                cmd_args.push("--eigenda.v2.network".to_string());
                cmd_args.push("sepolia_testnet".to_string());
                cmd_args.push(
                    "--eigenda.v2.cert-verifier-router-or-immutable-verifier-addr".to_string(),
                );
                // Latest CertVerifier on the Router: https://sepolia.etherscan.io/address/0x17ec4112c4BbD540E2c1fE0A49D264a280176F0D#readProxyContract
                // TODO(samlaf): make this lib support router
                cmd_args.push("0x19a469Ddb7199c7EB9E40455978b39894BB90974".to_string());
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
