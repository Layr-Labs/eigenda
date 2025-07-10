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

    function registerKey(string memory name) internal {
        AddressDirectoryStorage.Layout storage s = AddressDirectoryStorage.layout();
        bytes32 key = getKey(name);
        require(bytes(s.names[key]).length == 0, "Key already exists");
        s.names[key] = name;
        s.nameList.push(name);
    }

    function deregisterKey(string memory name) internal {
        AddressDirectoryStorage.Layout storage s = AddressDirectoryStorage.layout();
        bytes32 key = getKey(name);
        require(bytes(s.names[key]).length > 0, "Key does not exist");
        delete s.names[key];
        // Here we utilize a simple swap and pop to remove the name from the list.
        // There is no guarantee of preservation of ordering.
        for (uint256 i; i < s.nameList.length; i++) {
            if (getKey(s.nameList[i]) == key) {
                s.nameList[i] = s.nameList[s.nameList.length - 1];
                s.nameList.pop();
                break;
            }
        }
    }

    function getName(bytes32 key) internal view returns (string memory) {
        return AddressDirectoryStorage.layout().names[key];
    }

    function getNameList() internal view returns (string[] memory) {
        return AddressDirectoryStorage.layout().nameList;
    }
}
