import {
  assert,
  describe,
  test,
  clearStore,
  beforeAll,
  afterAll
} from "matchstick-as/assembly/index"
import { Address, Bytes, BigInt } from "@graphprotocol/graph-ts"
import { EjectorUpdated } from "../generated/schema"
import { EjectorUpdated as EjectorUpdatedEvent } from "../generated/EjectionManager/EjectionManager"
import { handleEjectorUpdated } from "../src/ejection-manager"
import { createEjectorUpdatedEvent } from "./ejection-manager-utils"

// Tests structure (matchstick-as >=0.5.0)
// https://thegraph.com/docs/en/developer/matchstick/#tests-structure-0-5-0

describe("Describe entity assertions", () => {
  beforeAll(() => {
    let ejector = Address.fromString(
      "0x0000000000000000000000000000000000000001"
    )
    let status = "boolean Not implemented"
    let newEjectorUpdatedEvent = createEjectorUpdatedEvent(ejector, status)
    handleEjectorUpdated(newEjectorUpdatedEvent)
  })

  afterAll(() => {
    clearStore()
  })

  // For more test scenarios, see:
  // https://thegraph.com/docs/en/developer/matchstick/#write-a-unit-test

  test("EjectorUpdated created and stored", () => {
    assert.entityCount("EjectorUpdated", 1)

    // 0xa16081f360e3847006db660bae1c6d1b2e17ec2a is the default address used in newMockEvent() function
    assert.fieldEquals(
      "EjectorUpdated",
      "0xa16081f360e3847006db660bae1c6d1b2e17ec2a-1",
      "ejector",
      "0x0000000000000000000000000000000000000001"
    )
    assert.fieldEquals(
      "EjectorUpdated",
      "0xa16081f360e3847006db660bae1c6d1b2e17ec2a-1",
      "status",
      "boolean Not implemented"
    )

    // More assert options:
    // https://thegraph.com/docs/en/developer/matchstick/#asserts
  })
})
