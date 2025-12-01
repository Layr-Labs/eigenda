// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDAEjectionTypes} from "./EigenDAEjectionTypes.sol";
import {EigenDAEjectionStorage} from "./EigenDAEjectionStorage.sol";

library EigenDAEjectionLib {
    event EjectionStarted(
        address operator,
        address ejector,
        bytes quorums,
        uint64 timestampStarted,
        uint64 ejectionTime,
        uint256 depositAmount
    );

    event EjectionCancelled(address operator);

    event EjectionCompleted(address operator, bytes quorums);

    event DelaySet(uint64 delay);

    event CooldownSet(uint64 cooldown);

    /// @notice Sets the delay for ejection processes.
    function setDelay(uint64 delay) internal {
        s().delay = delay;
        emit DelaySet(delay);
    }

    /// @notice Sets the cooldown for ejection processes.
    function setCooldown(uint64 cooldown) internal {
        s().cooldown = cooldown;
        emit CooldownSet(cooldown);
    }

    /// @notice Starts an ejection process for an operator.
    function startEjection(address operator, address ejector, bytes memory quorums, uint256 depositAmount) internal {
        EigenDAEjectionTypes.Ejectee storage ejectee = s().ejectionParams[operator];

        require(ejectee.params.proceedingTime == 0, "Ejection already in progress");
        require(ejectee.lastProceedingInitiated + s().cooldown <= block.timestamp, "Ejection cooldown not met");

        ejectee.params.ejector = ejector;
        ejectee.params.quorums = quorums;
        ejectee.params.proceedingTime = uint64(block.timestamp) + s().delay;
        ejectee.params.depositAmount = depositAmount;
        ejectee.lastProceedingInitiated = uint64(block.timestamp);
        emit EjectionStarted(
            operator, ejector, quorums, ejectee.lastProceedingInitiated, ejectee.params.proceedingTime, depositAmount
        );
    }

    /// @notice Cancels an ejection process for an operator.
    function cancelEjection(address operator) internal {
        EigenDAEjectionTypes.Ejectee storage ejectee = s().ejectionParams[operator];
        require(ejectee.params.proceedingTime > 0, "No ejection in progress");

        deleteEjection(operator);
        emit EjectionCancelled(operator);
    }

    /// @notice Completes an ejection process for an operator.
    function completeEjection(address operator, bytes memory quorums) internal {
        require(quorumsEqual(s().ejectionParams[operator].params.quorums, quorums), "Quorums do not match");
        EigenDAEjectionTypes.Ejectee storage ejectee = s().ejectionParams[operator];
        require(ejectee.params.proceedingTime > 0, "No proceeding in progress");

        require(block.timestamp >= ejectee.params.proceedingTime, "Proceeding not yet due");

        deleteEjection(operator);
        emit EjectionCompleted(operator, quorums);
    }

    /// @notice Helper function to clear an ejection.
    /// @dev The lastProceedingInitiated field is not cleared to allow cooldown enforcement.
    function deleteEjection(address operator) internal {
        EigenDAEjectionTypes.Ejectee storage ejectee = s().ejectionParams[operator];
        ejectee.params.ejector = address(0);
        ejectee.params.quorums = hex"";
        ejectee.params.proceedingTime = 0;
        ejectee.params.depositAmount = 0;
    }

    /// @notice Adds to the ejector's balance for ejection processes.
    /// @dev This function does not handle tokens
    function addEjectorBalance(address ejector, uint256 amount) internal {
        s().ejectorBalance[ejector] += amount;
    }

    /// @notice Subtracts from the ejector's balance for ejection processes.
    /// @dev This function does not handle tokens
    function subtractEjectorBalance(address ejector, uint256 amount) internal {
        require(s().ejectorBalance[ejector] >= amount, "Insufficient balance");
        // no underflow check needed
        unchecked {
            s().ejectorBalance[ejector] -= amount;
        }
    }

    /// @notice Returns the address of the ejector for a given operator.
    /// @dev If the address is zero, it means no ejection is in progress.
    function getEjector(address operator) internal view returns (address ejector) {
        return s().ejectionParams[operator].params.ejector;
    }

    function getDepositAmount(address operator) internal view returns (uint256 depositAmount) {
        return s().ejectionParams[operator].params.depositAmount;
    }

    function lastProceedingInitiated(address operator) internal view returns (uint64) {
        return s().ejectionParams[operator].lastProceedingInitiated;
    }

    /// @notice Compares two quorums to see if they are equal.
    function quorumsEqual(bytes memory quorums1, bytes memory quorums2) internal pure returns (bool) {
        return keccak256(quorums1) == keccak256(quorums2);
    }

    function ejectionParams(address operator) internal view returns (EigenDAEjectionTypes.EjectionParams storage) {
        return s().ejectionParams[operator].params;
    }

    /// @return The amount of time that must elapse from initialization before an ejection can be completed.
    function getDelay() internal view returns (uint64) {
        return s().delay;
    }

    /// @return The amount of time that must elapse after an ejection is initiated before another can be initiated for an operator.
    function getCooldown() internal view returns (uint64) {
        return s().cooldown;
    }

    /// @notice Returns the ejection storage.
    function s() private pure returns (EigenDAEjectionStorage.Layout storage) {
        return EigenDAEjectionStorage.layout();
    }
}
