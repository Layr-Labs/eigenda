// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "../interfaces/IEigenDAStructs.sol";
import {EigenDATypesV2} from "../libraries/V2/EigenDATypesV2.sol";

/**
 * @title Storage variables for the `EigenDADisperserRegistry` contract.
 * @author Layr Labs, Inc.
 * @notice This storage contract is separate from the logic to simplify the upgrade process.
 */
abstract contract EigenDADisperserRegistryStorage {
    mapping(uint32 => EigenDATypesV2.DisperserInfo) public disperserKeyToInfo;

    // storage gap for upgradeability
    // slither-disable-next-line shadowing-state
    uint256[49] private __GAP;
}
