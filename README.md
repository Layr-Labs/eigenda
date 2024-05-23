# EigenDA Sidecar Proxy

## Introduction
This simple DA server implementation is a side-car communication relay between different rollup frameworks and EigenDA. This allows us to keep existing protocol functions (i.e, batch submission, state derivation) lightweight in respect to modification since the server handles key security and data operations like:
* blob submission/retrieval to EigenDA
* data <--> blob encoding/decoding
* tamper resistance assurance (i.e, cryptographic verification of retrieved blobs)

This allows for deduplication of redundant logical flows into a single representation which can be used cross functionally across rollups.

## EigenDA Configuration
Additional cli args are provided for targeting an EigenDA network backend:
- `--eigenda-rpc`: RPC host of disperser service. (e.g, on holesky this is `disperser-holesky.eigenda.xyz:443`)
- `--eigenda-status-query-timeout`: (default: 25m) Duration for which a client will wait for a blob to finalize after being sent for dispersal.
- `--eigenda-status-query-retry-interval`: (default: 5s) How often a client will attempt a retry when awaiting network blob finalization. 
- `--eigenda-use-tls`: (default: true) Whether or not to use TLS for grpc communication with disperser.
- `--eigenda-g1-path`: Directory path to g1.point file
- `--eigenda-g2-power-of-tau`: Directory path to g2.point.powerOf2 file
- `--eigenda-cache-path`: Directory path to dump cached SRS tables

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

The `raw commitment` for EigenDA is encoding the following certificate and kzg fields:
```go
type Cert struct {
	BatchHeaderHash      []byte
	BlobIndex            uint32
	ReferenceBlockNumber uint32
	QuorumIDs            []uint32
	BlobCommitment *common.G1Commitment
}
```

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
- [op-stack](https://github.com/ethereum-optimism/optimism)
- [plasma spec](https://specs.optimism.io/experimental/plasma.html)
- [eigen da](https://github.com/Layr-Labs/eigenda)
