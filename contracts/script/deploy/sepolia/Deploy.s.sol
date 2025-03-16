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
import "forge-std/StdJson.sol";

contract DeploySepolia is Script {
    string deployConfigPath = string(bytes("./script/deploy/sepolia/config/placeholder.config.json"));
    address initialOwner;

    ProxyAdmin proxyAdmin;

    address rewardsCoordinator;
    address avsDirectory;
    address delegationManager;
    address pauserRegistry;
    address churnApprover;
    address ejector;

    DeployedAddresses proxies;
    DeployedAddresses implementations;

    uint256 initialPausedStatus;

    address emptyContract;
    address mockStakeRegistry;
    address mockRegistryCoordinator;

    DeploymentInitializer deploymentInitializer;

    string configData;

    function run() public {
        // READ JSON CONFIG DATA
        configData = vm.readFile(deployConfigPath);
        rewardsCoordinator = stdJson.readAddress(configData, ".rewardsCoordinator");
        avsDirectory = stdJson.readAddress(configData, ".avsDirectory");
        delegationManager = stdJson.readAddress(configData, ".delegationManager");
        initialOwner = stdJson.readAddress(configData, ".initialOwner");
        pauserRegistry = stdJson.readAddress(configData, ".pauserRegistry");

        vm.startBroadcast();
        proxyAdmin = new ProxyAdmin();
        emptyContract = address(new EmptyContract());

        // Deploy empty contracts to get addresses
        proxies.indexRegistry = address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));

        // The mock stake registry is needed by the service manager contract's constructor
        mockStakeRegistry = address(new MockStakeRegistry(IDelegationManager(delegationManager)));
        proxies.stakeRegistry = address(new TransparentUpgradeableProxy(mockStakeRegistry, address(proxyAdmin), ""));
        proxies.socketRegistry = address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));
        proxies.blsApkRegistry = address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));

        // The mock registry coordinator is needed by the service manager contract's constructor
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

        // Deploy implementations
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

        // Deploy deployment initializer
        deploymentInitializer = new DeploymentInitializer(initParams());

        // Transfer ownership of proxy admin to deployment initializer
        proxyAdmin.transferOwnership(address(deploymentInitializer));

        vm.stopBroadcast();

        doTests(configData);
    }

    function doTests(string memory cfg) internal {
        vm.startPrank(initialOwner);
        CalldataInitParams memory params = CalldataInitParamsLib.getCalldataInitParams(cfg);
        deploymentInitializer.initializeDeployment(params);

        // TODO: Add more tests

        vm.stopPrank();
    }

    function initParams() internal view returns (ImmutableInitParams memory) {
        ImmutableRegistryCoordinatorParams memory registryCoordinatorParams = ImmutableRegistryCoordinatorParams({
            churnApprover: stdJson.readAddress(configData, ".registryCoordinator.churnApprover"),
            ejector: stdJson.readAddress(configData, ".registryCoordinator.ejector")
        });
        ImmutablePaymentVaultParams memory paymentVaultParams = ImmutablePaymentVaultParams({
            minNumSymbols: uint64(stdJson.readUint(configData, ".paymentVault.minNumSymbols")),
            pricePerSymbol: uint64(stdJson.readUint(configData, ".paymentVault.pricePerSymbol")),
            priceUpdateCooldown: uint64(stdJson.readUint(configData, ".paymentVault.priceUpdateCooldown")),
            globalSymbolsPerPeriod: uint64(stdJson.readUint(configData, ".paymentVault.globalSymbolsPerPeriod")),
            reservationPeriodInterval: uint64(stdJson.readUint(configData, ".paymentVault.reservationPeriodInterval")),
            globalRatePeriodInterval: uint64(stdJson.readUint(configData, ".paymentVault.globalRatePeriodInterval"))
        });
        ImmutableServiceManagerParams memory serviceManagerParams = ImmutableServiceManagerParams({
            rewardsInitiator: stdJson.readAddress(configData, ".serviceManager.rewardsInitiator")
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
        bytes memory operatorConfigsRaw = stdJson.parseRaw(configData, ".registryCoordinator.operatorSetParams");
        return abi.decode(operatorConfigsRaw, (IRegistryCoordinator.OperatorSetParam[]));
    }

    function minimumStakes(string memory configData) internal pure returns (uint96[] memory) {
        bytes memory stakesConfigsRaw = stdJson.parseRaw(configData, ".registryCoordinator.minimumStakes");
        return abi.decode(stakesConfigsRaw, (uint96[]));
    }

    function strategyParams(string memory configData)
        internal
        pure
        returns (IStakeRegistry.StrategyParams[][] memory)
    {
        bytes memory strategyConfigsRaw = stdJson.parseRaw(configData, ".registryCoordinator.strategyParams");
        return abi.decode(strategyConfigsRaw, (IStakeRegistry.StrategyParams[][]));
    }

    function quorumAdversaryThresholdPercentages(string memory configData) internal pure returns (bytes memory) {
        return stdJson.readBytes(configData, ".thresholdRegistry.quorumAdversaryThresholdPercentages");
    }

    function quorumConfirmationThresholdPercentages(string memory configData) internal pure returns (bytes memory) {
        return stdJson.readBytes(configData, ".thresholdRegistry.quorumConfirmationThresholdPercentages");
    }

    function quorumNumbersRequired(string memory configData) internal pure returns (bytes memory) {
        return stdJson.readBytes(configData, ".thresholdRegistry.quorumNumbersRequired");
    }

    function versionedBlobParams(string memory configData) internal pure returns (VersionedBlobParams[] memory) {
        bytes memory versionedBlobParamsRaw = stdJson.parseRaw(configData, ".thresholdRegistry.versionedBlobParams");
        return abi.decode(versionedBlobParamsRaw, (VersionedBlobParams[]));
    }

    function batchConfirmers(string memory configData) internal pure returns (address[] memory) {
        return stdJson.readAddressArray(configData, ".serviceManager.batchConfirmers");
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
