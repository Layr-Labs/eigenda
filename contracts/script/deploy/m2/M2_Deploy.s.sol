// SPDX-License-Identifier: BUSL-1.1
/*
pragma solidity =0.8.12;

import "eigenlayer-scripts/utils/ExistingDeploymentParser.sol";

import "eigenlayer-middleware/BLSPublicKeyCompendium.sol";
import "eigenlayer-middleware/BLSRegistryCoordinatorWithIndices.sol";
import "eigenlayer-middleware/BLSPubkeyRegistry.sol";
import "eigenlayer-middleware/IndexRegistry.sol";
import "eigenlayer-middleware/StakeRegistry.sol";
import "eigenlayer-middleware/BLSOperatorStateRetriever.sol";

import "eigenlayer-core/contracts/permissions/PauserRegistry.sol";

import {EigenDAServiceManager} from "../../src/core/EigenDAServiceManager.sol";

//forge script script/m2/M2_Deploy.s.sol:Deployer_M2 --rpc-url $RPC_URL  --private-key $PRIVATE_KEY --broadcast -vvvv
contract Deployer_M2 is ExistingDeploymentParser {

    string public existingDeploymentInfoPath  = string(bytes("./script/m2/existing/M1_deployment_goerli_2023_3_23.json"));
    string public deployConfigPath = string(bytes("./script/m2/config/M2_deploy.config.json"));
    string public outputPath = "script/m2/output/M2_deployment_data.json";

    //permissioned addresses
    ProxyAdmin public eigenDAProxyAdmin;
    address public eigenDAOwner;
    address public eigenDAUpgrader;

    //non-upgradeable contracts
    BLSPublicKeyCompendium public pubkeyCompendium;
    BLSOperatorStateRetriever public blsOperatorStateRetriever;
    
    //upgradeable contracts
    EigenDAServiceManager public eigenDAServiceManager;
    BLSRegistryCoordinatorWithIndices public registryCoordinator;
    BLSPubkeyRegistry public blsPubkeyRegistry;
    IndexRegistry public indexRegistry;
    StakeRegistry public stakeRegistry;

    //upgradeable contract implementations
    EigenDAServiceManager public eigenDAServiceManagerImplementation;
    BLSRegistryCoordinatorWithIndices public registryCoordinatorImplementation;
    BLSPubkeyRegistry public blsPubkeyRegistryImplementation;
    IndexRegistry public indexRegistryImplementation;
    StakeRegistry public stakeRegistryImplementation;


    function run() external {
        // get info on all the already-deployed contracts
        _parseDeployedContracts(existingDeploymentInfoPath);

        // READ JSON CONFIG DATA
        string memory config_data = vm.readFile(deployConfigPath);

        // check that the chainID matches the one in the config
        uint256 currentChainId = block.chainid;
        uint256 configChainId = stdJson.readUint(config_data, ".chainInfo.chainId");
        emit log_named_uint("You are deploying on ChainID", currentChainId);
        require(configChainId == currentChainId, "You are on the wrong chain for this config");

        // parse initalization params and permissions from config data
        (uint96[] memory minimumStakeForQuourm, IVoteWeigher.StrategyAndWeightingMultiplier[][] memory strategyAndWeightingMultipliers) = _parseStakeRegistryParams(config_data);
        (IBLSRegistryCoordinatorWithIndices.OperatorSetParam[] memory operatorSetParams, address churner, address ejector) = _parseRegistryCoordinatorParams(config_data);

        eigenDAOwner = stdJson.readAddress(config_data, ".permissions.owner");
        eigenDAUpgrader = stdJson.readAddress(config_data, ".permissions.upgrader");

        // begin deployment
        vm.startBroadcast();

        // deploy proxy admin for ability to upgrade proxy contracts
        eigenDAProxyAdmin = new ProxyAdmin();

        //deploy non-upgradeable contracts
        pubkeyCompendium = new BLSPublicKeyCompendium();
        blsOperatorStateRetriever = new BLSOperatorStateRetriever();

        //Deploy upgradeable proxy contracts that point to empty contract implementations
        eigenDAServiceManager = EigenDAServiceManager(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        registryCoordinator = BLSRegistryCoordinatorWithIndices(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        blsPubkeyRegistry = BLSPubkeyRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        indexRegistry = IndexRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        stakeRegistry = StakeRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );

        // deploy StakeRegistry
        stakeRegistryImplementation = new StakeRegistry(
            registryCoordinator,
            strategyManager,
            IServiceManager(address(eigenDAServiceManager))
        );

        // upgrade stake registry proxy to implementation and initialbize
        eigenDAProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(stakeRegistry))),
            address(stakeRegistryImplementation),
            abi.encodeWithSelector(
                StakeRegistry.initialize.selector,
                minimumStakeForQuourm,
                strategyAndWeightingMultipliers
            )
        );

        // deploy RegistryCoordinator
        registryCoordinatorImplementation = new BLSRegistryCoordinatorWithIndices(
            slasher,
            IServiceManager(address(eigenDAServiceManager)),
            stakeRegistry,
            blsPubkeyRegistry,
            indexRegistry
        );

        // upgrade registry coordinator proxy to implementation and initialize
        eigenDAProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(registryCoordinator))),
            address(registryCoordinatorImplementation),
            abi.encodeWithSelector(
                BLSRegistryCoordinatorWithIndices.initialize.selector,
                churner,
                ejector,
                operatorSetParams,
                eigenLayerPauserReg,
                0
            )
        );

        // deploy BLSPubkeyRegistry
        blsPubkeyRegistryImplementation = new BLSPubkeyRegistry(
            registryCoordinator,
            pubkeyCompendium
        );

        // upgrade bls pubkey registry proxy to implementation
        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(blsPubkeyRegistry))),
            address(blsPubkeyRegistryImplementation)
        );

        //deploy IndexRegistry
        indexRegistryImplementation = new IndexRegistry(
            registryCoordinator
        );

        // upgrade index registry proxy to implementation
        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(indexRegistry))),
            address(indexRegistryImplementation)
        );

        //deploy EigenDAServiceManager
        eigenDAServiceManagerImplementation = new EigenDAServiceManager(
            registryCoordinator,
            strategyManager,
            delegation,
            slasher
        );

        // upgrade service manager proxy to implementation and initialize
        eigenDAProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenDAServiceManager))),
            address(eigenDAServiceManagerImplementation),
            abi.encodeWithSelector(
                EigenDAServiceManager.initialize.selector,
                eigenLayerPauserReg,
                eigenDAOwner
            )
        );

        // transfer ownership of proxy admin to upgrader
        eigenDAProxyAdmin.transferOwnership(eigenDAUpgrader);

        // end deployment
        vm.stopBroadcast();

        // sanity checks
        _verifyContractPointers(
            eigenDAServiceManager,
            registryCoordinator,
            blsPubkeyRegistry,
            indexRegistry,
            stakeRegistry
        );

        _verifyContractPointers(
            eigenDAServiceManagerImplementation,
            registryCoordinatorImplementation,
            blsPubkeyRegistryImplementation,
            indexRegistryImplementation,
            stakeRegistryImplementation
        );

        _verifyImplementations();
        _verifyInitalizations(churner, ejector, operatorSetParams, minimumStakeForQuourm, strategyAndWeightingMultipliers);

        //write output
        _writeOutput(churner, ejector);

    }

    function _parseStakeRegistryParams(string memory config_data) internal pure returns (uint96[] memory minimumStakeForQuourm, IVoteWeigher.StrategyAndWeightingMultiplier[][] memory strategyAndWeightingMultipliers) {
        bytes memory stakesConfigsRaw = stdJson.parseRaw(config_data, ".minimumStakes");
        minimumStakeForQuourm = abi.decode(stakesConfigsRaw, (uint96[]));
        
        bytes memory strategyConfigsRaw = stdJson.parseRaw(config_data, ".strategyWeights");
        strategyAndWeightingMultipliers = abi.decode(strategyConfigsRaw, (IVoteWeigher.StrategyAndWeightingMultiplier[][]));
    }

    function _parseRegistryCoordinatorParams(string memory config_data) internal returns (IBLSRegistryCoordinatorWithIndices.OperatorSetParam[] memory operatorSetParams, address churner, address ejector) {
        bytes memory operatorConfigsRaw = stdJson.parseRaw(config_data, ".operatorSetParams");
        operatorSetParams = abi.decode(operatorConfigsRaw, (IBLSRegistryCoordinatorWithIndices.OperatorSetParam[]));

        churner = stdJson.readAddress(config_data, ".permissions.churner");
        ejector = stdJson.readAddress(config_data, ".permissions.ejector");
    }

    function _verifyContractPointers(
        EigenDAServiceManager _eigenDAServiceManager,
        BLSRegistryCoordinatorWithIndices _registryCoordinator,
        BLSPubkeyRegistry _blsPubkeyRegistry,
        IndexRegistry _indexRegistry,
        StakeRegistry _stakeRegistry
    ) internal view {
        require(_eigenDAServiceManager.registryCoordinator() == registryCoordinator, "eigenDAServiceManager.registryCoordinator() != registryCoordinator");
        require(_eigenDAServiceManager.strategyManager() == strategyManager, "eigenDAServiceManager.strategyManager() != strategyManager");
        require(_eigenDAServiceManager.delegationManager() == delegation, "eigenDAServiceManager.delegationManager() != delegation");
        require(_eigenDAServiceManager.slasher() == slasher, "eigenDAServiceManager.slasher() != slasher");

        require(_registryCoordinator.slasher() == slasher, "registryCoordinator.slasher() != slasher");
        require(address(_registryCoordinator.serviceManager()) == address(eigenDAServiceManager), "registryCoordinator.eigenDAServiceManager() != eigenDAServiceManager");
        require(_registryCoordinator.stakeRegistry() == stakeRegistry, "registryCoordinator.stakeRegistry() != stakeRegistry");
        require(_registryCoordinator.blsPubkeyRegistry() == blsPubkeyRegistry, "registryCoordinator.blsPubkeyRegistry() != blsPubkeyRegistry");
        require(_registryCoordinator.indexRegistry() == indexRegistry, "registryCoordinator.indexRegistry() != indexRegistry");

        require(_blsPubkeyRegistry.registryCoordinator() == registryCoordinator, "blsPubkeyRegistry.registryCoordinator() != registryCoordinator");
        require(_blsPubkeyRegistry.pubkeyCompendium() == pubkeyCompendium, "blsPubkeyRegistry.pubkeyCompendium() != pubkeyCompendium");

        require(_indexRegistry.registryCoordinator() == registryCoordinator, "indexRegistry.registryCoordinator() != registryCoordinator");

        require(_stakeRegistry.registryCoordinator() == registryCoordinator, "stakeRegistry.registryCoordinator() != registryCoordinator");
        require(_stakeRegistry.strategyManager() == strategyManager, "stakeRegistry.strategyManager() != strategyManager");
        require(address(_stakeRegistry.serviceManager()) == address(eigenDAServiceManager), "stakeRegistry.eigenDAServiceManager() != eigenDAServiceManager");
    }

    function _verifyImplementations() internal view {
        require(eigenDAProxyAdmin.getProxyImplementation(
            TransparentUpgradeableProxy(payable(address(eigenDAServiceManager)))) == address(eigenDAServiceManagerImplementation),
            "eigenDAServiceManager: implementation set incorrectly");
        require(eigenDAProxyAdmin.getProxyImplementation(
            TransparentUpgradeableProxy(payable(address(registryCoordinator)))) == address(registryCoordinatorImplementation),
            "registryCoordinator: implementation set incorrectly");
        require(eigenDAProxyAdmin.getProxyImplementation(
            TransparentUpgradeableProxy(payable(address(blsPubkeyRegistry)))) == address(blsPubkeyRegistryImplementation),
            "blsPubkeyRegistry: implementation set incorrectly");
        require(eigenDAProxyAdmin.getProxyImplementation(
            TransparentUpgradeableProxy(payable(address(indexRegistry)))) == address(indexRegistryImplementation),
            "indexRegistry: implementation set incorrectly");
        require(eigenDAProxyAdmin.getProxyImplementation(
            TransparentUpgradeableProxy(payable(address(stakeRegistry)))) == address(stakeRegistryImplementation),
            "stakeRegistry: implementation set incorrectly");
    }

    function _verifyInitalizations(
        address churner, 
        address ejector, 
        IBLSRegistryCoordinatorWithIndices.OperatorSetParam[] memory operatorSetParams,
        uint96[] memory minimumStakeForQuourm, 
        IVoteWeigher.StrategyAndWeightingMultiplier[][] memory strategyAndWeightingMultipliers
        ) internal view {
        require(eigenDAServiceManager.owner() == eigenDAOwner, "eigenDAServiceManager.owner() != eigenDAOwner");
        require(eigenDAServiceManager.pauserRegistry() == eigenLayerPauserReg, "eigenDAServiceManager: pauser registry not set correctly");
        require(strategyManager.paused() == 0, "eigenDAServiceManager: init paused status set incorrectly");

        require(registryCoordinator.churnApprover() == churner, "registryCoordinator.churner() != churner");
        require(registryCoordinator.ejector() == ejector, "registryCoordinator.ejector() != ejector");
        require(registryCoordinator.pauserRegistry() == eigenLayerPauserReg, "registryCoordinator: pauser registry not set correctly");
        require(registryCoordinator.paused() == 0, "registryCoordinator: init paused status set incorrectly");
        
        for (uint8 i = 0; i < operatorSetParams.length; ++i) {
            require(keccak256(abi.encode(registryCoordinator.getOperatorSetParams(i))) == keccak256(abi.encode(operatorSetParams[i])), "registryCoordinator.operatorSetParams != operatorSetParams");
        }

        for (uint256 i = 0; i < minimumStakeForQuourm.length; ++i) {
            require(stakeRegistry.minimumStakeForQuorum(i) == minimumStakeForQuourm[i], "stakeRegistry.minimumStakeForQuourm != minimumStakeForQuourm");
        }

        for (uint8 i = 0; i < strategyAndWeightingMultipliers.length; ++i) {
            for(uint8 j = 0; j < strategyAndWeightingMultipliers[i].length; ++j) {
                (IStrategy strategy, uint96 multiplier) = stakeRegistry.strategiesConsideredAndMultipliers(i, j);
                require(address(strategy) == address(strategyAndWeightingMultipliers[i][j].strategy), "stakeRegistry.strategyAndWeightingMultipliers != strategyAndWeightingMultipliers");
                require(multiplier == strategyAndWeightingMultipliers[i][j].multiplier, "stakeRegistry.strategyAndWeightingMultipliers != strategyAndWeightingMultipliers");
            }
        }

        require(operatorSetParams.length == strategyAndWeightingMultipliers.length && operatorSetParams.length == minimumStakeForQuourm.length, "operatorSetParams, strategyAndWeightingMultipliers, and minimumStakeForQuourm must be the same length");
    }

    function _writeOutput(address churner, address ejector) internal {
        string memory parent_object = "parent object";

        string memory deployed_addresses = "addresses";
        vm.serializeAddress(deployed_addresses, "eigenDAProxyAdmin", address(eigenDAProxyAdmin));
        vm.serializeAddress(deployed_addresses, "blsPubKeyCompendium", address(pubkeyCompendium));
        vm.serializeAddress(deployed_addresses, "blsOperatorStateRetriever", address(blsOperatorStateRetriever));
        vm.serializeAddress(deployed_addresses, "eigenDAServiceManager", address(eigenDAServiceManager));
        vm.serializeAddress(deployed_addresses, "eigenDAServiceManagerImplementation", address(eigenDAServiceManagerImplementation));
        vm.serializeAddress(deployed_addresses, "registryCoordinator", address(registryCoordinator));
        vm.serializeAddress(deployed_addresses, "registryCoordinatorImplementation", address(registryCoordinatorImplementation));
        vm.serializeAddress(deployed_addresses, "blsPubkeyRegistry", address(blsPubkeyRegistry));
        vm.serializeAddress(deployed_addresses, "blsPubkeyRegistryImplementation", address(blsPubkeyRegistryImplementation));
        vm.serializeAddress(deployed_addresses, "indexRegistry", address(indexRegistry));
        vm.serializeAddress(deployed_addresses, "indexRegistryImplementation", address(indexRegistryImplementation));
        vm.serializeAddress(deployed_addresses, "stakeRegistry", address(stakeRegistry));
        string memory deployed_addresses_output = vm.serializeAddress(deployed_addresses, "stakeRegistryImplementation", address(stakeRegistryImplementation));

        string memory chain_info = "chainInfo";
        vm.serializeUint(chain_info, "deploymentBlock", block.number);
        string memory chain_info_output = vm.serializeUint(chain_info, "chainId", block.chainid);

        string memory permissions = "permissions";
        vm.serializeAddress(permissions, "eigenDAOwner", eigenDAOwner);
        vm.serializeAddress(permissions, "eigenDAUpgrader", eigenDAUpgrader);
        vm.serializeAddress(permissions, "eigenDAChurner", churner);
        string memory permissions_output = vm.serializeAddress(permissions, "eigenDAEjector", ejector);
        
        vm.serializeString(parent_object, chain_info, chain_info_output);
        vm.serializeString(parent_object, deployed_addresses, deployed_addresses_output);
        string memory finalJson = vm.serializeString(parent_object, permissions, permissions_output);
        vm.writeJson(finalJson, outputPath);
    } 
}
*/