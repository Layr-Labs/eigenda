// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import "../MockEigenDADeployer.sol";

contract EigenDABlobUtilsV1Unit is MockEigenDADeployer {
    using stdStorage for StdStorage;
    using BN254 for BN254.G1Point;
    using EigenDAHasher for BatchHeader;
    using EigenDAHasher for ReducedBatchHeader;
    using EigenDAHasher for BlobHeader;
    using EigenDAHasher for BatchMetadata;

    function setUp() virtual public {
        _deployDA();
    }

    function testVerifyBlob_TwoQuorums(uint256 pseudoRandomNumber) public {
        uint256 numQuorumBlobParams = 2;
        BlobHeader[] memory blobHeader = new BlobHeader[](2);
        blobHeader[0] = _generateRandomBlobHeader(pseudoRandomNumber, numQuorumBlobParams);
        uint256 anotherPseudoRandomNumber = uint256(keccak256(abi.encodePacked(pseudoRandomNumber)));
        blobHeader[1] = _generateRandomBlobHeader(anotherPseudoRandomNumber, numQuorumBlobParams);

        BatchHeader memory batchHeader;
        bytes memory firstBlobHash = abi.encodePacked(blobHeader[0].hashBlobHeader());
        bytes memory secondBlobHash = abi.encodePacked(blobHeader[1].hashBlobHeader());
        batchHeader.blobHeadersRoot = keccak256(abi.encodePacked(keccak256(firstBlobHash), keccak256(secondBlobHash)));
        for (uint i = 0; i < blobHeader[1].quorumBlobParams.length; i++) {
            batchHeader.quorumNumbers = abi.encodePacked(batchHeader.quorumNumbers, blobHeader[1].quorumBlobParams[i].quorumNumber);
            batchHeader.signedStakeForQuorums = abi.encodePacked(batchHeader.signedStakeForQuorums, blobHeader[1].quorumBlobParams[i].confirmationThresholdPercentage);
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

        uint256 gasBefore = gasleft();
        eigenDABlobVerifier.verifyBlobV1(blobHeader[1], blobVerificationProof);
        uint256 gasAfter = gasleft();
        emit log_named_uint("gas used", gasBefore - gasAfter);
    }

    function testVerifyBlobs_TwoBlobs(uint256 pseudoRandomNumber) public {
        uint256 numQuorumBlobParams = 2;
        BlobHeader[] memory blobHeader = new BlobHeader[](2);
        blobHeader[0] = _generateRandomBlobHeader(pseudoRandomNumber, numQuorumBlobParams);
        uint256 anotherPseudoRandomNumber = uint256(keccak256(abi.encodePacked(pseudoRandomNumber)));
        blobHeader[1] = _generateRandomBlobHeader(anotherPseudoRandomNumber, numQuorumBlobParams);
        BatchHeader memory batchHeader;
        bytes memory firstBlobHash = abi.encodePacked(blobHeader[0].hashBlobHeader());
        bytes memory secondBlobHash = abi.encodePacked(blobHeader[1].hashBlobHeader());
        batchHeader.blobHeadersRoot = keccak256(abi.encodePacked(keccak256(firstBlobHash), keccak256(secondBlobHash)));
        // add dummy quorum numbers and quorum threshold percentages making sure confirmationThresholdPercentage = adversaryThresholdPercentage + defaultCodingRatioPercentage
        for (uint i = 0; i < blobHeader[1].quorumBlobParams.length; i++) {
            batchHeader.quorumNumbers = abi.encodePacked(batchHeader.quorumNumbers, blobHeader[1].quorumBlobParams[i].quorumNumber);
            batchHeader.signedStakeForQuorums = abi.encodePacked(batchHeader.signedStakeForQuorums, blobHeader[1].quorumBlobParams[i].confirmationThresholdPercentage);        }
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
        BlobVerificationProof[] memory blobVerificationProofs = new BlobVerificationProof[](2);
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
        for (uint i = 0; i < batchHeader.quorumNumbers.length; i++) {
            blobVerificationProofs[0].quorumIndices[i] = bytes1(uint8(i));
            blobVerificationProofs[1].quorumIndices[i] = bytes1(uint8(i));
        }
        uint256 gasBefore = gasleft();
        eigenDABlobVerifier.verifyBlobsV1(blobHeader, blobVerificationProofs);
        uint256 gasAfter = gasleft();
        emit log_named_uint("gas used", gasBefore - gasAfter);
    }

    function testVerifyBlob_InvalidMetadataHash(uint256 pseudoRandomNumber) public {
        uint256 numQuorumBlobParams = pseudoRandomNumber % 192;
        BlobHeader[] memory blobHeader = new BlobHeader[](2);
        blobHeader[0] = _generateRandomBlobHeader(pseudoRandomNumber, numQuorumBlobParams);
        uint256 anotherPseudoRandomNumber = uint256(keccak256(abi.encodePacked(pseudoRandomNumber)));
        blobHeader[1] = _generateRandomBlobHeader(anotherPseudoRandomNumber, numQuorumBlobParams);

        BlobVerificationProof memory blobVerificationProof;
        blobVerificationProof.batchId = defaultBatchId;

        cheats.expectRevert("EigenDABlobVerificationUtils._verifyBlobForQuorums: batchMetadata does not match stored metadata");
        eigenDABlobVerifier.verifyBlobV1(blobHeader[1], blobVerificationProof);
    }

    function testVerifyBlob_InvalidMerkleProof(uint256 pseudoRandomNumber) public {
        uint256 numQuorumBlobParams = pseudoRandomNumber % 192;
        BlobHeader[] memory blobHeader = new BlobHeader[](2);
        blobHeader[0] = _generateRandomBlobHeader(pseudoRandomNumber, numQuorumBlobParams);
        uint256 anotherPseudoRandomNumber = uint256(keccak256(abi.encodePacked(pseudoRandomNumber)));
        blobHeader[1] = _generateRandomBlobHeader(anotherPseudoRandomNumber, numQuorumBlobParams);

        // add dummy batch metadata
        BatchMetadata memory batchMetadata;

        stdstore
            .target(address(eigenDAServiceManager))
            .sig("batchIdToBatchMetadataHash(uint32)")
            .with_key(defaultBatchId)
            .checked_write(batchMetadata.hashBatchMetadata());

        BlobVerificationProof memory blobVerificationProof;
        blobVerificationProof.batchId = defaultBatchId;
        blobVerificationProof.batchMetadata = batchMetadata;
        blobVerificationProof.inclusionProof = abi.encodePacked(bytes32(0));        
        blobVerificationProof.blobIndex = 1;

        cheats.expectRevert("EigenDABlobVerificationUtils._verifyBlobForQuorums: inclusion proof is invalid");
        eigenDABlobVerifier.verifyBlobV1(blobHeader[1], blobVerificationProof);
    }

    function testVerifyBlob_RequiredQuorumsNotMet(uint256 pseudoRandomNumber) public {
        uint256 numQuorumBlobParams = 1;
        BlobHeader[] memory blobHeader = new BlobHeader[](2);
        blobHeader[0] = _generateRandomBlobHeader(pseudoRandomNumber, numQuorumBlobParams);
        uint256 anotherPseudoRandomNumber = uint256(keccak256(abi.encodePacked(pseudoRandomNumber)));
        blobHeader[1] = _generateRandomBlobHeader(anotherPseudoRandomNumber, numQuorumBlobParams);

        BatchHeader memory batchHeader;
        bytes memory firstBlobHash = abi.encodePacked(blobHeader[0].hashBlobHeader());
        bytes memory secondBlobHash = abi.encodePacked(blobHeader[1].hashBlobHeader());
        batchHeader.blobHeadersRoot = keccak256(abi.encodePacked(keccak256(firstBlobHash), keccak256(secondBlobHash)));
        for (uint i = 0; i < blobHeader[1].quorumBlobParams.length; i++) {
            batchHeader.quorumNumbers = abi.encodePacked(batchHeader.quorumNumbers, blobHeader[1].quorumBlobParams[i].quorumNumber);
            batchHeader.signedStakeForQuorums = abi.encodePacked(batchHeader.signedStakeForQuorums, blobHeader[1].quorumBlobParams[i].confirmationThresholdPercentage);
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

        cheats.expectRevert("EigenDABlobVerificationUtils._verifyBlobForQuorums: required quorums are not a subset of the confirmed quorums");
        eigenDABlobVerifier.verifyBlobV1(blobHeader[1], blobVerificationProof);
    }

    function testVerifyBlob_QuorumNumberMismatch(uint256 pseudoRandomNumber) public {
        uint256 numQuorumBlobParams = 2;
        BlobHeader[] memory blobHeader = new BlobHeader[](2);
        blobHeader[0] = _generateRandomBlobHeader(pseudoRandomNumber, numQuorumBlobParams);
        uint256 anotherPseudoRandomNumber = uint256(keccak256(abi.encodePacked(pseudoRandomNumber)));
        blobHeader[1] = _generateRandomBlobHeader(anotherPseudoRandomNumber, numQuorumBlobParams);

        BatchHeader memory batchHeader;
        bytes memory firstBlobHash = abi.encodePacked(blobHeader[0].hashBlobHeader());
        bytes memory secondBlobHash = abi.encodePacked(blobHeader[1].hashBlobHeader());
        batchHeader.blobHeadersRoot = keccak256(abi.encodePacked(keccak256(firstBlobHash), keccak256(secondBlobHash)));
        for (uint i = 0; i < blobHeader[1].quorumBlobParams.length; i++) {
            batchHeader.quorumNumbers = abi.encodePacked(batchHeader.quorumNumbers, blobHeader[1].quorumBlobParams[i].quorumNumber);
            batchHeader.signedStakeForQuorums = abi.encodePacked(batchHeader.signedStakeForQuorums, blobHeader[1].quorumBlobParams[i].confirmationThresholdPercentage);
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
            // implant the incorrect quorumNumbers here
            blobVerificationProof.quorumIndices[i] = bytes1(uint8(batchHeader.quorumNumbers.length - 1 - i));
        }

        cheats.expectRevert("EigenDABlobVerificationUtils._verifyBlobForQuorums: quorumNumber does not match");
        eigenDABlobVerifier.verifyBlobV1(blobHeader[1], blobVerificationProof);
    }

    function testVerifyBlob_QuorumThresholdNotMet(uint256 pseudoRandomNumber) public {
        uint256 numQuorumBlobParams = 2;
        BlobHeader[] memory blobHeader = new BlobHeader[](2);
        blobHeader[0] = _generateRandomBlobHeader(pseudoRandomNumber, numQuorumBlobParams);
        uint256 anotherPseudoRandomNumber = uint256(keccak256(abi.encodePacked(pseudoRandomNumber)));
        blobHeader[1] = _generateRandomBlobHeader(anotherPseudoRandomNumber, numQuorumBlobParams);

        BatchHeader memory batchHeader;
        bytes memory firstBlobHash = abi.encodePacked(blobHeader[0].hashBlobHeader());
        bytes memory secondBlobHash = abi.encodePacked(blobHeader[1].hashBlobHeader());
        batchHeader.blobHeadersRoot = keccak256(abi.encodePacked(keccak256(firstBlobHash), keccak256(secondBlobHash)));
        // add dummy quorum numbers and quorum threshold percentages making sure confirmationThresholdPercentage = 100
        for (uint i = 0; i < blobHeader[1].quorumBlobParams.length; i++) {
            batchHeader.quorumNumbers = abi.encodePacked(batchHeader.quorumNumbers, blobHeader[1].quorumBlobParams[i].quorumNumber);
            batchHeader.signedStakeForQuorums = abi.encodePacked(batchHeader.signedStakeForQuorums, blobHeader[1].quorumBlobParams[i].confirmationThresholdPercentage - 1);
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
            // implant the incorrect quorumNumbers here
            blobVerificationProof.quorumIndices[i] = bytes1(uint8(i));
        }

        cheats.expectRevert("EigenDABlobVerificationUtils._verifyBlobForQuorums: confirmationThresholdPercentage is not met");
        eigenDABlobVerifier.verifyBlobV1(blobHeader[1], blobVerificationProof);
    }

    function testThresholds() public {
        require(eigenDABlobVerifier.getQuorumAdversaryThresholdPercentage(0) == 33, "getQuorumAdversaryThresholdPercentage failed");
        require(eigenDABlobVerifier.getQuorumAdversaryThresholdPercentage(1) == 33, "getQuorumAdversaryThresholdPercentage failed");
        require(eigenDABlobVerifier.getQuorumAdversaryThresholdPercentage(2) == 33, "getQuorumAdversaryThresholdPercentage failed");
        require(eigenDABlobVerifier.getQuorumConfirmationThresholdPercentage(0) == 55, "getQuorumConfirmationThresholdPercentage failed");
        require(eigenDABlobVerifier.getQuorumConfirmationThresholdPercentage(1) == 55, "getQuorumConfirmationThresholdPercentage failed");
        require(eigenDABlobVerifier.getQuorumConfirmationThresholdPercentage(2) == 55, "getQuorumConfirmationThresholdPercentage failed");
        require(eigenDABlobVerifier.getIsQuorumRequired(0) == true, "getIsQuorumRequired failed");
        require(eigenDABlobVerifier.getIsQuorumRequired(1) == true, "getIsQuorumRequired failed");
        require(eigenDABlobVerifier.getIsQuorumRequired(2) == false, "getIsQuorumRequired failed");
    }    
}