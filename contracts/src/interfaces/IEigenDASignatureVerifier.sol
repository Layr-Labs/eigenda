// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "./IEigenDAStructs.sol";

/**
 * @title IEigenDASignatureVerifier
 * @notice This contract is used for verifying the signatures of either V1 batches or V2 certificates
 * @dev This contract is deployed on L1 as the EigenDAServiceManager contract
 */
interface IEigenDASignatureVerifier {

    /**
     * @notice Verifies the BLS signature of a batch certificate
     * @param msgHash The hash of the message to verify
     * @param quorumNumbers The quorum numbers to verify for
     * @param referenceBlockNumber The reference block number of the signature
     * @param params The non-signer stakes and signatures needed to verify the signature
     */
    function checkSignatures(
        bytes32 msgHash,
        bytes calldata quorumNumbers,
        uint32 referenceBlockNumber,
        NonSignerStakesAndSignature memory params
    ) external view returns (QuorumStakeTotals memory, bytes32);
}