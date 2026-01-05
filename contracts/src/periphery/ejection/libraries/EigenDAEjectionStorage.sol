// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDAEjectionTypes} from "src/periphery/ejection/libraries/EigenDAEjectionTypes.sol";
import {IEigenDAEjectionManager} from "src/periphery/ejection/IEigenDAEjectionManager.sol";
import {IAccessControl} from "@openzeppelin/contracts/access/IAccessControl.sol";
import {IBLSApkRegistry} from "lib/eigenlayer-middleware/src/interfaces/IBLSApkRegistry.sol";
import {BLSSignatureChecker} from "lib/eigenlayer-middleware/src/BLSSignatureChecker.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";

abstract contract ImmutableEigenDAEjectionsStorage is IEigenDAEjectionManager {
    /// @dev callee dependencies
    IAccessControl public immutable accessControl;
    IBLSApkRegistry public immutable blsApkKeyRegistry;
    BLSSignatureChecker public immutable signatureChecker;
    IRegistryCoordinator public immutable registryCoordinator;

    constructor(
        IAccessControl accessControl_,
        IBLSApkRegistry blsApkKeyRegistry_,
        BLSSignatureChecker signatureChecker_,
        IRegistryCoordinator registryCoordinator_
    ) {
        accessControl = accessControl_;
        blsApkKeyRegistry = blsApkKeyRegistry_;
        signatureChecker = signatureChecker_;
        registryCoordinator = registryCoordinator_;
    }
}

library EigenDAEjectionStorage {
    string internal constant STORAGE_ID = "eigen.da.ejection";
    bytes32 internal constant STORAGE_POSITION =
        keccak256(abi.encode(uint256(keccak256(abi.encodePacked(STORAGE_ID))) - 1)) & ~bytes32(uint256(0xff));

    struct Layout {
        /// @dev ejection state
        mapping(address => EigenDAEjectionTypes.EjecteeState) ejectees;

        /// @dev protocol params
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
