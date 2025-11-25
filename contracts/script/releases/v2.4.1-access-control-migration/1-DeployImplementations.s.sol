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

contract DeployImplementations is EOADeployer {
    using Env for *;

    /// forgefmt: disable-next-item
    function _runAsEOA() internal override {
        // Get the directory to access deployed proxies
        EigenDADirectory directory = Env.proxy.directory();
        
        // Get EigenLayer contract addresses from environment
        // These should be set based on the network being deployed to
        address avsDirectory = vm.envOr("AVS_DIRECTORY", address(0));
        address rewardsCoordinator = vm.envOr("REWARDS_COORDINATOR", address(0));
        
        // If not set, try to read from a known EigenLayer deployment or config
        if (avsDirectory == address(0)) {
            // For mainnet: 0x135DDa560e946695d6f155dACaFC6f1F25C1F5AF
            // For holesky: 0x055733000064333CaDDbC92763c58BF0192fFeBf
            avsDirectory = vm.envOr("AVS_DIRECTORY", address(0));
        }
        if (rewardsCoordinator == address(0)) {
            rewardsCoordinator = vm.envOr("REWARDS_COORDINATOR", address(0));
        }

        vm.startBroadcast();

        // Deploy new ServiceManager implementation
        // Constructor requires references to EigenLayer and EigenDA contracts
        require(avsDirectory != address(0), "AVS_DIRECTORY not set");
        require(rewardsCoordinator != address(0), "REWARDS_COORDINATOR not set");
        
        deployImpl({
            name: type(EigenDAServiceManager).name,
            deployedTo: address(
                new EigenDAServiceManager(
                    IAVSDirectory(avsDirectory),
                    IRewardsCoordinator(rewardsCoordinator),
                    IRegistryCoordinator(address(Env.proxy.registryCoordinator())),
                    IStakeRegistry(address(Env.proxy.stakeRegistry())),
                    IEigenDAThresholdRegistry(address(Env.proxy.thresholdRegistry())),
                    IEigenDARelayRegistry(address(Env.proxy.relayRegistry())),
                    IPaymentVault(address(Env.proxy.paymentVault())),
                    IEigenDADisperserRegistry(address(Env.proxy.disperserRegistry()))
                )
            )
        });

        // Deploy new RegistryCoordinator implementation
        deployImpl({
            name: type(EigenDARegistryCoordinator).name,
            deployedTo: address(new EigenDARegistryCoordinator(address(directory)))
        });

        // Deploy new ThresholdRegistry implementation
        deployImpl({
            name: type(EigenDAThresholdRegistry).name,
            deployedTo: address(new EigenDAThresholdRegistry())
        });

        // Deploy new RelayRegistry implementation
        deployImpl({
            name: type(EigenDARelayRegistry).name,
            deployedTo: address(new EigenDARelayRegistry())
        });

        // Deploy new DisperserRegistry implementation
        deployImpl({
            name: type(EigenDADisperserRegistry).name,
            deployedTo: address(new EigenDADisperserRegistry())
        });

        // Deploy new PaymentVault implementation
        deployImpl({
            name: type(PaymentVault).name,
            deployedTo: address(new PaymentVault())
        });

        // Deploy new IndexRegistry implementation
        deployImpl({
            name: type(IndexRegistry).name,
            deployedTo: address(
                new IndexRegistry(IRegistryCoordinator(address(Env.proxy.registryCoordinator())))
            )
        });

        // Deploy new StakeRegistry implementation
        // Get DelegationManager from existing StakeRegistry
        StakeRegistry currentStakeRegistry = Env.proxy.stakeRegistry();
        address delegationManager = address(currentStakeRegistry.delegation());
        
        deployImpl({
            name: type(StakeRegistry).name,
            deployedTo: address(
                new StakeRegistry(
                    IRegistryCoordinator(address(Env.proxy.registryCoordinator())),
                    IDelegationManager(delegationManager)
                )
            )
        });

        // Deploy new BLSApkRegistry implementation
        deployImpl({
            name: type(BLSApkRegistry).name,
            deployedTo: address(
                new BLSApkRegistry(IRegistryCoordinator(address(Env.proxy.registryCoordinator())))
            )
        });

        // Deploy new SocketRegistry implementation
        deployImpl({
            name: type(SocketRegistry).name,
            deployedTo: address(
                new SocketRegistry(IRegistryCoordinator(address(Env.proxy.registryCoordinator())))
            )
        });

        // Deploy new EjectionManager implementation
        // Get parameters from existing ejection manager or environment
        address depositToken = vm.envOr("EJECTION_DEPOSIT_TOKEN", address(0));
        uint256 depositBaseFeeMultiplier = vm.envOr("EJECTION_DEPOSIT_BASE_FEE_MULTIPLIER", uint256(1));
        uint256 estimatedGasWithoutSig = vm.envOr("EJECTION_ESTIMATED_GAS_WITHOUT_SIG", uint256(100000));
        uint256 estimatedGasWithSig = vm.envOr("EJECTION_ESTIMATED_GAS_WITH_SIG", uint256(150000));
        
        deployImpl({
            name: type(EigenDAEjectionManager).name,
            deployedTo: address(
                new EigenDAEjectionManager(
                    depositToken,
                    depositBaseFeeMultiplier,
                    address(directory),
                    estimatedGasWithoutSig,
                    estimatedGasWithSig
                )
            )
        });

        // // Deploy new CertVerifier implementation
        // deployImpl({
        //     name: type(EigenDACertVerifier).name,
        //     deployedTo: address(
        //         new EigenDACertVerifier(
        //             IEigenDAThresholdRegistry(address(Env.proxy.thresholdRegistry())),
        //             IEigenDASignatureVerifier(address(Env.proxy.stakeRegistry())),
        //             new IEigenDACertVerifier.SecurityThresholds[](0), // Empty array, configured during initialization
        //             new uint8[](0) // Empty array, configured during initialization
        //         )
        //     )
        // });

        // Deploy new CertVerifierRouter implementation
        deployImpl({
            name: type(EigenDACertVerifierRouter).name,
            deployedTo: address(new EigenDACertVerifierRouter())
        });

        vm.stopBroadcast();
    }

    function testScript() public virtual {
        // Deploy the new implementations
        runAsEOA();

        // Validate implementations were deployed correctly
        _validateNewImplAddresses();
        _validateImplConstructors();
    }

    /// @dev Validate that new implementation addresses are non-zero and different from proxies
    function _validateNewImplAddresses() internal view {
        address serviceManagerImpl = address(Env.impl.serviceManager());
        address registryCoordinatorImpl = address(Env.impl.registryCoordinator());

        assertTrue(serviceManagerImpl != address(0), "ServiceManager implementation should be deployed");
        assertTrue(registryCoordinatorImpl != address(0), "RegistryCoordinator implementation should be deployed");

        // Ensure implementations are different from proxy addresses
        assertTrue(
            serviceManagerImpl != address(Env.proxy.serviceManager()),
            "ServiceManager implementation should differ from proxy"
        );
        assertTrue(
            registryCoordinatorImpl != address(Env.proxy.registryCoordinator()),
            "RegistryCoordinator implementation should differ from proxy"
        );

        // Validate other implementations
        address thresholdRegistryImpl = address(Env.impl.thresholdRegistry());
        assertTrue(thresholdRegistryImpl != address(0), "ThresholdRegistry implementation should be deployed");
    }

    /// @dev Validate implementation constructor values
    function _validateImplConstructors() internal view {
        // Validate ServiceManager constructor arguments
        EigenDAServiceManager serviceManager = Env.impl.serviceManager();

        // Verify immutable constructor parameters are set correctly
        assertTrue(
            address(serviceManager.registryCoordinator()) == address(Env.proxy.registryCoordinator()),
            "ServiceManager registryCoordinator should match"
        );
        assertTrue(
            address(serviceManager.stakeRegistry()) == address(Env.proxy.stakeRegistry()),
            "ServiceManager stakeRegistry should match"
        );

        // Validate RegistryCoordinator constructor arguments
        EigenDARegistryCoordinator registryCoordinator = Env.impl.registryCoordinator();
        EigenDADirectory directory = Env.proxy.directory();

        assertTrue(
            address(registryCoordinator.directory()) == address(directory), "RegistryCoordinator directory should match"
        );

        // Validate StakeRegistry constructor
        StakeRegistry stakeRegistry = Env.impl.stakeRegistry();
        assertTrue(
            address(stakeRegistry.registryCoordinator()) == address(Env.proxy.registryCoordinator()),
            "StakeRegistry registryCoordinator should match"
        );
    }
}
