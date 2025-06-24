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

    /// @inheritdoc IUsageAuthorizationRegistry
    function decreaseReservation(uint64 quorumId, UsageAuthorizationTypes.Reservation memory reservation) external {
        UsageAuthorizationLib.decreaseReservation(quorumId, msg.sender, reservation, SCHEDULE_PERIOD);
    }

    /// @inheritdoc IUsageAuthorizationRegistry
    function depositOnDemand(uint64 quorumId, address account, uint256 amount) external onlyOnDemandEnabled(quorumId) {
        UsageAuthorizationLib.depositOnDemand(quorumId, account, amount, msg.sender);
    }

    /// OWNER

    /// @inheritdoc IUsageAuthorizationRegistry
    function transferOwnership(address newOwner) external onlyOwner {
        if (newOwner == address(0)) {
            revert IUsageAuthorizationRegistry.OwnerIsZeroAddress();
        }
        AccessControlLib.transferRole(Constants.OWNER_ROLE, msg.sender, newOwner);
    }

    /// @inheritdoc IUsageAuthorizationRegistry
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

    /// @inheritdoc IUsageAuthorizationRegistry
    function setQuorumProtocolConfig(uint64 quorumId, UsageAuthorizationTypes.QuorumProtocolConfig memory protocolCfg)
        external
        onlyOwner
    {
        ps().quorum[quorumId].protocolCfg = protocolCfg;
    }

    /// QUORUM OWNER

    /// @inheritdoc IUsageAuthorizationRegistry
    function addReservation(uint64 quorumId, address account, UsageAuthorizationTypes.Reservation memory reservation)
        external
        onlyQuorumOwner(quorumId)
    {
        UsageAuthorizationLib.addReservation(quorumId, account, reservation, SCHEDULE_PERIOD);
    }

    /// @inheritdoc IUsageAuthorizationRegistry
    function increaseReservation(
        uint64 quorumId,
        address account,
        UsageAuthorizationTypes.Reservation memory reservation
    ) external onlyQuorumOwner(quorumId) {
        UsageAuthorizationLib.increaseReservation(quorumId, account, reservation, SCHEDULE_PERIOD);
    }

    /// @inheritdoc IUsageAuthorizationRegistry
    function setQuorumConfig(uint64 quorumId, UsageAuthorizationTypes.QuorumConfig memory paymentConfig)
        external
        onlyQuorumOwner(quorumId)
    {
        ps().quorum[quorumId].cfg = paymentConfig;
    }

    /// @inheritdoc IUsageAuthorizationRegistry
    function transferQuorumOwnership(uint64 quorumId, address newOwner) external onlyQuorumOwner(quorumId) {
        if (newOwner == address(0)) {
            revert IUsageAuthorizationRegistry.OwnerIsZeroAddress();
        }
        AccessControlLib.transferRole(Constants.QUORUM_OWNER_ROLE(quorumId), msg.sender, newOwner);
    }

    /// GETTERS

    /// @inheritdoc IUsageAuthorizationRegistry
    function getOnDemandDeposit(uint64 quorumId, address account) external view returns (uint256) {
        return ps().quorum[quorumId].user[account].deposit;
    }

    /// @inheritdoc IUsageAuthorizationRegistry
    function getReservation(uint64 quorumId, address account)
        external
        view
        returns (UsageAuthorizationTypes.Reservation memory)
    {
        return ps().quorum[quorumId].user[account].reservation;
    }

    /// @inheritdoc IUsageAuthorizationRegistry
    function getQuorumProtocolConfig(uint64 quorumId)
        external
        view
        returns (UsageAuthorizationTypes.QuorumProtocolConfig memory)
    {
        return ps().quorum[quorumId].protocolCfg;
    }

    /// @inheritdoc IUsageAuthorizationRegistry
    function getQuorumPaymentConfig(uint64 quorumId)
        external
        view
        returns (UsageAuthorizationTypes.QuorumConfig memory)
    {
        return ps().quorum[quorumId].cfg;
    }

    /// @inheritdoc IUsageAuthorizationRegistry
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
