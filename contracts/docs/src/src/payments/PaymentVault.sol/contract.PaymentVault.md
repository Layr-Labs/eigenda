# PaymentVault
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/payments/PaymentVault.sol)

**Inherits:**
OwnableUpgradeable, [PaymentVaultStorage](/src/payments/PaymentVaultStorage.sol/abstract.PaymentVaultStorage.md)

**Author:**
Layr Labs, Inc.


## Functions
### constructor


```solidity
constructor();
```

### receive


```solidity
receive() external payable;
```

### fallback


```solidity
fallback() external payable;
```

### initialize


```solidity
function initialize(
    address _initialOwner,
    uint64 _minNumSymbols,
    uint64 _pricePerSymbol,
    uint64 _priceUpdateCooldown,
    uint64 _globalSymbolsPerPeriod,
    uint64 _reservationPeriodInterval,
    uint64 _globalRatePeriodInterval
) public initializer;
```

### setReservation

This function is called by EigenDA governance to store reservations


```solidity
function setReservation(address _account, Reservation memory _reservation) external onlyOwner;
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


### setPriceParams


```solidity
function setPriceParams(uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown)
    external
    onlyOwner;
```

### setGlobalSymbolsPerPeriod


```solidity
function setGlobalSymbolsPerPeriod(uint64 _globalSymbolsPerPeriod) external onlyOwner;
```

### setReservationPeriodInterval


```solidity
function setReservationPeriodInterval(uint64 _reservationPeriodInterval) external onlyOwner;
```

### setGlobalRatePeriodInterval


```solidity
function setGlobalRatePeriodInterval(uint64 _globalRatePeriodInterval) external onlyOwner;
```

### withdraw


```solidity
function withdraw(uint256 _amount) external onlyOwner;
```

### withdrawERC20


```solidity
function withdrawERC20(IERC20 _token, uint256 _amount) external onlyOwner;
```

### _checkQuorumSplit


```solidity
function _checkQuorumSplit(bytes memory _quorumNumbers, bytes memory _quorumSplits) internal pure;
```

### _deposit


```solidity
function _deposit(address _account, uint256 _amount) internal;
```

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

