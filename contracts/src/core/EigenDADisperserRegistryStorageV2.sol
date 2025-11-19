// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";
import {EnumerableSet} from "lib/openzeppelin-contracts/contracts/utils/structs/EnumerableSet.sol";

abstract contract EigenDADisperserRegistryStorageV2 {
    using EnumerableSet for EnumerableSet.UintSet;

    /// -----------------------------------------------------------------------
    /// Constants
    /// -----------------------------------------------------------------------

    /// @notice The EIP-712 typehash signed by a disperser that signals intent to register.
    bytes32 public constant REGISTRATION_TYPEHASH =
        keccak256("Register(address disperser,string relayURL,uint256 nonce)");

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

    /// @dev Returns the ERC-7201 namespace for this contract.
    ///      https://eips.ethereum.org/EIPS/eip-7201
    bytes32 private constant DISPERSER_REGISTRY_STORAGE_LOCATION =
        keccak256(abi.encode(uint256(keccak256("eigenda.disperser.registry.storage")) - 1)) & ~bytes32(uint256(0xff));

    /// @dev Struct containing the storage layout for this contract (without dependencies).
    struct Layout {
        /// @dev Returns the total number of registrations.
        uint32 totalRegistrations;
        /// @dev Returns the nonce for a given disperser address.
        mapping(address disperser => uint256 nonce) nonces;
        /// @dev Set of disperser IDs for default dispersers.
        /// Validators should default to accepting dispersals from dispersers in this set.
        EnumerableSet.UintSet defaultDispersers;
        /// @dev Set of disperser IDs for on-demand dispersers.
        /// Dispersers in this set are authorized to use on-demand (pay-per-use) payments.
        EnumerableSet.UintSet onDemandDispersers;
        /// @dev Mapping from disperser ID to disperser info.
        mapping(uint32 disperserId => EigenDATypesV2.DisperserInfoV2 disperserInfo) disperserInfo;
    }

    /// @dev Returns the storage layout for this contract.
    /// Usage: `Layout storage $ = getDisperserRegistryStorage();`.
    function getDisperserRegistryStorage() internal pure returns (Layout storage $) {
        bytes32 ptr = DISPERSER_REGISTRY_STORAGE_LOCATION;
        /// @solidity memory-safe-assembly
        assembly {
            $.slot := ptr
        }
    }
}
