// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDADirectory} from "src/core/EigenDADirectory.sol";

/// @notice Demo contract to demonstrate a view contract on the periphery that handles more complex retrieval logic
///         so that the core contract remains simple and gas efficient
contract MockConfigRetriever {
    EigenDADirectory public directory;

    constructor(EigenDADirectory _directory) {
        directory = _directory;
    }

    /// @notice Gets the number of checkpoints for a block number-based configuration
    /// @param key The hash of the configuration name
    /// @return The number of checkpoints stored
    function getNumCheckpointsBlockNumber(bytes32 key) external view returns (uint256) {
        return directory.getNumCheckpointsBlockNumber(key);
    }

    /// @notice Gets the number of checkpoints for a timestamp-based configuration
    /// @param key The hash of the configuration name
    /// @return The number of checkpoints stored
    function getNumCheckpointsTimeStamp(bytes32 key) external view returns (uint256) {
        return directory.getNumCheckpointsTimeStamp(key);
    }

    /// @notice Gets the configuration value at a specific index for a block number-based configuration
    /// @param key The hash of the configuration name
    /// @param index The index of the checkpoint to retrieve
    /// @return The bytes configuration value at the specified index
    function getConfigBytesBlockNumber(bytes32 key, uint256 index) external view returns (bytes memory) {
        return directory.getConfigBlockNumber(key, index);
    }

    /// @notice Gets the configuration value at a specific index for a timestamp-based configuration
    /// @param key The hash of the configuration name
    /// @param index The index of the checkpoint to retrieve
    /// @return The bytes configuration value at the specified index
    function getConfigTimeStamp(bytes32 key, uint256 index) external view returns (bytes memory) {
        return directory.getConfigTimeStamp(key, index);
    }

    /// @notice Gets the activation block number at a specific index for a block number-based configuration
    /// @param key The hash of the configuration name
    /// @param index The index of the checkpoint to retrieve
    /// @return The activation block number at the specified index
    function getActivationKeyBlockNumber(bytes32 key, uint256 index) external view returns (uint256) {
        return directory.getActivationKeyBlockNumber(key, index);
    }

    /// @notice Gets the activation timestamp at a specific index for a timestamp-based configuration
    /// @param key The hash of the configuration name
    /// @param index The index of the checkpoint to retrieve
    /// @return The activation timestamp at the specified index
    function getActivationKeyTimeStamp(bytes32 key, uint256 index) external view returns (uint256) {
        return directory.getActivationKeyTimeStamp(key, index);
    }

    /// @notice Retrieves all block number-based configs with activation blocks greater than or equal to the specified block number
    /// @param key The hash of the configuration name
    /// @param activationBlock The minimum activation block number to filter by
    /// @return activationKeys Array of activation block numbers for matching checkpoints
    /// @return configs Array of bytes32 configuration values for matching checkpoints (converted from bytes)
    function getAllConfigsGeActivationKeyBlockNumber(bytes32 key, uint256 activationBlock)
        external
        view
        returns (uint256[] memory, bytes32[] memory)
    {
        uint256 numCheckpoints = directory.getNumCheckpointsBlockNumber(key);
        uint256 count = 0;
        for (uint256 i = 0; i < numCheckpoints; i++) {
            if (directory.getActivationKeyBlockNumber(key, i) >= activationBlock) {
                count++;
            }
        }

        uint256[] memory activationKeys = new uint256[](count);
        bytes32[] memory configs = new bytes32[](count);
        uint256 index = 0;
        for (uint256 i = 0; i < numCheckpoints; i++) {
            uint256 activationKeyAtIdx = directory.getActivationKeyBlockNumber(key, i);
            if (activationKeyAtIdx >= activationBlock) {
                activationKeys[index] = activationKeyAtIdx;
                bytes memory config = directory.getConfigBlockNumber(key, i);
                configs[index] = bytes32(config);
                index++;
            }
        }

        return (activationKeys, configs);
    }

    /// @notice Retrieves all timestamp-based configs with activation timestamps greater than or equal to the specified timestamp
    /// @param key The hash of the configuration name
    /// @param activationKey The minimum activation timestamp to filter by
    /// @return activationKeys Array of activation timestamps for matching checkpoints
    /// @return configs Array of bytes configuration values for matching checkpoints
    function getAllConfigsGeActivationKeyTimeStamp(bytes32 key, uint256 activationKey)
        external
        view
        returns (uint256[] memory, bytes[] memory)
    {
        uint256 numCheckpoints = directory.getNumCheckpointsTimeStamp(key);
        uint256 count = 0;
        for (uint256 i = 0; i < numCheckpoints; i++) {
            if (directory.getActivationKeyTimeStamp(key, i) >= activationKey) {
                count++;
            }
        }

        uint256[] memory activationKeys = new uint256[](count);
        bytes[] memory configs = new bytes[](count);
        uint256 index = 0;
        for (uint256 i = 0; i < numCheckpoints; i++) {
            uint256 activationKeyAtIdx = directory.getActivationKeyTimeStamp(key, i);
            if (activationKeyAtIdx >= activationKey) {
                activationKeys[index] = activationKeyAtIdx;
                configs[index] = directory.getConfigTimeStamp(key, i);
                index++;
            }
        }

        return (activationKeys, configs);
    }
}
