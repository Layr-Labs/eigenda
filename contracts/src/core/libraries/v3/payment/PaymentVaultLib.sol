// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {SafeERC20, IERC20} from "lib/openzeppelin-contracts/contracts/token/ERC20/utils/SafeERC20.sol";
import {IPaymentVault} from "src/core/interfaces/IPaymentVault.sol";
import {PaymentVaultTypes} from "src/core/libraries/v3/payment/PaymentVaultTypes.sol";
import {PaymentVaultStorage} from "src/core/libraries/v3/payment/PaymentVaultStorage.sol";

library PaymentVaultLib {
    using SafeERC20 for IERC20;

    function s() internal pure returns (PaymentVaultStorage.Layout storage) {
        return PaymentVaultStorage.layout();
    }

    function addReservation(
        uint64 quorumId,
        address account,
        PaymentVaultTypes.Reservation memory reservation,
        uint64 schedulePeriod
    ) internal {
        checkReservation(quorumId, reservation, schedulePeriod);

        if (reservation.startTimestamp < s().quorum[quorumId].user[account].reservation.endTimestamp) {
            revert IPaymentVault.ReservationStillActive(s().quorum[quorumId].user[account].reservation.endTimestamp);
        }
        if (reservation.startTimestamp < block.timestamp) {
            revert IPaymentVault.InvalidStartTimestamp(uint64(block.timestamp));
        }

        increaseReservedSymbols(
            quorumId, reservation.startTimestamp, reservation.endTimestamp, reservation.symbolsPerSecond, schedulePeriod
        );
        s().quorum[quorumId].user[account].reservation = reservation;
    }

    function increaseReservation(
        uint64 quorumId,
        address account,
        PaymentVaultTypes.Reservation memory reservation,
        uint64 schedulePeriod
    ) internal {
        PaymentVaultTypes.Reservation storage currentReservation = s().quorum[quorumId].user[account].reservation;
        checkReservation(quorumId, reservation, schedulePeriod);

        if (reservation.startTimestamp != currentReservation.startTimestamp) {
            revert IPaymentVault.StartTimestampMustMatch(currentReservation.startTimestamp);
        }
        if (
            reservation.endTimestamp < currentReservation.endTimestamp
                || reservation.symbolsPerSecond < currentReservation.symbolsPerSecond
        ) {
            revert IPaymentVault.ReservationMustIncrease(
                currentReservation.endTimestamp, currentReservation.symbolsPerSecond
            );
        }

        if (reservation.symbolsPerSecond > currentReservation.symbolsPerSecond) {
            // increase reservation symbols for the current reservation time.
            increaseReservedSymbols(
                quorumId,
                currentReservation.startTimestamp,
                currentReservation.endTimestamp,
                reservation.symbolsPerSecond - currentReservation.symbolsPerSecond,
                schedulePeriod
            );
        }
        if (reservation.endTimestamp > currentReservation.endTimestamp) {
            // increase reservation time with new symbols per second
            increaseReservedSymbols(
                quorumId,
                currentReservation.endTimestamp,
                reservation.endTimestamp,
                reservation.symbolsPerSecond,
                schedulePeriod
            );
        }
        s().quorum[quorumId].user[account].reservation = reservation;
    }

    function decreaseReservation(
        uint64 quorumId,
        address account,
        PaymentVaultTypes.Reservation memory reservation,
        uint64 schedulePeriod
    ) internal {
        PaymentVaultTypes.Reservation storage currentReservation = s().quorum[quorumId].user[account].reservation;
        checkReservation(quorumId, reservation, schedulePeriod);

        if (reservation.startTimestamp != currentReservation.startTimestamp) {
            revert IPaymentVault.StartTimestampMustMatch(currentReservation.startTimestamp);
        }

        if (
            reservation.endTimestamp > currentReservation.endTimestamp
                || reservation.symbolsPerSecond > currentReservation.symbolsPerSecond
        ) {
            revert IPaymentVault.ReservationMustDecrease(
                currentReservation.endTimestamp, currentReservation.symbolsPerSecond
            );
        }

        if (reservation.endTimestamp < currentReservation.endTimestamp) {
            // decrease reservation time
            decreaseReservedSymbols(
                quorumId,
                reservation.endTimestamp,
                currentReservation.endTimestamp,
                currentReservation.symbolsPerSecond,
                schedulePeriod
            );
        }
        if (reservation.symbolsPerSecond < currentReservation.symbolsPerSecond) {
            // decrease reservation symbols for the remaining reservation time
            decreaseReservedSymbols(
                quorumId,
                currentReservation.startTimestamp,
                currentReservation.endTimestamp,
                currentReservation.symbolsPerSecond - reservation.symbolsPerSecond,
                schedulePeriod
            );
        }
        s().quorum[quorumId].user[account].reservation = reservation;
    }

    /// @notice Does required checks on a reservation, and returns the starting timestamp for accounting for additional bandwidth.
    function checkReservation(uint64 quorumId, PaymentVaultTypes.Reservation memory reservation, uint64 schedulePeriod)
        internal
        view
    {
        PaymentVaultTypes.Quorum storage quorum = s().quorum[quorumId];
        if (reservation.startTimestamp % schedulePeriod != 0) {
            revert IPaymentVault.TimestampSchedulePeriodMismatch(reservation.startTimestamp, schedulePeriod);
        }
        if (reservation.endTimestamp % schedulePeriod != 0) {
            revert IPaymentVault.TimestampSchedulePeriodMismatch(reservation.endTimestamp, schedulePeriod);
        }
        if (reservation.startTimestamp <= reservation.endTimestamp) {
            revert IPaymentVault.InvalidStartTimestamp(uint64(block.timestamp));
        }
        if (reservation.endTimestamp - reservation.startTimestamp > quorum.protocolCfg.reservationAdvanceWindow) {
            revert IPaymentVault.ReservationTooLong(
                reservation.endTimestamp - reservation.startTimestamp, quorum.protocolCfg.reservationAdvanceWindow
            );
        }
    }

    function depositOnDemand(uint64 quorumId, address account, uint256 amount) internal {
        PaymentVaultTypes.Quorum storage quorum = s().quorum[quorumId];
        PaymentVaultTypes.User storage user = quorum.user[account];
        PaymentVaultTypes.QuorumConfig storage cfg = quorum.cfg;

        uint256 newAmount = user.deposit + amount;
        if (newAmount > type(uint80).max) {
            revert IPaymentVault.AmountTooLarge(newAmount, type(uint80).max);
        }

        IERC20(cfg.token).safeTransferFrom(account, cfg.recipient, amount);

        user.deposit = newAmount;
    }

    function increaseReservedSymbols(
        uint64 quorumId,
        uint64 startTimestamp,
        uint64 endTimestamp,
        uint64 symbolsPerSecond,
        uint64 schedulePeriod
    ) internal {
        if (startTimestamp % schedulePeriod != 0) {
            revert IPaymentVault.TimestampSchedulePeriodMismatch(startTimestamp, schedulePeriod);
        }
        if (endTimestamp % schedulePeriod != 0) {
            revert IPaymentVault.TimestampSchedulePeriodMismatch(endTimestamp, schedulePeriod);
        }
        if (endTimestamp <= startTimestamp) {
            revert IPaymentVault.InvalidReservationPeriod(startTimestamp, endTimestamp);
        }
        uint64 startPeriod = startTimestamp / schedulePeriod;
        uint64 endPeriod = endTimestamp / schedulePeriod;

        PaymentVaultTypes.Quorum storage quorum = s().quorum[quorumId];
        uint64 maxReservedSymbols = quorum.cfg.reservationSymbolsPerSecond;
        for (uint64 i = startPeriod; i < endPeriod; i++) {
            uint64 reservedSymbols = quorum.reservedSymbols[i] + symbolsPerSecond;
            require(reservedSymbols <= maxReservedSymbols, "Not enough symbols available");
            quorum.reservedSymbols[i] = reservedSymbols;
        }
    }

    function decreaseReservedSymbols(
        uint64 quorumId,
        uint64 startTimestamp,
        uint64 endTimestamp,
        uint64 symbolsPerSecond,
        uint64 schedulePeriod
    ) internal {
        if (startTimestamp % schedulePeriod != 0) {
            revert IPaymentVault.TimestampSchedulePeriodMismatch(startTimestamp, schedulePeriod);
        }
        if (endTimestamp % schedulePeriod != 0) {
            revert IPaymentVault.TimestampSchedulePeriodMismatch(endTimestamp, schedulePeriod);
        }
        if (endTimestamp <= startTimestamp) {
            revert IPaymentVault.InvalidReservationPeriod(startTimestamp, endTimestamp);
        }
        uint64 startPeriod = startTimestamp / schedulePeriod;
        uint64 endPeriod = endTimestamp / schedulePeriod;

        PaymentVaultTypes.Quorum storage quorum = s().quorum[quorumId];
        for (uint64 i = startPeriod; i < endPeriod; i++) {
            quorum.reservedSymbols[i] -= symbolsPerSecond; // Revert on underflow
        }
    }
}
