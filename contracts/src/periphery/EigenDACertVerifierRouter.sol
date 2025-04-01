// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifierBase} from "src/interfaces/IEigenDACertVerifier.sol";
import {Ownable} from "lib/openzeppelin-contracts/contracts/access/Ownable.sol";
import "src/interfaces/IEigenDAStructs.sol";

contract EigenDACertVerifierRouter is IEigenDACertVerifierBase, Ownable {
    /// @notice The mapping of reference block numbers to cert verifiers.
    mapping(uint32 => IEigenDACertVerifierBase) public certVerifiers;

    /// @dev This array should only be appended to in increasing order, and thus should be sorted.
    ///      These values contain all added indexes for the certVerifiers mapping.
    uint32[] public certVerifierRBNs;

    event CertVerifierAdded(uint32 indexed referenceBlockNumber, address indexed certVerifier);

    function addCertVerifier(uint32 referenceBlockNumber, address certVerifier) external onlyOwner {
        require(referenceBlockNumber > block.number, "Reference block number must be in the future");
        require(
            certVerifierRBNs.length == 0 || referenceBlockNumber > certVerifierRBNs[certVerifierRBNs.length - 1],
            "Reference block number must be greater than the last registered RBN"
        );
        certVerifiers[referenceBlockNumber] = IEigenDACertVerifierBase(certVerifier);
        certVerifierRBNs.push(referenceBlockNumber);
        emit CertVerifierAdded(referenceBlockNumber, certVerifier);
    }

    function verifyDACertV1(BlobHeader calldata blobHeader, BlobVerificationProof calldata blobVerificationProof)
        public
        view
    {
        uint32 referenceBlockNumber = blobVerificationProof.batchMetadata.batchHeader.referenceBlockNumber;
        uint32 closestRBN = _findClosestRegisteredRBN(referenceBlockNumber);
        certVerifiers[closestRBN].verifyDACertV1(blobHeader, blobVerificationProof);
    }

    function verifyDACertsV1(BlobHeader[] calldata blobHeaders, BlobVerificationProof[] calldata blobVerificationProofs)
        external
        view
    {
        require(blobHeaders.length == blobVerificationProofs.length, "Blob headers and proofs length mismatch");
        for (uint256 i; i < blobHeaders.length; i++) {
            verifyDACertV1(blobHeaders[i], blobVerificationProofs[i]);
        }
    }

    function verifyDACertV2(
        BatchHeaderV2 calldata batchHeader,
        BlobInclusionInfo calldata blobInclusionInfo,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature,
        bytes memory signedQuorumNumbers
    ) external view {
        uint32 referenceBlockNumber = batchHeader.referenceBlockNumber;
        uint32 closestRBN = _findClosestRegisteredRBN(referenceBlockNumber);
        certVerifiers[closestRBN].verifyDACertV2(
            batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers
        );
    }

    function verifyDACertV2FromSignedBatch(
        SignedBatch calldata signedBatch,
        BlobInclusionInfo calldata blobInclusionInfo
    ) external view {
        uint32 referenceBlockNumber = signedBatch.batchHeader.referenceBlockNumber;
        uint32 closestRBN = _findClosestRegisteredRBN(referenceBlockNumber);
        certVerifiers[closestRBN].verifyDACertV2FromSignedBatch(signedBatch, blobInclusionInfo);
    }

    function verifyDACertV2ForZKProof(
        BatchHeaderV2 calldata batchHeader,
        BlobInclusionInfo calldata blobInclusionInfo,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature,
        bytes memory signedQuorumNumbers
    ) external view returns (bool) {
        uint32 referenceBlockNumber = batchHeader.referenceBlockNumber;
        uint32 closestRBN = _findClosestRegisteredRBN(referenceBlockNumber);
        return certVerifiers[closestRBN].verifyDACertV2ForZKProof(
            batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers
        );
    }

    function getNonSignerStakesAndSignature(SignedBatch calldata signedBatch)
        external
        view
        returns (NonSignerStakesAndSignature memory)
    {
        uint32 referenceBlockNumber = signedBatch.batchHeader.referenceBlockNumber;
        uint32 closestRBN = _findClosestRegisteredRBN(referenceBlockNumber);
        return certVerifiers[closestRBN].getNonSignerStakesAndSignature(signedBatch);
    }

    /// @notice Given an RBN, find the closest RBN registered in this contract that is less than or equal to the given RBN.
    /// @param referenceBlockNumber The reference block number to find the closest RBN for
    /// @return closestRBN The closest RBN registered in this contract that is less than or equal to the given RBN.
    function _findClosestRegisteredRBN(uint32 referenceBlockNumber) internal view returns (uint32) {
        // It is assumed that the latest RBNs are the most likely to be used.
        require(certVerifierRBNs.length > 0, "No cert verifiers available");

        uint256 rbnMaxIndex = certVerifierRBNs.length - 1; // cache to memory
        for (uint256 i; i < certVerifierRBNs.length; i++) {
            uint32 certVerifierRBNMem = certVerifierRBNs[rbnMaxIndex - i];
            if (certVerifierRBNMem <= referenceBlockNumber) {
                return certVerifierRBNMem;
            }
        }
        revert("No cert verifier found for the given reference block number");
    }
}
