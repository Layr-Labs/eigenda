// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAEjectionManager} from "src/periphery/ejection/IEigenDAEjectionManager.sol";
import {EigenDAEjectionLib} from "src/periphery/ejection/libraries/EigenDAEjectionLib.sol";
import {EigenDAEjectionManagerStorage} from "src/periphery/ejection/libraries/EigenDAEjectionManagerStorage.sol";
import {SafeERC20, IERC20} from "lib/openzeppelin-contracts/contracts/token/ERC20/utils/SafeERC20.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry} from "lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {IBLSApkRegistry} from "lib/eigenlayer-middleware/src/interfaces/IBLSApkRegistry.sol";
import {BLSSignatureChecker} from "lib/eigenlayer-middleware/src/BLSSignatureChecker.sol";
import {BN254} from "lib/eigenlayer-middleware/src/libraries/BN254.sol";
import {AddressDirectoryLib} from "src/core/libraries/v3/address-directory/AddressDirectoryLib.sol";
import {AddressDirectoryConstants} from "src/core/libraries/v3/address-directory/AddressDirectoryConstants.sol";

import {AccessControlConstants} from "src/core/libraries/v3/access-control/AccessControlConstants.sol";
import {IAccessControl} from "@openzeppelin/contracts/access/IAccessControl.sol";
import {IEigenDASemVer} from "src/core/interfaces/IEigenDASemVer.sol";
import {InitializableLib} from "src/core/libraries/v3/initializable/InitializableLib.sol";

contract EigenDAEjectionManager is IEigenDAEjectionManager, IEigenDASemVer {
    using AddressDirectoryLib for string;
    using EigenDAEjectionLib for address;
    using SafeERC20 for IERC20;

    bytes32 internal constant CANCEL_EJECTION_MESSAGE_IDENTIFIER = keccak256(
        "CancelEjection(address operator,uint64 proceedingTime,uint64 lastProceedingInitiated,bytes quorums,address recipient)"
    );

    modifier initializer() {
        InitializableLib.initialize();
        _;
    }

    /// @notice Initializes the contract via setting the required parameters into storage
    /// @param accessControl_ the EigenDA access control contract used for checking caller role ownership
    ///                       for ejector and owner
    /// @param depositToken_ The ERC20 token used for deposits
    /// @param blsApkKeyRegistry_ The BLS agg pub key registry contract address
    /// @param serviceManager_ The EigenDA AVS ServiceManager contract address
    /// @param depositBaseFeeMultiplier_ The multiplier for calculating deposit amounts based on base fee
    /// @param accessControl_ The access control contract
    /// @param estimatedGasUsedWithoutSig_ Estimated gas for operations without signature verification
    /// @param estimatedGasUsedWithSig_ Estimated gas for operations with signature verification
    function initialize(
        address depositToken_,
        address accessControl_,
        address blsApkKeyRegistry_,
        address serviceManager_,
        address registryCoordinator_,
        uint256 depositBaseFeeMultiplier_,
        uint256 estimatedGasUsedWithoutSig_,
        uint256 estimatedGasUsedWithSig_
    ) external initializer {
        require(depositToken_ != address(0), "EigenDAEjectionManager: deposit token cannot be zero");
        require(accessControl_ != address(0), "EigenDAEjectionManager: access control cannot be zero");
        require(blsApkKeyRegistry_ != address(0), "EigenDAEjectionManager: bls apk key cannot be zero");
        require(serviceManager_ != address(0), "EigenDAEjectionManager: service manager cannot be zero");
        require(registryCoordinator_ != address(0), "EigenDAEjectionManager: registry coordinator cannot be zero");

        EigenDAEjectionManagerStorage.Layout storage s = EigenDAEjectionManagerStorage.layout();
        s.depositToken = depositToken_;
        s.accessControl = accessControl_;
        s.blsApkKeyRegistry = blsApkKeyRegistry_;
        s.serviceManager = serviceManager_;
        s.registryCoordinator = registryCoordinator_;

        s.depositBaseFeeMultiplier = depositBaseFeeMultiplier_;
        s.estimatedGasUsedWithoutSig = estimatedGasUsedWithoutSig_;
        s.estimatedGasUsedWithSig = estimatedGasUsedWithSig_;
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
        IERC20(EigenDAEjectionManagerStorage.layout().depositToken).safeTransferFrom(msg.sender, address(this), amount);
    }

    /// @inheritdoc IEigenDAEjectionManager
    function withdrawEjectorBalance(uint256 amount) external onlyEjector(msg.sender) {
        msg.sender.subtractEjectorBalance(amount);
        IERC20(EigenDAEjectionManagerStorage.layout().depositToken).safeTransfer(msg.sender, amount);
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
        operator.getEjector().addEjectorBalance(depositAmount);
        operator.cancelEjection();
    }

    /// @inheritdoc IEigenDAEjectionManager
    function completeEjection(address operator, bytes memory quorums) external onlyEjector(msg.sender) {
        uint256 depositAmount = operator.getDepositAmount();
        operator.getEjector().addEjectorBalance(depositAmount);
        operator.completeEjection(quorums);
        _tryEjectOperator(operator, quorums);
    }

    /// OPERATOR FUNCTIONS

    /// @inheritdoc IEigenDAEjectionManager
    function cancelEjectionWithSig(
        address operator,
        BN254.G2Point memory apkG2,
        BN254.G1Point memory sigma,
        address recipient
    ) external {
        address blsApkRegistry = EigenDAEjectionManagerStorage.layout().blsApkKeyRegistry;
        (BN254.G1Point memory apk,) = IBLSApkRegistry(blsApkRegistry).getRegisteredPubkey(operator);
        _verifySig(_cancelEjectionMessageHash(operator, recipient), apk, apkG2, sigma);

        uint256 depositAmount = EigenDAEjectionLib.getEjectionRecord(operator).depositAmount;
        operator.cancelEjection();
        _refundGas(recipient, EigenDAEjectionManagerStorage.layout().estimatedGasUsedWithSig, depositAmount);
    }

    /// @inheritdoc IEigenDAEjectionManager
    function cancelEjection() external {
        uint256 depositAmount = EigenDAEjectionLib.getEjectionRecord(msg.sender).depositAmount;
        msg.sender.cancelEjection();
        _refundGas(msg.sender, EigenDAEjectionManagerStorage.layout().estimatedGasUsedWithoutSig, depositAmount);
    }

    /// GETTERS

    /// @inheritdoc IEigenDAEjectionManager
    function getDepositToken() external view returns (address) {
        return EigenDAEjectionManagerStorage.layout().depositToken;
    }

    /// @inheritdoc IEigenDAEjectionManager
    function getEjector(address operator) external view returns (address) {
        return operator.getEjector();
    }

    /// @inheritdoc IEigenDAEjectionManager
    function getEjectorBalance(address ejector) external view returns (uint256) {
        return ejector.getEjectorBalance();
    }

    /// @inheritdoc IEigenDAEjectionManager
    function ejectionTime(address operator) external view returns (uint64) {
        return EigenDAEjectionLib.getEjectionRecord(operator).proceedingTime;
    }

    /// @inheritdoc IEigenDAEjectionManager
    function lastEjectionInitiated(address operator) external view returns (uint64) {
        return operator.lastProceedingInitiated();
    }

    /// @inheritdoc IEigenDAEjectionManager
    function ejectionQuorums(address operator) external view returns (bytes memory) {
        return EigenDAEjectionLib.getEjectionRecord(operator).quorums;
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
    /// @notice Returns the required deposit for initiating an ejection based on a multiple of the base fee of the block.
    function _depositAmount() internal virtual returns (uint256) {
        EigenDAEjectionManagerStorage.Layout storage s = EigenDAEjectionManagerStorage.layout();
        return s.estimatedGasUsedWithSig * block.basefee * s.depositBaseFeeMultiplier;
    }

    function _refundGas(address receiver, uint256 estimatedGasUsed, uint256 depositAmount) internal virtual {
        uint256 estimatedRefund = estimatedGasUsed * block.basefee;
        IERC20(EigenDAEjectionManagerStorage.layout().depositToken).safeTransfer(
            receiver, estimatedRefund > depositAmount ? depositAmount : estimatedRefund
        );
    }

    /// @notice Attempts to eject an operator. If the ejection fails, it catches the error and does nothing.
    function _tryEjectOperator(address operator, bytes memory quorums) internal {
        address registryCoordinator = EigenDAEjectionManagerStorage.layout().registryCoordinator;
        try IRegistryCoordinator(registryCoordinator).ejectOperator(operator, quorums) {} catch {}
    }

    /// @notice Defines a unique identifier for a cancel ejection message to be signed by an operator for the purpose of authorizing a cancellation.
    function _cancelEjectionMessageHash(address operator, address recipient) internal view returns (bytes32) {
        return keccak256(
            abi.encode(
                CANCEL_EJECTION_MESSAGE_IDENTIFIER,
                block.chainid,
                address(this),
                EigenDAEjectionLib.getEjectionRecord(operator),
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
        EigenDAEjectionManagerStorage.Layout storage s = EigenDAEjectionManagerStorage.layout();
        (bool paired, bool valid) =
            BLSSignatureChecker(s.serviceManager).trySignatureAndApkVerification(messageHash, apk, apkG2, sigma);
        require(paired, "EigenDAEjectionManager: Pairing failed");
        require(valid, "EigenDAEjectionManager: Invalid signature");
    }

    function _onlyOwner(address sender) internal view virtual {
        require(
            IAccessControl(EigenDAEjectionManagerStorage.layout().accessControl).hasRole(
                AccessControlConstants.OWNER_ROLE, sender
            ),
            "EigenDAEjectionManager: Caller is not the owner"
        );
    }

    function _onlyEjector(address sender) internal view virtual {
        require(
            IAccessControl(EigenDAEjectionManagerStorage.layout().accessControl).hasRole(
                AccessControlConstants.EJECTOR_ROLE, sender
            ),
            "EigenDAEjectionManager: Caller is not an ejector"
        );
    }
}
