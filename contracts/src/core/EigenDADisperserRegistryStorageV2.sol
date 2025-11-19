// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";
import {
    EnumerableSetUpgradeable
} from "lib/openzeppelin-contracts-upgradeable/contracts/utils/structs/EnumerableSetUpgradeable.sol";

abstract contract EigenDADisperserRegistryStorageV2 {
    using EnumerableSetUpgradeable for EnumerableSetUpgradeable.UintSet;

    /// -----------------------------------------------------------------------
    /// Constants
    /// -----------------------------------------------------------------------

    /// @notice The EIP-712 typehash signed by a disperser that signals intent to deregister.
    /// @dev Owner has the ability to censor deregistrations by not processing the signature.
    bytes32 public constant DEREGISTRATION_TYPEHASH = keccak256("Deregister(uint32 disperserId,uint256 nonce)");

    /// @notice The EIP-712 typehash signed by a disperser that signals intent to update their relay URL.
    /// @dev Owner has the ability to censor relay URL updates by not processing the signature.
    bytes32 public constant UPDATE_RELAY_URL_TYPEHASH =
        keccak256("UpdateRelayURL(uint32 disperserId,string newRelayURL,uint256 nonce)");

    /// -----------------------------------------------------------------------
    /// Mutable Storage
    /// -----------------------------------------------------------------------

    /// @notice Returns the total number of registrations.
    uint32 public totalRegistrations;
    /// @notice Returns the nonce for a given disperser address.
    mapping(address disperser => uint256 nonce) public nonces;

    /// @dev Set of disperser IDs for default dispersers.
    /// Validators should default to accepting dispersals from dispersers in this set.
    EnumerableSetUpgradeable.UintSet internal _defaultDispersers;

    /// @dev Set of disperser IDs for on-demand dispersers.
    /// Dispersers in this set are authorized to use on-demand (pay-per-use) payments.
    EnumerableSetUpgradeable.UintSet internal _onDemandDispersers;
    /// @dev Mapping from disperser ID to disperser info.
    mapping(uint32 disperserId => EigenDATypesV2.DisperserInfoV2 disperserInfo) internal _disperserInfo;

    /// -----------------------------------------------------------------------
    /// Storage Gap
    /// -----------------------------------------------------------------------

    // slither-disable-next-line shadowing-state
    uint256[41] private __GAP;
}
