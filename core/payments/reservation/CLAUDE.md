# Reservation Payments

The reservation package implements accounting logic for reservation-based EigenDA usage.

## Key Implementation Details

- Reservation accounting is performed with a LeakyBucket algorithm.
- Each instance of the LeakyBucket algorithm is configured with a boolean startFull parameter to determine 
  whether the bucket starts full (requiring leakage before use) or empty (available for immediate use).
- Each instance of the LeakyBucket algorithm is configured with an OverfillBehavior, which governs behavior when bucket
  capacity is exceeded.
