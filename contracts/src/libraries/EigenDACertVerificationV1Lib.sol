// SPDX-License-Identifier: MIT

pragma solidity ^0.8.9;

import {Merkle} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/libraries/Merkle.sol";
import {BN254} from "lib/eigenlayer-middleware/src/libraries/BN254.sol";
import {EigenDAHasher} from "src/libraries/EigenDAHasher.sol";
import {BitmapUtils} from "lib/eigenlayer-middleware/src/libraries/BitmapUtils.sol";
import {IEigenDABatchMetadataStorage} from "src/interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDAThresholdRegistry} from "src/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/interfaces/IEigenDASignatureVerifier.sol";

import {BlobHeader, BlobVerificationProof} from "src/interfaces/IEigenDAStructs.sol";

/**
 * @title Library of functions to be used by smart contracts wanting to verify submissions of blob certificates on EigenDA.
 * @author Layr Labs, Inc.
 */
library EigenDACertVerificationV1Lib {
    function _verifyDACertV1ForQuorums(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDABatchMetadataStorage batchMetadataStorage,
        BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof,
        bytes memory requiredQuorumNumbers
    ) internal view {
        require(
            EigenDAHasher.hashBatchMetadata(blobVerificationProof.batchMetadata)
                == IEigenDABatchMetadataStorage(batchMetadataStorage).batchIdToBatchMetadataHash(
                    blobVerificationProof.batchId
                ),
            "EigenDACertVerificationV1Lib._verifyDACertForQuorums: batchMetadata does not match stored metadata"
        );

        require(
            Merkle.verifyInclusionKeccak(
                blobVerificationProof.inclusionProof,
                blobVerificationProof.batchMetadata.batchHeader.blobHeadersRoot,
                keccak256(abi.encodePacked(EigenDAHasher.hashBlobHeader(blobHeader))),
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
        BlobHeader[] calldata blobHeaders,
        BlobVerificationProof[] calldata blobVerificationProofs,
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
                EigenDAHasher.hashBatchMetadata(blobVerificationProofs[i].batchMetadata)
                    == IEigenDABatchMetadataStorage(batchMetadataStorage).batchIdToBatchMetadataHash(
                        blobVerificationProofs[i].batchId
                    ),
                "EigenDACertVerificationV1Lib._verifyDACertsForQuorums: batchMetadata does not match stored metadata"
            );

            require(
                Merkle.verifyInclusionKeccak(
                    blobVerificationProofs[i].inclusionProof,
                    blobVerificationProofs[i].batchMetadata.batchHeader.blobHeadersRoot,
                    keccak256(abi.encodePacked(EigenDAHasher.hashBlobHeader(blobHeaders[i]))),
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
}
