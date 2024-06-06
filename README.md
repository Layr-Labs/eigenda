# EigenDA Sidecar Proxy

## Introduction

This service wraps the high-level EigenDA client, exposing endpoints for interacting with the EigenDA disperser in conformance to the [OP plasma server spec](https://specs.optimism.io/experimental/plasma.html), and adding disperser verification logic. This simplifies integrating EigenDA into various rollup frameworks by minimizing the footprint of changes needed within their respective services. Features of the EigenDA sidecar proxy include:

* blob submission/retrieval to EigenDA
* data <--> blob encoding/decoding for BN254 field element compatibility
* tamper resistance assurance (i.e, cryptographic verification of retrieved blobs) for avoiding a trust dependency on the EigenDA disperser

In order to disperse to the EigenDA network in production, or at high throughput on testnet, please register your authentication ethereum address through [this form](https://forms.gle/3QRNTYhSMacVFNcU8). Your EigenDA authentication keypair should not be associated with any funds anywhere.

## EigenDA Configuration

Additional CLI args are provided for targeting an EigenDA network backend:

* `--eigenda-rpc`: RPC host of disperser service. (e.g, on holesky this is `disperser-holesky.eigenda.xyz:443`, full network list [here](https://docs.eigenlayer.xyz/eigenda/networks/))
* `--eigenda-status-query-timeout`: (default: 30m) Duration for which a client will wait for a blob to finalize after being sent for dispersal.
* `--eigenda-status-query-retry-interval`: (default: 5s) How often a client will attempt a retry when awaiting network blob finalization.
* `--eigenda-disable-tls`: (default: false) Whether to disable TLS for grpc communication with disperser.
* `--eigenda-response-timeout`: (default: 10s) The total amount of time that the client will wait for a response from the EigenDA disperser.
* `--eigenda-custom-quorum-ids`: (default: []) The quorum IDs to write blobs to using this client. Should not include default quorums 0 or 1.
* `--eigenda-signer-private-key-hex`: Signer private key in hex encoded format. This key should not be associated with an Ethereum address holding any funds.
* `--eigenda-put-blob-encoding-version`: The blob encoding version to use when writing blobs from the high level interface.
* `--eigenda-disable-point-verification-mode`: Point verification mode does an IFFT on data before it is written, and does an FFT on data after it is read. This makes it possible to open points on the KZG commitment to prove that the field elements correspond to the commitment. With this mode disabled, you will need to supply the entire blob to perform a verification that any part of the data matches the KZG commitment.
* `--eigenda-g1-path`: Directory path to g1.point file
* `--eigenda-g2-power-of-tau`: Directory path to g2.point.powerOf2 file
* `--eigenda-cache-path`: Directory path to dump cached SRS tables
* `--eigenda-max-blob-length`: The maximum blob length that this EigenDA sidecar proxy should expect to be written or read from EigenDA. This configuration setting is used to determine how many SRS points should be loaded into memory for generating/verifying KZG commitments returned by the EigenDA disperser. Valid byte units are either base-2 or base-10 byte amounts (not bits), e.g. `30 MiB`, `4Kb`, `30MB`. The maximum blob size is a little more than `1GB`.

### In-Memory Storage

An ephemeral memory store backend can be used for faster feedback testing when performing rollup integrations. The following cli args can be used to target the feature:

* `--memstore.enabled`: Boolean feature flag
* `--memstore.expiration`: Duration for which a blob will exist

## Running Locally

1. Compile binary: `make eigenda-proxy`
2. Run binary; e.g: `./bin/eigenda-proxy --addr 127.0.0.1 --port 5050 --eigenda-rpc 127.0.0.1:443 --eigenda-status-query-timeout 45m --eigenda-g1-path test/resources/g1.point --eigenda-g2-tau-path test/resources/g2.point.powerOf2 --eigenda-use-tls true`

**Env File**
An env file can be provided to the binary for runtime process ingestion; e.g:

1. Create env: `cp .env.example .env`
2. Pass into binary: `ENV_PATH=.env ./bin/eigenda-proxy`

## Running via Docker

Container can be built via running `make build-docker`.

## Commitment Schemas

An `EigenDACommitment` layer type has been added that supports verification against its respective pre-images. The commitment is encoded via the following byte array:

```
            0        1        2        3        4                 N
            |--------|--------|--------|--------|-----------------|
             commit   da layer  ext da   version  raw commitment
             type       type    type      byte

```

The `raw commitment` for EigenDA is encoding certificate and kzg fields.

**NOTE:** Commitments are cryptographically verified against the data fetched from EigenDA for all `/get` calls. The server will respond with status `500` in the event where EigenDA were to lie and provide falsified data thats irrespective of the client provided commitment. This feature isn't flag guarded and is part of standard operation.

## Testing

Some unit tests have been introduced to assert the correctness of:

* DA Certificate encoding/decoding logic
* commitment verification logic

Unit tests can be ran via `make test`.

Otherwise E2E tests (`test/e2e_test.go`) exists which asserts that a commitment can be generated when inserting some arbitrary data to the server and can be read using the commitment for a key lookup via the client. These can be ran via `make e2e-test`. Please **note** that this test uses the EigenDA Holesky network which is subject to rate-limiting and slow confirmation times *(i.e, >10 minutes per blob confirmation)*. Please advise EigenDA's [inabox](https://github.com/Layr-Labs/eigenda/tree/master/inabox#readme) if you'd like to spin-up a local DA network for faster iteration testing.

## Downloading Mainnet SRS

KZG commitment verification requires constructing the SRS string from the proper trusted setup values (g1, g2, g2.power_of_tau). These values can be downloaded locally using the [operator-setup](https://github.com/Layr-Labs/eigenda-operator-setup) submodule via the following commands.

1. `make submodules`
2. `make srs`

## Hardware Requirements

The following specs are recommended for running on a single production server:

* 12 GB SSD (assuming SRS values are stored on instance)
* 16 GB RAM
* 1-2 cores CPU

## Resources

* [op-stack](https://github.com/ethereum-optimism/optimism)

* [plasma spec](https://specs.optimism.io/experimental/plasma.html)
* [eigen da](https://github.com/Layr-Labs/eigenda)
