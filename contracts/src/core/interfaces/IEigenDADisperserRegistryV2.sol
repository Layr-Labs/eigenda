// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";

/// @title IEigenDADisperserRegistryV2
/// @author Layr Labs, Inc.
/// @notice This interface defines the contract for managing EigenDA disperser registration and configuration.
/// @dev This registry manages disperser registration, deregistration, and relay URL updates.
interface IEigenDADisperserRegistryV2 {
    /// -----------------------------------------------------------------------
    /// Events
    /// -----------------------------------------------------------------------
    /// @notice This event is emitted when a new disperser is successfully registered with the registry.
    /// @param disperserId The unique ID that has been assigned to the newly registered disperser.
    /// @param disperser The address of the disperser that has been registered.
    event DisperserRegistered(uint32 indexed disperserId, address indexed disperser);

    /// @notice This event is emitted when a disperser is successfully deregistered from the registry.
    /// @param disperserId The ID of the disperser that has been deregistered.
    event DisperserDeregistered(uint32 indexed disperserId);

    /// @notice This event is emitted when a disperser successfully updates their relay URL.
    /// @param disperserId The ID of the disperser that updated their relay URL.
    /// @param relayURL The new relay URL that has been set for the disperser.
    event RelayURLUpdated(uint32 indexed disperserId, string relayURL);

    /// @notice This event is emitted when a disperser is added to the default dispersers set.
    /// @param disperserId The ID of the disperser that has been added to the default set.
    event DefaultDisperserAdded(uint32 indexed disperserId);

    /// @notice This event is emitted when a disperser is added to the on-demand dispersers set.
    /// @param disperserId The ID of the disperser that has been added to the on-demand set.
    event OnDemandDisperserAdded(uint32 indexed disperserId);

    /// @notice This event is emitted when a disperser is removed from the default dispersers set.
    /// @param disperserId The ID of the disperser that has been removed from the default set.
    event DefaultDisperserRemoved(uint32 indexed disperserId);

    /// @notice This event is emitted when a disperser is removed from the on-demand dispersers set.
    /// @param disperserId The ID of the disperser that has been removed from the on-demand set.
    event OnDemandDisperserRemoved(uint32 indexed disperserId);

    /// -----------------------------------------------------------------------
    /// Errors
    /// -----------------------------------------------------------------------

    /// @notice Thrown when a zero address is provided as input.
    error InputAddressZero();
    /// @notice Thrown when an invalid signature is provided for verification.
    error InvalidSignature();
    /// @notice Thrown when attempting to register a disperser that is already registered.
    error DisperserIsRegistered();
    /// @notice Thrown when attempting to deregister a disperser that is not registered.
    error DisperserIsNotRegistered();
    /// @notice Thrown when attempting to add an existing disperser to the specified set.
    error DisperserInSet();
    /// @notice Thrown when attempting to remove a disperser that is not in the specified set.
    error DisperserNotInSet();

    /// -----------------------------------------------------------------------
    /// External Functions
    /// -----------------------------------------------------------------------

    /// @notice This function registers a new disperser with the registry and assigns it a unique ID.
    /// @param disperser The address of the disperser that should be registered.
    /// @param relayURL The relay URL that will be associated with this disperser.
    /// @return disperserId The unique ID that has been assigned to the newly registered disperser.
    function registerDisperser(address disperser, string memory relayURL) external returns (uint32 disperserId);

    /// @notice This function deregisters a disperser from the registry using a valid signature.
    /// @dev The `signature` parameter can be empty if `msg.sender` is the registered disperser.
    /// Reverts if the signature is invalid and the caller is not the disperser.
    /// @param disperserId The ID of the disperser that should be deregistered.
    /// @param signature The signature from the disperser that authorizes this deregistration.
    function deregisterDisperser(uint32 disperserId, bytes memory signature) external;

    /// @notice This function updates the relay URL for a disperser using a valid signature.
    /// @dev The `signature` parameter can be empty if `msg.sender` is the registered disperser.
    /// Reverts if the signature is invalid and the caller is not the disperser.
    /// @param disperserId The ID of the disperser whose relay URL should be updated.
    /// @param relayURL The new relay URL that should be set for this disperser.
    /// @param signature The signature from the disperser that authorizes this update.
    function updateRelayURL(uint32 disperserId, string memory relayURL, bytes memory signature) external;

    /// @notice This function revokes the nonce for the caller.
    /// @dev This function is used to revoke the caller's current nonce.
    function revokeNonce() external;

    /// -----------------------------------------------------------------------
    /// Owner-only Functions
    /// -----------------------------------------------------------------------

    /// @notice This function adds a disperser to the default dispersers set (only callable by owner).
    /// @param disperserId The ID of the disperser that should be added to the default set.
    function addDefaultDisperser(uint32 disperserId) external;

    /// @notice This function adds a disperser to the on-demand dispersers set (only callable by owner).
    /// @param disperserId The ID of the disperser that should be added to the on-demand set.
    function addOnDemandDisperser(uint32 disperserId) external;

    /// @notice This function removes a disperser from the default dispersers set (only callable by owner).
    /// @param disperserId The ID of the disperser that should be removed from the default set.
    function removeDefaultDisperser(uint32 disperserId) external;

    /// @notice This function removes a disperser from the on-demand dispersers set (only callable by owner).
    /// @param disperserId The ID of the disperser that should be removed from the on-demand set.
    function removeOnDemandDisperser(uint32 disperserId) external;

    /// -----------------------------------------------------------------------
    /// View Functions
    /// -----------------------------------------------------------------------

    /// @notice This function returns the disperser information for the provided array of disperser IDs.
    /// @param ids The array of disperser IDs for which information should be retrieved.
    /// @return An array of DisperserInfoV2 structs containing the information for each requested disperser.
    function getDisperserInfo(uint32[] memory ids) external view returns (EigenDATypesV2.DisperserInfoV2[] memory);

    /// @notice This function returns all disperser IDs that are in the default dispersers set.
    /// @return An array containing all disperser IDs in the default dispersers set.
    function getDefaultDisperserIds() external view returns (uint32[] memory);

    /// @notice This function returns all disperser IDs that are in the on-demand dispersers set.
    /// @return An array containing all disperser IDs in the on-demand dispersers set.
    function getOnDemandDisperserIds() external view returns (uint32[] memory);

    /// @notice This function checks whether a disperser is in the default dispersers set.
    /// @param disperserId The ID of the disperser to check.
    /// @return A boolean value indicating whether the disperser is in the default set (true) or not (false).
    function isDefaultDisperserId(uint32 disperserId) external view returns (bool);

    /// @notice This function checks whether a disperser is in the on-demand dispersers set.
    /// @param disperserId The ID of the disperser to check.
    /// @return A boolean value indicating whether the disperser is in the on-demand set (true) or not (false).
    function isOnDemandDisperserId(uint32 disperserId) external view returns (bool);
}

