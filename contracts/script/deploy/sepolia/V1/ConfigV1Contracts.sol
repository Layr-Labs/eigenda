// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {Script} from "forge-std/Script.sol";
import {stdToml} from "forge-std/StdToml.sol";
import {IRegistryCoordinatorTest} from "./interfaces/IRegistryCoordinatorTest.sol";
import {IEjectionManagerTest} from "./interfaces/IEjectionManagerTest.sol";
import {IEigenDAServiceManagerTest} from "./interfaces/IEigenDAServiceManagerTest.sol";

contract ConfigV1Contracts is Script {
    using stdToml for string;

    struct ContractSource {
        string dir;
        string artifact;
    }

    struct ContractDeployment {
        address proxy;
        address implementation;
    }

    struct DeployedContracts {
        ContractDeployment stakeRegistry;
        ContractDeployment indexRegistry;
        ContractDeployment socketRegistry;
        ContractDeployment blsApkRegistry;
        ContractDeployment registryCoordinator;
        ContractDeployment ejectionManager;
        ContractDeployment eigenDAServiceManager;
    }

    struct DeployConfig {
        address delegationManager;
        address avsDirectory;
        address rewardsCoordinator;
        address initialOwner;
        uint256 initialPausedStatus;
        ContractSource stakeRegistry;
        ContractSource indexRegistry;
        ContractSource socketRegistry;
        ContractSource blsApkRegistry;
        ContractSource registryCoordinator;
        ContractSource ejectionManager;
        ContractSource eigenDAServiceManager;
        address[] pausers;
        address unpauser;
        address rcChurnApprover;
        address rcEjector;
        address rcPauserRegistry;
        address[] ejectors;
        uint256 smInitialPausedStatus;
        address[] smBatchConfirmers;
        address smRewardsInitiator;
    }

    DeployConfig internal config;
    DeployedContracts internal deployedContracts;

    function _parseConfig() internal view returns (DeployConfig memory deployConfig) {
        string memory configPath = vm.envOr(
            "DEPLOY_CONFIG_PATH",
            string("/workspaces/eigenda/contracts/script/deploy/sepolia/V1/config/sepolia.config.toml")
        );
        string memory tomlConfig = vm.readFile(configPath);
        deployConfig.delegationManager = tomlConfig.readAddress(".initParams.shared.delegationManager");
        deployConfig.avsDirectory = tomlConfig.readAddress(".initParams.shared.avsDirectory");
        deployConfig.rewardsCoordinator = tomlConfig.readAddress(".initParams.shared.rewardsCoordinator");
        deployConfig.initialOwner = tomlConfig.readAddress(".initialOwner");
        deployConfig.initialPausedStatus = tomlConfig.readUint(".initParams.shared.initialPausedStatus");
        deployConfig.stakeRegistry.dir = tomlConfig.readString(".sources.stakeRegistry");
        deployConfig.stakeRegistry.artifact = tomlConfig.readString(".sources.stakeRegistryArtifact");
        deployConfig.indexRegistry.dir = tomlConfig.readString(".sources.indexRegistry");
        deployConfig.indexRegistry.artifact = tomlConfig.readString(".sources.indexRegistryArtifact");
        deployConfig.socketRegistry.dir = tomlConfig.readString(".sources.socketRegistry");
        deployConfig.socketRegistry.artifact = tomlConfig.readString(".sources.socketRegistryArtifact");
        deployConfig.blsApkRegistry.dir = tomlConfig.readString(".sources.blsApkRegistry");
        deployConfig.blsApkRegistry.artifact = tomlConfig.readString(".sources.blsApkRegistryArtifact");
        deployConfig.registryCoordinator.dir = tomlConfig.readString(".sources.registryCoordinator");
        deployConfig.registryCoordinator.artifact = tomlConfig.readString(".sources.registryCoordinatorArtifact");
        deployConfig.ejectionManager.dir = tomlConfig.readString(".sources.ejectionManager");
        deployConfig.ejectionManager.artifact = tomlConfig.readString(".sources.ejectionManagerArtifact");
        deployConfig.eigenDAServiceManager.dir = tomlConfig.readString(".sources.eigenDAServiceManager");
        deployConfig.eigenDAServiceManager.artifact = tomlConfig.readString(".sources.eigenDAServiceManagerArtifact");
        deployConfig.pausers = tomlConfig.readAddressArray(".initParams.core.pauserRegistry.pausers");
        deployConfig.unpauser = tomlConfig.readAddress(".initParams.core.pauserRegistry.unpauser");
        deployConfig.rcChurnApprover =
            tomlConfig.readAddress(".initParams.middleware.registryCoordinator.churnApprover");
        deployConfig.rcEjector = tomlConfig.readAddress(".initParams.middleware.registryCoordinator.ejector");
        deployConfig.rcPauserRegistry =
            tomlConfig.readAddress(".initParams.middleware.registryCoordinator.pauserRegistry");
        deployConfig.ejectors = tomlConfig.readAddressArray(".initParams.middleware.ejectionManager.ejectors");
        deployConfig.smInitialPausedStatus =
            tomlConfig.readUint(".initParams.eigenDA.serviceManager.initialPausedStatus");
        deployConfig.smBatchConfirmers =
            tomlConfig.readAddressArray(".initParams.eigenDA.serviceManager.batchConfirmers");
        deployConfig.smRewardsInitiator = tomlConfig.readAddress(".initParams.eigenDA.serviceManager.rewardsInitiator");
        return deployConfig;
    }

    function _getRegistryCoordinatorInitData(address initialOwner, address pauserRegistry)
        internal
        view
        returns (bytes memory)
    {
        string memory configPath = vm.envOr(
            "DEPLOY_CONFIG_PATH",
            string("/workspaces/eigenda/contracts/script/deploy/sepolia/V1/config/sepolia.config.toml")
        );
        string memory tomlConfig = vm.readFile(configPath);
        bytes memory operatorSetParamsRaw =
            stdToml.parseRaw(tomlConfig, ".initParams.middleware.registryCoordinator.operatorSetParams");
        IRegistryCoordinatorTest.OperatorSetParam[] memory operatorSetParams =
            abi.decode(operatorSetParamsRaw, (IRegistryCoordinatorTest.OperatorSetParam[]));
        bytes memory minimumStakesRaw =
            stdToml.parseRaw(tomlConfig, ".initParams.middleware.registryCoordinator.minimumStakes");
        uint96[] memory minimumStakes = abi.decode(minimumStakesRaw, (uint96[]));
        bytes memory strategyParamsRaw =
            stdToml.parseRaw(tomlConfig, ".initParams.middleware.registryCoordinator.strategyParams");
        IRegistryCoordinatorTest.StrategyParams[][] memory strategyParams =
            abi.decode(strategyParamsRaw, (IRegistryCoordinatorTest.StrategyParams[][]));
        return abi.encodeCall(
            IRegistryCoordinatorTest.initialize,
            (
                initialOwner,
                config.rcChurnApprover,
                config.rcEjector,
                pauserRegistry,
                config.initialPausedStatus,
                operatorSetParams,
                minimumStakes,
                strategyParams
            )
        );
    }

    function _getEjectionManagerInitData(address initialOwner) internal view returns (bytes memory) {
        string memory configPath = vm.envOr(
            "DEPLOY_CONFIG_PATH",
            string("/workspaces/eigenda/contracts/script/deploy/sepolia/V1/config/sepolia.config.toml")
        );
        string memory tomlConfig = vm.readFile(configPath);
        bytes memory quorumEjectionParamsRaw =
            stdToml.parseRaw(tomlConfig, ".initParams.middleware.ejectionManager.quorumEjectionParams");
        IEjectionManagerTest.QuorumEjectionParams[] memory quorumEjectionParams =
            abi.decode(quorumEjectionParamsRaw, (IEjectionManagerTest.QuorumEjectionParams[]));
        return abi.encodeCall(IEjectionManagerTest.initialize, (initialOwner, config.ejectors, quorumEjectionParams));
    }

    function _getServiceManagerInitData(address initialOwner, address pauserRegistry)
        internal
        view
        returns (bytes memory)
    {
        return abi.encodeCall(
            IEigenDAServiceManagerTest.initialize,
            (
                pauserRegistry,
                config.smInitialPausedStatus,
                initialOwner,
                config.smBatchConfirmers,
                config.smRewardsInitiator
            )
        );
    }
}
