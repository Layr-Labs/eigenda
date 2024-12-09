// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDAThresholdRegistryStorage} from "./EigenDAThresholdRegistryStorage.sol";
import {IEigenDAThresholdRegistry} from "../interfaces/IEigenDAThresholdRegistry.sol";
import {OwnableUpgradeable} from "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";
import {BitmapUtils} from "eigenlayer-middleware/libraries/BitmapUtils.sol";
import "../interfaces/IEigenDAStructs.sol";

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
        VersionedBlobParams[] memory _versionedBlobParams,
        SecurityThresholds memory _defaultSecurityThresholdsV2
    ) external initializer {
        _transferOwnership(_initialOwner);

        quorumAdversaryThresholdPercentages = _quorumAdversaryThresholdPercentages;
        quorumConfirmationThresholdPercentages = _quorumConfirmationThresholdPercentages;
        quorumNumbersRequired = _quorumNumbersRequired;
        defaultSecurityThresholdsV2 = _defaultSecurityThresholdsV2;
        
        for (uint256 i = 0; i < _versionedBlobParams.length; ++i) {
            _addVersionedBlobParams(_versionedBlobParams[i]);
        }
    }

    function updateQuorumAdversaryThresholdPercentages(bytes memory _quorumAdversaryThresholdPercentages) external onlyOwner {
        emit QuorumAdversaryThresholdPercentagesUpdated(quorumAdversaryThresholdPercentages, _quorumAdversaryThresholdPercentages);
        quorumAdversaryThresholdPercentages = _quorumAdversaryThresholdPercentages;
    }

    function updateQuorumConfirmationThresholdPercentages(bytes memory _quorumConfirmationThresholdPercentages) external onlyOwner {
        emit QuorumConfirmationThresholdPercentagesUpdated(quorumConfirmationThresholdPercentages, _quorumConfirmationThresholdPercentages);
        quorumConfirmationThresholdPercentages = _quorumConfirmationThresholdPercentages;
    }

    function updateQuorumNumbersRequired(bytes memory _quorumNumbersRequired) external onlyOwner {
        emit QuorumNumbersRequiredUpdated(quorumNumbersRequired, _quorumNumbersRequired);
        quorumNumbersRequired = _quorumNumbersRequired;
    }

    function updateDefaultSecurityThresholdsV2(SecurityThresholds memory _defaultSecurityThresholdsV2) external onlyOwner {
        emit DefaultSecurityThresholdsV2Updated(defaultSecurityThresholdsV2, _defaultSecurityThresholdsV2);
        defaultSecurityThresholdsV2 = _defaultSecurityThresholdsV2;
    }

    function addVersionedBlobParams(VersionedBlobParams memory _versionedBlobParams) external onlyOwner returns (uint16) {
        return _addVersionedBlobParams(_versionedBlobParams);
    }

    function _addVersionedBlobParams(VersionedBlobParams memory _versionedBlobParams) internal returns (uint16) {
        versionedBlobParams[nextBlobVersion] = _versionedBlobParams;
        emit VersionedBlobParamsAdded(nextBlobVersion, _versionedBlobParams);
        return nextBlobVersion++;
    }

    ///////////////////////// V1 ///////////////////////////////

    /// @notice Gets the adversary threshold percentage for a quorum
    function getQuorumAdversaryThresholdPercentage(
        uint8 quorumNumber
    ) public view virtual returns (uint8 adversaryThresholdPercentage) {
        if(quorumAdversaryThresholdPercentages.length > quorumNumber){
            adversaryThresholdPercentage = uint8(quorumAdversaryThresholdPercentages[quorumNumber]);
        }
    }

    /// @notice Gets the confirmation threshold percentage for a quorum
    function getQuorumConfirmationThresholdPercentage(
        uint8 quorumNumber
    ) public view virtual returns (uint8 confirmationThresholdPercentage) {
        if(quorumConfirmationThresholdPercentages.length > quorumNumber){
            confirmationThresholdPercentage = uint8(quorumConfirmationThresholdPercentages[quorumNumber]);
        }
    }

    /// @notice Checks if a quorum is required
    function getIsQuorumRequired(
        uint8 quorumNumber
    ) public view virtual returns (bool) {
        uint256 quorumBitmap = BitmapUtils.setBit(0, quorumNumber);
        return (quorumBitmap & BitmapUtils.orderedBytesArrayToBitmap(quorumNumbersRequired) == quorumBitmap);
    }

    ///////////////////////// V2 ///////////////////////////////

    /// @notice Gets the default security thresholds for V2
    function getDefaultSecurityThresholdsV2() external view returns (SecurityThresholds memory) {
        return defaultSecurityThresholdsV2;
    }

    /// @notice Returns the blob params for a given blob version
    function getBlobParams(uint16 version) external view returns (VersionedBlobParams memory) {
        return versionedBlobParams[version];
    }
}