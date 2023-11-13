import {
  OperatorAddedToQuorums as OperatorAddedToQuorumsEvent,
  OperatorRemovedFromQuorums as OperatorRemovedFromQuorumsEvent
} from "../generated/BLSPubkeyRegistry/BLSPubkeyRegistry"
import {
  OperatorAddedToQuorum,
  OperatorRemovedFromQuorum
} from "../generated/schema"

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
