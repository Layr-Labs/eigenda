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
        EigenDAEjectionTypes.EjecteeState storage ejectee = getEjectee(operator);

        require(ejectee.record.proceedingTime == 0, "Ejection already in progress");
        require(ejectee.lastProceedingInitiated + s().cooldown <= block.timestamp, "Ejection cooldown not met");

        ejectee.record.ejector = ejector;
        ejectee.record.quorums = quorums;
        ejectee.record.proceedingTime = uint64(block.timestamp) + s().delay;
        ejectee.record.depositAmount = depositAmount;
        ejectee.lastProceedingInitiated = uint64(block.timestamp);
        emit EjectionStarted(
            operator, ejector, quorums, ejectee.lastProceedingInitiated, ejectee.record.proceedingTime, depositAmount
        );
    }

    /// @notice Cancels an ejection process for an operator.
    function cancelEjection(address operator) internal {
        EigenDAEjectionTypes.EjecteeState storage ejectee = getEjectee(operator);
        require(ejectee.record.proceedingTime > 0, "No ejection in progress");

        deleteEjectionRecord(operator);
        emit EjectionCancelled(operator);
    }

    /// @notice Completes an ejection process for an operator.
    function completeEjection(address operator, bytes memory quorums) internal {
        require(quorumsEqual(s().ejectees[operator].record.quorums, quorums), "Quorums do not match");
        EigenDAEjectionTypes.EjecteeState storage ejectee = s().ejectees[operator];
        require(ejectee.record.proceedingTime > 0, "No proceeding in progress");

        require(block.timestamp >= ejectee.record.proceedingTime, "Proceeding not yet due");

        deleteEjectionRecord(operator);
        emit EjectionCompleted(operator, quorums);
    }

    /// @notice Helper function to clear an ejection.
    /// @dev The lastProceedingInitiated field is not cleared to allow cooldown enforcement.
    function deleteEjectionRecord(address operator) internal {
        EigenDAEjectionTypes.EjecteeState storage ejectee = s().ejectees[operator];
        ejectee.record.ejector = address(0);
        ejectee.record.quorums = hex"";
        ejectee.record.proceedingTime = 0;
        ejectee.record.depositAmount = 0;
    }

    /// @notice Adds to the ejector's balance for ejection processes.
    /// @dev This function does not handle tokens
    function addEjectorBalance(address ejector, uint256 amount) internal {
        s().ejectorBalanceRecord[ejector] += amount;
    }

    /// @notice Subtracts from the ejector's balance for ejection processes.
    /// @dev This function does not handle tokens
    function subtractEjectorBalance(address ejector, uint256 amount) internal {
        require(s().ejectorBalanceRecord[ejector] >= amount, "Insufficient balance");
        // no underflow check needed
        unchecked {
            s().ejectorBalanceRecord[ejector] -= amount;
        }
    }

    /// @notice Returns the address of the ejector for a given operator.
    /// @dev If the address is zero, it means no ejection is in progress.
    function getEjector(address operator) internal view returns (address ejector) {
        return s().ejectees[operator].record.ejector;
    }

    function getEjectorBalance(address ejector) internal view returns (uint256) {
        return s().ejectorBalanceRecord[ejector];
    }

    function getDepositAmount(address operator) internal view returns (uint256 depositAmount) {
        return s().ejectees[operator].record.depositAmount;
    }

    function lastProceedingInitiated(address operator) internal view returns (uint64) {
        return s().ejectees[operator].lastProceedingInitiated;
    }

    /// @notice Compares two quorums to see if they are equal.
    function quorumsEqual(bytes memory quorums1, bytes memory quorums2) internal pure returns (bool) {
        return keccak256(quorums1) == keccak256(quorums2);
    }

    function getEjectee(address operator) internal view returns (EigenDAEjectionTypes.EjecteeState storage) {
        return s().ejectees[operator];
    }

    function getEjectionRecord(address operator) internal view returns (EigenDAEjectionTypes.EjectionRecord storage) {
        return s().ejectees[operator].record;
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
