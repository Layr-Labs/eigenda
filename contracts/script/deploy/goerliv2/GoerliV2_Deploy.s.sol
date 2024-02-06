// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {PauserRegistry} from "eigenlayer-core/contracts/permissions/PauserRegistry.sol";
import {EmptyContract} from "eigenlayer-core/test/mocks/EmptyContract.sol";

import {BLSApkRegistry} from "eigenlayer-middleware/BLSApkRegistry.sol";
import {IBLSApkRegistry} from "eigenlayer-middleware/interfaces/IBLSApkRegistry.sol";
import {RegistryCoordinator} from "eigenlayer-middleware/RegistryCoordinator.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IndexRegistry} from "eigenlayer-middleware/IndexRegistry.sol";
import {IIndexRegistry} from "eigenlayer-middleware/interfaces/IIndexRegistry.sol";
import {StakeRegistry} from "eigenlayer-middleware/StakeRegistry.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {EigenDAServiceManager} from "src/core/EigenDAServiceManager.sol";
import {IServiceManager} from "eigenlayer-middleware/interfaces/IServiceManager.sol";
import {OperatorStateRetriever} from "eigenlayer-middleware/OperatorStateRetriever.sol";

import "eigenlayer-scripts/utils/ExistingDeploymentParser.sol";
import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";


//forge script script/deploy/goerliv2/GoerliV2_Deploy.s.sol:Deployer_GV2 --rpc-url $RPC_URL  --private-key $PRIVATE_KEY -vvvv //--broadcast 
contract Deployer_GV2 is ExistingDeploymentParser {

    //string public existingDeploymentInfoPath  = string(bytes("./script/deploy/existing/M1_deployment_goerli_2023_3_23.json"));
    string public existingDeploymentInfoPath  = string(bytes("./script/deploy/existing/GV2_preprod_deployment_2024_30_1.json"));
    //string public deployConfigPath = string(bytes("./script/deploy/goerliv2/config/prod.config.json"));
    string public deployConfigPath = string(bytes("./script/deploy/goerliv2/config/preprod.config.json"));
    string public outputPath = string.concat("script/deploy/goerliv2/output/GV2_preprod_deployment_data.json");

    ProxyAdmin public eigenDAProxyAdmin;
    address public eigenDAOwner;
    address public eigenDAUpgrader;
    address public batchConfirmer;
    address public pauser;
    uint256 public initalPausedStatus;

    BLSApkRegistry public apkRegistry;
    EigenDAServiceManager public eigenDAServiceManager;
    RegistryCoordinator public registryCoordinator;
    IndexRegistry public indexRegistry;
    StakeRegistry public stakeRegistry;
    OperatorStateRetriever public operatorStateRetriever;

    BLSApkRegistry public apkRegistryImplementation;
    EigenDAServiceManager public eigenDAServiceManagerImplementation;
    RegistryCoordinator public registryCoordinatorImplementation;
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

        // parse the addresses of permissioned roles
        eigenDAOwner = stdJson.readAddress(config_data, ".permissions.owner");
        eigenDAUpgrader = stdJson.readAddress(config_data, ".permissions.upgrader");
        batchConfirmer = stdJson.readAddress(config_data, ".permissions.batchConfirmer");
        pauser = stdJson.readAddress(config_data, ".permissions.pauser");
        initalPausedStatus = stdJson.readUint(config_data, ".permissions.initalPausedStatus");

        if(pauser == address(0)) {
            pauser = address(eigenLayerPauserReg);
        }

        vm.startBroadcast();

        // deploy proxy admin for ability to upgrade proxy contracts
        eigenDAProxyAdmin = new ProxyAdmin();

        /**
         * First, deploy upgradeable proxy contracts that **will point** to the implementations. Since the implementation contracts are
         * not yet deployed, we give these proxies an empty contract as the initial implementation, to act as if they have no code.
         */
        eigenDAServiceManager = EigenDAServiceManager(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        registryCoordinator = RegistryCoordinator(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        indexRegistry = IndexRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        stakeRegistry = StakeRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        apkRegistry = BLSApkRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );

        //deploy index registry implementation
        indexRegistryImplementation = new IndexRegistry(
            registryCoordinator
        );

        //upgrade index registry proxy to implementation
        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(indexRegistry))),
            address(indexRegistryImplementation)
        );

        //deploy stake registry implementation
        stakeRegistryImplementation = new StakeRegistry(
            registryCoordinator,
            delegation
        );

        //upgrade stake registry proxy to implementation
        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(stakeRegistry))),
            address(stakeRegistryImplementation)
        );

        //deploy apk registry implementation
        apkRegistryImplementation = new BLSApkRegistry(
            registryCoordinator
        );

        //upgrade apk registry proxy to implementation
        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(apkRegistry))),
            address(apkRegistryImplementation)
        );

        //deploy the registry coordinator implementation.
        registryCoordinatorImplementation = new RegistryCoordinator(
            IServiceManager(address(eigenDAServiceManager)),
            stakeRegistry,
            apkRegistry,
            indexRegistry
        );

        {
        // parse initalization params and permissions from config data
        (
            uint96[] memory minimumStakeForQuourm, 
            IStakeRegistry.StrategyParams[][] memory strategyAndWeightingMultipliers
        ) = _parseStakeRegistryParams(config_data);
        (
            IRegistryCoordinator.OperatorSetParam[] memory operatorSetParams, 
            address churner, 
            address ejector
        ) = _parseRegistryCoordinatorParams(config_data);

        //upgrade the registry coordinator proxy to implementation
        eigenDAProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(registryCoordinator))),
            address(registryCoordinatorImplementation),
            abi.encodeWithSelector(
                RegistryCoordinator.initialize.selector,
                eigenDAOwner,
                churner,
                ejector,
                IPauserRegistry(pauser),
                initalPausedStatus, 
                operatorSetParams, 
                minimumStakeForQuourm,
                strategyAndWeightingMultipliers 
            )
        );
        }

        //deploy the eigenDA service manager implementation
        eigenDAServiceManagerImplementation = new EigenDAServiceManager(
            avsDirectory,
            registryCoordinator,
            stakeRegistry
        );

        //upgrade the eigenDA service manager proxy to implementation
        eigenDAProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenDAServiceManager))),
            address(eigenDAServiceManagerImplementation),
            abi.encodeWithSelector(
                EigenDAServiceManager.initialize.selector,
                IPauserRegistry(pauser),
                initalPausedStatus,
                eigenDAOwner,
                batchConfirmer
            )
        );

        //deploy the operator state retriever
        operatorStateRetriever = new OperatorStateRetriever();

        // transfer ownership of proxy admin to upgrader
        eigenDAProxyAdmin.transferOwnership(eigenDAUpgrader);

        vm.stopBroadcast();

        // sanity checks
        _verifyContractPointers(
            apkRegistry,
            eigenDAServiceManager,
            registryCoordinator,
            indexRegistry,
            stakeRegistry
        );

        _verifyContractPointers(
            apkRegistryImplementation,
            eigenDAServiceManagerImplementation,
            registryCoordinatorImplementation,
            indexRegistryImplementation,
            stakeRegistryImplementation
        );

        _verifyImplementations();
        _verifyInitalizations(config_data);

        //write output
        _writeOutput(config_data);
    }

    function _verifyContractPointers(
        BLSApkRegistry _apkRegistry,
        EigenDAServiceManager _eigenDAServiceManager,
        RegistryCoordinator _registryCoordinator,
        IndexRegistry _indexRegistry,
        StakeRegistry _stakeRegistry
    ) internal view {
        require(address(_apkRegistry.registryCoordinator()) == address(registryCoordinator), "blsApkRegistry.registryCoordinator() != registryCoordinator");

        require(address(_indexRegistry.registryCoordinator()) == address(registryCoordinator), "indexRegistry.registryCoordinator() != registryCoordinator");

        require(address(_stakeRegistry.registryCoordinator()) == address(registryCoordinator), "stakeRegistry.registryCoordinator() != registryCoordinator");
        require(address(_stakeRegistry.delegation()) == address(delegation), "stakeRegistry.delegationManager() != delegation");

        require(address(_eigenDAServiceManager.registryCoordinator()) == address(registryCoordinator), "eigenDAServiceManager.registryCoordinator() != registryCoordinator");
        require(address(_eigenDAServiceManager.stakeRegistry()) == address(stakeRegistry), "eigenDAServiceManager.stakeRegistry() != stakeRegistry");

        require(address(_registryCoordinator.serviceManager()) == address(eigenDAServiceManager), "registryCoordinator.eigenDAServiceManager() != eigenDAServiceManager");
        require(address(_registryCoordinator.stakeRegistry()) == address(stakeRegistry), "registryCoordinator.stakeRegistry() != stakeRegistry");
        require(address(_registryCoordinator.blsApkRegistry()) == address(apkRegistry), "registryCoordinator.blsApkRegistry() != blsPubkeyRegistry");
        require(address(_registryCoordinator.indexRegistry()) == address(indexRegistry), "registryCoordinator.indexRegistry() != indexRegistry");
    }

    function _verifyImplementations() internal view {
        require(eigenDAProxyAdmin.getProxyImplementation(
            TransparentUpgradeableProxy(payable(address(eigenDAServiceManager)))) == address(eigenDAServiceManagerImplementation),
            "eigenDAServiceManager: implementation set incorrectly");
        require(eigenDAProxyAdmin.getProxyImplementation(
            TransparentUpgradeableProxy(payable(address(registryCoordinator)))) == address(registryCoordinatorImplementation),
            "registryCoordinator: implementation set incorrectly");
        require(eigenDAProxyAdmin.getProxyImplementation(
            TransparentUpgradeableProxy(payable(address(apkRegistry)))) == address(apkRegistryImplementation),
            "blsApkRegistry: implementation set incorrectly");
        require(eigenDAProxyAdmin.getProxyImplementation(
            TransparentUpgradeableProxy(payable(address(indexRegistry)))) == address(indexRegistryImplementation),
            "indexRegistry: implementation set incorrectly");
        require(eigenDAProxyAdmin.getProxyImplementation(
            TransparentUpgradeableProxy(payable(address(stakeRegistry)))) == address(stakeRegistryImplementation),
            "stakeRegistry: implementation set incorrectly");
    }

    function _verifyInitalizations(string memory config_data) internal {
        (
            uint96[] memory minimumStakeForQuourm, 
            IStakeRegistry.StrategyParams[][] memory strategyAndWeightingMultipliers
        ) = _parseStakeRegistryParams(config_data);
        (
            IRegistryCoordinator.OperatorSetParam[] memory operatorSetParams, 
            address churner, 
            address ejector
        ) = _parseRegistryCoordinatorParams(config_data);

        require(eigenDAServiceManager.owner() == eigenDAOwner, "eigenDAServiceManager.owner() != eigenDAOwner");
        require(eigenDAServiceManager.pauserRegistry() == IPauserRegistry(pauser), "eigenDAServiceManager: pauser registry not set correctly");
        require(eigenDAServiceManager.batchConfirmer() == batchConfirmer, "eigenDAServiceManager.batchConfirmer() != batchConfirmer");
        require(eigenDAServiceManager.paused() == initalPausedStatus, "eigenDAServiceManager: init paused status set incorrectly");

        require(registryCoordinator.owner() == eigenDAOwner, "registryCoordinator.owner() != eigenDAOwner");
        require(registryCoordinator.churnApprover() == churner, "registryCoordinator.churner() != churner");
        require(registryCoordinator.ejector() == ejector, "registryCoordinator.ejector() != ejector");
        require(registryCoordinator.pauserRegistry() == IPauserRegistry(pauser), "registryCoordinator: pauser registry not set correctly");
        require(registryCoordinator.paused() == initalPausedStatus, "registryCoordinator: init paused status set incorrectly");
        
        for (uint8 i = 0; i < operatorSetParams.length; ++i) {
            require(keccak256(abi.encode(registryCoordinator.getOperatorSetParams(i))) == keccak256(abi.encode(operatorSetParams[i])), "registryCoordinator.operatorSetParams != operatorSetParams");
        }

        for (uint8 i = 0; i < minimumStakeForQuourm.length; ++i) {
            require(stakeRegistry.minimumStakeForQuorum(i) == minimumStakeForQuourm[i], "stakeRegistry.minimumStakeForQuourm != minimumStakeForQuourm");
        }

        for (uint8 i = 0; i < strategyAndWeightingMultipliers.length; ++i) {
            for(uint8 j = 0; j < strategyAndWeightingMultipliers[i].length; ++j) {
                IStakeRegistry.StrategyParams memory strategyParams = stakeRegistry.strategyParamsByIndex(i, j);
                require(address(strategyParams.strategy) == address(strategyAndWeightingMultipliers[i][j].strategy), "stakeRegistry.strategyAndWeightingMultipliers != strategyAndWeightingMultipliers");
                require(strategyParams.multiplier == strategyAndWeightingMultipliers[i][j].multiplier, "stakeRegistry.strategyAndWeightingMultipliers != strategyAndWeightingMultipliers");
            }
        }

        require(operatorSetParams.length == strategyAndWeightingMultipliers.length && operatorSetParams.length == minimumStakeForQuourm.length, "operatorSetParams, strategyAndWeightingMultipliers, and minimumStakeForQuourm must be the same length");
    }

    function _writeOutput(string memory config_data) internal {
        string memory parent_object = "parent object";

        string memory deployed_addresses = "addresses";
        vm.serializeAddress(deployed_addresses, "eigenDAProxyAdmin", address(eigenDAProxyAdmin));
        vm.serializeAddress(deployed_addresses, "operatorStateRetriever", address(operatorStateRetriever));
        vm.serializeAddress(deployed_addresses, "eigenDAServiceManager", address(eigenDAServiceManager));
        vm.serializeAddress(deployed_addresses, "eigenDAServiceManagerImplementation", address(eigenDAServiceManagerImplementation));
        vm.serializeAddress(deployed_addresses, "registryCoordinator", address(registryCoordinator));
        vm.serializeAddress(deployed_addresses, "registryCoordinatorImplementation", address(registryCoordinatorImplementation));
        vm.serializeAddress(deployed_addresses, "blsApkRegistry", address(apkRegistry));
        vm.serializeAddress(deployed_addresses, "blsApkRegistryImplementation", address(apkRegistryImplementation));
        vm.serializeAddress(deployed_addresses, "indexRegistry", address(indexRegistry));
        vm.serializeAddress(deployed_addresses, "indexRegistryImplementation", address(indexRegistryImplementation));
        vm.serializeAddress(deployed_addresses, "stakeRegistry", address(stakeRegistry));
        string memory deployed_addresses_output = vm.serializeAddress(deployed_addresses, "stakeRegistryImplementation", address(stakeRegistryImplementation));

        string memory chain_info = "chainInfo";
        vm.serializeUint(chain_info, "deploymentBlock", block.number);
        string memory chain_info_output = vm.serializeUint(chain_info, "chainId", block.chainid);

        address churner = stdJson.readAddress(config_data, ".permissions.churner");
        address ejector = stdJson.readAddress(config_data, ".permissions.ejector");
        string memory permissions = "permissions";
        vm.serializeAddress(permissions, "eigenDAOwner", eigenDAOwner);
        vm.serializeAddress(permissions, "eigenDAUpgrader", eigenDAUpgrader);
        vm.serializeAddress(permissions, "eigenDAChurner", churner);
        vm.serializeAddress(permissions, "eigenDABatchConfirmer", batchConfirmer);
        vm.serializeAddress(permissions, "pauserRegistry", pauser);
        string memory permissions_output = vm.serializeAddress(permissions, "eigenDAEjector", ejector);
        
        vm.serializeString(parent_object, chain_info, chain_info_output);
        vm.serializeString(parent_object, deployed_addresses, deployed_addresses_output);
        string memory finalJson = vm.serializeString(parent_object, permissions, permissions_output);
        vm.writeJson(finalJson, outputPath);
    } 

    function _parseStakeRegistryParams(string memory config_data) internal pure returns (uint96[] memory minimumStakeForQuourm, IStakeRegistry.StrategyParams[][] memory strategyAndWeightingMultipliers) {
        bytes memory stakesConfigsRaw = stdJson.parseRaw(config_data, ".minimumStakes");
        minimumStakeForQuourm = abi.decode(stakesConfigsRaw, (uint96[]));
        
        bytes memory strategyConfigsRaw = stdJson.parseRaw(config_data, ".strategyWeights");
        strategyAndWeightingMultipliers = abi.decode(strategyConfigsRaw, (IStakeRegistry.StrategyParams[][]));
    }

    function _parseRegistryCoordinatorParams(string memory config_data) internal returns (IRegistryCoordinator.OperatorSetParam[] memory operatorSetParams, address churner, address ejector) {
        bytes memory operatorConfigsRaw = stdJson.parseRaw(config_data, ".operatorSetParams");
        operatorSetParams = abi.decode(operatorConfigsRaw, (IRegistryCoordinator.OperatorSetParam[]));

        churner = stdJson.readAddress(config_data, ".permissions.churner");
        ejector = stdJson.readAddress(config_data, ".permissions.ejector");
    }
}
