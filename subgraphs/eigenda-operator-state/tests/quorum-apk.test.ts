import {
    assert,
    describe,
    test,
    clearStore,
    beforeAll,
    afterAll,
    newMockCall,
    createMockedFunction
  } from "matchstick-as/assembly/index"
  import { Address, BigInt, Bytes, ethereum, log } from "@graphprotocol/graph-ts"
  import { BLSPubkeyRegistry, BLSPubkeyRegistry__getApkForQuorumResultValue0Struct } from "../generated/BLSPubkeyRegistry_QuorumApkUpdates/BLSPubkeyRegistry"
  import { createNewOperatorAddedToQuorumsEvent, createNewOperatorRemovedFromQuorumsEvent } from "./quorum-apk-utils"
  import { handleOperatorAddedToQuorums, handleOperatorRemovedFromQuorums } from "../src/quorum-apk-updates"
  
  
  
  let operator: Address = Address.fromBytes(Bytes.fromHexString("0xa16081f360e3847006db660bae1c6d1b2e17ec2a"))
  
  function generateRandomPublicKeyFromSeed(seed: string): ethereum.Tuple {
    let pubkeyG1_X = BigInt.fromString(seed)
    let pubkeyG1_Y = BigInt.fromString(seed + "1")
  
    let apk = new BLSPubkeyRegistry__getApkForQuorumResultValue0Struct(2);
    apk[0] = ethereum.Value.fromUnsignedBigInt(pubkeyG1_X)
    apk[1] = ethereum.Value.fromUnsignedBigInt(pubkeyG1_Y)
    
    return apk
  }
  
  describe("Describe entity assertions", () => {
    beforeAll(() => {
  
    })
  
    afterAll(() => {
      clearStore()
    })
  
    // For more test scenarios, see:
    // https://thegraph.com/docs/en/developer/matchstick/#write-a-unit-test
  
    test("quorum apks updates on operators added", () => {
      let quorumNumbers1 = Bytes.fromHexString("0x0102030405")
      let quorumApks1: ethereum.Tuple[] = []
      for(let i = 0; i < quorumNumbers1.length; i++) {
        let quorumNumber = quorumNumbers1[i]
        quorumApks1.push(generateRandomPublicKeyFromSeed((quorumNumber + 128375).toString()))
      }
      let quorumNumbers2 = Bytes.fromHexString("0x01415379")
      let quorumApks2: ethereum.Tuple[] = []
      for(let i = 0; i < quorumNumbers2.length; i++) {
        let quorumNumber = quorumNumbers2[i]
        quorumApks2.push(generateRandomPublicKeyFromSeed((quorumNumber + 234612).toString()))
      }
      
      let newOperatorAddedToQuorumsEvent1 = createNewOperatorAddedToQuorumsEvent(
        operator,
        quorumNumbers1
      )
      
      // for each quroum in quorumNumbers, mock the call to getApkForQuorum
      for(let i = 0; i < quorumNumbers1.length; i++) {
        let quorumNumber = quorumNumbers1[i]
        let quorumNumberBigInt = BigInt.fromI32(quorumNumber)
        createMockedFunction(newOperatorAddedToQuorumsEvent1.address, 'getApkForQuorum', 'getApkForQuorum(uint8):((uint256,uint256))')
          .withArgs([ethereum.Value.fromUnsignedBigInt(quorumNumberBigInt)])
          .returns([ethereum.Value.fromTuple(quorumApks1[i])])
      }
      
      handleOperatorAddedToQuorums(newOperatorAddedToQuorumsEvent1)
  
      assert.entityCount("QuorumApk", quorumNumbers1.length)
  
      assert.entityCount("QuorumApk", quorumNumbers1.length)
      checkQuorumApkEntities(newOperatorAddedToQuorumsEvent1.transaction.hash, newOperatorAddedToQuorumsEvent1.logIndex, quorumNumbers1, quorumApks1)
  
      let newOperatorRemovedFromQuorumsEvent2 = createNewOperatorRemovedFromQuorumsEvent(operator, quorumNumbers2)
      newOperatorRemovedFromQuorumsEvent2.logIndex = newOperatorAddedToQuorumsEvent1.logIndex.plus(BigInt.fromI32(1))
  
      // for each quroum in quorumNumbers, mock the call to getApkForQuorum
      for(let i = 0; i < quorumNumbers2.length; i++) {
        let quorumNumber = quorumNumbers2[i]
        let quorumNumberBigInt = BigInt.fromI32(quorumNumber)
        createMockedFunction(newOperatorRemovedFromQuorumsEvent2.address, 'getApkForQuorum', 'getApkForQuorum(uint8):((uint256,uint256))')
          .withArgs([ethereum.Value.fromUnsignedBigInt(quorumNumberBigInt)])
          .returns([ethereum.Value.fromTuple(quorumApks2[i])])
      }
  
      handleOperatorRemovedFromQuorums(newOperatorRemovedFromQuorumsEvent2)
  
      assert.entityCount("QuorumApk", quorumNumbers1.length + quorumNumbers2.length)
      checkQuorumApkEntities(newOperatorAddedToQuorumsEvent1.transaction.hash, newOperatorAddedToQuorumsEvent1.logIndex, quorumNumbers1, quorumApks1)
      checkQuorumApkEntities(newOperatorRemovedFromQuorumsEvent2.transaction.hash, newOperatorRemovedFromQuorumsEvent2.logIndex, quorumNumbers2, quorumApks2)
    })
  })

function checkQuorumApkEntities(txHash: Bytes, logIndex: BigInt, quorumNumbers: Bytes, quorumApks: ethereum.Tuple[]): void {
    for(let i = 0; i < quorumNumbers.length; i++) {
        let quorumNumber = quorumNumbers[i]
        let apkId = txHash.concatI32(logIndex.toI32()).concatI32(quorumNumber).toHexString()
        assert.fieldEquals(
          "QuorumApk",
          apkId,
          "quorumNumber",
          quorumNumber.toString()
        )
        assert.fieldEquals(
            "QuorumApk",
            apkId,
            "apk_X",
            quorumApks[i][0].toBigInt().toString()
        )
        assert.fieldEquals(
            "QuorumApk",
            apkId,
            "apk_Y",
            quorumApks[i][1].toBigInt().toString()
        )
    }
}