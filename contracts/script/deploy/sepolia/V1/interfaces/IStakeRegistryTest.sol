// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

interface IStakeRegistryTest {

    struct StrategyParams {
        address strategy;
        uint96 multiplier;
    }
    function removeStrategies(uint8 quorumNumber, uint256[] calldata indicesToRemove) external;
    function addStrategies(
        uint8 quorumNumber,
        StrategyParams[] memory strategyParams
    ) external;
    function registryCoordinator() external view returns (address);
    function delegation() external view returns (address);
}
