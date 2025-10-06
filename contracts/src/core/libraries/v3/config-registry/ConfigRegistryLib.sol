// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {ConfigRegistryStorage as S} from "src/core/libraries/v3/config-registry/ConfigRegistryStorage.sol";
import {ConfigRegistryTypes as T} from "src/core/libraries/v3/config-registry/ConfigRegistryTypes.sol";

import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";

library ConfigRegistryLib {
    event ConfigBytes32Set(bytes32 key, bytes32 value);
    event ConfigBytesSet(bytes32 key, bytes value);

    error NoCheckpoints();

    function getKey(string memory name) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(name));
    }

    function getNumCheckpointsBytes32(bytes32 key) internal view returns (uint256) {
        return S.layout().bytes32Config.values[key].length;
    }

    function getNumCheckpointsBytes(bytes32 key) internal view returns (uint256) {
        return S.layout().bytesConfig.values[key].length;
    }

    function getConfigBytes32(bytes32 key, uint256 index) internal view returns (bytes32) {
        return S.layout().bytes32Config.values[key][index].value;
    }

    function getConfigBytes(bytes32 key, uint256 index) internal view returns (bytes memory) {
        return S.layout().bytesConfig.values[key][index].value;
    }

    function setConfigBytes32(bytes32 key, uint256 activationKey, bytes32 value) internal {
        T.Bytes32Cfg storage cfg = S.layout().bytes32Config;
        cfg.values[key].push(T.Bytes32Checkpoint({value: value, activationKey: activationKey}));
        emit ConfigBytes32Set(key, value);
    }

    function setConfigBytes(bytes32 key, uint256 activationKey, bytes memory value) internal {
        T.BytesCfg storage cfg = S.layout().bytesConfig;
        cfg.values[key].push(T.BytesCheckpoint({value: value, activationKey: activationKey}));
        emit ConfigBytesSet(key, value);
    }

    function registerKeyBytes32(string memory name) internal {
        registerKey(S.layout().bytes32Config.nameSet, name);
    }

    function registerKeyBytes(string memory name) internal {
        registerKey(S.layout().bytesConfig.nameSet, name);
    }

    function registerKey(T.NameSet storage nameSet, string memory name) internal {
        bytes32 key = getKey(name);
        if(bytes(nameSet.names[key]).length == 0){
            require(bytes(name).length > 0, "Name cannot be empty");
            nameSet.names[key] = name;
            nameSet.nameList.push(name);
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

    function getNumRegisteredKeysBytes32() internal view returns (uint256) {
        return S.layout().bytes32Config.nameSet.nameList.length;
    }

    function getNumRegisteredKeysBytes() internal view returns (uint256) {
        return S.layout().bytesConfig.nameSet.nameList.length;
    }

    function getRegisteredKeyBytes32(uint256 index) internal view returns (string memory) {
        return S.layout().bytes32Config.nameSet.nameList[index];
    }

    function getRegisteredKeyBytes(uint256 index) internal view returns (string memory) {
        return S.layout().bytesConfig.nameSet.nameList[index];
    }
}
