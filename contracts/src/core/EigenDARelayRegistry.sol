// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {EigenDARelayRegistryStorage} from "./EigenDARelayRegistryStorage.sol";
import {IEigenDARelayRegistry} from "../interfaces/IEigenDARelayRegistry.sol";
import "../interfaces/IEigenDAStructs.sol";

/**
 * @title Registry for EigenDA relay keys
 * @author Layr Labs, Inc.
 * @notice Permissionless registration with a deposit on registration and each update.
 *         Deregistration refunds the accumulated deposit minus a flat percentage tax.
 *         All users, including the owner, must pay deposits and taxes.
 */
contract EigenDARelayRegistry is OwnableUpgradeable, EigenDARelayRegistryStorage, IEigenDARelayRegistry {

    constructor() {
        _disableInitializers();
    }

    function initialize(
        address _initialOwner
    ) external initializer {
        _transferOwnership(_initialOwner);
    }

    /**
     * @notice Register a new relay.
     * @dev The sender must supply exactly REGISTRATION_DEPOSIT.
     *      All users, including the owner, must pay the full deposit.
     */
    function addRelayInfo(RelayInfo memory relayInfo) external payable returns (uint32) {
        // Require deposit for all users
        require(msg.value == REGISTRATION_DEPOSIT, "Deposit of 10 ether required for registration");
        
        // Store relay info
        relayKeyToInfo[nextRelayKey] = relayInfo;
        
        // Track registrant and deposit
        relayRegistrants[nextRelayKey] = msg.sender;
        relayDeposits[nextRelayKey] = msg.value;
        
        // Track if the relay was created by the owner
        if (msg.sender == owner()) {
            isOwnerCreatedRelay[nextRelayKey] = true;
        }
        
        emit RelayAdded(relayInfo.relayAddress, nextRelayKey, relayInfo.relayURL);
        return nextRelayKey++;
    }

    /**
     * @notice Deregister a relay entry.
     * @dev The caller must be the original registrant.
     *      A TAX_PERCENTAGE of the accumulated deposit is retained; the rest is refunded.
     *      All users, including the owner, pay the tax.
     */
    function deregisterRelay(uint32 _relayKey) external {
        require(relayRegistrants[_relayKey] == msg.sender, "Caller is not the registrant");
        
        uint256 totalDeposit = relayDeposits[_relayKey];
        require(totalDeposit > 0, "No deposit found");
        
        uint256 tax = (totalDeposit * TAX_PERCENTAGE) / 100;
        uint256 refund = totalDeposit - tax;

        // Clear registration data
        delete relayRegistrants[_relayKey];
        delete relayDeposits[_relayKey];
        delete relayKeyToInfo[_relayKey];
        delete isOwnerCreatedRelay[_relayKey];

        emit RelayRemoved(_relayKey, msg.sender);
        payable(msg.sender).transfer(refund);
    }

    function relayKeyToAddress(uint32 key) external view returns (address) {
        return relayKeyToInfo[key].relayAddress;
    }

    function relayKeyToUrl(uint32 key) external view returns (string memory) {
        return relayKeyToInfo[key].relayURL;
    }
}
