// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {InitializableLib} from "src/core/libraries/v3/initializable/InitializableLib.sol";
import {EigenDAEjectionLib, EigenDAEjectionTypes} from "src/periphery/ejection/libraries/EigenDAEjectionLib.sol";
import {SafeERC20, IERC20} from "lib/openzeppelin-contracts/contracts/token/ERC20/utils/SafeERC20.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {IIndexRegistry} from "lib/eigenlayer-middleware/src/interfaces/IIndexRegistry.sol";
import {IStakeRegistry} from "lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {IBLSApkRegistry} from "lib/eigenlayer-middleware/src/interfaces/IBLSApkRegistry.sol";
import {BLSSignatureChecker} from "lib/eigenlayer-middleware/src/BLSSignatureChecker.sol";
import {BN254} from "lib/eigenlayer-middleware/src/libraries/BN254.sol";

contract EigenDAEjectionManager {
    using EigenDAEjectionLib for address;
    using SafeERC20 for IERC20;

    address internal immutable _depositToken;
    uint256 internal immutable _depositAmount;
    address internal immutable _registryCoordinator;
    address internal immutable _signatureVerifier;

    bytes32 internal constant CANCEL_EJECTION_TYPEHASH = keccak256(
        "CancelEjection(address operator, uint64 proceedingTime, uint64 lastProceedingInitiated, bytes quorums, address recipient)"
    );

    constructor(
        address depositToken_,
        uint256 depositAmount_,
        address registryCoordinator_,
        address signatureVerifier_
    ) {
        _depositToken = depositToken_;
        _depositAmount = depositAmount_;
        _registryCoordinator = registryCoordinator_;
        _signatureVerifier = signatureVerifier_;
    }

    function initialize(
        uint64 delay,
        uint64 cooldown
    ) external {
        InitializableLib.setInitializedVersion(1);
        EigenDAEjectionLib.initialize(delay, cooldown);
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

    /// OPERATOR FUNCTIONS

    /// @notice Cancels the ejection process for a given operator with their signature.
    /// @param operator The address of the operator whose ejection is being cancelled.
    /// @param apkG2 The G2 point of the operator's public key.
    /// @param sigma The BLS signature of the operator.
    /// @param recipient The address to which the deposit will be returned.
    function cancelEjectionWithSig(
        address operator,
        BN254.G2Point memory apkG2,
        BN254.G1Point memory sigma,
        address recipient
    ) external {
        (BN254.G1Point memory apk,) =
            IRegistryCoordinator(_registryCoordinator).blsApkRegistry().getRegisteredPubkey(operator);
        _verifySig(_cancelEjectionMessageHash(operator, recipient), apk, apkG2, sigma);

        operator.cancelEjection();
        _returnDeposit(recipient);
    }

    /// @notice Cancels the ejection process initiated by the operator. Transfers the deposit to the operator.
    function cancelEjection() external {
        msg.sender.cancelEjection();
        _returnDeposit(msg.sender);
    }

    /// @notice Finds the lowest stake operator via an exhaustive search and eject them.
    function churn(uint8 quorumId) external {
        bytes32[] memory operators = IRegistryCoordinator(_registryCoordinator).indexRegistry()
            .getOperatorListAtBlockNumber(quorumId, uint32(block.number));

        uint96[] memory stakes = new uint96[](operators.length);
        for (uint256 i; i < operators.length; i++) {
            stakes[i] = IRegistryCoordinator(_registryCoordinator).stakeRegistry().getStakeAtBlockNumber(
                operators[i], 0, uint32(block.number)
            );
        }
        bytes32 lowestStakeOperator = operators[0];
        uint96 lowestStake = stakes[0];
        for (uint256 i = 1; i < operators.length; i++) {
            if (stakes[i] < lowestStake) {
                lowestStake = stakes[i];
                lowestStakeOperator = operators[i];
            }
        }
        _tryEjectOperator(
            IRegistryCoordinator(_registryCoordinator).blsApkRegistry().pubkeyHashToOperator(lowestStakeOperator),
            bytes(abi.encodePacked(quorumId))
        );
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

    function ejectionInitiated(address operator) external view returns (bool) {
        return operator.ejectionInitiated();
    }

    function ejectionTime(address operator) external view returns (uint64) {
        return EigenDAEjectionLib.ejectionParams(operator).proceedingTime;
    }

    function lastEjectionInitiated(address operator) external view returns (uint64) {
        return EigenDAEjectionLib.ejectionParams(operator).lastProceedingInitiated;
    }

    function ejectionQuorums(address operator) external view returns (bytes memory) {
        return EigenDAEjectionLib.ejectionParams(operator).quorums;
    }

    function ejectionDelay() external view returns (uint64) {
        return EigenDAEjectionLib.delay();
    }

    function ejectionCooldown() external view returns (uint64) {
        return EigenDAEjectionLib.cooldown();
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

    function _cancelEjectionMessageHash(address operator, address recipient) internal view returns (bytes32) {
        return keccak256(
            abi.encode(
                CANCEL_EJECTION_TYPEHASH, EigenDAEjectionLib.ejectionParams(operator), recipient
            )
        );
    }

    function _verifySig(
        bytes32 messageHash,
        BN254.G1Point memory apk,
        BN254.G2Point memory apkG2,
        BN254.G1Point memory sigma
    ) internal view {
        (bool paired, bool valid) =
            BLSSignatureChecker(_signatureVerifier).trySignatureAndApkVerification(messageHash, apk, apkG2, sigma);
        require(paired, "EigenDAEjectionManager: Pairing failed");
        require(valid, "EigenDAEjectionManager: Invalid signature");
    }

    function getInitializedVersion() external view returns (uint8) {
        return InitializableLib.getInitializedVersion();
    }
}
