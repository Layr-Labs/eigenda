import { newMockEvent } from "matchstick-as"
import { ethereum, Address, Bytes, BigInt } from "@graphprotocol/graph-ts"
import {
  EjectorUpdated,
  Initialized,
  OperatorEjected,
  OwnershipTransferred,
  QuorumEjection,
  QuorumEjectionParamsSet
} from "../generated/EjectionManager/EjectionManager"

export function createEjectorUpdatedEvent(
  ejector: Address,
  status: boolean
): EjectorUpdated {
  let ejectorUpdatedEvent = changetype<EjectorUpdated>(newMockEvent())

  ejectorUpdatedEvent.parameters = new Array()

  ejectorUpdatedEvent.parameters.push(
    new ethereum.EventParam("ejector", ethereum.Value.fromAddress(ejector))
  )
  ejectorUpdatedEvent.parameters.push(
    new ethereum.EventParam("status", ethereum.Value.fromBoolean(status))
  )

  return ejectorUpdatedEvent
}

export function createInitializedEvent(version: i32): Initialized {
  let initializedEvent = changetype<Initialized>(newMockEvent())

  initializedEvent.parameters = new Array()

  initializedEvent.parameters.push(
    new ethereum.EventParam(
      "version",
      ethereum.Value.fromUnsignedBigInt(BigInt.fromI32(version))
    )
  )

  return initializedEvent
}

export function createOperatorEjectedEvent(
  operatorId: Bytes,
  quorumNumber: i32
): OperatorEjected {
  let operatorEjectedEvent = changetype<OperatorEjected>(newMockEvent())

  operatorEjectedEvent.parameters = new Array()

  operatorEjectedEvent.parameters.push(
    new ethereum.EventParam(
      "operatorId",
      ethereum.Value.fromFixedBytes(operatorId)
    )
  )
  operatorEjectedEvent.parameters.push(
    new ethereum.EventParam(
      "quorumNumber",
      ethereum.Value.fromUnsignedBigInt(BigInt.fromI32(quorumNumber))
    )
  )

  return operatorEjectedEvent
}

export function createOwnershipTransferredEvent(
  previousOwner: Address,
  newOwner: Address
): OwnershipTransferred {
  let ownershipTransferredEvent = changetype<OwnershipTransferred>(
    newMockEvent()
  )

  ownershipTransferredEvent.parameters = new Array()

  ownershipTransferredEvent.parameters.push(
    new ethereum.EventParam(
      "previousOwner",
      ethereum.Value.fromAddress(previousOwner)
    )
  )
  ownershipTransferredEvent.parameters.push(
    new ethereum.EventParam("newOwner", ethereum.Value.fromAddress(newOwner))
  )

  return ownershipTransferredEvent
}

export function createQuorumEjectionEvent(
  ejectedOperators: BigInt,
  ratelimitHit: boolean
): QuorumEjection {
  let quorumEjectionEvent = changetype<QuorumEjection>(newMockEvent())

  quorumEjectionEvent.parameters = new Array()

  quorumEjectionEvent.parameters.push(
    new ethereum.EventParam(
      "ejectedOperators",
      ethereum.Value.fromUnsignedBigInt(ejectedOperators)
    )
  )
  quorumEjectionEvent.parameters.push(
    new ethereum.EventParam(
      "ratelimitHit",
      ethereum.Value.fromBoolean(ratelimitHit)
    )
  )

  return quorumEjectionEvent
}

export function createQuorumEjectionParamsSetEvent(
  quorumNumber: i32,
  rateLimitWindow: BigInt,
  ejectableStakePercent: i32
): QuorumEjectionParamsSet {
  let quorumEjectionParamsSetEvent = changetype<QuorumEjectionParamsSet>(
    newMockEvent()
  )

  quorumEjectionParamsSetEvent.parameters = new Array()

  quorumEjectionParamsSetEvent.parameters.push(
    new ethereum.EventParam(
      "quorumNumber",
      ethereum.Value.fromUnsignedBigInt(BigInt.fromI32(quorumNumber))
    )
  )
  quorumEjectionParamsSetEvent.parameters.push(
    new ethereum.EventParam(
      "rateLimitWindow",
      ethereum.Value.fromUnsignedBigInt(rateLimitWindow)
    )
  )
  quorumEjectionParamsSetEvent.parameters.push(
    new ethereum.EventParam(
      "ejectableStakePercent",
      ethereum.Value.fromUnsignedBigInt(BigInt.fromI32(ejectableStakePercent))
    )
  )

  return quorumEjectionParamsSetEvent
}
