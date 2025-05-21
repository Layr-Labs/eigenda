// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {EigenDATypesV2 as DATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";

/// @notice This interface is a placeholder to accommodate upgrading the threshold registry without having to upgrade the service manager
///         which inherits this interface. Eventually this contract will be deprecated.
interface IEigenDAThresholdRegistryBase {
    event VersionedBlobParamsAdded(uint16 indexed version, DATypesV1.VersionedBlobParams versionedBlobParams);

    event QuorumAdversaryThresholdPercentagesUpdated(
        bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages
    );

    event QuorumConfirmationThresholdPercentagesUpdated(
        bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages
    );

    event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired);

    event DefaultSecurityThresholdsV2Updated(
        DATypesV1.SecurityThresholds previousDefaultSecurityThresholdsV2,
        DATypesV1.SecurityThresholds newDefaultSecurityThresholdsV2
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

interface IEigenDAThresholdRegistry is IEigenDAThresholdRegistryBase {
    event VersionedBlobParamsV2Added(uint256 indexed version, DATypesV2.VersionedBlobParams versionedBlobParams);

    function getBlobParamsV2(uint256 version) external view returns (DATypesV2.VersionedBlobParams memory);

    function addVersionedBlobParamsV2(DATypesV2.VersionedBlobParams memory newVersionedBlobParams) external;
}
