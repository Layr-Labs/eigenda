// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import {PauserRegistry} from
    "../lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/permissions/PauserRegistry.sol";
import {EmptyContract} from "../lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/test/mocks/EmptyContract.sol";

import {BLSApkRegistry} from "../lib/eigenlayer-middleware/src/BLSApkRegistry.sol";
import {RegistryCoordinator} from "../lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {OperatorStateRetriever} from "../lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "../lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {IndexRegistry} from "../lib/eigenlayer-middleware/src/IndexRegistry.sol";
import {IIndexRegistry} from "../lib/eigenlayer-middleware/src/interfaces/IIndexRegistry.sol";
import {StakeRegistry, IStrategy} from "../lib/eigenlayer-middleware/src/StakeRegistry.sol";
import {IStakeRegistry, IDelegationManager} from "../lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {IServiceManager} from "../lib/eigenlayer-middleware/src/interfaces/IServiceManager.sol";
import {IBLSApkRegistry} from "../lib/eigenlayer-middleware/src/interfaces/IBLSApkRegistry.sol";
import {IRelayRegistry} from "src/core/interfaces/IRelayRegistry.sol";
import {RelayRegistry} from "src/core/RelayRegistry.sol";
import {EigenDAServiceManager, IAVSDirectory, IRewardsCoordinator} from "src/core/EigenDAServiceManager.sol";
import {EigenDAThresholdRegistry} from "src/core/EigenDAThresholdRegistry.sol";
import {EigenDACertVerifierV2} from "src/integrations/cert/legacy/v2/EigenDACertVerifierV2.sol";
import {EigenDACertVerifier} from "src/integrations/cert/EigenDACertVerifier.sol";
import {EigenDACertVerifierRouter} from "src/integrations/cert/router/EigenDACertVerifierRouter.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {EigenDATypesV2 as DATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";
import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "src/core/interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";
import {IEigenDARelayRegistry} from "src/core/interfaces/IEigenDARelayRegistry.sol";
import {IPaymentVault} from "src/core/interfaces/IPaymentVault.sol";
import {PaymentVault} from "src/core/PaymentVault.sol";
import {EigenDADisperserRegistry} from "src/core/EigenDADisperserRegistry.sol";
import {IEigenDADisperserRegistry} from "src/core/interfaces/IEigenDADisperserRegistry.sol";
import {EigenDARelayRegistry} from "src/core/EigenDARelayRegistry.sol";
import {ISocketRegistry, SocketRegistry} from "../lib/eigenlayer-middleware/src/SocketRegistry.sol";
import {IEigenDADirectory, EigenDADirectory} from "src/core/EigenDADirectory.sol";
import {
    DeployOpenEigenLayer,
    ProxyAdmin,
    ERC20PresetFixedSupply,
    TransparentUpgradeableProxy,
    IPauserRegistry
} from "./DeployOpenEigenLayer.s.sol";
import {AddressDirectoryConstants} from "src/core/libraries/v3/address-directory/AddressDirectoryConstants.sol";
import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";

// NOTE: This contract is used to deploy the EigenDA contracts to a local inabox environment. It is not meant to be used in production and is only used for testing purposes.
// # To load the variables in the .env file
// source .env
// # To deploy and verify our contract
// forge script script/Deployer.s.sol:EigenDADeployer --rpc-url $RPC_URL  --private-key $PRIVATE_KEY --broadcast -vvvv
contract EigenDADeployer is DeployOpenEigenLayer {
    // EigenDA contracts
    ProxyAdmin public eigenDAProxyAdmin;
    PauserRegistry public eigenDAPauserReg;

    EigenDADirectory public eigenDADirectory;
    BLSApkRegistry public apkRegistry;
    EigenDAServiceManager public eigenDAServiceManager;
    EigenDAThresholdRegistry public eigenDAThresholdRegistry;
    EigenDACertVerifierV2 public legacyEigenDACertVerifier;
    EigenDACertVerifier public eigenDACertVerifier;
    EigenDACertVerifierRouter public eigenDACertVerifierRouter;
    RegistryCoordinator public registryCoordinator;
    IIndexRegistry public indexRegistry;
    IStakeRegistry public stakeRegistry;
    ISocketRegistry public socketRegistry;
    OperatorStateRetriever public operatorStateRetriever;
    IPaymentVault public paymentVault;
    EigenDARelayRegistry public eigenDARelayRegistry;
    IEigenDADisperserRegistry public eigenDADisperserRegistry;
    IRelayRegistry public relayRegistry;

    EigenDADirectory public eigenDADirectoryImplementation;
    BLSApkRegistry public apkRegistryImplementation;
    EigenDAServiceManager public eigenDAServiceManagerImplementation;
    EigenDACertVerifierRouter public eigenDACertVerifierRouterImplementation;
    IRegistryCoordinator public registryCoordinatorImplementation;
    IIndexRegistry public indexRegistryImplementation;
    IStakeRegistry public stakeRegistryImplementation;
    EigenDAThresholdRegistry public eigenDAThresholdRegistryImplementation;
    EigenDARelayRegistry public eigenDARelayRegistryImplementation;
    ISocketRegistry public socketRegistryImplementation;
    IPaymentVault public paymentVaultImplementation;
    IEigenDADisperserRegistry public eigenDADisperserRegistryImplementation;
    IRelayRegistry public relayRegistryImplementation;

    uint64 _minNumSymbols = 4096;
    uint64 _pricePerSymbol = 0.447 gwei;
    uint64 _priceUpdateCooldown = 1;
    uint64 _globalSymbolsPerPeriod = 131072;
    uint64 _reservationPeriodInterval = 300;
    uint64 _globalRatePeriodInterval = 30;

    struct AddressConfig {
        address eigenLayerCommunityMultisig;
        address eigenLayerOperationsMultisig;
        address eigenLayerPauserMultisig;
        address eigenDACommunityMultisig;
        address eigenDAPauser;
        address churner;
        address ejector;
        address confirmer;
    }

    function _deployEigenDAAndEigenLayerContracts(
        AddressConfig memory addressConfig,
        uint8 numStrategies,
        uint256 initialSupply,
        address tokenOwner,
        uint256 maxOperatorCount
    ) internal {
        StrategyConfig[] memory strategyConfigs = new StrategyConfig[](numStrategies);
        // deploy a token and create a strategy config for each token
        for (uint8 i = 0; i < numStrategies; i++) {
            address tokenAddress = address(
                new ERC20PresetFixedSupply(
                    string(abi.encodePacked("Token", i)), string(abi.encodePacked("TOK", i)), initialSupply, tokenOwner
                )
            );
            strategyConfigs[i] = StrategyConfig({
                maxDeposits: type(uint256).max,
                maxPerDeposit: type(uint256).max,
                tokenAddress: tokenAddress,
                tokenSymbol: string(abi.encodePacked("TOK", i))
            });
        }

        _deployEigenLayer(
            addressConfig.eigenLayerCommunityMultisig,
            addressConfig.eigenLayerOperationsMultisig,
            addressConfig.eigenLayerPauserMultisig,
            strategyConfigs
        );

        // deploy proxy admin for ability to upgrade proxy contracts
        eigenDAProxyAdmin = new ProxyAdmin();

        // deploy pauser registry
        {
            address[] memory pausers = new address[](2);
            pausers[0] = addressConfig.eigenDAPauser;
            pausers[1] = addressConfig.eigenDACommunityMultisig;
            eigenDAPauserReg = new PauserRegistry(pausers, addressConfig.eigenDACommunityMultisig);
        }

        emptyContract = new EmptyContract();

        eigenDADirectoryImplementation = new EigenDADirectory();
        eigenDADirectory = EigenDADirectory(
            address(
                new TransparentUpgradeableProxy(
                    address(eigenDADirectoryImplementation),
                    address(eigenDAProxyAdmin),
                    abi.encodeWithSelector(EigenDADirectory.initialize.selector, msg.sender)
                )
            )
        );

        /**
         * First, deploy upgradeable proxy contracts that **will point** to the implementations. Since the implementation contracts are
         * not yet deployed, we give these proxies an empty contract as the initial implementation, to act as if they have no code.
         */
        eigenDAServiceManager = EigenDAServiceManager(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        eigenDADirectory.addAddress(AddressDirectoryConstants.SERVICE_MANAGER_NAME, address(eigenDAServiceManager));
        eigenDAThresholdRegistry = EigenDAThresholdRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        eigenDADirectory.addAddress(
            AddressDirectoryConstants.THRESHOLD_REGISTRY_NAME, address(eigenDAThresholdRegistry)
        );
        eigenDARelayRegistry = EigenDARelayRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        eigenDADirectory.addAddress(AddressDirectoryConstants.RELAY_REGISTRY_LEGACY_NAME, address(eigenDARelayRegistry));
        eigenDACertVerifierRouter = EigenDACertVerifierRouter(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        eigenDADirectory.addAddress(
            AddressDirectoryConstants.CERT_VERIFIER_ROUTER_NAME, address(eigenDACertVerifierRouter)
        );

        registryCoordinator = RegistryCoordinator(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        eigenDADirectory.addAddress(AddressDirectoryConstants.REGISTRY_COORDINATOR_NAME, address(registryCoordinator));
        indexRegistry = IIndexRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        eigenDADirectory.addAddress(AddressDirectoryConstants.INDEX_REGISTRY_NAME, address(indexRegistry));
        stakeRegistry = IStakeRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        eigenDADirectory.addAddress(AddressDirectoryConstants.STAKE_REGISTRY_NAME, address(stakeRegistry));
        apkRegistry = BLSApkRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        eigenDADirectory.addAddress(AddressDirectoryConstants.BLS_APK_REGISTRY_NAME, address(apkRegistry));
        socketRegistry = ISocketRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        eigenDADirectory.addAddress(AddressDirectoryConstants.SOCKET_REGISTRY_NAME, address(socketRegistry));
        relayRegistry = IRelayRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        eigenDADirectory.addAddress(AddressDirectoryConstants.RELAY_REGISTRY_NAME, address(relayRegistry));

        {
            paymentVault = IPaymentVault(
                address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
            );
            eigenDADirectory.addAddress(AddressDirectoryConstants.PAYMENT_VAULT_NAME, address(paymentVault));

            eigenDADisperserRegistry = IEigenDADisperserRegistry(
                address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
            );
            eigenDADirectory.addAddress(
                AddressDirectoryConstants.DISPERSER_REGISTRY_NAME, address(eigenDADisperserRegistry)
            );

            paymentVaultImplementation = new PaymentVault();

            eigenDAProxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(paymentVault))),
                address(paymentVaultImplementation),
                abi.encodeWithSelector(
                    PaymentVault.initialize.selector,
                    addressConfig.eigenDACommunityMultisig,
                    _minNumSymbols,
                    _pricePerSymbol,
                    _priceUpdateCooldown,
                    _globalSymbolsPerPeriod,
                    _reservationPeriodInterval,
                    _globalRatePeriodInterval
                )
            );
        }

        eigenDADisperserRegistryImplementation = new EigenDADisperserRegistry();

        eigenDAProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenDADisperserRegistry))),
            address(eigenDADisperserRegistryImplementation),
            abi.encodeWithSelector(EigenDADisperserRegistry.initialize.selector, addressConfig.eigenDACommunityMultisig)
        );

        indexRegistryImplementation = new IndexRegistry(registryCoordinator);

        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(indexRegistry))), address(indexRegistryImplementation)
        );

        stakeRegistryImplementation = new StakeRegistry(registryCoordinator, IDelegationManager(address(delegation)));

        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(stakeRegistry))), address(stakeRegistryImplementation)
        );

        apkRegistryImplementation = new BLSApkRegistry(registryCoordinator);

        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(apkRegistry))), address(apkRegistryImplementation)
        );

        socketRegistryImplementation = new SocketRegistry(registryCoordinator);

        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(socketRegistry))), address(socketRegistryImplementation)
        );

        relayRegistryImplementation = new RelayRegistry();
        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(relayRegistry))), address(relayRegistryImplementation)
        );

        registryCoordinatorImplementation = new RegistryCoordinator(
            IServiceManager(address(eigenDAServiceManager)), stakeRegistry, apkRegistry, indexRegistry, socketRegistry
        );

        {
            IRegistryCoordinator.OperatorSetParam[] memory operatorSetParams =
                new IRegistryCoordinator.OperatorSetParam[](numStrategies);
            for (uint256 i = 0; i < numStrategies; i++) {
                // hard code these for now
                operatorSetParams[i] = IRegistryCoordinator.OperatorSetParam({
                    maxOperatorCount: uint32(maxOperatorCount),
                    kickBIPsOfOperatorStake: 11000, // an operator needs to have kickBIPsOfOperatorStake / 10000 times the stake of the operator with the least stake to kick them out
                    kickBIPsOfTotalStake: 1001 // an operator needs to have less than kickBIPsOfTotalStake / 10000 of the total stake to be kicked out
                });
            }

            uint96[] memory minimumStakeForQuourm = new uint96[](numStrategies);
            IStakeRegistry.StrategyParams[][] memory strategyAndWeightingMultipliers =
                new IStakeRegistry.StrategyParams[][](numStrategies);
            for (uint256 i = 0; i < numStrategies; i++) {
                strategyAndWeightingMultipliers[i] = new IStakeRegistry.StrategyParams[](1);
                strategyAndWeightingMultipliers[i][0] = IStakeRegistry.StrategyParams({
                    strategy: IStrategy(address(deployedStrategyArray[i])),
                    multiplier: 1 ether
                });
            }

            eigenDAProxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(registryCoordinator))),
                address(registryCoordinatorImplementation),
                abi.encodeWithSelector(
                    RegistryCoordinator.initialize.selector,
                    addressConfig.eigenDACommunityMultisig,
                    addressConfig.churner,
                    addressConfig.ejector,
                    IPauserRegistry(address(eigenDAPauserReg)),
                    0, // initial paused status is nothing paused
                    operatorSetParams,
                    minimumStakeForQuourm,
                    strategyAndWeightingMultipliers
                )
            );
        }

        eigenDAServiceManagerImplementation = new EigenDAServiceManager(
            avsDirectory,
            rewardsCoordinator,
            registryCoordinator,
            stakeRegistry,
            eigenDAThresholdRegistry,
            eigenDARelayRegistry,
            paymentVault,
            eigenDADisperserRegistry
        );

        address[] memory confirmers = new address[](1);
        confirmers[0] = addressConfig.eigenDACommunityMultisig;

        // Third, upgrade the proxy contracts to use the correct implementation contracts and initialize them.
        eigenDAProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenDAServiceManager))),
            address(eigenDAServiceManagerImplementation),
            abi.encodeWithSelector(
                EigenDAServiceManager.initialize.selector,
                eigenDAPauserReg,
                0,
                addressConfig.eigenDACommunityMultisig,
                confirmers,
                addressConfig.eigenDACommunityMultisig
            )
        );

        eigenDAThresholdRegistryImplementation = new EigenDAThresholdRegistry();

        DATypesV1.VersionedBlobParams[] memory versionedBlobParams = new DATypesV1.VersionedBlobParams[](0);
        DATypesV1.SecurityThresholds memory defaultSecurityThresholds = DATypesV1.SecurityThresholds(55, 33);

        eigenDAProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenDAThresholdRegistry))),
            address(eigenDAThresholdRegistryImplementation),
            abi.encodeWithSelector(
                EigenDAThresholdRegistry.initialize.selector,
                addressConfig.eigenDACommunityMultisig,
                hex"212121",
                hex"373737",
                hex"0001",
                versionedBlobParams
            )
        );

        operatorStateRetriever = new OperatorStateRetriever();
        eigenDADirectory.addAddress(
            AddressDirectoryConstants.OPERATOR_STATE_RETRIEVER_NAME, address(operatorStateRetriever)
        );

        // NOTE: will be deprecated in the future with subsequent release
        //       which removes the legacy V2 cert verifier entirely
        legacyEigenDACertVerifier = new EigenDACertVerifierV2(
            IEigenDAThresholdRegistry(address(eigenDAThresholdRegistry)),
            IEigenDASignatureVerifier(address(eigenDAServiceManager)),
            OperatorStateRetriever(address(operatorStateRetriever)),
            IRegistryCoordinator(address(registryCoordinator)),
            defaultSecurityThresholds,
            hex"0001"
        );
        eigenDADirectory.addAddress(
            AddressDirectoryConstants.CERT_VERIFIER_LEGACY_V2_NAME, address(legacyEigenDACertVerifier)
        );

        eigenDACertVerifier = new EigenDACertVerifier(
            IEigenDAThresholdRegistry(address(eigenDAThresholdRegistry)),
            IEigenDASignatureVerifier(address(eigenDAServiceManager)),
            defaultSecurityThresholds,
            hex"0001"
        );
        eigenDADirectory.addAddress(AddressDirectoryConstants.CERT_VERIFIER_NAME, address(eigenDACertVerifier));

        eigenDACertVerifierRouterImplementation = new EigenDACertVerifierRouter();

        eigenDAProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenDACertVerifierRouter))),
            address(eigenDACertVerifierRouterImplementation),
            abi.encodeWithSelector(
                EigenDACertVerifierRouter.initialize.selector,
                addressConfig.eigenDACommunityMultisig,
                address(eigenDACertVerifier)
            )
        );
        eigenDARelayRegistryImplementation = new EigenDARelayRegistry();

        eigenDAProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenDARelayRegistry))),
            address(eigenDARelayRegistryImplementation),
            abi.encodeWithSelector(EigenDARelayRegistry.initialize.selector, addressConfig.eigenDACommunityMultisig)
        );

        eigenDADirectory.transferOwnership(addressConfig.eigenLayerCommunityMultisig);
    }
}
