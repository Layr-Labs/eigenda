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

    /**
     * @notice The maximum amount of blocks in the past that the service will consider stake amounts to still be 'valid'.
     * @dev To clarify edge cases, the middleware can look `BLOCK_STALE_MEASURE` blocks into the past, i.e. it may trust stakes from the interval
     * [block.number - BLOCK_STALE_MEASURE, block.number] (specifically, *inclusive* of the block that is `BLOCK_STALE_MEASURE` before the current one)
     * @dev BLOCK_STALE_MEASURE should be greater than the number of blocks till finalization, but not too much greater, as it is the amount of
     * time that nodes can be active after they have deregistered. The larger it is, the farther back stakes can be used, but the longer operators
     * have to serve after they've deregistered.
     */
    uint32 public constant BLOCK_STALE_MEASURE = 150;

    /** 
     * @notice The quorum adversary threshold percentages stored as an ordered bytes array 
     * this is the percentage of the total stake that must be adversarial to consider a blob invalid
     */
    bytes public constant quorumAdversaryThresholdPercentages = hex"2121";

    /** 
     * @notice The quorum confirmation threshold percentages stored as an ordered bytes array 
     * this is the percentage of the total stake needed to confirm a blob
     */
    bytes public constant quorumConfirmationThresholdPercentages = hex"4242";

    /** 
     * @notice The quorum numbers required for confirmation stored as an ordered bytes array 
     * these quorum numbers have respective canonical thresholds in the
     * quorumConfirmationThresholdPercentages and quorumAdversaryThresholdPercentages above
     */
    bytes public constant quorumsNumbersRequired = hex"0001";
    
    /// @notice The current batchId
    uint32 public batchId;

    /// @notice mapping between the batchId to the hash of the metadata of the corresponding Batch
    mapping(uint32 => bytes32) public batchIdToBatchMetadataHash;

    /// @notice address that is permissioned to confirm batches
    address public batchConfirmer;

     // storage gap for upgradeability
     // slither-disable-next-line shadowing-state
     uint256[47] private __GAP;

}
