// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";

library EigenDAEjectionTypes {
    /// @param proceedingTime Timestamp when the proceeding is set to complete
    /// @param lastProceedingInitiated Timestamp of when the last proceeding was initiated to enforce cooldowns
    /// @param quorums The quorums associated with the proceeding
    struct OperatorProceedingParams {
        uint64 proceedingTime;
        uint64 lastProceedingInitiated;
        bytes quorums;
    }

    /// @param operatorProceedingParams Mapping of operator addresses to their proceeding parameters
    /// @param delay Delay before the proceeding can be completed
    /// @param cooldown Cooldown period after a proceeding is completed before a new one can be initiated
    struct ProceedingParams {
        mapping(address => OperatorProceedingParams) operatorProceedingParams;
        uint64 delay;
        uint64 cooldown;
    }
}

library EigenDAEjectionStorage {
    string internal constant STORAGE_ID = "eigen.da.ejection";
    bytes32 internal constant STORAGE_POSITION =
        keccak256(abi.encode(uint256(keccak256(abi.encodePacked(STORAGE_ID))) - 1)) & ~bytes32(uint256(0xff));

    struct Layout {
        EigenDAEjectionTypes.ProceedingParams ejectionParams;
        EigenDAEjectionTypes.ProceedingParams churnParams;
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

    event ChurnStarted(address operator, bytes quorums, uint64 timestampStarted, uint64 ejectionTime);

    event ChurnCancelled(address operator);

    event ChurnCompleted(address operator, bytes quorums);

    /// @notice Starts a churning process for an operator.
    function startChurn(address operator, bytes memory quorums) internal {
        startProceeding(operator, quorums, churnParams());
        emit ChurnStarted(
            operator,
            quorums,
            churnParams().operatorProceedingParams[operator].lastProceedingInitiated,
            churnParams().operatorProceedingParams[operator].proceedingTime
        );
    }

    /// @notice Cancels a churning process for an operator.
    function cancelChurn(address operator) internal {
        cancelProceeding(operator, churnParams());
        emit ChurnCancelled(operator);
    }

    /// @notice Completes a churning process for an operator.
    function completeChurn(address operator, bytes memory quorums) internal {
        completeProceeding(operator, quorums, churnParams());
        emit ChurnCompleted(operator, quorums);
    }

    function startEjection(address operator, bytes memory quorums) internal {
        startProceeding(operator, quorums, ejectionParams());
        emit EjectionStarted(
            operator,
            quorums,
            ejectionParams().operatorProceedingParams[operator].lastProceedingInitiated,
            ejectionParams().operatorProceedingParams[operator].proceedingTime
        );
    }

    function cancelEjection(address operator) internal {
        cancelProceeding(operator, ejectionParams());
        emit EjectionCancelled(operator);
    }

    function completeEjection(address operator, bytes memory quorums) internal {
        completeProceeding(operator, quorums, ejectionParams());
        emit EjectionCompleted(operator, quorums);
    }

    function startProceeding(
        address operator,
        bytes memory quorums,
        EigenDAEjectionTypes.ProceedingParams storage params
    ) internal {
        EigenDAEjectionTypes.OperatorProceedingParams storage operatorParams = params.operatorProceedingParams[operator];

        require(operatorParams.proceedingTime == 0, "Proceeding already in progress");
        require(
            operatorParams.lastProceedingInitiated + params.cooldown <= block.timestamp, "Proceeding cooldown not met"
        );

        operatorParams.quorums = quorums;
        operatorParams.proceedingTime = uint64(block.timestamp) + params.delay;
        operatorParams.lastProceedingInitiated = uint64(block.timestamp);
    }

    function cancelProceeding(address operator, EigenDAEjectionTypes.ProceedingParams storage params) internal {
        EigenDAEjectionTypes.OperatorProceedingParams storage operatorParams = params.operatorProceedingParams[operator];
        require(operatorParams.proceedingTime > 0, "No proceeding in progress");

        operatorParams.proceedingTime = 0;
    }

    function completeProceeding(
        address operator,
        bytes memory quorums,
        EigenDAEjectionTypes.ProceedingParams storage params
    ) internal {
        require(quorumsEqual(params.operatorProceedingParams[operator].quorums, quorums), "Quorums do not match");
        EigenDAEjectionTypes.OperatorProceedingParams storage operatorParams = params.operatorProceedingParams[operator];
        require(operatorParams.proceedingTime > 0, "No proceeding in progress");
        // require(operatorParams.quorums == quorums, "Quorums do not match"); // TODO: FIX THIS

        require(block.timestamp >= operatorParams.proceedingTime, "Proceeding not yet due");

        operatorParams.quorums = hex"";
        operatorParams.proceedingTime = 0;
    }

    function ejectionInitiated(address operator) internal view returns (bool) {
        return ejectionParams().operatorProceedingParams[operator].proceedingTime > 0;
    }

    function churnInitiated(address operator) internal view returns (bool) {
        return churnParams().operatorProceedingParams[operator].proceedingTime > 0;
    }

    function proceedingInitiated(address operator, EigenDAEjectionTypes.ProceedingParams storage params)
        internal
        view
        returns (bool)
    {
        return params.operatorProceedingParams[operator].proceedingTime > 0;
    }

    function ejectionParams() internal view returns (EigenDAEjectionTypes.ProceedingParams storage) {
        return EigenDAEjectionStorage.layout().ejectionParams;
    }

    function churnParams() internal view returns (EigenDAEjectionTypes.ProceedingParams storage) {
        return EigenDAEjectionStorage.layout().churnParams;
    }

    function quorumsEqual(bytes memory quorums1, bytes memory quorums2) internal pure returns (bool) {
        return keccak256(quorums1) == keccak256(quorums2);
    }
}
