// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import "@eigenlayer-scripts/middleware/DeployOpenEigenLayer.s.sol";

import "@eigenlayer-core/contracts/permissions/PauserRegistry.sol";
import "@eigenlayer-core/test/mocks/EmptyContract.sol";

import "@eigenlayer-middleware/BLSPublicKeyCompendium.sol";
import "@eigenlayer-middleware/BLSRegistryCoordinatorWithIndices.sol";
import "@eigenlayer-middleware/BLSPubkeyRegistry.sol";
import "@eigenlayer-middleware/IndexRegistry.sol";
import "@eigenlayer-middleware/StakeRegistry.sol";
import "@eigenlayer-middleware/BLSOperatorStateRetriever.sol";

import {EigenDAServiceManager} from "../src/core/EigenDAServiceManager.sol";
import "../src/libraries/EigenDAHasher.sol";

import "forge-std/Test.sol";

import "forge-std/Script.sol";
import "forge-std/StdJson.sol";

// TODO: REVIEW AND FIX THIS ENTIRE SCRIPT

// # To load the variables in the .env file
// source .env

// # To deploy and verify our contract
// forge script script/Deployer.s.sol:EigenDADeployer --rpc-url $RPC_URL  --private-key $PRIVATE_KEY --broadcast -vvvv
contract EigenDADeployer is DeployOpenEigenLayer {
    // EigenDA contracts
    ProxyAdmin public eigenDAProxyAdmin;
    PauserRegistry public eigenDAPauserReg;

    BLSPublicKeyCompendium public pubkeyCompendium;
    EigenDAServiceManager public eigenDAServiceManager;
    BLSRegistryCoordinatorWithIndices public registryCoordinator;
    IBLSPubkeyRegistry public blsPubkeyRegistry;
    IIndexRegistry public indexRegistry;
    IStakeRegistry public stakeRegistry;
    BLSOperatorStateRetriever public blsOperatorStateRetriever;

    EigenDAServiceManager public eigenDAServiceManagerImplementation;
    IBLSRegistryCoordinatorWithIndices public registryCoordinatorImplementation;
    IBLSPubkeyRegistry public blsPubkeyRegistryImplementation;
    IIndexRegistry public indexRegistryImplementation;
    IStakeRegistry public stakeRegistryImplementation;

    struct AddressConfig {
        address eigenLayerCommunityMultisig;
        address eigenLayerOperationsMultisig;
        address eigenLayerPauserMultisig;
        address eigenDACommunityMultisig;
        address eigenDAPauser;
        address churner;
        address ejector;
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

        // hard-coded inputs

        /**
         * First, deploy upgradeable proxy contracts that **will point** to the implementations. Since the implementation contracts are
         * not yet deployed, we give these proxies an empty contract as the initial implementation, to act as if they have no code.
         */
        eigenDAServiceManager = EigenDAServiceManager(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        pubkeyCompendium = new BLSPublicKeyCompendium();
        registryCoordinator = BLSRegistryCoordinatorWithIndices(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        blsPubkeyRegistry = IBLSPubkeyRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        indexRegistry = IIndexRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        stakeRegistry = IStakeRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );

        // Second, deploy the *implementation* contracts, using the *proxy contracts* as inputs
        {
            stakeRegistryImplementation = new StakeRegistry(
                registryCoordinator,
                strategyManager,
                IServiceManager(address(eigenDAServiceManager))
            );

            // set up a quorum with each strategy that needs to be set up
            uint96[] memory minimumStakeForQuourm = new uint96[](numStrategies);
            IVoteWeigher.StrategyAndWeightingMultiplier[][] memory strategyAndWeightingMultipliers = new IVoteWeigher.StrategyAndWeightingMultiplier[][](numStrategies);
            for (uint i = 0; i < numStrategies; i++) {
                strategyAndWeightingMultipliers[i] = new IVoteWeigher.StrategyAndWeightingMultiplier[](1);
                strategyAndWeightingMultipliers[i][0] = IVoteWeigher.StrategyAndWeightingMultiplier({
                    strategy: deployedStrategyArray[i],
                    multiplier: 1 ether
                });
            }

            eigenDAProxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(stakeRegistry))),
                address(stakeRegistryImplementation),
                abi.encodeWithSelector(
                    StakeRegistry.initialize.selector,
                    minimumStakeForQuourm,
                    strategyAndWeightingMultipliers
                )
            );
        }

        registryCoordinatorImplementation = new BLSRegistryCoordinatorWithIndices(
            slasher,
            IServiceManager(address(eigenDAServiceManager)),
            stakeRegistry,
            blsPubkeyRegistry,
            indexRegistry
        );
        
        {
            IBLSRegistryCoordinatorWithIndices.OperatorSetParam[] memory operatorSetParams = new IBLSRegistryCoordinatorWithIndices.OperatorSetParam[](numStrategies);
            for (uint i = 0; i < numStrategies; i++) {
                // hard code these for now
                operatorSetParams[i] = IBLSRegistryCoordinatorWithIndices.OperatorSetParam({
                    maxOperatorCount: uint32(maxOperatorCount),
                    kickBIPsOfOperatorStake: 11000, // an operator needs to have kickBIPsOfOperatorStake / 10000 times the stake of the operator with the least stake to kick them out
                    kickBIPsOfTotalStake: 1001 // an operator needs to have less than kickBIPsOfTotalStake / 10000 of the total stake to be kicked out
                });
            }
            eigenDAProxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(registryCoordinator))),
                address(registryCoordinatorImplementation),
                abi.encodeWithSelector(
                    BLSRegistryCoordinatorWithIndices.initialize.selector,
                    addressConfig.churner,
                    addressConfig.ejector,
                    operatorSetParams,
                    IPauserRegistry(address(eigenDAPauserReg)),
                    0 // initial paused status is nothing paused
                )
            );
        }

        blsPubkeyRegistryImplementation = new BLSPubkeyRegistry(
            registryCoordinator,
            pubkeyCompendium
        );

        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(blsPubkeyRegistry))),
            address(blsPubkeyRegistryImplementation)
        );

        indexRegistryImplementation = new IndexRegistry(
            registryCoordinator
        );

        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(indexRegistry))),
            address(indexRegistryImplementation)
        );

        eigenDAServiceManagerImplementation = new EigenDAServiceManager(
            registryCoordinator,
            strategyManager,
            delegation,
            slasher
        );

        // Third, upgrade the proxy contracts to use the correct implementation contracts and initialize them.
        eigenDAProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenDAServiceManager))),
            address(eigenDAServiceManagerImplementation),
            abi.encodeWithSelector(
                EigenDAServiceManager.initialize.selector,
                eigenDAPauserReg,
                addressConfig.eigenDACommunityMultisig,
                0,
                addressConfig.eigenDACommunityMultisig
            )
        );

        blsOperatorStateRetriever = new BLSOperatorStateRetriever();
    }
}
