// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {AddressDirectoryLib} from "src/core/libraries/v3/address-directory/AddressDirectoryLib.sol";
import {IEigenDADirectory} from "src/core/interfaces/IEigenDADirectory.sol";

contract EigenDADirectory is OwnableUpgradeable, IEigenDADirectory {
    using AddressDirectoryLib for string;
    using AddressDirectoryLib for bytes32;

    function initialize(address _initialOwner) external initializer {
        _transferOwnership(_initialOwner);
    }

    /// @inheritdoc IEigenDADirectory
    function addAddress(string memory name, address value) external onlyOwner {
        bytes32 key = name.getKey();

        if (value == address(0)) {
            revert ZeroAddress();
        }
        if (key.getAddress() != address(0)) {
            revert AddressAlreadyExists(name);
        }

        key.setAddress(value);

        emit AddressAdded(name, key, value);
    }

    /// @inheritdoc IEigenDADirectory
    function replaceAddress(string memory name, address value) external onlyOwner {
        bytes32 key = name.getKey();
        address oldValue = key.getAddress();

        if (oldValue == address(0)) {
            revert AddressNotAdded(name);
        }
        if (value == address(0)) {
            revert ZeroAddress();
        }
        if (oldValue == value) {
            revert NewValueIsOldValue(value);
        }

        key.setAddress(value);

        emit AddressReplaced(name, key, oldValue, value);
    }

    /// @inheritdoc IEigenDADirectory
    function removeAddress(string memory name) external onlyOwner {
        bytes32 key = name.getKey();
        address existingAddress = key.getAddress();

        if (existingAddress == address(0)) {
            revert AddressNotAdded(name);
        }

        key.setAddress(address(0));

        emit AddressRemoved(name, key);
    }

    /// @inheritdoc IEigenDADirectory
    function getAddress(string memory name) external view returns (address) {
        return name.getKey().getAddress();
    }

    /// @inheritdoc IEigenDADirectory
    function getAddress(bytes32 key) external view returns (address) {
        return key.getAddress();
    }
}
