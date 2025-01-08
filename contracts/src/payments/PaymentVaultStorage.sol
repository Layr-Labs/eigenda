// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IPaymentVault} from "../interfaces/IPaymentVault.sol";

abstract contract PaymentVaultStorage is IPaymentVault {

    /// @notice minimum chargeable size for on-demand payments
    uint64 public minNumSymbols; 
    /// @notice price per symbol in wei
    uint64 public pricePerSymbol; 
    /// @notice cooldown period before the price can be updated again
    uint64 public priceUpdateCooldown;    
    /// @notice timestamp of the last price update
    uint64 public lastPriceUpdateTime; 

    /// @notice maximum number of symbols to disperse per second network-wide for on-demand payments (applied to only ETH and EIGEN)
    uint64 public globalSymbolsPerPeriod;  
    /// @notice reservation period interval 
    uint64 public reservationPeriodInterval;  
    /// @notice global rate period interval
    uint64 public globalRatePeriodInterval;

    /// @notice mapping from user address to current reservation 
    mapping(address => Reservation) public reservations;
    /// @notice mapping from user address to current on-demand payment
    mapping(address => OnDemandPayment) public onDemandPayments;

    uint256[46] private __GAP;
}