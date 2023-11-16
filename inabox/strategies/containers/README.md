# Inabox Docker Environment

A docker compose based local EigenDA environment.

### Goals

* Support local environment deploys
  * in other words, not production environments
* Use docker compose
* Be simple
  * Data flows in one direction: `dockerCompose(configLock(config.yaml, secrets/, resources/))`
  * Initialization containers use config.lock and mounts in /contracts, and /subgraphs

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
