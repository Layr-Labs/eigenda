// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

interface IPaymentVault {

    struct Reservation {
        uint64 symbolsPerSecond; // Number of symbols reserved per second
        uint64 startTimestamp;   // timestamp of epoch where reservation begins
        uint64 endTimestamp;     // timestamp of epoch where reservation ends
		bytes quorumNumbers;     // quorum numbers in an ordered bytes array
		bytes quorumSplits;      // quorum splits in a bytes array that correspond to the quorum numbers
    }

    /// @notice Emitted when a reservation is created or updated
    event ReservationUpdated(address indexed account, Reservation reservation);
    /// @notice Emitted when an on-demand payment is created or updated
    event OnDemandPaymentUpdated(address indexed account, uint256 onDemandPayment, uint256 totalDeposit);
    /// @notice Emitted when globalSymbolsPerBin is updated
    event GlobalSymbolsPerBinUpdated(uint256 previousValue, uint256 newValue);
    /// @notice Emitted when reservationBinInterval is updated
    event ReservationBinIntervalUpdated(uint256 previousValue, uint256 newValue);
    /// @notice Emitted when globalRateBinInterval is updated
    event GlobalRateBinIntervalUpdated(uint256 previousValue, uint256 newValue);
    /// @notice Emitted when priceParams are updated
    event PriceParamsUpdated(
        uint256 previousMinNumSymbols, 
        uint256 newMinNumSymbols, 
        uint256 previousPricePerSymbol, 
        uint256 newPricePerSymbol, 
        uint256 previousPriceUpdateCooldown, 
        uint256 newPriceUpdateCooldown
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

    /// @notice Fetches the current reservation for an account
    function getReservation(address _account) external view returns (Reservation memory);

    /// @notice Fetches the current reservations for a set of accounts
    function getReservations(address[] memory _accounts) external view returns (Reservation[] memory _reservations);

    /// @notice Fetches the current total on demand balance of an account
    function getOnDemandAmount(address _account) external view returns (uint256);

    /// @notice Fetches the current total on demand balances for a set of accounts
    function getOnDemandAmounts(address[] memory _accounts) external view returns (uint256[] memory _payments);
}