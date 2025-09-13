# EigenDA Payment System Migration

## 1. Overview

EigenDA is migrating from a fixed bin reservation accounting model to a leaky bucket algorithm. The new payment system
is being implemented to be compatible with permissionless dispersal. While making changes to support this new feature,
the opportunity to reduce accumulated tech debt is being seized.

**Key Changes:**
- Reservation accounting switches from fixed time bins to continuous leaky bucket rate limiting
- Validators become the source of truth for reservation metering (previously the EigenDA Disperser)
- On-demand payment logic remains unchanged
- Payment logic is being reorganized or reimplemented, to reduce tech debt

## 2. Legacy Payment System

The legacy implementation uses a **fixed bin model** where:
- Users disperse against reservation bandwidth allotted for the current fixed time bin
- Once capacity for the current bin is exhausted, users must wait for the next bin to arrive, to disperse more data
- Implementation split between [`core/meterer/`](../../../../../core/meterer/) and
  [`api/clients/v2/accountant.go`](../../../../../api/clients/v2/accountant.go)

**Weaknesses:**

- Bursty behavior at bin boundaries creates uneven load distribution
- Network-wide bin synchronization causes simultaneous bursts across all users, exacerbating the problem of bursts

## 3. New Payment System

Reservation payments will be managed with a [leaky bucket](../../../../../core/payments/reservation/leaky_bucket.go)
algorithm, instead of using fixed bins. This alternate algorithm smooths out bursts with smooth capacity recovery.

- Less severe bursts for each individual user: the maximum burst size from a single user is now limited by the size
of the leaky bucket, compared to the fixed bin algorithm where maximum burst is 2x bin size
- Network wide bursts are unlikely to be simultaneous, since there aren't synced bin boundaries

## 4. Migration Considerations

### Requirements
- **Backward compatibility:** Old clients will work seamlessly with new disperser logic
   - Users operating well below reservation limits will experience no interruption, and may choose to update clients
   whenever convenient
   - Users operating near reservation limits may experience some degraded behavior, if the local algorithm disagrees
   with the updated remote algorithm. Such users may resolve degraded behavior by updating client code to match the
   new algorithm.
- **Gradual rollout:** Phased deployment with feature flags for safety

## 5. Migration Rollout Process

### Phase 1: Client & Disperser Release
1. **Client release** with leaky bucket accounting
2. **Disperser release** with leaky bucket support

### Phase 2: Validator Release
- Deploy after client/disperser adoption complete
- Feature flag controls activation
- Once this phase is complete, validators have become authoritative source for reservation metering
   - This must occur before a second Disperser is brought online
