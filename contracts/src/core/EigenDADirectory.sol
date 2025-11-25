// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

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
import {IEigenDASemVer} from "src/core/interfaces/IEigenDASemVer.sol";

contract EigenDADirectory is IEigenDADirectory, IEigenDASemVer {
    using AddressDirectoryLib for string;
    using AddressDirectoryLib for bytes32;

    modifier reinitializer(uint8 version) {
        InitializableLib.reinitialize(version);
        _;
    }

    modifier onlyOwner() {
        require(
            IAccessControl(AddressDirectoryConstants.ACCESS_CONTROL_NAME.getKey().getAddress())
                .hasRole(AccessControlConstants.OWNER_ROLE, msg.sender),
            "Caller is not the owner"
        );
        _;
    }

    /// @dev If doing a fresh deployment, this contract should be deployed AFTER an access control contract has been deployed.
    function initialize(address accessControl) external reinitializer(2) {
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
    function getAddress(bytes32 nameDigest) external view returns (address) {
        return nameDigest.getAddress();
    }

    /// @inheritdoc IEigenDAAddressDirectory
    function getName(bytes32 nameDigest) external view returns (string memory) {
        return AddressDirectoryLib.getName(nameDigest);
    }

    /// @inheritdoc IEigenDAAddressDirectory
    function getAllNames() external view returns (string[] memory) {
        return AddressDirectoryLib.getNameList();
    }

    /// CONFIG REGISTRY FUNCTIONS ///

    /// @inheritdoc IEigenDAConfigRegistry
    function addConfigBytes32(string memory name, uint256 activationKey, bytes32 value) external onlyOwner {
        bytes32 nameDigest = ConfigRegistryLib.getNameDigest(name);
        ConfigRegistryLib.addConfigBytes32(nameDigest, activationKey, value);
        ConfigRegistryLib.registerNameBytes32(name);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function addConfigBytes(string memory name, uint256 activationKey, bytes memory value) external onlyOwner {
        bytes32 nameDigest = ConfigRegistryLib.getNameDigest(name);
        ConfigRegistryLib.addConfigBytes(nameDigest, activationKey, value);
        ConfigRegistryLib.registerNameBytes(name);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getNumCheckpointsBytes32(bytes32 nameDigest) external view returns (uint256) {
        return ConfigRegistryLib.getNumCheckpointsBytes32(nameDigest);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getNumCheckpointsBytes(bytes32 nameDigest) external view returns (uint256) {
        return ConfigRegistryLib.getNumCheckpointsBytes(nameDigest);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigBytes32(bytes32 nameDigest, uint256 index) external view returns (bytes32) {
        return ConfigRegistryLib.getConfigBytes32(nameDigest, index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigBytes(bytes32 nameDigest, uint256 index) external view returns (bytes memory) {
        return ConfigRegistryLib.getConfigBytes(nameDigest, index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getActivationKeyBytes32(bytes32 nameDigest, uint256 index) external view returns (uint256) {
        return ConfigRegistryLib.getActivationKeyBytes32(nameDigest, index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getActivationKeyBytes(bytes32 nameDigest, uint256 index) external view returns (uint256) {
        return ConfigRegistryLib.getActivationKeyBytes(nameDigest, index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getCheckpointBytes32(bytes32 nameDigest, uint256 index)
        external
        view
        returns (ConfigRegistryTypes.Bytes32Checkpoint memory)
    {
        return ConfigRegistryLib.getCheckpointBytes32(nameDigest, index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getCheckpointBytes(bytes32 nameDigest, uint256 index)
        external
        view
        returns (ConfigRegistryTypes.BytesCheckpoint memory)
    {
        return ConfigRegistryLib.getCheckpointBytes(nameDigest, index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigNameBytes32(bytes32 nameDigest) external view returns (string memory) {
        return ConfigRegistryLib.getNameBytes32(nameDigest);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigNameBytes(bytes32 nameDigest) external view returns (string memory) {
        return ConfigRegistryLib.getNameBytes(nameDigest);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getAllConfigNamesBytes32() external view returns (string[] memory) {
        return ConfigRegistryLib.getNameListBytes32();
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getAllConfigNamesBytes() external view returns (string[] memory) {
        return ConfigRegistryLib.getNameListBytes();
    }

    /// @inheritdoc IEigenDASemVer
    function semver() external pure returns (uint8 major, uint8 minor, uint8 patch) {
        major = 1;
        minor = 1;
        patch = 0;
    }
}
