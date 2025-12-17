// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifier} from "src/integrations/cert/interfaces/IEigenDACertVerifier.sol";
import {IEigenDACertVerifierBase} from "src/integrations/cert/interfaces/IEigenDACertVerifierBase.sol";
import {IVersionedEigenDACertVerifier} from "src/integrations/cert/interfaces/IVersionedEigenDACertVerifier.sol";

import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";

import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";

import {IEigenDASemVer} from "src/core/interfaces/IEigenDASemVer.sol";

import {EigenDACertVerificationLib as CertLib} from "src/integrations/cert/libraries/EigenDACertVerificationLib.sol";
import {EigenDACertTypes as CT} from "src/integrations/cert/EigenDACertTypes.sol";

/// @title EigenDACertVerifier
/// @notice Verifies EigenDA certificates
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

    /// @notice Security thresholds used by {checkDACert}.
    /// @dev Checked inside {EigenDACertVerificationLib-checkDACert}. Constraints to respect:
    ///      - confirmationThreshold > adversaryThreshold (constructor-enforced)
    ///      - confirmationThreshold - adversaryThreshold > reconstructionThreshold
    ///        (see eigenda/docs/spec/src/protocol/architecture/security-parameters.md
    ///         for the definition of reconstructionThreshold and more info)
    DATypesV1.SecurityThresholds internal _securityThresholds;

    bytes internal _quorumNumbersRequired;
    uint16 internal _offchainDerivationVersion;

    uint8 internal constant CERT_VERSION = 4;

    uint8 internal constant MAJOR_VERSION = 4;
    uint8 internal constant MINOR_VERSION = 0;
    uint8 internal constant PATCH_VERSION = 0;

    /// @notice Status codes for certificate verification results
    /// @dev checkDACert calls are classified into: success (200), invalid_cert (400), and internal_error (500).
    enum StatusCode {
        NULL_ERROR, // Unused error code. If this is returned, there is a bug in the code.
        SUCCESS, // 200: Verification succeeded
        // The below 4 status codes are kept for backwards compatibility, but are no longer used.
        // We previously had plans to have more granular error codes, but decided this was not necessary,
        // and the only signal useful to offchain is to separate certs into: success, invalid (400), and bugs (500).
        UNUSED_HISTORICAL_INVALID_INCLUSION_PROOF,
        UNUSED_HISTORICAL_SECURITY_ASSUMPTIONS_NOT_MET,
        UNUSED_HISTORICAL_BLOB_QUORUMS_NOT_SUBSET,
        UNUSED_HISTORICAL_REQUIRED_QUORUMS_NOT_SUBSET,
        INVALID_CERT, // 400: Certificate is invalid due to some revert from the verification library
        INTERNAL_ERROR // 500: Bug or misconfiguration in the CertVerifier contract itself. This includes solidity panics and evm reverts.
    }

    constructor(
        IEigenDAThresholdRegistry initEigenDAThresholdRegistry,
        IEigenDASignatureVerifier initEigenDASignatureVerifier,
        DATypesV1.SecurityThresholds memory initSecurityThresholds,
        bytes memory initQuorumNumbersRequired,
        uint16 initOffchainDerivationVersion
    ) {
        if (initSecurityThresholds.confirmationThreshold <= initSecurityThresholds.adversaryThreshold) {
            revert InvalidSecurityThresholds();
        }
        if (initQuorumNumbersRequired.length == 0 || initQuorumNumbersRequired.length > 256) {
            revert InvalidQuorumNumbersRequired(initQuorumNumbersRequired.length);
        }
        _eigenDAThresholdRegistry = initEigenDAThresholdRegistry;
        _eigenDASignatureVerifier = initEigenDASignatureVerifier;
        _securityThresholds = initSecurityThresholds;
        _quorumNumbersRequired = initQuorumNumbersRequired;
        _offchainDerivationVersion = initOffchainDerivationVersion;
    }

    /// @notice Decodes a certificate from bytes to an EigenDACertV4
    /// @dev This function is external for the purpose of try/catch'ing it inside checkDACert,
    /// and should be considered an implementation detail. Do not rely on this function being
    /// part of the public interface of this contract.
    function _decodeCert(bytes calldata data) external pure returns (CT.EigenDACertV4 memory cert) {
        return abi.decode(data, (CT.EigenDACertV4));
    }

    /// @inheritdoc IEigenDACertVerifierBase
    /// @dev checkDACert is designed to be zk provable by risczero's Steel library,
    /// which does not support zk proving reverting calls: https://github.com/risc0/risc0-ethereum/issues/438.
    /// It try catches checkDACertReverts, and maps any reverts to status codes.
    /// This means invalid certs can easily be proven so by looking at the status code returned,
    /// which is also useful for optimistic rollup one step prover contracts.
    /// @dev Make sure to call this at a block number that is > RBN, otherwise this function will
    /// return an INVALID_CERT status code because of a require in the BLSSignatureChecker library that we use.
    function checkDACert(bytes calldata abiEncodedCert) external view returns (uint8) {
        CT.EigenDACertV4 memory daCert;
        // We try catch this here because decoding error would appear as a Panic,
        // which we consider bugs in the try/catch for the checkDACertReverts call below.
        try this._decodeCert(abiEncodedCert) returns (CT.EigenDACertV4 memory _daCert) {
            daCert = _daCert;
        } catch {
            return uint8(StatusCode.INVALID_CERT);
        }

        // The try catch below is used to filter certs into 3 status codes:
        // 1. success
        // 2. invalid cert (any failing require statement; we assume all require statements return either a string or custom error)
        // 3. internal error (everything else, including solidity panics and low-level evm reverts, basically anything unexpected)
        // TODO(samlaf): certVerifier should be set with a maxGas param that will be passed here, to enforce deterministic behavior
        // between different execution environments: EVM running onchain during optimistic rollup fraud proofs, zkVM, eth-call with higher gas limit.
        try this.checkDACertReverts(daCert) {
            return uint8(StatusCode.SUCCESS);
        } catch Error(string memory) {
            /*reason*/
            // This matches any require(..., "string reason") revert that is pre custom errors,
            // which many of our current eigenlayer-middleware dependencies like the BLSSignatureChecker still use. See:
            // https://github.com/Layr-Labs/eigenlayer-middleware/blob/fe5834371caed60c1d26ab62b5519b0cbdcb42fa/src/BLSSignatureChecker.sol#L96
            return uint8(StatusCode.INVALID_CERT);
        } catch Panic(uint256) {
            /*errorCode*/
            // This matches any panic (e.g. arithmetic overflow, division by zero, invalid array access, etc.),
            // which means a bug or misconfiguration of the CertVerifier contract itself.
            return uint8(StatusCode.INTERNAL_ERROR);
        } catch (bytes memory reason) {
            if (reason.length == 0) {
                // This matches low-level evm reverts like out-of-gas or stack too few values.
                // See https://rareskills.io/post/try-catch-solidity for more info.
                return uint8(StatusCode.INTERNAL_ERROR);
            } else if (reason.length < 4) {
                // Don't think this is possible...
                return uint8(StatusCode.INTERNAL_ERROR);
            }
            // Any revert here is from custom errors coming from a failed require(..., SomeCustomError()) statement.
            // This mean that the cert is invalid.
            return uint8(StatusCode.INVALID_CERT);
        }
    }

    /// @notice Check a DA cert's validity
    /// @param daCert The EigenDA certificate
    /// @dev This function will revert if the certificate is invalid.
    function checkDACertReverts(CT.EigenDACertV4 calldata daCert) external view {
        CertLib.checkDACert(
            _eigenDAThresholdRegistry,
            _eigenDASignatureVerifier,
            daCert,
            _securityThresholds,
            _quorumNumbersRequired,
            _offchainDerivationVersion
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

    /// @inheritdoc IEigenDACertVerifier
    function offchainDerivationVersion() external view returns (uint16) {
        return _offchainDerivationVersion;
    }

    /// @inheritdoc IVersionedEigenDACertVerifier
    function certVersion() external pure returns (uint8) {
        return CERT_VERSION;
    }

    /// @inheritdoc IEigenDASemVer
    function semver() external pure returns (uint8 major, uint8 minor, uint8 patch) {
        major = MAJOR_VERSION;
        minor = MINOR_VERSION;
        patch = PATCH_VERSION;
    }
}
