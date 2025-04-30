// SPDX-License-Identifier: MIT

pragma solidity ^0.8.9;

import {Merkle} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/libraries/Merkle.sol";
import {BN254} from "lib/eigenlayer-middleware/src/libraries/BN254.sol";
import {BitmapUtils} from "lib/eigenlayer-middleware/src/libraries/BitmapUtils.sol";
import {IEigenDABatchMetadataStorage} from "src/interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDAThresholdRegistry} from "src/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/interfaces/IEigenDASignatureVerifier.sol";

import {EigenDATypesV1 as DATypesV1} from "src/libraries/V1/EigenDATypesV1.sol";

/**
 * @title Library of functions to be used by smart contracts wanting to verify submissions of blob certificates on EigenDA.
 * @author Layr Labs, Inc.
 */
library EigenDACertVerificationV1Lib {
    function _verifyDACertV1ForQuorums(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDABatchMetadataStorage batchMetadataStorage,
        DATypesV1.BlobHeader calldata blobHeader,
        DATypesV1.BlobVerificationProof calldata blobVerificationProof,
        bytes memory requiredQuorumNumbers
    ) internal view {
        require(
            hashBatchMetadata(blobVerificationProof.batchMetadata)
                == IEigenDABatchMetadataStorage(batchMetadataStorage).batchIdToBatchMetadataHash(
                    blobVerificationProof.batchId
                ),
            "EigenDACertVerificationV1Lib._verifyDACertForQuorums: batchMetadata does not match stored metadata"
        );

        require(
            Merkle.verifyInclusionKeccak(
                blobVerificationProof.inclusionProof,
                blobVerificationProof.batchMetadata.batchHeader.blobHeadersRoot,
                keccak256(abi.encodePacked(hashBlobHeader(blobHeader))),
                blobVerificationProof.blobIndex
            ),
            "EigenDACertVerificationV1Lib._verifyDACertForQuorums: inclusion proof is invalid"
        );

        uint256 confirmedQuorumsBitmap;

        for (uint256 i = 0; i < blobHeader.quorumBlobParams.length; i++) {
            require(
                uint8(
                    blobVerificationProof.batchMetadata.batchHeader.quorumNumbers[uint8(
                        blobVerificationProof.quorumIndices[i]
                    )]
                ) == blobHeader.quorumBlobParams[i].quorumNumber,
                "EigenDACertVerificationV1Lib._verifyDACertForQuorums: quorumNumber does not match"
            );

            require(
                blobHeader.quorumBlobParams[i].confirmationThresholdPercentage
                    > blobHeader.quorumBlobParams[i].adversaryThresholdPercentage,
                "EigenDACertVerificationV1Lib._verifyDACertForQuorums: threshold percentages are not valid"
            );

            require(
                blobHeader.quorumBlobParams[i].confirmationThresholdPercentage
                    >= eigenDAThresholdRegistry.getQuorumConfirmationThresholdPercentage(
                        blobHeader.quorumBlobParams[i].quorumNumber
                    ),
                "EigenDACertVerificationV1Lib._verifyDACertForQuorums: confirmationThresholdPercentage is not met"
            );

            require(
                uint8(
                    blobVerificationProof.batchMetadata.batchHeader.signedStakeForQuorums[uint8(
                        blobVerificationProof.quorumIndices[i]
                    )]
                ) >= blobHeader.quorumBlobParams[i].confirmationThresholdPercentage,
                "EigenDACertVerificationV1Lib._verifyDACertForQuorums: confirmationThresholdPercentage is not met"
            );

            confirmedQuorumsBitmap =
                BitmapUtils.setBit(confirmedQuorumsBitmap, blobHeader.quorumBlobParams[i].quorumNumber);
        }

        require(
            BitmapUtils.isSubsetOf(BitmapUtils.orderedBytesArrayToBitmap(requiredQuorumNumbers), confirmedQuorumsBitmap),
            "EigenDACertVerificationV1Lib._verifyDACertForQuorums: required quorums are not a subset of the confirmed quorums"
        );
    }

    function _verifyDACertsV1ForQuorums(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDABatchMetadataStorage batchMetadataStorage,
        DATypesV1.BlobHeader[] calldata blobHeaders,
        DATypesV1.BlobVerificationProof[] calldata blobVerificationProofs,
        bytes memory requiredQuorumNumbers
    ) internal view {
        require(
            blobHeaders.length == blobVerificationProofs.length,
            "EigenDACertVerificationV1Lib._verifyDACertsForQuorums: blobHeaders and blobVerificationProofs length mismatch"
        );

        bytes memory confirmationThresholdPercentages =
            eigenDAThresholdRegistry.quorumConfirmationThresholdPercentages();

        for (uint256 i = 0; i < blobHeaders.length; ++i) {
            require(
                hashBatchMetadata(blobVerificationProofs[i].batchMetadata)
                    == IEigenDABatchMetadataStorage(batchMetadataStorage).batchIdToBatchMetadataHash(
                        blobVerificationProofs[i].batchId
                    ),
                "EigenDACertVerificationV1Lib._verifyDACertsForQuorums: batchMetadata does not match stored metadata"
            );

            require(
                Merkle.verifyInclusionKeccak(
                    blobVerificationProofs[i].inclusionProof,
                    blobVerificationProofs[i].batchMetadata.batchHeader.blobHeadersRoot,
                    keccak256(abi.encodePacked(hashBlobHeader(blobHeaders[i]))),
                    blobVerificationProofs[i].blobIndex
                ),
                "EigenDACertVerificationV1Lib._verifyDACertsForQuorums: inclusion proof is invalid"
            );

            uint256 confirmedQuorumsBitmap;

            for (uint256 j = 0; j < blobHeaders[i].quorumBlobParams.length; j++) {
                require(
                    uint8(
                        blobVerificationProofs[i].batchMetadata.batchHeader.quorumNumbers[uint8(
                            blobVerificationProofs[i].quorumIndices[j]
                        )]
                    ) == blobHeaders[i].quorumBlobParams[j].quorumNumber,
                    "EigenDACertVerificationV1Lib._verifyDACertsForQuorums: quorumNumber does not match"
                );

                require(
                    blobHeaders[i].quorumBlobParams[j].confirmationThresholdPercentage
                        > blobHeaders[i].quorumBlobParams[j].adversaryThresholdPercentage,
                    "EigenDACertVerificationV1Lib._verifyDACertsForQuorums: threshold percentages are not valid"
                );

                require(
                    blobHeaders[i].quorumBlobParams[j].confirmationThresholdPercentage
                        >= uint8(confirmationThresholdPercentages[blobHeaders[i].quorumBlobParams[j].quorumNumber]),
                    "EigenDACertVerificationV1Lib._verifyDACertsForQuorums: confirmationThresholdPercentage is not met"
                );

                require(
                    uint8(
                        blobVerificationProofs[i].batchMetadata.batchHeader.signedStakeForQuorums[uint8(
                            blobVerificationProofs[i].quorumIndices[j]
                        )]
                    ) >= blobHeaders[i].quorumBlobParams[j].confirmationThresholdPercentage,
                    "EigenDACertVerificationV1Lib._verifyDACertsForQuorums: confirmationThresholdPercentage is not met"
                );

                confirmedQuorumsBitmap =
                    BitmapUtils.setBit(confirmedQuorumsBitmap, blobHeaders[i].quorumBlobParams[j].quorumNumber);
            }

            require(
                BitmapUtils.isSubsetOf(
                    BitmapUtils.orderedBytesArrayToBitmap(requiredQuorumNumbers), confirmedQuorumsBitmap
                ),
                "EigenDACertVerificationV1Lib._verifyDACertsForQuorums: required quorums are not a subset of the confirmed quorums"
            );
        }
    }

    
    /**
     * @notice hashes the given metdata into the commitment that will be stored in the contract
     * @param batchHeaderHash the hash of the batchHeader
     * @param signatoryRecordHash the hash of the signatory record
     * @param blockNumber the block number at which the batch was confirmed
     */
    function hashBatchHashedMetadata(bytes32 batchHeaderHash, bytes32 signatoryRecordHash, uint32 blockNumber)
        internal
        pure
        returns (bytes32)
    {
        return keccak256(abi.encodePacked(batchHeaderHash, signatoryRecordHash, blockNumber));
    }

    /**
     * @notice hashes the given metdata into the commitment that will be stored in the contract
     * @param batchHeaderHash the hash of the batchHeader
     * @param confirmationData the confirmation data of the batch
     * @param blockNumber the block number at which the batch was confirmed
     */
    function hashBatchHashedMetadata(bytes32 batchHeaderHash, bytes memory confirmationData, uint32 blockNumber)
        internal
        pure
        returns (bytes32)
    {
        return keccak256(abi.encodePacked(batchHeaderHash, confirmationData, blockNumber));
    }

    /**
     * @notice given the batchHeader in the provided metdata, calculates the hash of the batchMetadata
     * @param batchMetadata the metadata of the batch
     */
    function hashBatchMetadata(DATypesV1.BatchMetadata memory batchMetadata) internal pure returns (bytes32) {
        return hashBatchHashedMetadata(
            keccak256(abi.encode(batchMetadata.batchHeader)),
            batchMetadata.signatoryRecordHash,
            batchMetadata.confirmationBlockNumber
        );
    }

    /**
     * @notice hashes the given batch header
     * @param batchHeader the batch header to hash
     */
    function hashBatchHeaderMemory(DATypesV1.BatchHeader memory batchHeader) internal pure returns (bytes32) {
        return keccak256(abi.encode(batchHeader));
    }

    /**
     * @notice hashes the given batch header
     * @param batchHeader the batch header to hash
     */
    function hashBatchHeader(DATypesV1.BatchHeader calldata batchHeader) internal pure returns (bytes32) {
        return keccak256(abi.encode(batchHeader));
    }

    /**
     * @notice hashes the given reduced batch header
     * @param reducedBatchHeader the reduced batch header to hash
     */
    function hashReducedBatchHeader(DATypesV1.ReducedBatchHeader memory reducedBatchHeader) internal pure returns (bytes32) {
        return keccak256(abi.encode(reducedBatchHeader));
    }

    /**
     * @notice hashes the given blob header
     * @param blobHeader the blob header to hash
     */
    function hashBlobHeader(DATypesV1.BlobHeader memory blobHeader) internal pure returns (bytes32) {
        return keccak256(abi.encode(blobHeader));
    }

    /**
     * @notice converts a batch header to a reduced batch header
     * @param batchHeader the batch header to convert
     */
    function convertBatchHeaderToReducedBatchHeader(DATypesV1.BatchHeader memory batchHeader)
        internal
        pure
        returns (DATypesV1.ReducedBatchHeader memory)
    {
        return DATypesV1.ReducedBatchHeader({
            blobHeadersRoot: batchHeader.blobHeadersRoot,
            referenceBlockNumber: batchHeader.referenceBlockNumber
        });
    }

    /**
     * @notice converts the given batch header to a reduced batch header and then hashes it
     * @param batchHeader the batch header to hash
     */
    function hashBatchHeaderToReducedBatchHeader(DATypesV1.BatchHeader memory batchHeader) internal pure returns (bytes32) {
        return keccak256(abi.encode(convertBatchHeaderToReducedBatchHeader(batchHeader)));
    }
}