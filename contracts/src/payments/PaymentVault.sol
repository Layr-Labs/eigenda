// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";
import {PaymentVaultStorage} from "./PaymentVaultStorage.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/**
 * @title Entrypoint for making reservations and on demand payments for EigenDA.
 * @author Layr Labs, Inc.
**/
contract PaymentVault is PaymentVaultStorage, OwnableUpgradeable {
 
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
        uint256 _minNumSymbols,
        uint256 _globalSymbolsPerBin,
        uint256 _pricePerSymbol,
        uint256 _reservationBinInterval,
        uint256 _priceUpdateCooldown,
        uint256 _globalRateBinInterval
    ) public initializer {
        _transferOwnership(_initialOwner);
        
        minNumSymbols = _minNumSymbols;
        globalSymbolsPerBin = _globalSymbolsPerBin;
        pricePerSymbol = _pricePerSymbol;
        reservationBinInterval = _reservationBinInterval;
        priceUpdateCooldown = _priceUpdateCooldown;
        globalRateBinInterval = _globalRateBinInterval;

        lastPriceUpdateTime = block.timestamp;
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
        uint256 _minNumSymbols,
        uint256 _pricePerSymbol,
        uint256 _priceUpdateCooldown
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
        lastPriceUpdateTime = block.timestamp;
    }

    function setGlobalSymbolsPerBin(uint256 _globalSymbolsPerBin) external onlyOwner {
        emit GlobalSymbolsPerBinUpdated(globalSymbolsPerBin, _globalSymbolsPerBin);
        globalSymbolsPerBin = _globalSymbolsPerBin;
    }

    function setReservationBinInterval(uint256 _reservationBinInterval) external onlyOwner {
        emit ReservationBinIntervalUpdated(reservationBinInterval, _reservationBinInterval);
        reservationBinInterval = _reservationBinInterval;
    }

    function setGlobalRateBinInterval(uint256 _globalRateBinInterval) external onlyOwner {
        emit GlobalRateBinIntervalUpdated(globalRateBinInterval, _globalRateBinInterval);
        globalRateBinInterval = _globalRateBinInterval;
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
        onDemandPayments[_account] += _amount;
        emit OnDemandPaymentUpdated(_account, _amount, onDemandPayments[_account]);
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
    function getOnDemandAmount(address _account) external view returns (uint256) {
        return onDemandPayments[_account];
    }    

    /// @notice Fetches the current total on demand balances for a set of accounts
    function getOnDemandAmounts(address[] memory _accounts) external view returns (uint256[] memory _payments) {
        _payments = new uint256[](_accounts.length);
        for(uint256 i; i < _accounts.length; ++i){
            _payments[i] = onDemandPayments[_accounts[i]];
        }
    }
}