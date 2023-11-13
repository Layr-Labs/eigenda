import { BigInt, Bytes, log } from "@graphprotocol/graph-ts"
import { NewPubkeyRegistration as NewPubkeyRegistrationEvent } from "../generated/BLSPubkeyCompendium_Operator/BLSPubkeyCompendium"
import { Operator } from "../generated/schema"
import { BLSPubkeyCompendium } from "../generated/BLSPubkeyCompendium/BLSPubkeyCompendium"

export function handleNewPubkeyRegistration(
  event: NewPubkeyRegistrationEvent
): void {
  let pubkeyCompendium = BLSPubkeyCompendium.bind(event.address)

  let entity = new Operator(
    pubkeyCompendium.operatorToPubkeyHash(event.params.operator) // this is the operator id
  )

  entity.operator = event.params.operator
  entity.pubkeyG1_X = event.params.pubkeyG1.X
  entity.pubkeyG1_Y = event.params.pubkeyG1.Y
  entity.pubkeyG2_X = event.params.pubkeyG2.X
  entity.pubkeyG2_Y = event.params.pubkeyG2.Y
  entity.deregistrationBlockNumber = BigInt.fromI32(0)

  entity.save()
}
