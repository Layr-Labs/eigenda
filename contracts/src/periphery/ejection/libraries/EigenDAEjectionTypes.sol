// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

library EigenDAEjectionTypes {
    /// @param proceedingTime Timestamp when the proceeding is set to complete
    /// @param lastProceedingInitiated Timestamp of when the last proceeding was initiated to enforce cooldowns
    /// @param quorums The quorums associated with the proceeding
    struct EjectionParams {
        uint64 proceedingTime;
        uint64 lastProceedingInitiated;
        bytes quorums;
    }
}