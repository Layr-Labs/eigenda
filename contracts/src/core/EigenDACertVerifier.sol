// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifier} from "../interfaces/IEigenDACertVerifier.sol";
import {IEigenDAThresholdRegistry} from "../interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "../interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "../interfaces/IEigenDASignatureVerifier.sol";
import {EigenDACertVerificationUtils} from "../libraries/EigenDACertVerificationUtils.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {IEigenDARelayRegistry} from "../interfaces/IEigenDARelayRegistry.sol";
import "../interfaces/IEigenDAStructs.sol";

contract EigenDACertVerifier is IEigenDACertVerifier {

    IEigenDAThresholdRegistry public immutable eigenDAThresholdRegistry;
    IEigenDABatchMetadataStorage public immutable eigenDABatchMetadataStorage;
    IEigenDASignatureVerifier public immutable eigenDASignatureVerifier;
    IEigenDARelayRegistry public immutable eigenDARelayRegistry;

    OperatorStateRetriever public immutable operatorStateRetriever;
    IRegistryCoordinator public immutable registryCoordinator;

    constructor(
        IEigenDAThresholdRegistry _eigenDAThresholdRegistry,
        IEigenDABatchMetadataStorage _eigenDABatchMetadataStorage,
        IEigenDASignatureVerifier _eigenDASignatureVerifier,
        IEigenDARelayRegistry _eigenDARelayRegistry,
        OperatorStateRetriever _operatorStateRetriever,
        IRegistryCoordinator _registryCoordinator
    ) {
        eigenDAThresholdRegistry = _eigenDAThresholdRegistry;
        eigenDABatchMetadataStorage = _eigenDABatchMetadataStorage;
        eigenDASignatureVerifier = _eigenDASignatureVerifier;
        eigenDARelayRegistry = _eigenDARelayRegistry;
        operatorStateRetriever = _operatorStateRetriever;
        registryCoordinator = _registryCoordinator;
    }

    ///////////////////////// V1 ///////////////////////////////

    /**
     * @notice Verifies a the blob cert is valid for the required quorums
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob cert verification proof to verify
     */
    function verifyDACertV1(
        BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof
    ) external view {
        EigenDACertVerificationUtils._verifyDACertV1ForQuorums(
            eigenDAThresholdRegistry,
            eigenDABatchMetadataStorage, 
            blobHeader, 
            blobVerificationProof, 
            quorumNumbersRequired()
        );
    }

    /**
     * @notice Verifies a batch of blob certs for the required quorums
     * @param blobHeaders The blob headers to verify
     * @param blobVerificationProofs The blob cert verification proofs to verify against
     */
    function verifyDACertsV1(
        BlobHeader[] calldata blobHeaders,
        BlobVerificationProof[] calldata blobVerificationProofs
    ) external view {
        EigenDACertVerificationUtils._verifyDACertsV1ForQuorums(
            eigenDAThresholdRegistry,
            eigenDABatchMetadataStorage, 
            blobHeaders, 
            blobVerificationProofs, 
            quorumNumbersRequired()
        );
    }

    ///////////////////////// V2 ///////////////////////////////

    /**
     * @notice Verifies a blob cert for the specified quorums with the default security thresholds
     * @param batchHeader The batch header of the blob 
     * @param blobInclusionInfo The inclusion proof for the blob cert
     * @param nonSignerStakesAndSignature The nonSignerStakesAndSignature to verify the blob cert against
     */
    function verifyDACertV2(
        BatchHeaderV2 calldata batchHeader,
        BlobInclusionInfo calldata blobInclusionInfo,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature
    ) external view {
        EigenDACertVerificationUtils._verifyDACertV2ForQuorums(
            eigenDAThresholdRegistry,
            eigenDASignatureVerifier,
            eigenDARelayRegistry,
            batchHeader,
            blobInclusionInfo,
            nonSignerStakesAndSignature,
            getDefaultSecurityThresholdsV2(),
            blobInclusionInfo.blobCertificate.blobHeader.quorumNumbers
        );
    }

    /**
     * @notice Verifies a blob cert for the specified quorums with the default security thresholds
     * @param signedBatch The signed batch to verify the blob cert against
     * @param blobInclusionInfo The inclusion proof for the blob cert
     */
    function verifyDACertV2FromSignedBatch(
        SignedBatch calldata signedBatch,
        BlobInclusionInfo calldata blobInclusionInfo
    ) external view {
        EigenDACertVerificationUtils._verifyDACertV2ForQuorumsFromSignedBatch(
            eigenDAThresholdRegistry,
            eigenDASignatureVerifier,
            eigenDARelayRegistry,
            operatorStateRetriever,
            registryCoordinator,
            signedBatch,
            blobInclusionInfo,
            getDefaultSecurityThresholdsV2(),
            blobInclusionInfo.blobCertificate.blobHeader.quorumNumbers
        );
    }

    ///////////////////////// HELPER FUNCTIONS ///////////////////////////////

    /**
     * @notice Returns the nonSignerStakesAndSignature for a given blob cert and signed batch
     * @param signedBatch The signed batch to get the nonSignerStakesAndSignature for
     */
    function getNonSignerStakesAndSignature(
        SignedBatch calldata signedBatch
    ) external view returns (NonSignerStakesAndSignature memory) {
        return EigenDACertVerificationUtils._getNonSignerStakesAndSignature(
            operatorStateRetriever, 
            registryCoordinator, 
            signedBatch
        );
    }

    /**
     * @notice Verifies the security parameters for a blob cert
     * @param blobParams The blob params to verify 
     * @param securityThresholds The security thresholds to verify against
     */
    function verifyDACertSecurityParams(
        VersionedBlobParams memory blobParams,
        SecurityThresholds memory securityThresholds
    ) external view {
        EigenDACertVerificationUtils._verifyDACertSecurityParams(blobParams, securityThresholds);
    }

    /**
     * @notice Verifies the security parameters for a blob cert
     * @param version The version of the blob to verify
     * @param securityThresholds The security thresholds to verify against
     */
    function verifyDACertSecurityParams(
        uint16 version,
        SecurityThresholds memory securityThresholds
    ) external view {
        EigenDACertVerificationUtils._verifyDACertSecurityParams(getBlobParams(version), securityThresholds);
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

    /// @notice Returns the blob params for a given blob version
    function getBlobParams(uint16 version) public view returns (VersionedBlobParams memory) {
        return eigenDAThresholdRegistry.getBlobParams(version);
    }

    /// @notice Gets the default security thresholds for V2
    function getDefaultSecurityThresholdsV2() public view returns (SecurityThresholds memory) {
        return eigenDAThresholdRegistry.getDefaultSecurityThresholdsV2();
    }
}
