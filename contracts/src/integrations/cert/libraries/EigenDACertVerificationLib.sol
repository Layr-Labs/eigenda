// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";
import {IEigenDARelayRegistry} from "src/core/interfaces/IEigenDARelayRegistry.sol";
import {BN254} from "lib/eigenlayer-middleware/src/libraries/BN254.sol";
import {Merkle} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/libraries/Merkle.sol";
import {BitmapUtils} from "lib/eigenlayer-middleware/src/libraries/BitmapUtils.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry} from "lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {IBLSApkRegistry} from "lib/eigenlayer-middleware/src/interfaces/IBLSApkRegistry.sol";
import {EigenDATypesV2 as DATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";

import {EigenDACertTypes as CT} from "src/integrations/cert/EigenDACertTypes.sol";

/// @title EigenDACertVerificationLib
/// @notice Library for verifying EigenDA certificates
library EigenDACertVerificationLib {
    /// @notice Denominator used for threshold percentage calculations (100 for percentages)
    uint256 internal constant THRESHOLD_DENOMINATOR = 100;

    // TODO: discuss whether we really need these errors + revertOnError fct.
    // Given that we're trying to write a library that doesn't revert,
    // the only use case for these is to try/catch on custom errors in the CertVerifier,
    // but is that really how we want to architect this?
    // For example, we have the revertOnError internal function, but its never called,
    // so these errors are never thrown right now...

    /// @notice Thrown when the inclusion proof is not a multiple of 32
    /// @param length Length of the inclusion proof
    error InclusionProofNotMultipleOf32(uint256 length);

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

    /// @notice Thrown when certificate decoding from bytes into struct fails
    error CertDecodeRevert();

    /// @notice Thrown when the external call to the signature verifier reverts
    error SignatureVerificationCallRevert(bytes reason);

    /// @notice Status codes for certificate verification results.
    /// @dev Returned by checkDACert and checkDACertV2 functions, along with an errParam bytes array for additional context.
    ///      See the revertOnError function for what errParams contains for each error code.
    enum StatusCode {
        NULL_ERROR, // Unused error code. If this is returned, there is a bug in the code.
        SUCCESS, // Verification succeeded
        INVALID_INCLUSION_PROOF, // Merkle inclusion proof is invalid
        SECURITY_ASSUMPTIONS_NOT_MET, // Security assumptions not met
        BLOB_QUORUMS_NOT_SUBSET, // Blob quorums not a subset of confirmed quorums
        REQUIRED_QUORUMS_NOT_SUBSET, // Required quorums not a subset of blob quorums
        INCLUSION_PROOF_NOT_MULTIPLE_32, // Merkle inclusion proof not right length
        EMPTY_BLOB_QUORUMS, // Certificate quorums are empty.
        UNORDERED_OR_DUPLICATE_QUORUMS, // Quorum numbers are not ordered, or contain duplicates
        GREATER_THAN_256_QUORUMS, // Quorum numbers exceed 256
        CERT_DECODE_REVERT, // Certificate abi.decoding reverted
        SIGNATURE_VERIFICATION_CALL_REVERT, // External call to signature verifier reverted
        INVALID_BLOB_VERSION // Blob version does not exist onchain
    }

    /// @notice Checks a DA certificate using all parameters that a CertVerifier has registered, and returns a status.
    /// @dev Uses the same verification logic as verifyDACertV2. The only difference is that the certificate is ABI encoded bytes.
    /// @param eigenDAThresholdRegistry The threshold registry contract
    /// @param eigenDASignatureVerifier The signature verifier contract
    /// @param daCert The EigenDA certificate
    /// @param securityThresholds The security thresholds to verify against
    /// @param requiredQuorumNumbers The required quorum numbers
    /// @return status Status code (SUCCESS if verification succeeded)
    /// @return statusParams Additional status parameters
    function checkDACert(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDASignatureVerifier eigenDASignatureVerifier,
        CT.EigenDACertV3 memory daCert,
        DATypesV1.SecurityThresholds memory securityThresholds,
        bytes memory requiredQuorumNumbers
    ) internal view returns (StatusCode, bytes memory) {
        return checkDACertV2(
            eigenDAThresholdRegistry,
            eigenDASignatureVerifier,
            daCert.batchHeader,
            daCert.blobInclusionInfo,
            daCert.nonSignerStakesAndSignature,
            securityThresholds,
            requiredQuorumNumbers,
            daCert.signedQuorumNumbers
        );
    }

    /**
     * @notice Checks a complete blob certificate for V2 in a single call
     * @param eigenDAThresholdRegistry The threshold registry contract
     * @param signatureVerifier The signature verifier contract
     * @param batchHeader The batch header
     * @param blobInclusionInfo The blob inclusion info
     * @param nonSignerStakesAndSignature The non-signer stakes and signature
     * @param securityThresholds The security thresholds to verify against
     * @param requiredQuorumNumbers The required quorum numbers. This library does not require these to be non-empty.
     * Callers should ensure that the requiredQuorumNumbers passed are non-empty if required.
     * @param signedQuorumNumbers The signed quorum numbers
     * @return err Error code (SUCCESS if verification succeeded)
     * @return errParams Additional error parameters
     */
    function checkDACertV2(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDASignatureVerifier signatureVerifier,
        DATypesV2.BatchHeaderV2 memory batchHeader,
        DATypesV2.BlobInclusionInfo memory blobInclusionInfo,
        DATypesV1.NonSignerStakesAndSignature memory nonSignerStakesAndSignature,
        DATypesV1.SecurityThresholds memory securityThresholds,
        bytes memory requiredQuorumNumbers,
        bytes memory signedQuorumNumbers
    ) internal view returns (StatusCode err, bytes memory errParams) {
        (err, errParams) = checkBlobInclusion(batchHeader, blobInclusionInfo);
        if (err != StatusCode.SUCCESS) {
            return (err, errParams);
        }

        // We validate that the cert's blob_version is valid. Otherwise the getBlobParams call below
        // would return a codingRate=0 which will cause a divide by 0 error in checkSecurityParams.
        uint16 nextBlobVersion = eigenDAThresholdRegistry.nextBlobVersion();
        if (blobInclusionInfo.blobCertificate.blobHeader.version >= nextBlobVersion) {
            return (StatusCode.INVALID_BLOB_VERSION, abi.encode(blobInclusionInfo.blobCertificate.blobHeader.version));
        }

        (err, errParams) = checkSecurityParams(
            eigenDAThresholdRegistry.getBlobParams(blobInclusionInfo.blobCertificate.blobHeader.version),
            securityThresholds
        );
        if (err != StatusCode.SUCCESS) {
            return (err, errParams);
        }

        // Verify signatures and build confirmed quorums bitmap
        uint256 confirmedQuorumsBitmap;
        (err, errParams, confirmedQuorumsBitmap) = checkSignaturesAndBuildConfirmedQuorums(
            signatureVerifier,
            hashBatchHeaderV2(batchHeader),
            signedQuorumNumbers,
            batchHeader.referenceBlockNumber,
            nonSignerStakesAndSignature,
            securityThresholds
        );
        if (err != StatusCode.SUCCESS) {
            return (err, errParams);
        }

        // Verify blob quorums are a subset of confirmed quorums
        uint256 blobQuorumsBitmap;
        (err, errParams, blobQuorumsBitmap) =
            checkBlobQuorumsSubset(blobInclusionInfo.blobCertificate.blobHeader.quorumNumbers, confirmedQuorumsBitmap);
        if (err != StatusCode.SUCCESS) {
            return (err, errParams);
        }

        // Verify required quorums are a subset of blob quorums
        return checkRequiredQuorumsSubset(requiredQuorumNumbers, blobQuorumsBitmap);
    }

    /**
     * @notice Checks blob inclusion in the batch using Merkle proof
     * @param batchHeader The batch header
     * @param blobInclusionInfo The blob inclusion info
     * @return err Error code (SUCCESS if verification succeeded)
     * @return errParams Additional error parameters
     */
    function checkBlobInclusion(
        DATypesV2.BatchHeaderV2 memory batchHeader,
        DATypesV2.BlobInclusionInfo memory blobInclusionInfo
    ) internal pure returns (StatusCode err, bytes memory errParams) {
        bytes32 blobCertHash = hashBlobCertificate(blobInclusionInfo.blobCertificate);
        bytes32 encodedBlobHash = keccak256(abi.encodePacked(blobCertHash));
        bytes32 rootHash = batchHeader.batchRoot;

        // Explicitly check this before calling verifyInclusionKeccak which reverts.
        if (blobInclusionInfo.inclusionProof.length % 32 != 0) {
            return (StatusCode.INCLUSION_PROOF_NOT_MULTIPLE_32, abi.encode(blobInclusionInfo.inclusionProof.length));
        }
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
     * @return err Error code (SUCCESS if verification succeeded)
     * @return errParams Additional error parameters
     * @dev Checks the invariant `numChunks * (1 - 100/gamma/codingRate) >= maxNumOperators`
     */
    function checkSecurityParams(
        DATypesV1.VersionedBlobParams memory blobParams,
        DATypesV1.SecurityThresholds memory securityThresholds
    ) internal pure returns (StatusCode err, bytes memory errParams) {
        // In order to not revert, we need gamma > 0 and codingRate > 0.
        // We assume here that the CertVerifier constructor checked that confirmationThreshold > adversaryThreshold.
        // We also assume that the blobParams passed are from a valid version.
        // Thus, dividing by codingRate below will only panic if codingRate of a proper initialized version is 0,
        // which is either a configuration bug, or a malicious attack. In both cases, we cannot tell whether the
        // cert is valid or invalid, so it is ok to panic and let social consensus intervene (put a human debugger in the loop).
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
     * @param batchHashRoot The hash of the batch header
     * @param signedQuorumNumbers The signed quorum numbers
     * @param referenceBlockNumber The reference block number
     * @param nonSignerStakesAndSignature The non-signer stakes and signature
     * @param securityThresholds The security thresholds to verify against
     * @return err Error code (SUCCESS if verification succeeded)
     * @return errParams Additional error parameters
     * @return confirmedQuorumsBitmap The bitmap of confirmed quorums
     */
    function checkSignaturesAndBuildConfirmedQuorums(
        IEigenDASignatureVerifier signatureVerifier,
        bytes32 batchHashRoot,
        bytes memory signedQuorumNumbers,
        uint32 referenceBlockNumber,
        DATypesV1.NonSignerStakesAndSignature memory nonSignerStakesAndSignature,
        DATypesV1.SecurityThresholds memory securityThresholds
    ) internal view returns (StatusCode err, bytes memory errParams, uint256 confirmedQuorumsBitmap) {
        try signatureVerifier.checkSignatures(
            batchHashRoot, signedQuorumNumbers, referenceBlockNumber, nonSignerStakesAndSignature
        ) returns (DATypesV1.QuorumStakeTotals memory quorumStakeTotals, bytes32) {
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

            return (StatusCode.SUCCESS, "", confirmedQuorumsBitmap);
        } catch Error(string memory reason) {
            // This would match any require(..., "string reason") revert that is pre custom errors,
            // which earlier versions of BLSSignatureChecker used, and might still be deployed. See:
            // https://github.com/Layr-Labs/eigenlayer-middleware/blob/fe5834371caed60c1d26ab62b5519b0cbdcb42fa/src/BLSSignatureChecker.sol#L96
            return (StatusCode.SIGNATURE_VERIFICATION_CALL_REVERT, bytes(reason), 0);
        } catch (bytes memory reason) {
            if (reason.length < 4) {
                // We re-throw any non custom-error that was caught here. For example,
                // low-level evm reverts such as out-of-gas don't return any data.
                // See https://rareskills.io/post/try-catch-solidity#gdvnie-9-what-gets-returned-during-an-out-of-gas?
                // These generally mean there is a bug in our implementation, which should be addressed by a human debugger.
                // TODO: figure out whether we can programmatically deal with out of gas, since that might happen from
                // a maliciously constructed cert.
                revert(string(reason));
            }
            // We assume that any revert here is coming from a failed require(..., SomeCustomError()) statement
            // TODO: make sure that this doesn't catch failing asserts, panics, or other low-level evm reverts like out of gas.
            return (StatusCode.SIGNATURE_VERIFICATION_CALL_REVERT, reason, 0);
        }
    }

    /**
     * @notice Checks that blob quorums requested as part of the dispersal are a subset of confirmed quorums
     * @param blobQuorumNumbers The blob quorum numbers
     * @param confirmedQuorumsBitmap The bitmap of confirmed quorums
     * @return err Error code (SUCCESS if verification succeeded)
     * @return errParams Additional error parameters
     * @return blobQuorumsBitmap The bitmap of blob quorums
     */
    function checkBlobQuorumsSubset(bytes memory blobQuorumNumbers, uint256 confirmedQuorumsBitmap)
        internal
        pure
        returns (StatusCode err, bytes memory errParams, uint256 blobQuorumsBitmap)
    {
        // A blobCert containing no quorum numbers means the cert is invalid, so we return a status code.
        if (blobQuorumNumbers.length == 0) {
            return (StatusCode.EMPTY_BLOB_QUORUMS, "", 0);
        }

        (err, errParams, blobQuorumsBitmap) = orderedQuorumsToBitmap(blobQuorumNumbers);
        if (err != StatusCode.SUCCESS) {
            return (err, errParams, 0);
        }

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
     * @return err Error code (SUCCESS if verification succeeded)
     * @return errParams Additional error parameters
     */
    function checkRequiredQuorumsSubset(bytes memory requiredQuorumNumbers, uint256 blobQuorumsBitmap)
        internal
        pure
        returns (StatusCode err, bytes memory errParams)
    {
        uint256 requiredQuorumsBitmap;
        (err, errParams, requiredQuorumsBitmap) = orderedQuorumsToBitmap(requiredQuorumNumbers);
        if (err != StatusCode.SUCCESS) {
            return (err, errParams);
        }

        if (BitmapUtils.isSubsetOf(requiredQuorumsBitmap, blobQuorumsBitmap)) {
            return (StatusCode.SUCCESS, "");
        } else {
            return (StatusCode.REQUIRED_QUORUMS_NOT_SUBSET, abi.encode(requiredQuorumsBitmap, blobQuorumsBitmap));
        }
    }

    /**
     * @notice Gets nonSignerStakesAndSignature for a given signed batch
     * @param operatorStateRetriever The operator state retriever contract
     * @param registryCoordinator The registry coordinator contract
     * @param signedBatch The signed batch
     * @return nonSignerStakesAndSignature The non-signer stakes and signature
     * @return signedQuorumNumbers The signed quorum numbers
     */
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

        return (nonSignerStakesAndSignature, signedQuorumNumbers);
    }

    /**
     * @notice Handles error codes by reverting with appropriate custom errors
     * @param err The error code
     * @param errParams The error parameters
     * @dev This function is not meant to be called, but is exposed here as a schema of the different error codes
     *  returned by this library, and what data their accompanying errParams contain.
     */
    function revertOnError(StatusCode err, bytes memory errParams) internal pure {
        if (err == StatusCode.SUCCESS) {
            return; // No error to handle
        }

        if (err == StatusCode.INVALID_INCLUSION_PROOF) {
            (uint256 blobIndex, bytes32 blobHash, bytes32 rootHash) = abi.decode(errParams, (uint256, bytes32, bytes32));
            revert InvalidInclusionProof(blobIndex, blobHash, rootHash);
        } else if (err == StatusCode.INCLUSION_PROOF_NOT_MULTIPLE_32) {
            uint256 inclusionProofLen = abi.decode(errParams, (uint256));
            revert InclusionProofNotMultipleOf32(inclusionProofLen);
        } else if (err == StatusCode.SECURITY_ASSUMPTIONS_NOT_MET) {
            (uint256 gamma, uint256 n, uint256 minRequired) = abi.decode(errParams, (uint256, uint256, uint256));
            revert SecurityAssumptionsNotMet(gamma, n, minRequired);
        } else if (err == StatusCode.BLOB_QUORUMS_NOT_SUBSET) {
            (uint256 blobQuorumsBitmap, uint256 confirmedQuorumsBitmap) = abi.decode(errParams, (uint256, uint256));
            revert BlobQuorumsNotSubset(blobQuorumsBitmap, confirmedQuorumsBitmap);
        } else if (err == StatusCode.REQUIRED_QUORUMS_NOT_SUBSET) {
            (uint256 requiredQuorumsBitmap, uint256 blobQuorumsBitmap) = abi.decode(errParams, (uint256, uint256));
            revert RequiredQuorumsNotSubset(requiredQuorumsBitmap, blobQuorumsBitmap);
        } else if (err == StatusCode.CERT_DECODE_REVERT) {
            revert CertDecodeRevert();
        } else if (err == StatusCode.SIGNATURE_VERIFICATION_CALL_REVERT) {
            // errParams contains the revert reason
            revert SignatureVerificationCallRevert(errParams);
        } else {
            // TODO: add EMPTY_BLOB_QUORUMS, UNORDERED_OR_DUPLICATE_QUORUMS, GREATER_THAN_256_QUORUMS, INVALID_BLOB_VERSION
            revert("Unknown error code");
        }
    }

    /**
     * @notice hashes the given V2 batch header
     * @param batchHeader the V2 batch header to hash
     */
    function hashBatchHeaderV2(DATypesV2.BatchHeaderV2 memory batchHeader) internal pure returns (bytes32) {
        return keccak256(abi.encode(batchHeader));
    }

    /**
     * @notice hashes the given V2 blob header
     * @param blobHeader the V2 blob header to hash
     */
    function hashBlobHeaderV2(DATypesV2.BlobHeaderV2 memory blobHeader) internal pure returns (bytes32) {
        return keccak256(
            abi.encode(
                keccak256(abi.encode(blobHeader.version, blobHeader.quorumNumbers, blobHeader.commitment)),
                blobHeader.paymentHeaderHash
            )
        );
    }

    /**
     * @notice hashes the given V2 blob certificate
     * @param blobCertificate the V2 blob certificate to hash
     */
    function hashBlobCertificate(DATypesV2.BlobCertificate memory blobCertificate) internal pure returns (bytes32) {
        return keccak256(
            abi.encode(
                hashBlobHeaderV2(blobCertificate.blobHeader), blobCertificate.signature, blobCertificate.relayKeys
            )
        );
    }

    /**
     * @notice Converts an ordered array of quorum numbers into a bitmap.
     * @param orderedQuorums The array of quorum numbers to convert/compress into a bitmap. Must be in strictly ascending order.
     * @return err The status code indicating success or failure.
     * @return errParams Additional parameters related to the error, if any.
     * @return The resulting bitmap.
     * @dev Each byte in the input is processed as indicating a single bit to flip in the bitmap.
     * @dev This function returns an error status code if there >256 quorums, or the quorums are not ordered or contain duplicates.
     */
    function orderedQuorumsToBitmap(bytes memory orderedQuorums) internal pure returns (StatusCode err, bytes memory errParams, uint256) {
        // sanity-check on input. a too-long input would fail later on due to having duplicate entry(s)
        if (orderedQuorums.length > 256) {
            return (StatusCode.GREATER_THAN_256_QUORUMS, abi.encode(orderedQuorums.length), 0);
        }

        // Return empty bitmap early if length of array is 0.
        // This could be used for example if a caller doesn't want to enforce any required quorums.
        if (orderedQuorums.length == 0) {
            return (StatusCode.SUCCESS, "", uint256(0));
        }

        // initialize the empty bitmap, to be built inside the loop
        uint256 bitmap;
        // initialize an empty uint256 to be used as a bitmask inside the loop
        uint256 bitMask;

        // perform the 0-th loop iteration with the ordering check *omitted* (since it is unnecessary / will always pass)
        // construct a single-bit mask from the numerical value of the 0th byte of the array, and immediately add it to the bitmap
        bitmap = uint256(1 << uint8(orderedQuorums[0]));

        // loop through each byte in the array to construct the bitmap
        for (uint256 i = 1; i < orderedQuorums.length; ++i) {
            // construct a single-bit mask from the numerical value of the next byte of the array
            bitMask = uint256(1 << uint8(orderedQuorums[i]));
            // check strictly ascending array ordering by comparing the mask to the bitmap so far (revert if mask isn't greater than bitmap)
            if (bitMask <= bitmap) {
                return (StatusCode.UNORDERED_OR_DUPLICATE_QUORUMS, abi.encode(orderedQuorums), 0);
            }
            // add the entry to the bitmap
            bitmap = (bitmap | bitMask);
        }
        return (StatusCode.SUCCESS, "", bitmap);
    }
}
