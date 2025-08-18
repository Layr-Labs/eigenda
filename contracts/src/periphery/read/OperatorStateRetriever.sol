// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import {IRegistryCoordinator} from "./interfaces/IRegistryCoordinator.sol";
import {ISocketRegistry} from "./interfaces/ISocketRegistry.sol";
import {IBLSApkRegistry} from "./interfaces/IBLSApkRegistry.sol";
import {IStakeRegistry} from "./interfaces/IStakeRegistry.sol";
import {IIndexRegistry} from "./interfaces/IIndexRegistry.sol";

import {BitmapUtils} from "./libraries/BitmapUtils.sol";

/**
 * @title OperatorStateRetriever with view functions that allow to retrieve the state of an AVSs registry system.
 * @author Layr Labs Inc.
 */
contract OperatorStateRetriever {
    struct Operator {
        address operator;
        bytes32 operatorId;
        uint96 stake;
    }

    struct CheckSignaturesIndices {
        uint32[] nonSignerQuorumBitmapIndices;
        uint32[] quorumApkIndices;
        uint32[] totalStakeIndices;
        uint32[][] nonSignerStakeIndices; // nonSignerStakeIndices[quorumNumberIndex][nonSignerIndex]
    }

    /**
     * @notice This function is intended to to be called by AVS operators every time a new task is created (i.e.)
     * the AVS coordinator makes a request to AVS operators. Since all of the crucial information is kept onchain,
     * operators don't need to run indexers to fetch the data.
     * @param registryCoordinator is the registry coordinator to fetch the AVS registry information from
     * @param operatorId the id of the operator to fetch the quorums lists
     * @param blockNumber is the block number to get the operator state for
     * @return 1) the quorumBitmap of the operator at the given blockNumber
     *         2) 2d array of Operator structs. For each quorum the provided operator
     *            was a part of at `blockNumber`, an ordered list of operators.
     */
    function getOperatorState(IRegistryCoordinator registryCoordinator, bytes32 operatorId, uint32 blockNumber)
        external
        view
        returns (uint256, Operator[][] memory)
    {
        bytes32[] memory operatorIds = new bytes32[](1);
        operatorIds[0] = operatorId;
        uint256 index = registryCoordinator.getQuorumBitmapIndicesAtBlockNumber(blockNumber, operatorIds)[0];

        uint256 quorumBitmap = registryCoordinator.getQuorumBitmapAtBlockNumberByIndex(operatorId, blockNumber, index);

        bytes memory quorumNumbers = BitmapUtils.bitmapToBytesArray(quorumBitmap);

        return (quorumBitmap, getOperatorState(registryCoordinator, quorumNumbers, blockNumber));
    }

    /**
     * @notice returns the ordered list of operators (id and stake) for each quorum. The AVS coordinator
     * may call this function directly to get the operator state for a given block number
     * @param registryCoordinator is the registry coordinator to fetch the AVS registry information from
     * @param quorumNumbers are the ids of the quorums to get the operator state for
     * @param blockNumber is the block number to get the operator state for
     * @return 2d array of Operators. For each quorum, an ordered list of Operators
     */
    function getOperatorState(IRegistryCoordinator registryCoordinator, bytes memory quorumNumbers, uint32 blockNumber)
        public
        view
        returns (Operator[][] memory)
    {
        IStakeRegistry stakeRegistry = registryCoordinator.stakeRegistry();
        IIndexRegistry indexRegistry = registryCoordinator.indexRegistry();
        IBLSApkRegistry blsApkRegistry = registryCoordinator.blsApkRegistry();

        Operator[][] memory operators = new Operator[][](quorumNumbers.length);
        for (uint256 i = 0; i < quorumNumbers.length; i++) {
            uint8 quorumNumber = uint8(quorumNumbers[i]);
            bytes32[] memory operatorIds = indexRegistry.getOperatorListAtBlockNumber(quorumNumber, blockNumber);
            operators[i] = new Operator[](operatorIds.length);
            for (uint256 j = 0; j < operatorIds.length; j++) {
                operators[i][j] = Operator({
                    operator: blsApkRegistry.getOperatorFromPubkeyHash(operatorIds[j]),
                    operatorId: bytes32(operatorIds[j]),
                    stake: stakeRegistry.getStakeAtBlockNumber(bytes32(operatorIds[j]), quorumNumber, blockNumber)
                });
            }
        }

        return operators;
    }

    /**
     * @notice This function is intended to to be called by AVS operators every time a new task is created (i.e.)
     * the AVS coordinator makes a request to AVS operators. Since all of the crucial information is kept onchain,
     * operators don't need to run indexers to fetch the data.
     * @param registryCoordinator is the registry coordinator to fetch the AVS registry information from
     * @param operatorId the id of the operator to fetch the quorums lists
     * @param blockNumber is the block number to get the operator state for
     * @return quorumBitmap the quorumBitmap of the operator at the given blockNumber
     * @return operators a 2d array of Operators. For each quorum, an ordered list of Operators
     * @return sockets a 2d array of sockets. For each quorum, an ordered list of sockets
     */
    function getOperatorStateWithSocket(
        IRegistryCoordinator registryCoordinator,
        bytes32 operatorId,
        uint32 blockNumber
    ) external view returns (uint256 quorumBitmap, Operator[][] memory operators, string[][] memory sockets) {
        bytes32[] memory operatorIds = new bytes32[](1);
        operatorIds[0] = operatorId;
        uint256 index = registryCoordinator.getQuorumBitmapIndicesAtBlockNumber(blockNumber, operatorIds)[0];

        quorumBitmap = registryCoordinator.getQuorumBitmapAtBlockNumberByIndex(operatorId, blockNumber, index);

        bytes memory quorumNumbers = BitmapUtils.bitmapToBytesArray(quorumBitmap);

        (operators, sockets) = getOperatorStateWithSocket(registryCoordinator, quorumNumbers, blockNumber);
    }

    /// @dev Used below to avoid stack too deep.
    struct Registries {
        IStakeRegistry stakeRegistry;
        IIndexRegistry indexRegistry;
        IBLSApkRegistry blsApkRegistry;
        ISocketRegistry socketRegistry;
    }

    /**
     * @notice returns the ordered list of operators (id, stake, socket) for each quorum. The AVS coordinator
     * may call this function directly to get the operator state for a given block number
     * @param registryCoordinator is the registry coordinator to fetch the AVS registry information from
     * @param quorumNumbers are the ids of the quorums to get the operator state for
     * @param blockNumber is the block number to get the operator state for
     * @return operators a 2d array of Operators. For each quorum, an ordered list of Operators
     * @return sockets a 2d array of sockets. For each quorum, an ordered list of sockets
     */
    function getOperatorStateWithSocket(
        IRegistryCoordinator registryCoordinator,
        bytes memory quorumNumbers,
        uint32 blockNumber
    ) public view returns (Operator[][] memory operators, string[][] memory sockets) {
        Registries memory registries = Registries({
            stakeRegistry: registryCoordinator.stakeRegistry(),
            indexRegistry: registryCoordinator.indexRegistry(),
            blsApkRegistry: registryCoordinator.blsApkRegistry(),
            socketRegistry: registryCoordinator.socketRegistry()
        });

        operators = new Operator[][](quorumNumbers.length);
        sockets = new string[][](quorumNumbers.length);
        for (uint256 i = 0; i < quorumNumbers.length; i++) {
            uint8 quorumNumber = uint8(quorumNumbers[i]);
            bytes32[] memory operatorIds =
                registries.indexRegistry.getOperatorListAtBlockNumber(quorumNumber, blockNumber);
            operators[i] = new Operator[](operatorIds.length);
            sockets[i] = new string[](operatorIds.length);
            for (uint256 j = 0; j < operatorIds.length; j++) {
                operators[i][j] = Operator({
                    operator: registries.blsApkRegistry.getOperatorFromPubkeyHash(operatorIds[j]),
                    operatorId: bytes32(operatorIds[j]),
                    stake: registries.stakeRegistry.getStakeAtBlockNumber(
                        bytes32(operatorIds[j]), quorumNumber, blockNumber
                    )
                });
                sockets[i][j] = registries.socketRegistry.getOperatorSocket(bytes32(operatorIds[j]));
            }
        }
    }

    /**
     * @notice this is called by the AVS operator to get the relevant indices for the checkSignatures function
     * if they are not running an indexer
     * @param registryCoordinator is the registry coordinator to fetch the AVS registry information from
     * @param referenceBlockNumber is the block number to get the indices for
     * @param quorumNumbers are the ids of the quorums to get the operator state for
     * @param nonSignerOperatorIds are the ids of the nonsigning operators
     * @return 1) the indices of the quorumBitmaps for each of the operators in the @param nonSignerOperatorIds array at the given blocknumber
     *         2) the indices of the total stakes entries for the given quorums at the given blocknumber
     *         3) the indices of the stakes of each of the nonsigners in each of the quorums they were a
     *            part of (for each nonsigner, an array of length the number of quorums they were a part of
     *            that are also part of the provided quorumNumbers) at the given blocknumber
     *         4) the indices of the quorum apks for each of the provided quorums at the given blocknumber
     */
    function getCheckSignaturesIndices(
        IRegistryCoordinator registryCoordinator,
        uint32 referenceBlockNumber,
        bytes calldata quorumNumbers,
        bytes32[] calldata nonSignerOperatorIds
    ) external view returns (CheckSignaturesIndices memory) {
        IStakeRegistry stakeRegistry = registryCoordinator.stakeRegistry();
        CheckSignaturesIndices memory checkSignaturesIndices;

        // get the indices of the quorumBitmap updates for each of the operators in the nonSignerOperatorIds array
        checkSignaturesIndices.nonSignerQuorumBitmapIndices =
            registryCoordinator.getQuorumBitmapIndicesAtBlockNumber(referenceBlockNumber, nonSignerOperatorIds);

        // get the indices of the totalStake updates for each of the quorums in the quorumNumbers array
        checkSignaturesIndices.totalStakeIndices =
            stakeRegistry.getTotalStakeIndicesAtBlockNumber(referenceBlockNumber, quorumNumbers);

        checkSignaturesIndices.nonSignerStakeIndices = new uint32[][](quorumNumbers.length);
        for (uint8 quorumNumberIndex = 0; quorumNumberIndex < quorumNumbers.length; quorumNumberIndex++) {
            uint256 numNonSignersForQuorum = 0;
            // this array's length will be at most the number of nonSignerOperatorIds, this will be trimmed after it is filled
            checkSignaturesIndices.nonSignerStakeIndices[quorumNumberIndex] = new uint32[](nonSignerOperatorIds.length);

            for (uint256 i = 0; i < nonSignerOperatorIds.length; i++) {
                // get the quorumBitmap for the operator at the given blocknumber and index
                uint192 nonSignerQuorumBitmap = registryCoordinator.getQuorumBitmapAtBlockNumberByIndex(
                    nonSignerOperatorIds[i],
                    referenceBlockNumber,
                    checkSignaturesIndices.nonSignerQuorumBitmapIndices[i]
                );

                require(
                    nonSignerQuorumBitmap != 0,
                    "OperatorStateRetriever.getCheckSignaturesIndices: operator must be registered at blocknumber"
                );

                // if the operator was a part of the quorum and the quorum is a part of the provided quorumNumbers
                if ((nonSignerQuorumBitmap >> uint8(quorumNumbers[quorumNumberIndex])) & 1 == 1) {
                    // get the index of the stake update for the operator at the given blocknumber and quorum number
                    checkSignaturesIndices.nonSignerStakeIndices[quorumNumberIndex][numNonSignersForQuorum] =
                    stakeRegistry.getStakeUpdateIndexAtBlockNumber(
                        nonSignerOperatorIds[i], uint8(quorumNumbers[quorumNumberIndex]), referenceBlockNumber
                    );
                    numNonSignersForQuorum++;
                }
            }

            // resize the array to the number of nonSigners for this quorum
            uint32[] memory nonSignerStakeIndicesForQuorum = new uint32[](numNonSignersForQuorum);
            for (uint256 i = 0; i < numNonSignersForQuorum; i++) {
                nonSignerStakeIndicesForQuorum[i] = checkSignaturesIndices.nonSignerStakeIndices[quorumNumberIndex][i];
            }
            checkSignaturesIndices.nonSignerStakeIndices[quorumNumberIndex] = nonSignerStakeIndicesForQuorum;
        }

        IBLSApkRegistry blsApkRegistry = registryCoordinator.blsApkRegistry();
        // get the indices of the quorum apks for each of the provided quorums at the given blocknumber
        checkSignaturesIndices.quorumApkIndices =
            blsApkRegistry.getApkIndicesAtBlockNumber(quorumNumbers, referenceBlockNumber);

        return checkSignaturesIndices;
    }

    /**
     * @notice this function returns the quorumBitmaps for each of the operators in the operatorIds array at the given blocknumber
     * @param registryCoordinator is the AVS registry coordinator to fetch the operator information from
     * @param operatorIds are the ids of the operators to get the quorumBitmaps for
     * @param blockNumber is the block number to get the quorumBitmaps for
     */
    function getQuorumBitmapsAtBlockNumber(
        IRegistryCoordinator registryCoordinator,
        bytes32[] memory operatorIds,
        uint32 blockNumber
    ) external view returns (uint256[] memory) {
        uint32[] memory quorumBitmapIndices =
            registryCoordinator.getQuorumBitmapIndicesAtBlockNumber(blockNumber, operatorIds);
        uint256[] memory quorumBitmaps = new uint256[](operatorIds.length);
        for (uint256 i = 0; i < operatorIds.length; i++) {
            quorumBitmaps[i] = registryCoordinator.getQuorumBitmapAtBlockNumberByIndex(
                operatorIds[i], blockNumber, quorumBitmapIndices[i]
            );
        }
        return quorumBitmaps;
    }

    /**
     * @notice This function returns the operatorIds for each of the operators in the operators array
     * @param registryCoordinator is the AVS registry coordinator to fetch the operator information from
     * @param operators is the array of operator address to get corresponding operatorIds for
     * @dev if an operator is not registered, the operatorId will be 0
     */
    function getBatchOperatorId(IRegistryCoordinator registryCoordinator, address[] memory operators)
        external
        view
        returns (bytes32[] memory operatorIds)
    {
        operatorIds = new bytes32[](operators.length);
        for (uint256 i = 0; i < operators.length; ++i) {
            operatorIds[i] = registryCoordinator.getOperatorId(operators[i]);
        }
    }

    /**
     * @notice This function returns the operator addresses for each of the operators in the operatorIds array
     * @param registryCoordinator is the AVS registry coordinator to fetch the operator information from
     * @param operators is the array of operatorIds to get corresponding operator addresses for
     * @dev if an operator is not registered, the operator address will be 0
     */
    function getBatchOperatorFromId(IRegistryCoordinator registryCoordinator, bytes32[] memory operatorIds)
        external
        view
        returns (address[] memory operators)
    {
        operators = new address[](operatorIds.length);
        for (uint256 i = 0; i < operatorIds.length; ++i) {
            operators[i] = registryCoordinator.getOperatorFromId(operatorIds[i]);
        }
    }
}
