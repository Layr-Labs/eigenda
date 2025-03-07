// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {Pausable} from "../../lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/permissions/Pausable.sol";
import {IPauserRegistry} from "../../lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IPauserRegistry.sol";

import {ServiceManagerBase, IAVSDirectory, IRewardsCoordinator, IServiceManager} from "../../lib/eigenlayer-middleware/src/ServiceManagerBase.sol";
import {BLSSignatureChecker} from "../../lib/eigenlayer-middleware/src/BLSSignatureChecker.sol";
import {IRegistryCoordinator} from "../../lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry} from "../../lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {IEigenDAThresholdRegistry} from "../interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDARelayRegistry} from "../interfaces/IEigenDARelayRegistry.sol";
import {IPaymentVault} from "../interfaces/IPaymentVault.sol";
import {IEigenDADisperserRegistry} from "../interfaces/IEigenDADisperserRegistry.sol";
import {EigenDAServiceManagerStorage} from "./EigenDAServiceManagerStorage.sol";
import {EigenDAHasher} from "../libraries/EigenDAHasher.sol";
import "../interfaces/IEigenDAStructs.sol";

/**
 * @title EigenDAServiceManager
 * @notice The Service Manager is the central contract of the EigenDA AVS and is responsible for:
 * - accepting and confirming the signature of bridged V1 batches
 * - routing rewards submissions to operators
 * - setting metadata for the AVS
 */
contract EigenDAServiceManager is EigenDAServiceManagerStorage, ServiceManagerBase, BLSSignatureChecker, Pausable {
    using EigenDAHasher for BatchHeader;
    using EigenDAHasher for ReducedBatchHeader;

    uint8 internal constant PAUSED_CONFIRM_BATCH = 0;

    /// @notice when applied to a function, ensures that the function is only callable by the `batchConfirmer`.
    modifier onlyBatchConfirmer() {
        require(isBatchConfirmer[msg.sender]);
        _;
    }

    constructor(
        IAVSDirectory __avsDirectory,
        IRewardsCoordinator __rewardsCoordinator,
        IRegistryCoordinator __registryCoordinator,
        IStakeRegistry __stakeRegistry,
        IEigenDAThresholdRegistry __eigenDAThresholdRegistry,
        IEigenDARelayRegistry __eigenDARelayRegistry,
        IPaymentVault __paymentVault,
        IEigenDADisperserRegistry __eigenDADisperserRegistry
    )
        BLSSignatureChecker(__registryCoordinator)
        ServiceManagerBase(__avsDirectory, __rewardsCoordinator, __registryCoordinator, __stakeRegistry)
        EigenDAServiceManagerStorage(__eigenDAThresholdRegistry, __eigenDARelayRegistry, __paymentVault, __eigenDADisperserRegistry)
    {
        _disableInitializers();
    }

    function initialize(
        IPauserRegistry _pauserRegistry,
        uint256 _initialPausedStatus,
        address _initialOwner,
        address[] memory _batchConfirmers,
        address _rewardsInitiator
    )
        public
        initializer
    {
        _initializePauser(_pauserRegistry, _initialPausedStatus);
        _transferOwnership(_initialOwner);
        _setRewardsInitiator(_rewardsInitiator);
        for (uint i = 0; i < _batchConfirmers.length; ++i) {
            _setBatchConfirmer(_batchConfirmers[i]);
        }
    }

    /**
     * @notice Accepts a batch from the disperser and confirms its signature for V1 bridging
     * @param batchHeader The batch header to confirm
     * @param nonSignerStakesAndSignature The non-signer stakes and signature to confirm the batch with
     */
    function confirmBatch(
        BatchHeader calldata batchHeader,
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature
    ) external onlyWhenNotPaused(PAUSED_CONFIRM_BATCH) onlyBatchConfirmer() {
        // make sure the information needed to derive the non-signers and batch is in calldata to avoid emitting events
        require(tx.origin == msg.sender, "header and nonsigner data must be in calldata");
        // make sure the stakes against which the Batch is being confirmed are not stale
        require(
            batchHeader.referenceBlockNumber < block.number, "specified referenceBlockNumber is in future"
        );

        require(
            (batchHeader.referenceBlockNumber + BLOCK_STALE_MEASURE) >= uint32(block.number),
            "specified referenceBlockNumber is too far in past"
        );

        //make sure that the quorumNumbers and signedStakeForQuorums are of the same length
        require(
            batchHeader.quorumNumbers.length == batchHeader.signedStakeForQuorums.length,
            "quorumNumbers and signedStakeForQuorums must be same length"
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
                "signatories do not own threshold percentage of a quorum"
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

    /**
     * @notice Toggles a batch confirmer role to allow them to confirm batches
     * @param _batchConfirmer The address of the batch confirmer to set
     */
    function setBatchConfirmer(address _batchConfirmer) external onlyOwner() {
        _setBatchConfirmer(_batchConfirmer);
    }

    /// @notice internal function to set a batch confirmer
    function _setBatchConfirmer(address _batchConfirmer) internal {
        isBatchConfirmer[_batchConfirmer] = !isBatchConfirmer[_batchConfirmer];
        emit BatchConfirmerStatusChanged(_batchConfirmer, isBatchConfirmer[_batchConfirmer]);
    }

    /// @notice Returns the current batchId
    function taskNumber() external view returns (uint32) {
        return batchId;
    }

    /**
     * @notice Given a reference block number, returns the block until which operators must serve.
     * @param referenceBlockNumber The reference block number to get the serve until block for
     */
    function latestServeUntilBlock(uint32 referenceBlockNumber) external view returns (uint32) {
        return referenceBlockNumber + STORE_DURATION_BLOCKS + BLOCK_STALE_MEASURE;
    }

    /// @notice Returns an array of bytes where each byte represents the adversary threshold percentage of the quorum at that index for V1 verification
    function quorumAdversaryThresholdPercentages() external view returns (bytes memory) {
        return eigenDAThresholdRegistry.quorumAdversaryThresholdPercentages();
    }

    /// @notice Returns an array of bytes where each byte represents the confirmation threshold percentage of the quorum at that index for V1 verification
    function quorumConfirmationThresholdPercentages() external view returns (bytes memory) {
        return eigenDAThresholdRegistry.quorumConfirmationThresholdPercentages();
    }

    /// @notice Returns an array of bytes where each byte represents the number of a required quorum for V1 verification
    function quorumNumbersRequired() external view returns (bytes memory) {
        return eigenDAThresholdRegistry.quorumNumbersRequired();
    }

    /// @notice Returns the adversary threshold percentage for a quorum for V1 verification
    /// @param quorumNumber The number of the quorum to get the adversary threshold percentage for
    function getQuorumAdversaryThresholdPercentage(
        uint8 quorumNumber
    ) external view returns (uint8){
        return eigenDAThresholdRegistry.getQuorumAdversaryThresholdPercentage(quorumNumber);
    }

    /// @notice Returns the confirmation threshold percentage for a quorum for V1 verification
    /// @param quorumNumber The number of the quorum to get the confirmation threshold percentage for
    function getQuorumConfirmationThresholdPercentage(
        uint8 quorumNumber
    ) external view returns (uint8){
        return eigenDAThresholdRegistry.getQuorumConfirmationThresholdPercentage(quorumNumber);
    }

    /// @notice Returns true if a quorum is required for V1 verification
    /// @param quorumNumber The number of the quorum to check if it is required for V1 verification
    function getIsQuorumRequired(
        uint8 quorumNumber
    ) external view returns (bool){
        return eigenDAThresholdRegistry.getIsQuorumRequired(quorumNumber);
    }

    /// @notice Returns the blob params for a given blob version
    /// @param version The version of the blob to get the params for
    function getBlobParams(uint16 version) external view returns (VersionedBlobParams memory) {
        return eigenDAThresholdRegistry.getBlobParams(version);
    }
}