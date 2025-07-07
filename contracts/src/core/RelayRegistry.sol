// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IRelayRegistry} from "src/core/interfaces/IRelayRegistry.sol";

import {RelayRegistryLib} from "src/core/libraries/v3/relay-registry/RelayRegistryLib.sol";

contract RelayRegistry is IRelayRegistry {
    using RelayRegistryLib for uint32;

    function addRelay(address relay, string memory url, uint32[] memory dispersers)
        external
        override
        returns (uint32)
    {
        uint32 relayId = RelayRegistryLib.addRelay(relay, url, dispersers);
        emit RelayAdded(relayId, relay, url, dispersers);
        return relayId;
    }

    function getRelayAddress(uint32 relayId) external view override returns (address) {
        return relayId.getRelayAddress();
    }

    function getRelayUrl(uint32 relayId) external view override returns (string memory) {
        return relayId.getRelayUrl();
    }

    function getRelayDispersers(uint32 relayId) external view override returns (uint32[] memory) {
        return relayId.getRelayDispersers();
    }
}
