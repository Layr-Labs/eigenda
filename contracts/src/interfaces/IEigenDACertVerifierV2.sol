// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAThresholdRegistry} from "./IEigenDAThresholdRegistry.sol";
import "./IEigenDAStructs.sol";

interface IEigenDACertVerifierV2 {
    function verifyDACertV2(EigenDACertV2 calldata cert) external view;
}
