// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {EigenDACertVerifier} from "src/core/EigenDACertVerifier.sol";
import {RegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {EigenDAServiceManager} from "src/core/EigenDAServiceManager.sol";
import {IEigenDAServiceManager} from "src/interfaces/IEigenDAServiceManager.sol";
import {EigenDAThresholdRegistry} from "src/core/EigenDAThresholdRegistry.sol";
import {IEigenDAThresholdRegistry} from "src/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "src/interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "src/interfaces/IEigenDASignatureVerifier.sol";
import {EigenDARelayRegistry} from "src/core/EigenDARelayRegistry.sol";
import {IEigenDARelayRegistry} from "src/interfaces/IEigenDARelayRegistry.sol";
import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";
import "src/interfaces/IEigenDAStructs.sol";

//forge script script/deploy/certverifier/CertVerifierDeployer.s.sol:CertVerifierDeployer --sig "run(string, string)" <config.json> <output.json> --rpc-url $RPC --private-key $PRIVATE_KEY -vvvv --etherscan-api-key $ETHERSCAN_API_KEY --verify --broadcast
contract CertVerifierDeployer is Script, Test {

    address eigenDAServiceManager; // v1 && v2
    address eigenDAThresholdRegistry; // v1 && v2
    address eigenDARelayRegistry; // v2
    address registryCoordinator; // v2
    address operatorStateRetriever;  // v2
    
    SecurityThresholds defaultSecurityThresholds;
    bytes quorumNumbersRequired;

    // deployment target
    address eigenDACertVerifier;


    /**
     * @TODO Add precursor checks using casted interfaces to ensure addreses actually map to what
     *       the cert verifier expects! This is especially critical since rollups will be using
     *       this for individual deployments in the future. This may not make sense to add until
     *       EigenDA V2 contracts are deployed on mainnet.
     * @dev loads addreses from json env to deploy new EigenDACertVerifier contract
     * @param json json file used for extracting deployment address context
     * @param outputPath target destination for outputing deployment results
     */
    function run(string memory json, string memory outputPath) external {

        // 1 - load dependency contracts from env
        string memory path = string.concat("./script/deploy/certverifier/config/", json);
        string memory data = vm.readFile(path);

        bytes memory raw = stdJson.parseRaw(data, ".eigenDAServiceManager");
        eigenDAServiceManager = abi.decode(raw, (address));

        raw = stdJson.parseRaw(data, ".eigenDAThresholdRegistry");
        eigenDAThresholdRegistry = abi.decode(raw, (address));

        raw = stdJson.parseRaw(data, ".eigenDARelayRegistry");
        eigenDARelayRegistry = abi.decode(raw, (address));

        raw = stdJson.parseRaw(data, ".registryCoordinator");
        registryCoordinator = abi.decode(raw, (address));

        raw = stdJson.parseRaw(data, ".operatorStateRetriever");
        operatorStateRetriever = abi.decode(raw, (address));

        raw = stdJson.parseRaw(data, ".defaultSecurityThresholds");
        defaultSecurityThresholds = abi.decode(raw, (SecurityThresholds));

        raw = stdJson.parseRaw(data, ".quorumNumbersRequired");
        quorumNumbersRequired = abi.decode(raw, (bytes));

        // 2 - deploy cert verifier

        vm.startBroadcast();

        eigenDACertVerifier = address(
            new EigenDACertVerifier(
                IEigenDAThresholdRegistry(eigenDAThresholdRegistry),
                IEigenDABatchMetadataStorage(eigenDAServiceManager),
                IEigenDASignatureVerifier(eigenDAServiceManager),
                IEigenDARelayRegistry(eigenDARelayRegistry),
                OperatorStateRetriever(operatorStateRetriever),
                IRegistryCoordinator(registryCoordinator),
                defaultSecurityThresholds,
                quorumNumbersRequired
        ));

        vm.stopBroadcast();
        
        console.log("Deployed new EigenDACertVerifier at address: ", eigenDACertVerifier);

        // 3 - output deployment context

        string memory outputPath = string.concat("./script/deploy/certverifier/output/", outputPath);
        string memory parentObject = "parent object";
        string memory finalJson = vm.serializeAddress(parentObject, "eigenDACertVerifier", address(eigenDACertVerifier));
        vm.writeJson(finalJson, outputPath);
    }
}