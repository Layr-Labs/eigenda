# Reservation Payments

The reservation package implements accounting logic for reservation-based EigenDA usage.

## Key Files

- `reservation.go` - Describes parameters of a single account's reservation
- `reservation_ledger.go` - Tracks usage of a single account's reservation
- `leaky_bucket.go` - Rate limiting algorithm utility, utilized by the `ReservationLedger`
- `reservation_ledger_config.go` - Configures a `ReservationLedger`
- `errors.go` - Sentinel errors for reservation related failures
