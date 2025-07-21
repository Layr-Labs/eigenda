// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {AccessControlConstants} from "src/core/libraries/v3/access-control/AccessControlConstants.sol";

/// @title EigenDAAccessControl
/// @notice This contract is to serve as the centralized source of truth for access control in all EigenDA contracts.
contract EigenDAAccessControl is AccessControl {
    constructor(address owner) {
        _grantRole(AccessControlConstants.OWNER_ROLE, owner);
    }

    function setupRole(bytes32 role, address account) external {
        require(hasRole(AccessControlConstants.OWNER_ROLE, msg.sender), "Caller is not the owner");
        _grantRole(role, account);
    }
}
