# Metrics Documentation for namespace 'relay'

This documentation was automatically generated at time `2024-11-27T12:13:08-06:00`

There are a total of `34` registered metrics.

---

## average_get_blob_data_bytes

Average data size of requested blobs

|   |   |
|---|---|
| **Name** | `average_get_blob_data` |
| **Unit** | `bytes` |
| **Type** | `running average` |
| **Time Window** | `1m0s` |
| **Fully Qualified Name** | `relay_average_get_blob_data_bytes` |
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

## blob_cache_size

Number of items in the blob cache

|   |   |
|---|---|
| **Name** | `blob_cache` |
| **Unit** | `size` |
| **Type** | `gauge` |
| **Fully Qualified Name** | `relay_blob_cache_size` |
---

## blob_cache_weight

Total weight of items in the blob cache

|   |   |
|---|---|
| **Name** | `blob_cache` |
| **Unit** | `weight` |
| **Type** | `gauge` |
| **Fully Qualified Name** | `relay_blob_cache_weight` |
---

## blob_cache_average_weight

Average weight of items currently in the blob cache

|   |   |
|---|---|
| **Name** | `blob_cache_average` |
| **Unit** | `weight` |
| **Type** | `gauge` |
| **Fully Qualified Name** | `relay_blob_cache_average_weight` |
---

## blob_cache_average_lifespan_ms

Average time an item remains in the blob cache before being evicted.

|   |   |
|---|---|
| **Name** | `blob_cache_average_lifespan` |
| **Unit** | `ms` |
| **Type** | `running average` |
| **Time Window** | `1m0s` |
| **Fully Qualified Name** | `relay_blob_cache_average_lifespan_ms` |
---

## blob_cache_hit_count

Number of cache hits in the blob cache

|   |   |
|---|---|
| **Name** | `blob_cache_hit` |
| **Unit** | `count` |
| **Type** | `counter` |
| **Fully Qualified Name** | `relay_blob_cache_hit_count` |
---

## blob_cache_miss_count

Number of cache misses in the blob cache

|   |   |
|---|---|
| **Name** | `blob_cache_miss` |
| **Unit** | `count` |
| **Type** | `counter` |
| **Fully Qualified Name** | `relay_blob_cache_miss_count` |
---

## blob_cache_miss_latency_ms

Latency of cache misses in the blob cache

|   |   |
|---|---|
| **Name** | `blob_cache_miss_latency` |
| **Unit** | `ms` |
| **Type** | `latency` |
| **Quantiles** | `0.500`, `0.900`, `0.990` |
| **Fully Qualified Name** | `relay_blob_cache_miss_latency_ms` |
---

## chunk_cache_size

Number of items in the chunk cache

|   |   |
|---|---|
| **Name** | `chunk_cache` |
| **Unit** | `size` |
| **Type** | `gauge` |
| **Fully Qualified Name** | `relay_chunk_cache_size` |
---

## chunk_cache_weight

Total weight of items in the chunk cache

|   |   |
|---|---|
| **Name** | `chunk_cache` |
| **Unit** | `weight` |
| **Type** | `gauge` |
| **Fully Qualified Name** | `relay_chunk_cache_weight` |
---

## chunk_cache_average_weight

Average weight of items currently in the chunk cache

|   |   |
|---|---|
| **Name** | `chunk_cache_average` |
| **Unit** | `weight` |
| **Type** | `gauge` |
| **Fully Qualified Name** | `relay_chunk_cache_average_weight` |
---

## chunk_cache_average_lifespan_ms

Average time an item remains in the chunk cache before being evicted.

|   |   |
|---|---|
| **Name** | `chunk_cache_average_lifespan` |
| **Unit** | `ms` |
| **Type** | `running average` |
| **Time Window** | `1m0s` |
| **Fully Qualified Name** | `relay_chunk_cache_average_lifespan_ms` |
---

## chunk_cache_hit_count

Number of cache hits in the chunk cache

|   |   |
|---|---|
| **Name** | `chunk_cache_hit` |
| **Unit** | `count` |
| **Type** | `counter` |
| **Fully Qualified Name** | `relay_chunk_cache_hit_count` |
---

## chunk_cache_miss_count

Number of cache misses in the chunk cache

|   |   |
|---|---|
| **Name** | `chunk_cache_miss` |
| **Unit** | `count` |
| **Type** | `counter` |
| **Fully Qualified Name** | `relay_chunk_cache_miss_count` |
---

## chunk_cache_miss_latency_ms

Latency of cache misses in the chunk cache

|   |   |
|---|---|
| **Name** | `chunk_cache_miss_latency` |
| **Unit** | `ms` |
| **Type** | `latency` |
| **Quantiles** | `0.500`, `0.900`, `0.990` |
| **Fully Qualified Name** | `relay_chunk_cache_miss_latency_ms` |
---

## get_blob_data_latency_ms

Latency of the GetBlob RPC data retrieval

|   |   |
|---|---|
| **Name** | `get_blob_data_latency` |
| **Unit** | `ms` |
| **Type** | `latency` |
| **Quantiles** | `0.500`, `0.900`, `0.990` |
| **Fully Qualified Name** | `relay_get_blob_data_latency_ms` |
---

## get_blob_latency_ms

Latency of the GetBlob RPC

|   |   |
|---|---|
| **Name** | `get_blob_latency` |
| **Unit** | `ms` |
| **Type** | `latency` |
| **Quantiles** | `0.500`, `0.900`, `0.990` |
| **Fully Qualified Name** | `relay_get_blob_latency_ms` |
---

## get_blob_metadata_latency_ms

Latency of the GetBlob RPC metadata retrieval

|   |   |
|---|---|
| **Name** | `get_blob_metadata_latency` |
| **Unit** | `ms` |
| **Type** | `latency` |
| **Quantiles** | `0.500`, `0.900`, `0.990` |
| **Fully Qualified Name** | `relay_get_blob_metadata_latency_ms` |
---

## get_blob_rate_limited_count

Number of GetBlob RPC rate limited

|   |   |
|---|---|
| **Name** | `get_blob_rate_limited` |
| **Unit** | `count` |
| **Labels** | `reason` |
| **Type** | `counter` |
| **Fully Qualified Name** | `relay_get_blob_rate_limited_count` |
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
---

## metadata_cache_size

Number of items in the metadata cache

|   |   |
|---|---|
| **Name** | `metadata_cache` |
| **Unit** | `size` |
| **Type** | `gauge` |
| **Fully Qualified Name** | `relay_metadata_cache_size` |
---

## metadata_cache_weight

Total weight of items in the metadata cache

|   |   |
|---|---|
| **Name** | `metadata_cache` |
| **Unit** | `weight` |
| **Type** | `gauge` |
| **Fully Qualified Name** | `relay_metadata_cache_weight` |
---

## metadata_cache_average_weight

Average weight of items currently in the metadata cache

|   |   |
|---|---|
| **Name** | `metadata_cache_average` |
| **Unit** | `weight` |
| **Type** | `gauge` |
| **Fully Qualified Name** | `relay_metadata_cache_average_weight` |
---

## metadata_cache_average_lifespan_ms

Average time an item remains in the metadata cache before being evicted.

|   |   |
|---|---|
| **Name** | `metadata_cache_average_lifespan` |
| **Unit** | `ms` |
| **Type** | `running average` |
| **Time Window** | `1m0s` |
| **Fully Qualified Name** | `relay_metadata_cache_average_lifespan_ms` |
---

## metadata_cache_hit_count

Number of cache hits in the metadata cache

|   |   |
|---|---|
| **Name** | `metadata_cache_hit` |
| **Unit** | `count` |
| **Type** | `counter` |
| **Fully Qualified Name** | `relay_metadata_cache_hit_count` |
---

## metadata_cache_miss_count

Number of cache misses in the metadata cache

|   |   |
|---|---|
| **Name** | `metadata_cache_miss` |
| **Unit** | `count` |
| **Type** | `counter` |
| **Fully Qualified Name** | `relay_metadata_cache_miss_count` |
---

## metadata_cache_miss_latency_ms

Latency of cache misses in the metadata cache

|   |   |
|---|---|
| **Name** | `metadata_cache_miss_latency` |
| **Unit** | `ms` |
| **Type** | `latency` |
| **Quantiles** | `0.500`, `0.900`, `0.990` |
| **Fully Qualified Name** | `relay_metadata_cache_miss_latency_ms` |
