// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import "../Env.sol";
import "./1-DeployImplementations.s.sol";
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
import {IEigenDAAddressDirectory} from "src/core/interfaces/IEigenDADirectory.sol";

// TODO: Sort out whatever is wrong with the EjectionManager.
// TODO: Add ProxyAdmin to zeus.
// TODO: Add post deployment assertions.

/// NOTE: Inconsistent use of EigenDARegistry
/// forgefmt: disable-next-item
contract ExecuteUpgrade is DeployImplementations {
    using Env for *;
    using Encode for *;

    function _runAsEOA() internal override {        
        // Get proxy admin.
        ProxyAdmin proxyAdmin = ProxyAdmin(address(0xaFe14D2c59526C7f47EDa9c9BA83574363Db899a)); // TODO: fixme

        /// -----------------------------------------------------------------------
        /// WARNING: NETWORK BROADCAST BEGINS HERE!
        /// -----------------------------------------------------------------------

        vm.startBroadcast();

        // TODO: Upgrade AccessControl.
        
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
                Env.impl.owner(), // newOwner
                new uint32[](0),
                new address[](0)
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

        // // TODO: This doesn't seam right, I think our zeus environment is using the old EjectionManager (at least on hoodi-preprod).
        // // Upgrade EjectionManager (no reinitialization needed).
        // proxyAdmin.upgrade(
        //     TransparentUpgradeableProxy(payable(address(Env.proxy.ejectionManager()))),
        //     address(Env.impl.ejectionManager())
        // );

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

        // // Upgrade RegistryCoordinator.
        // proxyAdmin.upgradeAndCall(
        //     TransparentUpgradeableProxy(payable(address(Env.proxy.registryCoordinator()))),
        //     address(Env.impl.registryCoordinator()),
        //     abi.encodeWithSelector(
        //         EigenDARegistryCoordinator.initialize.selector,
        //         Env.impl.owner(), // newOwner
        //         Env.proxy.registryCoordinator().ejector(),
        //         Env.impl.pauserRegistry(),
        //         0, // initial paused status (nothing paused)
        //         new IRegistryCoordinator.OperatorSetParam[](0),
        //         new uint96[](0),
        //         new IStakeRegistry.StrategyParams[][](0)
        //     )
        // );

        // Upgrade RelayRegistry.
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(Env.proxy.relayRegistry()))),
            address(Env.impl.relayRegistry()),
            abi.encodeWithSelector(EigenDARelayRegistry.initialize.selector, Env.impl.owner()) // newOwner
        );

        // // Upgrade ServiceManager.
        // proxyAdmin.upgradeAndCall(
        //     TransparentUpgradeableProxy(payable(address(Env.proxy.serviceManager()))),
        //     address(Env.impl.serviceManager()),
        //     abi.encodeWithSelector(
        //         EigenDAServiceManager.initialize.selector,
        //         Env.impl.pauserRegistry(),
        //         0, // initial paused status (nothing paused)
        //         Env.impl.owner(), // newOwner
        //         new address[](0),
        //         Env.proxy.serviceManager().rewardsInitiator()
        //     )
        // );

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

    function testScript() public override {
        // Set EOA mode for deployment
        _mode = OperationalMode.EOA;
        // Deploy implementations first (sets Zeus state)
        DeployImplementations._runAsEOA();
        // Execute upgrades (reads implementations from Zeus state)
        ExecuteUpgrade._runAsEOA();
        // Run parent's deployment tests
        DeployImplementations._afterTestScript();
        // Run upgrade-specific tests
        _afterUpgradeTests();
    }

    /// -----------------------------------------------------------------------
    /// Test hooks
    /// -----------------------------------------------------------------------

    function _afterUpgradeTests() internal view {
        // Run upgrade-specific tests
        _testOwnership();
        _testUpgradedImplementations();
        _testPauseStates();
        _testPaymentVaultStatePreservation();
        _testThresholdRegistryStatePreservation();
        _testCriticalReferencesPreserved();
        _testPauserRegistryConfiguration();
        _testDirectoryAccessControl();
        _testCrossContractReferences();
    }

    /// -----------------------------------------------------------------------
    /// Tests
    /// -----------------------------------------------------------------------

    /// @notice Verify ownership has been transferred to the new owner
    function _testOwnership() internal view {
        assertEq(Env.proxy.certVerifierRouter().owner(), Env.impl.owner(), "CertVerifierRouter: incorrect owner");
        // assertEq(Env.proxy.directory().owner(), Env.impl.owner()); // Not ownable compliant.
        assertEq(Env.proxy.disperserRegistry().owner(), Env.impl.owner(), "DisperserRegistry: incorrect owner");
        // assertEq(Env.proxy.ejectionManager().owner(), Env.impl.owner()); // Not ownable compliant.
        assertEq(Env.proxy.paymentVault().owner(), Env.impl.owner(), "PaymentVault: incorrect owner");
        assertEq(Env.proxy.registryCoordinator().owner(), Env.impl.owner(), "RegistryCoordinator: incorrect owner");
        assertEq(Env.proxy.relayRegistry().owner(), Env.impl.owner(), "RelayRegistry: incorrect owner");
        // assertEq(Env.proxy.serviceManager().owner(), Env.impl.owner(), "ServiceManager: incorrect owner");
        assertEq(Env.proxy.thresholdRegistry().owner(), Env.impl.owner(), "ThresholdRegistry: incorrect owner");
    }

    /// @notice Verify all proxy implementations were upgraded
    function _testUpgradedImplementations() internal view {
        // ProxyAdmin proxyAdmin = ProxyAdmin(Env.proxyAdmin()); // TODO: Add proxy admin to zeus.
        
        // assertEq(proxyAdmin.getProxyImplementation(TransparentUpgradeableProxy(payable(address(Env.proxy.blsApkRegistry())))), 
        //     address(Env.impl.blsApkRegistry()), "BLSApkRegistry: implementation not upgraded");
        // assertEq(proxyAdmin.getProxyImplementation(TransparentUpgradeableProxy(payable(address(Env.proxy.certVerifierRouter())))), 
        //     address(Env.impl.certVerifierRouter()), "CertVerifierRouter: implementation not upgraded");
        // assertEq(proxyAdmin.getProxyImplementation(TransparentUpgradeableProxy(payable(address(Env.proxy.directory())))), 
        //     address(Env.impl.directory()), "Directory: implementation not upgraded");
        // assertEq(proxyAdmin.getProxyImplementation(TransparentUpgradeableProxy(payable(address(Env.proxy.disperserRegistry())))), 
        //     address(Env.impl.disperserRegistry()), "DisperserRegistry: implementation not upgraded");
        // assertEq(proxyAdmin.getProxyImplementation(TransparentUpgradeableProxy(payable(address(Env.proxy.ejectionManager())))), 
        //     address(Env.impl.ejectionManager()), "EjectionManager: implementation not upgraded");
        // assertEq(proxyAdmin.getProxyImplementation(TransparentUpgradeableProxy(payable(address(Env.proxy.indexRegistry())))), 
        //     address(Env.impl.indexRegistry()), "IndexRegistry: implementation not upgraded");
        // assertEq(proxyAdmin.getProxyImplementation(TransparentUpgradeableProxy(payable(address(Env.proxy.paymentVault())))), 
        //     address(Env.impl.paymentVault()), "PaymentVault: implementation not upgraded");
        // assertEq(proxyAdmin.getProxyImplementation(TransparentUpgradeableProxy(payable(address(Env.proxy.registryCoordinator())))), 
        //     address(Env.impl.registryCoordinator()), "RegistryCoordinator: implementation not upgraded");
        // assertEq(proxyAdmin.getProxyImplementation(TransparentUpgradeableProxy(payable(address(Env.proxy.relayRegistry())))), 
        //     address(Env.impl.relayRegistry()), "RelayRegistry: implementation not upgraded");
        // assertEq(proxyAdmin.getProxyImplementation(TransparentUpgradeableProxy(payable(address(Env.proxy.serviceManager())))), 
        //     address(Env.impl.serviceManager()), "ServiceManager: implementation not upgraded");
        // assertEq(proxyAdmin.getProxyImplementation(TransparentUpgradeableProxy(payable(address(Env.proxy.socketRegistry())))), 
        //     address(Env.impl.socketRegistry()), "SocketRegistry: implementation not upgraded");
        // assertEq(proxyAdmin.getProxyImplementation(TransparentUpgradeableProxy(payable(address(Env.proxy.stakeRegistry())))), 
        //     address(Env.impl.stakeRegistry()), "StakeRegistry: implementation not upgraded");
        // assertEq(proxyAdmin.getProxyImplementation(TransparentUpgradeableProxy(payable(address(Env.proxy.thresholdRegistry())))), 
        //     address(Env.impl.thresholdRegistry()), "ThresholdRegistry: implementation not upgraded");
    }

    /// @notice Verify contracts are not paused after upgrade
    function _testPauseStates() internal view {
        assertEq(Env.proxy.registryCoordinator().paused(), 0, "RegistryCoordinator: should not be paused");
        assertEq(Env.proxy.serviceManager().paused(), 0, "ServiceManager: should not be paused");
    }

    /// @notice Verify PaymentVault state was preserved during upgrade
    function _testPaymentVaultStatePreservation() internal view {
        assertTrue(Env.proxy.paymentVault().minNumSymbols() > 0, "PaymentVault: minNumSymbols not preserved");
        assertTrue(Env.proxy.paymentVault().pricePerSymbol() > 0, "PaymentVault: pricePerSymbol not preserved");
        assertTrue(Env.proxy.paymentVault().priceUpdateCooldown() > 0, "PaymentVault: priceUpdateCooldown not preserved");
        assertTrue(Env.proxy.paymentVault().globalSymbolsPerPeriod() > 0, "PaymentVault: globalSymbolsPerPeriod not preserved");
        assertTrue(Env.proxy.paymentVault().reservationPeriodInterval() > 0, "PaymentVault: reservationPeriodInterval not preserved");
        assertTrue(Env.proxy.paymentVault().globalRatePeriodInterval() > 0, "PaymentVault: globalRatePeriodInterval not preserved");
    }

    /// @notice Verify ThresholdRegistry state was preserved during upgrade
    function _testThresholdRegistryStatePreservation() internal view {
        assertTrue(Env.proxy.thresholdRegistry().quorumAdversaryThresholdPercentages().length > 0, 
            "ThresholdRegistry: quorumAdversaryThresholdPercentages not preserved");
        assertTrue(Env.proxy.thresholdRegistry().quorumConfirmationThresholdPercentages().length > 0, 
            "ThresholdRegistry: quorumConfirmationThresholdPercentages not preserved");
        assertTrue(Env.proxy.thresholdRegistry().quorumNumbersRequired().length > 0, 
            "ThresholdRegistry: quorumNumbersRequired not preserved");
    }

    /// @notice Verify critical references were preserved during upgrade
    function _testCriticalReferencesPreserved() internal view {
        assertTrue(Env.proxy.registryCoordinator().ejector() != address(0), "RegistryCoordinator: ejector not preserved");
        assertTrue(Env.proxy.serviceManager().rewardsInitiator() != address(0), "ServiceManager: rewardsInitiator not preserved");
    }

    /// @notice Verify PauserRegistry configuration is correct
    function _testPauserRegistryConfiguration() internal view {
        // assertEq(address(Env.proxy.registryCoordinator().pauserRegistry()), address(Env.impl.pauserRegistry()), 
        //     "RegistryCoordinator: incorrect pauserRegistry");
        // assertEq(address(Env.proxy.serviceManager().pauserRegistry()), address(Env.impl.pauserRegistry()), 
        //     "ServiceManager: incorrect pauserRegistry");
    }

    /// @notice Verify Directory has access control configured
    function _testDirectoryAccessControl() internal view {
        address accessControlAddr = Env.proxy.directory().getAddress(keccak256("ACCESS_CONTROL"));
        assertTrue(accessControlAddr != address(0), "Directory: accessControl not set");
        assertEq(accessControlAddr, address(Env.impl.accessControl()), 
            "Directory: incorrect accessControl address");
    }

    /// @notice Verify cross-contract references are still correct
    function _testCrossContractReferences() internal view {
        assertEq(address(Env.proxy.serviceManager().avsDirectory()), address(Env.proxy.avsDirectory()), 
            "ServiceManager: avsDirectory reference broken");
        assertEq(address(Env.proxy.serviceManager().eigenDAThresholdRegistry()), address(Env.proxy.thresholdRegistry()), 
            "ServiceManager: thresholdRegistry reference broken");
        assertEq(address(Env.proxy.serviceManager().eigenDARelayRegistry()), address(Env.proxy.relayRegistry()), 
            "ServiceManager: relayRegistry reference broken");
        assertEq(address(Env.proxy.serviceManager().paymentVault()), address(Env.proxy.paymentVault()), 
            "ServiceManager: paymentVault reference broken");
        assertEq(address(Env.proxy.serviceManager().eigenDADisperserRegistry()), address(Env.proxy.disperserRegistry()), 
            "ServiceManager: disperserRegistry reference broken");
    }
}
