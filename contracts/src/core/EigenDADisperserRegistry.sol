// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {EigenDADisperserRegistryStorage} from "./EigenDADisperserRegistryStorage.sol";
import {IEigenDADisperserRegistry} from "../interfaces/IEigenDADisperserRegistry.sol";
import "../interfaces/IEigenDAStructs.sol";

/**
 * @title EigenDADisperserRegistry
 * @notice A registry for EigenDA disperser info
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
     * @notice Sets the disperser info for a given disperser key
     * @param _disperserKey The key of the disperser to set the info for
     * @param _disperserInfo The info to set for the disperser
     */
    function setDisperserInfo(uint32 _disperserKey, DisperserInfo memory _disperserInfo) external onlyOwner {
        disperserKeyToInfo[_disperserKey] = _disperserInfo;
        emit DisperserAdded(_disperserKey, _disperserInfo.disperserAddress);
    }

    /**
     * @notice Returns the disperser address for a given disperser key
     * @param _key The key of the disperser to get the address for
     */
    function disperserKeyToAddress(uint32 _key) external view returns (address) {
        return disperserKeyToInfo[_key].disperserAddress;
    }
}
