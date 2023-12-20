// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.9;

import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import {BLSMockAVSDeployer} from "../lib/eigenlayer-middleware/test/utils/BLSMockAVSDeployer.sol";
import {MockRollup, BN254} from "./mocks/MockRollup.sol";
import {EigenDAHasher} from "../src/libraries/EigenDAHasher.sol";
import {EigenDAServiceManager, IEigenDAServiceManager} from "../src/core/EigenDAServiceManager.sol";
import {EigenDABlobUtils} from "../src/libraries/EigenDABlobUtils.sol";
//import {BN254} from "../lib/eigenlayer-middleware/src/libraries/BN254.sol";

import "forge-std/StdStorage.sol";

contract MockRollupTest is BLSMockAVSDeployer {
    using stdStorage for StdStorage;
    using BN254 for BN254.G1Point;
    using EigenDAHasher for IEigenDAServiceManager.BatchHeader;
    using EigenDAHasher for IEigenDAServiceManager.ReducedBatchHeader;
    using EigenDAHasher for IEigenDAServiceManager.BlobHeader;
    using EigenDAHasher for IEigenDAServiceManager.BatchMetadata;

    EigenDAServiceManager eigenDAServiceManager;
    EigenDAServiceManager eigenDAServiceManagerImplementation;

    uint256 feePerBytePerTime = 0;
    uint8 defaultCodingRatioPercentage = 10;
    uint32 defaultReferenceBlockNumber = 100;
    uint32 defaultConfirmationBlockNumber = 1000;
    uint32 defaultBatchId = 0;
    uint256 defaultStakeRequired = 1 ether;

    mapping(uint8 => bool) public quorumNumbersUsed;

    address alice = address(0x101);
    address bob = address(0x202);

    MockRollup mockRollup;

    //powers of tau
    BN254.G1Point s0 = BN254.generatorG1().scalar_mul(1);
    BN254.G1Point s1 = BN254.generatorG1().scalar_mul(2);
    BN254.G1Point s2 = BN254.generatorG1().scalar_mul(4);
    BN254.G1Point s3 = BN254.generatorG1().scalar_mul(8);
    BN254.G1Point s4 = BN254.generatorG1().scalar_mul(16);

    uint256 illegalPoint = 6;
    uint256 illegalValue = 1555;
    BN254.G2Point illegalProof;
    bytes32 quorumBlobParamsHash = keccak256(abi.encodePacked("quorumBlobParamsHash"));

    function setUp() public {
        _setUpBLSMockAVSDeployer();

        eigenDAServiceManagerImplementation = new EigenDAServiceManager(
            registryCoordinator,
            strategyManagerMock,
            delegationMock,
            slasher
        );

        eigenDAServiceManager = EigenDAServiceManager(
            address(
                new TransparentUpgradeableProxy(
                    address(eigenDAServiceManagerImplementation),
                    address(proxyAdmin),
                    abi.encodeWithSelector(
                        EigenDAServiceManager.initialize.selector,
                        pauserRegistry,
                        serviceManagerOwner,
                        feePerBytePerTime,
                        serviceManagerOwner
                    )
                )
            )
        );

        mockRollup = new MockRollup(eigenDAServiceManager, s1, illegalValue, quorumBlobParamsHash, defaultStakeRequired);

        //hardcode g2 proof
        illegalProof.X[1] = 11151623676041303181597631684634074376466382703418354161831688442589830350329;
        illegalProof.X[0] = 21587740443732524623985464356760343072434825248946003815467233999912459579351;
        illegalProof.Y[1] = 4222041728992406478862708226745479381252734858741080790666424175645694456140;
        illegalProof.Y[0] = 17511259870083276759899704237100059449000397154439723516103658719937845846446;

    }

    function testChallenge(uint256 pseudoRandomNumber) public {
        //register alice with rollup
        vm.deal(alice, 1 ether);
        vm.prank(alice);
        mockRollup.registerValidator{value: 1 ether}();

        //get commitment with illegal value
        (IEigenDAServiceManager.BlobHeader memory blobHeader, EigenDABlobUtils.BlobVerificationProof memory blobVerificationProof) = _getCommitment(pseudoRandomNumber);

        IEigenDAServiceManager.QuorumBlobParam[] memory quorumBlobParamsCopy = new IEigenDAServiceManager.QuorumBlobParam[](2);
        for (uint i = 0; i < blobHeader.quorumBlobParams.length; i++) {
            quorumBlobParamsCopy[i].quorumNumber = blobHeader.quorumBlobParams[i].quorumNumber;
            quorumBlobParamsCopy[i].adversaryThresholdPercentage = blobHeader.quorumBlobParams[i].adversaryThresholdPercentage;
            quorumBlobParamsCopy[i].quorumThresholdPercentage = blobHeader.quorumBlobParams[i].quorumThresholdPercentage;
        }
        
        stdstore
            .target(address(mockRollup))
            .sig("quorumBlobParamsHash()")
            .checked_write(keccak256(abi.encode(quorumBlobParamsCopy)));

        //post commitment
        vm.prank(alice);
        mockRollup.postCommitment(blobHeader, blobVerificationProof);

        //challenge commitment
        vm.prank(bob);
        mockRollup.challengeCommitment(block.timestamp, illegalPoint, illegalProof);

        //check that kzg proof was verified
        assertEq(mockRollup.blacklist(alice), true);
        assertEq(mockRollup.validators(alice), false);
        assertEq(bob.balance, 1 ether);
    }

    function _getIllegalCommitment() internal view returns (BN254.G1Point memory illegalCommitment) {
        illegalCommitment = s0.scalar_mul(1).plus(s1.scalar_mul(1)).plus(s2.scalar_mul(1)).plus(s3.scalar_mul(1)).plus(s4.scalar_mul(1));
    }

    function _getCommitment(uint256 pseudoRandomNumber) internal returns (IEigenDAServiceManager.BlobHeader memory, EigenDABlobUtils.BlobVerificationProof memory){
        uint256 numQuorumBlobParams = 2;
        IEigenDAServiceManager.BlobHeader[] memory blobHeader = new IEigenDAServiceManager.BlobHeader[](2);
        blobHeader[0] = _generateBlobHeader(pseudoRandomNumber, numQuorumBlobParams, defaultCodingRatioPercentage);
        uint256 anotherPseudoRandomNumber = uint256(keccak256(abi.encodePacked(pseudoRandomNumber)));
        blobHeader[1] = _generateBlobHeader(anotherPseudoRandomNumber, numQuorumBlobParams, defaultCodingRatioPercentage);

        IEigenDAServiceManager.BatchHeader memory batchHeader;
        bytes memory firstBlobHash = abi.encodePacked(blobHeader[0].hashBlobHeader());
        bytes memory secondBlobHash = abi.encodePacked(blobHeader[1].hashBlobHeader());
        batchHeader.blobHeadersRoot = keccak256(abi.encodePacked(keccak256(firstBlobHash), keccak256(secondBlobHash)));
        // add dummy quorum numbers and quorum threshold percentages making sure quorumThresholdPercentage = adversaryThresholdPercentage + defaultCodingRatioPercentage
        for (uint i = 0; i < blobHeader[1].quorumBlobParams.length; i++) {
            batchHeader.quorumNumbers = abi.encodePacked(batchHeader.quorumNumbers, blobHeader[1].quorumBlobParams[i].quorumNumber);
            batchHeader.quorumThresholdPercentages = abi.encodePacked(batchHeader.quorumThresholdPercentages, blobHeader[1].quorumBlobParams[i].adversaryThresholdPercentage + defaultCodingRatioPercentage);
        }
        batchHeader.referenceBlockNumber = uint32(block.number);

        // add dummy batch metadata
        IEigenDAServiceManager.BatchMetadata memory batchMetadata;
        batchMetadata.batchHeader = batchHeader;
        batchMetadata.signatoryRecordHash = keccak256(abi.encodePacked("signatoryRecordHash"));
        batchMetadata.fee = 100;
        batchMetadata.confirmationBlockNumber = defaultConfirmationBlockNumber;

        stdstore
            .target(address(eigenDAServiceManager))
            .sig("batchIdToBatchMetadataHash(uint32)")
            .with_key(defaultBatchId)
            .checked_write(batchMetadata.hashBatchMetadata());

        EigenDABlobUtils.BlobVerificationProof memory blobVerificationProof;
        blobVerificationProof.batchId = defaultBatchId;
        blobVerificationProof.batchMetadata = batchMetadata;
        blobVerificationProof.inclusionProof = abi.encodePacked(keccak256(firstBlobHash));
        blobVerificationProof.blobIndex = 1;
        blobVerificationProof.quorumThresholdIndexes = new bytes(batchHeader.quorumNumbers.length);
        for (uint i = 0; i < batchHeader.quorumNumbers.length; i++) {
            blobVerificationProof.quorumThresholdIndexes[i] = bytes1(uint8(i));
        }

        return (blobHeader[1], blobVerificationProof);
    }

    function _generateBlobHeader(uint256 pseudoRandomNumber, uint256 numQuorumsBlobParams, uint8 codingRatioPercentage) internal returns (IEigenDAServiceManager.BlobHeader memory) {
        if(pseudoRandomNumber == 0) {
            pseudoRandomNumber = 1;
        }

        IEigenDAServiceManager.BlobHeader memory blobHeader;
        blobHeader.commitment = _getIllegalCommitment();

        blobHeader.dataLength = uint32(uint256(keccak256(abi.encodePacked(pseudoRandomNumber, "blobHeader.dataLength"))));

        blobHeader.quorumBlobParams = new IEigenDAServiceManager.QuorumBlobParam[](numQuorumsBlobParams);
        for (uint i = 0; i < numQuorumsBlobParams; i++) {
            blobHeader.quorumBlobParams[i].quorumNumber = uint8(uint256(keccak256(abi.encodePacked(pseudoRandomNumber, "blobHeader.quorumBlobParams[i].quorumNumber", i)))) % 192;
            // make sure it isn't already used
            while(quorumNumbersUsed[blobHeader.quorumBlobParams[i].quorumNumber]) {
                blobHeader.quorumBlobParams[i].quorumNumber = uint8(uint256(blobHeader.quorumBlobParams[i].quorumNumber) + 1) % 192;
            }
            quorumNumbersUsed[blobHeader.quorumBlobParams[i].quorumNumber] = true;
            blobHeader.quorumBlobParams[i].adversaryThresholdPercentage = uint8(uint256(keccak256(abi.encodePacked(pseudoRandomNumber, "blobHeader.quorumBlobParams[i].adversaryThresholdPercentage", i)))) % 100;
            // make the adversaryRatioPercentage at most 100 - codingRatioPercentage
            uint256 j = uint256(keccak256(abi.encodePacked(pseudoRandomNumber, "blobHeader.quorumBlobParams[i].adversaryThresholdPercentage nonce", i)));
            while(blobHeader.quorumBlobParams[i].adversaryThresholdPercentage > 100 - codingRatioPercentage) {
                blobHeader.quorumBlobParams[i].adversaryThresholdPercentage = uint8(uint256(keccak256(abi.encodePacked(pseudoRandomNumber, "blobHeader.quorumBlobParams[i].adversaryThresholdPercentage", j)))) % 100;
                j++;
            }
            blobHeader.quorumBlobParams[i].chunkLength = uint32(uint256(keccak256(abi.encodePacked(pseudoRandomNumber, "blobHeader.quorumBlobParams[i].chunkLength", i))));
            blobHeader.quorumBlobParams[i].quorumThresholdPercentage = blobHeader.quorumBlobParams[i].adversaryThresholdPercentage + 1;
        }
        // mark all quorum numbers as unused
        for (uint i = 0; i < numQuorumsBlobParams; i++) {
            quorumNumbersUsed[blobHeader.quorumBlobParams[i].quorumNumber] = false;
        }

        return blobHeader;
    }

}