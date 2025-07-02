import {
  assert,
  describe,
  test,
  clearStore,
  beforeAll,
  afterAll
} from "matchstick-as/assembly/index"
import { BigInt, Address } from "@graphprotocol/graph-ts"
import { GlobalRatePeriodIntervalUpdated } from "../generated/schema"
import { GlobalRatePeriodIntervalUpdated as GlobalRatePeriodIntervalUpdatedEvent } from "../generated/PaymentVault/PaymentVault"
import { handleGlobalRatePeriodIntervalUpdated } from "../src/payment-vault"
import { createGlobalRatePeriodIntervalUpdatedEvent } from "./payment-vault-utils"

// Tests structure (matchstick-as >=0.5.0)
// https://thegraph.com/docs/en/developer/matchstick/#tests-structure-0-5-0

describe("Describe entity assertions", () => {
  beforeAll(() => {
    let previousValue = BigInt.fromI32(234)
    let newValue = BigInt.fromI32(234)
    let newGlobalRatePeriodIntervalUpdatedEvent = createGlobalRatePeriodIntervalUpdatedEvent(
      previousValue,
      newValue
    )
    handleGlobalRatePeriodIntervalUpdated(
      newGlobalRatePeriodIntervalUpdatedEvent
    )
  })

  afterAll(() => {
    clearStore()
  })

  // For more test scenarios, see:
  // https://thegraph.com/docs/en/developer/matchstick/#write-a-unit-test

  test("GlobalRatePeriodIntervalUpdated created and stored", () => {
    assert.entityCount("GlobalRatePeriodIntervalUpdated", 1)

    // 0xa16081f360e3847006db660bae1c6d1b2e17ec2a is the default address used in newMockEvent() function
    assert.fieldEquals(
      "GlobalRatePeriodIntervalUpdated",
      "0xa16081f360e3847006db660bae1c6d1b2e17ec2a-1",
      "previousValue",
      "234"
    )
    assert.fieldEquals(
      "GlobalRatePeriodIntervalUpdated",
      "0xa16081f360e3847006db660bae1c6d1b2e17ec2a-1",
      "newValue",
      "234"
    )

    // More assert options:
    // https://thegraph.com/docs/en/developer/matchstick/#asserts
  })
})
