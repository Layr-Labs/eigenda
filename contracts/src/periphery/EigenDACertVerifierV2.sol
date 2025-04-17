// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifierV2} from "src/periphery/interfaces/IEigenDACertVerifierV2.sol";
import {IEigenDAThresholdRegistry} from "src/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/interfaces/IEigenDASignatureVerifier.sol";
import {IEigenDARelayRegistry} from "src/interfaces/IEigenDARelayRegistry.sol";
import {EigenDACertVerificationV2Lib as CertLib} from "src/periphery/libraries/EigenDACertVerificationV2Lib.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import "src/interfaces/IEigenDAStructs.sol";

/**
 * @title EigenDACertVerifierV2 - EigenDA V2 certificate verification
 * @author Layr Labs, Inc.
 * @notice Contains all V2-specific verification functionality
 */
contract EigenDACertVerifierV2 is IEigenDACertVerifierV2 {
    error InvalidSecurityThresholds();

    /// @notice The EigenDAThresholdRegistry contract address
    IEigenDAThresholdRegistry public immutable eigenDAThresholdRegistryV2;

    /// @notice The EigenDASignatureVerifier contract address
    IEigenDASignatureVerifier public immutable eigenDASignatureVerifierV2;

    /// @notice The EigenDARelayRegistry contract address
    IEigenDARelayRegistry public immutable eigenDARelayRegistryV2;

    /// @notice The EigenDA middleware OperatorStateRetriever contract address
    OperatorStateRetriever public immutable operatorStateRetrieverV2;

    /// @notice The EigenDA middleware RegistryCoordinator contract address
    IRegistryCoordinator public immutable registryCoordinatorV2;

    /// @notice Security thresholds for V2 certificate verification
    /// @notice Minimum percentage of stake required for confirmation of a blob (must be higher than adversaryThreshold)
    uint8 public immutable securityThresholdsConfirmation;
    /// @notice Maximum percentage of stake an adversary is assumed to control (must be lower than confirmationThreshold)
    uint8 public immutable securityThresholdsAdversary;

    /**
     * @notice Constructor for the EigenDA V2 certificate verifier
     * @param _eigenDAThresholdRegistryV2 The address of the EigenDAThresholdRegistry contract
     * @param _eigenDASignatureVerifierV2 The address of the EigenDASignatureVerifier contract
     * @param _eigenDARelayRegistryV2 The address of the EigenDARelayRegistry contract
     * @param _operatorStateRetrieverV2 The address of the OperatorStateRetriever contract
     * @param _registryCoordinatorV2 The address of the RegistryCoordinator contract
     * @param _securityThresholdsV2 The security thresholds for verification
     */
    constructor(
        IEigenDAThresholdRegistry _eigenDAThresholdRegistryV2,
        IEigenDASignatureVerifier _eigenDASignatureVerifierV2,
        IEigenDARelayRegistry _eigenDARelayRegistryV2,
        OperatorStateRetriever _operatorStateRetrieverV2,
        IRegistryCoordinator _registryCoordinatorV2,
        SecurityThresholds memory _securityThresholdsV2
    ) {
        eigenDAThresholdRegistryV2 = _eigenDAThresholdRegistryV2;
        eigenDASignatureVerifierV2 = _eigenDASignatureVerifierV2;
        eigenDARelayRegistryV2 = _eigenDARelayRegistryV2;
        operatorStateRetrieverV2 = _operatorStateRetrieverV2;
        registryCoordinatorV2 = _registryCoordinatorV2;

        // Store security thresholds as individual fields
        if (_securityThresholdsV2.confirmationThreshold <= _securityThresholdsV2.adversaryThreshold) {
            revert InvalidSecurityThresholds();
        }
        securityThresholdsConfirmation = _securityThresholdsV2.confirmationThreshold;
        securityThresholdsAdversary = _securityThresholdsV2.adversaryThreshold;
    }

    /**
     * @notice Verifies a blob cert using the immutable required quorums and security thresholds
     * @param batchHeader The batch header of the blob
     * @param blobInclusionInfo The inclusion proof for the blob cert
     * @param nonSignerStakesAndSignature The nonSignerStakesAndSignature to verify the blob cert against
     * @param signedQuorumNumbers The signed quorum numbers corresponding to the nonSignerStakesAndSignature
     */
    function verifyDACertV2(
        BatchHeaderV2 calldata batchHeader,
        BlobInclusionInfo calldata blobInclusionInfo,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature,
        bytes memory signedQuorumNumbers
    ) public view {
        (CertLib.ErrorCode err, bytes memory errParams) = CertLib.verifyDACertV2(
            _thresholdRegistry(),
            _signatureVerifier(),
            _relayRegistry(),
            batchHeader,
            blobInclusionInfo,
            nonSignerStakesAndSignature,
            _securityThresholds(),
            _quorumNumbersRequired(),
            signedQuorumNumbers
        );
        CertLib.revertOnError(err, errParams);
    }

    /**
     * @notice Verifies a blob cert using a signed batch
     * @param signedBatch The signed batch to verify the blob cert against
     * @param blobInclusionInfo The inclusion proof for the blob cert
     */
    function verifyDACertV2FromSignedBatch(
        SignedBatch calldata signedBatch,
        BlobInclusionInfo calldata blobInclusionInfo
    ) external view {
        (CertLib.ErrorCode err, bytes memory errParams) = CertLib.verifyDACertV2FromSignedBatch(
            _thresholdRegistry(),
            _signatureVerifier(),
            _relayRegistry(),
            _operatorStateRetriever(),
            _registryCoordinator(),
            signedBatch,
            blobInclusionInfo,
            _securityThresholds(),
            _quorumNumbersRequired()
        );
        CertLib.revertOnError(err, errParams);
    }

    /**
     * @notice Thin try/catch wrapper around verifyDACertV2 that returns false instead of reverting
     * @dev Useful for ZK proof systems that require a return value
     * @param batchHeader The batch header of the blob
     * @param blobInclusionInfo The inclusion proof for the blob cert
     * @param nonSignerStakesAndSignature The nonSignerStakesAndSignature to verify the blob cert against
     * @param signedQuorumNumbers The signed quorum numbers corresponding to the nonSignerStakesAndSignature
     * @return success True if verification succeeded
     */
    function verifyDACertV2ForZKProof(
        BatchHeaderV2 calldata batchHeader,
        BlobInclusionInfo calldata blobInclusionInfo,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature,
        bytes memory signedQuorumNumbers
    ) external view returns (bool success) {
        (CertLib.ErrorCode errorCode,) = CertLib.verifyDACertV2(
            _thresholdRegistry(),
            _signatureVerifier(),
            _relayRegistry(),
            batchHeader,
            blobInclusionInfo,
            nonSignerStakesAndSignature,
            _securityThresholds(),
            _quorumNumbersRequired(),
            signedQuorumNumbers
        );
        if (errorCode == CertLib.ErrorCode.SUCCESS) {
            return true;
        } else {
            return false;
        }
    }

    /**
     * @notice Returns the nonSignerStakesAndSignature for a given blob cert and signed batch
     * @param signedBatch The signed batch to get the nonSignerStakesAndSignature for
     * @return nonSignerStakesAndSignature The nonSignerStakesAndSignature for the given signed batch attestation
     */
    function getNonSignerStakesAndSignature(SignedBatch calldata signedBatch)
        external
        view
        returns (NonSignerStakesAndSignature memory nonSignerStakesAndSignature)
    {
        (nonSignerStakesAndSignature,) =
            CertLib.getNonSignerStakesAndSignature(_operatorStateRetriever(), _registryCoordinator(), signedBatch);
    }

    /**
     * @notice Verifies the security parameters for a blob cert
     * @param blobParams The blob params to verify
     * @param securityThresholds The security thresholds to verify against
     */
    function verifyDACertSecurityParams(
        VersionedBlobParams memory blobParams,
        SecurityThresholds memory securityThresholds
    ) external pure {
        (CertLib.ErrorCode err, bytes memory errParams) = CertLib.verifySecurityParams(blobParams, securityThresholds);
        CertLib.revertOnError(err, errParams);
    }

    /**
     * @notice Verifies the security parameters for a blob version
     * @param version The version of the blob to verify
     * @param securityThresholds The security thresholds to verify against
     */
    function verifyDACertSecurityParams(uint16 version, SecurityThresholds memory securityThresholds) external view {
        (CertLib.ErrorCode err, bytes memory errParams) =
            CertLib.verifySecurityParams(_thresholdRegistry().getBlobParams(version), securityThresholds);
        CertLib.revertOnError(err, errParams);
    }

    // Virtual accessor methods for dependency injection in derived contracts

    /**
     * @notice Returns the threshold registry contract
     * @return The IEigenDAThresholdRegistry contract
     * @dev Can be overridden by derived contracts
     */
    function _thresholdRegistry() internal view virtual returns (IEigenDAThresholdRegistry) {
        return eigenDAThresholdRegistryV2;
    }

    /**
     * @notice Returns the signature verifier contract
     * @return The IEigenDASignatureVerifier contract
     * @dev Can be overridden by derived contracts
     */
    function _signatureVerifier() internal view virtual returns (IEigenDASignatureVerifier) {
        return eigenDASignatureVerifierV2;
    }

    /**
     * @notice Returns the relay registry contract
     * @return The IEigenDARelayRegistry contract
     * @dev Can be overridden by derived contracts
     */
    function _relayRegistry() internal view virtual returns (IEigenDARelayRegistry) {
        return eigenDARelayRegistryV2;
    }

    /**
     * @notice Returns the operator state retriever contract
     * @return The OperatorStateRetriever contract
     * @dev Can be overridden by derived contracts
     */
    function _operatorStateRetriever() internal view virtual returns (OperatorStateRetriever) {
        return operatorStateRetrieverV2;
    }

    /**
     * @notice Returns the registry coordinator contract
     * @return The IRegistryCoordinator contract
     * @dev Can be overridden by derived contracts
     */
    function _registryCoordinator() internal view virtual returns (IRegistryCoordinator) {
        return registryCoordinatorV2;
    }

    /**
     * @notice Returns the security thresholds used for verification
     * @return The SecurityThresholds struct with confirmation and adversary thresholds
     * @dev Can be overridden by derived contracts
     */
    function _securityThresholds() internal view virtual returns (SecurityThresholds memory) {
        return SecurityThresholds({
            confirmationThreshold: securityThresholdsConfirmation,
            adversaryThreshold: securityThresholdsAdversary
        });
    }

    /**
     * @notice Returns the quorum numbers required for verification
     * @return bytes The required quorum numbers
     * @dev Can be overridden by derived contracts
     */
    function _quorumNumbersRequired() internal view virtual returns (bytes memory) {
        return _thresholdRegistry().quorumNumbersRequired();
    }
}
