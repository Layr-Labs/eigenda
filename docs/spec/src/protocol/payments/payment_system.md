# EigenDA Payment System

## 1. Overview

The EigenDA payment system allows users to pay for blob dispersals through two methods: reservations and on-demand
payments. All payment logic is implemented in the [`core/payments`](../../../../../core/payments/) package.

**Key Concepts:**
- Blob sizes are measured in *symbols*, where each symbol contains 32 bytes of data
- Blob sizes are measured **post-blob encoding**: user payloads expand to some degree during blob encoding
- Blob sizes are constrained to powers-of-two: dispersals are rounded up to the next power-of-two number of symbols
  when computing size
- The [PaymentVault](https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/core/PaymentVault.sol) contract
  stores all on-chain payment-related data:
  - User reservation parameters
  - User on-demand deposits
  - Global payment parameters, including but not limited to:
    - `minNumSymbols`: dispersals smaller than this threshold are billed as if they were `minNumSymbols` in size
    - `pricePerSymbol`: the price per symbol (in wei) for on-demand payments

## 2. Payment Methods

### 2.1 Reservation Payments

- Reservations provide guaranteed bandwidth for a specified time period.
   - Users reserve capacity in advance, and must "use it or lose it".
   - Reservations are procured out-of-band, through EigenDA
- The system uses a [leaky bucket algorithm](../../../../../core/payments/reservation/leaky_bucket.go) to manage usage:
  symbols are added to the bucket each time a blob is dispersed, and these leak out over time. A user can only make a
  dispersal if the leaky bucket has available capacity.
   - The total capacity of the leaky bucket is defined in terms of a *duration*. The size of the bucket in symbols is
`reservationRate * bucketDuration`. This calculation controls the burstiness of reservation usage.
- Parameters describing active user reservations are kept in the `PaymentVault` contract
   - `symbolsPerSecond`: the reservation bandwidth rate
   - `startTimestamp` and `endTimestamp`: define when the reservation is active
   - `quorumNumbers`: which quorums the reservation can be used for

#### 2.1.1 Source of Truth

- Validator nodes are the source of truth for reservation usage.
   - Each validator keeps track of the dispersals from each user account, and will reject dispersals if the user doesn't
   have enough capacity
   - Clients keep a local reckoning of their own reservation usage, so that they can stay within the bounds of their
   reserved bandwidth
   - Dispersers also keep a local reckoning of client reservation usage, but a malicious client can bypass this check
   by intentionally dispersing too much data spread out over multiple dispersers. From the perspective of any given
   disperser, the client is within reservation limits. But in total, the client is over the limit. This isn't a problem,
   because validator nodes will catch the misbehavior. By having dispersers keep track of reservation usage, we are 
   imposing a limit on how severely a client can misbehave in this way: in a system with N dispersers, a malicious
   client can disperse at most N * reservation rate
- Reservation usage state agreement: since clients keep a local reckoning of reservation usage without any input from
validators, it's all but guaranteed that their local state will differ (at least slightly) from the state on any given
validator. This actually doesn't present a problem, so long as these key invariants are maintained:
   - A client behaving honestly must be able to disperse blobs without payment failures
   - The amount of "free" dispersals that can be stolen by a dishonest client must be tightly limited

#### 2.1.2 Bucket Capacity Configuration

- We can achieve these invariants by using buckets of differing sizes between clients and validator nodes. If we make
validator buckets larger than client buckets by some multiple, then slight discrepancies between client and
validator are naturally smoothed out.
- If a dishonest client tries to disperse more data than allowed, the behavior will be permitted by validators for
a short time, but eventually even the larger validator bucket will fill. At that point, validators will limit new
dispersals from the dishonest client to the rate of the reservation, and no additional dispersals may be stolen.
- The capacity difference between client and validator buckets must be chosen to accommodate the maximum
expected latency of the system. Specifically:
`validatorBucketCapacity - clientBucketCapacity = reservationRate * maxSystemLatency`
- This ensures that honest clients operating at full capacity won't be rejected due to timing discrepancies.
- Proposed bucket sizing (actual configuration may vary): Client buckets will use 1 minute duration while validator
buckets use 6 minutes, accommodating up to 5 minutes of system latency.
- Validators may potentially be configured to leak buckets slightly faster (e.g., 1% faster) than the actual reservation
rate. This causes validator bucket states to converge toward empty for honest clients operating within their reservation
limits, as the faster leak rate ensures buckets tracked by the validators drain over time.
   - The trade-off is that dishonest clients could potentially disperse up to 1% more data than their allotted
   reservation.

#### 2.1.3 Leaky Bucket Overfill

The reservation leaky bucket implementation permits clients to overfill their buckets, with certain constraints:
- If a client has *any* available capacity in their bucket, they may make a single dispersal up to the maximum blob
size, even if that dispersal causes the bucket to exceed its maximum capacity
- When this happens, the bucket level actually goes above the maximum capacity, and the client must wait for the
bucket to leak back down below full capacity before making the next dispersal
- This feature exists to solve a problem with small reservations: without overfill, a reservation might be so small
that its total bucket capacity is less than the max blob size, which would prevent the user from dispersing blobs up
to max size.
- By permitting a single overfill, even the smallest reservation can disperse blobs of maximum size

#### 2.1.4 Reservation Usage Persistence

The leaky bucket algorithm does not require persisting reservation usage state across system restarts. Different
system components initialize their buckets with opposing biases to maintain system integrity without persistence:

**Client Initialization (Conservative Bias)**
- Clients initialize their leaky bucket as completely full (no capacity available) upon restart
- They must wait for symbols to leak out before dispersing, guaranteeing compliance with reservation rate limits
- While this may result in slight underutilization if usage was low before restart, it prevents violation of
reservation limits

**Validator Initialization (Permissive Bias)**
- Validators initialize leaky buckets as completely empty (full capacity available) upon restart
- This ensures they never incorrectly deny service to users entitled to a reservation
- In the worst case, a malicious client timing dispersals with validator restarts might be able to cause a small amount
of extra work for that specific validator

This dual-bias approach eliminates the complexity of distributed reservation state persistence.

### 2.2 On-Demand Payments

- On-demand payments allow users to pay per dispersal from funds deposited in the PaymentVault contract
   - Once deposited, funds cannot be withdrawn - they can only be used for dispersals or abandoned
- Limited to quorums 0 (ETH) and 1 (EIGEN)
   - Custom quorums are not supported for on-demand payments because quorum resources are closely tailored to expected
   usage. Allowing on-demand payments could enable third parties to overutilize these limited resources.
- Costs are calculated based on blob size (in symbols) multiplied by the `pricePerSymbol` parameter in PaymentVault
- Payment usage is not tracked on-chain; instead, the EigenDA Disperser maintains a DynamoDB table recording total
historical usage for all clients
- When processing a dispersal, the Disperser compares a user's total historical usage against their on-chain deposits
in the PaymentVault to determine if they have sufficient funds
- Clients fetch the latest cumulative payment state from the EigenDA Disperser on startup via the `GetPaymentState` RPC

#### 2.2.1 Why Only the EigenDA Disperser?

- On-demand payments are supported only through the EigenDA Disperser
- Since EigenDA currently lacks a consensus mechanism, validators cannot easily coordinate to limit total on-demand
throughput across the network
- Therefore, the EigenDA Disperser fills the role of arbiter, ensuring that total network throughput doesn't exceed
configured levels

#### 2.2.2 Cumulative Payment

The cumulative payment is a field set in the PaymentHeader by the client when making a dispersal. It represents the
total cost (in wei) of all previous dispersals, plus the new dispersal.

- **Historical context:** In a prior payments implementation, the cumulative payment field included by the client in
  the PaymentHeader had to be monotonically increasing, and the disperser would verify that each new cumulative payment
  received exceeded the previous one by at least the cost of the new blob. This severely limited concurrency, since
  clients had to make sure that all on-demand dispersals were handled by the disperser in strict order. In practice,
  that meant waiting for the entire network roundtrip, for dispersal N to be confirmed before submitting dispersal N+1.
- **Current implementation:** The system has been simplified to improve concurrency:
   - Clients still populate the `cumulative_payment` field with their local calculation of cumulative payment
   - However, the Disperser now only checks if this field is non-zero (to determine payment type) and ignores the exact
   value
   - The Disperser tracks each account's on-demand usage in DynamoDB, incrementing by the blob cost for each dispersal
   - This removes the strict ordering requirement and allows for highly concurrent dispersals
- **Why clients still populate the field:** Although currently unused beyond the zero/non-zero check, clients continue
  to populate this field with meaningful values. This preserves the option to reintroduce cumulative payment validation
  in the future if needed.

## 3. Client Payment Strategy

### 3.1 Payment Header

Each dispersal request includes a [PaymentHeader](../../../../../api/proto/common/v2/common_v2.proto#L111) containing:
- `account_id`: Ethereum address identifying the payment account
- `timestamp`: Nanosecond UNIX timestamp (serves as nonce)
- `cumulative_payment`: Variable-length big-endian uint for on-demand dispersal (or empty for reservation dispersal)

The payment header implicitly specifies which payment mechanism is being used:
- If `cumulative_payment` is empty/zero → reservation payment
- If `cumulative_payment` is non-zero → on-demand payment

### 3.2 Client Configuration Options

Clients can configure their payment strategy in three ways:

1. **Reservation-only:** Client exclusively uses reservation payments
   - `cumulative_payment` field always left empty
   - Dispersals fail if reservation capacity is exhausted

2. **On-demand-only:** Client exclusively uses on-demand payments
   - `cumulative_payment` field always populated
   - All dispersals charged against deposited balance

3. **Hybrid with fallback:** Client uses both payment methods
   - Primary: Uses reservation while capacity is available
   - Fallback: Automatically switches to on-demand when reservation is exhausted
   - Ensures continuous operation without manual intervention
