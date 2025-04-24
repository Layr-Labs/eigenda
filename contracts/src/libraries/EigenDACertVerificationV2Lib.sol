// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAThresholdRegistry} from "src/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/interfaces/IEigenDASignatureVerifier.sol";
import {IEigenDARelayRegistry} from "src/interfaces/IEigenDARelayRegistry.sol";
import {EigenDAHasher} from "src/libraries/EigenDAHasher.sol";
import {BN254} from "lib/eigenlayer-middleware/src/libraries/BN254.sol";
import {Merkle} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/libraries/Merkle.sol";
import {BitmapUtils} from "lib/eigenlayer-middleware/src/libraries/BitmapUtils.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry} from "lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {IBLSApkRegistry} from "lib/eigenlayer-middleware/src/interfaces/IBLSApkRegistry.sol";

import {
    BatchHeaderV2,
    BlobInclusionInfo,
    NonSignerStakesAndSignature,
    SecurityThresholds,
    QuorumStakeTotals,
    VersionedBlobParams,
    SignedBatch,
    EigenDACertV2
} from "src/interfaces/IEigenDAStructs.sol";

/**
 * @title EigenDACertVerificationV2Lib - EigenDA V2 certificate verification library
 * @author Layr Labs, Inc.
 * @notice Library of functions for verifying EigenDA V2 certificates
 * @dev Provides functions for verifying blob certificates, inclusion proofs, signatures, and security parameters
 */
library EigenDACertVerificationV2Lib {
    using BN254 for BN254.G1Point;

    /// @notice Denominator used for threshold percentage calculations (100 for percentages)
    uint256 internal constant THRESHOLD_DENOMINATOR = 100;

    /// @notice Thrown when the inclusion proof is invalid
    /// @param blobIndex The index of the blob in the batch
    /// @param blobHash The hash of the blob certificate
    /// @param rootHash The root hash of the merkle tree
    error InvalidInclusionProof(uint256 blobIndex, bytes32 blobHash, bytes32 rootHash);

    /// @notice Thrown when security assumptions are not met
    /// @param gamma The difference between confirmation and adversary thresholds
    /// @param n The calculated security parameter
    /// @param minRequired The minimum required value for n
    error SecurityAssumptionsNotMet(uint256 gamma, uint256 n, uint256 minRequired);

    /// @notice Thrown when blob quorums are not a subset of confirmed quorums
    /// @param blobQuorumsBitmap The bitmap of blob quorums
    /// @param confirmedQuorumsBitmap The bitmap of confirmed quorums
    error BlobQuorumsNotSubset(uint256 blobQuorumsBitmap, uint256 confirmedQuorumsBitmap);

    /// @notice Thrown when required quorums are not a subset of blob quorums
    /// @param requiredQuorumsBitmap The bitmap of required quorums
    /// @param blobQuorumsBitmap The bitmap of blob quorums
    error RequiredQuorumsNotSubset(uint256 requiredQuorumsBitmap, uint256 blobQuorumsBitmap);

    /// @notice Status codes for certificate verification results
    enum StatusCode {
        SUCCESS, // Verification succeeded
        INVALID_INCLUSION_PROOF, // Merkle inclusion proof is invalid
        SECURITY_ASSUMPTIONS_NOT_MET, // Security assumptions not met
        BLOB_QUORUMS_NOT_SUBSET, // Blob quorums not a subset of confirmed quorums
        REQUIRED_QUORUMS_NOT_SUBSET // Required quorums not a subset of blob quorums

    }

    function verifyDACertV2(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDASignatureVerifier signatureVerifier,
        EigenDACertV2 calldata cert,
        SecurityThresholds memory securityThresholds,
        bytes memory requiredQuorumNumbers
    ) internal view {
        (StatusCode status, bytes memory statusParams) =
            checkDACertV2(eigenDAThresholdRegistry, signatureVerifier, cert, securityThresholds, requiredQuorumNumbers);
        revertOnError(status, statusParams);
    }

    /**
     * @notice Checks a complete blob certificate for V2 in a single call
     * @param eigenDAThresholdRegistry The threshold registry contract
     * @param signatureVerifier The signature verifier contract
     * @param cert The blob certificate to verify
     * @param securityThresholds The security thresholds to verify against
     * @param requiredQuorumNumbers The required quorum numbers
     * @return status Error code (SUCCESS if verification succeeded)
     * @return statusParams Additional error parameters
     */
    function checkDACertV2(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDASignatureVerifier signatureVerifier,
        EigenDACertV2 calldata cert,
        SecurityThresholds memory securityThresholds,
        bytes memory requiredQuorumNumbers
    ) internal view returns (StatusCode status, bytes memory statusParams) {
        (status, statusParams) = checkBlobInclusion(cert.batchHeader, cert.blobInclusionInfo);
        if (status != StatusCode.SUCCESS) {
            return (status, statusParams);
        }

        (status, statusParams) = checkSecurityParams(
            eigenDAThresholdRegistry.getBlobParams(cert.blobInclusionInfo.blobCertificate.blobHeader.version),
            securityThresholds
        );
        if (status != StatusCode.SUCCESS) {
            return (status, statusParams);
        }

        // Verify signatures and build confirmed quorums bitmap
        uint256 confirmedQuorumsBitmap;
        (status, statusParams, confirmedQuorumsBitmap) =
            checkSignaturesAndBuildConfirmedQuorums(signatureVerifier, cert, securityThresholds);
        if (status != StatusCode.SUCCESS) {
            return (status, statusParams);
        }

        // Verify blob quorums are a subset of confirmed quorums
        uint256 blobQuorumsBitmap;
        (status, statusParams, blobQuorumsBitmap) = checkBlobQuorumsSubset(
            cert.blobInclusionInfo.blobCertificate.blobHeader.quorumNumbers, confirmedQuorumsBitmap
        );
        if (status != StatusCode.SUCCESS) {
            return (status, statusParams);
        }

        // Verify required quorums are a subset of blob quorums
        return checkRequiredQuorumsSubset(requiredQuorumNumbers, blobQuorumsBitmap);
    }

    /**
     * @notice Checks blob inclusion in the batch using Merkle proof
     * @param batchHeader The batch header
     * @param blobInclusionInfo The blob inclusion info
     * @return status Error code (SUCCESS if verification succeeded)
     * @return statusParams Additional error parameters
     */
    function checkBlobInclusion(BatchHeaderV2 memory batchHeader, BlobInclusionInfo memory blobInclusionInfo)
        internal
        pure
        returns (StatusCode status, bytes memory statusParams)
    {
        bytes32 blobCertHash = EigenDAHasher.hashBlobCertificate(blobInclusionInfo.blobCertificate);
        bytes32 encodedBlobHash = keccak256(abi.encodePacked(blobCertHash));
        bytes32 rootHash = batchHeader.batchRoot;

        bool isValid = Merkle.verifyInclusionKeccak(
            blobInclusionInfo.inclusionProof, rootHash, encodedBlobHash, blobInclusionInfo.blobIndex
        );

        if (isValid) {
            return (StatusCode.SUCCESS, "");
        } else {
            return
                (StatusCode.INVALID_INCLUSION_PROOF, abi.encode(blobInclusionInfo.blobIndex, encodedBlobHash, rootHash));
        }
    }

    /**
     * @notice Checks the security parameters for a blob cert
     * @param blobParams The blob params to verify
     * @param securityThresholds The security thresholds to verify against
     * @return status Error code (SUCCESS if verification succeeded)
     * @return statusParams Additional error parameters
     */
    function checkSecurityParams(VersionedBlobParams memory blobParams, SecurityThresholds memory securityThresholds)
        internal
        pure
        returns (StatusCode status, bytes memory statusParams)
    {
        uint256 gamma = securityThresholds.confirmationThreshold - securityThresholds.adversaryThreshold;
        uint256 n = (10000 - ((1_000_000 / gamma) / uint256(blobParams.codingRate))) * uint256(blobParams.numChunks);
        uint256 minRequired = blobParams.maxNumOperators * 10000;

        if (n >= minRequired) {
            return (StatusCode.SUCCESS, "");
        } else {
            return (StatusCode.SECURITY_ASSUMPTIONS_NOT_MET, abi.encode(gamma, n, minRequired));
        }
    }

    /**
     * @notice Checks quorum signatures and builds a bitmap of confirmed quorums
     * @param signatureVerifier The signature verifier contract
     * @param cert The blob certificate to verify
     * @param securityThresholds The security thresholds to verify against
     * @return status Error code (SUCCESS if verification succeeded)
     * @return statusParams Additional error parameters
     * @return confirmedQuorumsBitmap The bitmap of confirmed quorums
     */
    function checkSignaturesAndBuildConfirmedQuorums(
        IEigenDASignatureVerifier signatureVerifier,
        EigenDACertV2 calldata cert,
        SecurityThresholds memory securityThresholds
    ) internal view returns (StatusCode status, bytes memory statusParams, uint256 confirmedQuorumsBitmap) {
        bytes32 batchHashRoot = EigenDAHasher.hashBatchHeaderV2(cert.batchHeader);

        (QuorumStakeTotals memory quorumStakeTotals,) = signatureVerifier.checkSignatures(
            batchHashRoot,
            cert.signedQuorumNumbers,
            cert.batchHeader.referenceBlockNumber,
            cert.nonSignerStakesAndSignature
        );

        confirmedQuorumsBitmap = 0;

        // Record confirmed quorums where signatories own at least the threshold percentage of the quorum
        for (uint256 i = 0; i < cert.signedQuorumNumbers.length; i++) {
            if (
                quorumStakeTotals.signedStakeForQuorum[i] * THRESHOLD_DENOMINATOR
                    >= quorumStakeTotals.totalStakeForQuorum[i] * securityThresholds.confirmationThreshold
            ) {
                confirmedQuorumsBitmap = BitmapUtils.setBit(confirmedQuorumsBitmap, uint8(cert.signedQuorumNumbers[i]));
            }
        }

        return (StatusCode.SUCCESS, "", confirmedQuorumsBitmap);
    }

    /**
     * @notice Checks that blob quorums are a subset of confirmed quorums
     * @param blobQuorumNumbers The blob quorum numbers
     * @param confirmedQuorumsBitmap The bitmap of confirmed quorums
     * @return status Error code (SUCCESS if verification succeeded)
     * @return statusParams Additional error parameters
     * @return blobQuorumsBitmap The bitmap of blob quorums
     */
    function checkBlobQuorumsSubset(bytes memory blobQuorumNumbers, uint256 confirmedQuorumsBitmap)
        internal
        pure
        returns (StatusCode status, bytes memory statusParams, uint256 blobQuorumsBitmap)
    {
        blobQuorumsBitmap = BitmapUtils.orderedBytesArrayToBitmap(blobQuorumNumbers);

        if (BitmapUtils.isSubsetOf(blobQuorumsBitmap, confirmedQuorumsBitmap)) {
            return (StatusCode.SUCCESS, "", blobQuorumsBitmap);
        } else {
            return (StatusCode.BLOB_QUORUMS_NOT_SUBSET, abi.encode(blobQuorumsBitmap, confirmedQuorumsBitmap), 0);
        }
    }

    /**
     * @notice Checks that required quorums are a subset of blob quorums
     * @param requiredQuorumNumbers The required quorum numbers
     * @param blobQuorumsBitmap The bitmap of blob quorums
     * @return status Error code (SUCCESS if verification succeeded)
     * @return statusParams Additional error parameters
     */
    function checkRequiredQuorumsSubset(bytes memory requiredQuorumNumbers, uint256 blobQuorumsBitmap)
        internal
        pure
        returns (StatusCode status, bytes memory statusParams)
    {
        uint256 requiredQuorumsBitmap = BitmapUtils.orderedBytesArrayToBitmap(requiredQuorumNumbers);

        if (BitmapUtils.isSubsetOf(requiredQuorumsBitmap, blobQuorumsBitmap)) {
            return (StatusCode.SUCCESS, "");
        } else {
            return (StatusCode.REQUIRED_QUORUMS_NOT_SUBSET, abi.encode(requiredQuorumsBitmap, blobQuorumsBitmap));
        }
    }

    /**
     * @notice Gets nonSignerStakesAndSignature for a given signed batch
     * @param registryCoordinator The registry coordinator contract
     * @param signedBatch The signed batch
     * @return nonSignerStakesAndSignature The non-signer stakes and signature
     * @return signedQuorumNumbers The signed quorum numbers
     */
    function getNonSignerStakesAndSignature(IRegistryCoordinator registryCoordinator, SignedBatch memory signedBatch)
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

        CheckSignaturesIndices memory checkSignaturesIndices = getCheckSignaturesIndices(
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

        return (nonSignerStakesAndSignature, signedQuorumNumbers);
    }

    /**
     * @notice Handles error codes by reverting with appropriate custom errors
     * @param status The error code
     * @param statusParams The error parameters
     */
    function revertOnError(StatusCode status, bytes memory statusParams) internal pure {
        if (status == StatusCode.SUCCESS) {
            return; // No error to handle
        }

        if (status == StatusCode.INVALID_INCLUSION_PROOF) {
            (uint256 blobIndex, bytes32 blobHash, bytes32 rootHash) =
                abi.decode(statusParams, (uint256, bytes32, bytes32));
            revert InvalidInclusionProof(blobIndex, blobHash, rootHash);
        } else if (status == StatusCode.SECURITY_ASSUMPTIONS_NOT_MET) {
            (uint256 gamma, uint256 n, uint256 minRequired) = abi.decode(statusParams, (uint256, uint256, uint256));
            revert SecurityAssumptionsNotMet(gamma, n, minRequired);
        } else if (status == StatusCode.BLOB_QUORUMS_NOT_SUBSET) {
            (uint256 blobQuorumsBitmap, uint256 confirmedQuorumsBitmap) = abi.decode(statusParams, (uint256, uint256));
            revert BlobQuorumsNotSubset(blobQuorumsBitmap, confirmedQuorumsBitmap);
        } else if (status == StatusCode.REQUIRED_QUORUMS_NOT_SUBSET) {
            (uint256 requiredQuorumsBitmap, uint256 blobQuorumsBitmap) = abi.decode(statusParams, (uint256, uint256));
            revert RequiredQuorumsNotSubset(requiredQuorumsBitmap, blobQuorumsBitmap);
        } else {
            revert("Unknown error code");
        }
    }

    struct CheckSignaturesIndices {
        uint32[] nonSignerQuorumBitmapIndices;
        uint32[] quorumApkIndices;
        uint32[] totalStakeIndices;
        uint32[][] nonSignerStakeIndices; // nonSignerStakeIndices[quorumNumberIndex][nonSignerIndex]
    }

    /**
     * @notice copied from OperatorStateRetriever.sol (which should probably be a library in the first place)
     * @param registryCoordinator is the registry coordinator to fetch the AVS registry information from
     * @param referenceBlockNumber is the block number to get the indices for
     * @param quorumNumbers are the ids of the quorums to get the operator state for
     * @param nonSignerOperatorIds are the ids of the nonsigning operators
     * @return 1) the indices of the quorumBitmaps for each of the operators in the @param nonSignerOperatorIds array at the given blocknumber
     *         2) the indices of the total stakes entries for the given quorums at the given blocknumber
     *         3) the indices of the stakes of each of the nonsigners in each of the quorums they were a
     *            part of (for each nonsigner, an array of length the number of quorums they were a part of
     *            that are also part of the provided quorumNumbers) at the given blocknumber
     *         4) the indices of the quorum apks for each of the provided quorums at the given blocknumber
     */
    function getCheckSignaturesIndices(
        IRegistryCoordinator registryCoordinator,
        uint32 referenceBlockNumber,
        bytes memory quorumNumbers,
        bytes32[] memory nonSignerOperatorIds
    ) internal view returns (CheckSignaturesIndices memory) {
        IStakeRegistry stakeRegistry = registryCoordinator.stakeRegistry();
        CheckSignaturesIndices memory checkSignaturesIndices;

        // get the indices of the quorumBitmap updates for each of the operators in the nonSignerOperatorIds array
        checkSignaturesIndices.nonSignerQuorumBitmapIndices =
            registryCoordinator.getQuorumBitmapIndicesAtBlockNumber(referenceBlockNumber, nonSignerOperatorIds);

        // get the indices of the totalStake updates for each of the quorums in the quorumNumbers array
        checkSignaturesIndices.totalStakeIndices =
            stakeRegistry.getTotalStakeIndicesAtBlockNumber(referenceBlockNumber, quorumNumbers);

        checkSignaturesIndices.nonSignerStakeIndices = new uint32[][](quorumNumbers.length);
        for (uint8 quorumNumberIndex = 0; quorumNumberIndex < quorumNumbers.length; quorumNumberIndex++) {
            uint256 numNonSignersForQuorum = 0;
            // this array's length will be at most the number of nonSignerOperatorIds, this will be trimmed after it is filled
            checkSignaturesIndices.nonSignerStakeIndices[quorumNumberIndex] = new uint32[](nonSignerOperatorIds.length);

            for (uint256 i = 0; i < nonSignerOperatorIds.length; i++) {
                // get the quorumBitmap for the operator at the given blocknumber and index
                uint192 nonSignerQuorumBitmap = registryCoordinator.getQuorumBitmapAtBlockNumberByIndex(
                    nonSignerOperatorIds[i],
                    referenceBlockNumber,
                    checkSignaturesIndices.nonSignerQuorumBitmapIndices[i]
                );

                require(
                    nonSignerQuorumBitmap != 0,
                    "EigenDACertVerificationV2Lib.getCheckSignaturesIndices: operator must be registered at blocknumber"
                );

                // if the operator was a part of the quorum and the quorum is a part of the provided quorumNumbers
                if ((nonSignerQuorumBitmap >> uint8(quorumNumbers[quorumNumberIndex])) & 1 == 1) {
                    // get the index of the stake update for the operator at the given blocknumber and quorum number
                    checkSignaturesIndices.nonSignerStakeIndices[quorumNumberIndex][numNonSignersForQuorum] =
                    stakeRegistry.getStakeUpdateIndexAtBlockNumber(
                        nonSignerOperatorIds[i], uint8(quorumNumbers[quorumNumberIndex]), referenceBlockNumber
                    );
                    numNonSignersForQuorum++;
                }
            }

            // resize the array to the number of nonSigners for this quorum
            uint32[] memory nonSignerStakeIndicesForQuorum = new uint32[](numNonSignersForQuorum);
            for (uint256 i = 0; i < numNonSignersForQuorum; i++) {
                nonSignerStakeIndicesForQuorum[i] = checkSignaturesIndices.nonSignerStakeIndices[quorumNumberIndex][i];
            }
            checkSignaturesIndices.nonSignerStakeIndices[quorumNumberIndex] = nonSignerStakeIndicesForQuorum;
        }

        IBLSApkRegistry blsApkRegistry = registryCoordinator.blsApkRegistry();
        // get the indices of the quorum apks for each of the provided quorums at the given blocknumber
        checkSignaturesIndices.quorumApkIndices =
            blsApkRegistry.getApkIndicesAtBlockNumber(quorumNumbers, referenceBlockNumber);

        return checkSignaturesIndices;
    }
}
