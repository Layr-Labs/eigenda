// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {Script} from "forge-std/Script.sol";

import {TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {UsageAuthorizationRegistry} from "src/core/UsageAuthorizationRegistry.sol";
import {ProxyAdmin} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";

import {console2} from "forge-std/console2.sol";

contract DeployUsageAuthorizationRegistry is Script {
    function run() external {
        vm.startBroadcast();

        uint64 schedulePeriod = 1 hours; // Example schedule period

        ProxyAdmin proxyAdmin = new ProxyAdmin();
        UsageAuthorizationRegistry usageAuthorizationRegistry = new UsageAuthorizationRegistry(schedulePeriod);

        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(usageAuthorizationRegistry),
            address(proxyAdmin),
            abi.encodeWithSelector(UsageAuthorizationRegistry.initialize.selector, msg.sender)
        );

        vm.stopBroadcast();

        console2.log("ProxyAdmin deployed at:", address(proxyAdmin));
        console2.log("UsageAuthorizationRegistry deployed at:", address(proxy));
        console2.log("UsageAuthorizationRegistry implementation at:", address(usageAuthorizationRegistry));
    }
}
