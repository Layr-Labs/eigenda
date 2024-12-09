// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAThresholdRegistry} from "./IEigenDAThresholdRegistry.sol";
import "./IEigenDAStructs.sol";

interface IEigenDABlobVerifier is IEigenDAThresholdRegistry {

    /**
     * @notice Verifies a the blob is valid for the required quorums
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob verification proof to verify the blob against
     */
    function verifyBlobV1(
        BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof
    ) external view;


    /**
     * @notice Verifies a batch of blobs for the required quorums
     * @param blobHeaders The blob headers to verify
     * @param blobVerificationProofs The blob verification proofs to verify the blobs against
     */
    function verifyBlobsV1(
        BlobHeader[] calldata blobHeaders,
        BlobVerificationProof[] calldata blobVerificationProofs
    ) external view;

    /**
     * @notice Verifies a blob for the required quorums and the default security thresholds
     * @param batchHeader The batch header of the blob
     * @param blobVerificationProof The blob verification proof for the blob
     * @param nonSignerStakesAndSignature The nonSignerStakesAndSignature to verify the blob against
     */
    function verifyBlobV2(
        BatchHeaderV2 calldata batchHeader,
        BlobVerificationProofV2 calldata blobVerificationProof,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature
    ) external view;

    /**
     * @notice Verifies a blob for the base required quorums and the default security thresholds
     * @param signedBatch The signed batch to verify the blob against
     * @param blobVerificationProof The blob verification proof for the blob
     */
    function verifyBlobV2FromSignedBatch(
        SignedBatch calldata signedBatch,
        BlobVerificationProofV2 calldata blobVerificationProof
    ) external view;

    /**
     * @notice Returns the nonSignerStakesAndSignature for a given blob and signed batch
     * @param signedBatch The signed batch to get the nonSignerStakesAndSignature for
     */
    function getNonSignerStakesAndSignature(
        SignedBatch calldata signedBatch
    ) external view returns (NonSignerStakesAndSignature memory);

    /**
     * @notice Verifies the security parameters for a blob
     * @param blobParams The blob params to verify 
     * @param securityThresholds The security thresholds to verify against
     */
    function verifyBlobSecurityParams(
        VersionedBlobParams memory blobParams,
        SecurityThresholds memory securityThresholds
    ) external view;

    /**
     * @notice Verifies the security parameters for a blob
     * @param version The version of the blob to verify
     * @param securityThresholds The security thresholds to verify against
     */
    function verifyBlobSecurityParams(
        uint16 version,
        SecurityThresholds memory securityThresholds
    ) external view;
}