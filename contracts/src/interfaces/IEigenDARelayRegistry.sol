// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "./IEigenDAStructs.sol";

/**
 * @title IEigenDARelayRegistry
 * @notice A registry for EigenDA relay info
 */
interface IEigenDARelayRegistry {

    /// @notice Emitted when a relay is added to the registry
    event RelayAdded(address indexed relay, uint32 indexed key, string relayURL);

    /**
     * @notice Appends a relay info to the registry and returns the relay key
     * @param relayInfo The relay info to add
     */
    function addRelayInfo(RelayInfo memory relayInfo) external returns (uint32);

    /**
     * @notice Returns the relay address for a given relay key
     * @param key The key of the relay to get the address for
     */
    function relayKeyToAddress(uint32 key) external view returns (address);

    /**
     * @notice Returns the relay URL for a given relay key
     * @param key The key of the relay to get the URL for
     */
    function relayKeyToUrl(uint32 key) external view returns (string memory);
}
