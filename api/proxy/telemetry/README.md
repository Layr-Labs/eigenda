# OpenTelemetry Tracing for EigenDA Proxy

This package provides OpenTelemetry tracing integration for the EigenDA proxy, allowing distributed tracing of requests through the v2 clients (DisperserClient, PayloadDisperser, RelayPayloadRetriever, and ValidatorPayloadRetriever).

## Configuration

The proxy can be configured to export traces using the following flags/environment variables:

| Flag | Environment Variable | Default | Description |
|------|---------------------|---------|-------------|
| `--otel.enabled` | `EIGENDA_PROXY_OTEL_ENABLED` | `false` | Enable OpenTelemetry tracing |
| `--otel.service-name` | `EIGENDA_PROXY_OTEL_SERVICE_NAME` | `eigenda-proxy` | Service name for tracing |
| `--otel.exporter.otlp.endpoint` | `EIGENDA_PROXY_OTEL_EXPORTER_OTLP_ENDPOINT` | `localhost:4318` | OTLP exporter endpoint (HTTP) |
| `--otel.exporter.otlp.insecure` | `EIGENDA_PROXY_OTEL_EXPORTER_OTLP_INSECURE` | `true` | Use insecure connection for OTLP exporter |
| `--otel.trace.sample-rate` | `EIGENDA_PROXY_OTEL_TRACE_SAMPLE_RATE` | `1.0` | Trace sampling rate (0.0 to 1.0) |

## Quick Start with Docker Compose

The easiest way to get started is using the included docker-compose setup:

### 1. Start all services (Proxy + Tempo + Grafana)

```bash
cd api/proxy
docker-compose up -d
```

This will start:
- **EigenDA Proxy** with tracing enabled (port 4242)
- **Grafana Tempo** for trace collection (ports 3200, 4317, 4318)
- **Grafana** for visualization (port 3000)
- **Prometheus** for metrics (port 9090)
- **MinIO** for S3-compatible storage

### 2. View traces in Grafana

1. Open Grafana at http://localhost:3000 (credentials: admin/admin)
2. The Tempo datasource is pre-configured automatically
3. Go to **Explore** → Select **Tempo** datasource
4. Search for traces by:
   - Service name: `eigenda-proxy`
   - Operation name (e.g., `DisperserClient.DisperseBlob`)
   - Trace ID

### 3. View the distributed trace timeline

Each trace will show the complete request flow:
- HTTP request received by proxy
- Dispersal operations
- Status polling
- Certificate building and verification
- Retrieval operations

## Manual Setup (without Docker Compose)

If you prefer to run services separately:

### 1. Start Grafana Tempo

```bash
docker run -d --name tempo \
  -p 4318:4318 \
  -p 3200:3200 \
  grafana/tempo:latest \
  -config.file=/etc/tempo.yaml
```

### 2. Start Grafana

```bash
docker run -d --name grafana \
  -p 3000:3000 \
  grafana/grafana:latest
```

Then configure Tempo as a data source in Grafana:
1. Go to Configuration > Data Sources
2. Add Tempo data source
3. Set URL to `http://tempo:3200` (or `http://localhost:3200` if running on host)
4. Save & Test

### 3. Start the EigenDA Proxy with tracing enabled

```bash
# Using flags
./eigenda-proxy \
  --otel.enabled=true \
  --otel.service-name=eigenda-proxy \
  --otel.exporter.otlp.endpoint=localhost:4318 \
  --otel.exporter.otlp.insecure=true \
  --otel.trace.sample-rate=1.0 \
  # ... other proxy flags ...

# Or using environment variables
export EIGENDA_PROXY_OTEL_ENABLED=true
export EIGENDA_PROXY_OTEL_SERVICE_NAME=eigenda-proxy
export EIGENDA_PROXY_OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318
export EIGENDA_PROXY_OTEL_EXPORTER_OTLP_INSECURE=true
export EIGENDA_PROXY_OTEL_TRACE_SAMPLE_RATE=1.0
./eigenda-proxy # ... other proxy flags ...
```

### 5. View traces in Grafana

1. Go to Explore in Grafana
2. Select the Tempo data source
3. Use the search to find traces by service name (`eigenda-proxy`)
4. View distributed traces showing the full request flow through:
   - DisperserClient operations (DisperseBlob, GetBlobStatus, etc.)
   - PayloadDisperser operations (SendPayload, polling)
   - Retrieval operations (GetPayload, GetEncodedPayload)

## Trace Spans

The following operations are automatically instrumented with traces:

### DisperserClient
- `DisperserClient.DisperseBlob` - Blob dispersal with blob size, quorum count, and blob key
- `DisperserClient.GetBlobStatus` - Status retrieval with blob key and status
- `DisperserClient.GetPaymentState` - Payment state retrieval
- `DisperserClient.GetBlobCommitment` - Commitment calculation

### PayloadDisperser
- `PayloadDisperser.SendPayload` - Full payload dispersal flow
- `PayloadDisperser.pollBlobStatusUntilSigned` - Polling for blob signatures

### Retrieval Clients
- `RelayPayloadRetriever.GetPayload` / `GetEncodedPayload` - Retrieval from relays
- `ValidatorPayloadRetriever.GetPayload` / `GetEncodedPayload` - Retrieval from validators

Each span includes relevant attributes like blob keys, payload sizes, quorum information, and error details.

## Production Considerations

### Sampling

For high-throughput production environments, consider reducing the sample rate to avoid overwhelming your tracing backend:

```bash
--otel.trace.sample-rate=0.1  # Sample 10% of traces
```

### Security

For production deployments, use secure connections:

```bash
--otel.exporter.otlp.insecure=false
```

And ensure your OTLP endpoint supports TLS.

### Alternative Backends

The OTLP exporter works with any OpenTelemetry-compatible backend, including:
- Grafana Tempo
- Jaeger (v1.35+)
- Zipkin (via OTLP)
- Cloud providers (AWS X-Ray, Google Cloud Trace, Azure Monitor)
- Commercial solutions (Datadog, New Relic, Honeycomb, etc.)
