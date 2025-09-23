# Vault

The `vault` package contains utilities for interacting with the `PaymentVault` contract.

## Concepts

- `PaymentVault`: This is the [EigenDA ethereum contract](../../../../contracts/src/core/PaymentVault.sol) that defines
global payment parameters, reservations that have been allocated to users, and keeps track of user deposits that can be
used for on-demand dispersal.

## Files

- `payment_vault.go` - Provides methods for interacting with the `PaymentVault` contract
- `test_payment_vault.go` - Test implementation of `PaymentVault`
