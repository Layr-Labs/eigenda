// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";

interface IEigenDADisperserRegistry {
    event DisperserAdded(uint32 indexed key, address indexed disperser);

    function setDisperserInfo(uint32 _disperserKey, EigenDATypesV2.DisperserInfo memory _disperserInfo) external;

    function disperserKeyToAddress(uint32 key) external view returns (address);
}
