# On-Demand Payments

The `ondemand` package implements accounting logic for on-demand EigenDA usage.

## Concepts

- On-demand payments: users deposit funds on-chain in the `PaymentVault` contract, and these funds are used
to pay for blobs as they are dispersed.
- Source of truth: the EigenDA Disperser is the source of truth for on-demand payments. Validator nodes do not validate
on-demand payments. *Only* the EigenDA Disperser supports on-demand payments: all other Dispersers are limited to
reservation payments. When a client starts up, it must fetch the latest on-demand payment state from the EigenDA
Disperser to be able to make on-demand dispersals.

## Subpackages

- `ondemandvalidation` - Contains utilities used by Dispersers and Validators, for validating on-demand payments for
multiple accounts at the same time.

## Files

- `on_demand_ledger.go` - Tracks cumulative payment state for on-demand dispersals for a single account
- `on_demand_vault_monitor.go` - Monitors `PaymentVault` contract for deposit updates
- `cumulative_payment_store.go` - Struct for storing and retrieving cumulative payment state in/from DynamoDB
- `errors.go` - Error types for on-demand payment failures
