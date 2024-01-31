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
     * @notice Emitted when the batch confirmer is changed.
     * @param previousAddress The address of the previous batch confirmer
     * @param newAddress The address of the new batch confirmer
     */
    event BatchConfirmerChanged(address previousAddress, address newAddress);

    // STRUCTS

    struct QuorumBlobParam {
        uint8 quorumNumber;
        uint8 adversaryThresholdPercentage;
        uint8 quorumThresholdPercentage; 
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
        bytes quorumThresholdPercentages; // every bytes is an amount less than 100 specifying the percentage of stake 
                                          // the must have signed in the corresponding quorum in `quorumNumbers` 
        uint32 referenceBlockNumber;
    }
    
    // Relevant metadata for a given datastore
    struct BatchMetadata {
        BatchHeader batchHeader; // the header of the data store
        bytes32 signatoryRecordHash; // the hash of the signatory record
        uint96 fee; // the amount of paymentToken paid for the datastore
        uint32 confirmationBlockNumber; // the block number at which the batch was confirmed
    }

    // Relevant metadata for a given datastore
    struct BatchMetadataWithSignatoryRecord {
        bytes32 batchHeaderHash; // the header hash of the data store
        uint32 referenceBlockNumber; // the block number at which stakes 
        bytes32[] nonSignerPubkeyHashes; // the pubkeyHashes of all of the nonSigners
        uint96 fee; // the amount of paymentToken paid for the datastore
        uint32 blockNumber; // the block number at which the datastore was confirmed
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

    /// @notice Returns the block until which operators must serve.
    function latestServeUntilBlock() external view returns (uint32);

    /// @notice The maximum amount of blocks in the past that the service will consider stake amounts to still be 'valid'.
    function BLOCK_STALE_MEASURE() external view returns (uint32);
}
