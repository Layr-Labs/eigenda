// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

interface IRelayRegistry {
    event RelayAdded(uint32 indexed relayId, address indexed relay, string url, uint32[] dispersers);

    function addRelay(address relay, string memory url, uint32[] memory dispersers) external returns (uint32);

    function getRelayAddress(uint32 relayId) external view returns (address);

    function getRelayUrl(uint32 relayId) external view returns (string memory);

    function getRelayDispersers(uint32 relayId) external view returns (uint32[] memory);
}
