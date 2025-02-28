# IPaymentVault
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/interfaces/IPaymentVault.sol)


## Functions
### setReservation

This function is called by EigenDA governance to store reservations


```solidity
function setReservation(address _account, Reservation memory _reservation) external;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_account`|`address`|is the address to submit the reservation for|
|`_reservation`|`Reservation`|is the Reservation struct containing details of the reservation|


### depositOnDemand

This function is called to deposit funds for on demand payment


```solidity
function depositOnDemand(address _account) external payable;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_account`|`address`|is the address to deposit the funds for|


### getReservation

Fetches the current reservation for an account


```solidity
function getReservation(address _account) external view returns (Reservation memory);
```

### getReservations

Fetches the current reservations for a set of accounts


```solidity
function getReservations(address[] memory _accounts) external view returns (Reservation[] memory _reservations);
```

### getOnDemandTotalDeposit

Fetches the current total on demand balance of an account


```solidity
function getOnDemandTotalDeposit(address _account) external view returns (uint80);
```

### getOnDemandTotalDeposits

Fetches the current total on demand balances for a set of accounts


```solidity
function getOnDemandTotalDeposits(address[] memory _accounts) external view returns (uint80[] memory _payments);
```

## Events
### ReservationUpdated
Emitted when a reservation is created or updated


```solidity
event ReservationUpdated(address indexed account, Reservation reservation);
```

### OnDemandPaymentUpdated
Emitted when an on-demand payment is created or updated


```solidity
event OnDemandPaymentUpdated(address indexed account, uint80 onDemandPayment, uint80 totalDeposit);
```

### GlobalSymbolsPerPeriodUpdated
Emitted when globalSymbolsPerPeriod is updated


```solidity
event GlobalSymbolsPerPeriodUpdated(uint64 previousValue, uint64 newValue);
```

### ReservationPeriodIntervalUpdated
Emitted when reservationPeriodInterval is updated


```solidity
event ReservationPeriodIntervalUpdated(uint64 previousValue, uint64 newValue);
```

### GlobalRatePeriodIntervalUpdated
Emitted when globalRatePeriodInterval is updated


```solidity
event GlobalRatePeriodIntervalUpdated(uint64 previousValue, uint64 newValue);
```

### PriceParamsUpdated
Emitted when priceParams are updated


```solidity
event PriceParamsUpdated(
    uint64 previousMinNumSymbols,
    uint64 newMinNumSymbols,
    uint64 previousPricePerSymbol,
    uint64 newPricePerSymbol,
    uint64 previousPriceUpdateCooldown,
    uint64 newPriceUpdateCooldown
);
```

## Structs
### Reservation

```solidity
struct Reservation {
    uint64 symbolsPerSecond;
    uint64 startTimestamp;
    uint64 endTimestamp;
    bytes quorumNumbers;
    bytes quorumSplits;
}
```

### OnDemandPayment

```solidity
struct OnDemandPayment {
    uint80 totalDeposit;
}
```

