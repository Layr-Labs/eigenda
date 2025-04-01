// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "./IEigenDAStructs.sol";

/// @dev It is impossible to implement all of IEigenDACertVerifier in the CertVerifierRouter, and so a subset of the interface is defined here.
interface IEigenDACertVerifierBase {
    /**
     * @notice Verifies a the blob cert is valid for the required quorums
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob cert verification proof to verify against
     */
    function verifyDACertV1(BlobHeader calldata blobHeader, BlobVerificationProof calldata blobVerificationProof)
        external
        view;

    /**
     * @notice Verifies a batch of blob certs for the required quorums
     * @param blobHeaders The blob headers to verify
     * @param blobVerificationProofs The blob cert verification proofs to verify against
     */
    function verifyDACertsV1(BlobHeader[] calldata blobHeaders, BlobVerificationProof[] calldata blobVerificationProofs)
        external
        view;

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
}

interface IEigenDACertVerifier is IEigenDACertVerifierBase {
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
<<<<<<< HEAD
<<<<<<< HEAD
    function verifyDACertSecurityParams(uint16 version, SecurityThresholds memory securityThresholds) external view;
=======
    function verifyDACertSecurityParams(
        uint16 version,
        SecurityThresholds memory securityThresholds
    ) external view;
>>>>>>> de14d145 (feat: cert verifier router)
=======
    function verifyDACertSecurityParams(uint16 version, SecurityThresholds memory securityThresholds) external view;
>>>>>>> d72f87a1 (chore: forge fmt)
}
