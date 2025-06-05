// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {AccessControlLib} from "src/core/libraries/AccessControlLib.sol";
import {PaymentVaultStorage} from "src/core/libraries/v3/payment/PaymentVaultStorage.sol";
import {PaymentVaultTypes} from "src/core/libraries/v3/payment/PaymentVaultTypes.sol";
import {PaymentVaultLib} from "src/core/libraries/v3/payment/PaymentVaultLib.sol";
import {Constants} from "src/core/libraries/Constants.sol";
import {IPaymentVault} from "src/core/interfaces/IPaymentVault.sol";
import {InitializableLib} from "src/core/libraries/InitializableLib.sol";

contract PaymentVault is IPaymentVault {
    uint64 public immutable SCHEDULE_PERIOD;
    using PaymentVaultLib for PaymentVaultTypes.Reservation;

    modifier onlyOwner() {
        _onlyOwner();
        _;
    }

    modifier onlyQuorumOwner(uint64 quorumId) {
        _onlyQuorumOwner(quorumId);
        _;
    }

    modifier initializer() {
        InitializableLib.setInitializedVersion(1);
        _;
    }

    constructor(uint64 schedulePeriod) {
        if (schedulePeriod == 0) {
            revert IPaymentVault.SchedulePeriodCannotBeZero();
        }
        SCHEDULE_PERIOD = schedulePeriod;
    }

    function initialize(address owner) external initializer {
        AccessControlLib.grantRole(Constants.OWNER_ROLE, owner);
    }

    /// USER

    function decreaseReservation(uint64 quorumId, PaymentVaultTypes.Reservation memory reservation) external {
        PaymentVaultLib.decreaseReservation(quorumId, msg.sender, reservation, SCHEDULE_PERIOD);
    }

    function depositOnDemand(uint64 quorumId, uint256 amount) external {
        PaymentVaultLib.depositOnDemand(quorumId, msg.sender, amount);
    }

    /// OWNER

    function transferOwnership(address newOwner) external onlyOwner {
        if (newOwner == address(0)) {
            revert IPaymentVault.OwnerIsZeroAddress();
        }
        AccessControlLib.transferRole(Constants.OWNER_ROLE, msg.sender, newOwner);
    }

    function initializeQuorum(
        uint64 quorumId,
        address newOwner,
        PaymentVaultTypes.QuorumProtocolConfig memory protocolCfg
    ) external onlyOwner {
        if (AccessControlLib.getRoleMemberCount(Constants.QUORUM_OWNER_ROLE(quorumId)) > 0) {
            revert IPaymentVault.QuorumOwnerAlreadySet(quorumId);
        }
        AccessControlLib.grantRole(Constants.QUORUM_OWNER_ROLE(quorumId), newOwner);
        ps().quorum[quorumId].protocolCfg = protocolCfg;
    }

    function setReservationAdvanceWindow(uint64 quorumId, PaymentVaultTypes.QuorumProtocolConfig memory protocolCfg)
        external
        onlyQuorumOwner(quorumId)
    {
        ps().quorum[quorumId].protocolCfg.reservationAdvanceWindow = protocolCfg.reservationAdvanceWindow;
    }

    function setOnDemandEnabled(uint64 quorumId, PaymentVaultTypes.QuorumProtocolConfig memory protocolCfg)
        external
        onlyQuorumOwner(quorumId)
    {
        ps().quorum[quorumId].protocolCfg.onDemandEnabled = protocolCfg.onDemandEnabled;
    }

    /// QUORUM OWNER

    function createReservation(uint64 quorumId, address account, PaymentVaultTypes.Reservation memory reservation)
        external
        onlyQuorumOwner(quorumId)
    {
        PaymentVaultLib.addReservation(quorumId, account, reservation, SCHEDULE_PERIOD);
    }

    function increaseReservation(uint64 quorumId, address account, PaymentVaultTypes.Reservation memory reservation)
        external
        onlyQuorumOwner(quorumId)
    {
        PaymentVaultLib.increaseReservation(quorumId, account, reservation, SCHEDULE_PERIOD);
    }

    function setQuorumPaymentConfig(uint64 quorumId, PaymentVaultTypes.QuorumConfig memory paymentConfig)
        external
        onlyQuorumOwner(quorumId)
    {
        ps().quorum[quorumId].cfg = paymentConfig;
    }

    function transferQuorumOwnership(uint64 quorumId, address newOwner) external onlyQuorumOwner(quorumId) {
        if (newOwner == address(0)) {
            revert IPaymentVault.OwnerIsZeroAddress();
        }
        AccessControlLib.transferRole(Constants.QUORUM_OWNER_ROLE(quorumId), msg.sender, newOwner);
    }

    /// GETTERS

    function getOnDemandDeposit(uint64 quorumId, address account) external view returns (uint256) {
        return ps().quorum[quorumId].user[account].deposit;
    }

    function getReservation(uint64 quorumId, address account)
        external
        view
        returns (PaymentVaultTypes.Reservation memory)
    {
        return ps().quorum[quorumId].user[account].reservation;
    }

    function getQuorumProtocolConfig(uint64 quorumId)
        external
        view
        returns (PaymentVaultTypes.QuorumProtocolConfig memory)
    {
        return ps().quorum[quorumId].protocolCfg;
    }

    function getQuorumPaymentConfig(uint64 quorumId) external view returns (PaymentVaultTypes.QuorumConfig memory) {
        return ps().quorum[quorumId].cfg;
    }

    function getQuorumReservedSymbols(uint64 quorumId, uint64 period) external view returns (uint64) {
        return ps().quorum[quorumId].reservedSymbols[period];
    }

    /// HELPER

    function ps() internal pure returns (PaymentVaultStorage.Layout storage) {
        return PaymentVaultStorage.layout();
    }

    function _onlyOwner() internal view virtual {
        AccessControlLib.checkRole(Constants.OWNER_ROLE, msg.sender);
    }

    function _onlyQuorumOwner(uint64 quorumId) internal view virtual {
        AccessControlLib.checkRole(Constants.QUORUM_OWNER_ROLE(quorumId), msg.sender);
    }
}
