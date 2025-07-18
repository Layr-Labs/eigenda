// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

library AccessControlConstants {
    bytes32 internal constant OWNER_ROLE = keccak256("OWNER");
    bytes32 internal constant QUORUM_OWNER_SEED = keccak256("QUORUM_OWNER");

    function QUORUM_OWNER_ROLE(uint64 quorumId) internal pure returns (bytes32) {
        return bytes32(uint256(QUORUM_OWNER_SEED) + quorumId);
    }
}
