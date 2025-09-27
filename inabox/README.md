# Inabox Devnet + E2E Tests

Inabox is a local eigenda devnet, that can be used in 2 modes:
1. short-running devnet for [e2e-tests](#run-e2e-tests-against-inabox)
2. long-running devnet for [local interactions](#run-long-lived-local-inabox-devnet)

Make sure to look at the Makefile, which is well documented.

## Dependencies
- Ensure all submodules are initialized and checked out
    ```
    $ git submodule update --init --recursive
    ```
- Docker is installed. [Instructions for installing docker](https://www.docker.com/products/docker-desktop/).
- We use mise as a dependency manager. Most dependencies are defined in our [mise.toml](../mise.toml) file. [Install mise](https://mise.jdx.dev/getting-started.html) and run `mise install` to install them.
- Two dependencies are not available via mise, so need to be installed independently: 
  - Localstack CLI is installed (simulates AWS stack on local machine; we also provide instructions for running localstack from docker without the CLI):
      ```
      $ brew install localstack/tap/localstack-cli
      ```
  - `aws` CLI  (install instructions [here](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html))

## Run E2E tests against inabox

You can run the end-to-end test suite by running the following command:
```
make run-e2e-tests
```

## Run long-lived local inabox devnet

You can run a long-lived local inabox devnet by running the following command:
```
make start-inabox
```
This will start the devnet and print this log output:
```
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@                     INABOX IS RUNNING!                         @
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

export these variables:
ETH_RPC_URL=http://localhost:8545
EIGENDA_DIRECTORY_ADDR=0x1613beB3B2C4f22Ee086B2b38C1476A3cE7f78E8
EIGENDA_DISPERSER_V1_URL=localhost:32004
EIGENDA_DISPERSER_V2_URL=localhost:32005

You can query other contract addresses from the directory:
cast call $EIGENDA_DIRECTORY_ADDR "getAddress(string)(address)" "CERT_VERIFIER_ROUTER"
You can query the disperser v2 by using:
grpcurl -plaintext $EIGENDA_DISPERSER_V2_URL list

Infra components (anvil, graph, aws localstack) are managed by docker.
Run 'docker ps' to see and manage them.

EigenDA services (disperser, validators, etc) are ran as local processes.
Their config is available under /Users/samlaf/devel/eigenda/inabox/testdata/_latest/envs
Their logs are available under /Users/samlaf/devel/eigenda/inabox/testdata/_latest/logs
```

It can also be stopped by running:
```
make stop-inabox
```

### Custom inabox devnet

If you need to make modifications to the template config file used by inabox, then instead run:
```
make new-testdata-dir
# make modifications to `./testdata/_latest/config.yaml`
make start-infra
make start-services
```


### Send traffic via grpcurl

Disperse a blob:
```
# This command uses `grpcurl`, a tool to send gRPC request in cli, and `kzgpad` to encode payloads into blobs.
# To install `grpcurl`, run `brew install grpcurl` or `go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest`.
# To install `kzgpad`, run `go install github.com/Layr-Labs/eigenda/tools/kzgpad@latest`

# From top level eigenda directory
$ grpcurl -plaintext -d '{"data": "'$(kzgpad -e hello)'"}' \
  localhost:32003 disperser.Disperser/DisperseBlob
```

This will return a message in the following form:

```
{
  "result": "PROCESSING",
  "requestId": "$REQUEST_ID"
}
```

Look for logs such as the following to indicate that the disperser has successfully confirmed the batch:
```
TRACE[10-12|22:02:13.365] [batcher] Aggregating signatures...      caller=batcher.go:178
DEBUG[10-12|22:02:13.371] Exiting process batch                    duration=110ns caller=node.go:222
DEBUG[10-12|22:02:13.371] Exiting process batch                    duration=80ns  caller=node.go:222
DEBUG[10-12|22:02:13.373] Exiting process batch                    duration=100ns caller=node.go:222
DEBUG[10-12|22:02:13.373] Exiting process batch                    duration=160ns caller=node.go:222
TRACE[10-12|22:02:13.376] [batcher] AggregateSignatures took       duration=10.609723ms  caller=batcher.go:195
TRACE[10-12|22:02:13.376] [batcher] Confirming batch...            caller=batcher.go:198
```

To check the status of that same blob (replace `$REQUEST_ID` with the request ID from the prior step):

```
grpcurl -plaintext -d '{"request_id": "$REQUEST_ID"}' \
  localhost:32005 disperser.Disperser/GetBlobStatus
```

