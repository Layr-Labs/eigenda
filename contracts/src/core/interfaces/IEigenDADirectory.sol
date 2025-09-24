// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

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
/// @notice Interface for a configuration registry that allows adding, replacing, removing, and retrieving configuration entries by name.
///         Supports both bytes32 and bytes types for configuration values, along with optional extra information.
/// @dev    Contracts should use the bytes32 functions for gas efficiency.
interface IEigenDAConfigRegistry {
    error ConfigAlreadyExists(string name);
    error ConfigDoesNotExist(string name);

    event ConfigBytes32Added(string name, bytes32 indexed key, bytes32 value, string extraInfo);
    event ConfigBytesAdded(string name, bytes32 indexed key, bytes value, string extraInfo);
    event ConfigBytes32Replaced(string name, bytes32 indexed key, bytes32 oldValue, bytes32 newValue, string extraInfo);
    event ConfigBytesReplaced(string name, bytes32 indexed key, bytes oldValue, bytes newValue, string extraInfo);
    event ConfigBytes32Removed(string name, bytes32 indexed key);
    event ConfigBytesRemoved(string name, bytes32 indexed key);

    /// @notice Adds a new bytes32 configuration entry to the registry by name, along with optional extra info.
    /// @dev Reverts if the entry already exists.
    function addConfigBytes32(string memory name, bytes32 value, string memory extraInfo) external;

    /// @notice Adds a new bytes configuration entry to the registry by name, along with optional extra info.
    /// @dev Reverts if the entry already exists.
    function addConfigBytes(string memory name, bytes memory value, string memory extraInfo) external;

    /// @notice Replaces an existing bytes32 configuration entry in the registry by name, along with optional extra info.
    /// @dev Reverts if the entry does not exist.
    function replaceConfigBytes32(string memory name, bytes32 value, string memory extraInfo) external;

    /// @notice Replaces an existing bytes configuration entry in the registry by name, along with optional extra info.
    /// @dev Reverts if the entry does not exist.
    function replaceConfigBytes(string memory name, bytes memory value, string memory extraInfo) external;

    /// @notice Removes an existing bytes32 configuration entry from the registry by name.
    /// @dev Reverts if the entry does not exist.
    function removeConfigBytes32(string memory name) external;

    /// @notice Removes an existing bytes configuration entry from the registry by name.
    /// @dev Reverts if the entry does not exist.
    function removeConfigBytes(string memory name) external;

    /// @notice Gets the bytes32 configuration entry by name.
    function getConfigBytes32(string memory name) external view returns (bytes32);

    /// @notice Gets the bytes configuration entry by name.
    function getConfigBytes(string memory name) external view returns (bytes memory);

    /// @notice Gets the extra info associated with a bytes32 configuration entry by name.
    function getConfigBytes32ExtraInfo(string memory name) external view returns (string memory);

    /// @notice Gets the extra info associated with a bytes configuration entry by name.
    function getConfigBytesExtraInfo(string memory name) external view returns (string memory);

    /// @notice Gets the bytes32 configuration entry by keccak256 hash of the name.
    function getConfigBytes32(bytes32 key) external view returns (bytes32);

    /// @notice Gets the bytes configuration entry by keccak256 hash of the name.
    function getConfigBytes(bytes32 key) external view returns (bytes memory);

    /// @notice Gets the extra info associated with a bytes32 configuration entry by keccak256 hash of the name.
    function getConfigBytes32ExtraInfo(bytes32 key) external view returns (string memory);

    /// @notice Gets the extra info associated with a bytes configuration entry by keccak256 hash of the name.
    function getConfigBytesExtraInfo(bytes32 key) external view returns (string memory);

    /// @notice Gets the number of registered keys for bytes32 configuration entries.
    function getNumRegisteredKeysBytes32() external view returns (uint256);

    /// @notice Gets the number of registered keys for bytes configuration entries.
    function getNumRegisteredKeysBytes() external view returns (uint256);

    /// @notice Gets the registered key for bytes32 configuration entries at a specific index.
    function getRegisteredKeyBytes32(uint256 index) external view returns (string memory);

    /// @notice Gets the registered key for bytes configuration entries at a specific index.
    function getRegisteredKeyBytes(uint256 index) external view returns (string memory);
}

/// @notice Interface for the EigenDA Directory
interface IEigenDADirectory is IEigenDAAddressDirectory, IEigenDAConfigRegistry {}
