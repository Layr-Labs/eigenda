// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {ConfigRegistryTypes as T} from "src/core/libraries/v3/config-registry/ConfigRegistryTypes.sol";

/// @notice Defines the storage layout for a config registry based on ERC-7201
///         https://eips.ethereum.org/EIPS/eip-7201
library ConfigRegistryStorage {
    ///@custom: storage-location erc7201:config.registry.storage
    struct Layout {
        T.BlockNumberCfg BlockNumberCfg;
        T.TimeStampCfg TimeStampCfg;
    }
    
    /// v3 suffix is appended to migrate away from legacy layout that used
    /// bytes32 and bytes mapping types
    string internal constant STORAGE_ID = "config.registry.storage-v2";
    bytes32 internal constant STORAGE_POSITION =
        keccak256(abi.encode(uint256(keccak256(abi.encodePacked(STORAGE_ID))) - 1)) & ~bytes32(uint256(0xff));

    function layout() internal pure returns (Layout storage s) {
        bytes32 position = STORAGE_POSITION;
        assembly {
            s.slot := position
        }
    }
}
