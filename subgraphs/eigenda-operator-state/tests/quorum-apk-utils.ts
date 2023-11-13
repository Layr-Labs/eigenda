import { newMockEvent } from "matchstick-as"
import { ethereum, BigInt, Bytes, Address } from "@graphprotocol/graph-ts"
import { OperatorAddedToQuorums as OperatorAddedToQuorumsEvent, OperatorRemovedFromQuorums as OperatorRemovedFromQuorumsEvent } from "../generated/BLSPubkeyRegistry_QuorumApkUpdates/BLSPubkeyRegistry"

export function createNewOperatorAddedToQuorumsEvent(
  operator: Address,
  quorumNumbers: Bytes
): OperatorAddedToQuorumsEvent {
  let newOperatorAddedToQuorumsEvent = changetype<
    OperatorAddedToQuorumsEvent
  >(newMockEvent())

  newOperatorAddedToQuorumsEvent.parameters = new Array()

  newOperatorAddedToQuorumsEvent.parameters.push(
    new ethereum.EventParam("operator", ethereum.Value.fromAddress(operator))
  )
  newOperatorAddedToQuorumsEvent.parameters.push(
    new ethereum.EventParam("quorumNumbers", ethereum.Value.fromBytes(quorumNumbers))
  )

  return newOperatorAddedToQuorumsEvent
}

export function createNewOperatorRemovedFromQuorumsEvent(
  operator: Address,
  quorumNumbers: Bytes
): OperatorRemovedFromQuorumsEvent {
  let newOperatorRemovedFromQuorumsEvent = changetype<
    OperatorRemovedFromQuorumsEvent
  >(newMockEvent())

  newOperatorRemovedFromQuorumsEvent.parameters = new Array()

  newOperatorRemovedFromQuorumsEvent.parameters.push(
    new ethereum.EventParam("operator", ethereum.Value.fromAddress(operator))
  )
  newOperatorRemovedFromQuorumsEvent.parameters.push(
    new ethereum.EventParam("quorumNumbers", ethereum.Value.fromBytes(quorumNumbers))
  )

  return newOperatorRemovedFromQuorumsEvent
}



