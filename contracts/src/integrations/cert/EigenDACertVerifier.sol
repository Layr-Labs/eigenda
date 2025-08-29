// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifier, IEigenDACertVerifierBase, IVersionedEigenDACertVerifier} from "src/integrations/cert/interfaces/IEigenDACertVerifier.sol";

import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";

import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {EigenDATypesV2 as DATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";

import {IEigenDASemVer} from "src/core/interfaces/IEigenDASemVer.sol";

import {EigenDACertVerificationLib as CertLib} from "src/integrations/cert/libraries/EigenDACertVerificationLib.sol";
import {EigenDACertTypes as CT} from "src/integrations/cert/EigenDACertTypes.sol";

/// @title EigenDACertVerifier
/// @notice Verifies EigenDA certificates
/// @dev This contract's checkDACert function is designed to be zk provable by risczero's Steel library,
/// which does not support zk proving reverting calls: https://github.com/risc0/risc0-ethereum/issues/438.
/// For this reason, we avoid using revert statements and instead return error codes.
/// The only acceptable reverts are from bugs, such as contract misconfiguration.
/// The goal is to be able to perfectly classify into 3 categories: valid certs, invalid certs, and bugs.
/// Bugs could be configuration mistakes (_eigenDAThresholdRegistry points to 0x0), or other unexpected logic issues.
contract EigenDACertVerifier is
    IEigenDACertVerifier,
    IEigenDACertVerifierBase,
    IVersionedEigenDACertVerifier,
    IEigenDASemVer
{
    error InvalidSecurityThresholds();
    error InvalidQuorumNumbersRequired(uint256 length);

    IEigenDAThresholdRegistry internal immutable _eigenDAThresholdRegistry;

    IEigenDASignatureVerifier internal immutable _eigenDASignatureVerifier;

    DATypesV1.SecurityThresholds internal _securityThresholds;

    bytes internal _quorumNumbersRequired;

    uint8 internal constant MAJOR_VERSION = 3;
    uint8 internal constant MINOR_VERSION = 1;
    uint8 internal constant PATCH_VERSION = 0;

    constructor(
        IEigenDAThresholdRegistry initEigenDAThresholdRegistry,
        IEigenDASignatureVerifier initEigenDASignatureVerifier,
        DATypesV1.SecurityThresholds memory initSecurityThresholds,
        bytes memory initQuorumNumbersRequired
    ) {
        if (
            initSecurityThresholds.confirmationThreshold <=
            initSecurityThresholds.adversaryThreshold
        ) {
            revert InvalidSecurityThresholds();
        }
        if (initQuorumNumbersRequired.length == 0 || initQuorumNumbersRequired.length > 256) {
            revert InvalidQuorumNumbersRequired(initQuorumNumbersRequired.length);
        }
        _eigenDAThresholdRegistry = initEigenDAThresholdRegistry;
        _eigenDASignatureVerifier = initEigenDASignatureVerifier;
        _securityThresholds = initSecurityThresholds;
        _quorumNumbersRequired = initQuorumNumbersRequired;
    }

    /// @notice Decodes a certificate from bytes to an EigenDACertV3
    /// @dev This function is external for the purpose of try/catch'ing it inside checkDACert.
    function decodeCert(
        bytes calldata data
    ) external pure returns (CT.EigenDACertV3 memory cert) {
        return abi.decode(data, (CT.EigenDACertV3));
    }

    /// @inheritdoc IEigenDACertVerifierBase
    function checkDACert(
        bytes calldata abiEncodedCert
    ) external view returns (uint8) {
        CT.EigenDACertV3 memory daCert;

        try this.decodeCert(abiEncodedCert) returns (
            CT.EigenDACertV3 memory cert
        ) {
            daCert = cert;
        } catch {
            return uint8(CertLib.StatusCode.CERT_DECODE_REVERT);
        }

        (CertLib.StatusCode status, ) = CertLib.checkDACert(
            _eigenDAThresholdRegistry,
            _eigenDASignatureVerifier,
            daCert,
            _securityThresholds,
            _quorumNumbersRequired
        );
        return uint8(status);
    }

    /// @inheritdoc IEigenDACertVerifier
    function eigenDAThresholdRegistry()
        external
        view
        returns (IEigenDAThresholdRegistry)
    {
        return _eigenDAThresholdRegistry;
    }

    /// @inheritdoc IEigenDACertVerifier
    function eigenDASignatureVerifier()
        external
        view
        returns (IEigenDASignatureVerifier)
    {
        return _eigenDASignatureVerifier;
    }

    /// @inheritdoc IEigenDACertVerifier
    function securityThresholds()
        external
        view
        returns (DATypesV1.SecurityThresholds memory)
    {
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
    function semver()
        external
        pure
        returns (uint8 major, uint8 minor, uint8 patch)
    {
        major = MAJOR_VERSION;
        minor = MINOR_VERSION;
        patch = PATCH_VERSION;
    }
}
