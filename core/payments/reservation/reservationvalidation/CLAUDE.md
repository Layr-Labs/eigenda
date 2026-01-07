# Reservation Payment Validation

The `reservationvalidation` package contains utilities used by Dispersers and Validators, for validating reservation
payments for multiple accounts at the same time.

## Files

- `reservation_payment_validator.go` - Validates reservation payments for multiple accounts
- `reservation_ledger_cache.go` - LRU cache for storing a collection of `ReservationLedger`s, used by the
`ReservationPaymentValidator`
- `reservation_ledger_cache_config.go` - Configuration parameters for the `ReservationLedgerCache`
- `reservation_validator_metrics.go` - Metrics for reservation payment validation
- `reservation_cache_metrics.go` - Metrics for the LRU ledger cache
