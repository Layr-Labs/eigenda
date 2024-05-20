# EigenDA Plasma DA Server

## Introduction
This simple DA server implementation supports ephemeral storage via EigenDA.

## EigenDA Configuration
Additional cli args are provided for targeting an EigenDA network backend:
- `--eigenda-rpc`: RPC host of disperser service. (e.g, on holesky this is `disperser-holesky.eigenda.xyz:443`)
- `--eigenda-status-query-timeout`: (default: 25m) Duration for which a client will wait for a blob to finalize after being sent for dispersal.
- `--eigenda-status-query-retry-interval`: (default: 5s) How often a client will attempt a retry when awaiting network blob finalization. 
- `--eigenda-use-tls`: (default: true) Whether or not to use TLS for grpc communication with disperser.
- `eigenda-g1-path`: Directory path to g1.point file
- `eigenda-g2-power-of-tau`: Directory path to g2.point.powerOf2 file
- `eigenda-cache-path`: Directory path to dump cached SRS tables

## Running Locally
1. Compile binary: `make da-server`
2. Run binary; e.g: `./bin/da-server --addr 127.0.0.1 --port 5050 --eigenda-rpc 127.0.0.1:443 --eigenda-status-query-timeout 45m --eigenda-g1-path test/resources/g1.point --eigenda-g2-tau-path test/resources/g2.point.powerOf2 --eigenda-use-tls true`

### Commitment Schemas
An `EigenDACommitment` layer type has been added that supports verification against its respective pre-images. Otherwise this logic is pseudo-identical to the existing `Keccak256` commitment type. The commitment is encoded via the following byte array:
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

**NOTE:** Commitments are cryptographically verified against the data fetched from EigenDA for all `/get` calls.

## Testing
Some unit tests have been introduced to assert correctness of encoding/decoding logic and mocked server interactions. These can be ran via `make test`.

Otherwise E2E tests (`test/e2e_test.go`) exists which asserts that a commitment can be generated when inserting some arbitrary data to the server and can be read using the commitment for a key lookup via the client. These can be ran via `make e2e-test`. Please **note** that this test uses the EigenDA Holesky network which is subject to rate-limiting and slow confirmation times *(i.e, >10 minutes per blob confirmation)*. Please advise EigenDA's [inabox](https://github.com/Layr-Labs/eigenda/tree/master/inabox#readme) if you'd like to spin-up a local DA network for quicker iteration testing. 


## Downloading Mainnet SRS
KZG commitment verification requires constructing the SRS string from the proper trusted setup values (g1, g2, g2.power_of_tau). These values can be downloaded locally using the [operator-setup](https://github.com/Layr-Labs/eigenda-operator-setup) submodule via the following commands.

1. `make submodules`
2. `make srs`


## Resources
- [op-stack](https://github.com/ethereum-optimism/optimism)
- [plasma spec](https://specs.optimism.io/experimental/plasma.html)
- [eigen da](https://github.com/Layr-Labs/eigenda)


## Hardware Requirements
The following specs are recommended for running on a single production server:
* 12 GB SSD (assuming SRS values are stored on instance)
* 16 GB RAM
* 1-2 cores CPU
