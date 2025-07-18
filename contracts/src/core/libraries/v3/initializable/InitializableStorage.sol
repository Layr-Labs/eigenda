// SPDX-License-Identifier: MIT
pragma solidity ^0.8.12;

/// @notice Defines a storage layout based on ERC-7201
///         https://eips.ethereum.org/EIPS/eip-7201
library InitializableStorage {
    /// @custom: storage-location erc7201:initializable.storage
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
