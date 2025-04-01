// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {EigenDACertVerifierRouter} from "src/periphery/EigenDACertVerifierRouter.sol";
import {IEigenDACertVerifierBase} from "src/interfaces/IEigenDACertVerifier.sol";
import {TransparentUpgradeableProxy} from
    "openzeppelin-contracts/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {ProxyAdmin} from "openzeppelin-contracts/contracts/proxy/transparent/ProxyAdmin.sol";

import {
    BlobHeader,
    BlobVerificationProof,
    BatchHeaderV2,
    BlobInclusionInfo,
    NonSignerStakesAndSignature,
    SignedBatch
} from "src/interfaces/IEigenDAStructs.sol";

import "forge-std/Test.sol";

contract TestEigenDACertVerifierRouter is Test {
    EigenDACertVerifierRouter certVerifierRouter;
    address routerImpl;
    EigenDACertVerifierMock[] certVerifierMocks;

    uint256 numCertVerifiers = 3;
    uint32 startBlockNumber = 100;

    event CertVerifierAdded(uint32 indexed referenceBlockNumber, address indexed certVerifier);

    function setUp() public {
        vm.roll(startBlockNumber);
        routerImpl = address(new EigenDACertVerifierRouter());
        ProxyAdmin proxyAdmin = new ProxyAdmin();
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            routerImpl, address(proxyAdmin), abi.encodeCall(EigenDACertVerifierRouter.initialize, address(this))
        );
        certVerifierRouter = EigenDACertVerifierRouter(address(proxy));
        for (uint256 i; i < numCertVerifiers; i++) {
            certVerifierMocks.push(new EigenDACertVerifierMock());
        }
    }

    function testAddCertVerifier() public {
        uint32 referenceBlockNumber = uint32(startBlockNumber + 1);
        vm.expectEmit(address(certVerifierRouter));
        emit CertVerifierAdded(referenceBlockNumber, address(certVerifierMocks[0]));
        certVerifierRouter.addCertVerifier(referenceBlockNumber, address(certVerifierMocks[0]));
        assertEq(address(certVerifierRouter.certVerifiers(referenceBlockNumber)), address(certVerifierMocks[0]));
        assertEq(certVerifierRouter.certVerifierRBNs(0), referenceBlockNumber);
    }

    function testAddCertVerifierReverts() public {
        vm.roll(startBlockNumber);
        vm.expectRevert("Reference block number must be in the future");
        certVerifierRouter.addCertVerifier(startBlockNumber - 1, address(certVerifierMocks[0]));

        certVerifierRouter.addCertVerifier(startBlockNumber + 1, address(certVerifierMocks[0]));
        vm.expectRevert("Reference block number must be greater than the last registered RBN");
        certVerifierRouter.addCertVerifier(101, address(certVerifierMocks[0]));
    }

    function setupMultipleCertVerifiers() internal {
        for (uint32 i; i < numCertVerifiers; i++) {
            certVerifierRouter.addCertVerifier(startBlockNumber + i + 1, address(certVerifierMocks[i]));
        }
    }

    function testMultipleCertVerifiers() public {
        setupMultipleCertVerifiers();
        for (uint32 i; i < numCertVerifiers; i++) {
            assertEq(address(certVerifierRouter.certVerifiers(startBlockNumber + i + 1)), address(certVerifierMocks[i]));
            assertEq(certVerifierRouter.certVerifierRBNs(i), startBlockNumber + i + 1);
        }
    }

    /// @dev just makes tests less verbose as we are just trying to test routing.
    function verifyRoutingHelper(uint32 rbn) internal view {
        BlobHeader memory blobHeader;
        BlobVerificationProof memory blobVerificationProof;
        blobVerificationProof.batchMetadata.batchHeader.referenceBlockNumber = rbn;
        certVerifierRouter.verifyDACertV1(blobHeader, blobVerificationProof);
    }

    function testRouting(uint32 blockIncrement) public {
        blockIncrement = uint32(bound(blockIncrement, 0, numCertVerifiers - 1));
        setupMultipleCertVerifiers();
        certVerifierMocks[blockIncrement].setRevertOnCall(true);

        vm.expectRevert("Mock: verifyDACertV1 reverted");
        verifyRoutingHelper(startBlockNumber + blockIncrement + 1);
    }

    function testRoutingFails() public {
        vm.expectRevert("No cert verifiers available");
        verifyRoutingHelper(startBlockNumber);

        setupMultipleCertVerifiers();
        vm.expectRevert("No cert verifier found for the given reference block number");
        verifyRoutingHelper(startBlockNumber);
    }

    function testVerifyDACertV1() public {
        BlobHeader memory blobHeader;
        BlobVerificationProof memory blobVerificationProof;
        blobVerificationProof.batchMetadata.batchHeader.referenceBlockNumber = uint32(block.number + 1);
        certVerifierRouter.addCertVerifier(
            blobVerificationProof.batchMetadata.batchHeader.referenceBlockNumber, address(certVerifierMocks[0])
        );
        certVerifierRouter.verifyDACertV1(blobHeader, blobVerificationProof);
    }

    function testVerifyDACertsV1() public {
        BlobHeader[] memory blobHeaders = new BlobHeader[](2);
        BlobVerificationProof[] memory blobVerificationProofs = new BlobVerificationProof[](2);
        blobVerificationProofs[0].batchMetadata.batchHeader.referenceBlockNumber = uint32(block.number + 1);
        blobVerificationProofs[1].batchMetadata.batchHeader.referenceBlockNumber = uint32(block.number + 2);
        certVerifierRouter.addCertVerifier(
            blobVerificationProofs[0].batchMetadata.batchHeader.referenceBlockNumber, address(certVerifierMocks[0])
        );
        certVerifierRouter.addCertVerifier(
            blobVerificationProofs[1].batchMetadata.batchHeader.referenceBlockNumber, address(certVerifierMocks[1])
        );
        certVerifierRouter.verifyDACertsV1(blobHeaders, blobVerificationProofs);
    }

    function testVerifyDACertsV1RevertsForSingleRouteFailure(uint256 x) public {
        BlobHeader[] memory blobHeaders = new BlobHeader[](2);
        BlobVerificationProof[] memory blobVerificationProofs = new BlobVerificationProof[](2);
        blobVerificationProofs[0].batchMetadata.batchHeader.referenceBlockNumber = uint32(block.number + 1);
        blobVerificationProofs[1].batchMetadata.batchHeader.referenceBlockNumber = uint32(block.number + 2);
        certVerifierRouter.addCertVerifier(
            blobVerificationProofs[0].batchMetadata.batchHeader.referenceBlockNumber, address(certVerifierMocks[0])
        );
        certVerifierRouter.addCertVerifier(
            blobVerificationProofs[1].batchMetadata.batchHeader.referenceBlockNumber, address(certVerifierMocks[1])
        );
        certVerifierMocks[x % 2].setRevertOnCall(true);
        vm.expectRevert("Mock: verifyDACertV1 reverted");
        certVerifierRouter.verifyDACertsV1(blobHeaders, blobVerificationProofs);
    }

    function testVerifyDACertV2() public {
        BatchHeaderV2 memory batchHeader;
        BlobInclusionInfo memory blobInclusionInfo;
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature;
        bytes memory signedQuorumNumbers = new bytes(0);
        batchHeader.referenceBlockNumber = uint32(block.number + 1);
        certVerifierRouter.addCertVerifier(batchHeader.referenceBlockNumber, address(certVerifierMocks[0]));
        certVerifierRouter.verifyDACertV2(
            batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers
        );
    }

    function testVerifyDACertV2FromSignedBatch() public {
        SignedBatch memory signedBatch;
        BlobInclusionInfo memory blobInclusionInfo;
        signedBatch.batchHeader.referenceBlockNumber = uint32(block.number + 1);
        certVerifierRouter.addCertVerifier(signedBatch.batchHeader.referenceBlockNumber, address(certVerifierMocks[0]));
        certVerifierRouter.verifyDACertV2FromSignedBatch(signedBatch, blobInclusionInfo);
    }

    function testVerifyDACertV2ForZKProof() public {
        BatchHeaderV2 memory batchHeader;
        BlobInclusionInfo memory blobInclusionInfo;
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature;
        bytes memory signedQuorumNumbers = new bytes(0);
        batchHeader.referenceBlockNumber = uint32(block.number + 1);
        certVerifierRouter.addCertVerifier(batchHeader.referenceBlockNumber, address(certVerifierMocks[0]));
        certVerifierRouter.verifyDACertV2ForZKProof(
            batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers
        );
    }

    function testGetNonSignerStakesAndSignature() public {
        SignedBatch memory signedBatch;
        signedBatch.batchHeader.referenceBlockNumber = uint32(block.number + 1);
        certVerifierRouter.addCertVerifier(signedBatch.batchHeader.referenceBlockNumber, address(certVerifierMocks[0]));
        certVerifierRouter.getNonSignerStakesAndSignature(signedBatch);
    }
}

contract EigenDACertVerifierMock is IEigenDACertVerifierBase {
    bool revertOnCall;

    function setRevertOnCall(bool _revertOnCall) external {
        revertOnCall = _revertOnCall;
    }

    function verifyDACertV1(BlobHeader calldata, BlobVerificationProof calldata) external view override {
        require(!revertOnCall, "Mock: verifyDACertV1 reverted");
    }

    function verifyDACertsV1(BlobHeader[] calldata, BlobVerificationProof[] calldata) external view override {
        require(!revertOnCall, "Mock: verifyDACertsV1 reverted");
    }

    function verifyDACertV2(
        BatchHeaderV2 calldata,
        BlobInclusionInfo calldata,
        NonSignerStakesAndSignature calldata,
        bytes memory
    ) external view override {
        require(!revertOnCall, "Mock: verifyDACertV2 reverted");
    }

    function verifyDACertV2FromSignedBatch(SignedBatch calldata, BlobInclusionInfo calldata) external view override {
        require(!revertOnCall, "Mock: verifyDACertV2FromSignedBatch reverted");
    }

    function verifyDACertV2ForZKProof(
        BatchHeaderV2 calldata,
        BlobInclusionInfo calldata,
        NonSignerStakesAndSignature calldata,
        bytes memory
    ) external view override returns (bool success) {
        require(!revertOnCall, "Mock: verifyDACertV2ForZKProof reverted");
        return true;
    }

    function getNonSignerStakesAndSignature(SignedBatch calldata)
        external
        view
        override
        returns (NonSignerStakesAndSignature memory nonSignerStakesAndSignature)
    {
        require(!revertOnCall, "Mock: getNonSignerStakesAndSignature reverted");
        return nonSignerStakesAndSignature;
    }
}
