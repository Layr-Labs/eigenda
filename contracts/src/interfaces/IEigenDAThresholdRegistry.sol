// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "../interfaces/IEigenDAStructs.sol";

interface IEigenDAThresholdRegistry {

    /// @notice Emitted when a new blob version is added to the registry
    event VersionedBlobParamsAdded(uint16 indexed version, VersionedBlobParams versionedBlobParams);

    ///////////////////////// V1 ///////////////////////////////

    /// @notice Returns an array of bytes where each byte represents the adversary threshold percentage of the quorum at that index
    function quorumAdversaryThresholdPercentages() external view returns (bytes memory);

    /// @notice Returns an array of bytes where each byte represents the confirmation threshold percentage of the quorum at that index
    function quorumConfirmationThresholdPercentages() external view returns (bytes memory);

    /// @notice Returns an array of bytes where each byte represents the number of a required quorum 
    function quorumNumbersRequired() external view returns (bytes memory);

    /// @notice Returns the adversary threshold percentage for a quorum for V1 verification
    /// @param quorumNumber The number of the quorum to get the adversary threshold percentage for
    function getQuorumAdversaryThresholdPercentage(
        uint8 quorumNumber
    ) external view returns (uint8);

    /// @notice Returns the confirmation threshold percentage for a quorum for V1 verification
    /// @param quorumNumber The number of the quorum to get the confirmation threshold percentage for
    function getQuorumConfirmationThresholdPercentage(
        uint8 quorumNumber
    ) external view returns (uint8);

    /// @notice Returns true if a quorum is required for V1 verification
    /// @param quorumNumber The number of the quorum to check if it is required for V1 verification
    function getIsQuorumRequired(
        uint8 quorumNumber
    ) external view returns (bool);

    ///////////////////////// V2 ///////////////////////////////

    /// @notice Returns the blob params for a given blob version
    /// @param version The version of the blob to get the params for
    function getBlobParams(uint16 version) external view returns (VersionedBlobParams memory);
}