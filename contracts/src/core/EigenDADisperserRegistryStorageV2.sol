// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";
import {
    EnumerableMapUpgradeable
} from "lib/openzeppelin-contracts-upgradeable/contracts/utils/structs/EnumerableMapUpgradeable.sol";

abstract contract EigenDADisperserRegistryStorageV2 {
    using EnumerableMapUpgradeable for EnumerableMapUpgradeable.UintToAddressMap;

    /// -----------------------------------------------------------------------
    /// Constants
    /// -----------------------------------------------------------------------

    /// @notice The EIP-712 typehash signed by a disperser that signals intent to deregister.
    /// @dev Owner has the ability to censor deregistrations by not processing the signature.
    bytes32 public constant DEREGISTRATION_TYPEHASH = keccak256("Deregister(uint32 disperserId)");

    /// @notice The EIP-712 typehash signed by a disperser that signals intent to update their relay URL.
    /// @dev Owner has the ability to censor relay URL updates by not processing the signature.
    bytes32 public constant UPDATE_RELAY_URL_TYPEHASH =
        keccak256("UpdateRelayURL(uint32 disperserId,string newRelayURL)");

    /// -----------------------------------------------------------------------
    /// Mutable Storage
    /// -----------------------------------------------------------------------

    /// @notice Returns the total number of registered dispersers.
    uint32 public totalDispersers;
    /// @dev Mapping from disperser ID to disperser address for default dispersers.
    EnumerableMapUpgradeable.UintToAddressMap internal _defaultDispersers;
    /// @dev Mapping from disperser ID to disperser address for on-demand dispersers.
    EnumerableMapUpgradeable.UintToAddressMap internal _onDemandDispersers;
    /// @dev Mapping from disperser ID to disperser info.
    mapping(uint32 disperserId => EigenDATypesV2.DisperserInfoV2 disperserInfo) public _disperserInfo;

    /// -----------------------------------------------------------------------
    /// Storage Gap
    /// -----------------------------------------------------------------------

    // slither-disable-next-line shadowing-state
    uint256[44] private __GAP; // TODO: Update gap to accounts for enumerable maps.
}
