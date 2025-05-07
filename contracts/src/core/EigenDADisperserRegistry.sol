// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDADisperserRegistry} from "src/core/interfaces/IEigenDADisperserRegistry.sol";
import {EigenDATypesV3} from "src/core/libraries/v3/EigenDATypesV3.sol";

/**
 * @title Registry for EigenDA disperser info
 * @author Layr Labs, Inc.
 */
contract EigenDADisperserRegistry is IEigenDADisperserRegistry {

    function setDisperserInfo(uint32 disperserKey, address disperser, string memory disperserURL)
        external payable
    {
        _disperserInfo[disperserKey] = EigenDATypesV3.DisperserInfo({
            disperser: disperser,
            disperserURL: disperserURL,
            registered: true,
            withdrawalUnlock: type(uint64).max
        });
        emit DisperserAdded(disperserKey, disperser);
    }

    function deregisterDisperser(uint32 _disperserKey) external {
        require(_disperserInfo[_disperserKey].registered, "Disperser not registered");
        _disperserInfo[_disperserKey].registered = false;
        _disperserInfo[_disperserKey].withdrawalUnlock = uint64(block.timestamp + 1 days);

    }

    function withdrawDeposit(uint32 _disperserKey) external {
    }

    function disperserInfo(uint32 _key) external view returns (EigenDATypesV3.DisperserInfo memory) {
        return _disperserInfo[_key];
    }
}
