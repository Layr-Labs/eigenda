// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {Test} from "forge-std/Test.sol";
import {PaymentVault} from "../../src/payments/PaymentVault.sol";
import {IPaymentVault} from "../../src/interfaces/IPaymentVault.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract PaymentVaultTest is Test {
    PaymentVault private _paymentVault;
    address private _user;
    address private _user2;
    bytes private _quorumNumbers;
    bytes private _quorumSplits;

    function setUp() public {
        _paymentVault = new PaymentVault();
        _user = address(0x1);
        _user2 = address(0x2);
        _quorumNumbers = hex"0001";
        _quorumSplits = hex"3232";
    }

    function test_initialize() public {
        _paymentVault.initialize(
            address(0), // _initialOwner
            1, // _minNumSymbols
            1, // _pricePerSymbol
            1, // _priceUpdateCooldown
            1, // _globalSymbolsPerPeriod
            1, // _reservationPeriodInterval
            1, // _globalRatePeriodInterval
            1, // _maxAdvanceWindow
            1, // _maxPermissionlessReservationSymbolsPerSecond
            1, // _reservationPricePerSymbol
            1, // _reservationAdvanceWindow
            1 // _reservationSchedulePeriod
        );
    }

    function test_initialize_revert_zeroOwner() public {
        vm.expectRevert("Initial owner cannot be zero address");
        _paymentVault.initialize(
            address(0), // _initialOwner
            1, // _minNumSymbols
            1, // _pricePerSymbol
            1, // _priceUpdateCooldown
            1, // _globalSymbolsPerPeriod
            1, // _reservationPeriodInterval
            1, // _globalRatePeriodInterval
            1, // _maxAdvanceWindow
            1, // _maxPermissionlessReservationSymbolsPerSecond
            1, // _reservationPricePerSymbol
            1, // _reservationAdvanceWindow
            1 // _reservationSchedulePeriod
        );
    }

    function test_setReservation() public {
        IPaymentVault.Reservation memory reservation = IPaymentVault.Reservation({
            symbolsPerSecond: 1,
            startTimestamp: uint64(block.timestamp),
            endTimestamp: uint64(block.timestamp + 1 days),
            quorumNumber: 1
        });

        vm.deal(_user, 1 ether);
        vm.prank(_user);
        _paymentVault.setReservation{value: 1 ether}(_user, reservation);

        IPaymentVault.Reservation memory storedReservation = _paymentVault.getReservation(1, _user);
        assertEq(storedReservation.symbolsPerSecond, reservation.symbolsPerSecond);
        assertEq(storedReservation.startTimestamp, reservation.startTimestamp);
        assertEq(storedReservation.endTimestamp, reservation.endTimestamp);
        assertEq(storedReservation.quorumNumber, reservation.quorumNumber);
    }

    function test_setReservation_revert_incorrectPayment() public {
        IPaymentVault.Reservation memory reservation = IPaymentVault.Reservation({
            symbolsPerSecond: 1,
            startTimestamp: uint64(block.timestamp),
            endTimestamp: uint64(block.timestamp + 1 days),
            quorumNumber: 1
        });

        vm.deal(_user, 1 ether);
        vm.prank(_user);
        vm.expectRevert("Incorrect payment amount");
        _paymentVault.setReservation{value: 0.5 ether}(_user, reservation);
    }

    function test_setReservation_revert_notAuthorized() public {
        IPaymentVault.Reservation memory reservation = IPaymentVault.Reservation({
            symbolsPerSecond: 1,
            startTimestamp: uint64(block.timestamp),
            endTimestamp: uint64(block.timestamp + 1 days),
            quorumNumber: 1
        });

        vm.deal(_user2, 1 ether);
        vm.prank(_user2);
        vm.expectRevert("Not authorized");
        _paymentVault.setReservation{value: 1 ether}(_user, reservation);
    }

    function test_getReservations() public {
        IPaymentVault.Reservation memory reservation = IPaymentVault.Reservation({
            symbolsPerSecond: 1,
            startTimestamp: uint64(block.timestamp),
            endTimestamp: uint64(block.timestamp + 1 days),
            quorumNumber: 1
        });

        vm.deal(_user, 1 ether);
        vm.prank(_user);
        _paymentVault.setReservation{value: 1 ether}(_user, reservation);

        uint64[] memory quorums = new uint64[](1);
        quorums[0] = 1;
        address[] memory accounts = new address[](1);
        accounts[0] = _user;
        IPaymentVault.Reservation[][] memory reservations = _paymentVault.getReservations(quorums, accounts);
        
        assertEq(reservations[0][0].symbolsPerSecond, reservation.symbolsPerSecond);
        assertEq(reservations[0][0].startTimestamp, reservation.startTimestamp);
        assertEq(reservations[0][0].endTimestamp, reservation.endTimestamp);
        assertEq(reservations[0][0].quorumNumber, reservation.quorumNumber);
    }

    function test_toggleNewReservations() public {
        vm.prank(address(0)); // owner
        _paymentVault.toggleNewReservations(false);
        assertEq(_paymentVault.newReservationsEnabled(), false);

        vm.prank(address(0)); // owner
        _paymentVault.toggleNewReservations(true);
        assertEq(_paymentVault.newReservationsEnabled(), true);
    }

    function test_toggleNewReservations_revert_notOwner() public {
        vm.prank(_user);
        vm.expectRevert("Ownable: caller is not the owner");
        _paymentVault.toggleNewReservations(false);
    }

    function test_setQuorumOwner() public {
        vm.prank(address(0)); // owner
        _paymentVault.setQuorumOwner(1, _user);
        assertEq(_paymentVault.quorumOwner(1), abi.encodePacked(_user));
    }

    function test_setQuorumOwner_revert_notOwner() public {
        vm.prank(_user);
        vm.expectRevert("Ownable: caller is not the owner");
        _paymentVault.setQuorumOwner(1, _user);
    }

    function test_setPricePerSymbol() public {
        vm.prank(address(0)); // owner
        _paymentVault.setPricePerSymbol(2);
        assertEq(_paymentVault.reservationPricePerSymbol(), 2);
    }

    function test_setPricePerSymbol_revert_notOwner() public {
        vm.prank(_user);
        vm.expectRevert("Ownable: caller is not the owner");
        _paymentVault.setPricePerSymbol(2);
    }

    function test_getOnDemandTotalDeposits() public {
        vm.deal(_user, 1 ether);
        vm.deal(_user2, 2 ether);
        
        vm.prank(_user);
        _paymentVault.depositOnDemand{value: 1 ether}(_user);
        
        vm.prank(_user2);
        _paymentVault.depositOnDemand{value: 2 ether}(_user2);

        address[] memory accounts = new address[](2);
        accounts[0] = _user;
        accounts[1] = _user2;
        uint80[] memory deposits = _paymentVault.getOnDemandTotalDeposits(accounts);
        
        assertEq(deposits[0], 1 ether);
        assertEq(deposits[1], 2 ether);
    }
}