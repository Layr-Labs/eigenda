// SPDX-License-Identifier: MIT

pragma solidity ^0.8.9;

import {Merkle} from "../../lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/libraries/Merkle.sol";
import {BN254} from "../../lib/eigenlayer-middleware/src/libraries/BN254.sol";
import {EigenDAHasher} from "./EigenDAHasher.sol";
import {BitmapUtils} from "../../lib/eigenlayer-middleware/src/libraries/BitmapUtils.sol";
import {IEigenDABatchMetadataStorage} from "../interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDAThresholdRegistry} from "../interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "../interfaces/IEigenDASignatureVerifier.sol";
import {OperatorStateRetriever} from "../../lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "../../lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {IEigenDARelayRegistry} from "../interfaces/IEigenDARelayRegistry.sol";
import "../interfaces/IEigenDAStructs.sol";

/**
 * @title EigenDACertVerificationUtils
 * @notice Library of functions to be used by smart contracts wanting to verify submissions of blob certificates on EigenDA.
 */
library EigenDACertVerificationUtils {
    using BN254 for BN254.G1Point;

    uint256 public constant THRESHOLD_DENOMINATOR = 100;
    
    /**
     * @notice Verifies a V1 blob certificate for a set of quorums
     * @param eigenDAThresholdRegistry is the ThresholdRegistry contract address
     * @param batchMetadataStorage is the BatchMetadataStorage contract address
     * @param blobHeader is the blob header to verify
     * @param blobVerificationProof is the blob verification proof to verify
     * @param requiredQuorumNumbers is the required quorum numbers to verify against
     */
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

    /**
     * @notice Verifies a set of V1 blob certificates for a set of quorums
     * @param eigenDAThresholdRegistry is the ThresholdRegistry contract address
     * @param batchMetadataStorage is the BatchMetadataStorage contract address
     * @param blobHeaders is the set of blob headers to verify
     * @param blobVerificationProofs is the set of blob verification proofs to verify for each blob header
     * @param requiredQuorumNumbers is the required quorum numbers to verify against
     */
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

    /**
     * @notice Verifies a V2 blob certificate for a set of quorums
     * @param eigenDAThresholdRegistry is the ThresholdRegistry contract address
     * @param signatureVerifier is the SignatureVerifier contract address
     * @param eigenDARelayRegistry is the RelayRegistry contract address
     * @param batchHeader is the batch header to verify
     * @param blobInclusionInfo is the blob inclusion proof to verify against the batch
     * @param nonSignerStakesAndSignature is the non-signer stakes and signatures to check the signature against
     * @param securityThresholds are the confirmation and adversary thresholds to verify
     * @param requiredQuorumNumbers is the required quorum numbers to verify against 
     * @param signedQuorumNumbers are the quorum numbers that signed on the batch
     */
    function _verifyDACertV2ForQuorums(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDASignatureVerifier signatureVerifier,
        IEigenDARelayRegistry eigenDARelayRegistry,
        BatchHeaderV2 memory batchHeader,
        BlobInclusionInfo memory blobInclusionInfo,
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature,
        SecurityThresholds memory securityThresholds,
        bytes memory requiredQuorumNumbers,
        bytes memory signedQuorumNumbers
    ) internal view {
        // check blob inclusion in the batch from merkle proof
        require(
            Merkle.verifyInclusionKeccak(
                blobInclusionInfo.inclusionProof, 
                batchHeader.batchRoot, 
                keccak256(abi.encodePacked(EigenDAHasher.hashBlobCertificate(blobInclusionInfo.blobCertificate))),
                blobInclusionInfo.blobIndex
            ),
            "EigenDACertVerificationUtils._verifyDACertV2ForQuorums: inclusion proof is invalid"
        );

        // check BLS signature and get stake signed for batch quorums
        (
            QuorumStakeTotals memory quorumStakeTotals,
            bytes32 signatoryRecordHash
        ) = signatureVerifier.checkSignatures(
            EigenDAHasher.hashBatchHeaderV2(batchHeader),
            signedQuorumNumbers,
            batchHeader.referenceBlockNumber,
            nonSignerStakesAndSignature
        );

        // check relay keys are set
        _verifyRelayKeysSet(
            eigenDARelayRegistry,
            blobInclusionInfo.blobCertificate.relayKeys
        );

        // check the blob version is valid with security thresholds
        _verifyDACertSecurityParams(
            eigenDAThresholdRegistry.getBlobParams(blobInclusionInfo.blobCertificate.blobHeader.version),
            securityThresholds
        );

        uint256 confirmedQuorumsBitmap;

        // record confirmed quorums where signatories own at least the threshold percentage of the quorum
        for (uint i = 0; i < signedQuorumNumbers.length; i++) {
            if(
                quorumStakeTotals.signedStakeForQuorum[i] * THRESHOLD_DENOMINATOR >=
                quorumStakeTotals.totalStakeForQuorum[i] * securityThresholds.confirmationThreshold
            ) {
                confirmedQuorumsBitmap = BitmapUtils.setBit(
                    confirmedQuorumsBitmap, 
                    uint8(signedQuorumNumbers[i])
                );
            }
        }

        uint256 blobQuorumsBitmap = BitmapUtils.orderedBytesArrayToBitmap(blobInclusionInfo.blobCertificate.blobHeader.quorumNumbers);

        // check if the blob quorums are a subset of the confirmed quorums
        require(
            BitmapUtils.isSubsetOf(
                blobQuorumsBitmap,
                confirmedQuorumsBitmap
            ),
            "EigenDACertVerificationUtils._verifyDACertV2ForQuorums: blob quorums are not a subset of the confirmed quorums"
        );

        // check if the required quorums are a subset of the blob quorums
        require(
            BitmapUtils.isSubsetOf(
                BitmapUtils.orderedBytesArrayToBitmap(requiredQuorumNumbers),
                blobQuorumsBitmap
            ),
            "EigenDACertVerificationUtils._verifyDACertV2ForQuorums: required quorums are not a subset of the blob quorums"
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
        bytes memory _requiredQuorumNumbers,
        bytes memory _signedQuorumNumbers
    ) external view {
        EigenDACertVerificationUtils._verifyDACertV2ForQuorums(
            _eigenDAThresholdRegistry,
            _signatureVerifier,
            _eigenDARelayRegistry,
            _batchHeader,
            _blobInclusionInfo,
            _nonSignerStakesAndSignature,
            _securityThresholds,
            _requiredQuorumNumbers,
            _signedQuorumNumbers
        );
    }

    /**
     * @notice Verifies a V2 blob certificate for a set of quorums from a signed batch
     * @param eigenDAThresholdRegistry is the ThresholdRegistry contract address
     * @param signatureVerifier is the SignatureVerifier contract address
     * @param eigenDARelayRegistry is the RelayRegistry contract address
     * @param operatorStateRetriever is the OperatorStateRetriever contract address
     * @param registryCoordinator is the RegistryCoordinator contract address
     * @param signedBatch is the signed batch to verify
     * @param blobInclusionInfo is the blob inclusion proof to verify against the batch
     * @param securityThresholds are the confirmation and adversary thresholds to verify
     * @param requiredQuorumNumbers is the required quorum numbers to verify against 
     */
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
        (
            NonSignerStakesAndSignature memory nonSignerStakesAndSignature,
            bytes memory signedQuorumNumbers
        ) = _getNonSignerStakesAndSignature(
            operatorStateRetriever,
            registryCoordinator,
            signedBatch
        );

        _verifyDACertV2ForQuorums(
            eigenDAThresholdRegistry,
            signatureVerifier,
            eigenDARelayRegistry,
            signedBatch.batchHeader,
            blobInclusionInfo,
            nonSignerStakesAndSignature,
            securityThresholds,
            requiredQuorumNumbers,
            signedQuorumNumbers
        );
    }

    /// @dev Internal function to get the non-signer stakes and signature from the Attestation of a signed batch
    function _getNonSignerStakesAndSignature(
        OperatorStateRetriever operatorStateRetriever,
        IRegistryCoordinator registryCoordinator,
        SignedBatch memory signedBatch
    ) internal view returns (NonSignerStakesAndSignature memory nonSignerStakesAndSignature, bytes memory signedQuorumNumbers) {
        bytes32[] memory nonSignerOperatorIds = new bytes32[](signedBatch.attestation.nonSignerPubkeys.length);
        for (uint i = 0; i < signedBatch.attestation.nonSignerPubkeys.length; ++i) {
            nonSignerOperatorIds[i] = BN254.hashG1Point(signedBatch.attestation.nonSignerPubkeys[i]);
        }
      
        for (uint i = 0; i < signedBatch.attestation.quorumNumbers.length; ++i) {
            signedQuorumNumbers = abi.encodePacked(signedQuorumNumbers, uint8(signedBatch.attestation.quorumNumbers[i]));
        }

        OperatorStateRetriever.CheckSignaturesIndices memory checkSignaturesIndices = operatorStateRetriever.getCheckSignaturesIndices(
            registryCoordinator,
            signedBatch.batchHeader.referenceBlockNumber,
            signedQuorumNumbers,
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

    /// @dev Internal function to verify the security parameters of a V2 blob certificate
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

    /// @dev Internal function to verify that the provided relay keys are set on the RelayRegistry
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
