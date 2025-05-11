// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";

interface IEigenDASignatureVerifier {
    function checkSignatures(
        bytes32 msgHash,
        bytes calldata quorumNumbers,
        uint32 referenceBlockNumber,
        EigenDATypesV1.NonSignerStakesAndSignature memory params
    ) external view returns (EigenDATypesV1.QuorumStakeTotals memory, bytes32);
}
