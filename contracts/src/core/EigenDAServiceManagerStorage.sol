// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {IEigenDAServiceManager} from "../interfaces/IEigenDAServiceManager.sol";

/**
 * @title Storage variables for the `EigenDAServiceManager` contract.
 * @author Layr Labs, Inc.
 * @notice This storage contract is separate from the logic to simplify the upgrade process.
 */
abstract contract EigenDAServiceManagerStorage is IEigenDAServiceManager {
    // CONSTANTS
    uint256 public constant THRESHOLD_DENOMINATOR = 100;

    //TODO: mechanism to change any of these values?
    /// @notice Unit of measure (in blocks) for which data will be stored for after confirmation.
    uint32 public constant STORE_DURATION_BLOCKS = 2 weeks / 12 seconds;

    /// @notice Minimum Batch size, in bytes.
    uint32 internal constant MIN_STORE_SIZE = 32;
    /// @notice Maximum Batch size, in bytes.
    uint32 internal constant MAX_STORE_SIZE = 4e9;
    /**
     * @notice The maximum amount of blocks in the past that the service will consider stake amounts to still be 'valid'.
     * @dev To clarify edge cases, the middleware can look `BLOCK_STALE_MEASURE` blocks into the past, i.e. it may trust stakes from the interval
     * [block.number - BLOCK_STALE_MEASURE, block.number] (specifically, *inclusive* of the block that is `BLOCK_STALE_MEASURE` before the current one)
     * @dev BLOCK_STALE_MEASURE should be greater than the number of blocks till finalization, but not too much greater, as it is the amount of
     * time that nodes can be active after they have deregistered. The larger it is, the farther back stakes can be used, but the longer operators
     * have to serve after they've deregistered.
     */
    uint32 public constant BLOCK_STALE_MEASURE = 150;
    
    /// @notice The current batchId
    uint32 public batchId;

    /// @notice mapping between the batchId to the hash of the metadata of the corresponding Batch
    mapping(uint32 => bytes32) public batchIdToBatchMetadataHash;

    /// @notice address that is permissioned to confirm batches
    address public batchConfirmer;

    /// @notice metadata URI for the EigenDA AVS
    string public metadataURI;
}
