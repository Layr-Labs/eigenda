// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDABlobVerifier} from "../interfaces/IEigenDABlobVerifier.sol";
import {IEigenDAThresholdRegistry} from "../interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "../interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "../interfaces/IEigenDASignatureVerifier.sol";
import {EigenDABlobVerificationUtils} from "../libraries/EigenDABlobVerificationUtils.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import "../interfaces/IEigenDAStructs.sol";

contract EigenDABlobVerifier is IEigenDABlobVerifier {

    IEigenDAThresholdRegistry public immutable eigenDAThresholdRegistry;
    IEigenDABatchMetadataStorage public immutable eigenDABatchMetadataStorage;
    IEigenDASignatureVerifier public immutable eigenDASignatureVerifier;

    OperatorStateRetriever public immutable operatorStateRetriever;
    IRegistryCoordinator public immutable registryCoordinator;

    constructor(
        IEigenDAThresholdRegistry _eigenDAThresholdRegistry,
        IEigenDABatchMetadataStorage _eigenDABatchMetadataStorage,
        IEigenDASignatureVerifier _eigenDASignatureVerifier,
        OperatorStateRetriever _operatorStateRetriever,
        IRegistryCoordinator _registryCoordinator
    ) {
        eigenDAThresholdRegistry = _eigenDAThresholdRegistry;
        eigenDABatchMetadataStorage = _eigenDABatchMetadataStorage;
        eigenDASignatureVerifier = _eigenDASignatureVerifier;

        operatorStateRetriever = _operatorStateRetriever;
        registryCoordinator = _registryCoordinator;
    }

    /**
     * @notice Verifies a the blob is valid for the required quorums
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob verification proof to verify the blob against
     */
    function verifyBlobV1(
        BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof
    ) external view {
        EigenDABlobVerificationUtils._verifyBlobV1ForQuorums(
            eigenDAThresholdRegistry,
            eigenDABatchMetadataStorage, 
            blobHeader, 
            blobVerificationProof, 
            quorumNumbersRequired()
        );
    }

    /**
     * @notice Verifies that a blob is valid for the required quorums and additional quorums
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob verification proof to verify the blob against
     * @param additionalQuorumNumbersRequired The additional required quorum numbers 
     */
    function verifyBlobV1(
        BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof,
        bytes calldata additionalQuorumNumbersRequired
    ) external view {
        EigenDABlobVerificationUtils._verifyBlobV1ForQuorums(
            eigenDAThresholdRegistry,
            eigenDABatchMetadataStorage, 
            blobHeader, 
            blobVerificationProof, 
            bytes.concat(quorumNumbersRequired(), additionalQuorumNumbersRequired)
        );
    }

    /**
     * @notice Verifies a batch of blobs for the required quorums
     * @param blobHeaders The blob headers to verify
     * @param blobVerificationProofs The blob verification proofs to verify the blobs against
     */
    function verifyBlobsV1(
        BlobHeader[] calldata blobHeaders,
        BlobVerificationProof[] calldata blobVerificationProofs
    ) external view {
        EigenDABlobVerificationUtils._verifyBlobsV1ForQuorums(
            eigenDAThresholdRegistry,
            eigenDABatchMetadataStorage, 
            blobHeaders, 
            blobVerificationProofs, 
            quorumNumbersRequired()
        );
    }

    /**
     * @notice Verifies a batch of blobs for the required quorums and additional quorums
     * @param blobHeaders The blob headers to verify
     * @param blobVerificationProofs The blob verification proofs to verify the blobs against
     * @param additionalQuorumNumbersRequired The additional required quorum numbers
     */
    function verifyBlobsV1(
        BlobHeader[] calldata blobHeaders,
        BlobVerificationProof[] calldata blobVerificationProofs,
        bytes calldata additionalQuorumNumbersRequired
    ) external view {
        EigenDABlobVerificationUtils._verifyBlobsV1ForQuorums(
            eigenDAThresholdRegistry,
            eigenDABatchMetadataStorage, 
            blobHeaders, 
            blobVerificationProofs, 
            bytes.concat(quorumNumbersRequired(), additionalQuorumNumbersRequired)
        );
    }

    /**
     * @notice Verifies a blob for the required quorums and the default security thresholds
     * @param batchHeader The batch header of the blob
     * @param blobVerificationProof The blob verification proof for the blob
     * @param nonSignerStakesAndSignature The nonSignerStakesAndSignature for the blob
     */
    function verifyBlobV2(
        BatchHeaderV2 calldata batchHeader,
        BlobVerificationProofV2 calldata blobVerificationProof,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature
    ) external view {
        EigenDABlobVerificationUtils._verifyBlobV2ForQuorums(
            eigenDAThresholdRegistry,
            eigenDASignatureVerifier,
            batchHeader,
            blobVerificationProof,
            nonSignerStakesAndSignature,
            getDefaultSecurityThresholdsV2(),
            quorumNumbersRequired()
        );
    }

    /**
     * @notice Verifies a blob for the required quorums and additional quorums and the default security thresholds
     * @param batchHeader The batch header of the blob
     * @param blobVerificationProof The blob verification proof for the blob
     * @param nonSignerStakesAndSignature The nonSignerStakesAndSignature for the blob
     * @param additionalQuorumNumbersRequired The additional required quorum numbers
     */
    function verifyBlobV2(
        BatchHeaderV2 calldata batchHeader,
        BlobVerificationProofV2 calldata blobVerificationProof,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature,
        bytes calldata additionalQuorumNumbersRequired
    ) external view {
        EigenDABlobVerificationUtils._verifyBlobV2ForQuorums(
            eigenDAThresholdRegistry,
            eigenDASignatureVerifier,
            batchHeader,
            blobVerificationProof,
            nonSignerStakesAndSignature,
            getDefaultSecurityThresholdsV2(),
            bytes.concat(quorumNumbersRequired(), additionalQuorumNumbersRequired)
        );
    }

    /**
     * @notice Verifies a blob for the required quorums and additional quorums and a custom security threshold
     * @param batchHeader The batch header of the blob
     * @param blobVerificationProof The blob verification proof for the blob
     * @param nonSignerStakesAndSignature The nonSignerStakesAndSignature for the blob
     * @param securityThreshold The custom security threshold to verify the blob against
     * @param additionalQuorumNumbersRequired The additional required quorum numbers
     */
    function verifyBlobV2(
        BatchHeaderV2 calldata batchHeader,
        BlobVerificationProofV2 calldata blobVerificationProof,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature,
        SecurityThresholds memory securityThreshold,
        bytes calldata additionalQuorumNumbersRequired
    ) external view {
        EigenDABlobVerificationUtils._verifyBlobV2ForQuorums(
            eigenDAThresholdRegistry,
            eigenDASignatureVerifier,
            batchHeader,
            blobVerificationProof,
            nonSignerStakesAndSignature,
            securityThreshold,
            bytes.concat(quorumNumbersRequired(), additionalQuorumNumbersRequired)
        );
    }

    /**
     * @notice Verifies a blob for the required quorums and additional quorums and a set of custom security thresholds
     * @param batchHeader The batch header of the blob
     * @param blobVerificationProof The blob verification proof for the blob
     * @param nonSignerStakesAndSignature The nonSignerStakesAndSignature for the blob
     * @param securityThresholds The set of custom security thresholds to verify the blob against
     * @param additionalQuorumNumbersRequired The additional required quorum numbers
     */
    function verifyBlobV2(
        BatchHeaderV2 calldata batchHeader,
        BlobVerificationProofV2 calldata blobVerificationProof,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature,
        SecurityThresholds[] memory securityThresholds,
        bytes calldata additionalQuorumNumbersRequired
    ) external view {
        EigenDABlobVerificationUtils._verifyBlobV2ForQuorumsForThresholds(
            eigenDAThresholdRegistry,
            eigenDASignatureVerifier,
            batchHeader,
            blobVerificationProof,
            nonSignerStakesAndSignature,
            securityThresholds,
            bytes.concat(quorumNumbersRequired(), additionalQuorumNumbersRequired)
        );
    }

    /**
     * @notice Returns the nonSignerStakesAndSignature for a given blob and signed batch
     * @param signedBatch The signed batch to get the nonSignerStakesAndSignature for
     * @param blobHeader The blob header of the blob in the signed batch
     */
    function getNonSignerStakesAndSignature(
        SignedBatch calldata signedBatch,
        BlobHeaderV2 calldata blobHeader
    ) external view returns (NonSignerStakesAndSignature memory) {
        return EigenDABlobVerificationUtils._getNonSignerStakesAndSignature(
            operatorStateRetriever, 
            registryCoordinator, 
            signedBatch.attestation, 
            blobHeader.quorumNumbers
        );
    }

    /**
     * @notice Verifies the security parameters for a blob
     * @param blobParams The blob params to verify 
     * @param securityThresholds The security thresholds to verify against
     */
    function verifyBlobSecurityParams(
        VersionedBlobParams memory blobParams,
        SecurityThresholds memory securityThresholds
    ) external view {
        EigenDABlobVerificationUtils._verifyBlobSecurityParams(blobParams, securityThresholds);
    }

    /**
     * @notice Verifies the security parameters for a blob
     * @param version The version of the blob to verify
     * @param securityThresholds The security thresholds to verify against
     */
    function verifyBlobSecurityParams(
        uint16 version,
        SecurityThresholds memory securityThresholds
    ) external view {
        EigenDABlobVerificationUtils._verifyBlobSecurityParams(getBlobParams(version), securityThresholds);
    }

    /// @notice Returns the blob params for a given blob version
    function getBlobParams(uint16 version) public view returns (VersionedBlobParams memory) {
        return eigenDAThresholdRegistry.getBlobParams(version);
    }

    /// @notice Returns an array of bytes where each byte represents the adversary threshold percentage of the quorum at that index
    function quorumAdversaryThresholdPercentages() external view returns (bytes memory) {
        return eigenDAThresholdRegistry.quorumAdversaryThresholdPercentages();
    }

    /// @notice Returns an array of bytes where each byte represents the confirmation threshold percentage of the quorum at that index
    function quorumConfirmationThresholdPercentages() external view returns (bytes memory) {
        return eigenDAThresholdRegistry.quorumConfirmationThresholdPercentages();
    }

    /// @notice Returns an array of bytes where each byte represents the number of a required quorum 
    function quorumNumbersRequired() public view returns (bytes memory) {
        return eigenDAThresholdRegistry.quorumNumbersRequired();
    }

    function getQuorumAdversaryThresholdPercentage(
        uint8 quorumNumber
    ) external view returns (uint8){
        return eigenDAThresholdRegistry.getQuorumAdversaryThresholdPercentage(quorumNumber);
    }

    /// @notice Gets the confirmation threshold percentage for a quorum
    function getQuorumConfirmationThresholdPercentage(
        uint8 quorumNumber
    ) external view returns (uint8){
        return eigenDAThresholdRegistry.getQuorumConfirmationThresholdPercentage(quorumNumber);
    }

    /// @notice Checks if a quorum is required
    function getIsQuorumRequired(
        uint8 quorumNumber
    ) external view returns (bool){
        return eigenDAThresholdRegistry.getIsQuorumRequired(quorumNumber);
    }   

    /// @notice Gets the default security thresholds for V2
    function getDefaultSecurityThresholdsV2() public view returns (SecurityThresholds memory) {
        return eigenDAThresholdRegistry.getDefaultSecurityThresholdsV2();
    }
}
