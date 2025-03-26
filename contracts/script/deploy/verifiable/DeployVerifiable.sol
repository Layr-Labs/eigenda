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
import {PaymentVault} from "src/payments/PaymentVault.sol";
import {IPaymentVault} from "src/interfaces/IPaymentVault.sol";
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
import {IServiceManager} from "src/core/EigenDAServiceManager.sol";
import {VersionedBlobParams} from "src/interfaces/IEigenDAStructs.sol";

import {MockStakeRegistry} from "./mocks/MockStakeRegistry.sol";
import {MockRegistryCoordinator} from "./mocks/MockRegistryCoordinator.sol";

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
    CalldataServiceManagerParams
} from "./DeploymentInitializer.sol";

import "forge-std/Script.sol";
import "forge-std/StdToml.sol";

contract DeployVerifiable is Script {
    // The intended owner of all the contracts after the full deployment process is completed.
    address initialOwner;

    // All proxies and implementations to upgrade them to, namespaced in structs.
    ProxyAdmin proxyAdmin;
    DeployedAddresses proxies;
    DeployedAddresses implementations;

    // Contracts deployed without a proxy
    IPauserRegistry pauserRegistry;

    // Configuration parameters for construction of implementations
    address rewardsCoordinator;
    address avsDirectory;
    address delegationManager;
    address churnApprover;
    address ejector;

    uint256 initialPausedStatus;

    // Inert implementation contracts used as initial implementations before the proxies are initialized
    address emptyContract;
    address mockStakeRegistry;
    address mockRegistryCoordinator;

    // A contract that would be called after this deployment by the initialOwner to initialize the contracts.
    DeploymentInitializer deploymentInitializer;

    // Script config
    string configData;

    /// @dev override this if you don't want to use the environment to get the config path
    function _configPath() internal view virtual returns (string memory) {
        return vm.envString("DEPLOY_CONFIG_PATH");
    }

    function run() public {
        // Read JSON config
        configData = vm.readFile(_configPath());
        initialOwner = stdToml.readAddress(configData, ".initialOwner");
        rewardsCoordinator = stdToml.readAddress(configData, ".initParams.shared.rewardsCoordinator");
        avsDirectory = stdToml.readAddress(configData, ".initParams.shared.avsDirectory");
        delegationManager = stdToml.readAddress(configData, ".initParams.shared.delegationManager");
        initialPausedStatus = stdToml.readUint(configData, ".initParams.shared.initialPausedStatus");

        vm.startBroadcast();
        proxyAdmin = new ProxyAdmin();
        emptyContract = address(new EmptyContract());
        mockStakeRegistry = address(new MockStakeRegistry(IDelegationManager(delegationManager)));
        pauserRegistry = IPauserRegistry(
            new PauserRegistry(
                stdToml.readAddressArray(configData, ".initParams.core.pauserRegistry.pausers"),
                stdToml.readAddress(configData, ".initParams.core.pauserRegistry.unpauser")
            )
        );

        _deployInertProxies();
        _deployImplementations();

        // Deploy deployment initializer
        deploymentInitializer = new DeploymentInitializer(_immutableInitParams());

        // Transfer ownership of proxy admin to deployment initializer
        proxyAdmin.transferOwnership(address(deploymentInitializer));

        vm.stopBroadcast();

        _doTests(configData);
    }

    function _doTests(string memory cfg) internal {
        vm.startPrank(initialOwner);
        CalldataInitParams memory params = CalldataInitParamsLib.getCalldataInitParams(cfg);
        deploymentInitializer.initializeDeployment(params);

        // TODO: Add more tests

        vm.stopPrank();
    }
    function _deployInertProxies() internal virtual {
        proxyAdmin = new ProxyAdmin();
        emptyContract = address(new EmptyContract());
        mockStakeRegistry = address(new MockStakeRegistry(IDelegationManager(delegationManager)));

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
            new StakeRegistry(IRegistryCoordinator(proxies.registryCoordinator), IDelegationManager(delegationManager))
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
        implementations.paymentVault = address(new PaymentVault());
        implementations.disperserRegistry = address(new EigenDADisperserRegistry());
        implementations.serviceManager = address(
            new EigenDAServiceManager(
                IAVSDirectory(avsDirectory),
                IRewardsCoordinator(rewardsCoordinator),
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
        ImmutableRegistryCoordinatorParams memory registryCoordinatorParams = ImmutableRegistryCoordinatorParams({
            churnApprover: stdToml.readAddress(configData, ".initParams.middleware.registryCoordinator.churnApprover"),
            ejector: stdToml.readAddress(configData, ".initParams.middleware.registryCoordinator.ejector")
        });
        ImmutablePaymentVaultParams memory paymentVaultParams = ImmutablePaymentVaultParams({
            minNumSymbols: uint64(stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.minNumSymbols")),
            pricePerSymbol: uint64(stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.pricePerSymbol")),
            priceUpdateCooldown: uint64(
                stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.priceUpdateCooldown")
            ),
            globalSymbolsPerPeriod: uint64(
                stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.globalSymbolsPerPeriod")
            ),
            reservationPeriodInterval: uint64(
                stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.reservationPeriodInterval")
            ),
            globalRatePeriodInterval: uint64(
                stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.globalRatePeriodInterval")
            )
        });
        ImmutableServiceManagerParams memory serviceManagerParams = ImmutableServiceManagerParams({
            rewardsInitiator: stdToml.readAddress(configData, ".initParams.eigenDA.serviceManager.rewardsInitiator")
        });

        return ImmutableInitParams({
            proxyAdmin: proxyAdmin,
            initialOwner: initialOwner,
            pauserRegistry: pauserRegistry,
            initialPausedStatus: initialPausedStatus,
            proxies: proxies,
            implementations: implementations,
            registryCoordinatorParams: registryCoordinatorParams,
            paymentVaultParams: paymentVaultParams,
            serviceManagerParams: serviceManagerParams
        });
    }
}

library CalldataInitParamsLib {
    function operatorSetParams(string memory configData)
        internal
        pure
        returns (IRegistryCoordinator.OperatorSetParam[] memory)
    {
        bytes memory operatorConfigsRaw =
            stdToml.parseRaw(configData, ".initParams.middleware.registryCoordinator.operatorSetParams");
        return abi.decode(operatorConfigsRaw, (IRegistryCoordinator.OperatorSetParam[]));
    }

    function minimumStakes(string memory configData) internal pure returns (uint96[] memory) {
        uint256[] memory stakesConfigs256 =
            stdToml.readUintArray(configData, ".initParams.middleware.registryCoordinator.minimumStakes");
        uint96[] memory stakesConfigs = new uint96[](stakesConfigs256.length);
        for (uint256 i; i < stakesConfigs.length; i++) {
            stakesConfigs[i] = uint96(stakesConfigs[i]);
        }
        return stakesConfigs;
    }

    function strategyParams(string memory configData)
        internal
        pure
        returns (IStakeRegistry.StrategyParams[][] memory)
    {
        bytes memory strategyConfigsRaw =
            stdToml.parseRaw(configData, ".initParams.middleware.registryCoordinator.strategyParams");
        return abi.decode(strategyConfigsRaw, (IStakeRegistry.StrategyParams[][]));
    }

    function quorumAdversaryThresholdPercentages(string memory configData) internal pure returns (bytes memory) {
        return
            stdToml.readBytes(configData, ".initParams.eigenDA.thresholdRegistry.quorumAdversaryThresholdPercentages");
    }

    function quorumConfirmationThresholdPercentages(string memory configData) internal pure returns (bytes memory) {
        return stdToml.readBytes(
            configData, ".initParams.eigenDA.thresholdRegistry.quorumConfirmationThresholdPercentages"
        );
    }

    function quorumNumbersRequired(string memory configData) internal pure returns (bytes memory) {
        return stdToml.readBytes(configData, ".initParams.eigenDA.thresholdRegistry.quorumNumbersRequired");
    }

    function versionedBlobParams(string memory configData) internal pure returns (VersionedBlobParams[] memory) {
        bytes memory versionedBlobParamsRaw =
            stdToml.parseRaw(configData, ".initParams.eigenDA.thresholdRegistry.versionedBlobParams");
        return abi.decode(versionedBlobParamsRaw, (VersionedBlobParams[]));
    }

    function batchConfirmers(string memory configData) internal pure returns (address[] memory) {
        return stdToml.readAddressArray(configData, ".initParams.eigenDA.serviceManager.batchConfirmers");
    }

    function getCalldataInitParams(string memory configData) internal pure returns (CalldataInitParams memory) {
        return CalldataInitParams({
            registryCoordinatorParams: CalldataRegistryCoordinatorParams({
                operatorSetParams: operatorSetParams(configData),
                minimumStakes: minimumStakes(configData),
                strategyParams: strategyParams(configData)
            }),
            thresholdRegistryParams: CalldataThresholdRegistryParams({
                quorumAdversaryThresholdPercentages: quorumAdversaryThresholdPercentages(configData),
                quorumConfirmationThresholdPercentages: quorumConfirmationThresholdPercentages(configData),
                quorumNumbersRequired: quorumNumbersRequired(configData),
                versionedBlobParams: versionedBlobParams(configData)
            }),
            serviceManagerParams: CalldataServiceManagerParams({batchConfirmers: batchConfirmers(configData)})
        });
    }
}
