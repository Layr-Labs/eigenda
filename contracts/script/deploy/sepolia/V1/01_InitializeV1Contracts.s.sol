// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {Script, console2} from "forge-std/Script.sol";
import {IRegistryCoordinatorTest} from "./interfaces/IRegistryCoordinatorTest.sol";
import {IEjectionManagerTest} from "./interfaces/IEjectionManagerTest.sol";
import {IEigenDAServiceManagerTest} from "./interfaces/IEigenDAServiceManagerTest.sol";

contract InitializeV1Contracts is Script {
    address constant OWNER = 0xDF291ebfe90eF9187c3f45609603E366a21a16Ea;

    address constant STETH_STRATEGY = 0x8b29d91e67b013e855EaFe0ad704aC4Ab086a574;
    address constant WETH_STRATEGY = 0x424246eF71b01ee33aA33aC590fd9a0855F5eFbc;
    address constant EIGEN_BEIGEN_STRATEGY = 0x8E93249a6C37a32024756aaBd813E6139b17D1d5;

    address constant REGISTRY_COORDINATOR = 0xAF21d3811B5d23D5466AC83BA7a9c34c261A8D81;
    address constant EJECTION_MANAGER = 0xc9d4541C409f15C0408c022D7e8C3F37Ac960f66;
    address constant SERVICE_MANAGER = 0x3a5acf46ba6890B8536420F4900AC9BC45Df4764;

    function run() external {
        vm.startBroadcast(OWNER);

        createQuorums();

        IEjectionManagerTest(EJECTION_MANAGER).setEjector(0x1424CC9b5013d3902487E54Aa9F13C74134e6637, true);
        IEigenDAServiceManagerTest(SERVICE_MANAGER).setBatchConfirmer(0xd6F7B4fEbE08454BFE5aF0bE5fD40A2F8Ee80b6c);

        vm.stopBroadcast();
    }

    struct QuorumInitParams {
        IRegistryCoordinatorTest.OperatorSetParam operatorSetParam;
        uint96 minStake;
        IRegistryCoordinatorTest.StrategyParams[] strategyParams;
        IEjectionManagerTest.QuorumEjectionParams ejectorParams;
    }

    function createQuorums() internal {
        QuorumInitParams memory params0 = quorum0();
        QuorumInitParams memory params1 = quorum1();
        IRegistryCoordinatorTest(REGISTRY_COORDINATOR).createQuorum(
            params0.operatorSetParam, params0.minStake, params0.strategyParams
        );
        IEjectionManagerTest(EJECTION_MANAGER).setQuorumEjectionParams(0, params0.ejectorParams);
        IRegistryCoordinatorTest(REGISTRY_COORDINATOR).createQuorum(
            params1.operatorSetParam, params1.minStake, params1.strategyParams
        );
        IEjectionManagerTest(EJECTION_MANAGER).setQuorumEjectionParams(1, params1.ejectorParams);
    }

    function quorum0() internal pure returns (QuorumInitParams memory params) {
        IRegistryCoordinatorTest.StrategyParams[] memory strategyParams =
            new IRegistryCoordinatorTest.StrategyParams[](2);
        strategyParams[0] = IRegistryCoordinatorTest.StrategyParams({strategy: STETH_STRATEGY, multiplier: 1 ether});
        strategyParams[1] = IRegistryCoordinatorTest.StrategyParams({strategy: WETH_STRATEGY, multiplier: 1 ether});

        params = QuorumInitParams({
            operatorSetParam: IRegistryCoordinatorTest.OperatorSetParam({
                maxOperatorCount: 200,
                kickBIPsOfOperatorStake: 11000,
                kickBIPsOfTotalStake: 50
            }),
            minStake: 1 ether,
            strategyParams: strategyParams,
            ejectorParams: IEjectionManagerTest.QuorumEjectionParams({rateLimitWindow: 7 days, ejectableStakePercent: 3000})
        });
    }

    function quorum1() internal pure returns (QuorumInitParams memory params) {
        IRegistryCoordinatorTest.StrategyParams[] memory strategyParams =
            new IRegistryCoordinatorTest.StrategyParams[](1);
        strategyParams[0] =
            IRegistryCoordinatorTest.StrategyParams({strategy: EIGEN_BEIGEN_STRATEGY, multiplier: 1 ether});

        params = QuorumInitParams({
            operatorSetParam: IRegistryCoordinatorTest.OperatorSetParam({
                maxOperatorCount: 30,
                kickBIPsOfOperatorStake: 11000,
                kickBIPsOfTotalStake: 334
            }),
            minStake: 1 ether,
            strategyParams: strategyParams,
            ejectorParams: IEjectionManagerTest.QuorumEjectionParams({rateLimitWindow: 7 days, ejectableStakePercent: 3000})
        });
    }
}
