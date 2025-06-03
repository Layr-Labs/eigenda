// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";

library EigenDAEjectionTypes {
    struct OperatorProceedingParams {
        mapping(bytes32 => bool) salts;
        uint64 proceedingTime;
        uint64 lastProceedingInitiated;
        bool churn;
    }

    struct ProceedingParams {
        mapping(bytes32 => OperatorProceedingParams) operatorProceedingParams;
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
    event EjectionStarted(bytes32 operatorId, uint64 timestampStarted, uint64 ejectionTime);

    event EjectionCancelled(bytes32 operatorId);

    event EjectionCompleted(bytes32 operatorId);

    event ChurnStarted(bytes32 operatorId, uint64 timestampStarted, uint64 ejectionTime);

    event ChurnCancelled(bytes32 operatorId);

    event ChurnCompleted(bytes32 operatorId);

    function startChurn(bytes32 operatorId) internal {
        startProceeding(operatorId, churnParams());
        emit ChurnStarted(
            operatorId,
            churnParams().operatorProceedingParams[operatorId].lastProceedingInitiated,
            churnParams().operatorProceedingParams[operatorId].proceedingTime
        );
    }

    function cancelChurn(bytes32 operatorId) internal {
        cancelProceeding(operatorId, churnParams());
        emit ChurnCancelled(operatorId);
    }

    function completeChurn(bytes32 operatorId) internal {
        completeProceeding(operatorId, churnParams());
        emit ChurnCompleted(operatorId);
    }

    function startEjection(bytes32 operatorId) internal {
        startProceeding(operatorId, ejectionParams());
        emit EjectionStarted(
            operatorId,
            ejectionParams().operatorProceedingParams[operatorId].lastProceedingInitiated,
            ejectionParams().operatorProceedingParams[operatorId].proceedingTime
        );
    }

    function cancelEjection(bytes32 operatorId) internal {
        cancelProceeding(operatorId, ejectionParams());
        emit EjectionCancelled(operatorId);
    }

    function completeEjection(bytes32 operatorId) internal {
        completeProceeding(operatorId, ejectionParams());
        emit EjectionCompleted(operatorId);
    }

    function startProceeding(bytes32 operatorId, EigenDAEjectionTypes.ProceedingParams storage params) internal {
        EigenDAEjectionTypes.OperatorProceedingParams storage operatorParams =
            params.operatorProceedingParams[operatorId];

        require(operatorParams.proceedingTime == 0, "Proceeding already in progress");
        require(
            operatorParams.lastProceedingInitiated + params.proceedingCooldown <= block.timestamp,
            "Proceeding cooldown not met"
        );

        operatorParams.proceedingTime = uint64(block.timestamp) + params.proceedingDelay;
        operatorParams.lastProceedingInitiated = uint64(block.timestamp);
        emit EjectionStarted(operatorId, operatorParams.lastProceedingInitiated, operatorParams.proceedingTime);
    }

    function cancelProceeding(bytes32 operatorId, EigenDAEjectionTypes.ProceedingParams storage params) internal {
        EigenDAEjectionTypes.OperatorProceedingParams storage operatorParams =
            params.operatorProceedingParams[operatorId];
        require(operatorParams.proceedingTime > 0, "No proceeding in progress");

        operatorParams.proceedingTime = 0;
        emit EjectionCancelled(operatorId);
    }

    function completeProceeding(bytes32 operatorId, EigenDAEjectionTypes.ProceedingParams storage params) internal {
        EigenDAEjectionTypes.OperatorProceedingParams storage operatorParams =
            params.operatorProceedingParams[operatorId];
        require(operatorParams.proceedingTime > 0, "No proceeding in progress");

        require(block.timestamp >= operatorParams.proceedingTime, "Proceeding not yet due");

        operatorParams.proceedingTime = 0;
        emit EjectionCompleted(operatorId);
    }

    function ejectionInitiated(bytes32 operatorId) internal view returns (bool) {
        return ejectionParams().operatorProceedingParams[operatorId].proceedingTime > 0;
    }

    function churnInitiated(bytes32 operatorId) internal view returns (bool) {
        return churnParams().operatorProceedingParams[operatorId].proceedingTime > 0;
    }

    function proceedingInitiated(bytes32 operatorId, EigenDAEjectionTypes.ProceedingParams storage params)
        internal
        view
        returns (bool)
    {
        return params.operatorProceedingParams[operatorId].proceedingTime > 0;
    }

    function consumeSignature(
        bytes32 operatorId,
        address recipient,
        bytes32 salt,
        bytes memory signature,
        EigenDAEjectionTypes.ProceedingParams storage params
    ) internal {
        EigenDAEjectionTypes.OperatorProceedingParams storage operatorParams =
            params.operatorProceedingParams[operatorId];

        require(!operatorParams.salts[salt], "Signature already consumed");
        // Placeholder for signature verification logic
        // This should be replaced with actual signature verification logic
        recipient;
        signature;
        operatorParams.salts[salt] = true;
    }

    function ejectionParams() internal view returns (EigenDAEjectionTypes.ProceedingParams storage) {
        return EigenDAEjectionStorage.layout().ejectionParams;
    }

    function churnParams() internal view returns (EigenDAEjectionTypes.ProceedingParams storage) {
        return EigenDAEjectionStorage.layout().churnParams;
    }
}
