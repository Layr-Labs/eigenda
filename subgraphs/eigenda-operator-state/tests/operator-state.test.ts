import {
    assert,
    describe,
    test,
    clearStore,
    beforeAll,
    afterAll,
    createMockedFunction
  } from "matchstick-as/assembly/index"
  import { Address, BigInt, Bytes, ethereum } from "@graphprotocol/graph-ts"
  import { createNewOperatorDeregisteredEvent, createNewOperatorRegisteredEvent, createNewOperatorSocketUpdateEvent, createNewPubkeyRegistrationEvent, createNewOperatorEjectedEvent } from "./operator-state-utils"
  import { handleNewPubkeyRegistration } from "../src/operator-creation"
  import { handleOperatorDeregistered, handleOperatorRegistered } from "../src/operator-registration-status"
  import { handleOperatorSocketUpdate } from "../src/registry-coordinator"
  import { handleOperatorEjected } from "../src/ejection-manager"
  
  let operator: Address = Address.fromBytes(Bytes.fromHexString("0xa16081f360e3847006db660bae1c6d1b2e17ec2a"))
  let pubkeyG1_X = BigInt.fromI32(123)
  let pubkeyG1_Y = BigInt.fromI32(456)
  let pubkeyG2_X = [BigInt.fromI32(789), BigInt.fromI32(234)]
  let pubkeyG2_Y = [BigInt.fromI32(345), BigInt.fromI32(678)]
  
  let pubkeyHash = Bytes.fromHexString("0x1234567890123124125325832000000999900000000004106127096123760321")
  
  let socket1 = "0.0.0.0:1234"
  let socket2 = "1.1.1.1:4321"
  
  describe("Operators", () => {
    beforeAll(() => {
  
      let newPubkeyRegistrationEvent = createNewPubkeyRegistrationEvent(
        operator,
        pubkeyG1_X,
        pubkeyG1_Y,
        pubkeyG2_X,
        pubkeyG2_Y
      )
  
      // mock the call to operatorToPubkeyHash
      createMockedFunction(newPubkeyRegistrationEvent.address, 'operatorToPubkeyHash', 'operatorToPubkeyHash(address):(bytes32)')
        .withArgs([ethereum.Value.fromAddress(operator)])
        .returns([ethereum.Value.fromBytes(pubkeyHash)])
  
      handleNewPubkeyRegistration(newPubkeyRegistrationEvent)
    })
  
    afterAll(() => {
      clearStore()
    })
  
    // For more test scenarios, see:
    // https://thegraph.com/docs/en/developer/matchstick/#write-a-unit-test
  
    test("can be created and stored", () => {
      assert.entityCount("Operator", 1)
  
      assert.fieldEquals(
        "Operator",
        pubkeyHash.toHexString(),
        "operator",
        operator.toHexString()
      )
      assert.fieldEquals(
        "Operator",
        pubkeyHash.toHexString(),
        "pubkeyG1_X",
        pubkeyG1_X.toString()
      )
      assert.fieldEquals(
        "Operator",
        pubkeyHash.toHexString(),
        "pubkeyG1_Y",
        pubkeyG1_Y.toString()
      )
      assert.fieldEquals(
        "Operator",
        pubkeyHash.toHexString(),
        "deregistrationBlockNumber",
        "0"
      )
    })
  
    test("update deregistrationBlockNumber on registration/deregistration", () => {
      assert.fieldEquals(
        "Operator",
        pubkeyHash.toHexString(),
        "deregistrationBlockNumber",
        "0"
      )
  
      let operatorRegisteredEvent = createNewOperatorRegisteredEvent(
        operator,
        pubkeyHash
      )
  
      handleOperatorRegistered(operatorRegisteredEvent)
  
      assert.fieldEquals(
        "Operator",
        pubkeyHash.toHexString(),
        "deregistrationBlockNumber",
        "4294967295"
      )
  
      let operatorDeregisteredEvent = createNewOperatorDeregisteredEvent(
        operator,
        pubkeyHash
      )
  
      handleOperatorDeregistered(operatorDeregisteredEvent)
  
      assert.fieldEquals(
        "Operator",
        pubkeyHash.toHexString(),
        "deregistrationBlockNumber",
        operatorDeregisteredEvent.block.number.toString()
      )
    })

    test("have their sockets updated", () => {
      let operatorSocketUpdatedEvent = createNewOperatorSocketUpdateEvent(
        pubkeyHash,
        socket1
      )

      handleOperatorSocketUpdate(operatorSocketUpdatedEvent)

      assert.entityCount(
        "OperatorSocketUpdate",
        1
      )

      assert.fieldEquals(
        "OperatorSocketUpdate",
        operatorSocketUpdatedEvent.transaction.hash.concatI32(operatorSocketUpdatedEvent.logIndex.toI32()).toHexString(),
        "operatorId",
        pubkeyHash.toHexString()
      )

      assert.fieldEquals(
        "OperatorSocketUpdate",
        operatorSocketUpdatedEvent.transaction.hash.concatI32(operatorSocketUpdatedEvent.logIndex.toI32()).toHexString(),
        "socket",
        socket1
      )

    })

    test("operator registered", () => {
      assert.fieldEquals("Operator", operator.toHex(), "id", operator.toHex())
      assert.fieldEquals("Operator", operator.toHex(), "pubkeyG1_X", pubkeyG1_X.toString())
      assert.fieldEquals("Operator", operator.toHex(), "pubkeyG1_Y", pubkeyG1_Y.toString())
    })

    test("operator ejected", () => {
      let ejectionEvent = createNewOperatorEjectedEvent(operator)
      handleOperatorEjected(ejectionEvent)

      assert.fieldEquals("Operator", operator.toHex(), "status", "ejected")
    })
})
