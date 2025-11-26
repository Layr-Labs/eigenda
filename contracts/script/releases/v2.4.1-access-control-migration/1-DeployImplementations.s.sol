// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import {EOADeployer} from "zeus-templates/templates/EOADeployer.sol";
import "../Env.sol";

// Import EigenDA core contracts
import {EigenDAServiceManager} from "src/core/EigenDAServiceManager.sol";
import {EigenDARegistryCoordinator} from "src/core/EigenDARegistryCoordinator.sol";
import {EigenDAThresholdRegistry} from "src/core/EigenDAThresholdRegistry.sol";
import {EigenDARelayRegistry} from "src/core/EigenDARelayRegistry.sol";
import {EigenDADisperserRegistry} from "src/core/EigenDADisperserRegistry.sol";
import {PaymentVault} from "src/core/PaymentVault.sol";
import {EigenDADirectory} from "src/core/EigenDADirectory.sol";
import {EigenDAAccessControl} from "src/core/EigenDAAccessControl.sol";

// Import middleware contracts
import {BLSApkRegistry} from "lib/eigenlayer-middleware/src/BLSApkRegistry.sol";
import {IndexRegistry} from "lib/eigenlayer-middleware/src/IndexRegistry.sol";
import {StakeRegistry} from "lib/eigenlayer-middleware/src/StakeRegistry.sol";
import {SocketRegistry} from "lib/eigenlayer-middleware/src/SocketRegistry.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {
    PauserRegistry
} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/permissions/PauserRegistry.sol";

// Import periphery contracts
import {EigenDAEjectionManager} from "src/periphery/ejection/EigenDAEjectionManager.sol";

// Import certificate verification
import {EigenDACertVerifier} from "src/integrations/cert/EigenDACertVerifier.sol";
import {EigenDACertVerifierRouter} from "src/integrations/cert/router/EigenDACertVerifierRouter.sol";
import {EigenDATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";

// Import interfaces
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry} from "lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDARelayRegistry} from "src/core/interfaces/IEigenDARelayRegistry.sol";
import {IPaymentVault} from "src/core/interfaces/IPaymentVault.sol";
import {IEigenDADisperserRegistry} from "src/core/interfaces/IEigenDADisperserRegistry.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";
import {
    IAVSDirectory
} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IAVSDirectory.sol";
import {
    IRewardsCoordinator
} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IRewardsCoordinator.sol";
import {
    IDelegationManager
} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IDelegationManager.sol";

// TODO: Fetch CertVerifier and EjectionManager constructor parameters.
// TODO: Add DelegationManager to zeus.
// TODO: Figure out what initEigenDASignatureVerifier should be.
// TODO: Vender Pausable.
// TODO: Directory updates.

contract DeployImplementations is EOADeployer {
    using Env for *;

    /// -----------------------------------------------------------------------
    /// 1) Deploy new implementations as EOA
    /// -----------------------------------------------------------------------

    /// forgefmt: disable-next-item
    function _runAsEOA() internal virtual override {
        /// -----------------------------------------------------------------------
        /// Constructor parameters
        /// -----------------------------------------------------------------------

        // CertVerifier.
        EigenDATypesV1.SecurityThresholds memory initSecurityThresholds = EigenDATypesV1.SecurityThresholds({
            confirmationThreshold: 100, // 100% confirmation
            adversaryThreshold: 33 // 33% adversary
        });
        bytes memory initQuorumNumbersRequired = hex"00";
        
        // // EjectionManager.
        // uint256 depositBaseFeeMultiplier;
        // uint256 estimatedGasUsedWithoutSig;
        // uint256 estimatedGasUsedWithSig;

        // PauserRegistry.
        address[] memory initPausers = new address[](1);
        initPausers[0] = Env.impl.owner();

        /// -----------------------------------------------------------------------
        /// WARNING: NETWORK BROADCAST BEGINS HERE!
        /// -----------------------------------------------------------------------

        vm.startBroadcast();
        
        // Deploy new AccessControl implementation.
        deployImpl({
            name: "AccessControl", 
            deployedTo: address(new EigenDAAccessControl(Env.impl.owner()))
        });
        // Deploy new BLSApkRegistry implementation.
        deployImpl({
            name: "BlsApkRegistry", 
            deployedTo: address(new BLSApkRegistry(Env.proxy.registryCoordinator()))
        });
        // Deploy new CertVerifierRouter implementation.
        deployImpl({
            name: "CertVerifierRouter", 
            deployedTo: address(new EigenDACertVerifierRouter())
        });
        // Deploy new CertVerifier implementation.
        deployImpl({ // TODO: Likely needs removed + from Directory too.
            name: "CertVerifier",
            deployedTo: address(
                new EigenDACertVerifier({
                    initEigenDAThresholdRegistry: IEigenDAThresholdRegistry(address(Env.proxy.thresholdRegistry())),
                    initEigenDASignatureVerifier: IEigenDASignatureVerifier(address(0)), // TODO
                    initSecurityThresholds: initSecurityThresholds,
                    initQuorumNumbersRequired: initQuorumNumbersRequired
                })
            )
        });
        // Deploy new Directory implementation.
        deployImpl({
            name: "Directory", 
            deployedTo: address(new EigenDADirectory())
        });
        // Deploy new DisperserRegistry implementation.
        deployImpl({
            name: "DisperserRegistry", 
            deployedTo: address(new EigenDADisperserRegistry())
        });
        // // Deploy new EjectionManager implementation.
        // deployImpl({
        //     name: "EjectionManager",
        //     deployedTo: address(
        //         new EigenDAEjectionManager({
        //             depositToken_: Env.proxy.ejectionManager().getDepositToken(),
        //             depositBaseFeeMultiplier_: depositBaseFeeMultiplier,
        //             addressDirectory_: address(Env.proxy.directory()),
        //             estimatedGasUsedWithoutSig_: estimatedGasUsedWithoutSig,
        //             estimatedGasUsedWithSig_: estimatedGasUsedWithSig
        //         })
        //     )
        // });
        // Deploy new IndexRegistry implementation.
        deployImpl({
            name: type(IndexRegistry).name, 
            deployedTo: address(new IndexRegistry(Env.proxy.registryCoordinator()))
        });
        // Deploy new OperatorStateRetriever implementation.
        deployImpl({
            name: type(OperatorStateRetriever).name, 
            deployedTo: address(new OperatorStateRetriever())}
        );
        // Deploy new PauserRegistry implementation.
        deployImpl({
            name: type(PauserRegistry).name, 
            deployedTo: address(new PauserRegistry({_pausers: initPausers, _unpauser: Env.impl.owner()}))
        });
        // Deploy new PaymentVault implementation.
        deployImpl({
            name: type(PaymentVault).name, 
            deployedTo: address(new PaymentVault())
        });
        // Deploy new RegistryCoordinator implementation.
        deployImpl({
            name: "RegistryCoordinator",
            deployedTo: address(new EigenDARegistryCoordinator({_directory: address(Env.proxy.directory())}))
        });
        // Deploy new RelayRegistry implementation.
        deployImpl({
            name: "RelayRegistry", 
            deployedTo: address(new EigenDARelayRegistry())
        });
        // Deploy new ServiceManager implementation.
        deployImpl({
            name: "ServiceManager",
            deployedTo: address(
                new EigenDAServiceManager({
                    __avsDirectory: Env.proxy.avsDirectory(),
                    __rewardsCoordinator: Env.proxy.rewardsCoordinator(),
                    __registryCoordinator: Env.proxy.registryCoordinator(),
                    __stakeRegistry: Env.proxy.stakeRegistry(),
                    __eigenDAThresholdRegistry: Env.proxy.thresholdRegistry(),
                    __eigenDARelayRegistry: Env.proxy.relayRegistry(),
                    __paymentVault: Env.proxy.paymentVault(),
                    __eigenDADisperserRegistry: Env.proxy.disperserRegistry()
                })
            )
        });
        // Deploy new SocketRegistry implementation.
        deployImpl({
            name: type(SocketRegistry).name, 
            deployedTo: address(new SocketRegistry({_registryCoordinator: Env.proxy.registryCoordinator()}))
        });
        // Deploy new StakeRegistry implementation.
        deployImpl({
            name: type(StakeRegistry).name,
            deployedTo: address(
                new StakeRegistry({
                    _registryCoordinator: Env.proxy.registryCoordinator(),
                    _delegationManager: IDelegationManager(Env.proxy.stakeRegistry().delegation())
                })
            )
        });
        // Deploy new ThresholdRegistry implementation.
        deployImpl({
            name: "ThresholdRegistry", 
            deployedTo: address(new EigenDAThresholdRegistry())
        });

        vm.stopBroadcast();
    }

    /// -----------------------------------------------------------------------
    /// 2) Post-deployment assertions
    /// -----------------------------------------------------------------------

    function testScript() public virtual {
        // Deploy new implementations as EOA.
        runAsEOA();
        // Hook for post-deployment assertions.
        _afterTestScript();
    }

    /// -----------------------------------------------------------------------
    /// Test hooks
    /// -----------------------------------------------------------------------

    function _afterTestScript() internal view {
        _testDeploymentAddresses();
        _testImplementationCode();
        _testRegistryImmutables();
        _testServiceManagerImmutables();
        _testOtherImmutables();
    }

    /// -----------------------------------------------------------------------
    /// Tests
    /// -----------------------------------------------------------------------

    /// @notice Verify all implementations deployed to non-zero addresses
    function _testDeploymentAddresses() internal view {
        assertTrue(address(Env.impl.blsApkRegistry()) != address(0), "BLSApkRegistry not deployed");
        assertTrue(address(Env.impl.certVerifierRouter()) != address(0), "CertVerifierRouter not deployed");
        assertTrue(address(Env.impl.certVerifier()) != address(0), "CertVerifier not deployed");
        assertTrue(address(Env.impl.directory()) != address(0), "Directory not deployed");
        assertTrue(address(Env.impl.disperserRegistry()) != address(0), "DisperserRegistry not deployed");
        assertTrue(address(Env.impl.ejectionManager()) != address(0), "EjectionManager not deployed");
        assertTrue(address(Env.impl.indexRegistry()) != address(0), "IndexRegistry not deployed");
        assertTrue(address(Env.impl.operatorStateRetriever()) != address(0), "OperatorStateRetriever not deployed");
        assertTrue(address(Env.impl.paymentVault()) != address(0), "PaymentVault not deployed");
        assertTrue(address(Env.impl.registryCoordinator()) != address(0), "RegistryCoordinator not deployed");
        assertTrue(address(Env.impl.relayRegistry()) != address(0), "RelayRegistry not deployed");
        assertTrue(address(Env.impl.serviceManager()) != address(0), "ServiceManager not deployed");
        assertTrue(address(Env.impl.socketRegistry()) != address(0), "SocketRegistry not deployed");
        assertTrue(address(Env.impl.stakeRegistry()) != address(0), "StakeRegistry not deployed");
        assertTrue(address(Env.impl.thresholdRegistry()) != address(0), "ThresholdRegistry not deployed");
    }

    /// @notice Verify implementations have bytecode deployed
    function _testImplementationCode() internal view {
        assertTrue(address(Env.impl.registryCoordinator()).code.length > 0, "RegistryCoordinator has no code");
        assertTrue(address(Env.impl.serviceManager()).code.length > 0, "ServiceManager has no code");
        assertTrue(address(Env.impl.paymentVault()).code.length > 0, "PaymentVault has no code");
        assertTrue(address(Env.impl.thresholdRegistry()).code.length > 0, "ThresholdRegistry has no code");
        assertTrue(address(Env.impl.directory()).code.length > 0, "Directory has no code");
    }

    /// @notice Verify immutable constructor parameters for registry implementations
    function _testRegistryImmutables() internal view {
        // Verify individual registry immutables
        assertEq(
            address(Env.impl.blsApkRegistry().registryCoordinator()),
            address(Env.proxy.registryCoordinator()),
            "BLSApkRegistry: incorrect registryCoordinator"
        );
        assertEq(
            address(Env.impl.indexRegistry().registryCoordinator()),
            address(Env.proxy.registryCoordinator()),
            "IndexRegistry: incorrect registryCoordinator"
        );
        assertEq(
            address(Env.impl.socketRegistry().registryCoordinator()),
            address(Env.proxy.registryCoordinator()),
            "SocketRegistry: incorrect registryCoordinator"
        );
        assertEq(
            address(Env.impl.stakeRegistry().registryCoordinator()),
            address(Env.proxy.registryCoordinator()),
            "StakeRegistry: incorrect registryCoordinator"
        );
        assertEq(
            address(Env.impl.stakeRegistry().delegation()),
            address(Env.proxy.stakeRegistry().delegation()),
            "StakeRegistry: incorrect delegation"
        );

        // Verify RegistryCoordinator references to other registries
        assertEq(
            address(Env.impl.registryCoordinator().stakeRegistry()),
            address(Env.proxy.stakeRegistry()),
            "RegistryCoordinator: incorrect stakeRegistry"
        );
        assertEq(
            address(Env.impl.registryCoordinator().blsApkRegistry()),
            address(Env.proxy.blsApkRegistry()),
            "RegistryCoordinator: incorrect blsApkRegistry"
        );
        assertEq(
            address(Env.impl.registryCoordinator().indexRegistry()),
            address(Env.proxy.indexRegistry()),
            "RegistryCoordinator: incorrect indexRegistry"
        );
        assertEq(
            address(Env.impl.registryCoordinator().socketRegistry()),
            address(Env.proxy.socketRegistry()),
            "RegistryCoordinator: incorrect socketRegistry"
        );
        assertEq(
            address(Env.impl.registryCoordinator().directory()),
            address(Env.proxy.directory()),
            "RegistryCoordinator: incorrect directory"
        );
    }

    /// @notice Verify ServiceManager immutable parameters
    function _testServiceManagerImmutables() internal view {
        assertEq(
            address(Env.impl.serviceManager().avsDirectory()),
            address(Env.proxy.avsDirectory()),
            "ServiceManager: incorrect avsDirectory"
        );
        assertEq(
            address(Env.impl.serviceManager().eigenDAThresholdRegistry()),
            address(Env.proxy.thresholdRegistry()),
            "ServiceManager: incorrect eigenDAThresholdRegistry"
        );
        assertEq(
            address(Env.impl.serviceManager().eigenDARelayRegistry()),
            address(Env.proxy.relayRegistry()),
            "ServiceManager: incorrect eigenDARelayRegistry"
        );
        assertEq(
            address(Env.impl.serviceManager().paymentVault()),
            address(Env.proxy.paymentVault()),
            "ServiceManager: incorrect paymentVault"
        );
        assertEq(
            address(Env.impl.serviceManager().eigenDADisperserRegistry()),
            address(Env.proxy.disperserRegistry()),
            "ServiceManager: incorrect eigenDADisperserRegistry"
        );
    }

    /// @notice Verify other contract immutables
    function _testOtherImmutables() internal view {
        // CertVerifier
        assertEq(
            address(Env.impl.certVerifier().eigenDAThresholdRegistry()),
            address(Env.proxy.thresholdRegistry()),
            "CertVerifier: incorrect thresholdRegistry"
        );

        // PauserRegistry
        assertEq(Env.impl.pauserRegistry().unpauser(), Env.impl.owner(), "PauserRegistry: incorrect unpauser");
        assertTrue(Env.impl.pauserRegistry().isPauser(Env.impl.owner()), "PauserRegistry: owner not set as pauser");

        // // EjectionManager
        // assertEq(
        //     Env.impl.ejectionManager().getDepositToken(),
        //     Env.proxy.ejectionManager().getDepositToken(),
        //     "EjectionManager: incorrect depositToken"
        // );
    }
}
