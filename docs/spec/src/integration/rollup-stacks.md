# Rollup Stacks

## OP Stack

Links:
- [Our OP Fork](https://github.com/Layr-Labs/optimism)
- [Fork Diff](https://layr-labs.github.io/optimism/)

## Arbitrum Orbit

Our up-to-date Arbitrum Orbit docs for EigenDA V1 are available at [docs.eigenda.xyz](https://docs.eigenda.xyz/integrations-guides/rollup-guides/orbit/overview). EigenDA V2 support is currently work-in-progress; technical design updates can be found [here](./rollup-stacks/4-arbitrum-secure-integration.md).

We maintain fork diffs for the different arbitrum orbit repos that we fork:
- [nitro](https://layr-labs.github.io/nitro/)
- [nitro-contracts](https://layr-labs.github.io/nitro-contracts/)
- [nitro-testnode](https://layr-labs.github.io/nitro-testnode/)
- [nitro-go-ethereum](https://layr-labs.github.io/nitro-go-ethereum/)

## ZKsync ZK Stack

ZKSync-era currently supports and maintains a [validium mode](https://docs.zksync.io/zk-stack/running/validium), which means we don't need to fork ZKSync, unlike the other stacks.

The zksync eigenda client is implemented [here](https://github.com/matter-labs/zksync-era/tree/8ce774d20865a2b5223d26e10e227f0ea7cb3693/core/node/da_clients/src/eigen). It makes use of our [eigenda-client-rs](https://github.com/Layr-Labs/eigenda-client-rs) repo.