// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {Pausable} from "eigenlayer-core/contracts/permissions/Pausable.sol";
import {IPauserRegistry} from "eigenlayer-core/contracts/interfaces/IPauserRegistry.sol";

import {ServiceManagerBase, IAVSDirectory, IRewardsCoordinator, IServiceManager} from "eigenlayer-middleware/ServiceManagerBase.sol";
import {BLSSignatureChecker} from "eigenlayer-middleware/BLSSignatureChecker.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {IEigenDAThresholdRegistry} from "../interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDARelayRegistry} from "../interfaces/IEigenDARelayRegistry.sol";
import {IPaymentVault} from "../interfaces/IPaymentVault.sol";
import {IEigenDADisperserRegistry} from "../interfaces/IEigenDADisperserRegistry.sol";
import {EigenDAServiceManagerStorage} from "./EigenDAServiceManagerStorage.sol";
import {EigenDAHasher} from "../libraries/EigenDAHasher.sol";
import "../interfaces/IEigenDAStructs.sol";

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

    /// @notice Returns the bytes array of quorumAdversaryThresholdPercentages
    function quorumAdversaryThresholdPercentages() external view returns (bytes memory) {
        return eigenDAThresholdRegistry.quorumAdversaryThresholdPercentages();
    }

    /// @notice Returns the bytes array of quorumAdversaryThresholdPercentages
    function quorumConfirmationThresholdPercentages() external view returns (bytes memory) {
        return eigenDAThresholdRegistry.quorumConfirmationThresholdPercentages();
    }

    /// @notice Returns the bytes array of quorumsNumbersRequired
    function quorumNumbersRequired() external view returns (bytes memory) {
        return eigenDAThresholdRegistry.quorumNumbersRequired();
    }

    function getQuorumAdversaryThresholdPercentage(
        uint8 quorumNumber
    ) external view returns (uint8){
        return eigenDAThresholdRegistry.getQuorumAdversaryThresholdPercentage(quorumNumber);
    }

    /// @notice Gets the confirmation threshold percentage for a quorum
    function getQuorumConfirmationThresholdPercentage(
        uint8 quorumNumber
    ) external view returns (uint8){
        return eigenDAThresholdRegistry.getQuorumConfirmationThresholdPercentage(quorumNumber);
    }

    /// @notice Checks if a quorum is required
    function getIsQuorumRequired(
        uint8 quorumNumber
    ) external view returns (bool){
        return eigenDAThresholdRegistry.getIsQuorumRequired(quorumNumber);
    }

    /// @notice Gets the default security thresholds for V2
    function getDefaultSecurityThresholdsV2() external view returns (SecurityThresholds memory) {
        return eigenDAThresholdRegistry.getDefaultSecurityThresholdsV2();
    }

    /// @notice Returns the blob params for a given blob version
    function getBlobParams(uint16 version) external view returns (VersionedBlobParams memory) {
        return eigenDAThresholdRegistry.getBlobParams(version);
    }
}