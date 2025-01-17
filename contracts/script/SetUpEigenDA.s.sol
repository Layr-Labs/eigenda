// SPDX-License-Identifier: UNLICENSED 
pragma solidity ^0.8.9;

import {PauserRegistry} from "eigenlayer-core/contracts/permissions/PauserRegistry.sol";
import {EmptyContract} from "eigenlayer-core/test/mocks/EmptyContract.sol";

import {RegistryCoordinator} from "eigenlayer-middleware/RegistryCoordinator.sol";
import {IndexRegistry} from "eigenlayer-middleware/IndexRegistry.sol";
import {StakeRegistry} from "eigenlayer-middleware/StakeRegistry.sol";
import {IIndexRegistry} from "eigenlayer-middleware/interfaces/IIndexRegistry.sol";

import {EigenDAServiceManager} from "../src/core/EigenDAServiceManager.sol";
import {PaymentVault} from "../src/payments/PaymentVault.sol";
import {IPaymentVault} from "../src/interfaces/IPaymentVault.sol";
import {EigenDAHasher} from "../src/libraries/EigenDAHasher.sol";
import {EigenDADeployer} from "./EigenDADeployer.s.sol";
import {EigenLayerUtils} from "./EigenLayerUtils.s.sol";

import "./DeployOpenEigenLayer.s.sol";
import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";


// Helper function to create single-element arrays
function toArray(address element) pure returns (address[] memory) {
    address[] memory arr = new address[](1);
    arr[0] = element;
    return arr;
}

function toArray(uint256 element) pure returns (uint256[] memory) {
    uint256[] memory arr = new uint256[](1);
    arr[0] = element;
    return arr;
}


// # To load the variables in the .env file
// source .env
// # To deploy and verify our contract
// forge script script/Deployer.s.sol:SetupEigenDA --rpc-url $RPC_URL  --private-key $PRIVATE_KEY --broadcast -vvvv
contract SetupEigenDA is EigenDADeployer, EigenLayerUtils {

    string deployConfigPath = "script/input/eigenda_deploy_config.json";

    // deploy all the EigenDA contracts. Relies on many EL contracts having already been deployed.
    function run() external {
        

        // READ JSON CONFIG DATA
        string memory config_data = vm.readFile(deployConfigPath);

        
        uint8 numStrategies = uint8(stdJson.readUint(config_data, ".numStrategies"));
        {
            AddressConfig memory addressConfig;
            addressConfig.eigenLayerCommunityMultisig = msg.sender;
            addressConfig.eigenLayerOperationsMultisig = msg.sender;
            addressConfig.eigenLayerPauserMultisig = msg.sender;
            addressConfig.eigenDACommunityMultisig = msg.sender;
            addressConfig.eigenDAPauser = msg.sender;
            addressConfig.churner = msg.sender;
            addressConfig.ejector = msg.sender;
            addressConfig.confirmer = msg.sender;

            uint256 initialSupply = 1000 ether;
            address tokenOwner = msg.sender;
            uint256 maxOperatorCount = 3;
            // bytes memory parsedData = vm.parseJson(config_data);
            bool useDefaults = stdJson.readBool(config_data, ".useDefaults");
            if(!useDefaults) {
                addressConfig.eigenLayerCommunityMultisig = stdJson.readAddress(config_data, ".eigenLayerCommunityMultisig");
                addressConfig.eigenLayerOperationsMultisig = stdJson.readAddress(config_data, ".eigenLayerOperationsMultisig");
                addressConfig.eigenLayerPauserMultisig = stdJson.readAddress(config_data, ".eigenLayerPauserMultisig");
                addressConfig.eigenDACommunityMultisig = stdJson.readAddress(config_data, ".eigenDACommunityMultisig");
                addressConfig.eigenDAPauser = stdJson.readAddress(config_data, ".eigenDAPauser");
                addressConfig.churner = stdJson.readAddress(config_data, ".churner");
                addressConfig.ejector = stdJson.readAddress(config_data, ".ejector");

                initialSupply = stdJson.readUint(config_data, ".initialSupply");
                tokenOwner = stdJson.readAddress(config_data, ".tokenOwner");
                maxOperatorCount = stdJson.readUint(config_data, ".maxOperatorCount");
            }

            
            addressConfig.confirmer = vm.addr(stdJson.readUint(config_data, ".confirmerPrivateKey"));


            vm.startBroadcast();

            _deployEigenDAAndEigenLayerContracts(
                addressConfig,
                numStrategies,
                initialSupply,
                tokenOwner,
                maxOperatorCount
            );
            
            eigenDAServiceManager.setBatchConfirmer(addressConfig.confirmer);

            vm.stopBroadcast();
        }

        uint256[] memory stakerPrivateKeys = stdJson.readUintArray(config_data, ".stakerPrivateKeys");
        address[] memory stakers = new address[](stakerPrivateKeys.length);
        for (uint i = 0; i < stakers.length; i++) {
            stakers[i] = vm.addr(stakerPrivateKeys[i]);
        }
        uint256[] memory stakerETHAmounts = new uint256[](stakers.length);
        // 0.1 eth each
        for (uint i = 0; i < stakerETHAmounts.length; i++) {
            stakerETHAmounts[i] = 0.1 ether;
        }

        // stakerTokenAmount[i][j] is the amount of token i that staker j will receive
        bytes memory stakerTokenAmountsRaw = stdJson.parseRaw(config_data, ".stakerTokenAmounts");
        uint256[][] memory stakerTokenAmounts = abi.decode(stakerTokenAmountsRaw, (uint256[][]));

        uint256[] memory operatorPrivateKeys = stdJson.readUintArray(config_data, ".operatorPrivateKeys");
        address[] memory operators = new address[](operatorPrivateKeys.length);
        for (uint i = 0; i < operators.length; i++) {
            operators[i] = vm.addr(operatorPrivateKeys[i]);
        }
        uint256[] memory operatorETHAmounts = new uint256[](operators.length);
        // 5 eth each
        for (uint i = 0; i < operatorETHAmounts.length; i++) {
            operatorETHAmounts[i] = 5 ether;
        }

        vm.startBroadcast();
        // Allocate eth to stakers, operators, dispserser clients
        _allocate(
            IERC20(address(0)),
            stakers,
            stakerETHAmounts
        );

        _allocate(
            IERC20(address(0)),
            operators,
            operatorETHAmounts
        );

        // Allocate tokens to stakers
        for (uint8 i = 0; i < numStrategies; i++) {
            _allocate(
                IERC20(deployedStrategyArray[i].underlyingToken()),
                stakers,
                stakerTokenAmounts[i]
            );
        }

        {
            IStrategy[] memory strategies = new IStrategy[](numStrategies);
            bool[] memory transferLocks = new bool[](numStrategies);
            for (uint8 i = 0; i < numStrategies; i++) {
                strategies[i] = deployedStrategyArray[i];
            }
            strategyManager.addStrategiesToDepositWhitelist(strategies, transferLocks);
        }

        vm.stopBroadcast();

        // Register operators with EigenLayer
        for (uint256 i = 0; i < operatorPrivateKeys.length; i++) {
            vm.broadcast(operatorPrivateKeys[i]);
            address earningsReceiver = address(uint160(uint256(keccak256(abi.encodePacked(operatorPrivateKeys[i])))));
            address delegationApprover = address(0); //address(uint160(uint256(keccak256(abi.encodePacked(earningsReceiver)))));
            uint32 stakerOptOutWindowBlocks = 100;
            string memory metadataURI = string.concat("https://urmom.com/operator/", vm.toString(i));
            delegation.registerAsOperator(IDelegationManager.OperatorDetails(earningsReceiver, delegationApprover, stakerOptOutWindowBlocks), metadataURI);
        }


        // Register Reservations for client as the eigenDACommunityMultisig
        IPaymentVault.Reservation memory reservation = IPaymentVault.Reservation({
            symbolsPerSecond: 452198,
            startTimestamp: uint64(block.timestamp),
            endTimestamp: uint64(block.timestamp + 1000000000),
            quorumNumbers: hex"0001",
            quorumSplits: hex"3232"
        });
        address clientAddress = address(0x1aa8226f6d354380dDE75eE6B634875c4203e522);
        vm.startBroadcast(msg.sender);
        paymentVault.setReservation(clientAddress, reservation);
        // Deposit OnDemand 
        paymentVault.depositOnDemand{value: 0.1 ether}(clientAddress);
        vm.stopBroadcast();

        // Deposit stakers into EigenLayer and delegate to operators
        for (uint256 i = 0; i < stakerPrivateKeys.length; i++) {
            vm.startBroadcast(stakerPrivateKeys[i]);
            for (uint j = 0; j < numStrategies; j++) {
                if(stakerTokenAmounts[j][i] > 0) {
                    deployedStrategyArray[j].underlyingToken().approve(address(strategyManager), stakerTokenAmounts[j][i]);
                    strategyManager.depositIntoStrategy(
                        deployedStrategyArray[j],
                        deployedStrategyArray[j].underlyingToken(),
                        stakerTokenAmounts[j][i]
                    );
                }
            }
            IDelegationManager.SignatureWithExpiry memory approverSignatureAndExpiry;
            delegation.delegateTo(operators[i], approverSignatureAndExpiry, bytes32(0));
            vm.stopBroadcast();
        }

        string memory output = "eigenDA deployment output";
        vm.serializeAddress(output, "eigenDAServiceManager", address(eigenDAServiceManager));
        vm.serializeAddress(output, "operatorStateRetriever", address(operatorStateRetriever));
        vm.serializeAddress(output, "blsApkRegistry" , address(apkRegistry));
        vm.serializeAddress(output, "registryCoordinator", address(registryCoordinator));
        vm.serializeAddress(output, "blobVerifier", address(eigenDABlobVerifier));

        string memory finalJson = vm.serializeString(output, "object", output);

        vm.createDir("./script/output", true);
        vm.writeJson(finalJson, "./script/output/eigenda_deploy_output.json");        
    }
}
