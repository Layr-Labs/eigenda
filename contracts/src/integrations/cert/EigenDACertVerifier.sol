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
import {EigenDACertTypes as CT} from "src/integrations/cert/EigenDACertTypes.sol";

contract EigenDACertVerifier is
    IEigenDACertVerifier,
    IEigenDACertVerifierBase,
    IVersionedEigenDACertVerifier,
    IEigenDASemVer
{
    error InvalidSecurityThresholds();

    IEigenDAThresholdRegistry internal immutable _eigenDAThresholdRegistry;

    IEigenDASignatureVerifier internal immutable _eigenDASignatureVerifier;

    DATypesV1.SecurityThresholds internal _securityThresholds;

    bytes internal _quorumNumbersRequired;

    uint8 internal constant MAJOR_VERSION = 3;
    uint8 internal constant MINOR_VERSION = 0;
    uint8 internal constant PATCH_VERSION = 0;

    /// @notice Status codes for certificate verification results
    enum StatusCode {
        NULL_ERROR, // Unused error code. If this is returned, there is a bug in the code.
        SUCCESS, // Verification succeeded
        INVALID_INCLUSION_PROOF, // Merkle inclusion proof is invalid
        SECURITY_ASSUMPTIONS_NOT_MET, // Security assumptions not met
        BLOB_QUORUMS_NOT_SUBSET, // Blob quorums not a subset of confirmed quorums
        REQUIRED_QUORUMS_NOT_SUBSET, // Required quorums not a subset of blob quorums
        INVALID_CERT // Certificate is invalid, due to some low level library revert having been caught
    }

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

    /// @notice Decodes a certificate from bytes to an EigenDACertV3
    /// @dev This function is external for the purpose of try/catch'ing it inside checkDACert.
    function decodeCert(bytes calldata data) external pure returns (CT.EigenDACertV3 memory cert) {
        return abi.decode(data, (CT.EigenDACertV3));
    }

    /// @inheritdoc IEigenDACertVerifierBase
    /// @dev This function try catches checkDACertReverts, and maps any reverts to status codes.
    /// TODO: we should return (uint8, bytes) instead and include the revert reason.
    function checkDACert(bytes calldata abiEncodedCert) external view returns (uint8) {
        CT.EigenDACertV3 memory daCert;
        // We try catch this here because decoding error would appear as a Panic,
        // which we consider bugs in the try/catch for the checkDACertReverts call below.
        try this.decodeCert(abiEncodedCert) returns (CT.EigenDACertV3 memory _daCert) {
            daCert = _daCert;
        } catch {
            return uint8(StatusCode.INVALID_CERT);
        }

        try this.checkDACertReverts(daCert) {
            return uint8(StatusCode.SUCCESS);
        } catch Error(string memory /*reason*/) {
            // This matches any require(..., "string reason") revert that is pre custom errors,
            // which many of our current eigenlayer-middleware dependencies like the BLSSignatureChecker still use. See:
            // https://github.com/Layr-Labs/eigenlayer-middleware/blob/fe5834371caed60c1d26ab62b5519b0cbdcb42fa/src/BLSSignatureChecker.sol#L96
            return uint8(StatusCode.INVALID_CERT);
        } catch Panic(uint errorCode) {
            // This matches any panic (e.g. arithmetic overflow, division by zero, invalid array access, etc.)
            // We pattern match these only to 
            revert(string(abi.encode("panic", errorCode)));
        } catch (bytes memory reason) {
            if (reason.length < 4) {
                // We re-throw any non custom-error that was caught here. For example,
                // low-level evm reverts such as out-of-gas don't return any data.
                // See https://rareskills.io/post/try-catch-solidity#gdvnie-9-what-gets-returned-during-an-out-of-gas?
                // These generally mean there is a bug in our implementation, which should be addressed by a human debugger.
                // TODO: figure out whether we can programmatically deal with out of gas, since that might happen from
                // a maliciously constructed cert.
                revert(string(reason));
            }
            // Any revert here is from custom errors coming from a failed require(..., SomeCustomError()) statement.
            // This mean that the cert is invalid.
            return uint8(StatusCode.INVALID_CERT);
        }
    }

    /// @notice Check a DA cert's validity
    /// @param daCert The EigenDA certificate
    /// @dev This function will revert if the certificate is invalid.
    function checkDACertReverts(CT.EigenDACertV3 calldata daCert) external view {
        CertLib.checkDACert(
            _eigenDAThresholdRegistry,
            _eigenDASignatureVerifier,
            daCert,
            _securityThresholds,
            _quorumNumbersRequired
        );
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
