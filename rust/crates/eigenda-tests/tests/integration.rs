mod common;

use bytes::Bytes;
use dotenvy::dotenv;
use rand::RngCore;
use std::str::FromStr;
use tracing::info;

use crate::common::proxy::{ProxyNetwork, start_proxy};
use alloy_signer_local::LocalSigner;
use eigenda_ethereum::provider::{EigenDaProvider, EigenDaProviderConfig, Network};
use eigenda_proxy::{EigenDaProxyConfig, ProxyClient};
use eigenda_verification::verification::{
    blob::codec::decode_encoded_payload, cert, verify_blob, verify_cert_recency,
};

#[tokio::test]
async fn post_payload_to_proxy() {
    common::tracing::init_tracing();

    dotenv().ok();
    let signer_sk_hex = std::env::var("EIGENDA_SIGNER_PRIVATE_KEY_HEX")
        .expect("EIGENDA_SIGNER_PRIVATE_KEY_HEX env var must be exported or set in .env file");

    let (url, _container) = start_proxy(ProxyNetwork::Sepolia, &signer_sk_hex)
        .await
        .unwrap();
    info!(%url, "Started eigenda-proxy for testing");

    let proxy_client = ProxyClient::new(&EigenDaProxyConfig {
        url,
        min_retry_delay: None,
        max_retry_delay: None,
        max_retry_times: None,
    })
    .unwrap();
    info!(?proxy_client, "proxy-client initialized");

    let payload = {
        let mut payload = vec![0u8; 1024];
        rand::thread_rng().fill_bytes(&mut payload);
        Bytes::from(payload)
    };

    let std_commitment = proxy_client.store_payload(&payload).await.unwrap();
    info!(?std_commitment, "successfully submitted payload to proxy");

    // Setup Ethereum client
    // TODO(samlaf): would be ideal if we didn't need a signer.. since its only needed to submit certs to ethereum as a batcher would.
    // prob want to separate eigenda-ethereum crate into a reader and writer.
    let batcher_signer =
        LocalSigner::from_str("0x0000000000000000000000000000000000000000000000000000000000000001")
            .unwrap();
    let provider_config = EigenDaProviderConfig {
        network: Network::Sepolia,
        rpc_url: "wss://ethereum-sepolia-rpc.publicnode.com".to_string(),
        compute_units: None,
        max_retry_times: None,
        initial_backoff: None,
    };
    let provider = EigenDaProvider::new(&provider_config, batcher_signer.clone())
        .await
        .unwrap();

    let cert_state = provider
        .fetch_cert_state(std_commitment.reference_block(), &std_commitment)
        .await
        .unwrap();

    let rbn_u32: u32 = std_commitment.reference_block().try_into().unwrap();
    let cur_block_num = rbn_u32 + 1;

    let inputs = cert_state.extract(&std_commitment, cur_block_num).unwrap();

    // TODO(samlaf): should just encode it locally rather than needing to go through the proxy
    let encoded_payload = proxy_client
        .get_encoded_payload(&std_commitment)
        .await
        .unwrap();

    // or call verify_and_extract_payload() to do all of these together.
    verify_cert_recency(cur_block_num as u64, rbn_u32 as u64, 1_000_000).unwrap();
    cert::verify(inputs).unwrap();
    verify_blob(&std_commitment, &encoded_payload).unwrap();
    let _payload = decode_encoded_payload(&encoded_payload).unwrap();
}
