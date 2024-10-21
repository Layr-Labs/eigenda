// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAServiceManager} from "./IEigenDAServiceManager.sol";
import {EigenDABlobVerificationUtils} from "../libraries/EigenDABlobVerificationUtils.sol";
import {IEigenDAThresholdRegistry} from "./IEigenDAThresholdRegistry.sol";

interface IEigenDABlobVerifier is IEigenDAThresholdRegistry {

    /**
     * @notice Verifies a the blob is valid for the required quorums
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob verification proof to verify the blob against
     */
    function verifyBlobV1(
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        EigenDABlobVerificationUtils.BlobVerificationProof calldata blobVerificationProof
    ) external view;

    /**
     * @notice Verifies that a blob is valid for the required quorums and additional quorums
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob verification proof to verify the blob against
     * @param additionalQuorumNumbersRequired The additional required quorum numbers 
     */
    function verifyBlobV1(
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        EigenDABlobVerificationUtils.BlobVerificationProof calldata blobVerificationProof,
        bytes calldata additionalQuorumNumbersRequired
    ) external view;

    function verifyBlobV2(
        EigenDABlobVerificationUtils.SignedCertificate calldata signedCertificate
    ) external view;

    function verifyBlobV2(
        EigenDABlobVerificationUtils.SignedCertificate calldata signedCertificate,
        bytes calldata additionalQuorumNumbersRequired
    ) external view;
}
