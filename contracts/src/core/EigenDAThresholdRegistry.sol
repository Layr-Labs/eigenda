// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDAThresholdRegistryStorage} from "./EigenDAThresholdRegistryStorage.sol";
import {IEigenDAThresholdRegistry} from "../interfaces/IEigenDAThresholdRegistry.sol";
import {OwnableUpgradeable} from "lib/openzeppelin-contracts-upgradeable/contracts/access/OwnableUpgradeable.sol";
import {BitmapUtils} from "lib/eigenlayer-middleware/src/libraries/BitmapUtils.sol";
import "../interfaces/IEigenDAStructs.sol";
import {EigenDATypesV1 as DATypesV1} from "../libraries/V1/EigenDATypesV1.sol";

/**
 * @title The `EigenDAThresholdRegistry` contract.
 * @author Layr Labs, Inc.
 */
contract EigenDAThresholdRegistry is EigenDAThresholdRegistryStorage, OwnableUpgradeable {
    constructor() {
        _disableInitializers();
    }

    function initialize(
        address _initialOwner,
        bytes memory _quorumAdversaryThresholdPercentages,
        bytes memory _quorumConfirmationThresholdPercentages,
        bytes memory _quorumNumbersRequired,
        DATypesV1.VersionedBlobParams[] memory _versionedBlobParams
    ) external initializer {
        _transferOwnership(_initialOwner);

        quorumAdversaryThresholdPercentages = _quorumAdversaryThresholdPercentages;
        quorumConfirmationThresholdPercentages = _quorumConfirmationThresholdPercentages;
        quorumNumbersRequired = _quorumNumbersRequired;

        for (uint256 i = 0; i < _versionedBlobParams.length; ++i) {
            _addVersionedBlobParams(_versionedBlobParams[i]);
        }
    }

    function addVersionedBlobParams(DATypesV1.VersionedBlobParams memory _versionedBlobParams)
        external
        onlyOwner
        returns (uint16)
    {
        return _addVersionedBlobParams(_versionedBlobParams);
    }

    function _addVersionedBlobParams(DATypesV1.VersionedBlobParams memory _versionedBlobParams)
        internal
        returns (uint16)
    {
        versionedBlobParams[nextBlobVersion] = _versionedBlobParams;
        emit VersionedBlobParamsAdded(nextBlobVersion, _versionedBlobParams);
        return nextBlobVersion++;
    }

    ///////////////////////// V1 ///////////////////////////////

    /// @notice Gets the adversary threshold percentage for a quorum
    function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber)
        public
        view
        virtual
        returns (uint8 adversaryThresholdPercentage)
    {
        if (quorumAdversaryThresholdPercentages.length > quorumNumber) {
            adversaryThresholdPercentage = uint8(quorumAdversaryThresholdPercentages[quorumNumber]);
        }
    }

    /// @notice Gets the confirmation threshold percentage for a quorum
    function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber)
        public
        view
        virtual
        returns (uint8 confirmationThresholdPercentage)
    {
        if (quorumConfirmationThresholdPercentages.length > quorumNumber) {
            confirmationThresholdPercentage = uint8(quorumConfirmationThresholdPercentages[quorumNumber]);
        }
    }

    /// @notice Checks if a quorum is required
    function getIsQuorumRequired(uint8 quorumNumber) public view virtual returns (bool) {
        uint256 quorumBitmap = BitmapUtils.setBit(0, quorumNumber);
        return (quorumBitmap & BitmapUtils.orderedBytesArrayToBitmap(quorumNumbersRequired) == quorumBitmap);
    }

    ///////////////////////// V2 ///////////////////////////////

    /// @notice Returns the blob params for a given blob version
    function getBlobParams(uint16 version) external view returns (DATypesV1.VersionedBlobParams memory) {
        return versionedBlobParams[version];
    }
}
