# EigenDA Metadata Store Benchmark Suite

A comprehensive benchmark suite for comparing the performance of DynamoDB and PostgreSQL implementations of the EigenDA metadata store.

## Features

- **Multi-operation support**: Benchmark different operations simultaneously at different rates
- **Configurable load generation**: Control operations per second, duration, and concurrency
- **Comprehensive metrics**: Latency percentiles (p50, p90, p95, p99), throughput, error rates
- **Warmup phase**: Optional warmup period before collecting metrics
- **Real-time reporting**: Periodic progress updates during benchmark execution
- **Easy setup**: Docker Compose for local testing environments

## Quick Start

### 1. Setup Local Test Environment

```bash
# Start local DynamoDB and PostgreSQL
docker-compose up -d

# Wait for services to be ready
docker-compose ps

# View logs if needed
docker-compose logs -f
```

### 2. Build the Benchmark Tool

```bash
go build -o benchmark cmd/benchmark/main.go
```

### 3. Run Basic Benchmarks

```bash
# DynamoDB: UpdateBlobStatus at 200 ops/s for 30 seconds
./benchmark -store=dynamodb \
  -dynamo-endpoint=http://localhost:8000 \
  -ops="UpdateBlobStatus:200:30s"

# PostgreSQL: Same test
./benchmark -store=postgresql \
  -pg-host=localhost \
  -pg-password=benchmark123 \
  -ops="UpdateBlobStatus:200:30s"
```

## Operation Types

The benchmark suite supports the following operations:

| Operation | Description |
|-----------|-------------|
| `UpdateBlobStatus` | Updates the status of a blob |
| `PutBlobMetadata` | Stores new blob metadata |
| `GetBlobMetadata` | Retrieves blob metadata |
| `PutBlobCertificate` | Stores blob certificate |
| `GetBlobCertificate` | Retrieves blob certificate |
| `PutBatch` | Stores batch information |
| `GetBatch` | Retrieves batch information |
| `PutDispersalRequest` | Stores dispersal request |
| `GetDispersalRequest` | Retrieves dispersal request |
| `PutAttestation` | Stores attestation |
| `GetAttestation` | Retrieves attestation |
| `PutBlobInclusionInfo` | Stores blob inclusion info |
| `GetBlobInclusionInfo` | Retrieves blob inclusion info |

## Configuration

### Command Line Options

```
-store          Store type: "dynamodb" or "postgresql" (default: "dynamodb")
-ops            Operations to benchmark (format: op1:rate:duration,op2:rate:duration)
-warmup         Warmup duration (default: 10s)
-report         Reporting interval (default: 5s)
-workers        Number of workers per operation (default: 10)

DynamoDB Options:
-dynamo-table    Table name (default: "test-metadata")
-dynamo-region   AWS region (default: "us-east-1")
-dynamo-endpoint DynamoDB endpoint for local testing

PostgreSQL Options:
-pg-host         PostgreSQL host (default: "localhost")
-pg-port         PostgreSQL port (default: 5432)
-pg-user         PostgreSQL username (default: "postgres")
-pg-password     PostgreSQL password
-pg-database     PostgreSQL database (default: "eigenda_benchmark")
-pg-sslmode      PostgreSQL SSL mode (default: "disable")

Logging:
-log-level       Log level: debug, info, warn, error (default: "info")
```

## Advanced Examples

### Multi-Operation Benchmark

Test multiple operations simultaneously with different rates:

```bash
./benchmark -store=postgresql \
  -pg-password=benchmark123 \
  -ops="UpdateBlobStatus:200:60s,PutBlobMetadata:100:60s,GetBlobMetadata:300:60s" \
  -warmup=20s \
  -workers=20
```

### High-Load Testing

Test with higher operation rates and more workers:

```bash
./benchmark -store=dynamodb \
  -dynamo-endpoint=http://localhost:8000 \
  -ops="UpdateBlobStatus:1000:2m,PutBlobCertificate:500:2m" \
  -workers=50 \
  -warmup=30s
```

### Production-like Testing

Test against production-like configurations:

```bash
# DynamoDB (using real AWS)
./benchmark -store=dynamodb \
  -dynamo-table=eigenda-metadata-prod \
  -dynamo-region=us-east-1 \
  -ops="UpdateBlobStatus:500:5m,GetBlobMetadata:2000:5m" \
  -workers=100

# PostgreSQL (remote server)
./benchmark -store=postgresql \
  -pg-host=db.production.example.com \
  -pg-port=5432 \
  -pg-user=eigenda \
  -pg-password=$DB_PASSWORD \
  -pg-database=eigenda \
  -pg-sslmode=require \
  -ops="UpdateBlobStatus:500:5m,GetBlobMetadata:2000:5m" \
  -workers=100
```

## Understanding Results

The benchmark provides detailed metrics for each operation:

```
=== Final Benchmark Results ===
Store Type: dynamodb
Total Duration: 30.5s

--- UpdateBlobStatus ---
Target Rate: 200 ops/sec
Total Operations: 6000 (Success: 5985, Failed: 15)
Actual Throughput: 196.72 ops/sec
Error Rate: 0.25%
Latency (μs):
  Min: 234
  Avg: 4567
  P50: 3234
  P90: 8901
  P95: 11234
  P99: 18901
  Max: 45678
```

### Key Metrics Explained

- **Total Operations**: Number of operations attempted
- **Success/Failed**: Count of successful vs failed operations
- **Actual Throughput**: Achieved operations per second
- **Error Rate**: Percentage of failed operations
- **Latency Percentiles**: Response time distribution
  - P50: 50% of requests completed within this time
  - P90: 90% of requests completed within this time
  - P95: 95% of requests completed within this time
  - P99: 99% of requests completed within this time

## Makefile Usage

```bash
# Setup local environment
make setup

# Run DynamoDB benchmark
make bench-dynamodb

# Run PostgreSQL benchmark
make bench-postgres

# Run comparison benchmark
make bench-compare

# Cleanup
make clean
```

## Performance Tuning Tips

### DynamoDB
- Ensure sufficient read/write capacity units
- Consider on-demand billing for variable workloads
- Use consistent reads only when necessary
- Monitor throttling metrics

### PostgreSQL
- Tune connection pool settings
- Adjust `max_connections` in postgresql.conf
- Consider using connection pooler like PgBouncer
- Monitor slow queries and create appropriate indexes

### General
- Start with lower rates and gradually increase
- Use warmup period to stabilize connections
- Monitor system resources (CPU, memory, network)
- Run benchmarks from same region/network as database

## Troubleshooting

### DynamoDB Issues

```bash
# Check if local DynamoDB is running
curl http://localhost:8000

# View DynamoDB Admin UI
open http://localhost:8001

# Check table creation
aws dynamodb list-tables --endpoint-url http://localhost:8000
```

### PostgreSQL Issues

```bash
# Check PostgreSQL connection
psql -h localhost -U benchmark -d eigenda_benchmark -c "SELECT 1"

# View pgAdmin
open http://localhost:8080
# Login: admin@benchmark.local / admin123

# Check table creation
psql -h localhost -U benchmark -d eigenda_benchmark -c "\dt"
```

## Contributing

When adding new operations:

1. Add the operation type to `OperationType` enum
2. Implement the operation in `executeOperation` method
3. Add test data generation if needed
4. Update the README with the new operation

## License

[Your License Here]