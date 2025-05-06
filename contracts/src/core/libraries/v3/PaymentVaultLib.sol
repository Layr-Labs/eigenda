// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV3} from "src/core/libraries/v3/EigenDATypesV3.sol";
import {SafeERC20, IERC20} from "lib/openzeppelin-contracts/contracts/token/ERC20/utils/SafeERC20.sol";

library PaymentVaultStorage {
    string internal constant STORAGE_ID = "eigen.da.payment.vault";
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

library PaymentVaultLib {
    using SafeERC20 for IERC20;

    event ReservationCreated(uint64 indexed quorumId, address indexed account, EigenDATypesV3.Reservation reservation);

    function s() internal pure returns (PaymentVaultStorage.Layout storage) {
        return PaymentVaultStorage.layout();
    }

    function addReservation(
        uint64 quorumId,
        address account,
        EigenDATypesV3.Reservation memory reservation,
        uint64 schedulePeriod
    ) internal {
        checkReservation(quorumId, account, reservation, schedulePeriod);

        require(
            reservation.startTimestamp >= s().quorum[quorumId].user[account].reservation.endTimestamp,
            "Invalid start timestamp"
        );
        require(reservation.startTimestamp >= block.timestamp, "Invalid start timestamp");

        populateReservedSymbols(
            quorumId, reservation.startTimestamp, reservation.endTimestamp, reservation.symbolsPerSecond, schedulePeriod
        );
        s().quorum[quorumId].user[account].reservation = reservation;
    }

    function increaseReservation(
        uint64 quorumId,
        address account,
        EigenDATypesV3.Reservation memory reservation,
        uint64 schedulePeriod
    ) internal {
        EigenDATypesV3.Reservation storage currentReservation = s().quorum[quorumId].user[account].reservation;
        checkReservation(quorumId, account, reservation, schedulePeriod);

        require(reservation.startTimestamp == currentReservation.startTimestamp, "Invalid start timestamp");
        require(reservation.endTimestamp > currentReservation.endTimestamp, "Invalid end timestamp");
        require(reservation.symbolsPerSecond >= currentReservation.symbolsPerSecond, "Invalid symbols per second");

        populateReservedSymbols(
            quorumId,
            currentReservation.endTimestamp,
            reservation.endTimestamp,
            reservation.symbolsPerSecond,
            schedulePeriod
        );
        s().quorum[quorumId].user[account].reservation = reservation;
    }

    function decreaseReservation(
        uint64 quorumId,
        address account,
        EigenDATypesV3.Reservation memory reservation,
        uint64 schedulePeriod
    ) internal {
        EigenDATypesV3.Reservation storage currentReservation = s().quorum[quorumId].user[account].reservation;
        checkReservation(quorumId, account, reservation, schedulePeriod);

        require(reservation.startTimestamp == currentReservation.startTimestamp, "Invalid start timestamp");
        require(reservation.endTimestamp <= currentReservation.endTimestamp, "Invalid end timestamp");
        require(reservation.symbolsPerSecond <= currentReservation.symbolsPerSecond, "Invalid symbols per second");

        populateReservedSymbols(
            quorumId,
            currentReservation.endTimestamp,
            reservation.endTimestamp,
            reservation.symbolsPerSecond,
            schedulePeriod
        );
        s().quorum[quorumId].user[account].reservation = reservation;
    }

    /// @notice Does required checks on a reservation, and returns the starting timestamp for accounting for additional bandwidth.
    function checkReservation(
        uint64 quorumId,
        address account,
        EigenDATypesV3.Reservation memory reservation,
        uint64 schedulePeriod
    ) internal view returns (uint64 startTimestamp) {
        PaymentVaultStorage.Quorum storage quorum = s().quorum[quorumId];
        PaymentVaultStorage.User storage user = quorum.user[account];
        require(reservation.startTimestamp % schedulePeriod == 0, "Invalid start timestamp");
        require(reservation.endTimestamp % schedulePeriod == 0, "Invalid end timestamp");
        require(reservation.endTimestamp > reservation.startTimestamp, "Invalid reservation period");
        require(
            reservation.endTimestamp - reservation.startTimestamp <= quorum.protocolCfg.reservationAdvanceWindow,
            "Reservation period too long"
        );
        // If the reservation is not expired, the reservation can only be updated favorably to the user.
        if (block.timestamp <= reservation.endTimestamp) {
            EigenDATypesV3.Reservation memory currentReservation = user.reservation;
            require(
                reservation.startTimestamp == currentReservation.startTimestamp
                    && reservation.endTimestamp > currentReservation.endTimestamp,
                "Invalid reservation update"
            );
            require(reservation.symbolsPerSecond >= currentReservation.symbolsPerSecond, "Invalid symbols per second");
            return currentReservation.endTimestamp;
        }
        return reservation.startTimestamp;
    }

    function depositOnDemand(uint64 quorumId, address account, uint256 amount) internal {
        PaymentVaultStorage.Quorum storage quorum = s().quorum[quorumId];
        PaymentVaultStorage.User storage user = quorum.user[account];
        EigenDATypesV3.QuorumPaymentConfig storage cfg = quorum.cfg;

        uint256 newAmount = user.deposit + amount;
        require(newAmount <= type(uint80).max, "Amount too large");

        IERC20(cfg.token).safeTransferFrom(account, cfg.recipient, amount);

        user.deposit = newAmount;
    }

    function populateReservedSymbols(
        uint64 quorumId,
        uint64 startTimestamp,
        uint64 endTimestamp,
        uint64 symbolsPerSecond,
        uint64 schedulePeriod
    ) internal {
        require(startTimestamp % schedulePeriod == 0, "Invalid start timestamp");
        require(endTimestamp % schedulePeriod == 0, "Invalid end timestamp");
        uint64 startPeriod = startTimestamp / schedulePeriod;
        uint64 endPeriod = endTimestamp / schedulePeriod;

        PaymentVaultStorage.Quorum storage quorum = s().quorum[quorumId];
        uint64 maxReservedSymbols = quorum.cfg.reservationSymbolsPerSecond;
        for (uint64 i = startPeriod; i < endPeriod; i++) {
            uint64 reservedSymbols = quorum.reservedSymbols[i] + symbolsPerSecond;
            require(reservedSymbols <= maxReservedSymbols, "Not enough symbols available");
            quorum.reservedSymbols[i] = reservedSymbols;
        }
    }
}
