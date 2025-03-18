// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "./IEigenDAStructs.sol";

/**
 * @title IEigenDADisperserRegistry
 * @notice A registry for EigenDA disperser info
 */
interface IEigenDADisperserRegistry {

    /// @notice Emitted when a disperser is added to the registry
    event DisperserAdded(uint32 indexed key, address indexed disperser);

    /**
     * @notice Sets the disperser info for a given disperser key
     * @param _disperserKey The key of the disperser to set the info for
     * @param _disperserInfo The info to set for the disperser
     */
    function setDisperserInfo(uint32 _disperserKey, DisperserInfo memory _disperserInfo) external;

    /**
     * @notice Returns the disperser address for a given disperser key
     * @param key The key of the disperser to get the address for
     */
    function disperserKeyToAddress(uint32 key) external view returns (address);
}