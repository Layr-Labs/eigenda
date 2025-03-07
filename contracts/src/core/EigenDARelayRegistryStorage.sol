// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "../interfaces/IEigenDAStructs.sol";

/**
 * @title EigenDARelayRegistryStorage
 * @notice This storage contract is separated from the logic to simplify the upgrade process.
 */
abstract contract EigenDARelayRegistryStorage {

    /// @notice A mapping of relay keys to relay info
    mapping(uint32 => RelayInfo) public relayKeyToInfo;

    /// @notice The next relay key to be used
    uint32 public nextRelayKey;

    /// @notice Storage gap for upgradeability
    // slither-disable-next-line shadowing-state
    uint256[48] private __GAP;
}