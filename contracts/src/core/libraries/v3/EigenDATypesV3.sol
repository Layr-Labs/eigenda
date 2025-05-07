// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";

library EigenDATypesV3 {
    struct LockedDisperserDeposit {
        uint256 deposit;
        uint256 refund;
        address token;
        uint64 lockPeriod;
    }

    struct DisperserInfo {
        address disperser;
        bool registered;
        string disperserURL;
    }

    struct QuorumPaymentConfig {
        address token; // the address of the token used for on-demand payments.
        address recipient; // the address of the recipient of the on-demand payments.
        uint64 reservationSymbolsPerSecond;
        uint64 onDemandSymbolsPerPeriod;
        uint64 onDemandPricePerSymbol;
    }

    struct QuorumPaymentProtocolConfig {
        uint64 reservationAdvanceWindow;
    }

    struct Reservation {
        uint64 symbolsPerSecond;
        uint64 startTimestamp;
        uint64 endTimestamp;
    }
}
