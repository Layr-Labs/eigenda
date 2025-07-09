// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

interface IRelayRegistry {
    event RelayAdded(uint32 indexed relayId, address indexed relay, string url, uint32[] dispersers);

    function addRelayInfo(address relay, string memory url, uint32[] memory dispersers) external returns (uint32);

    function addRelayInfo(address relay, string memory url) external returns (uint32);

    function relayKeyToAddress(uint32 relayId) external view returns (address);

    function relayKeyToUrl(uint32 relayId) external view returns (string memory);

    function relayKeyToDispersers(uint32 relayId) external view returns (uint32[] memory);
}
