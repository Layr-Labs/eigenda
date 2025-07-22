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

import {MockStakeRegistry} from "test/mock/MockStakeRegistry.sol";
import {MockRegistryCoordinator} from "test/mock/MockRegistryCoordinator.sol";

import {InitParamsLib} from "script/deploy/eigenda/DeployEigenDAConfig.sol";

import {Script} from "forge-std/Script.sol";
import {console2} from "forge-std/console2.sol";

/// @notice This script deploys EigenDA contracts and should eventually replace the other deployment scripts,
///         which cannot currently be removed due to CI depending on them.
contract DeployEigenDA is Script {
    using InitParamsLib for string;

    string constant PROXY_ADMIN = "PROXY_ADMIN";
    string constant INDEX_REGISTRY = "INDEX_REGISTRY";
    string constant STAKE_REGISTRY = "STAKE_REGISTRY";
    string constant SOCKET_REGISTRY = "SOCKET_REGISTRY";
    string constant BLS_APK_REGISTRY = "BLS_APK_REGISTRY";
    string constant REGISTRY_COORDINATOR = "REGISTRY_COORDINATOR";
    string constant THRESHOLD_REGISTRY = "THRESHOLD_REGISTRY";
    string constant RELAY_REGISTRY = "RELAY_REGISTRY";
    string constant PAYMENT_VAULT = "PAYMENT_VAULT";
    string constant DISPERSER_REGISTRY = "DISPERSER_REGISTRY";
    string constant SERVICE_MANAGER = "SERVICE_MANAGER";
    string constant EJECTION_MANAGER = "EJECTION_MANAGER";
    string constant OPERATOR_STATE_RETRIEVER = "OPERATOR_STATE_RETRIEVER";
    string constant CERT_VERIFIER = "CERT_VERIFIER";
    string constant PAUSER_REGISTRY = "PAUSER_REGISTRY";
    string constant EMPTY_CONTRACT = "EMPTY_CONTRACT";
    string constant MOCK_STAKE_REGISTRY = "MOCK_STAKE_REGISTRY";
    string constant MOCK_REGISTRY_COORDINATOR = "MOCK_REGISTRY_COORDINATOR";

    mapping(string => address) deployed; // Addresses of the deployed contracts, whether they be a proxy or not.
    mapping(string => address) impl; // Implementation addresses of the deployed contracts.
    mapping(string => bool) upgraded; // Whether the deployment of a contract is upgraded to its final implementation. Should beTrue if the contract is not a proxy

    string cfg;

    function initConfig() internal virtual {
        cfg = vm.readFile(vm.envString("DEPLOY_CONFIG_PATH"));
    }

    function run() public virtual {
        initConfig();
        vm.startBroadcast();

        // DEPLOY PROXY ADMIN
        deployed[PROXY_ADMIN] = address(new ProxyAdmin());
        impl[PROXY_ADMIN] = deployed[PROXY_ADMIN];
        upgraded[PROXY_ADMIN] = true;

        // DEPLOY MOCK IMPLEMENTATION
        impl[EMPTY_CONTRACT] = address(new EmptyContract());

        // DEPLOY PAUSER
        deployed[PAUSER_REGISTRY] = address(new PauserRegistry(cfg.pausers(), cfg.unpauser()));
        impl[PAUSER_REGISTRY] = deployed[PAUSER_REGISTRY];
        upgraded[PAUSER_REGISTRY] = true;

        // Registry coordinator requires these contracts as constructor arguments for implementation deployment
        // However, these contracts also require knowing the registry coordinator address
        // before they can be deployed, so we deploy them as inert proxies first.
        // INDEX REGISTRY
        // STAKE REGISTRY
        // SOCKET REGISTRY
        // BLS APK REGISTRY
        // SERVICE MANAGER
        deployed[INDEX_REGISTRY] =
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], deployed[PROXY_ADMIN], ""));
        deployed[SOCKET_REGISTRY] =
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], deployed[PROXY_ADMIN], ""));
        deployed[BLS_APK_REGISTRY] =
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], deployed[PROXY_ADMIN], ""));
        impl[MOCK_STAKE_REGISTRY] = address(new MockStakeRegistry(IDelegationManager(cfg.delegationManager())));
        // The service manager implementation requires the stake registry to expose the delegation manager on construction.
        deployed[STAKE_REGISTRY] =
            address(new TransparentUpgradeableProxy(impl[MOCK_STAKE_REGISTRY], deployed[PROXY_ADMIN], ""));
        // The service manager implementation requires the registry coordinator to expose the stake registry and bls APK registry on construction.
        // And this can only be done after the stake registry and bls APK registry proxies are known.
        impl[MOCK_REGISTRY_COORDINATOR] = address(
            new MockRegistryCoordinator(
                IStakeRegistry(deployed[STAKE_REGISTRY]), IBLSApkRegistry(deployed[BLS_APK_REGISTRY])
            )
        );
        deployed[REGISTRY_COORDINATOR] =
            address(new TransparentUpgradeableProxy(impl[MOCK_REGISTRY_COORDINATOR], deployed[PROXY_ADMIN], ""));
        deployed[THRESHOLD_REGISTRY] =
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], deployed[PROXY_ADMIN], ""));
        deployed[RELAY_REGISTRY] =
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], deployed[PROXY_ADMIN], ""));
        deployed[DISPERSER_REGISTRY] =
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], deployed[PROXY_ADMIN], ""));
        deployed[PAYMENT_VAULT] =
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], deployed[PROXY_ADMIN], ""));
        deployed[SERVICE_MANAGER] =
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], deployed[PROXY_ADMIN], ""));
        deployed[EJECTION_MANAGER] =
            address(new TransparentUpgradeableProxy(impl[EMPTY_CONTRACT], deployed[PROXY_ADMIN], ""));

        impl[INDEX_REGISTRY] = address(new IndexRegistry(IRegistryCoordinator(deployed[REGISTRY_COORDINATOR])));
        upgrade(INDEX_REGISTRY, "");

        impl[STAKE_REGISTRY] = address(
            new StakeRegistry(
                IRegistryCoordinator(deployed[REGISTRY_COORDINATOR]), IDelegationManager(cfg.delegationManager())
            )
        );
        upgrade(STAKE_REGISTRY, "");

        impl[SOCKET_REGISTRY] = address(new SocketRegistry(IRegistryCoordinator(deployed[REGISTRY_COORDINATOR])));
        upgrade(SOCKET_REGISTRY, "");

        impl[BLS_APK_REGISTRY] = address(new BLSApkRegistry(IRegistryCoordinator(deployed[REGISTRY_COORDINATOR])));
        upgrade(BLS_APK_REGISTRY, "");

        impl[REGISTRY_COORDINATOR] = address(
            new EigenDARegistryCoordinator(
                IServiceManager(deployed[SERVICE_MANAGER]),
                IStakeRegistry(deployed[STAKE_REGISTRY]),
                IBLSApkRegistry(deployed[BLS_APK_REGISTRY]),
                IIndexRegistry(deployed[INDEX_REGISTRY]),
                ISocketRegistry(deployed[SOCKET_REGISTRY])
            )
        );
        upgrade(
            REGISTRY_COORDINATOR,
            abi.encodeCall(
                EigenDARegistryCoordinator.initialize,
                (
                    cfg.initialOwner(),
                    cfg.churnApprover(),
                    deployed[EJECTION_MANAGER],
                    IPauserRegistry(deployed[PAUSER_REGISTRY]),
                    cfg.initialPausedStatus(),
                    cfg.operatorSetParams(),
                    cfg.minimumStakes(),
                    cfg.strategyParams()
                )
            )
        );

        impl[SERVICE_MANAGER] = address(
            new EigenDAServiceManager(
                IAVSDirectory(cfg.avsDirectory()),
                IRewardsCoordinator(cfg.rewardsCoordinator()),
                IRegistryCoordinator(deployed[REGISTRY_COORDINATOR]),
                IStakeRegistry(deployed[STAKE_REGISTRY]),
                IEigenDAThresholdRegistry(deployed[THRESHOLD_REGISTRY]),
                IEigenDARelayRegistry(deployed[RELAY_REGISTRY]),
                IPaymentVault(deployed[PAYMENT_VAULT]),
                IEigenDADisperserRegistry(deployed[DISPERSER_REGISTRY])
            )
        );
        upgrade(
            SERVICE_MANAGER,
            abi.encodeCall(
                EigenDAServiceManager.initialize,
                (
                    IPauserRegistry(deployed[PAUSER_REGISTRY]),
                    cfg.initialPausedStatus(),
                    cfg.initialOwner(),
                    cfg.batchConfirmers(),
                    cfg.rewardsInitiator()
                )
            )
        );

        impl[EJECTION_MANAGER] = address(
            new EjectionManager(
                IRegistryCoordinator(deployed[REGISTRY_COORDINATOR]), IStakeRegistry(deployed[STAKE_REGISTRY])
            )
        );
        upgrade(
            EJECTION_MANAGER,
            abi.encodeCall(EjectionManager.initialize, (cfg.initialOwner(), cfg.ejectors(), cfg.quorumEjectionParams()))
        );

        impl[THRESHOLD_REGISTRY] = address(new EigenDAThresholdRegistry());
        upgrade(
            THRESHOLD_REGISTRY,
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

        impl[RELAY_REGISTRY] = address(new EigenDARelayRegistry());
        upgrade(RELAY_REGISTRY, abi.encodeCall(EigenDARelayRegistry.initialize, (cfg.initialOwner())));

        impl[DISPERSER_REGISTRY] = address(new EigenDADisperserRegistry());
        upgrade(DISPERSER_REGISTRY, abi.encodeCall(EigenDADisperserRegistry.initialize, (cfg.initialOwner())));

        impl[PAYMENT_VAULT] = address(new PaymentVault());
        upgrade(
            PAYMENT_VAULT,
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

        deployed[OPERATOR_STATE_RETRIEVER] = address(new OperatorStateRetriever());
        impl[OPERATOR_STATE_RETRIEVER] = deployed[OPERATOR_STATE_RETRIEVER];
        upgraded[OPERATOR_STATE_RETRIEVER] = true;

        deployed[CERT_VERIFIER] = address(
            new EigenDACertVerifier(
                IEigenDAThresholdRegistry(deployed[THRESHOLD_REGISTRY]),
                IEigenDASignatureVerifier(deployed[STAKE_REGISTRY]),
                cfg.certVerifierSecurityThresholds(),
                cfg.certVerifierQuorumNumbersRequired()
            )
        );
        impl[CERT_VERIFIER] = deployed[CERT_VERIFIER];
        upgraded[CERT_VERIFIER] = true;

        ProxyAdmin(deployed[PROXY_ADMIN]).transferOwnership(cfg.initialOwner());

        vm.stopBroadcast();

        _sanityTest(PROXY_ADMIN);
        _sanityTest(INDEX_REGISTRY);
        _sanityTest(STAKE_REGISTRY);
        _sanityTest(SOCKET_REGISTRY);
        _sanityTest(BLS_APK_REGISTRY);
        _sanityTest(REGISTRY_COORDINATOR);
        _sanityTest(THRESHOLD_REGISTRY);
        _sanityTest(RELAY_REGISTRY);
        _sanityTest(PAYMENT_VAULT);
        _sanityTest(DISPERSER_REGISTRY);
        _sanityTest(SERVICE_MANAGER);
        _sanityTest(EJECTION_MANAGER);
        _sanityTest(OPERATOR_STATE_RETRIEVER);
        _sanityTest(CERT_VERIFIER);
        _sanityTest(PAUSER_REGISTRY);
    }

    function _sanityTest(string memory contractName) internal view {
        require(deployed[contractName] != address(0), string.concat("Contract not deployed: ", contractName));
        require(impl[contractName] != address(0), string.concat("Implementation not set: ", contractName));
        require(upgraded[contractName], string.concat("Contract not upgraded: ", contractName));
    }

    function upgrade(string memory contractName, bytes memory initData) internal {
        require(!upgraded[contractName], string.concat("Contract already upgraded: ", contractName));
        require(deployed[contractName] != address(0), string.concat("Contract not deployed: ", contractName));

        ProxyAdmin proxyAdmin = ProxyAdmin(deployed[PROXY_ADMIN]);
        address implementation = impl[contractName];
        TransparentUpgradeableProxy proxy = TransparentUpgradeableProxy(payable(deployed[contractName]));

        proxyAdmin.upgrade(proxy, implementation);
        if (initData.length > 0) {
            proxyAdmin.upgradeAndCall(proxy, implementation, initData);
        }
        upgraded[contractName] = true;
    }
}
