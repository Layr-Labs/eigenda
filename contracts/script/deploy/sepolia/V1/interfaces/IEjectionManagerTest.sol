// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

interface IEjectionManagerTest {
    function registryCoordinator() external view returns (address);
    function stakeRegistry() external view returns (address);

    struct QuorumEjectionParams {
        uint32 threshold;
        uint16 ejectionTimestamp;
    }

    function initialize(address _owner, address[] memory _ejectors, QuorumEjectionParams[] memory _quorumEjectionParams)
        external;
}
