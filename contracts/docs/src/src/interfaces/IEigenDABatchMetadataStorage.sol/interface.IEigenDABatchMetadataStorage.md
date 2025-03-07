# IEigenDABatchMetadataStorage
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/interfaces/IEigenDABatchMetadataStorage.sol)

This contract is used for storing the batch metadata for a bridged V1 batch

*This contract is deployed on L1 as the EigenDAServiceManager contract*


## Functions
### batchIdToBatchMetadataHash

Returns the batch metadata hash for a given batch id


```solidity
function batchIdToBatchMetadataHash(uint32 batchId) external view returns (bytes32);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`batchId`|`uint32`|The id of the batch to get the metadata hash for|


