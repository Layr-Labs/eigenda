// SPDX-License-Identifier: BUSL-1.1

pragma solidity =0.8.12;

import "forge-std/Script.sol";
import "forge-std/StdToml.sol";

interface Ownable {
    function transferOwnership(address newOwner) external;
}

/// @title Upgrade Mainnet V1 to V2 Phase 2
contract UpgradeMainnet_V1_V2_P2_EXECUTOR is Script {
    using stdToml for string;

    struct InitParams {
        address executorMsig;
        address daProxyAdmin;
        address daOpsMsig;
    }

    function run() external {
        InitParams memory initParams = _initParams();

        vm.startBroadcast(initParams.executorMsig);
        
        Ownable(initParams.daProxyAdmin).transferOwnership(initParams.daOpsMsig);

        vm.stopBroadcast();
    }

    /// @dev override this if you don't want to use the environment to get the config path
    function _cfg() internal virtual returns (string memory) {
        return vm.readFile(vm.envString("UPGRADE_MAINNET_V1_V2_P1_CONFIG"));
    }

    function _initParams() internal virtual returns (InitParams memory) {
        string memory cfg = _cfg();
        return InitParams({
            executorMsig: cfg.readAddress(".initParams.existing.executorMsig"),
            daProxyAdmin: cfg.readAddress(".initParams.existing.daProxyAdmin"),
            daOpsMsig: cfg.readAddress(".initParams.existing.daOpsMsig")
        });
    }
}
