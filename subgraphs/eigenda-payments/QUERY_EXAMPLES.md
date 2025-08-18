# EigenDA Payments Subgraph Query Examples

## Event-based Queries for Reservations

Querying reservation updates can be done using the `reservationUpdateds` event. This event captures all updates to reservations, including changes in symbols per second, start and end timestamps, and more.

### Query All Reservations
To retrieve all reservations updates for a given account, you can use the following GraphQL query:

```graphql
query ReservationUpdatesForAccount($account: Bytes!) {
  reservationUpdateds(
  	where: {
  		account: $account
  	}
  )
  {
    transactionHash
    blockNumber
    reservation_startTimestamp
    reservation_endTimestamp
    reservation_quorumSplits
    reservation_quorumNumbers
    reservation_symbolsPerSecond
  }
}
```

## Timestamp-based Reservation Filtering

Since reservations never get deleted on-chain we use timestamp-based filtering in queries to determine reservation status. Note: the `reservations` entity
is used to represent the latest state of reservations, which is updated based on the latest reservation update events.

### Query Active Reservations

To find all currently active reservations, filter by comparing the current timestamp with start and end times:

```graphql
query ActiveReservations($currentTime: BigInt!) {
  reservations(
    where: { 
      startTimestamp_lte: $currentTime,
      endTimestamp_gt: $currentTime 
    }
  ) {
    account
    symbolsPerSecond
    startTimestamp
    endTimestamp
    quorumNumbers
    quorumSplits
  }
}
```

Variables:
```json
{
  "currentTime": "1699564800"  // Unix timestamp in seconds
}
```

### Query Pending Reservations

To find reservations that haven't started yet:

```graphql
query PendingReservations($currentTime: BigInt!) {
  reservations(
    where: { 
      startTimestamp_gt: $currentTime 
    }
  ) {
    account
    startTimestamp
    endTimestamp
    symbolsPerSecond
  }
}
```

### Query Expired Reservations

To find reservations that have already ended:

```graphql
query ExpiredReservations($currentTime: BigInt!) {
  reservations(
    where: { 
      endTimestamp_lte: $currentTime 
    }
  ) {
    account
    startTimestamp
    endTimestamp
    lastUpdatedTimestamp
  }
}
```

### Query Reservations by Account

To get a specific account's reservation:

```graphql
query AccountReservation($account: Bytes!) {
  reservation(id: $account) {
    symbolsPerSecond
    startTimestamp
    endTimestamp
    quorumNumbers
    quorumSplits
    lastUpdatedBlock
    lastUpdatedTimestamp
  }
}
```
