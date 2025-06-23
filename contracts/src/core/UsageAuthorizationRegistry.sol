// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {AccessControlLib} from "src/core/libraries/AccessControlLib.sol";
import {UsageAuthorizationStorage} from "src/core/libraries/v3/usage-authorization/UsageAuthorizationStorage.sol";
import {UsageAuthorizationTypes} from "src/core/libraries/v3/usage-authorization/UsageAuthorizationTypes.sol";
import {UsageAuthorizationLib} from "src/core/libraries/v3/usage-authorization/UsageAuthorizationLib.sol";
import {Constants} from "src/core/libraries/Constants.sol";
import {IUsageAuthorizationRegistry} from "src/core/interfaces/IUsageAuthorizationRegistry.sol";
import {InitializableLib} from "src/core/libraries/InitializableLib.sol";

/// @notice UsageAuthorizationRegistry is a contract that replaces the PaymentVault contract to allow for more flexible usage authorization management.
contract UsageAuthorizationRegistry is IUsageAuthorizationRegistry {
    uint64 public immutable SCHEDULE_PERIOD;

    using UsageAuthorizationLib for UsageAuthorizationTypes.Reservation;

    modifier onlyOwner() {
        _onlyOwner();
        _;
    }

    modifier onlyQuorumOwner(uint64 quorumId) {
        _onlyQuorumOwner(quorumId);
        _;
    }

    modifier onlyOnDemandEnabled(uint64 quorumId) {
        if (!ps().quorum[quorumId].protocolCfg.onDemandEnabled) {
            revert IUsageAuthorizationRegistry.OnDemandDisabled(quorumId);
        }
        _;
    }

    modifier initializer() {
        InitializableLib.setInitializedVersion(1);
        _;
    }

    constructor(uint64 schedulePeriod) {
        if (schedulePeriod == 0) {
            revert IUsageAuthorizationRegistry.SchedulePeriodCannotBeZero();
        }
        SCHEDULE_PERIOD = schedulePeriod;
    }

    function initialize(address owner) external initializer {
        AccessControlLib.grantRole(Constants.OWNER_ROLE, owner);
    }

    /// USER

    function decreaseReservation(uint64 quorumId, UsageAuthorizationTypes.Reservation memory reservation) external {
        UsageAuthorizationLib.decreaseReservation(quorumId, msg.sender, reservation, SCHEDULE_PERIOD);
    }

    function depositOnDemand(uint64 quorumId, address account, uint256 amount) external onlyOnDemandEnabled(quorumId) {
        UsageAuthorizationLib.depositOnDemand(quorumId, account, amount, msg.sender);
    }

    /// OWNER

    function transferOwnership(address newOwner) external onlyOwner {
        if (newOwner == address(0)) {
            revert IUsageAuthorizationRegistry.OwnerIsZeroAddress();
        }
        AccessControlLib.transferRole(Constants.OWNER_ROLE, msg.sender, newOwner);
    }

    function initializeQuorum(
        uint64 quorumId,
        address newOwner,
        UsageAuthorizationTypes.QuorumProtocolConfig memory protocolCfg
    ) external onlyOwner {
        if (AccessControlLib.getRoleMemberCount(Constants.QUORUM_OWNER_ROLE(quorumId)) > 0) {
            revert IUsageAuthorizationRegistry.QuorumOwnerAlreadySet(quorumId);
        }
        AccessControlLib.grantRole(Constants.QUORUM_OWNER_ROLE(quorumId), newOwner);
        ps().quorum[quorumId].protocolCfg = protocolCfg;
    }

    function setReservationAdvanceWindow(uint64 quorumId, uint64 reservationAdvanceWindow)
        external
        onlyQuorumOwner(quorumId)
    {
        ps().quorum[quorumId].protocolCfg.reservationAdvanceWindow = reservationAdvanceWindow;
    }

    function setOnDemandEnabled(uint64 quorumId, bool enabled) external onlyQuorumOwner(quorumId) {
        ps().quorum[quorumId].protocolCfg.onDemandEnabled = enabled;
    }

    /// QUORUM OWNER

    function addReservation(uint64 quorumId, address account, UsageAuthorizationTypes.Reservation memory reservation)
        external
        onlyQuorumOwner(quorumId)
    {
        UsageAuthorizationLib.addReservation(quorumId, account, reservation, SCHEDULE_PERIOD);
    }

    function increaseReservation(
        uint64 quorumId,
        address account,
        UsageAuthorizationTypes.Reservation memory reservation
    ) external onlyQuorumOwner(quorumId) {
        UsageAuthorizationLib.increaseReservation(quorumId, account, reservation, SCHEDULE_PERIOD);
    }

    function setQuorumPaymentConfig(uint64 quorumId, UsageAuthorizationTypes.QuorumConfig memory paymentConfig)
        external
        onlyQuorumOwner(quorumId)
    {
        ps().quorum[quorumId].cfg = paymentConfig;
    }

    function transferQuorumOwnership(uint64 quorumId, address newOwner) external onlyQuorumOwner(quorumId) {
        if (newOwner == address(0)) {
            revert IUsageAuthorizationRegistry.OwnerIsZeroAddress();
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
        returns (UsageAuthorizationTypes.Reservation memory)
    {
        return ps().quorum[quorumId].user[account].reservation;
    }

    function getQuorumProtocolConfig(uint64 quorumId)
        external
        view
        returns (UsageAuthorizationTypes.QuorumProtocolConfig memory)
    {
        return ps().quorum[quorumId].protocolCfg;
    }

    function getQuorumPaymentConfig(uint64 quorumId)
        external
        view
        returns (UsageAuthorizationTypes.QuorumConfig memory)
    {
        return ps().quorum[quorumId].cfg;
    }

    function getQuorumReservedSymbols(uint64 quorumId, uint64 period) external view returns (uint64) {
        return ps().quorum[quorumId].reservedSymbols[period];
    }

    /// HELPER

    function ps() internal pure returns (UsageAuthorizationStorage.Layout storage) {
        return UsageAuthorizationStorage.layout();
    }

    function _onlyOwner() internal view virtual {
        AccessControlLib.checkRole(Constants.OWNER_ROLE, msg.sender);
    }

    function _onlyQuorumOwner(uint64 quorumId) internal view virtual {
        AccessControlLib.checkRole(Constants.QUORUM_OWNER_ROLE(quorumId), msg.sender);
    }
}
