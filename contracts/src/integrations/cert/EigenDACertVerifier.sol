// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {
    IEigenDACertVerifier,
    IEigenDACertVerifierBase,
    IVersionedEigenDACertVerifier
} from "src/integrations/cert/interfaces/IEigenDACertVerifier.sol";

import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";

import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {EigenDATypesV2 as DATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";

import {IEigenDASemVer} from "src/core/interfaces/IEigenDASemVer.sol";

import {EigenDACertVerificationLib as CertLib} from "src/integrations/cert/libraries/EigenDACertVerificationLib.sol";

contract EigenDACertVerifier is
    IEigenDACertVerifier,
    IEigenDACertVerifierBase,
    IVersionedEigenDACertVerifier,
    IEigenDASemVer
{
    error InvalidSecurityThresholds();

    IEigenDAThresholdRegistry internal immutable _eigenDAThresholdRegistry;

    IEigenDASignatureVerifier internal immutable _eigenDASignatureVerifier;

    /// @notice Default security thresholds used by {checkDACert}.
    /// @dev Checked inside {EigenDACertVerificationLib-checkDACert}. Changing this default
    ///      at deployment affects verification outcomes. Constraints to respect:
    ///      - confirmationThreshold > adversaryThreshold (constructor-enforced)
    ///      - confirmationThreshold - adversaryThreshold > reconstructionThreshold
    ///        (see src/integrations/cert/libraries/EigenDACertVerificationLib.sol)
    DATypesV1.SecurityThresholds internal _securityThresholds;

    bytes internal _quorumNumbersRequired;

    uint8 internal constant MAJOR_VERSION = 3;
    uint8 internal constant MINOR_VERSION = 0;
    uint8 internal constant PATCH_VERSION = 0;

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
    function checkDACert(bytes calldata abiEncodedCert) external view returns (uint8) {
        (CertLib.StatusCode status,) = CertLib.checkDACert(
            _eigenDAThresholdRegistry,
            _eigenDASignatureVerifier,
            abiEncodedCert,
            _securityThresholds,
            _quorumNumbersRequired
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

    /// @inheritdoc IVersionedEigenDACertVerifier
    function certVersion() external pure returns (uint8) {
        return MAJOR_VERSION;
    }

    /// @inheritdoc IEigenDASemVer
    function semver() external pure returns (uint8 major, uint8 minor, uint8 patch) {
        major = MAJOR_VERSION;
        minor = MINOR_VERSION;
        patch = PATCH_VERSION;
    }
}
