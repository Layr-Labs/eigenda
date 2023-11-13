// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin-upgrades/contracts/proxy/utils/Initializable.sol";
import "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";

import "@eigenlayer-core/contracts/interfaces/IDelegationManager.sol";

import "@eigenlayer-core/contracts/libraries/BytesLib.sol";
import "@eigenlayer-core/contracts/libraries/Merkle.sol";
import "@eigenlayer-core/contracts/permissions/Pausable.sol";

import "../libraries/EigenDAHasher.sol";

import "./EigenDAServiceManagerStorage.sol";

/**
 * @title Primary entrypoint for procuring services from EigenDA.
 * @author Layr Labs, Inc.
 * @notice This contract is used for:
 * - initializing the data store by the disperser
 * - confirming the data store by the disperser with inferred aggregated signatures of the quorum
 * - freezing operators as the result of various "challenges"
 */
contract EigenDAServiceManager is Initializable, OwnableUpgradeable, EigenDAServiceManagerStorage, BLSSignatureChecker, Pausable {
    using BytesLib for bytes;
    using EigenDAHasher for BatchHeader;
    using EigenDAHasher for ReducedBatchHeader;

    uint8 internal constant PAUSED_CONFIRM_BATCH = 0;

    /**
     * @notice The EigenLayer delegation contract for this EigenDA which is primarily used by
     * delegators to delegate their stake to operators who would serve as EigenDA
     * nodes and so on.
     * @dev For more details, see DelegationManager.sol.
     */
    IDelegationManager public immutable delegationManager;

    IStrategyManager public immutable strategyManager;

    ISlasher public immutable slasher;

    /// @notice when applied to a function, ensures that the function is only callable by the `registryCoordinator`.
    modifier onlyRegistryCoordinator() {
        require(msg.sender == address(registryCoordinator), "onlyRegistryCoordinator: not from registry coordinator");
        _;
    }

    constructor(
        IBLSRegistryCoordinatorWithIndices _registryCoordinator,
        IStrategyManager _strategyManager,
        IDelegationManager _delegationMananger,
        ISlasher _slasher
    )
        BLSSignatureChecker(_registryCoordinator)
    {
        strategyManager = _strategyManager;
        delegationManager = _delegationMananger;
        slasher = _slasher;
        _disableInitializers();
    }

    function initialize(
        IPauserRegistry _pauserRegistry,
        address initialOwner
    )
        public
        initializer
    {
        _initializePauser(_pauserRegistry, UNPAUSE_ALL);
        _transferOwnership(initialOwner);
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
    ) external onlyWhenNotPaused(PAUSED_CONFIRM_BATCH) {
        // make sure the information needed to derive the non-signers and batch is in calldata to avoid emitting events
        require(tx.origin == msg.sender, "EigenDAServiceManager.confirmBatch: header and nonsigner data must be in calldata");
        // make sure the stakes against which the Batch is being confirmed are not stale
        require(
            batchHeader.referenceBlockNumber <= block.number, "EigenDAServiceManager.confirmBatch: specified referenceBlockNumber is in future"
        );

        require(
            (batchHeader.referenceBlockNumber + BLOCK_STALE_MEASURE) >= uint32(block.number),
            "EigenDAServiceManager.confirmBatch: specified referenceBlockNumber is too far in past"
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
        for (uint i = 0; i < batchHeader.quorumThresholdPercentages.length; i++) {
            // we don't check that the quorumThresholdPercentages are not >100 because a greater value would trivially fail the check, implying 
            // signed stake > total stake
            require(
                quorumStakeTotals.signedStakeForQuorum[i] * THRESHOLD_DENOMINATOR >= 
                    quorumStakeTotals.totalStakeForQuorum[i] * uint8(batchHeader.quorumThresholdPercentages[i]),
                "EigenDAServiceManager.confirmBatch: signatories do not own at least threshold percentage of a quorum"
            );
        }

        // store the metadata hash
        uint96 fee = 0;
        uint32 batchIdMemory = batchId;
        bytes32 batchHeaderHash = batchHeader.hashBatchHeader();
        batchIdToBatchMetadataHash[batchIdMemory] = EigenDAHasher.hashBatchHashedMetadata(batchHeaderHash, signatoryRecordHash, fee, uint32(block.number));

        emit BatchConfirmed(reducedBatchHeaderHash, batchIdMemory, fee);

        // increment the batchId
        batchId = batchIdMemory + 1;
    }

    /// @notice Called in the event of challenge resolution, in order to forward a call to the Slasher, which 'freezes' the `operator`.
    function freezeOperator(address /*operator*/) external {
        revert("EigenDAServiceManager.freezeOperator: not implemented");
        // require(
        //     msg.sender == address(eigenDAChallenge)
        //         || msg.sender == address(eigenDABombVerifier),
        //     "EigenDAServiceManager.freezeOperator: Only challenge resolvers can slash operators"
        // );
        // slasher.freezeOperator(operator);
    }

    /**
     * @notice Called by the Registry in the event of a new registration, to forward a call to the Slasher
     * @param operator The operator whose stake is being updated
     * @param serveUntilBlock The block until which the stake accounted for in the first update is slashable by this middleware
     */ 
    function recordFirstStakeUpdate(address operator, uint32 serveUntilBlock) external onlyRegistryCoordinator {
        // slasher.recordFirstStakeUpdate(operator, serveUntilBlock);
    }

    /** 
     * @notice Called by the registryCoordinator, in order to forward a call to the Slasher, informing it of a stake update
     * @param operator The operator whose stake is being updated
     * @param updateBlock The block at which the update is being made
     * @param serveUntilBlock The block until which the stake withdrawn from the operator in this update is slashable by this middleware
     * @param prevElement The value of the previous element in the linked list of stake updates (generated offchain)
     */
    function recordStakeUpdate(address operator, uint32 updateBlock, uint32 serveUntilBlock, uint256 prevElement) external onlyRegistryCoordinator {
        // slasher.recordStakeUpdate(operator, updateBlock, serveUntilBlock, prevElement);
    }

    /**
     * @notice Called by the registryCoordinator in the event of deregistration, to forward a call to the Slasher
     * @param operator The operator being deregistered
     * @param serveUntilBlock The block until which the stake delegated to the operator is slashable by this middleware
     */ 
    function recordLastStakeUpdateAndRevokeSlashingAbility(address operator, uint32 serveUntilBlock) external onlyRegistryCoordinator {
        // slasher.recordLastStakeUpdateAndRevokeSlashingAbility(operator, serveUntilBlock);
    }

    // VIEW FUNCTIONS
    function taskNumber() external view returns (uint32) {
        return batchId;
    }

    /// @notice Returns the block until which operators must serve.
    function latestServeUntilBlock() external view returns (uint32) {
        return uint32(block.number) + STORE_DURATION_BLOCKS + BLOCK_STALE_MEASURE;
    }

    /// @dev need to override function here since its defined in both these contracts
    function owner() public view override(OwnableUpgradeable, IServiceManager) returns (address) {
        return OwnableUpgradeable.owner();
    }
}