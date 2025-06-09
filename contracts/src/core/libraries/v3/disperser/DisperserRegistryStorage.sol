// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {DisperserRegistryTypes} from "src/core/libraries/v3/disperser/DisperserRegistryTypes.sol";

library DisperserRegistryStorage {
    string internal constant STORAGE_ID = "eigen.da.disperser.registry";
    bytes32 internal constant STORAGE_POSITION =
        keccak256(abi.encode(uint256(keccak256(abi.encodePacked(STORAGE_ID))) - 1)) & ~bytes32(uint256(0xff));

    struct Layout {
        mapping(uint32 => DisperserRegistryTypes.DisperserInfo) disperser;
        mapping(address => uint256) excess; // deposits + fees - refunds
        DisperserRegistryTypes.LockedDisperserDeposit depositParams;
        uint256 updateFee;
        uint32 nextDisperserKey;
    }

    function layout() internal pure returns (Layout storage s) {
        bytes32 position = STORAGE_POSITION;
        assembly {
            s.slot := position
        }
    }
}
