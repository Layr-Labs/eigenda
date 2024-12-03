// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "../interfaces/IEigenDAStructs.sol";

interface IEigenDAThresholdRegistry {

    event VersionedBlobParamsAdded(uint16 indexed version, VersionedBlobParams versionedBlobParams);

    ///////////////////////// V1 ///////////////////////////////

    /// @notice Returns an array of bytes where each byte represents the adversary threshold percentage of the quorum at that index
    function quorumAdversaryThresholdPercentages() external view returns (bytes memory);

    /// @notice Returns an array of bytes where each byte represents the confirmation threshold percentage of the quorum at that index
    function quorumConfirmationThresholdPercentages() external view returns (bytes memory);

    /// @notice Returns an array of bytes where each byte represents the number of a required quorum 
    function quorumNumbersRequired() external view returns (bytes memory);

    /// @notice Gets the adversary threshold percentage for a quorum
    function getQuorumAdversaryThresholdPercentage(
        uint8 quorumNumber
    ) external view returns (uint8);

    /// @notice Gets the confirmation threshold percentage for a quorum
    function getQuorumConfirmationThresholdPercentage(
        uint8 quorumNumber
    ) external view returns (uint8);

    /// @notice Checks if a quorum is required
    function getIsQuorumRequired(
        uint8 quorumNumber
    ) external view returns (bool);

    ///////////////////////// V2 ///////////////////////////////

    /// @notice Gets the default security thresholds for V2
    function getDefaultSecurityThresholdsV2() external view returns (SecurityThresholds memory);

    /// @notice Returns the blob params for a given blob version
    function getBlobParams(uint16 version) external view returns (VersionedBlobParams memory);
}