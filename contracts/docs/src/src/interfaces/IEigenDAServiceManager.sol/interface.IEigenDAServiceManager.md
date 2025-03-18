# IEigenDAServiceManager
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/interfaces/IEigenDAServiceManager.sol)

**Inherits:**
IServiceManager, [IEigenDAThresholdRegistry](/src/interfaces/IEigenDAThresholdRegistry.sol/interface.IEigenDAThresholdRegistry.md)

The Service Manager is the central contract of the EigenDA AVS and is responsible for:
- accepting and confirming the signature of bridged V1 batches
- routing rewards submissions to operators
- setting metadata for the AVS


## Functions
### confirmBatch

This function is used for
- submitting data availabilty certificates,
- check that the aggregate signature is valid,
- and check whether quorum has been achieved or not.


```solidity
function confirmBatch(
    BatchHeader calldata batchHeader,
    BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature
) external;
```

### batchIdToBatchMetadataHash

mapping between the batchId to the hash of the metadata of the corresponding Batch


```solidity
function batchIdToBatchMetadataHash(uint32 batchId) external view returns (bytes32);
```

### taskNumber

Returns the current batchId


```solidity
function taskNumber() external view returns (uint32);
```

### latestServeUntilBlock

Given a reference block number, returns the block until which operators must serve.


```solidity
function latestServeUntilBlock(uint32 referenceBlockNumber) external view returns (uint32);
```

### BLOCK_STALE_MEASURE

The maximum amount of blocks in the past that the service will consider stake amounts to still be 'valid'.


```solidity
function BLOCK_STALE_MEASURE() external view returns (uint32);
```

## Events
### BatchConfirmed
Emitted when a Batch is confirmed.


```solidity
event BatchConfirmed(bytes32 indexed batchHeaderHash, uint32 batchId);
```

**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`batchHeaderHash`|`bytes32`|The hash of the batch header|
|`batchId`|`uint32`|The ID for the Batch inside of the specified duration (i.e. *not* the globalBatchId)|

### BatchConfirmerStatusChanged
Emitted when a batch confirmer status is updated.


```solidity
event BatchConfirmerStatusChanged(address batchConfirmer, bool status);
```

**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`batchConfirmer`|`address`|The address of the batch confirmer|
|`status`|`bool`|The new status of the batch confirmer|

