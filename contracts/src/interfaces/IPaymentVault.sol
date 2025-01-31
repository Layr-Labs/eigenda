// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

interface IPaymentVault {

    struct Reservation {
        uint64 symbolsPerSecond; // Number of symbols reserved per second
        uint64 startTimestamp;   // timestamp of epoch where reservation begins
        uint64 endTimestamp;     // timestamp of epoch where reservation ends
        uint64 quorumNumber;     // quorum number for the reservation
    }

    struct OnDemandPayment {
        uint80 totalDeposit;
    }

    /// @notice Emitted when a reservation is created or updated
    event ReservationUpdated(
        address indexed account,
        uint64 indexed quorumNumber,
        Reservation reservation,
        uint64 startPeriod,
        uint64 endPeriod
    );
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
    /// @notice Emitted when maxAdvanceWindow is updated
    event MaxAdvanceWindowUpdated(uint256 oldWindow, uint256 newWindow);
    /// @notice Emitted when maxPermissionlessReservationSymbolsPerSecond is updated
    event MaxSymbolsPerSecondUpdated(uint256 oldRate, uint256 newRate);
    /// @notice Emitted when reservationPricePerSymbol is updated
    event ReservationPricePerSymbolUpdated(uint256 oldPrice, uint256 newPrice);
    /// @notice Emitted when new reservations status is changed
    event NewReservationsStatusChange(bool newStatus);

    /**
     * @notice This function is called by EigenDA governance to store reservations
     * @param _account is the address to submit the reservation for
     * @param _reservation is the Reservation struct containing details of the reservation
     */
    function setReservation(
        address _account, 
        Reservation memory _reservation
    ) external payable;

    /**
     * @notice This function is called to deposit funds for on demand payment
     * @param _account is the address to deposit the funds for
     */
    function depositOnDemand(address _account) external payable;

    /**
     * @notice Calculate required payment for symbol rate reservation
     * @param symbolsPerSecond Number of symbols per second to reserve
     * @return payment Required payment in wei
     */
    function calculateReservationPayment(uint256 symbolsPerSecond, uint256 numPeriods) external view returns (uint256);

    /**
     * @notice Control if new reservations can be created
     * @param status The new status for reservations
     */
    function toggleNewReservations(bool status) external;

    /**
     * @notice Set the owner for a specific quorum
     * @param _quorumNumber The quorum number to set the owner for
     * @param _newOwner The new owner address
     */
    function setQuorumOwner(uint64 _quorumNumber, address _newOwner) external;

    /**
     * @notice Set the price per symbol
     * @param newPrice The new price per symbol
     */
    function setPricePerSymbol(uint256 newPrice) external;

    /**
     * @notice Set the maximum advance window
     * @param newWindow The new maximum advance window
     */
    function setMaxAdvanceWindow(uint256 newWindow) external;

    /**
     * @notice Set the maximum symbols per second
     * @param newRate The new maximum symbols per second
     */
    function setMaxSymbolsPerSecond(uint256 newRate) external;

    /**
     * @notice Get the current period
     * @return The current period
     */
    function getCurrentPeriod() external view returns (uint256);

    /**
     * @notice Set price parameters
     * @param _minNumSymbols Minimum number of symbols
     * @param _pricePerSymbol Price per symbol
     * @param _priceUpdateCooldown Price update cooldown
     */
    function setPriceParams(
        uint64 _minNumSymbols,
        uint64 _pricePerSymbol,
        uint64 _priceUpdateCooldown
    ) external;

    /**
     * @notice Set global symbols per period
     * @param _globalSymbolsPerPeriod New global symbols per period
     */
    function setGlobalSymbolsPerPeriod(uint64 _globalSymbolsPerPeriod) external;

    /**
     * @notice Set reservation period interval
     * @param _reservationPeriodInterval New reservation period interval
     */
    function setReservationPeriodInterval(uint64 _reservationPeriodInterval) external;

    /**
     * @notice Set global rate period interval
     * @param _globalRatePeriodInterval New global rate period interval
     */
    function setGlobalRatePeriodInterval(uint64 _globalRatePeriodInterval) external;

    /**
     * @notice Withdraw funds from the contract
     * @param _amount Amount to withdraw
     */
    function withdraw(uint256 _amount) external;

    /**
     * @notice Withdraw ERC20 tokens from the contract
     * @param _token Token to withdraw
     * @param _amount Amount to withdraw
     */
    function withdrawERC20(IERC20 _token, uint256 _amount) external;

    /// @notice Fetches the current reservation for a quorum and account
    function getReservation(uint64 _quorumNumber, address _account) external view returns (Reservation memory);

    /// @notice Fetches the current reservations for a set of quorums and accounts
    function getReservations(uint64[] memory _quorums, address[] memory _accounts) external view returns (Reservation[][] memory _reservations);

    /// @notice Fetches the current total on demand balance of an account
    function getOnDemandTotalDeposit(address _account) external view returns (uint80);

    /// @notice Fetches the current total on demand balances for a set of accounts
    function getOnDemandTotalDeposits(address[] memory _accounts) external view returns (uint80[] memory _payments);
}
