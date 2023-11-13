import {
    assert,
    describe,
    test,
    clearStore,
    beforeAll,
    afterAll,
    newMockCall,
    createMockedFunction
  } from "matchstick-as"
import { Address, BigInt, Bytes, ethereum, log } from "@graphprotocol/graph-ts"
import { handleBatchConfirmed, BATCH_PREFIX_BYTES, BATCH_GAS_FEES_PREFIX_BYTES, handleConfirmBatchCall, BATCH_HEADER_PREFIX_BYTES, NON_SIGNING_PREFIX_BYTES, OPERATOR_PREFIX_BYTES, hash2BigInts, bytesToBigIntArray } from "../src/edasm"
import { createNewBatchConfirmedEvent, createNewConfimBatchCall } from "./edasm-utils"

let blobHeadersRoot: Bytes = Bytes.fromHexString("0x1111000011110000111100001111000011110000111100001111000011110000")
let quorumNumbers: Bytes = Bytes.fromHexString("0x000112")
let quorumThresholdPercentages: Bytes = Bytes.fromHexString("0x646464")
let referenceBlockNumber: BigInt = BigInt.fromI32(123123)
let nonSignerPubkeysBigInts: Array<Array<BigInt>> = [
  [BigInt.fromI32(123), BigInt.fromI32(456)],
  [BigInt.fromI32(789), BigInt.fromI32(234)]
]

// 64 bytes
let batchHeaderHash: Bytes = Bytes.fromHexString("0x1234567890123456789012345678901234567890123456789012345678901234")
let batchId: BigInt = BigInt.fromI32(123)
let fee: BigInt = BigInt.fromI32(123123123)

describe("EigenDASM", () => {
  beforeAll(() => {

  })

  afterAll(() => {
    clearStore()
  })

  // For more test scenarios, see:
  // https://thegraph.com/docs/en/developer/matchstick/#write-a-unit-test

  test("has batchheader, nonsigners, and operators created", () => {
    let confirmBatchCall = createNewConfimBatchCall(
      blobHeadersRoot,
      quorumNumbers,
      quorumThresholdPercentages,
      referenceBlockNumber,
      nonSignerPubkeysBigInts
    )
    
    handleConfirmBatchCall(confirmBatchCall)

    assert.entityCount("BatchHeader", 1)
    let batchHeaderEntityId = BATCH_HEADER_PREFIX_BYTES.concat(confirmBatchCall.transaction.hash)

    assert.entityCount("NonSigning", 1)
    let nonSigningEntityId = NON_SIGNING_PREFIX_BYTES.concat(confirmBatchCall.transaction.hash)

    assert.entityCount("Operator", 2)
    let operatorId1 = hash2BigInts(nonSignerPubkeysBigInts[0][0], nonSignerPubkeysBigInts[0][1])
    let operatorEntityId1 = OPERATOR_PREFIX_BYTES.concat(operatorId1)

    let operatorId2 = hash2BigInts(nonSignerPubkeysBigInts[1][0], nonSignerPubkeysBigInts[1][1])
    let operatorEntityId2 = OPERATOR_PREFIX_BYTES.concat(operatorId2)

    assert.fieldEquals(
      "BatchHeader",
      batchHeaderEntityId.toHexString(),
      "blobHeadersRoot",
      blobHeadersRoot.toHexString()
    )

    assert.fieldEquals(
      "BatchHeader",
      batchHeaderEntityId.toHexString(),
      "quorumNumbers",
      convertArraySringToAssertString(bytesToBigIntArray(quorumNumbers).toString())
    )

    assert.fieldEquals(
      "BatchHeader",
      batchHeaderEntityId.toHexString(),
      "quorumThresholdPercentages",
      convertArraySringToAssertString(bytesToBigIntArray(quorumThresholdPercentages).toString())
    )

    assert.fieldEquals(
      "BatchHeader",
      batchHeaderEntityId.toHexString(),
      "referenceBlockNumber",
      referenceBlockNumber.toString()
    )

    assert.fieldEquals(
      "BatchHeader",
      batchHeaderEntityId.toHexString(),
      "batch",
      BATCH_PREFIX_BYTES.concat(confirmBatchCall.transaction.hash).toHexString()
    )

    assert.fieldEquals(
      "NonSigning",
      nonSigningEntityId.toHexString(),
      "batch",
      confirmBatchCall.transaction.hash.toHexString()
    )

    assert.fieldEquals(
      "NonSigning",
      nonSigningEntityId.toHexString(),
      "nonSigners",
      convertArraySringToAssertString([operatorEntityId1.toHexString(), operatorEntityId2.toHexString()].toString())
    )

    assert.fieldEquals(
      "Operator",
      operatorEntityId1.toHexString(),
      "operatorId",
      operatorId1.toHexString()
    )

  })

  test("has batch and gas fees created", () => {
    let batchConfirmedEvent = createNewBatchConfirmedEvent(
      batchHeaderHash,
      batchId,
      fee
    )

    handleBatchConfirmed(batchConfirmedEvent)

    assert.entityCount("Batch", 1)
    let batchEntityId = BATCH_PREFIX_BYTES.concat(batchConfirmedEvent.transaction.hash)
    assert.entityCount("GasFees", 1)
    let gasFeesEntityId = BATCH_GAS_FEES_PREFIX_BYTES.concat(batchConfirmedEvent.transaction.hash)

    assert.fieldEquals(
      "Batch",
      batchEntityId.toHexString(),
      "batchId",
      batchId.toString()
    )

    assert.fieldEquals(
      "Batch",
      batchEntityId.toHexString(),
      "batchHeaderHash",
      batchHeaderHash.toHexString()
    )

    assert.fieldEquals(
      "Batch",
      batchEntityId.toHexString(),
      "gasFees",
      gasFeesEntityId.toHexString()
    )

    assert.fieldEquals(
      "Batch",
      batchEntityId.toHexString(),
      "blockNumber",
      batchConfirmedEvent.block.number.toString()
    )

    assert.fieldEquals(
      "Batch",
      batchEntityId.toHexString(),
      "blockTimestamp",
      batchConfirmedEvent.block.timestamp.toString()
    )

    // type GasFees @entity(immutable: true) {
    //   id: Bytes!
    //   gasUsed: BigInt!
    //   gasPrice: BigInt!
    //   txFee: BigInt!
    // }

    assert.fieldEquals(
      "GasFees",
      gasFeesEntityId.toHexString(),
      "gasUsed",
      batchConfirmedEvent.receipt!.gasUsed.toString()
    )

    assert.fieldEquals(
      "GasFees",
      gasFeesEntityId.toHexString(),
      "gasPrice",
      batchConfirmedEvent.transaction.gasPrice.toString()
    )

    assert.fieldEquals(
      "GasFees",
      gasFeesEntityId.toHexString(),
      "txFee",
      batchConfirmedEvent.transaction.gasPrice.times(batchConfirmedEvent.receipt!.gasUsed).toString()
    )
  })


})

function convertArraySringToAssertString(arrString: string): string {
  return "["  + arrString.split(",").join(", ") + "]"
}