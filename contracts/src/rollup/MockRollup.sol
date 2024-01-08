// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.9;

import {EigenDAKZGUtils} from "../libraries/EigenDAKZGUtils.sol";
import {EigenDABlobUtils} from "../libraries/EigenDABlobUtils.sol";
import {EigenDAServiceManager, IEigenDAServiceManager, BN254} from "../core/EigenDAServiceManager.sol";


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
    uint256 public stakeRequired; // amount of stake required to register as a validator

    ///@notice mapping of validators who have registered
    mapping(address => bool) public validators;
    ///@notice mapping of validators who have been blacklisted
    mapping(address => bool) public blacklist;
    ///@notice mapping of timestamps to commitments
    mapping(uint256 => Commitment) public commitments;

    constructor(IEigenDAServiceManager _eigenDAServiceManager, BN254.G1Point memory _tau, uint256 _illegalValue, uint256 _stakeRequired) {
        eigenDAServiceManager = _eigenDAServiceManager;
        tau = _tau;
        illegalValue = _illegalValue;
        stakeRequired = _stakeRequired;
    }

    ///@notice registers msg.sender as validator by putting up 1 ether of stake
    function registerValidator() external payable {
        require(msg.value == stakeRequired, "MockRollup.registerValidator: Must send stake required to register");
        require(!validators[msg.sender], "MockRollup.registerValidator: Validator already registered");
        require(!blacklist[msg.sender], "MockRollup.registerValidator: Validator blacklisted");
        validators[msg.sender] = true;
    }

    ///@notice deregisters msg.sender as validator
    function deRegisterValidator() external {
        require(validators[msg.sender], "MockRollup.registerValidator: Validator already registered");
        require(!blacklist[msg.sender], "MockRollup.registerValidator: Validator blacklisted");
        validators[msg.sender] = false;
    }

    /**
     * @notice a function for validators to post a commitment to a blob on behalf of the rollup
     * @param blobHeader the blob header
     * @param blobVerificationProof the blob verification proof
     */
    function postCommitment(
        IEigenDAServiceManager.BlobHeader memory blobHeader, 
        EigenDABlobUtils.BlobVerificationProof memory blobVerificationProof,
        bytes32 blobParamsHashInput
    ) external { 
        require(validators[msg.sender], "MockRollup.postCommitment: Validator not registered");
        require(commitments[block.timestamp].validator == address(0), "MockRollup.postCommitment: Commitment already posted");

        // Calculate Hash of QuorumBlobParams
        bytes32 quorumBlobParamsHash = computeQuorumBlobParamsHash(blobHeader.quorumBlobParams);

        // verify that the blob header matches computed blob params hash
        require(quorumBlobParamsHash == blobParamsHashInput, "MockRollup.postCommitment: QuorumBlobParams do not match quorumBlobParamsHash");

        
        // verify that the blob was included in the batch
        EigenDABlobUtils.verifyBlob(blobHeader, eigenDAServiceManager, blobVerificationProof);

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

    // Helper function to compute the keccak256 hash of QuorumBlobParam
    function computeQuorumBlobParamsHash(IEigenDAServiceManager.QuorumBlobParam[] memory quorumBlobParams) public pure returns (bytes32) {
        bytes memory serializedData = abi.encode(quorumBlobParams.length);

        for (uint i = 0; i < quorumBlobParams.length; i++) {
            // Serializing each QuorumBlobParam
            serializedData = abi.encodePacked(
                serializedData,
                quorumBlobParams[i].quorumNumber,
                quorumBlobParams[i].adversaryThresholdPercentage,
                quorumBlobParams[i].quorumThresholdPercentage,
                quorumBlobParams[i].quantizationParameter
            );
        }

        // Compute and return the keccak256 hash of the serialized data
        return keccak256(serializedData);
    }

}