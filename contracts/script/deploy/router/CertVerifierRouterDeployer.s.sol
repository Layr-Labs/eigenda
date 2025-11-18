// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {IEigenDACertVerifier} from "src/integrations/cert/interfaces/IEigenDACertVerifier.sol";
import {EigenDACertVerifierRouter} from "src/integrations/cert/router/EigenDACertVerifierRouter.sol";
import {IEigenDAServiceManager} from "src/core/interfaces/IEigenDAServiceManager.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";

struct ABNConfig {
    uint32 blockNumber;
    address certVerifier;
}

/**
 * @title CertVerifierRouterDeployer
 * @notice Deployment script for upgradable EigenDACertVerifierRouter
 * @dev This script deploys the EigenDACertVerifierRouter contract and initializes it through the proxy
 *      with the initial owner and cert verifier.
 * @dev Run with:
 *      forge script script/deploy/router/CertVerifierRouterDeployer.s.sol:CertVerifierRouterDeployer \
 *      --sig "run(string, string)" <config.json> <output.json> \
 *      --rpc-url $RPC \
 *      --private-key $PRIVATE_KEY \
 *      -vvvv \
 *      --etherscan-api-key $ETHERSCAN_API_KEY \
 *      --verify \
 *      --broadcast
 */
contract CertVerifierRouterDeployer is Script, Test {
    // Configuration parameters
    address initialOwner;
    address proxyAdmin;
    uint32[] initABNs;
    address[] initCertVerifiers;

    // Mappings for efficient duplicate detection
    mapping(uint32 => bool) private seenBlockNumbers;
    mapping(address => bool) private seenCertVerifiers;

    function run(string memory inputJSONFile, string memory outputJSONFile) external {
        // 1. Read the configuration from the JSON input file
        string memory configPath = string.concat("./script/deploy/router/config/", inputJSONFile);
        string memory configData = vm.readFile(configPath);

        // Parse configuration parameters
        initialOwner = stdJson.readAddress(configData, ".initialOwner");
        setABNConfigs(configData);
        proxyAdmin = stdJson.readAddress(configData, ".proxyAdmin");

        // 2. Deploy the implementation and proxy contracts
        vm.startBroadcast();

        EigenDACertVerifierRouter implementation = new EigenDACertVerifierRouter();

        // Deploy proxy and initialize in one step
        bytes memory initData =
            abi.encodeCall(EigenDACertVerifierRouter.initialize, (initialOwner, initABNs, initCertVerifiers));

        TransparentUpgradeableProxy proxy =
            new TransparentUpgradeableProxy(address(implementation), address(proxyAdmin), initData);

        vm.stopBroadcast();

        // 4. Output the deployed addresses to a JSON file

        string memory outputPath =
            string.concat("./script/deploy/router/output/", vm.toString(block.chainid), "/", outputJSONFile);
        string memory parent = "parent object";
        string memory finalJson = vm.serializeAddress(parent, "eigenDACertVerifierRouter", address(proxy));
        finalJson = vm.serializeAddress(parent, "eigenDACertVerifierRouterImplementation", address(implementation));

        vm.writeJson(finalJson, outputPath);
    }

    function setABNConfigs(string memory configData) internal {
        bytes memory raw = stdJson.parseRaw(configData, ".initABNConfigs");
        ABNConfig[] memory configs = abi.decode(raw, (ABNConfig[]));
        for (uint256 i; i < configs.length; i++) {
            uint32 blockNumber = configs[i].blockNumber;
            address certVerifier = configs[i].certVerifier;

            // run user input safety checks
            //
            // 1) the cert verifier's dependencies appear correctly initialized
            address thresholdRegistry = address(IEigenDACertVerifier(certVerifier).eigenDAThresholdRegistry());
            IEigenDAThresholdRegistry(thresholdRegistry).nextBlobVersion();

            address serviceManager = address(IEigenDACertVerifier(certVerifier).eigenDASignatureVerifier());
            // 2) the signature verifier address can be cast to IServiceManager
            IEigenDAServiceManager(serviceManager).taskNumber();

            // 3) ensure no duplicate block numbers
            assertFalse(seenBlockNumbers[blockNumber], "Duplicate block number detected");
            seenBlockNumbers[blockNumber] = true;

            // 4) ensure no duplicate cert verifiers
            assertFalse(seenCertVerifiers[certVerifier], "Duplicate cert verifier detected");
            seenCertVerifiers[certVerifier] = true;

            initABNs.push(blockNumber);
            initCertVerifiers.push(certVerifier);
        }
    }
}
