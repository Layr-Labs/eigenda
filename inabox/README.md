# Inabox Tests

## First time setup

- Ensure [docker](https://www.docker.com/get-started/) is installed and running

- Ensure all submodules are initialized and checked out

    ```
    git submodule update --init --recursive
    ```

- Ensure foundry is installed (comes with `anvil` which we use as a test chain and `forge` which we use for deployment scripting):

    ```
    curl -L https://foundry.paradigm.xyz | bash
    foundryup -C $(cat .foundryrc)
    ```

- grpcurl is installed:

    ```
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
    ```

## Run a complete end-to-end test

You can run a complete end-to-end test by running the following command:

```
cd inabox
make run-e2e
```

# Inabox Docker Environment

A docker compose based local EigenDA environment.

### Goals

- Support local environment deploys
- Use docker compose
- Be simple
  - Data flows in one direction: `dockerCompose(configLock(config.yaml, secrets/, resources/))`
  - Initialization containers use config.lock and mounts in /contracts, and /subgraphs

### Usage

1. `make new-anvil`
2. Modify new configuration as needed
3. `make config`
4. `make devnet-up`
5. `make devnet-logs`
6. Clean up with `make devnet-down clean`

### How it works

`make config` does a few things:

1. Locally runs the gen-env/cmd script to combine a configuration's config.yaml and the global inabox resources/ and secrets/
into a config.lock and a docker-compose.gen.yaml.
2. Combines the outputted docker-compose.gen.yaml with the "common" docker-compose.yaml to generate the configuration's docker-compose.yaml

`make devnet-up` does a simple `docker-compose up --build`. Under the hood this is doing a few things:

1. After the "common" services come up, the `eth-init-script` container starts and deploys AWS resources and the EigenDA smart contracts
2. After that the `graph-init-script` container starts and deploys subgraphs.
3. Finally, the rest of the EigenDA services start.

Steps 1 and 2 are separate because the `forge script --broadcast` command within step 1 will only work on a (possibly emulated) amd64 image and
the graph deploy will only work on an image in the native-architecture, so for arm64 support these must be separate images.

### Troubleshooting

- If your docker builds hang, try enabling containerd pulls in your docker settings.

### How to Demo

Disperse a blob:

```
# This command uses `grpcurl`, a tool to send gRPC request in cli
# To install `grpcurl`, run `brew install grpcurl` or `go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest`

$ grpcurl -plaintext -d '{"data": "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000011111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", "security_params": [{"quorum_id": 0, "adversary_threshold": 50, "quorum_threshold": 100}]}' localhost:32001 disperser.Disperser/DisperseBlob
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

Then you can get the confirmation status of the blob like this:

```
grpcurl -plaintext -d '{"request_id":"<request id from disperse blob response>" }' localhost:32001 disperser.Disperser/GetBlobStatus
```

### Notes

- Notice: The scripts for setting up a local geth chain are currently broken. The instructions below use anvil instead
