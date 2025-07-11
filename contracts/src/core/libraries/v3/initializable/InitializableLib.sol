// SPDX-License-Identifier: MIT
pragma solidity ^0.8.12;

/// @notice Defines the storage layout for an address directory based on ERC-7201
///         https://eips.ethereum.org/EIPS/eip-7201
library InitializableStorage {
    /// @custom: storage-location erc7201:address.directory.storage
    struct Layout {
        uint8 initialized;
    }

    string internal constant STORAGE_ID = "initializable.storage";
    bytes32 internal constant STORAGE_POSITION =
        keccak256(abi.encode(uint256(keccak256(abi.encodePacked(STORAGE_ID))) - 1)) & ~bytes32(uint256(0xff));

    function layout() internal pure returns (Layout storage s) {
        bytes32 position = STORAGE_POSITION;
        assembly {
            s.slot := position
        }
    }
}

library InitializableLib {
    event Initialized(uint8 version);

    error AlreadyInitialized();

    function s() private pure returns (InitializableStorage.Layout storage) {
        return InitializableStorage.layout();
    }

    function setInitializedVersion(uint8 version) internal {
        if (s().initialized >= version) {
            revert AlreadyInitialized();
        }

        s().initialized = version;
        emit Initialized(version);
    }

    function getInitializedVersion() internal view returns (uint8 version) {
        version = s().initialized;
    }
}
