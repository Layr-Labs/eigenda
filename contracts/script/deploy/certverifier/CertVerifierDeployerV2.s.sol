// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {EigenDACertVerifier} from "src/integrations/cert/EigenDACertVerifier.sol";
import {EigenDAServiceManager} from "src/core/EigenDAServiceManager.sol";
import {IEigenDAServiceManager} from "src/core/interfaces/IEigenDAServiceManager.sol";
import {EigenDAThresholdRegistry} from "src/core/EigenDAThresholdRegistry.sol";
import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "src/core/interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";
import {EigenDARelayRegistry} from "src/core/EigenDARelayRegistry.sol";
import {IEigenDARelayRegistry} from "src/core/interfaces/IEigenDARelayRegistry.sol";
import {IEigenDADirectory} from "src/core/interfaces/IEigenDADirectory.sol";
import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";

//forge script script/deploy/certverifier/CertVerifierDeployerV2.s.sol:CertVerifierDeployerV2 --sig "run(string, string)" <config.json> <output.json> --rpc-url $RPC --private-key $PRIVATE_KEY -vvvv --etherscan-api-key $ETHERSCAN_API_KEY --verify --broadcast
contract CertVerifierDeployerV2 is Script, Test {
    // CertVerifierDeployerV2 is a foundry deployment contract used for deploying EigenDACertVerifier contracts
    // compatible with the EigenDA V2 protocol. 
    // 
    // There's loose correctness assumptions provided by the inabox testing framework which calls into this script
    // for deploying a verifier which is used for testing the E2E correctness of the eigenda V2 client's 
    // dispersal, VERIFICATION, & retrieval logics

    address eigenDACertVerifier;

    address eigenDADirectory;

    DATypesV1.SecurityThresholds defaultSecurityThresholds;
    bytes quorumNumbersRequired;

    // TODO(ethenotethan): is it worth detecting if these keys exist in the directory?
    string directoryServiceManagerKey = "SERVICE_MANAGER";
    string directoryThresholdRegistryKey = "THRESHOLD_REGISTRY";


    function run(string memory inputJSONFile, string memory outputJSONFile) external {
        // 1 - ingest JSON config file as string and extract dependency fields used for
        //     EigenDACertVerifier constructor params
        string memory path = string.concat("./script/deploy/certverifier/config/v2/", inputJSONFile);
        string memory data = vm.readFile(path);

        bytes memory raw = stdJson.parseRaw(data, ".eigenDADirectory");
        eigenDADirectory = abi.decode(raw, (address));
        raw = stdJson.parseRaw(data, ".defaultSecurityThresholds");
        defaultSecurityThresholds = abi.decode(raw, (DATypesV1.SecurityThresholds));

        raw = stdJson.parseRaw(data, ".quorumNumbersRequired");
        quorumNumbersRequired = abi.decode(raw, (bytes));

        // 2 - read dependency contract addresses from EigenDA Directory namespaced resolution
        //     contract and ensure that addresses are correct w.r.t their intended interfaces

        address eigenDAServiceManager = IEigenDADirectory(eigenDADirectory).getAddress(directoryServiceManagerKey);
        if (eigenDAServiceManager == address(0)) {
            revert("EigenDAServiceManager contract address cannot be nil");
        }

        // 2.a - assume we can read a batch number that's greater than zero
        uint32 batchNumber = IEigenDAServiceManager(eigenDAServiceManager).taskNumber();
        if(batchNumber == 0) {
            revert("Expected to have batch ID > 0 in EigenDAServiceManager contract storage");
        }

        // 2.b - assume we can read the blob params at version index 0 and that the struct
        //       is initialized
        address eigenDAThresholdRegistry = IEigenDADirectory(eigenDADirectory).getAddress(directoryThresholdRegistryKey);
        if (eigenDAThresholdRegistry == address(0)) {
            revert("EigenDAThresholdRegistry contract address cannot be nil");
        }

        DATypesV1.VersionedBlobParams memory blobParams = IEigenDAThresholdRegistry(eigenDAThresholdRegistry).getBlobParams(0);

        if (blobParams.codingRate == 0) {
            revert("EigenDAThresholdRegistry contract should return blob params that have been initialized at version index 0");
        }

        // 3 - validate arbitrary user input for correctness
        //
        //     these checks are done in constructor but saves some user some gas if caught here
        if (quorumNumbersRequired.length == 0 || quorumNumbersRequired.length > 256) {
            revert("quorumNumbersRequired must be in size range (0, 256]");
        }

        if (defaultSecurityThresholds.adversaryThreshold > defaultSecurityThresholds.confirmationThreshold) {
            revert("adversaryThreshold cannot be greter than the confirmationThreshold");
        }

        // 4 - broadcast single deploy tx which constructs the immutable EigenDACertVerifier contract
        //     using standard CREATE
        vm.startBroadcast();

        eigenDACertVerifier = address(
            new EigenDACertVerifier(
                IEigenDAThresholdRegistry(eigenDAThresholdRegistry),
                IEigenDASignatureVerifier(eigenDAServiceManager),
                defaultSecurityThresholds,
                quorumNumbersRequired
            )
        );

        vm.stopBroadcast();

        // 5 - output deployment context to a user named output JSON file
        console.log("Deployed new EigenDACertVerifier at address: ", eigenDACertVerifier);

        string memory outputPath = string.concat("./script/deploy/certverifier/output/", outputJSONFile);
        string memory parent_object = "parent object";
        string memory finalJson =
            vm.serializeAddress(parent_object, "eigenDACertVerifier", address(eigenDACertVerifier));
        vm.writeJson(finalJson, outputPath);
    }
}
