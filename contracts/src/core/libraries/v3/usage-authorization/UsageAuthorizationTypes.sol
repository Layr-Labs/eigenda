// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

library UsageAuthorizationTypes {
    /// @param deposit The total on-demand deposit of the user.
    /// @param reservation The reservation of the user.
    struct User {
        uint256 deposit;
        Reservation reservation;
    }

    /// @param protocolCfg The protocol configuration for the quorum, only settable by the contract owner.
    /// @param cfg The quorum configuration, settable by the quorum owner.
    /// @param user The mapping of users to their on-demand deposits and reservations.
    /// @param reservedSymbols The mapping of reserved symbols per period in this quorum.
    struct Quorum {
        QuorumProtocolConfig protocolCfg;
        QuorumConfig cfg;
        mapping(address => User) user;
        mapping(uint64 => uint64) reservedSymbols;
    }

    /// @param token The address of the token used for on-demand payments.
    /// @param recipient The address of the recipient of the on-demand payments.
    /// @param reservationSymbolsPerSecond The number of symbols reserved per second.
    /// @param onDemandSymbolsPerSecond The number of symbols available for on-demand payments per second.
    /// @param onDemandPricePerSymbol The price per symbol for on-demand payments.
    struct QuorumConfig {
        address token;
        address recipient;
        uint64 reservationSymbolsPerSecond;
        uint64 onDemandSymbolsPerSecond;
        uint64 onDemandPricePerSymbol;
    }

    /// @param minNumSymbols The minimum number of symbols required for the quorum.
    /// @param reservationAdvanceWindow The time window for which reservations can be advanced.
    /// @param reservationRateLimitWindow The time window for which reservations are rate-limited.
    /// @param onDemandRateLimitWindow The time window for which on-demand payments are rate-limited.
    /// @param onDemandEnabled Whether on-demand payments are enabled for the quorum.
    struct QuorumProtocolConfig {
        uint64 minNumSymbols;
        uint64 reservationAdvanceWindow;
        uint64 reservationRateLimitWindow;
        uint64 onDemandRateLimitWindow;
        bool onDemandEnabled;
    }

    /// @param symbolsPerSecond The number of symbols reserved per second.
    /// @param startTimestamp The start timestamp of the reservation.
    /// @param endTimestamp The end timestamp of the reservation.
    struct Reservation {
        uint64 symbolsPerSecond;
        uint64 startTimestamp;
        uint64 endTimestamp;
    }
}
