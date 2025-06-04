// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

library PaymentVaultTypes {
    struct User {
        uint256 deposit; // the total on demand deposit of the user
        Reservation reservation;
    }

    struct Quorum {
        QuorumProtocolConfig protocolCfg;
        QuorumConfig cfg;
        mapping(address => User) user;
        mapping(uint64 => uint64) reservedSymbols; // reserved symbols per period in this quorum
    }

    struct QuorumConfig {
        address token; // the address of the token used for on-demand payments.
        address recipient; // the address of the recipient of the on-demand payments.
        uint64 reservationSymbolsPerSecond;
        uint64 onDemandSymbolsPerSecond;
        uint64 onDemandPricePerSymbol;
    }

    struct QuorumProtocolConfig {
        uint64 minNumSymbols;
        uint64 reservationAdvanceWindow;
        uint64 reservationRateLimitWindow;
        uint64 onDemandRateLimitWindow;
        bool onDemandEnabled;
    }

    struct Reservation {
        uint64 symbolsPerSecond;
        uint64 startTimestamp;
        uint64 endTimestamp;
    }
}
