// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "src/interfaces/IEigenDAStructs.sol";

/**
 * @title IEigenDACertVerifierV2 - Interface for EigenDA V2 certificate verification
 * @author Layr Labs, Inc.
 * @notice Interface for V2-specific certificate verification methods
 */
interface IEigenDACertVerifierV2 {
    /**
     * @notice Verifies a blob cert for the specified quorums with the default security thresholds
     * @param batchHeader The batch header of the blob
     * @param blobInclusionInfo The inclusion proof for the blob cert
     * @param nonSignerStakesAndSignature The nonSignerStakesAndSignature to verify the blob cert against
     * @param signedQuorumNumbers The signed quorum numbers corresponding to the nonSignerStakesAndSignature
     */
    function verifyDACertV2(
        BatchHeaderV2 calldata batchHeader,
        BlobInclusionInfo calldata blobInclusionInfo,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature,
        bytes memory signedQuorumNumbers
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
     * @notice Thin try/catch wrapper around verifyDACertV2 that returns false instead of panicing
     * @dev The Steel library (https://github.com/risc0/risc0-ethereum/tree/main/crates/steel)
     *      currently has a limitation that it can only create zk proofs for functions that return a value
     * @param batchHeader The batch header of the blob
     * @param blobInclusionInfo The inclusion proof for the blob cert
     * @param nonSignerStakesAndSignature The nonSignerStakesAndSignature to verify the blob cert against
     * @param signedQuorumNumbers The signed quorum numbers corresponding to the nonSignerStakesAndSignature
     */
    function verifyDACertV2ForZKProof(
        BatchHeaderV2 calldata batchHeader,
        BlobInclusionInfo calldata blobInclusionInfo,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature,
        bytes memory signedQuorumNumbers
    ) external view returns (bool);

    /**
     * @notice Returns the nonSignerStakesAndSignature for a given blob cert and signed batch
     * @param signedBatch The signed batch to get the nonSignerStakesAndSignature for
     * @return nonSignerStakesAndSignature The nonSignerStakesAndSignature for the given signed batch attestation
     */
    function getNonSignerStakesAndSignature(SignedBatch calldata signedBatch)
        external
        view
        returns (NonSignerStakesAndSignature memory);

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
    function verifyDACertSecurityParams(uint16 version, SecurityThresholds memory securityThresholds) external view;
}