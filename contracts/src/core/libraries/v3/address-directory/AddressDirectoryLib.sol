// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {AddressDirectoryStorage} from "src/core/libraries/v3/address-directory/AddressDirectoryStorage.sol";

library AddressDirectoryLib {
    event AddressSet(bytes32 key, address indexed value);

    function getKey(string memory name) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(name));
    }

    function getAddress(bytes32 key) internal view returns (address) {
        return AddressDirectoryStorage.layout().addresses[key];
    }

    function setAddress(bytes32 key, address value) internal {
        AddressDirectoryStorage.layout().addresses[key] = value;
        emit AddressSet(key, value);
    }
}
