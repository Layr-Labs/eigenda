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
import {ConfigRegistryTypes} from "src/core/libraries/v3/config-registry/ConfigRegistryTypes.sol";

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
    function addConfigBytes32(string memory name, uint256 activationKey, bytes32 value) external onlyOwner {
        bytes32 key = ConfigRegistryLib.getKey(name);
        if (ConfigRegistryLib.isKeyRegisteredBytes32(key)) {
            revert ConfigAlreadyExists(name);
        }
        ConfigRegistryLib.addConfigBytes32(key, activationKey, value);
        ConfigRegistryLib.registerKeyBytes32(name);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function addConfigBytes(string memory name, uint256 activationKey, bytes memory value) external onlyOwner {
        bytes32 key = ConfigRegistryLib.getKey(name);
        if (ConfigRegistryLib.isKeyRegisteredBytes(key)) {
            revert ConfigAlreadyExists(name);
        }
        ConfigRegistryLib.addConfigBytes(key, activationKey, value);
        ConfigRegistryLib.registerKeyBytes(name);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getNumCheckpointsBytes32(bytes32 key) external view returns (uint256) {
        return ConfigRegistryLib.getNumCheckpointsBytes32(key);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getNumCheckpointsBytes(bytes32 key) external view returns (uint256) {
        return ConfigRegistryLib.getNumCheckpointsBytes(key);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigBytes32(bytes32 key, uint256 index) external view returns (bytes32) {
        return ConfigRegistryLib.getConfigBytes32(key, index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigBytes(bytes32 key, uint256 index) external view returns (bytes memory) {
        return ConfigRegistryLib.getConfigBytes(key, index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getActivationKeyBytes32(bytes32 key, uint256 index) external view returns (uint256) {
        return ConfigRegistryLib.getActivationKeyBytes32(key, index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getActivationKeyBytes(bytes32 key, uint256 index) external view returns (uint256) {
        return ConfigRegistryLib.getActivationKeyBytes(key, index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getCheckpointBytes32(bytes32 key, uint256 index)
        external
        view
        returns (ConfigRegistryTypes.Bytes32Checkpoint memory)
    {
        return ConfigRegistryLib.getCheckpointBytes32(key, index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getCheckpointBytes(bytes32 key, uint256 index)
        external
        view
        returns (ConfigRegistryTypes.BytesCheckpoint memory)
    {
        return ConfigRegistryLib.getCheckpointBytes(key, index);
    }
}
