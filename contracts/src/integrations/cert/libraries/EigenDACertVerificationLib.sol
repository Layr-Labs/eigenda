// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";
import {BN254} from "lib/eigenlayer-middleware/src/libraries/BN254.sol";
import {Merkle} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/libraries/Merkle.sol";
import {BitmapUtils} from "lib/eigenlayer-middleware/src/libraries/BitmapUtils.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {EigenDATypesV2 as DATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";

import {EigenDACertTypes as CT} from "src/integrations/cert/EigenDACertTypes.sol";

/// @title EigenDACertVerificationLib
/// @notice Library for verifying EigenDA certificates
library EigenDACertVerificationLib {
    /// @notice Denominator used for threshold percentage calculations (100 for percentages)
    uint256 internal constant THRESHOLD_DENOMINATOR = 100;

    /// @notice The maximum number of quorums this contract supports
    uint256 internal constant MAX_QUORUM_COUNT = 5;

    /// @notice The maximum number of nonsigner this contract supports. The count can contain duplicate count if
    ///         an operator belongs to multiple quorums
    uint256 internal constant MAX_NONSIGNER_COUNT_ALL_QUORUM = 415;

    /// @notice Thrown when the inclusion proof is invalid
    /// @param blobIndex The index of the blob in the batch
    /// @param blobHash The hash of the blob certificate
    /// @param rootHash The root hash of the merkle tree
    error InvalidInclusionProof(uint32 blobIndex, bytes32 blobHash, bytes32 rootHash);

    /// @notice Thrown when security assumptions are not met
    /// @param confirmationThreshold The confirmation threshold percentage
    /// @param adversaryThreshold The safety threshold percentage
    /// @param codingRate The coding rate for the blob
    /// @param numChunks The number of chunks in the blob
    /// @param maxNumOperators The maximum number of operators
    error SecurityAssumptionsNotMet(
        uint8 confirmationThreshold,
        uint8 adversaryThreshold,
        uint8 codingRate,
        uint32 numChunks,
        uint32 maxNumOperators
    );

    /// @notice Thrown when blob quorums are not a subset of confirmed quorums
    /// @param blobQuorumsBitmap The bitmap of blob quorums
    /// @param confirmedQuorumsBitmap The bitmap of confirmed quorums
    error BlobQuorumsNotSubset(uint256 blobQuorumsBitmap, uint256 confirmedQuorumsBitmap);

    /// @notice Thrown when required quorums are not a subset of blob quorums
    /// @param requiredQuorumsBitmap The bitmap of required quorums
    /// @param blobQuorumsBitmap The bitmap of blob quorums
    error RequiredQuorumsNotSubset(uint256 requiredQuorumsBitmap, uint256 blobQuorumsBitmap);

    /// @notice Thrown when the blob version is invalid (doesn't exist in the ThresholdRegistry contract)
    /// @param blobVersion The invalid blob version
    /// @param nextBlobVersion The next blob version (valid versions need to be less than this number)
    error InvalidBlobVersion(uint16 blobVersion, uint16 nextBlobVersion);

    /// @notice Thrown when the offchain derivation version is invalid
    /// @param certDerivationVer The offchain derivation version in the certificate
    /// @param requiredDerivationVer The required offchain derivation version
    error InvalidOffchainDerivationVersion(uint16 certDerivationVer, uint16 requiredDerivationVer);

    /// @notice Thrown when the number of signed quorums exceeds the maximum allowed
    /// @param count The actual number of signed quorums provided
    error QuorumCountExceedsMaximum(uint256 count);

    /// @notice Thrown when the total number of non-signers across all quorums exceeds the maximum allowed
    /// @param count The total number of non-signers counted
    error NonSignerCountExceedsMaximum(uint256 count);

    /// @notice Checks a DA certificate using all parameters that a CertVerifier has registered, and returns a status.
    /// @param eigenDAThresholdRegistry The threshold registry contract
    /// @param eigenDASignatureVerifier The signature verifier contract
    /// @param daCert The EigenDA certificate
    /// @param securityThresholds The security thresholds to verify against
    /// Callers should ensure that the requiredQuorumNumbers passed are non-empty if needed.
    /// @param requiredQuorumNumbers The required quorum numbers. Can be empty if not required.
    /// @param offchainDerivationVersion The offchain derivation version to verify against
    function checkDACert(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDASignatureVerifier eigenDASignatureVerifier,
        CT.EigenDACertV4 memory daCert,
        DATypesV1.SecurityThresholds memory securityThresholds,
        bytes memory requiredQuorumNumbers,
        uint16 offchainDerivationVersion
    ) internal view {
        checkOffchainDerivationVersion(daCert.offchainDerivationVersion, offchainDerivationVersion);

        // verifying merkle inclusion proof is very efficient, even assuming the worst depth 256.
        // A single depth verification takes about 300 gas for KECCAK256 and CALLDATALOAD
        // so at worst 80K Gas.
        checkBlobInclusion(daCert.batchHeader, daCert.blobInclusionInfo);

        checkSecurityParams(
            eigenDAThresholdRegistry, daCert.blobInclusionInfo.blobCertificate.blobHeader.version, securityThresholds
        );

        // Verify signatures and build confirmed quorums bitmap
        uint256 confirmedQuorumsBitmap = checkSignaturesAndBuildConfirmedQuorums(
            eigenDASignatureVerifier,
            hashBatchHeaderV2(daCert.batchHeader),
            daCert.signedQuorumNumbers,
            daCert.batchHeader.referenceBlockNumber,
            daCert.nonSignerStakesAndSignature,
            securityThresholds
        );

        // The different quorums are related by: requiredQuorums ⊆ blobQuorums ⊆ confirmedQuorums ⊆ signedQuorums
        // checkSignaturesAndBuildConfirmedQuorums checked the last inequality. We now verify the other two.
        checkQuorumSubsets(
            requiredQuorumNumbers,
            daCert.blobInclusionInfo.blobCertificate.blobHeader.quorumNumbers,
            confirmedQuorumsBitmap
        );
    }

    /// @notice Checks blob inclusion in the batch using Merkle proof
    /// @param batchHeader The batch header
    /// @param blobInclusionInfo The blob inclusion info
    function checkBlobInclusion(
        DATypesV2.BatchHeaderV2 memory batchHeader,
        DATypesV2.BlobInclusionInfo memory blobInclusionInfo
    ) internal pure {
        bytes32 blobCertHash = hashBlobCertificate(blobInclusionInfo.blobCertificate);
        bytes32 encodedBlobHash = keccak256(abi.encodePacked(blobCertHash));
        bytes32 rootHash = batchHeader.batchRoot;

        bool isValid = Merkle.verifyInclusionKeccak(
            blobInclusionInfo.inclusionProof, rootHash, encodedBlobHash, blobInclusionInfo.blobIndex
        );

        if (!isValid) {
            revert InvalidInclusionProof(blobInclusionInfo.blobIndex, encodedBlobHash, rootHash);
        }
    }

    /// @notice Checks the security parameters for a blob cert
    /// @dev Verifies that the security condition
    ///      (confirmationThreshold - adversaryThreshold > reconstructionThreshold)
    ///      holds, by checking an invariant.
    ///      If the inequality fails, the blob is considered insecure.
    /// @param eigenDAThresholdRegistry The threshold registry contract
    /// @param blobVersion The blob version to verify
    /// @param securityThresholds The security thresholds to verify against
    function checkSecurityParams(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        uint16 blobVersion,
        DATypesV1.SecurityThresholds memory securityThresholds
    ) internal view {
        // We validate that the cert's blob_version is valid. Otherwise the getBlobParams call below
        // would return a codingRate=0 which will cause a divide by 0 error below.
        uint16 nextBlobVersion = eigenDAThresholdRegistry.nextBlobVersion();
        if (blobVersion >= nextBlobVersion) {
            revert InvalidBlobVersion(blobVersion, nextBlobVersion);
        }
        DATypesV1.VersionedBlobParams memory blobParams = eigenDAThresholdRegistry.getBlobParams(blobVersion);

        // Check for potential underflow:
        // maxNumOperators must not exceed numChunks
        //
        if (
            blobParams.maxNumOperators > blobParams.numChunks
                || securityThresholds.confirmationThreshold < securityThresholds.adversaryThreshold
        ) {
            revert SecurityAssumptionsNotMet(
                securityThresholds.confirmationThreshold,
                securityThresholds.adversaryThreshold,
                blobParams.codingRate,
                blobParams.numChunks,
                blobParams.maxNumOperators
            );
        }

        uint256 lhs = blobParams.codingRate * (blobParams.numChunks - blobParams.maxNumOperators)
            * (securityThresholds.confirmationThreshold - securityThresholds.adversaryThreshold);
        uint256 rhs = 100 * blobParams.numChunks;

        if (!(lhs >= rhs)) {
            revert SecurityAssumptionsNotMet(
                securityThresholds.confirmationThreshold,
                securityThresholds.adversaryThreshold,
                blobParams.codingRate,
                blobParams.numChunks,
                blobParams.maxNumOperators
            );
        }
    }

    /// @notice Checks quorum signatures and builds a bitmap of confirmed quorums
    /// @param signatureVerifier The signature verifier contract
    /// @param batchHashRoot The hash of the batch header
    /// @param signedQuorumNumbers The signed quorum numbers
    /// @param referenceBlockNumber The reference block number
    /// @param nonSignerStakesAndSignature The non-signer stakes and signature
    /// @param securityThresholds The security thresholds to verify against
    /// @return confirmedQuorumsBitmap The bitmap of confirmed quorums
    function checkSignaturesAndBuildConfirmedQuorums(
        IEigenDASignatureVerifier signatureVerifier,
        bytes32 batchHashRoot,
        bytes memory signedQuorumNumbers,
        uint32 referenceBlockNumber,
        DATypesV1.NonSignerStakesAndSignature memory nonSignerStakesAndSignature,
        DATypesV1.SecurityThresholds memory securityThresholds
    ) internal view returns (uint256 confirmedQuorumsBitmap) {
        // the maximal supported number of quorum is 192
        // https://github.com/Layr-Labs/eigenda/blob/00cc8868b7e2d742fc6584dc1dea312193c8d4c2/contracts/src/core/EigenDARegistryCoordinatorStorage.sol#L36
        if (signedQuorumNumbers.length > MAX_QUORUM_COUNT) {
            revert QuorumCountExceedsMaximum(signedQuorumNumbers.length);
        }

        // if an nonsigning operator belongs to multiple quorums, the totalNonsignersCount counts it multiple times
        uint256 totalNonsignersCount = 0;
        for (uint256 i = 0; i < nonSignerStakesAndSignature.nonSignerStakeIndices.length; i++) {
            totalNonsignersCount += nonSignerStakesAndSignature.nonSignerStakeIndices[i].length;
        }
        if (totalNonsignersCount > MAX_NONSIGNER_COUNT_ALL_QUORUM) {
            revert NonSignerCountExceedsMaximum(totalNonsignersCount);
        }

        (DATypesV1.QuorumStakeTotals memory quorumStakeTotals,) = signatureVerifier.checkSignatures(
            batchHashRoot, signedQuorumNumbers, referenceBlockNumber, nonSignerStakesAndSignature
        );

        confirmedQuorumsBitmap = 0;

        // Record confirmed quorums where signatories own at least the threshold percentage of the quorum
        for (uint256 i = 0; i < signedQuorumNumbers.length; i++) {
            if (
                quorumStakeTotals.signedStakeForQuorum[i] * THRESHOLD_DENOMINATOR
                    >= quorumStakeTotals.totalStakeForQuorum[i] * securityThresholds.confirmationThreshold
            ) {
                confirmedQuorumsBitmap = BitmapUtils.setBit(confirmedQuorumsBitmap, uint8(signedQuorumNumbers[i]));
            }
        }

        return confirmedQuorumsBitmap;
    }

    /// @notice Checks that requiredQuorums ⊆ blobQuorums ⊆ confirmedQuorums
    /// @param requiredQuorumNumbers The required quorum numbers
    /// @param blobQuorumNumbers The blob quorum numbers, which are the quorums requested in the blobHeader part of the dispersal
    /// @param confirmedQuorumsBitmap The bitmap of confirmed quorums, which are signed quorums that meet the confirmationThreshold
    function checkQuorumSubsets(
        bytes memory requiredQuorumNumbers,
        bytes memory blobQuorumNumbers,
        uint256 confirmedQuorumsBitmap
    ) internal pure {
        uint256 blobQuorumsBitmap = BitmapUtils.orderedBytesArrayToBitmap(blobQuorumNumbers);
        if (!BitmapUtils.isSubsetOf(blobQuorumsBitmap, confirmedQuorumsBitmap)) {
            revert BlobQuorumsNotSubset(blobQuorumsBitmap, confirmedQuorumsBitmap);
        }

        uint256 requiredQuorumsBitmap = BitmapUtils.orderedBytesArrayToBitmap(requiredQuorumNumbers);
        if (!BitmapUtils.isSubsetOf(requiredQuorumsBitmap, blobQuorumsBitmap)) {
            revert RequiredQuorumsNotSubset(requiredQuorumsBitmap, blobQuorumsBitmap);
        }
    }

    /// @notice Checks that the offchain derivation version matches the required version
    /// @param certDerivationVer The offchain derivation version in the certificate
    /// @param requiredDerivationVer The required offchain derivation version
    function checkOffchainDerivationVersion(uint16 certDerivationVer, uint16 requiredDerivationVer) internal pure {
        if (certDerivationVer != requiredDerivationVer) {
            revert InvalidOffchainDerivationVersion(certDerivationVer, requiredDerivationVer);
        }
    }

    /// @notice Gets nonSignerStakesAndSignature for a given signed batch
    /// @param operatorStateRetriever The operator state retriever contract
    /// @param registryCoordinator The registry coordinator contract
    /// @param signedBatch The signed batch
    /// @return nonSignerStakesAndSignature The non-signer stakes and signature
    /// @return signedQuorumNumbers The signed quorum numbers
    function getNonSignerStakesAndSignature(
        OperatorStateRetriever operatorStateRetriever,
        IRegistryCoordinator registryCoordinator,
        DATypesV2.SignedBatch memory signedBatch
    )
        internal
        view
        returns (
            DATypesV1.NonSignerStakesAndSignature memory nonSignerStakesAndSignature,
            bytes memory signedQuorumNumbers
        )
    {
        bytes32[] memory nonSignerOperatorIds = new bytes32[](signedBatch.attestation.nonSignerPubkeys.length);
        for (uint256 i = 0; i < signedBatch.attestation.nonSignerPubkeys.length; ++i) {
            nonSignerOperatorIds[i] = BN254.hashG1Point(signedBatch.attestation.nonSignerPubkeys[i]);
        }

        for (uint256 i = 0; i < signedBatch.attestation.quorumNumbers.length; ++i) {
            signedQuorumNumbers = abi.encodePacked(signedQuorumNumbers, uint8(signedBatch.attestation.quorumNumbers[i]));
        }

        OperatorStateRetriever.CheckSignaturesIndices memory checkSignaturesIndices =
            operatorStateRetriever.getCheckSignaturesIndices(
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

        return (nonSignerStakesAndSignature, signedQuorumNumbers);
    }

    /// @notice hashes the given V2 batch header
    /// @param batchHeader the V2 batch header to hash
    function hashBatchHeaderV2(DATypesV2.BatchHeaderV2 memory batchHeader) internal pure returns (bytes32) {
        return keccak256(abi.encode(batchHeader));
    }

    /// @notice hashes the given V2 blob header
    /// @param blobHeader the V2 blob header to hash
    function hashBlobHeaderV2(DATypesV2.BlobHeaderV2 memory blobHeader) internal pure returns (bytes32) {
        return keccak256(
            abi.encode(
                keccak256(abi.encode(blobHeader.version, blobHeader.quorumNumbers, blobHeader.commitment)),
                blobHeader.paymentHeaderHash
            )
        );
    }

    /// @notice hashes the given V2 blob certificate
    /// @param blobCertificate the V2 blob certificate to hash
    function hashBlobCertificate(DATypesV2.BlobCertificate memory blobCertificate) internal pure returns (bytes32) {
        return keccak256(
            abi.encode(
                hashBlobHeaderV2(blobCertificate.blobHeader), blobCertificate.signature, blobCertificate.relayKeys
            )
        );
    }
}
