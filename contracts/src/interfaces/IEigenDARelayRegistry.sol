// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "./IEigenDAStructs.sol";

interface IEigenDARelayRegistry {

    event RelayAdded(address indexed relay, uint32 indexed key, string relayURL);

    function addRelayInfo(RelayInfo memory relayInfo) external returns (uint32);

    function relayKeyToAddress(uint32 key) external view returns (address);

    function relayKeyToUrl(uint32 key) external view returns (string memory);
}
