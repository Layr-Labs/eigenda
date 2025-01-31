// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "../interfaces/IEigenDAStructs.sol";

/**
 * @title Storage variables for the `EigenDADisperserRegistry` contract.
 * @author Layr Labs, Inc.
 * @notice This storage contract is separate from the logic to simplify the upgrade process.
 */
abstract contract EigenDADisperserRegistryStorage {

    mapping(uint32 => DisperserInfo) public disperserKeyToInfo;
    
    // Deposit required for initial registration and for each update.
    uint256 public constant REGISTRATION_DEPOSIT = 10 ether;
    // Tax percentage applied on deregistration.
    uint256 public constant TAX_PERCENTAGE = 10; // 10%

    // Track the registrant and the total deposit for each disperser key.
    mapping(uint32 => address) public disperserRegistrants;
    mapping(uint32 => uint256) public disperserDeposits;
    
    // Add mapping to track owner-created dispersers
    mapping(uint32 => bool) public isOwnerCreatedDisperser;

    // storage gap for upgradeability
    // slither-disable-next-line shadowing-state
    uint256[44] private __GAP; // reduced from 45 to 44 to account for the new mapping
}