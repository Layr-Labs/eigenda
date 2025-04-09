// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {EigenDADisperserRegistryStorage} from "./EigenDADisperserRegistryStorage.sol";
import {IEigenDADisperserRegistry} from "../interfaces/IEigenDADisperserRegistry.sol";
import "../interfaces/IEigenDAStructs.sol";

/**
 * @title Registry for EigenDA disperser info
 * @author Layr Labs, Inc.
 * @notice Permissionless registration with a deposit on registration and each update.
 *         Deregistration refunds the accumulated deposit minus a flat percentage tax.
 *         All users, including the owner, must pay deposits and taxes.
 */
contract EigenDADisperserRegistry is OwnableUpgradeable, EigenDADisperserRegistryStorage, IEigenDADisperserRegistry {

    constructor() {
        _disableInitializers();
    }

    function initialize(
        address _initialOwner
    ) external initializer {
        _transferOwnership(_initialOwner);
    }

    /**
     * @notice Register or update disperser info.
     * @dev If the key is new, the sender must supply exactly REGISTRATION_DEPOSIT.
     *      If already registered, only the original registrant may update and must supply an additional REGISTRATION_DEPOSIT.
     *      All users, including the owner, must pay the full deposit.
     */
    function setDisperserInfo(uint32 _disperserKey, DisperserInfo memory _disperserInfo) external payable {
        bool isOwner = msg.sender == owner();
        
        if (disperserRegistrants[_disperserKey] == address(0)) {
            // New registration - all users must pay deposit
            require(msg.value == REGISTRATION_DEPOSIT, "Deposit of 0.1 ether required for registration");
            disperserRegistrants[_disperserKey] = msg.sender;
            disperserDeposits[_disperserKey] = msg.value;
        } else {
            // Update – only the original registrant can update
            require(disperserRegistrants[_disperserKey] == msg.sender, "Caller is not the registrant");
            
            // All users must pay additional deposit for updates
            require(msg.value == REGISTRATION_DEPOSIT, "Additional deposit of 0.1 ether required for update");
            disperserDeposits[_disperserKey] += msg.value;
        }
        
        // Set disperser info
        disperserKeyToInfo[_disperserKey] = _disperserInfo;
        
        // If the sender is the owner, mark this disperser as owner-created
        if (isOwner) {
            isOwnerCreatedDisperser[_disperserKey] = true;
        }
        
        emit DisperserAdded(_disperserKey, _disperserInfo.disperserAddress);
    }

    /**
     * @notice Deregister a disperser entry.
     * @dev The caller must be the original registrant.
     *      A TAX_PERCENTAGE of the accumulated deposit is retained; the rest is refunded.
     *      All users, including the owner, pay the tax.
     */
    function deregisterDisperser(uint32 _disperserKey) external {
        require(disperserRegistrants[_disperserKey] == msg.sender, "Caller is not the registrant");
        
        uint256 totalDeposit = disperserDeposits[_disperserKey];
        require(totalDeposit > 0, "No deposit found");
        
        uint256 tax = (totalDeposit * TAX_PERCENTAGE) / 100;
        uint256 refund = totalDeposit - tax;

        // Clear registration data
        delete disperserRegistrants[_disperserKey];
        delete disperserDeposits[_disperserKey];
        delete disperserKeyToInfo[_disperserKey];
        delete isOwnerCreatedDisperser[_disperserKey];

        emit DisperserRemoved(_disperserKey, msg.sender);
        payable(msg.sender).transfer(refund);
    }

    function disperserKeyToAddress(uint32 _key) external view returns (address) {
        return disperserKeyToInfo[_key].disperserAddress;
    }
    
    function disperserKeysToAddresses(uint32[] memory _keys) external view returns (address[] memory) {
        address[] memory addresses = new address[](_keys.length);
        for (uint32 i = 0; i < _keys.length; i++) {
            addresses[i] = disperserKeyToInfo[_keys[i]].disperserAddress;
        }
        return addresses;
    }
}
