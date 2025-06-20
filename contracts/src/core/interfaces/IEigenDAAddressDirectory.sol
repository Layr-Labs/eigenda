// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

interface IEigenDAAddressDirectory {
    event AddressSet(string name, address indexed value);

    /// @notice Sets the address of a contract.
    function setAddress(string memory name, address value) external;

    /// @notice Gets the address by keccak256 hash of the name.
    function getAddress(bytes32 key) external view returns (address);

    /// @notice Gets the address by name.
    function getAddress(string memory name) external view returns (address);
}
