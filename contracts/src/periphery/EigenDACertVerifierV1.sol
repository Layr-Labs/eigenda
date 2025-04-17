// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifierV1} from "src/periphery/interfaces/IEigenDACertVerifierV1.sol";
import {IEigenDAThresholdRegistry} from "src/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "src/interfaces/IEigenDABatchMetadataStorage.sol";
import {EigenDACertVerificationV1Lib as CertLib} from "src/periphery/libraries/EigenDACertVerificationV1Lib.sol";
import "src/interfaces/IEigenDAStructs.sol";

/**
 * @title CertVerifierV1 - EigenDA V1 certificate verification
 * @author Layr Labs, Inc.
 * @notice Contains all V1-specific verification functionality
 */
contract EigenDACertVerifierV1 is IEigenDACertVerifierV1 {

    IEigenDAThresholdRegistry public immutable eigenDAThresholdRegistryV1;

    IEigenDABatchMetadataStorage public immutable eigenDABatchMetadataStorageV1;

    /**
     * @notice Constructor for the EigenDA V1 certificate verifier
     * @param _eigenDAThresholdRegistryV1 The address of the EigenDAThresholdRegistry contract
     * @param _eigenDABatchMetadataStorageV1 The address of the EigenDABatchMetadataStorage contract
     */
    constructor(
        IEigenDAThresholdRegistry _eigenDAThresholdRegistryV1,
        IEigenDABatchMetadataStorage _eigenDABatchMetadataStorageV1
    ) {
        eigenDAThresholdRegistryV1 = _eigenDAThresholdRegistryV1;
        eigenDABatchMetadataStorageV1 = _eigenDABatchMetadataStorageV1;
    }

    /**
     * @notice Verifies a the blob cert is valid for the required quorums. Reverts on verification failure.
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
     * @notice Checks a blob cert and returns result without reverting
     * @param blobHeader Pointer to the blob header in calldata
     * @param blobVerificationProof Pointer to the blob cert verification proof in calldata
     * @return success True if verification succeeded, false otherwise
     */
    function checkDACertV1(
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

    /**
     * @notice Returns the threshold registry contract
     * @return The IEigenDAThresholdRegistry contract
     * @dev Can be overridden by derived contracts
     */
    function _thresholdRegistry() internal view virtual returns (IEigenDAThresholdRegistry) {
        return eigenDAThresholdRegistryV1;
    }

    /**
     * @notice Returns the batch metadata storage contract
     * @return The IEigenDABatchMetadataStorage contract
     * @dev Can be overridden by derived contracts
     */
    function _batchMetadataStorage() internal view virtual returns (IEigenDABatchMetadataStorage) {
        return eigenDABatchMetadataStorageV1;
    }

    /**
     * @notice Returns the quorum numbers required for verification
     * @return bytes The required quorum numbers
     * @dev Can be overridden by derived contracts
     */
    function _quorumNumbersRequired() internal view virtual returns (bytes memory) {
        return _thresholdRegistry().quorumNumbersRequired();
    }

    /**
     * @notice Returns the quorum confirmation threshold percentages
     * @return bytes The confirmation threshold percentages for each quorum
     * @dev Can be overridden by derived contracts
     */
    function _quorumConfirmationThresholdPercentages() internal view virtual returns (bytes memory) {
        return _thresholdRegistry().quorumConfirmationThresholdPercentages();
    }

    /**
     * @notice Returns the stored batch metadata hash for a given blob verification proof
     * @param blobVerificationProof The blob verification proof
     * @return bytes32 The stored batch metadata hash
     * @dev Can be overridden by derived contracts
     */
    function _storedBatchMetadataHash(BlobVerificationProof calldata blobVerificationProof)
        internal
        view
        virtual
        returns (bytes32)
    {
        return _batchMetadataStorage().batchIdToBatchMetadataHash(blobVerificationProof.batchId);
    }
}
