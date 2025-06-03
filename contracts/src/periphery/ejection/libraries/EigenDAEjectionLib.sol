// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";

library EigenDAEjectionTypes {
    struct OperatorProceedingParams {
        mapping(address => bool) salts;
        uint64 proceedingTime;
        uint64 lastProceedingInitiated;
        bytes quorums;
    }

    struct ProceedingParams {
        mapping(address => OperatorProceedingParams) operatorProceedingParams;
        uint64 proceedingDelay;
        uint64 proceedingCooldown;
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

    function startChurn(address operator, bytes memory quorums) internal {
        startProceeding(operator, quorums, churnParams());
        emit ChurnStarted(
            operator,
            quorums,
            churnParams().operatorProceedingParams[operator].lastProceedingInitiated,
            churnParams().operatorProceedingParams[operator].proceedingTime
        );
    }

    function cancelChurn(address operator) internal {
        cancelProceeding(operator, churnParams());
        emit ChurnCancelled(operator);
    }

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
            operatorParams.lastProceedingInitiated + params.proceedingCooldown <= block.timestamp,
            "Proceeding cooldown not met"
        );

        operatorParams.quorums = quorums;
        operatorParams.proceedingTime = uint64(block.timestamp) + params.proceedingDelay;
        operatorParams.lastProceedingInitiated = uint64(block.timestamp);
    }

    function cancelProceeding(address operator, EigenDAEjectionTypes.ProceedingParams storage params) internal {
        EigenDAEjectionTypes.OperatorProceedingParams storage operatorParams = params.operatorProceedingParams[operator];
        require(operatorParams.proceedingTime > 0, "No proceeding in progress");

        operatorParams.proceedingTime = 0;
    }

    function completeProceeding(address operator, bytes memory, EigenDAEjectionTypes.ProceedingParams storage params)
        internal
    {
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
}
