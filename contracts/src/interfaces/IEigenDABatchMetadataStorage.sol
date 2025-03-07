// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

/**
 * @title IEigenDABatchMetadataStorage
 * @notice This contract is used for storing the batch metadata for a bridged V1 batch
 * @dev This contract is deployed on L1 as the EigenDAServiceManager contract
 */
interface IEigenDABatchMetadataStorage {

    /**
     * @notice Returns the batch metadata hash for a given batch id
     * @param batchId The id of the batch to get the metadata hash for
     */
    function batchIdToBatchMetadataHash(uint32 batchId) external view returns (bytes32);
}