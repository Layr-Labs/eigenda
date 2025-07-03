// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {UsageAuthorizationTypes} from "src/core/libraries/v3/usage-authorization/UsageAuthorizationTypes.sol";

library UsageAuthorizationStorage {
    string internal constant STORAGE_ID = "eigen.da.usage.authorization.storage";
    bytes32 internal constant STORAGE_POSITION =
        keccak256(abi.encode(uint256(keccak256(abi.encodePacked(STORAGE_ID))) - 1)) & ~bytes32(uint256(0xff));

    struct Layout {
        mapping(uint64 => UsageAuthorizationTypes.Quorum) quorum;
    }

    function layout() internal pure returns (Layout storage s) {
        bytes32 position = STORAGE_POSITION;
        assembly {
            s.slot := position
        }
    }
}
