// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "./IEigenDAStructs.sol";

interface IEigenDADisperserRegistry {

    event DisperserAdded(uint32 indexed disperserKey, address disperserAddress);
    event DisperserRemoved(uint32 indexed disperserKey, address registrant);

    function setDisperserInfo(uint32 _disperserKey, DisperserInfo memory _disperserInfo) external payable;

    function deregisterDisperser(uint32 _disperserKey) external;
    
    function disperserKeyToAddress(uint32 key) external view returns (address);
}