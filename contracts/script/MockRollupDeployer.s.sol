// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.9;

import "forge-std/Script.sol";
import "../test/rollup/MockRollup.sol";
import {IEigenDAServiceManager} from "../src/interfaces/IEigenDAServiceManager.sol";

contract MockRollupDeployer is Script {
    using BN254 for BN254.G1Point;

    MockRollup public mockRollup;

    BN254.G1Point public s1 = BN254.generatorG1().scalar_mul(2);
    
    // forge script script/MockRollupDeployer.s.sol:MockRollupDeployer --sig "run(address)" <DASM address> --rpc-url $RPC_URL --private-key $PRIVATE_KEY -vvvv // --broadcast
    function run(address _eigenDAServiceManager) external {
        vm.startBroadcast();

        mockRollup = new MockRollup(
            IEigenDAServiceManager(_eigenDAServiceManager),
            s1
        );

        vm.stopBroadcast();

        string memory output = "eigenDA mock rollup deployment output";
        vm.serializeAddress(output, "mockRollup", address(mockRollup));
        string memory finalJson = vm.serializeString(output, "object", output);
        vm.createDir("./script/output", true);
        vm.writeJson(finalJson, "./script/output/mock_rollup_deploy_output.json");
    }
}
