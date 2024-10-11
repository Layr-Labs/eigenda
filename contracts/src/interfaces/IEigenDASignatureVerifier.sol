// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {BN254} from "eigenlayer-middleware/libraries/BN254.sol";

interface IEigenDASignatureVerifier {
    struct NonSignerStakesAndSignature {
        uint32[] nonSignerQuorumBitmapIndices; 
        BN254.G1Point[] nonSignerPubkeys; 
        BN254.G1Point[] quorumApks; 
        BN254.G2Point apkG2; 
        BN254.G1Point sigma; 
        uint32[] quorumApkIndices; 
        uint32[] totalStakeIndices; 
        uint32[][] nonSignerStakeIndices; 
    }

    struct QuorumStakeTotals {
        uint96[] signedStakeForQuorum;
        uint96[] totalStakeForQuorum;
    }

    function checkSignatures(
        bytes32 msgHash,
        bytes calldata quorumNumbers,
        uint32 referenceBlockNumber,
        NonSignerStakesAndSignature memory params
    ) external view returns (QuorumStakeTotals memory, bytes32);
}