// SPDX-License-Identifier: BUSL-1.1

pragma solidity =0.8.12;

import "forge-std/Script.sol";
import "forge-std/StdToml.sol";

interface Ownable {
    function transferOwnership(address newOwner) external;
}

/// @title Upgrade Mainnet V1 to V2 Phase 2
contract UpgradeMainnet_V1_V2_P2_CORE_OPS is Script {
    using stdToml for string;

    struct InitParams {
        address coreOpsMsig;
        address ejectionManager;
        address registryCoordinator;
        address serviceManager;
        address daOpsMsig;
    }

    function run() external {
        InitParams memory initParams = _initParams();

        vm.startBroadcast(initParams.coreOpsMsig);
        
        Ownable(initParams.ejectionManager).transferOwnership(initParams.daOpsMsig);
        Ownable(initParams.registryCoordinator).transferOwnership(initParams.daOpsMsig);
        Ownable(initParams.serviceManager).transferOwnership(initParams.daOpsMsig);

        vm.stopBroadcast();
    }

    /// @dev override this if you don't want to use the environment to get the config path
    function _cfg() internal virtual returns (string memory) {
        return vm.readFile(vm.envString("UPGRADE_MAINNET_V1_V2_P1_CONFIG"));
    }

    function _initParams() internal virtual returns (InitParams memory) {
        string memory cfg = _cfg();
        return InitParams({
            coreOpsMsig: cfg.readAddress(".initParams.existing.coreOpsMsig"),
            ejectionManager: cfg.readAddress(".initParams.existing.ejectionManager"),
            registryCoordinator: cfg.readAddress(".initParams.existing.registryCoordinator"),
            serviceManager: cfg.readAddress(".initParams.existing.serviceManager"),
            daOpsMsig: cfg.readAddress(".initParams.existing.daOpsMsig")
        });
    }
}
