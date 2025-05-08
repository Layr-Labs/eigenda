// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifier} from "src/periphery/cert/interfaces/IEigenDACertVerifier.sol";

import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";

import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {EigenDATypesV2 as DATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";

import {EigenDACertVerificationLib as CertLib} from "src/periphery/cert/libraries/EigenDACertVerificationLib.sol";

contract EigenDACertVerifier is IEigenDACertVerifier {
    error InvalidSecurityThresholds();

    IEigenDAThresholdRegistry public immutable eigenDAThresholdRegistry;

    IEigenDASignatureVerifier public immutable eigenDASignatureVerifier;

    DATypesV1.SecurityThresholds public securityThresholds;

    bytes public quorumNumbersRequired;

    uint64 internal constant CERT_VERSION = 3;

    constructor(
        IEigenDAThresholdRegistry _eigenDAThresholdRegistry,
        IEigenDASignatureVerifier _eigenDASignatureVerifier,
        DATypesV1.SecurityThresholds memory _securityThresholds,
        bytes memory _quorumNumbersRequired
    ) {
        if (_securityThresholds.confirmationThreshold <= _securityThresholds.adversaryThreshold) {
            revert InvalidSecurityThresholds();
        }
        eigenDAThresholdRegistry = _eigenDAThresholdRegistry;
        eigenDASignatureVerifier = _eigenDASignatureVerifier;
        securityThresholds = _securityThresholds;
        quorumNumbersRequired = _quorumNumbersRequired;
    }

    function checkDACert(bytes calldata certBytes) external view returns (uint8) {
        (CertLib.StatusCode status,) = CertLib.checkDACert(
            eigenDAThresholdRegistry, eigenDASignatureVerifier, certBytes, securityThresholds, quorumNumbersRequired
        );
        return uint8(status);
    }

    function certVersion() external pure returns (uint64) {
        return CERT_VERSION;
    }
}
