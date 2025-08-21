# Reservation Payments

The reservation package implements accounting logic for reservation-based EigenDA usage.

## Concepts

- Reservation payments: User reservation parameters are recorded on-chain in the `PaymentVault` contract. A
reservation represents a conceptual "leaky bucket", where each user dispersal adds tokens that leak out over time.
Dispersals can only be made when there is enough available capacity in the bucket.

## Files

- `reservation.go` - Describes parameters of a single account's reservation
- `reservation_ledger.go` - Tracks usage of a single account's reservation
- `leaky_bucket.go` - Rate limiting algorithm utility, utilized by the `ReservationLedger`
- `reservation_ledger_config.go` - Configures a `ReservationLedger`
- `errors.go` - Sentinel errors for reservation related failures
