# Subgraphs

## Build the subgraph
```shell
yarn install
yarn prepare:preprod-hoodi
yarn codegen
yarn build
```

## Creating new subgraph
Get the ABI of the contract you want to index either get it from the build, e.g.

```shell
yq ".abi" contracts-dir/out/Contract.sol/Contract.json > subgraphs/abis/Contract.json
```

## Run the graph CLI command
```shell
# install on Linux
yarn global add @graphprotocol/graph-cli # install if u haven't
# or install on MacOS
npm install -g @graphprotocol/graph-cli

graph init --from-contract <contract_addr> --abi abis/Contract.json 
```

## Reference documentation
- [goldsky docs](https://docs.goldsky.com/subgraphs/introduction)
- [thegraph docs](https://thegraph.com/docs/en/network/overview/)
