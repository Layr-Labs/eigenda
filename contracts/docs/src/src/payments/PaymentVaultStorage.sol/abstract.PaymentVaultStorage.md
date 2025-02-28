# PaymentVaultStorage
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/payments/PaymentVaultStorage.sol)

**Inherits:**
[IPaymentVault](/src/interfaces/IPaymentVault.sol/interface.IPaymentVault.md)


## State Variables
### minNumSymbols
minimum chargeable size for on-demand payments


```solidity
uint64 public minNumSymbols;
```


### pricePerSymbol
price per symbol in wei


```solidity
uint64 public pricePerSymbol;
```


### priceUpdateCooldown
cooldown period before the price can be updated again


```solidity
uint64 public priceUpdateCooldown;
```


### lastPriceUpdateTime
timestamp of the last price update


```solidity
uint64 public lastPriceUpdateTime;
```


### globalSymbolsPerPeriod
maximum number of symbols to disperse per second network-wide for on-demand payments (applied to only ETH and EIGEN)


```solidity
uint64 public globalSymbolsPerPeriod;
```


### reservationPeriodInterval
reservation period interval


```solidity
uint64 public reservationPeriodInterval;
```


### globalRatePeriodInterval
global rate period interval


```solidity
uint64 public globalRatePeriodInterval;
```


### reservations
mapping from user address to current reservation


```solidity
mapping(address => Reservation) public reservations;
```


### onDemandPayments
mapping from user address to current on-demand payment


```solidity
mapping(address => OnDemandPayment) public onDemandPayments;
```


### __GAP

```solidity
uint256[46] private __GAP;
```


