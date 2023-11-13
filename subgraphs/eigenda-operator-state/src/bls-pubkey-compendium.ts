import { NewPubkeyRegistration as NewPubkeyRegistrationEvent } from "../generated/BLSPubkeyCompendium/BLSPubkeyCompendium"
import { NewPubkeyRegistration } from "../generated/schema"

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
