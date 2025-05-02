// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import "src/core/interfaces/IEigenDAServiceManager.sol";

import "forge-std/Script.sol";
import "forge-std/console.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {BN254} from "lib/eigenlayer-middleware/src/libraries/BN254.sol";

// # To generate the hashes needed for core/serialization_test.go:
// forge script script/GenerateUnitTestHashes.s.sol  -v

contract GenerateHashes is Script {
    string deployConfigPath = "script/input/eigenda_deploy_config.json";

    function run() external pure {
        DATypesV1.QuorumBlobParam[] memory quorumBlobParam = new DATypesV1.QuorumBlobParam[](1);

        quorumBlobParam[0] = DATypesV1.QuorumBlobParam({
            quorumNumber: 0,
            adversaryThresholdPercentage: 80,
            confirmationThresholdPercentage: 100,
            chunkLength: 10
        });

        bytes32 quorumBlobParamsHash = keccak256(abi.encode(quorumBlobParam));
        console.logBytes32(quorumBlobParamsHash);

        BN254.G1Point memory commitment = BN254.G1Point({X: 1, Y: 2});

        quorumBlobParam[0] = DATypesV1.QuorumBlobParam({
            quorumNumber: 1,
            adversaryThresholdPercentage: 80,
            confirmationThresholdPercentage: 100,
            chunkLength: 10
        });

        DATypesV1.BlobHeader memory header =
            DATypesV1.BlobHeader({commitment: commitment, dataLength: 10, quorumBlobParams: quorumBlobParam});

        console.logBytes(abi.encode(header));

        bytes32 blobHeaderHash = keccak256(abi.encode(header));

        console.logBytes32(blobHeaderHash);
    }
}
