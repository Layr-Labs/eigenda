// SPDX-License-Identifier: MIT

pragma solidity ^0.8.9;

import {Merkle} from "eigenlayer-core/contracts/libraries/Merkle.sol";
import {BN254} from "eigenlayer-middleware/libraries/BN254.sol";
import {EigenDAHasher} from "./EigenDAHasher.sol";
import {IEigenDAServiceManager} from "../interfaces/IEigenDAServiceManager.sol";
import {BitmapUtils} from "eigenlayer-middleware/libraries/BitmapUtils.sol";
import {IEigenDABatchMetadataStorage} from "../interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDAThresholdRegistry} from "../interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "../interfaces/IEigenDASignatureVerifier.sol";

/**
 * @title Library of functions to be used by smart contracts wanting to verify submissions of blobs on EigenDA.
 * @author Layr Labs, Inc.
 */
library EigenDABlobVerificationUtils {
    using BN254 for BN254.G1Point;

    struct BlobVerificationProof {
        uint32 batchId;
        uint32 blobIndex;
        IEigenDAServiceManager.BatchMetadata batchMetadata;
        bytes inclusionProof;
        bytes quorumIndices;
    }

    struct SignedCertificate {
        BlobCertificate blobCertificate;
        Attestation nonSignerStakesAndSignature;
    }

    struct BlobCertificate {
        bytes blobKey;
        IEigenDAServiceManager.BlobHeader blobHeader;
        uint32 referenceBlockNumber;
        string[] relayKeys;
    }

    struct Attestation {
        uint32[] nonSignerQuorumBitmapIndices;
        BN254.G1Point[] nonSignerPubkeys;
        BN254.G1Point[] quorumApks;
        BN254.G2Point apkG2;
        BN254.G1Point sigma;
        uint32[] quorumApkIndices;
        uint32[] totalStakeIndices;
        uint32[][] nonSignerStakeIndices;
    }
    
    function _verifyBlobV1ForQuorums(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDABatchMetadataStorage batchMetadataStorage,
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof,
        bytes memory requiredQuorumNumbers
    ) internal view {
        require(
            EigenDAHasher.hashBatchMetadata(blobVerificationProof.batchMetadata) ==
            IEigenDABatchMetadataStorage(batchMetadataStorage).batchIdToBatchMetadataHash(blobVerificationProof.batchId),
            "EigenDABlobVerificationUtils._verifyBlobForQuorums: batchMetadata does not match stored metadata"
        );

        require(
            Merkle.verifyInclusionKeccak(
                blobVerificationProof.inclusionProof, 
                blobVerificationProof.batchMetadata.batchHeader.blobHeadersRoot, 
                keccak256(abi.encodePacked(EigenDAHasher.hashBlobHeader(blobHeader))),
                blobVerificationProof.blobIndex
            ),
            "EigenDABlobVerificationUtils._verifyBlobForQuorums: inclusion proof is invalid"
        );

        uint256 confirmedQuorumsBitmap;

        for (uint i = 0; i < blobHeader.quorumBlobParams.length; i++) {

            require(
                uint8(blobVerificationProof.batchMetadata.batchHeader.quorumNumbers[uint8(blobVerificationProof.quorumIndices[i])]) == 
                blobHeader.quorumBlobParams[i].quorumNumber, 
                "EigenDABlobVerificationUtils._verifyBlobForQuorums: quorumNumber does not match"
            );

            require(
                blobHeader.quorumBlobParams[i].confirmationThresholdPercentage >
                blobHeader.quorumBlobParams[i].adversaryThresholdPercentage, 
                "EigenDABlobVerificationUtils._verifyBlobForQuorums: threshold percentages are not valid"
            );

            require(
                blobHeader.quorumBlobParams[i].confirmationThresholdPercentage >= 
                eigenDAThresholdRegistry.getQuorumConfirmationThresholdPercentage(blobHeader.quorumBlobParams[i].quorumNumber), 
                "EigenDABlobVerificationUtils._verifyBlobForQuorums: confirmationThresholdPercentage is not met"
            );

            require(
                uint8(blobVerificationProof.batchMetadata.batchHeader.signedStakeForQuorums[uint8(blobVerificationProof.quorumIndices[i])]) >= 
                blobHeader.quorumBlobParams[i].confirmationThresholdPercentage, 
                "EigenDABlobVerificationUtils._verifyBlobForQuorums: confirmationThresholdPercentage is not met"
            );

            confirmedQuorumsBitmap = BitmapUtils.setBit(confirmedQuorumsBitmap, blobHeader.quorumBlobParams[i].quorumNumber);
        }

        require(
            BitmapUtils.isSubsetOf(
                BitmapUtils.orderedBytesArrayToBitmap(requiredQuorumNumbers),
                confirmedQuorumsBitmap
            ),
            "EigenDABlobVerificationUtils._verifyBlobForQuorums: required quorums are not a subset of the confirmed quorums"
        );
    }

    function _verifyBlobV2ForQuorums(
        IEigenDASignatureVerifier eigenDASignatureVerifier,
        SignedCertificate calldata signedCertificate,
        bytes memory requiredQuorumNumbers
    ) internal view {}

}
