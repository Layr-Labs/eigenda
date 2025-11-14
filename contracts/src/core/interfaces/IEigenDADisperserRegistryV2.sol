// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";

interface IEigenDADisperserRegistryV2 {
    // Custom Errors
    error InputAddressZero();
    error InvalidRelayURL();
    error InvalidPublicKey();
    error DisperserAlreadyRegistered();
    error DisperserNotFound();
    error InvalidSignature();
    error AlreadyInDefaultSet();
    error NotInDefaultSet();
    error AlreadyInOnDemandSet();
    error NotInOnDemandSet();

    // Events
    event DisperserRegistered(uint32 indexed id, address indexed addr, string relayURL);
    event DisperserDeregistered(uint32 indexed id);
    event RelayURLUpdated(uint32 indexed id, string newURL);
    event DefaultDisperserAdded(uint32 indexed id);
    event DefaultDisperserRemoved(uint32 indexed id);
    event OnDemandDisperserAdded(uint32 indexed id);
    event OnDemandDisperserRemoved(uint32 indexed id);

    // Registration and updates
    function registerDisperser(address _disperserAddress, string memory _relayURL, bytes memory _pubKey)
        external
        returns (uint32);
    function deregisterDisperser(uint32 _disperserId, bytes memory _signature) external;
    function updateRelayURL(uint32 _disperserId, string memory _newRelayURL, bytes memory _signature) external;

    // Getters
    function getDisperserInfo(uint32[] memory _ids) external view returns (EigenDATypesV2.DisperserInfoV2[] memory);
    function disperserKeyToAddress(uint32 _key) external view returns (address);
    function getDisperserIdByAddress(address _disperserAddress) external view returns (uint32);
    function getAllDisperserIds() external view returns (uint32[] memory);

    // Default dispersers management
    function addDefaultDisperser(uint32 _disperserId) external;
    function removeDefaultDisperser(uint32 _disperserId) external;
    function getDefaultDispersers() external view returns (uint32[] memory);
    function isDefaultDisperser(uint32 _disperserId) external view returns (bool);

    // On-demand dispersers management
    function addOnDemandDisperser(uint32 _disperserId) external;
    function removeOnDemandDisperser(uint32 _disperserId) external;
    function getOnDemandDispersers() external view returns (uint32[] memory);
    function isOnDemandDisperser(uint32 _disperserId) external view returns (bool);
}

