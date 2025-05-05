// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV3} from "src/core/libraries/v3/EigenDATypesV3.sol";

interface IEigenDADisperserRegistry {
    event DisperserAdded(uint32 indexed key, address indexed disperser);

    function setDisperserInfo(uint32 _disperserKey, EigenDATypesV3.DisperserInfo memory _disperserInfo) external;

    function disperserKeyToAddress(uint32 key) external view returns (address);
}
