// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import "forge-std/Script.sol";
import "../src/rollup/MockRollup.sol";
import {IEigenDAServiceManager} from "../src/interfaces/IEigenDAServiceManager.sol";

contract MockRollupDeployer is Script {
    using BN254 for BN254.G1Point;

    MockRollup public mockRollup;

    BN254.G1Point public s1 = BN254.generatorG1().scalar_mul(2);
    uint256 public illegalValue = 1555;
    
    // forge script script/MockRollupDeployer.s.sol:MockRollupDeployer --sig "run(address, uint256)" <DASM address> <stake> --rpc-url $RPC_URL --private-key $PRIVATE_KEY --broadcast -vvvv
    // <security hash> = keccak256(abi.encode(blobHeader.quorumBlobParams))
    function run(address _eigenDAServiceManager, uint256 _stakeRequired) external {
        vm.startBroadcast();

        mockRollup = new MockRollup(
            IEigenDAServiceManager(_eigenDAServiceManager),
            s1,
            illegalValue,
            _stakeRequired
        );

        string memory output = "eigenDA mock rollup deployment output";
        vm.serializeAddress(output, "mockRollup", address(mockRollup));

        string memory finalJson = vm.serializeString(output, "object", output);

        vm.createDir("./script/output", true);
        vm.writeJson(finalJson, "./script/output/mock_rollup_deploy_output.json");
        vm.stopBroadcast();
    }
}
