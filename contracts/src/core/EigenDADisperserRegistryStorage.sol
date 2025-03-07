// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "../interfaces/IEigenDAStructs.sol";

/**
 * @title EigenDADisperserRegistryStorage
 * @notice This storage contract is separated from the logic to simplify the upgrade process.
 */
abstract contract EigenDADisperserRegistryStorage {

    /// @notice A mapping of disperser keys to disperser info
    mapping(uint32 => DisperserInfo) public disperserKeyToInfo;

    /// @notice Storage gap for upgradeability
    // slither-disable-next-line shadowing-state
    uint256[49] private __GAP;
}