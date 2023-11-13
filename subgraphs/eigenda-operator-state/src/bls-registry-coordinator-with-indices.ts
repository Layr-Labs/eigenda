import {
  ChurnApproverUpdated as ChurnApproverUpdatedEvent,
  Initialized as InitializedEvent,
  OperatorDeregistered as OperatorDeregisteredEvent,
  OperatorRegistered as OperatorRegisteredEvent,
  OperatorSetParamsUpdated as OperatorSetParamsUpdatedEvent,
  OperatorSocketUpdate as OperatorSocketUpdateEvent
} from "../generated/BLSRegistryCoordinatorWithIndices/BLSRegistryCoordinatorWithIndices"
import {
  ChurnApproverUpdated,
  OperatorDeregistered,
  OperatorRegistered,
  OperatorSetParamsUpdated,
  OperatorSocketUpdate
} from "../generated/schema"

export function handleChurnApproverUpdated(
  event: ChurnApproverUpdatedEvent
): void {
  let entity = new ChurnApproverUpdated(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.prevChurnApprover = event.params.prevChurnApprover
  entity.newChurnApprover = event.params.newChurnApprover

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleOperatorDeregistered(
  event: OperatorDeregisteredEvent
): void {
  let entity = new OperatorDeregistered(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.operator = event.params.operator
  entity.operatorId = event.params.operatorId

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleOperatorRegistered(event: OperatorRegisteredEvent): void {
  let entity = new OperatorRegistered(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.operator = event.params.operator
  entity.operatorId = event.params.operatorId

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleOperatorSetParamsUpdated(
  event: OperatorSetParamsUpdatedEvent
): void {
  let entity = new OperatorSetParamsUpdated(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.quorumNumber = event.params.quorumNumber
  entity.operatorSetParams_maxOperatorCount =
    event.params.operatorSetParams.maxOperatorCount
  entity.operatorSetParams_kickBIPsOfOperatorStake =
    event.params.operatorSetParams.kickBIPsOfOperatorStake
  entity.operatorSetParams_kickBIPsOfTotalStake =
    event.params.operatorSetParams.kickBIPsOfTotalStake

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleOperatorSocketUpdate(
  event: OperatorSocketUpdateEvent
): void {
  let entity = new OperatorSocketUpdate(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.operatorId = event.params.operatorId
  entity.socket = event.params.socket

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}