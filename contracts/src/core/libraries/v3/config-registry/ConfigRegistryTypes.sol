// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

library ConfigRegistryTypes {
    struct NameSet {
        mapping(bytes32 => string) names;
        string[] nameList;
    }

    struct Bytes32Checkpoint {
        uint256 activationKey;
        bytes32 value;
    }

    struct BytesCheckpoint {
        uint256 activationKey;
        bytes value;
    }

    struct Bytes32Cfg {
        mapping(bytes32 => Bytes32Checkpoint[]) values;
        NameSet nameSet;
    }

    struct BytesCfg {
        mapping(bytes32 => BytesCheckpoint[]) values;
        NameSet nameSet;
    }
}
