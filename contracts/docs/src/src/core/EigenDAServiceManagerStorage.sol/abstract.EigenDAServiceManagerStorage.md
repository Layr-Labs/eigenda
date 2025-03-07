# EigenDAServiceManagerStorage
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/core/EigenDAServiceManagerStorage.sol)

**Inherits:**
[IEigenDAServiceManager](/src/interfaces/IEigenDAServiceManager.sol/interface.IEigenDAServiceManager.md)

This storage contract is separated from the logic to simplify the upgrade process.


## State Variables
### THRESHOLD_DENOMINATOR
The denominator for the threshold percentages


```solidity
uint256 public constant THRESHOLD_DENOMINATOR = 100;
```


### STORE_DURATION_BLOCKS
Unit of measure (in blocks) for which data will be stored for after confirmation.


```solidity
uint32 public constant STORE_DURATION_BLOCKS = 2 weeks / 12 seconds;
```


### BLOCK_STALE_MEASURE
The maximum amount of blocks in the past that the service will consider stake amounts to still be 'valid'.

*To clarify edge cases, the middleware can look `BLOCK_STALE_MEASURE` blocks into the past, i.e. it may trust stakes from the interval
[block.number - BLOCK_STALE_MEASURE, block.number] (specifically, *inclusive* of the block that is `BLOCK_STALE_MEASURE` before the current one)*

*BLOCK_STALE_MEASURE should be greater than the number of blocks till finalization, but not too much greater, as it is the amount of
time that nodes can be active after they have deregistered. The larger it is, the farther back stakes can be used, but the longer operators
have to serve after they've deregistered.
Note that this parameter needs to accommodate the delays which are introduced by the disperser, which are of two types:
- FinalizationBlockDelay: when initializing a batch, the disperser will use a ReferenceBlockNumber which is this many
blocks behind the current block number. This is to ensure that the operator state associated with the reference block
will be stable.
- BatchInterval: the batch itself will only be confirmed after the batch interval has passed.
Currently, we use a FinalizationBlockDelay of 75 blocks and a BatchInterval of 50 blocks,
So using a BLOCK_STALE_MEASURE of 300 should be sufficient to ensure that the batch is not
stale when it is confirmed.*


```solidity
uint32 public constant BLOCK_STALE_MEASURE = 300;
```


### eigenDAThresholdRegistry
The EigenDAThresholdRegistry contract address


```solidity
IEigenDAThresholdRegistry public immutable eigenDAThresholdRegistry;
```


### eigenDARelayRegistry
The EigenDARelayRegistry contract address


```solidity
IEigenDARelayRegistry public immutable eigenDARelayRegistry;
```


### paymentVault
The PaymentVault contract address


```solidity
IPaymentVault public immutable paymentVault;
```


### eigenDADisperserRegistry
The EigenDADisperserRegistry contract address


```solidity
IEigenDADisperserRegistry public immutable eigenDADisperserRegistry;
```


### batchId
The current batchId


```solidity
uint32 public batchId;
```


### batchIdToBatchMetadataHash
mapping between the batchId to the hash of the metadata of the corresponding Batch


```solidity
mapping(uint32 => bytes32) public batchIdToBatchMetadataHash;
```


### isBatchConfirmer
mapping of addressed that are permissioned to confirm batches


```solidity
mapping(address => bool) public isBatchConfirmer;
```


### __GAP
Storage gap for upgradeability


```solidity
uint256[47] private __GAP;
```


## Functions
### constructor


```solidity
constructor(
    IEigenDAThresholdRegistry _eigenDAThresholdRegistry,
    IEigenDARelayRegistry _eigenDARelayRegistry,
    IPaymentVault _paymentVault,
    IEigenDADisperserRegistry _eigenDADisperserRegistry
);
```

