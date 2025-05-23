// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {IRegistryCoordinator, RegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {IPauserRegistry} from
    "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IPauserRegistry.sol";
import {IStakeRegistry} from "lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {IEigenDAThresholdRegistry, EigenDAThresholdRegistry} from "src/core/EigenDAThresholdRegistry.sol";
import {IEigenDARelayRegistry, EigenDARelayRegistry} from "src/core/EigenDARelayRegistry.sol";
import {PaymentVault} from "src/core/PaymentVault.sol";
import {IEigenDADisperserRegistry, EigenDADisperserRegistry} from "src/core/EigenDADisperserRegistry.sol";
import {EigenDAServiceManager} from "src/core/EigenDAServiceManager.sol";
import {IRewardsCoordinator} from
    "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IRewardsCoordinator.sol";
import {ProxyAdmin, TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "./DeploymentTypes.sol";

/**
 * @title DeploymentInitializer
 * @author Layr Labs, Inc.
 * @notice This contract is intended to be used by a multisig to verify a deployment.
 * It accomplishes this by storing all the statically-sized parameters needed to initialize the contracts.
 * The dynamically sized parameters are passed in as calldata by the multisig to avoid using storage and save deployment costs.
 */
contract DeploymentInitializer {
    bool public initialized;
    /// The ProxyAdmin contract that owns the proxies should be owned by this contract.
    ProxyAdmin public immutable PROXY_ADMIN;

    /// The owner of all the contracts and proxies after this contract is called.
    address public immutable INITIAL_OWNER;

    /// Initialization parameters that are shared between contracts
    uint256 public immutable INITIAL_PAUSED_STATUS;
    IPauserRegistry public immutable PAUSER_REGISTRY;

    /// Contracts that need to be upgraded and initialized.
    address public immutable INDEX_REGISTRY;
    address public immutable INDEX_REGISTRY_IMPL;

    address public immutable STAKE_REGISTRY;
    address public immutable STAKE_REGISTRY_IMPL;

    address public immutable SOCKET_REGISTRY;
    address public immutable SOCKET_REGISTRY_IMPL;

    address public immutable BLS_APK_REGISTRY;
    address public immutable BLS_APK_REGISTRY_IMPL;

    address public immutable REGISTRY_COORDINATOR;
    address public immutable REGISTRY_COORDINATOR_IMPL;

    address public immutable THRESHOLD_REGISTRY;
    address public immutable THRESHOLD_REGISTRY_IMPL;

    address public immutable RELAY_REGISTRY;
    address public immutable RELAY_REGISTRY_IMPL;

    address public immutable PAYMENT_VAULT;
    address public immutable PAYMENT_VAULT_IMPL;

    address public immutable DISPERSER_REGISTRY;
    address public immutable DISPERSER_REGISTRY_IMPL;

    address public immutable SERVICE_MANAGER;
    address public immutable SERVICE_MANAGER_IMPL;

    /// Registry Coordinator Immutables
    address public immutable CHURN_APPROVER;
    address public immutable EJECTOR;

    /// Payment Vault Immutables
    uint64 public immutable MIN_NUM_SYMBOLS;
    uint64 public immutable PRICE_PER_SYMBOL;
    uint64 public immutable PRICE_UPDATE_COOLDOWN;
    uint64 public immutable GLOBAL_SYMBOLS_PER_PERIOD;
    uint64 public immutable RESERVATION_PERIOD_INTERVAL;
    uint64 public immutable GLOBAL_RATE_PERIOD_INTERVAL;

    /// Service Manager Immutables
    address public immutable REWARDS_INITIATOR;

    constructor(ImmutableInitParams memory initParams) {
        {
            PROXY_ADMIN = initParams.proxyAdmin;
            INITIAL_OWNER = initParams.initialOwner;
            PAUSER_REGISTRY = initParams.pauserRegistry;
            INITIAL_PAUSED_STATUS = initParams.initialPausedStatus;
        }
        {
            // Proxies
            INDEX_REGISTRY = initParams.proxies.indexRegistry;
            STAKE_REGISTRY = initParams.proxies.stakeRegistry;
            SOCKET_REGISTRY = initParams.proxies.socketRegistry;
            BLS_APK_REGISTRY = initParams.proxies.blsApkRegistry;
            REGISTRY_COORDINATOR = initParams.proxies.registryCoordinator;
            THRESHOLD_REGISTRY = initParams.proxies.thresholdRegistry;
            RELAY_REGISTRY = initParams.proxies.relayRegistry;
            PAYMENT_VAULT = initParams.proxies.paymentVault;
            DISPERSER_REGISTRY = initParams.proxies.disperserRegistry;
            SERVICE_MANAGER = initParams.proxies.serviceManager;
        }
        {
            // Implementations
            INDEX_REGISTRY_IMPL = initParams.implementations.indexRegistry;
            STAKE_REGISTRY_IMPL = initParams.implementations.stakeRegistry;
            SOCKET_REGISTRY_IMPL = initParams.implementations.socketRegistry;
            BLS_APK_REGISTRY_IMPL = initParams.implementations.blsApkRegistry;
            REGISTRY_COORDINATOR_IMPL = initParams.implementations.registryCoordinator;
            THRESHOLD_REGISTRY_IMPL = initParams.implementations.thresholdRegistry;
            RELAY_REGISTRY_IMPL = initParams.implementations.relayRegistry;
            PAYMENT_VAULT_IMPL = initParams.implementations.paymentVault;
            DISPERSER_REGISTRY_IMPL = initParams.implementations.disperserRegistry;
            SERVICE_MANAGER_IMPL = initParams.implementations.serviceManager;
        }
        {
            // Registry Coordinator
            CHURN_APPROVER = initParams.registryCoordinatorParams.churnApprover;
            EJECTOR = initParams.registryCoordinatorParams.ejector;
        }
        {
            // Payment Vault
            MIN_NUM_SYMBOLS = initParams.paymentVaultParams.minNumSymbols;
            PRICE_PER_SYMBOL = initParams.paymentVaultParams.pricePerSymbol;
            PRICE_UPDATE_COOLDOWN = initParams.paymentVaultParams.priceUpdateCooldown;
            GLOBAL_SYMBOLS_PER_PERIOD = initParams.paymentVaultParams.globalSymbolsPerPeriod;
            RESERVATION_PERIOD_INTERVAL = initParams.paymentVaultParams.reservationPeriodInterval;
            GLOBAL_RATE_PERIOD_INTERVAL = initParams.paymentVaultParams.globalRatePeriodInterval;
        }
        {
            // Service Manager
            REWARDS_INITIATOR = initParams.serviceManagerParams.rewardsInitiator;
        }
    }

    function upgrade(address proxy, address implementation) internal {
        PROXY_ADMIN.upgrade(TransparentUpgradeableProxy(payable(proxy)), implementation);
    }

    /**
     * @dev It should be verified that the deployer has already done the following:
     *      - Deployed the implementation contracts
     *      - Deployed the proxy contracts with empty implementations
     *      - Set the proxy admin owner to this contract
     */
    function initializeDeployment(CalldataInitParams calldata initParams) external {
        require(msg.sender == INITIAL_OWNER, "Only the owner can initialize the deployment");
        require(!initialized, "Already initialized");

        upgrade(INDEX_REGISTRY, INDEX_REGISTRY_IMPL);
        upgrade(STAKE_REGISTRY, STAKE_REGISTRY_IMPL);
        upgrade(SOCKET_REGISTRY, SOCKET_REGISTRY_IMPL);
        upgrade(BLS_APK_REGISTRY, BLS_APK_REGISTRY_IMPL);
        upgrade(REGISTRY_COORDINATOR, REGISTRY_COORDINATOR_IMPL);
        RegistryCoordinator(REGISTRY_COORDINATOR).initialize(
            INITIAL_OWNER,
            CHURN_APPROVER,
            EJECTOR,
            IPauserRegistry(PAUSER_REGISTRY),
            INITIAL_PAUSED_STATUS,
            initParams.registryCoordinatorParams.operatorSetParams,
            initParams.registryCoordinatorParams.minimumStakes,
            initParams.registryCoordinatorParams.strategyParams
        );

        upgrade(THRESHOLD_REGISTRY, THRESHOLD_REGISTRY_IMPL);
        EigenDAThresholdRegistry(THRESHOLD_REGISTRY).initialize(
            INITIAL_OWNER,
            initParams.thresholdRegistryParams.quorumAdversaryThresholdPercentages,
            initParams.thresholdRegistryParams.quorumConfirmationThresholdPercentages,
            initParams.thresholdRegistryParams.quorumNumbersRequired,
            initParams.thresholdRegistryParams.versionedBlobParams
        );

        upgrade(RELAY_REGISTRY, RELAY_REGISTRY_IMPL);
        EigenDARelayRegistry(RELAY_REGISTRY).initialize(INITIAL_OWNER);

        upgrade(PAYMENT_VAULT, PAYMENT_VAULT_IMPL);

        upgrade(DISPERSER_REGISTRY, DISPERSER_REGISTRY_IMPL);
        EigenDADisperserRegistry(DISPERSER_REGISTRY).initialize(INITIAL_OWNER);

        upgrade(SERVICE_MANAGER, SERVICE_MANAGER_IMPL);
        EigenDAServiceManager(SERVICE_MANAGER).initialize(
            PAUSER_REGISTRY,
            INITIAL_PAUSED_STATUS,
            INITIAL_OWNER,
            initParams.serviceManagerParams.batchConfirmers,
            REWARDS_INITIATOR
        );
        PROXY_ADMIN.transferOwnership(INITIAL_OWNER);

        initialized = true;
    }
}
