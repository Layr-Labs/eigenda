// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import "../src/interfaces/IEigenDAServiceManager.sol";

import "forge-std/Script.sol";
import "forge-std/console.sol";
import "../src/interfaces/IEigenDAStructs.sol";

// # To generate the hashes needed for core/serialization_test.go:
// forge script script/GenerateUnitTestHashes.s.sol  -v

contract GenerateHashes is Script {

    string deployConfigPath = "script/input/eigenda_deploy_config.json";

    // deploy all the EigenDA contracts. Relies on many EL contracts having already been deployed.
    function run() external {
        

        QuorumBlobParam[] memory quorumBlobParam = new QuorumBlobParam[](1);
        
        quorumBlobParam[0] = QuorumBlobParam({
            quorumNumber: 0,
            adversaryThresholdPercentage: 80,
            confirmationThresholdPercentage: 100,
            chunkLength: 10
        });


        bytes32 quorumBlobParamsHash = keccak256(abi.encode(quorumBlobParam));
        console.logBytes32(quorumBlobParamsHash);

        BN254.G1Point memory commitment = BN254.G1Point({
            X: 1,
            Y: 2
        });


        quorumBlobParam[0] = QuorumBlobParam({
            quorumNumber: 1,
            adversaryThresholdPercentage: 80,
            confirmationThresholdPercentage: 100,
            chunkLength: 10
        });

        BlobHeader memory header = BlobHeader({
            commitment: commitment,
            dataLength: 10, 
            quorumBlobParams: quorumBlobParam
        });

        
        console.logBytes(abi.encode(header));

        bytes32 blobHeaderHash = keccak256(abi.encode(header));

        console.logBytes32(blobHeaderHash);


    }
}