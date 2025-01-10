// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import "../../lib/eigenlayer-middleware/test/utils/BLSMockAVSDeployer.sol";
import {EigenDAServiceManager, IRewardsCoordinator} from "../../src/core/EigenDAServiceManager.sol";
import {EigenDAHasher} from "../../src/libraries/EigenDAHasher.sol";
import {EigenDAServiceManager} from "../../src/core/EigenDAServiceManager.sol";
import {IEigenDAServiceManager} from "../../src/interfaces/IEigenDAServiceManager.sol";
import {EigenDABlobVerifier} from "../../src/core/EigenDABlobVerifier.sol";
import {EigenDAThresholdRegistry, IEigenDAThresholdRegistry} from "../../src/core/EigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "../../src/interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "../../src/interfaces/IEigenDASignatureVerifier.sol";
import {IRegistryCoordinator} from "../../lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {IEigenDARelayRegistry} from "../../src/interfaces/IEigenDARelayRegistry.sol";
import {IPaymentVault} from "../../src/interfaces/IPaymentVault.sol";
import {EigenDARelayRegistry} from "../../src/core/EigenDARelayRegistry.sol";
import {IEigenDADisperserRegistry} from "../../src/interfaces/IEigenDADisperserRegistry.sol";
import "../../src/interfaces/IEigenDAStructs.sol";

contract EigenDAServiceManagerUnit is BLSMockAVSDeployer {
    using BN254 for BN254.G1Point;
    using EigenDAHasher for BatchHeader;
    using EigenDAHasher for ReducedBatchHeader;

    address confirmer = address(uint160(uint256(keccak256(abi.encodePacked("confirmer")))));
    address notConfirmer = address(uint160(uint256(keccak256(abi.encodePacked("notConfirmer")))));
    address newFeeSetter = address(uint160(uint256(keccak256(abi.encodePacked("newFeeSetter")))));
    address rewardsInitiator = address(uint160(uint256(keccak256(abi.encodePacked("rewardsInitiator")))));

    EigenDAServiceManager eigenDAServiceManager;
    EigenDAServiceManager eigenDAServiceManagerImplementation;
    EigenDABlobVerifier eigenDABlobVerifier;
    EigenDARelayRegistry eigenDARelayRegistry;
    EigenDARelayRegistry eigenDARelayRegistryImplementation;

    EigenDAThresholdRegistry eigenDAThresholdRegistry;
    EigenDAThresholdRegistry eigenDAThresholdRegistryImplementation;
    bytes quorumAdversaryThresholdPercentages = hex"212121";
    bytes quorumConfirmationThresholdPercentages = hex"373737";
    bytes quorumNumbersRequired = hex"0001";
    SecurityThresholds defaultSecurityThresholds = SecurityThresholds(55, 33);

    uint256 feePerBytePerTime = 0;

    event BatchConfirmed(bytes32 indexed batchHeaderHash, uint32 batchId);
    event FeePerBytePerTimeSet(uint256 previousValue, uint256 newValue);
    event FeeSetterChanged(address previousAddress, address newAddress);

    function setUp() virtual public {
        _setUpBLSMockAVSDeployer();

        eigenDAServiceManager = EigenDAServiceManager(
            address(
                new TransparentUpgradeableProxy(address(emptyContract), address(proxyAdmin), "")
            )
        );

        eigenDAThresholdRegistry = EigenDAThresholdRegistry(
            address(
                new TransparentUpgradeableProxy(address(emptyContract), address(proxyAdmin), "")
            )
        );

        eigenDARelayRegistry = EigenDARelayRegistry(
            address(
                new TransparentUpgradeableProxy(address(emptyContract), address(proxyAdmin), "")
            )
        );

        eigenDAThresholdRegistryImplementation = new EigenDAThresholdRegistry();

        eigenDAServiceManagerImplementation = new EigenDAServiceManager(
            avsDirectory,
            rewardsCoordinator,
            registryCoordinator,
            stakeRegistry,
            eigenDAThresholdRegistry,
            eigenDARelayRegistry,
            IPaymentVault(address(0)),
            IEigenDADisperserRegistry(address(0))
        );

        address[] memory confirmers = new address[](1);
        confirmers[0] = confirmer;

        cheats.prank(proxyAdminOwner);
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenDAServiceManager))),
            address(eigenDAServiceManagerImplementation),
            abi.encodeWithSelector(
                EigenDAServiceManager.initialize.selector,
                pauserRegistry,
                0,
                registryCoordinatorOwner,
                confirmers,
                registryCoordinatorOwner
            )
        );

        VersionedBlobParams[] memory versionedBlobParams = new VersionedBlobParams[](0);

        cheats.prank(proxyAdminOwner);
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenDAThresholdRegistry))),
            address(eigenDAThresholdRegistryImplementation),
            abi.encodeWithSelector(
                EigenDAThresholdRegistry.initialize.selector,
                registryCoordinatorOwner,
                quorumAdversaryThresholdPercentages,
                quorumConfirmationThresholdPercentages,
                quorumNumbersRequired,
                versionedBlobParams,
                defaultSecurityThresholds
            )
        );

        eigenDABlobVerifier = new EigenDABlobVerifier(
            IEigenDAThresholdRegistry(address(eigenDAThresholdRegistry)),
            IEigenDABatchMetadataStorage(address(eigenDAServiceManager)),
            IEigenDASignatureVerifier(address(eigenDAServiceManager)),
            IEigenDARelayRegistry(address(eigenDARelayRegistry)),
            OperatorStateRetriever(address(operatorStateRetriever)),
            IRegistryCoordinator(address(registryCoordinator))
        );

        eigenDARelayRegistryImplementation = new EigenDARelayRegistry();

        cheats.prank(proxyAdminOwner);
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenDARelayRegistry))),
            address(eigenDARelayRegistryImplementation),
            abi.encodeWithSelector(EigenDARelayRegistry.initialize.selector, registryCoordinatorOwner)
        );
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

    function _getHeaderandNonSigners(uint256 _nonSigners, uint256 _pseudoRandomNumber, uint8 _threshold) internal returns (BatchHeader memory, BLSSignatureChecker.NonSignerStakesAndSignature memory) {
        // register a bunch of operators
        uint256 quorumBitmap = 1;
        bytes memory quorumNumbers = BitmapUtils.bitmapToBytesArray(quorumBitmap);

        // 0 nonSigners
        (uint32 referenceBlockNumber, BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature) = 
            _registerSignatoriesAndGetNonSignerStakeAndSignatureRandom(_pseudoRandomNumber, _nonSigners, quorumBitmap);

        // get a random batch header
        BatchHeader memory batchHeader = _getRandomBatchHeader(_pseudoRandomNumber, quorumNumbers, referenceBlockNumber, _threshold);

        // set batch specific signature
        bytes32 reducedBatchHeaderHash = batchHeader.hashBatchHeaderToReducedBatchHeader();
        nonSignerStakesAndSignature.sigma = BN254.hashToG1(reducedBatchHeaderHash).scalar_mul(aggSignerPrivKey);

        return (batchHeader, nonSignerStakesAndSignature);
    }

    function _getRandomBatchHeader(uint256 pseudoRandomNumber, bytes memory quorumNumbers, uint32 referenceBlockNumber, uint8 threshold) internal pure returns(BatchHeader memory) {
        BatchHeader memory batchHeader;
        batchHeader.blobHeadersRoot = keccak256(abi.encodePacked("blobHeadersRoot", pseudoRandomNumber));
        batchHeader.quorumNumbers = quorumNumbers;
        batchHeader.signedStakeForQuorums = new bytes(quorumNumbers.length);
        for (uint256 i = 0; i < quorumNumbers.length; i++) {
            batchHeader.signedStakeForQuorums[i] = bytes1(threshold);
        }
        batchHeader.referenceBlockNumber = referenceBlockNumber;
        return batchHeader;
    }
}