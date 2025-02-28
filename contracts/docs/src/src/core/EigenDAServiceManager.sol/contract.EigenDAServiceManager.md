# EigenDAServiceManager
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/core/EigenDAServiceManager.sol)

**Inherits:**
[EigenDAServiceManagerStorage](/src/core/EigenDAServiceManagerStorage.sol/abstract.EigenDAServiceManagerStorage.md), ServiceManagerBase, BLSSignatureChecker, Pausable

**Author:**
Layr Labs, Inc.

This contract is used for:
- initializing the data store by the disperser
- confirming the data store by the disperser with inferred aggregated signatures of the quorum
- freezing operators as the result of various "challenges"


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

This function is used for
- submitting data availabilty certificates,
- check that the aggregate signature is valid,
- and check whether quorum has been achieved or not.


```solidity
function confirmBatch(BatchHeader calldata batchHeader, NonSignerStakesAndSignature memory nonSignerStakesAndSignature)
    external
    onlyWhenNotPaused(PAUSED_CONFIRM_BATCH)
    onlyBatchConfirmer;
```

### setBatchConfirmer

This function is used for changing the batch confirmer


```solidity
function setBatchConfirmer(address _batchConfirmer) external onlyOwner;
```

### _setBatchConfirmer

changes the batch confirmer


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

### quorumAdversaryThresholdPercentages

Returns the bytes array of quorumAdversaryThresholdPercentages


```solidity
function quorumAdversaryThresholdPercentages() external view returns (bytes memory);
```

### quorumConfirmationThresholdPercentages

Returns the bytes array of quorumAdversaryThresholdPercentages


```solidity
function quorumConfirmationThresholdPercentages() external view returns (bytes memory);
```

### quorumNumbersRequired

Returns the bytes array of quorumsNumbersRequired


```solidity
function quorumNumbersRequired() external view returns (bytes memory);
```

### getQuorumAdversaryThresholdPercentage


```solidity
function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) external view returns (uint8);
```

### getQuorumConfirmationThresholdPercentage

Gets the confirmation threshold percentage for a quorum


```solidity
function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) external view returns (uint8);
```

### getIsQuorumRequired

Checks if a quorum is required


```solidity
function getIsQuorumRequired(uint8 quorumNumber) external view returns (bool);
```

### getBlobParams

Returns the blob params for a given blob version


```solidity
function getBlobParams(uint16 version) external view returns (VersionedBlobParams memory);
```

