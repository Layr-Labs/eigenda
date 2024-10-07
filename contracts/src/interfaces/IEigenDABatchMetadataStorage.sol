// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

interface IEigenDABatchMetadataStorage {
    function batchIdToBatchMetadataHash(uint32 batchId) external view returns (bytes32);
}