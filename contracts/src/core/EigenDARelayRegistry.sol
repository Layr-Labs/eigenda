// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";
import {EigenDARelayRegistryStorage} from "./EigenDARelayRegistryStorage.sol";
import {IEigenDARelayRegistry} from "../interfaces/IEigenDARelayRegistry.sol";
import "../interfaces/IEigenDAStructs.sol";

/**
 * @title Registry for EigenDA relay keys
 * @author Layr Labs, Inc.
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

    function addRelayInfo(RelayInfo memory relayInfo) external onlyOwner returns (uint32) {
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
