import {
  assert,
  describe,
  test,
  clearStore,
  beforeAll,
  afterAll
} from "matchstick-as/assembly/index"
import { BigInt, Address, Bytes } from "@graphprotocol/graph-ts"
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

describe("CurrentReservation entity", () => {
  afterAll(() => {
    clearStore()
  })

  test("CurrentReservation created and updated on ReservationUpdated event", () => {
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

    // Check that CurrentReservation was created
    assert.entityCount("CurrentReservation", 1)
    
    // Verify the CurrentReservation fields
    let accountId = account.toHexString()
    assert.fieldEquals("CurrentReservation", accountId, "account", accountId)
    assert.fieldEquals("CurrentReservation", accountId, "symbolsPerSecond", "1000")
    assert.fieldEquals("CurrentReservation", accountId, "startTimestamp", "1000000")
    assert.fieldEquals("CurrentReservation", accountId, "endTimestamp", "2000000")
    assert.fieldEquals("CurrentReservation", accountId, "quorumNumbers", "0x01")
    assert.fieldEquals("CurrentReservation", accountId, "quorumSplits", "0x64")

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

    // Check that we still have only one CurrentReservation (it was updated, not created new)
    assert.entityCount("CurrentReservation", 1)
    
    // Verify the updated fields
    assert.fieldEquals("CurrentReservation", accountId, "symbolsPerSecond", "2000")
    assert.fieldEquals("CurrentReservation", accountId, "endTimestamp", "3000000")
  })

  test("Multiple accounts have separate CurrentReservations with different time ranges", () => {
    clearStore()
    
    // Create three accounts with different reservation time ranges
    let accounts = [
      Address.fromString("0x1111111111111111111111111111111111111111"), // Past (expired)
      Address.fromString("0x2222222222222222222222222222222222222222"), // Current (would be active)
      Address.fromString("0x3333333333333333333333333333333333333333")  // Future (would be pending)
    ]
    
    // Past reservation (expired) - ended at timestamp 200000
    let event1 = createReservationUpdatedEvent(
      accounts[0],
      BigInt.fromI32(1000),
      BigInt.fromI32(100000),
      BigInt.fromI32(200000),
      Bytes.fromHexString("0x01"),
      Bytes.fromHexString("0x64")
    )
    handleReservationUpdated(event1)
    
    // Current reservation (active) - from 150000 to 250000
    let event2 = createReservationUpdatedEvent(
      accounts[1],
      BigInt.fromI32(2000),
      BigInt.fromI32(150000),
      BigInt.fromI32(250000),
      Bytes.fromHexString("0x02"),
      Bytes.fromHexString("0x32")
    )
    handleReservationUpdated(event2)
    
    // Future reservation (pending) - starts at 300000
    let event3 = createReservationUpdatedEvent(
      accounts[2],
      BigInt.fromI32(3000),
      BigInt.fromI32(300000),
      BigInt.fromI32(400000),
      Bytes.fromHexString("0x03"),
      Bytes.fromHexString("0x50")
    )
    handleReservationUpdated(event3)
    
    // Verify we have three CurrentReservations
    assert.entityCount("CurrentReservation", 3)
    
    // Verify each account has its own reservation with correct data
    assert.fieldEquals("CurrentReservation", accounts[0].toHexString(), "symbolsPerSecond", "1000")
    assert.fieldEquals("CurrentReservation", accounts[0].toHexString(), "startTimestamp", "100000")
    assert.fieldEquals("CurrentReservation", accounts[0].toHexString(), "endTimestamp", "200000")
    
    assert.fieldEquals("CurrentReservation", accounts[1].toHexString(), "symbolsPerSecond", "2000")
    assert.fieldEquals("CurrentReservation", accounts[1].toHexString(), "startTimestamp", "150000")
    assert.fieldEquals("CurrentReservation", accounts[1].toHexString(), "endTimestamp", "250000")
    
    assert.fieldEquals("CurrentReservation", accounts[2].toHexString(), "symbolsPerSecond", "3000")
    assert.fieldEquals("CurrentReservation", accounts[2].toHexString(), "startTimestamp", "300000")
    assert.fieldEquals("CurrentReservation", accounts[2].toHexString(), "endTimestamp", "400000")
  })
})
