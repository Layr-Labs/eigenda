// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {ConfigRegistryStorage as S} from "src/core/libraries/v3/config-registry/ConfigRegistryStorage.sol";
import {ConfigRegistryTypes as T} from "src/core/libraries/v3/config-registry/ConfigRegistryTypes.sol";

library ConfigRegistryLib {
    event ConfigBytes32Set(bytes32 key, bytes32 value, string extraInfo);

    event ConfigBytesSet(bytes32 key, bytes value, string extraInfo);

    function getKey(string memory name) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(name));
    }

    function getConfigBytes32(bytes32 key) internal view returns (bytes32) {
        return S.layout().bytes32Config.values[key];
    }

    function getConfigBytes(bytes32 key) internal view returns (bytes memory) {
        return S.layout().bytesConfig.values[key];
    }

    function getConfigBytes32ExtraInfo(bytes32 key) internal view returns (string memory) {
        return S.layout().bytes32Config.extraInfo[key];
    }

    function getConfigBytesExtraInfo(bytes32 key) internal view returns (string memory) {
        return S.layout().bytesConfig.extraInfo[key];
    }

    function setConfigBytes32(bytes32 key, bytes32 value, string memory extraInfo) internal {
        T.Bytes32Cfg storage cfg = S.layout().bytes32Config;
        cfg.values[key] = value;
        cfg.extraInfo[key] = extraInfo;
        emit ConfigBytes32Set(key, value, extraInfo);
    }

    function setConfigBytes(bytes32 key, bytes memory value, string memory extraInfo) internal {
        T.BytesCfg storage cfg = S.layout().bytesConfig;
        cfg.values[key] = value;
        cfg.extraInfo[key] = extraInfo;
        emit ConfigBytesSet(key, value, extraInfo);
    }

    function registerKeyBytes32(string memory name) internal {
        registerKey(S.layout().bytes32Config.nameSet, name);
    }

    function registerKeyBytes(string memory name) internal {
        registerKey(S.layout().bytesConfig.nameSet, name);
    }

    function deregisterKeyBytes32(string memory name) internal {
        deregisterKey(S.layout().bytes32Config.nameSet, name);
    }

    function deregisterKeyBytes(string memory name) internal {
        deregisterKey(S.layout().bytesConfig.nameSet, name);
    }

    function registerKey(T.NameSet storage nameSet, string memory name) internal {
        bytes32 key = getKey(name);
        require(bytes(nameSet.names[key]).length == 0, "Key already exists");
        require(bytes(name).length > 0, "Name cannot be empty");
        nameSet.names[key] = name;
        nameSet.nameList.push(name);
    }

    function deregisterKey(T.NameSet storage nameSet, string memory name) internal {
        bytes32 key = getKey(name);
        require(bytes(nameSet.names[key]).length > 0, "Key does not exist");
        delete nameSet.names[key];
        // Here we utilize a simple swap and pop to remove the name from the list.
        // There is no guarantee of preservation of ordering.
        for (uint256 i; i < nameSet.nameList.length; i++) {
            if (getKey(nameSet.nameList[i]) == key) {
                nameSet.nameList[i] = nameSet.nameList[nameSet.nameList.length - 1];
                nameSet.nameList.pop();
                break;
            }
        }
    }

    function isKeyRegistered(T.NameSet storage nameSet, bytes32 key) internal view returns (bool) {
        return bytes(nameSet.names[key]).length > 0;
    }

    function isKeyRegisteredBytes32(bytes32 key) internal view returns (bool) {
        return isKeyRegistered(S.layout().bytes32Config.nameSet, key);
    }

    function isKeyRegisteredBytes(bytes32 key) internal view returns (bool) {
        return isKeyRegistered(S.layout().bytesConfig.nameSet, key);
    }
}
