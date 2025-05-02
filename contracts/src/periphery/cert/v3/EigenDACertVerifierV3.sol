// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifier} from "src/periphery/cert/interfaces/IEigenDACertVerifier.sol";

import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";

import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {EigenDATypesV2 as DATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";

import {EigenDACertVerificationV2Lib as CertV2Lib} from "src/periphery/cert/v2/EigenDACertVerificationV2Lib.sol";
import {EigenDACertVerificationV3Lib as CertV3Lib} from "src/periphery/cert/v3/EigenDACertVerificationV3Lib.sol";

contract EigenDACertVerifierV3 is IEigenDACertVerifier {
    error InvalidSecurityThresholds();

    IEigenDAThresholdRegistry public immutable eigenDAThresholdRegistry;

    IEigenDASignatureVerifier public immutable eigenDASignatureVerifier;

    DATypesV1.SecurityThresholds public securityThresholds;

    bytes public quorumNumbersRequired;

    uint8 internal constant CERT_VERSION = 3;

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
        (CertV2Lib.StatusCode status,) = CertV3Lib.checkDACert(
            eigenDAThresholdRegistry, eigenDASignatureVerifier, certBytes, securityThresholds, quorumNumbersRequired
        );
        return uint8(status);
    }

    function certVersion() external pure returns (uint8) {
        return CERT_VERSION;
    }
}
