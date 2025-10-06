// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDADirectory} from "src/core/EigenDADirectory.sol";

contract MockConfigRetriever {
    EigenDADirectory public directory;

    constructor(EigenDADirectory _directory) {
        directory = _directory;
    }

    function getNumCheckpointsBytes32(bytes32 key) external view returns (uint256) {
        return directory.getNumCheckpointsBytes32(key);
    }

    function getNumCheckpointsBytes(bytes32 key) external view returns (uint256) {
        return directory.getNumCheckpointsBytes(key);
    }

    function getConfigBytes32(bytes32 key, uint256 index) external view returns (bytes32) {
        return directory.getConfigBytes32(key, index);
    }

    function getConfigBytes(bytes32 key, uint256 index) external view returns (bytes memory) {
        return directory.getConfigBytes(key, index);
    }

    function getActivationKey(bytes32 key, uint256 index) external view returns (uint256) {
        return directory.getActivationKeyBytes32(key, index);
    }

    function getActivationKeyBytes(bytes32 key, uint256 index) external view returns (bytes32) {
        return directory.getActivationKeyBytes(key, index);
    }

    /// @notice Retrieves all configs with activation keys greater than or equal to the specified activation key.
    function getAllConfigsGeActivationKeyBytes32(bytes32 key, uint256 activationKey) external view returns (uint256[] memory, bytes32[] memory) {
        uint256 numCheckpoints = directory.getNumCheckpointsBytes32(key);
        uint256 count = 0;
        for (uint256 i = 0; i < numCheckpoints; i++) {
            if (directory.getActivationKeyBytes32(key, i) >= activationKey) {
                count++;
            }
        }

        uint256[] memory activationKeys = new uint256[](count);
        bytes32[] memory configs = new bytes32[](count);
        uint256 index = 0;
        for (uint256 i = 0; i < numCheckpoints; i++) {
            if (directory.getActivationKeyBytes32(key, i) >= activationKey) {
                activationKeys[index] = directory.getActivationKeyBytes32(key, i);
                configs[index] = directory.getConfigBytes32(key, i);
                index++;
            }
        }

        return (activationKeys, configs);
    }

    /// @notice Retrieves all configs with activation keys greater than or equal to the specified activation key.
    function getAllConfigsGeActivationKeyBytes(bytes32 key, uint256 activationKey) external view returns (uint256[] memory, bytes[] memory) {
        uint256 numCheckpoints = directory.getNumCheckpointsBytes(key);
        uint256 count = 0;
        for (uint256 i = 0; i < numCheckpoints; i++) {
            if (directory.getActivationKeyBytes(key, i) >= activationKey) {
                count++;
            }
        }

        uint256[] memory activationKeys = new uint256[](count);
        bytes[] memory configs = new bytes[](count);
        uint256 index = 0;
        for (uint256 i = 0; i < numCheckpoints; i++) {
            if (directory.getActivationKeyBytes(key, i) >= activationKey) {
                activationKeys[index] = directory.getActivationKeyBytes(key, i);
                configs[index] = directory.getConfigBytes(key, i);
                index++;
            }
        }

        return (activationKeys, configs);
    }
}