// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";
import {EigenDARelayRegistryStorage} from "./EigenDARelayRegistryStorage.sol";
import {IEigenDARelayRegistry} from "../interfaces/IEigenDARelayRegistry.sol";

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

    function addRelayURL(address relay, string memory relayURL) external onlyOwner returns (uint32) {
        relayKeyToURL[nextRelayKey] = relayURL;
        relayAddressToKey[relay] = nextRelayKey;
        relayKeyToAddress[nextRelayKey] = relay;

        emit RelayAdded(relay, nextRelayKey, relayURL);
        return nextRelayKey++;
    }

    function getRelayURL(uint32 key) external view returns (string memory) {
        return relayKeyToURL[key];
    }

    function getRelayKey(address relay) external view returns (uint32) {
        return relayAddressToKey[relay];
    }

    function getRelayAddress(uint32 key) external view returns (address) {
        return relayKeyToAddress[key];
    }
}