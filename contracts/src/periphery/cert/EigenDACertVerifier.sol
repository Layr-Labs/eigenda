// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifier} from "src/periphery/cert/interfaces/IEigenDACertVerifier.sol";
import {IEigenDACertVerifierBase} from "src/periphery/cert/interfaces/IEigenDACertVerifierBase.sol";

import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";

import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {EigenDATypesV2 as DATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";

import {EigenDACertVerificationLib as CertLib} from "src/periphery/cert/libraries/EigenDACertVerificationLib.sol";

contract EigenDACertVerifier is IEigenDACertVerifier {
    error InvalidSecurityThresholds();

    IEigenDAThresholdRegistry internal immutable _eigenDAThresholdRegistry;

    IEigenDASignatureVerifier internal immutable _eigenDASignatureVerifier;

    DATypesV1.SecurityThresholds internal _securityThresholds;

    bytes internal _quorumNumbersRequired;

    uint64 internal constant CERT_VERSION = 3;

    constructor(
        IEigenDAThresholdRegistry initEigenDAThresholdRegistry,
        IEigenDASignatureVerifier initEigenDASignatureVerifier,
        DATypesV1.SecurityThresholds memory initSecurityThresholds,
        bytes memory initQuorumNumbersRequired
    ) {
        if (initSecurityThresholds.confirmationThreshold <= initSecurityThresholds.adversaryThreshold) {
            revert InvalidSecurityThresholds();
        }
        _eigenDAThresholdRegistry = initEigenDAThresholdRegistry;
        _eigenDASignatureVerifier = initEigenDASignatureVerifier;
        _securityThresholds = initSecurityThresholds;
        _quorumNumbersRequired = initQuorumNumbersRequired;
    }

    /// @inheritdoc IEigenDACertVerifierBase
    function checkDACert(bytes calldata certBytes) external view returns (uint8) {
        (CertLib.StatusCode status,) = CertLib.checkDACert(
            _eigenDAThresholdRegistry, _eigenDASignatureVerifier, certBytes, _securityThresholds, _quorumNumbersRequired
        );
        return uint8(status);
    }

    /// @inheritdoc IEigenDACertVerifier
    function eigenDAThresholdRegistry() external view returns (IEigenDAThresholdRegistry) {
        return _eigenDAThresholdRegistry;
    }

    /// @inheritdoc IEigenDACertVerifier
    function eigenDASignatureVerifier() external view returns (IEigenDASignatureVerifier) {
        return _eigenDASignatureVerifier;
    }

    /// @inheritdoc IEigenDACertVerifier
    function securityThresholds() external view returns (DATypesV1.SecurityThresholds memory) {
        return _securityThresholds;
    }

    /// @inheritdoc IEigenDACertVerifier
    function quorumNumbersRequired() external view returns (bytes memory) {
        return _quorumNumbersRequired;
    }

    /// @inheritdoc IEigenDACertVerifier
    function certVersion() external pure returns (uint64) {
        return CERT_VERSION;
    }
}
