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

    function startEjection(address operator, bytes memory quorums) external onlyWatcher(msg.sender) {
        _takeDeposit(msg.sender);
        operator.startEjection(quorums);
    }

    function cancelEjectionByWatcher(address operator) external onlyWatcher(msg.sender) {
        _returnDeposit(msg.sender);
        operator.cancelEjection();
    }

    function completeEjection(address operator, bytes memory quorums) external onlyWatcher(msg.sender) {
        operator.completeEjection(quorums);
        IRegistryCoordinator(_registryCoordinator).ejectOperator(operator, quorums);
        _returnDeposit(msg.sender);
    }

    function startChurn(address operator, bytes memory quorums) external onlyWatcher(msg.sender) {
        _takeDeposit(msg.sender);
        operator.startChurn(quorums);
    }

    function cancelChurnByWatcher(address operator) external onlyWatcher(msg.sender) {
        operator.cancelChurn();
        _returnDeposit(msg.sender);
    }

    function completeChurn(address operator, bytes memory quorums) external onlyWatcher(msg.sender) {
        operator.completeChurn(quorums);
        IRegistryCoordinator(_registryCoordinator).ejectOperator(operator, quorums);
        _returnDeposit(msg.sender);
    }

    /// OPERATOR FUNCTIONS

    function cancelEjection() external {
        msg.sender.cancelEjection();
        _returnDeposit(msg.sender);
    }

    function cancelChurn(address lowerStakeOperator) external {
        lowerStakeOperator; // TODO: OPERATOR MUST PROVE THAT CHURN IS INVALID BY PROVIDING A LOWER STAKE OPERATOR
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
        return EigenDAEjectionLib.churnParams().operatorProceedingParams[operator].proceedingTime;
    }

    function lastChurnInitiated(address operator) external view returns (uint64) {
        return EigenDAEjectionLib.churnParams().operatorProceedingParams[operator].lastProceedingInitiated;
    }

    function churnDelay() external view returns (uint64) {
        return EigenDAEjectionLib.churnParams().proceedingDelay;
    }

    function churnCooldown() external view returns (uint64) {
        return EigenDAEjectionLib.churnParams().proceedingCooldown;
    }

    function ejectionInitiated(address operator) external view returns (bool) {
        return operator.ejectionInitiated();
    }

    function ejectionTime(address operator) external view returns (uint64) {
        return EigenDAEjectionLib.ejectionParams().operatorProceedingParams[operator].proceedingTime;
    }

    function lastEjectionInitiated(address operator) external view returns (uint64) {
        return EigenDAEjectionLib.ejectionParams().operatorProceedingParams[operator].lastProceedingInitiated;
    }

    function ejectionDelay() external view returns (uint64) {
        return EigenDAEjectionLib.ejectionParams().proceedingDelay;
    }

    function ejectionCooldown() external view returns (uint64) {
        return EigenDAEjectionLib.ejectionParams().proceedingCooldown;
    }

    /// INTERNAL FUNCTIONS

    function _takeDeposit(address sender) internal virtual {
        IERC20(_depositToken).safeTransferFrom(sender, address(this), _depositAmount);
    }

    function _returnDeposit(address receiver) internal virtual {
        IERC20(_depositToken).safeTransfer(receiver, _depositAmount);
    }

    function _onlyWatcher(address sender) internal view virtual {
        sender; // TODO: PLACEHOLDER UNTIL ACCESS CONTROL IS IMPLEMENTED
    }
}
