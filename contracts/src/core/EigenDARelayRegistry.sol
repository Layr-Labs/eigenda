// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "lib/openzeppelin-contracts-upgradeable/contracts/access/OwnableUpgradeable.sol";
import {EigenDARelayRegistryStorage} from "./EigenDARelayRegistryStorage.sol";
import {IEigenDARelayRegistry} from "../interfaces/IEigenDARelayRegistry.sol";
import "../interfaces/IEigenDAStructs.sol";

/**
 * @title Registry for EigenDA relay keys
 * @author Layr Labs, Inc.
 */
contract EigenDARelayRegistry is OwnableUpgradeable, EigenDARelayRegistryStorage, IEigenDARelayRegistry {


    // Add mapping to track owner-created relays
    mapping(uint32 => bool) public isOwnerCreatedRelay;

    constructor() {
        _disableInitializers();
    }

    function initialize(
        address _initialOwner
    ) external initializer {
        _transferOwnership(_initialOwner);
    }

    function addRelayInfo(RelayInfo memory relayInfo) external returns (uint32) {
        // Estimate the cost of removing the relay info 
        // require msg.sender to add value at least as much as the cost of removing the relay info 
        // put the value in contract wallet

        relayKeyToInfo[nextRelayKey] = relayInfo;
        
        // Track if the relay was created by the owner
        if (msg.sender == owner()) {
            isOwnerCreatedRelay[nextRelayKey] = true;
        }
        
        emit RelayAdded(relayInfo.relayAddress, nextRelayKey, relayInfo.relayURL);
        return nextRelayKey++;
    }

    function deregisterRelayInfo(uint32 _relayKey) external onlyOwner {
        delete relayKeyToInfo[_relayKey];
        delete isOwnerCreatedRelay[_relayKey];
        
        emit RelayRemoved(_relayKey);
        // taking funds from the wallet, refund gas cost to msg.sender
    }

    function relayKeyToAddress(uint32 key) external view returns (address) {
        return relayKeyToInfo[key].relayAddress;
    }

    function relayKeyToUrl(uint32 key) external view returns (string memory) {
        return relayKeyToInfo[key].relayURL;
    }
}
