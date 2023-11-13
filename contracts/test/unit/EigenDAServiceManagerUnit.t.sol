// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import "../../lib/eigenlayer-middleware/test/utils/BLSMockAVSDeployer.sol";
import {EigenDAServiceManager} from "../../src/core/EigenDAServiceManager.sol";
import {EigenDAHasher} from "../../src/libraries/EigenDAHasher.sol";
import {EigenDAServiceManager, IEigenDAServiceManager} from "../../src/core/EigenDAServiceManager.sol";

contract EigenDAServiceManagerUnit is BLSMockAVSDeployer {
    using BN254 for BN254.G1Point;
    using EigenDAHasher for IEigenDAServiceManager.BatchHeader;
    using EigenDAHasher for IEigenDAServiceManager.ReducedBatchHeader;

    address confirmer = address(uint160(uint256(keccak256(abi.encodePacked("confirmer")))));
    address notConfirmer = address(uint160(uint256(keccak256(abi.encodePacked("notConfirmer")))));
    address newFeeSetter = address(uint160(uint256(keccak256(abi.encodePacked("newFeeSetter")))));

    EigenDAServiceManager eigenDAServiceManager;
    EigenDAServiceManager eigenDAServiceManagerImplementation;

    uint256 feePerBytePerTime = 0;

    event BatchConfirmed(bytes32 indexed batchHeaderHash, uint32 batchId, uint96 fee);
    event FeePerBytePerTimeSet(uint256 previousValue, uint256 newValue);
    event FeeSetterChanged(address previousAddress, address newAddress);

    function setUp() virtual public {
        _setUpBLSMockAVSDeployer();

        eigenDAServiceManagerImplementation = new EigenDAServiceManager(
            registryCoordinator,
            strategyManagerMock,
            delegationMock,
            slasher
        );

        // Third, upgrade the proxy contracts to use the correct implementation contracts and initialize them.
        eigenDAServiceManager = EigenDAServiceManager(
            address(
                new TransparentUpgradeableProxy(
                    address(eigenDAServiceManagerImplementation),
                    address(proxyAdmin),
                    abi.encodeWithSelector(
                        EigenDAServiceManager.initialize.selector,
                        pauserRegistry,
                        serviceManagerOwner
                    )
                )
            )
        );
    }

    function testConfirmBatch_AllSigning_Valid(uint256 pseudoRandomNumber) public {
        (IEigenDAServiceManager.BatchHeader memory batchHeader, BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature) 
            = _getHeaderandNonSigners(0, pseudoRandomNumber, 100);
        
        uint32 batchIdToConfirm = eigenDAServiceManager.batchId();
        bytes32 batchHeaderHash = batchHeader.hashBatchHeaderToReducedBatchHeader();

        cheats.prank(confirmer, confirmer);
        cheats.expectEmit(true, true, true, true, address(eigenDAServiceManager));
        emit BatchConfirmed(batchHeaderHash, batchIdToConfirm, 0);
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
        (IEigenDAServiceManager.BatchHeader memory batchHeader, BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature) 
            = _getHeaderandNonSigners(0, pseudoRandomNumber, 100);

        cheats.expectRevert(bytes("EigenDAServiceManager.confirmBatch: header and nonsigner data must be in calldata"));
        cheats.prank(confirmer, notConfirmer);
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

        IEigenDAServiceManager.BatchHeader memory batchHeader = 
            _getRandomBatchHeader(pseudoRandomNumber, quorumNumbers, uint32(block.number + 1), 100);

        bytes32 batchHeaderHash = batchHeader.hashBatchHeaderMemory();
        nonSignerStakesAndSignature.sigma = BN254.hashToG1(batchHeaderHash).scalar_mul(aggSignerPrivKey);

        cheats.expectRevert(bytes("EigenDAServiceManager.confirmBatch: specified referenceBlockNumber is in future"));
        cheats.prank(confirmer, confirmer);
        eigenDAServiceManager.confirmBatch(
            batchHeader,
            nonSignerStakesAndSignature
        );
    }

    function testConfirmBatch_Revert_PastBlocknumber(uint256 pseudoRandomNumber) public {
        (IEigenDAServiceManager.BatchHeader memory batchHeader, BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature) 
            = _getHeaderandNonSigners(0, pseudoRandomNumber, 100);

        cheats.roll(block.number + eigenDAServiceManager.BLOCK_STALE_MEASURE());
        cheats.expectRevert(bytes("EigenDAServiceManager.confirmBatch: specified referenceBlockNumber is too far in past"));
        cheats.prank(confirmer, confirmer);
        eigenDAServiceManager.confirmBatch(
            batchHeader,
            nonSignerStakesAndSignature
        );
    }

    function testConfirmBatch_Revert_Threshold(uint256 pseudoRandomNumber) public {
        (IEigenDAServiceManager.BatchHeader memory batchHeader, BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature) 
            = _getHeaderandNonSigners(1, pseudoRandomNumber, 100);

        cheats.expectRevert(bytes("EigenDAServiceManager.confirmBatch: signatories do not own at least threshold percentage of a quorum"));
        cheats.prank(confirmer, confirmer);
        eigenDAServiceManager.confirmBatch(
            batchHeader,
            nonSignerStakesAndSignature
        );
    }

    function testConfirmBatch_NonSigner_Valid(uint256 pseudoRandomNumber) public {
        (IEigenDAServiceManager.BatchHeader memory batchHeader, BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature) 
            = _getHeaderandNonSigners(1, pseudoRandomNumber, 75);

        uint32 batchIdToConfirm = eigenDAServiceManager.batchId();
        bytes32 batchHeaderHash = batchHeader.hashBatchHeaderToReducedBatchHeader();

        cheats.prank(confirmer, confirmer);
        cheats.expectEmit(true, true, true, true, address(eigenDAServiceManager));
        emit BatchConfirmed(batchHeaderHash, batchIdToConfirm, 0);
        uint256 gasBefore = gasleft();
        eigenDAServiceManager.confirmBatch(
            batchHeader,
            nonSignerStakesAndSignature
        );
        uint256 gasAfter = gasleft();
        emit log_named_uint("gasUsed", gasBefore - gasAfter);

        assertEq(eigenDAServiceManager.batchId(), batchIdToConfirm + 1);
    }

    function testFreezeOperator_Revert() public {
        cheats.expectRevert(bytes("EigenDAServiceManager.freezeOperator: not implemented"));
        eigenDAServiceManager.freezeOperator(address(0));
    }


    function _getHeaderandNonSigners(uint256 _nonSigners, uint256 _pseudoRandomNumber, uint8 _threshold) internal returns (IEigenDAServiceManager.BatchHeader memory, BLSSignatureChecker.NonSignerStakesAndSignature memory) {
        // register a bunch of operators
        uint256 quorumBitmap = 1;
        bytes memory quorumNumbers = BitmapUtils.bitmapToBytesArray(quorumBitmap);

        // 0 nonSigners
        (uint32 referenceBlockNumber, BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature) = 
            _registerSignatoriesAndGetNonSignerStakeAndSignatureRandom(_pseudoRandomNumber, _nonSigners, quorumBitmap);

        // get a random batch header
        IEigenDAServiceManager.BatchHeader memory batchHeader = _getRandomBatchHeader(_pseudoRandomNumber, quorumNumbers, referenceBlockNumber, _threshold);

        // set batch specific signature
        bytes32 reducedBatchHeaderHash = batchHeader.hashBatchHeaderToReducedBatchHeader();
        nonSignerStakesAndSignature.sigma = BN254.hashToG1(reducedBatchHeaderHash).scalar_mul(aggSignerPrivKey);

        return (batchHeader, nonSignerStakesAndSignature);
    }

    function _getRandomBatchHeader(uint256 pseudoRandomNumber, bytes memory quorumNumbers, uint32 referenceBlockNumber, uint8 threshold) internal pure returns(IEigenDAServiceManager.BatchHeader memory) {
        IEigenDAServiceManager.BatchHeader memory batchHeader;
        batchHeader.blobHeadersRoot = keccak256(abi.encodePacked("blobHeadersRoot", pseudoRandomNumber));
        batchHeader.quorumNumbers = quorumNumbers;
        batchHeader.quorumThresholdPercentages = new bytes(quorumNumbers.length);
        for (uint256 i = 0; i < quorumNumbers.length; i++) {
            batchHeader.quorumThresholdPercentages[i] = bytes1(threshold);
        }
        batchHeader.referenceBlockNumber = referenceBlockNumber;
        return batchHeader;
    }
}