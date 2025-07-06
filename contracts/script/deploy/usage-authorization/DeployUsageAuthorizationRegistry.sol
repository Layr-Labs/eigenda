// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {Script} from "forge-std/Script.sol";

import {TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {UsageAuthorizationRegistry, UsageAuthorizationTypes} from "src/core/UsageAuthorizationRegistry.sol";
import {ProxyAdmin} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

import {console2} from "forge-std/console2.sol";
import {stdToml} from "forge-std/StdToml.sol";

contract DeployUsageAuthorizationRegistry is Script {
    using stdToml for string;

    UsageAuthorizationRegistry usageAuthorizationRegistry;
    UsageAuthorizationRegistry usageAuthorizationRegistryImpl;

    string cfg;

    function _initConfig() internal virtual {
        cfg = vm.readFile(vm.envString("DEPLOY_CONFIG_PATH"));
        console2.logString(cfg);
    }

    function run() external {
        _initConfig();
        vm.startBroadcast();

        usageAuthorizationRegistryImpl = new UsageAuthorizationRegistry(uint64(cfg.readUint(".schedulePeriod")));

        usageAuthorizationRegistry = UsageAuthorizationRegistry(
            address(
                new TransparentUpgradeableProxy(
                    address(usageAuthorizationRegistryImpl),
                    address(cfg.readAddress(".proxyAdmin")),
                    abi.encodeWithSelector(UsageAuthorizationRegistry.initialize.selector, msg.sender)
                )
            )
        );

        vm.stopBroadcast();

        console2.log("UsageAuthorizationRegistry deployed at:", address(usageAuthorizationRegistry));
        console2.log("UsageAuthorizationRegistry implementation at:", address(usageAuthorizationRegistryImpl));
    }
}
