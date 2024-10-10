// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin-upgrades/contracts/access/OwnableUpgradeable.sol";
import {PaymentVaultStorage} from "./PaymentVaultStorage.sol";

/**
 * @title Entrypoint for making reservations and on demand payments for EigenDA.
 * @author Layr Labs, Inc.
**/
contract PaymentVault is PaymentVaultStorage, OwnableUpgradeable {
 
    constructor(
        uint256 _reservationBinInterval,
        uint256 _reservationBinStartTimestamp,
        uint256 _priceUpdateCooldown
    ) PaymentVaultStorage(
        _reservationBinInterval,
        _reservationBinStartTimestamp,
        _priceUpdateCooldown
    ){
        _disableInitializers();
    }

    function initialize(
        address _initialOwner,
        uint256 _minChargeableSize,
        uint256 _globalBytesPerSecond
    ) public initializer {
        transferOwnership(_initialOwner);
        minChargeableSize = _minChargeableSize;
        globalBytesPerSecond = _globalBytesPerSecond;
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
        reservations[_account] = _reservation;
        emit ReservationUpdated(_account, _reservation);
    }

    /**
     * @notice This function is called to deposit funds for on demand payment
     * @param _account is the address to deposit the funds for
     */
    function depositOnDemand(address _account) external payable {
		onDemandPayments[_account] += msg.value;
        emit OnDemandPaymentUpdated(_account, msg.value, onDemandPayments[_account]);
    }

    function setMinChargeableSize(uint256 _minChargeableSize) external onlyOwner {
        require(block.timestamp >= lastPriceUpdateTime + priceUpdateCooldown, "price update cooldown not surpassed");
        emit MinChargeableSizeUpdated(minChargeableSize, _minChargeableSize);
        lastPriceUpdateTime = block.timestamp;
        minChargeableSize = _minChargeableSize;
    }

    function setGlobalBytesPerSec(uint256 _globalBytesPerSecond) external onlyOwner {
        emit GlobalBytesPerSecondUpdated(globalBytesPerSecond, _globalBytesPerSecond);
        globalBytesPerSecond = _globalBytesPerSecond;
    }

    function withdraw(uint256 _amount) external onlyOwner {
        (bool success,) = payable(owner()).call{value: _amount}("");
        require(success);
    }

    function _checkQuorumSplit(bytes memory _quorumNumbers, bytes memory _quorumSplits) internal pure {
        require(_quorumNumbers.length == _quorumSplits.length, "arrays must have the same length");
        uint8 total;
        for(uint256 i; i < _quorumSplits.length; ++i) total += uint8(_quorumSplits[i]);
        require(total == 100, "sum of quorumSplits must be 100");
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