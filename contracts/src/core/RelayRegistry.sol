// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IRelayRegistry} from "src/core/interfaces/IRelayRegistry.sol";

import {RelayRegistryLib} from "src/core/libraries/v3/relay-registry/RelayRegistryLib.sol";

contract RelayRegistry is IRelayRegistry {
    using RelayRegistryLib for uint32;

    function addRelayInfo(address relay, string memory url, uint32[] memory dispersers) external returns (uint32) {
        uint32 relayId = RelayRegistryLib.addRelay(relay, url, dispersers);
        emit RelayAdded(relayId, relay, url, dispersers);
        return relayId;
    }

    function addRelayInfo(address relay, string memory url) external returns (uint32) {
        uint32 relayId = RelayRegistryLib.addRelay(relay, url, new uint32[](0));
        emit RelayAdded(relayId, relay, url, new uint32[](0));
        return relayId;
    }

    function relayKeyToAddress(uint32 key) external view returns (address) {
        return key.getRelayAddress();
    }

    function relayKeyToUrl(uint32 key) external view returns (string memory) {
        return key.getRelayUrl();
    }

    function relayKeyToDispersers(uint32 key) external view override returns (uint32[] memory) {
        return key.getRelayDispersers();
    }
}
