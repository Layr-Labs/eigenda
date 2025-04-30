// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {EigenDACertVerifierV3} from "src/periphery/EigenDACertVerifierV3.sol";
import {RegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
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
import {EigenDATypesV1 as DATypesV1} from "src/libraries/V1/EigenDATypesV1.sol";

//forge script script/deploy/certverifier/CertVerifierDeployer.s.sol:CertVerifierDeployer --sig "run(string, string)" <config.json> <output.json> --rpc-url $RPC --private-key $PRIVATE_KEY -vvvv --etherscan-api-key $ETHERSCAN_API_KEY --verify --broadcast
contract CertVerifierDeployer is Script, Test {
    address eigenDACertVerifier;

    address eigenDAServiceManager;
    address eigenDAThresholdRegistry;
    address eigenDARelayRegistry;
    address registryCoordinator;

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

        raw = stdJson.parseRaw(data, ".defaultSecurityThresholds");
        defaultSecurityThresholds = abi.decode(raw, (DATypesV1.SecurityThresholds));

        raw = stdJson.parseRaw(data, ".quorumNumbersRequired");
        quorumNumbersRequired = abi.decode(raw, (bytes));

        vm.startBroadcast();

        eigenDACertVerifier = address(
            new EigenDACertVerifierV3(
                IEigenDAThresholdRegistry(eigenDAThresholdRegistry),
                IEigenDASignatureVerifier(eigenDAServiceManager),
                IRegistryCoordinator(registryCoordinator),
                defaultSecurityThresholds,
                quorumNumbersRequired
            )
        );

        vm.stopBroadcast();

        console.log("Deployed new EigenDACertVerifier at address: ", eigenDACertVerifier);

        string memory outputPath = string.concat("./script/deploy/certverifier/output/", outputJSONFile);
        string memory parent_object = "parent object";
        string memory finalJson =
            vm.serializeAddress(parent_object, "eigenDACertVerifier", address(eigenDACertVerifier));
        vm.writeJson(finalJson, outputPath);
    }
}
