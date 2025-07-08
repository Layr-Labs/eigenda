// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {EigenDACertVerifierRouter} from "src/integrations/cert/router/EigenDACertVerifierRouter.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";

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
    address initialCertVerifier;

    function run(string memory inputJSONFile, string memory outputJSONFile) external {
        // 1. Read the configuration from the JSON input file
        string memory configPath = string.concat("./script/deploy/router/config/", inputJSONFile);
        string memory configData = vm.readFile(configPath);

        // Parse configuration parameters
        bytes memory raw = stdJson.parseRaw(configData, ".initialOwner");
        initialOwner = abi.decode(raw, (address));

        raw = stdJson.parseRaw(configData, ".initialCertVerifier");
        initialCertVerifier = abi.decode(raw, (address));

        raw = stdJson.parseRaw(configData, ".proxyAdmin");
        proxyAdmin = abi.decode(raw, (address));

        // 2. Deploy the implementation and proxy contracts
        vm.startBroadcast();

        EigenDACertVerifierRouter implementation = new EigenDACertVerifierRouter();

        // Deploy proxy and initialize in one step
        bytes memory initData =
            abi.encodeWithSignature("initialize(address,address)", initialOwner, initialCertVerifier);

        TransparentUpgradeableProxy proxy =
            new TransparentUpgradeableProxy(address(implementation), address(tx.origin), initData);

        // 3. Transfer proxy admin to the specified address
        proxy.changeAdmin(proxyAdmin);

        vm.stopBroadcast();

        // 4. Output the deployed addresses to a JSON file

        string memory outputPath =
            string.concat("./script/deploy/router/output/", vm.toString(block.chainid), "/", outputJSONFile);
        string memory parent = "parent object";
        string memory finalJson = vm.serializeAddress(parent, "eigenDACertVerifierRouter", address(proxy));
        finalJson = vm.serializeAddress(parent, "eigenDACertVerifierRouterImplementation", address(implementation));

        vm.writeJson(finalJson, outputPath);
    }
}
