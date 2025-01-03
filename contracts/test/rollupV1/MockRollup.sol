// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDARollupUtils} from "./EigenDARollupUtils.sol";
import {EigenDAServiceManager} from "../../src/core/EigenDAServiceManager.sol";
import {IEigenDAServiceManager} from "../../src/interfaces/IEigenDAServiceManager.sol";
import {BN254} from "eigenlayer-middleware/libraries/BN254.sol";
import "../../src/interfaces/IEigenDAStructs.sol";

struct Commitment {
    address confirmer; // confirmer who posted the commitment
    uint32 dataLength; // length of the data
    BN254.G1Point polynomialCommitment; // commitment to the polynomial
}

contract MockRollup {

    IEigenDAServiceManager public eigenDAServiceManager; // EigenDASM contract
    BN254.G1Point public tau; //power of tau

    ///@notice mapping of timestamps to commitments
    mapping(uint256 => Commitment) public commitments;

    constructor(IEigenDAServiceManager _eigenDAServiceManager, BN254.G1Point memory _tau) {
        eigenDAServiceManager = _eigenDAServiceManager;
        tau = _tau;
    }

    /**
     * @notice a function for a confirmer to post a commitment to a blob and verfiy it on EigenDA
     * @param blobHeader the blob header
     * @param blobVerificationProof the blob verification proof
     */
    function postCommitment(
        BlobHeader memory blobHeader, 
        BlobVerificationProof memory blobVerificationProof
    ) external { 
        // require commitment has not already been posted
        // require(commitments[block.timestamp].confirmer == address(0), "MockRollup.postCommitment: Commitment already posted");

        // verify that the blob was included in the batch
        EigenDARollupUtils.verifyBlob(blobHeader, eigenDAServiceManager, blobVerificationProof);

        // store the commitment
        commitments[block.timestamp] = Commitment(msg.sender, blobHeader.dataLength, blobHeader.commitment);
    }

    /**
     * @notice a function for users to challenge a commitment against a provided value
     * @param timestamp the timestamp of the commitment being challenged
     * @param point the point on the polynomial to evaluate
     * @param proof revelvant KZG proof 
     * @param challengeValue The value expected upon opening the commitment
     */
    function challengeCommitment(uint256 timestamp, uint256 point, BN254.G2Point memory proof, uint256 challengeValue) external returns (bool) {
        Commitment memory commitment = commitments[timestamp];
        // require the commitment exists
        require(commitment.confirmer != address(0), "MockRollup.challengeCommitment: Commitment not posted");

        // point on the polynomial must be less than the length of the data stored
        require(point < commitment.dataLength, "MockRollup.challengeCommitment: Point must be less than data length");

        // verify that the commitment contains the challenge value
        return EigenDARollupUtils.openCommitment(point, challengeValue, tau, commitment.polynomialCommitment, proof);
    }

}