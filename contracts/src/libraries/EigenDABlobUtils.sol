// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.9;

import "@eigenlayer-core/contracts/libraries/Merkle.sol";
import "./EigenDAHasher.sol";
import "../interfaces/IEigenDAServiceManager.sol";

/**
 * @title Library of functions to be used by smart contracts wanting to prove blobs on EigenDA.
 * @author Layr Labs, Inc.
 */
library EigenDABlobUtils {
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
            "EigenDABlobUtils.verifyBlob: batchMetadata does not match stored metadata"
        );

        require(
            Merkle.verifyInclusionKeccak(
                blobVerificationProof.inclusionProof, 
                blobVerificationProof.batchMetadata.batchHeader.blobHeadersRoot, 
                keccak256(abi.encodePacked(EigenDAHasher.hashBlobHeader(blobHeader))),
                blobVerificationProof.blobIndex
            ),
            "EigenDABlobUtils.verifyBlob: inclusion proof is invalid"
        );

        // require that the security param in each blob is met
        for (uint i = 0; i < blobHeader.quorumBlobParams.length; i++) {
            // make sure that the quorumIndex matches the given quorumNumber
            require(uint8(blobVerificationProof.batchMetadata.batchHeader.quorumNumbers[uint8(blobVerificationProof.quorumThresholdIndexes[i])]) == blobHeader.quorumBlobParams[i].quorumNumber, 
                "EigenDABlobUtils.verifyBlob: quorumNumber does not match"
            );

            // make sure that the adversaryThresholdPercentage is less than the given quorumThresholdPercentage
            require(blobHeader.quorumBlobParams[i].adversaryThresholdPercentage 
                < blobHeader.quorumBlobParams[i].quorumThresholdPercentage, 
                "EigenDABlobUtils.verifyBlob: adversaryThresholdPercentage is not valid"
            );

            // make sure that the stake signed for is greater than the given quorumThresholdPercentage
            require(uint8(blobVerificationProof.batchMetadata.batchHeader.quorumThresholdPercentages[uint8(blobVerificationProof.quorumThresholdIndexes[i])]) 
                >= blobHeader.quorumBlobParams[i].quorumThresholdPercentage, 
                "EigenDABlobUtils.verifyBlob: quorumThresholdPercentage is not met"
            );

        }
    }
}
