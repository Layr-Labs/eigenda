import {
  assert,
  describe,
  test,
  clearStore,
  beforeAll,
  afterAll
} from "matchstick-as/assembly/index"
import { BigInt, Address, Bytes } from "@graphprotocol/graph-ts"
import { GlobalRatePeriodIntervalUpdated, ActiveReservation } from "../generated/schema"
import { handleGlobalRatePeriodIntervalUpdated, handleReservationUpdated } from "../src/payment-vault"
import { createGlobalRatePeriodIntervalUpdatedEvent, createReservationUpdatedEvent } from "./payment-vault-utils"

// Tests structure (matchstick-as >=0.5.0)
// https://thegraph.com/docs/en/subgraphs/developing/creating/unit-testing-framework/#tests-structure

describe("Describe entity assertions", () => {
  beforeAll(() => {
    let previousValue = BigInt.fromI32(234)
    let newValue = BigInt.fromI32(234)
    let newGlobalRatePeriodIntervalUpdatedEvent =
      createGlobalRatePeriodIntervalUpdatedEvent(previousValue, newValue)
    handleGlobalRatePeriodIntervalUpdated(
      newGlobalRatePeriodIntervalUpdatedEvent
    )
  })

  afterAll(() => {
    clearStore()
  })

  // For more test scenarios, see:
  // https://thegraph.com/docs/en/subgraphs/developing/creating/unit-testing-framework/#write-a-unit-test

  test("GlobalRatePeriodIntervalUpdated created and stored", () => {
    assert.entityCount("GlobalRatePeriodIntervalUpdated", 1)

    // Create a new event to get the same entity ID format
    let mockEvent = createGlobalRatePeriodIntervalUpdatedEvent(
      BigInt.fromI32(234),
      BigInt.fromI32(234)
    )
    // The entity ID is created by concatenating transaction hash with log index
    let entityId = mockEvent.transaction.hash.concatI32(mockEvent.logIndex.toI32()).toHexString()
    
    assert.fieldEquals(
      "GlobalRatePeriodIntervalUpdated",
      entityId,
      "previousValue",
      "234"
    )
    assert.fieldEquals(
      "GlobalRatePeriodIntervalUpdated",
      entityId,
      "newValue",
      "234"
    )
  })
})

describe("ActiveReservation entity", () => {
  afterAll(() => {
    clearStore()
  })

  test("ActiveReservation created and updated on ReservationUpdated event", () => {
    // Create test data
    let account = Address.fromString("0x1234567890123456789012345678901234567890")
    let symbolsPerSecond = BigInt.fromI32(1000)
    let startTimestamp = BigInt.fromI32(1000000)
    let endTimestamp = BigInt.fromI32(2000000)
    let quorumNumbers = Bytes.fromHexString("0x01")
    let quorumSplits = Bytes.fromHexString("0x64")

    // Create and handle first reservation event
    let event1 = createReservationUpdatedEvent(
      account,
      symbolsPerSecond,
      startTimestamp,
      endTimestamp,
      quorumNumbers,
      quorumSplits
    )
    handleReservationUpdated(event1)

    // Check that ActiveReservation was created
    assert.entityCount("ActiveReservation", 1)
    
    // Verify the ActiveReservation fields
    let accountId = account.toHexString()
    assert.fieldEquals("ActiveReservation", accountId, "account", accountId)
    assert.fieldEquals("ActiveReservation", accountId, "symbolsPerSecond", "1000")
    assert.fieldEquals("ActiveReservation", accountId, "startTimestamp", "1000000")
    assert.fieldEquals("ActiveReservation", accountId, "endTimestamp", "2000000")
    assert.fieldEquals("ActiveReservation", accountId, "quorumNumbers", "0x01")
    assert.fieldEquals("ActiveReservation", accountId, "quorumSplits", "0x64")

    // Create and handle updated reservation event
    let newSymbolsPerSecond = BigInt.fromI32(2000)
    let newEndTimestamp = BigInt.fromI32(3000000)
    
    let event2 = createReservationUpdatedEvent(
      account,
      newSymbolsPerSecond,
      startTimestamp,
      newEndTimestamp,
      quorumNumbers,
      quorumSplits
    )
    handleReservationUpdated(event2)

    // Check that we still have only one ActiveReservation (it was updated, not created new)
    assert.entityCount("ActiveReservation", 1)
    
    // Verify the updated fields
    assert.fieldEquals("ActiveReservation", accountId, "symbolsPerSecond", "2000")
    assert.fieldEquals("ActiveReservation", accountId, "endTimestamp", "3000000")
  })

  test("Multiple accounts have separate ActiveReservations", () => {
    clearStore()
    
    let account1 = Address.fromString("0x1111111111111111111111111111111111111111")
    let account2 = Address.fromString("0x2222222222222222222222222222222222222222")
    
    // Create reservation for account1
    let event1 = createReservationUpdatedEvent(
      account1,
      BigInt.fromI32(1000),
      BigInt.fromI32(1000000),
      BigInt.fromI32(2000000),
      Bytes.fromHexString("0x01"),
      Bytes.fromHexString("0x64")
    )
    handleReservationUpdated(event1)
    
    // Create reservation for account2
    let event2 = createReservationUpdatedEvent(
      account2,
      BigInt.fromI32(2000),
      BigInt.fromI32(1500000),
      BigInt.fromI32(2500000),
      Bytes.fromHexString("0x02"),
      Bytes.fromHexString("0x32")
    )
    handleReservationUpdated(event2)
    
    // Check that we have two ActiveReservations
    assert.entityCount("ActiveReservation", 2)
    
    // Verify each account has its own reservation
    assert.fieldEquals("ActiveReservation", account1.toHexString(), "symbolsPerSecond", "1000")
    assert.fieldEquals("ActiveReservation", account2.toHexString(), "symbolsPerSecond", "2000")
  })
})
