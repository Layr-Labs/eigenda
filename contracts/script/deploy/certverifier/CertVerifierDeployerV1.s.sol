// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {EigenDACertVerifierV1} from "src/periphery/cert/legacy/v1/EigenDACertVerifierV1.sol";
import {RegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {EigenDAServiceManager} from "src/core/EigenDAServiceManager.sol";
import {IEigenDAServiceManager} from "src/core/interfaces/IEigenDAServiceManager.sol";
import {EigenDAThresholdRegistry} from "src/core/EigenDAThresholdRegistry.sol";
import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "src/core/interfaces/IEigenDABatchMetadataStorage.sol";
import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";

//forge script script/deploy/certverifier/CertVerifierDeployerV1.s.sol:CertVerifierDeployerV1 --sig "run(string, string)" <config.json> <output.json> --rpc-url $RPC --private-key $PRIVATE_KEY -vvvv --etherscan-api-key $ETHERSCAN_API_KEY --verify --broadcast
contract CertVerifierDeployerV1 is Script, Test {
    address eigenDACertVerifier;

    address eigenDAServiceManager;
    address eigenDAThresholdRegistry;

    function run(string memory inputJSONFile, string memory outputJSONFile) external {
        string memory path = string.concat("./script/deploy/certverifier/config/", inputJSONFile);
        string memory data = vm.readFile(path);

        bytes memory raw = stdJson.parseRaw(data, ".eigenDAServiceManager");
        eigenDAServiceManager = abi.decode(raw, (address));

        raw = stdJson.parseRaw(data, ".eigenDAThresholdRegistry");
        eigenDAThresholdRegistry = abi.decode(raw, (address));

        vm.startBroadcast();

        eigenDACertVerifier = address(
            new EigenDACertVerifierV1(
                IEigenDAThresholdRegistry(eigenDAThresholdRegistry), IEigenDABatchMetadataStorage(eigenDAServiceManager)
            )
        );

        vm.stopBroadcast();

        console.log("Deployed new EigenDACertVerifierV1 at address: ", eigenDACertVerifier);

        string memory outputPath = string.concat("./script/deploy/certverifier/output/", outputJSONFile);
        string memory parent_object = "parent object";
        string memory finalJson =
            vm.serializeAddress(parent_object, "eigenDACertVerifier", address(eigenDACertVerifier));
        vm.writeJson(finalJson, outputPath);
    }
}
