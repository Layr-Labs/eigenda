# EigenDA Contracts
This package contains all smart contracts used to power EigenDA's on-chain operations. This includes both core protocol logic and verification constructs that a rollup can leverage to verify certificate integrity. This project uses both NPM and local submodules for dependency management. Most recently published NPM release artifacts can be found [here](https://www.npmjs.com/package/@eigenda/contracts).

This project is divided into core and integrations. Versions in core represent changes in the EigenDA protocol, while versions in integrations/cert represent changes in EigenDA blob verification certificate types.

### Install
Please ensure you've installed latest [foundry nightly](https://book.getfoundry.sh/getting-started/installation) as well as [yarn](https://classic.yarnpkg.com/lang/en/docs/install). To install dependencies, run the following commands:
```
cd contracts
yarn install
forge install
```


### Compile

To compile contracts, run the following:
```
make compile
```

## Generate Golang Bindings

To generate golang ABI bindings (both ABI V1 and V2), run the following (which will compile the contracts as a dependency):
```
make bindings
```

The `generate-bindings.sh` script specifies which contracts to build bindings for. It must be manually updated whenever contract targets are added or removed.

### Testing
Tests are all written using foundry and can be ran via the following commands:
```
yarn run test
```
or 
```
forge test -v
```

## (ERC-7201) Namespaced Storage Schemas

Some EigenDA core contracts implement the [ERC-7201](https://eips.ethereum.org/EIPS/eip-7201) namespaced storage standard. Contracts using this standard follow the namespace pattern based on the storage library name, e.g a library named `ConfigRegistryStorage` would have a namespaced storage ID preimage of `config.registry.storage`.