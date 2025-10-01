// SPDX-License-Identifier: BUSL-1.1

pragma solidity =0.8.12;

import {EmptyContract} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/test/mocks/EmptyContract.sol";
import {ProxyAdmin, TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";

import {IDelegationManager} from
    "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IDelegationManager.sol";
import {ISocketRegistry, SocketRegistry} from "lib/eigenlayer-middleware/src/SocketRegistry.sol";
import {IIndexRegistry} from "lib/eigenlayer-middleware/src/interfaces/IIndexRegistry.sol";
import {IndexRegistry} from "lib/eigenlayer-middleware/src/IndexRegistry.sol";
import {IStakeRegistry, StakeRegistry} from "lib/eigenlayer-middleware/src/StakeRegistry.sol";
import {IBLSApkRegistry} from "lib/eigenlayer-middleware/src/interfaces/IBLSApkRegistry.sol";
import {BLSApkRegistry} from "lib/eigenlayer-middleware/src/BLSApkRegistry.sol";
import {EigenDARegistryCoordinator, IRegistryCoordinator} from "src/core/EigenDARegistryCoordinator.sol";
import {IEigenDAThresholdRegistry, EigenDAThresholdRegistry} from "src/core/EigenDAThresholdRegistry.sol";
import {IEigenDARelayRegistry, EigenDARelayRegistry} from "src/core/EigenDARelayRegistry.sol";
import {PaymentVault} from "src/core/PaymentVault.sol";
import {IPaymentVault} from "src/core/interfaces/IPaymentVault.sol";
import {IEigenDADisperserRegistry, EigenDADisperserRegistry} from "src/core/EigenDADisperserRegistry.sol";
import {EigenDAServiceManager, IServiceManager} from "src/core/EigenDAServiceManager.sol";
import {IAVSDirectory} from
    "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IAVSDirectory.sol";
import {IRewardsCoordinator} from
    "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IRewardsCoordinator.sol";
import {
    IPauserRegistry,
    PauserRegistry
} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/permissions/PauserRegistry.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";
import {EjectionManager} from "lib/eigenlayer-middleware/src/EjectionManager.sol";
import {IServiceManager} from "lib/eigenlayer-middleware/src/interfaces/IServiceManager.sol";
import {EigenDATypesV2 as DATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {EigenDACertVerifier} from "src/integrations/cert/EigenDACertVerifier.sol";
import {EigenDACertVerifierRouter} from "src/integrations/cert/router/EigenDACertVerifierRouter.sol";

import {MockStakeRegistry} from "test/mock/MockStakeRegistry.sol";
import {MockRegistryCoordinator} from "test/mock/MockRegistryCoordinator.sol";

import {InitParamsLib} from "script/deploy/eigenda/DeployEigenDAConfig.sol";

import {AddressDirectoryConstants} from "src/core/libraries/v3/address-directory/AddressDirectoryConstants.sol";
import {AccessControlConstants} from "src/core/libraries/v3/access-control/AccessControlConstants.sol";
import {EigenDADirectory} from "src/core/EigenDADirectory.sol";
import {EigenDAAccessControl} from "src/core/EigenDAAccessControl.sol";

import {Script} from "forge-std/Script.sol";
import {console2} from "forge-std/console2.sol";

/// @notice This script deploys EigenDA contracts and should eventually replace the other deployment scripts,
///         which cannot currently be removed due to CI depending on them.
contract DeployEigenDA is Script {
    using InitParamsLib for string;

    EigenDADirectory directory;
    address proxyAdmin;

    mapping(string => address) impl; // Implementation addresses of the deployed contracts.
    mapping(string => bool) upgraded; // Whether the deployment of a contract is upgraded to its final implementation. Should beTrue if the contract is not a proxy

    string cfg;

    string constant EMPTY_CONTRACT = "EMPTY_CONTRACT";
    string constant MOCK_STAKE_REGISTRY = "MOCK_STAKE_REGISTRY";
    string constant MOCK_REGISTRY_COORDINATOR = "MOCK_REGISTRY_COORDINATOR";

    function initConfig() internal virtual {
        cfg = vm.readFile(vm.envString("DEPLOY_CONFIG_PATH"));
    }

    function run() public virtual {
        initConfig();
        vm.startBroadcast();

        // DEPLOY PROXY ADMIN
        proxyAdmin = address(new ProxyAdmin());

        /// These steps are done after the main deployment because not all eigenDA contracts use these contracts yet.
        /// So these contracts can be considered to live somewhere in the "periphery" of the eigenDA system for now.

        /// DEPLOY EIGENDA DIRECTORY AND ACCESS CONTROL

        directory = EigenDADirectory(
            address(
                new TransparentUpgradeableProxy(
                    address(new EigenDADirectory()),
                    proxyAdmin,
                    abi.encodeCall(EigenDADirectory.initialize, (address(new EigenDAAccessControl(msg.sender))))
                )
            )
        );

        console2.log("DIRECTORY:", address(directory));

        // DEPLOY MOCK IMPLEMENTATION
        impl[EMPTY_CONTRACT] = address(new EmptyContract());

        // DEPLOY PAUSER
        directory.addAddress(
            AddressDirectoryConstants.PAUSER_REGISTRY_NAME, address(new PauserRegistry(cfg.pausers(), cfg.unpauser()))
        );

        // Registry coordinator requires these contracts as constructor arguments for implementation deployment
        // However, these contracts also require knowing the registry coordinator address
        // before they can be deployed, so we deploy them as inert proxies first.
        // INDEX REGISTRY
        // STAKE REGISTRY
        // SOCKET REGISTRY
        // BLS APK REGISTRY
        // SERVICE MANAGER
        directory.addAddress(
            AddressDirectoryConstants.INDEX_REGISTRY_NAME,
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], proxyAdmin, ""))
        );

        directory.addAddress(
            AddressDirectoryConstants.SOCKET_REGISTRY_NAME,
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], proxyAdmin, ""))
        );

        directory.addAddress(
            AddressDirectoryConstants.BLS_APK_REGISTRY_NAME,
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], proxyAdmin, ""))
        );
        impl[MOCK_STAKE_REGISTRY] = address(new MockStakeRegistry(IDelegationManager(cfg.delegationManager())));
        // The service manager implementation requires the stake registry to expose the delegation manager on construction.
        directory.addAddress(
            AddressDirectoryConstants.STAKE_REGISTRY_NAME,
            address(new TransparentUpgradeableProxy(impl[MOCK_STAKE_REGISTRY], proxyAdmin, ""))
        );
        // The service manager implementation requires the registry coordinator to expose the stake registry and bls APK registry on construction.
        // And this can only be done after the stake registry and bls APK registry proxies are known.
        impl[MOCK_REGISTRY_COORDINATOR] = address(
            new MockRegistryCoordinator(
                IStakeRegistry(directory.getAddress(AddressDirectoryConstants.STAKE_REGISTRY_NAME)),
                IBLSApkRegistry(directory.getAddress(AddressDirectoryConstants.BLS_APK_REGISTRY_NAME))
            )
        );
        directory.addAddress(
            AddressDirectoryConstants.REGISTRY_COORDINATOR_NAME,
            address(new TransparentUpgradeableProxy(impl[MOCK_REGISTRY_COORDINATOR], proxyAdmin, ""))
        );
        directory.addAddress(
            AddressDirectoryConstants.THRESHOLD_REGISTRY_NAME,
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], proxyAdmin, ""))
        );
        directory.addAddress(
            AddressDirectoryConstants.RELAY_REGISTRY_NAME,
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], proxyAdmin, ""))
        );
        directory.addAddress(
            AddressDirectoryConstants.DISPERSER_REGISTRY_NAME,
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], proxyAdmin, ""))
        );
        directory.addAddress(
            AddressDirectoryConstants.SERVICE_MANAGER_NAME,
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], proxyAdmin, ""))
        );
        directory.addAddress(
            AddressDirectoryConstants.PAYMENT_VAULT_NAME,
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], proxyAdmin, ""))
        );
        directory.addAddress(
            AddressDirectoryConstants.EJECTION_MANAGER_NAME,
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], proxyAdmin, ""))
        );

        impl[AddressDirectoryConstants.INDEX_REGISTRY_NAME] = address(
            new IndexRegistry(
                IRegistryCoordinator(directory.getAddress(AddressDirectoryConstants.REGISTRY_COORDINATOR_NAME))
            )
        );
        upgrade(AddressDirectoryConstants.INDEX_REGISTRY_NAME, "");

        impl[AddressDirectoryConstants.STAKE_REGISTRY_NAME] = address(
            new StakeRegistry(
                IRegistryCoordinator(directory.getAddress(AddressDirectoryConstants.REGISTRY_COORDINATOR_NAME)),
                IDelegationManager(cfg.delegationManager())
            )
        );
        upgrade(AddressDirectoryConstants.STAKE_REGISTRY_NAME, "");

        impl[AddressDirectoryConstants.SOCKET_REGISTRY_NAME] = address(
            new SocketRegistry(
                IRegistryCoordinator(directory.getAddress(AddressDirectoryConstants.REGISTRY_COORDINATOR_NAME))
            )
        );
        upgrade(AddressDirectoryConstants.SOCKET_REGISTRY_NAME, "");

        impl[AddressDirectoryConstants.BLS_APK_REGISTRY_NAME] = address(
            new BLSApkRegistry(
                IRegistryCoordinator(directory.getAddress(AddressDirectoryConstants.REGISTRY_COORDINATOR_NAME))
            )
        );
        upgrade(AddressDirectoryConstants.BLS_APK_REGISTRY_NAME, "");

        impl[AddressDirectoryConstants.REGISTRY_COORDINATOR_NAME] = address(
            new EigenDARegistryCoordinator(
                IServiceManager(directory.getAddress(AddressDirectoryConstants.SERVICE_MANAGER_NAME)),
                IStakeRegistry(directory.getAddress(AddressDirectoryConstants.STAKE_REGISTRY_NAME)),
                IBLSApkRegistry(directory.getAddress(AddressDirectoryConstants.BLS_APK_REGISTRY_NAME)),
                IIndexRegistry(directory.getAddress(AddressDirectoryConstants.INDEX_REGISTRY_NAME)),
                ISocketRegistry(directory.getAddress(AddressDirectoryConstants.SOCKET_REGISTRY_NAME))
            )
        );
        upgrade(
            AddressDirectoryConstants.REGISTRY_COORDINATOR_NAME,
            abi.encodeCall(
                EigenDARegistryCoordinator.initialize,
                (
                    cfg.initialOwner(),
                    directory.getAddress(AddressDirectoryConstants.EJECTION_MANAGER_NAME),
                    IPauserRegistry(directory.getAddress(AddressDirectoryConstants.PAUSER_REGISTRY_NAME)),
                    cfg.initialPausedStatus(),
                    cfg.operatorSetParams(),
                    cfg.minimumStakes(),
                    cfg.strategyParams()
                )
            )
        );

        impl[AddressDirectoryConstants.SERVICE_MANAGER_NAME] = address(
            new EigenDAServiceManager(
                IAVSDirectory(cfg.avsDirectory()),
                IRewardsCoordinator(cfg.rewardsCoordinator()),
                IRegistryCoordinator(directory.getAddress(AddressDirectoryConstants.REGISTRY_COORDINATOR_NAME)),
                IStakeRegistry(directory.getAddress(AddressDirectoryConstants.STAKE_REGISTRY_NAME)),
                IEigenDAThresholdRegistry(directory.getAddress(AddressDirectoryConstants.THRESHOLD_REGISTRY_NAME)),
                IEigenDARelayRegistry(directory.getAddress(AddressDirectoryConstants.RELAY_REGISTRY_NAME)),
                IPaymentVault(directory.getAddress(AddressDirectoryConstants.PAYMENT_VAULT_NAME)),
                IEigenDADisperserRegistry(directory.getAddress(AddressDirectoryConstants.DISPERSER_REGISTRY_NAME))
            )
        );
        upgrade(
            AddressDirectoryConstants.SERVICE_MANAGER_NAME,
            abi.encodeCall(
                EigenDAServiceManager.initialize,
                (
                    IPauserRegistry(directory.getAddress(AddressDirectoryConstants.PAUSER_REGISTRY_NAME)),
                    cfg.initialPausedStatus(),
                    cfg.initialOwner(),
                    cfg.batchConfirmers(),
                    cfg.rewardsInitiator()
                )
            )
        );

        impl[AddressDirectoryConstants.EJECTION_MANAGER_NAME] = address(
            new EjectionManager(
                IRegistryCoordinator(directory.getAddress(AddressDirectoryConstants.REGISTRY_COORDINATOR_NAME)),
                IStakeRegistry(directory.getAddress(AddressDirectoryConstants.STAKE_REGISTRY_NAME))
            )
        );
        upgrade(
            AddressDirectoryConstants.EJECTION_MANAGER_NAME,
            abi.encodeCall(EjectionManager.initialize, (cfg.initialOwner(), cfg.ejectors(), cfg.quorumEjectionParams()))
        );

        impl[AddressDirectoryConstants.THRESHOLD_REGISTRY_NAME] = address(new EigenDAThresholdRegistry());
        upgrade(
            AddressDirectoryConstants.THRESHOLD_REGISTRY_NAME,
            abi.encodeCall(
                EigenDAThresholdRegistry.initialize,
                (
                    cfg.initialOwner(),
                    cfg.quorumAdversaryThresholdPercentages(),
                    cfg.quorumConfirmationThresholdPercentages(),
                    cfg.quorumNumbersRequired(),
                    cfg.versionedBlobParams()
                )
            )
        );

        impl[AddressDirectoryConstants.RELAY_REGISTRY_NAME] = address(new EigenDARelayRegistry());
        upgrade(
            AddressDirectoryConstants.RELAY_REGISTRY_NAME, abi.encodeCall(EigenDARelayRegistry.initialize, (msg.sender))
        );

        impl[AddressDirectoryConstants.DISPERSER_REGISTRY_NAME] = address(new EigenDADisperserRegistry());
        upgrade(
            AddressDirectoryConstants.DISPERSER_REGISTRY_NAME,
            abi.encodeCall(EigenDADisperserRegistry.initialize, (msg.sender))
        );

        impl[AddressDirectoryConstants.PAYMENT_VAULT_NAME] = address(new PaymentVault());
        upgrade(
            AddressDirectoryConstants.PAYMENT_VAULT_NAME,
            abi.encodeCall(
                PaymentVault.initialize,
                (
                    cfg.initialOwner(),
                    cfg.minNumSymbols(),
                    cfg.pricePerSymbol(),
                    cfg.priceUpdateCooldown(),
                    cfg.globalSymbolsPerPeriod(),
                    cfg.reservationPeriodInterval(),
                    cfg.globalRatePeriodInterval()
                )
            )
        );
        directory.addAddress(
            AddressDirectoryConstants.OPERATOR_STATE_RETRIEVER_NAME, address(new OperatorStateRetriever())
        );

        directory.addAddress(
            AddressDirectoryConstants.CERT_VERIFIER_NAME,
            address(
                new EigenDACertVerifier(
                    IEigenDAThresholdRegistry(directory.getAddress(AddressDirectoryConstants.THRESHOLD_REGISTRY_NAME)),
                    IEigenDASignatureVerifier(directory.getAddress(AddressDirectoryConstants.STAKE_REGISTRY_NAME)),
                    cfg.certVerifierSecurityThresholds(),
                    cfg.certVerifierQuorumNumbersRequired()
                )
            )
        );

        address routerImpl = address(new EigenDACertVerifierRouter());
        address[] memory certVerifiers = new address[](1);

        certVerifiers[0] = directory.getAddress(AddressDirectoryConstants.CERT_VERIFIER_NAME);

        directory.addAddress(
            AddressDirectoryConstants.CERT_VERIFIER_ROUTER_NAME,
            address(
                new TransparentUpgradeableProxy(
                    routerImpl,
                    proxyAdmin,
                    abi.encodeWithSelector(
                        EigenDACertVerifierRouter.initialize.selector,
                        cfg.initialOwner(),
                        new uint32[](1), // equivalent to [0]
                        certVerifiers
                    )
                )
            )
        );

        ProxyAdmin(proxyAdmin).transferOwnership(cfg.initialOwner());
        EigenDAAccessControl accessControl =
            EigenDAAccessControl(directory.getAddress(AddressDirectoryConstants.ACCESS_CONTROL_NAME));

        for (uint256 i; i < cfg.dispersers().length; i++) {
            IEigenDADisperserRegistry(directory.getAddress(AddressDirectoryConstants.DISPERSER_REGISTRY_NAME))
                .setDisperserInfo(uint32(i), DATypesV2.DisperserInfo(cfg.dispersers()[i]));
        }

        for (uint256 i; i < cfg.relayInfos().length; i++) {
            IEigenDARelayRegistry(directory.getAddress(AddressDirectoryConstants.RELAY_REGISTRY_NAME)).addRelayInfo(
                cfg.relayInfos()[i]
            );
        }

        if (msg.sender != cfg.initialOwner()) {
            accessControl.grantRole(accessControl.DEFAULT_ADMIN_ROLE(), cfg.initialOwner());
            accessControl.grantRole(AccessControlConstants.OWNER_ROLE, cfg.initialOwner());
            accessControl.revokeRole(AccessControlConstants.OWNER_ROLE, msg.sender);
            accessControl.revokeRole(accessControl.DEFAULT_ADMIN_ROLE(), msg.sender);
            EigenDADisperserRegistry(directory.getAddress(AddressDirectoryConstants.DISPERSER_REGISTRY_NAME))
                .transferOwnership(cfg.initialOwner());
            EigenDARelayRegistry(directory.getAddress(AddressDirectoryConstants.RELAY_REGISTRY_NAME)).transferOwnership(
                cfg.initialOwner()
            );
        }

        vm.stopBroadcast();
    }

    function upgrade(string memory contractName, bytes memory initData) internal {
        address implementation = impl[contractName];
        TransparentUpgradeableProxy proxy = TransparentUpgradeableProxy(payable(directory.getAddress(contractName)));

        ProxyAdmin(proxyAdmin).upgrade(proxy, implementation);
        if (initData.length > 0) {
            ProxyAdmin(proxyAdmin).upgradeAndCall(proxy, implementation, initData);
        }
        upgraded[contractName] = true;
    }
}
