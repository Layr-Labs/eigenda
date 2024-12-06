// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

/**
 * @title Storage variables for the `EigenDARelayRegistry` contract.
 * @author Layr Labs, Inc.
 * @notice This storage contract is separate from the logic to simplify the upgrade process.
 */
abstract contract EigenDARelayRegistryStorage {

    mapping(uint32 => string) public relayKeyToURL;

    mapping(address => uint32) public relayAddressToKey;

    mapping(uint32 => address) public relayKeyToAddress;

    uint32 public nextRelayKey;

    // storage gap for upgradeability
    // slither-disable-next-line shadowing-state
    uint256[46] private __GAP;
}