// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAThresholdRegistry} from "../interfaces/IEigenDAThresholdRegistry.sol";
import "../interfaces/IEigenDAStructs.sol";

abstract contract EigenDAThresholdRegistryStorage is IEigenDAThresholdRegistry {

    bytes public quorumAdversaryThresholdPercentages;

    bytes public quorumConfirmationThresholdPercentages;

    bytes public quorumNumbersRequired;

    mapping(uint16 => VersionedBlobParams) public versionedBlobParams;

    SecurityThresholds public defaultSecurityThresholdsV2;

    uint256[45] private __GAP;
}