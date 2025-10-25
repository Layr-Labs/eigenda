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

        // Count how many checkpoints are after the given activation key
        uint256 count = 0;
        for (uint256 i = 0; i < numCheckpoints; i++) {
            if (configRegistry.getActivationKeyBytes32(nameDigest, i) > activationKey) {
                count++;
            }
        }

        // Collect the checkpoints after the given activation key
        ConfigRegistryTypes.Bytes32Checkpoint[] memory result = new ConfigRegistryTypes.Bytes32Checkpoint[](count);
        uint256 index = 0;
        for (uint256 i = 0; i < numCheckpoints; i++) {
            uint256 checkpointActivationKey = configRegistry.getActivationKeyBytes32(nameDigest, i);
            if (checkpointActivationKey > activationKey) {
                result[index] = configRegistry.getCheckpointBytes32(nameDigest, i);
                index++;
            }
        }

        return result;
    }

    /// @notice Retrieves all bytes configuration checkpoints for a given name that have an activation key greater than the provided activation key.
    function getAllBytesConfigsAfterKey(string memory name, uint256 activationKey)
        external
        view
        returns (ConfigRegistryTypes.BytesCheckpoint[] memory)
    {
        bytes32 nameDigest = ConfigRegistryLib.getNameDigest(name);
        uint256 numCheckpoints = configRegistry.getNumCheckpointsBytes(nameDigest);

        // Count how many checkpoints are after the given activation key
        uint256 count = 0;
        for (uint256 i = 0; i < numCheckpoints; i++) {
            if (configRegistry.getActivationKeyBytes(nameDigest, i) > activationKey) {
                count++;
            }
        }
        // Collect the checkpoints after the given activation key
        ConfigRegistryTypes.BytesCheckpoint[] memory result = new ConfigRegistryTypes.BytesCheckpoint[](count);
        uint256 index = 0;
        for (uint256 i = 0; i < numCheckpoints; i++) {
            uint256 checkpointActivationKey = configRegistry.getActivationKeyBytes(nameDigest, i);
            if (checkpointActivationKey > activationKey) {
                result[index] = configRegistry.getCheckpointBytes(nameDigest, i);
                index++;
            }
        }
        return result;
    }
}
