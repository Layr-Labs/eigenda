// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {AccessControlConstants} from "src/core/libraries/v3/access-control/AccessControlConstants.sol";

/// @title EigenDAAccessControl
/// @notice This contract is to serve as the centralized source of truth for access control in all EigenDA contracts.
contract EigenDAAccessControl is AccessControl {
    constructor(address owner) {
        // The DEFAULT_ADMIN_ROLE can set the admin role for all other roles, and should be put behind a timelock.
        _grantRole(DEFAULT_ADMIN_ROLE, owner);
        // The OWNER_ROLE is the default ownership role for EigenDA contracts.
        _grantRole(AccessControlConstants.OWNER_ROLE, owner);
    }
}
