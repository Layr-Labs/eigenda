// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";
import {PaymentVaultStorage} from "./PaymentVaultStorage.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/**
 * @title Entrypoint for making reservations and on demand payments for EigenDA.
 * @author Layr Labs, Inc.
**/
contract PaymentVault is OwnableUpgradeable, PaymentVaultStorage {
 
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
        uint64 _globalRatePeriodInterval
    ) public initializer {
        _transferOwnership(_initialOwner);
        
        minNumSymbols = _minNumSymbols;
        pricePerSymbol = _pricePerSymbol;
        priceUpdateCooldown = _priceUpdateCooldown;
        lastPriceUpdateTime = uint64(block.timestamp);

        globalSymbolsPerPeriod = _globalSymbolsPerPeriod;
        reservationPeriodInterval = _reservationPeriodInterval;
        globalRatePeriodInterval = _globalRatePeriodInterval;
    }

    /**
     * @notice This function is called by EigenDA governance to store reservations
     * @param _account is the address to submit the reservation for
     * @param _reservation is the Reservation struct containing details of the reservation
     */
    function setReservation(
        address _account, 
        Reservation memory _reservation
    ) external onlyOwner { 
        _checkQuorumSplit(_reservation.quorumNumbers, _reservation.quorumSplits);
        require(_reservation.endTimestamp > _reservation.startTimestamp, "end timestamp must be greater than start timestamp");
        reservations[_account] = _reservation;
        emit ReservationUpdated(_account, _reservation);
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
        require(block.timestamp >= lastPriceUpdateTime + priceUpdateCooldown, "price update cooldown not surpassed");

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
        emit GlobalSymbolsPerPeriodUpdated(globalSymbolsPerPeriod, _globalSymbolsPerPeriod);
        globalSymbolsPerPeriod = _globalSymbolsPerPeriod;
    }

    function setReservationPeriodInterval(uint64 _reservationPeriodInterval) external onlyOwner {
        emit ReservationPeriodIntervalUpdated(reservationPeriodInterval, _reservationPeriodInterval);
        reservationPeriodInterval = _reservationPeriodInterval;
    }

    function setGlobalRatePeriodInterval(uint64 _globalRatePeriodInterval) external onlyOwner {
        emit GlobalRatePeriodIntervalUpdated(globalRatePeriodInterval, _globalRatePeriodInterval);
        globalRatePeriodInterval = _globalRatePeriodInterval;
    }

    function withdraw(uint256 _amount) external onlyOwner {
        (bool success,) = payable(owner()).call{value: _amount}("");
        require(success);
    }

    function withdrawERC20(address _token, uint256 _amount) external onlyOwner {
        IERC20(_token).transfer(owner(), _amount);
    }

    function _checkQuorumSplit(bytes memory _quorumNumbers, bytes memory _quorumSplits) internal pure {
        require(_quorumNumbers.length == _quorumSplits.length, "arrays must have the same length");
        uint8 total;
        for(uint256 i; i < _quorumSplits.length; ++i) total += uint8(_quorumSplits[i]);
        require(total == 100, "sum of quorumSplits must be 100");
    }

    function _deposit(address _account, uint256 _amount) internal {
        require(_amount <= type(uint80).max, "amount must be less than or equal to 80 bits");
        onDemandPayments[_account].totalDeposit += uint80(_amount);
        emit OnDemandPaymentUpdated(_account, uint80(_amount), onDemandPayments[_account].totalDeposit);
    }

    /// @notice Fetches the current reservation for an account
    function getReservation(address _account) external view returns (Reservation memory) {
        return reservations[_account];
    }

    /// @notice Fetches the current reservations for a set of accounts
    function getReservations(address[] memory _accounts) external view returns (Reservation[] memory _reservations) {
        _reservations = new Reservation[](_accounts.length);
        for(uint256 i; i < _accounts.length; ++i){
            _reservations[i] = reservations[_accounts[i]];
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