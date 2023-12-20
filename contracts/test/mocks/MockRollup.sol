// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.9;

import {EigenDAKZGUtils} from "../../src/libraries/EigenDAKZGUtils.sol";
import {EigenDABlobUtils} from "../../src/libraries/EigenDABlobUtils.sol";
import {EigenDAServiceManager, IEigenDAServiceManager, BN254} from "../../src/core/EigenDAServiceManager.sol";
//import {BN254} from "../../lib/eigenlayer-middleware/src/libraries/BN254.sol";


struct Commitment {
    address validator; // validator who posted the commitment
    uint32 dataLength; // length of the data
    BN254.G1Point polynomialCommitment; // commitment to the polynomial
}

/**
 * @title MockRollup
 * @author Layr Labs, Inc.
 * @notice This contract is used to emulate a rollup contract for the purpose of testing the rollup interface.
 */
contract MockRollup {
    
    IEigenDAServiceManager public eigenDAServiceManager; // EigenDASM contract
    BN254.G1Point public tau; //power of tau
    uint256 public illegalValue; // special "illegal" value that should not be included in blob
    bytes32 public quorumBlobParamsHash; // hash of the security parameters
    uint256 public stakeRequired; // amount of stake required to register as a validator

    ///@notice mapping of validators who have registered
    mapping(address => bool) public validators;
    ///@notice mapping of validators who have been blacklisted
    mapping(address => bool) public blacklist;
    ///@notice mapping of timestamps to commitments
    mapping(uint256 => Commitment) public commitments;

    constructor(IEigenDAServiceManager _eigenDAServiceManager, BN254.G1Point memory _tau, uint256 _illegalValue, bytes32 _quorumBlobParamsHash, uint256 _stakeRequired) {
        eigenDAServiceManager = _eigenDAServiceManager;
        tau = _tau;
        illegalValue = _illegalValue;
        quorumBlobParamsHash = _quorumBlobParamsHash;
        stakeRequired = _stakeRequired;
    }

    ///@notice registers msg.sender as validator by putting up 1 ether of stake
    function registerValidator() external payable {
        require(msg.value == stakeRequired, "MockRollup.registerValidator: Must send stake required to register");
        require(!validators[msg.sender], "MockRollup.registerValidator: Validator already registered");
        require(!blacklist[msg.sender], "MockRollup.registerValidator: Validator blacklisted");
        validators[msg.sender] = true;
    }

    /**
     * @notice a function for validators to post a commitment to a blob on behalf of the rollup
     * @param blobHeader the blob header
     * @param blobVerificationProof the blob verification proof
     */
    function postCommitment(
        IEigenDAServiceManager.BlobHeader memory blobHeader, 
        EigenDABlobUtils.BlobVerificationProof memory blobVerificationProof
    ) external { 
        require(validators[msg.sender], "MockRollup.postCommitment: Validator not registered");
        require(commitments[block.timestamp].validator == address(0), "MockRollup.postCommitment: Commitment already posted");

        // verify that the blob was included in the batch
        EigenDABlobUtils.verifyBlob(blobHeader, eigenDAServiceManager, blobVerificationProof);

        // zero out the chunkLengths (this a temporary hack)
        blobHeader.quorumBlobParams[0].chunkLength = 0;

        // verify that the blob header contains the correct quorumBlobParams
        require(keccak256(abi.encode(blobHeader.quorumBlobParams)) == quorumBlobParamsHash, "MockRollup.postCommitment: QuorumBlobParams do not match quorumBlobParamsHash");


        commitments[block.timestamp] = Commitment(msg.sender, blobHeader.dataLength, blobHeader.commitment);
    }

    /**
     * @notice a function for users to challenge a commitment that contains the illegal value
     * @param timestamp the timestamp of the commitment being challenged
     * @param point the point on the polynomial to evaluate
     * @param proof revelvant KZG proof 
     */
    function challengeCommitment(uint256 timestamp, uint256 point, BN254.G2Point memory proof) external {
        Commitment memory commitment = commitments[timestamp];
        require(commitment.validator != address(0), "MockRollup.challengeCommitment: Commitment not posted");

        // point on the polynomial must be less than the length of the data stored
        require(point < commitment.dataLength, "MockRollup.challengeCommitment: Point must be less than data length");

        // verify that the commitment contains the illegal value
        require(EigenDAKZGUtils.openCommitment(point, illegalValue, tau, commitment.polynomialCommitment, proof), "MockRollup.challengeCommitment: Does not evaluate to illegal value");
            
        // blacklist the validator
        validators[commitment.validator] = false;
        blacklist[commitment.validator] = true;

        // send validators stake to the user who challenged the commitment
        (bool success, ) = msg.sender.call{value: 1 ether}("");
        require(success);
        
    }

}