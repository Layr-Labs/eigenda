// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {SafeERC20, IERC20} from "lib/openzeppelin-contracts/contracts/token/ERC20/utils/SafeERC20.sol";
import {IUsageAuthorizationRegistry} from "src/core/interfaces/IUsageAuthorizationRegistry.sol";
import {UsageAuthorizationTypes} from "src/core/libraries/v3/usage-authorization/UsageAuthorizationTypes.sol";
import {UsageAuthorizationStorage} from "src/core/libraries/v3/usage-authorization/UsageAuthorizationStorage.sol";

library UsageAuthorizationLib {
    using SafeERC20 for IERC20;

    event ReservationAdded(
        uint64 indexed quorumId, address indexed account, UsageAuthorizationTypes.Reservation reservation
    );

    event ReservationIncreased(
        uint64 indexed quorumId, address indexed account, UsageAuthorizationTypes.Reservation reservation
    );

    event ReservationDecreased(
        uint64 indexed quorumId, address indexed account, UsageAuthorizationTypes.Reservation reservation
    );

    event DepositOnDemand(uint64 indexed quorumId, address indexed account, uint256 amount, address indexed payer);

    function s() internal pure returns (UsageAuthorizationStorage.Layout storage) {
        return UsageAuthorizationStorage.layout();
    }

    /// @notice Adds a reservation for a user in a quorum. Requires that any previous reservation has ended, and the start of the new reservation is in the future.
    function addReservation(
        uint64 quorumId,
        address account,
        UsageAuthorizationTypes.Reservation memory reservation,
        uint64 schedulePeriod
    ) internal {
        checkReservation(quorumId, reservation, schedulePeriod);

        if (reservation.startTimestamp < s().quorum[quorumId].user[account].reservation.endTimestamp) {
            revert IUsageAuthorizationRegistry.ReservationStillActive(
                s().quorum[quorumId].user[account].reservation.endTimestamp
            );
        }

        increaseReservedSymbols(
            quorumId, reservation.startTimestamp, reservation.endTimestamp, reservation.symbolsPerSecond, schedulePeriod
        );
        s().quorum[quorumId].user[account].reservation = reservation;
        emit ReservationAdded(quorumId, account, reservation);
    }

    /// @notice Updates a reservation for a user in a quorum. Requires that the start timestamp matches the current reservation,
    ///         and that the new reservation is either the same or an increase over the current reservation.
    function increaseReservation(
        uint64 quorumId,
        address account,
        UsageAuthorizationTypes.Reservation memory reservation,
        uint64 schedulePeriod
    ) internal {
        checkReservation(quorumId, reservation, schedulePeriod);

        UsageAuthorizationTypes.Reservation storage currentReservation = s().quorum[quorumId].user[account].reservation;
        if (reservation.startTimestamp != currentReservation.startTimestamp) {
            revert IUsageAuthorizationRegistry.StartTimestampMustMatch(currentReservation.startTimestamp);
        }
        if (
            reservation.endTimestamp < currentReservation.endTimestamp
                || reservation.symbolsPerSecond < currentReservation.symbolsPerSecond
        ) {
            revert IUsageAuthorizationRegistry.ReservationMustIncrease(
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
        emit ReservationIncreased(quorumId, account, reservation);
    }

    /// @notice Updates a reservation for a user in a quorum. Requires that the start timestamp matches the current reservation,
    ///         and that the new reservation is either the same or a decrease over the current reservation.
    function decreaseReservation(
        uint64 quorumId,
        address account,
        UsageAuthorizationTypes.Reservation memory reservation,
        uint64 schedulePeriod
    ) internal {
        UsageAuthorizationTypes.Reservation storage currentReservation = s().quorum[quorumId].user[account].reservation;
        checkReservation(quorumId, reservation, schedulePeriod);

        if (reservation.startTimestamp != currentReservation.startTimestamp) {
            revert IUsageAuthorizationRegistry.StartTimestampMustMatch(currentReservation.startTimestamp);
        }

        if (
            reservation.endTimestamp > currentReservation.endTimestamp
                || reservation.symbolsPerSecond > currentReservation.symbolsPerSecond
        ) {
            revert IUsageAuthorizationRegistry.ReservationMustDecrease(
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
        emit ReservationDecreased(quorumId, account, reservation);
    }

    /// @notice Does required checks on a reservation
    function checkReservation(
        uint64 quorumId,
        UsageAuthorizationTypes.Reservation memory reservation,
        uint64 schedulePeriod
    ) internal view {
        UsageAuthorizationTypes.Quorum storage quorum = s().quorum[quorumId];
        uint64 roundedCurrentTimestamp = uint64(block.timestamp / schedulePeriod) * schedulePeriod;
        if (reservation.startTimestamp % schedulePeriod != 0) {
            revert IUsageAuthorizationRegistry.TimestampSchedulePeriodMismatch(
                reservation.startTimestamp, schedulePeriod
            );
        }
        if (reservation.endTimestamp % schedulePeriod != 0) {
            revert IUsageAuthorizationRegistry.TimestampSchedulePeriodMismatch(reservation.endTimestamp, schedulePeriod);
        }
        if (reservation.startTimestamp >= reservation.endTimestamp) {
            revert IUsageAuthorizationRegistry.InvalidReservationPeriod(
                reservation.startTimestamp, reservation.endTimestamp
            );
        }
        if (reservation.endTimestamp - roundedCurrentTimestamp > quorum.protocolCfg.reservationAdvanceWindow) {
            revert IUsageAuthorizationRegistry.ReservationTooLong(
                uint64(reservation.endTimestamp - roundedCurrentTimestamp), quorum.protocolCfg.reservationAdvanceWindow
            );
        }
    }

    /// @notice Deposits an amount on-demand for a user in a quorum. Requires that the amount does not exceed the maximum allowed deposit.
    function depositOnDemand(uint64 quorumId, address account, uint256 amount, address payer) internal {
        UsageAuthorizationTypes.Quorum storage quorum = s().quorum[quorumId];
        UsageAuthorizationTypes.User storage user = quorum.user[account];
        UsageAuthorizationTypes.QuorumConfig storage cfg = quorum.cfg;

        uint256 newAmount = user.deposit + amount;
        if (newAmount > type(uint80).max) {
            revert IUsageAuthorizationRegistry.AmountTooLarge(newAmount, type(uint80).max);
        }

        IERC20(cfg.token).safeTransferFrom(payer, cfg.recipient, amount);

        user.deposit = newAmount;
        emit DepositOnDemand(quorumId, account, amount, payer);
    }

    /// @notice Increases the reserved symbols for a quorum in a given period. Requires that the start and end timestamps are multiples of the schedule period.
    /// @dev Assumes that the timestamps are already checked to align with the schedule period.
    function increaseReservedSymbols(
        uint64 quorumId,
        uint64 startTimestamp,
        uint64 endTimestamp,
        uint64 symbolsPerSecond,
        uint64 schedulePeriod
    ) internal {
        uint64 startPeriod = startTimestamp / schedulePeriod;
        uint64 endPeriod = endTimestamp / schedulePeriod;

        UsageAuthorizationTypes.Quorum storage quorum = s().quorum[quorumId];
        uint64 maxReservedSymbols = quorum.cfg.reservationSymbolsPerSecond;
        for (uint64 i = startPeriod; i < endPeriod; i++) {
            uint64 reservedSymbols = quorum.reservedSymbols[i] + symbolsPerSecond;
            if (reservedSymbols > maxReservedSymbols) {
                revert IUsageAuthorizationRegistry.NotEnoughSymbolsAvailable(
                    i * schedulePeriod, reservedSymbols, maxReservedSymbols
                );
            }
            quorum.reservedSymbols[i] = reservedSymbols;
        }
    }

    /// @notice Decreases the reserved symbols for a quorum in a given period. Requires that the start and end timestamps are multiples of the schedule period.
    /// @dev Assumes that the timestamps are already checked to align with the schedule period.
    function decreaseReservedSymbols(
        uint64 quorumId,
        uint64 startTimestamp,
        uint64 endTimestamp,
        uint64 symbolsPerSecond,
        uint64 schedulePeriod
    ) internal {
        uint64 startPeriod = startTimestamp / schedulePeriod;
        uint64 endPeriod = endTimestamp / schedulePeriod;

        UsageAuthorizationTypes.Quorum storage quorum = s().quorum[quorumId];
        for (uint64 i = startPeriod; i < endPeriod; i++) {
            quorum.reservedSymbols[i] -= symbolsPerSecond; // Revert on underflow
        }
    }
}
