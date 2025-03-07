# QuorumStakeTotals
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/interfaces/IEigenDAStructs.sol)

The total stakes signed for each quorum returned by the signature verifier


```solidity
struct QuorumStakeTotals {
    uint96[] signedStakeForQuorum;
    uint96[] totalStakeForQuorum;
}
```

