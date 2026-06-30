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
- **Configurable**: Flexible configuration via files, environment variables, or CLI flags

## Building

```bash
# Build from the chainstate directory
cd chainstate
make build

# Or build from the root directory
make build
```

The binary will be created at `chainstate/bin/chainstate-indexer`.

## Configuration

The indexer uses the EigenDA documented config framework. Configuration can be provided via:
1. Configuration files (YAML, TOML, or JSON)
2. Environment variables (prefixed with `CHAINSTATE_INDEXER_`)
3. CLI flags

### Required Configuration

| Field | Description | Environment Variable |
|-------|-------------|---------------------|
| `RegistryCoordinatorAddr` | Address of the RegistryCoordinator contract | `CHAINSTATE_INDEXER_REGISTRY_COORDINATOR_ADDR` |
| `BLSApkRegistryAddr` | Address of the BLSApkRegistry contract | `CHAINSTATE_INDEXER_BLS_APK_REGISTRY_ADDR` |
| `EjectionManagerAddr` | Address of the EjectionManager contract | `CHAINSTATE_INDEXER_EJECTION_MANAGER_ADDR` |
| `PersistencePath` | Path to JSON file for state persistence | `CHAINSTATE_INDEXER_PERSISTENCE_PATH` |
| `HTTPPort` | Port for the HTTP API server | `CHAINSTATE_INDEXER_HTTP_PORT` |
| `EthRpcUrls` | Ethereum RPC endpoint URLs (must be archive nodes, see note) | `CHAINSTATE_INDEXER_ETH_RPC_URLS` |

> **Archive node required.** When recording quorum APK snapshots, the indexer
> reads each quorum's aggregate public key and total stake *as of the block of
> the membership-change event* (matching the operator-state subgraph it
> replaces). Serving these historical `eth_call`s requires an archive node. If
> pointed at a non-archive (full) node, the indexer will fail to read historical
> state and stall on the affected block range rather than record incorrect data.

### Optional Configuration

| Field | Default | Description |
|-------|---------|-------------|
| `StartBlockNumber` | 0 | Starting block for indexing (0 = current block) |
| `BlockBatchSize` | 1000 | Number of blocks to process per batch |
| `PollInterval` | 12s | Interval between polling for new blocks |
| `PersistInterval` | 30s | Interval for persisting state snapshots |
| `MetricsHTTPPort` | 9090 | Port for Prometheus metrics endpoint |
| `EnableMetrics` | true | Enable metrics collection |

### Example Configuration File

```yaml
# config.yaml
registry_coordinator_addr: "0x1234..."
bls_apk_registry_addr: "0x5678..."
ejection_manager_addr: "0x9abc..."
persistence_path: "/data/chainstate.json"
http_port: "8080"
start_block_number: 0
block_batch_size: 1000
poll_interval: 12s
persist_interval: 30s
enable_metrics: true
metrics_http_port: "9090"

# eth_client_config is nested
eth_client_config:
  rpc_url: "http://localhost:8545"

# logger_config is nested
logger_config:
  log_level: "info"
  format: "json"
```

### Example Using Environment Variables

```bash
export CHAINSTATE_INDEXER_REGISTRY_COORDINATOR_ADDR="0x1234..."
export CHAINSTATE_INDEXER_BLS_APK_REGISTRY_ADDR="0x5678..."
export CHAINSTATE_INDEXER_EJECTION_MANAGER_ADDR="0x9abc..."
export CHAINSTATE_INDEXER_PERSISTENCE_PATH="/data/chainstate.json"
export CHAINSTATE_INDEXER_HTTP_PORT="8080"
export CHAINSTATE_INDEXER_ETH_RPC_URLS="http://localhost:8545"

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
      "bls_pub_key_g1": {...},
      "bls_pub_key_g2": {...},
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
- `quorum_id` (uint8): Filter by quorum ID
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

Prometheus metrics are exposed at `/metrics` on the configured metrics port (default: 9090).

### Available Metrics

**Indexer Metrics:**
- `eigenda_chainstate_indexer_last_indexed_block`: Last block number indexed
- `eigenda_chainstate_indexer_blocks_indexed_total`: Total blocks indexed
- `eigenda_chainstate_indexer_index_errors_total`: Total indexing errors
- `eigenda_chainstate_indexer_operators_registered_total`: Total operator registrations
- `eigenda_chainstate_indexer_operators_deregistered_total`: Total operator deregistrations
- `eigenda_chainstate_indexer_socket_updates_total`: Total socket updates
- `eigenda_chainstate_indexer_ejections_recorded_total`: Total ejections recorded

**API Metrics:**
- `eigenda_chainstate_indexer_api_requests_total`: Total API requests by method, endpoint, and status
- `eigenda_chainstate_indexer_api_latency_seconds`: API request latency histogram
- `eigenda_chainstate_indexer_api_errors_total`: Total API errors by endpoint and type

**Standard Metrics:**
- Process metrics (CPU, memory, file descriptors, etc.)
- Go runtime metrics (goroutines, GC stats, etc.)

## Architecture

### Components

1. **Indexer**: Core component that polls Ethereum for new blocks and indexes events
2. **Store**: In-memory storage with JSON persistence (interface-based for future database support)
3. **API Server**: HTTP server providing REST endpoints for querying indexed data
4. **Metrics**: Prometheus metrics for observability

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
> not this indexer, except for the standalone `TestChainStateIndexerE2E` which
> runs the indexer as an observer against the same chain. Migrating consumers off
> `thegraph.IndexedChainState` onto the REST API is separate, future work.

If you're migrating from the operator state subgraph:

1. Deploy the chainstate indexer with the same contract addresses
2. Let it index from the desired start block (or current block)
3. Update clients to use the new REST API instead of GraphQL
4. The API response format differs from GraphQL, so client code will need updates

### Key Differences from Subgraph

| Aspect | Subgraph | ChainState Indexer |
|--------|----------|-------------------|
| Query Language | GraphQL | REST |
| Storage | PostgreSQL | In-memory + JSON |
| Deployment | Requires graph-node infrastructure | Single binary |
| Real-time updates | WebSocket subscriptions | Polling (WebSocket planned) |
| Query flexibility | High (GraphQL) | Medium (REST endpoints) |
| Performance | Depends on graph-node | Very fast (in-memory) |

## Development

### Running Tests

```bash
make test
```

### Running End-to-End Tests

An end-to-end test (`TestChainStateIndexerE2E`) runs the indexer against a live
[inabox](../inabox/README.md) devnet (a local chain with the EigenDA contracts
deployed and operators registered on-chain) and asserts on the results through
the REST API. It lives in `inabox/tests` and runs as part of the inabox suite:

```bash
cd inabox && make run-e2e-tests
```

To run only this test (the suite still brings up the full devnet, which requires
Docker):

```bash
cd inabox && go test ./tests -v -config=../templates/testconfig-anvil.yaml -run TestChainStateIndexerE2E
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
- **Reorg handling**: Detect and handle chain reorganizations
- **WebSocket support**: Real-time updates via WebSocket subscriptions
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
3. Monitor memory usage via Prometheus metrics

## Support

For issues, questions, or feature requests, please open an issue on GitHub.
