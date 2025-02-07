// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAThresholdRegistry} from "./IEigenDAThresholdRegistry.sol";
import "./IEigenDAStructs.sol";

interface IEigenDACertVerifier is IEigenDAThresholdRegistry {

    /**
     * @notice Verifies a the blob cert is valid for the required quorums
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob cert verification proof to verify against
     */
    function verifyDACertV1(
        BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof
    ) external view;


    /**
     * @notice Verifies a batch of blob certs for the required quorums
     * @param blobHeaders The blob headers to verify
     * @param blobVerificationProofs The blob cert verification proofs to verify against
     */
    function verifyDACertsV1(
        BlobHeader[] calldata blobHeaders,
        BlobVerificationProof[] calldata blobVerificationProofs
    ) external view;

    /**
     * @notice Verifies a blob cert for the specified quorums with the default security thresholds
     * @param batchHeader The batch header of the blob 
     * @param blobInclusionInfo The inclusion proof for the blob cert
     * @param nonSignerStakesAndSignature The nonSignerStakesAndSignature to verify the blob cert against
     */
    function verifyDACertV2(
        BatchHeaderV2 calldata batchHeader,
        BlobInclusionInfo calldata blobInclusionInfo,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature
    ) external view;

    /**
     * @notice Verifies a blob cert for the specified quorums with the default security thresholds
     * @param signedBatch The signed batch to verify the blob cert against
     * @param blobInclusionInfo The inclusion proof for the blob cert
     */
    function verifyDACertV2FromSignedBatch(
        SignedBatch calldata signedBatch,
        BlobInclusionInfo calldata blobInclusionInfo
    ) external view;

    /**
     * @notice Returns the nonSignerStakesAndSignature for a given blob cert and signed batch
     * @param signedBatch The signed batch to get the nonSignerStakesAndSignature for
     */
    function getNonSignerStakesAndSignature(
        SignedBatch calldata signedBatch
    ) external view returns (NonSignerStakesAndSignature memory);

    /**
     * @notice Verifies the security parameters for a blob cert
     * @param blobParams The blob params to verify 
     * @param securityThresholds The security thresholds to verify against
     */
    function verifyDACertSecurityParams(
        VersionedBlobParams memory blobParams,
        SecurityThresholds memory securityThresholds
    ) external view;

    /**
     * @notice Verifies the security parameters for a blob cert
     * @param version The version of the blob to verify
     * @param securityThresholds The security thresholds to verify against
     */
    function verifyDACertSecurityParams(
        uint16 version,
        SecurityThresholds memory securityThresholds
    ) external view;
}