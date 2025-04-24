// SPDX-License-Identifier: MIT

pragma solidity ^0.8.9;

import {BN254} from "lib/eigenlayer-middleware/src/libraries/BN254.sol";
import {Merkle} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/libraries/Merkle.sol";
import {IEigenDAThresholdRegistry} from "src/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/interfaces/IEigenDASignatureVerifier.sol";
import {IEigenDARelayRegistry} from "src/interfaces/IEigenDARelayRegistry.sol";

import {EigenDAHasher} from "src/libraries/EigenDAHasher.sol";
import {BitmapUtils} from "lib/eigenlayer-middleware/src/libraries/BitmapUtils.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {BatchHeaderV2, BlobInclusionInfo, SignedBatch} from "src/interfaces/IEigenDAStructs.sol";
import {
    BlobHeader,
    BlobVerificationProof,
    QuorumStakeTotals,
    NonSignerStakesAndSignature,
    SecurityThresholds,
    VersionedBlobParams
} from "src/interfaces/IEigenDAStructs.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";

/**
 * @title Library of functions to be used by smart contracts wanting to verify submissions of blob certificates on EigenDA.
 * @author Layr Labs, Inc.
 */
library EigenDACertVerificationV2Lib {
    using BN254 for BN254.G1Point;

    uint256 public constant THRESHOLD_DENOMINATOR = 100;

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
            "EigenDACertVerificationV2Lib._verifyDACertV2ForQuorums: inclusion proof is invalid"
        );

        // check BLS signature and get stake signed for batch quorums
        (QuorumStakeTotals memory quorumStakeTotals,) = signatureVerifier.checkSignatures(
            EigenDAHasher.hashBatchHeaderV2(batchHeader),
            signedQuorumNumbers,
            batchHeader.referenceBlockNumber,
            nonSignerStakesAndSignature
        );

        // check relay keys are set
        _verifyRelayKeysSet(eigenDARelayRegistry, blobInclusionInfo.blobCertificate.relayKeys);

        // check the blob version is valid with security thresholds
        _verifyDACertSecurityParams(
            eigenDAThresholdRegistry.getBlobParams(blobInclusionInfo.blobCertificate.blobHeader.version),
            securityThresholds
        );

        uint256 confirmedQuorumsBitmap;

        // record confirmed quorums where signatories own at least the threshold percentage of the quorum
        for (uint256 i = 0; i < signedQuorumNumbers.length; i++) {
            if (
                quorumStakeTotals.signedStakeForQuorum[i] * THRESHOLD_DENOMINATOR
                    >= quorumStakeTotals.totalStakeForQuorum[i] * securityThresholds.confirmationThreshold
            ) {
                confirmedQuorumsBitmap = BitmapUtils.setBit(confirmedQuorumsBitmap, uint8(signedQuorumNumbers[i]));
            }
        }

        uint256 blobQuorumsBitmap =
            BitmapUtils.orderedBytesArrayToBitmap(blobInclusionInfo.blobCertificate.blobHeader.quorumNumbers);

        // check if the blob quorums are a subset of the confirmed quorums
        require(
            BitmapUtils.isSubsetOf(blobQuorumsBitmap, confirmedQuorumsBitmap),
            "EigenDACertVerificationV2Lib._verifyDACertV2ForQuorums: blob quorums are not a subset of the confirmed quorums"
        );

        // check if the required quorums are a subset of the blob quorums
        require(
            BitmapUtils.isSubsetOf(BitmapUtils.orderedBytesArrayToBitmap(requiredQuorumNumbers), blobQuorumsBitmap),
            "EigenDACertVerificationV2Lib._verifyDACertV2ForQuorums: required quorums are not a subset of the blob quorums"
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
        _verifyDACertV2ForQuorums(
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
        (NonSignerStakesAndSignature memory nonSignerStakesAndSignature, bytes memory signedQuorumNumbers) =
            _getNonSignerStakesAndSignature(operatorStateRetriever, registryCoordinator, signedBatch);

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

    function _getNonSignerStakesAndSignature(
        OperatorStateRetriever operatorStateRetriever,
        IRegistryCoordinator registryCoordinator,
        SignedBatch memory signedBatch
    )
        internal
        view
        returns (NonSignerStakesAndSignature memory nonSignerStakesAndSignature, bytes memory signedQuorumNumbers)
    {
        bytes32[] memory nonSignerOperatorIds = new bytes32[](signedBatch.attestation.nonSignerPubkeys.length);
        for (uint256 i = 0; i < signedBatch.attestation.nonSignerPubkeys.length; ++i) {
            nonSignerOperatorIds[i] = BN254.hashG1Point(signedBatch.attestation.nonSignerPubkeys[i]);
        }

        for (uint256 i = 0; i < signedBatch.attestation.quorumNumbers.length; ++i) {
            signedQuorumNumbers = abi.encodePacked(signedQuorumNumbers, uint8(signedBatch.attestation.quorumNumbers[i]));
        }

        OperatorStateRetriever.CheckSignaturesIndices memory checkSignaturesIndices = operatorStateRetriever
            .getCheckSignaturesIndices(
            registryCoordinator, signedBatch.batchHeader.referenceBlockNumber, signedQuorumNumbers, nonSignerOperatorIds
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
            "EigenDACertVerificationV2Lib._verifyDACertSecurityParams: confirmationThreshold must be greater than adversaryThreshold"
        );
        uint256 gamma = securityThresholds.confirmationThreshold - securityThresholds.adversaryThreshold;
        uint256 n = (10000 - ((1_000_000 / gamma) / uint256(blobParams.codingRate))) * uint256(blobParams.numChunks);
        require(
            n >= blobParams.maxNumOperators * 10000,
            "EigenDACertVerificationV2Lib._verifyDACertSecurityParams: security assumptions are not met"
        );
    }

    function _verifyRelayKeysSet(IEigenDARelayRegistry eigenDARelayRegistry, uint32[] memory relayKeys) internal view {
        for (uint256 i = 0; i < relayKeys.length; ++i) {
            require(
                eigenDARelayRegistry.relayKeyToAddress(relayKeys[i]) != address(0),
                "EigenDACertVerificationV2Lib._verifyRelayKeysSet: relay key is not set"
            );
        }
    }
}
