// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

library ConfigRegistryTypes {
    struct NameSet {
        mapping(bytes32 => string) names;
        string[] nameList;
    }

    struct Bytes32Cfg {
        mapping(bytes32 => bytes32) values;
        mapping(bytes32 => string) extraInfo;
        NameSet nameSet;
    }

    struct BytesCfg {
        mapping(bytes32 => bytes) values;
        mapping(bytes32 => string) extraInfo;
        NameSet nameSet;
    }
}
