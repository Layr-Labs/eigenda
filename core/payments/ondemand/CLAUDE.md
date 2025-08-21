# On-Demand Payments

The on-demand package implements accounting logic for on-demand EigenDA usage.

## Concepts

- On-demand payments: users deposit funds on-chain in the `PaymentVault` contract, and these funds are used
to pay for blobs as they are dispersed.

## Files

- `on_demand_ledger.go` - Tracks cumulative payment state for on-demand dispersals for a single account
- `on_demand_ledgers.go` - Tracks cumulative payment state for on-demand dispersals for multiple accounts
- `cumulative_payment_store.go` - Interface for storing and retrieving cumulative payment state
- `errors.go` - Error types for on-demand payment failures

## Sub-packages

- `dynamostore` - contains a DynamoDB implementation of the `CumulativePaymentStore`
- `ephemeral` - contains an in-memory implementation of the `CumulativePaymentStore`
