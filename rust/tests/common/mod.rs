pub mod proxy;

use std::str::FromStr;

use sov_eigenda_adapter::{
    service::{
        EigenDaService, EigenDaServiceError,
        config::{EigenDaConfig, Network},
    },
    spec::{NamespaceId, RollupParams},
    verifier::EigenDaVerifier,
};
use sov_rollup_interface::da::DaVerifier;

pub static SEQUENCER_SIGNER: &str =
    "0x354945e623e9a9070ef2be9dec2a71c49784a6e8348f4bfb6ace91622df91d83";

// TODO: Change to custom addresses. These are from the dev net. The keys are known to public.
pub static ROLLUP_BATCH_NAMESPACE: &str = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8";
pub static ROLLUP_PROOF_NAMESPACE: &str = "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC";
pub static CERT_RECENCY_WINDOW: u64 = 3200;

pub async fn setup_adapter(
    proxy_url: String,
) -> Result<(EigenDaService, EigenDaVerifier), EigenDaServiceError> {
    let config = EigenDaConfig {
        network: Network::Holesky,
        ethereum_rpc_url: "wss://ethereum-holesky-rpc.publicnode.com".to_string(),
        sequencer_signer: SEQUENCER_SIGNER.to_string(),
        ethereum_compute_units: None,
        ethereum_max_retry_times: None,
        ethereum_initial_backoff: None,
        proxy_url,
        proxy_min_retry_delay: None,
        proxy_max_retry_delay: None,
        proxy_max_retry_times: None,
    };
    let params = RollupParams {
        rollup_batch_namespace: NamespaceId::from_str(ROLLUP_BATCH_NAMESPACE).unwrap(),
        rollup_proof_namespace: NamespaceId::from_str(ROLLUP_PROOF_NAMESPACE).unwrap(),
        cert_recency_window: CERT_RECENCY_WINDOW,
    };

    let service = EigenDaService::new(config, params).await?;
    let verifier = EigenDaVerifier::new(params);

    Ok((service, verifier))
}
