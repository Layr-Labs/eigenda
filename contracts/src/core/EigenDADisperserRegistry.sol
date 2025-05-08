// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV3} from "src/core/libraries/v3/EigenDATypesV3.sol";
import {DisperserRegistryLib, DisperserRegistryStorage} from "src/core/libraries/v3/DisperserRegistryLib.sol";
import {IEigenDADisperserRegistry} from "src/core/interfaces/IEigenDADisperserRegistry.sol";
import {AccessControlLib} from "src/core/libraries/AccessControlLib.sol";

import {InitializableLib} from "src/core/libraries/InitializableLib.sol";
import {Constants} from "src/core/libraries/Constants.sol";

/**
 * @title Registry for EigenDA disperser info
 * @author Layr Labs, Inc.
 */
contract EigenDADisperserRegistry is IEigenDADisperserRegistry {
    modifier initializer() {
        InitializableLib.setInitializedVersion(1);
        _;
    }

    modifier onlyOwner() {
        AccessControlLib.checkRole(Constants.OWNER_ROLE, msg.sender);
        _;
    }

    modifier onlyDisperser(uint32 disperserKey) {
        require(msg.sender == DisperserRegistryLib.getDisperser(disperserKey), "Caller is not the disperser");
        _;
    }

    function initialize(address owner, EigenDATypesV3.LockedDisperserDeposit memory depositParams)
        external
        initializer
    {
        AccessControlLib.grantRole(Constants.OWNER_ROLE, owner);
        DisperserRegistryLib.setDepositParams(depositParams);
    }

    function registerDisperser(address disperserAddress, string memory disperserURL)
        external
        returns (uint32 disperserKey)
    {
        return DisperserRegistryLib.registerDisperser(disperserAddress, disperserURL);
    }

    /// DISPERSER

    function transferDisperserOwnership(uint32 disperserKey, address disperserAddress)
        external
        onlyDisperser(disperserKey)
    {
        DisperserRegistryLib.transferDisperserOwnership(disperserKey, disperserAddress);
    }

    function updateDisperserURL(uint32 disperserKey, string memory disperserURL) external onlyDisperser(disperserKey) {
        DisperserRegistryLib.updateDisperserURL(disperserKey, disperserURL);
    }

    function deregisterDisperser(uint32 disperserKey) external onlyDisperser(disperserKey) {
        DisperserRegistryLib.deregisterDisperser(disperserKey);
    }

    function withdraw(uint32 disperserKey) external onlyDisperser(disperserKey) {
        DisperserRegistryLib.withdraw(disperserKey);
    }

    /// OWNER

    function setDepositParams(EigenDATypesV3.LockedDisperserDeposit memory depositParams) external {
        DisperserRegistryLib.setDepositParams(depositParams);
    }

    function transferOwnership(address newOwner) external onlyOwner {
        AccessControlLib.transferRole(Constants.OWNER_ROLE, msg.sender, newOwner);
    }

    /// GETTERS

    function getDepositParams() external view returns (EigenDATypesV3.LockedDisperserDeposit memory) {
        return DisperserRegistryLib.getDepositParams();
    }

    function getDisperserInfo(uint32 disperserKey) external view returns (EigenDATypesV3.DisperserInfo memory) {
        return DisperserRegistryLib.getDisperserInfo(disperserKey);
    }

    function getLockedDeposit(uint32 disperserKey)
        external
        view
        returns (EigenDATypesV3.LockedDisperserDeposit memory, uint64 unlockTimestamp)
    {
        return DisperserRegistryLib.getLockedDeposit(disperserKey);
    }
}
