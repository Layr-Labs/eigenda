// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

library EigenDAEjectionTypes {
    /// @param ejector The address initiating the ejection
    /// @param proceedingTime Timestamp when the proceeding is set to complete
    /// @param lastProceedingInitiated Timestamp of when the last proceeding was initiated to enforce cooldowns
    /// @param depositAmount The amount of deposit the ejector has commmitted to initiating the ejection.
    /// @param quorums The quorums associated with the proceeding.
    struct EjectionParams {
        address ejector;
        uint64 proceedingTime;
        uint256 depositAmount;
        bytes quorums;
    }

    /// @param params The ejection parameters
    /// @param lastProceedingInitiated Timestamp of when the last proceeding was initiated to enforce cooldowns.
    /// @dev The parameters are separated to make the ejection parameters safer to delete.
    struct Ejectee {
        EjectionParams params;
        uint64 lastProceedingInitiated;
    }
}
