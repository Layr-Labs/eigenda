// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAConfigRegistry} from "src/core/interfaces/IEigenDADirectory.sol";
import {ConfigRegistryTypes} from "src/core/libraries/v3/config-registry/ConfigRegistryTypes.sol";
import {ConfigRegistryLib} from "src/core/libraries/v3/config-registry/ConfigRegistryLib.sol";

/// @title EigenDAConfigRetriever
/// @notice A stateless contract that defines getter functions using the underlying configuration registry's primitive API.
contract EigenDAConfigRetriever {
    IEigenDAConfigRegistry public immutable configRegistry;

    constructor(address configRegistryAddress) {
        configRegistry = IEigenDAConfigRegistry(configRegistryAddress);
    }

    /// @notice Retrieves all bytes32 configuration checkpoints for a given name that have an activation key greater than the provided activation key.
    function getAllBytes32ConfigsAfterKey(string memory name, uint256 activationKey)
        external
        view
        returns (ConfigRegistryTypes.Bytes32Checkpoint[] memory)
    {
        bytes32 nameDigest = ConfigRegistryLib.getNameDigest(name);
        uint256 numCheckpoints = configRegistry.getNumCheckpointsBytes32(nameDigest);

        // There are 3 cases to handle:
        // 1. If no checkpoints have activation keys greater than the provided activation key, we return an empty array.
        // 2. If all checkpoints have activation keys greater than the provided activation key, we return the entire array.
        // 3. If some checkpoints have activation keys greater than the provided activation key, we return the relevant subset.

        uint256 startIndex = numCheckpoints; // Default to numCheckpoints (case 1)
        for (uint256 i = 0; i < numCheckpoints; i++) {
            uint256 checkpointActivationKey = configRegistry.getActivationKeyBytes32(nameDigest, numCheckpoints - 1 - i);
            if (checkpointActivationKey <= activationKey) {
                startIndex = numCheckpoints - i; // Found the first checkpoint with activation key <= provided key
                break;
            }
        }
        // Collect the checkpoints from startIndex to the end
        uint256 resultCount = numCheckpoints - startIndex;
        ConfigRegistryTypes.Bytes32Checkpoint[] memory results = new ConfigRegistryTypes.Bytes32Checkpoint[](
            resultCount
        );
        for (uint256 i = 0; i < resultCount; i++) {
            results[i] = configRegistry.getCheckpointBytes32(nameDigest, startIndex + i);
        }
        return results;
    }

    /// @notice Retrieves all bytes configuration checkpoints for a given name that have an activation key greater than the provided activation key.
    function getAllBytesConfigsAfterKey(string memory name, uint256 activationKey)
        external
        view
        returns (ConfigRegistryTypes.BytesCheckpoint[] memory)
    {
        bytes32 nameDigest = ConfigRegistryLib.getNameDigest(name);
        uint256 numCheckpoints = configRegistry.getNumCheckpointsBytes(nameDigest);
        
        // There are 3 cases to handle:
        // 1. If no checkpoints have activation keys greater than the provided activation key, we return an empty array.
        // 2. If all checkpoints have activation keys greater than the provided activation key, we return the entire array.
        // 3. If some checkpoints have activation keys greater than the provided activation key, we return the relevant subset.

        uint256 startIndex = numCheckpoints; // Default to numCheckpoints (case 1)
        for (uint256 i = 0; i < numCheckpoints; i++) {
            uint256 checkpointActivationKey = configRegistry.getActivationKeyBytes(nameDigest, numCheckpoints - 1 - i);
            if (checkpointActivationKey <= activationKey) {
                startIndex = numCheckpoints - i; // Found the first checkpoint with activation key <= provided key
                break;
            }
        }
        // Collect the checkpoints from startIndex to the end
        uint256 resultCount = numCheckpoints - startIndex;
        ConfigRegistryTypes.BytesCheckpoint[] memory results = new ConfigRegistryTypes.BytesCheckpoint[](
            resultCount
        );
        for (uint256 i = 0; i < resultCount; i++) {
            results[i] = configRegistry.getCheckpointBytes(nameDigest, startIndex + i);
        }
        return results;
    }
}
