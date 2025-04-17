// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "src/interfaces/IEigenDAStructs.sol";

/**
 * @title IEigenDACertVerifierV1 - Interface for EigenDA V1 certificate verification
 * @author Layr Labs, Inc.
 * @notice Interface for V1-specific certificate verification methods
 */
interface IEigenDACertVerifierV1 {
    /**
     * @notice Verifies a the blob cert is valid for the required quorums
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob cert verification proof to verify against
     */
    function verifyDACertV1(BlobHeader calldata blobHeader, BlobVerificationProof calldata blobVerificationProof)
        external
        view;

    /**
     * @notice Checks a blob cert and returns result without reverting
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob cert verification proof to verify against
     * @return success True if verification succeeded, false otherwise
     */
    function checkDACertV1(
        BlobHeader calldata blobHeader, 
        BlobVerificationProof calldata blobVerificationProof
    ) external view returns (bool success);
}
