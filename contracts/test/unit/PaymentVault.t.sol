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

    /// @notice Tests that we can add a reservation successfully.
    function test_AddReservation() public {

    }

    /// @notice Tests that adding a reservation reverts if a reservation is currently active.
    function test_AddReservationRevertsIfReservationStillActive() public {

    }

    /// @notice Tests that adding a reservation reverts if the start timestamp is in the past.
    function test_AddReservationRevertsIfInvalidStartTimestamp() public {

    }

    /// @notice Tests that we can successfully pass a reservation check.
    function test_CheckReservation() public {

    }

    /// @notice Tests that a start timestamp not in the schedule period reverts.
    function test_CheckReservationRevertsIfStartTimestampNotInSchedulePeriod() public {

    }

    /// @notice Tests that an end timestamp not in the schedule period reverts.
    function test_CheckReservationRevertsIfEndTimestampNotInSchedulePeriod() public {

    }

    /// @notice Tests that start timestamp must not be greater than the end timestamp.
    function test_CheckReservationRevertsIfStartTimestampGreaterThanEndTimestamp() public {

    }

    /// @notice Tests that reservation length cannot exceed the quorum's reservation advance window
    function test_CheckReservationRevertsIfReservationTooLong() public {

    }

    /// @notice Tests that increasing a reservation's reserved symbols successfully increases the quorum's reserved symbols
    function test_IncreaseReservedSymbols() public {

    }

    /// @notice Tests that increasing a reservation's reserved symbols reverts if not enough symbols are available
    function test_IncreaseReservedSymbolsRevertsIfNotEnoughSymbolsAvailable() public {

    }

    /// @notice Tests that decreasing a reservation's reserved symbols successfully decreases the quorum's reserved symbols
    function test_DecreaseReservedSymbols() public {

    }

    /// @notice Tests that a reservation can be increased successfully.
    function test_IncreaseReservation() public {

    }

    /// @notice Tests that increasing a reservation reverts if the start timestamp does not match the reservation's start timestamp.
    function test_IncreaseReservationRevertsIfStartTimestampDoesNotMatch() public {

    }

    /// @notice Tests that increasing a reservation reverts if the reservation decreases.
    function test_IncreaseReservationRevertsIfReservationDecreases() public {

    }

    /// @notice Tests that a reservation can be decreased successfully.
    function test_DecreaseReservation() public {

    }

    /// @notice Tests that decreasing a reservation reverts if the start timestamp does not match the reservation's start timestamp.
    function test_DecreaseReservationRevertsIfStartTimestampDoesNotMatch() public {

    }

    /// @notice Tests that decreasing a reservation reverts if the reservation increases.
    function test_DecreaseReservationRevertsIfReservationIncreases() public {

    }

    /// @notice Tests that a user can deposit on demand successfully.
    function test_DepositOnDemand() public {
        
    }

    /// @notice Tests that the contract initializes correctly
    function test_Initialize() public {
    }

    /// @notice Tests that decreaseReservation and depositOnDemand use the msg sender as the target account.
    function test_DecreaseReservationAndDepositOnDemandUseMsgSender() public {
        
    }

    /// @notice Tests that these functions are properly gated to the owner.
    function testOnlyOwnerFunctions() public {

    }

    /// @notice Tests that these functions are properly gated to the quorum owner.
    function testOnlyQuorumOwnerFunctions() public {

    }

    /// @notice Tests that all the getters are there.
    function testGetters() public {
        
    }

}
