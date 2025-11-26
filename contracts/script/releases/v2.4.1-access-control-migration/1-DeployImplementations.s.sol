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

// NOTE: The names in deployImpl must match the names in in zeus,
// will likely have to correct all contracts with "EigenDA" prefix.

// TODO: Fetch CertVerifier and EjectionManager constructor parameters.
// TODO: Add post deployment assertions.

contract DeployImplementations is EOADeployer {
    using Env for *;

    /// -----------------------------------------------------------------------
    /// 1) Deploy new implementations as EOA
    /// -----------------------------------------------------------------------

    /// forgefmt: disable-next-item
    function _runAsEOA() internal override {
        // CertVerifier constructor parameters.
        EigenDATypesV1.SecurityThresholds memory initSecurityThresholds = EigenDATypesV1.SecurityThresholds({
            confirmationThreshold: 100, // 100% confirmation
            adversaryThreshold: 33 // 33% adversary
        });
        bytes memory initQuorumNumbersRequired = hex"00";
        
        // EjectionManager constructor parameters.
        uint256 depositBaseFeeMultiplier;
        uint256 estimatedGasUsedWithoutSig;
        uint256 estimatedGasUsedWithSig;

        // PauserRegistry constructor parameters.
        address[] memory initPausers = new address[](1);
        initPausers[0] = Env.impl.owner();

        vm.startBroadcast();

        // Deploy new BLSApkRegistry implementation.
        deployImpl({
            name: type(BLSApkRegistry).name, 
            deployedTo: address(new BLSApkRegistry(Env.proxy.registryCoordinator()))
        });
        // Deploy new CertVerifierRouter implementation.
        deployImpl({
            name: "CertVerifierRouter", 
            deployedTo: address(new EigenDACertVerifierRouter())
        });
        // Deploy new CertVerifier implementation.
        deployImpl({
            name: "CertVerifier",
            deployedTo: address(
                new EigenDACertVerifier({
                    initEigenDAThresholdRegistry: IEigenDAThresholdRegistry(address(Env.proxy.thresholdRegistry())),
                    initEigenDASignatureVerifier: IEigenDASignatureVerifier(address(0)), // TODO: Figure out what contract this should be.
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

        // Deploy new EjectionManager implementation.
        deployImpl({
            name: "EjectionManager",
            deployedTo: address(
                new EigenDAEjectionManager({
                    depositToken_: Env.proxy.ejectionManager().getDepositToken(),
                    depositBaseFeeMultiplier_: depositBaseFeeMultiplier,
                    addressDirectory_: address(Env.proxy.directory()),
                    estimatedGasUsedWithoutSig_: estimatedGasUsedWithoutSig,
                    estimatedGasUsedWithSig_: estimatedGasUsedWithSig
                })
            )
        });
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
                    _delegationManager: IDelegationManager(Env.proxy.stakeRegistry().delegation()) // TODO: Add DM to zeus.
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
        // Hook for pre-test setup.
        _beforeTestScript();
        // Deploy new implementations as EOA.
        runAsEOA();
        // Hook for post-deployment assertions.
        _afterTestScript();
    }

    /// -----------------------------------------------------------------------
    /// Test hooks
    /// -----------------------------------------------------------------------

    function _beforeTestScript() internal view {}

    function _afterTestScript() internal view {}
}
