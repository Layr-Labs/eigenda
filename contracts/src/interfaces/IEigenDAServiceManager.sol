// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IServiceManager} from "eigenlayer-middleware/interfaces/IServiceManager.sol";
import {BLSSignatureChecker} from "eigenlayer-middleware/BLSSignatureChecker.sol";
import {BN254} from "eigenlayer-middleware/libraries/BN254.sol";
import {IEigenDAThresholdRegistry} from "./IEigenDAThresholdRegistry.sol";
import "./IEigenDAStructs.sol";

interface IEigenDAServiceManager is IServiceManager, IEigenDAThresholdRegistry {
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

    /// @notice mapping between the batchId to the hash of the metadata of the corresponding Batch
    function batchIdToBatchMetadataHash(uint32 batchId) external view returns(bytes32);

    /// @notice Returns the current batchId
    function taskNumber() external view returns (uint32);

    /// @notice Given a reference block number, returns the block until which operators must serve.
    function latestServeUntilBlock(uint32 referenceBlockNumber) external view returns (uint32);

    /// @notice The maximum amount of blocks in the past that the service will consider stake amounts to still be 'valid'.
    function BLOCK_STALE_MEASURE() external view returns (uint32);
}