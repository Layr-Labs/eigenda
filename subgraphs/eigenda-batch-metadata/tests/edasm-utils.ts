import { newMockEvent, newMockCall } from "matchstick-as"
import { ethereum, BigInt, Bytes, Address } from "@graphprotocol/graph-ts"
import { BatchConfirmed as BatchConfirmedEvent, ConfirmBatchCall, ConfirmBatchCallBatchHeaderStruct, ConfirmBatchCallNonSignerStakesAndSignatureNonSignerPubkeysStruct, ConfirmBatchCallNonSignerStakesAndSignatureStruct } from "../generated/EigenDAServiceManager/EigenDAServiceManager"
import { BatchHeader } from "../generated/schema"

export function createNewConfimBatchCall(
  blobHeadersRoot: Bytes,
  quorumNumbers: Bytes,
  quorumThresholdPercentages: Bytes,
  referenceBlockNumber: BigInt,
  nonSignerPubkeysBigInts: Array<Array<BigInt>>,
): ConfirmBatchCall {
  let confirmBatchCall = changetype<
    ConfirmBatchCall
  >(newMockCall())

  let batchHeader = new ConfirmBatchCallBatchHeaderStruct(4)
  batchHeader[0] = ethereum.Value.fromBytes(blobHeadersRoot)
  batchHeader[1] = ethereum.Value.fromBytes(quorumNumbers)
  batchHeader[2] = ethereum.Value.fromBytes(quorumThresholdPercentages)
  batchHeader[3] = ethereum.Value.fromUnsignedBigInt(referenceBlockNumber)

  let nonSignerPubkeys: ethereum.Tuple[] = []
  for (let index = 0; index < nonSignerPubkeysBigInts.length; index++) {
    const pubkey = nonSignerPubkeysBigInts[index];
    let nonSignerPubkey = new ConfirmBatchCallNonSignerStakesAndSignatureNonSignerPubkeysStruct(2)
    nonSignerPubkey[0] = ethereum.Value.fromUnsignedBigInt(pubkey[0])
    nonSignerPubkey[1] = ethereum.Value.fromUnsignedBigInt(pubkey[1])
    nonSignerPubkeys.push(nonSignerPubkey)
  }
  
  let emptyTuple = new ethereum.Tuple(0)
  let nonSignerStakesAndSignature = new ConfirmBatchCallNonSignerStakesAndSignatureStruct(8)
  nonSignerStakesAndSignature[0] = ethereum.Value.fromUnsignedBigIntArray([]),
  nonSignerStakesAndSignature[1] = ethereum.Value.fromTupleArray(nonSignerPubkeys),
  nonSignerStakesAndSignature[2] = ethereum.Value.fromTupleArray([]),
  nonSignerStakesAndSignature[3] = ethereum.Value.fromTuple(emptyTuple),
  nonSignerStakesAndSignature[4] = ethereum.Value.fromTuple(emptyTuple),
  nonSignerStakesAndSignature[5] = ethereum.Value.fromUnsignedBigIntArray([]),
  nonSignerStakesAndSignature[6] = ethereum.Value.fromUnsignedBigIntArray([]),
  nonSignerStakesAndSignature[7] = ethereum.Value.fromUnsignedBigIntMatrix([])
  
  
  confirmBatchCall.inputValues.push(
    new ethereum.EventParam("batchHeader", ethereum.Value.fromTuple(batchHeader)),
  )

  confirmBatchCall.inputValues.push(
    new ethereum.EventParam("nonSignerStakesAndSignature", ethereum.Value.fromTuple(nonSignerStakesAndSignature))
  )

  return confirmBatchCall
}

export function createNewBatchConfirmedEvent(
  batchHeaderHash: Bytes,
  batchId: BigInt,
  fee: BigInt
): BatchConfirmedEvent {
  let batchConfirmedEvent = changetype<
    BatchConfirmedEvent
  >(newMockEvent())

  // get batchHeaderHash(): Bytes {
  //   return this._event.parameters[0].value.toBytes();
  // }

  // get batchId(): BigInt {
  //   return this._event.parameters[1].value.toBigInt();
  // }

  // get fee(): BigInt {
  //   return this._event.parameters[2].value.toBigInt();
  // }

  batchConfirmedEvent.parameters = new Array()
  batchConfirmedEvent.parameters.push(
    new ethereum.EventParam("batchHeaderHash", ethereum.Value.fromFixedBytes(batchHeaderHash))
  )
  batchConfirmedEvent.parameters.push(
    new ethereum.EventParam("batchId", ethereum.Value.fromUnsignedBigInt(batchId))
  )
  batchConfirmedEvent.parameters.push(
    new ethereum.EventParam("fee", ethereum.Value.fromUnsignedBigInt(fee))
  )
  
  return batchConfirmedEvent
}

