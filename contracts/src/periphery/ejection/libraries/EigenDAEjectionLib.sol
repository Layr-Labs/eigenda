// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDAEjectionTypes} from "src/periphery/ejection/libraries/EigenDAEjectionTypes.sol";
import {EigenDAEjectionStorage} from "src/periphery/ejection/libraries/EigenDAEjectionStorage.sol";

library EigenDAEjectionLib {
    event EjectionStarted(address operator, bytes quorums, uint64 timestampStarted, uint64 ejectionTime);

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
    function startEjection(address operator, bytes memory quorums) internal {
        EigenDAEjectionTypes.EjectionParams storage operatorParams = s().ejectionParams[operator];

        require(operatorParams.proceedingTime == 0, "Ejection already in progress");
        require(operatorParams.lastProceedingInitiated + s().cooldown <= block.timestamp, "Ejection cooldown not met");

        operatorParams.quorums = quorums;
        operatorParams.proceedingTime = uint64(block.timestamp) + s().delay;
        operatorParams.lastProceedingInitiated = uint64(block.timestamp);
        emit EjectionStarted(operator, quorums, operatorParams.lastProceedingInitiated, operatorParams.proceedingTime);
    }

    /// @notice Cancels an ejection process for an operator.
    function cancelEjection(address operator) internal {
        EigenDAEjectionTypes.EjectionParams storage operatorParams = s().ejectionParams[operator];
        require(operatorParams.proceedingTime > 0, "No ejection in progress");

        operatorParams.quorums = hex"";
        operatorParams.proceedingTime = 0;
        emit EjectionCancelled(operator);
    }

    /// @notice Completes an ejection process for an operator.
    function completeEjection(address operator, bytes memory quorums) internal {
        require(quorumsEqual(s().ejectionParams[operator].quorums, quorums), "Quorums do not match");
        EigenDAEjectionTypes.EjectionParams storage operatorParams = s().ejectionParams[operator];
        require(operatorParams.proceedingTime > 0, "No proceeding in progress");

        require(block.timestamp >= operatorParams.proceedingTime, "Proceeding not yet due");

        operatorParams.quorums = hex"";
        operatorParams.proceedingTime = 0;
        emit EjectionCompleted(operator, quorums);
    }

    /// @notice Checks if an ejection or churn process has been initiated for the operator.
    function ejectionInitiated(address operator) internal view returns (bool) {
        return s().ejectionParams[operator].proceedingTime > 0;
    }

    /// @notice Compares two quorums to see if they are equal.
    function quorumsEqual(bytes memory quorums1, bytes memory quorums2) internal pure returns (bool) {
        return keccak256(quorums1) == keccak256(quorums2);
    }

    function ejectionParams(address operator) internal view returns (EigenDAEjectionTypes.EjectionParams storage) {
        return s().ejectionParams[operator];
    }

    function getDelay() internal view returns (uint64) {
        return s().delay;
    }

    function getCooldown() internal view returns (uint64) {
        return s().cooldown;
    }

    /// @notice Returns the ejection storage.
    function s() private pure returns (EigenDAEjectionStorage.Layout storage) {
        return EigenDAEjectionStorage.layout();
    }
}
