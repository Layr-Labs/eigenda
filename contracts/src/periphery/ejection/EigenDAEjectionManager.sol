// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDAEjectionLib, EigenDAEjectionTypes} from "src/periphery/ejection/libraries/EigenDAEjectionLib.sol";
import {SafeERC20, IERC20} from "lib/openzeppelin-contracts/contracts/token/ERC20/utils/SafeERC20.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";

contract EigenDAEjectionManager {
    using EigenDAEjectionLib for bytes32;
    using SafeERC20 for IERC20;

    address internal immutable _depositToken;
    uint256 internal immutable _depositAmount;
    address internal immutable _registryCoordinator;

    constructor(address depositToken_, uint256 depositAmount_, address registryCoordinator_) {
        _depositToken = depositToken_;
        _depositAmount = depositAmount_;
        _registryCoordinator = registryCoordinator_;
    }

    modifier onlyOwner() {
        _onlyOwner();
        _;
    }

    /// OWNER FUNCTIONS

    function startEjection(bytes32 operatorId) external onlyOwner {
        _takeDeposit(msg.sender);
        operatorId.startEjection();
    }

    function completeEjection(bytes32 operatorId) external onlyOwner {
        _returnDeposit(msg.sender);
        operatorId.completeEjection();
    }

    function startChurn(bytes32 operatorId) external onlyOwner {
        _takeDeposit(msg.sender);
        operatorId.startChurn();
    }

    function completeChurn(bytes32 operatorId) external onlyOwner {
        _returnDeposit(msg.sender);
        operatorId.completeChurn();
    }

    /// OPERATOR FUNCTIONS

    function cancelEjection(bytes32 operatorId, address recipient, bytes32 salt, bytes memory signature) external {
        operatorId.consumeSignature(recipient, salt, signature, EigenDAEjectionLib.ejectionParams());
        _returnDeposit(recipient);
        operatorId.cancelEjection();
    }

    function cancelChurn(bytes32 operatorId, address recipient, bytes32 salt, bytes memory signature) external {
        operatorId.consumeSignature(recipient, salt, signature, EigenDAEjectionLib.churnParams());
        _returnDeposit(recipient);
        operatorId.cancelChurn();
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

    function churnInitiated(bytes32 operatorId) external view returns (bool) {
        return operatorId.churnInitiated();
    }

    function churnTime(bytes32 operatorId) external view returns (uint64) {
        return EigenDAEjectionLib.churnParams().operatorProceedingParams[operatorId].proceedingTime;
    }

    function lastChurnInitiated(bytes32 operatorId) external view returns (uint64) {
        return EigenDAEjectionLib.churnParams().operatorProceedingParams[operatorId].lastProceedingInitiated;
    }

    function churnDelay() external view returns (uint64) {
        return EigenDAEjectionLib.churnParams().proceedingDelay;
    }

    function churnCooldown() external view returns (uint64) {
        return EigenDAEjectionLib.churnParams().proceedingCooldown;
    }

    function churnSaltConsumed(bytes32 operatorId, bytes32 salt) external view returns (bool) {
        return EigenDAEjectionLib.churnParams().operatorProceedingParams[operatorId].salts[salt];
    }

    function ejectionInitiated(bytes32 operatorId) external view returns (bool) {
        return operatorId.ejectionInitiated();
    }

    function ejectionTime(bytes32 operatorId) external view returns (uint64) {
        return EigenDAEjectionLib.ejectionParams().operatorProceedingParams[operatorId].proceedingTime;
    }

    function lastEjectionInitiated(bytes32 operatorId) external view returns (uint64) {
        return EigenDAEjectionLib.ejectionParams().operatorProceedingParams[operatorId].lastProceedingInitiated;
    }

    function ejectionDelay() external view returns (uint64) {
        return EigenDAEjectionLib.ejectionParams().proceedingDelay;
    }

    function ejectionCooldown() external view returns (uint64) {
        return EigenDAEjectionLib.ejectionParams().proceedingCooldown;
    }

    function ejectionSaltConsumed(bytes32 operatorId, bytes32 salt) external view returns (bool) {
        return EigenDAEjectionLib.ejectionParams().operatorProceedingParams[operatorId].salts[salt];
    }

    /// INTERNAL FUNCTIONS

    function _takeDeposit(address sender) internal virtual {
        IERC20(_depositToken).safeTransferFrom(sender, address(this), _depositAmount);
    }

    function _returnDeposit(address receiver) internal virtual {
        IERC20(_depositToken).safeTransfer(receiver, _depositAmount);
    }

    function _onlyOwner() internal view virtual {
        this; // PLACEHOLDER UNTIL ACCESS CONTROL IS IMPLEMENTED
    }
}
