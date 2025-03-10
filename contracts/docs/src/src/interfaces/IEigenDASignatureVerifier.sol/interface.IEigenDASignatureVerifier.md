# IEigenDASignatureVerifier
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/interfaces/IEigenDASignatureVerifier.sol)

This contract is used for verifying the signatures of either V1 batches or V2 certificates

*This contract is deployed on L1 as the EigenDAServiceManager contract*


## Functions
### checkSignatures

Verifies the BLS signature of a batch certificate


```solidity
function checkSignatures(
    bytes32 msgHash,
    bytes calldata quorumNumbers,
    uint32 referenceBlockNumber,
    NonSignerStakesAndSignature memory params
) external view returns (QuorumStakeTotals memory, bytes32);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`msgHash`|`bytes32`|The hash of the message to verify|
|`quorumNumbers`|`bytes`|The quorum numbers to verify for|
|`referenceBlockNumber`|`uint32`|The reference block number of the signature|
|`params`|`NonSignerStakesAndSignature`|The non-signer stakes and signatures needed to verify the signature|


