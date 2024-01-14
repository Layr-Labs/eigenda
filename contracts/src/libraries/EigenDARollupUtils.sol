// SPDX-License-Identifier: MIT

pragma solidity ^0.8.9;

import {Merkle} from "eigenlayer-core/contracts/libraries/Merkle.sol";
import {BN254} from "eigenlayer-middleware/libraries/BN254.sol";
import {EigenDAHasher} from "./EigenDAHasher.sol";
import {IEigenDAServiceManager} from "../interfaces/IEigenDAServiceManager.sol";

/**
 * @title Library of functions to be used by smart contracts wanting to prove blobs on EigenDA and open KZG commitments.
 * @author Layr Labs, Inc.
 */
library EigenDARollupUtils {
    using BN254 for BN254.G1Point;

    // STRUCTS
    struct BlobVerificationProof {
        uint32 batchId;
        uint8 blobIndex;
        IEigenDAServiceManager.BatchMetadata batchMetadata;
        bytes inclusionProof;
        bytes quorumThresholdIndexes;
    }
    
    /**
     * @notice Verifies the inclusion of a blob within a batch confirmed in `eigenDAServiceManager` and its trust assumptions
     * @param blobHeader the header of the blob containing relevant attributes of the blob
     * @param eigenDAServiceManager the contract in which the batch was confirmed 
     * @param blobVerificationProof the relevant data needed to prove inclusion of the blob and that the trust assumptions were as expected
     */
    function verifyBlob(
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        IEigenDAServiceManager eigenDAServiceManager,
        BlobVerificationProof calldata blobVerificationProof
    ) external view {
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

        // require that the security param in each blob is met
        for (uint i = 0; i < blobHeader.quorumBlobParams.length; i++) {
            // make sure that the quorumIndex matches the given quorumNumber
            require(uint8(blobVerificationProof.batchMetadata.batchHeader.quorumNumbers[uint8(blobVerificationProof.quorumThresholdIndexes[i])]) == blobHeader.quorumBlobParams[i].quorumNumber, 
                "EigenDARollupUtils.verifyBlob: quorumNumber does not match"
            );

            // make sure that the adversaryThresholdPercentage is less than the given quorumThresholdPercentage
            require(blobHeader.quorumBlobParams[i].adversaryThresholdPercentage 
                < blobHeader.quorumBlobParams[i].quorumThresholdPercentage, 
                "EigenDARollupUtils.verifyBlob: adversaryThresholdPercentage is not valid"
            );

            // make sure that the stake signed for is greater than the given quorumThresholdPercentage
            require(uint8(blobVerificationProof.batchMetadata.batchHeader.quorumThresholdPercentages[uint8(blobVerificationProof.quorumThresholdIndexes[i])]) 
                >= blobHeader.quorumBlobParams[i].quorumThresholdPercentage, 
                "EigenDARollupUtils.verifyBlob: quorumThresholdPercentage is not met"
            );

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
