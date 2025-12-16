// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAEjectionManager} from "src/periphery/ejection/IEigenDAEjectionManager.sol";
import {EigenDAEjectionLib} from "src/periphery/ejection/libraries/EigenDAEjectionLib.sol";
import {EigenDAEjectionStorage} from "src/periphery/ejection/libraries/EigenDAEjectionStorage.sol";
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
    /// @param blsApkKeyRegistry_ The BLS agg pub key registry contract address
    /// @param serviceManager_ The EigenDA AVS ServiceManager contract address
    /// @param registryCoordinator_ The EigenDA Registy Coordinator contract address
    function initialize(
        address accessControl_,
        address blsApkKeyRegistry_,
        address serviceManager_,
        address registryCoordinator_
    ) external initializer {
        require(accessControl_ != address(0), "EigenDAEjectionManager: access control cannot be zero");
        require(blsApkKeyRegistry_ != address(0), "EigenDAEjectionManager: bls apk key cannot be zero");
        require(serviceManager_ != address(0), "EigenDAEjectionManager: service manager cannot be zero");
        require(registryCoordinator_ != address(0), "EigenDAEjectionManager: registry coordinator cannot be zero");

        EigenDAEjectionStorage.Layout storage s = EigenDAEjectionStorage.layout();
        s.accessControl = accessControl_;
        s.blsApkKeyRegistry = blsApkKeyRegistry_;
        s.serviceManager = serviceManager_;
        s.registryCoordinator = registryCoordinator_;
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
    function startEjection(address operator, bytes memory quorums) external onlyEjector(msg.sender) {
        operator.startEjection(msg.sender, quorums);
    }

    /// @inheritdoc IEigenDAEjectionManager
    function cancelEjectionByEjector(address operator) external onlyEjector(msg.sender) {
        operator.cancelEjection();
    }

    /// @inheritdoc IEigenDAEjectionManager
    function completeEjection(address operator, bytes memory quorums) external onlyEjector(msg.sender) {
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
        address blsApkRegistry = EigenDAEjectionStorage.layout().blsApkKeyRegistry;
        (BN254.G1Point memory apk,) = IBLSApkRegistry(blsApkRegistry).getRegisteredPubkey(operator);
        _verifySig(_cancelEjectionMessageHash(operator, recipient), apk, apkG2, sigma);

        operator.cancelEjection();
    }

    /// @inheritdoc IEigenDAEjectionManager
    function cancelEjection() external {
        msg.sender.cancelEjection();
    }

    /// GETTERS

    /// @inheritdoc IEigenDAEjectionManager
    function getEjector(address operator) external view returns (address) {
        return operator.getEjector();
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

    /// @notice Attempts to eject an operator. If the ejection fails, it catches the error and does nothing.
    function _tryEjectOperator(address operator, bytes memory quorums) internal {
        address registryCoordinator = EigenDAEjectionStorage.layout().registryCoordinator;
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
        EigenDAEjectionStorage.Layout storage s = EigenDAEjectionStorage.layout();
        (bool paired, bool valid) =
            BLSSignatureChecker(s.serviceManager).trySignatureAndApkVerification(messageHash, apk, apkG2, sigma);
        require(paired, "EigenDAEjectionManager: Pairing failed");
        require(valid, "EigenDAEjectionManager: Invalid signature");
    }

    function _onlyOwner(address sender) internal view virtual {
        require(
            IAccessControl(EigenDAEjectionStorage.layout().accessControl)
                .hasRole(AccessControlConstants.OWNER_ROLE, sender),
            "EigenDAEjectionManager: Caller is not the owner"
        );
    }

    function _onlyEjector(address sender) internal view virtual {
        require(
            IAccessControl(EigenDAEjectionStorage.layout().accessControl)
                .hasRole(AccessControlConstants.EJECTOR_ROLE, sender),
            "EigenDAEjectionManager: Caller is not an ejector"
        );
    }
}
