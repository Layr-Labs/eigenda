// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import "../MockEigenDADeployer.sol";

contract EigenDAServiceManagerUnit is MockEigenDADeployer {
    using BN254 for BN254.G1Point;
    using EigenDAHasher for BatchHeader;
    using EigenDAHasher for ReducedBatchHeader;

    event BatchConfirmed(bytes32 indexed batchHeaderHash, uint32 batchId);

    function setUp() virtual public {
        _deployDA();
    }

    function testConfirmBatch_AllSigning_Valid(uint256 pseudoRandomNumber) public {
        (BatchHeader memory batchHeader, BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature) 
            = _getHeaderandNonSigners(0, pseudoRandomNumber, 100);
        
        uint32 batchIdToConfirm = eigenDAServiceManager.batchId();
        bytes32 batchHeaderHash = batchHeader.hashBatchHeaderToReducedBatchHeader();

        cheats.prank(confirmer, confirmer);
        cheats.expectEmit(true, true, true, true, address(eigenDAServiceManager));
        emit BatchConfirmed(batchHeaderHash, batchIdToConfirm);
        uint256 gasBefore = gasleft();
        eigenDAServiceManager.confirmBatch(
            batchHeader,
            nonSignerStakesAndSignature
        );
        uint256 gasAfter = gasleft();
        emit log_named_uint("gasUsed", gasBefore - gasAfter);

        assertEq(eigenDAServiceManager.batchId(), batchIdToConfirm + 1);
    }

    function testConfirmBatch_Revert_NotEOA(uint256 pseudoRandomNumber) public {
        (BatchHeader memory batchHeader, BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature) 
            = _getHeaderandNonSigners(0, pseudoRandomNumber, 100);

        cheats.expectRevert(bytes("header and nonsigner data must be in calldata"));
        cheats.prank(confirmer, notConfirmer);
        eigenDAServiceManager.confirmBatch(
            batchHeader,
            nonSignerStakesAndSignature
        );
    }

    function testConfirmBatch_Revert_NotConfirmer(uint256 pseudoRandomNumber) public {
        (BatchHeader memory batchHeader, BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature) 
            = _getHeaderandNonSigners(0, pseudoRandomNumber, 100);

        cheats.expectRevert();
        cheats.prank(notConfirmer, notConfirmer);
        eigenDAServiceManager.confirmBatch(
            batchHeader,
            nonSignerStakesAndSignature
        );
    }

    function testConfirmBatch_Revert_FutureBlocknumber(uint256 pseudoRandomNumber) public {
        uint256 quorumBitmap = 1;
        bytes memory quorumNumbers = BitmapUtils.bitmapToBytesArray(quorumBitmap);

        (, BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature) = 
            _registerSignatoriesAndGetNonSignerStakeAndSignatureRandom(pseudoRandomNumber, 0, quorumBitmap);

        BatchHeader memory batchHeader = 
            _getRandomBatchHeader(pseudoRandomNumber, quorumNumbers, uint32(block.number + 1), 100);

        bytes32 batchHeaderHash = batchHeader.hashBatchHeaderMemory();
        nonSignerStakesAndSignature.sigma = BN254.hashToG1(batchHeaderHash).scalar_mul(aggSignerPrivKey);

        cheats.expectRevert(bytes("specified referenceBlockNumber is in future"));
        cheats.prank(confirmer, confirmer);
        eigenDAServiceManager.confirmBatch(
            batchHeader,
            nonSignerStakesAndSignature
        );
    }

    function testConfirmBatch_Revert_PastBlocknumber(uint256 pseudoRandomNumber) public {
        (BatchHeader memory batchHeader, BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature) 
            = _getHeaderandNonSigners(0, pseudoRandomNumber, 100);

        cheats.roll(block.number + eigenDAServiceManager.BLOCK_STALE_MEASURE());
        cheats.expectRevert(bytes("specified referenceBlockNumber is too far in past"));
        cheats.prank(confirmer, confirmer);
        eigenDAServiceManager.confirmBatch(
            batchHeader,
            nonSignerStakesAndSignature
        );
    }

    function testConfirmBatch_Revert_Threshold(uint256 pseudoRandomNumber) public {
        (BatchHeader memory batchHeader, BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature) 
            = _getHeaderandNonSigners(1, pseudoRandomNumber, 100);

        cheats.expectRevert(bytes("signatories do not own threshold percentage of a quorum"));
        cheats.prank(confirmer, confirmer);
        eigenDAServiceManager.confirmBatch(
            batchHeader,
            nonSignerStakesAndSignature
        );
    }

    function testConfirmBatch_NonSigner_Valid(uint256 pseudoRandomNumber) public {
        (BatchHeader memory batchHeader, BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature) 
            = _getHeaderandNonSigners(1, pseudoRandomNumber, 75);

        uint32 batchIdToConfirm = eigenDAServiceManager.batchId();
        bytes32 batchHeaderHash = batchHeader.hashBatchHeaderToReducedBatchHeader();

        cheats.prank(confirmer, confirmer);
        cheats.expectEmit(true, true, true, true, address(eigenDAServiceManager));
        emit BatchConfirmed(batchHeaderHash, batchIdToConfirm);
        uint256 gasBefore = gasleft();
        eigenDAServiceManager.confirmBatch(
            batchHeader,
            nonSignerStakesAndSignature
        );
        uint256 gasAfter = gasleft();
        emit log_named_uint("gasUsed", gasBefore - gasAfter);

        assertEq(eigenDAServiceManager.batchId(), batchIdToConfirm + 1);
    }

    function testConfirmBatch_Revert_LengthMismatch(uint256 pseudoRandomNumber) public {
        (BatchHeader memory batchHeader, BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature) 
            = _getHeaderandNonSigners(0, pseudoRandomNumber, 100);
        batchHeader.signedStakeForQuorums = new bytes(0);

        cheats.expectRevert(bytes("quorumNumbers and signedStakeForQuorums must be same length"));
        cheats.prank(confirmer, confirmer);
        eigenDAServiceManager.confirmBatch(
            batchHeader,
            nonSignerStakesAndSignature
        );
    }
}