// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

library EigenDAEjectionTypes {
    /// @param ejector The address initiating the ejection
    /// @param proceedingTime Timestamp when the proceeding is set to complete
    /// @param depositAmount The amount of deposit the ejector has committed to initiating the ejection.
    /// @param quorums The quorums associated with the proceeding.
    struct EjectionRecord {
        address ejector;
        uint64 proceedingTime;
        uint256 depositAmount;
        bytes quorums;
    }

    /// @dev stateful storage entry for an ejectee - first constructed when the ejectee being targeted for ejection
    ///      hasn't been challenged before and is preserved after a cancellation for cooldown enforcements to stop
    ///      a malicious ejector from spam attacks
    ///
    /// @param record The ejection record (can be empty if previous ejection attempt was cancelled or successful).
    /// @param lastProceedingInitiated Timestamp of when the last proceeding was initiated to enforce cooldowns.
    /// @dev The parameters are separated to make the ejection record safer to delete
    struct EjecteeState {
        EjectionRecord record;
        uint64 lastProceedingInitiated;
    }
}
