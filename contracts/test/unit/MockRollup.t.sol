// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.9;

import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import {BLSMockAVSDeployer} from "../../lib/eigenlayer-middleware/test/utils/BLSMockAVSDeployer.sol";
import {MockRollup} from "../../src/rollup/MockRollup.sol";
import {EigenDAHasher} from "../../src/libraries/EigenDAHasher.sol";
import {EigenDAServiceManager, IRewardsCoordinator} from "../../src/core/EigenDAServiceManager.sol";
import {IEigenDAServiceManager} from "../../src/interfaces/IEigenDAServiceManager.sol";
import {EigenDARollupUtils} from "../../src/libraries/EigenDARollupUtils.sol";
import {BN254} from "eigenlayer-middleware/libraries/BN254.sol";
import "../../src/interfaces/IEigenDAStructs.sol";
import {IEigenDAThresholdRegistry} from "../../src/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDARelayRegistry} from "../../src/interfaces/IEigenDARelayRegistry.sol";
import "forge-std/StdStorage.sol";

contract MockRollupTest is BLSMockAVSDeployer {
    using stdStorage for StdStorage;
    using BN254 for BN254.G1Point;
    using EigenDAHasher for BatchHeader;
    using EigenDAHasher for ReducedBatchHeader;
    using EigenDAHasher for BlobHeader;
    using EigenDAHasher for BatchMetadata;

    EigenDAServiceManager eigenDAServiceManager;
    EigenDAServiceManager eigenDAServiceManagerImplementation;

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

    function setUp() public {
        _setUpBLSMockAVSDeployer();

        eigenDAServiceManagerImplementation = new EigenDAServiceManager(
            avsDirectory,
            rewardsCoordinator,
            registryCoordinator,
            stakeRegistry,
            IEigenDAThresholdRegistry(address(0)),
            IEigenDARelayRegistry(address(0))
        );

        address[] memory confirmers = new address[](1);
        confirmers[0] = registryCoordinatorOwner;

        // Third, upgrade the proxy contracts to use the correct implementation contracts and initialize them.
        eigenDAServiceManager = EigenDAServiceManager(
            address(
                new TransparentUpgradeableProxy(
                    address(eigenDAServiceManagerImplementation),
                    address(proxyAdmin),
                    abi.encodeWithSelector(
                        EigenDAServiceManager.initialize.selector,
                        pauserRegistry,
                        0,
                        registryCoordinatorOwner,
                        confirmers,
                        registryCoordinatorOwner
                    )
                )
            )
        );

        mockRollup = new MockRollup(eigenDAServiceManager, s1);

        //hardcode g2 proof
        illegalProof.X[1] = 11151623676041303181597631684634074376466382703418354161831688442589830350329;
        illegalProof.X[0] = 21587740443732524623985464356760343072434825248946003815467233999912459579351;
        illegalProof.Y[1] = 4222041728992406478862708226745479381252734858741080790666424175645694456140;
        illegalProof.Y[0] = 17511259870083276759899704237100059449000397154439723516103658719937845846446;

    }

    function xtestChallenge(uint256 pseudoRandomNumber) public {
        //get commitment with illegal value
        (BlobHeader memory blobHeader, BlobVerificationProof memory blobVerificationProof) = _getCommitment(pseudoRandomNumber);

        mockRollup.postCommitment(blobHeader, blobVerificationProof);

        bool success = mockRollup.challengeCommitment(block.timestamp, illegalPoint, illegalProof, illegalValue);
        assertEq(success, true);

        success = mockRollup.challengeCommitment(block.timestamp, illegalPoint, illegalProof, 69);
        assertEq(success, false);
    }

    function _getIllegalCommitment() internal view returns (BN254.G1Point memory illegalCommitment) {
        illegalCommitment = s0.scalar_mul(1).plus(s1.scalar_mul(1)).plus(s2.scalar_mul(1)).plus(s3.scalar_mul(1)).plus(s4.scalar_mul(1));
    }

    function _getCommitment(uint256 pseudoRandomNumber) internal returns (BlobHeader memory, BlobVerificationProof memory){
        uint256 numQuorumBlobParams = 2;
        BlobHeader[] memory blobHeader = new BlobHeader[](2);
        blobHeader[0] = _generateBlobHeader(pseudoRandomNumber, numQuorumBlobParams);
        uint256 anotherPseudoRandomNumber = uint256(keccak256(abi.encodePacked(pseudoRandomNumber)));
        blobHeader[1] = _generateBlobHeader(anotherPseudoRandomNumber, numQuorumBlobParams);

        BatchHeader memory batchHeader;
        bytes memory firstBlobHash = abi.encodePacked(blobHeader[0].hashBlobHeader());
        bytes memory secondBlobHash = abi.encodePacked(blobHeader[1].hashBlobHeader());
        batchHeader.blobHeadersRoot = keccak256(abi.encodePacked(keccak256(firstBlobHash), keccak256(secondBlobHash)));
        // add dummy quorum numbers and quorum threshold percentages making sure confirmationThresholdPercentage = adversaryThresholdPercentage + defaultCodingRatioPercentage
        for (uint i = 0; i < blobHeader[1].quorumBlobParams.length; i++) {
            batchHeader.quorumNumbers = abi.encodePacked(batchHeader.quorumNumbers, blobHeader[1].quorumBlobParams[i].quorumNumber);
            batchHeader.signedStakeForQuorums = abi.encodePacked(batchHeader.signedStakeForQuorums, blobHeader[1].quorumBlobParams[i].adversaryThresholdPercentage + defaultCodingRatioPercentage);
        }
        batchHeader.referenceBlockNumber = uint32(block.number);

        // add dummy batch metadata
        BatchMetadata memory batchMetadata;
        batchMetadata.batchHeader = batchHeader;
        batchMetadata.signatoryRecordHash = keccak256(abi.encodePacked("signatoryRecordHash"));
        batchMetadata.confirmationBlockNumber = defaultConfirmationBlockNumber;

        stdstore
            .target(address(eigenDAServiceManager))
            .sig("batchIdToBatchMetadataHash(uint32)")
            .with_key(defaultBatchId)
            .checked_write(batchMetadata.hashBatchMetadata());

        BlobVerificationProof memory blobVerificationProof;
        blobVerificationProof.batchId = defaultBatchId;
        blobVerificationProof.batchMetadata = batchMetadata;
        blobVerificationProof.inclusionProof = abi.encodePacked(keccak256(firstBlobHash));
        blobVerificationProof.blobIndex = 1;
        blobVerificationProof.quorumIndices = new bytes(batchHeader.quorumNumbers.length);
        for (uint i = 0; i < batchHeader.quorumNumbers.length; i++) {
            blobVerificationProof.quorumIndices[i] = bytes1(uint8(i));
        }

        return (blobHeader[1], blobVerificationProof);
    }

    function _generateBlobHeader(uint256 pseudoRandomNumber, uint256 numQuorumsBlobParams) internal returns (BlobHeader memory) {
        if(pseudoRandomNumber == 0) {
            pseudoRandomNumber = 1;
        }

        BlobHeader memory blobHeader;
        blobHeader.commitment = _getIllegalCommitment();

        blobHeader.dataLength = uint32(uint256(keccak256(abi.encodePacked(pseudoRandomNumber, "blobHeader.dataLength"))));

        blobHeader.quorumBlobParams = new QuorumBlobParam[](numQuorumsBlobParams);
        for (uint i = 0; i < numQuorumsBlobParams; i++) {
            if(i < 2){
                blobHeader.quorumBlobParams[i].quorumNumber = uint8(i);
            } else {
                blobHeader.quorumBlobParams[i].quorumNumber = uint8(uint256(keccak256(abi.encodePacked(pseudoRandomNumber, "blobHeader.quorumBlobParams[i].quorumNumber", i)))) % 192;

                // make sure it isn't already used
                while(quorumNumbersUsed[blobHeader.quorumBlobParams[i].quorumNumber]) {
                    blobHeader.quorumBlobParams[i].quorumNumber = uint8(uint256(blobHeader.quorumBlobParams[i].quorumNumber) + 1) % 192;
                }
                quorumNumbersUsed[blobHeader.quorumBlobParams[i].quorumNumber] = true;
            }
            blobHeader.quorumBlobParams[i].adversaryThresholdPercentage = EigenDARollupUtils.getQuorumAdversaryThreshold(eigenDAServiceManager, blobHeader.quorumBlobParams[i].quorumNumber);
            blobHeader.quorumBlobParams[i].chunkLength = uint32(uint256(keccak256(abi.encodePacked(pseudoRandomNumber, "blobHeader.quorumBlobParams[i].chunkLength", i))));
            blobHeader.quorumBlobParams[i].confirmationThresholdPercentage = blobHeader.quorumBlobParams[i].adversaryThresholdPercentage + 1;
        }
        // mark all quorum numbers as unused
        for (uint i = 0; i < numQuorumsBlobParams; i++) {
            quorumNumbersUsed[blobHeader.quorumBlobParams[i].quorumNumber] = false;
        }

        return blobHeader;
    }

    function xtestGetQuorumAdversaryThreshold () public {
        require(EigenDARollupUtils.getQuorumAdversaryThreshold(eigenDAServiceManager, 0) == 33, "getQuorumAdversaryThreshold failed");
        //require(EigenDARollupUtils.getQuorumAdversaryThreshold(eigenDAServiceManager, 1) == 33, "getQuorumAdversaryThreshold failed");
    }

}