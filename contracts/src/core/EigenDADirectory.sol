// SPDX-License-Identifier: MIT
pragma solidity ^0.8.12;

import {AddressDirectoryLib} from "src/core/libraries/AddressDirectoryLib.sol";
import {AccessControlLib} from "src/core/libraries/AccessControlLib.sol";
import {InitializableLib} from "src/core/libraries/InitializableLib.sol";
import {Constants} from "src/core/libraries/Constants.sol";

contract EigenDADirectory {
    modifier onlyOwner() {
        _onlyOwner();
        _;
    }

    modifier onlyQuorumOwner(uint64 quorumId) {
        _onlyQuorumOwner(quorumId);
        _;
    }

    function initialize(address owner) external {
        InitializableLib.setInitializedVersion(1);
        AccessControlLib.grantRole(Constants.OWNER_ROLE, owner);
    }

    /// SETTERS

    function addOwner(address owner) external onlyOwner {
        AccessControlLib.grantRole(Constants.OWNER_ROLE, owner);
    }

    function removeOwner(address owner) external onlyOwner {
        AccessControlLib.revokeRole(Constants.OWNER_ROLE, owner);
    }

    function addQuorumOwner(uint64 quorumId, address owner) external onlyOwner {
        AccessControlLib.grantRole(Constants.QUORUM_OWNER_ROLE(quorumId), owner);
    }

    function removeQuorumOwner(uint64 quorumId, address owner) external onlyOwner {
        AccessControlLib.revokeRole(Constants.QUORUM_OWNER_ROLE(quorumId), owner);
    }

    function setAddress(bytes32 key, address value) external onlyOwner {
        AddressDirectoryLib.setAddress(key, value);
    }

    /// GETTERS

    function getAddress(bytes32 key) external view returns (address) {
        return AddressDirectoryLib.getAddress(key);
    }

    function getRoleMember(bytes32 role, uint256 index) external view returns (address) {
        return AccessControlLib.getRoleMember(role, index);
    }

    function getRoleMemberCount(bytes32 role) external view returns (uint256) {
        return AccessControlLib.getRoleMemberCount(role);
    }

    /// INTERNAL

    function _onlyOwner() internal view virtual {
        require(AccessControlLib.hasRole(Constants.OWNER_ROLE, msg.sender), "Not owner");
    }

    function _onlyQuorumOwner(uint64 quorumId) internal view virtual {
        require(AccessControlLib.hasRole(Constants.QUORUM_OWNER_ROLE(quorumId), msg.sender), "Not quorum owner");
    }
}
