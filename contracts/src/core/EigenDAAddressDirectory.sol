// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {AddressDirectoryLib} from "src/core/libraries/v3/address-directory/AddressDirectoryLib.sol";
import {IEigenDAAddressDirectory} from "src/core/interfaces/IEigenDAAddressDirectory.sol";

contract EigenDAAddressDirectory is OwnableUpgradeable, IEigenDAAddressDirectory {
    using AddressDirectoryLib for string;
    using AddressDirectoryLib for bytes32;

    function initialize(address _initialOwner) external initializer {
        _transferOwnership(_initialOwner);
    }

    /// @inheritdoc IEigenDAAddressDirectory
    function addAddress(string memory name, address value) external onlyOwner {
        bytes32 key = name.getKey();

        if (value == address(0)) {
            revert InvalidAddress(name);
        }
        if (key.getAddress() != address(0)) {
            revert AddressAlreadyExists(name);
        }

        key.setAddress(value);

        emit AddressAdded(name, key, value);
    }

    /// @inheritdoc IEigenDAAddressDirectory
    function replaceAddress(string memory name, address value) external onlyOwner {
        bytes32 key = name.getKey();
        address oldValue = key.getAddress();

        require(oldValue != address(0), "Address does not exist");
        require(value != address(0), "Invalid address");
        require(oldValue != value, "Address already set");

        key.setAddress(value);

        emit AddressReplaced(name, key, oldValue, value);
    }

    /// @inheritdoc IEigenDAAddressDirectory
    function removeAddress(string memory name) external onlyOwner {
        bytes32 key = name.getKey();
        address existingAddress = key.getAddress();

        require(existingAddress != address(0), "Address does not exist");

        key.setAddress(address(0));

        emit AddressRemoved(name, key);
    }

    /// @inheritdoc IEigenDAAddressDirectory
    function getAddress(string memory name) external view returns (address) {
        return name.getKey().getAddress();
    }

    /// @inheritdoc IEigenDAAddressDirectory
    function getAddress(bytes32 key) external view returns (address) {
        return key.getAddress();
    }
}
