// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {BitmapUtils} from "lib/eigenlayer-middleware/src/libraries/BitmapUtils.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";

/**
 * @title The `EigenDAThresholdRegistryImmutableV1` contract.
 * @author Layr Labs, Inc.
 * @notice this contract is an immutable version of the `EigenDAThresholdRegistry` contract and is only
 *         intended to be used for enabling custom quorums/thresholds for rollups using EigenDAV1.
 *         The lifespan of this contract is expected to be short, as it is intended to be used
 *         for a soon-to-be deprecated protocol version.
 */
contract EigenDAThresholdRegistryImmutableV1 is IEigenDAThresholdRegistry {
    /// @notice The adversary threshold percentage for the quorum at position `quorumNumber`
    bytes public quorumAdversaryThresholdPercentages;

    /// @notice The confirmation threshold percentage for the quorum at position `quorumNumber`
    bytes public quorumConfirmationThresholdPercentages;

    /// @notice The set of quorum numbers that are required
    bytes public quorumNumbersRequired;

    constructor(
        bytes memory _quorumAdversaryThresholdPercentages,
        bytes memory _quorumConfirmationThresholdPercentages,
        bytes memory _quorumNumbersRequired
    ) {
        quorumAdversaryThresholdPercentages = _quorumAdversaryThresholdPercentages;
        quorumConfirmationThresholdPercentages = _quorumConfirmationThresholdPercentages;
        quorumNumbersRequired = _quorumNumbersRequired;
    }

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

    // @notice Gets the quorum numbers that are required. Disabled for this immutable version since its only
    // usable for EigenDA V2.
    function getBlobParams(uint16) public pure returns (DATypesV1.VersionedBlobParams memory) {
        revert("EigenDAThresholdRegistryImmutableV1: Blob params not supported");
    }
}
