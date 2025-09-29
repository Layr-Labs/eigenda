# Payments

The payments package contains the logic for how clients pay for blob dispersals.

## Concepts

There are two possible ways to pay for a blob dispersal:

1. Reservation (logic in the `reservation` sub-package)
2. On-demand (logic in the `ondemand` sub-package)

## Historical Context

- The logic in this package is a reimplementation of pre-existing payment logic
- The new implementation is being added parallel to the old implementation
- The old implementation will be removed once migration to the new system is complete
- Old payment logic exists in the following places:
    - `api/clients/v2/accountant.go`
    - `core/meterer/`
