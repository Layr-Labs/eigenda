# Reservation Payments

The reservation package implements accounting logic for reservation-based EigenDA usage.

## Key Implementation Details

- Reservation accounting is performed with a LeakyBucket algorithm.
- Each instance of the LeakyBucket algorithm is configured with a BiasBehavior, to determine whether to err
on the side of permitting more or less throughput.
- Each instance of the LeakyBucket algorithm is configured with an OverfillBehavior, which governs behavior when bucket
capacity is exceeded.
