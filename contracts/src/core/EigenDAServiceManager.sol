// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {Initializable} from "@openzeppelin-upgrades/contracts/proxy/utils/Initializable.sol";
import {OwnableUpgradeable} from "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";

import {IDelegationManager} from "eigenlayer-core/contracts/interfaces/IDelegationManager.sol";
import {BytesLib} from "eigenlayer-core/contracts/libraries/BytesLib.sol";
import {Merkle} from "eigenlayer-core/contracts/libraries/Merkle.sol";
import {Pausable} from "eigenlayer-core/contracts/permissions/Pausable.sol";
import {IStrategyManager} from "eigenlayer-core/contracts/interfaces/IStrategyManager.sol";
import {ISlasher} from "eigenlayer-core/contracts/interfaces/ISlasher.sol";
import {IPauserRegistry} from "eigenlayer-core/contracts/interfaces/IPauserRegistry.sol";
import {ISignatureUtils} from "eigenlayer-core/contracts/interfaces/ISignatureUtils.sol";

import {BLSSignatureChecker, IRegistryCoordinator} from "eigenlayer-middleware/BLSSignatureChecker.sol";
import {IServiceManager} from "eigenlayer-middleware/interfaces/IServiceManager.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {BitmapUtils} from "eigenlayer-middleware/libraries/BitmapUtils.sol";
import {EigenDAServiceManagerStorage} from "./EigenDAServiceManagerStorage.sol";
import {EigenDAHasher} from "../libraries/EigenDAHasher.sol";


/**b
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
        IRegistryCoordinator _registryCoordinator,
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

    /**
     * @notice Forwards a call to EigenLayer's DelegationManager contract to confirm operator registration with the AVS
     * @param operator The address of the operator to register.
     * @param operatorSignature The signature, salt, and expiry of the operator's signature.
     */
    function registerOperatorToAVS(
        address operator,
        ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature
    ) external {
        delegationManager.registerOperatorToAVS(operator, operatorSignature);
    }

    /**
     * @notice Forwards a call to EigenLayer's DelegationManager contract to confirm operator deregistration from the AVS
     * @param operator The address of the operator to deregister.
     */
    function deregisterOperatorFromAVS(address operator) external {
        delegationManager.deregisterOperatorFromAVS(operator);
    }

    /**
     * @notice Sets the metadata URI for the AVS
     * @param _metadataURI is the metadata URI for the AVS
     */
    function setMetadataURI(string memory _metadataURI) external onlyOwner() {
        metadataURI = _metadataURI;
    }

    /// @notice Returns the current batchId
    function taskNumber() external view returns (uint32) {
        return batchId;
    }

    /// @notice Returns the block until which operators must serve.
    function latestServeUntilBlock() external view returns (uint32) {
        return uint32(block.number) + STORE_DURATION_BLOCKS + BLOCK_STALE_MEASURE;
    }

    /**
     * @notice Returns the list of strategies that the operator has potentially restaked on the AVS
     * @param operator The address of the operator to get restaked strategies for
     * @dev This function is intended to be called off-chain
     * @dev No guarantee is made on whether the operator has shares for a strategy in a quorum or uniqueness 
     *      of each element in the returned array. The off-chain service should do that validation separately
     */
    function getOperatorRestakedStrategies(address operator) external view returns (address[] memory) {
        bytes32 operatorId = registryCoordinator.getOperatorId(operator);
        uint256 quorumBitmap = registryCoordinator.getCurrentQuorumBitmap(operatorId);
        bytes memory quorumBytesArray = BitmapUtils.bitmapToBytesArray(quorumBitmap);

        uint256 strategiesLength;
        for (uint i = 0; i < quorumBytesArray.length; i++) {
            uint8 quorumNumber = uint8(quorumBytesArray[i]);
            strategiesLength += stakeRegistry.strategyParamsLength(quorumNumber);
        }

        address[] memory restakedStrategies = new address[](strategiesLength);
        uint256 index;
        for (uint i = 0; i < quorumBytesArray.length; i++) {
            uint8 quorumNumber = uint8(quorumBytesArray[i]);
            uint256 strategyParamsLength = stakeRegistry.strategyParamsLength(quorumNumber);
            for (uint j = 0; j < strategyParamsLength; j++) {
                IStakeRegistry.StrategyParams memory strategyParams = stakeRegistry.strategyParamsByIndex(quorumNumber, j);
                restakedStrategies[index] = address(strategyParams.strategy);
                ++index;
            }
        }

        return restakedStrategies;
    }

    /**
     * @notice Returns the list of strategies that the AVS supports for restaking
     * @dev This function is intended to be called off-chain
     * @dev No guarantee is made on uniqueness of each element in the returned array. 
     *      The off-chain service should do that validabution separately
     */
    function getRestakeableStrategies() external view returns (address[] memory) {
        uint256 quorumBitmap;
        uint256 strategiesLength;
        for (uint8 i = 0; i < type(uint8).max; i++) {
            if(stakeRegistry.minimumStakeForQuorum(i) > 0) {
                quorumBitmap = BitmapUtils.setBit(quorumBitmap, i);
                strategiesLength += stakeRegistry.strategyParamsLength(i);
            }
        }

        bytes memory quorumBytesArray = BitmapUtils.bitmapToBytesArray(quorumBitmap);
        address[] memory restakedStrategies = new address[](strategiesLength);
        uint256 index;
        for (uint i = 0; i < quorumBytesArray.length; i++) {
            uint8 quorumNumber = uint8(quorumBytesArray[i]);
            uint256 strategyParamsLength = stakeRegistry.strategyParamsLength(quorumNumber);
            for (uint j = 0; j < strategyParamsLength; j++) {
                IStakeRegistry.StrategyParams memory strategyParams = stakeRegistry.strategyParamsByIndex(quorumNumber, j);
                restakedStrategies[index] = address(strategyParams.strategy);
                ++index;
            }
        }

        return restakedStrategies;
    }
}