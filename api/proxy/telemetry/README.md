# OpenTelemetry Tracing for EigenDA Proxy

This package provides OpenTelemetry tracing integration for the EigenDA proxy, allowing distributed tracing of requests through the v2 clients (currently only DisperserClient and PayloadDisperser).

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
