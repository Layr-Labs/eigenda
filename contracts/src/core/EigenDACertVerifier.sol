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

/**
 * @title A CertVerifier is an immutable contract that is used by a consumer to verify EigenDA blob certificates 
 * @notice For V2 verification this contract is deployed with immutable security thresholds and required quorum numbers,
 *         to change these values or verification behavior a new CertVerifier must be deployed
 */
contract EigenDACertVerifier is IEigenDACertVerifier {

    IEigenDAThresholdRegistry public immutable eigenDAThresholdRegistry;
    IEigenDABatchMetadataStorage public immutable eigenDABatchMetadataStorage;
    IEigenDASignatureVerifier public immutable eigenDASignatureVerifier;
    IEigenDARelayRegistry public immutable eigenDARelayRegistry;

    OperatorStateRetriever public immutable operatorStateRetriever;
    IRegistryCoordinator public immutable registryCoordinator;

    SecurityThresholds public securityThresholdsV2;
    bytes public quorumNumbersRequiredV2;

    constructor(
        IEigenDAThresholdRegistry _eigenDAThresholdRegistry,
        IEigenDABatchMetadataStorage _eigenDABatchMetadataStorage,
        IEigenDASignatureVerifier _eigenDASignatureVerifier,
        IEigenDARelayRegistry _eigenDARelayRegistry,
        OperatorStateRetriever _operatorStateRetriever,
        IRegistryCoordinator _registryCoordinator,
        SecurityThresholds memory _securityThresholdsV2,
        bytes memory _quorumNumbersRequiredV2
    ) {
        eigenDAThresholdRegistry = _eigenDAThresholdRegistry;
        eigenDABatchMetadataStorage = _eigenDABatchMetadataStorage;
        eigenDASignatureVerifier = _eigenDASignatureVerifier;
        eigenDARelayRegistry = _eigenDARelayRegistry;
        operatorStateRetriever = _operatorStateRetriever;
        registryCoordinator = _registryCoordinator;

        // confirmation and adversary signing thresholds that must be met for all quorums in a V2 certificate
        securityThresholdsV2 = _securityThresholdsV2;

        // quorum numbers that must be validated in a V2 certificate
        quorumNumbersRequiredV2 = _quorumNumbersRequiredV2;
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
     * @notice Verifies a blob cert using the immutable required quorums and security thresholds set in the constructor
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
            securityThresholdsV2,
            quorumNumbersRequiredV2
        );
    }

    /**
     * @notice Verifies a blob cert using the immutable required quorums and security thresholds set in the constructor
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
            securityThresholdsV2,
            quorumNumbersRequiredV2
        );
    }

    /**
     * @notice Thin try/catch wrapper around verifyDACertV2 that returns false instead of panicing
     * @dev The Steel library (https://github.com/risc0/risc0-ethereum/tree/main/crates/steel) 
     *      currently has a limitation that it can only create zk proofs for functions that return a value
     * @param batchHeader The batch header of the blob 
     * @param blobInclusionInfo The inclusion proof for the blob cert
     * @param nonSignerStakesAndSignature The nonSignerStakesAndSignature to verify the blob cert against
     */
    function verifyDACertV2ForZKProof(
        BatchHeaderV2 calldata batchHeader,
        BlobInclusionInfo calldata blobInclusionInfo,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature
    ) external view returns (bool) {
        try EigenDACertVerificationUtils.verifyDACertV2ForQuorumsExternal(
            eigenDAThresholdRegistry,
            eigenDASignatureVerifier,
            eigenDARelayRegistry,
            batchHeader,
            blobInclusionInfo,
            nonSignerStakesAndSignature,
            securityThresholdsV2,
            quorumNumbersRequiredV2
        ) {
            return true;
        } catch {
            return false;
        }
    }

    ///////////////////////// HELPER FUNCTIONS ///////////////////////////////

    /**
     * @notice Returns the nonSignerStakesAndSignature for a given signed batch
     * @param signedBatch The signed batch to get the nonSignerStakesAndSignature for
     */
    function getNonSignerStakesAndSignature(
        SignedBatch calldata signedBatch
    ) external view returns (NonSignerStakesAndSignature memory) {
        bytes memory quorumNumbers;
        for (uint i = 0; i < signedBatch.attestation.quorumNumbers.length; ++i) {
            quorumNumbers = abi.encodePacked(quorumNumbers, uint8(signedBatch.attestation.quorumNumbers[i]));
        }
    
        return EigenDACertVerificationUtils._getNonSignerStakesAndSignature(
            operatorStateRetriever, 
            registryCoordinator, 
            signedBatch,
            quorumNumbers
        );
    }

    /**
     * @notice Returns the nonSignerStakesAndSignature for a given signed batch with specific quorum numbers
     * @param signedBatch The signed batch to get the nonSignerStakesAndSignature for
     * @param quorumNumbers The quorum numbers to get the nonSignerStakesAndSignature for
     */
    function getNonSignerStakesAndSignature(
        SignedBatch calldata signedBatch,
        bytes memory quorumNumbers
    ) external view returns (NonSignerStakesAndSignature memory) {
        return EigenDACertVerificationUtils._getNonSignerStakesAndSignature(
            operatorStateRetriever, 
            registryCoordinator, 
            signedBatch,
            quorumNumbers
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
}
