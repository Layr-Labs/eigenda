// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "./IEigenDAStructs.sol";

interface IEigenDASignatureVerifier {
    function checkSignatures(
        bytes32 msgHash,
        bytes calldata quorumNumbers,
        uint32 referenceBlockNumber,
        NonSignerStakesAndSignature memory params
    ) external view returns (QuorumStakeTotals memory, bytes32);
}