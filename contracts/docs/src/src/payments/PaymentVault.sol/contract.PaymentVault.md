# PaymentVault
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/payments/PaymentVault.sol)

**Inherits:**
OwnableUpgradeable, [PaymentVaultStorage](/src/payments/PaymentVaultStorage.sol/abstract.PaymentVaultStorage.md)

Entrypoint for making reservations and on demand payments for EigenDA.


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

This function is called by EigenDA governance to set the price parameters


```solidity
function setPriceParams(uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown)
    external
    onlyOwner;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_minNumSymbols`|`uint64`|is the minimum number of symbols to charge for|
|`_pricePerSymbol`|`uint64`|is the price per symbol in wei|
|`_priceUpdateCooldown`|`uint64`|is the cooldown period before the price can be updated again|


### setGlobalSymbolsPerPeriod

This function is called by EigenDA governance to set the global symbols per period


```solidity
function setGlobalSymbolsPerPeriod(uint64 _globalSymbolsPerPeriod) external onlyOwner;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_globalSymbolsPerPeriod`|`uint64`|is the global symbols per period|


### setReservationPeriodInterval

This function is called by EigenDA governance to set the reservation period interval


```solidity
function setReservationPeriodInterval(uint64 _reservationPeriodInterval) external onlyOwner;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_reservationPeriodInterval`|`uint64`|is the reservation period interval|


### setGlobalRatePeriodInterval

This function is called by EigenDA governance to set the global rate period interval


```solidity
function setGlobalRatePeriodInterval(uint64 _globalRatePeriodInterval) external onlyOwner;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_globalRatePeriodInterval`|`uint64`|is the global rate period interval|


### withdraw

This function is called by EigenDA governance to withdraw funds


```solidity
function withdraw(uint256 _amount) external onlyOwner;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_amount`|`uint256`|is the amount to withdraw|


### withdrawERC20

This function is called by EigenDA governance to withdraw ERC20 tokens


```solidity
function withdrawERC20(IERC20 _token, uint256 _amount) external onlyOwner;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_token`|`IERC20`|is the token to withdraw|
|`_amount`|`uint256`|is the amount to withdraw|


### _checkQuorumSplit

Internal function to check that the quorum split is valid


```solidity
function _checkQuorumSplit(bytes memory _quorumNumbers, bytes memory _quorumSplits) internal pure;
```

### _deposit

Internal function to deposit funds for on demand payment


```solidity
function _deposit(address _account, uint256 _amount) internal;
```

### getReservation

Returns the current reservation for an account


```solidity
function getReservation(address _account) external view returns (Reservation memory);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_account`|`address`|is the account to get the reservation for|


### getReservations

Returns the current reservations for a set of accounts


```solidity
function getReservations(address[] memory _accounts) external view returns (Reservation[] memory _reservations);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_accounts`|`address[]`|is the set of accounts to get the reservations for|


### getOnDemandTotalDeposit

Returns the current total on demand balance of an account


```solidity
function getOnDemandTotalDeposit(address _account) external view returns (uint80);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_account`|`address`|is the account to get the total on demand balance for|


### getOnDemandTotalDeposits

Returns the current total on demand balances for a set of accounts


```solidity
function getOnDemandTotalDeposits(address[] memory _accounts) external view returns (uint80[] memory _payments);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_accounts`|`address[]`|is the set of accounts to get the total on demand balances for|


