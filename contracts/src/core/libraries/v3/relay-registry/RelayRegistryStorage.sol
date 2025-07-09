// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {RelayRegistryLib} from "src/core/libraries/v3/relay-registry/RelayRegistryLib.sol";

/// @notice Defines the storage layout for a relay registry based on ERC-7201
///         https://eips.ethereum.org/EIPS/eip-7201
library RelayRegistryStorage {
    /// @custom: storage-location erc7201:relay.registry.storage
    struct Layout {
        mapping(uint32 => RelayRegistryLib.RelayInfo) relay;
        uint32 nextRelayId; // Used to track the next available relay ID
    }

    string internal constant STORAGE_ID = "relay.registry.storage";
    bytes32 internal constant STORAGE_POSITION =
        keccak256(abi.encode(uint256(keccak256(abi.encodePacked(STORAGE_ID))) - 1)) & ~bytes32(uint256(0xff));

    function layout() internal pure returns (Layout storage s) {
        bytes32 position = STORAGE_POSITION;
        assembly {
            s.slot := position
        }
    }
}
