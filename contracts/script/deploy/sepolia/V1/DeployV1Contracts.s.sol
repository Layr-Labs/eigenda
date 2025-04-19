// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {Script, console2} from "forge-std/Script.sol";
import {ProxyAdmin, TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {EmptyContract} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/test/mocks/EmptyContract.sol";
import {
    IPauserRegistry,
    PauserRegistry
} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/permissions/PauserRegistry.sol";

import {ConfigV1Contracts} from "./ConfigV1Contracts.sol";
import {IStakeRegistryTest} from "./interfaces/IStakeRegistryTest.sol";
import {IIndexRegistryTest} from "./interfaces/IIndexRegistryTest.sol";
import {IBLSApkRegistryTest} from "./interfaces/IBLSApkRegistryTest.sol";
import {ISocketRegistryTest} from "./interfaces/ISocketRegistryTest.sol";
import {IRegistryCoordinatorTest} from "./interfaces/IRegistryCoordinatorTest.sol";
import {IEjectionManagerTest} from "./interfaces/IEjectionManagerTest.sol";
import {IEigenDAServiceManagerTest} from "./interfaces/IEigenDAServiceManagerTest.sol";

contract DeployV1Contracts is ConfigV1Contracts {
    function run() public virtual {
        config = _parseConfig();

        vm.startBroadcast();

        // Deploy the proxy admin, empty contract, and pauser registry which is not behind a proxy.
        (ProxyAdmin proxyAdmin, address emptyContract, IPauserRegistry pauserRegistry) = _deployInfrastructure();

        // These need to be deployed first since other contracts rely on knowing its address for construction (including a circular reference)
        deployedContracts.registryCoordinator.proxy =
            address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));
        deployedContracts.eigenDAServiceManager.proxy =
            address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));

        deployedContracts.stakeRegistry = _deployStakeRegistry(proxyAdmin, deployedContracts.registryCoordinator.proxy);
        deployedContracts.indexRegistry = _deployIndexRegistry(proxyAdmin, deployedContracts.registryCoordinator.proxy);
        deployedContracts.blsApkRegistry =
            _deployBLSApkRegistry(proxyAdmin, deployedContracts.registryCoordinator.proxy);
        deployedContracts.socketRegistry =
            _deploySocketRegistry(proxyAdmin, deployedContracts.registryCoordinator.proxy);

        bytes memory rcBytecode = _readContractBytecode(config.registryCoordinator);
        bytes memory rcConstructorArgs = abi.encode(
            deployedContracts.eigenDAServiceManager.proxy,
            deployedContracts.stakeRegistry.proxy,
            deployedContracts.blsApkRegistry.proxy,
            deployedContracts.indexRegistry.proxy
        );
        deployedContracts.registryCoordinator.implementation = _deployContract(rcBytecode, rcConstructorArgs);

        bytes memory rcInitData = _getRegistryCoordinatorInitData(config.initialOwner, address(pauserRegistry));
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(deployedContracts.registryCoordinator.proxy)),
            deployedContracts.registryCoordinator.implementation,
            rcInitData
        );

        bytes memory smBytecode = _readContractBytecode(config.eigenDAServiceManager);
        bytes memory smConstructorArgs = abi.encode(
            config.avsDirectory,
            config.rewardsCoordinator,
            deployedContracts.registryCoordinator.proxy,
            deployedContracts.stakeRegistry.proxy
        );
        deployedContracts.eigenDAServiceManager.implementation = _deployContract(smBytecode, smConstructorArgs);
        bytes memory smInitData = _getServiceManagerInitData(config.initialOwner, address(pauserRegistry));
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(deployedContracts.eigenDAServiceManager.proxy)),
            deployedContracts.eigenDAServiceManager.implementation,
            smInitData
        );

        deployedContracts.ejectionManager = _deployEjectionManager(
            proxyAdmin, deployedContracts.registryCoordinator.proxy, deployedContracts.stakeRegistry.proxy
        );
        if (config.initialOwner != address(0)) {
            proxyAdmin.transferOwnership(config.initialOwner);
        }
        _logDeployedAddresses(proxyAdmin, emptyContract, pauserRegistry);
        vm.stopBroadcast();
    }

    function _deployInfrastructure()
        private
        returns (ProxyAdmin proxyAdmin, address emptyContract, IPauserRegistry pauserRegistry)
    {
        proxyAdmin = new ProxyAdmin();
        emptyContract = address(new EmptyContract());
        pauserRegistry = new PauserRegistry(config.pausers, config.unpauser);
        return (proxyAdmin, emptyContract, pauserRegistry);
    }

    function _deployStakeRegistry(ProxyAdmin proxyAdmin, address registryCoordinatorAddr)
        private
        returns (ContractDeployment memory deployment)
    {
        bytes memory bytecode = _readContractBytecode(config.stakeRegistry);
        bytes memory constructorArgs = abi.encode(registryCoordinatorAddr, config.delegationManager);
        deployment.implementation = _deployContract(bytecode, constructorArgs);
        bytes memory initData = "";
        deployment.proxy =
            address(new TransparentUpgradeableProxy(deployment.implementation, address(proxyAdmin), initData));
        return deployment;
    }

    function _readContractBytecode(ContractSource memory source) private view returns (bytes memory) {
        string memory fullPath = string(abi.encodePacked(source.dir, "/", source.artifact));
        string memory artifact = vm.readFile(fullPath);
        bytes memory bytecode = vm.parseJsonBytes(artifact, ".bytecode.object");
        return bytecode;
    }

    function _deployContract(bytes memory bytecode, bytes memory constructorArgs) private returns (address addr) {
        bytes memory deploymentBytecode = abi.encodePacked(bytecode, constructorArgs);
        assembly {
            addr := create(0, add(deploymentBytecode, 0x20), mload(deploymentBytecode))
            if iszero(extcodesize(addr)) { revert(0, 0) }
        }
    }

    function _deployIndexRegistry(ProxyAdmin proxyAdmin, address registryCoordinatorAddr)
        private
        returns (ContractDeployment memory deployment)
    {
        bytes memory bytecode = _readContractBytecode(config.indexRegistry);
        bytes memory constructorArgs = abi.encode(registryCoordinatorAddr);
        deployment.implementation = _deployContract(bytecode, constructorArgs);
        bytes memory initData = "";
        deployment.proxy =
            address(new TransparentUpgradeableProxy(deployment.implementation, address(proxyAdmin), initData));
        return deployment;
    }

    function _deployBLSApkRegistry(ProxyAdmin proxyAdmin, address registryCoordinatorAddr)
        private
        returns (ContractDeployment memory deployment)
    {
        bytes memory bytecode = _readContractBytecode(config.blsApkRegistry);
        bytes memory constructorArgs = abi.encode(registryCoordinatorAddr);
        deployment.implementation = _deployContract(bytecode, constructorArgs);
        bytes memory initData = "";
        deployment.proxy =
            address(new TransparentUpgradeableProxy(deployment.implementation, address(proxyAdmin), initData));
        return deployment;
    }

    function _deploySocketRegistry(ProxyAdmin proxyAdmin, address registryCoordinatorAddr)
        private
        returns (ContractDeployment memory deployment)
    {
        bytes memory bytecode = _readContractBytecode(config.socketRegistry);
        bytes memory constructorArgs = abi.encode(registryCoordinatorAddr);
        deployment.implementation = _deployContract(bytecode, constructorArgs);
        bytes memory initData = "";
        deployment.proxy =
            address(new TransparentUpgradeableProxy(deployment.implementation, address(proxyAdmin), initData));
        return deployment;
    }

    function _deployEjectionManager(ProxyAdmin proxyAdmin, address registryCoordinatorAddr, address stakeRegistryAddr)
        private
        returns (ContractDeployment memory deployment)
    {
        bytes memory bytecode = _readContractBytecode(config.ejectionManager);
        bytes memory constructorArgs = abi.encode(registryCoordinatorAddr, stakeRegistryAddr);
        deployment.implementation = _deployContract(bytecode, constructorArgs);
        bytes memory initData = _getEjectionManagerInitData(config.initialOwner);
        deployment.proxy =
            address(new TransparentUpgradeableProxy(deployment.implementation, address(proxyAdmin), initData));
        return deployment;
    }

    function _logDeployedAddresses(ProxyAdmin proxyAdmin, address emptyContract, IPauserRegistry pauserRegistry)
        private
        view
    {
        console2.log("\n=== DEPLOYMENT SUMMARY ===\n");
        console2.log("Infrastructure:");
        console2.log("- ProxyAdmin:", address(proxyAdmin));
        console2.log("- Empty Contract:", emptyContract);
        console2.log("- PauserRegistry:", address(pauserRegistry));
        console2.log("\nCore Registry Contracts:");
        console2.log("- StakeRegistry Proxy:", deployedContracts.stakeRegistry.proxy);
        console2.log("- StakeRegistry Implementation:", deployedContracts.stakeRegistry.implementation);
        console2.log("- IndexRegistry Proxy:", deployedContracts.indexRegistry.proxy);
        console2.log("- IndexRegistry Implementation:", deployedContracts.indexRegistry.implementation);
        console2.log("- BLSApkRegistry Proxy:", deployedContracts.blsApkRegistry.proxy);
        console2.log("- BLSApkRegistry Implementation:", deployedContracts.blsApkRegistry.implementation);
        console2.log("- SocketRegistry Proxy:", deployedContracts.socketRegistry.proxy);
        console2.log("- SocketRegistry Implementation:", deployedContracts.socketRegistry.implementation);
        console2.log("\nCoordinator:");
        console2.log("- RegistryCoordinator Proxy:", deployedContracts.registryCoordinator.proxy);
        console2.log("- RegistryCoordinator Implementation:", deployedContracts.registryCoordinator.implementation);
        console2.log("\nOther Contracts:");
        console2.log("- EjectionManager Proxy:", deployedContracts.ejectionManager.proxy);
        console2.log("- EjectionManager Implementation:", deployedContracts.ejectionManager.implementation);
        console2.log("- EigenDAServiceManager Proxy:", deployedContracts.eigenDAServiceManager.proxy);
        console2.log("- EigenDAServiceManager Implementation:", deployedContracts.eigenDAServiceManager.implementation);
    }
}
