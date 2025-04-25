// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

interface IEigenDAServiceManagerTest {
    function registryCoordinator() external view returns (address);
    function avsDirectory() external view returns (address);

    function setBatchConfirmer(address _batchConfirmer) external;
    function initialize(
        address _pauserRegistry,
        uint256 _initialPausedStatus,
        address _initialOwner,
        address[] memory _batchConfirmers,
        address _rewardsInitiator
    ) external;
}
