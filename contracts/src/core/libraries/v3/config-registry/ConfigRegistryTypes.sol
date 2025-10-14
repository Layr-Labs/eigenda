// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

library ConfigRegistryTypes {
    /// @notice Struct to keep track of names associated with name digests
    /// @param names Mapping from name digest to name
    /// @param nameList List of all config names
    struct NameSet {
        mapping(bytes32 => string) names;
        string[] nameList;
    }

    /// @notice Struct to represent checkpoints for fixed-size byte32 configurations
    /// @param activationKey The activation key (e.g., block number or timestamp) for the checkpoint
    /// @param value The bytes32 configuration value at this checkpoint
    struct Bytes32Checkpoint {
        uint256 activationKey;
        bytes32 value;
    }

    /// @notice Struct to represent checkpoints for variable-size bytes configurations
    /// @param activationKey The activation key (e.g., block number or timestamp) for the checkpoint
    /// @param value The bytes configuration value at this checkpoint
    struct BytesCheckpoint {
        uint256 activationKey;
        bytes value;
    }

    /// @notice Struct to hold all bytes32 configuration checkpoints and associated names
    /// @param values Mapping from name digest to array of Bytes32Checkpoint structs. This entire structure is meant to be able to be queried.
    /// @param nameSet The NameSet struct to manage names associated with the configuration entries
    /// @dev See docs for the structs for more information
    struct Bytes32Cfg {
        mapping(bytes32 => Bytes32Checkpoint[]) values;
        NameSet nameSet;
    }

    /// @notice Struct to hold all bytes configuration checkpoints and associated names
    /// @dev See docs for the structs for more information
    struct BytesCfg {
        mapping(bytes32 => BytesCheckpoint[]) values;
        NameSet nameSet;
    }
}
