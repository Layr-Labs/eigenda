// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {AddressDirectoryStorage} from "src/core/libraries/v3/address-directory/AddressDirectoryStorage.sol";
import {IEigenDAAddressDirectory} from "src/core/interfaces/IEigenDAAddressDirectory.sol";

contract EigenDAAddressDirectory is OwnableUpgradeable, IEigenDAAddressDirectory {
    function initialize(address _initialOwner) external initializer {
        _transferOwnership(_initialOwner);
    }

    /// @inheritdoc IEigenDAAddressDirectory
    function setAddress(string memory name, address value) external onlyOwner {
        AddressDirectoryStorage.Layout storage s = AddressDirectoryStorage.layout();
        bytes32 key = keccak256(abi.encodePacked(name));
        s.addresses[key] = value;
        emit AddressSet(name, value);
    }

    /// @inheritdoc IEigenDAAddressDirectory
    function getAddress(string memory name) external view returns (address) {
        AddressDirectoryStorage.Layout storage s = AddressDirectoryStorage.layout();
        bytes32 key = keccak256(abi.encodePacked(name));
        return s.addresses[key];
    }

    /// @inheritdoc IEigenDAAddressDirectory
    function getAddress(bytes32 key) external view returns (address) {
        AddressDirectoryStorage.Layout storage s = AddressDirectoryStorage.layout();
        return s.addresses[key];
    }
}
