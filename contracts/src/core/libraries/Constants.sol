// SPDX-License-Identifier: MIT
pragma solidity ^0.8.12;

library Constants {
    /// ROLES

    bytes32 internal constant OWNER_ROLE = keccak256("eigen.da.owner");
    bytes32 internal constant QUORUM_OWNER_SEED = keccak256("eigen.da.quorum.owner");

    function QUORUM_OWNER_ROLE(uint64 quorumId) internal pure returns (bytes32) {
        return bytes32(uint256(QUORUM_OWNER_SEED) + quorumId);
    }

    /// ADDRESS DIRECTORY

    bytes32 internal constant EIGEN_DA_CERT_VERIFIER_ROUTER = keccak256("eigen.da.blob.params.registry");
}
