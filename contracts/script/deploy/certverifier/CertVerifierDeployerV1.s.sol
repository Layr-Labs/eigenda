// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {EigenDACertVerifierV1} from "src/integrations/cert/legacy/v1/EigenDACertVerifierV1.sol";
import {EigenDARegistryCoordinator} from "src/core/EigenDARegistryCoordinator.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {EigenDAServiceManager} from "src/core/EigenDAServiceManager.sol";
import {IEigenDAServiceManager} from "src/core/interfaces/IEigenDAServiceManager.sol";
import {EigenDAThresholdRegistryImmutableV1} from "src/core/EigenDAThresholdRegistryImmutableV1.sol";
import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "src/core/interfaces/IEigenDABatchMetadataStorage.sol";
import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {BitmapUtils} from "lib/eigenlayer-middleware/src/libraries/BitmapUtils.sol";

//forge script script/deploy/certverifier/CertVerifierDeployerV1.s.sol:CertVerifierDeployerV1 --sig "run(string, string)" <config.json> <output.json> --rpc-url $RPC --private-key $PRIVATE_KEY -vvvv --etherscan-api-key $ETHERSCAN_API_KEY --verify --broadcast
contract CertVerifierDeployerV1 is Script, Test {
    address eigenDACertVerifier;

    address eigenDAServiceManager;
    address eigenDAThresholdRegistry;
    bytes quorumAdversaryThresholdPercentages;
    bytes quorumConfirmationThresholdPercentages;
    bytes quorumNumbersRequired;

    function run(string memory inputJSONFile, string memory outputJSONFile) external {
        // 1 - Read the input JSON file to get the EigenDAServiceManager address and thresholds
        string memory path = string.concat("./script/deploy/certverifier/config/v1/", inputJSONFile);
        string memory data = vm.readFile(path);

        bytes memory raw = stdJson.parseRaw(data, ".eigenDAServiceManager");
        eigenDAServiceManager = abi.decode(raw, (address));

        // 1.a - Parse thresholds from config as uint8[] arrays and convert to bytes
        uint8[] memory adversaryThresholds = abi.decode(stdJson.parseRaw(data, ".adversaryThresholds"), (uint8[]));
        uint8[] memory confirmationThresholds = abi.decode(stdJson.parseRaw(data, ".confirmationThresholds"), (uint8[]));
        uint8[] memory requiredQuorums = abi.decode(stdJson.parseRaw(data, ".requiredQuorums"), (uint8[]));

        // 1.b - Convert uint8[] arrays to bytes for EigenDAThresholdRegistryImmutableV1 constructor
        quorumAdversaryThresholdPercentages = uint8ArrayToBytes(adversaryThresholds);
        quorumConfirmationThresholdPercentages = uint8ArrayToBytes(confirmationThresholds);
        quorumNumbersRequired = uint8ArrayToBytes(requiredQuorums);

        // 1.c - Validate user input lengths (i.e, # of adversial/confirmation threshold value is equal to # of required quorums)
        require(
            quorumAdversaryThresholdPercentages.length == quorumNumbersRequired.length,
            "CertVerifierDeployerV1: Adversary threshold length mismatch"
        );

        require(
            quorumConfirmationThresholdPercentages.length == quorumNumbersRequired.length,
            "CertVerifierDeployerV1: Confirmation threshold length mismatch"
        );

        // 2 - Deploy the immutable threshold registry and v1 cert verifier contracts
        vm.startBroadcast();

        eigenDAThresholdRegistry = address(
            new EigenDAThresholdRegistryImmutableV1(
                quorumAdversaryThresholdPercentages, quorumConfirmationThresholdPercentages, quorumNumbersRequired
            )
        );

        eigenDACertVerifier = address(
            new EigenDACertVerifierV1(
                IEigenDAThresholdRegistry(eigenDAThresholdRegistry), IEigenDABatchMetadataStorage(eigenDAServiceManager)
            )
        );

        vm.stopBroadcast();

        // 3 - Log the deployment details and write to output JSON file

        console.log("Deployed new EigenDAThresholdRegistryImmutableV1 at address: ", address(eigenDAThresholdRegistry));
        console.log("Deployed new EigenDACertVerifierV1 at address: ", eigenDACertVerifier);

        string memory outputPath = string.concat("./script/deploy/certverifier/output/", outputJSONFile);
        string memory output = "cert verifier v1 deployment output";

        vm.serializeAddress(output, "eigenDACertVerifier", address(eigenDACertVerifier));
        vm.serializeAddress(output, "eigenDAThresholdRegistry", address(eigenDAThresholdRegistry));

        string memory finalJson = vm.serializeString(output, "object", output);
        vm.writeJson(finalJson, outputPath);
    }

    /**
     * @notice Helper function to convert uint8[] to bytes
     * @param arr The uint8 array to convert
     * @return result The bytes representation of the array
     */
    function uint8ArrayToBytes(uint8[] memory arr) internal pure returns (bytes memory) {
        bytes memory result = new bytes(arr.length);
        for (uint256 i = 0; i < arr.length; i++) {
            result[i] = bytes1(arr[i]);
        }
        return result;
    }
}
