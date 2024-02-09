import {
  OperatorAddedToQuorums as OperatorAddedToQuorumsEvent,
  OperatorRemovedFromQuorums as OperatorRemovedFromQuorumsEvent
} from "../generated/BLSApkRegistry/BLSApkRegistry"
import {
  OperatorAddedToQuorum,
  OperatorRemovedFromQuorum
} from "../generated/schema"

import { NewPubkeyRegistration as NewPubkeyRegistrationEvent } from "../generated/BLSApkRegistry/BLSApkRegistry"
import { NewPubkeyRegistration } from "../generated/schema"




export function handleOperatorAddedToQuorums(
  event: OperatorAddedToQuorumsEvent
): void {
  let entity = new OperatorAddedToQuorum(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.operator = event.params.operator
  entity.quorumNumbers = event.params.quorumNumbers

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleOperatorRemovedFromQuorums(
  event: OperatorRemovedFromQuorumsEvent
): void {
  let entity = new OperatorRemovedFromQuorum(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.operator = event.params.operator
  entity.quorumNumbers = event.params.quorumNumbers

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleNewPubkeyRegistration(
  event: NewPubkeyRegistrationEvent
): void {
  let entity = new NewPubkeyRegistration(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.operator = event.params.operator
  entity.pubkeyG1_X = event.params.pubkeyG1.X
  entity.pubkeyG1_Y = event.params.pubkeyG1.Y
  entity.pubkeyG2_X = event.params.pubkeyG2.X
  entity.pubkeyG2_Y = event.params.pubkeyG2.Y

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}