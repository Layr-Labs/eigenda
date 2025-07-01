// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {Script} from "forge-std/Script.sol";
import {EigenDADirectory} from "src/core/EigenDADirectory.sol";

import {TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import {stdToml} from "forge-std/StdToml.sol";
import {console2} from "forge-std/console2.sol";

contract DeployEigenDADirectory is Script {
    using stdToml for string;

    EigenDADirectory directory;
    EigenDADirectory directoryImpl;

    struct AddressEntry {
        string name;
        address value;
    }

    // Script config
    string cfg;
    AddressEntry[] entries;

    function run() external virtual {
        _initConfig();

        vm.startBroadcast();
        _deployDirectory();
        _populateDirectory();
        directory.transferOwnership(cfg.readAddress(".owner"));
        vm.stopBroadcast();
    }

    function _initConfig() internal virtual {
        cfg = vm.readFile(vm.envString("DEPLOY_CONFIG_PATH"));
        console2.logString(cfg);
    }

    function _deployDirectory() internal virtual {
        directoryImpl = new EigenDADirectory();
        directory = EigenDADirectory(
            address(
                new TransparentUpgradeableProxy(
                    address(directoryImpl),
                    cfg.readAddress(".proxyAdmin"),
                    abi.encodeCall(EigenDADirectory.initialize, msg.sender)
                )
            )
        );
    }

    function _populateDirectory() internal virtual {
        string[] memory names = cfg.readStringArray(".contractNames");
        address[] memory values = cfg.readAddressArray(".contractAddresses");

        require(names.length == values.length, "Names and addresses length mismatch");

        for (uint256 i; i < names.length; i++) {
            directory.addAddress(names[i], values[i]);
        }
    }
}
