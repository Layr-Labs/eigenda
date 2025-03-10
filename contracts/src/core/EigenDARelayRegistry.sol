// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "lib/openzeppelin-contracts-upgradeable/contracts/access/OwnableUpgradeable.sol";
import {EigenDARelayRegistryStorage} from "./EigenDARelayRegistryStorage.sol";
import {IEigenDARelayRegistry} from "../interfaces/IEigenDARelayRegistry.sol";
import "../interfaces/IEigenDAStructs.sol";

/**
 * @title EigenDARelayRegistry
 * @notice A registry for EigenDA relay info
 * @dev This contract is append only and does not support updating or removing relay info
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
     * @notice Appends a relay info to the registry and returns the relay key
     * @param relayInfo The relay info to add
     */
    function addRelayInfo(RelayInfo memory relayInfo) external onlyOwner returns (uint32) {
        relayKeyToInfo[nextRelayKey] = relayInfo;
        emit RelayAdded(relayInfo.relayAddress, nextRelayKey, relayInfo.relayURL);
        return nextRelayKey++;
    }

    /**
     * @notice Returns the relay address for a given relay key
     * @param key The key of the relay to get the address for
     */
    function relayKeyToAddress(uint32 key) external view returns (address) {
        return relayKeyToInfo[key].relayAddress;
    }

    /**
     * @notice Returns the relay URL for a given relay key
     * @param key The key of the relay to get the URL for
     */
    function relayKeyToUrl(uint32 key) external view returns (string memory) {
        return relayKeyToInfo[key].relayURL;
    }
}
