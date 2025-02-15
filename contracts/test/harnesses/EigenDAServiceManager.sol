// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {BLSSignatureChecker} from "eigenlayer-middleware/BLSSignatureChecker.sol";

import {IEigenDAServiceManager} from "../../src/interfaces/IEigenDAServiceManager.sol";
import {ISignatureUtils} from "eigenlayer-core/contracts/interfaces/ISignatureUtils.sol";


import {IPaymentCoordinator} from "eigenlayer-core/contracts/interfaces/IPaymentCoordinator.sol";


interface IDummyServiceManager is IEigenDAServiceManager {
    // params used to define EigenDA blob verification behaviors
    struct SimulationParams {
        // The offset at which the verification should fail
        uint32 failureOffset;
        // Whether the verification should always fail
        bool alwaysFail;
    }

    function verifyReturn() external payable returns (bool);
}


contract EigenDAServiceManager is IDummyServiceManager {
    SimulationParams public simulationParams;
    uint64 private verifyCallCount;

    constructor(
        SimulationParams memory _simulationParams
    )
    {
        simulationParams = _simulationParams;
    }

    function verifyReturn() external payable override returns (bool) {
        verifyCallCount++;

        if(!simulationParams.alwaysFail && simulationParams.failureOffset == 0) {
            return true;
        }

        if (simulationParams.alwaysFail) {
            return false;
        }

        return verifyCallCount % simulationParams.failureOffset != 0;

    }

    function confirmBatch(
        BatchHeader calldata batchHeader,
        BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature
    ) external {
        revert("A rollup should never attempt to confirm a batch.");
    }

    function setBatchConfirmer(address) external {
        revert("A rollup sequencer contract should never try setting a batch confirmer.");
    }

    /// @notice Returns the current batchId
    function taskNumber() external view returns (uint32) {
        revert("this function should never be triggered by a rollup sequencer contract");
    }

    /// @notice Given a reference block number, returns the block until which operators must serve.
    function latestServeUntilBlock(uint32 referenceBlockNumber) external view returns (uint32) {
        revert("this function should never be triggered by a rollup sequencer contract");
    }

    function BLOCK_STALE_MEASURE() external view returns (uint32) {
        revert("this function should never be triggered by a rollup sequencer contract");
    }
    function avsDirectory() external view returns (address) {
        revert("this function should never be triggered by a rollup sequencer contract");
    }
    function batchIdToBatchMetadataHash(uint32 batchId) external view returns(bytes32) {
revert("this function should never be triggered by a rollup sequencer contract");
    }
    function deregisterOperatorFromAVS(address operator) external {
        revert("this function should never be triggered by a rollup sequencer contract");
    }
    function getOperatorRestakedStrategies(address operator) external view returns (address[] memory) {
        revert("this function should never be triggered by a rollup sequencer contract");
    }
    function getRestakeableStrategies() external view returns (address[] memory) {
        revert("this function should never be triggered by a rollup sequencer contract");
    }
    function payForRange(IPaymentCoordinator.RangePayment[] calldata rangePayments) external {
        revert("this function should never be triggered by a rollup sequencer contract");
    }
    function quorumAdversaryThresholdPercentages() external view returns (bytes memory) {
        revert("this function should never be triggered by a rollup sequencer contract");
    }
    function quorumConfirmationThresholdPercentages() external view returns (bytes memory) {
        revert("this function should never be triggered by a rollup sequencer contract");
    }
    function quorumNumbersRequired() external view returns (bytes memory) {
        revert("this function should never be triggered by a rollup sequencer contract");
    }
    function registerOperatorToAVS(
        address operator,
        ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature
    ) external {
        revert("this function should never be triggered by a rollup sequencer contract");
    }
    function updateAVSMetadataURI(string memory _metadataURI) external {
                revert("this function should never be triggered by a rollup sequencer contract");
    }

}