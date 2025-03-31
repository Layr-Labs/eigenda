// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifierBase} from "src/interfaces/IEigenDACertVerifier.sol";
import {Ownable} from "lib/openzeppelin-contracts/contracts/access/Ownable.sol";
import "src/interfaces/IEigenDAStructs.sol";

contract CertVerifierRouter is IEigenDACertVerifierBase, Ownable {
    mapping(uint64 => IEigenDACertVerifierBase) public certVerifiers;
    uint64[] public certVerifierRBNs;

    uint256 public immutable DELAY_BLOCKS;

    event CertVerifierAdded(uint64 indexed referenceBlockNumber, address indexed certVerifier);

    constructor(uint256 delayBlocks) Ownable() {
        DELAY_BLOCKS = delayBlocks;
    }

    function addCertVerifier(uint64 referenceBlockNumber, address certVerifier) external onlyOwner {
        require(referenceBlockNumber > block.number + DELAY_BLOCKS, "Reference block number must be in the future");
        require(
            referenceBlockNumber > certVerifierRBNs[certVerifierRBNs.length - 1],
            "Reference block number must be greater than the last registered RBN"
        );
        require(
            certVerifiers[referenceBlockNumber] == IEigenDACertVerifierBase(address(0)), "Cert verifier already exists"
        );
        certVerifiers[referenceBlockNumber] = IEigenDACertVerifierBase(certVerifier);
        certVerifierRBNs.push(referenceBlockNumber);
        emit CertVerifierAdded(referenceBlockNumber, certVerifier);
    }

    function verifyDACertV1(BlobHeader calldata blobHeader, BlobVerificationProof calldata blobVerificationProof)
        external
        view
    {
        uint64 referenceBlockNumber = blobVerificationProof.batchMetadata.batchHeader.referenceBlockNumber;
        uint64 closestRBN = _findClosestRegisteredRBN(referenceBlockNumber);
        certVerifiers[closestRBN].verifyDACertV1(blobHeader, blobVerificationProof);
    }

    function verifyDACertsV1(BlobHeader[] calldata blobHeaders, BlobVerificationProof[] calldata blobVerificationProofs)
        external
        view
    {
        require(blobHeaders.length == blobVerificationProofs.length, "Blob headers and proofs length mismatch");
        uint64 referenceBlockNumber = blobVerificationProofs[0].batchMetadata.batchHeader.referenceBlockNumber;
        uint64 closestRBN = _findClosestRegisteredRBN(referenceBlockNumber);
        certVerifiers[closestRBN].verifyDACertsV1(blobHeaders, blobVerificationProofs);
    }

    function verifyDACertV2(
        BatchHeaderV2 calldata batchHeader,
        BlobInclusionInfo calldata blobInclusionInfo,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature,
        bytes memory signedQuorumNumbers
    ) external view {
        uint64 referenceBlockNumber = batchHeader.referenceBlockNumber;
        uint64 closestRBN = _findClosestRegisteredRBN(referenceBlockNumber);
        certVerifiers[closestRBN].verifyDACertV2(
            batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers
        );
    }

    function verifyDACertV2FromSignedBatch(
        SignedBatch calldata signedBatch,
        BlobInclusionInfo calldata blobInclusionInfo
    ) external view {
        uint64 referenceBlockNumber = signedBatch.batchHeader.referenceBlockNumber;
        uint64 closestRBN = _findClosestRegisteredRBN(referenceBlockNumber);
        certVerifiers[closestRBN].verifyDACertV2FromSignedBatch(signedBatch, blobInclusionInfo);
    }

    function verifyDACertV2ForZKProof(
        BatchHeaderV2 calldata batchHeader,
        BlobInclusionInfo calldata blobInclusionInfo,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature,
        bytes memory signedQuorumNumbers
    ) external view returns (bool) {
        uint64 referenceBlockNumber = batchHeader.referenceBlockNumber;
        uint64 closestRBN = _findClosestRegisteredRBN(referenceBlockNumber);
        return certVerifiers[closestRBN].verifyDACertV2ForZKProof(
            batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers
        );
    }

    function getNonSignerStakesAndSignature(SignedBatch calldata signedBatch)
        external
        view
        returns (NonSignerStakesAndSignature memory)
    {
        uint64 referenceBlockNumber = signedBatch.batchHeader.referenceBlockNumber;
        uint64 closestRBN = _findClosestRegisteredRBN(referenceBlockNumber);
        return certVerifiers[closestRBN].getNonSignerStakesAndSignature(signedBatch);
    }

    /// @notice Given an RBN, find the closest RBN registered in this contract that is less than or equal to the given RBN.
    /// @param referenceBlockNumber The reference block number to find the closest RBN for
    /// @return closestRBN The closest RBN registered in this contract that is less than or equal to the given RBN.
    function _findClosestRegisteredRBN(uint64 referenceBlockNumber) internal view returns (uint64) {
        // It is assumed that the latest RBNs are the most likely to be used.
        require(certVerifierRBNs.length > 0, "No cert verifiers available");

        for (uint256 i = certVerifierRBNs.length - 1; i >= 0; i--) {
            uint64 certVerifierRBNMem = certVerifierRBNs[i];
            if (certVerifierRBNMem == referenceBlockNumber) {
                return certVerifierRBNMem;
            }
        }
        revert("No cert verifier found for the given reference block number");
    }
}
