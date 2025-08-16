// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDAEjectionTypes} from "src/periphery/ejection/libraries/EigenDAEjectionTypes.sol";
import {EigenDAEjectionStorage} from "src/periphery/ejection/libraries/EigenDAEjectionStorage.sol";

library EigenDAEjectionLib {
    event EjectionStarted(
        address operator, address ejector, bytes quorums, uint64 timestampStarted, uint64 ejectionTime
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
    function startEjection(address operator, address ejector, bytes memory quorums) internal {
        EigenDAEjectionTypes.EjectionParams storage operatorParams = s().ejectionParams[operator];

        require(operatorParams.proceedingTime == 0, "Ejection already in progress");
        require(operatorParams.lastProceedingInitiated + s().cooldown <= block.timestamp, "Ejection cooldown not met");

        operatorParams.ejector = ejector;
        operatorParams.quorums = quorums;
        operatorParams.proceedingTime = uint64(block.timestamp) + s().delay;
        operatorParams.lastProceedingInitiated = uint64(block.timestamp);
        emit EjectionStarted(
            operator, ejector, quorums, operatorParams.lastProceedingInitiated, operatorParams.proceedingTime
        );
    }

    /// @notice Cancels an ejection process for an operator.
    function cancelEjection(address operator) internal {
        EigenDAEjectionTypes.EjectionParams storage operatorParams = s().ejectionParams[operator];
        require(operatorParams.proceedingTime > 0, "No ejection in progress");

        deleteEjection(operator);
        emit EjectionCancelled(operator);
    }

    /// @notice Completes an ejection process for an operator.
    function completeEjection(address operator, bytes memory quorums) internal {
        require(quorumsEqual(s().ejectionParams[operator].quorums, quorums), "Quorums do not match");
        EigenDAEjectionTypes.EjectionParams storage operatorParams = s().ejectionParams[operator];
        require(operatorParams.proceedingTime > 0, "No proceeding in progress");

        require(block.timestamp >= operatorParams.proceedingTime, "Proceeding not yet due");

        deleteEjection(operator);
        emit EjectionCompleted(operator, quorums);
    }

    /// @notice Helper function to clear an ejection.
    function deleteEjection(address operator) internal {
        EigenDAEjectionTypes.EjectionParams storage operatorParams = s().ejectionParams[operator];
        operatorParams.ejector = address(0);
        operatorParams.quorums = hex"";
        operatorParams.proceedingTime = 0;
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
        return s().ejectionParams[operator].ejector;
    }

    /// @notice Compares two quorums to see if they are equal.
    function quorumsEqual(bytes memory quorums1, bytes memory quorums2) internal pure returns (bool) {
        return keccak256(quorums1) == keccak256(quorums2);
    }

    function ejectionParams(address operator) internal view returns (EigenDAEjectionTypes.EjectionParams storage) {
        return s().ejectionParams[operator];
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
