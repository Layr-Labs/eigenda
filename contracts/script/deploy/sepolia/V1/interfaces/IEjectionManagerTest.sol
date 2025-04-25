// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

interface IEjectionManagerTest {
    function registryCoordinator() external view returns (address);
    function stakeRegistry() external view returns (address);
    function setQuorumEjectionParams(uint8 _quorumNumber, QuorumEjectionParams memory _quorumEjectionParams) external;
    function setEjector(address _ejector, bool _status) external;

    struct QuorumEjectionParams {
        uint32 rateLimitWindow; // Time delta to track ejection over
        uint16 ejectableStakePercent; // Max stake to be ejectable per time delta
    }

    function initialize(address _owner, address[] memory _ejectors, QuorumEjectionParams[] memory _quorumEjectionParams)
        external;
}
