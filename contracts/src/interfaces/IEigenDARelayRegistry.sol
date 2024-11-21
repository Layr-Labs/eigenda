// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

interface IEigenDARelayRegistry {

    event RelayAdded(address indexed relay, uint32 indexed id, string relayURL);

    function setRelayURL(address relay, uint32 id, string memory relayURL) external;

    function getRelayURL(uint32 id) external view returns (string memory);

    function getRelayId(address relay) external view returns (uint32);

    function getRelayAddress(uint32 id) external view returns (address);
}
