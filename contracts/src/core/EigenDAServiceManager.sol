// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {Pausable} from "eigenlayer-core/contracts/permissions/Pausable.sol";
import {IPauserRegistry} from "eigenlayer-core/contracts/interfaces/IPauserRegistry.sol";

import {ServiceManagerBase, IAVSDirectory, IPaymentCoordinator} from "eigenlayer-middleware/ServiceManagerBase.sol";
import {BLSSignatureChecker} from "eigenlayer-middleware/BLSSignatureChecker.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";

import {EigenDAServiceManagerStorage} from "./EigenDAServiceManagerStorage.sol";
import {EigenDAHasher} from "../libraries/EigenDAHasher.sol";

/**
 * @title Primary entrypoint for procuring services from EigenDA.
 * @author Layr Labs, Inc.
 * @notice This contract is used for:
 * - initializing the data store by the disperser
 * - confirming the data store by the disperser with inferred aggregated signatures of the quorum
 * - freezing operators as the result of various "challenges"
 */
contract EigenDAServiceManager is EigenDAServiceManagerStorage, ServiceManagerBase, BLSSignatureChecker, Pausable {
    using EigenDAHasher for BatchHeader;
    using EigenDAHasher for ReducedBatchHeader;

    uint8 internal constant PAUSED_CONFIRM_BATCH = 0;

    /// @notice when applied to a function, ensures that the function is only callable by the `batchConfirmer`.
    modifier onlyBatchConfirmer() {
        require(isBatchConfirmer[msg.sender], "onlyBatchConfirmer: not from batch confirmer");
        _;
    }

    constructor(
        IAVSDirectory __avsDirectory,
        IPaymentCoordinator __paymentCoordinator,
        IRegistryCoordinator __registryCoordinator,
        IStakeRegistry __stakeRegistry
    )
        BLSSignatureChecker(__registryCoordinator)
        ServiceManagerBase(__avsDirectory, __paymentCoordinator, __registryCoordinator, __stakeRegistry)
    {
        _disableInitializers();
    }

    function initialize(
        IPauserRegistry _pauserRegistry,
        uint256 _initialPausedStatus,
        address _initialOwner,
        address[] memory _batchConfirmers
    )
        public
        initializer
    {
        _initializePauser(_pauserRegistry, _initialPausedStatus);
        _transferOwnership(_initialOwner);
        for (uint i = 0; i < _batchConfirmers.length; ++i) {
            _setBatchConfirmer(_batchConfirmers[i]);
        }
    }

    /**
     * @notice This function is used for
     * - submitting data availabilty certificates,
     * - check that the aggregate signature is valid,
     * - and check whether quorum has been achieved or not.
     */
    function confirmBatch(
        BatchHeader calldata batchHeader,
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature
    ) external onlyWhenNotPaused(PAUSED_CONFIRM_BATCH) onlyBatchConfirmer() {
        // make sure the information needed to derive the non-signers and batch is in calldata to avoid emitting events
        require(tx.origin == msg.sender, "EigenDAServiceManager.confirmBatch: header and nonsigner data must be in calldata");
        // make sure the stakes against which the Batch is being confirmed are not stale
        require(
            batchHeader.referenceBlockNumber < block.number, "EigenDAServiceManager.confirmBatch: specified referenceBlockNumber is in future"
        );

        require(
            (batchHeader.referenceBlockNumber + BLOCK_STALE_MEASURE) >= uint32(block.number),
            "EigenDAServiceManager.confirmBatch: specified referenceBlockNumber is too far in past"
        );

        //make sure that the quorumNumbers and signedStakeForQuorums are of the same length
        require(
            batchHeader.quorumNumbers.length == batchHeader.signedStakeForQuorums.length,
            "EigenDAServiceManager.confirmBatch: quorumNumbers and signedStakeForQuorums must be of the same length"
        );

        // calculate reducedBatchHeaderHash which nodes signed
        bytes32 reducedBatchHeaderHash = batchHeader.hashBatchHeaderToReducedBatchHeader();

        // check the signature
        (
            QuorumStakeTotals memory quorumStakeTotals,
            bytes32 signatoryRecordHash
        ) = checkSignatures(
            reducedBatchHeaderHash, 
            batchHeader.quorumNumbers, // use list of uint8s instead of uint256 bitmap to not iterate 256 times
            batchHeader.referenceBlockNumber, 
            nonSignerStakesAndSignature
        );

        // check that signatories own at least a threshold percentage of each quourm
        for (uint i = 0; i < batchHeader.signedStakeForQuorums.length; i++) {
            // we don't check that the signedStakeForQuorums are not >100 because a greater value would trivially fail the check, implying 
            // signed stake > total stake
            require(
                quorumStakeTotals.signedStakeForQuorum[i] * THRESHOLD_DENOMINATOR >= 
                    quorumStakeTotals.totalStakeForQuorum[i] * uint8(batchHeader.signedStakeForQuorums[i]),
                "EigenDAServiceManager.confirmBatch: signatories do not own at least threshold percentage of a quorum"
            );
        }

        // store the metadata hash
        uint32 batchIdMemory = batchId;
        bytes32 batchHeaderHash = batchHeader.hashBatchHeader();
        batchIdToBatchMetadataHash[batchIdMemory] = EigenDAHasher.hashBatchHashedMetadata(batchHeaderHash, signatoryRecordHash, uint32(block.number));

        emit BatchConfirmed(reducedBatchHeaderHash, batchIdMemory);

        // increment the batchId
        batchId = batchIdMemory + 1;
    }

    /// @notice This function is used for changing the batch confirmer
    function setBatchConfirmer(address _batchConfirmer) external onlyOwner() {
        _setBatchConfirmer(_batchConfirmer);
    }

    /// @notice changes the batch confirmer
    function _setBatchConfirmer(address _batchConfirmer) internal {
        isBatchConfirmer[_batchConfirmer] = !isBatchConfirmer[_batchConfirmer];
        emit BatchConfirmerStatusChanged(_batchConfirmer, isBatchConfirmer[_batchConfirmer]);
    }

    /// @notice Returns the current batchId
    function taskNumber() external view returns (uint32) {
        return batchId;
    }

    /// @notice Given a reference block number, returns the block until which operators must serve.
    function latestServeUntilBlock(uint32 referenceBlockNumber) external view returns (uint32) {
        return referenceBlockNumber + STORE_DURATION_BLOCKS + BLOCK_STALE_MEASURE;
    }

}