# Metrics Documentation for namespace 'eigenda_node'

This documentation was automatically generated at time `2024-12-03T11:40:32-06:00`

There are a total of `5` registered metrics.

---

## db_size_bytes

The size of the leveldb database.

|   |   |
|---|---|
| **Name** | `db_size` |
| **Unit** | `bytes` |
| **Type** | `gauge` |
| **Fully Qualified Name** | `eigenda_node_db_size_bytes` |
---

## get_chunks_data_size_bytes

The size of the data requested to be retrieved by GetChunks() RPC calls.

|   |   |
|---|---|
| **Name** | `get_chunks_data_size` |
| **Unit** | `bytes` |
| **Type** | `gauge` |
| **Fully Qualified Name** | `eigenda_node_get_chunks_data_size_bytes` |
---

## get_chunks_latency_ms

The latency of a GetChunks() RPC call.

|   |   |
|---|---|
| **Name** | `get_chunks_latency` |
| **Unit** | `ms` |
| **Type** | `latency` |
| **Quantiles** | `0.500`, `0.900`, `0.990` |
| **Fully Qualified Name** | `eigenda_node_get_chunks_latency_ms` |
---

## store_chunks_data_size_bytes

The size of the data requested to be stored by StoreChunks() RPC calls.

|   |   |
|---|---|
| **Name** | `store_chunks_data_size` |
| **Unit** | `bytes` |
| **Type** | `gauge` |
| **Fully Qualified Name** | `eigenda_node_store_chunks_data_size_bytes` |
---

## store_chunks_latency_ms

The latency of a StoreChunks() RPC call.

|   |   |
|---|---|
| **Name** | `store_chunks_latency` |
| **Unit** | `ms` |
| **Type** | `latency` |
| **Quantiles** | `0.500`, `0.900`, `0.990` |
| **Fully Qualified Name** | `eigenda_node_store_chunks_latency_ms` |
