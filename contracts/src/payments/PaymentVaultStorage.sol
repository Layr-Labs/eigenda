// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "../interfaces/IEigenDAStructs.sol";

/**
 * @title Storage variables for the `PaymentVault` contract.
 * @author Layr Labs, Inc.
 * @notice This storage contract is separate from the logic to simplify the upgrade process.
 */
abstract contract PaymentVaultStorage {
    // Quorum-specific reservation data
    mapping(uint64 => mapping(address => Reservation)) public reservations;
    mapping(uint64 => mapping(uint64 => uint64)) public quorumPeriodUsage;
    mapping(uint64 => uint64) public quorumReservationSymbolsPerPeriod;
    mapping(uint64 => bytes) public quorumOwner;
    mapping(uint256 => address) public QuorumOwner;

    // General config
    uint64 public reservationAdvanceWindow;
    uint64 public reservationSchedulePeriod;
    bool public newReservationsEnabled;
    bool private _locked; // Reentrancy guard

    // On-demand payment data
    struct OnDemandPayment {
        uint80 totalDeposit;
    }
    mapping(address => OnDemandPayment) public onDemandPayments;

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

    // Events
    event ReservationUpdated(
        address indexed account,
        uint64 indexed quorumNumber,
        Reservation reservation,
        uint64 startPeriod,
        uint64 endPeriod
    );
    event PriceParamsUpdated(
        uint64 oldMinNumSymbols, uint64 newMinNumSymbols,
        uint64 oldPricePerSymbol, uint64 newPricePerSymbol,
        uint64 oldPriceUpdateCooldown, uint64 newPriceUpdateCooldown
    );
    event GlobalSymbolsPerPeriodUpdated(uint64 oldValue, uint64 newValue);
    event ReservationPeriodIntervalUpdated(uint64 oldValue, uint64 newValue);
    event GlobalRatePeriodIntervalUpdated(uint64 oldValue, uint64 newValue);
    event OnDemandPaymentUpdated(address indexed account, uint80 amount, uint80 totalDeposit);
    event ReservationPricePerSymbolUpdated(uint256 oldValue, uint256 newValue);
    event MaxAdvanceWindowUpdated(uint256 oldValue, uint256 newValue);
    event MaxSymbolsPerSecondUpdated(uint256 oldValue, uint256 newValue);
    event NewReservationsStatusChange(bool newStatus);

    // storage gap for upgradeability
    // slither-disable-next-line shadowing-state
    uint256[39] private __GAP;

    // Reentrancy guard modifiers
    modifier nonReentrant() {
        require(!_locked, "Reentrant call");
        _locked = true;
        _;
        _locked = false;
    }
}