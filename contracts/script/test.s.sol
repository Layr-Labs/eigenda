// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import "../src/core/EigenDAServiceManager.sol";
import "@eigenlayer-middleware/BLSOperatorStateRetriever.sol";
import "forge-std/Test.sol";

contract AAA is Test {
    using BN254 for BN254.G1Point;
    
    //forge test --match-test AAA --fork-url https://eth-goerli.g.alchemy.com/v2/nlFCarNvkza_SiVvw_l30OwkA4RPth1N
    function run() public {

        address dasm = address(0x9FcE30E01a740660189bD8CbEaA48Abd36040010);
        address churner = address(0xC0996A3Cc9ECF2A96115C117f6Da99FA80F525eB);

        IEigenDAServiceManager.BatchHeader memory batchHeader = IEigenDAServiceManager.BatchHeader({
            blobHeadersRoot: 0x06b37930d223188f8333dbe5d685dbc8f62906911e5f55e0f6dd7c3a00347a56,
            quorumNumbers: hex"00",
            quorumThresholdPercentages: hex"01",
            referenceBlockNumber: 10050612
        });

        uint32[] memory nsqbi = new uint32[](2);
        nsqbi[0] = 0;
        nsqbi[1] = 0;

        BN254.G1Point[] memory nspk = new BN254.G1Point[](2);
        nspk[0] = BN254.G1Point({X: 18640509892842556982483418152281752986979222329843183887797739632948938639891, Y: 12109273595002599203795592657066988712649206626110784817759094968276119587723});
        nspk[1] = BN254.G1Point({X: 8307250515012432645297735748562320972030003417271615691733500238974556229182, Y: 975942404230366830556184798459567435371298538490041668565570472597801325302});


        BN254.G1Point[] memory qa = new BN254.G1Point[](1);
        qa[0] = BN254.G1Point({X: 2310612090003801016696802914614438543181656173319641536064685846546040517238, Y: 17590614933768748865253655782571288820538767826882424195527080820300100234399});

        BN254.G2Point memory apkg22 = BN254.G2Point({X: [10863966866581104870598650905323087788543152324143724545808765994326885039537, 9416851004264125172476766859775404187731743736408875143950936833804374836551], Y: [9239714728935900206842727570378648411274905327222577522230236001530547176447, 14380719517706952250177033681997062225978229368453843066916053645031326845670]});

        BN254.G1Point memory sig = BN254.G1Point({X: 15587483973917301238715993224826968930619273131434990419379847737439932896088, Y: 17840911269102125334501453432074817402378948232619942392525762522328209396484});

        uint32[] memory qai = new uint32[](1);
        qai[0] = 3;

        uint32[] memory tsi = new uint32[](1);
        tsi[0] = 10;

        //[3,0]
        uint32[][] memory nssi = new uint32[][](1);
        nssi[0] = new uint32[](2);
        nssi[0][0] = 3;
        nssi[0][1] = 0;


        IBLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature = IBLSSignatureChecker.NonSignerStakesAndSignature({
            nonSignerQuorumBitmapIndices: nsqbi,
            nonSignerPubkeys: nspk,
            quorumApks: qa,
            apkG2: apkg22,
            sigma: sig,
            quorumApkIndices: qai,
            totalStakeIndices: tsi,
            nonSignerStakeIndices: nssi
        });

        bytes32[] memory xyz = new bytes32[](2);
        xyz[0] = bytes32(0x29f7bf28855f652da9f0c4316466b92781560b24bb69e2156093b4cdc93b1ddc);
        xyz[1] = bytes32(0x29f7bf28855f652da9f0c4316466b92781560b24bb69e2156093b4cdc93b1ddc);

        vm.startPrank(churner, churner);
        // BLSOperatorStateRetriever(0x737Dd62816a9392e84Fa21C531aF77C00816A3a3).getCheckSignaturesIndices(
        //     IBLSRegistryCoordinatorWithIndices(0x0b30a3427765f136754368a4500bAca8d2a54C0B),
        //     10050612,
        //     hex"00",
        //     xyz
        // );
        EigenDAServiceManager(dasm).confirmBatch(
            batchHeader,
            nonSignerStakesAndSignature
        );
        vm.stopPrank();
    }
}