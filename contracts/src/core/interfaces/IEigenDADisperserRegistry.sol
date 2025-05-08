// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV3} from "src/core/libraries/v3/EigenDATypesV3.sol";

interface IEigenDADisperserRegistry {
    function registerDisperser(address disperserAddress, string memory disperserURL)
        external
        returns (uint32 disperserKey);

    function transferDisperserOwnership(uint32 disperserKey, address newOwner) external;

    function updateDisperserInfo(uint32 disperserKey, address disperser, string memory disperserURL) external;

    function deregisterDisperser(uint32 disperserKey) external;

    function withdraw(uint32 disperserKey) external;

    function getDisperserInfo(uint32 disperserKey) external view returns (EigenDATypesV3.DisperserInfo memory);

    function getLockedDeposit(uint32 disperserKey)
        external
        view
        returns (EigenDATypesV3.LockedDisperserDeposit memory, uint64);

    function getDepositParams() external view returns (EigenDATypesV3.LockedDisperserDeposit memory);
}
