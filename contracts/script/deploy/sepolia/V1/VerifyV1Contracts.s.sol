// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {Script} from "forge-std/Script.sol";
import {DeployV1Contracts} from "./DeployV1Contracts.s.sol";
import {IStakeRegistryTest} from "./interfaces/IStakeRegistryTest.sol";
import {IIndexRegistryTest} from "./interfaces/IIndexRegistryTest.sol";
import {IBLSApkRegistryTest} from "./interfaces/IBLSApkRegistryTest.sol";
import {ISocketRegistryTest} from "./interfaces/ISocketRegistryTest.sol";
import {IEjectionManagerTest} from "./interfaces/IEjectionManagerTest.sol";
import {IEigenDAServiceManagerTest} from "./interfaces/IEigenDAServiceManagerTest.sol";

contract VerifyV1Contracts is DeployV1Contracts {
    function run() public override {
        super.run();

        _verifyStakeRegistry(
            deployedContracts.stakeRegistry.proxy, deployedContracts.registryCoordinator.proxy, config.delegationManager
        );

        _verifyIndexRegistry(deployedContracts.indexRegistry.proxy, deployedContracts.registryCoordinator.proxy);

        _verifyBLSApkRegistry(deployedContracts.blsApkRegistry.proxy, deployedContracts.registryCoordinator.proxy);

        _verifySocketRegistry(deployedContracts.socketRegistry.proxy, deployedContracts.registryCoordinator.proxy);

        _verifyEjectionManager(deployedContracts.ejectionManager.proxy, deployedContracts.registryCoordinator.proxy);

        _verifyServiceManager(
            deployedContracts.eigenDAServiceManager.proxy,
            deployedContracts.registryCoordinator.proxy,
            config.avsDirectory
        );
    }

    function _verifyStakeRegistry(
        address stakeRegistryProxy,
        address expectedRegistryCoordinator,
        address expectedDelegationManager
    ) private view {
        IStakeRegistryTest stakeRegistry = IStakeRegistryTest(stakeRegistryProxy);
        address actualRegistryCoordinator = stakeRegistry.registryCoordinator();
        address actualDelegationManager = stakeRegistry.delegation();
        require(
            actualRegistryCoordinator == expectedRegistryCoordinator,
            "StakeRegistry has incorrect registryCoordinator address"
        );
        require(
            actualDelegationManager == expectedDelegationManager,
            "StakeRegistry has incorrect delegationManager address"
        );
    }

    function _verifyIndexRegistry(address indexRegistryProxy, address expectedRegistryCoordinator) private view {
        IIndexRegistryTest indexRegistry = IIndexRegistryTest(indexRegistryProxy);
        address actualRegistryCoordinator = indexRegistry.registryCoordinator();
        require(
            actualRegistryCoordinator == expectedRegistryCoordinator,
            "IndexRegistry has incorrect registryCoordinator address"
        );
    }

    function _verifyBLSApkRegistry(address blsApkRegistryProxy, address expectedRegistryCoordinator) private view {
        IBLSApkRegistryTest blsApkRegistry = IBLSApkRegistryTest(blsApkRegistryProxy);
        address actualRegistryCoordinator = blsApkRegistry.registryCoordinator();
        require(
            actualRegistryCoordinator == expectedRegistryCoordinator,
            "BLSApkRegistry has incorrect registryCoordinator address"
        );
    }

    function _verifySocketRegistry(address socketRegistryProxy, address expectedRegistryCoordinator) private view {
        ISocketRegistryTest socketRegistry = ISocketRegistryTest(socketRegistryProxy);
        address actualRegistryCoordinator = socketRegistry.registryCoordinator();
        require(
            actualRegistryCoordinator == expectedRegistryCoordinator,
            "SocketRegistry has incorrect registryCoordinator address"
        );
    }

    function _verifyEjectionManager(address ejectionManagerProxy, address expectedRegistryCoordinator) private view {
        IEjectionManagerTest ejectionManager = IEjectionManagerTest(ejectionManagerProxy);
        address actualRegistryCoordinator = ejectionManager.registryCoordinator();
        address actualStakeRegistry = ejectionManager.stakeRegistry();
        require(
            actualRegistryCoordinator == expectedRegistryCoordinator,
            "EjectionManager has incorrect registryCoordinator address"
        );
        require(
            actualStakeRegistry == deployedContracts.stakeRegistry.proxy,
            "EjectionManager has incorrect stakeRegistry address"
        );
    }

    function _verifyServiceManager(
        address serviceManagerProxy,
        address expectedRegistryCoordinator,
        address expectedAVSDirectory
    ) private view {
        IEigenDAServiceManagerTest serviceManager = IEigenDAServiceManagerTest(serviceManagerProxy);
        address actualRegistryCoordinator = serviceManager.registryCoordinator();
        address actualAVSDirectory = serviceManager.avsDirectory();

        require(
            actualRegistryCoordinator == expectedRegistryCoordinator,
            "EigenDAServiceManager has incorrect registryCoordinator address"
        );
        require(actualAVSDirectory == expectedAVSDirectory, "EigenDAServiceManager has incorrect avsDirectory address");
        // Note: rewardsCoordinator and stakeRegistry are internal variables without public getters,
        // so we can't verify them directly from the contract. We have to trust that the constructor
        // and deployment set them correctly.
    }
}
