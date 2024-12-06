// SPDX-License-Identifier: MIT

pragma solidity ^0.8.9;

import {Merkle} from "eigenlayer-core/contracts/libraries/Merkle.sol";
import {BN254} from "eigenlayer-middleware/libraries/BN254.sol";
import {EigenDAHasher} from "../../src/libraries/EigenDAHasher.sol";
import {IEigenDAServiceManager} from "../../src/interfaces/IEigenDAServiceManager.sol";
import {BitmapUtils} from "eigenlayer-middleware/libraries/BitmapUtils.sol";
import "../../src/interfaces/IEigenDAStructs.sol";

/**
 * @title Library of functions to be used by smart contracts wanting to prove blobs on EigenDA and open KZG commitments.
 * @author Layr Labs, Inc.
 */
library EigenDARollupUtils {
    using BN254 for BN254.G1Point;
    
    /**
     * @notice Verifies the inclusion of a blob within a batch confirmed in `eigenDAServiceManager` and its trust assumptions
     * @param blobHeader the header of the blob containing relevant attributes of the blob
     * @param eigenDAServiceManager the contract in which the batch was confirmed 
     * @param blobVerificationProof the relevant data needed to prove inclusion of the blob and that the trust assumptions were as expected
     */
    function verifyBlob(
        BlobHeader memory blobHeader,
        IEigenDAServiceManager eigenDAServiceManager,
        BlobVerificationProof memory blobVerificationProof
    ) internal view {
        require(
            EigenDAHasher.hashBatchMetadata(blobVerificationProof.batchMetadata) 
                == eigenDAServiceManager.batchIdToBatchMetadataHash(blobVerificationProof.batchId),
            "EigenDARollupUtils.verifyBlob: batchMetadata does not match stored metadata"
        );

        require(
            Merkle.verifyInclusionKeccak(
                blobVerificationProof.inclusionProof, 
                blobVerificationProof.batchMetadata.batchHeader.blobHeadersRoot, 
                keccak256(abi.encodePacked(EigenDAHasher.hashBlobHeader(blobHeader))),
                blobVerificationProof.blobIndex
            ),
            "EigenDARollupUtils.verifyBlob: inclusion proof is invalid"
        );

        // bitmap of quorum numbers in all quorumBlobParams
        uint256 confirmedQuorumsBitmap;

        // require that the security param in each blob is met
        for (uint i = 0; i < blobHeader.quorumBlobParams.length; i++) {
            // make sure that the quorumIndex matches the given quorumNumber
            require(uint8(blobVerificationProof.batchMetadata.batchHeader.quorumNumbers[uint8(blobVerificationProof.quorumIndices[i])]) == blobHeader.quorumBlobParams[i].quorumNumber, 
                "EigenDARollupUtils.verifyBlob: quorumNumber does not match"
            );

            // make sure that the adversaryThresholdPercentage is less than the given confirmationThresholdPercentage
            require(blobHeader.quorumBlobParams[i].adversaryThresholdPercentage 
                < blobHeader.quorumBlobParams[i].confirmationThresholdPercentage, 
                "EigenDARollupUtils.verifyBlob: adversaryThresholdPercentage is not valid"
            );

            // make sure that the adversaryThresholdPercentage is at least the given quorumAdversaryThresholdPercentage
            uint8 _adversaryThresholdPercentage = getQuorumAdversaryThreshold(eigenDAServiceManager, blobHeader.quorumBlobParams[i].quorumNumber);
            if(_adversaryThresholdPercentage > 0){
                require(blobHeader.quorumBlobParams[i].adversaryThresholdPercentage >= _adversaryThresholdPercentage, 
                    "EigenDARollupUtils.verifyBlob: adversaryThresholdPercentage is not met"
                );
            }

            // make sure that the stake signed for is greater than the given confirmationThresholdPercentage
            require(uint8(blobVerificationProof.batchMetadata.batchHeader.signedStakeForQuorums[uint8(blobVerificationProof.quorumIndices[i])]) 
                >= blobHeader.quorumBlobParams[i].confirmationThresholdPercentage, 
                "EigenDARollupUtils.verifyBlob: confirmationThresholdPercentage is not met"
            );

            // mark confirmed quorum in the bitmap
            confirmedQuorumsBitmap = BitmapUtils.setBit(confirmedQuorumsBitmap, blobHeader.quorumBlobParams[i].quorumNumber);
        }

        // check that required quorums are a subset of the confirmed quorums
        require(
            BitmapUtils.isSubsetOf(
                BitmapUtils.orderedBytesArrayToBitmap(
                    eigenDAServiceManager.quorumNumbersRequired()
                ),
                confirmedQuorumsBitmap
            ),
            "EigenDARollupUtils.verifyBlob: required quorums are not a subset of the confirmed quorums"
        );
    }

    /**
     * @notice Verifies the inclusion of a blob within a batch confirmed in `eigenDAServiceManager` and its trust assumptions
     * @param blobHeaders the headers of the blobs containing relevant attributes of the blobs
     * @param eigenDAServiceManager the contract in which the batch was confirmed 
     * @param blobVerificationProofs the relevant data needed to prove inclusion of the blobs and that the trust assumptions were as expected
     */
    function verifyBlobs(
        BlobHeader[] memory blobHeaders,
        IEigenDAServiceManager eigenDAServiceManager,
        BlobVerificationProof[] memory blobVerificationProofs
    ) internal view {
        require(blobHeaders.length == blobVerificationProofs.length, "EigenDARollupUtils.verifyBlobs: blobHeaders and blobVerificationProofs must have the same length");

        bytes memory quorumAdversaryThresholdPercentages = eigenDAServiceManager.quorumAdversaryThresholdPercentages();
        uint256 quorumNumbersRequiredBitmap = BitmapUtils.orderedBytesArrayToBitmap(eigenDAServiceManager.quorumNumbersRequired());

        for (uint i = 0; i < blobHeaders.length; i++) {
            require(
                EigenDAHasher.hashBatchMetadata(blobVerificationProofs[i].batchMetadata) 
                    == eigenDAServiceManager.batchIdToBatchMetadataHash(blobVerificationProofs[i].batchId),
                "EigenDARollupUtils.verifyBlob: batchMetadata does not match stored metadata"
            );

            require(
                Merkle.verifyInclusionKeccak(
                    blobVerificationProofs[i].inclusionProof, 
                    blobVerificationProofs[i].batchMetadata.batchHeader.blobHeadersRoot, 
                    keccak256(abi.encodePacked(EigenDAHasher.hashBlobHeader(blobHeaders[i]))),
                    blobVerificationProofs[i].blobIndex
                ),
                "EigenDARollupUtils.verifyBlob: inclusion proof is invalid"
            );

            // bitmap of quorum numbers in all quorumBlobParams
            uint256 confirmedQuorumsBitmap;

            // require that the security param in each blob is met
            for (uint j = 0; j < blobHeaders[i].quorumBlobParams.length; j++) {
                // make sure that the quorumIndex matches the given quorumNumber
                require(uint8(blobVerificationProofs[i].batchMetadata.batchHeader.quorumNumbers[uint8(blobVerificationProofs[i].quorumIndices[i])]) == blobHeaders[i].quorumBlobParams[i].quorumNumber, 
                    "EigenDARollupUtils.verifyBlob: quorumNumber does not match"
                );

                // make sure that the adversaryThresholdPercentage is less than the given confirmationThresholdPercentage
                require(blobHeaders[i].quorumBlobParams[i].adversaryThresholdPercentage 
                    < blobHeaders[i].quorumBlobParams[i].confirmationThresholdPercentage, 
                    "EigenDARollupUtils.verifyBlob: adversaryThresholdPercentage is not valid"
                );

                // make sure that the adversaryThresholdPercentage is at least the given quorumAdversaryThresholdPercentage
                uint8 _adversaryThresholdPercentage = uint8(quorumAdversaryThresholdPercentages[blobHeaders[i].quorumBlobParams[j].quorumNumber]);
                if(_adversaryThresholdPercentage > 0){
                    require(blobHeaders[i].quorumBlobParams[j].adversaryThresholdPercentage >= _adversaryThresholdPercentage, 
                        "EigenDARollupUtils.verifyBlob: adversaryThresholdPercentage is not met"
                    );
                }

                // make sure that the stake signed for is greater than the given confirmationThresholdPercentage
                require(uint8(blobVerificationProofs[i].batchMetadata.batchHeader.signedStakeForQuorums[uint8(blobVerificationProofs[i].quorumIndices[j])]) 
                    >= blobHeaders[i].quorumBlobParams[j].confirmationThresholdPercentage, 
                    "EigenDARollupUtils.verifyBlob: confirmationThresholdPercentage is not met"
                );

                // mark confirmed quorum in the bitmap
                confirmedQuorumsBitmap = BitmapUtils.setBit(confirmedQuorumsBitmap, blobHeaders[i].quorumBlobParams[j].quorumNumber);
            }

            // check that required quorums are a subset of the confirmed quorums
            require(
                BitmapUtils.isSubsetOf(
                    quorumNumbersRequiredBitmap,
                    confirmedQuorumsBitmap
                ),
                "EigenDARollupUtils.verifyBlob: required quorums are not a subset of the confirmed quorums"
            );
        }
    }

    /**
     * @notice gets the adversary threshold percentage for a given quorum
     * @param eigenDAServiceManager the contract in which the batch was confirmed 
     * @param quorumNumber the quorum number to get the adversary threshold percentage for
     * @dev returns 0 if the quorumNumber is not found
     */
    function getQuorumAdversaryThreshold(
        IEigenDAServiceManager eigenDAServiceManager,
        uint256 quorumNumber
    ) internal view returns(uint8 adversaryThresholdPercentage) {
        if(eigenDAServiceManager.quorumAdversaryThresholdPercentages().length > quorumNumber){
            adversaryThresholdPercentage = uint8(eigenDAServiceManager.quorumAdversaryThresholdPercentages()[quorumNumber]);
        }
    }

    /**
     * @notice opens the KZG commitment at a point
     * @param point the point to evaluate the polynomial at
     * @param evaluation the evaluation of the polynomial at the point
     * @param tau the power of tau
     * @param commitment the commitment to the polynomial
     * @param proof the proof of the commitment
     */
    function openCommitment(
        uint256 point, 
        uint256 evaluation,
        BN254.G1Point memory tau, 
        BN254.G1Point memory commitment, 
        BN254.G2Point memory proof 
    ) internal view returns(bool) {
        BN254.G1Point memory negGeneratorG1 = BN254.generatorG1().negate();

        //e([s]_1 - w[1]_1, [pi(x)]_2) = e([p(x)]_1 - p(w)[1]_1, [1]_2)
        return BN254.pairing(
            tau.plus(negGeneratorG1.scalar_mul(point)), 
            proof, 
            commitment.plus(negGeneratorG1.scalar_mul(evaluation)), 
            BN254.negGeneratorG2()
        );
    }
}