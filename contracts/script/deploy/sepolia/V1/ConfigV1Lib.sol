// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {stdToml} from "forge-std/StdToml.sol";
import {IRegistryCoordinatorTest} from "./interfaces/IRegistryCoordinatorTest.sol";
import {IEjectionManagerTest} from "./interfaces/IEjectionManagerTest.sol";
import {IEigenDAServiceManagerTest} from "./interfaces/IEigenDAServiceManagerTest.sol";

struct ContractSource {
    string dir;
    string artifact;
}

library ConfigV1Lib {
    using stdToml for string;

    // Basic configuration fields
    function delegationManager(string memory config) internal pure returns (address) {
        return stdToml.readAddress(config, ".initParams.existing.delegationManager");
    }

    function avsDirectory(string memory config) internal pure returns (address) {
        return stdToml.readAddress(config, ".initParams.existing.avsDirectory");
    }

    function rewardsCoordinator(string memory config) internal pure returns (address) {
        return stdToml.readAddress(config, ".initParams.existing.rewardsCoordinator");
    }

    function daProxyAdmin(string memory config) internal pure returns (address) {
        return stdToml.readAddress(config, ".initParams.existing.daProxyAdmin");
    }

    function registryCoordinator(string memory config) internal pure returns (address) {
        return stdToml.readAddress(config, ".initParams.existing.registryCoordinator");
    }

    function stakeRegistry(string memory config) internal pure returns (address) {
        return stdToml.readAddress(config, ".initParams.existing.stakeRegistry");
    }

    function serviceManager(string memory config) internal pure returns (address) {
        return stdToml.readAddress(config, ".initParams.existing.serviceManager");
    }

    function blsApkRegistry(string memory config) internal pure returns (address) {
        return stdToml.readAddress(config, ".initParams.existing.blsApkRegistry");
    }

    function indexRegistry(string memory config) internal pure returns (address) {
        return stdToml.readAddress(config, ".initParams.existing.indexRegistry");
    }

    function socketRegistry(string memory config) internal pure returns (address) {
        return stdToml.readAddress(config, ".initParams.existing.socketRegistry");
    }

    function ejectionManager(string memory config) internal pure returns (address) {
        return stdToml.readAddress(config, ".initParams.existing.ejectionManager");
    }

    function initialOwner(string memory config) internal pure returns (address) {
        return stdToml.readAddress(config, ".owner");
    }

    function initialPausedStatus(string memory config) internal pure returns (uint256) {
        return stdToml.readUint(config, ".initParams.core.pauserRegistry.initialPausedStatus");
    }

    // Contract sources
    function stakeRegistrySource(string memory config) internal pure returns (ContractSource memory) {
        return ContractSource({
            dir: stdToml.readString(config, ".sources.stakeRegistry"),
            artifact: stdToml.readString(config, ".sources.stakeRegistryArtifact")
        });
    }

    function indexRegistrySource(string memory config) internal pure returns (ContractSource memory) {
        return ContractSource({
            dir: stdToml.readString(config, ".sources.indexRegistry"),
            artifact: stdToml.readString(config, ".sources.indexRegistryArtifact")
        });
    }

    function socketRegistrySource(string memory config) internal pure returns (ContractSource memory) {
        return ContractSource({
            dir: stdToml.readString(config, ".sources.socketRegistry"),
            artifact: stdToml.readString(config, ".sources.socketRegistryArtifact")
        });
    }

    function blsApkRegistrySource(string memory config) internal pure returns (ContractSource memory) {
        return ContractSource({
            dir: stdToml.readString(config, ".sources.blsApkRegistry"),
            artifact: stdToml.readString(config, ".sources.blsApkRegistryArtifact")
        });
    }

    function registryCoordinatorSource(string memory config) internal pure returns (ContractSource memory) {
        return ContractSource({
            dir: stdToml.readString(config, ".sources.registryCoordinator"),
            artifact: stdToml.readString(config, ".sources.registryCoordinatorArtifact")
        });
    }

    function ejectionManagerSource(string memory config) internal pure returns (ContractSource memory) {
        return ContractSource({
            dir: stdToml.readString(config, ".sources.ejectionManager"),
            artifact: stdToml.readString(config, ".sources.ejectionManagerArtifact")
        });
    }

    function eigenDAServiceManagerSource(string memory config) internal pure returns (ContractSource memory) {
        return ContractSource({
            dir: stdToml.readString(config, ".sources.eigenDAServiceManager"),
            artifact: stdToml.readString(config, ".sources.eigenDAServiceManagerArtifact")
        });
    }

    // PauserRegistry config
    function pausers(string memory config) internal pure returns (address[] memory) {
        return stdToml.readAddressArray(config, ".initParams.core.pauserRegistry.pausers");
    }

    function unpauser(string memory config) internal pure returns (address) {
        return stdToml.readAddress(config, ".initParams.core.pauserRegistry.unpauser");
    }

    // Registry Coordinator config
    function rcChurnApprover(string memory config) internal pure returns (address) {
        return stdToml.readAddress(config, ".initParams.middleware.registryCoordinator.churnApprover");
    }

    function rcEjector(string memory config) internal pure returns (address) {
        return stdToml.readAddress(config, ".initParams.middleware.registryCoordinator.ejector");
    }

    function rcPauserRegistry(string memory config) internal pure returns (address) {
        return stdToml.readAddress(config, ".initParams.middleware.registryCoordinator.pauserRegistry");
    }

    // Registry Coordinator init params
    function operatorSetParams(string memory config)
        internal
        pure
        returns (IRegistryCoordinatorTest.OperatorSetParam[] memory)
    {
        bytes memory operatorSetParamsRaw =
            stdToml.parseRaw(config, ".initParams.middleware.registryCoordinator.operatorSetParams");
        return abi.decode(operatorSetParamsRaw, (IRegistryCoordinatorTest.OperatorSetParam[]));
    }

    function minimumStakes(string memory config) internal pure returns (uint96[] memory) {
        bytes memory minimumStakesRaw =
            stdToml.parseRaw(config, ".initParams.middleware.registryCoordinator.minimumStakes");
        return abi.decode(minimumStakesRaw, (uint96[]));
    }

    function strategyParams(string memory config)
        internal
        pure
        returns (IRegistryCoordinatorTest.StrategyParams[][] memory)
    {
        bytes memory strategyParamsRaw =
            stdToml.parseRaw(config, ".initParams.middleware.registryCoordinator.strategyParams");
        return abi.decode(strategyParamsRaw, (IRegistryCoordinatorTest.StrategyParams[][]));
    }

    // Ejection Manager config
    function ejectors(string memory config) internal pure returns (address[] memory) {
        return stdToml.readAddressArray(config, ".initParams.middleware.ejectionManager.ejectors");
    }

    // Ejection Manager init params
    function quorumEjectionParams(string memory config)
        internal
        pure
        returns (IEjectionManagerTest.QuorumEjectionParams[] memory)
    {
        bytes memory quorumEjectionParamsRaw =
            stdToml.parseRaw(config, ".initParams.middleware.ejectionManager.quorumEjectionParams");
        return abi.decode(quorumEjectionParamsRaw, (IEjectionManagerTest.QuorumEjectionParams[]));
    }

    // Service Manager config
    function smInitialPausedStatus(string memory config) internal pure returns (uint256) {
        return stdToml.readUint(config, ".initParams.eigenDA.serviceManager.initialPausedStatus");
    }

    function smBatchConfirmers(string memory config) internal pure returns (address[] memory) {
        return stdToml.readAddressArray(config, ".initParams.eigenDA.serviceManager.batchConfirmers");
    }

    function smRewardsInitiator(string memory config) internal pure returns (address) {
        return stdToml.readAddress(config, ".initParams.eigenDA.serviceManager.rewardsInitiator");
    }

    // Helper functions for getting initialization data
    function getRegistryCoordinatorInitData(string memory config, address ownerAddr, address pauserRegistryAddr)
        internal
        pure
        returns (bytes memory)
    {
        return abi.encodeCall(
            IRegistryCoordinatorTest.initialize,
            (
                ownerAddr,
                rcChurnApprover(config),
                rcEjector(config),
                pauserRegistryAddr,
                initialPausedStatus(config),
                operatorSetParams(config),
                minimumStakes(config),
                strategyParams(config)
            )
        );
    }

    function getEjectionManagerInitData(string memory config, address ownerAddr) internal pure returns (bytes memory) {
        return
            abi.encodeCall(IEjectionManagerTest.initialize, (ownerAddr, ejectors(config), quorumEjectionParams(config)));
    }

    function getServiceManagerInitData(string memory config, address ownerAddr, address pauserRegistryAddr)
        internal
        pure
        returns (bytes memory)
    {
        return abi.encodeCall(
            IEigenDAServiceManagerTest.initialize,
            (
                pauserRegistryAddr,
                smInitialPausedStatus(config),
                ownerAddr,
                smBatchConfirmers(config),
                smRewardsInitiator(config)
            )
        );
    }
}
