// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import {IRegistry} from "./IRegistry.sol";

/**
 * @title Interface for a `Registry`-type contract that keeps track of an ordered list of operators for up to 256 quorums.
 * @author Layr Labs, Inc.
 */
interface IIndexRegistry is IRegistry {
    // EVENTS
    
    // emitted when an operator's index in the ordered operator list for the quorum with number `quorumNumber` is updated
    event QuorumIndexUpdate(bytes32 indexed operatorId, uint8 quorumNumber, uint32 newOperatorIndex);

    // DATA STRUCTURES

    // struct used to give definitive ordering to operators at each blockNumber. 
    struct OperatorUpdate {
        // blockNumber number from which `operatorIndex` was the operators index
        // the operator's index is the first entry such that `blockNumber >= entry.fromBlockNumber`
        uint32 fromBlockNumber;
        // the operator at this index
        bytes32 operatorId;
    }

    // struct used to denote the number of operators in a quorum at a given blockNumber
    struct QuorumUpdate {
        // The total number of operators at a `blockNumber` is the first entry such that `blockNumber >= entry.fromBlockNumber`
        uint32 fromBlockNumber;
        // The number of operators at `fromBlockNumber`
        uint32 numOperators;
    }

    /**
     * @notice Registers the operator with the specified `operatorId` for the quorums specified by `quorumNumbers`.
     * @param operatorId is the id of the operator that is being registered
     * @param quorumNumbers is the quorum numbers the operator is registered for
     * @return numOperatorsPerQuorum is a list of the number of operators (including the registering operator) in each of the quorums the operator is registered for
     * @dev access restricted to the RegistryCoordinator
     * @dev Preconditions (these are assumed, not validated in this contract):
     *         1) `quorumNumbers` has no duplicates
     *         2) `quorumNumbers.length` != 0
     *         3) `quorumNumbers` is ordered in ascending order
     *         4) the operator is not already registered
     */
    function registerOperator(bytes32 operatorId, bytes calldata quorumNumbers) external returns(uint32[] memory);

    /**
     * @notice Deregisters the operator with the specified `operatorId` for the quorums specified by `quorumNumbers`.
     * @param operatorId is the id of the operator that is being deregistered
     * @param quorumNumbers is the quorum numbers the operator is deregistered for
     * @dev access restricted to the RegistryCoordinator
     * @dev Preconditions (these are assumed, not validated in this contract):
     *         1) `quorumNumbers` has no duplicates
     *         2) `quorumNumbers.length` != 0
     *         3) `quorumNumbers` is ordered in ascending order
     *         4) the operator is not already deregistered
     *         5) `quorumNumbers` is a subset of the quorumNumbers that the operator is registered for
     */
    function deregisterOperator(bytes32 operatorId, bytes calldata quorumNumbers) external;

    /**
     * @notice Initialize a quorum by pushing its first quorum update
     * @param quorumNumber The number of the new quorum
     */
    function initializeQuorum(uint8 quorumNumber) external;

    /// @notice Returns the OperatorUpdate entry for the specified `operatorIndex` and `quorumNumber` at the specified `arrayIndex`
    function getOperatorUpdateAtIndex(
        uint8 quorumNumber,
        uint32 operatorIndex,
        uint32 arrayIndex
    ) external view returns (OperatorUpdate memory);

    /// @notice Returns the QuorumUpdate entry for the specified `quorumNumber` at the specified `quorumIndex`
    function getQuorumUpdateAtIndex(uint8 quorumNumber, uint32 quorumIndex) external view returns (QuorumUpdate memory);

    /// @notice Returns the most recent OperatorUpdate entry for the specified quorumNumber and operatorIndex
    function getLatestOperatorUpdate(uint8 quorumNumber, uint32 operatorIndex) external view returns (OperatorUpdate memory);

    /// @notice Returns the most recent QuorumUpdate entry for the specified quorumNumber
    function getLatestQuorumUpdate(uint8 quorumNumber) external view returns (QuorumUpdate memory);

    /// @notice Returns the current number of operators of this service for `quorumNumber`.
    function totalOperatorsForQuorum(uint8 quorumNumber) external view returns (uint32);

    /// @notice Returns an ordered list of operators of the services for the given `quorumNumber` at the given `blockNumber`
    function getOperatorListAtBlockNumber(uint8 quorumNumber, uint32 blockNumber) external view returns (bytes32[] memory);
}
