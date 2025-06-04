// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {Test} from "forge-std/Test.sol";

import {PaymentVaultLib} from "src/core/libraries/v3/payment/PaymentVaultLib.sol";
import {PaymentVault} from "src/core/PaymentVault.sol";

contract PaymentVaultUnit is Test {
    PaymentVault paymentVault;

    uint64 constant SCHEDULE_PERIOD = 1 days;

    function setUp() public virtual {
        paymentVault = new PaymentVault(SCHEDULE_PERIOD);
    }
}
