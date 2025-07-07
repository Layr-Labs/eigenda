// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {RelayRegistryStorage} from "src/core/libraries/v3/relay-registry/RelayRegistryStorage.sol";

library RelayRegistryLib {
    struct RelayInfo {
        address relay;
        string url;
        uint32[] dispersers;
    }

    function getRelayInfo(uint32 relayId) internal view returns (RelayInfo storage) {
        return RelayRegistryStorage.layout().relay[relayId];
    }

    function getRelayAddress(uint32 relayId) internal view returns (address) {
        return getRelayInfo(relayId).relay;
    }

    function getRelayUrl(uint32 relayId) internal view returns (string memory) {
        return getRelayInfo(relayId).url;
    }

    function getRelayDispersers(uint32 relayId) internal view returns (uint32[] memory) {
        return getRelayInfo(relayId).dispersers;
    }

    function addRelay(address relay, string memory url, uint32[] memory dispersers) internal returns (uint32) {
        RelayRegistryStorage.Layout storage s = RelayRegistryStorage.layout();
        uint32 relayId = s.nextRelayId++;
        RelayInfo storage info = s.relay[relayId];
        info.relay = relay;
        info.url = url;
        info.dispersers = dispersers;
        return relayId;
    }
}
