// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "../EigenDABlobVerifier.sol";
import "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";

abstract contract EigenDABlobVerifierL2 is OwnableUpgradeable, EigenDABlobVerifier {

    IEigenDASignatureVerifier public immutable signatureVerifier;
    IEigenDABatchMetadataStorage public immutable batchMetadataStorage;

    constructor(
        IEigenDASignatureVerifier _signatureVerifier, 
        IEigenDABatchMetadataStorage _batchMetadataStorage
    ) {
        signatureVerifier = _signatureVerifier;
        batchMetadataStorage = _batchMetadataStorage;
        _disableInitializers();
    }

    function initialize(address initialOwner) external initializer {
        _transferOwnership(initialOwner);
    }

    /**
     * @notice Verifies a the blob is valid for the required quorums
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob verification proof to verify the blob against
     */
    function verifyBlob(
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof
    ) external view override {
        /*
        _verifyBlobForQuorums(
            batchMetadataStorage, 
            blobHeader, 
            blobVerificationProof, 
            quorumNumbersRequired
        );
        */
    }

    /**
     * @notice Verifies that a blob is valid for the required quorums and additional quorums
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob verification proof to verify the blob against
     * @param additionalQuorumNumbersRequired The additional required quorum numbers 
     */
    function verifyBlob(
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof,
        bytes memory additionalQuorumNumbersRequired
    ) external view override {
        /*
        _verifyBlobForQuorums(
            batchMetadataStorage, 
            blobHeader, 
            blobVerificationProof, 
            bytes.concat(quorumNumbersRequired, additionalQuorumNumbersRequired)
        );
        */
    }

    /**
     * @notice Verifies that a blob preconfirmation is valid for the required quorums
     * @param miniBatchHeader The mini batch header to verify
     * @param blobHeader The blob header to verify
     * @param nonSignerStakesAndSignature The operator signatures returned as the preconfirmation
     */
    function verifyPreconfirmation(
        IEigenDAServiceManager.BatchHeader calldata miniBatchHeader,
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        IEigenDASignatureVerifier.NonSignerStakesAndSignature memory nonSignerStakesAndSignature
        ) external view override {
        /*
        _verifyPreconfirmationForQuorums(
            signatureVerifier, 
            miniBatchHeader, 
            blobHeader,
            nonSignerStakesAndSignature, 
            quorumNumbersRequired
        );
        */
    }

    /**
     * @notice Verifies that a blob preconfirmation is valid for the required quorums and additional quorums
     * @param miniBatchHeader The mini batch header to verify
     * @param blobHeader The blob header to verify
     * @param nonSignerStakesAndSignature The operator signatures returned as the preconfirmation
     * @param additionalQuorumNumbersRequired The additional required quorum numbers 
     */
    function verifyPreconfirmation(
        IEigenDAServiceManager.BatchHeader calldata miniBatchHeader,
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        IEigenDASignatureVerifier.NonSignerStakesAndSignature memory nonSignerStakesAndSignature,
        bytes memory additionalQuorumNumbersRequired
    ) external view override {
        /*
        _verifyPreconfirmationForQuorums(
            signatureVerifier, 
            miniBatchHeader, 
            blobHeader,
            nonSignerStakesAndSignature, 
            bytes.concat(quorumNumbersRequired, additionalQuorumNumbersRequired)
        );
        */
    }

}