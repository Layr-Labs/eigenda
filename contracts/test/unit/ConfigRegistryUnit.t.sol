// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {Test} from "lib/forge-std/src/Test.sol";

import {ConfigRegistryStorage} from "src/core/libraries/v3/config-registry/ConfigRegistryStorage.sol";
import {ConfigRegistryLib} from "src/core/libraries/v3/config-registry/ConfigRegistryLib.sol";

contract ConfigRegistryUnit is Test {
    function testSetAndGetBytes32Config(bytes32 value) public {
        bytes32 key = ConfigRegistryLib.getKey("testKey");
        string memory extraInfo = "extraInfo";

        ConfigRegistryLib.setConfigBytes32(key, value, extraInfo);

        bytes32 retrievedValue = ConfigRegistryLib.getConfigBytes32(key);
        string memory retrievedExtraInfo = ConfigRegistryLib.getConfigBytes32ExtraInfo(key);

        assertEq(retrievedValue, value);
        assertEq(retrievedExtraInfo, extraInfo);
    }

    function testSetAndGetBytesConfig(bytes memory value) public {
        bytes32 key = ConfigRegistryLib.getKey("testKeyBytes");
        string memory extraInfo = "extraInfoBytes";

        ConfigRegistryLib.setConfigBytes(key, value, extraInfo);

        bytes memory retrievedValue = ConfigRegistryLib.getConfigBytes(key);
        string memory retrievedExtraInfo = ConfigRegistryLib.getConfigBytesExtraInfo(key);

        assertEq(keccak256(retrievedValue), keccak256(value));
        assertEq(retrievedExtraInfo, extraInfo);
    }

    function testRegisterKeyBytes(string memory key) public {
        vm.assume(bytes(key).length > 0);
        ConfigRegistryLib.registerKeyBytes(key);
        ConfigRegistryStorage.Layout storage layout = ConfigRegistryStorage.layout();
        bool found = false;
        for (uint256 i; i < layout.bytesConfig.nameSet.nameList.length; i++) {
            if (keccak256(bytes(layout.bytesConfig.nameSet.nameList[i])) == keccak256(bytes(key))) {
                found = true;
                break;
            }
        }
        assertTrue(found, "Key was not registered");
    }

    function testDeregisterKeyBytes(string memory key) public {
        vm.assume(bytes(key).length > 0);
        ConfigRegistryLib.registerKeyBytes(key);
        ConfigRegistryLib.deregisterKeyBytes(key);
        ConfigRegistryStorage.Layout storage layout = ConfigRegistryStorage.layout();
        bool found = false;
        for (uint256 i; i < layout.bytesConfig.nameSet.nameList.length; i++) {
            if (keccak256(bytes(layout.bytesConfig.nameSet.nameList[i])) == keccak256(bytes(key))) {
                found = true;
                break;
            }
        }
        assertFalse(found, "Key was not deregistered");
    }
}
