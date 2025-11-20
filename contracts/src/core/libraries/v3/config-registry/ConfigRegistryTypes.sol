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
    /// @param activationTime The activation timestamp for the checkpoint
    /// @param value The bytes32 configuration value at this checkpoint
    struct TimeStampCheckpoint {
        uint256 activationTime;
        bytes value;
    }

    /// @notice Struct to represent checkpoints for variable-size bytes configurations
    /// @param activationBlock The activation block number for the checkpoint
    /// @param value The bytes configuration value at this checkpoint
    struct BlockNumberCheckpoint {
        uint256 activationBlock;
        bytes value;
    }

    /// @notice Struct to hold all bytes32 configuration checkpoints and associated names
    /// @param values Mapping from name digest to array of TimeStampCheckpoint structs. This entire structure is meant to be able to be queried.
    /// @param nameSet The NameSet struct to manage names associated with the configuration entries
    /// @dev See docs for the structs for more information
    struct TimeStampCfg {
        mapping(bytes32 => TimeStampCheckpoint[]) values;
        NameSet nameSet;
    }

    /// @notice Struct to hold all bytes configuration checkpoints and associated names
    /// @dev See docs for the structs for more information
    struct BlockNumberCfg {
        mapping(bytes32 => BlockNumberCheckpoint[]) values;
        NameSet nameSet;
    }
}
