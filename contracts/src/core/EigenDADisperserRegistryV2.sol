// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "lib/openzeppelin-contracts-upgradeable/contracts/access/OwnableUpgradeable.sol";
import {
    EnumerableMapUpgradeable,
    EnumerableSetUpgradeable
} from "lib/openzeppelin-contracts-upgradeable/contracts/utils/structs/EnumerableMapUpgradeable.sol";
import {
    EIP712Upgradeable
} from "lib/openzeppelin-contracts-upgradeable/contracts/utils/cryptography/draft-EIP712Upgradeable.sol";
import {
    SignatureCheckerUpgradeable
} from "lib/openzeppelin-contracts-upgradeable/contracts/utils/cryptography/SignatureCheckerUpgradeable.sol";

import {EigenDATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";
import {IEigenDADisperserRegistryV2} from "src/core/interfaces/IEigenDADisperserRegistryV2.sol";
import {EigenDADisperserRegistryStorageV2} from "src/core/EigenDADisperserRegistryStorageV2.sol";

/// @title Registry for EigenDA disperser info V2
/// @author Layr Labs, Inc.
contract EigenDADisperserRegistryV2 is
    OwnableUpgradeable,
    EIP712Upgradeable,
    EigenDADisperserRegistryStorageV2,
    IEigenDADisperserRegistryV2
{
    using SignatureCheckerUpgradeable for address;
    using EnumerableMapUpgradeable for *;
    using EnumerableSetUpgradeable for *;

    /// -----------------------------------------------------------------------
    /// Initialization
    /// -----------------------------------------------------------------------

    constructor() {
        _disableInitializers();
    }

    function initialize(address _initialOwner) external initializer {
        __Ownable_init();
        __EIP712_init("EigenDADisperserRegistry", "2");
        _transferOwnership(_initialOwner);
    }

    /// -----------------------------------------------------------------------
    /// External Logic
    /// -----------------------------------------------------------------------

    /// @inheritdoc IEigenDADisperserRegistryV2
    function registerDisperser(address disperser, string memory relayURL)
        external
        virtual
        returns (uint32 disperserId)
    {
        // Increment and assign the next available disperserId, starting at 1 (not 0).
        disperserId = ++totalDispersers;

        EigenDATypesV2.DisperserInfoV2 storage disperserInfo = _disperserInfo[disperserId];

        // Assert that the disperser address is non-zero.
        if (disperser == address(0)) revert InputAddressZero();
        // Assert that the disperser is not already registered.
        if (disperserInfo.disperser != address(0)) revert AlreadyRegistered();

        // Set the disperser info.
        disperserInfo.disperser = disperser;
        disperserInfo.relayURL = relayURL;

        emit DisperserRegistered(disperserId, disperser);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function deregisterDisperser(uint32 disperserId, bytes memory signature) external virtual {
        EigenDATypesV2.DisperserInfoV2 storage disperserInfo = _disperserInfo[disperserId];

        bytes32 digest = keccak256(abi.encode(DEREGISTRATION_TYPEHASH, disperserId));

        // Assert that the disperser is not in the default or on-demand dispersers sets.
        // This means owner can only deregister dispersers after revoking their default or on-demand status.
        if (_defaultDispersers.contains(disperserId) || _onDemandDispersers.contains(disperserId)) {
            revert DisperserInSet();
        }
        // Assert that the signature is valid (supports EIP-1271).
        if (!disperserInfo.disperser.isValidSignatureNow(digest, signature)) revert InvalidSignature();
        // Assert that the disperser is registered.
        if (disperserInfo.disperser == address(0)) revert NotRegistered();

        // Delete the disperser info.
        delete disperserInfo.disperser;
        delete disperserInfo.relayURL;

        emit DisperserDeregistered(disperserId);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function updateRelayURL(uint32 disperserId, string memory relayURL, bytes memory signature) external {
        EigenDATypesV2.DisperserInfoV2 storage disperserInfo = _disperserInfo[disperserId];

        // TODO: Skip verification if the disperser is msg.sender.

        bytes32 digest = keccak256(abi.encode(UPDATE_RELAY_URL_TYPEHASH, disperserId, relayURL));

        // Assert that the signature is valid (supports EIP-1271).
        if (!disperserInfo.disperser.isValidSignatureNow(digest, signature)) revert InvalidSignature();
        // Assert that the disperser is registered.
        if (disperserInfo.disperser == address(0)) revert NotRegistered();

        // Set the relay URL.
        disperserInfo.relayURL = relayURL;

        emit RelayURLUpdated(disperserId, relayURL);
    }

    /// -----------------------------------------------------------------------
    /// Owner-only Logic
    /// -----------------------------------------------------------------------

    // TODO: Extra checks

    /// @inheritdoc IEigenDADisperserRegistryV2
    function addDefaultDisperser(address disperser, uint32 disperserId) external onlyOwner {
        if (!_defaultDispersers.set(disperserId, disperser)) revert DisperserAlreadyExists();
        emit DefaultDisperserAdded(disperserId);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function addOnDemandDisperser(address disperser, uint32 disperserId) external onlyOwner {
        if (!_onDemandDispersers.set(disperserId, disperser)) revert DisperserAlreadyExists();
        emit OnDemandDisperserAdded(disperserId);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function removeDefaultDisperser(uint32 disperserId) external onlyOwner {
        if (!_defaultDispersers.remove(disperserId)) revert DisperserNotFound();
        emit DefaultDisperserRemoved(disperserId);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function removeOnDemandDisperser(uint32 disperserId) external onlyOwner {
        if (!_onDemandDispersers.remove(disperserId)) revert DisperserNotFound();
        emit OnDemandDisperserRemoved(disperserId);
    }

    /// -----------------------------------------------------------------------
    /// View Logic
    /// -----------------------------------------------------------------------

    /// @inheritdoc IEigenDADisperserRegistryV2
    function getDisperserInfo(uint32[] memory ids) external view returns (EigenDATypesV2.DisperserInfoV2[] memory) {}

    // /// @inheritdoc IEigenDADisperserRegistryV2
    // function disperserKeyToAddress(uint32 key) external view returns (address) {}

    // /// @inheritdoc IEigenDADisperserRegistryV2
    // function getDisperserIdByAddress(address disperser) external view returns (uint32) {
    //     return _disperserInfo[disperser].disperserId;
    // }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function getDefaultDisperserIds() external view returns (uint32[] memory ids) {
        bytes32[] memory keys = _defaultDispersers._inner._keys.values();
        /// @solidity memory-safe-assembly
        assembly {
            ids := keys
        }
        return ids;
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function getOnDemandDisperserIds() external view returns (uint32[] memory ids) {
        bytes32[] memory keys = _onDemandDispersers._inner._keys.values();
        /// @solidity memory-safe-assembly
        assembly {
            ids := keys
        }
        return ids;
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function isDefaultDisperserId(uint32 disperserId) external view returns (bool) {
        return _defaultDispersers.contains(disperserId);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function isOnDemandDisperserId(uint32 disperserId) external view returns (bool) {
        return _onDemandDispersers.contains(disperserId);
    }
}
