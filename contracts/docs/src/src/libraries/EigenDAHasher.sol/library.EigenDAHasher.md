# EigenDAHasher
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/libraries/EigenDAHasher.sol)

Library of functions for hashing various EigenDA structs.


## Functions
### hashBatchHashedMetadata

hashes the given metdata into the commitment that will be stored in the contract


```solidity
function hashBatchHashedMetadata(bytes32 batchHeaderHash, bytes32 signatoryRecordHash, uint32 blockNumber)
    internal
    pure
    returns (bytes32);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`batchHeaderHash`|`bytes32`|the hash of the batchHeader|
|`signatoryRecordHash`|`bytes32`|the hash of the signatory record|
|`blockNumber`|`uint32`|the block number at which the batch was confirmed|


### hashBatchHashedMetadata

hashes the given metdata into the commitment that will be stored in the contract


```solidity
function hashBatchHashedMetadata(bytes32 batchHeaderHash, bytes memory confirmationData, uint32 blockNumber)
    internal
    pure
    returns (bytes32);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`batchHeaderHash`|`bytes32`|the hash of the batchHeader|
|`confirmationData`|`bytes`|the confirmation data of the batch|
|`blockNumber`|`uint32`|the block number at which the batch was confirmed|


### hashBatchMetadata

given the batchHeader in the provided metdata, calculates the hash of the batchMetadata


```solidity
function hashBatchMetadata(BatchMetadata memory batchMetadata) internal pure returns (bytes32);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`batchMetadata`|`BatchMetadata`|the metadata of the batch|


### hashBatchHeaderMemory

hashes the given batch header


```solidity
function hashBatchHeaderMemory(BatchHeader memory batchHeader) internal pure returns (bytes32);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`batchHeader`|`BatchHeader`|the batch header to hash|


### hashBatchHeader

hashes the given batch header


```solidity
function hashBatchHeader(BatchHeader calldata batchHeader) internal pure returns (bytes32);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`batchHeader`|`BatchHeader`|the batch header to hash|


### hashReducedBatchHeader

hashes the given reduced batch header


```solidity
function hashReducedBatchHeader(ReducedBatchHeader memory reducedBatchHeader) internal pure returns (bytes32);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`reducedBatchHeader`|`ReducedBatchHeader`|the reduced batch header to hash|


### hashBlobHeader

hashes the given blob header


```solidity
function hashBlobHeader(BlobHeader memory blobHeader) internal pure returns (bytes32);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`blobHeader`|`BlobHeader`|the blob header to hash|


### convertBatchHeaderToReducedBatchHeader

converts a batch header to a reduced batch header


```solidity
function convertBatchHeaderToReducedBatchHeader(BatchHeader memory batchHeader)
    internal
    pure
    returns (ReducedBatchHeader memory);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`batchHeader`|`BatchHeader`|the batch header to convert|


### hashBatchHeaderToReducedBatchHeader

converts the given batch header to a reduced batch header and then hashes it


```solidity
function hashBatchHeaderToReducedBatchHeader(BatchHeader memory batchHeader) internal pure returns (bytes32);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`batchHeader`|`BatchHeader`|the batch header to hash|


### hashBatchHeaderV2

hashes the given V2 batch header


```solidity
function hashBatchHeaderV2(BatchHeaderV2 memory batchHeader) internal pure returns (bytes32);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`batchHeader`|`BatchHeaderV2`|the V2 batch header to hash|


### hashBlobHeaderV2

hashes the given V2 blob header


```solidity
function hashBlobHeaderV2(BlobHeaderV2 memory blobHeader) internal pure returns (bytes32);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`blobHeader`|`BlobHeaderV2`|the V2 blob header to hash|


### hashBlobCertificate

hashes the given V2 blob certificate


```solidity
function hashBlobCertificate(BlobCertificate memory blobCertificate) internal pure returns (bytes32);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`blobCertificate`|`BlobCertificate`|the V2 blob certificate to hash|


