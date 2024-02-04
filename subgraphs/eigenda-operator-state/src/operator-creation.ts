import { BigInt, Bytes, log } from "@graphprotocol/graph-ts"
import { NewPubkeyRegistration as NewPubkeyRegistrationEvent } from "../generated/BLSApkRegistry_Operator/BLSApkRegistry"
import { Operator } from "../generated/schema"
import { BLSApkRegistry } from "../generated/BLSApkRegistry/BLSApkRegistry"

export function handleNewPubkeyRegistration(
  event: NewPubkeyRegistrationEvent
): void {
  let apkRegistry = BLSApkRegistry.bind(event.address)

  let entity = new Operator(
    apkRegistry.operatorToPubkeyHash(event.params.operator) // this is the operator id
  )

  entity.operator = event.params.operator
  entity.pubkeyG1_X = event.params.pubkeyG1.X
  entity.pubkeyG1_Y = event.params.pubkeyG1.Y
  entity.pubkeyG2_X = event.params.pubkeyG2.X
  entity.pubkeyG2_Y = event.params.pubkeyG2.Y
  entity.deregistrationBlockNumber = BigInt.fromI32(0)

  entity.save()
}
