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

    modifier initializer() {
        InitializableLib.initialize();
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
    function initialize(address accessControl) external initializer {
        require(accessControl != address(0), "Access control address cannot be zero");
        bytes32 key = AddressDirectoryConstants.ACCESS_CONTROL_NAME.getKey();
        key.setAddress(accessControl);
        AddressDirectoryLib.registerKey(AddressDirectoryConstants.ACCESS_CONTROL_NAME);
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
    function addConfigBlockNumber(string memory name, uint256 abn, bytes memory value) external onlyOwner {
        bytes32 nameDigest = ConfigRegistryLib.getNameDigest(name);
        ConfigRegistryLib.addConfigBlockNumber(nameDigest, abn, value);
        ConfigRegistryLib.registerNameBlockNumber(name);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function addConfigTimeStamp(string memory name, uint256 activationTimeStamp, bytes memory value)
        external
        onlyOwner
    {
        bytes32 nameDigest = ConfigRegistryLib.getNameDigest(name);
        ConfigRegistryLib.addConfigTimeStamp(nameDigest, activationTimeStamp, value);
        ConfigRegistryLib.registerNameTimeStamp(name);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getNumCheckpointsBlockNumber(bytes32 nameDigest) external view returns (uint256) {
        return ConfigRegistryLib.getNumCheckpointsBlockNumber(nameDigest);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getNumCheckpointsTimeStamp(bytes32 nameDigest) external view returns (uint256) {
        return ConfigRegistryLib.getNumCheckpointsTimeStamp(nameDigest);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigBlockNumber(bytes32 nameDigest, uint256 index) external view returns (bytes memory) {
        return ConfigRegistryLib.getConfigBlockNumber(nameDigest, index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigTimeStamp(bytes32 nameDigest, uint256 index) external view returns (bytes memory) {
        return ConfigRegistryLib.getConfigTimeStamp(nameDigest, index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getActivationBlockNumber(bytes32 nameDigest, uint256 index) external view returns (uint256) {
        return ConfigRegistryLib.getActivationBlockNumber(nameDigest, index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getActivationTimeStamp(bytes32 nameDigest, uint256 index) external view returns (uint256) {
        return ConfigRegistryLib.getActivationTimeStamp(nameDigest, index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getCheckpointBlockNumber(bytes32 nameDigest, uint256 index)
        external
        view
        returns (ConfigRegistryTypes.BlockNumberCheckpoint memory)
    {
        return ConfigRegistryLib.getCheckpointBlockNumber(nameDigest, index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getCheckpointTimeStamp(bytes32 nameDigest, uint256 index)
        external
        view
        returns (ConfigRegistryTypes.TimeStampCheckpoint memory)
    {
        return ConfigRegistryLib.getCheckpointTimeStamp(nameDigest, index);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigNameBlockNumber(bytes32 nameDigest) external view returns (string memory) {
        return ConfigRegistryLib.getNameBlockNumber(nameDigest);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getConfigNameTimeStamp(bytes32 nameDigest) external view returns (string memory) {
        return ConfigRegistryLib.getNameTimeStamp(nameDigest);
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getAllConfigNamesBlockNumber() external view returns (string[] memory) {
        return ConfigRegistryLib.getNameListBlockNumber();
    }

    /// @inheritdoc IEigenDAConfigRegistry
    function getAllConfigNamesTimeStamp() external view returns (string[] memory) {
        return ConfigRegistryLib.getNameListTimeStamp();
    }


    /// @notice Retrieves the currently active block number config checkpoint and all future checkpoints for a given name.
    /// @dev Returns the checkpoint with the highest activation block that is less than or equal to the provided reference block,
    ///      plus all checkpoints with activation block numbers greater than the provided reference block.
    ///      This allows offchain clients to know the current configuration value and plan ahead for upcoming updates.
    function getActiveAndFutureBlockNumberConfigs(string memory name, uint256 referenceBlockNumber)
        external
        view
        returns (ConfigRegistryTypes.BlockNumberCheckpoint[] memory)
    {
        return ConfigRegistryLib.getActiveAndFutureBlockNumberConfigs(name, referenceBlockNumber);
    }

    /// @notice Retrieves the currently active timestamp config checkpoint and all future checkpoints for a given name.
    /// @dev Returns the checkpoint with the highest activation timestamp that is less than or equal to the provided reference timestamp,
    ///      plus all checkpoints with activation timestamps greater than the provided reference timestamp.
    ///      This allows offchain clients to know the current configuration value and plan ahead for upcoming updates.
    function getActiveAndFutureTimestampConfigs(string memory name, uint256 referenceTimestamp)
        external
        view
        returns (ConfigRegistryTypes.TimeStampCheckpoint[] memory)
    {
        return ConfigRegistryLib.getActiveAndFutureTimestampConfigs(name, referenceTimestamp);
    }

    /// @inheritdoc IEigenDASemVer
    function semver() external pure returns (uint8 major, uint8 minor, uint8 patch) {
        major = 2;
        minor = 0;
        patch = 0;
    }
}
