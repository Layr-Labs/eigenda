// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

// TODO: We MUST use OZ v5.5.0 in order to use EIP712Upgradeable.

import {OwnableUpgradeable} from "lib/openzeppelin-contracts-upgradeable/contracts/access/OwnableUpgradeable.sol";
import {EnumerableSet} from "lib/openzeppelin-contracts/contracts/utils/structs/EnumerableSet.sol";
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
/// @notice This contract manages disperser registration and authorization for the EigenDA network.
/// @dev Supports permissionless registration with owner-managed authorization sets:
///      - Default dispersers: Validators should default to accepting dispersals from these dispersers
///      - On-demand dispersers: Authorized to use on-demand (pay-per-use) payments
contract EigenDADisperserRegistryV2 is
    OwnableUpgradeable,
    EIP712Upgradeable,
    EigenDADisperserRegistryStorageV2,
    IEigenDADisperserRegistryV2
{
    using SignatureCheckerUpgradeable for address;
    using EnumerableSet for *;

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
    function registerDisperser(address disperser, string memory relayURL, bytes memory signature)
        external
        virtual
        returns (uint32 disperserId)
    {
        // Create a reference to the contract's storage.
        Layout storage $ = getDisperserRegistryStorage();
        // Create a reference to the disperser info.
        EigenDATypesV2.DisperserInfoV2 storage disperserInfo = $.disperserInfo[disperserId];

        // Increment and assign the next available disperserId, starting at 1 (not 0).
        disperserId = $.totalRegistrations + 1; // Monotonically increasing.

        bytes32 digest = keccak256(abi.encode(REGISTRATION_TYPEHASH, disperser, relayURL, $.nonces[disperser]++));

        // Assert that the signature is valid (supports EIP-1271).
        _checkSignature(disperser, digest, signature);
        // Assert that the disperser address is non-zero.
        if (disperser == address(0)) revert InputAddressZero();
        // Assert that the disperser is not already registered.
        if (disperserInfo.disperser != address(0)) revert DisperserIsRegistered();

        // Increment total registrations.
        $.totalRegistrations = disperserId;
        // Set the disperser info.
        disperserInfo.disperser = disperser;
        disperserInfo.relayURL = relayURL;

        emit DisperserRegistered(disperserId, disperser);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function deregisterDisperser(uint32 disperserId, bytes memory signature) external virtual {
        // Create a reference to the contract's storage.
        Layout storage $ = getDisperserRegistryStorage();
        // Create a reference to the disperser info.
        EigenDATypesV2.DisperserInfoV2 storage disperserInfo = $.disperserInfo[disperserId];

        bytes32 digest =
            keccak256(abi.encode(DEREGISTRATION_TYPEHASH, disperserId, $.nonces[disperserInfo.disperser]++));

        // TODO: Instead we should simply remove them from the sets.

        // Assert that the disperser is not in the default or on-demand dispersers sets.
        // This means owner can only deregister dispersers after revoking their default or on-demand status.
        if ($.defaultDispersers.contains(disperserId) || $.onDemandDispersers.contains(disperserId)) {
            revert DisperserInSet();
        }

        _checkSignature(disperserInfo.disperser, digest, signature);
        // Assert that the disperser is registered.
        if (disperserInfo.disperser == address(0)) revert DisperserIsNotRegistered();

        // Total registrations are monotonically increasing, no need to decrement.

        // Delete the disperser info.
        delete disperserInfo.disperser;
        delete disperserInfo.relayURL;

        emit DisperserDeregistered(disperserId);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function updateRelayURL(uint32 disperserId, string memory relayURL, bytes memory signature) external {
        Layout storage $ = getDisperserRegistryStorage();

        EigenDATypesV2.DisperserInfoV2 storage disperserInfo = $.disperserInfo[disperserId];

        bytes32 digest = keccak256(
            abi.encode(UPDATE_RELAY_URL_TYPEHASH, disperserId, relayURL, $.nonces[disperserInfo.disperser]++)
        );

        // Assert that the signature is valid (supports EIP-1271).
        _checkSignature(disperserInfo.disperser, digest, signature);
        // Assert that the disperser is registered.
        if (disperserInfo.disperser == address(0)) revert DisperserIsNotRegistered();

        // Set the relay URL.
        disperserInfo.relayURL = relayURL;

        emit RelayURLUpdated(disperserId, relayURL);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function revokeNonce() external {
        Layout storage $ = getDisperserRegistryStorage();
        unchecked {
            ++$.nonces[msg.sender];
        }
    }

    /// @dev Verifies that the signature is valid if caller is not the disperser.
    /// @dev Reverts if the signature is invalid and the caller is not the disperser.
    /// @param disperser The address of the disperser signaling an intent.
    /// @param digest The digest of the update relay URL transaction.
    /// @param signature The signature of the disperser signaling an intent.
    function _checkSignature(address disperser, bytes32 digest, bytes memory signature) internal view {
        if (msg.sender != disperser) {
            if (!disperser.isValidSignatureNow(_hashTypedDataV4(digest), signature)) revert InvalidSignature();
        }
    }

    /// -----------------------------------------------------------------------
    /// Owner-only Logic
    /// -----------------------------------------------------------------------

    /// @inheritdoc IEigenDADisperserRegistryV2
    function addDefaultDisperser(uint32 disperserId) external onlyOwner {
        Layout storage $ = getDisperserRegistryStorage();
        if ($.disperserInfo[disperserId].disperser == address(0)) revert DisperserIsNotRegistered();
        if (!$.defaultDispersers.add(disperserId)) revert DisperserInSet();
        emit DefaultDisperserAdded(disperserId);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function addOnDemandDisperser(uint32 disperserId) external onlyOwner {
        Layout storage $ = getDisperserRegistryStorage();
        if ($.disperserInfo[disperserId].disperser == address(0)) revert DisperserIsNotRegistered();
        if (!$.onDemandDispersers.add(disperserId)) revert DisperserInSet();
        emit OnDemandDisperserAdded(disperserId);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function removeDefaultDisperser(uint32 disperserId) external onlyOwner {
        Layout storage $ = getDisperserRegistryStorage();
        if (!$.defaultDispersers.remove(disperserId)) revert DisperserNotInSet();
        emit DefaultDisperserRemoved(disperserId);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function removeOnDemandDisperser(uint32 disperserId) external onlyOwner {
        Layout storage $ = getDisperserRegistryStorage();
        if (!$.onDemandDispersers.remove(disperserId)) revert DisperserNotInSet();
        emit OnDemandDisperserRemoved(disperserId);
    }

    /// -----------------------------------------------------------------------
    /// View Logic
    /// -----------------------------------------------------------------------

    // /// @inheritdoc IEigenDADisperserRegistryV2
    // function getTotalRegistrations() external view returns (uint32) {
    //     Layout storage $ = getDisperserRegistryStorage();
    //     return $.totalRegistrations;
    // }

    // /// @inheritdoc IEigenDADisperserRegistryV2
    // function getNonce(address disperser) external view returns (uint256) {
    //     Layout storage $ = getDisperserRegistryStorage();
    //     return $.nonces[disperser];
    // }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function getDisperserInfo(uint32[] memory ids)
        external
        view
        returns (EigenDATypesV2.DisperserInfoV2[] memory info)
    {
        Layout storage $ = getDisperserRegistryStorage();
        uint256 len = ids.length;
        info = new EigenDATypesV2.DisperserInfoV2[](len);
        for (uint256 i = 0; i < len; ++i) {
            info[i] = $.disperserInfo[ids[i]];
        }
        return info;
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function getDefaultDisperserIds() external view returns (uint32[] memory ids) {
        Layout storage $ = getDisperserRegistryStorage();
        uint256[] memory keys = $.defaultDispersers.values();
        /// @solidity memory-safe-assembly
        assembly {
            ids := keys
        }
        return ids;
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function getOnDemandDisperserIds() external view returns (uint32[] memory ids) {
        Layout storage $ = getDisperserRegistryStorage();
        uint256[] memory keys = $.onDemandDispersers.values();
        /// @solidity memory-safe-assembly
        assembly {
            ids := keys
        }
        return ids;
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function isDefaultDisperserId(uint32 disperserId) external view returns (bool) {
        Layout storage $ = getDisperserRegistryStorage();
        return $.defaultDispersers.contains(disperserId);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function isOnDemandDisperserId(uint32 disperserId) external view returns (bool) {
        Layout storage $ = getDisperserRegistryStorage();
        return $.onDemandDispersers.contains(disperserId);
    }
}
