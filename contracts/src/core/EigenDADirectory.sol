// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {AddressDirectoryLib} from "src/core/libraries/v3/address-directory/AddressDirectoryLib.sol";
import {IEigenDADirectory} from "src/core/interfaces/IEigenDADirectory.sol";
import {AccessControlConstants} from "src/core/libraries/v3/access-control/AccessControlConstants.sol";
import {AddressDirectoryConstants} from "src/core/libraries/v3/address-directory/AddressDirectoryConstants.sol";
import {IAccessControl} from "@openzeppelin/contracts/access/IAccessControl.sol";
import {InitializableLib} from "src/core/libraries/v3/initializable/InitializableLib.sol";

contract EigenDADirectory is IEigenDADirectory {
    using AddressDirectoryLib for string;
    using AddressDirectoryLib for bytes32;

    modifier initializer() {
        InitializableLib.initialize();
        _;
    }

    modifier onlyOwner() {
        require(
            IAccessControl(AddressDirectoryConstants.ACCESS_CONTROL_NAME.getKey().getAddress()).hasRole(
                AccessControlConstants.OWNER_ROLE, msg.sender
            ),
            "Caller is not the owner"
        );
        _;
    }

    /// @dev If doing a fresh deployment, this contract should be deployed AFTER an access control contract has been deployed.
    function initialize(address accessControl) external initializer {
        require(accessControl != address(0), "Access control address cannot be zero");
        bytes32 key = AddressDirectoryConstants.ACCESS_CONTROL_NAME.getKey();
        AddressDirectoryConstants.ACCESS_CONTROL_NAME.getKey().setAddress(accessControl);
        emit AddressAdded(AddressDirectoryConstants.ACCESS_CONTROL_NAME, key, accessControl);
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
        AddressDirectoryLib.registerKey(name);

        emit AddressAdded(name, key, value);
    }

    /// @inheritdoc IEigenDADirectory
    function replaceAddress(string memory name, address value) external onlyOwner {
        bytes32 key = name.getKey();
        address oldValue = key.getAddress();

        if (oldValue == address(0)) {
            revert AddressDoesNotExist(name);
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
            revert AddressDoesNotExist(name);
        }

        key.setAddress(address(0));
        AddressDirectoryLib.deregisterKey(name);

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

    function getName(bytes32 key) external view returns (string memory) {
        return AddressDirectoryLib.getName(key);
    }

    function getAllNames() external view returns (string[] memory) {
        return AddressDirectoryLib.getNameList();
    }
}
