// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifier} from "src/periphery/cert/interfaces/IEigenDACertVerifier.sol";
import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "src/core/interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";
import {EigenDACertVerificationV1Lib as CertV1Lib} from "src/periphery/cert/legacy/v1/EigenDACertVerificationV1Lib.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {EigenDATypesV2 as DATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {IEigenDARelayRegistry} from "src/core/interfaces/IEigenDARelayRegistry.sol";

/**
 * @title A CertVerifier is an immutable contract that is used by a consumer to verify EigenDA blob certificates
 *         to change these values or verification behavior a new CertVerifier must be deployed
 */
contract EigenDACertVerifierV1 {
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
     * @notice Verifies that the blob cert is valid for the required quorums
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob cert verification proof to verify
     */
    function verifyDACertV1(
        DATypesV1.BlobHeader calldata blobHeader,
        DATypesV1.BlobVerificationProof calldata blobVerificationProof
    ) external view {
        CertV1Lib._verifyDACertV1ForQuorums(
            _thresholdRegistry(), _batchMetadataStorage(), blobHeader, blobVerificationProof, quorumNumbersRequired()
        );
    }

    /**
     * @notice Verifies a batch of blob certs for the required quorums
     * @param blobHeaders The blob headers to verify
     * @param blobVerificationProofs The blob cert verification proofs to verify against
     */
    function verifyDACertsV1(
        DATypesV1.BlobHeader[] calldata blobHeaders,
        DATypesV1.BlobVerificationProof[] calldata blobVerificationProofs
    ) external view {
        CertV1Lib._verifyDACertsV1ForQuorums(
            _thresholdRegistry(), _batchMetadataStorage(), blobHeaders, blobVerificationProofs, quorumNumbersRequired()
        );
    }

    /// @notice Returns an array of bytes where each byte represents the adversary threshold percentage of the quorum at that index
    function quorumAdversaryThresholdPercentages() external view returns (bytes memory) {
        return _thresholdRegistry().quorumAdversaryThresholdPercentages();
    }

    /// @notice Returns an array of bytes where each byte represents the confirmation threshold percentage of the quorum at that index
    function quorumConfirmationThresholdPercentages() external view returns (bytes memory) {
        return _thresholdRegistry().quorumConfirmationThresholdPercentages();
    }

    /// @notice Returns an array of bytes where each byte represents the number of a required quorum
    function quorumNumbersRequired() public view returns (bytes memory) {
        return _thresholdRegistry().quorumNumbersRequired();
    }

    function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) external view returns (uint8) {
        return _thresholdRegistry().getQuorumAdversaryThresholdPercentage(quorumNumber);
    }

    /// @notice Gets the confirmation threshold percentage for a quorum
    function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) external view returns (uint8) {
        return _thresholdRegistry().getQuorumConfirmationThresholdPercentage(quorumNumber);
    }

    /// @notice Checks if a quorum is required
    function getIsQuorumRequired(uint8 quorumNumber) external view returns (bool) {
        return _thresholdRegistry().getIsQuorumRequired(quorumNumber);
    }

    /// @notice Returns the blob params for a given blob version
    function getBlobParams(uint16 version) public view returns (DATypesV1.VersionedBlobParams memory) {
        return _thresholdRegistry().getBlobParams(version);
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
}
