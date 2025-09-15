# Client Ledger

The `clientledger` package manages payment state for clients making dispersal requests.

## Concepts

- Client Ledger: Each client is responsible for tracking EigenDA usage for their own account. Depending on the
configured payment mode, a client may have to keep track of reservation usage, on-demand payments, or both.
- Sources of truth: The payment tracking performed by a client represents a local view of the "actual" payment state,
which is maintained by the Validator Nodes (for reservation payments), and the EigenDA Disperser (for on-demand
payments). Clients maintain a local reckoning of payment state to be able to decide which payment method to utilize
for any given dispersal, and to be able to know how much data can be dispersed.

## Files

- `client_ledger.go` - Manages payment state for a single client account, for both reservation and on-demand payments
- `client_ledger_mode.go` - Defines which payments are configured for a given client ledger
