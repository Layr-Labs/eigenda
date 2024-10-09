import {
  EjectorUpdated as EjectorUpdatedEvent,
  Initialized as InitializedEvent,
  OperatorEjected as OperatorEjectedEvent,
  OwnershipTransferred as OwnershipTransferredEvent,
  QuorumEjection as QuorumEjectionEvent,
  QuorumEjectionParamsSet as QuorumEjectionParamsSetEvent
} from "../generated/EjectionManager/EjectionManager"
import {
  EjectorUpdated,
  Initialized,
  OperatorEjected,
  OwnershipTransferred,
  QuorumEjection,
  QuorumEjectionParamsSet
} from "../generated/schema"

export function handleEjectorUpdated(event: EjectorUpdatedEvent): void {
  let entity = new EjectorUpdated(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.ejector = event.params.ejector
  entity.status = event.params.status

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleInitialized(event: InitializedEvent): void {
  let entity = new Initialized(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.version = event.params.version

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleOperatorEjected(event: OperatorEjectedEvent): void {
  let entity = new OperatorEjected(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.operatorId = event.params.operatorId
  entity.quorumNumber = event.params.quorumNumber

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleOwnershipTransferred(
  event: OwnershipTransferredEvent
): void {
  let entity = new OwnershipTransferred(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.previousOwner = event.params.previousOwner
  entity.newOwner = event.params.newOwner

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleQuorumEjection(event: QuorumEjectionEvent): void {
  let entity = new QuorumEjection(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.ejectedOperators = event.params.ejectedOperators
  entity.ratelimitHit = event.params.ratelimitHit

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleQuorumEjectionParamsSet(
  event: QuorumEjectionParamsSetEvent
): void {
  let entity = new QuorumEjectionParamsSet(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.quorumNumber = event.params.quorumNumber
  entity.rateLimitWindow = event.params.rateLimitWindow
  entity.ejectableStakePercent = event.params.ejectableStakePercent

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}
