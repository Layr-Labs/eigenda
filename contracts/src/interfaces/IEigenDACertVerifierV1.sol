// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAThresholdRegistry} from "./IEigenDAThresholdRegistry.sol";
import "./IEigenDAStructs.sol";

interface IEigenDACertVerifierV1 is IEigenDAThresholdRegistry {
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
}
