// SPDX-License-Identifier: MIT
pragma solidity ^0.8.12;

import {InitializableStorage} from "src/core/libraries/v3/initializable/InitializableStorage.sol";

library InitializableLib {
    event Initialized(uint8 version);

    error AlreadyInitialized();

    function s() private pure returns (InitializableStorage.Layout storage) {
        return InitializableStorage.layout();
    }

    function initialize() internal {
        setInitializedVersion(1);
    }

    function reinitialize(uint8 version) internal {
        setInitializedVersion(version);
    }

    function setInitializedVersion(uint8 version) internal {
        if (s().initialized >= version) {
            revert AlreadyInitialized();
        }

        s().initialized = version;
        emit Initialized(version);
    }

    function getInitializedVersion() internal view returns (uint8 version) {
        version = s().initialized;
    }
}
