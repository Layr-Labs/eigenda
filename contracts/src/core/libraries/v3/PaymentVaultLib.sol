// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV3} from "src/core/libraries/v3/EigenDATypesV3.sol";

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
        uint64 schedulePeriod;
    }

    function layout() internal pure returns (Layout storage s) {
        bytes32 position = STORAGE_POSITION;
        assembly {
            s.slot := position
        }
    }
}

library PaymentVaultLib {
    event ReservationCreated(uint64 indexed quorumId, address indexed account, EigenDATypesV3.Reservation reservation);

    function s() internal pure returns (PaymentVaultStorage.Layout storage) {
        return PaymentVaultStorage.layout();
    }

    function checkReservation(
        uint64 quorumId,
        address account,
        EigenDATypesV3.Reservation memory reservation,
        uint64 schedulePeriod
    ) internal view {
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
        }
    }

    function createReservation(
        uint64 quorumId,
        address account,
        EigenDATypesV3.Reservation memory reservation,
        uint64 schedulePeriod
    ) internal {
        checkReservation(quorumId, account, reservation, schedulePeriod);
        uint64 startPeriod = reservation.startTimestamp / schedulePeriod;
        uint64 endPeriod = reservation.endTimestamp / schedulePeriod;

        populateReservedSymbols(quorumId, startPeriod, endPeriod, reservation.symbolsPerSecond);
        s().quorum[quorumId].user[account].reservation = reservation;
    }

    function populateReservedSymbols(uint64 quorumId, uint64 startPeriod, uint64 endPeriod, uint64 symbolsPerSecond)
        internal
    {
        PaymentVaultStorage.Quorum storage quorum = s().quorum[quorumId];
        uint64 maxReservedSymbols = quorum.cfg.reservationSymbolsPerSecond;
        for (uint64 i = startPeriod; i < endPeriod; i++) {
            uint64 reservedSymbols = quorum.reservedSymbols[i] + symbolsPerSecond;
            require(reservedSymbols <= maxReservedSymbols, "Not enough symbols available");
            quorum.reservedSymbols[i] = reservedSymbols;
        }
    }
}
