import { BigInt, Bytes, log } from "@graphprotocol/graph-ts"
import {
  OperatorRegistered as OperatorRegisteredEvent,
  OperatorDeregistered as OperatorDeregisteredEvent
} from "../generated/BLSRegistryCoordinatorWithIndices_Operator/BLSRegistryCoordinatorWithIndices"
import { NewPubkeyRegistration as NewPubkeyRegistrationEvent } from "../generated/BLSPubkeyCompendium/BLSPubkeyCompendium"
import { Operator } from "../generated/schema"
import { BLSPubkeyCompendium } from "../generated/BLSPubkeyCompendium/BLSPubkeyCompendium"

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

