// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";
import {PaymentVaultStorage} from "./PaymentVaultStorage.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IPaymentVault} from "../interfaces/IPaymentVault.sol";


/**
 * @title Entrypoint for making reservations and on demand payments for EigenDA.
 * @author Layr Labs, Inc.
**/
contract PaymentVault is OwnableUpgradeable, ReentrancyGuardUpgradeable, PaymentVaultStorage, IPaymentVault {
 
    constructor() {
        _disableInitializers();
    }

    receive() external payable {
        _deposit(msg.sender, msg.value);
    }

    fallback() external payable {
        _deposit(msg.sender, msg.value);
    }

    function initialize(
        address _initialOwner,
        uint64 _minNumSymbols,
        uint64 _pricePerSymbol,
        uint64 _priceUpdateCooldown,
        uint64 _globalSymbolsPerPeriod,
        uint64 _reservationPeriodInterval,
        uint64 _globalRatePeriodInterval,
        uint256 _maxAdvanceWindow,
        uint256 _maxPermissionlessReservationSymbolsPerSecond,
        uint256 _reservationPricePerSymbol,
        uint64 _reservationAdvanceWindow,
        uint64 _reservationSchedulePeriod
    ) public initializer {
        require(_initialOwner != address(0), "Initial owner cannot be zero address");
        require(_reservationSchedulePeriod > 0, "Schedule period must be positive");
        require(_reservationSchedulePeriod <= 1 days, "Schedule period too large");
        require(_reservationAdvanceWindow > 0, "Advance window must be positive");
        require(_minNumSymbols > 0, "Min symbols must be positive");
        require(_pricePerSymbol > 0, "Price per symbol must be positive");
        require(_globalSymbolsPerPeriod > 0, "Global symbols per period must be positive");
        require(_reservationPeriodInterval > 0, "Reservation period interval must be positive");
        require(_reservationPeriodInterval <= 1 days, "Reservation period interval too large");
        require(_globalRatePeriodInterval > 0, "Global rate period interval must be positive");
        require(_globalRatePeriodInterval <= 1 days, "Global rate period interval too large");
        
        __Ownable_init();
        __ReentrancyGuard_init();
        _transferOwnership(_initialOwner);
        
        minNumSymbols = _minNumSymbols;
        pricePerSymbol = _pricePerSymbol;
        priceUpdateCooldown = _priceUpdateCooldown;
        lastPriceUpdateTime = uint64(block.timestamp);

        globalSymbolsPerPeriod = _globalSymbolsPerPeriod;
        reservationPeriodInterval = _reservationPeriodInterval;
        globalRatePeriodInterval = _globalRatePeriodInterval;

        maxAdvanceWindow = _maxAdvanceWindow;
        maxPermissionlessReservationSymbolsPerSecond = _maxPermissionlessReservationSymbolsPerSecond;
        reservationPricePerSymbol = _reservationPricePerSymbol;

        reservationAdvanceWindow = _reservationAdvanceWindow;
        reservationSchedulePeriod = _reservationSchedulePeriod;
        newReservationsEnabled = true;
    }

    /**
     * @notice This function is called to set reservations for a quorum
     * @param _account is the address to submit the reservation for
     * @param _reservation is the Reservation struct containing details of the reservation
     */
    function setReservation(
        address _account, 
        IPaymentVault.Reservation memory _reservation
    ) external { 
        require(newReservationsEnabled, "New reservations are currently disabled");
        require(_reservation.symbolsPerSecond > 0, "Symbols per second must be positive");
        require(quorumOwner[_reservation.quorumNumber].length > 0, "Quorum does not exist");
        
        // Ensure the sender is the designated owner for the quorum
        bytes memory ownerBytes = quorumOwner[_reservation.quorumNumber];
        address owner;
        assembly {
            owner := mload(add(ownerBytes, 20))
        }
        require(msg.sender == owner, "Not authorized");

        require(_reservation.endTimestamp > _reservation.startTimestamp, "End timestamp must be greater than start timestamp");

        // Check if the end timestamp exceeds the maximum allowed advance window
        require(
            _reservation.endTimestamp <= block.timestamp + reservationAdvanceWindow, 
            "End timestamp exceeds maximum advance window"
        );

        // Validate timestamp overflow
        require(_reservation.startTimestamp <= type(uint64).max / reservationSchedulePeriod, "Start timestamp too large");
        require(_reservation.endTimestamp - _reservation.startTimestamp <= type(uint256).max / (1000 * 365 days), "Reservation period too long");

        // Retrieve the existing reservation (if any)
        IPaymentVault.Reservation storage existingReservation = reservations[_reservation.quorumNumber][_account];

        // Calculate the new reservation periods *before* modifying anything
        uint64 startPeriod = _reservation.startTimestamp / reservationSchedulePeriod * reservationSchedulePeriod;
        
        // Validate number of periods
        require(_reservation.endTimestamp >= startPeriod, "End period before start period");
        require((_reservation.endTimestamp - startPeriod) / reservationSchedulePeriod <= 1000, "Too many periods");

        mapping(uint64 => uint64) storage usageMap = quorumPeriodUsage[_reservation.quorumNumber];
        
        // Validate and update in a single loop for efficiency
        for (uint64 currentPeriod = startPeriod; currentPeriod < _reservation.endTimestamp; currentPeriod += reservationSchedulePeriod) {
            uint64 periodEnd = currentPeriod + reservationSchedulePeriod;
            if (periodEnd > _reservation.endTimestamp) {
                periodEnd = _reservation.endTimestamp;
            }

            uint64 newReservedSymbols = _reservation.symbolsPerSecond * (periodEnd - currentPeriod);
            uint64 existingReservedSymbols = 0;

            // Only consider existing reservation if it overlaps with this period
            if (existingReservation.endTimestamp > currentPeriod) {
                uint64 existingPeriodEnd = existingReservation.endTimestamp < periodEnd ? existingReservation.endTimestamp : periodEnd;
                existingReservedSymbols = existingReservation.symbolsPerSecond * (existingPeriodEnd - currentPeriod);
            }

            // Ensure the update does not exceed the quorum limit
            require(
                usageMap[currentPeriod] + newReservedSymbols - existingReservedSymbols <= quorumReservationSymbolsPerPeriod[_reservation.quorumNumber],
                "Exceeds quorum symbols reservation limit for this period"
            );

            // Apply the reservation update
            usageMap[currentPeriod] = usageMap[currentPeriod] + newReservedSymbols - existingReservedSymbols;
        }

        // Store the new reservation only after successful validation
        reservations[_reservation.quorumNumber][_account] = _reservation;

        emit ReservationUpdated(
            _account,
            _reservation.quorumNumber,
            _reservation,
            startPeriod,
            _reservation.endTimestamp
        );
    }

    /**
     * @notice Control if new reservations can be created.
     * @dev Callable only by the contract owner.
     */
    function toggleNewReservations(bool status) external onlyOwner {
        require(newReservationsEnabled != status, "Reservations status already set to the target");
        newReservationsEnabled = status;
        emit NewReservationsStatusChange(status);
    }

    /**
     * @notice Set the owner for a specific quorum
     * @param _quorumNumber The quorum number to set the owner for
     * @param _newOwner The new owner address
     */
    function setQuorumOwner(uint64 _quorumNumber, address _newOwner) external onlyOwner {
        require(_newOwner != address(0), "New owner cannot be zero address");
        require(quorumOwner[_quorumNumber].length > 0, "Quorum does not exist");
        quorumOwner[_quorumNumber] = abi.encodePacked(_newOwner);
    }

    function setPricePerSymbol(uint256 newPrice) external onlyOwner {
        require(newPrice > 0, "Price must be positive");
        uint256 oldPrice = reservationPricePerSymbol;
        reservationPricePerSymbol = newPrice;
        emit ReservationPricePerSymbolUpdated(oldPrice, newPrice);
    }

    function setMaxAdvanceWindow(uint256 newWindow) external onlyOwner {
        require(newWindow > 0, "Window must be positive");
        uint256 oldWindow = maxAdvanceWindow;
        maxAdvanceWindow = newWindow;
        emit MaxAdvanceWindowUpdated(oldWindow, newWindow);
    }

    function setMaxSymbolsPerSecond(uint256 newRate) external onlyOwner {
        require(newRate > 0, "Rate must be positive");
        uint256 oldRate = maxPermissionlessReservationSymbolsPerSecond;
        maxPermissionlessReservationSymbolsPerSecond = newRate;
        emit MaxSymbolsPerSecondUpdated(oldRate, newRate);
    }

    function getCurrentPeriod() public view returns (uint256) {
        return block.timestamp / reservationPeriodInterval;
    }

    /**
     * @notice This function is called to deposit funds for on demand payment
     * @param _account is the address to deposit the funds for
     */
    function depositOnDemand(address _account) external payable {
        _deposit(_account, msg.value);
    }

    function setPriceParams(
        uint64 _minNumSymbols,
        uint64 _pricePerSymbol,
        uint64 _priceUpdateCooldown
    ) external onlyOwner {
        require(_minNumSymbols > 0, "Min symbols must be positive");
        require(_pricePerSymbol > 0, "Price per symbol must be positive");
        require(_priceUpdateCooldown >= 1 hours, "Price update cooldown too short");
        require(block.timestamp >= lastPriceUpdateTime + priceUpdateCooldown, "price update cooldown not surpassed");
        require(_pricePerSymbol <= pricePerSymbol * 2, "Price increase too large"); // Max 100% increase
        require(_pricePerSymbol >= pricePerSymbol / 2, "Price decrease too large"); // Max 50% decrease

        emit PriceParamsUpdated(
            minNumSymbols, _minNumSymbols, 
            pricePerSymbol, _pricePerSymbol, 
            priceUpdateCooldown, _priceUpdateCooldown
        );

        pricePerSymbol = _pricePerSymbol;
        minNumSymbols = _minNumSymbols;
        priceUpdateCooldown = _priceUpdateCooldown;
        lastPriceUpdateTime = uint64(block.timestamp);
    }

    function setGlobalSymbolsPerPeriod(uint64 _globalSymbolsPerPeriod) external onlyOwner {
        require(_globalSymbolsPerPeriod > 0, "Global symbols per period must be positive");
        emit GlobalSymbolsPerPeriodUpdated(globalSymbolsPerPeriod, _globalSymbolsPerPeriod);
        globalSymbolsPerPeriod = _globalSymbolsPerPeriod;
    }

    function setReservationPeriodInterval(uint64 _reservationPeriodInterval) external onlyOwner {
        require(_reservationPeriodInterval > 0, "Reservation period interval must be positive");
        require(_reservationPeriodInterval <= 1 days, "Reservation period interval too large");
        emit ReservationPeriodIntervalUpdated(reservationPeriodInterval, _reservationPeriodInterval);
        reservationPeriodInterval = _reservationPeriodInterval;
    }

    function setGlobalRatePeriodInterval(uint64 _globalRatePeriodInterval) external onlyOwner {
        require(_globalRatePeriodInterval > 0, "Global rate period interval must be positive");
        require(_globalRatePeriodInterval <= 1 days, "Global rate period interval too large");
        emit GlobalRatePeriodIntervalUpdated(globalRatePeriodInterval, _globalRatePeriodInterval);
        globalRatePeriodInterval = _globalRatePeriodInterval;
    }

    function withdraw(uint256 _amount) external onlyOwner nonReentrant {
        require(_amount <= address(this).balance, "Insufficient balance");
        (bool success,) = payable(owner()).call{value: _amount}("");
        require(success, "Transfer failed");
    }

    function _checkQuorumSplit(bytes memory _quorumNumbers, bytes memory _quorumSplits) public pure {
        require(_quorumNumbers.length == _quorumSplits.length, "arrays must have the same length");
        uint8 total;
        for(uint256 i; i < _quorumSplits.length; ++i) total += uint8(_quorumSplits[i]);
        require(total == 100, "sum of quorumSplits must be 100");
    }

    function _deposit(address _account, uint256 _amount) internal {
        require(_amount <= type(uint80).max, "amount must be less than or equal to 80 bits");
        require(_amount > 0, "Deposit amount must be positive");
        onDemandPayments[_account].totalDeposit += uint80(_amount);
        emit OnDemandPaymentUpdated(_account, uint80(_amount), onDemandPayments[_account].totalDeposit);
    }

    /// @notice Fetches the current reservation for a quorum and account
    function getReservation(uint64 _quorumNumber, address _account) external view returns (IPaymentVault.Reservation memory) {
        return reservations[_quorumNumber][_account];
    }

    /// @notice Fetches the current reservations for a set of quorums and accounts
    function getReservations(uint64[] memory _quorums, address[] memory _accounts) external view returns (IPaymentVault.Reservation[][] memory _reservations) {
        _reservations = new IPaymentVault.Reservation[][](_quorums.length);
        for(uint256 i; i < _quorums.length; ++i) {
            _reservations[i] = new IPaymentVault.Reservation[](_accounts.length);
            for(uint256 j; j < _accounts.length; ++j) {
                _reservations[i][j] = reservations[_quorums[i]][_accounts[j]];
            }
        }
    }

    /// @notice Fetches the current total on demand balance of an account
    function getOnDemandTotalDeposit(address _account) external view returns (uint80) {
        return onDemandPayments[_account].totalDeposit;
    }    

    /// @notice Fetches the current total on demand balances for a set of accounts
    function getOnDemandTotalDeposits(address[] memory _accounts) external view returns (uint80[] memory _payments) {
        _payments = new uint80[](_accounts.length);
        for(uint256 i; i < _accounts.length; ++i){
            _payments[i] = onDemandPayments[_accounts[i]].totalDeposit;
        }
    }
}