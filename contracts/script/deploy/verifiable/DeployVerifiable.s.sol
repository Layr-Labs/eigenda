// SPDX-License-Identifier: BUSL-1.1

pragma solidity =0.8.12;

import {EmptyContract} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/test/mocks/EmptyContract.sol";
import {ProxyAdmin, TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";

import {IDelegationManager} from
    "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IDelegationManager.sol";
import {ISocketRegistry, SocketRegistry} from "lib/eigenlayer-middleware/src/SocketRegistry.sol";
import {IIndexRegistry} from "lib/eigenlayer-middleware/src/interfaces/IIndexRegistry.sol";
import {IndexRegistry} from "lib/eigenlayer-middleware/src/IndexRegistry.sol";
import {IStakeRegistry, StakeRegistry} from "lib/eigenlayer-middleware/src/StakeRegistry.sol";
import {IBLSApkRegistry} from "lib/eigenlayer-middleware/src/interfaces/IBLSApkRegistry.sol";
import {BLSApkRegistry} from "lib/eigenlayer-middleware/src/BLSApkRegistry.sol";
import {RegistryCoordinator, IRegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {IEigenDAThresholdRegistry, EigenDAThresholdRegistry} from "src/core/EigenDAThresholdRegistry.sol";
import {IEigenDARelayRegistry, EigenDARelayRegistry} from "src/core/EigenDARelayRegistry.sol";
import {PaymentVault} from "src/core/PaymentVault.sol";
import {IPaymentVault} from "src/core/interfaces/IPaymentVault.sol";
import {IEigenDADisperserRegistry, EigenDADisperserRegistry} from "src/core/EigenDADisperserRegistry.sol";
import {EigenDAServiceManager, IServiceManager} from "src/core/EigenDAServiceManager.sol";
import {IAVSDirectory} from
    "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IAVSDirectory.sol";
import {IRewardsCoordinator} from
    "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IRewardsCoordinator.sol";
import {
    IPauserRegistry,
    PauserRegistry
} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/permissions/PauserRegistry.sol";
import {IServiceManager} from "lib/eigenlayer-middleware/src/interfaces/IServiceManager.sol";
import {EigenDATypesV2 as DATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";

import {MockStakeRegistry} from "./mocks/MockStakeRegistry.sol";
import {MockRegistryCoordinator} from "./mocks/MockRegistryCoordinator.sol";

import {BeforeVerifiableDeploymentInitialization} from "./test/BeforeVerifiableDeploymentInitialization.s.sol";

import {
    DeploymentInitializer,
    ImmutableInitParams,
    DeployedAddresses,
    ImmutableRegistryCoordinatorParams,
    ImmutablePaymentVaultParams,
    ImmutableServiceManagerParams,
    CalldataInitParams,
    CalldataRegistryCoordinatorParams,
    CalldataThresholdRegistryParams,
    CalldataServiceManagerParams,
    InitParamsLib
} from "./DeploymentInitializer.sol";

import "forge-std/Script.sol";
import {console2} from "forge-std/console2.sol";

contract DeployVerifiable is Script {
    using InitParamsLib for string;

    // All proxies and implementations to upgrade them to, namespaced in structs.
    ProxyAdmin proxyAdmin;
    DeployedAddresses proxies;
    DeployedAddresses implementations;

    // Contracts deployed without a proxy
    IPauserRegistry pauserRegistry;

    // Inert implementation contracts used as initial implementations before the proxies are initialized
    address emptyContract;
    address mockStakeRegistry;
    address mockRegistryCoordinator;

    // A contract that would be called after this deployment by the initialOwner to initialize the contracts.
    DeploymentInitializer deploymentInitializer;

    // Script config
    string cfg;

    function run() public {
        // Read JSON config
        _initConfig();

        vm.startBroadcast();
        proxyAdmin = new ProxyAdmin();
        emptyContract = address(new EmptyContract());
        mockStakeRegistry = address(new MockStakeRegistry(IDelegationManager(cfg.delegationManager())));
        pauserRegistry = IPauserRegistry(new PauserRegistry(cfg.pausers(), cfg.unpauser()));

        _deployInertProxies();
        _deployImplementations();

        // Deploy deployment initializer
        deploymentInitializer = new DeploymentInitializer(_immutableInitParams());

        // Transfer ownership of proxy admin to deployment initializer
        proxyAdmin.transferOwnership(address(deploymentInitializer));

        _logs();

        vm.stopBroadcast();

        _doTests();
    }

    /// @dev override this if you don't want to use the environment to get the config path
    function _initConfig() internal virtual {
        cfg = vm.readFile(vm.envString("DEPLOY_CONFIG_PATH"));
    }

    function _logs() internal virtual {
        console2.log("Deployment addresses: ");
        console2.log("Deployment Initializer: ", address(deploymentInitializer));
        console2.log("Empty Contract Implementation", emptyContract);
        console2.log("Mock Stake Registry Implementation", mockStakeRegistry);
        console2.log("Mock Registry Coordinator Implementation", mockRegistryCoordinator);

        console2.log(
            "\n\nAll other relevant deployment addresses should be queried from the DeploymentInitializer contract."
        );
    }

    /// @dev This function does the same tests that a verifier of the deployment should be doing after the deployment.
    function _doTests() internal {
        BeforeVerifiableDeploymentInitialization beforeTest = new BeforeVerifiableDeploymentInitialization();
        beforeTest.doBeforeInitializationTests(
            cfg, deploymentInitializer, emptyContract, mockStakeRegistry, mockRegistryCoordinator
        );

        vm.startPrank(cfg.initialOwner());
        CalldataInitParams memory params = InitParamsLib.calldataInitParams(cfg);
        deploymentInitializer.initializeDeployment(params);

        // TODO: Add more tests

        vm.stopPrank();
    }

    function _deployInertProxies() internal virtual {
        proxyAdmin = new ProxyAdmin();
        emptyContract = address(new EmptyContract());
        mockStakeRegistry = address(new MockStakeRegistry(IDelegationManager(cfg.delegationManager())));

        // Deploy empty contracts to get addresses
        proxies.indexRegistry = address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));
        // The mock stake registry is needed by the service manager contract's constructor
        proxies.stakeRegistry = address(new TransparentUpgradeableProxy(mockStakeRegistry, address(proxyAdmin), ""));
        proxies.socketRegistry = address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));
        proxies.blsApkRegistry = address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));
        // The mock registry coordinator is needed by the service manager contract's constructor.
        // It is deployed here after the needed proxy addresses are known.
        mockRegistryCoordinator = address(
            new MockRegistryCoordinator(IStakeRegistry(proxies.stakeRegistry), IBLSApkRegistry(proxies.blsApkRegistry))
        );
        proxies.registryCoordinator =
            address(new TransparentUpgradeableProxy(mockRegistryCoordinator, address(proxyAdmin), ""));
        proxies.thresholdRegistry = address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));
        proxies.relayRegistry = address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));
        proxies.paymentVault = address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));
        proxies.disperserRegistry = address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));
        proxies.serviceManager = address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));
    }

    function _deployImplementations() internal virtual {
        implementations.indexRegistry = address(new IndexRegistry(IRegistryCoordinator(proxies.registryCoordinator)));
        implementations.stakeRegistry = address(
            new StakeRegistry(
                IRegistryCoordinator(proxies.registryCoordinator), IDelegationManager(cfg.delegationManager())
            )
        );
        implementations.socketRegistry = address(new SocketRegistry(IRegistryCoordinator(proxies.registryCoordinator)));
        implementations.blsApkRegistry = address(new BLSApkRegistry(IRegistryCoordinator(proxies.registryCoordinator)));
        implementations.registryCoordinator = address(
            new RegistryCoordinator(
                IServiceManager(proxies.serviceManager),
                IStakeRegistry(proxies.stakeRegistry),
                IBLSApkRegistry(proxies.blsApkRegistry),
                IIndexRegistry(proxies.indexRegistry),
                ISocketRegistry(proxies.socketRegistry)
            )
        );
        implementations.thresholdRegistry = address(new EigenDAThresholdRegistry());
        implementations.relayRegistry = address(new EigenDARelayRegistry());
        implementations.paymentVault = address(new PaymentVault(100));
        implementations.disperserRegistry = address(new EigenDADisperserRegistry());
        implementations.serviceManager = address(
            new EigenDAServiceManager(
                IAVSDirectory(cfg.avsDirectory()),
                IRewardsCoordinator(cfg.rewardsCoordinator()),
                IRegistryCoordinator(proxies.registryCoordinator),
                IStakeRegistry(proxies.stakeRegistry),
                IEigenDAThresholdRegistry(proxies.thresholdRegistry),
                IEigenDARelayRegistry(proxies.relayRegistry),
                IPaymentVault(proxies.paymentVault),
                IEigenDADisperserRegistry(proxies.disperserRegistry)
            )
        );
    }

    function _immutableInitParams() internal view returns (ImmutableInitParams memory) {
        return ImmutableInitParams({
            proxyAdmin: proxyAdmin,
            initialOwner: cfg.initialOwner(),
            pauserRegistry: pauserRegistry,
            initialPausedStatus: cfg.initialPausedStatus(),
            proxies: proxies,
            implementations: implementations,
            registryCoordinatorParams: cfg.registryCoordinatorParams(),
            paymentVaultParams: cfg.paymentVaultParams(),
            serviceManagerParams: cfg.serviceManagerParams()
        });
    }
}
