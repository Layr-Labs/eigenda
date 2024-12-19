// SPDX-License-Identifier: MIT

pragma solidity ^0.8.9;

import {IEigenDAServiceManager} from "../interfaces/IEigenDAServiceManager.sol";
import "../interfaces/IEigenDAStructs.sol";

/**
 * @title Library of functions for hashing various EigenDA structs.
 * @author Layr Labs, Inc.
 */
library EigenDAHasher {

    /**
     * @notice hashes the given metdata into the commitment that will be stored in the contract
     * @param batchHeaderHash the hash of the batchHeader
     * @param signatoryRecordHash the hash of the signatory record
     * @param blockNumber the block number at which the batch was confirmed
     */
    function hashBatchHashedMetadata(
        bytes32 batchHeaderHash,
        bytes32 signatoryRecordHash,
        uint32 blockNumber
    ) internal pure returns(bytes32) {
        return keccak256(abi.encodePacked(batchHeaderHash, signatoryRecordHash, blockNumber));
    }

    /**
     * @notice hashes the given metdata into the commitment that will be stored in the contract
     * @param batchHeaderHash the hash of the batchHeader
     * @param confirmationData the confirmation data of the batch
     * @param blockNumber the block number at which the batch was confirmed
     */
    function hashBatchHashedMetadata(
        bytes32 batchHeaderHash,
        bytes memory confirmationData,
        uint32 blockNumber
    ) internal pure returns(bytes32) {
        return keccak256(abi.encodePacked(batchHeaderHash, confirmationData, blockNumber));
    }

    /**
     * @notice given the batchHeader in the provided metdata, calculates the hash of the batchMetadata
     * @param batchMetadata the metadata of the batch
     */
    function hashBatchMetadata(
        BatchMetadata memory batchMetadata
    ) internal pure returns(bytes32) {
        return hashBatchHashedMetadata(
            keccak256(abi.encode(batchMetadata.batchHeader)),
            batchMetadata.signatoryRecordHash,
            batchMetadata.confirmationBlockNumber
        );
    }

    /**
     * @notice hashes the given batch header
     * @param batchHeader the batch header to hash
     */
    function hashBatchHeaderMemory(BatchHeader memory batchHeader) internal pure returns(bytes32) {
        return keccak256(abi.encode(batchHeader));
    }

    /**
     * @notice hashes the given batch header
     * @param batchHeader the batch header to hash
     */
    function hashBatchHeader(BatchHeader calldata batchHeader) internal pure returns(bytes32) {
        return keccak256(abi.encode(batchHeader));
    }

    /**
     * @notice hashes the given reduced batch header
     * @param reducedBatchHeader the reduced batch header to hash
     */
    function hashReducedBatchHeader(ReducedBatchHeader memory reducedBatchHeader) internal pure returns(bytes32) {
        return keccak256(abi.encode(reducedBatchHeader));
    }

    /**
     * @notice hashes the given blob header
     * @param blobHeader the blob header to hash
     */
    function hashBlobHeader(BlobHeader memory blobHeader) internal pure returns(bytes32) {
        return keccak256(abi.encode(blobHeader));
    }

    /**
     * @notice converts a batch header to a reduced batch header
     * @param batchHeader the batch header to convert
     */
    function convertBatchHeaderToReducedBatchHeader(BatchHeader memory batchHeader) internal pure returns(ReducedBatchHeader memory) {
        return ReducedBatchHeader({
            blobHeadersRoot: batchHeader.blobHeadersRoot,
            referenceBlockNumber: batchHeader.referenceBlockNumber
        });
    }

    /**
     * @notice converts the given batch header to a reduced batch header and then hashes it
     * @param batchHeader the batch header to hash
     */
    function hashBatchHeaderToReducedBatchHeader(BatchHeader memory batchHeader) internal pure returns(bytes32) {
        return keccak256(abi.encode(convertBatchHeaderToReducedBatchHeader(batchHeader)));
    }

    /**
     * @notice hashes the given V2 batch header 
     * @param batchHeader the V2 batch header to hash
     */
    function hashBatchHeaderV2(BatchHeaderV2 memory batchHeader) internal pure returns(bytes32) {
        return keccak256(abi.encode(batchHeader));
    }

    /**
     * @notice hashes the given V2 blob header
     * @param blobHeader the V2 blob header to hash
     */
    function hashBlobHeaderV2(BlobHeaderV2 memory blobHeader) internal pure returns(bytes32) {
        return keccak256(
            abi.encode(
                keccak256(
                    abi.encode(
                        blobHeader.version,
                        blobHeader.quorumNumbers,
                        blobHeader.commitment
                    )
                ),
                blobHeader.paymentHeaderHash
            )
        );
    }

    /**
     * @notice hashes the given V2 blob certificate
     * @param blobCertificate the V2 blob certificate to hash
     */
    function hashBlobCertificate(BlobCertificate memory blobCertificate) internal pure returns(bytes32) {
        return keccak256(
            abi.encode(
                hashBlobHeaderV2(blobCertificate.blobHeader),
                blobCertificate.relayKeys
            )
        );
    }
}
