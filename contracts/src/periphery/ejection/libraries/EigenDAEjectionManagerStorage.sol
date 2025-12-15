// SPDX-License-Identifier: MIT
pragma solidity ^0.8.12;

/// @notice Defines a storage layout based on ERC-7201
///         https://eips.ethereum.org/EIPS/eip-7201
library EigenDAEjectionManagerStorage {
    /// @custom: storage-location erc7201:eigenda.ejection.manager.storage
    struct Layout {
        address accessControl;
        address depositToken;
        address blsApkKeyRegistry;
        address serviceManager;
        address registryCoordinator;
        uint256 estimatedGasUsedWithoutSig;
        uint256 estimatedGasUsedWithSig;
        uint256 depositBaseFeeMultiplier;
    }

    string internal constant STORAGE_ID = "eigenda.ejection.manager.storage";
    bytes32 internal constant STORAGE_POSITION =
        keccak256(abi.encode(uint256(keccak256(abi.encodePacked(STORAGE_ID))) - 1)) & ~bytes32(uint256(0xff));

    function layout() internal pure returns (Layout storage s) {
        bytes32 position = STORAGE_POSITION;
        assembly {
            s.slot := position
        }
    }
}
