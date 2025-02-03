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
 * @title Library of functions to be used by smart contracts wanting to verify submissions of blob certificates on EigenDA.
 * @author Layr Labs, Inc.
 */
library EigenDACertVerificationUtils {
    using BN254 for BN254.G1Point;

    uint256 public constant THRESHOLD_DENOMINATOR = 100;
    
    function _verifyDACertV1ForQuorums(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDABatchMetadataStorage batchMetadataStorage,
        BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof,
        bytes memory requiredQuorumNumbers
    ) internal view {
        require(
            EigenDAHasher.hashBatchMetadata(blobVerificationProof.batchMetadata) ==
            IEigenDABatchMetadataStorage(batchMetadataStorage).batchIdToBatchMetadataHash(blobVerificationProof.batchId),
            "EigenDACertVerificationUtils._verifyDACertForQuorums: batchMetadata does not match stored metadata"
        );

        require(
            Merkle.verifyInclusionKeccak(
                blobVerificationProof.inclusionProof, 
                blobVerificationProof.batchMetadata.batchHeader.blobHeadersRoot, 
                keccak256(abi.encodePacked(EigenDAHasher.hashBlobHeader(blobHeader))),
                blobVerificationProof.blobIndex
            ),
            "EigenDACertVerificationUtils._verifyDACertForQuorums: inclusion proof is invalid"
        );

        uint256 confirmedQuorumsBitmap;

        for (uint i = 0; i < blobHeader.quorumBlobParams.length; i++) {

            require(
                uint8(blobVerificationProof.batchMetadata.batchHeader.quorumNumbers[uint8(blobVerificationProof.quorumIndices[i])]) == 
                blobHeader.quorumBlobParams[i].quorumNumber, 
                "EigenDACertVerificationUtils._verifyDACertForQuorums: quorumNumber does not match"
            );

            require(
                blobHeader.quorumBlobParams[i].confirmationThresholdPercentage >
                blobHeader.quorumBlobParams[i].adversaryThresholdPercentage, 
                "EigenDACertVerificationUtils._verifyDACertForQuorums: threshold percentages are not valid"
            );

            require(
                blobHeader.quorumBlobParams[i].confirmationThresholdPercentage >= 
                eigenDAThresholdRegistry.getQuorumConfirmationThresholdPercentage(blobHeader.quorumBlobParams[i].quorumNumber), 
                "EigenDACertVerificationUtils._verifyDACertForQuorums: confirmationThresholdPercentage is not met"
            );

            require(
                uint8(blobVerificationProof.batchMetadata.batchHeader.signedStakeForQuorums[uint8(blobVerificationProof.quorumIndices[i])]) >= 
                blobHeader.quorumBlobParams[i].confirmationThresholdPercentage, 
                "EigenDACertVerificationUtils._verifyDACertForQuorums: confirmationThresholdPercentage is not met"
            );

            confirmedQuorumsBitmap = BitmapUtils.setBit(confirmedQuorumsBitmap, blobHeader.quorumBlobParams[i].quorumNumber);
        }

        require(
            BitmapUtils.isSubsetOf(
                BitmapUtils.orderedBytesArrayToBitmap(requiredQuorumNumbers),
                confirmedQuorumsBitmap
            ),
            "EigenDACertVerificationUtils._verifyDACertForQuorums: required quorums are not a subset of the confirmed quorums"
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
            "EigenDACertVerificationUtils._verifyDACertsForQuorums: blobHeaders and blobVerificationProofs length mismatch"
        );

        bytes memory confirmationThresholdPercentages = eigenDAThresholdRegistry.quorumConfirmationThresholdPercentages();

        for (uint i = 0; i < blobHeaders.length; ++i) {
            require(
                EigenDAHasher.hashBatchMetadata(blobVerificationProofs[i].batchMetadata) ==
                IEigenDABatchMetadataStorage(batchMetadataStorage).batchIdToBatchMetadataHash(blobVerificationProofs[i].batchId),
                "EigenDACertVerificationUtils._verifyDACertsForQuorums: batchMetadata does not match stored metadata"
            );

            require(
                Merkle.verifyInclusionKeccak(
                    blobVerificationProofs[i].inclusionProof, 
                    blobVerificationProofs[i].batchMetadata.batchHeader.blobHeadersRoot, 
                    keccak256(abi.encodePacked(EigenDAHasher.hashBlobHeader(blobHeaders[i]))),
                    blobVerificationProofs[i].blobIndex
                ),
                "EigenDACertVerificationUtils._verifyDACertsForQuorums: inclusion proof is invalid"
            );

            uint256 confirmedQuorumsBitmap;

            for (uint j = 0; j < blobHeaders[i].quorumBlobParams.length; j++) {

                require(
                    uint8(blobVerificationProofs[i].batchMetadata.batchHeader.quorumNumbers[uint8(blobVerificationProofs[i].quorumIndices[j])]) == 
                    blobHeaders[i].quorumBlobParams[j].quorumNumber, 
                    "EigenDACertVerificationUtils._verifyDACertsForQuorums: quorumNumber does not match"
                );

                require(
                    blobHeaders[i].quorumBlobParams[j].confirmationThresholdPercentage >
                    blobHeaders[i].quorumBlobParams[j].adversaryThresholdPercentage, 
                    "EigenDACertVerificationUtils._verifyDACertsForQuorums: threshold percentages are not valid"
                );

                require(
                    blobHeaders[i].quorumBlobParams[j].confirmationThresholdPercentage >= 
                    uint8(confirmationThresholdPercentages[blobHeaders[i].quorumBlobParams[j].quorumNumber]), 
                    "EigenDACertVerificationUtils._verifyDACertsForQuorums: confirmationThresholdPercentage is not met"
                );

                require(
                    uint8(blobVerificationProofs[i].batchMetadata.batchHeader.signedStakeForQuorums[uint8(blobVerificationProofs[i].quorumIndices[j])]) >= 
                    blobHeaders[i].quorumBlobParams[j].confirmationThresholdPercentage, 
                    "EigenDACertVerificationUtils._verifyDACertsForQuorums: confirmationThresholdPercentage is not met"
                );

                confirmedQuorumsBitmap = BitmapUtils.setBit(confirmedQuorumsBitmap, blobHeaders[i].quorumBlobParams[j].quorumNumber);
            }

            require(
                BitmapUtils.isSubsetOf(
                    BitmapUtils.orderedBytesArrayToBitmap(requiredQuorumNumbers),
                    confirmedQuorumsBitmap
                ),
                "EigenDACertVerificationUtils._verifyDACertsForQuorums: required quorums are not a subset of the confirmed quorums"
            );

        }
    }

    function _verifyDACertV2ForQuorums(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDASignatureVerifier signatureVerifier,
        IEigenDARelayRegistry eigenDARelayRegistry,
        BatchHeaderV2 memory batchHeader,
        BlobInclusionInfo memory blobInclusionInfo,
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature,
        SecurityThresholds memory securityThresholds,
        bytes memory requiredQuorumNumbers
    ) internal view {
        require(
            Merkle.verifyInclusionKeccak(
                blobInclusionInfo.inclusionProof, 
                batchHeader.batchRoot, 
                keccak256(abi.encodePacked(EigenDAHasher.hashBlobCertificate(blobInclusionInfo.blobCertificate))),
                blobInclusionInfo.blobIndex
            ),
            "EigenDACertVerificationUtils._verifyDACertV2ForQuorums: inclusion proof is invalid"
        );

        (
            QuorumStakeTotals memory quorumStakeTotals,
            bytes32 signatoryRecordHash
        ) = signatureVerifier.checkSignatures(
            EigenDAHasher.hashBatchHeaderV2(batchHeader),
            blobInclusionInfo.blobCertificate.blobHeader.quorumNumbers,
            batchHeader.referenceBlockNumber,
            nonSignerStakesAndSignature
        );

        _verifyRelayKeysSet(
            eigenDARelayRegistry,
            blobInclusionInfo.blobCertificate.relayKeys
        );

        _verifyDACertSecurityParams(
            eigenDAThresholdRegistry.getBlobParams(blobInclusionInfo.blobCertificate.blobHeader.version),
            securityThresholds
        );

        uint256 confirmedQuorumsBitmap;

        for (uint i = 0; i < blobInclusionInfo.blobCertificate.blobHeader.quorumNumbers.length; i++) {
            require(
                quorumStakeTotals.signedStakeForQuorum[i] * THRESHOLD_DENOMINATOR >= 
                quorumStakeTotals.totalStakeForQuorum[i] * securityThresholds.confirmationThreshold,
                "EigenDACertVerificationUtils._verifyDACertV2ForQuorums: signatories do not own at least threshold percentage of a quorum"
            );

            confirmedQuorumsBitmap = BitmapUtils.setBit(
                confirmedQuorumsBitmap, 
                uint8(blobInclusionInfo.blobCertificate.blobHeader.quorumNumbers[i])
            );
        }

        require(
            BitmapUtils.isSubsetOf(
                BitmapUtils.orderedBytesArrayToBitmap(requiredQuorumNumbers),
                confirmedQuorumsBitmap
            ),
            "EigenDACertVerificationUtils._verifyDACertV2ForQuorums: required quorums are not a subset of the confirmed quorums"
        );
    }

    /// @dev External function needed for try-catch wrapper
    function verifyDACertV2ForQuorumsExternal(
        IEigenDAThresholdRegistry _eigenDAThresholdRegistry,
        IEigenDASignatureVerifier _signatureVerifier,
        IEigenDARelayRegistry _eigenDARelayRegistry,
        BatchHeaderV2 memory _batchHeader,
        BlobInclusionInfo memory _blobInclusionInfo,
        NonSignerStakesAndSignature memory _nonSignerStakesAndSignature,
        SecurityThresholds memory _securityThresholds,
        bytes memory _requiredQuorumNumbers
    ) external view {
        EigenDACertVerificationUtils._verifyDACertV2ForQuorums(
            _eigenDAThresholdRegistry,
            _signatureVerifier,
            _eigenDARelayRegistry,
            _batchHeader,
            _blobInclusionInfo,
            _nonSignerStakesAndSignature,
            _securityThresholds,
            _requiredQuorumNumbers
        );
    }

    function _verifyDACertV2ForQuorumsForThresholds(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDASignatureVerifier signatureVerifier,
        IEigenDARelayRegistry eigenDARelayRegistry,
        BatchHeaderV2 memory batchHeader,
        BlobInclusionInfo memory blobInclusionInfo,
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature,
        SecurityThresholds[] memory securityThresholds,
        bytes memory requiredQuorumNumbers
    ) internal view {
        require(
            securityThresholds.length == blobInclusionInfo.blobCertificate.blobHeader.quorumNumbers.length,
            "EigenDACertVerificationUtils._verifyDACertV2ForQuorums: securityThresholds length does not match quorumNumbers"
        );

        require(
            Merkle.verifyInclusionKeccak(
                blobInclusionInfo.inclusionProof, 
                batchHeader.batchRoot, 
                keccak256(abi.encodePacked(EigenDAHasher.hashBlobCertificate(blobInclusionInfo.blobCertificate))),
                blobInclusionInfo.blobIndex
            ),
            "EigenDACertVerificationUtils._verifyDACertV2ForQuorums: inclusion proof is invalid"
        );

        (
            QuorumStakeTotals memory quorumStakeTotals,
            bytes32 signatoryRecordHash
        ) = signatureVerifier.checkSignatures(
            EigenDAHasher.hashBatchHeaderV2(batchHeader),
            blobInclusionInfo.blobCertificate.blobHeader.quorumNumbers,
            batchHeader.referenceBlockNumber,
            nonSignerStakesAndSignature
        );

        _verifyRelayKeysSet(
            eigenDARelayRegistry,
            blobInclusionInfo.blobCertificate.relayKeys
        );

        uint256 confirmedQuorumsBitmap;
        VersionedBlobParams memory blobParams = eigenDAThresholdRegistry.getBlobParams(blobInclusionInfo.blobCertificate.blobHeader.version);

        for (uint i = 0; i < blobInclusionInfo.blobCertificate.blobHeader.quorumNumbers.length; i++) {
            _verifyDACertSecurityParams(
                blobParams,
                securityThresholds[i]
            );

            require(
                quorumStakeTotals.signedStakeForQuorum[i] * THRESHOLD_DENOMINATOR >= 
                quorumStakeTotals.totalStakeForQuorum[i] * securityThresholds[i].confirmationThreshold,
                "EigenDACertVerificationUtils._verifyDACertV2ForQuorums: signatories do not own at least threshold percentage of a quorum"
            );

            confirmedQuorumsBitmap = BitmapUtils.setBit(
                confirmedQuorumsBitmap, 
                uint8(blobInclusionInfo.blobCertificate.blobHeader.quorumNumbers[i])
            );
        }

        require(
            BitmapUtils.isSubsetOf(
                BitmapUtils.orderedBytesArrayToBitmap(requiredQuorumNumbers),
                confirmedQuorumsBitmap
            ),
            "EigenDACertVerificationUtils._verifyDACertV2ForQuorums: required quorums are not a subset of the confirmed quorums"
        );
    }

    function _verifyDACertV2ForQuorumsFromSignedBatch(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDASignatureVerifier signatureVerifier,
        IEigenDARelayRegistry eigenDARelayRegistry,
        OperatorStateRetriever operatorStateRetriever,
        IRegistryCoordinator registryCoordinator,
        SignedBatch memory signedBatch,
        BlobInclusionInfo memory blobInclusionInfo,
        SecurityThresholds memory securityThresholds,
        bytes memory requiredQuorumNumbers
    ) internal view {
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature = _getNonSignerStakesAndSignature(
            operatorStateRetriever,
            registryCoordinator,
            signedBatch,
            blobInclusionInfo.blobCertificate.blobHeader.quorumNumbers
        );

        _verifyDACertV2ForQuorums(
            eigenDAThresholdRegistry,
            signatureVerifier,
            eigenDARelayRegistry,
            signedBatch.batchHeader,
            blobInclusionInfo,
            nonSignerStakesAndSignature,
            securityThresholds,
            requiredQuorumNumbers
        );
    }

    function _verifyDACertV2ForQuorumsForThresholdsFromSignedBatch(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDASignatureVerifier signatureVerifier,
        IEigenDARelayRegistry eigenDARelayRegistry,
        OperatorStateRetriever operatorStateRetriever,
        IRegistryCoordinator registryCoordinator,
        SignedBatch memory signedBatch,
        BlobInclusionInfo memory blobInclusionInfo,
        SecurityThresholds[] memory securityThresholds,
        bytes memory requiredQuorumNumbers
    ) internal view {
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature = _getNonSignerStakesAndSignature(
            operatorStateRetriever,
            registryCoordinator,
            signedBatch,
            blobInclusionInfo.blobCertificate.blobHeader.quorumNumbers
        );

        _verifyDACertV2ForQuorumsForThresholds(
            eigenDAThresholdRegistry,
            signatureVerifier,
            eigenDARelayRegistry,
            signedBatch.batchHeader,
            blobInclusionInfo,
            nonSignerStakesAndSignature,
            securityThresholds,
            requiredQuorumNumbers
        );
    }

    function _getNonSignerStakesAndSignature(
        OperatorStateRetriever operatorStateRetriever,
        IRegistryCoordinator registryCoordinator,
        SignedBatch memory signedBatch,
        bytes memory quorumNumbers
    ) internal view returns (NonSignerStakesAndSignature memory nonSignerStakesAndSignature) {
        bytes32[] memory nonSignerOperatorIds = new bytes32[](signedBatch.attestation.nonSignerPubkeys.length);
        for (uint i = 0; i < signedBatch.attestation.nonSignerPubkeys.length; ++i) {
            nonSignerOperatorIds[i] = BN254.hashG1Point(signedBatch.attestation.nonSignerPubkeys[i]);
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

    function _verifyDACertSecurityParams(
        VersionedBlobParams memory blobParams,
        SecurityThresholds memory securityThresholds
    ) internal pure {
        require(
            securityThresholds.confirmationThreshold > securityThresholds.adversaryThreshold,
            "EigenDACertVerificationUtils._verifyDACertSecurityParams: confirmationThreshold must be greater than adversaryThreshold"
        );
        uint256 gamma = securityThresholds.confirmationThreshold - securityThresholds.adversaryThreshold;
        uint256 n = (10000 - ((1_000_000 / gamma) / uint256(blobParams.codingRate))) * uint256(blobParams.numChunks);
        require(n >= blobParams.maxNumOperators * 10000, "EigenDACertVerificationUtils._verifyDACertSecurityParams: security assumptions are not met");
    }

    function _verifyRelayKeysSet(
        IEigenDARelayRegistry eigenDARelayRegistry,
        uint32[] memory relayKeys
    ) internal view {
        for (uint i = 0; i < relayKeys.length; ++i) {
            require(
                eigenDARelayRegistry.relayKeyToAddress(relayKeys[i]) != address(0),
                "EigenDACertVerificationUtils._verifyRelayKeysSet: relay key is not set"
            );
        }
    }
}
