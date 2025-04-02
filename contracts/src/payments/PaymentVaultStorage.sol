// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "../interfaces/IEigenDAStructs.sol";
import "../interfaces/IPaymentVault.sol";

/**
 * @title Storage variables for the `PaymentVault` contract.
 * @author Layr Labs, Inc.
 * @notice This storage contract is separate from the logic to simplify the upgrade process.
 */
abstract contract PaymentVaultStorage {
    // Quorum-specific reservation data
    mapping(uint64 => mapping(address => IPaymentVault.Reservation)) public reservations;
    mapping(uint64 => mapping(uint64 => uint64)) public quorumPeriodUsage;
    mapping(uint64 => uint64) public quorumReservationSymbolsPerPeriod;
    mapping(uint64 => bytes) public quorumOwner;
    mapping(uint256 => address) public quorumOwnerAddress;

    // General config
    uint64 public reservationAdvanceWindow;
    uint64 public reservationSchedulePeriod;
    bool public newReservationsEnabled;

    // On-demand payment data
    mapping(address => IPaymentVault.OnDemandPayment) public onDemandPayments;

    // Reservation parameters
    uint64 public minNumSymbols;
    uint64 public pricePerSymbol;
    uint64 public priceUpdateCooldown;
    uint64 public lastPriceUpdateTime;
    uint64 public globalSymbolsPerPeriod;
    uint64 public reservationPeriodInterval;
    uint64 public globalRatePeriodInterval;
    uint256 public maxAdvanceWindow;
    uint256 public maxPermissionlessReservationSymbolsPerSecond;
    uint256 public reservationPricePerSymbol;

    // storage gap for upgradeability
    // slither-disable-next-line shadowing-state
    uint256[39] private __gap;
}