# On-Demand Payments

The on-demand package implements accounting logic for on-demand EigenDA usage.

## Key Files

- `on_demand_ledger.go` - Tracks cumulative payment state for on-demand dispersals for a single account
- `cumulative_payment_store.go` - Struct for storing and retrieving cumulative payment state in/from DynamoDB
- `errors.go` - Error types for on-demand payment failures
