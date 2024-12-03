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

    bytes public quorumAdversaryThresholdPercentages;

    bytes public quorumConfirmationThresholdPercentages;

    bytes public quorumNumbersRequired;

    uint16 public nextBlobVersion;

    mapping(uint16 => VersionedBlobParams) public versionedBlobParams;

    SecurityThresholds public defaultSecurityThresholdsV2;

    uint256[44] private __GAP;
}