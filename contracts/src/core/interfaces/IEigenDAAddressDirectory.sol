// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

/// @notice Interface for the EigenDA Directory
///         This interface currently only includes functions for managing a directory of addresses by name.
///         In the future, it may be extended to include access control as well.
interface IEigenDADirectory {
    error AddressAlreadyExists(string name);
    error AddressNotAdded(string name);
    error ZeroAddress();
    error NewValueIsOldValue(address value);

    event AddressAdded(string name, bytes32 indexed key, address indexed value);
    event AddressReplaced(string name, bytes32 indexed key, address indexed oldValue, address indexed newValue);
    event AddressRemoved(string name, bytes32 indexed key);

    /// @notice Adds a new address to the directory by name.
    function addAddress(string memory name, address value) external;

    /// @notice Replaces an existing address in the directory by name.
    function replaceAddress(string memory name, address value) external;

    /// @notice Removes an address from the directory by name.
    function removeAddress(string memory name) external;

    /// @notice Gets the address by keccak256 hash of the name.
    function getAddress(bytes32 key) external view returns (address);

    /// @notice Gets the address by name.
    function getAddress(string memory name) external view returns (address);
}
