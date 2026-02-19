// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "lib/openzeppelin-contracts-upgradeable/contracts/access/OwnableUpgradeable.sol";
import {EigenDARelayRegistryStorage} from "./EigenDARelayRegistryStorage.sol";
import {IEigenDARelayRegistry} from "src/core/interfaces/IEigenDARelayRegistry.sol";
import {EigenDATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";

/// @title Registry for EigenDA relay keys
/// @author Layr Labs, Inc.
contract EigenDARelayRegistry is OwnableUpgradeable, EigenDARelayRegistryStorage, IEigenDARelayRegistry {
    constructor() {
        _disableInitializers();
    }

    function initialize(address _initialOwner) external initializer {
        _transferOwnership(_initialOwner);
    }

    /// @notice Reinitializer to set a new owner during proxy upgrade, replacing the lost EOA owner.
    /// @param _newOwner The address of the new owner.
    function initializeV2(address _newOwner) external reinitializer(2) {
        _transferOwnership(_newOwner);
    }

    function addRelayInfo(EigenDATypesV2.RelayInfo memory relayInfo) external onlyOwner returns (uint32) {
        relayKeyToInfo[nextRelayKey] = relayInfo;
        emit RelayAdded(relayInfo.relayAddress, nextRelayKey, relayInfo.relayURL);
        return nextRelayKey++;
    }

    function relayKeyToAddress(uint32 key) external view returns (address) {
        return relayKeyToInfo[key].relayAddress;
    }

    function relayKeyToUrl(uint32 key) external view returns (string memory) {
        return relayKeyToInfo[key].relayURL;
    }
}
