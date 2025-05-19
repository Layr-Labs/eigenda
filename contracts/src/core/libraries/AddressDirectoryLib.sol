// SPDX-License-Identifier: MIT
pragma solidity ^0.8.12;

library AddressDirectoryStorage {
    struct Layout {
        mapping(bytes32 => address) values;
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

library AddressDirectoryLib {
    event AddressSet(bytes32 indexed key, address value);

    function ads() private pure returns (AddressDirectoryStorage.Layout storage) {
        return AddressDirectoryStorage.layout();
    }

    function setAddress(bytes32 key, address value) internal {
        ads().values[key] = value;
        emit AddressSet(key, value);
    }

    function getAddress(bytes32 key) internal view returns (address) {
        return ads().values[key];
    }
}
