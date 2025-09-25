// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "lib/openzeppelin-contracts-upgradeable/contracts/access/OwnableUpgradeable.sol";
import {AddressDirectoryLib} from "src/core/libraries/v3/address-directory/AddressDirectoryLib.sol";
import {
    IEigenDADirectory,
    IEigenDAAddressDirectory,
    IEigenDAConfigRegistry
} from "src/core/interfaces/IEigenDADirectory.sol";
import {AccessControlConstants} from "src/core/libraries/v3/access-control/AccessControlConstants.sol";
import {AddressDirectoryConstants} from "src/core/libraries/v3/address-directory/AddressDirectoryConstants.sol";
import {IAccessControl} from "@openzeppelin/contracts/access/IAccessControl.sol";
import {InitializableLib} from "src/core/libraries/v3/initializable/InitializableLib.sol";
import {ConfigRegistryLib} from "src/core/libraries/v3/config-registry/ConfigRegistryLib.sol";

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

    /// ADDRESS DIRECTORY FUNCTIONS ///

    /// @inheritdoc IEigenDAAddressDirectory
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

    /// @inheritdoc IEigenDAAddressDirectory
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

    /// @inheritdoc IEigenDAAddressDirectory
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

    /// @inheritdoc IEigenDAAddressDirectory
    function getAddress(string memory name) external view returns (address) {
        return name.getKey().getAddress();
    }

    /// @inheritdoc IEigenDAAddressDirectory
    function getAddress(bytes32 key) external view returns (address) {
        return key.getAddress();
    }

    /// @inheritdoc IEigenDAAddressDirectory
    function getName(bytes32 key) external view returns (string memory) {
        return AddressDirectoryLib.getName(key);
    }

    /// @inheritdoc IEigenDAAddressDirectory
    function getAllNames() external view returns (string[] memory) {
        return AddressDirectoryLib.getNameList();
    }

    /// CONFIG REGISTRY FUNCTIONS ///

    /// @inheritdoc IEigenDAConfigRegistry
    function addConfigBytes32(string memory name, bytes32 value, string memory extraInfo) external onlyOwner {
        bytes32 key = ConfigRegistryLib.getKey(name);
        if (ConfigRegistryLib.isKeyRegisteredBytes32(key)) {
            revert ConfigAlreadyExists(name);
        }
        ConfigRegistryLib.setConfigBytes32(key, value, extraInfo);
        ConfigRegistryLib.registerKeyBytes32(name);
        emit ConfigBytes32Added(name, key, value, extraInfo);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function addConfigBytes(string memory name, bytes memory value, string memory extraInfo) external onlyOwner {
        bytes32 key = ConfigRegistryLib.getKey(name);
        if (ConfigRegistryLib.isKeyRegisteredBytes(key)) {
            revert ConfigAlreadyExists(name);
        }
        ConfigRegistryLib.setConfigBytes(key, value, extraInfo);
        ConfigRegistryLib.registerKeyBytes(name);
        emit ConfigBytesAdded(name, key, value, extraInfo);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function replaceConfigBytes32(string memory name, bytes32 value, string memory extraInfo) external onlyOwner {
        bytes32 key = ConfigRegistryLib.getKey(name);
        if (!ConfigRegistryLib.isKeyRegisteredBytes32(key)) {
            revert ConfigDoesNotExist(name);
        }
        bytes32 oldValue = ConfigRegistryLib.getConfigBytes32(key);
        ConfigRegistryLib.setConfigBytes32(key, value, extraInfo);
        emit ConfigBytes32Replaced(name, key, oldValue, value, extraInfo);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function replaceConfigBytes(string memory name, bytes memory value, string memory extraInfo) external onlyOwner {
        bytes32 key = ConfigRegistryLib.getKey(name);
        if (!ConfigRegistryLib.isKeyRegisteredBytes(key)) {
            revert ConfigDoesNotExist(name);
        }
        bytes memory oldValue = ConfigRegistryLib.getConfigBytes(key);
        ConfigRegistryLib.setConfigBytes(key, value, extraInfo);
        emit ConfigBytesReplaced(name, key, oldValue, value, extraInfo);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function removeConfigBytes32(string memory name) external onlyOwner {
        bytes32 key = ConfigRegistryLib.getKey(name);
        if (!ConfigRegistryLib.isKeyRegisteredBytes32(key)) {
            revert ConfigDoesNotExist(name);
        }
        ConfigRegistryLib.setConfigBytes32(key, bytes32(0), "");
        ConfigRegistryLib.deregisterKeyBytes32(name);
        emit ConfigBytes32Removed(name, key);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function removeConfigBytes(string memory name) external onlyOwner {
        bytes32 key = ConfigRegistryLib.getKey(name);
        if (!ConfigRegistryLib.isKeyRegisteredBytes(key)) {
            revert ConfigDoesNotExist(name);
        }
        ConfigRegistryLib.setConfigBytes(key, "", "");
        ConfigRegistryLib.deregisterKeyBytes(name);
        emit ConfigBytesRemoved(name, key);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigBytes32(string memory name) external view returns (bytes32) {
        bytes32 key = ConfigRegistryLib.getKey(name);
        return ConfigRegistryLib.getConfigBytes32(key);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigBytes32(bytes32 key) external view returns (bytes32) {
        return ConfigRegistryLib.getConfigBytes32(key);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigBytes32ExtraInfo(string memory name) external view returns (string memory) {
        bytes32 key = ConfigRegistryLib.getKey(name);
        return ConfigRegistryLib.getConfigBytes32ExtraInfo(key);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigBytes32ExtraInfo(bytes32 key) external view returns (string memory) {
        return ConfigRegistryLib.getConfigBytes32ExtraInfo(key);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigBytes(string memory name) external view returns (bytes memory) {
        bytes32 key = ConfigRegistryLib.getKey(name);
        return ConfigRegistryLib.getConfigBytes(key);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigBytes(bytes32 key) external view returns (bytes memory) {
        return ConfigRegistryLib.getConfigBytes(key);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigBytesExtraInfo(string memory name) external view returns (string memory) {
        bytes32 key = ConfigRegistryLib.getKey(name);
        return ConfigRegistryLib.getConfigBytesExtraInfo(key);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigBytesExtraInfo(bytes32 key) external view returns (string memory) {
        return ConfigRegistryLib.getConfigBytesExtraInfo(key);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getNumRegisteredKeysBytes32() external view returns (uint256) {
        return ConfigRegistryLib.getNumRegisteredKeysBytes32();
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getNumRegisteredKeysBytes() external view returns (uint256) {
        return ConfigRegistryLib.getNumRegisteredKeysBytes();
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getRegisteredKeyBytes32(uint256 index) external view returns (string memory) {
        return ConfigRegistryLib.getRegisteredKeyBytes32(index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getRegisteredKeyBytes(uint256 index) external view returns (string memory) {
        return ConfigRegistryLib.getRegisteredKeyBytes(index);
    }
}
