// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {UsageAuthorizationTypes} from "src/core/libraries/v3/usage-authorization/UsageAuthorizationTypes.sol";

/// @title IUsageAuthorizationRegistry
/// @notice Interface for the Usage Authorization Registry, which manages usage authorizations and reservations for quorums.
/// @dev This interface does not contain events since solidity libraries cannot use events from interfaces.
///      Events are defined in the UsageAuthorizationLib library, and can be retrieved off-chain from the ABI of the emitting contract.
interface IUsageAuthorizationRegistry {
    error OnDemandDisabled(uint64 quorumId);

    error ReservationStillActive(uint64 endTimestamp);

    error InvalidStartTimestamp(uint64 startTimestamp);

    error StartTimestampMustMatch(uint64 startTimestamp);

    error ReservationMustIncrease(uint64 endTimestamp, uint64 symbolsPerSecond);

    error ReservationMustDecrease(uint64 endTimestamp, uint64 symbolsPerSecond);

    error TimestampSchedulePeriodMismatch(uint64 timestamp, uint64 schedulePeriod);

    error InvalidReservationPeriod(uint64 startTimestamp, uint64 endTimestamp);

    error ReservationTooLong(uint64 length, uint64 maxLength);

    error NotEnoughSymbolsAvailable(uint64 timestamp, uint64 requiredSymbols, uint64 availableSymbols);

    error AmountTooLarge(uint256 amount, uint256 maxAmount);

    error SchedulePeriodCannotBeZero();

    error OwnerIsZeroAddress();

    error QuorumOwnerAlreadySet(uint64 quorumId);

    /// @notice Increases the on demand deposit balance for the given account in the specified quorum.
    /// @param quorumId The ID of the quorum. The ERC20 token used for on-demand payments can be found in the quorum configuration.
    /// @param account The address of the account to increase the deposit for.
    /// @param amount The amount to increase the deposit by.
    /// @dev This function takes ERC20 tokens from the caller.
    function depositOnDemand(uint64 quorumId, address account, uint256 amount) external;

    /// @notice Decreases the reservation for the caller in the specified quorum.
    /// @param quorumId The ID of the quorum.
    /// @param reservation The proposed new reservation.
    /// @dev An example of a use case of this function is to facilitate renegotiating a reservation with the quorum owner.
    function decreaseReservation(uint64 quorumId, UsageAuthorizationTypes.Reservation memory reservation) external;

    /// @notice Adds a new reservation in the specified quorum.
    /// @param quorumId The ID of the quorum.
    /// @param account The address of the account being reserved for.
    /// @param reservation The reservation to add.
    /// @dev This call will fail if there is a currently active reservation.
    function addReservation(uint64 quorumId, address account, UsageAuthorizationTypes.Reservation memory reservation)
        external;

    /// @notice Increases a reservation in the specified quorum.
    /// @param quorumId The ID of the quorum.
    /// @param account The address of the account being reserved for.
    /// @param reservation The reservation to increase.
    /// @dev This call will fail if the reservation does not exist or if the new reservation
    ///       does not increase the end timestamp or the symbols per second.
    function increaseReservation(
        uint64 quorumId,
        address account,
        UsageAuthorizationTypes.Reservation memory reservation
    ) external;

    /// @notice Used by the quorum owner to set the quorum configuration.
    /// @param quorumId The ID of the quorum.
    /// @param config The new protocol configuration for the quorum.
    function setQuorumPaymentConfig(uint64 quorumId, UsageAuthorizationTypes.QuorumConfig memory config) external;

    /// @notice Used by the quorum owner to transfer ownership of the quorum to a new owner.
    /// @param quorumId The ID of the quorum.
    /// @param newOwner The address of the new owner of the quorum.
    function transferQuorumOwnership(uint64 quorumId, address newOwner) external;

    /// @notice Gets the on demand deposit of a given account for a given quorum.
    /// @param quorumId The ID of the quorum.
    /// @param account The address of the account to get the deposit balance for.
    /// @return The on demand deposit balance of the account in the specified quorum.
    function getOnDemandDeposit(uint64 quorumId, address account) external view returns (uint256);

    /// @notice Gets the reservation of the account that is currently written to contract storage.
    /// @param quorumId The ID of the quorum.
    /// @param account The address of the account to get the reservation for.
    /// @return The reservation for the account in the specified quorum.
    function getReservation(uint64 quorumId, address account)
        external
        view
        returns (UsageAuthorizationTypes.Reservation memory);

    /// @notice Gets the current configuration of the quorum set by the protocol owner.
    /// @param quorumId The ID of the quorum.
    /// @return The protocol configuration for the specified quorum.
    function getQuorumProtocolConfig(uint64 quorumId)
        external
        view
        returns (UsageAuthorizationTypes.QuorumProtocolConfig memory);

    /// @notice Gets the current configuration of the quorum set by the quorum owner.
    /// @param quorumId The ID of the quorum.
    /// @return The payment configuration for the specified quorum.
    function getQuorumPaymentConfig(uint64 quorumId)
        external
        view
        returns (UsageAuthorizationTypes.QuorumConfig memory);

    /// @notice Gets the number of reserved symbols for a given quorum and period.
    /// @param quorumId The ID of the quorum.
    /// @param period The period for which to get the reserved symbols.
    /// @return The number of reserved symbols for the specified quorum and period.
    function getQuorumReservedSymbols(uint64 quorumId, uint64 period) external view returns (uint64);

    /// @notice Transfers ownership of the contract to a new owner.
    /// @param newOwner The address of the new owner.
    function transferOwnership(address newOwner) external;

    /// @notice Initializes a new quorum with the given ID, owner, and protocol configuration.
    /// @param quorumId The ID of the new quorum.
    /// @param newOwner The address of the new owner of the quorum.
    /// @param protocolCfg The protocol configuration for the new quorum.
    function initializeQuorum(
        uint64 quorumId,
        address newOwner,
        UsageAuthorizationTypes.QuorumProtocolConfig memory protocolCfg
    ) external;

    /// @notice Sets the reservation advance window for a quorum. Reservations cannot extend beyond the current time plus this window.
    /// @param quorumId The ID of the quorum.
    /// @param reservationAdvanceWindow The new reservation advance window in seconds.
    function setReservationAdvanceWindow(uint64 quorumId, uint64 reservationAdvanceWindow) external;

    /// @notice Sets the on-demand enabled status for a quorum.
    /// @param quorumId The ID of the quorum.
    /// @param enabled Whether on-demand usage is enabled for the quorum.
    function setOnDemandEnabled(uint64 quorumId, bool enabled) external;
}
