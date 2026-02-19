// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import {ProxyAdmin, TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {EigenDADisperserRegistry} from "src/core/EigenDADisperserRegistry.sol";
import {EigenDARelayRegistry} from "src/core/EigenDARelayRegistry.sol";
import {EigenDATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";

import {Script} from "forge-std/Script.sol";
import {console2} from "forge-std/console2.sol";

/// @notice Upgrades the EigenDADisperserRegistry and EigenDARelayRegistry proxies to new implementations
///         and calls initializeV2 on each to set new owners (replacing lost EOA owners).
///         Also updates disperser key 1 to a new address.
/// @dev The broadcaster (msg.sender) must be:
///      1. The ProxyAdmin owner (to call upgradeAndCall)
///      2. The DISPERSER_REGISTRY_NEW_OWNER (to call setDisperserInfo after ownership transfer)
///
///   PROXY_ADMIN=<proxyAdmin> \
///   DISPERSER_REGISTRY_PROXY=<proxy> \
///   DISPERSER_REGISTRY_NEW_OWNER=<newOwner> \
///   RELAY_REGISTRY_PROXY=<proxy> \
///   RELAY_REGISTRY_NEW_OWNER=<newOwner> \
///   forge script script/UpgradeRegistries.s.sol --rpc-url <rpc> --broadcast
contract UpgradeRegistries is Script {
    function run() external {
        address proxyAdmin = vm.envAddress("PROXY_ADMIN");

        address disperserRegistryProxy = vm.envAddress("DISPERSER_REGISTRY_PROXY");
        address disperserNewOwner = vm.envAddress("DISPERSER_REGISTRY_NEW_OWNER");

        address relayRegistryProxy = vm.envAddress("RELAY_REGISTRY_PROXY");
        address relayNewOwner = vm.envAddress("RELAY_REGISTRY_NEW_OWNER");

        require(disperserNewOwner != address(0), "DISPERSER_REGISTRY_NEW_OWNER must not be zero address");
        require(relayNewOwner != address(0), "RELAY_REGISTRY_NEW_OWNER must not be zero address");

        vm.startBroadcast();

        // --- Disperser Registry ---
        EigenDADisperserRegistry disperserImpl = new EigenDADisperserRegistry();
        console2.log("DisperserRegistry new implementation:", address(disperserImpl));

        ProxyAdmin(proxyAdmin).upgradeAndCall(
            TransparentUpgradeableProxy(payable(disperserRegistryProxy)),
            address(disperserImpl),
            abi.encodeCall(EigenDADisperserRegistry.initializeV2, (disperserNewOwner))
        );
        console2.log("DisperserRegistry upgraded. New owner:", disperserNewOwner);

        // Update disperser key 1 to new address (requires msg.sender == disperserNewOwner)
        EigenDADisperserRegistry(disperserRegistryProxy).setDisperserInfo(
            1, EigenDATypesV2.DisperserInfo(0xa44c6C843CFD5b7720F5D8e241D207D50E125993)
        );
        console2.log("DisperserRegistry: updated disperser key 1");

        // --- Relay Registry ---
        EigenDARelayRegistry relayImpl = new EigenDARelayRegistry();
        console2.log("RelayRegistry new implementation:", address(relayImpl));

        ProxyAdmin(proxyAdmin).upgradeAndCall(
            TransparentUpgradeableProxy(payable(relayRegistryProxy)),
            address(relayImpl),
            abi.encodeCall(EigenDARelayRegistry.initializeV2, (relayNewOwner))
        );
        console2.log("RelayRegistry upgraded. New owner:", relayNewOwner);

        // Add relay info for key 1 (requires msg.sender == relayNewOwner)
        uint32 relayKey = EigenDARelayRegistry(relayRegistryProxy).addRelayInfo(
            EigenDATypesV2.RelayInfo(
                0xA3f41F215E06De8439e9F8b767976647dE8C44cc,
                "relay-1-testnet-ussj1-hoodi.eigenda.xyz"
            )
        );
        console2.log("RelayRegistry: added relay key", relayKey);

        vm.stopBroadcast();
    }
}
