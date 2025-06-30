# EigenDA Contracts
This package contains all smart contracts used to power EigenDA's on-chain operations. This includes both core protocol logic and verification constructs that a rollup can leverage to verify certificate integrity. This project uses both NPM and local submodules for dependency management. Most recently published NPM release artifacts can be found [here](https://www.npmjs.com/package/@eigenda/contracts).

This project is divided into core and periphery. Versions in core represent changes in the eigenDA protocol, while versions in periphery/cert represent changes in an eigenDA blob verification certificate.

### Install
Please ensure you've installed latest [foundry nightly](https://book.getfoundry.sh/getting-started/installation) as well as [yarn](https://classic.yarnpkg.com/lang/en/docs/install). To install dependencies, run the following commands:
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

To just compile contracts, run the following:
```
yarn run build
```

### Testing
Tests are all written using foundry and can be ran via the following commands:
```
yarn run test
```
or 
```
forge test -v
```