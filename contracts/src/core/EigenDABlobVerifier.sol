// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {Merkle} from "eigenlayer-core/contracts/libraries/Merkle.sol";
import {BitmapUtils} from "eigenlayer-middleware/libraries/BitmapUtils.sol";
import {EigenDAHasher} from "../libraries/EigenDAHasher.sol";
import {IEigenDAServiceManager} from "../interfaces/IEigenDAServiceManager.sol";
import {IEigenDABlobVerifier} from "../interfaces/IEigenDABlobVerifier.sol";
import {IEigenDASignatureVerifier} from "../interfaces/IEigenDASignatureVerifier.sol";
import {IEigenDABatchMetadataStorage} from "../interfaces/IEigenDABatchMetadataStorage.sol";

abstract contract EigenDABlobVerifier is IEigenDABlobVerifier {
    using EigenDAHasher for IEigenDAServiceManager.BatchHeader;
    using EigenDAHasher for IEigenDAServiceManager.ReducedBatchHeader;

    /// @notice The denominator of a threshold percentage
    uint256 public constant THRESHOLD_DENOMINATOR = 100;

    /**
     * @notice The quorum adversary threshold percentages stored as an ordered bytes array
     * this is the percentage of the total stake that must be adversarial to consider a blob invalid.
     * The first byte is the threshold for quorum 0, the second byte is the threshold for quorum 1, etc.
     */
    bytes public constant quorumAdversaryThresholdPercentages = hex"2121";

    /**
     * @notice The quorum confirmation threshold percentages stored as an ordered bytes array
     * this is the percentage of the total stake needed to confirm a blob.
     * The first byte is the threshold for quorum 0, the second byte is the threshold for quorum 1, etc.
     */
    bytes public constant quorumConfirmationThresholdPercentages = hex"3737";

    /**
     * @notice The quorum numbers required for confirmation stored as an ordered bytes array
     * these quorum numbers have respective canonical thresholds in the
     * quorumConfirmationThresholdPercentages and quorumAdversaryThresholdPercentages above.
     */
    bytes public constant quorumNumbersRequired = hex"0001";

    /**
     * @notice Verifies a the blob is valid for the required quorums
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob verification proof to verify the blob against
     */
    function verifyBlob(
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof
    ) external view virtual {
        _verifyBlobForQuorums(
            IEigenDABatchMetadataStorage(address(this)), 
            blobHeader, 
            blobVerificationProof, 
            quorumNumbersRequired
        );
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
    ) external view virtual {
        _verifyBlobForQuorums(
            IEigenDABatchMetadataStorage(address(this)), 
            blobHeader, 
            blobVerificationProof, 
            bytes.concat(quorumNumbersRequired, additionalQuorumNumbersRequired)
        );
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
    ) external view virtual {
        _verifyPreconfirmationForQuorums(
            IEigenDASignatureVerifier(address(this)), 
            miniBatchHeader, 
            blobHeader,
            nonSignerStakesAndSignature, 
            quorumNumbersRequired
        );
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
    ) external view virtual {
        _verifyPreconfirmationForQuorums(
            IEigenDASignatureVerifier(address(this)), 
            miniBatchHeader, 
            blobHeader,
            nonSignerStakesAndSignature, 
            bytes.concat(quorumNumbersRequired, additionalQuorumNumbersRequired)
        );
    }

    /// @notice Verifies that a blob is valid for the a set of quorums
    function _verifyBlobForQuorums(
        IEigenDABatchMetadataStorage batchMetadataStorage,
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof,
        bytes memory requiredQuorumNumbers
    ) internal view virtual {
        require(
            EigenDAHasher.hashBatchMetadata(blobVerificationProof.batchMetadata) 
                == IEigenDABatchMetadataStorage(batchMetadataStorage).batchIdToBatchMetadataHash(blobVerificationProof.batchId),
            "EigenDABlobVerifier._verifyBlobForQuorums: batchMetadata does not match stored metadata"
        );

        require(
            Merkle.verifyInclusionKeccak(
                blobVerificationProof.inclusionProof, 
                blobVerificationProof.batchMetadata.batchHeader.blobHeadersRoot, 
                keccak256(abi.encodePacked(EigenDAHasher.hashBlobHeader(blobHeader))),
                blobVerificationProof.blobIndex
            ),
            "EigenDABlobVerifier._verifyBlobForQuorums: inclusion proof is invalid"
        );

        uint256 confirmedQuorumsBitmap;

        for (uint i = 0; i < blobHeader.quorumBlobParams.length; i++) {

            require(
                uint8(blobVerificationProof.batchMetadata.batchHeader.quorumNumbers[uint8(blobVerificationProof.quorumIndices[i])]) == 
                blobHeader.quorumBlobParams[i].quorumNumber, 
                "EigenDABlobVerifier._verifyBlobForQuorums: quorumNumber does not match"
            );

            require(
                blobHeader.quorumBlobParams[i].adversaryThresholdPercentage < 
                blobHeader.quorumBlobParams[i].confirmationThresholdPercentage, 
                "EigenDABlobVerifier._verifyBlobForQuorums: threshold percentages are not valid"
            );

            uint8 _adversaryThresholdPercentage = getQuorumAdversaryThresholdPercentage(blobHeader.quorumBlobParams[i].quorumNumber);
            if(_adversaryThresholdPercentage > 0){
                require(
                    blobHeader.quorumBlobParams[i].adversaryThresholdPercentage >= 
                    _adversaryThresholdPercentage, 
                    "EigenDABlobVerifier._verifyBlobForQuorums: adversaryThresholdPercentage is not met"
                );
            }

            require(
                uint8(blobVerificationProof.batchMetadata.batchHeader.signedStakeForQuorums[uint8(blobVerificationProof.quorumIndices[i])]) >= 
                blobHeader.quorumBlobParams[i].confirmationThresholdPercentage, 
                "EigenDABlobVerifier._verifyBlobForQuorums: confirmationThresholdPercentage is not met"
            );

            confirmedQuorumsBitmap = BitmapUtils.setBit(confirmedQuorumsBitmap, blobHeader.quorumBlobParams[i].quorumNumber);
        }

        require(
            BitmapUtils.isSubsetOf(
                BitmapUtils.orderedBytesArrayToBitmap(requiredQuorumNumbers),
                confirmedQuorumsBitmap
            ),
            "EigenDABlobVerifier._verifyBlobForQuorums: required quorums are not a subset of the confirmed quorums"
        );
    }

    /// @notice Verifies that a blob preconfirmation is valid for a set of quorums
    function _verifyPreconfirmationForQuorums(
        IEigenDASignatureVerifier signatureVerifier,
        IEigenDAServiceManager.BatchHeader calldata miniBatchHeader,
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        IEigenDASignatureVerifier.NonSignerStakesAndSignature memory nonSignerStakesAndSignature,
        bytes memory requiredQuorumNumbers
    ) internal view virtual {
        require(
            miniBatchHeader.quorumNumbers.length == miniBatchHeader.signedStakeForQuorums.length,
            "EigenDABlobVerifier._verifyPreconfirmationForQuorums: quorumNumbers and signedStakeForQuorums must be of the same length"
        );

        require(
            miniBatchHeader.blobHeadersRoot == keccak256(abi.encodePacked(EigenDAHasher.hashBlobHeader(blobHeader))),
            "EigenDABlobVerifier._verifyPreconfirmationForQuorums: blobHeadersRoot does not match blobHeader"
        );

        bytes32 reducedBatchHeaderHash = miniBatchHeader.hashBatchHeaderToReducedBatchHeader();

        (
            IEigenDASignatureVerifier.QuorumStakeTotals memory quorumStakeTotals,
            bytes32 signatoryRecordHash
        ) = IEigenDASignatureVerifier(signatureVerifier).checkSignatures(
            reducedBatchHeaderHash, 
            miniBatchHeader.quorumNumbers, 
            miniBatchHeader.referenceBlockNumber, 
            nonSignerStakesAndSignature
        );

        for (uint i = 0; i < miniBatchHeader.signedStakeForQuorums.length; i++) {
            require(
                quorumStakeTotals.signedStakeForQuorum[i] * THRESHOLD_DENOMINATOR >= 
                quorumStakeTotals.totalStakeForQuorum[i] * getQuorumConfirmationThresholdPercentage(uint8(miniBatchHeader.quorumNumbers[i])),
                "EigenDABlobVerifier._verifyPreconfirmationForQuorums: signatories do not own at least confirmation threshold percentage of a quorum"
            );
        }

        uint256 confirmedQuorumsBitmap;

        for (uint i = 0; i < blobHeader.quorumBlobParams.length; i++) {

            require(
                blobHeader.quorumBlobParams[i].adversaryThresholdPercentage < 
                blobHeader.quorumBlobParams[i].confirmationThresholdPercentage, 
                "EigenDABlobVerifier._verifyPreconfirmationForQuorums: threshold percentages are not valid"
            );

            uint8 _adversaryThresholdPercentage = getQuorumAdversaryThresholdPercentage(blobHeader.quorumBlobParams[i].quorumNumber);
            if(_adversaryThresholdPercentage > 0){
                require(
                    blobHeader.quorumBlobParams[i].adversaryThresholdPercentage >= 
                    _adversaryThresholdPercentage, 
                    "EigenDABlobVerifier._verifyPreconfirmationForQuorums: adversaryThresholdPercentage is not met"
                );
            }

            uint8 _confirmationThresholdPercentage = getQuorumConfirmationThresholdPercentage(blobHeader.quorumBlobParams[i].quorumNumber);
            if(_confirmationThresholdPercentage > 0){
                require(blobHeader.quorumBlobParams[i].confirmationThresholdPercentage >= _confirmationThresholdPercentage, 
                    "EigenDABlobVerifier._verifyPreconfirmationForQuorums: confirmationThresholdPercentage is not met"
                );
            }

            confirmedQuorumsBitmap = BitmapUtils.setBit(confirmedQuorumsBitmap, blobHeader.quorumBlobParams[i].quorumNumber);
        }

        require(
            BitmapUtils.isSubsetOf(
                BitmapUtils.orderedBytesArrayToBitmap(requiredQuorumNumbers),
                confirmedQuorumsBitmap
            ),
            "EigenDABlobVerifier._verifyBlobForQuorums: required quorums are not a subset of the confirmed quorums"
        );
    }

    /// @notice Gets the adversary threshold percentage for a quorum
    function getQuorumAdversaryThresholdPercentage(
        uint8 quorumNumber
    ) public view virtual returns (uint8 adversaryThresholdPercentage) {
        if(quorumAdversaryThresholdPercentages.length > quorumNumber){
            adversaryThresholdPercentage = uint8(quorumAdversaryThresholdPercentages[quorumNumber]);
        }
    }

    /// @notice Gets the confirmation threshold percentage for a quorum
    function getQuorumConfirmationThresholdPercentage(
        uint8 quorumNumber
    ) public view virtual returns (uint8 confirmationThresholdPercentage) {
        if(quorumConfirmationThresholdPercentages.length > quorumNumber){
            confirmationThresholdPercentage = uint8(quorumConfirmationThresholdPercentages[quorumNumber]);
        }
    }

    /// @notice Checks if a quorum is required
    function getIsQuorumRequired(
        uint8 quorumNumber
    ) public view virtual returns (bool) {
        uint256 quorumBitmap = BitmapUtils.setBit(0, quorumNumber);
        return (quorumBitmap & BitmapUtils.orderedBytesArrayToBitmap(quorumNumbersRequired) == quorumBitmap);
    }
}