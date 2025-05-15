// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {EigenDACertVerifierV2} from "src/periphery/cert/legacy/v2/EigenDACertVerifierV2.sol";
import {RegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {EigenDAServiceManager} from "src/core/EigenDAServiceManager.sol";
import {IEigenDAServiceManager} from "src/core/interfaces/IEigenDAServiceManager.sol";
import {EigenDAThresholdRegistry} from "src/core/EigenDAThresholdRegistry.sol";
import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "src/core/interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";
import {EigenDARelayRegistry} from "src/core/EigenDARelayRegistry.sol";
import {IEigenDARelayRegistry} from "src/core/interfaces/IEigenDARelayRegistry.sol";
import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";

//forge script script/deploy/certverifier/CertVerifierDeployerV2.s.sol:CertVerifierDeployerV2 --sig "run(string, string)" <config.json> <output.json> --rpc-url $RPC --private-key $PRIVATE_KEY -vvvv --etherscan-api-key $ETHERSCAN_API_KEY --verify --broadcast
contract CertVerifierDeployerV2 is Script, Test {
    address eigenDACertVerifier;

    address eigenDAServiceManager;
    address eigenDAThresholdRegistry;
    address eigenDARelayRegistry;
    address registryCoordinator;
    address operatorStateRetriever;

    DATypesV1.SecurityThresholds defaultSecurityThresholds;
    bytes quorumNumbersRequired;

    function run(string memory inputJSONFile, string memory outputJSONFile) external {
        string memory path = string.concat("./script/deploy/certverifier/config/", inputJSONFile);
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
        defaultSecurityThresholds = abi.decode(raw, (DATypesV1.SecurityThresholds));

        raw = stdJson.parseRaw(data, ".quorumNumbersRequired");
        quorumNumbersRequired = abi.decode(raw, (bytes));

        vm.startBroadcast();

        eigenDACertVerifier = address(
            new EigenDACertVerifierV2(
                IEigenDAThresholdRegistry(eigenDAThresholdRegistry),
                IEigenDASignatureVerifier(eigenDAServiceManager),
                OperatorStateRetriever(operatorStateRetriever),
                IRegistryCoordinator(registryCoordinator),
                defaultSecurityThresholds,
                quorumNumbersRequired
            )
        );

        vm.stopBroadcast();

        console.log("Deployed new EigenDACertVerifierV2 at address: ", eigenDACertVerifier);

        string memory outputPath = string.concat("./script/deploy/certverifier/output/", outputJSONFile);
        string memory parent_object = "parent object";
        string memory finalJson =
            vm.serializeAddress(parent_object, "eigenDACertVerifier", address(eigenDACertVerifier));
        vm.writeJson(finalJson, outputPath);
    }
}
