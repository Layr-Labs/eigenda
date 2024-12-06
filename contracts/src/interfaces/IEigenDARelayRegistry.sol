// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

interface IEigenDARelayRegistry {

    event RelayAdded(address indexed relay, uint32 indexed key, string relayURL);

    function addRelayURL(address relay, string memory relayURL) external returns (uint32);

    function getRelayURL(uint32 key) external view returns (string memory);

    function getRelayKey(address relay) external view returns (uint32);

    function getRelayAddress(uint32 key) external view returns (address);
}