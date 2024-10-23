// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAThresholdRegistry} from "../interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDAServiceManager} from "../interfaces/IEigenDAServiceManager.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {BitmapUtils} from "eigenlayer-middleware/libraries/BitmapUtils.sol";

contract EigenDAThresholdRegistry is IEigenDAThresholdRegistry, Ownable {

    address public immutable eigenDAServiceManager;

    bytes public quorumAdversaryThresholdPercentages = hex"212121";

    bytes public quorumConfirmationThresholdPercentages = hex"373737";

    bytes public quorumNumbersRequired = hex"0001";

    modifier onlyServiceManagerOwner() {
        require(msg.sender == Ownable(eigenDAServiceManager).owner(), "EigenDAThresholdRegistry: only the service manager owner can call this function");
        _;
    }

    constructor(address _eigenDAServiceManager) {
        eigenDAServiceManager = _eigenDAServiceManager;
    }

    function updateQuorumAdversaryThresholdPercentages(bytes memory _quorumAdversaryThresholdPercentages) external onlyServiceManagerOwner {
        quorumAdversaryThresholdPercentages = _quorumAdversaryThresholdPercentages;
    }

    function updateQuorumConfirmationThresholdPercentages(bytes memory _quorumConfirmationThresholdPercentages) external onlyServiceManagerOwner {
        quorumConfirmationThresholdPercentages = _quorumConfirmationThresholdPercentages;
    }

    function updateQuorumNumbersRequired(bytes memory _quorumNumbersRequired) external onlyServiceManagerOwner {
        quorumNumbersRequired = _quorumNumbersRequired;
    }

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

}