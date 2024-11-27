# Metrics Documentation for namespace 'relay'

This documentation was automatically generated at time `2024-11-27T10:08:20-06:00`

There are a total of `8` registered metrics.

---

## average_get_chunks_data_bytes

Average data size in a GetChunks request

|   |   |
|---|---|
| **Name** | `average_get_chunks_data` |
| **Unit** | `bytes` |
| **Type** | `running average` |
| **Time Window** | `1m0s` |
| **Fully Qualified Name** | `relay_average_get_chunks_data_bytes` |
---

## average_get_chunks_key_count

Average number of keys in a GetChunks request

|   |   |
|---|---|
| **Name** | `average_get_chunks_key` |
| **Unit** | `count` |
| **Type** | `running average` |
| **Time Window** | `1m0s` |
| **Fully Qualified Name** | `relay_average_get_chunks_key_count` |
---

## get_chunks_auth_failure_count

Number of GetChunks RPC authentication failures

|   |   |
|---|---|
| **Name** | `get_chunks_auth_failure` |
| **Unit** | `count` |
| **Type** | `counter` |
| **Fully Qualified Name** | `relay_get_chunks_auth_failure_count` |
---

## get_chunks_authentication_latency_ms

Latency of the GetChunks RPC client authentication

|   |   |
|---|---|
| **Name** | `get_chunks_authentication_latency` |
| **Unit** | `ms` |
| **Type** | `latency` |
| **Quantiles** | `0.500`, `0.900`, `0.990` |
| **Fully Qualified Name** | `relay_get_chunks_authentication_latency_ms` |
---

## get_chunks_data_latency_ms

Latency of the GetChunks RPC data retrieval

|   |   |
|---|---|
| **Name** | `get_chunks_data_latency` |
| **Unit** | `ms` |
| **Type** | `latency` |
| **Quantiles** | `0.500`, `0.900`, `0.990` |
| **Fully Qualified Name** | `relay_get_chunks_data_latency_ms` |
---

## get_chunks_latency_ms

Latency of the GetChunks RPC

|   |   |
|---|---|
| **Name** | `get_chunks_latency` |
| **Unit** | `ms` |
| **Type** | `latency` |
| **Quantiles** | `0.500`, `0.900`, `0.990` |
| **Fully Qualified Name** | `relay_get_chunks_latency_ms` |
---

## get_chunks_metadata_latency_ms

Latency of the GetChunks RPC metadata retrieval

|   |   |
|---|---|
| **Name** | `get_chunks_metadata_latency` |
| **Unit** | `ms` |
| **Type** | `latency` |
| **Quantiles** | `0.500`, `0.900`, `0.990` |
| **Fully Qualified Name** | `relay_get_chunks_metadata_latency_ms` |
---

## get_chunks_rate_limited_count

Number of GetChunks RPC rate limited

|   |   |
|---|---|
| **Name** | `get_chunks_rate_limited` |
| **Unit** | `count` |
| **Labels** | `reason` |
| **Type** | `counter` |
| **Fully Qualified Name** | `relay_get_chunks_rate_limited_count` |
