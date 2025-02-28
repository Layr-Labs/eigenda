# EigenDA Contracts
This package contains all smart contracts used to power EigenDA's on-chain operations. This includes both core protocol logic and verification constructs that a rollup can leverage to verify certificate integrity. This project uses both npm and local submodules for dependency management. Most recently published release artifacts can be found [here](https://www.npmjs.com/package/@eigenda/contracts).


### Install
Please ensure you've installed latest [foundry nightly](https://book.getfoundry.sh/getting-started/installation) as well as [yarn](https://classic.yarnpkg.com/lang/en/docs/install). To install packages, run the following commands:
```
cd contracts
yarn install
forge install
```


### Compile
To compile contracts and generate golang ABI bindings, run the following:
```
make compile-contracts

```