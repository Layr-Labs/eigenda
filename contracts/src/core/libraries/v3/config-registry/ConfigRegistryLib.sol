// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {ConfigRegistryStorage as S} from "src/core/libraries/v3/config-registry/ConfigRegistryStorage.sol";
import {ConfigRegistryTypes as T} from "src/core/libraries/v3/config-registry/ConfigRegistryTypes.sol";

library ConfigRegistryLib {
    event ConfigBytes32Set(bytes32 nameDigest, uint256 activationKey, bytes32 value);
    event ConfigBytesSet(bytes32 nameDigest, uint256 activationKey, bytes value);

    error NameDigestNotRegistered(bytes32 nameDigest);
    error NotIncreasingActivationKey(uint256 previousActivationKey, uint256 newActivationKey);

    function getNameDigest(string memory name) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(name));
    }

    function getNumCheckpointsBytes32(bytes32 nameDigest) internal view returns (uint256) {
        return S.layout().bytes32Config.values[nameDigest].length;
    }

    function getNumCheckpointsBytes(bytes32 nameDigest) internal view returns (uint256) {
        return S.layout().bytesConfig.values[nameDigest].length;
    }

    function getConfigBytes32(bytes32 nameDigest, uint256 index) internal view returns (bytes32) {
        return S.layout().bytes32Config.values[nameDigest][index].value;
    }

    function getConfigBytes(bytes32 nameDigest, uint256 index) internal view returns (bytes memory) {
        return S.layout().bytesConfig.values[nameDigest][index].value;
    }

    function getActivationKeyBytes32(bytes32 nameDigest, uint256 index) internal view returns (uint256) {
        return S.layout().bytes32Config.values[nameDigest][index].activationKey;
    }

    function getActivationKeyBytes(bytes32 nameDigest, uint256 index) internal view returns (uint256) {
        return S.layout().bytesConfig.values[nameDigest][index].activationKey;
    }

    function getCheckpointBytes32(bytes32 nameDigest, uint256 index)
        internal
        view
        returns (T.Bytes32Checkpoint memory)
    {
        return S.layout().bytes32Config.values[nameDigest][index];
    }

    function getCheckpointBytes(bytes32 nameDigest, uint256 index) internal view returns (T.BytesCheckpoint memory) {
        return S.layout().bytesConfig.values[nameDigest][index];
    }

    function addConfigBytes32(bytes32 nameDigest, uint256 activationKey, bytes32 value) internal {
        T.Bytes32Cfg storage cfg = S.layout().bytes32Config;
        if (cfg.values[nameDigest].length > 0) {
            uint256 lastActivationKey = cfg.values[nameDigest][cfg.values[nameDigest].length - 1].activationKey;
            if (activationKey <= lastActivationKey) {
                revert NotIncreasingActivationKey(lastActivationKey, activationKey);
            }
        }
        cfg.values[nameDigest].push(T.Bytes32Checkpoint({value: value, activationKey: activationKey}));
        emit ConfigBytes32Set(nameDigest, activationKey, value);
    }

    function addConfigBytes(bytes32 nameDigest, uint256 activationKey, bytes memory value) internal {
        T.BytesCfg storage cfg = S.layout().bytesConfig;
        if (cfg.values[nameDigest].length > 0) {
            uint256 lastActivationKey = cfg.values[nameDigest][cfg.values[nameDigest].length - 1].activationKey;
            if (activationKey <= lastActivationKey) {
                revert NotIncreasingActivationKey(lastActivationKey, activationKey);
            }
        }
        cfg.values[nameDigest].push(T.BytesCheckpoint({value: value, activationKey: activationKey}));
        emit ConfigBytesSet(nameDigest, activationKey, value);
    }

    function registerNameBytes32(string memory name) internal {
        registerName(S.layout().bytes32Config.nameSet, name);
    }

    function registerNameBytes(string memory name) internal {
        registerName(S.layout().bytesConfig.nameSet, name);
    }

    function registerName(T.NameSet storage nameSet, string memory name) internal {
        bytes32 nameDigest = getNameDigest(name);
        if (bytes(nameSet.names[nameDigest]).length == 0) {
            require(bytes(name).length > 0, "Name cannot be empty");
            nameSet.names[nameDigest] = name;
            nameSet.nameList.push(name);
        }
    }

    function isNameDigestRegistered(T.NameSet storage nameSet, bytes32 nameDigest) internal view returns (bool) {
        return bytes(nameSet.names[nameDigest]).length > 0;
    }

    function isNameRegisteredBytes32(bytes32 nameDigest) internal view returns (bool) {
        return isNameDigestRegistered(S.layout().bytes32Config.nameSet, nameDigest);
    }

    function isNameRegisteredBytes(bytes32 nameDigest) internal view returns (bool) {
        return isNameDigestRegistered(S.layout().bytesConfig.nameSet, nameDigest);
    }

    function getNumRegisteredNamesBytes32() internal view returns (uint256) {
        return S.layout().bytes32Config.nameSet.nameList.length;
    }

    function getNumRegisteredNamesBytes() internal view returns (uint256) {
        return S.layout().bytesConfig.nameSet.nameList.length;
    }

    function getRegisteredNameBytes32(uint256 index) internal view returns (string memory) {
        return S.layout().bytes32Config.nameSet.nameList[index];
    }

    function getRegisteredNameBytes(uint256 index) internal view returns (string memory) {
        return S.layout().bytesConfig.nameSet.nameList[index];
    }

    function getNameBytes32(bytes32 nameDigest) internal view returns (string memory) {
        string memory name = S.layout().bytes32Config.nameSet.names[nameDigest];
        if (bytes(name).length == 0) {
            revert NameDigestNotRegistered(nameDigest);
        }
        return name;
    }

    function getNameBytes(bytes32 nameDigest) internal view returns (string memory) {
        string memory name = S.layout().bytesConfig.nameSet.names[nameDigest];
        if (bytes(name).length == 0) {
            revert NameDigestNotRegistered(nameDigest);
        }
        return name;
    }

    function getNameListBytes32() internal view returns (string[] memory) {
        return S.layout().bytes32Config.nameSet.nameList;
    }

    function getNameListBytes() internal view returns (string[] memory) {
        return S.layout().bytesConfig.nameSet.nameList;
    }
}
