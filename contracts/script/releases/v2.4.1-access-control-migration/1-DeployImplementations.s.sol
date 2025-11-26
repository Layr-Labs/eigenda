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

// TODO: Fetch remaining constructor/initialize parameters from existing contracts on chain.
// TODO: Ensure all relevant contracts are included in the upgrade.
// TODO: Add post deployment assertions.

contract DeployImplementations is EOADeployer {
    using Env for *;

    /// -----------------------------------------------------------------------
    /// 1) Deploy new implementations as EOA
    /// -----------------------------------------------------------------------

    function _runAsEOA() internal override {
        // Deploy new ServiceManager implementation.
        deployImpl({
            name: "ServiceManager",
            deployedTo: address(
                new EigenDAServiceManager(
                    IAVSDirectory(Env.proxy.avsDirectory()),
                    IRewardsCoordinator(Env.proxy.rewardsCoordinator()),
                    IRegistryCoordinator(Env.proxy.registryCoordinator()),
                    IStakeRegistry(Env.proxy.stakeRegistry()),
                    IEigenDAThresholdRegistry(Env.proxy.thresholdRegistry()),
                    IEigenDARelayRegistry(Env.proxy.relayRegistry()),
                    IPaymentVault(Env.proxy.paymentVault()),
                    IEigenDADisperserRegistry(Env.proxy.disperserRegistry())
                )
            )
        });

        // Deploy new RegistryCoordinator implementation.
        deployImpl({
            name: "RegistryCoordinator",
            deployedTo: address(new EigenDARegistryCoordinator(address(Env.proxy.directory())))
        });

        // Deploy new ThresholdRegistry implementation.
        deployImpl({name: "ThresholdRegistry", deployedTo: address(new EigenDAThresholdRegistry())});
        // Deploy new RelayRegistry implementation.
        deployImpl({name: "RelayRegistry", deployedTo: address(new EigenDARelayRegistry())});
        // Deploy new DisperserRegistry implementation.
        deployImpl({name: "DisperserRegistry", deployedTo: address(new EigenDADisperserRegistry())});
        // Deploy new PaymentVault implementation.
        deployImpl({name: type(PaymentVault).name, deployedTo: address(new PaymentVault())});
        // Deploy new IndexRegistry implementation.
        deployImpl({
            name: type(IndexRegistry).name, deployedTo: address(new IndexRegistry(Env.proxy.registryCoordinator()))
        });

        // Deploy new StakeRegistry implementation.
        deployImpl({
            name: type(StakeRegistry).name,
            deployedTo: address(
                new StakeRegistry(
                    IRegistryCoordinator(Env.proxy.registryCoordinator()),
                    IDelegationManager(Env.proxy.stakeRegistry().delegation())
                )
            )
        });

        // Deploy new BLSApkRegistry implementation.
        deployImpl({
            name: type(BLSApkRegistry).name, deployedTo: address(new BLSApkRegistry(Env.proxy.registryCoordinator()))
        });
        // Deploy new SocketRegistry implementation.
        deployImpl({
            name: type(SocketRegistry).name, deployedTo: address(new SocketRegistry(Env.proxy.registryCoordinator()))
        });

        // TODO: Get parameters from existing ejection manager.
        uint256 depositBaseFeeMultiplier = vm.envOr("EJECTION_DEPOSIT_BASE_FEE_MULTIPLIER", uint256(1));
        uint256 estimatedGasWithoutSig = vm.envOr("EJECTION_ESTIMATED_GAS_WITHOUT_SIG", uint256(100_000));
        uint256 estimatedGasWithSig = vm.envOr("EJECTION_ESTIMATED_GAS_WITH_SIG", uint256(150_000));

        // Deploy new EjectionManager implementation.
        deployImpl({
            name: "EjectionManager",
            deployedTo: address(
                new EigenDAEjectionManager(
                    Env.proxy.ejectionManager().getDepositToken(),
                    depositBaseFeeMultiplier,
                    address(Env.proxy.directory()),
                    estimatedGasWithoutSig,
                    estimatedGasWithSig
                )
            )
        });

        // TODO: Get values from team or onchain state.
        EigenDATypesV1.SecurityThresholds memory defaultThresholds = EigenDATypesV1.SecurityThresholds({
            confirmationThreshold: 100, // 100% confirmation
            adversaryThreshold: 33 // 33% adversary
        });

        deployImpl({
            name: "CertVerifier",
            deployedTo: address(
                new EigenDACertVerifier(
                    IEigenDAThresholdRegistry(address(Env.proxy.thresholdRegistry())),
                    IEigenDASignatureVerifier(address(Env.proxy.stakeRegistry())),
                    defaultThresholds,
                    hex"00" // Default quorum numbers required
                )
            )
        });

        // Deploy new CertVerifierRouter implementation.
        deployImpl({name: "CertVerifierRouter", deployedTo: address(new EigenDACertVerifierRouter())});

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
