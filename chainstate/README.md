# ChainState Indexer

The ChainState Indexer is a standalone service that indexes operator state events from Ethereum smart contracts and provides a query API for accessing the indexed data. It is designed to replace the operator state subgraph with a more performant and maintainable solution.

## Overview

The indexer monitors the following smart contracts:
- **RegistryCoordinator**: Operator registration, deregistration, and socket updates
- **BLSApkRegistry**: BLS public key registrations and quorum membership changes
- **EjectionManager**: Operator ejection events

It stores all indexed data in memory and periodically persists snapshots to disk as JSON files. The storage layer is abstracted through an interface, allowing for easy migration to a database backend in the future.

## Features

- **Real-time indexing**: Continuously polls for new blocks and indexes events
- **Persistent storage**: Automatic periodic snapshots to JSON files
- **REST API**: Query indexed data via HTTP endpoints
- **Graceful shutdown**: Ensures final state snapshot on shutdown
- **Configurable**: Flexible configuration via files and environment variables

## Building

```bash
# Build from the chainstate directory
cd chainstate
make build
```

The binary will be created at `chainstate/bin/chainstate-indexer`.

Note: the repository root's `make build` does not currently include this
service.

## Configuration

The indexer uses the EigenDA documented config framework. Configuration can be provided via:
1. Configuration files (YAML, TOML, or JSON; pass one or more with `--config`)
2. Environment variables (prefixed with `CHAINSTATE_INDEXER_`)

### Required Configuration

| Field | Description | Environment Variable |
|-------|-------------|---------------------|
| `Config.EigenDADirectory` | Address of the EigenDADirectory contract, from which all other contract addresses (RegistryCoordinator, BLSApkRegistry, EjectionManager, StakeRegistry) are resolved | `CHAINSTATE_INDEXER_CONFIG_EIGEN_DA_DIRECTORY` |
| `Config.PersistencePath` | Path to JSON file for state persistence | `CHAINSTATE_INDEXER_CONFIG_PERSISTENCE_PATH` |
| `Config.HTTPPort` | Port for the HTTP API server | `CHAINSTATE_INDEXER_CONFIG_HTTP_PORT` |
| `Secret.EthRpcUrls` | Ethereum RPC endpoint URLs (must be archive nodes, see note); extra URLs act as failover fallbacks | `CHAINSTATE_INDEXER_SECRET_ETH_RPC_URLS` |

> **Archive node required.** When recording quorum APK snapshots, the indexer
> reads each quorum's aggregate public key and total stake *as of the block of
> the membership-change event* (matching the operator-state subgraph it
> replaces). Serving these historical `eth_call`s requires an archive node. If
> pointed at a non-archive (full) node, the indexer will fail to read historical
> state and stall on the affected block range rather than record incorrect data.

### Optional Configuration

| Field | Default | Description |
|-------|---------|-------------|
| `Config.StartBlockNumber` | 0 | First block to index, inclusive. **0 means the current chain head**, which skips all historical events — including the one-time BLS pubkey registrations. Set it to the contract deployment block to index full history. |
| `Config.BlockBatchSize` | 1000 | Number of blocks to process per batch |
| `Config.PollInterval` | 12s | Interval between polling for new blocks |
| `Config.PersistInterval` | 30s | Interval for persisting state snapshots |
| `Config.EthClientConfig` | see `geth.DefaultEthClientConfig` | Ethereum client settings (retries, confirmations); the RPC URLs come from `Secret.EthRpcUrls` |
| `Config.LoggerConfig` | see `common.DefaultLoggerConfig` | Logging configuration |

### Example Configuration File

```yaml
# config.yaml
Config:
  EigenDADirectory: "0x1234..."
  PersistencePath: "/data/chainstate.json"
  HTTPPort: "8080"
  StartBlockNumber: 0
  BlockBatchSize: 1000
  PollInterval: 12s
  PersistInterval: 30s
  LoggerConfig:
    Format: json

Secret:
  EthRpcUrls:
    - "http://localhost:8545"
```

### Example Using Environment Variables

```bash
export CHAINSTATE_INDEXER_CONFIG_EIGEN_DA_DIRECTORY="0x1234..."
export CHAINSTATE_INDEXER_CONFIG_PERSISTENCE_PATH="/data/chainstate.json"
export CHAINSTATE_INDEXER_CONFIG_HTTP_PORT="8080"
export CHAINSTATE_INDEXER_SECRET_ETH_RPC_URLS="http://localhost:8545"

./bin/chainstate-indexer
```

## Running

### Using a Configuration File

```bash
./bin/chainstate-indexer --config config.yaml
```

### Using Environment Variables

```bash
# Set environment variables as shown above
./bin/chainstate-indexer
```

### Verify Configuration Without Running

```bash
./bin/chainstate-indexer --config config.yaml --only-verify-config
```

## API Endpoints

The indexer provides a REST API for querying indexed data.

### Health Check

```
GET /api/v1/health
```

Returns the health status of the API server.

### Status

```
GET /api/v1/status
```

Returns the current indexer status including the last indexed block.

**Response:**
```json
{
  "last_indexed_block": 12345678,
  "time": "2025-01-28T12:00:00Z"
}
```

### List Operators

```
GET /api/v1/operators?registered=true&limit=100&offset=0
```

Returns a paginated list of operators.

**Query Parameters:**
- `registered` (bool): Filter to only registered operators
- `deregistered` (bool): Filter to only deregistered operators
- `quorum_id` (uint8): Filter by quorum ID
- `min_block` (uint64): Filter by minimum registration block
- `max_block` (uint64): Filter by maximum registration block
- `limit` (int): Maximum number of results (default: 100)
- `offset` (int): Pagination offset (default: 0)

**Response:**
```json
{
  "operators": [
    {
      "id": "0x1234...",
      "address": "0x5678...",
      "bls_pubkey_g1": {...},
      "bls_pubkey_g2": {...},
      "socket": "example.com:32004",
      "registered_at_block_number": 12345000,
      "deregistered_at_block_number": null,
      "quorum_ids": [0, 1],
      "registered_tx_hash": "0xabcd...",
      "deregistered_tx_hash": null
    }
  ],
  "count": 1,
  "limit": 100,
  "offset": 0
}
```

### Get Operator

```
GET /api/v1/operators/:id
```

Returns a single operator by ID (32-byte hex string with or without "0x" prefix).

**Response:** Same as a single operator object from the list endpoint.

### Get Quorum APK

```
GET /api/v1/quorum-apk?quorum_id=0&block_number=12345678
```

Returns the aggregate public key for a quorum at a specific block.

**Query Parameters:**
- `quorum_id` (uint8, required): The quorum identifier
- `block_number` (uint64, required): The block number

**Response:**
```json
{
  "quorum_id": 0,
  "block_number": 12345678,
  "apk": {...},
  "total_stake": "1000000000000000000",
  "updated_at": "2025-01-28T12:00:00Z"
}
```

### List Quorum APK History

```
GET /api/v1/quorum-apk/history?quorum_id=0&min_block=12340000&max_block=12350000
```

Returns a list of quorum APK snapshots.

**Query Parameters:**
- `quorum_id` (uint8, required): The quorum whose snapshots to list
- `block_number` (uint64): Get APK for specific block
- `min_block` (uint64): Minimum block number
- `max_block` (uint64): Maximum block number

### List Ejections

```
GET /api/v1/ejections?limit=100&offset=0
```

Returns a paginated list of all ejection events.

**Query Parameters:**
- `limit` (int): Maximum number of results (default: 100)
- `offset` (int): Pagination offset (default: 0)

### List Operator Ejections

```
GET /api/v1/ejections/:operator_id?limit=100&offset=0
```

Returns ejection events for a specific operator.

### List Socket Updates

```
GET /api/v1/socket-updates/:operator_id?limit=100&offset=0
```

Returns socket update events for a specific operator.

## Metrics

Prometheus metrics are not yet implemented (see Future Enhancements). The
`/api/v1/status` endpoint exposes the last indexed block for basic progress
monitoring in the meantime.

## Architecture

### Components

1. **Indexer**: Core component that polls Ethereum for new blocks and indexes events
2. **Store**: In-memory storage with JSON persistence (interface-based for future database support)
3. **API Server**: HTTP server providing REST endpoints for querying indexed data

### Data Flow

```
Ethereum RPC → Indexer → MemoryStore → JSONPersister → Disk
                              ↓
                         API Server → HTTP Client
```

### Storage

The indexer uses an in-memory store with periodic JSON snapshots for persistence. This design provides:
- Fast query performance
- Simple deployment (no database required)
- Easy backup and recovery (just copy the JSON file)
- Future migration path (implement database Store interface)

## Migration from Subgraph

> **Status: not yet wired into any consumer.** The chainstate indexer is a
> candidate replacement for the operator-state subgraph, but no EigenDA service
> reads from it yet. Components that need indexed operator state (churner,
> disperser/controller, etc.) still query the subgraph via the `core/thegraph`
> client. The inabox test suite reflects this: running it exercises the subgraph,
> not this indexer, except for the standalone `TestChainStateIndexerE2E` and
> `TestChainStateSubgraphParity` tests, which run the indexer as an observer
> against the same chain. Migrating consumers off `thegraph.IndexedChainState`
> is separate, future work.

If you're migrating from the operator state subgraph:

1. Deploy the chainstate indexer pointed at the network's EigenDADirectory
   contract (the individual contract addresses are resolved from it)
2. Let it index from the desired start block (or current block)
3. Update clients to use the new REST API instead of GraphQL
4. The API response format differs from GraphQL, so client code will need updates

See [TESTING_AND_MIGRATION.md](./TESTING_AND_MIGRATION.md) for the full
validation and phased-migration plan.

### Key Differences from Subgraph

| Aspect | Subgraph | ChainState Indexer |
|--------|----------|-------------------|
| Query Language | GraphQL | REST |
| Storage | PostgreSQL | In-memory + JSON |
| Deployment | Requires graph-node infrastructure | Single binary |
| Real-time updates | WebSocket subscriptions | Polling |
| Query flexibility | High (GraphQL) | Medium (REST endpoints) |
| Performance | Depends on graph-node | Very fast (in-memory) |

## Development

### Running Tests

```bash
make test
```

### Running End-to-End Tests

Two end-to-end tests run the indexer against a live
[inabox](../inabox/README.md) devnet (a local chain with the EigenDA contracts
deployed and operators registered on-chain): `TestChainStateIndexerE2E`
asserts on the results through the REST API, and
`TestChainStateSubgraphParity` compares the `core.IndexedChainState`
implementation against the operator-state subgraph at every
membership-change block. Both live in `inabox/tests` and run as part of the
inabox suite:

```bash
cd inabox && make run-e2e-tests
```

To run only these tests (the suite still brings up the full devnet, which
requires Docker):

```bash
cd inabox && go test ./tests -v -config=../templates/testconfig-anvil.yaml -run 'TestChainStateIndexerE2E|TestChainStateSubgraphParity'
```

### Building

```bash
make build
```

### Cleaning

```bash
make clean
```

## Future Enhancements

- **Database backend**: PostgreSQL store implementation for larger datasets
- **GraphQL API**: Add GraphQL endpoint for backward compatibility
- **Metrics**: Prometheus metrics for indexing progress and API usage
- **Reorg handling**: Detect and handle chain reorganizations
- **Multi-chain support**: Index from multiple chains simultaneously
- **Historical import**: Tool to backfill from existing subgraph data

## Troubleshooting

### Indexer is not progressing

1. Check Ethereum RPC connectivity
2. Verify contract addresses are correct
3. Check logs for errors
4. Ensure sufficient disk space for state snapshots

### API returns 404 for operators

1. Verify the indexer has indexed the blocks containing the operator registration
2. Check the `last_indexed_block` via the `/api/v1/status` endpoint
3. Ensure the operator ID is correctly formatted (32-byte hex)

### High memory usage

The in-memory store will grow over time as more data is indexed. For large datasets:
1. Consider implementing a database backend
2. Add TTL or archival logic for old data
3. Monitor process memory usage externally (Prometheus metrics are not yet implemented)

## Support

For issues, questions, or feature requests, please open an issue on GitHub.
