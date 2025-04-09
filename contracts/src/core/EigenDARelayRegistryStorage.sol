// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "../interfaces/IEigenDAStructs.sol";

/**
 * @title Storage variables for the `EigenDARelayRegistry` contract.
 * @author Layr Labs, Inc.
 * @notice This storage contract is separate from the logic to simplify the upgrade process.
 */
abstract contract EigenDARelayRegistryStorage {

    mapping(uint32 => RelayInfo) public relayKeyToInfo;

    uint32 public nextRelayKey;
    
    // Deposit required for initial registration and for each update.
    uint256 public constant REGISTRATION_DEPOSIT = 10 ether;
    // Tax percentage applied on deregistration.
    uint256 public constant TAX_PERCENTAGE = 10; // 10%

    // Track the registrant and the total deposit for each relay key.
    mapping(uint32 => address) public relayRegistrants;
    mapping(uint32 => uint256) public relayDeposits;
    
    // Add mapping to track owner-created relays
    mapping(uint32 => bool) public isOwnerCreatedRelay;

    // storage gap for upgradeability
    // slither-disable-next-line shadowing-state
    uint256[44] private __GAP; // reduced from 49 to 44 to account for the new storage slots
}