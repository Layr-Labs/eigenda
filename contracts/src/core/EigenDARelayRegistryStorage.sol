// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";

/**
 * @title Storage variables for the `EigenDARelayRegistry` contract.
 * @author Layr Labs, Inc.
 * @notice This storage contract is separate from the logic to simplify the upgrade process.
 */
abstract contract EigenDARelayRegistryStorage {
    mapping(uint32 => EigenDATypesV2.RelayInfo) public relayKeyToInfo;

    uint32 public nextRelayKey;

    // storage gap for upgradeability
    // slither-disable-next-line shadowing-state
    uint256[48] private __GAP;
}
