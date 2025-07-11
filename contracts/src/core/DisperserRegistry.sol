// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {
    DisperserRegistryLib,
    DisperserRegistryStorage,
    DisperserRegistryTypes
} from "src/core/libraries/v3/disperser/DisperserRegistryLib.sol";
import {IDisperserRegistry} from "src/core/interfaces/IDisperserRegistry.sol";

import {AccessControlLib} from "src/core/libraries/AccessControlLib.sol";
import {InitializableLib} from "src/core/libraries/InitializableLib.sol";
import {Constants} from "src/core/libraries/Constants.sol";

/**
 * @title Registry for EigenDA disperser info
 * @author Layr Labs, Inc.
 */
contract DisperserRegistry is IDisperserRegistry {
    modifier initializer() {
        InitializableLib.setInitializedVersion(1);
        _;
    }

    modifier onlyOwner() {
        AccessControlLib.checkRole(Constants.OWNER_ROLE, msg.sender);
        _;
    }

    modifier onlyDisperserOwner(uint32 disperserKey) {
        if (msg.sender != DisperserRegistryLib.getDisperserOwner(disperserKey)) {
            revert IDisperserRegistry.NotDisperserOwner(
                disperserKey, DisperserRegistryLib.getDisperserOwner(disperserKey)
            );
        }
        _;
    }

    function initialize(
        address initialOwner,
        DisperserRegistryTypes.LockedDisperserDeposit memory depositParams,
        uint256 updateFee
    ) external initializer {
        AccessControlLib.grantRole(Constants.OWNER_ROLE, initialOwner);
        DisperserRegistryLib.setDepositParams(depositParams);
        DisperserRegistryLib.setUpdateFee(updateFee);
    }

    /// @inheritdoc IDisperserRegistry
    function registerDisperser(address disperserAddress, string memory disperserURL)
        external
        returns (uint32 disperserKey)
    {
        return DisperserRegistryLib.registerDisperser(disperserAddress, disperserURL);
    }

    /// DISPERSER

    /// @inheritdoc IDisperserRegistry
    function transferDisperserOwnership(uint32 disperserKey, address disperserAddress)
        external
        onlyDisperserOwner(disperserKey)
    {
        DisperserRegistryLib.transferDisperserOwnership(disperserKey, disperserAddress);
    }

    /// @inheritdoc IDisperserRegistry
    function updateDisperserInfo(uint32 disperserKey, address disperser, string memory disperserURL)
        external
        onlyDisperserOwner(disperserKey)
    {
        DisperserRegistryLib.updateDisperserInfo(disperserKey, disperser, disperserURL);
    }

    /// @inheritdoc IDisperserRegistry
    function deregisterDisperser(uint32 disperserKey) external onlyDisperserOwner(disperserKey) {
        DisperserRegistryLib.deregisterDisperser(disperserKey);
    }

    /// @inheritdoc IDisperserRegistry
    function withdrawDisperserDeposit(uint32 disperserKey) external onlyDisperserOwner(disperserKey) {
        DisperserRegistryLib.withdrawDisperserDeposit(disperserKey);
    }

    /// OWNER

    function setDepositParams(DisperserRegistryTypes.LockedDisperserDeposit memory depositParams) external onlyOwner {
        DisperserRegistryLib.setDepositParams(depositParams);
    }

    function setUpdateFee(uint256 updateFee) external onlyOwner {
        DisperserRegistryLib.setUpdateFee(updateFee);
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

    function getDisperserOwner(uint32 disperserKey) external view returns (address) {
        return DisperserRegistryLib.getDisperserOwner(disperserKey);
    }

    function getDisperserURL(uint32 disperserKey) external view returns (string memory) {
        return DisperserRegistryLib.getDisperserURL(disperserKey);
    }

    function getDisperserDepositUnlockTime(uint32 disperserKey) external view returns (uint64) {
        return DisperserRegistryLib.getDisperserUnlockTime(disperserKey);
    }

    function getDisperserDepositParams(uint32 disperserKey)
        external
        view
        returns (DisperserRegistryTypes.LockedDisperserDeposit memory)
    {
        return DisperserRegistryLib.getDisperserDepositParams(disperserKey);
    }

    function getDepositParams() external view returns (DisperserRegistryTypes.LockedDisperserDeposit memory) {
        return DisperserRegistryLib.getDepositParams();
    }

    function getNextDisperserKey() external view returns (uint32) {
        return DisperserRegistryLib.getNextDisperserKey();
    }

    function getExcessBalance(address token) external view returns (uint256) {
        return DisperserRegistryLib.getExcessBalance(token);
    }

    function getUpdateFee() external view returns (uint256) {
        return DisperserRegistryLib.getUpdateFee();
    }
}
