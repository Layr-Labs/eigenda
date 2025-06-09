// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {
    DisperserRegistryLib,
    DisperserRegistryStorage,
    DisperserRegistryTypes
} from "src/core/libraries/v3/disperser/DisperserRegistryLib.sol";
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

    modifier onlyDisperserOwner(uint32 disperserKey) {
        require(msg.sender == DisperserRegistryLib.getDisperserOwner(disperserKey), "Caller is not the disperser");
        _;
    }

    function initialize(address initialOwner, DisperserRegistryTypes.LockedDisperserDeposit memory depositParams)
        external
        initializer
    {
        AccessControlLib.grantRole(Constants.OWNER_ROLE, initialOwner);
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
        onlyDisperserOwner(disperserKey)
    {
        DisperserRegistryLib.transferDisperserOwnership(disperserKey, disperserAddress);
    }

    function updateDisperserInfo(uint32 disperserKey, address disperser, string memory disperserURL)
        external
        onlyDisperserOwner(disperserKey)
    {
        DisperserRegistryLib.updateDisperserInfo(disperserKey, disperser, disperserURL);
    }

    function deregisterDisperser(uint32 disperserKey) external onlyDisperserOwner(disperserKey) {
        DisperserRegistryLib.deregisterDisperser(disperserKey);
    }

    function withdraw(uint32 disperserKey) external onlyDisperserOwner(disperserKey) {
        DisperserRegistryLib.withdraw(disperserKey);
    }

    /// OWNER

    function setDepositParams(DisperserRegistryTypes.LockedDisperserDeposit memory depositParams) external {
        DisperserRegistryLib.setDepositParams(depositParams);
    }

    function transferOwnership(address newOwner) external onlyOwner {
        AccessControlLib.transferRole(Constants.OWNER_ROLE, msg.sender, newOwner);
    }

    /// GETTERS

    function owner() external view returns (address) {
        return AccessControlLib.getRoleMember(Constants.OWNER_ROLE, 0);
    }

    function getDisperserAddress(uint32 disperserKey) external view returns (address) {
        return DisperserRegistryLib.getDisperserAddress(disperserKey);
    }

    function getDepositParams() external view returns (DisperserRegistryTypes.LockedDisperserDeposit memory) {
        return DisperserRegistryLib.getDepositParams();
    }
}
