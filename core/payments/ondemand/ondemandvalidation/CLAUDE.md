# On-Demand Payment Validation

The `ondemandvalidation` package contains utilities used by Dispersers and Validators, for validating on-demand
payments for multiple accounts at the same time.

## Files

- `on_demand_payment_validator.go` - Validates on-demand payments for multiple accounts
- `on_demand_ledger_cache.go` - LRU cache for storing a collection of `OnDemandLedger`s, used by the
`OnDemandPaymentValidator`
- `on_demand_ledger_cache_config.go` - Configuration parameters for the `OnDemandLedgerCache`
- `on_demand_validator_metrics.go` - Metrics for on-demand payment validation
- `on_demand_cache_metrics.go` - Metrics for the LRU ledger cache
