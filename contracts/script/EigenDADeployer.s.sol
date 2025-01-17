// SPDX-License-Identifier: UNLICENSED 
pragma solidity ^0.8.9;

import {PauserRegistry} from "eigenlayer-core/contracts/permissions/PauserRegistry.sol";
import {EmptyContract} from "eigenlayer-core/test/mocks/EmptyContract.sol";

import {BLSApkRegistry} from "eigenlayer-middleware/BLSApkRegistry.sol";
import {RegistryCoordinator} from "eigenlayer-middleware/RegistryCoordinator.sol";
import {OperatorStateRetriever} from "eigenlayer-middleware/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IndexRegistry} from "eigenlayer-middleware/IndexRegistry.sol";
import {IIndexRegistry} from "eigenlayer-middleware/interfaces/IIndexRegistry.sol";
import {StakeRegistry, IStrategy} from "eigenlayer-middleware/StakeRegistry.sol";
import {IStakeRegistry, IDelegationManager} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {IServiceManager} from "eigenlayer-middleware/interfaces/IServiceManager.sol";
import {IBLSApkRegistry} from "eigenlayer-middleware/interfaces/IBLSApkRegistry.sol";
import {EigenDAServiceManager, IAVSDirectory, IRewardsCoordinator} from "../src/core/EigenDAServiceManager.sol";
import {EigenDAHasher} from "../src/libraries/EigenDAHasher.sol";
import {EigenDAThresholdRegistry} from "../src/core/EigenDAThresholdRegistry.sol";
import {EigenDABlobVerifier} from "../src/core/EigenDABlobVerifier.sol";
import {IEigenDAThresholdRegistry} from "../src/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "../src/interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "../src/interfaces/IEigenDASignatureVerifier.sol";
import {IEigenDARelayRegistry} from "../src/interfaces/IEigenDARelayRegistry.sol";
import {IPaymentVault} from "../src/interfaces/IPaymentVault.sol";
import {PaymentVault} from "../src/payments/PaymentVault.sol";
import {EigenDADisperserRegistry} from "../src/core/EigenDADisperserRegistry.sol";
import {IEigenDADisperserRegistry} from "../src/interfaces/IEigenDADisperserRegistry.sol";
import {EigenDARelayRegistry} from "../src/core/EigenDARelayRegistry.sol";
import {ISocketRegistry, SocketRegistry} from "eigenlayer-middleware/SocketRegistry.sol";
import {DeployOpenEigenLayer, ProxyAdmin, ERC20PresetFixedSupply, TransparentUpgradeableProxy, IPauserRegistry} from "./DeployOpenEigenLayer.s.sol";
import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";
import "../src/interfaces/IEigenDAStructs.sol";

// # To load the variables in the .env file
// source .env
// # To deploy and verify our contract
// forge script script/Deployer.s.sol:EigenDADeployer --rpc-url $RPC_URL  --private-key $PRIVATE_KEY --broadcast -vvvv
contract EigenDADeployer is DeployOpenEigenLayer {
    // EigenDA contracts
    ProxyAdmin public eigenDAProxyAdmin;
    PauserRegistry public eigenDAPauserReg;

    BLSApkRegistry public apkRegistry;
    EigenDAServiceManager public eigenDAServiceManager;
    EigenDAThresholdRegistry public eigenDAThresholdRegistry;
    EigenDABlobVerifier public eigenDABlobVerifier;
    RegistryCoordinator public registryCoordinator;
    IIndexRegistry public indexRegistry;
    IStakeRegistry public stakeRegistry;
    ISocketRegistry public socketRegistry;
    OperatorStateRetriever public operatorStateRetriever;
    IPaymentVault public paymentVault;
    EigenDARelayRegistry public eigenDARelayRegistry;
    IEigenDADisperserRegistry public eigenDADisperserRegistry;

    BLSApkRegistry public apkRegistryImplementation;
    EigenDAServiceManager public eigenDAServiceManagerImplementation;
    IRegistryCoordinator public registryCoordinatorImplementation;
    IIndexRegistry public indexRegistryImplementation;
    IStakeRegistry public stakeRegistryImplementation;
    EigenDAThresholdRegistry public eigenDAThresholdRegistryImplementation;
    EigenDARelayRegistry public eigenDARelayRegistryImplementation;
    ISocketRegistry public socketRegistryImplementation;
    IPaymentVault public paymentVaultImplementation;
    IEigenDADisperserRegistry public eigenDADisperserRegistryImplementation;

    uint64 _minNumSymbols = 4096;
    uint64 _pricePerSymbol = 0.4470 gwei;
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
            address tokenAddress = address(new ERC20PresetFixedSupply(string(abi.encodePacked("Token", i)), string(abi.encodePacked("TOK", i)), initialSupply, tokenOwner));
            strategyConfigs[i] = StrategyConfig({
                maxDeposits: type(uint256).max,
                maxPerDeposit: type(uint256).max,
                tokenAddress: tokenAddress,
                tokenSymbol: string(abi.encodePacked("TOK", i))
            });
        }

        _deployEigenLayer(addressConfig.eigenLayerCommunityMultisig, addressConfig.eigenLayerOperationsMultisig, addressConfig.eigenLayerPauserMultisig, strategyConfigs);

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
        
        /**
         * First, deploy upgradeable proxy contracts that **will point** to the implementations. Since the implementation contracts are
         * not yet deployed, we give these proxies an empty contract as the initial implementation, to act as if they have no code.
         */
        eigenDAServiceManager = EigenDAServiceManager(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        eigenDAThresholdRegistry = EigenDAThresholdRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        eigenDARelayRegistry = EigenDARelayRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );

        registryCoordinator = RegistryCoordinator(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        indexRegistry = IIndexRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        stakeRegistry = IStakeRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        apkRegistry = BLSApkRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        socketRegistry = ISocketRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );

        {
        paymentVault = IPaymentVault(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );

        eigenDADisperserRegistry = IEigenDADisperserRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
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
            abi.encodeWithSelector(
                EigenDADisperserRegistry.initialize.selector,
                addressConfig.eigenDACommunityMultisig
            )
        );

        indexRegistryImplementation = new IndexRegistry(
            registryCoordinator
        );

        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(indexRegistry))),
            address(indexRegistryImplementation)
        );

        stakeRegistryImplementation = new StakeRegistry(
            registryCoordinator,
            IDelegationManager(address(delegation))
        );

        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(stakeRegistry))),
            address(stakeRegistryImplementation)
        );

        apkRegistryImplementation = new BLSApkRegistry(
            registryCoordinator
        );

        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(apkRegistry))),
            address(apkRegistryImplementation)
        );

        socketRegistryImplementation = new SocketRegistry(registryCoordinator);

        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(socketRegistry))),
            address(socketRegistryImplementation)
        );

        registryCoordinatorImplementation = new RegistryCoordinator(
                IServiceManager(address(eigenDAServiceManager)),
                stakeRegistry,
                apkRegistry,
                indexRegistry,
                socketRegistry
            );

        {
            IRegistryCoordinator.OperatorSetParam[] memory operatorSetParams = new IRegistryCoordinator.OperatorSetParam[](numStrategies);
            for (uint i = 0; i < numStrategies; i++) {
                // hard code these for now
                operatorSetParams[i] = IRegistryCoordinator.OperatorSetParam({
                    maxOperatorCount: uint32(maxOperatorCount),
                    kickBIPsOfOperatorStake: 11000, // an operator needs to have kickBIPsOfOperatorStake / 10000 times the stake of the operator with the least stake to kick them out
                    kickBIPsOfTotalStake: 1001 // an operator needs to have less than kickBIPsOfTotalStake / 10000 of the total stake to be kicked out
                });
            }

            uint96[] memory minimumStakeForQuourm = new uint96[](numStrategies);
            IStakeRegistry.StrategyParams[][] memory strategyAndWeightingMultipliers = new IStakeRegistry.StrategyParams[][](numStrategies);
            for (uint i = 0; i < numStrategies; i++) {
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

        VersionedBlobParams[] memory versionedBlobParams = new VersionedBlobParams[](0);
        SecurityThresholds memory defaultSecurityThresholds = SecurityThresholds(55, 33);

        eigenDAProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenDAThresholdRegistry))),
            address(eigenDAThresholdRegistryImplementation),
            abi.encodeWithSelector(
                EigenDAThresholdRegistry.initialize.selector,
                addressConfig.eigenDACommunityMultisig,
                hex"212121",
                hex"373737",
                hex"0001",
                versionedBlobParams,
                defaultSecurityThresholds
            )
        );

        operatorStateRetriever = new OperatorStateRetriever();

        eigenDABlobVerifier = new EigenDABlobVerifier(
            IEigenDAThresholdRegistry(address(eigenDAThresholdRegistry)),
            IEigenDABatchMetadataStorage(address(eigenDAServiceManager)),
            IEigenDASignatureVerifier(address(eigenDAServiceManager)),
            IEigenDARelayRegistry(address(eigenDARelayRegistry)),
            OperatorStateRetriever(address(operatorStateRetriever)),
            IRegistryCoordinator(address(registryCoordinator))
        );

        eigenDARelayRegistryImplementation = new EigenDARelayRegistry();

        eigenDAProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenDARelayRegistry))),
            address(eigenDARelayRegistryImplementation),
            abi.encodeWithSelector(EigenDARelayRegistry.initialize.selector, addressConfig.eigenDACommunityMultisig)
        );
    }
}