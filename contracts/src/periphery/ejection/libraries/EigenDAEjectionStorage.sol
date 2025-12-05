// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDAEjectionTypes} from "src/periphery/ejection/libraries/EigenDAEjectionTypes.sol";

library EigenDAEjectionStorage {
    string internal constant STORAGE_ID = "eigen.da.ejection";
    bytes32 internal constant STORAGE_POSITION =
        keccak256(abi.encode(uint256(keccak256(abi.encodePacked(STORAGE_ID))) - 1)) & ~bytes32(uint256(0xff));

    struct Layout {
        mapping(address => EigenDAEjectionTypes.EjecteeState) ejectees;
        /// @dev ejectorBalanceRecord is a book-keeping value of the ejector's balance
        ///      which reflects total_ejector_amount_added - ∑(ejector_ejection_deposit_i)
        ///      where some ejector_ejection_deposit_i can either be reclaimed by the ejector OR lost
        ///      in the event of an ejectee cancellation
        mapping(address => uint256) ejectorBalanceRecord;
        uint64 delay;
        uint64 cooldown;
    }

    function layout() internal pure returns (Layout storage s) {
        bytes32 position = STORAGE_POSITION;
        assembly {
            s.slot := position
        }
    }
}
