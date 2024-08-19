// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAServiceManager} from "./IEigenDAServiceManager.sol";
import {IEigenDASignatureVerifier} from "./IEigenDASignatureVerifier.sol";

interface IEigenDABlobVerifier {

    struct BlobVerificationProof {
        uint32 batchId;
        uint32 blobIndex;
        IEigenDAServiceManager.BatchMetadata batchMetadata;
        bytes inclusionProof;
        bytes quorumIndices;
    }

    /// @notice Returns an array of bytes where each byte represents the adversary threshold percentage of the quorum at that index
    function quorumAdversaryThresholdPercentages() external view returns (bytes memory);

    /// @notice Returns an array of bytes where each byte represents the confirmation threshold percentage of the quorum at that index
    function quorumConfirmationThresholdPercentages() external view returns (bytes memory);

    /// @notice Returns an array of bytes where each byte represents the number of a required quorum 
    function quorumNumbersRequired() external view returns (bytes memory);

    /**
     * @notice Verifies a the blob is valid for the required quorums
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob verification proof to verify the blob against
     */
    function verifyBlob(
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof
    ) external view;

    /**
     * @notice Verifies that a blob is valid for the required quorums and additional quorums
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob verification proof to verify the blob against
     * @param additionalQuorumNumbersRequired The additional required quorum numbers 
     */
    function verifyBlobForAdditionalQuorums(
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof,
        bytes calldata additionalQuorumNumbersRequired
    ) external view;

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
    ) external view;

    /**
     * @notice Verifies that a blob preconfirmation is valid for the required quorums and additional quorums
     * @param miniBatchHeader The mini batch header to verify
     * @param blobHeader The blob header to verify
     * @param nonSignerStakesAndSignature The operator signatures returned as the preconfirmation
     * @param additionalQuorumNumbersRequired The additional required quorum numbers 
     */
    function verifyPreconfirmationForAdditionalQuorums(
        IEigenDAServiceManager.BatchHeader calldata miniBatchHeader,
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        IEigenDASignatureVerifier.NonSignerStakesAndSignature memory nonSignerStakesAndSignature,
        bytes memory additionalQuorumNumbersRequired
    ) external view;

    /// @notice Gets the adversary threshold percentage for a quorum
    function getQuorumAdversaryThresholdPercentage(
        uint8 quorumNumber
    ) external view returns (uint8);

    /// @notice Gets the confirmation threshold percentage for a quorum
    function getQuorumConfirmationThresholdPercentage(
        uint8 quorumNumber
    ) external view returns (uint8);

    /// @notice Checks if a quorum is required
    function getIsQuorumRequired(
        uint8 quorumNumber
    ) external view returns (bool);
}