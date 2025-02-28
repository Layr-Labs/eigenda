import { newMockEvent, newMockCall } from "matchstick-as"
import { ethereum, BigInt, Bytes, Address } from "@graphprotocol/graph-ts"
import { NewPubkeyRegistration as NewPubkeyRegistrationEvent, NewPubkeyRegistrationPubkeyG1Struct, NewPubkeyRegistrationPubkeyG2Struct } from "../generated/BLSApkRegistry_Operator/BLSApkRegistry"
import { OperatorRegistered as OperatorRegisteredEvent, OperatorDeregistered as OperatorDeregisteredEvent } from "../generated/RegistryCoordinator_Operator/RegistryCoordinator"
import { OperatorSocketUpdate as OperatorSocketUpdateEvent } from "../generated/RegistryCoordinator/RegistryCoordinator"
import { OperatorEjected } from "../generated/EjectionManager/EjectionManager" 

export function createNewPubkeyRegistrationEvent(
  operator: Address,
  pubkeyG1_X: BigInt,
  pubkeyG1_Y: BigInt,
  pubkeyG2_X: Array<BigInt>,
  pubkeyG2_Y: Array<BigInt>
): NewPubkeyRegistrationEvent {
  let newPubkeyRegistrationEvent = changetype<
    NewPubkeyRegistrationEvent
  >(newMockEvent())

  let g1Pubkey = new NewPubkeyRegistrationPubkeyG1Struct(2)
  g1Pubkey[0] = ethereum.Value.fromUnsignedBigInt(pubkeyG1_X)
  g1Pubkey[1] = ethereum.Value.fromUnsignedBigInt(pubkeyG1_Y)

  let g2Pubkey = new NewPubkeyRegistrationPubkeyG2Struct(2)
  g2Pubkey[0] = ethereum.Value.fromUnsignedBigIntArray(pubkeyG2_X)
  g2Pubkey[1] = ethereum.Value.fromUnsignedBigIntArray(pubkeyG2_Y)

  newPubkeyRegistrationEvent.parameters = new Array()

  newPubkeyRegistrationEvent.parameters.push(
    new ethereum.EventParam("operator", ethereum.Value.fromAddress(operator))
  )
  newPubkeyRegistrationEvent.parameters.push(
    new ethereum.EventParam(
      "pubkeyG1",
      ethereum.Value.fromTuple(g1Pubkey)
    )
  )
  newPubkeyRegistrationEvent.parameters.push(
    new ethereum.EventParam(
      "pubkeyG2",
      ethereum.Value.fromTuple(g2Pubkey)
    )
  )

  return newPubkeyRegistrationEvent
}

export function createNewOperatorSocketUpdateEvent(
  operatorId: Bytes,
  socket: string
): OperatorSocketUpdateEvent {
  let newOperatorSocketUpdateEvent = changetype<
    OperatorSocketUpdateEvent
  >(newMockEvent())

  newOperatorSocketUpdateEvent.parameters = new Array()

  newOperatorSocketUpdateEvent.parameters.push(
    new ethereum.EventParam("operatorId", ethereum.Value.fromFixedBytes(operatorId))
  )

  newOperatorSocketUpdateEvent.parameters.push(
    new ethereum.EventParam("socket", ethereum.Value.fromString(socket))
  )

  return newOperatorSocketUpdateEvent
}

export function createNewOperatorRegisteredEvent(
  operator: Address,
  operatorId: Bytes
): OperatorRegisteredEvent {
  let newOperatorRegisteredEvent = changetype<
    OperatorRegisteredEvent
  >(newMockEvent())

  newOperatorRegisteredEvent.parameters = new Array()

  newOperatorRegisteredEvent.parameters.push(
    new ethereum.EventParam("operator", ethereum.Value.fromAddress(operator))
  )

  newOperatorRegisteredEvent.parameters.push(
    new ethereum.EventParam("operatorId", ethereum.Value.fromFixedBytes(operatorId))
  )

  return newOperatorRegisteredEvent
}

export function createNewOperatorDeregisteredEvent(
  operator: Address,
  operatorId: Bytes
): OperatorDeregisteredEvent {
  let newOperatorDeregisteredEvent = changetype<
    OperatorDeregisteredEvent
  >(newMockEvent())

  newOperatorDeregisteredEvent.parameters = new Array()

  newOperatorDeregisteredEvent.parameters.push(
    new ethereum.EventParam("operator", ethereum.Value.fromAddress(operator))
  )

  newOperatorDeregisteredEvent.parameters.push(
    new ethereum.EventParam("operatorId", ethereum.Value.fromFixedBytes(operatorId))
  )

  return newOperatorDeregisteredEvent
}

export function createNewOperatorEjectedEvent(operator: Address): OperatorEjected {
  let newOperatorEjectedEvent = changetype<OperatorEjected>(newMockCall())
  newOperatorEjectedEvent.parameters = new Array()

  let operatorParam = new ethereum.EventParam("operator", ethereum.Value.fromAddress(operator))
  newOperatorEjectedEvent.parameters.push(operatorParam)

  return newOperatorEjectedEvent
}