// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {ConfigRegistryStorage as S} from "./ConfigRegistryStorage.sol";
import {ConfigRegistryTypes as T} from "./ConfigRegistryTypes.sol";

library ConfigRegistryLib {
    event TimestampConfigBytesSet(bytes32 nameDigest, uint256 activationTS, bytes value);
    event BlockNumberConfigBytesSet(bytes32 nameDigest, uint256 abn, bytes value);

    /// @notice Thrown when attempting to retrieve a configuration by an unregistered name digest
    /// @param nameDigest The unregistered name digest
    error NameDigestNotRegistered(bytes32 nameDigest);

    /// @notice Thrown when trying to add a configuration with a timestamp that is not strictly increasing
    /// @param prevTS The last activation timestamp for this configuration
    /// @param newTS The timestamp being added (must be > prevTS)
    error NotIncreasingTimestamp(uint256 prevTS, uint256 newTS);

    /// @notice Thrown when trying to add a configuration with a block number that is not strictly increasing
    /// @param prevABN The last activation block number for this configuration
    /// @param newABN The activation block number being added (must be > prevABN)
    error NotIncreasingBlockNumber(uint256 prevABN, uint256 newABN);

    /// @notice Thrown when adding the first block number configuration with an activation block in the past
    /// @param currBlock The current block number (sourced via block.number)
    /// @param abn The activation block number being added (must be >= currBlock)
    error BlockNumberActivationInPast(uint256 currBlock, uint256 abn);

    /// @notice Thrown when adding the first timestamp configuration with an activation timestamp in the past
    /// @param currTS The current timestamp (sourced via block.timestamp)
    /// @param activationTS The activation timestamp being added (must be >= currTS)
    error TimeStampActivationInPast(uint256 currTS, uint256 activationTS);

    /// @notice Computes the keccak256 hash of a configuration name
    /// @param name The configuration name
    /// @return The keccak256 hash of the packed name
    function getNameDigest(string memory name) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(name));
    }

    /// @notice Gets the number of checkpoints for a timestamp-based configuration entry
    /// @param nameDigest The hash of the configuration name
    /// @return The number of checkpoints stored for this configuration
    function getNumCheckpointsTimeStamp(bytes32 nameDigest) internal view returns (uint256) {
        return S.layout().timestampCfg.values[nameDigest].length;
    }

    /// @notice Gets the number of checkpoints for a block number-based configuration entry
    /// @param nameDigest The hash of the configuration name
    /// @return The number of checkpoints stored for this configuration
    function getNumCheckpointsBlockNumber(bytes32 nameDigest) internal view returns (uint256) {
        return S.layout().blockNumberCfg.values[nameDigest].length;
    }

    /// @notice Gets the configuration value at a specific index for a timestamp-based configuration
    /// @param nameDigest The hash of the configuration name
    /// @param index The index of the checkpoint to retrieve
    /// @return The bytes configuration value at the specified index
    function getConfigTimeStamp(bytes32 nameDigest, uint256 index) internal view returns (bytes memory) {
        return S.layout().timestampCfg.values[nameDigest][index].value;
    }

    /// @notice Gets the configuration value at a specific index for a block number-based configuration
    /// @param nameDigest The hash of the configuration name
    /// @param index The index of the checkpoint to retrieve
    /// @return The bytes configuration value at the specified index
    function getConfigBlockNumber(bytes32 nameDigest, uint256 index) internal view returns (bytes memory) {
        return S.layout().blockNumberCfg.values[nameDigest][index].value;
    }

    /// @notice Gets the activation timestamp at a specific index for a timestamp-based configuration
    /// @param nameDigest The hash of the configuration name
    /// @param index The index of the checkpoint to retrieve
    /// @return The activation timestamp at the specified index
    function getActivationTimeStamp(bytes32 nameDigest, uint256 index) internal view returns (uint256) {
        return S.layout().timestampCfg.values[nameDigest][index].activationTime;
    }

    /// @notice Gets the activation block number at a specific index for a block number-based configuration
    /// @param nameDigest The hash of the configuration name
    /// @param index The index of the checkpoint to retrieve
    /// @return The activation block number at the specified index
    function getActivationBlockNumber(bytes32 nameDigest, uint256 index) internal view returns (uint256) {
        return S.layout().blockNumberCfg.values[nameDigest][index].activationBlock;
    }

    /// @notice Gets the full checkpoint at a specific index for a timestamp-based configuration
    /// @param nameDigest The hash of the configuration name
    /// @param index The index of the checkpoint to retrieve
    /// @return The TimeStampCheckpoint containing both value and activation timestamp
    function getCheckpointTimeStamp(bytes32 nameDigest, uint256 index)
        internal
        view
        returns (T.TimeStampCheckpoint memory)
    {
        return S.layout().timestampCfg.values[nameDigest][index];
    }

    /// @notice Gets the full checkpoint at a specific index for a block number-based configuration
    /// @param nameDigest The hash of the configuration name
    /// @param index The index of the checkpoint to retrieve
    /// @return The BlockNumberCheckpoint containing both value and activation block number
    function getCheckpointBlockNumber(bytes32 nameDigest, uint256 index)
        internal
        view
        returns (T.BlockNumberCheckpoint memory)
    {
        return S.layout().blockNumberCfg.values[nameDigest][index];
    }

    /// @notice Adds a new timestamp-based configuration checkpoint
    /// @param nameDigest The hash of the configuration name
    /// @param activationTS The activation timestamp (must be > last activation timestamp for this config)
    /// @param value The bytes configuration value
    /// @dev For the first checkpoint, activationTS must be >= block.timestamp
    /// @dev Subsequent checkpoints must have strictly increasing activation timestamps
    function addConfigTimeStamp(bytes32 nameDigest, uint256 activationTS, bytes memory value) internal {
        T.TimestampConfig storage cfg = S.layout().timestampCfg;
        if (cfg.values[nameDigest].length > 0) {
            uint256 lastActivationTS = cfg.values[nameDigest][cfg.values[nameDigest].length - 1].activationTime;
            if (activationTS <= lastActivationTS) {
                revert NotIncreasingTimestamp(lastActivationTS, activationTS);
            }
        }

        /// @dev activation timestamps being provided must always be at a future timestamp
        if (activationTS < block.timestamp) {
            revert TimeStampActivationInPast(block.timestamp, activationTS);
        }

        cfg.values[nameDigest].push(T.TimeStampCheckpoint({value: value, activationTime: activationTS}));
        emit TimestampConfigBytesSet(nameDigest, activationTS, value);
    }

    /// @notice Adds a new block number-based configuration checkpoint
    /// @param nameDigest The hash of the configuration name
    /// @param abn The activation block number (must be > last activation block for this config)
    /// @param value The bytes configuration value
    /// @dev For the first checkpoint, abn must be >= block.number
    /// @dev Subsequent checkpoints must have strictly increasing activation block numbers
    function addConfigBlockNumber(bytes32 nameDigest, uint256 abn, bytes memory value) internal {
        T.BlockNumberConfig storage cfg = S.layout().blockNumberCfg;
        if (cfg.values[nameDigest].length > 0) {
            uint256 lastABN = cfg.values[nameDigest][cfg.values[nameDigest].length - 1].activationBlock;
            if (abn <= lastABN) {
                revert NotIncreasingBlockNumber(lastABN, abn);
            }
        }

        /// @dev abn being provided must always be at a future block
        if (abn < block.number) {
            revert BlockNumberActivationInPast(block.number, abn);
        }

        cfg.values[nameDigest].push(T.BlockNumberCheckpoint({value: value, activationBlock: abn}));
        emit BlockNumberConfigBytesSet(nameDigest, abn, value);
    }

    /// @notice Registers a configuration name for timestamp-based configurations
    /// @param name The configuration name to register
    /// @dev Idempotent - safe to call multiple times with the same name
    function registerNameTimeStamp(string memory name) internal {
        registerName(S.layout().timestampCfg.nameSet, name);
    }

    /// @notice Registers a configuration name for block number-based configurations
    /// @param name The configuration name to register
    /// @dev Idempotent - safe to call multiple times with the same name
    function registerNameBlockNumber(string memory name) internal {
        registerName(S.layout().blockNumberCfg.nameSet, name);
    }

    /// @notice Internal function to register a configuration name in a name set
    /// @param nameSet The name set to register the name in
    /// @param name The configuration name to register
    /// @dev Only adds the name if it hasn't been registered before
    function registerName(T.NameSet storage nameSet, string memory name) internal {
        bytes32 nameDigest = getNameDigest(name);
        if (bytes(nameSet.names[nameDigest]).length == 0) {
            require(bytes(name).length > 0, "Name cannot be empty");
            nameSet.names[nameDigest] = name;
            nameSet.nameList.push(name);
        }
    }

    /// @notice Checks if a name digest is registered in a given name set
    /// @param nameSet The name set to check
    /// @param nameDigest The hash of the name to check
    /// @return True if the name digest is registered, false otherwise
    function isNameDigestRegistered(T.NameSet storage nameSet, bytes32 nameDigest) internal view returns (bool) {
        return bytes(nameSet.names[nameDigest]).length > 0;
    }

    /// @notice Checks if a name digest is registered for timestamp-based configurations
    /// @param nameDigest The hash of the name to check
    /// @return True if registered, false otherwise
    function isNameRegisteredTimeStamp(bytes32 nameDigest) internal view returns (bool) {
        return isNameDigestRegistered(S.layout().timestampCfg.nameSet, nameDigest);
    }

    /// @notice Checks if a name digest is registered for block number-based configurations
    /// @param nameDigest The hash of the name to check
    /// @return True if registered, false otherwise
    function isNameRegisteredBlockNumber(bytes32 nameDigest) internal view returns (bool) {
        return isNameDigestRegistered(S.layout().blockNumberCfg.nameSet, nameDigest);
    }

    /// @notice Gets the total number of registered timestamp-based configuration names
    /// @return The count of registered timestamp-based configuration names
    function getNumRegisteredNamesTimeStamp() internal view returns (uint256) {
        return S.layout().timestampCfg.nameSet.nameList.length;
    }

    /// @notice Gets the total number of registered block number-based configuration names
    /// @return The count of registered block number-based configuration names
    function getNumRegisteredNamesBlockNumber() internal view returns (uint256) {
        return S.layout().blockNumberCfg.nameSet.nameList.length;
    }

    /// @notice Gets a registered timestamp-based configuration name by its index in the name list
    /// @param index The index of the name to retrieve
    /// @return The configuration name at the specified index
    function getRegisteredNameTimeStamp(uint256 index) internal view returns (string memory) {
        return S.layout().timestampCfg.nameSet.nameList[index];
    }

    /// @notice Gets a registered block number-based configuration name by its index in the name list
    /// @param index The index of the name to retrieve
    /// @return The configuration name at the specified index
    function getRegisteredNameBlockNumber(uint256 index) internal view returns (string memory) {
        return S.layout().blockNumberCfg.nameSet.nameList[index];
    }

    /// @notice Gets the configuration name for a timestamp-based configuration by its name digest
    /// @param nameDigest The hash of the configuration name
    /// @return The configuration name
    /// @dev Reverts with NameDigestNotRegistered if the name digest is not registered
    function getNameTimeStamp(bytes32 nameDigest) internal view returns (string memory) {
        string memory name = S.layout().timestampCfg.nameSet.names[nameDigest];
        if (bytes(name).length == 0) {
            revert NameDigestNotRegistered(nameDigest);
        }
        return name;
    }

    /// @notice Gets the configuration name for a block number-based configuration by its name digest
    /// @param nameDigest The hash of the configuration name
    /// @return The configuration name
    /// @dev Reverts with NameDigestNotRegistered if the name digest is not registered
    function getNameBlockNumber(bytes32 nameDigest) internal view returns (string memory) {
        string memory name = S.layout().blockNumberCfg.nameSet.names[nameDigest];
        if (bytes(name).length == 0) {
            revert NameDigestNotRegistered(nameDigest);
        }
        return name;
    }

    /// @notice Gets the list of all registered timestamp-based configuration names
    /// @return An array containing all registered timestamp-based configuration names
    function getNameListTimeStamp() internal view returns (string[] memory) {
        return S.layout().timestampCfg.nameSet.nameList;
    }

    /// @notice Gets the list of all registered block number-based configuration names
    /// @return An array containing all registered block number-based configuration names
    function getNameListBlockNumber() internal view returns (string[] memory) {
        return S.layout().blockNumberCfg.nameSet.nameList;
    }
}
