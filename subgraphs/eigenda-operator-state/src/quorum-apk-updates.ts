import { Address, BigInt, Bytes } from "@graphprotocol/graph-ts"
import {
    BLSPubkeyRegistry,
    OperatorAddedToQuorums as OperatorAddedToQuorumsEvent,
    OperatorRemovedFromQuorums as OperatorRemovedFromQuorumsEvent
  } from "../generated/BLSPubkeyRegistry_QuorumApkUpdates/BLSPubkeyRegistry"
import {
    QuorumApk
} from "../generated/schema"

export function handleOperatorAddedToQuorums(
    event: OperatorAddedToQuorumsEvent
): void {
    updateApks(event.address, event.transaction.hash.concatI32(event.logIndex.toI32()), event.params.quorumNumbers, event.block.number, event.block.timestamp);
}

export function handleOperatorRemovedFromQuorums(
    event: OperatorRemovedFromQuorumsEvent
): void {
    updateApks(event.address, event.transaction.hash.concatI32(event.logIndex.toI32()), event.params.quorumNumbers, event.block.number, event.block.timestamp);
}

function updateApks(blsPubkeyRegistryAddress: Address, quorumApkIdPrefix: Bytes, quorumNumbers: Bytes, blockNumber: BigInt, blockTimestamp: BigInt): void {
    // create a binding for blspubkeyregistry
    let blsPubkeyRegistry = BLSPubkeyRegistry.bind(blsPubkeyRegistryAddress)
    // for each quorum, get the apk from the contract and store it as an entity
    for (let i = 0; i < quorumNumbers.length; i++) {
        let quorumNumber = quorumNumbers[i]
        let quorumApk = new QuorumApk(
            quorumApkIdPrefix.concatI32(quorumNumber)
        )
        quorumApk.quorumNumber = quorumNumber
        // get the apk from the contract
        let apk = blsPubkeyRegistry.getApkForQuorum(quorumNumber)
        quorumApk.apk_X = apk.X
        quorumApk.apk_Y = apk.Y

        quorumApk.blockNumber = blockNumber
        quorumApk.blockTimestamp = blockTimestamp
        quorumApk.save()
    }
}