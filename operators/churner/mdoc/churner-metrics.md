# Metrics Documentation for namespace 'eigenda_churner'

This documentation was automatically generated at time `2024-11-26T14:29:13-06:00`

There are a total of `2` registered metrics.

---

## latency_ms

latency summary in milliseconds

|   |   |
|---|---|
| **Name** | `latency` |
| **Unit** | `ms` |
| **Labels** | `method` |
| **Type** | `latency` |
| **Quantiles** | `0.500`, `0.900`, `0.950`, `0.990` |
| **Fully Qualified Name** | `eigenda_churner_latency_ms` |
---

## request_count

the number of requests

|   |   |
|---|---|
| **Name** | `request` |
| **Unit** | `count` |
| **Labels** | `status`, `method`, `reason` |
| **Type** | `counter` |
| **Fully Qualified Name** | `eigenda_churner_request_count` |
