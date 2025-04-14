// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAThresholdRegistry} from "src/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "src/interfaces/IEigenDABatchMetadataStorage.sol";
import {EigenDACertVerificationV1Lib as CertLib} from "src/periphery/libraries/EigenDACertVerificationV1Lib.sol";
import "src/interfaces/IEigenDAStructs.sol";

/**
 * @title CertVerifierV1 - EigenDA V1 certificate verification
 * @author Layr Labs, Inc.
 * @notice Contains all V1-specific verification functionality
 */
contract EigenDACertVerifierV1 {
    /// @notice Thrown when there is a length mismatch
    error LengthMismatch();

    /// @notice The EigenDAThresholdRegistry contract address
    IEigenDAThresholdRegistry public immutable eigenDAThresholdRegistryV1;

    /// @notice The EigenDABatchMetadataStorage contract address
    /// @dev On L1 this contract is the EigenDA Service Manager contract
    IEigenDABatchMetadataStorage public immutable eigenDABatchMetadataStorageV1;

    constructor(
        IEigenDAThresholdRegistry _eigenDAThresholdRegistryV1,
        IEigenDABatchMetadataStorage _eigenDABatchMetadataStorageV1
    ) {
        eigenDAThresholdRegistryV1 = _eigenDAThresholdRegistryV1;
        eigenDABatchMetadataStorageV1 = _eigenDABatchMetadataStorageV1;
    }

    /**
     * @notice Verifies a the blob cert is valid for the required quorums
     * @param blobHeader Pointer to the blob header in calldata
     * @param blobVerificationProof Pointer to the blob cert verification proof in calldata
     */
    function verifyDACertV1(BlobHeader calldata blobHeader, BlobVerificationProof calldata blobVerificationProof)
        public
        view
    {
        (CertLib.ErrorCode err, bytes memory errParams) = CertLib.verifyDACert(
            _quorumConfirmationThresholdPercentages(),
            _storedBatchMetadataHash(blobVerificationProof),
            blobHeader,
            blobVerificationProof,
            _quorumNumbersRequired()
        );
        CertLib.revertOnError(err, errParams);
    }

    /**
     * @notice Verifies a batch of blob certs for the required quorums
     * @param blobHeaders Pointer to array of blob headers in calldata
     * @param blobVerificationProofs Pointer to array of blob cert verification proofs in calldata
     */
    function verifyDACertsV1(BlobHeader[] calldata blobHeaders, BlobVerificationProof[] calldata blobVerificationProofs)
        external
        view
    {
        // Check length match
        if (blobHeaders.length != blobVerificationProofs.length) {
            revert LengthMismatch();
        }

        // Verify each blob
        for (uint256 i; i < blobHeaders.length; ++i) {
            verifyDACertV1(blobHeaders[i], blobVerificationProofs[i]);
        }
    }

    /**
     * @notice Verifies a blob cert and returns result without reverting
     * @param blobHeader Pointer to the blob header in calldata
     * @param blobVerificationProof Pointer to the blob cert verification proof in calldata
     * @return success True if verification succeeded
     */
    function verifyDACertV1ForZkProof(
        BlobHeader calldata blobHeader,
        BlobVerificationProof calldata blobVerificationProof
    ) external view returns (bool success) {
        (CertLib.ErrorCode errorCode,) = CertLib.verifyDACert(
            _quorumConfirmationThresholdPercentages(),
            _storedBatchMetadataHash(blobVerificationProof),
            blobHeader,
            blobVerificationProof,
            _quorumNumbersRequired()
        );
        if (errorCode == CertLib.ErrorCode.SUCCESS) {
            return true;
        } else {
            return false;
        }
    }

    function _thresholdRegistry() internal view virtual returns (IEigenDAThresholdRegistry) {
        return eigenDAThresholdRegistryV1;
    }

    function _batchMetadataStorage() internal view virtual returns (IEigenDABatchMetadataStorage) {
        return eigenDABatchMetadataStorageV1;
    }

    function _quorumNumbersRequired() internal view virtual returns (bytes memory) {
        return _thresholdRegistry().quorumNumbersRequired();
    }

    function _quorumConfirmationThresholdPercentages() internal view virtual returns (bytes memory) {
        return _thresholdRegistry().quorumConfirmationThresholdPercentages();
    }

    function _storedBatchMetadataHash(BlobVerificationProof calldata blobVerificationProof)
        internal
        view
        virtual
        returns (bytes32)
    {
        return _batchMetadataStorage().batchIdToBatchMetadataHash(blobVerificationProof.batchId);
    }
}
