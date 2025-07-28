import { BigInt, Bytes, log } from "@graphprotocol/graph-ts"
import {
  OperatorRegistered as OperatorRegisteredEvent,
  OperatorDeregistered as OperatorDeregisteredEvent
} from "../generated/RegistryCoordinator_Operator/RegistryCoordinator"
import { NewPubkeyRegistration as NewPubkeyRegistrationEvent } from "../generated/BLSApkRegistry/BLSApkRegistry"
import { Operator } from "../generated/schema"
import { BLSApkRegistry } from "../generated/BLSApkRegistry/BLSApkRegistry"

export function handleOperatorDeregistered(event: OperatorDeregisteredEvent) : void {
  let entity = Operator.load(event.params.operatorId)
  if (entity == null) {
    log.error("Operator {} not found", [event.params.operatorId.toString()])
    return
  }

  entity.deregistrationBlockNumber = event.block.number

  entity.save()
}

export function handleOperatorRegistered(event: OperatorRegisteredEvent) : void {
  let entity = Operator.load(event.params.operatorId)
  if (entity == null) {
    log.error("Operator {} not found", [event.params.operatorId.toString()])
    return
  }

  entity.deregistrationBlockNumber = BigInt.fromU32(4294967295)

  entity.save()
}

