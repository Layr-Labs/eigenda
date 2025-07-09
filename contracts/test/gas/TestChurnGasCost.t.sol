// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {IIndexRegistry} from "lib/eigenlayer-middleware/src/interfaces/IIndexRegistry.sol";
import {IStakeRegistry} from "lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";

import {console2} from "forge-std/console2.sol";

import {Test} from "forge-std/Test.sol";

/// @notice This test is meant to be run on mainnet.
contract TestChurnGasCost is Test {
    IIndexRegistry constant INDEX_REGISTRY = IIndexRegistry(0xBd35a7a1CDeF403a6a99e4E8BA0974D198455030);
    IStakeRegistry constant STAKE_REGISTRY = IStakeRegistry(0x006124Ae7976137266feeBFb3F4D2BE4C073139D);

    function test_churnGasCost() public {
        // get a list of all operators
        vm.startSnapshotGas("OPERATOR_LIST");
        bytes32[] memory operators = INDEX_REGISTRY.getOperatorListAtBlockNumber(0, uint32(block.number));
        console2.log("Number of operators: ", operators.length);
        vm.stopSnapshotGas("OPERATOR_LIST");


        // fetch the stakes of each operator
        vm.startSnapshotGas("OPERATOR_STAKES");
        uint96[] memory stakes = new uint96[](operators.length);
        for (uint256 i; i < operators.length; i++) {
            stakes[i] = STAKE_REGISTRY.getStakeAtBlockNumber(operators[i], 0, uint32(block.number));
        }
        vm.stopSnapshotGas("OPERATOR_STAKES");
    }
}
