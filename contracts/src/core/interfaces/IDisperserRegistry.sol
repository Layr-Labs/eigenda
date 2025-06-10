// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {DisperserRegistryTypes} from "src/core/libraries/v3/disperser/DisperserRegistryTypes.sol";

interface IDisperserRegistry {
    error NotDisperserOwner(uint32 disperserKey, address owner);
    error InvalidDisperserAddress(address disperserAddress);
    error DisperserNotRegistered(uint32 disperserKey);
    error InvalidNewOwner(address newOwner);
    error DisperserNotDeregistered(uint32 disperserKey);
    error RefundLocked(uint32 disperserKey, uint64 unlockTimestamp);
    error ZeroRefund(uint32 disperserKey);
    error DepositMustBeAtLeastRefund(uint256 deposit, uint256 refund);
    error InvalidTokenAddress(address token);

    /// @notice Registers a disperser with the given address and URL.
    /// @param disperserAddress The address of the disperser.
    /// @param disperserURL The URL of the disperser.
    /// @return disperserKey The key assigned to the registered disperser.
    /// @dev The caller must have sufficient balance to cover the deposit amount, detailed in getDepositParams.
    function registerDisperser(address disperserAddress, string memory disperserURL)
        external
        returns (uint32 disperserKey);

    /// @notice Transfers ownership of a disperser to a new owner.
    /// @param disperserKey The key of the disperser to transfer ownership of.
    /// @param newOwner The address of the new owner.
    /// @dev The caller must be the current owner of the disperser.
    function transferDisperserOwnership(uint32 disperserKey, address newOwner) external;

    /// @notice Updates the information of a registered disperser.
    /// @param disperserKey The key of the disperser to update.
    /// @param disperser The new address of the disperser.
    /// @param disperserURL The new URL of the disperser.
    /// @dev The caller must be the current owner of the disperser.
    function updateDisperserInfo(uint32 disperserKey, address disperser, string memory disperserURL) external;

    /// @notice Deregisters a disperser, allowing the owner to withdraw their deposit after the lock period.
    /// @param disperserKey The key of the disperser to deregister.
    /// @dev The caller must be the current owner of the disperser. The deposit will be locked for a period defined in getDepositParams.
    function deregisterDisperser(uint32 disperserKey) external;

    /// @notice Withdraws the deposit of a deregistered disperser after the lock period.
    /// @param disperserKey The key of the disperser to withdraw from.
    /// @dev The caller must be the owner of the disperser and the lock period must have expired.
    function withdrawDisperserDeposit(uint32 disperserKey) external;

    /// @notice Returns the address of a registered disperser.
    /// @param disperserKey The key of the disperser to query.
    /// @return The address of the registered disperser.
    /// @dev Returns the zero address if the disperser is not registered.
    function getDisperserAddress(uint32 disperserKey) external view returns (address);

    /// @notice Returns the information of a registered disperser.
    /// @param disperserKey The key of the disperser to query.
    /// @return The information of the registered disperser, including address, owner, URL, and deposit parameters.
    function getDisperserOwner(uint32 disperserKey) external view returns (address);

    /// @notice Returns the URL of a registered disperser.
    /// @param disperserKey The key of the disperser to query.
    /// @return The URL of the registered disperser.
    function getDisperserURL(uint32 disperserKey) external view returns (string memory);

    /// @notice Returns the unlock timestamp of a registered disperser's deposit.
    /// @param disperserKey The key of the disperser to query.
    /// @return The unlock timestamp of the registered disperser's deposit.
    function getDisperserDepositUnlockTime(uint32 disperserKey) external view returns (uint64);

    function getDisperserDepositParams(uint32 disperserKey)
        external
        view
        returns (DisperserRegistryTypes.LockedDisperserDeposit memory);

    /// @notice Returns the deposit parameters for registering a disperser.
    /// @return The deposit parameters for registering a disperser.
    function getDepositParams() external view returns (DisperserRegistryTypes.LockedDisperserDeposit memory);
}
