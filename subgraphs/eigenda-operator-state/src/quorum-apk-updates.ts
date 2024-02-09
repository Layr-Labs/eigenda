import { Address, BigInt, Bytes } from "@graphprotocol/graph-ts"
import {
    BLSApkRegistry,
    OperatorAddedToQuorums as OperatorAddedToQuorumsEvent,
    OperatorRemovedFromQuorums as OperatorRemovedFromQuorumsEvent
  } from "../generated/BLSApkRegistry_QuorumApkUpdates/BLSApkRegistry"
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

function updateApks(blsApkRegistryAddress: Address, quorumApkIdPrefix: Bytes, quorumNumbers: Bytes, blockNumber: BigInt, blockTimestamp: BigInt): void {
    // create a binding for blspubkeyregistry
    let blsApkRegistry = BLSApkRegistry.bind(blsApkRegistryAddress)
    // for each quorum, get the apk from the contract and store it as an entity
    for (let i = 0; i < quorumNumbers.length; i++) {
        let quorumNumber = quorumNumbers[i]
        let quorumApk = new QuorumApk(
            quorumApkIdPrefix.concatI32(quorumNumber)
        )
        quorumApk.quorumNumber = quorumNumber
        // get the apk from the contract
        let apk = blsApkRegistry.getApk(quorumNumber)
        quorumApk.apk_X = apk.X
        quorumApk.apk_Y = apk.Y

        quorumApk.blockNumber = blockNumber
        quorumApk.blockTimestamp = blockTimestamp
        quorumApk.save()
    }
}