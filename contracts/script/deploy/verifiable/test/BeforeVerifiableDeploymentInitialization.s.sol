// SPDX-License-Identifier: BUSL-1.1

pragma solidity =0.8.12;

import {DeploymentInitializer, InitParamsLib} from "../DeploymentInitializer.sol";
import {ProxyAdmin, TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {IPauserRegistry} from
    "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IPauserRegistry.sol";
import {MockStakeRegistry} from "../mocks/MockStakeRegistry.sol";
import {MockRegistryCoordinator} from "../mocks/MockRegistryCoordinator.sol";
import {EmptyContract} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/test/mocks/EmptyContract.sol";
import {IndexRegistry} from "lib/eigenlayer-middleware/src/IndexRegistry.sol";
import {StakeRegistry} from "lib/eigenlayer-middleware/src/StakeRegistry.sol";
import {SocketRegistry} from "lib/eigenlayer-middleware/src/SocketRegistry.sol";
import {BLSApkRegistry} from "lib/eigenlayer-middleware/src/BLSApkRegistry.sol";
import {RegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {EigenDAServiceManager} from "src/core/EigenDAServiceManager.sol";

import "forge-std/Script.sol";
import "forge-std/StdJson.sol";
import "forge-std/Test.sol";

/// @dev This script is used to test the state of the contracts after the deployment initializer is deployed.
/// It is written as a script so that it does not run with this repo's test suite because of the dependency on env variables
/// (May be integrated into the test suite in the future)
contract BeforeVerifiableDeploymentInitialization is Script, Test {
    using InitParamsLib for string;

    string cfg;
    DeploymentInitializer deploymentInitializer;
    address emptyContract;
    address mockStakeRegistry;
    address mockRegistryCoordinator;

    function run() external virtual {
        _initConfigEnv();
        _testExpectedState();
    }

    function runWithConfig(
        string memory _cfg,
        DeploymentInitializer _deploymentInitializer,
        address _emptyContract,
        address _mockStakeRegistry,
        address _mockRegistryCoordinator
    ) external virtual {
        cfg = _cfg;
        deploymentInitializer = _deploymentInitializer;
        emptyContract = _emptyContract;
        mockStakeRegistry = _mockStakeRegistry;
        mockRegistryCoordinator = _mockRegistryCoordinator;
        _testExpectedState();
    }

    // For usage by the deployment script without needing to provide the deployment initializer in env.
    function doBeforeInitializationTests(
        string memory _cfg,
        DeploymentInitializer _deploymentInitializer,
        address _emptyContract,
        address _mockStakeRegistry,
        address _mockRegistryCoordinator
    ) external {
        cfg = _cfg;
        deploymentInitializer = _deploymentInitializer;
        emptyContract = _emptyContract;
        mockStakeRegistry = _mockStakeRegistry;
        mockRegistryCoordinator = _mockRegistryCoordinator;
    }

    function _testExpectedState() internal view {
        _testProxyAdmin();
        _testInitialOwner();
        _testInitialPauseStatus();
        _testPauserRegistry();
        _testIndexRegistry();
        _testStakeRegistry();
        _testSocketRegistry();
        _testBlsApkRegistry();
        _testRegistryCoordinator();
        _testThresholdRegistry();
        _testRelayRegistry();
        _testPaymentVault();
        _testDisperserRegistry();
        _testServiceManager();
    }

    function _testProxyAdmin() internal view {
        assertEq(deploymentInitializer.PROXY_ADMIN().owner(), address(deploymentInitializer));
    }

    function _testInitialOwner() internal view {
        assertEq(deploymentInitializer.INITIAL_OWNER(), cfg.initialOwner());
    }

    function _testInitialPauseStatus() internal view {
        assertEq(deploymentInitializer.INITIAL_PAUSED_STATUS(), cfg.initialPausedStatus());
    }

    function _testPauserRegistry() internal view {
        IPauserRegistry pauserRegistry = deploymentInitializer.PAUSER_REGISTRY();
        address[] memory pausers = cfg.pausers();
        for (uint256 i; i < pausers.length; i++) {
            /// @dev There is no way to check in this test that only the cfg pausers are pausers.
            assertTrue(pauserRegistry.isPauser(pausers[i]));
        }
        assertEq(pauserRegistry.unpauser(), cfg.unpauser());
    }

    function _testIndexRegistry() internal view {
        assertTrue(_proxyHasImplementation(deploymentInitializer.INDEX_REGISTRY(), emptyContract));
        IndexRegistry impl = IndexRegistry(deploymentInitializer.INDEX_REGISTRY_IMPL());
        assertEq(impl.registryCoordinator(), deploymentInitializer.REGISTRY_COORDINATOR());
    }

    function _testStakeRegistry() internal view {
        assertTrue(_proxyHasImplementation(deploymentInitializer.STAKE_REGISTRY(), mockStakeRegistry));
        StakeRegistry impl = StakeRegistry(deploymentInitializer.STAKE_REGISTRY_IMPL());
        assertEq(impl.registryCoordinator(), deploymentInitializer.REGISTRY_COORDINATOR());
        assertEq(address(impl.delegation()), cfg.delegationManager());
    }

    function _testSocketRegistry() internal view {
        assertTrue(_proxyHasImplementation(deploymentInitializer.SOCKET_REGISTRY(), emptyContract));
        SocketRegistry impl = SocketRegistry(deploymentInitializer.SOCKET_REGISTRY_IMPL());
        assertEq(impl.registryCoordinator(), deploymentInitializer.REGISTRY_COORDINATOR());
    }

    function _testBlsApkRegistry() internal view {
        assertTrue(_proxyHasImplementation(deploymentInitializer.BLS_APK_REGISTRY(), emptyContract));
        BLSApkRegistry impl = BLSApkRegistry(deploymentInitializer.BLS_APK_REGISTRY_IMPL());
        assertEq(impl.registryCoordinator(), deploymentInitializer.REGISTRY_COORDINATOR());
    }

    function _testRegistryCoordinator() internal view {
        assertTrue(_proxyHasImplementation(deploymentInitializer.REGISTRY_COORDINATOR(), mockRegistryCoordinator));
        RegistryCoordinator impl = RegistryCoordinator(deploymentInitializer.REGISTRY_COORDINATOR_IMPL());
        assertEq(address(impl.serviceManager()), deploymentInitializer.SERVICE_MANAGER());
        assertEq(address(impl.blsApkRegistry()), deploymentInitializer.BLS_APK_REGISTRY());
        assertEq(address(impl.stakeRegistry()), deploymentInitializer.STAKE_REGISTRY());
        assertEq(address(impl.indexRegistry()), deploymentInitializer.INDEX_REGISTRY());
        assertEq(address(impl.socketRegistry()), deploymentInitializer.SOCKET_REGISTRY());
    }

    function _testThresholdRegistry() internal view {
        assertTrue(_proxyHasImplementation(deploymentInitializer.THRESHOLD_REGISTRY(), emptyContract));
    }

    function _testRelayRegistry() internal view {
        assertTrue(_proxyHasImplementation(deploymentInitializer.RELAY_REGISTRY(), emptyContract));
    }

    function _testPaymentVault() internal view {
        assertTrue(_proxyHasImplementation(deploymentInitializer.PAYMENT_VAULT(), emptyContract));
    }

    function _testDisperserRegistry() internal view {
        assertTrue(_proxyHasImplementation(deploymentInitializer.DISPERSER_REGISTRY(), emptyContract));
    }

    function _testServiceManager() internal view {
        assertTrue(_proxyHasImplementation(deploymentInitializer.SERVICE_MANAGER(), emptyContract));
        EigenDAServiceManager impl = EigenDAServiceManager(deploymentInitializer.SERVICE_MANAGER_IMPL());
        assertEq(address(impl.avsDirectory()), cfg.avsDirectory());
        // Unfortunately, the service manager base contract does not expose a rewards coordinator getter!

        // assertEq(address(impl.rewardsCoordinator()), cfg.rewardsCoordinator());
        assertEq(address(impl.registryCoordinator()), deploymentInitializer.REGISTRY_COORDINATOR());
        assertEq(address(impl.stakeRegistry()), deploymentInitializer.STAKE_REGISTRY());
        assertEq(address(impl.eigenDAThresholdRegistry()), deploymentInitializer.THRESHOLD_REGISTRY());
        assertEq(address(impl.eigenDARelayRegistry()), deploymentInitializer.RELAY_REGISTRY());
        assertEq(address(impl.paymentVault()), deploymentInitializer.PAYMENT_VAULT());
        assertEq(address(impl.disperserRegistry()), deploymentInitializer.DISPERSER_REGISTRY());
    }

    function _proxyHasImplementation(address proxy, address implementation) internal view returns (bool) {
        ProxyAdmin proxyAdmin = deploymentInitializer.PROXY_ADMIN();
        return proxyAdmin.getProxyImplementation(TransparentUpgradeableProxy(payable(proxy))) == implementation;
    }

    /// @dev override this if you don't want to use the environment
    function _initConfigEnv() internal virtual {
        deploymentInitializer = DeploymentInitializer(vm.envAddress("DEPLOYMENT_INITIALIZER"));
        emptyContract = vm.envAddress("EMPTY_CONTRACT");
        mockStakeRegistry = vm.envAddress("MOCK_STAKE_REGISTRY");
        mockRegistryCoordinator = vm.envAddress("MOCK_REGISTRY_COORDINATOR");
        cfg = vm.readFile(vm.envString("DEPLOY_CONFIG_PATH"));
    }
}
