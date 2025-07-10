// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";

library EigenDAEjectionTypes {
    /// @param proceedingTime Timestamp when the proceeding is set to complete
    /// @param lastProceedingInitiated Timestamp of when the last proceeding was initiated to enforce cooldowns
    /// @param quorums The quorums associated with the proceeding
    struct EjectionProceedingParams {
        uint64 proceedingTime;
        uint64 lastProceedingInitiated;
        bytes quorums;
    }
}

library EigenDAEjectionStorage {
    string internal constant STORAGE_ID = "eigen.da.ejection";
    bytes32 internal constant STORAGE_POSITION =
        keccak256(abi.encode(uint256(keccak256(abi.encodePacked(STORAGE_ID))) - 1)) & ~bytes32(uint256(0xff));

    struct Layout {
        mapping(address => EigenDAEjectionTypes.EjectionProceedingParams) proceedingParams;
        uint64 delay;
        uint64 cooldown;
    }

    function layout() internal pure returns (Layout storage s) {
        bytes32 position = STORAGE_POSITION;
        assembly {
            s.slot := position
        }
    }
}

library EigenDAEjectionLib {
    event EjectionStarted(address operator, bytes quorums, uint64 timestampStarted, uint64 ejectionTime);

    event EjectionCancelled(address operator);

    event EjectionCompleted(address operator, bytes quorums);

    /// @notice Starts an ejection process for an operator.
    function startEjection(address operator, bytes memory quorums) internal {
        startProceeding(operator, quorums);
        emit EjectionStarted(
            operator,
            quorums,
            ejectionStorage().proceedingParams[operator].lastProceedingInitiated,
            ejectionStorage().proceedingParams[operator].proceedingTime
        );
    }

    /// @notice Cancels an ejection process for an operator.
    function cancelEjection(address operator) internal {
        cancelProceeding(operator);
        emit EjectionCancelled(operator);
    }

    /// @notice Completes an ejection process for an operator.
    function completeEjection(address operator, bytes memory quorums) internal {
        completeProceeding(operator, quorums);
        emit EjectionCompleted(operator, quorums);
    }

    /// @notice Starts a proceeding process for an operator.
    function startProceeding(address operator, bytes memory quorums) internal {
        EigenDAEjectionTypes.EjectionProceedingParams storage operatorParams =
            ejectionStorage().proceedingParams[operator];

        require(operatorParams.proceedingTime == 0, "Proceeding already in progress");
        require(
            operatorParams.lastProceedingInitiated + ejectionStorage().cooldown <= block.timestamp,
            "Proceeding cooldown not met"
        );

        operatorParams.quorums = quorums;
        operatorParams.proceedingTime = uint64(block.timestamp) + ejectionStorage().delay;
        operatorParams.lastProceedingInitiated = uint64(block.timestamp);
    }

    /// @notice Cancels a proceeding process for an operator.
    function cancelProceeding(address operator) internal {
        EigenDAEjectionTypes.EjectionProceedingParams storage operatorParams =
            ejectionStorage().proceedingParams[operator];
        require(operatorParams.proceedingTime > 0, "No proceeding in progress");

        operatorParams.quorums = hex"";
        operatorParams.proceedingTime = 0;
    }

    /// @notice Completes a proceeding process for an operator.
    function completeProceeding(address operator, bytes memory quorums) internal {
        require(quorumsEqual(ejectionStorage().proceedingParams[operator].quorums, quorums), "Quorums do not match");
        EigenDAEjectionTypes.EjectionProceedingParams storage operatorParams =
            ejectionStorage().proceedingParams[operator];
        require(operatorParams.proceedingTime > 0, "No proceeding in progress");

        require(block.timestamp >= operatorParams.proceedingTime, "Proceeding not yet due");

        operatorParams.quorums = hex"";
        operatorParams.proceedingTime = 0;
    }

    /// @notice Checks if an ejection or churn process has been initiated for the operator.
    function ejectionInitiated(address operator) internal view returns (bool) {
        return ejectionStorage().proceedingParams[operator].proceedingTime > 0;
    }

    /// @notice Checks if a proceeding has been initiated for the operator.
    function proceedingInitiated(address operator) internal view returns (bool) {
        return ejectionStorage().proceedingParams[operator].proceedingTime > 0;
    }

    /// @notice Returns the ejection storage.
    function ejectionStorage() internal pure returns (EigenDAEjectionStorage.Layout storage) {
        return EigenDAEjectionStorage.layout();
    }

    /// @notice Compares two quorums to see if they are equal.
    function quorumsEqual(bytes memory quorums1, bytes memory quorums2) internal pure returns (bool) {
        return keccak256(quorums1) == keccak256(quorums2);
    }
}
