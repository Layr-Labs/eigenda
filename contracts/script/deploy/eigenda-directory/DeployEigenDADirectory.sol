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

    struct AddressConfig {
        string name;
        address value;
    }

    function _populateDirectory() internal virtual {
        // Dynamically read all contract names from the [contracts] table
        string[] memory contractNames = vm.parseTomlKeys(cfg, ".contracts");

        for (uint256 i; i < contractNames.length; i++) {
            string memory name = contractNames[i];
            address contractAddress = cfg.readAddress(string.concat(".contracts.", name));
            directory.addAddress(name, contractAddress);
        }
    }
}
