# IEigenDAThresholdRegistry
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/interfaces/IEigenDAThresholdRegistry.sol)


## Functions
### quorumAdversaryThresholdPercentages

Returns an array of bytes where each byte represents the adversary threshold percentage of the quorum at that index


```solidity
function quorumAdversaryThresholdPercentages() external view returns (bytes memory);
```

### quorumConfirmationThresholdPercentages

Returns an array of bytes where each byte represents the confirmation threshold percentage of the quorum at that index


```solidity
function quorumConfirmationThresholdPercentages() external view returns (bytes memory);
```

### quorumNumbersRequired

Returns an array of bytes where each byte represents the number of a required quorum


```solidity
function quorumNumbersRequired() external view returns (bytes memory);
```

### getQuorumAdversaryThresholdPercentage

Gets the adversary threshold percentage for a quorum


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

## Events
### VersionedBlobParamsAdded

```solidity
event VersionedBlobParamsAdded(uint16 indexed version, VersionedBlobParams versionedBlobParams);
```

### QuorumAdversaryThresholdPercentagesUpdated

```solidity
event QuorumAdversaryThresholdPercentagesUpdated(
    bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages
);
```

### QuorumConfirmationThresholdPercentagesUpdated

```solidity
event QuorumConfirmationThresholdPercentagesUpdated(
    bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages
);
```

### QuorumNumbersRequiredUpdated

```solidity
event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired);
```

### DefaultSecurityThresholdsV2Updated

```solidity
event DefaultSecurityThresholdsV2Updated(
    SecurityThresholds previousDefaultSecurityThresholdsV2, SecurityThresholds newDefaultSecurityThresholdsV2
);
```

