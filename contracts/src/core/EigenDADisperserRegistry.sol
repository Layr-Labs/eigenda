// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";
import {EigenDADisperserRegistryStorage} from "./EigenDADisperserRegistryStorage.sol";
import {IEigenDADisperserRegistry} from "../interfaces/IEigenDADisperserRegistry.sol";
import "../interfaces/IEigenDAStructs.sol";

/**
 * @title Registry for EigenDA disperser info
 * @author Layr Labs, Inc.
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

    function setDisperserInfo(uint32 _disperserKey, DisperserInfo memory _disperserInfo) external onlyOwner {
        disperserKeyToInfo[_disperserKey] = _disperserInfo;
        emit DisperserAdded(_disperserKey, _disperserInfo.disperserAddress);
    }

    function disperserKeyToAddress(uint32 _key) external view returns (address) {
        return disperserKeyToInfo[_key].disperserAddress;
    }
}
