# How to create a Subgraph

## Get the ABI of the contract you want to index

either get it from the build, e.g.

```shell
yq ".abi" contracts-dir/out/Contract.sol/Contract.json > subgraphs/abis/Contract.json
```

## Run the graph CLI command

```shell
cd subgraphs

# install on Linux
yarn global add @graphprotocol/graph-cli # install if u haven't
# or install on MacOS
npm install -g @graphprotocol/graph-cli

graph init --from-contract <contract_addr> --network {goerli,mainnet} --abi abis/Contract.json 
```

And go through the dialog.

## Remove the git files in folder

Suppose you created the subgraph in the folder named `contract-indexing`

```shell
cd contract-indexing

rm -rf .git
```

## Generate bindings and build subgraph

```shell
# on Linux
yarn codegen
yarn build

# on MacOS
npm run codegen
```

## Develop more

Check out the [graph docs](https://thegraph.com/docs/en/network/overview/).
