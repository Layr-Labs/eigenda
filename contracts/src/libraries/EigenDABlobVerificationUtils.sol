// SPDX-License-Identifier: MIT

pragma solidity ^0.8.9;

import {Merkle} from "eigenlayer-core/contracts/libraries/Merkle.sol";
import {BN254} from "eigenlayer-middleware/libraries/BN254.sol";
import {EigenDAHasher} from "./EigenDAHasher.sol";
import {BitmapUtils} from "eigenlayer-middleware/libraries/BitmapUtils.sol";
import {IEigenDABatchMetadataStorage} from "../interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDAThresholdRegistry} from "../interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "../interfaces/IEigenDASignatureVerifier.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {IEigenDARelayRegistry} from "../interfaces/IEigenDARelayRegistry.sol";
import "../interfaces/IEigenDAStructs.sol";

/**
 * @title Library of functions to be used by smart contracts wanting to verify submissions of blobs on EigenDA.
 * @author Layr Labs, Inc.
 */
library EigenDABlobVerificationUtils {
    using BN254 for BN254.G1Point;

    uint256 public constant THRESHOLD_DENOMINATOR = 100;
    
    function _verifyBlobV1ForQuorums(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDABatchMetadataStorage batchMetadataStorage,
        BlobHeader calldata blobHeader,
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

    function _verifyBlobsV1ForQuorums(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDABatchMetadataStorage batchMetadataStorage,
        BlobHeader[] calldata blobHeaders,
        BlobVerificationProof[] calldata blobVerificationProofs,
        bytes memory requiredQuorumNumbers
    ) internal view {
        require(
            blobHeaders.length == blobVerificationProofs.length,
            "EigenDABlobVerificationUtils._verifyBlobsForQuorums: blobHeaders and blobVerificationProofs length mismatch"
        );

        bytes memory confirmationThresholdPercentages = eigenDAThresholdRegistry.quorumConfirmationThresholdPercentages();

        for (uint i = 0; i < blobHeaders.length; ++i) {
            require(
                EigenDAHasher.hashBatchMetadata(blobVerificationProofs[i].batchMetadata) ==
                IEigenDABatchMetadataStorage(batchMetadataStorage).batchIdToBatchMetadataHash(blobVerificationProofs[i].batchId),
                "EigenDABlobVerificationUtils._verifyBlobForQuorums: batchMetadata does not match stored metadata"
            );

            require(
                Merkle.verifyInclusionKeccak(
                    blobVerificationProofs[i].inclusionProof, 
                    blobVerificationProofs[i].batchMetadata.batchHeader.blobHeadersRoot, 
                    keccak256(abi.encodePacked(EigenDAHasher.hashBlobHeader(blobHeaders[i]))),
                    blobVerificationProofs[i].blobIndex
                ),
                "EigenDABlobVerificationUtils._verifyBlobForQuorums: inclusion proof is invalid"
            );

            uint256 confirmedQuorumsBitmap;

            for (uint j = 0; j < blobHeaders[i].quorumBlobParams.length; j++) {

                require(
                    uint8(blobVerificationProofs[i].batchMetadata.batchHeader.quorumNumbers[uint8(blobVerificationProofs[i].quorumIndices[j])]) == 
                    blobHeaders[i].quorumBlobParams[j].quorumNumber, 
                    "EigenDABlobVerificationUtils._verifyBlobForQuorums: quorumNumber does not match"
                );

                require(
                    blobHeaders[i].quorumBlobParams[j].confirmationThresholdPercentage >
                    blobHeaders[i].quorumBlobParams[j].adversaryThresholdPercentage, 
                    "EigenDABlobVerificationUtils._verifyBlobForQuorums: threshold percentages are not valid"
                );

                require(
                    blobHeaders[i].quorumBlobParams[j].confirmationThresholdPercentage >= 
                    uint8(confirmationThresholdPercentages[blobHeaders[i].quorumBlobParams[j].quorumNumber]), 
                    "EigenDABlobVerificationUtils._verifyBlobForQuorums: confirmationThresholdPercentage is not met"
                );

                require(
                    uint8(blobVerificationProofs[i].batchMetadata.batchHeader.signedStakeForQuorums[uint8(blobVerificationProofs[i].quorumIndices[j])]) >= 
                    blobHeaders[i].quorumBlobParams[j].confirmationThresholdPercentage, 
                    "EigenDABlobVerificationUtils._verifyBlobForQuorums: confirmationThresholdPercentage is not met"
                );

                confirmedQuorumsBitmap = BitmapUtils.setBit(confirmedQuorumsBitmap, blobHeaders[i].quorumBlobParams[j].quorumNumber);
            }

            require(
                BitmapUtils.isSubsetOf(
                    BitmapUtils.orderedBytesArrayToBitmap(requiredQuorumNumbers),
                    confirmedQuorumsBitmap
                ),
                "EigenDABlobVerificationUtils._verifyBlobForQuorums: required quorums are not a subset of the confirmed quorums"
            );

        }
    }

    function _verifyBlobV2ForQuorums(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDASignatureVerifier signatureVerifier,
        IEigenDARelayRegistry eigenDARelayRegistry,
        BatchHeaderV2 memory batchHeader,
        BlobVerificationProofV2 memory blobVerificationProof,
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature,
        SecurityThresholds memory securityThresholds,
        bytes memory requiredQuorumNumbers
    ) internal view {
        require(
            Merkle.verifyInclusionKeccak(
                blobVerificationProof.inclusionProof, 
                batchHeader.batchRoot, 
                keccak256(abi.encodePacked(EigenDAHasher.hashBlobCertificate(blobVerificationProof.blobCertificate))),
                blobVerificationProof.blobIndex
            ),
            "EigenDABlobVerificationUtils._verifyBlobV2ForQuorums: inclusion proof is invalid"
        );

        (
            QuorumStakeTotals memory quorumStakeTotals,
            bytes32 signatoryRecordHash
        ) = signatureVerifier.checkSignatures(
            EigenDAHasher.hashBatchHeaderV2(batchHeader),
            blobVerificationProof.blobCertificate.blobHeader.quorumNumbers,
            batchHeader.referenceBlockNumber,
            nonSignerStakesAndSignature
        );

        _verifyRelayKeysSet(
            eigenDARelayRegistry,
            blobVerificationProof.blobCertificate.relayKeys
        );

        _verifyBlobSecurityParams(
            eigenDAThresholdRegistry.getBlobParams(blobVerificationProof.blobCertificate.blobHeader.version),
            securityThresholds
        );

        uint256 confirmedQuorumsBitmap;

        for (uint i = 0; i < blobVerificationProof.blobCertificate.blobHeader.quorumNumbers.length; i++) {
            require(
                quorumStakeTotals.signedStakeForQuorum[i] * THRESHOLD_DENOMINATOR >= 
                quorumStakeTotals.totalStakeForQuorum[i] * securityThresholds.confirmationThreshold,
                "EigenDABlobVerificationUtils._verifyBlobV2ForQuorums: signatories do not own at least threshold percentage of a quorum"
            );

            confirmedQuorumsBitmap = BitmapUtils.setBit(
                confirmedQuorumsBitmap, 
                uint8(blobVerificationProof.blobCertificate.blobHeader.quorumNumbers[i])
            );
        }

        require(
            BitmapUtils.isSubsetOf(
                BitmapUtils.orderedBytesArrayToBitmap(requiredQuorumNumbers),
                confirmedQuorumsBitmap
            ),
            "EigenDABlobVerificationUtils._verifyBlobV2ForQuorums: required quorums are not a subset of the confirmed quorums"
        );
    }

    function _verifyBlobV2ForQuorumsForThresholds(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDASignatureVerifier signatureVerifier,
        IEigenDARelayRegistry eigenDARelayRegistry,
        BatchHeaderV2 memory batchHeader,
        BlobVerificationProofV2 memory blobVerificationProof,
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature,
        SecurityThresholds[] memory securityThresholds,
        bytes memory requiredQuorumNumbers
    ) internal view {
        require(
            securityThresholds.length == blobVerificationProof.blobCertificate.blobHeader.quorumNumbers.length,
            "EigenDABlobVerificationUtils._verifyBlobV2ForQuorums: securityThresholds length does not match quorumNumbers"
        );

        require(
            Merkle.verifyInclusionKeccak(
                blobVerificationProof.inclusionProof, 
                batchHeader.batchRoot, 
                keccak256(abi.encodePacked(EigenDAHasher.hashBlobCertificate(blobVerificationProof.blobCertificate))),
                blobVerificationProof.blobIndex
            ),
            "EigenDABlobVerificationUtils._verifyBlobV2ForQuorums: inclusion proof is invalid"
        );

        (
            QuorumStakeTotals memory quorumStakeTotals,
            bytes32 signatoryRecordHash
        ) = signatureVerifier.checkSignatures(
            EigenDAHasher.hashBatchHeaderV2(batchHeader),
            blobVerificationProof.blobCertificate.blobHeader.quorumNumbers,
            batchHeader.referenceBlockNumber,
            nonSignerStakesAndSignature
        );

        _verifyRelayKeysSet(
            eigenDARelayRegistry,
            blobVerificationProof.blobCertificate.relayKeys
        );

        uint256 confirmedQuorumsBitmap;
        VersionedBlobParams memory blobParams = eigenDAThresholdRegistry.getBlobParams(blobVerificationProof.blobCertificate.blobHeader.version);

        for (uint i = 0; i < blobVerificationProof.blobCertificate.blobHeader.quorumNumbers.length; i++) {
            _verifyBlobSecurityParams(
                blobParams,
                securityThresholds[i]
            );

            require(
                quorumStakeTotals.signedStakeForQuorum[i] * THRESHOLD_DENOMINATOR >= 
                quorumStakeTotals.totalStakeForQuorum[i] * securityThresholds[i].confirmationThreshold,
                "EigenDABlobVerificationUtils._verifyBlobV2ForQuorums: signatories do not own at least threshold percentage of a quorum"
            );

            confirmedQuorumsBitmap = BitmapUtils.setBit(
                confirmedQuorumsBitmap, 
                uint8(blobVerificationProof.blobCertificate.blobHeader.quorumNumbers[i])
            );
        }

        require(
            BitmapUtils.isSubsetOf(
                BitmapUtils.orderedBytesArrayToBitmap(requiredQuorumNumbers),
                confirmedQuorumsBitmap
            ),
            "EigenDABlobVerificationUtils._verifyBlobV2ForQuorums: required quorums are not a subset of the confirmed quorums"
        );
    }

    function _verifyBlobV2ForQuorumsFromSignedBatch(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDASignatureVerifier signatureVerifier,
        IEigenDARelayRegistry eigenDARelayRegistry,
        OperatorStateRetriever operatorStateRetriever,
        IRegistryCoordinator registryCoordinator,
        SignedBatch memory signedBatch,
        BlobVerificationProofV2 memory blobVerificationProof,
        SecurityThresholds memory securityThresholds,
        bytes memory requiredQuorumNumbers
    ) internal view {
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature = _getNonSignerStakesAndSignature(
            operatorStateRetriever,
            registryCoordinator,
            signedBatch
        );

        _verifyBlobV2ForQuorums(
            eigenDAThresholdRegistry,
            signatureVerifier,
            eigenDARelayRegistry,
            signedBatch.batchHeader,
            blobVerificationProof,
            nonSignerStakesAndSignature,
            securityThresholds,
            requiredQuorumNumbers
        );
    }

    function _verifyBlobV2ForQuorumsForThresholdsFromSignedBatch(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDASignatureVerifier signatureVerifier,
        IEigenDARelayRegistry eigenDARelayRegistry,
        OperatorStateRetriever operatorStateRetriever,
        IRegistryCoordinator registryCoordinator,
        SignedBatch memory signedBatch,
        BlobVerificationProofV2 memory blobVerificationProof,
        SecurityThresholds[] memory securityThresholds,
        bytes memory requiredQuorumNumbers
    ) internal view {
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature = _getNonSignerStakesAndSignature(
            operatorStateRetriever,
            registryCoordinator,
            signedBatch
        );

        _verifyBlobV2ForQuorumsForThresholds(
            eigenDAThresholdRegistry,
            signatureVerifier,
            eigenDARelayRegistry,
            signedBatch.batchHeader,
            blobVerificationProof,
            nonSignerStakesAndSignature,
            securityThresholds,
            requiredQuorumNumbers
        );
    }

    function _getNonSignerStakesAndSignature(
        OperatorStateRetriever operatorStateRetriever,
        IRegistryCoordinator registryCoordinator,
        SignedBatch memory signedBatch
    ) internal view returns (NonSignerStakesAndSignature memory nonSignerStakesAndSignature) {
        bytes32[] memory nonSignerOperatorIds = new bytes32[](signedBatch.attestation.nonSignerPubkeys.length);
        for (uint i = 0; i < signedBatch.attestation.nonSignerPubkeys.length; ++i) {
            nonSignerOperatorIds[i] = BN254.hashG1Point(signedBatch.attestation.nonSignerPubkeys[i]);
        }

        bytes memory quorumNumbers;
        for (uint i = 0; i < signedBatch.attestation.quorumNumbers.length; ++i) {
            quorumNumbers = abi.encodePacked(quorumNumbers, uint8(signedBatch.attestation.quorumNumbers[i]));
        }

        OperatorStateRetriever.CheckSignaturesIndices memory checkSignaturesIndices = operatorStateRetriever.getCheckSignaturesIndices(
            registryCoordinator,
            signedBatch.batchHeader.referenceBlockNumber,
            quorumNumbers,
            nonSignerOperatorIds
        );

        nonSignerStakesAndSignature.nonSignerQuorumBitmapIndices = checkSignaturesIndices.nonSignerQuorumBitmapIndices; 
        nonSignerStakesAndSignature.nonSignerPubkeys = signedBatch.attestation.nonSignerPubkeys; 
        nonSignerStakesAndSignature.quorumApks = signedBatch.attestation.quorumApks; 
        nonSignerStakesAndSignature.apkG2 = signedBatch.attestation.apkG2; 
        nonSignerStakesAndSignature.sigma = signedBatch.attestation.sigma; 
        nonSignerStakesAndSignature.quorumApkIndices = checkSignaturesIndices.quorumApkIndices; 
        nonSignerStakesAndSignature.totalStakeIndices = checkSignaturesIndices.totalStakeIndices; 
        nonSignerStakesAndSignature.nonSignerStakeIndices = checkSignaturesIndices.nonSignerStakeIndices; 
    }

    function _verifyBlobSecurityParams(
        VersionedBlobParams memory blobParams,
        SecurityThresholds memory securityThresholds
    ) internal pure {
        require(
            securityThresholds.confirmationThreshold > securityThresholds.adversaryThreshold,
            "EigenDABlobVerificationUtils._verifyBlobSecurityParams: confirmationThreshold must be greater than adversaryThreshold"
        );
        uint256 gamma = securityThresholds.confirmationThreshold - securityThresholds.adversaryThreshold;
        uint256 n = (10000 - ((1_000_000 / gamma) / uint256(blobParams.codingRate))) * uint256(blobParams.numChunks);
        require(n >= blobParams.maxNumOperators * 10000, "EigenDABlobVerificationUtils._verifyBlobSecurityParams: security assumptions are not met");
    }

    function _verifyRelayKeysSet(
        IEigenDARelayRegistry eigenDARelayRegistry,
        uint32[] memory relayKeys
    ) internal view {
        for (uint i = 0; i < relayKeys.length; ++i) {
            require(
                eigenDARelayRegistry.relayKeyToAddress(relayKeys[i]) != address(0),
                "EigenDABlobVerificationUtils._verifyRelayKeysSet: relay key is not set"
            );
        }
    }
}
