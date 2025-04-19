// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

interface IRegistryCoordinatorTest {
    function stakeRegistry() external view returns (address);
    function indexRegistry() external view returns (address);
    function blsApkRegistry() external view returns (address);

    struct OperatorSetParam {
        uint32 maxOperatorCount;
        uint16 kickBIPsOfOperatorStake;
        uint16 kickBIPsOfTotalStake;
    }

    struct StrategyParams {
        address strategy;
        uint96 multiplier;
    }

    function initialize(
        address _initialOwner,
        address _churnApprover,
        address _ejector,
        address _pauserRegistry,
        uint256 _initialPausedStatus,
        OperatorSetParam[] memory _operatorSetParams,
        uint96[] memory _minimumStakes,
        StrategyParams[][] memory _strategyParams
    ) external;
}
