// SPDX-License-Identifier: BUSL-1.1

import {IDelegationManager} from
    "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IDelegationManager.sol";

pragma solidity =0.8.12;

// This mock is needed by the service manager contract's constructor
contract MockStakeRegistry {
    IDelegationManager public immutable delegation;

    constructor(IDelegationManager delegationManager) {
        delegation = delegationManager;
    }
}
