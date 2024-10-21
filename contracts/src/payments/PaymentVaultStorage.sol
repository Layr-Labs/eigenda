// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IPaymentVault} from "../interfaces/IPaymentVault.sol";

abstract contract PaymentVaultStorage is IPaymentVault {
    
    /// @notice reservation bin duration 
    uint256 public immutable reservationBinInterval; 
    /// @notice start timestamp of reservation bins
    uint256 public immutable reservationBinStartTimestamp;
    /// @notice cooldown period before the price can be updated again
    uint256 public immutable priceUpdateCooldown; 

    constructor(
        uint256 _reservationBinInterval,
        uint256 _reservationBinStartTimestamp,
        uint256 _priceUpdateCooldown
    ){
        reservationBinInterval = _reservationBinInterval;
        reservationBinStartTimestamp = _reservationBinStartTimestamp;
        priceUpdateCooldown = _priceUpdateCooldown;
    }

    /// @notice minimum chargeable size for on-demand payments
    uint256 public minChargeableSize;    
    /// @notice maximum number of symbols to disperse per second network-wide for on-demand payments (applied to only ETH and EIGEN)
    uint256 public globalSymbolsPerSecond;     
    /// @notice price per symbol in wei
    uint256 public pricePerSymbol; 
    /// @notice timestamp of the last price update
    uint256 public lastPriceUpdateTime;             

    /// @notice mapping from user address to current reservation 
    mapping(address => Reservation) public reservations;
    /// @notice mapping from user address to current on-demand payment
    mapping(address => uint256) public onDemandPayments;

    uint256[44] private __GAP;
}