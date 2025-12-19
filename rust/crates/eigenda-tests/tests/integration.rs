//! Integration tests combining all other eigenda-related crates.

mod common;

use bytes::Bytes;
use dotenvy::dotenv;
use rand::RngCore;
use std::str::FromStr;
use tracing::info;

use crate::common::proxy::start_proxy;
use alloy_primitives::B256;
use alloy_signer_local::LocalSigner;
use eigenda_ethereum::provider::{EigenDaProvider, EigenDaProviderConfig, Network};
use eigenda_proxy::{EigenDaProxyConfig, ProxyClient};
use eigenda_verification::verification::verify_and_extract_payload;

#[tokio::test]
// #[ignore = "Test that runs against sepolia network"]
async fn post_payload_and_verify_returned_cert_sepolia() {
    common::tracing::init_tracing();

    dotenv().ok();
    let signer_sk_hex = std::env::var("SEPOLIA_EIGENDA_SIGNER_PRIVATE_KEY_HEX").expect(
        "SEPOLIA_EIGENDA_SIGNER_PRIVATE_KEY_HEX env var must be exported or set in .env file",
    );
    let rpc_url = "wss://ethereum-sepolia-rpc.publicnode.com".to_string();

    post_payload_and_verify_returned_cert(Network::Sepolia, &signer_sk_hex, rpc_url).await;
}

#[tokio::test]
#[ignore = "Test that runs against inabox"]
async fn post_payload_and_verify_returned_cert_inabox() {
    common::tracing::init_tracing();

    dotenv().ok();
    // Inabox local dev signer private key, which matches the public key registered in:
    // https://github.com/Layr-Labs/eigenda/blob/bff1f8ab9c1841e6d05bc61225f66cfff508b751/contracts/script/SetUpEigenDA.s.sol#L168
    // It is safe to use for local development and testing only. Do not use this key in production or any other context.
    let signer_sk_hex =
        "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded".to_string();
    let rpc_url = "http://localhost:8545".to_string();
    post_payload_and_verify_returned_cert(Network::Inabox, &signer_sk_hex, rpc_url).await;
}

async fn post_payload_and_verify_returned_cert(
    network: Network,
    signer_sk_hex: &str,
    rpc_url: String,
) {
    let (url, _container) = start_proxy(network, signer_sk_hex).await.unwrap();
    info!(%url, "Started eigenda-proxy for testing");

    let proxy_client = ProxyClient::new(&EigenDaProxyConfig {
        url,
        min_retry_delay: None,
        max_retry_delay: None,
        max_retry_times: None,
    })
    .unwrap();

    let payload = {
        let mut payload = vec![0u8; 1024];
        rand::thread_rng().fill_bytes(&mut payload);
        Bytes::from(payload)
    };

    let std_commitment = proxy_client.store_payload(&payload).await.unwrap();

    // Setup Ethereum client
    // TODO(samlaf): would be ideal if we didn't need a signer.. since its only needed to submit certs to ethereum as a batcher would.
    // prob want to separate eigenda-ethereum crate into a reader and writer.
    let batcher_signer =
        LocalSigner::from_str("0x0000000000000000000000000000000000000000000000000000000000000001")
            .unwrap();
    let provider_config = EigenDaProviderConfig {
        network,
        rpc_url,
        cert_verifier_router_address: None,
        compute_units: None,
        max_retry_times: None,
        initial_backoff: None,
    };
    let provider = EigenDaProvider::new(&provider_config, batcher_signer.clone())
        .await
        .unwrap();

    let rbn = std_commitment.reference_block();
    // we pretend the std commitment was posted to a rollup's inbox 100 blocks after the reference block.
    let inclusion_block_num = rbn + 100;
    let recency_window = 1_000;

    let cert_state = provider
        .fetch_cert_state(std_commitment.reference_block(), &std_commitment)
        .await
        .unwrap();
    let rbn_state_root = provider
        .get_block_by_number(rbn.into())
        .await
        .unwrap()
        .unwrap()
        .header
        .state_root;

    // TODO(samlaf): should just encode it locally rather than needing to go through the proxy
    let encoded_payload = proxy_client
        .get_encoded_payload(&std_commitment)
        .await
        .unwrap();

    let _payload = verify_and_extract_payload(
        B256::ZERO,
        &std_commitment,
        Some(&cert_state),
        rbn_state_root,
        inclusion_block_num,
        rbn,
        recency_window,
        Some(&encoded_payload),
    )
    .unwrap()
    .unwrap();
}
