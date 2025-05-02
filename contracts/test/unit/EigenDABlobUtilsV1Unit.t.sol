// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import "../MockEigenDADeployer.sol";
import {EigenDACertVerifierV1} from "src/periphery/cert/v1/EigenDACertVerifierV1.sol";
import {EigenDACertVerificationV1Lib as CertV1Lib} from "src/periphery/cert/v1/EigenDACertVerificationV1Lib.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {IEigenDABatchMetadataStorage} from "src/core/interfaces/IEigenDABatchMetadataStorage.sol";

contract EigenDABlobUtilsV1Unit is MockEigenDADeployer {
    using stdStorage for StdStorage;
    using BN254 for BN254.G1Point;

    EigenDACertVerifierV1 eigenDACertVerifierV1;

    function setUp() public virtual {
        _deployDA();

        eigenDACertVerifierV1 = new EigenDACertVerifierV1(
            IEigenDAThresholdRegistry(address(eigenDAServiceManager)),
            IEigenDABatchMetadataStorage(address(eigenDAServiceManager))
        );
    }

    function testVerifyBlob_TwoQuorums(uint256 pseudoRandomNumber) public {
        uint256 numQuorumBlobParams = 2;
        DATypesV1.BlobHeader[] memory blobHeader = new DATypesV1.BlobHeader[](2);
        blobHeader[0] = _generateRandomBlobHeader(pseudoRandomNumber, numQuorumBlobParams);
        uint256 anotherPseudoRandomNumber = uint256(keccak256(abi.encodePacked(pseudoRandomNumber)));
        blobHeader[1] = _generateRandomBlobHeader(anotherPseudoRandomNumber, numQuorumBlobParams);

        DATypesV1.BatchHeader memory batchHeader;
        bytes memory firstBlobHash = abi.encodePacked(CertV1Lib.hashBlobHeader(blobHeader[0]));
        bytes memory secondBlobHash = abi.encodePacked(CertV1Lib.hashBlobHeader(blobHeader[1]));
        batchHeader.blobHeadersRoot = keccak256(abi.encodePacked(keccak256(firstBlobHash), keccak256(secondBlobHash)));
        for (uint256 i = 0; i < blobHeader[1].quorumBlobParams.length; i++) {
            batchHeader.quorumNumbers =
                abi.encodePacked(batchHeader.quorumNumbers, blobHeader[1].quorumBlobParams[i].quorumNumber);
            batchHeader.signedStakeForQuorums = abi.encodePacked(
                batchHeader.signedStakeForQuorums, blobHeader[1].quorumBlobParams[i].confirmationThresholdPercentage
            );
        }
        batchHeader.referenceBlockNumber = uint32(block.number);

        // add dummy batch metadata
        DATypesV1.BatchMetadata memory batchMetadata;
        batchMetadata.batchHeader = batchHeader;
        batchMetadata.signatoryRecordHash = keccak256(abi.encodePacked("signatoryRecordHash"));
        batchMetadata.confirmationBlockNumber = defaultConfirmationBlockNumber;

        stdstore.target(address(eigenDAServiceManager)).sig("batchIdToBatchMetadataHash(uint32)").with_key(
            defaultBatchId
        ).checked_write(CertV1Lib.hashBatchMetadata(batchMetadata));

        DATypesV1.BlobVerificationProof memory blobVerificationProof;
        blobVerificationProof.batchId = defaultBatchId;
        blobVerificationProof.batchMetadata = batchMetadata;
        blobVerificationProof.inclusionProof = abi.encodePacked(keccak256(firstBlobHash));
        blobVerificationProof.blobIndex = 1;
        blobVerificationProof.quorumIndices = new bytes(batchHeader.quorumNumbers.length);
        for (uint256 i = 0; i < batchHeader.quorumNumbers.length; i++) {
            blobVerificationProof.quorumIndices[i] = bytes1(uint8(i));
        }

        uint256 gasBefore = gasleft();
        eigenDACertVerifierV1.verifyDACertV1(blobHeader[1], blobVerificationProof);
        uint256 gasAfter = gasleft();
        emit log_named_uint("gas used", gasBefore - gasAfter);
    }

    function testVerifyBlobs_TwoBlobs(uint256 pseudoRandomNumber) public {
        uint256 numQuorumBlobParams = 2;
        DATypesV1.BlobHeader[] memory blobHeader = new DATypesV1.BlobHeader[](2);
        blobHeader[0] = _generateRandomBlobHeader(pseudoRandomNumber, numQuorumBlobParams);
        uint256 anotherPseudoRandomNumber = uint256(keccak256(abi.encodePacked(pseudoRandomNumber)));
        blobHeader[1] = _generateRandomBlobHeader(anotherPseudoRandomNumber, numQuorumBlobParams);
        DATypesV1.BatchHeader memory batchHeader;
        bytes memory firstBlobHash = abi.encodePacked(CertV1Lib.hashBlobHeader(blobHeader[0]));
        bytes memory secondBlobHash = abi.encodePacked(CertV1Lib.hashBlobHeader(blobHeader[1]));
        batchHeader.blobHeadersRoot = keccak256(abi.encodePacked(keccak256(firstBlobHash), keccak256(secondBlobHash)));
        // add dummy quorum numbers and quorum threshold percentages making sure confirmationThresholdPercentage = adversaryThresholdPercentage + defaultCodingRatioPercentage
        for (uint256 i = 0; i < blobHeader[1].quorumBlobParams.length; i++) {
            batchHeader.quorumNumbers =
                abi.encodePacked(batchHeader.quorumNumbers, blobHeader[1].quorumBlobParams[i].quorumNumber);
            batchHeader.signedStakeForQuorums = abi.encodePacked(
                batchHeader.signedStakeForQuorums, blobHeader[1].quorumBlobParams[i].confirmationThresholdPercentage
            );
        }
        batchHeader.referenceBlockNumber = uint32(block.number);
        // add dummy batch metadata
        DATypesV1.BatchMetadata memory batchMetadata;
        batchMetadata.batchHeader = batchHeader;
        batchMetadata.signatoryRecordHash = keccak256(abi.encodePacked("signatoryRecordHash"));
        batchMetadata.confirmationBlockNumber = defaultConfirmationBlockNumber;
        stdstore.target(address(eigenDAServiceManager)).sig("batchIdToBatchMetadataHash(uint32)").with_key(
            defaultBatchId
        ).checked_write(CertV1Lib.hashBatchMetadata(batchMetadata));
        DATypesV1.BlobVerificationProof[] memory blobVerificationProofs = new DATypesV1.BlobVerificationProof[](2);
        blobVerificationProofs[0].batchId = defaultBatchId;
        blobVerificationProofs[1].batchId = defaultBatchId;
        blobVerificationProofs[0].batchMetadata = batchMetadata;
        blobVerificationProofs[1].batchMetadata = batchMetadata;
        blobVerificationProofs[0].inclusionProof = abi.encodePacked(keccak256(secondBlobHash));
        blobVerificationProofs[1].inclusionProof = abi.encodePacked(keccak256(firstBlobHash));
        blobVerificationProofs[0].blobIndex = 0;
        blobVerificationProofs[1].blobIndex = 1;
        blobVerificationProofs[0].quorumIndices = new bytes(batchHeader.quorumNumbers.length);
        blobVerificationProofs[1].quorumIndices = new bytes(batchHeader.quorumNumbers.length);
        for (uint256 i = 0; i < batchHeader.quorumNumbers.length; i++) {
            blobVerificationProofs[0].quorumIndices[i] = bytes1(uint8(i));
            blobVerificationProofs[1].quorumIndices[i] = bytes1(uint8(i));
        }
        uint256 gasBefore = gasleft();
        eigenDACertVerifierV1.verifyDACertsV1(blobHeader, blobVerificationProofs);
        uint256 gasAfter = gasleft();
        emit log_named_uint("gas used", gasBefore - gasAfter);
    }

    function testVerifyBlob_InvalidMetadataHash(uint256 pseudoRandomNumber) public {
        uint256 numQuorumBlobParams = pseudoRandomNumber % 192;
        DATypesV1.BlobHeader[] memory blobHeader = new DATypesV1.BlobHeader[](2);
        blobHeader[0] = _generateRandomBlobHeader(pseudoRandomNumber, numQuorumBlobParams);
        uint256 anotherPseudoRandomNumber = uint256(keccak256(abi.encodePacked(pseudoRandomNumber)));
        blobHeader[1] = _generateRandomBlobHeader(anotherPseudoRandomNumber, numQuorumBlobParams);

        DATypesV1.BlobVerificationProof memory blobVerificationProof;
        blobVerificationProof.batchId = defaultBatchId;

        cheats.expectRevert(
            "EigenDACertVerificationV1Lib._verifyDACertForQuorums: batchMetadata does not match stored metadata"
        );
        eigenDACertVerifierV1.verifyDACertV1(blobHeader[1], blobVerificationProof);
    }

    function testVerifyBlob_InvalidMerkleProof(uint256 pseudoRandomNumber) public {
        uint256 numQuorumBlobParams = pseudoRandomNumber % 192;
        DATypesV1.BlobHeader[] memory blobHeader = new DATypesV1.BlobHeader[](2);
        blobHeader[0] = _generateRandomBlobHeader(pseudoRandomNumber, numQuorumBlobParams);
        uint256 anotherPseudoRandomNumber = uint256(keccak256(abi.encodePacked(pseudoRandomNumber)));
        blobHeader[1] = _generateRandomBlobHeader(anotherPseudoRandomNumber, numQuorumBlobParams);

        // add dummy batch metadata
        DATypesV1.BatchMetadata memory batchMetadata;

        stdstore.target(address(eigenDAServiceManager)).sig("batchIdToBatchMetadataHash(uint32)").with_key(
            defaultBatchId
        ).checked_write(CertV1Lib.hashBatchMetadata(batchMetadata));

        DATypesV1.BlobVerificationProof memory blobVerificationProof;
        blobVerificationProof.batchId = defaultBatchId;
        blobVerificationProof.batchMetadata = batchMetadata;
        blobVerificationProof.inclusionProof = abi.encodePacked(bytes32(0));
        blobVerificationProof.blobIndex = 1;

        cheats.expectRevert("EigenDACertVerificationV1Lib._verifyDACertForQuorums: inclusion proof is invalid");
        eigenDACertVerifierV1.verifyDACertV1(blobHeader[1], blobVerificationProof);
    }

    function testVerifyBlob_RequiredQuorumsNotMet(uint256 pseudoRandomNumber) public {
        uint256 numQuorumBlobParams = 1;
        DATypesV1.BlobHeader[] memory blobHeader = new DATypesV1.BlobHeader[](2);
        blobHeader[0] = _generateRandomBlobHeader(pseudoRandomNumber, numQuorumBlobParams);
        uint256 anotherPseudoRandomNumber = uint256(keccak256(abi.encodePacked(pseudoRandomNumber)));
        blobHeader[1] = _generateRandomBlobHeader(anotherPseudoRandomNumber, numQuorumBlobParams);

        DATypesV1.BatchHeader memory batchHeader;
        bytes memory firstBlobHash = abi.encodePacked(CertV1Lib.hashBlobHeader(blobHeader[0]));
        bytes memory secondBlobHash = abi.encodePacked(CertV1Lib.hashBlobHeader(blobHeader[1]));
        batchHeader.blobHeadersRoot = keccak256(abi.encodePacked(keccak256(firstBlobHash), keccak256(secondBlobHash)));
        for (uint256 i = 0; i < blobHeader[1].quorumBlobParams.length; i++) {
            batchHeader.quorumNumbers =
                abi.encodePacked(batchHeader.quorumNumbers, blobHeader[1].quorumBlobParams[i].quorumNumber);
            batchHeader.signedStakeForQuorums = abi.encodePacked(
                batchHeader.signedStakeForQuorums, blobHeader[1].quorumBlobParams[i].confirmationThresholdPercentage
            );
        }
        batchHeader.referenceBlockNumber = uint32(block.number);

        // add dummy batch metadata
        DATypesV1.BatchMetadata memory batchMetadata;
        batchMetadata.batchHeader = batchHeader;
        batchMetadata.signatoryRecordHash = keccak256(abi.encodePacked("signatoryRecordHash"));
        batchMetadata.confirmationBlockNumber = defaultConfirmationBlockNumber;

        stdstore.target(address(eigenDAServiceManager)).sig("batchIdToBatchMetadataHash(uint32)").with_key(
            defaultBatchId
        ).checked_write(CertV1Lib.hashBatchMetadata(batchMetadata));

        DATypesV1.BlobVerificationProof memory blobVerificationProof;
        blobVerificationProof.batchId = defaultBatchId;
        blobVerificationProof.batchMetadata = batchMetadata;
        blobVerificationProof.inclusionProof = abi.encodePacked(keccak256(firstBlobHash));
        blobVerificationProof.blobIndex = 1;
        blobVerificationProof.quorumIndices = new bytes(batchHeader.quorumNumbers.length);
        for (uint256 i = 0; i < batchHeader.quorumNumbers.length; i++) {
            blobVerificationProof.quorumIndices[i] = bytes1(uint8(i));
        }

        cheats.expectRevert(
            "EigenDACertVerificationV1Lib._verifyDACertForQuorums: required quorums are not a subset of the confirmed quorums"
        );
        eigenDACertVerifierV1.verifyDACertV1(blobHeader[1], blobVerificationProof);
    }

    function testVerifyBlob_QuorumNumberMismatch(uint256 pseudoRandomNumber) public {
        uint256 numQuorumBlobParams = 2;
        DATypesV1.BlobHeader[] memory blobHeader = new DATypesV1.BlobHeader[](2);
        blobHeader[0] = _generateRandomBlobHeader(pseudoRandomNumber, numQuorumBlobParams);
        uint256 anotherPseudoRandomNumber = uint256(keccak256(abi.encodePacked(pseudoRandomNumber)));
        blobHeader[1] = _generateRandomBlobHeader(anotherPseudoRandomNumber, numQuorumBlobParams);

        DATypesV1.BatchHeader memory batchHeader;
        bytes memory firstBlobHash = abi.encodePacked(CertV1Lib.hashBlobHeader(blobHeader[0]));
        bytes memory secondBlobHash = abi.encodePacked(CertV1Lib.hashBlobHeader(blobHeader[1]));
        batchHeader.blobHeadersRoot = keccak256(abi.encodePacked(keccak256(firstBlobHash), keccak256(secondBlobHash)));
        for (uint256 i = 0; i < blobHeader[1].quorumBlobParams.length; i++) {
            batchHeader.quorumNumbers =
                abi.encodePacked(batchHeader.quorumNumbers, blobHeader[1].quorumBlobParams[i].quorumNumber);
            batchHeader.signedStakeForQuorums = abi.encodePacked(
                batchHeader.signedStakeForQuorums, blobHeader[1].quorumBlobParams[i].confirmationThresholdPercentage
            );
        }
        batchHeader.referenceBlockNumber = uint32(block.number);

        // add dummy batch metadata
        DATypesV1.BatchMetadata memory batchMetadata;
        batchMetadata.batchHeader = batchHeader;
        batchMetadata.signatoryRecordHash = keccak256(abi.encodePacked("signatoryRecordHash"));
        batchMetadata.confirmationBlockNumber = defaultConfirmationBlockNumber;

        stdstore.target(address(eigenDAServiceManager)).sig("batchIdToBatchMetadataHash(uint32)").with_key(
            defaultBatchId
        ).checked_write(CertV1Lib.hashBatchMetadata(batchMetadata));

        DATypesV1.BlobVerificationProof memory blobVerificationProof;
        blobVerificationProof.batchId = defaultBatchId;
        blobVerificationProof.batchMetadata = batchMetadata;
        blobVerificationProof.inclusionProof = abi.encodePacked(keccak256(firstBlobHash));
        blobVerificationProof.blobIndex = 1;
        blobVerificationProof.quorumIndices = new bytes(batchHeader.quorumNumbers.length);
        for (uint256 i = 0; i < batchHeader.quorumNumbers.length; i++) {
            // implant the incorrect quorumNumbers here
            blobVerificationProof.quorumIndices[i] = bytes1(uint8(batchHeader.quorumNumbers.length - 1 - i));
        }

        cheats.expectRevert("EigenDACertVerificationV1Lib._verifyDACertForQuorums: quorumNumber does not match");
        eigenDACertVerifierV1.verifyDACertV1(blobHeader[1], blobVerificationProof);
    }

    function testVerifyBlob_QuorumThresholdNotMet(uint256 pseudoRandomNumber) public {
        uint256 numQuorumBlobParams = 2;
        DATypesV1.BlobHeader[] memory blobHeader = new DATypesV1.BlobHeader[](2);
        blobHeader[0] = _generateRandomBlobHeader(pseudoRandomNumber, numQuorumBlobParams);
        uint256 anotherPseudoRandomNumber = uint256(keccak256(abi.encodePacked(pseudoRandomNumber)));
        blobHeader[1] = _generateRandomBlobHeader(anotherPseudoRandomNumber, numQuorumBlobParams);

        DATypesV1.BatchHeader memory batchHeader;
        bytes memory firstBlobHash = abi.encodePacked(CertV1Lib.hashBlobHeader(blobHeader[0]));
        bytes memory secondBlobHash = abi.encodePacked(CertV1Lib.hashBlobHeader(blobHeader[1]));
        batchHeader.blobHeadersRoot = keccak256(abi.encodePacked(keccak256(firstBlobHash), keccak256(secondBlobHash)));
        // add dummy quorum numbers and quorum threshold percentages making sure confirmationThresholdPercentage = 100
        for (uint256 i = 0; i < blobHeader[1].quorumBlobParams.length; i++) {
            batchHeader.quorumNumbers =
                abi.encodePacked(batchHeader.quorumNumbers, blobHeader[1].quorumBlobParams[i].quorumNumber);
            batchHeader.signedStakeForQuorums = abi.encodePacked(
                batchHeader.signedStakeForQuorums, blobHeader[1].quorumBlobParams[i].confirmationThresholdPercentage - 1
            );
        }
        batchHeader.referenceBlockNumber = uint32(block.number);

        // add dummy batch metadata
        DATypesV1.BatchMetadata memory batchMetadata;
        batchMetadata.batchHeader = batchHeader;
        batchMetadata.signatoryRecordHash = keccak256(abi.encodePacked("signatoryRecordHash"));
        batchMetadata.confirmationBlockNumber = defaultConfirmationBlockNumber;

        stdstore.target(address(eigenDAServiceManager)).sig("batchIdToBatchMetadataHash(uint32)").with_key(
            defaultBatchId
        ).checked_write(CertV1Lib.hashBatchMetadata(batchMetadata));

        DATypesV1.BlobVerificationProof memory blobVerificationProof;
        blobVerificationProof.batchId = defaultBatchId;
        blobVerificationProof.batchMetadata = batchMetadata;
        blobVerificationProof.inclusionProof = abi.encodePacked(keccak256(firstBlobHash));
        blobVerificationProof.blobIndex = 1;
        blobVerificationProof.quorumIndices = new bytes(batchHeader.quorumNumbers.length);
        for (uint256 i = 0; i < batchHeader.quorumNumbers.length; i++) {
            // implant the incorrect quorumNumbers here
            blobVerificationProof.quorumIndices[i] = bytes1(uint8(i));
        }

        cheats.expectRevert(
            "EigenDACertVerificationV1Lib._verifyDACertForQuorums: confirmationThresholdPercentage is not met"
        );
        eigenDACertVerifierV1.verifyDACertV1(blobHeader[1], blobVerificationProof);
    }

    function testThresholds() public view {
        require(
            eigenDACertVerifierV1.getQuorumAdversaryThresholdPercentage(0) == 33,
            "getQuorumAdversaryThresholdPercentage failed"
        );
        require(
            eigenDACertVerifierV1.getQuorumAdversaryThresholdPercentage(1) == 33,
            "getQuorumAdversaryThresholdPercentage failed"
        );
        require(
            eigenDACertVerifierV1.getQuorumAdversaryThresholdPercentage(2) == 33,
            "getQuorumAdversaryThresholdPercentage failed"
        );
        require(
            eigenDACertVerifierV1.getQuorumConfirmationThresholdPercentage(0) == 55,
            "getQuorumConfirmationThresholdPercentage failed"
        );
        require(
            eigenDACertVerifierV1.getQuorumConfirmationThresholdPercentage(1) == 55,
            "getQuorumConfirmationThresholdPercentage failed"
        );
        require(
            eigenDACertVerifierV1.getQuorumConfirmationThresholdPercentage(2) == 55,
            "getQuorumConfirmationThresholdPercentage failed"
        );
        require(eigenDACertVerifierV1.getIsQuorumRequired(0) == true, "getIsQuorumRequired failed");
        require(eigenDACertVerifierV1.getIsQuorumRequired(1) == true, "getIsQuorumRequired failed");
        require(eigenDACertVerifierV1.getIsQuorumRequired(2) == false, "getIsQuorumRequired failed");
    }
}
