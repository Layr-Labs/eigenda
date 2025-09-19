#[cfg(feature = "native")]
pub mod provider;

use alloy_consensus::{EthereumTxEnvelope, Transaction, TxEip4844};
use tracing::instrument;

use crate::eigenda::cert::StandardCommitment;

/// Extract certificate from the transaction. Return None if no parsable
/// certificate exists.
#[instrument(skip_all)]
pub fn extract_certificate(tx: &EthereumTxEnvelope<TxEip4844>) -> Option<StandardCommitment> {
    let raw_cert = tx.as_eip1559()?.input();
    StandardCommitment::from_rlp_bytes(raw_cert).ok()
}

#[cfg(test)]
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

    #[derive(Debug, Default, Clone, Copy)]
    pub enum MiningKind {
        // Mining interval in seconds.
        #[allow(dead_code)]
        Interval(u64),
        // Mine the block after each submitted transaction.
        #[default]
        EachTransaction,
        // The blocks should be mined manually by the user.
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
