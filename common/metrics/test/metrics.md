# Metrics Documentation for namespace 'test'

This documentation was automatically generated at time `2024-11-25T12:46:49-06:00`

There are a total of `5` registered metrics.

---

## c1_count

this metric shows the number of times the sleep cycle has been executed

|   |   |
|---|---|
| **Name** | `c1` |
| **Unit** | `count` |
| **Labels** | `X`, `Y`, `Z` |
| **Type** | `counter` |
| **Fully Qualified Name** | `test_c1_count` |
---

## c2_count

the purpose of this counter is to test what happens if we don't provide a label template

|   |   |
|---|---|
| **Name** | `c2` |
| **Unit** | `count` |
| **Type** | `counter` |
| **Fully Qualified Name** | `test_c2_count` |
---

## g1_milliseconds

this metric shows the duration of the most recent sleep cycle

|   |   |
|---|---|
| **Name** | `g1` |
| **Unit** | `milliseconds` |
| **Labels** | `foo`, `bar`, `baz` |
| **Type** | `gauge` |
| **Fully Qualified Name** | `test_g1_milliseconds` |
---

## g2_milliseconds

this metric shows the sum of all sleep cycles

|   |   |
|---|---|
| **Name** | `g2` |
| **Unit** | `milliseconds` |
| **Labels** | `X`, `Y`, `Z` |
| **Type** | `gauge` |
| **Fully Qualified Name** | `test_g2_milliseconds` |
---

## l1_ms

this metric shows the latency of the sleep cycle

|   |   |
|---|---|
| **Name** | `l1` |
| **Unit** | `ms` |
| **Labels** | `foo`, `bar`, `baz` |
| **Type** | `latency` |
| **Quantiles** | `0.500`, `0.900`, `0.990` |
| **Fully Qualified Name** | `test_l1_ms` |
