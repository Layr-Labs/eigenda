// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "./IEigenDAStructs.sol";

interface IEigenDADisperserRegistry {

    event DisperserAdded(uint32 indexed key, address indexed disperser);

    function setDisperserInfo(uint32 _disperserKey, DisperserInfo memory _disperserInfo) external;

    function disperserKeyToAddress(uint32 key) external view returns (address);
}