// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {DisperserRegistryTypes} from "src/core/libraries/v3/disperser/DisperserRegistryTypes.sol";

interface IEigenDADisperserRegistry {
    function registerDisperser(address disperserAddress, string memory disperserURL)
        external
        returns (uint32 disperserKey);

    function transferDisperserOwnership(uint32 disperserKey, address newOwner) external;

    function updateDisperserInfo(uint32 disperserKey, address disperser, string memory disperserURL) external;

    function deregisterDisperser(uint32 disperserKey) external;

    function withdraw(uint32 disperserKey) external;

    function getDepositParams() external view returns (DisperserRegistryTypes.LockedDisperserDeposit memory);
}
