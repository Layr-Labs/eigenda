import { Address, BigInt, Bytes, crypto, ethereum, log } from "@graphprotocol/graph-ts"
import {
  BatchConfirmed as BatchConfirmedEvent,
  ConfirmBatchCall
} from "../generated/EigenDAServiceManager/EigenDAServiceManager"

import {
  Batch, BatchHeader, GasFees, Operator, NonSigning
} from "../generated/schema"

export const BATCH_HEADER_PREFIX_BYTES = Bytes.fromHexString("0x0001")
export const NON_SIGNING_PREFIX_BYTES = Bytes.fromHexString("0x0002")
export const OPERATOR_PREFIX_BYTES = Bytes.fromHexString("0x0003")
export const G1_POINT_PREFIX_BYTES = Bytes.fromHexString("0x0004")
export const G2_POINT_PREFIX_BYTES = Bytes.fromHexString("0x0005")
export const BATCH_GAS_FEES_PREFIX_BYTES = Bytes.fromHexString("0x0006")
export const BATCH_PREFIX_BYTES = Bytes.fromHexString("0x0007")

export function handleConfirmBatchCall(confirmBatchCall: ConfirmBatchCall): void {
  let batchHeader = new BatchHeader(BATCH_HEADER_PREFIX_BYTES.concat(confirmBatchCall.transaction.hash))
  batchHeader.blobHeadersRoot = confirmBatchCall.inputs.batchHeader.blobHeadersRoot
  batchHeader.blobHeadersRoot = confirmBatchCall.inputs.batchHeader.blobHeadersRoot
  batchHeader.quorumNumbers = bytesToBigIntArray(confirmBatchCall.inputs.batchHeader.quorumNumbers)
  batchHeader.signedStakeForQuorums = bytesToBigIntArray(confirmBatchCall.inputs.batchHeader.signedStakeForQuorums)
  batchHeader.referenceBlockNumber = confirmBatchCall.inputs.batchHeader.referenceBlockNumber
  batchHeader.batch = BATCH_PREFIX_BYTES.concat(confirmBatchCall.transaction.hash) // only one batch per tx
  batchHeader.save()

  let nonSignerStakesAndSignatures = new NonSigning(NON_SIGNING_PREFIX_BYTES.concat(confirmBatchCall.transaction.hash))
  
  // create the nonSigners
  let nonSigners: Bytes[] = [] 
  for (let index = 0; index < confirmBatchCall.inputs.nonSignerStakesAndSignature.nonSignerPubkeys.length; index++) {
    const pubkey = confirmBatchCall.inputs.nonSignerStakesAndSignature.nonSignerPubkeys[index];
    let operatorId = hash2BigInts(pubkey.X, pubkey.Y) // note: this is the operatorId in the contracts
    let operatorEntityId = OPERATOR_PREFIX_BYTES.concat(operatorId)
    let operator = Operator.load(operatorEntityId)
    if (operator == null) {
      operator = new Operator(operatorEntityId)
      operator.operatorId = operatorId
      operator.save()
    }
    // add the operator to the nonSigners list
    nonSigners.push(operatorEntityId)
  } 
  // link the nonSigners to the nonSignerStakesAndSignatures
  nonSignerStakesAndSignatures.nonSigners = nonSigners
  
  nonSignerStakesAndSignatures.batch = BATCH_PREFIX_BYTES.concat(confirmBatchCall.transaction.hash) // only one batch per tx
  nonSignerStakesAndSignatures.save()
}

export function handleBatchConfirmed(batchConfirmedEvent: BatchConfirmedEvent): void {
  if (batchConfirmedEvent.receipt == null) {
    log.error("handleBatchConfirmed: batchConfirmedEvent.receipt is null", [batchConfirmedEvent.transaction.hash.toHex()])
    return
  }
  let batchGasFees = new GasFees(BATCH_GAS_FEES_PREFIX_BYTES.concat(batchConfirmedEvent.transaction.hash)) // only one batch per tx
  batchGasFees.gasPrice = batchConfirmedEvent.transaction.gasPrice
  batchGasFees.gasUsed = batchConfirmedEvent.receipt!.gasUsed
  batchGasFees.txFee = batchGasFees.gasPrice.times(batchGasFees.gasUsed)
  batchGasFees.save()

  let batch = new Batch(BATCH_PREFIX_BYTES.concat(batchConfirmedEvent.transaction.hash)) // only one batch per tx
  batch.batchId = batchConfirmedEvent.params.batchId
  batch.batchHeaderHash = batchConfirmedEvent.params.batchHeaderHash
  batch.gasFees = batchGasFees.id
  batch.blockNumber = batchConfirmedEvent.block.number
  batch.blockTimestamp = batchConfirmedEvent.block.timestamp
  batch.txHash = batchConfirmedEvent.transaction.hash
  batch.save()
}

export function bytesToBigIntArray(bytes: Bytes): BigInt[] {
  let hex = bytes.toHex().substring(2);
  let result: BigInt[] = [];
  for (let i = 0; i < hex.length / 2; i++) {
    let byte = hex.substring(i * 2, (i+1) * 2 );
    let hexByteValue = Bytes.fromHexString(byte)
    let bigIntByte = BigInt.fromUnsignedBytes(hexByteValue)
    result.push(bigIntByte);
  }
  return result;
}

export function hash2BigInts(x: BigInt, y: BigInt): Bytes {
  // pad to 32 bytes
  let xBytes = x.toHex().substring(2).padStart(64, "0")
  let yBytes = y.toHex().substring(2).padStart(64, "0")
  let xy = Bytes.fromHexString(xBytes.concat(yBytes))
  return Bytes.fromByteArray(crypto.keccak256(xy))
}