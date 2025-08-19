pub mod proxy;

use std::str::FromStr;

use sov_eigenda_adapter::{
    service::{
        EigenDaService, EigenDaServiceError,
        config::{EigenDaConfig, EigenDaContracts},
    },
    spec::{EthereumAddress, NamespaceId, RollupParams},
    verifier::EigenDaVerifier,
};
use sov_rollup_interface::da::DaVerifier;

pub static SEQUENCER_SIGNER: &str =
    "0x354945e623e9a9070ef2be9dec2a71c49784a6e8348f4bfb6ace91622df91d83";

// TODO: Change to custom addresses. These are from the dev net. The keys are known to public.
pub static ROLLUP_BATCH_NAMESPACE: &str = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8";
pub static ROLLUP_PROOF_NAMESPACE: &str = "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC";

// Contracts used. Instructions on how to retrieve them: https://docs.eigencloud.xyz/products/eigenda/networks/holesky#contract-addresses
pub static REGISTRY_COORDINATOR_ADDRESS: &str = "0x53012C69A189cfA2D9d29eb6F19B32e0A2EA3490";
pub static BLS_APK_REGISTRY_ADDRESS: &str = "0x066cF95c1bf0927124DFB8B02B401bc23A79730D";
pub static STAKE_REGISTRY_ADDRESS: &str = "0xBDACD5998989Eec814ac7A0f0f6596088AA2a270";
pub static SERVICE_MANAGER_ADDRESS: &str = "0xD4A7E1Bd8015057293f0D0A557088c286942e84b";
pub const EIGEN_DA_RELAY_REGISTRY_ADDRESS: &str = "0xaC8C6C7Ee7572975454E2f0b5c720f9E74989254";
pub const EIGEN_DA_THRESHOLD_REGISTRY_ADDRESS: &str = "0x76d131CFBD900dA12f859a363Fb952eEDD1d1Ec1";
pub const EIGEN_DA_CERT_VERIFIER_V2_ADDRESS: &str = "0xFe52fE1940858DCb6e12153E2104aD0fDFbE1162";
// ServiceManager is BlsSignatureChecker
pub const BLS_SIGNATURE_CHECKER_ADDRESS: &str = "0xD4A7E1Bd8015057293f0D0A557088c286942e84b";

pub async fn setup_adapter(
    proxy_url: String,
) -> Result<(EigenDaService, EigenDaVerifier), EigenDaServiceError> {
    let config = EigenDaConfig {
        ethereum_rpc_url: "wss://ethereum-holesky-rpc.publicnode.com".to_string(),
        sequencer_signer: SEQUENCER_SIGNER.to_string(),
        ethereum_compute_units: None,
        ethereum_max_retry_times: None,
        ethereum_initial_backoff: None,
        ethereum_max_cache_items: None,
        proxy_url,
        proxy_min_retry_delay: None,
        proxy_max_retry_delay: None,
        proxy_max_retry_times: None,
        contracts: EigenDaContracts {
            eigen_da_relay_registry: EthereumAddress::from_str(EIGEN_DA_RELAY_REGISTRY_ADDRESS)
                .unwrap(),
            eigen_da_threshold_registry: EthereumAddress::from_str(
                EIGEN_DA_THRESHOLD_REGISTRY_ADDRESS,
            )
            .unwrap(),
            registry_coordinator: EthereumAddress::from_str(REGISTRY_COORDINATOR_ADDRESS).unwrap(),
            bls_signature_checker: EthereumAddress::from_str(BLS_SIGNATURE_CHECKER_ADDRESS)
                .unwrap(),
            bls_apk_registry: EthereumAddress::from_str(BLS_APK_REGISTRY_ADDRESS).unwrap(),
            stake_registry: EthereumAddress::from_str(STAKE_REGISTRY_ADDRESS).unwrap(),
            delegation_manager: EthereumAddress::from_str(SERVICE_MANAGER_ADDRESS).unwrap(),
            eigen_da_cert_verifier: EthereumAddress::from_str(EIGEN_DA_CERT_VERIFIER_V2_ADDRESS)
                .unwrap(),
        },
    };
    let params = RollupParams {
        rollup_batch_namespace: NamespaceId::from_str(ROLLUP_BATCH_NAMESPACE).unwrap(),
        rollup_proof_namespace: NamespaceId::from_str(ROLLUP_PROOF_NAMESPACE).unwrap(),
        cert_recency_window: 3600,
    };

    let service = EigenDaService::new(config, params).await?;
    let verifier = EigenDaVerifier::new(params);

    Ok((service, verifier))
}
