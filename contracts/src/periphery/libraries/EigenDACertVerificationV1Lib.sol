// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDABatchMetadataStorage} from "src/interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDAThresholdRegistry} from "src/interfaces/IEigenDAThresholdRegistry.sol";
import {EigenDAHasher} from "src/libraries/EigenDAHasher.sol";
import {BN254} from "lib/eigenlayer-middleware/src/libraries/BN254.sol";
import {Merkle} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/libraries/Merkle.sol";
import {BitmapUtils} from "lib/eigenlayer-middleware/src/libraries/BitmapUtils.sol";

import {BlobHeader, BlobVerificationProof, QuorumBlobParam} from "src/interfaces/IEigenDAStructs.sol";

/**
 * @title EigenDACertVerificationV1Lib - EigenDA V1 certificate verification library
 * @author Layr Labs, Inc.
 * @notice Library of functions for verifying EigenDA V1 certificates
 */
library EigenDACertVerificationV1Lib {
    using BN254 for BN254.G1Point;

    /// @notice Denominator used for threshold percentage calculations (100 for percentages)
    uint256 internal constant THRESHOLD_DENOMINATOR = 100;

    /// @notice Thrown when the batch metadata does not match stored metadata
    /// @param actualHash The hash computed from the provided metadata
    /// @param expectedHash The hash stored in the metadata storage
    error BatchMetadataMismatch(bytes32 actualHash, bytes32 expectedHash);

    /// @notice Thrown when the inclusion proof is invalid
    /// @param blobIndex The index of the blob in the batch
    /// @param blobHash The hash of the blob
    /// @param rootHash The root hash of the merkle tree
    error InvalidInclusionProof(uint256 blobIndex, bytes32 blobHash, bytes32 rootHash);

    /// @notice Thrown when the quorum number does not match
    /// @param expected The expected quorum number
    /// @param actual The actual quorum number
    error QuorumNumberMismatch(uint8 expected, uint8 actual);

    /// @notice Thrown when threshold percentages are not valid
    /// @param confirmationThreshold The confirmation threshold percentage
    /// @param adversaryThreshold The adversary threshold percentage
    error InvalidThresholdPercentages(uint8 confirmationThreshold, uint8 adversaryThreshold);

    /// @notice Thrown when signed stake percentage is not met
    /// @param quorumNumber The quorum number
    /// @param requiredThreshold The required threshold percentage
    /// @param actualThreshold The actual threshold percentage
    error StakeThresholdNotMet(uint8 quorumNumber, uint8 requiredThreshold, uint8 actualThreshold);

    /// @notice Thrown when confirmation threshold percentage is not met
    /// @param quorumNumber The quorum number
    /// @param requiredThreshold The required threshold percentage
    /// @param actualThreshold The actual threshold percentage
    error ConfirmationThresholdNotMet(uint8 quorumNumber, uint8 requiredThreshold, uint8 actualThreshold);

    /// @notice Thrown when required quorums are not a subset of confirmed quorums
    /// @param requiredQuorumsBitmap The bitmap of required quorums
    /// @param confirmedQuorumsBitmap The bitmap of confirmed quorums
    error RequiredQuorumsNotSubset(uint256 requiredQuorumsBitmap, uint256 confirmedQuorumsBitmap);

    /// @notice Thrown when security assumptions are not met
    /// @param errParams Additional error parameters
    error SecurityAssumptionsNotMet(bytes errParams);

    /// @notice Thrown when blob quorums are not a subset of confirmed quorums
    /// @param blobQuorumsBitmap The bitmap of blob quorums
    /// @param confirmedQuorumsBitmap The bitmap of confirmed quorums
    error BlobQuorumsNotSubset(uint256 blobQuorumsBitmap, uint256 confirmedQuorumsBitmap);

    /// @notice Thrown when there is a length mismatch
    /// @param expected The expected length
    /// @param actual The actual length
    error LengthMismatch(uint256 expected, uint256 actual);

    /// @notice Thrown when a relay key is not set
    /// @param relayKey The relay key that was not set
    error RelayKeyNotSet(uint32 relayKey);

    /// @notice Error codes for certificate verification results
    enum ErrorCode {
        SUCCESS, // Verification succeeded
        BATCH_METADATA_MISMATCH, // Batch metadata hash doesn't match stored hash
        INVALID_INCLUSION_PROOF, // Merkle inclusion proof is invalid
        QUORUM_NUMBER_MISMATCH, // Quorum number doesn't match expected value
        INVALID_THRESHOLD_PERCENTAGES, // Threshold percentages are invalid
        CONFIRMATION_THRESHOLD_NOT_MET, // Confirmation threshold not met
        STAKE_THRESHOLD_NOT_MET, // Stake threshold not met
        REQUIRED_QUORUMS_NOT_SUBSET, // Required quorums not a subset of confirmed quorums
        SECURITY_ASSUMPTIONS_NOT_MET, // Security assumptions not met
        BLOB_QUORUMS_NOT_SUBSET, // Blob quorums not a subset of confirmed quorums
        RELAY_KEY_NOT_SET // Relay key not set

    }

    /**
     * @notice Verifies that batch metadata matches the stored metadata.
     * @param storedBatchMetadataHash The batch metadata hash stored in the service manager.
     * @param blobVerificationProof Pointer to the blob verification proof in calldata.
     * @return err Error code (SUCCESS if verification succeeded).
     * @return errParams Additional error parameters.
     */
    function verifyBatchMetadata(bytes32 storedBatchMetadataHash, BlobVerificationProof calldata blobVerificationProof)
        internal
        pure
        returns (ErrorCode err, bytes memory errParams)
    {
        bytes32 batchMetadataHash = EigenDAHasher.hashBatchMetadata(blobVerificationProof.batchMetadata);

        if (batchMetadataHash == storedBatchMetadataHash) {
            return (ErrorCode.SUCCESS, "");
        } else {
            return (ErrorCode.BATCH_METADATA_MISMATCH, abi.encode(batchMetadataHash, storedBatchMetadataHash));
        }
    }

    /**
     * @notice Verifies blob inclusion in the batch using Merkle proof.
     * @param blobHeader Pointer to the blob header in calldata.
     * @param blobVerificationProof Pointer to the blob verification proof in calldata.
     * @return err Error code (SUCCESS if verification succeeded).
     * @return errParams Additional error parameters.
     */
    function verifyBlobInclusion(BlobHeader calldata blobHeader, BlobVerificationProof calldata blobVerificationProof)
        internal
        pure
        returns (ErrorCode err, bytes memory errParams)
    {
        bytes32 blobHeaderHash = EigenDAHasher.hashBlobHeader(blobHeader);
        bytes32 encodedBlobHash = keccak256(abi.encodePacked(blobHeaderHash));
        bytes32 rootHash = blobVerificationProof.batchMetadata.batchHeader.blobHeadersRoot;

        bool isValid = Merkle.verifyInclusionKeccak(
            blobVerificationProof.inclusionProof, rootHash, encodedBlobHash, blobVerificationProof.blobIndex
        );

        if (isValid) {
            return (ErrorCode.SUCCESS, "");
        } else {
            return (
                ErrorCode.INVALID_INCLUSION_PROOF,
                abi.encode(blobVerificationProof.blobIndex, encodedBlobHash, rootHash)
            );
        }
    }

    /**
     * @notice Verifies a single quorum parameter.
     * @param quorumConfirmationThresholdPercentages The quorum confirmation threshold percentages.
     * @param blobHeader Pointer to the blob header in calldata.
     * @param blobVerificationProof Pointer to the blob verification proof in calldata.
     * @param quorumIndex The index of the quorum to verify.
     * @param verificationQuorumIndex The index in the verification proof.
     * @return err Error code (SUCCESS if verification succeeded).
     * @return errParams Additional error parameters.
     * @return quorumNumber The verified quorum number.
     */
    function verifyQuorumParameter(
        bytes memory quorumConfirmationThresholdPercentages,
        BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof,
        uint256 quorumIndex,
        uint256 verificationQuorumIndex
    ) internal pure returns (ErrorCode err, bytes memory errParams, uint8 quorumNumber) {
        // Access quorumBlobParams directly from calldata
        QuorumBlobParam calldata quorumParams = blobHeader.quorumBlobParams[quorumIndex];
        quorumNumber = quorumParams.quorumNumber;

        // Get quorum indices from calldata
        uint8 batchQuorumIndex = uint8(blobVerificationProof.quorumIndices[verificationQuorumIndex]);
        uint8 batchQuorumNumber = uint8(blobVerificationProof.batchMetadata.batchHeader.quorumNumbers[batchQuorumIndex]);

        // Check quorum number matches
        if (batchQuorumNumber != quorumNumber) {
            return (ErrorCode.QUORUM_NUMBER_MISMATCH, abi.encode(quorumNumber, batchQuorumNumber), 0);
        }

        // Check threshold percentages are valid
        uint8 confirmationThreshold = quorumParams.confirmationThresholdPercentage;
        uint8 adversaryThreshold = quorumParams.adversaryThresholdPercentage;
        if (confirmationThreshold <= adversaryThreshold) {
            return (ErrorCode.INVALID_THRESHOLD_PERCENTAGES, abi.encode(confirmationThreshold, adversaryThreshold), 0);
        }

        // Check confirmation threshold percentage meets minimum requirement
        uint8 requiredThreshold = uint8(quorumConfirmationThresholdPercentages[quorumNumber]);
        if (confirmationThreshold < requiredThreshold) {
            return (
                ErrorCode.CONFIRMATION_THRESHOLD_NOT_MET,
                abi.encode(quorumNumber, requiredThreshold, confirmationThreshold),
                0
            );
        }

        // Check signed stake meets the confirmation threshold
        uint8 signedStake =
            uint8(blobVerificationProof.batchMetadata.batchHeader.signedStakeForQuorums[batchQuorumIndex]);
        if (signedStake < confirmationThreshold) {
            return (ErrorCode.STAKE_THRESHOLD_NOT_MET, abi.encode(quorumNumber, confirmationThreshold, signedStake), 0);
        }

        return (ErrorCode.SUCCESS, "", quorumNumber);
    }

    /**
     * @notice Verifies all quorum parameters for a blob and builds a bitmap of confirmed quorums.
     * @param quorumConfirmationThresholdPercentages The quorum confirmation threshold percentages.
     * @param blobHeader Pointer to the blob header in calldata.
     * @param blobVerificationProof Pointer to the blob verification proof in calldata.
     * @return err Error code (SUCCESS if verification succeeded).
     * @return errParams Additional error parameters.
     * @return confirmedQuorumsBitmap The bitmap of confirmed quorums.
     */
    function verifyQuorumParameters(
        bytes memory quorumConfirmationThresholdPercentages,
        BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof
    ) internal pure returns (ErrorCode err, bytes memory errParams, uint256 confirmedQuorumsBitmap) {
        confirmedQuorumsBitmap = 0;
        uint256 quorumCount = blobHeader.quorumBlobParams.length;

        for (uint256 i = 0; i < quorumCount; i++) {
            (ErrorCode quorumErr, bytes memory quorumErrParams, uint8 quorumNumber) =
                verifyQuorumParameter(quorumConfirmationThresholdPercentages, blobHeader, blobVerificationProof, i, i);

            if (quorumErr != ErrorCode.SUCCESS) {
                return (quorumErr, quorumErrParams, 0);
            }

            // Add to confirmed quorums bitmap
            confirmedQuorumsBitmap = BitmapUtils.setBit(confirmedQuorumsBitmap, quorumNumber);
        }

        return (ErrorCode.SUCCESS, "", confirmedQuorumsBitmap);
    }

    /**
     * @notice Verifies that required quorums are a subset of confirmed quorums.
     * @param requiredQuorumNumbers Pointer to the required quorum numbers in calldata.
     * @param confirmedQuorumsBitmap The bitmap of confirmed quorums.
     * @return err Error code (SUCCESS if verification succeeded).
     * @return errParams Additional error parameters.
     */
    function verifyRequiredQuorumsSubset(bytes memory requiredQuorumNumbers, uint256 confirmedQuorumsBitmap)
        internal
        pure
        returns (ErrorCode err, bytes memory errParams)
    {
        uint256 requiredQuorumsBitmap = BitmapUtils.orderedBytesArrayToBitmap(requiredQuorumNumbers);

        if (BitmapUtils.isSubsetOf(requiredQuorumsBitmap, confirmedQuorumsBitmap)) {
            return (ErrorCode.SUCCESS, "");
        } else {
            return (ErrorCode.REQUIRED_QUORUMS_NOT_SUBSET, abi.encode(requiredQuorumsBitmap, confirmedQuorumsBitmap));
        }
    }

    /**
     * @notice Verifies a complete blob certificate in a single call.
     * @param quorumConfirmationThresholdPercentages The quorum confirmation threshold percentages.
     * @param storedBatchMetadataHash The batch metadata hash stored in the service manager.
     * @param blobHeader The blob header to verify.
     * @param blobVerificationProof The blob cert verification proof to verify.
     * @param requiredQuorumNumbers The required quorum numbers.
     * @return err Error code (SUCCESS if verification succeeded).
     * @return errParams Additional error parameters.
     */
    function verifyDACert(
        bytes memory quorumConfirmationThresholdPercentages,
        bytes32 storedBatchMetadataHash,
        BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof,
        bytes memory requiredQuorumNumbers
    ) internal pure returns (ErrorCode err, bytes memory errParams) {
        // Verify batch metadata
        (err, errParams) = verifyBatchMetadata(storedBatchMetadataHash, blobVerificationProof);
        if (err != ErrorCode.SUCCESS) {
            return (err, errParams);
        }

        // Verify blob inclusion
        (err, errParams) = verifyBlobInclusion(blobHeader, blobVerificationProof);
        if (err != ErrorCode.SUCCESS) {
            return (err, errParams);
        }

        // Verify quorum parameters
        uint256 confirmedQuorumsBitmap;
        (err, errParams, confirmedQuorumsBitmap) =
            verifyQuorumParameters(quorumConfirmationThresholdPercentages, blobHeader, blobVerificationProof);
        if (err != ErrorCode.SUCCESS) {
            return (err, errParams);
        }

        // Verify required quorums are a subset of confirmed quorums
        return verifyRequiredQuorumsSubset(requiredQuorumNumbers, confirmedQuorumsBitmap);
    }

    /**
     * @notice Handles error codes by reverting with appropriate custom errors
     * @param err The error code
     * @param errParams The error parameters
     */
    function revertOnError(ErrorCode err, bytes memory errParams) internal pure {
        if (err == ErrorCode.SUCCESS) {
            return; // No error to handle
        }

        if (err == ErrorCode.BATCH_METADATA_MISMATCH) {
            (bytes32 actualHash, bytes32 expectedHash) = abi.decode(errParams, (bytes32, bytes32));
            revert BatchMetadataMismatch(actualHash, expectedHash);
        } else if (err == ErrorCode.INVALID_INCLUSION_PROOF) {
            (uint256 blobIndex, bytes32 blobHash, bytes32 rootHash) = abi.decode(errParams, (uint256, bytes32, bytes32));
            revert InvalidInclusionProof(blobIndex, blobHash, rootHash);
        } else if (err == ErrorCode.QUORUM_NUMBER_MISMATCH) {
            (uint8 expected, uint8 actual) = abi.decode(errParams, (uint8, uint8));
            revert QuorumNumberMismatch(expected, actual);
        } else if (err == ErrorCode.INVALID_THRESHOLD_PERCENTAGES) {
            (uint8 confirmationThreshold, uint8 adversaryThreshold) = abi.decode(errParams, (uint8, uint8));
            revert InvalidThresholdPercentages(confirmationThreshold, adversaryThreshold);
        } else if (err == ErrorCode.CONFIRMATION_THRESHOLD_NOT_MET) {
            (uint8 quorumNumber, uint8 requiredThreshold, uint8 actualThreshold) =
                abi.decode(errParams, (uint8, uint8, uint8));
            revert ConfirmationThresholdNotMet(quorumNumber, requiredThreshold, actualThreshold);
        } else if (err == ErrorCode.STAKE_THRESHOLD_NOT_MET) {
            (uint8 quorumNumber, uint8 requiredThreshold, uint8 actualThreshold) =
                abi.decode(errParams, (uint8, uint8, uint8));
            revert StakeThresholdNotMet(quorumNumber, requiredThreshold, actualThreshold);
        } else if (err == ErrorCode.REQUIRED_QUORUMS_NOT_SUBSET) {
            (uint256 requiredQuorumsBitmap, uint256 confirmedQuorumsBitmap) = abi.decode(errParams, (uint256, uint256));
            revert RequiredQuorumsNotSubset(requiredQuorumsBitmap, confirmedQuorumsBitmap);
        } else if (err == ErrorCode.RELAY_KEY_NOT_SET) {
            uint32 relayKey = abi.decode(errParams, (uint32));
            revert RelayKeyNotSet(relayKey);
        } else if (err == ErrorCode.SECURITY_ASSUMPTIONS_NOT_MET) {
            revert SecurityAssumptionsNotMet(errParams);
        } else if (err == ErrorCode.BLOB_QUORUMS_NOT_SUBSET) {
            (uint256 blobQuorumsBitmap, uint256 confirmedQuorumsBitmap) = abi.decode(errParams, (uint256, uint256));
            revert BlobQuorumsNotSubset(blobQuorumsBitmap, confirmedQuorumsBitmap);
        } else {
            revert("Unknown error code");
        }
    }
}
