// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

/**
 * @title IPaymentVault
 * @notice Entrypoint for making reservations and on demand payments for EigenDA.
 */
interface IPaymentVault {

    /**
     * @title Reservation
     * @notice A reservation for a set of quorums
     */
    struct Reservation {
        uint64 symbolsPerSecond; // Number of symbols reserved per second
        uint64 startTimestamp;   // timestamp of epoch where reservation begins
        uint64 endTimestamp;     // timestamp of epoch where reservation ends
        bytes quorumNumbers;     // quorum numbers in an ordered bytes array
        bytes quorumSplits;      // quorum splits in a bytes array that correspond to the quorum numbers
    }

    /**
     * @title OnDemandPayment
     * @notice An on demand payment 
     */
    struct OnDemandPayment {
        uint80 totalDeposit;
    }

    /// @notice Emitted when a reservation is created or updated
    event ReservationUpdated(address indexed account, Reservation reservation);
    /// @notice Emitted when an on-demand payment is created or updated
    event OnDemandPaymentUpdated(address indexed account, uint80 onDemandPayment, uint80 totalDeposit);
    /// @notice Emitted when globalSymbolsPerPeriod is updated
    event GlobalSymbolsPerPeriodUpdated(uint64 previousValue, uint64 newValue);
    /// @notice Emitted when reservationPeriodInterval is updated
    event ReservationPeriodIntervalUpdated(uint64 previousValue, uint64 newValue);
    /// @notice Emitted when globalRatePeriodInterval is updated
    event GlobalRatePeriodIntervalUpdated(uint64 previousValue, uint64 newValue);
    /// @notice Emitted when priceParams are updated
    event PriceParamsUpdated(
        uint64 previousMinNumSymbols, 
        uint64 newMinNumSymbols, 
        uint64 previousPricePerSymbol, 
        uint64 newPricePerSymbol, 
        uint64 previousPriceUpdateCooldown, 
        uint64 newPriceUpdateCooldown
    );

    /**
     * @notice This function is called by EigenDA governance to store reservations
     * @param _account is the address to submit the reservation for
     * @param _reservation is the Reservation struct containing details of the reservation
     */
    function setReservation(
        address _account, 
        Reservation memory _reservation
    ) external;

    /**
     * @notice This function is called to deposit funds for on demand payment
     * @param _account is the address to deposit the funds for
     */
    function depositOnDemand(address _account) external payable;

    /**
     * @notice Returns the current reservation for an account
     * @param _account is the address to get the reservation for
     */
    function getReservation(address _account) external view returns (Reservation memory);

    /**
     * @notice Returns the current reservations for a set of accounts
     * @param _accounts is the set of accounts to get the reservations for
     */
    function getReservations(address[] memory _accounts) external view returns (Reservation[] memory _reservations);

    /**
     * @notice Returns the current total on demand balance of an account
     * @param _account is the address to get the total on demand balance for
     */
    function getOnDemandTotalDeposit(address _account) external view returns (uint80);

    /**
     * @notice Returns the current total on demand balances for a set of accounts
     * @param _accounts is the set of accounts to get the total on demand balances for
     */
    function getOnDemandTotalDeposits(address[] memory _accounts) external view returns (uint80[] memory _payments);
}
