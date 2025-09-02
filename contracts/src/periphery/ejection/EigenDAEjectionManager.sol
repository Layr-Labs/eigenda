// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAEjectionManager} from "src/periphery/ejection/IEigenDAEjectionManager.sol";
import {EigenDAEjectionLib, EigenDAEjectionTypes} from "src/periphery/ejection/libraries/EigenDAEjectionLib.sol";
import {SafeERC20, IERC20} from "lib/openzeppelin-contracts/contracts/token/ERC20/utils/SafeERC20.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {IEigenDADirectory} from "src/core/interfaces/IEigenDADirectory.sol";
import {IIndexRegistry} from "lib/eigenlayer-middleware/src/interfaces/IIndexRegistry.sol";
import {IStakeRegistry} from "lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {IBLSApkRegistry} from "lib/eigenlayer-middleware/src/interfaces/IBLSApkRegistry.sol";
import {BLSSignatureChecker} from "lib/eigenlayer-middleware/src/BLSSignatureChecker.sol";
import {BN254} from "lib/eigenlayer-middleware/src/libraries/BN254.sol";
import {AddressDirectoryLib} from "src/core/libraries/v3/address-directory/AddressDirectoryLib.sol";
import {AddressDirectoryConstants} from "src/core/libraries/v3/address-directory/AddressDirectoryConstants.sol";

import {AccessControlConstants} from "src/core/libraries/v3/access-control/AccessControlConstants.sol";
import {IAccessControl} from "@openzeppelin/contracts/access/IAccessControl.sol";
import {IEigenDASemVer} from "src/core/interfaces/IEigenDASemVer.sol";

contract EigenDAEjectionManager is IEigenDAEjectionManager, IEigenDASemVer {
    using AddressDirectoryLib for string;
    using EigenDAEjectionLib for address;
    using SafeERC20 for IERC20;

    address internal immutable _depositToken;
    address internal immutable _addressDirectory;
    uint256 internal immutable _estimatedGasUsedWithoutSig;
    uint256 internal immutable _estimatedGasUsedWithSig;
    uint256 internal immutable _depositBaseFeeMultiplier;

    bytes32 internal constant CANCEL_EJECTION_MESSAGE_IDENTIFIER = keccak256(
        "CancelEjection(address operator,uint64 proceedingTime,uint64 lastProceedingInitiated,bytes quorums,address recipient)"
    );

    constructor(
        address depositToken_,
        uint256 depositBaseFeeMultiplier_,
        address addressDirectory_,
        uint256 estimatedGasUsedWithoutSig_,
        uint256 estimatedGasUsedWithSig_
    ) {
        _depositToken = depositToken_;
        _depositBaseFeeMultiplier = depositBaseFeeMultiplier_;
        _addressDirectory = addressDirectory_;
        _estimatedGasUsedWithoutSig = estimatedGasUsedWithoutSig_;
        _estimatedGasUsedWithSig = estimatedGasUsedWithSig_;
    }

    modifier onlyOwner(address sender) {
        _onlyOwner(sender);
        _;
    }

    modifier onlyEjector(address sender) {
        _onlyEjector(sender);
        _;
    }

    /// OWNER FUNCTIONS

    /// @inheritdoc IEigenDAEjectionManager
    function setDelay(uint64 delay) external onlyOwner(msg.sender) {
        EigenDAEjectionLib.setDelay(delay);
    }

    /// @inheritdoc IEigenDAEjectionManager
    function setCooldown(uint64 cooldown) external onlyOwner(msg.sender) {
        EigenDAEjectionLib.setCooldown(cooldown);
    }

    /// EJECTOR FUNCTIONS

    /// @inheritdoc IEigenDAEjectionManager
    function addEjectorBalance(uint256 amount) external onlyEjector(msg.sender) {
        msg.sender.addEjectorBalance(amount);
        IERC20(_depositToken).safeTransferFrom(msg.sender, address(this), amount);
    }

    /// @inheritdoc IEigenDAEjectionManager
    function withdrawEjectorBalance(uint256 amount) external onlyEjector(msg.sender) {
        msg.sender.subtractEjectorBalance(amount);
        IERC20(_depositToken).safeTransfer(msg.sender, amount);
    }

    /// @inheritdoc IEigenDAEjectionManager
    function startEjection(address operator, bytes memory quorums) external onlyEjector(msg.sender) {
        uint256 depositAmount = _depositAmount();
        msg.sender.subtractEjectorBalance(depositAmount);
        operator.startEjection(msg.sender, quorums, depositAmount);
    }

    /// @inheritdoc IEigenDAEjectionManager
    function cancelEjectionByEjector(address operator) external onlyEjector(msg.sender) {
        uint256 depositAmount = operator.getDepositAmount();
        operator.cancelEjection();
        operator.getEjector().addEjectorBalance(depositAmount);
    }

    /// @inheritdoc IEigenDAEjectionManager
    function completeEjection(address operator, bytes memory quorums) external onlyEjector(msg.sender) {
        uint256 depositAmount = operator.getDepositAmount();
        operator.completeEjection(quorums);
        _tryEjectOperator(operator, quorums);
        operator.getEjector().addEjectorBalance(depositAmount);
    }

    /// OPERATOR FUNCTIONS

    /// @inheritdoc IEigenDAEjectionManager
    function cancelEjectionWithSig(
        address operator,
        BN254.G2Point memory apkG2,
        BN254.G1Point memory sigma,
        address recipient
    ) external {
        address blsApkRegistry =
            IEigenDADirectory(_addressDirectory).getAddress(AddressDirectoryConstants.BLS_APK_REGISTRY_NAME.getKey());

        (BN254.G1Point memory apk,) = IBLSApkRegistry(blsApkRegistry).getRegisteredPubkey(operator);
        _verifySig(_cancelEjectionMessageHash(operator, recipient), apk, apkG2, sigma);

        operator.cancelEjection();
        _refundGas(recipient, _estimatedGasUsedWithSig);
    }

    /// @inheritdoc IEigenDAEjectionManager
    function cancelEjection() external {
        msg.sender.cancelEjection();
        _refundGas(msg.sender, _estimatedGasUsedWithoutSig);
    }

    /// GETTERS

    /// @inheritdoc IEigenDAEjectionManager
    function getDepositToken() external view returns (address) {
        return _depositToken;
    }

    /// @inheritdoc IEigenDAEjectionManager
    function getEjector(address operator) external view returns (address) {
        return operator.getEjector();
    }

    /// @inheritdoc IEigenDAEjectionManager
    function ejectionTime(address operator) external view returns (uint64) {
        return EigenDAEjectionLib.ejectionParams(operator).proceedingTime;
    }

    /// @inheritdoc IEigenDAEjectionManager
    function lastEjectionInitiated(address operator) external view returns (uint64) {
        return operator.lastProceedingInitiated();
    }

    /// @inheritdoc IEigenDAEjectionManager
    function ejectionQuorums(address operator) external view returns (bytes memory) {
        return EigenDAEjectionLib.ejectionParams(operator).quorums;
    }

    /// @inheritdoc IEigenDAEjectionManager
    function ejectionDelay() external view returns (uint64) {
        return EigenDAEjectionLib.getDelay();
    }

    /// @inheritdoc IEigenDAEjectionManager
    function ejectionCooldown() external view returns (uint64) {
        return EigenDAEjectionLib.getCooldown();
    }

    /// @inheritdoc IEigenDASemVer
    function semver() external pure returns (uint8 major, uint8 minor, uint8 patch) {
        return (3, 0, 0);
    }

    /// INTERNAL FUNCTIONS

    function _isOperatorWeightsGreater(address operator1, address operator2, bytes memory quorumNumbers)
        internal
        view
        returns (bool)
    {
        uint96[] memory weights1 = _getOperatorWeights(operator1, quorumNumbers);
        uint96[] memory weights2 = _getOperatorWeights(operator2, quorumNumbers);

        for (uint256 i; i < weights1.length; i++) {
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
        address stakeRegistry =
            IEigenDADirectory(_addressDirectory).getAddress(AddressDirectoryConstants.STAKE_REGISTRY_NAME.getKey());
        weights = new uint96[](quorumNumbers.length);
        for (uint256 i; i < quorumNumbers.length; i++) {
            uint8 quorumNumber = uint8(quorumNumbers[i]);

            weights[i] = IStakeRegistry(stakeRegistry).weightOfOperatorForQuorum(quorumNumber, operator);
        }
    }

    /// @notice Returns the required deposit for initiating an ejection based on a multiple of the base fee of the block.
    function _depositAmount() internal virtual returns (uint256) {
        return _estimatedGasUsedWithSig * block.basefee * _depositBaseFeeMultiplier;
    }

    function _refundGas(address receiver, uint256 estimatedGasUsed) internal virtual {
        uint256 estimatedRefund = estimatedGasUsed * block.basefee;
        uint256 depositAmount = EigenDAEjectionLib.ejectionParams(receiver).depositAmount;
        IERC20(_depositToken).safeTransfer(receiver, estimatedRefund > depositAmount ? depositAmount : estimatedRefund);
    }

    /// @notice Attempts to eject an operator. If the ejection fails, it catches the error and does nothing.
    function _tryEjectOperator(address operator, bytes memory quorums) internal {
        address registryCoordinator = IEigenDADirectory(_addressDirectory).getAddress(
            AddressDirectoryConstants.REGISTRY_COORDINATOR_NAME.getKey()
        );
        try IRegistryCoordinator(registryCoordinator).ejectOperator(operator, quorums) {} catch {}
    }

    /// @notice Defines a unique identifier for a cancel ejection message to be signed by an operator for the purpose of authorizing a cancellation.
    function _cancelEjectionMessageHash(address operator, address recipient) internal view returns (bytes32) {
        return keccak256(
            abi.encode(
                CANCEL_EJECTION_MESSAGE_IDENTIFIER,
                block.chainid,
                address(this),
                EigenDAEjectionLib.ejectionParams(operator),
                recipient
            )
        );
    }

    function _verifySig(
        bytes32 messageHash,
        BN254.G1Point memory apk,
        BN254.G2Point memory apkG2,
        BN254.G1Point memory sigma
    ) internal view {
        address signatureVerifier =
            IEigenDADirectory(_addressDirectory).getAddress(AddressDirectoryConstants.SERVICE_MANAGER_NAME.getKey());
        (bool paired, bool valid) =
            BLSSignatureChecker(signatureVerifier).trySignatureAndApkVerification(messageHash, apk, apkG2, sigma);
        require(paired, "EigenDAEjectionManager: Pairing failed");
        require(valid, "EigenDAEjectionManager: Invalid signature");
    }

    function _onlyOwner(address sender) internal view virtual {
        require(
            IAccessControl(
                IEigenDADirectory(_addressDirectory).getAddress(AddressDirectoryConstants.ACCESS_CONTROL_NAME.getKey())
            ).hasRole(AccessControlConstants.OWNER_ROLE, sender),
            "EigenDAEjectionManager: Caller is not the owner"
        );
    }

    function _onlyEjector(address sender) internal view virtual {
        require(
            IAccessControl(
                IEigenDADirectory(_addressDirectory).getAddress(AddressDirectoryConstants.ACCESS_CONTROL_NAME.getKey())
            ).hasRole(AccessControlConstants.EJECTOR_ROLE, sender),
            "EigenDAEjectionManager: Caller is not an ejector"
        );
    }
}
