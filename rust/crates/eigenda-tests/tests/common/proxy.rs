use std::borrow::Cow;
use std::time::Duration;

use eigenda_ethereum::provider::Network;
use testcontainers::core::{ContainerPort, WaitFor};
use testcontainers::runners::AsyncRunner;
use testcontainers::{ContainerAsync, Image, ImageExt};

const NAME: &str = "ghcr.io/layr-labs/eigenda-proxy";
const TAG: &str = "2.4.1";
// We use 3101 since inabox starts a proxy on 3100 already.
const PORT: ContainerPort = ContainerPort::Tcp(3101);
const READY_MSG: &str = "Started EigenDA Proxy REST ALT DA server";

/// Start the proxy server.
pub async fn start_proxy(
    network: Network,
    // In order to disperse payloads, signer_sk_hex must have a reservation and/or on-demand deposit in the PaymentVault contract.
    signer_sk_hex: &str,
) -> Result<(String, ContainerAsync<EigenDaProxy>), anyhow::Error> {
    let container = EigenDaProxy::new(network, signer_sk_hex)
        .with_startup_timeout(Duration::from_secs(30))
        // relay URLs are registered with localhost hostname, so we need to be on host network to access them.
        .with_network("host")
        .start()
        .await?;
    let url = format!("http://127.0.0.1:{PORT}");

    Ok((url, container))
}

/// EigenDAProxy image for testcontainers
#[derive(Debug)]
pub struct EigenDaProxy {
    cmd_args: Vec<String>,
}

impl EigenDaProxy {
    pub fn new(network: Network, signer_sk_hex: &str) -> Self {
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

        match network {
            Network::Sepolia => {
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
            Network::Inabox => {
                cmd_args.push("--eigenda.v2.eigenda-directory".to_string());
                cmd_args.push("0x1613beB3B2C4f22Ee086B2b38C1476A3cE7f78E8".to_string());
                cmd_args.push("--eigenda.v2.disperser-rpc".to_string());
                cmd_args.push("localhost:32005".to_string());
                cmd_args.push("--eigenda.v2.disable-tls".to_string());
                cmd_args.push(
                    "--eigenda.v2.cert-verifier-router-or-immutable-verifier-addr".to_string(),
                );
                // Local Inabox CertVerifier address
                cmd_args.push("0x99bbA657f2BbC93c02D617f8bA121cB8Fc104Acf".to_string());
                cmd_args.push("--eigenda.v2.eth-rpc".to_string());
                cmd_args.push("http://localhost:8545".to_string());
            }
            Network::Mainnet => {
                panic!("Mainnet network support not implemented");
            }
            Network::Hoodi => {
                panic!("Hoodi network support not implemented");
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
