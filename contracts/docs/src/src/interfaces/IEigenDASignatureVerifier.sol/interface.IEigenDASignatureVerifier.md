# IEigenDASignatureVerifier
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/interfaces/IEigenDASignatureVerifier.sol)


## Functions
### checkSignatures


```solidity
function checkSignatures(
    bytes32 msgHash,
    bytes calldata quorumNumbers,
    uint32 referenceBlockNumber,
    NonSignerStakesAndSignature memory params
) external view returns (QuorumStakeTotals memory, bytes32);
```

