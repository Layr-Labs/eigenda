// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

/// @notice This library defines constants for access control to use in solidity contracts. Off-chain users should derive the same constants defined here.
library AccessControlConstants {
    /// @notice This role manages all other roles, and is all powerful.
    bytes32 internal constant OWNER_ROLE = keccak256("OWNER");

    /// @notice This is the seed used to derive the quorum owner role for each quorum.
    bytes32 internal constant QUORUM_OWNER_SEED = keccak256("QUORUM_OWNER");

    /// @dev We simply add the quorum ID to the seed to derive a unique role for each quorum.
    function QUORUM_OWNER_ROLE(uint64 quorumId) internal pure returns (bytes32) {
        return bytes32(uint256(QUORUM_OWNER_SEED) + quorumId);
    }
}
