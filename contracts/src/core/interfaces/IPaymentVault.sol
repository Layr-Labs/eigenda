// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {PaymentVaultTypes} from "src/core/libraries/v3/payment/PaymentVaultTypes.sol";

interface IPaymentVault {
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

    function getOnDemandDeposit(uint64 quorumId, address account) external view returns (uint256);

    function getReservation(uint64 quorumId, address account)
        external
        view
        returns (PaymentVaultTypes.Reservation memory);

    function getQuorumProtocolConfig(uint64 quorumId)
        external
        view
        returns (PaymentVaultTypes.QuorumProtocolConfig memory);

    function getQuorumPaymentConfig(uint64 quorumId) external view returns (PaymentVaultTypes.QuorumConfig memory);

    function getQuorumReservedSymbols(uint64 quorumId, uint64 period) external view returns (uint64);
}
