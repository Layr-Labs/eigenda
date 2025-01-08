// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAThresholdRegistry} from "../interfaces/IEigenDAThresholdRegistry.sol";
import "../interfaces/IEigenDAStructs.sol";

/**
 * @title Storage variables for the `EigenDAThresholdRegistry` contract.
 * @author Layr Labs, Inc.
 * @notice This storage contract is separate from the logic to simplify the upgrade process.
 */
abstract contract EigenDAThresholdRegistryStorage is IEigenDAThresholdRegistry {

    /// @notice The adversary threshold percentage for the quorum at position `quorumNumber`
    bytes public quorumAdversaryThresholdPercentages;

    /// @notice The confirmation threshold percentage for the quorum at position `quorumNumber`
    bytes public quorumConfirmationThresholdPercentages;

    /// @notice The set of quorum numbers that are required
    bytes public quorumNumbersRequired;

    /// @notice The next blob version id to be added
    uint16 public nextBlobVersion;

    /// @notice mapping of blob version id to the params of the blob version
    mapping(uint16 => VersionedBlobParams) public versionedBlobParams;

    /// @notice Default security thresholds for quorums
    SecurityThresholds public defaultSecurityThresholdsV2;

    // storage gap for upgradeability
    // slither-disable-next-line shadowing-state
    uint256[44] private __GAP;
}