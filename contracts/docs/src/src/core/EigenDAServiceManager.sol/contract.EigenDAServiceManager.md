# EigenDAServiceManager
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/core/EigenDAServiceManager.sol)

**Inherits:**
[EigenDAServiceManagerStorage](/src/core/EigenDAServiceManagerStorage.sol/abstract.EigenDAServiceManagerStorage.md), ServiceManagerBase, BLSSignatureChecker, Pausable

The Service Manager is the central contract of the EigenDA AVS and is responsible for:
- accepting and confirming the signature of bridged V1 batches
- routing rewards submissions to operators
- setting metadata for the AVS


## State Variables
### PAUSED_CONFIRM_BATCH

```solidity
uint8 internal constant PAUSED_CONFIRM_BATCH = 0;
```


## Functions
### onlyBatchConfirmer

when applied to a function, ensures that the function is only callable by the `batchConfirmer`.


```solidity
modifier onlyBatchConfirmer();
```

### constructor


```solidity
constructor(
    IAVSDirectory __avsDirectory,
    IRewardsCoordinator __rewardsCoordinator,
    IRegistryCoordinator __registryCoordinator,
    IStakeRegistry __stakeRegistry,
    IEigenDAThresholdRegistry __eigenDAThresholdRegistry,
    IEigenDARelayRegistry __eigenDARelayRegistry,
    IPaymentVault __paymentVault,
    IEigenDADisperserRegistry __eigenDADisperserRegistry
)
    BLSSignatureChecker(__registryCoordinator)
    ServiceManagerBase(__avsDirectory, __rewardsCoordinator, __registryCoordinator, __stakeRegistry)
    EigenDAServiceManagerStorage(
        __eigenDAThresholdRegistry,
        __eigenDARelayRegistry,
        __paymentVault,
        __eigenDADisperserRegistry
    );
```

### initialize


```solidity
function initialize(
    IPauserRegistry _pauserRegistry,
    uint256 _initialPausedStatus,
    address _initialOwner,
    address[] memory _batchConfirmers,
    address _rewardsInitiator
) public initializer;
```

### confirmBatch

Accepts a batch from the disperser and confirms its signature for V1 bridging


```solidity
function confirmBatch(BatchHeader calldata batchHeader, NonSignerStakesAndSignature memory nonSignerStakesAndSignature)
    external
    onlyWhenNotPaused(PAUSED_CONFIRM_BATCH)
    onlyBatchConfirmer;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`batchHeader`|`BatchHeader`|The batch header to confirm|
|`nonSignerStakesAndSignature`|`NonSignerStakesAndSignature`|The non-signer stakes and signature to confirm the batch with|


### setBatchConfirmer

Toggles a batch confirmer role to allow them to confirm batches


```solidity
function setBatchConfirmer(address _batchConfirmer) external onlyOwner;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_batchConfirmer`|`address`|The address of the batch confirmer to set|


### _setBatchConfirmer

internal function to set a batch confirmer


```solidity
function _setBatchConfirmer(address _batchConfirmer) internal;
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
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`referenceBlockNumber`|`uint32`|The reference block number to get the serve until block for|


### quorumAdversaryThresholdPercentages

Returns an array of bytes where each byte represents the adversary threshold percentage of the quorum at that index for V1 verification


```solidity
function quorumAdversaryThresholdPercentages() external view returns (bytes memory);
```

### quorumConfirmationThresholdPercentages

Returns an array of bytes where each byte represents the confirmation threshold percentage of the quorum at that index for V1 verification


```solidity
function quorumConfirmationThresholdPercentages() external view returns (bytes memory);
```

### quorumNumbersRequired

Returns an array of bytes where each byte represents the number of a required quorum for V1 verification


```solidity
function quorumNumbersRequired() external view returns (bytes memory);
```

### getQuorumAdversaryThresholdPercentage

Returns the adversary threshold percentage for a quorum for V1 verification


```solidity
function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) external view returns (uint8);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`quorumNumber`|`uint8`|The number of the quorum to get the adversary threshold percentage for|


### getQuorumConfirmationThresholdPercentage

Returns the confirmation threshold percentage for a quorum for V1 verification


```solidity
function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) external view returns (uint8);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`quorumNumber`|`uint8`|The number of the quorum to get the confirmation threshold percentage for|


### getIsQuorumRequired

Returns true if a quorum is required for V1 verification


```solidity
function getIsQuorumRequired(uint8 quorumNumber) external view returns (bool);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`quorumNumber`|`uint8`|The number of the quorum to check if it is required for V1 verification|


### getBlobParams

Returns the blob params for a given blob version


```solidity
function getBlobParams(uint16 version) external view returns (VersionedBlobParams memory);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`version`|`uint16`|The version of the blob to get the params for|


