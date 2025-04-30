// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "../interfaces/IEigenDAStructs.sol";
import {EigenDATypesV1 as DATypesV1} from "../libraries/V1/EigenDATypesV1.sol";
import {EigenDATypesV2 as DATypesV2} from "../libraries/V2/EigenDATypesV2.sol";

interface IEigenDAThresholdRegistry {
    event VersionedBlobParamsAdded(uint16 indexed version, DATypesV1.VersionedBlobParams versionedBlobParams);

    event QuorumAdversaryThresholdPercentagesUpdated(
        bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages
    );

    event QuorumConfirmationThresholdPercentagesUpdated(
        bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages
    );

    event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired);

    event DefaultSecurityThresholdsV2Updated(
        DATypesV1.SecurityThresholds previousDefaultSecurityThresholdsV2, DATypesV1.SecurityThresholds newDefaultSecurityThresholdsV2
    );

    ///////////////////////// V1 ///////////////////////////////

    /// @notice Returns an array of bytes where each byte represents the adversary threshold percentage of the quorum at that index
    function quorumAdversaryThresholdPercentages() external view returns (bytes memory);

    /// @notice Returns an array of bytes where each byte represents the confirmation threshold percentage of the quorum at that index
    function quorumConfirmationThresholdPercentages() external view returns (bytes memory);

    /// @notice Returns an array of bytes where each byte represents the number of a required quorum
    function quorumNumbersRequired() external view returns (bytes memory);

    /// @notice Gets the adversary threshold percentage for a quorum
    function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) external view returns (uint8);

    /// @notice Gets the confirmation threshold percentage for a quorum
    function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) external view returns (uint8);

    /// @notice Checks if a quorum is required
    function getIsQuorumRequired(uint8 quorumNumber) external view returns (bool);

    ///////////////////////// V2 ///////////////////////////////

    /// @notice Returns the blob params for a given blob version
    function getBlobParams(uint16 version) external view returns (DATypesV1.VersionedBlobParams memory);
}
