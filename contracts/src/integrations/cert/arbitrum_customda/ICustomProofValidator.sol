// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

/*
    Interface vendored from offchainlabs/nitro-contracts
    https://github.com/OffchainLabs/nitro-contracts/blob/b85e20c22ce9d140c155a9ad51051e08d1031899/src/osp/ICustomDAProofValidator.sol
*/

/**
 * @title ICustomDAProofValidator
 * @notice Interface for custom data availability proof validators
 */
interface ICustomDAProofValidator {
    /**
     * @notice Validates a custom DA proof and returns the preimage chunk
     * @param certHash The keccak256 hash of the certificate (from machine's proven state)
     * @param offset The offset into the preimage to read from (from machine's proven state)
     * @param proof The proof data starting with [certSize(8), certificate, customData...]
     * @return preimageChunk The 32-byte chunk of preimage data at the specified offset
     */
    function validateReadPreimage(bytes32 certHash, uint256 offset, bytes calldata proof)
        external
        view
        returns (bytes memory preimageChunk);

    /**
     * @notice Validates whether a certificate is well-formed and legitimate
     * @dev This function MUST NOT revert. It should return false for malformed or invalid certificates.
     *      The security model requires that the prover's validity claim matches what this function returns.
     *      If they disagree (e.g., prover claims valid but this returns false), the OSP will revert.
     *
     *      The proof format is: [certSize(8), certificate, claimedValid(1), validityProof...]
     *      The validityProof section can contain additional verification data such as:
     *      - Cryptographic signatures
     *      - Merkle proofs
     *      - Timestamped attestations
     *      - Or other authentication mechanisms
     * @param proof The proof data starting with [certSize(8), certificate, validityProof...]
     * @return isValid True if the certificate is valid, false otherwise
     */
    function validateCertificate(bytes calldata proof) external view returns (bool isValid);
}
