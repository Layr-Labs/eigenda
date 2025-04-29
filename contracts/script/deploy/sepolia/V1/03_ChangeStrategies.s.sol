// SPDX-License-Identifier: BUSL-1.1

pragma solidity =0.8.12;

import {Script, console2} from "forge-std/Script.sol";
import {IStakeRegistryTest} from "./interfaces/IStakeRegistryTest.sol";

contract ChangeStrategy is Script {

    function run() external {
        IStakeRegistryTest stakeRegistry = IStakeRegistryTest(0x2743FaA0df103f7c6D5c339f85c8eE51147C462e);
        IStakeRegistryTest.StrategyParams[] memory strategies = new IStakeRegistryTest.StrategyParams[](1);
        strategies[0] = IStakeRegistryTest.StrategyParams({
            strategy: 0x424246eF71b01ee33aA33aC590fd9a0855F5eFbc,
            multiplier: 1 ether
        });


        vm.startBroadcast();
        stakeRegistry.addStrategies(
            1,
            strategies
        );

        vm.stopBroadcast();

    }
}
