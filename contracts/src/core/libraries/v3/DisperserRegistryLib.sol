// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV3} from "src/core/libraries/v3/EigenDATypesV3.sol";
import {SafeERC20, IERC20} from "lib/openzeppelin-contracts/contracts/token/ERC20/utils/SafeERC20.sol";

library DisperserRegistryStorage {
    string internal constant STORAGE_ID = "eigen.da.disperser.registry";
    bytes32 internal constant STORAGE_POSITION =
        keccak256(abi.encode(uint256(keccak256(abi.encodePacked(STORAGE_ID))) - 1)) & ~bytes32(uint256(0xff));

    struct User {
        uint256 deposit; // the total on demand deposit of the user
        EigenDATypesV3.Reservation reservation;
    }

    struct Quorum {
        EigenDATypesV3.QuorumPaymentProtocolConfig protocolCfg;
        EigenDATypesV3.QuorumPaymentConfig cfg;
        mapping(address => User) user;
        mapping(uint64 => uint64) reservedSymbols; // reserved symbols per period in this quorum
    }

    struct Layout {
        mapping(uint64 => Quorum) quorum;
    }

    function layout() internal pure returns (Layout storage s) {
        bytes32 position = STORAGE_POSITION;
        assembly {
            s.slot := position
        }
    }
}

library DisperserRegistryLib {
    using SafeERC20 for IERC20;

    function s() internal pure returns (DisperserRegistryStorage.Layout storage) {
        return DisperserRegistryStorage.layout();
    }
}