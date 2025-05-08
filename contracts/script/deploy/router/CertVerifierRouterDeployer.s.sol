// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {EigenDACertVerifierRouter} from "src/periphery/cert/router/EigenDACertVerifierRouter.sol";
import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";

/**
 * @notice Deployment script for EigenDACertVerifierRouter
 * @dev Run with:
 * forge script script/deploy/router/CertVerifierRouterDeployer.s.sol:CertVerifierRouterDeployer --sig "run(string, string)" <config.json> <output.json> --rpc-url $RPC --private-key $PRIVATE_KEY -vvvv --etherscan-api-key $ETHERSCAN_API_KEY --verify --broadcast
 */
contract CertVerifierRouterDeployer is Script, Test {
    address eigenDACertVerifierRouter;
    address initialOwner;

    // TODO: Consider adding an optional verifier parameter to instantiate one using closely upcoming block
    function run(string memory inputJSONFile, string memory outputJSONFile) external {
        string memory path = string.concat("./script/deploy/router/config/", inputJSONFile);
        string memory data = vm.readFile(path);

        bytes memory raw = stdJson.parseRaw(data, ".initialOwner");
        initialOwner = abi.decode(raw, (address));
        raw = stdJson.parseRaw(data, ".initialCertVerifier");
        address initialCertVerifier = abi.decode(raw, (address));

        vm.startBroadcast();

        // Deploy the EigenDACertVerifierRouter contract
        EigenDACertVerifierRouter router = new EigenDACertVerifierRouter();
        
        // Initialize the router with the initial owner
        router.initialize(initialOwner, initialCertVerifier);
        
        eigenDACertVerifierRouter = address(router);

        vm.stopBroadcast();

        console.log("Deployed new EigenDACertVerifierRouter at address: ", eigenDACertVerifierRouter);

        string memory outputPath = string.concat("./script/deploy/router/output/", outputJSONFile);
        string memory parent_object = "parent object";
        string memory finalJson =
            vm.serializeAddress(parent_object, "eigenDACertVerifierRouter", address(eigenDACertVerifierRouter));
        vm.writeJson(finalJson, outputPath);
    }
}