# Common

Shared utilities and code used across multiple packages in the EigenDA system.

## Subdirectories

| Subdirectory     | Description                                                                      |
|------------------|----------------------------------------------------------------------------------|
| ./aws            | AWS client config and utilities for DynamoDB, KMS, and secrets                   |
| ./cache          | Generic in-memory cache interfaces with weight-based capacity and FIFO eviction  |
| ./config         | Configuration parsing from files and environment variables with validation       |
| ./disperser      | Interface for querying disperser registry contract information                   |
| ./enforce        | Assertion functions that panic with descriptive messages on failure              |
| ./geth           | Ethereum client wrappers with multi-node failover and transaction signing        |
| ./healthcheck    | Heartbeat monitoring to detect stalled components                                |
| ./kvstore        | Key-value store interface backed by LevelDB with batch operations                |
| ./math           | Generic math utilities not in Go's standard library                              |
| ./memory         | Container memory limit detection and GC tuning to prevent OOM                    |
| ./metrics        | Prometheus metrics factory with automatic documentation                          |
| ./nameremapping  | YAML-based account address to human-readable name mapping                        |
| ./pprof          | HTTP server exposing Go runtime profiling endpoints                              |
| ./pubip          | Public IP address resolution with multiple provider fallback                     |
| ./ratelimit      | Leaky bucket rate limiter with KV store backend and metrics                      |
| ./replay         | Replay attack protection via request hash tracking with time windows             |
| ./reputation     | Entity reliability tracking using exponential moving average                     |
| ./s3             | S3 client interface supporting AWS and S3-compatible services                    |
| ./store          | Generic KV store implementations backed by DynamoDB or local files               |
| ./structures     | Data structures and algorithm utilities                                          |
| ./version        | Semantic versioning parsing and comparison                                       |
