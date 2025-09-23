# Inabox Tests

Notice: The scripts for setting up a local geth chain are currently broken. The instructions below use anvil instead

## First time setup
- Go path is in system path. [Instructions for installing go](https://go.dev/doc/install).
- Ensure all submodules are initialized and checked out
    ```
    $ git submodule update --init --recursive
    ```
- Docker is installed. [Instructions for installing docker](https://www.docker.com/products/docker-desktop/).
- Ensure foundry is installed (comes with `anvil` which we use as a test chain and `forge` which we use for deployment scripting):
    ```
    $ curl -L https://foundry.paradigm.xyz | bash
    $ foundryup
    ```
- `brew` is installed, see instructions [here](https://brew.sh/).
- Localstack CLI is installed (simulates AWS stack on local machine; we also provide instructions for running localstack from docker without the CLI):
    ```
    $ brew install localstack/tap/localstack-cli
    ```
- `grpcurl` is installed:
    ```
    $ brew install grpcurl
    ```
- `aws` is installed, follow instructions [here](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html).
- `npm` is installed
   ```
   $ brew install node
   ```
- `yarn` is installed
   ```
   $ npm install --global yarn
   ```
- Install contracts

  Follow the installation instructions under [`contracts/README.md`](../contracts/README.md)

- The Graph is installed
   ```
   $ npm install -g @graphprotocol/graph-cli@latest
   ```

## Run a complete end-to-end test

You can run a complete end-to-end test by running the following command:
```
cd inabox
make run-e2e
```

## Manually deploy the experiment and interact with the services

The following instructions will guide you through the process of manually deploying the test infrastructure, deploying the contracts, configuring the services, and interacting with the services. Currently this is done through the `inabox/deploy/cmd` tool. It will soon be replaced with a more user-friendly CLI.

### Preliminary setup steps

Ensure that all submodules (e.g. EigenLayer smart contracts) are checked out to the correct branch, and then build the binaries.
```
$ git submodule update --init --recursive
$ make build
```

#### Create a new configuration file:
```
cd inabox
make new-anvil
```

This will create a new file, e.g. `./testdata/12D-07M-2023Y-14H-41M-19S/config.yaml`. Please feel free to inspect the file and make any desired configuration changes at this point. After you have deployed the experiment, changes will not go into effect. 

> [!IMPORTANT]
If you shut down the test infrastructure (e.g. anvil, localstack, graph node), you will need to create a new configuration file by running `make new-anvil` again. This is because the state of the configuration file is tied to the state of the anvil chain, so restarting the chain will require a new configuration file.


### Provision the test infrastructure, deploy contracts, and configure services

The deployment process can be run using the new deploy command with various subcommands. The deploy command is located at `./deploy/cmd` and provides more granular control over the deployment process.

#### Available Deploy Commands

| Command | Description |
|---------|-------------|
| `chain` | Deploy the chain infrastructure (anvil, graph) for the inabox test |
| `localstack` | Deploy localstack and create the AWS resources needed for the inabox test |
| `exp` | Deploy the contracts and create configurations for all EigenDA components |
| `env` | Generate the environment variables for the inabox test |
| `eigenda` | Deploy EigenDA infrastructure with churner via testbed and other components via StartBinaries |
| `all` | Deploy all infrastructure, resources, contracts |

TODO(dmanc): Is there ever a reason we would want to run these commands separately? If not, we can simplify the instructions to just running `all`.

#### Command-line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--testname` | Name of the test to run (in `inabox/testdata`) | Auto-generated or latest |
| `--root-path` | Path to the root of repo | `../` |
| `--localstack-port` | Port for localstack services | `""` |
| `--deploy-resources` | Whether to deploy AWS resources on localstack | `""` |

#### Option 1 (simplest): Deploy everything in one command

Run the following command (from inabox directory)
```
make deploy-all
```

This will:
- Start all test infrastructure (localstack, graph node, anvil chain)
- Create the necessary AWS resources on localstack
- Deploy the smart contracts to anvil
- Deploy subgraphs to the graph node (if configured)
- Create configurations for the eigenda services (located in `inabox/testdata/DATETIME/envs`)

If there are any deployment errors, look at `inabox/testdata/DATETIME/deploy.log` for a detailed log. 

To view the configurations created for the EigenDA service components, look in `inabox/testdata/DATETIME/envs`

### Run the binaries and send traffic

Run the binaries:
```
cd inabox
./bin.sh start-detached
```
This will start all the EigenDA services in detached mode. The logs will be saved to `inabox/testdata/DATETIME/logs`.

Alternatively, you can start and stop the EigenDA services in the foreground by running:
```cd inabox
./bin.sh start
```
This will start all the EigenDA services in the foreground. The logs will be displayed in the terminal. It is highly recommended to run the services in detached mode using `./bin.sh start-detached` as described above.

The `eigenda` command automatically starts all EigenDA services including:
- Churner running as a testbed container on port 32002
- All other EigenDA components via StartBinaries

The churner logs can be viewed using docker logs on the churner container. Other component logs will be displayed or saved according to the StartBinaries configuration.

### Interact with EigenDA V1

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
  localhost:32003 disperser.Disperser/GetBlobStatus
```

## TODO: Interact with the EigenDA V2

### Cleanup

If you followed [Option 1](#option-1-simplest) above, you can run the following command in order to clean up the test infra:
```
cd inabox
make stop-infra
```

If you followed [Option 2](#option-2), you can stop the infra services by `Ctrl-C`'ing in each terminal. For the graph, it's also important to run `docker compose down -v` from within the `inabox/thegraph` directory to make sure that the containers are fully removed. 


