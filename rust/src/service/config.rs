use schemars::JsonSchema;

use crate::spec::EthereumAddress;

/// Configuration for the [`crate::service::EigenDaService`].
#[derive(Debug, JsonSchema, PartialEq)]
pub struct EigenDaConfig {
    /// URL of the Ethereum RPC node.
    pub ethereum_rpc_url: String,
    /// The number of compute units per second for the provider. Used in cases
    /// when the Ethereum node is hosted at providers like Alchemy that track
    /// compute units used when making a requests. If None, it means the node is
    /// not tracking compute units.
    pub ethereum_compute_units: Option<u64>,
    /// The maximal number of times we retry requests to the node before
    /// returning the error.
    pub ethereum_max_retry_times: Option<u32>,
    /// The initial backoff in milliseconds used when retrying Ethereum
    /// requests. It is increased on each subsequent retry.
    pub ethereum_initial_backoff: Option<u64>,
    /// Maximal number of responses that the cache stores. The
    pub ethereum_max_cache_items: Option<u32>,
    /// URL of the EigenDA proxy node.
    pub proxy_url: String,
    /// The initial backoff in milliseconds used when retrying EigenDA proxy
    /// requests. It is increased on each subsequent retry.
    pub proxy_min_retry_delay: Option<u64>,
    /// The maximal backoff in milliseconds used when retrying EigenDA proxy requests.
    pub proxy_max_retry_delay: Option<u64>,
    /// The maximal number of times we retry requests to the EigenDA proxy
    /// before returning the error.
    pub proxy_max_retry_times: Option<u64>,
    /// Private key of the sequencer. The account with corresponding private key
    /// is used by the sequencer to persist the certificates to Ethereum.
    /// Expected private key in the HEX format.
    pub sequencer_signer: String,
    /// EigenDA relevant contracts.
    pub contracts: EigenDaContracts,
}

/// EigenDA relevant contracts.
#[derive(Debug, Clone, JsonSchema, PartialEq)]
pub struct EigenDaContracts {
    /// # Ethereum description
    ///
    /// Registry for EigenDA relay keys
    ///
    /// # Details
    ///
    /// The `relayKeyToInfo` mapping is read from it. See [crate::eigenda::verification::cert::types::Storage]'s `relay_key_to_relay_address`
    pub eigen_da_relay_registry: EthereumAddress,

    /// # Ethereum description
    ///
    /// The `EigenDAThresholdRegistry` contract.
    ///
    /// # Details
    ///
    /// The `versionedBlobParams` mapping is read from it. See [crate::eigenda::verification::cert::types::Storage]'s `versioned_blob_params`
    pub eigen_da_threshold_registry: EthereumAddress,

    /// # Ethereum description
    ///
    /// A `RegistryCoordinator` that has three registries:
    ///   1. a `StakeRegistry` that keeps track of operators' stakes
    ///   2. a `BLSApkRegistry` that keeps track of operators' BLS public keys and aggregate BLS public keys for each quorum
    ///   3. an `IndexRegistry` that keeps track of an ordered list of operators for each quorum
    ///
    /// # Details
    ///
    /// The quorumCount variable is read from it. See [eigenda_cert_verifier::types::Storage]'s `quorum_count`
    /// The _operatorBitmapHistory mapping is read from it. See [crate::eigenda::verification::cert::types::Storage]'s `operator_bitmap_history`
    /// The quorumUpdateBlockNumber mapping is read from it. See [crate::eigenda::verification::cert::types::Storage]'s `quorum_update_block_number`
    pub registry_coordinator: EthereumAddress,

    /// # Ethereum description
    ///
    /// Used for checking BLS aggregate signatures from the operators of a `BLSRegistry`.
    ///
    /// # Details
    ///
    /// The staleStakesForbidden variable is read from it. See [crate::eigenda::verification::cert::types::Storage]'s `stale_stakes_forbidden`
    pub bls_signature_checker: EthereumAddress,

    /// # Ethereum description
    ///
    /// This is the contract for delegation in EigenLayer. The main functionalities of this contract are
    ///   - enabling anyone to register as an operator in EigenLayer
    ///   - allowing operators to specify parameters related to stakers who delegate to them
    ///   - enabling any staker to delegate its stake to the operator of its choice (a given staker can only delegate to a single operator at a time)
    ///   - enabling a staker to undelegate its assets from the operator it is delegated to (performed as part of the withdrawal process, initiated through the StrategyManager)
    ///
    /// # Details
    ///
    /// The minWithdrawalDelayBlocks variable is read from it. See [crate::eigenda::verification::cert::types::Storage]'s `min_withdrawal_delay_blocks`
    pub delegation_manager: EthereumAddress,

    /// # Ethereum description
    ///
    /// The `BlsApkRegistry` contract.
    ///
    /// # Details
    ///
    /// The apkHistory mapping is read from it. See [crate::eigenda::verification::cert::types::Storage]'s `apk_history`
    pub bls_apk_registry: EthereumAddress,

    /// # Ethereum description
    ///
    /// A `Registry` that keeps track of stakes of operators for up to 256 quorums.
    /// Specifically, it keeps track of
    ///   1. The stake of each operator in all the quorums they are a part of for block ranges
    ///   2. The total stake of all operators in each quorum for block ranges
    ///   3. The minimum stake required to register for each quorum
    ///
    /// It allows an additional functionality (in addition to registering and deregistering) to update the stake of an operator.
    ///
    /// # Details
    ///
    /// The _totalStakeHistory mapping is read from it. See [crate::eigenda::verification::cert::types::Storage]'s `total_stake_history`
    /// The operatorStakeHistory mapping is read from it. See [crate::eigenda::verification::cert::types::Storage]'s `operator_stake_history`
    pub stake_registry: EthereumAddress,

    /// # Ethereum description
    ///
    /// A CertVerifier is an immutable contract that is used by a consumer to verify EigenDA blob certificates
    /// For V2 verification this contract is deployed with immutable security thresholds and required quorum numbers,
    /// to change these values or verification behavior a new CertVerifier must be deployed
    ///
    /// # Details
    ///
    /// The quorumNumbersRequiredV2 variable is read from it. See [crate::eigenda::verification::cert::CertVerificationInputs]'s `required_quorum_numbers`
    /// The securityThresholdsV2 variable is read from it. See [crate::eigenda::verification::cert::CertVerificationInputs]'s `security_thresholds`
    pub eigen_da_cert_verifier: EthereumAddress,
}
