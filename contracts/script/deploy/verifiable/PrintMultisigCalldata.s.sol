// SPDX-License-Identifier: BUSL-1.1

pragma solidity =0.8.12;

import {
    DeploymentInitializer,
    CalldataInitParams,
    CalldataRegistryCoordinatorParams,
    CalldataThresholdRegistryParams,
    CalldataServiceManagerParams,
    InitParamsLib
} from "./DeploymentInitializer.sol";

import "forge-std/Script.sol";

import {console2} from "forge-std/console2.sol";

contract PrintMultisigCalldata is Script {
    using InitParamsLib for string;

    string cfg;

    function run() external {
        _initConfig();
        CalldataInitParams memory params = cfg.calldataInitParams();
        // Unfortunately, the parameters being calldata breaks the use of abi.encodeCall. So we manually form the calldata here.
        bytes memory encodedParams = abi.encode(params);
        bytes4 selector = DeploymentInitializer.initializeDeployment.selector;

        bytes memory calldataParams = abi.encodeWithSelector(selector, encodedParams);
        console2.log("Calldata: ");
        console2.logBytes(calldataParams);
    }

    /// @dev override this if you don't want to use the environment to get the config path
    function _initConfig() internal virtual {
        cfg = vm.readFile(vm.envString("DEPLOY_CONFIG_PATH"));
    }
}
