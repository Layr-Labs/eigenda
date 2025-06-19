// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

library AddressDirectoryStorage {
    struct Layout {
        mapping(bytes32 => address) addresses;
    }

    string internal constant STORAGE_ID = "address.directory.storage";
    bytes32 internal constant STORAGE_POSITION =
        keccak256(abi.encode(uint256(keccak256(abi.encodePacked(STORAGE_ID))) - 1)) & ~bytes32(uint256(0xff));

    function layout() internal pure returns (Layout storage s) {
        bytes32 position = STORAGE_POSITION;
        assembly {
            s.slot := position
        }
    }
}
