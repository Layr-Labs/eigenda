use alloy_provider::{DynProvider, Provider, ProviderBuilder};
use alloy_rpc_client::RpcClient;
use alloy_signer_local::PrivateKeySigner;
use alloy_transport::layers::RetryBackoffLayer;
use rustls::crypto::{CryptoProvider, aws_lc_rs};

use crate::service::{EigenDaServiceError, config::EigenDaConfig};

/// Default maximal number of times we retry requests.
const DEFAULT_MAX_RETRY_TIMES: u32 = 10;
/// Default starting delay at which requests will be retried. In milliseconds.
const DEFAULT_INITIAL_BACKOFF: u64 = 1000;
/// Default compute units per second.
const DEFAULT_COMPUTE_UNITS: u64 = u64::MAX;

/// Initialize Ethereum provider used by the [`crate::service::EigenDaService`].
pub async fn init_ethereum_provider(
    config: &EigenDaConfig,
    signer: PrivateKeySigner,
) -> Result<DynProvider, EigenDaServiceError> {
    let _ = CryptoProvider::install_default(aws_lc_rs::default_provider());

    let max_retry_times = config
        .ethereum_max_retry_times
        .unwrap_or(DEFAULT_MAX_RETRY_TIMES);

    let backoff = config
        .ethereum_initial_backoff
        .unwrap_or(DEFAULT_INITIAL_BACKOFF);

    let compute_units_per_second = config
        .ethereum_compute_units
        .unwrap_or(DEFAULT_COMPUTE_UNITS);

    let retry_layer = RetryBackoffLayer::new(max_retry_times, backoff, compute_units_per_second);

    let client = RpcClient::builder()
        .layer(retry_layer)
        .connect(&config.ethereum_rpc_url)
        .await?;

    let provider = ProviderBuilder::new()
        .wallet(signer)
        .on_client(client)
        .erased();

    Ok(provider)
}
