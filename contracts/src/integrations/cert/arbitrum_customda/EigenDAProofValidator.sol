// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "src/integrations/cert/arbitrum_customda/ICustomProofValidator.sol";
import {IEigenDACertVerifierRouter} from "src/integrations/cert/interfaces/IEigenDACertVerifierRouter.sol";

/**
 * @title EigenDAProofValidator
 * @notice Implementation of an EigenDA proof validator contract. This is a shallow implementation and is expected
 *        to change dramatically over the coming months.
 *
 * TODO: Understand what the version byte is that's being passed to validateReadPreimage
 * TODO: Define offchain kzg opening proof serialization standard which is then deserialized and verified against the
 *       DA Cert commitment in the validateReadPreimage function
 *
 * TODO: Add forge tests which assert the logical correctness of validateReadPreimage and validateCertificate
 *       under happy/unhappy cases
 *
 * TODO: Add an E2E test (probably integrated into inabox) which calls the daprovider_generateProof RPC method to serialize a correct ReadPreimage proof
 *       and ensure that it passes when calling the validateReadPreimage opcode
 *
 * TODO: Add an E2E test (probably integrated into inabox) that calls the daprovider_Store method
 *       to generate a CustomDA Commitment which is then passed against validateCertificate function for correctness
 *
 * TODO: Add a deployment foundry deployment script which allows customers to safely deploy this contract
 *
 */
contract EigenDAProofValidator is ICustomDAProofValidator {
    address immutable eigenDACertVeriferRouter;

    constructor(address _eigenDACertVerifierRouter) {
        eigenDACertVeriferRouter = _eigenDACertVerifierRouter;
    }
    /**
     * @notice Validates a EigenDA preimage proof and returns the preimage chunk
     * @param certHash The keccak256 hash of the certificate (from machine's proven state)
     * @param offset The offset into the preimage to read from (from machine's proven state)
     * @param proof The proof data: [certSize(8), certificate, version(1), preimageSize(8), preimageData]
     * @return preimageChunk The 32-byte chunk at the specified offset
     */

    function validateReadPreimage(bytes32 daCommitHash, uint256 offset, bytes calldata proof)
        external
        pure
        override
        returns (bytes memory preimageChunk)
    {
        // Extract DACommit size from proof
        uint256 commitSize;
        assembly {
            commitSize := shr(192, calldataload(add(proof.offset, 0))) // Read 8 bytes
        }

        require(proof.length >= 8 + commitSize, "Proof too short for certificate");
        bytes calldata daCommitment = proof[8:8 + commitSize];

        // Verify daCommitment hash matches what OSP validated
        require(keccak256(daCommitment) == daCommitHash, "DA Commit hash mismatch");

        // First byte must be 0x01 (CustomDA message header flag)
        require(certificate[0] == 0x01, "Invalid certificate header");

        // Second byte must be 0x00 (EigenDA DA Layer byte)
        require(certificate[1] == 0x00, "Invalid DA Layer byte");

        // TODO: Implement kzg proof deserialization and pairing check here.
        //       This will require reading the kzg data commitment from the DA Cert
        //       Blob Header which will require deserializing the cert into a structured
        //       Solidity type for adequate extraction

        // Ensure we aren't proving past the max blob size
        // TODO: will MaxBlobSize ever be defined onchain in the ConfigRegistry?
        require(offset < 16_233_000);

        bytes memory dummyReturnValue = hex"";
        return dummyReturnValue;
    }

    /**
     * @notice Validates whether a daCommit is well-formed and legitimate
     * @dev The proof format is: [daCommitSize(8), daCommit, claimedValid(1)]
     *
     *
     *      Return vs Revert behavior:
     *      - Reverts when:
     *        - Provided daCommit matches proven hash in the instruction (checked in hostio)
     *        - Claimed valid but is invalid and vice versa (checked in hostio)
     *      - Returns false when:
     *        - daCommit is malformed, including wrong length
     *        - checkDACert call against EigenDACertVeriferRouter returns a status code != SUCCESS
     *
     *      - Returns true when:
     *        - checkDACert call against EigenDACertVeriferRouter returns a status code == SUCCESS
     *
     * @param proof The proof data starting with [daCommitSize(8), certificate, claimedValid(1)]
     * @return isValid True if the certificate is valid, false otherwise
     */
    function validateCertificate(bytes calldata proof) external view override returns (bool isValid) {
        // Extract daCommitment size from first 8 bytes of proof
        require(proof.length >= 8, "Proof too short");

        uint256 commitSize;
        assembly {
            commitSize := shr(192, calldataload(add(proof.offset, 0)))
        }

        bytes calldata daCommitment = proof[8:8 + commitSize];

        // DACommit format is: [prefix(1), da_commitment_version(1), eigenda_cert_version(1), eigenda_cert_bytes(N)]
        // First byte must be 0x01 (CustomDA message header flag)
        // Second byte must be 0x00 (EigenDA DA Layer byte flag)
        // Third byte must be the EigenDA Cert version byte (dictated by the EigenDACertVerifier contract)
        // ... Could be beneficial to add an invariant against the cert verifier being used wrt the cert
        //     version being passed here
        //
        // The remaining N bytes are the EigenDA Certificate

        // First byte must be 0x01 (CustomDA message header flag)
        require(certificate[0] == 0x01, "Invalid certificate header");

        // Second byte must be 0x00 (EigenDA DA Layer byte)
        require(certificate[1] == 0x00, "Invalid DA Layer byte");
        //
        // Note: We return false for invalid certificates instead of reverting
        // because the certificate is already onchain. An honest validator must be able
        // to win a challenge to prove that ValidatePreImage should return false
        // so that an invalid cert can be skipped. If this call were to revert then the fraud proof's
        // correctness would be violated.

        uint8 statusCode = IEigenDACertVerifierRouter(eigenDACertVeriferRouter).checkDACert(certificate[3:]);

        // TODO: it'd be nice to compare against actual enum defined in EigenDACertVerifier
        return statusCode == 1;
    }
}
