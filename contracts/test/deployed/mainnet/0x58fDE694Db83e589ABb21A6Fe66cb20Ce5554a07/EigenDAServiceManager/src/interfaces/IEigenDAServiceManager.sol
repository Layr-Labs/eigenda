// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {IServiceManager} from "eigenlayer-middleware/interfaces/IServiceManager.sol";
import {BLSSignatureChecker} from "eigenlayer-middleware/BLSSignatureChecker.sol";
import {BN254} from "eigenlayer-middleware/libraries/BN254.sol";

interface IEigenDAServiceManager is IServiceManager {
    // EVENTS
    
    /**
     * @notice Emitted when a Batch is confirmed.
     * @param batchHeaderHash The hash of the batch header
     * @param batchId The ID for the Batch inside of the specified duration (i.e. *not* the globalBatchId)
     */
    event BatchConfirmed(bytes32 indexed batchHeaderHash, uint32 batchId);

    /**
     * @notice Emitted when a batch confirmer status is updated.
     * @param batchConfirmer The address of the batch confirmer
     * @param status The new status of the batch confirmer
     */
    event BatchConfirmerStatusChanged(address batchConfirmer, bool status);

    // STRUCTS

    struct QuorumBlobParam {
        uint8 quorumNumber;
        uint8 adversaryThresholdPercentage;
        uint8 confirmationThresholdPercentage; 
        uint32 chunkLength; // the length of the chunks in the quorum
    }

    struct BlobHeader {
        BN254.G1Point commitment; // the kzg commitment to the blob
        uint32 dataLength; // the length of the blob in coefficients of the polynomial
        QuorumBlobParam[] quorumBlobParams; // the quorumBlobParams for each quorum
    }

    struct ReducedBatchHeader {
        bytes32 blobHeadersRoot;
        uint32 referenceBlockNumber;
    }

    struct BatchHeader {
        bytes32 blobHeadersRoot;
        bytes quorumNumbers; // each byte is a different quorum number
        bytes signedStakeForQuorums; // every bytes is an amount less than 100 specifying the percentage of stake 
                                     // that has signed in the corresponding quorum in `quorumNumbers` 
        uint32 referenceBlockNumber;
    }
    
    // Relevant metadata for a given datastore
    struct BatchMetadata {
        BatchHeader batchHeader; // the header of the data store
        bytes32 signatoryRecordHash; // the hash of the signatory record
        uint32 confirmationBlockNumber; // the block number at which the batch was confirmed
    }

    // FUNCTIONS

    /// @notice mapping between the batchId to the hash of the metadata of the corresponding Batch
    function batchIdToBatchMetadataHash(uint32 batchId) external view returns(bytes32);

    /**
     * @notice This function is used for
     * - submitting data availabilty certificates,
     * - check that the aggregate signature is valid,
     * - and check whether quorum has been achieved or not.
     */
    function confirmBatch(
        BatchHeader calldata batchHeader,
        BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature
    ) external;

    /// @notice This function is used for changing the batch confirmer
    function setBatchConfirmer(address _batchConfirmer) external;

    /// @notice Returns the current batchId
    function taskNumber() external view returns (uint32);

    /// @notice Given a reference block number, returns the block until which operators must serve.
    function latestServeUntilBlock(uint32 referenceBlockNumber) external view returns (uint32);

    /// @notice The maximum amount of blocks in the past that the service will consider stake amounts to still be 'valid'.
    function BLOCK_STALE_MEASURE() external view returns (uint32);

    /// @notice Returns the bytes array of quorumAdversaryThresholdPercentages
    function quorumAdversaryThresholdPercentages() external view returns (bytes memory);

    /// @notice Returns the bytes array of quorumAdversaryThresholdPercentages
    function quorumConfirmationThresholdPercentages() external view returns (bytes memory);

    /// @notice Returns the bytes array of quorumsNumbersRequired
    function quorumNumbersRequired() external view returns (bytes memory);
}
