# Metrics Documentation for namespace 'test'

This documentation was automatically generated at time `2024-11-25T10:13:43-06:00`

There are a total of `7` registered metrics.

---

## c1_count

this metric shows the number of times the sleep cycle has been executed

|   |   |
|---|---|
| **Name** | `c1` |
| **Unit** | `count` |
| **Label** | - |
| **Type** | `counter` |
| **Fully Qualified Name** | `test_c1_count` |
---

## c1_count: DOUBLE

this metric shows the number of times the sleep cycle has been executed, doubled

|   |   |
|---|---|
| **Name** | `c1` |
| **Unit** | `count` |
| **Label** | `DOUBLE` |
| **Type** | `counter` |
| **Fully Qualified Name** | `test_c1_count` |
---

## g1_milliseconds

this metric shows the duration of the most recent sleep cycle

|   |   |
|---|---|
| **Name** | `g1` |
| **Unit** | `milliseconds` |
| **Label** | - |
| **Type** | `gauge` |
| **Fully Qualified Name** | `test_g1_milliseconds` |
---

## g1_milliseconds: autoPoll

this metric shows the sum of all sleep cycles

|   |   |
|---|---|
| **Name** | `g1` |
| **Unit** | `milliseconds` |
| **Label** | `autoPoll` |
| **Type** | `gauge` |
| **Fully Qualified Name** | `test_g1_milliseconds` |
---

## g1_milliseconds: previous

this metric shows the duration of the second most recent sleep cycle

|   |   |
|---|---|
| **Name** | `g1` |
| **Unit** | `milliseconds` |
| **Label** | `previous` |
| **Type** | `gauge` |
| **Fully Qualified Name** | `test_g1_milliseconds` |
---

## l1_ms

this metric shows the latency of the sleep cycle

|   |   |
|---|---|
| **Name** | `l1` |
| **Unit** | `seconds` |
| **Label** | - |
| **Type** | `latency` |
| **Quantiles** | `0.500`, `0.900`, `0.990` |
| **Fully Qualified Name** | `test_l1_ms` |
---

## l1_ms: HALF

this metric shows the latency of the sleep cycle, divided by two

|   |   |
|---|---|
| **Name** | `l1` |
| **Unit** | `seconds` |
| **Label** | `HALF` |
| **Type** | `latency` |
| **Quantiles** | `0.500`, `0.900`, `0.990` |
| **Fully Qualified Name** | `test_l1_ms` |
