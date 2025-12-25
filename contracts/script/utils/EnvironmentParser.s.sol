// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {StdStyle} from "forge-std/StdStyle.sol";

interface IEigenDADirectory {
    function getAllNames() external view returns (string[] memory);
    function getAddress(string memory name) external view returns (address);
}

interface IOwnable {
    function owner() external view returns (address);
}

/// Warning: Assumes the directory correctly contains all contracts for a given network.

contract EnvironmentParser is Script {
    // ERC1967 storage slots
    bytes32 constant IMPLEMENTATION_SLOT = 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;
    bytes32 constant ADMIN_SLOT = 0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;

    struct ContractInfo {
        string name;
        address addr;
        bool isProxy;
        bool isOwnable;
        address implementation;
        address admin;
        address owner;
    }

    function run(address directoryAddress) external {
        console.log(StdStyle.bold(StdStyle.blue("=== EigenDA Registry Audit ===")));
        console.log(StdStyle.cyan("Directory:"), directoryAddress);

        IEigenDADirectory directory = IEigenDADirectory(directoryAddress);

        // Step 1: Get all names
        string[] memory names = directory.getAllNames();
        console.log(StdStyle.green("Found"), StdStyle.bold(names.length), StdStyle.green("registered contracts"));
        console.log("");

        // Step 2-5: Iterate and audit each
        ContractInfo[] memory results = new ContractInfo[](names.length);

        for (uint256 i = 0; i < names.length; i++) {
            console.log(
                StdStyle.magenta(string.concat("[", vm.toString(i + 1), "/", vm.toString(names.length), "]")),
                StdStyle.bold(names[i])
            );

            ContractInfo memory info;

            info.name = names[i];

            // Step 2: Get address for this name
            info.addr = directory.getAddress(names[i]);

            // Step 4: Check ERC1967 slots (tells us if it's a proxy)
            info.implementation = address(uint160(uint256(vm.load(info.addr, IMPLEMENTATION_SLOT))));
            info.admin = address(uint160(uint256(vm.load(info.addr, ADMIN_SLOT))));
            info.isProxy = info.implementation != address(0);

            // Step 5: Query owner
            try IOwnable(info.addr).owner() returns (address _owner) {
                info.owner = _owner;
                info.isOwnable = true;
            } catch {}

            console.log(StdStyle.cyan("  Address:"), info.addr);

            if (info.isProxy) {
                console.log(StdStyle.green("  Implementation:"), info.implementation);
                console.log(StdStyle.green("  Admin:"), info.admin);
            } else {
                console.log(StdStyle.yellow("  Not a proxy..."));
            }

            if (info.isOwnable) {
                console.log(StdStyle.green("  Owner:"), info.owner);
            } else {
                console.log(StdStyle.yellow("  Not ownable..."));
            }

            console.log("");

            results[i] = info;
        }

        // Output JSON
        outputJSON(results);
    }

    function outputJSON(ContractInfo[] memory results) internal {
        string memory json = "audit";

        vm.serializeUint(json, "timestamp", block.timestamp);
        vm.serializeUint(json, "block", block.number);
        vm.serializeUint(json, "totalContracts", results.length);

        // Build array of serialized contracts
        string[] memory contractJsons = new string[](results.length);
        for (uint256 i = 0; i < results.length; i++) {
            contractJsons[i] = serializeContract(results[i], i);
        }

        // Serialize the array
        string memory output = vm.serializeString(json, "contracts", contractJsons);

        // Write to file
        string memory outputPath = "./audit_output.json";
        vm.writeJson(output, outputPath);

        console.log(StdStyle.bold(StdStyle.green("=== Audit Complete ===")));
        console.log(StdStyle.cyan("JSON written to:"), StdStyle.bold(outputPath));
    }

    function serializeContract(ContractInfo memory info, uint256 index) internal returns (string memory) {
        string memory obj = string.concat("contract_", vm.toString(index));

        vm.serializeString(obj, "name", info.name);
        vm.serializeAddress(obj, "address", info.addr);
        vm.serializeBool(obj, "isProxy", info.isProxy);
        vm.serializeBool(obj, "isOwnable", info.isOwnable);
        vm.serializeAddress(obj, "implementation", info.implementation);
        vm.serializeAddress(obj, "admin", info.admin);

        return vm.serializeAddress(obj, "owner", info.owner);
    }
}
