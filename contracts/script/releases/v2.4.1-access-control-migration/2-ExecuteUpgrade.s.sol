// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import "../Env.sol";
import {DeployImplementations} from "./1-DeployImplementations.s.sol";
import {MultisigBuilder} from "zeus-templates/templates/MultisigBuilder.sol";
import {Encode} from "zeus-templates/utils/Encode.sol";
import {ProxyAdmin} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

/// @title ExecuteUpgrade
/// @notice Execute upgrade of EigenDA implementations via timelock controller
contract ExecuteUpgrade is MultisigBuilder, DeployImplementations {
    using Env for *;
    using Encode for *;

    function _runAsMultisig() internal override prank(Env.executorMultisig()) {
        // Get proxy admin
        ProxyAdmin proxyAdmin = ProxyAdmin(Env.proxyAdmin());

        // Upgrade ServiceManager proxy
        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(Env.proxy.serviceManager()))),
            address(Env.impl.serviceManager())
        );

        // Upgrade RegistryCoordinator proxy
        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(Env.proxy.registryCoordinator()))),
            address(Env.impl.registryCoordinator())
        );

        // Upgrade other proxies as needed
        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(Env.proxy.thresholdRegistry()))),
            address(Env.impl.thresholdRegistry())
        );

        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(Env.proxy.relayRegistry()))), address(Env.impl.relayRegistry())
        );

        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(Env.proxy.disperserRegistry()))),
            address(Env.impl.disperserRegistry())
        );

        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(Env.proxy.paymentVault()))), address(Env.impl.paymentVault())
        );
    }

    function testScript() public virtual override {
        // 1 - Deploy implementations
        runAsEOA();

        // 2 - Execute upgrades
        execute();

        // 3 - Validate upgrades
        _validateUpgrades();
    }

    /// @dev Validate that all upgrades were successful
    function _validateUpgrades() internal view {
        // Validate ServiceManager upgrade
        address serviceManagerImpl = Env._getProxyImpl(address(Env.proxy.serviceManager()));
        assertTrue(
            serviceManagerImpl == address(Env.impl.serviceManager()),
            "ServiceManager proxy should point to new implementation"
        );

        // Validate RegistryCoordinator upgrade
        address registryCoordinatorImpl = Env._getProxyImpl(address(Env.proxy.registryCoordinator()));
        assertTrue(
            registryCoordinatorImpl == address(Env.impl.registryCoordinator()),
            "RegistryCoordinator proxy should point to new implementation"
        );

        // Validate ThresholdRegistry upgrade
        address thresholdRegistryImpl = Env._getProxyImpl(address(Env.proxy.thresholdRegistry()));
        assertTrue(
            thresholdRegistryImpl == address(Env.impl.thresholdRegistry()),
            "ThresholdRegistry proxy should point to new implementation"
        );

        // Validate RelayRegistry upgrade
        address relayRegistryImpl = Env._getProxyImpl(address(Env.proxy.relayRegistry()));
        assertTrue(
            relayRegistryImpl == address(Env.impl.relayRegistry()),
            "RelayRegistry proxy should point to new implementation"
        );

        // Validate ProxyAdmin ownership
        address proxyAdmin = Env._getProxyAdmin(address(Env.proxy.serviceManager()));
        assertTrue(proxyAdmin == Env.proxyAdmin(), "ServiceManager proxy admin should be correct");

        proxyAdmin = Env._getProxyAdmin(address(Env.proxy.registryCoordinator()));
        assertTrue(proxyAdmin == Env.proxyAdmin(), "RegistryCoordinator proxy admin should be correct");
    }
}
