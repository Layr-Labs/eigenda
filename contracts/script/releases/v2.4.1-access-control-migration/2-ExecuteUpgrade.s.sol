// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import "../Env.sol";
import "./1-DeployImplementations.s.sol";
import {EOADeployer} from "zeus-templates/templates/EOADeployer.sol";
import {Encode} from "zeus-templates/utils/Encode.sol";
import {ProxyAdmin} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {
    IPauserRegistry
} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IPauserRegistry.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry} from "lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {Pausable} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/permissions/Pausable.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";

// TODO: Sort out whatever is wrong with the EjectionManager.
// TODO: Add ProxyAdmin to zeus.
// TODO: Add post deployment assertions.

/// NOTE: Inconsistent use of EigenDARegistry
contract ExecuteUpgrade is EOADeployer {
    using Env for *;
    using Encode for *;

    /// forgefmt: disable-next-item
    function _runAsEOA() internal override {
        // Get proxy admin.
        ProxyAdmin proxyAdmin = ProxyAdmin(Env.proxyAdmin());

        /// -----------------------------------------------------------------------
        /// WARNING: NETWORK BROADCAST BEGINS HERE!
        /// -----------------------------------------------------------------------

        vm.startBroadcast();

        // Upgrade BlsApkRegistry (no reinitialization needed).
        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(Env.proxy.blsApkRegistry()))),
            address(Env.impl.blsApkRegistry())
        );

        // Upgrade CertVerifierRouter.
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(Env.proxy.certVerifierRouter()))),
            address(Env.impl.certVerifierRouter()),
            abi.encodeWithSelector(
                EigenDACertVerifierRouter.initialize.selector,
                Env.impl.owner() // newOwner
            )
        );

        // NOTE: CertVerifier (Not a proxy no upgrade or initialization needed).

        // Upgrade Directory.
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(Env.proxy.directory()))),
            address(Env.impl.directory()),
            abi.encodeWithSelector(
                EigenDADirectory.initialize.selector,
                Env.impl.accessControl()
            )
        );

        // Upgrade DisperserRegistry.
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(Env.proxy.disperserRegistry()))),
            address(Env.impl.disperserRegistry()),
            abi.encodeWithSelector(
                EigenDADisperserRegistry.initialize.selector, 
                Env.impl.owner() // newOwner
            )
        );

        // TODO: This doesn't seam right, I think our zeus environment is using the old EjectionManager (at least on hoodi-preprod).
        // Upgrade EjectionManager (no reinitialization needed).
        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(Env.proxy.ejectionManager()))),
            address(Env.impl.ejectionManager())
        );

        // Upgrade IndexRegistry (no reinitialization needed).
        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(Env.proxy.indexRegistry()))), 
            address(Env.impl.indexRegistry())
        );

        // NOTE: OperatorStateRetriever (not a proxy no upgrade or initialization needed).

        // NOTE: PauserRegistry (not a proxy no upgrade or initialization needed).

        // Upgrade PaymentVault.
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(Env.proxy.paymentVault()))),
            address(Env.impl.paymentVault()),
            abi.encodeWithSelector(
                PaymentVault.initialize.selector,
                Env.impl.owner(), // newOwner
                Env.proxy.paymentVault().minNumSymbols(),
                Env.proxy.paymentVault().pricePerSymbol(),
                Env.proxy.paymentVault().priceUpdateCooldown(),
                Env.proxy.paymentVault().globalSymbolsPerPeriod(),
                Env.proxy.paymentVault().reservationPeriodInterval(),
                Env.proxy.paymentVault().globalRatePeriodInterval()
            )
        );

        // Upgrade RegistryCoordinator.
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(Env.proxy.registryCoordinator()))),
            address(Env.impl.registryCoordinator()),
            abi.encodeWithSelector(
                EigenDARegistryCoordinator.initialize.selector,
                Env.impl.owner(), // newOwner
                Env.proxy.registryCoordinator().ejector(),
                Env.impl.pauserRegistry(),
                0, // initial paused status (nothing paused)
                new IRegistryCoordinator.OperatorSetParam[](0),
                new uint96[](0),
                new IStakeRegistry.StrategyParams[][](0)
            )
        );

        // Upgrade RelayRegistry.
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(Env.proxy.relayRegistry()))),
            address(Env.impl.relayRegistry()),
            abi.encodeWithSelector(EigenDARelayRegistry.initialize.selector, Env.impl.owner()) // newOwner
        );

        // Upgrade ServiceManager.
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(Env.proxy.serviceManager()))),
            address(Env.impl.serviceManager()),
            abi.encodeWithSelector(
                EigenDAServiceManager.initialize.selector,
                Env.impl.pauserRegistry(),
                0, // initial paused status (nothing paused)
                Env.impl.owner(), // newOwner
                new address[](0),
                Env.proxy.serviceManager().rewardsInitiator()
            )
        );

        // Upgrade SocketRegistry (no reinitialization needed).
        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(Env.proxy.socketRegistry()))),
            address(Env.impl.socketRegistry())
        );

        // Upgrade StakeRegistry (no reinitialization needed).
        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(Env.proxy.stakeRegistry()))), address(Env.impl.stakeRegistry())
        );

        // Upgrade ThresholdRegistry.
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(Env.proxy.thresholdRegistry()))),
            address(Env.impl.thresholdRegistry()),
            abi.encodeWithSelector(
                EigenDAThresholdRegistry.initialize.selector,
                Env.impl.owner(), // newOwner
                Env.proxy.thresholdRegistry().quorumAdversaryThresholdPercentages(),
                Env.proxy.thresholdRegistry().quorumConfirmationThresholdPercentages(),
                Env.proxy.thresholdRegistry().quorumNumbersRequired(),
                new DATypesV1.VersionedBlobParams[](0) // no additional blobs needed
            )
        );

        vm.stopBroadcast();
    }

    /// -----------------------------------------------------------------------
    /// 2) Post-upgrade assertions
    /// -----------------------------------------------------------------------

    function testScript() public virtual {
        DeployImplementations deployImplementations = new DeployImplementations();
        // Hook for pre-test setup.
        _beforeTestScript();
        // Deploy implementations.
        deployImplementations.runAsEOA();
        // Execute upgrade.
        runAsEOA();
        // Hook for post-upgrade assertions.
        _afterTestScript();
    }

    /// -----------------------------------------------------------------------
    /// Test hooks
    /// -----------------------------------------------------------------------

    function _beforeTestScript() internal view {}

    function _afterTestScript() internal view {
        // Assert ownership has been transferred to the new owner.
        assertEq(Env.proxy.certVerifierRouter().owner(), Env.impl.owner());
        // assertEq(Env.proxy.directory().owner(), Env.impl.owner()); // Not ownable compliant.
        assertEq(Env.proxy.disperserRegistry().owner(), Env.impl.owner());
        // assertEq(Env.proxy.ejectionManager().owner(), Env.impl.owner()); // Not ownable compliant.
        assertEq(Env.proxy.paymentVault().owner(), Env.impl.owner());
        assertEq(Env.proxy.registryCoordinator().owner(), Env.impl.owner());
        assertEq(Env.proxy.relayRegistry().owner(), Env.impl.owner());
        assertEq(Env.proxy.serviceManager().owner(), Env.impl.owner());
        assertEq(Env.proxy.thresholdRegistry().owner(), Env.impl.owner());
    }
}
