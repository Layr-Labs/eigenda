// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "./IEigenDAStructs.sol";

interface IEigenDARelayRegistry {

    event RelayAdded(address indexed relayAddress, uint32 indexed relayKey, string relayURL);
    event RelayRemoved(uint32 indexed relayKey, address registrant);

    function addRelayInfo(RelayInfo memory relayInfo) external payable returns (uint32);

    function deregisterRelay(uint32 _relayKey) external;

    function relayKeyToAddress(uint32 key) external view returns (address);

    function relayKeyToUrl(uint32 key) external view returns (string memory);
}
