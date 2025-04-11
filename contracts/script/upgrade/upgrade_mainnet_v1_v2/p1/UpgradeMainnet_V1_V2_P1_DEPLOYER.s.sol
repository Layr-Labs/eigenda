// SPDX-License-Identifier: BUSL-1.1

pragma solidity =0.8.12;

import {ProxyAdmin, TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {IEigenDAThresholdRegistry, EigenDAThresholdRegistry} from "src/core/EigenDAThresholdRegistry.sol";
import {IEigenDADisperserRegistry, EigenDADisperserRegistry} from "src/core/EigenDADisperserRegistry.sol";
import {IEigenDARelayRegistry, EigenDARelayRegistry} from "src/core/EigenDARelayRegistry.sol";
import {PaymentVault} from "src/payments/PaymentVault.sol";
import {IPaymentVault} from "src/interfaces/IPaymentVault.sol";

import {EjectionManager} from "lib/eigenlayer-middleware/src/EjectionManager.sol";
import {IRegistryCoordinator, RegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {EigenDAServiceManager} from "src/core/EigenDAServiceManager.sol";

import {IStakeRegistry} from "lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {IServiceManager} from "lib/eigenlayer-middleware/src/interfaces/IServiceManager.sol";
import {IIndexRegistry} from "lib/eigenlayer-middleware/src/interfaces/IIndexRegistry.sol";
import {IBLSApkRegistry} from "lib/eigenlayer-middleware/src/interfaces/IBLSApkRegistry.sol";
import {ISocketRegistry} from "lib/eigenlayer-middleware/src/interfaces/ISocketRegistry.sol";
import {IAVSDirectory} from
    "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IAVSDirectory.sol";
import {VersionedBlobParams} from "src/interfaces/IEigenDAStructs.sol";
import {IRewardsCoordinator} from
    "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IRewardsCoordinator.sol";

import "forge-std/Script.sol";
import "forge-std/StdToml.sol";
import {console2} from "forge-std/console2.sol";

/// @title Upgrade Mainnet V1 to V2 Phase 1
contract UpgradeMainnet_V1_V2_P1_DEPLOYER is Script {
    using stdToml for string;

    struct InitParams {
        ExistingDeployments existing;
        ThresholdRegistryParams thresholdRegistry;
        PaymentVaultParams paymentVault;
    }

    struct ExistingDeployments {
        address daOpsMsig;
        address registryCoordinator;
        address stakeRegistry;
        address serviceManager;
        address blsApkRegistry;
        address indexRegistry;
        address socketRegistry;
        address avsDirectory;
        address rewardsCoordinator;
    }

    struct CreatedContracts {
        address proxyAdmin;
        address thresholdRegistry;
        address thresholdRegistryImpl;
        address relayRegistry;
        address relayRegistryImpl;
        address disperserRegistry;
        address disperserRegistryImpl;
        address paymentVault;
        address paymentVaultImpl;
        address registryCoordinatorImpl;
        address ejectionManagerImpl;
        address eigenDAServiceManagerImpl;
    }

    struct ThresholdRegistryParams {
        bytes quorumAdversaryThresholdPercentages;
        bytes quorumConfirmationThresholdPercentages;
        bytes quorumNumbersRequired;
        VersionedBlobParams[] versionedBlobParams;
    }

    struct PaymentVaultParams {
        uint64 minNumSymbols;
        uint64 pricePerSymbol;
        uint64 priceUpdateCooldown;
        uint64 globalSymbolsPerPeriod;
        uint64 reservationPeriodInterval;
        uint64 globalRatePeriodInterval;
    }

    CreatedContracts createdContracts;

    /// @dev override this if you don't want to use the environment to get the config path
    function _cfg() internal virtual returns (string memory) {
        return vm.readFile(vm.envString("UPGRADE_MAINNET_V1_V2_P1_CONFIG"));
    }

    function _initParams() internal virtual returns (InitParams memory) {
        string memory cfg = _cfg();
        return InitParams({
            existing: ExistingDeployments({
                daOpsMsig: cfg.readAddress(".initParams.existing.daOpsMsig"),
                registryCoordinator: cfg.readAddress(".initParams.existing.registryCoordinator"),
                stakeRegistry: cfg.readAddress(".initParams.existing.stakeRegistry"),
                serviceManager: cfg.readAddress(".initParams.existing.serviceManager"),
                blsApkRegistry: cfg.readAddress(".initParams.existing.blsApkRegistry"),
                indexRegistry: cfg.readAddress(".initParams.existing.indexRegistry"),
                socketRegistry: cfg.readAddress(".initParams.existing.socketRegistry"),
                avsDirectory: cfg.readAddress(".initParams.existing.avsDirectory"),
                rewardsCoordinator: cfg.readAddress(".initParams.existing.rewardsCoordinator")
            }),
            thresholdRegistry: ThresholdRegistryParams({
                quorumAdversaryThresholdPercentages: cfg.readBytes(
                    ".initParams.thresholdRegistry.quorumAdversaryThresholdPercentages"
                ),
                quorumConfirmationThresholdPercentages: cfg.readBytes(
                    ".initParams.thresholdRegistry.quorumConfirmationThresholdPercentages"
                ),
                quorumNumbersRequired: cfg.readBytes(".initParams.thresholdRegistry.quorumNumbersRequired"),
                versionedBlobParams: abi.decode(
                    cfg.parseRaw(".initParams.thresholdRegistry.versionedBlobParams"), (VersionedBlobParams[])
                )
            }),
            paymentVault: PaymentVaultParams({
                minNumSymbols: uint64(cfg.readUint(".initParams.paymentVault.minNumSymbols")),
                pricePerSymbol: uint64(cfg.readUint(".initParams.paymentVault.pricePerSymbol")),
                priceUpdateCooldown: uint64(cfg.readUint(".initParams.paymentVault.priceUpdateCooldown")),
                globalSymbolsPerPeriod: uint64(cfg.readUint(".initParams.paymentVault.globalSymbolsPerPeriod")),
                reservationPeriodInterval: uint64(cfg.readUint(".initParams.paymentVault.reservationPeriodInterval")),
                globalRatePeriodInterval: uint64(cfg.readUint(".initParams.paymentVault.globalRatePeriodInterval"))
            })
        });
    }

    function run() external {
        InitParams memory initParams = _initParams(); // doing all cheatcodes before broadcast

        vm.startBroadcast();

        _deployContracts(initParams);

        // transfer ownership of the proxy admin to the multisig. The deployer should own no contracts after this script is run.
        ProxyAdmin(createdContracts.proxyAdmin).transferOwnership(initParams.existing.daOpsMsig);

        vm.stopBroadcast();

        _logContracts();
    }

    function _logContracts() internal view {
        console2.log("DEPLOYED CONTRACTS");
        console2.log("ProxyAdmin: ", createdContracts.proxyAdmin);
        console2.log();

        console2.log("PROXIES");
        console2.log("ThresholdRegistry: ", createdContracts.thresholdRegistry);
        console2.log("RelayRegistry: ", createdContracts.relayRegistry);
        console2.log("DisperserRegistry: ", createdContracts.disperserRegistry);
        console2.log("PaymentVault: ", createdContracts.paymentVault);
        console2.log();

        console2.log("IMPLEMENTATIONS");
        console2.log("ThresholdRegistryImpl: ", createdContracts.thresholdRegistryImpl);
        console2.log("RelayRegistryImpl: ", createdContracts.relayRegistryImpl);
        console2.log("DisperserRegistryImpl: ", createdContracts.disperserRegistryImpl);
        console2.log("PaymentVaultImpl: ", createdContracts.paymentVaultImpl);
        console2.log("RegistryCoordinatorImpl: ", createdContracts.registryCoordinatorImpl);
        console2.log("EjectionManagerImpl: ", createdContracts.ejectionManagerImpl);
        console2.log("EigenDAServiceManagerImpl: ", createdContracts.eigenDAServiceManagerImpl);
    }

    function _deployContracts(InitParams memory initParams) internal {
        createdContracts.proxyAdmin = address(new ProxyAdmin());
        _deployThresholdRegistry(initParams);
        _deployRelayRegistry(initParams);
        _deployDisperserRegistry(initParams);
        _deployPaymentVault(initParams);
        _deployEjectionManager(initParams);
        _deployRegistryCoordinator(initParams);
        _deployEigenDAServiceManager(initParams);
    }

    function _deployThresholdRegistry(InitParams memory initParams) internal {
        createdContracts.thresholdRegistryImpl = address(new EigenDAThresholdRegistry());
        createdContracts.thresholdRegistry = address(
            new TransparentUpgradeableProxy(
                createdContracts.thresholdRegistryImpl,
                createdContracts.proxyAdmin,
                abi.encodeCall(
                    EigenDAThresholdRegistry.initialize,
                    (
                        initParams.existing.daOpsMsig,
                        initParams.thresholdRegistry.quorumAdversaryThresholdPercentages,
                        initParams.thresholdRegistry.quorumConfirmationThresholdPercentages,
                        initParams.thresholdRegistry.quorumNumbersRequired,
                        initParams.thresholdRegistry.versionedBlobParams
                    )
                )
            )
        );
    }

    function _deployRelayRegistry(InitParams memory initParams) internal {
        createdContracts.relayRegistryImpl = address(new EigenDARelayRegistry());
        createdContracts.relayRegistry = address(
            new TransparentUpgradeableProxy(
                createdContracts.relayRegistryImpl,
                createdContracts.proxyAdmin,
                abi.encodeCall(EigenDARelayRegistry.initialize, (initParams.existing.daOpsMsig))
            )
        );
    }

    function _deployDisperserRegistry(InitParams memory initParams) internal {
        createdContracts.disperserRegistryImpl = address(new EigenDADisperserRegistry());
        createdContracts.disperserRegistry = address(
            new TransparentUpgradeableProxy(
                createdContracts.disperserRegistryImpl,
                createdContracts.proxyAdmin,
                abi.encodeCall(EigenDADisperserRegistry.initialize, (initParams.existing.daOpsMsig))
            )
        );
    }

    function _deployPaymentVault(InitParams memory initParams) internal {
        createdContracts.paymentVaultImpl = address(new PaymentVault());
        createdContracts.paymentVault = address(
            new TransparentUpgradeableProxy(
                createdContracts.paymentVaultImpl,
                createdContracts.proxyAdmin,
                abi.encodeCall(
                    PaymentVault.initialize,
                    (
                        initParams.existing.daOpsMsig,
                        initParams.paymentVault.minNumSymbols,
                        initParams.paymentVault.pricePerSymbol,
                        initParams.paymentVault.priceUpdateCooldown,
                        initParams.paymentVault.globalSymbolsPerPeriod,
                        initParams.paymentVault.reservationPeriodInterval,
                        initParams.paymentVault.globalRatePeriodInterval
                    )
                )
            )
        );
    }

    function _deployEjectionManager(InitParams memory initParams) internal {
        createdContracts.ejectionManagerImpl = address(
            new EjectionManager(
                IRegistryCoordinator(initParams.existing.registryCoordinator),
                IStakeRegistry(initParams.existing.stakeRegistry)
            )
        );
    }

    function _deployRegistryCoordinator(InitParams memory initParams) internal {
        createdContracts.registryCoordinatorImpl = address(
            new RegistryCoordinator(
                IServiceManager(initParams.existing.serviceManager),
                IStakeRegistry(initParams.existing.stakeRegistry),
                IBLSApkRegistry(initParams.existing.blsApkRegistry),
                IIndexRegistry(initParams.existing.indexRegistry),
                ISocketRegistry(initParams.existing.socketRegistry)
            )
        );
    }

    function _deployEigenDAServiceManager(InitParams memory initParams) internal {
        createdContracts.eigenDAServiceManagerImpl = address(
            new EigenDAServiceManager(
                IAVSDirectory(initParams.existing.avsDirectory),
                IRewardsCoordinator(initParams.existing.rewardsCoordinator),
                IRegistryCoordinator(initParams.existing.registryCoordinator),
                IStakeRegistry(initParams.existing.stakeRegistry),
                IEigenDAThresholdRegistry(createdContracts.thresholdRegistry),
                IEigenDARelayRegistry(createdContracts.relayRegistry),
                IPaymentVault(createdContracts.paymentVault),
                IEigenDADisperserRegistry(createdContracts.disperserRegistry)
            )
        );
    }
}
