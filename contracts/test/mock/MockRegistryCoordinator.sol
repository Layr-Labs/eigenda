// SPDX-License-Identifier: BUSL-1.1

import {IStakeRegistry} from "lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {IBLSApkRegistry} from "lib/eigenlayer-middleware/src/interfaces/IBLSApkRegistry.sol";

pragma solidity =0.8.12;

// This mock is needed by the service manager contract's constructor
contract MockRegistryCoordinator {
    IStakeRegistry public immutable stakeRegistry;
    IBLSApkRegistry public immutable blsApkRegistry;

    constructor(IStakeRegistry _stakeRegistry, IBLSApkRegistry _blsApkRegistry) {
        stakeRegistry = _stakeRegistry;
        blsApkRegistry = _blsApkRegistry;
    }
}
