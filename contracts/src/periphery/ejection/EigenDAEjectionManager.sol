// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDAEjectionLib, EigenDAEjectionTypes} from "src/periphery/ejection/libraries/EigenDAEjectionLib.sol";
import {SafeERC20, IERC20} from "lib/openzeppelin-contracts/contracts/token/ERC20/utils/SafeERC20.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";

contract EigenDAEjectionManager {
    using EigenDAEjectionLib for address;
    using SafeERC20 for IERC20;

    address internal immutable _depositToken;
    uint256 internal immutable _depositAmount;
    address internal immutable _registryCoordinator;

    constructor(address depositToken_, uint256 depositAmount_, address registryCoordinator_) {
        _depositToken = depositToken_;
        _depositAmount = depositAmount_;
        _registryCoordinator = registryCoordinator_;
    }

    modifier onlyWatcher(address sender) {
        _onlyWatcher(sender);
        _;
    }

    /// WATCHER FUNCTIONS

    /// @notice Starts the ejection process for an operator. Takes a deposit from the watcher.
    function startEjection(address operator, bytes memory quorums) external onlyWatcher(msg.sender) {
        _takeDeposit(msg.sender);
        operator.startEjection(quorums);
    }

    /// @notice Cancels the ejection process initiated by a watcher.
    function cancelEjectionByWatcher(address operator) external onlyWatcher(msg.sender) {
        _returnDeposit(msg.sender);
        operator.cancelEjection();
    }

    /// @notice Completes the ejection process for an operator. Transfers the deposit back to the watcher.
    function completeEjection(address operator, bytes memory quorums) external onlyWatcher(msg.sender) {
        operator.completeEjection(quorums);
        _tryEjectOperator(operator, quorums);
        _returnDeposit(msg.sender);
    }

    /// @notice Starts the churn process for an operator. Takes a deposit from the watcher.
    function startChurn(address operator, bytes memory quorums) external onlyWatcher(msg.sender) {
        _takeDeposit(msg.sender);
        operator.startChurn(quorums);
    }

    /// @notice Cancels the churn process initiated by a watcher.
    function cancelChurnByWatcher(address operator) external onlyWatcher(msg.sender) {
        operator.cancelChurn();
        _returnDeposit(msg.sender);
    }

    /// @notice Completes the churn process for an operator. Transfers the deposit back to the watcher.
    function completeChurn(address operator, bytes memory quorums) external onlyWatcher(msg.sender) {
        operator.completeChurn(quorums);
        _tryEjectOperator(operator, quorums);
        _returnDeposit(msg.sender);
    }

    /// OPERATOR FUNCTIONS

    /// @notice Cancels the ejection process initiated by the operator. Transfers the deposit to the operator.
    function cancelEjection() external {
        msg.sender.cancelEjection();
        _returnDeposit(msg.sender);
    }

    /// @notice Completes the ejection process for the operator. Transfers the deposit to the operator.
    function cancelChurn(address lowerStakeOperator, bytes memory quorums) external {
        require(
            _isOperatorWeightsGreater(msg.sender, lowerStakeOperator, quorums),
            "EigenDAEjectionManager: Operator does not have greater weights"
        );
        msg.sender.cancelChurn();
        _returnDeposit(msg.sender);
    }

    /// GETTERS

    function getRegistryCoordinator() external view returns (address) {
        return _registryCoordinator;
    }

    function getDepositToken() external view returns (address) {
        return _depositToken;
    }

    function getDepositAmount() external view returns (uint256) {
        return _depositAmount;
    }

    function churnInitiated(address operator) external view returns (bool) {
        return operator.churnInitiated();
    }

    function churnTime(address operator) external view returns (uint64) {
        return EigenDAEjectionLib.churnStorage().operatorProceedingParams[operator].proceedingTime;
    }

    function lastChurnInitiated(address operator) external view returns (uint64) {
        return EigenDAEjectionLib.churnStorage().operatorProceedingParams[operator].lastProceedingInitiated;
    }

    function churnDelay() external view returns (uint64) {
        return EigenDAEjectionLib.churnStorage().delay;
    }

    function churnCooldown() external view returns (uint64) {
        return EigenDAEjectionLib.churnStorage().cooldown;
    }

    function ejectionInitiated(address operator) external view returns (bool) {
        return operator.ejectionInitiated();
    }

    function ejectionTime(address operator) external view returns (uint64) {
        return EigenDAEjectionLib.ejectionStorage().operatorProceedingParams[operator].proceedingTime;
    }

    function lastEjectionInitiated(address operator) external view returns (uint64) {
        return EigenDAEjectionLib.ejectionStorage().operatorProceedingParams[operator].lastProceedingInitiated;
    }

    function ejectionDelay() external view returns (uint64) {
        return EigenDAEjectionLib.ejectionStorage().delay;
    }

    function ejectionCooldown() external view returns (uint64) {
        return EigenDAEjectionLib.ejectionStorage().cooldown;
    }

    /// INTERNAL FUNCTIONS

    function _isOperatorWeightsGreater(address operator1, address operator2, bytes memory quorumNumbers)
        internal
        view
        returns (bool)
    {
        uint96[] memory weights1 = _getOperatorWeights(operator1, quorumNumbers);
        uint96[] memory weights2 = _getOperatorWeights(operator2, quorumNumbers);

        for (uint256 i = 0; i < weights1.length; i++) {
            if (weights1[i] <= weights2[i]) {
                return false;
            }
        }
        return true;
    }

    function _getOperatorWeights(address operator, bytes memory quorumNumbers)
        internal
        view
        returns (uint96[] memory weights)
    {
        weights = new uint96[](quorumNumbers.length);
        for (uint256 i; i < quorumNumbers.length; i++) {
            uint8 quorumNumber = uint8(quorumNumbers[i]);
            weights[i] = IRegistryCoordinator(_registryCoordinator).stakeRegistry().weightOfOperatorForQuorum(
                quorumNumber, operator
            );
        }
    }

    function _takeDeposit(address sender) internal virtual {
        IERC20(_depositToken).safeTransferFrom(sender, address(this), _depositAmount);
    }

    function _returnDeposit(address receiver) internal virtual {
        IERC20(_depositToken).safeTransfer(receiver, _depositAmount);
    }

    function _onlyWatcher(address sender) internal view virtual {
        sender; // TODO: PLACEHOLDER UNTIL ACCESS CONTROL IS IMPLEMENTED
    }

    /// @notice Attempts to eject an operator. If the ejection fails, it catches the error and does nothing.
    function _tryEjectOperator(address operator, bytes memory quorums) internal {
        try IRegistryCoordinator(_registryCoordinator).ejectOperator(operator, quorums) {} catch {}
    }
}
