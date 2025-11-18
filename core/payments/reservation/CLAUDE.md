# Reservation Payments

The `reservation` package implements accounting logic for reservation-based EigenDA usage.

## Concepts

- Reservation payments: User reservation parameters are recorded on-chain in the `PaymentVault` contract. A
reservation represents a conceptual "leaky bucket", where each blob dispersal adds tokens that leak out over time.
Dispersals can only be made when there is enough available capacity in the bucket.
- Source of truth: Validator nodes are the source of truth for reservation payments. Clients and dispersers keep
a local reckoning of reservation data usage which approximates the source of truth that exists within the Validator
network. The reservation payment system is designed and implemented in such a way that an approximation is sufficient
to be able to make reservation-based dispersals to the EigenDA network.

## Subpackages

- `reservationvalidation` - Contains utilities used by Dispersers and Validators, for validating reservation payments
for multiple accounts at the same time.

## Files

- `reservation.go` - Describes parameters of a single account's reservation
- `reservation_ledger.go` - Tracks usage of a single account's reservation
- `reservation_vault_monitor.go` - Monitors `PaymentVault` contract for reservation updates
- `leaky_bucket.go` - Rate limiting algorithm utility, utilized by the `ReservationLedger`
- `reservation_ledger_config.go` - Configures a `ReservationLedger`
- `overfill_behavior.go` - Defines how bucket overfills are handled
- `errors.go` - Sentinel errors for reservation related failures
