// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAServiceManager} from "../interfaces/IEigenDAServiceManager.sol";
import {IEigenDABlobVerifier} from "../interfaces/IEigenDABlobVerifier.sol";
import {IEigenDAThresholdRegistry} from "../interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "../interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "../interfaces/IEigenDASignatureVerifier.sol";
import {EigenDABlobVerificationUtils} from "../libraries/EigenDABlobVerificationUtils.sol";

contract EigenDABlobVerifier is IEigenDABlobVerifier {

    IEigenDAThresholdRegistry public immutable eigenDAThresholdRegistry;
    IEigenDABatchMetadataStorage public immutable eigenDABatchMetadataStorage;
    IEigenDASignatureVerifier public immutable eigenDASignatureVerifier;

    constructor(
        IEigenDAThresholdRegistry _eigenDAThresholdRegistry,
        IEigenDABatchMetadataStorage _eigenDABatchMetadataStorage,
        IEigenDASignatureVerifier _eigenDASignatureVerifier
    ) {
        eigenDAThresholdRegistry = _eigenDAThresholdRegistry;
        eigenDABatchMetadataStorage = _eigenDABatchMetadataStorage;
        eigenDASignatureVerifier = _eigenDASignatureVerifier;
    }

    /**
     * @notice Verifies a the blob is valid for the required quorums
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob verification proof to verify the blob against
     */
    function verifyBlobV1(
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        EigenDABlobVerificationUtils.BlobVerificationProof calldata blobVerificationProof
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
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        EigenDABlobVerificationUtils.BlobVerificationProof calldata blobVerificationProof,
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

    function verifyBlobV2(
        EigenDABlobVerificationUtils.SignedCertificate calldata signedCertificate
    ) external view {
        EigenDABlobVerificationUtils._verifyBlobV2ForQuorums(
            eigenDASignatureVerifier,
            signedCertificate,
            quorumNumbersRequired()
        );
    }

    function verifyBlobV2(
        EigenDABlobVerificationUtils.SignedCertificate calldata signedCertificate,
        bytes calldata additionalQuorumNumbersRequired
    ) external view {
        EigenDABlobVerificationUtils._verifyBlobV2ForQuorums(
            eigenDASignatureVerifier,
            signedCertificate,
            bytes.concat(quorumNumbersRequired(), additionalQuorumNumbersRequired)
        );
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
}