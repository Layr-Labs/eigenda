// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";
import {EnumerableSet} from "lib/openzeppelin-contracts/contracts/utils/structs/EnumerableSet.sol";

abstract contract EigenDADisperserRegistryStorageV2 {
    using EnumerableSet for EnumerableSet.UintSet;

    /// -----------------------------------------------------------------------
    /// Constants
    /// -----------------------------------------------------------------------

    /// @notice The EIP-712 typehash signed by a disperser that signals intent to deregister.
    /// @dev Owner has the ability to censor deregistrations by not processing the signature.
    bytes32 public constant DEREGISTRATION_TYPEHASH = 
        keccak256("Deregister(uint32 disperserId)");

    /// @notice The EIP-712 typehash signed by a disperser that signals intent to update their relay URL.
    /// @dev Owner has the ability to censor relay URL updates by not processing the signature.
    bytes32 public constant UPDATE_RELAY_URL_TYPEHASH =
        keccak256("UpdateRelayURL(uint32 disperserId,string newRelayURL)");

    /// -----------------------------------------------------------------------
    /// Mutable Storage
    /// -----------------------------------------------------------------------

    /// @notice Given `disperserId`, returns the disperser's info.
    /// @dev Returns an empty struct for non-existent disperser IDs.
    mapping(uint32 disperserId => EigenDATypesV2.DisperserInfoV2 disperserInfo) public disperserIdToInfo;

    /// @notice Mapping from disperser address to disperser ID
    mapping(address => uint32) public disperserAddressToId;

    /// @notice Counter for the next disperser ID
    uint32 public nextDisperserId;

    /// @notice Set of default disperser IDs
    EnumerableSet.UintSet internal defaultDispersersSet;

    /// @notice Set of on-demand disperser IDs
    EnumerableSet.UintSet internal onDemandDispersersSet;

    /// -----------------------------------------------------------------------
    /// Storage Gap
    /// -----------------------------------------------------------------------

    // slither-disable-next-line shadowing-state
    uint256[44] private __GAP;
}
