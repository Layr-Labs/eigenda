// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@eigenlayer-middleware/interfaces/IServiceManager.sol";
import "@eigenlayer-middleware/interfaces/IDelayedService.sol";
import "@eigenlayer-middleware/BLSSignatureChecker.sol";
import "@eigenlayer-core/contracts/interfaces/IDelegationManager.sol";

interface IEigenDAServiceManager is IServiceManager, IDelayedService {
    // EVENTS
    
    /**b
     * @notice Emitted when a Batch is confirmed.
     * @param batchHeaderHash The hash of the batch header
     * @param batchId The ID for the Batch inside of the specified duration (i.e. *not* the globalBatchId)
     */
    event BatchConfirmed(bytes32 indexed batchHeaderHash, uint32 batchId, uint96 fee);

    event FeePerBytePerTimeSet(uint256 previousValue, uint256 newValue);

    event PaymentManagerSet(address previousAddress, address newAddress);

    event FeeSetterChanged(address previousAddress, address newAddress);

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
}
