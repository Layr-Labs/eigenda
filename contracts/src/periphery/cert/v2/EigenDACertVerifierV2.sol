// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "src/core/interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";
import {EigenDACertVerificationV1Lib as CertV1Lib} from "src/periphery/cert/v1/EigenDACertVerificationV1Lib.sol";
import {EigenDACertVerificationV2Lib as CertV2Lib} from "src/periphery/cert/v2/EigenDACertVerificationV2Lib.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {EigenDATypesV2 as DATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";

/**
 * @title A CertVerifier is an immutable contract that is used by a consumer to verify EigenDA blob certificates
 * @notice For V2 verification this contract is deployed with immutable security thresholds and required quorum numbers,
 *         to change these values or verification behavior a new CertVerifier must be deployed
 */
contract EigenDACertVerifierV2 {
    error InvalidSecurityThresholds();

    /// @notice The EigenDAThresholdRegistry contract address
    IEigenDAThresholdRegistry public immutable eigenDAThresholdRegistryV2;

    /// @notice The EigenDASignatureVerifier contract address
    IEigenDASignatureVerifier public immutable eigenDASignatureVerifierV2;

    /// @notice The EigenDA middleware OperatorStateRetriever contract address
    OperatorStateRetriever public immutable operatorStateRetrieverV2;

    /// @notice The EigenDA middleware RegistryCoordinator contract address
    IRegistryCoordinator public immutable registryCoordinatorV2;

    DATypesV1.SecurityThresholds public securityThresholdsV2;

    bytes public quorumNumbersRequiredV2;

    /**
     * @notice Constructor for the EigenDA V2 certificate verifier
     * @param _eigenDAThresholdRegistryV2 The address of the EigenDAThresholdRegistry contract
     * @param _eigenDASignatureVerifierV2 The address of the EigenDASignatureVerifier contract
     * @param _operatorStateRetrieverV2 The address of the OperatorStateRetriever contract
     * @param _registryCoordinatorV2 The address of the RegistryCoordinator contract
     * @param _securityThresholdsV2 The security thresholds for verification
     */
    constructor(
        IEigenDAThresholdRegistry _eigenDAThresholdRegistryV2,
        IEigenDASignatureVerifier _eigenDASignatureVerifierV2,
        OperatorStateRetriever _operatorStateRetrieverV2,
        IRegistryCoordinator _registryCoordinatorV2,
        DATypesV1.SecurityThresholds memory _securityThresholdsV2,
        bytes memory _quorumNumbersRequiredV2
    ) {
        if (_securityThresholdsV2.confirmationThreshold <= _securityThresholdsV2.adversaryThreshold) {
            revert InvalidSecurityThresholds();
        }
        eigenDAThresholdRegistryV2 = _eigenDAThresholdRegistryV2;
        eigenDASignatureVerifierV2 = _eigenDASignatureVerifierV2;
        operatorStateRetrieverV2 = _operatorStateRetrieverV2;
        registryCoordinatorV2 = _registryCoordinatorV2;
        securityThresholdsV2 = _securityThresholdsV2;
        quorumNumbersRequiredV2 = _quorumNumbersRequiredV2;
    }

    /**
     * @notice Verifies a blob cert using the immutable required quorums and security thresholds set in the constructor
     * @param batchHeader The batch header of the blob
     * @param blobInclusionInfo The inclusion proof for the blob cert
     * @param nonSignerStakesAndSignature The nonSignerStakesAndSignature to verify the blob cert against
     * @param signedQuorumNumbers The signed quorum numbers corresponding to the nonSignerStakesAndSignature
     */
    function verifyDACertV2(
        DATypesV2.BatchHeaderV2 calldata batchHeader,
        DATypesV2.BlobInclusionInfo calldata blobInclusionInfo,
        DATypesV1.NonSignerStakesAndSignature calldata nonSignerStakesAndSignature,
        bytes memory signedQuorumNumbers
    ) external view {
        CertV2Lib.verifyDACertV2(
            _thresholdRegistry(),
            _signatureVerifier(),
            batchHeader,
            blobInclusionInfo,
            nonSignerStakesAndSignature,
            _securityThresholds(),
            _quorumNumbersRequired(),
            signedQuorumNumbers
        );
    }

    /**
     * @notice Verifies a blob cert using the immutable required quorums and security thresholds set in the constructor
     * @param signedBatch The signed batch to verify the blob cert against
     * @param blobInclusionInfo The inclusion proof for the blob cert
     */
    function verifyDACertV2FromSignedBatch(
        DATypesV2.SignedBatch calldata signedBatch,
        DATypesV2.BlobInclusionInfo calldata blobInclusionInfo
    ) external view {
        CertV2Lib.verifyDACertV2FromSignedBatch(
            _thresholdRegistry(),
            _signatureVerifier(),
            _operatorStateRetriever(),
            _registryCoordinator(),
            signedBatch,
            blobInclusionInfo,
            _securityThresholds(),
            _quorumNumbersRequired()
        );
    }

    /**
     * @notice Thin try/catch wrapper around verifyDACertV2 that returns false instead of panicing
     * @dev The Steel library (https://github.com/risc0/risc0-ethereum/tree/main/crates/steel)
     *      currently has a limitation that it can only create zk proofs for functions that return a value
     * @param batchHeader The batch header of the blob
     * @param blobInclusionInfo The inclusion proof for the blob cert
     * @param nonSignerStakesAndSignature The nonSignerStakesAndSignature to verify the blob cert against
     * @param signedQuorumNumbers The signed quorum numbers corresponding to the nonSignerStakesAndSignature
     */
    function verifyDACertV2ForZKProof(
        DATypesV2.BatchHeaderV2 calldata batchHeader,
        DATypesV2.BlobInclusionInfo calldata blobInclusionInfo,
        DATypesV1.NonSignerStakesAndSignature calldata nonSignerStakesAndSignature,
        bytes memory signedQuorumNumbers
    ) external view returns (bool) {
        (CertV2Lib.StatusCode status,) = CertV2Lib.checkDACertV2(
            _thresholdRegistry(),
            _signatureVerifier(),
            batchHeader,
            blobInclusionInfo,
            nonSignerStakesAndSignature,
            _securityThresholds(),
            _quorumNumbersRequired(),
            signedQuorumNumbers
        );
        if (status == CertV2Lib.StatusCode.SUCCESS) {
            return true;
        } else {
            return false;
        }
    }

    function getNonSignerStakesAndSignature(DATypesV2.SignedBatch calldata signedBatch)
        external
        view
        returns (DATypesV1.NonSignerStakesAndSignature memory)
    {
        (DATypesV1.NonSignerStakesAndSignature memory nonSignerStakesAndSignature,) =
            CertV2Lib.getNonSignerStakesAndSignature(operatorStateRetrieverV2, registryCoordinatorV2, signedBatch);
        return nonSignerStakesAndSignature;
    }

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
    function _securityThresholds() internal view virtual returns (DATypesV1.SecurityThresholds memory) {
        return securityThresholdsV2;
    }

    /**
     * @notice Returns the quorum numbers required for verification
     * @return bytes The required quorum numbers
     * @dev Can be overridden by derived contracts
     */
    function _quorumNumbersRequired() internal view virtual returns (bytes memory) {
        return quorumNumbersRequiredV2;
    }
}
