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

    event ConfigBytes32Added(string name, bytes32 indexed key, uint256 activationKey, bytes32 value);
    event ConfigBytesAdded(string name, bytes32 indexed key, uint256 activationKey, bytes value);
    event ConfigBytes32Replaced(string name, bytes32 indexed key, uint256 activationKey, bytes32 oldValue, bytes32 newValue);
    event ConfigBytesReplaced(string name, bytes32 indexed key, uint256 activationKey, bytes oldValue, bytes newValue);
    event ConfigBytes32Removed(string name, bytes32 indexed key);
    event ConfigBytesRemoved(string name, bytes32 indexed key);

    function getConfigKey(string memory name) external pure returns (bytes32);

    function getConfigNameBytes32(bytes32 key) external view returns (string memory);

    function getConfigNameBytes(bytes32 key) external view returns (string memory);

    function isKeyRegisteredBytes32(bytes32 key) external view returns (bool);

    function isKeyRegisteredBytes(bytes32 key) external view returns (bool);

    function addConfigBytes32(string memory name, uint256 activationKey, bytes32 value) external;

    function addConfigBytes(string memory name, uint256 activationKey, bytes memory value) external;

    function getNumConfigCheckpointsBytes32(bytes32 key) external view returns (uint256);

    function getNumConfigCheckpointsBytes(bytes32 key) external view returns (uint256);

    function getConfigBytes32(bytes32 key, uint256 index) external view returns (bytes32);

    function getConfigBytes(bytes32 key, uint256 index) external view returns (bytes memory);

    function getNumRegisteredConfigsBytes32() external view returns (uint256);

    function getNumRegisteredConfigsBytes() external view returns (uint256);

    function getRegisteredConfigBytes32(uint256 index) external view returns (string memory);

    function getRegisteredConfigBytes(uint256 index) external view returns (string memory);
}

/// @notice Interface for the EigenDA Directory
interface IEigenDADirectory is IEigenDAAddressDirectory, IEigenDAConfigRegistry {}
