// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "lib/openzeppelin-contracts-upgradeable/contracts/access/OwnableUpgradeable.sol";
import {EIP712Upgradeable} from "lib/openzeppelin-contracts-upgradeable/contracts/utils/cryptography/draft-EIP712Upgradeable.sol";
import {ECDSA} from "lib/openzeppelin-contracts/contracts/utils/cryptography/ECDSA.sol";
import {EigenDADisperserRegistryStorageV2} from "./EigenDADisperserRegistryStorageV2.sol";
import {IEigenDADisperserRegistryV2} from "src/core/interfaces/IEigenDADisperserRegistryV2.sol";
import {EigenDATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";
import {EnumerableSet} from "lib/openzeppelin-contracts/contracts/utils/structs/EnumerableSet.sol";
 
/**
 * @title Registry for EigenDA disperser info V2
 * @author Layr Labs, Inc.
 */
contract EigenDADisperserRegistryV2 is
    OwnableUpgradeable,
    EIP712Upgradeable,
    EigenDADisperserRegistryStorageV2,
    IEigenDADisperserRegistryV2
{
    using EnumerableSet for EnumerableSet.UintSet;

    constructor() {
        _disableInitializers();
    }

    function initialize(address _initialOwner) external initializer {
        __Ownable_init();
        __EIP712_init("EigenDADisperserRegistry", "2");
        _transferOwnership(_initialOwner);
        nextDisperserId = 1; // Start IDs from 1, 0 is reserved for "not found"
    }

    /// -----------------------------------------------------------------------
    /// External Logic
    /// -----------------------------------------------------------------------

    /// @inheritdoc IEigenDADisperserRegistryV2
    function registerDisperser(address disperserAddress, string memory relayURL, bytes memory pubKey)
        external
        returns (uint32)
    {
        if (disperserAddress == address(0)) revert InputAddressZero();

        if (bytes(relayURL).length == 0) revert InvalidRelayURL();
        
        if (pubKey.length == 0) revert InvalidPublicKey();
        
        if (disperserAddressToId[disperserAddress] != 0) revert DisperserAlreadyRegistered();

        uint32 disperserId = nextDisperserId;
        nextDisperserId++;

        disperserKeyToInfo[disperserId] =
            EigenDATypesV2.DisperserInfoV2({disperserAddress: disperserAddress, relayURL: relayURL, pubKey: pubKey});

        disperserAddressToId[disperserAddress] = disperserId;

        emit DisperserRegistered(disperserId, disperserAddress, relayURL);
        return disperserId;
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function deregisterDisperser(uint32 disperserId, bytes memory signature) external {
        if (disperserKeyToInfo[disperserId].disperserAddress == address(0)) revert DisperserNotFound();

        // Verify EIP712 signature
        bytes32 structHash = keccak256(abi.encode(DEREGISTRATIONTYPEHASH, disperserId));
        bytes32 digest = hashTypedDataV4(structHash);
        address signer = ECDSA.recover(digest, signature);

        if (signer != disperserKeyToInfo[disperserId].disperserAddress) revert InvalidSignature();

        // Remove from sets if present
        defaultDispersersSet.remove(disperserId);
        onDemandDispersersSet.remove(disperserId);

        // Clear mappings
        address disperserAddr = disperserKeyToInfo[disperserId].disperserAddress;
        delete disperserAddressToId[disperserAddr];
        delete disperserKeyToInfo[disperserId];

        emit DisperserDeregistered(disperserId);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function updateRelayURL(uint32 disperserId, string memory newRelayURL, bytes memory signature) external {
        if (disperserKeyToInfo[disperserId].disperserAddress == address(0)) revert DisperserNotFound();
        if (bytes(newRelayURL).length == 0) revert InvalidRelayURL();

        // Verify EIP712 signature
        bytes32 structHash =
            keccak256(abi.encode(UPDATERELAYURLTYPEHASH, disperserId, keccak256(bytes(newRelayURL))));
        bytes32 digest = hashTypedDataV4(structHash);
        address signer = ECDSA.recover(digest, signature);

        if (signer != disperserKeyToInfo[disperserId].disperserAddress) revert InvalidSignature();

        disperserKeyToInfo[disperserId].relayURL = newRelayURL;

        emit RelayURLUpdated(disperserId, newRelayURL);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function getDisperserInfo(uint32[] memory ids) external view returns (EigenDATypesV2.DisperserInfoV2[] memory) {
        EigenDATypesV2.DisperserInfoV2[] memory infos = new EigenDATypesV2.DisperserInfoV2[](ids.length);
        for (uint256 i = 0; i < ids.length; i++) {
            infos[i] = disperserKeyToInfo[ids[i]];
        }
        return infos;
    }

    /// -----------------------------------------------------------------------
    /// Owner-only Logic
    /// -----------------------------------------------------------------------

    /// @inheritdoc IEigenDADisperserRegistryV2
    function addDefaultDisperser(uint32 disperserId) external onlyOwner {
        if (disperserKeyToInfo[disperserId].disperserAddress == address(0)) revert DisperserNotFound();
        if (!defaultDispersersSet.add(disperserId)) revert AlreadyInDefaultSet();
        emit DefaultDisperserAdded(disperserId);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function addOnDemandDisperser(uint32 disperserId) external onlyOwner {
        if (disperserKeyToInfo[disperserId].disperserAddress == address(0)) revert DisperserNotFound();
        if (!onDemandDispersersSet.add(disperserId)) revert AlreadyInOnDemandSet();
        emit OnDemandDisperserAdded(disperserId);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function removeDefaultDisperser(uint32 disperserId) external onlyOwner {
        if (!defaultDispersersSet.remove(disperserId)) revert NotInDefaultSet();
        emit DefaultDisperserRemoved(disperserId);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function removeOnDemandDisperser(uint32 disperserId) external onlyOwner {
        if (!onDemandDispersersSet.remove(disperserId)) revert NotInOnDemandSet();
        emit OnDemandDisperserRemoved(disperserId);
    }

    /// -----------------------------------------------------------------------
    /// View Logic
    /// -----------------------------------------------------------------------

    /// @inheritdoc IEigenDADisperserRegistryV2
    /// TODO:
    function getAllDisperserIds() external view returns (uint32[] memory) {}

    /// @inheritdoc IEigenDADisperserRegistryV2
    function disperserKeyToAddress(uint32 key) external view returns (address) {
        return disperserKeyToInfo[key].disperserAddress;
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function getDisperserIdByAddress(address disperserAddress) external view returns (uint32) {
        return disperserAddressToId[disperserAddress];
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function getDefaultDispersers() external view returns (uint32[] memory) {
        uint256 length = defaultDispersersSet.length();
        uint32[] memory dispersers = new uint32[](length);
        for (uint256 i = 0; i < length; i++) {
            dispersers[i] = uint32(defaultDispersersSet.at(i));
        }
        return dispersers;
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function getOnDemandDispersers() external view returns (uint32[] memory) {
        uint256 length = onDemandDispersersSet.length();
        uint32[] memory dispersers = new uint32[](length);
        for (uint256 i = 0; i < length; i++) {
            dispersers[i] = uint32(onDemandDispersersSet.at(i));
        }
        return dispersers;
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function isDefaultDisperser(uint32 disperserId) external view returns (bool) {
        return defaultDispersersSet.contains(disperserId);
    }

    /// @inheritdoc IEigenDADisperserRegistryV2
    function isOnDemandDisperser(uint32 disperserId) external view returns (bool) {
        return onDemandDispersersSet.contains(disperserId);
    }
}
