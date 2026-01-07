// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {ConfigRegistryTypes} from "src/core/libraries/v3/config-registry/ConfigRegistryTypes.sol";

interface IEigenDAAddressDirectory {
    error AddressAlreadyExists(string name);
    error AddressDoesNotExist(string name);
    error ZeroAddress();
    error NewValueIsOldValue(address value);

    event AddressAdded(string name, bytes32 indexed key, address indexed value);
    event AddressReplaced(string name, bytes32 indexed key, address indexed oldValue, address indexed newValue);
    event AddressRemoved(string name, bytes32 indexed key);

    /// @notice Adds a new address to the directory by name.
    /// @dev Fails if the address is zero or if an address with the same name already exists.
    ///      Emits an AddressAdded event on success.
    function addAddress(string memory name, address value) external;

    /// @notice Replaces an existing address in the directory by name.
    /// @dev Fails if the address is zero, if the address with the name does not exist, or if the new value is the same as the old value.
    ///      Emits an AddressReplaced event on success.
    function replaceAddress(string memory name, address value) external;

    /// @notice Removes an address from the directory by name.
    /// @dev Fails if the address with the name does not exist.
    ///      Emits an AddressRemoved event on success.
    function removeAddress(string memory name) external;

    /// @notice Gets the address by keccak256 hash of the name.
    /// @dev    This entry point is cheaper in gas because it avoids needing to compute the key from the name.
    function getAddress(bytes32 key) external view returns (address);

    /// @notice Gets the address by name.
    function getAddress(string memory name) external view returns (address);

    /// @notice Gets the name by keccak256 hash of the name.
    function getName(bytes32 key) external view returns (string memory);

    /// @notice Gets all names in the directory.
    function getAllNames() external view returns (string[] memory);
}

/// @title IEigenDAConfigRegistry
/// @notice Interface for a configuration registry that allows adding and retrieving configuration entries by name.
///         Supports bytes types for configuration values, and maintains a checkpointed structure for each configuration entry
///         by an arbitrary activation key.
interface IEigenDAConfigRegistry {
    /// @notice Adds a variable length byte configuration value to the configuration registry using block number as activation key.
    /// @param name The name of the configuration entry.
    /// @param abn The activation block number for the configuration entry.
    /// @param value The variable length byte configuration value.
    /// @dev The abn must be strictly greater than the last abn for the same name and must be greater than the current block number.
    function addConfigBlockNumber(string memory name, uint256 abn, bytes memory value) external;

    /// @notice Adds a variable length byte configuration value to the configuration registry using timestamp as activation key.
    /// @param name The name of the configuration entry.
    /// @param activationTS The activation timestamp for the configuration entry.
    /// @param value The variable length byte configuration value.
    /// @dev The activationTS must be strictly greater than the last activationTS for the same name and greater than the current block timestamp.
    function addConfigTimeStamp(string memory name, uint256 activationTS, bytes memory value) external;

    /// @notice Gets the number of checkpoints for a block number configuration entry.
    /// @param nameDigest The hash of the name of the configuration entry.
    /// @return The number of checkpoints for the configuration entry.
    function getNumCheckpointsBlockNumber(bytes32 nameDigest) external view returns (uint256);

    /// @notice Gets the number of checkpoints for a timestamp configuration entry.
    /// @param nameDigest The hash of the name of the configuration entry.
    /// @return The number of checkpoints for the configuration entry.
    function getNumCheckpointsTimeStamp(bytes32 nameDigest) external view returns (uint256);

    /// @notice Gets the block number configuration value at a specific index for a configuration entry.
    /// @param nameDigest The hash of the name of the configuration entry.
    /// @param index The index of the configuration value to retrieve.
    /// @return The variable length byte configuration value at the specified index.
    function getConfigBlockNumber(bytes32 nameDigest, uint256 index) external view returns (bytes memory);

    /// @notice Gets the timestamp configuration value at a specific index for a configuration entry.
    /// @param nameDigest The hash of the name of the configuration entry.
    /// @param index The index of the configuration value to retrieve.
    /// @return The variable length byte configuration value at the specified index.
    function getConfigTimeStamp(bytes32 nameDigest, uint256 index) external view returns (bytes memory);

    /// @notice Gets the activation key for a block number configuration entry at a specific index.
    /// @param nameDigest The hash of the name of the configuration entry.
    /// @param index The index of the configuration value to retrieve the activation key for.
    /// @return The activation key at the specified index.
    function getActivationBlockNumber(bytes32 nameDigest, uint256 index) external view returns (uint256);

    /// @notice Gets the activation key for a timestamp configuration entry at a specific index.
    /// @param nameDigest The hash of the name of the configuration entry.
    /// @param index The index of the configuration value to retrieve the activation key for.
    /// @return The activation key at the specified index.
    function getActivationTimeStamp(bytes32 nameDigest, uint256 index) external view returns (uint256);

    /// @notice Gets the full checkpoint (value and activation key) for a timestamp configuration entry at a specific index.
    /// @param nameDigest The hash of the name of the configuration entry.
    /// @param index The index of the configuration value to retrieve the checkpoint for.
    /// @return The full checkpoint (value and activation key) at the specified index.
    function getCheckpointTimeStamp(bytes32 nameDigest, uint256 index)
        external
        view
        returns (ConfigRegistryTypes.TimeStampCheckpoint memory);

    /// @notice Gets the full checkpoint (value and activation key) for a block number configuration entry at a specific index.
    /// @param nameDigest The hash of the name of the configuration entry.
    /// @param index The index of the configuration value to retrieve the checkpoint for.
    /// @return The full checkpoint (value and activation key) at the specified index.
    function getCheckpointBlockNumber(bytes32 nameDigest, uint256 index)
        external
        view
        returns (ConfigRegistryTypes.BlockNumberCheckpoint memory);

    /// @notice Gets the name of a block number configuration entry by its name digest.
    /// @param nameDigest The hash of the name of the configuration entry.
    /// @return The name of the configuration entry.
    function getConfigNameBlockNumber(bytes32 nameDigest) external view returns (string memory);

    /// @notice Gets the name of a timestamp configuration entry by its name digest.
    /// @param nameDigest The hash of the name of the configuration entry.
    /// @return The name of the configuration entry.
    function getConfigNameTimeStamp(bytes32 nameDigest) external view returns (string memory);

    /// @notice Gets all names of block number configuration entries.
    /// @return An array of all configuration entry names.
    function getAllConfigNamesBlockNumber() external view returns (string[] memory);

    /// @notice Gets all names of timestamp configuration entries.
    /// @return An array of all configuration entry names.
    function getAllConfigNamesTimeStamp() external view returns (string[] memory);

    /// @notice Retrieves the currently active block number config checkpoint and all future checkpoints for a given name.
    ///         this is only expected to be used via eth_calls by offchain EigenDA services.
    /// @param name the config string name
    /// @param referenceBlockNumber the reference block number used for filtered lookups against the checkpoints
    /// @return checkpoints with the highest activation block that is less than or equal to the provided reference block,
    ///      plus all checkpoints with activation block numbers greater than the provided reference block.
    ///      This allows offchain clients to know the current configuration value and plan ahead for upcoming updates.
    function getActiveAndFutureBlockNumberConfigs(string memory name, uint256 referenceBlockNumber)
        external
        view
        returns (ConfigRegistryTypes.BlockNumberCheckpoint[] memory);

    /// @notice Retrieves the currently active timestamp config checkpoint and all future checkpoints for a given name.
    ///         this is only expected to be used via eth_calls by offchain EigenDA services.
    /// @param name the config string name
    /// @param referenceTimestamp the reference timestamp used for filtered lookups against the checkpoints
    /// @return checkpoints with the highest activation timestamp that is less than or equal to the provided reference timestamp,
    ///      plus all checkpoints with activation timestamps greater than the provided reference timestamp.
    ///      This allows offchain clients to know the current configuration value and plan ahead for upcoming updates.
    function getActiveAndFutureTimestampConfigs(string memory name, uint256 referenceTimestamp)
        external
        view
        returns (ConfigRegistryTypes.TimeStampCheckpoint[] memory);
}

/// @notice Interface for the EigenDA Directory
interface IEigenDADirectory is IEigenDAAddressDirectory, IEigenDAConfigRegistry {}
