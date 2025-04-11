// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import {IRegistry} from "./IRegistry.sol";

import {BN254} from "../libraries/BN254.sol";

/**
 * @title Minimal interface for a registry that keeps track of aggregate operator public keys across many quorums.
 * @author Layr Labs, Inc.
 */
interface IBLSApkRegistry is IRegistry {
    // STRUCTS
    /// @notice Data structure used to track the history of the Aggregate Public Key of all operators
    struct ApkUpdate {
        // first 24 bytes of keccak256(apk_x0, apk_x1, apk_y0, apk_y1)
        bytes24 apkHash;
        // block number at which the update occurred
        uint32 updateBlockNumber;
        // block number at which the next update occurred
        uint32 nextUpdateBlockNumber;
    }

    /**
     * @notice Struct used when registering a new public key
     * @param pubkeyRegistrationSignature is the registration message signed by the private key of the operator
     * @param pubkeyG1 is the corresponding G1 public key of the operator 
     * @param pubkeyG2 is the corresponding G2 public key of the operator
     */     
    struct PubkeyRegistrationParams {
        BN254.G1Point pubkeyRegistrationSignature;
        BN254.G1Point pubkeyG1;
        BN254.G2Point pubkeyG2;
    }

    // EVENTS
    /// @notice Emitted when `operator` registers with the public keys `pubkeyG1` and `pubkeyG2`.
    event NewPubkeyRegistration(address indexed operator, BN254.G1Point pubkeyG1, BN254.G2Point pubkeyG2);

    // @notice Emitted when a new operator pubkey is registered for a set of quorums
    event OperatorAddedToQuorums(
        address operator,
        bytes32 operatorId,
        bytes quorumNumbers
    );

    // @notice Emitted when an operator pubkey is removed from a set of quorums
    event OperatorRemovedFromQuorums(
        address operator, 
        bytes32 operatorId,
        bytes quorumNumbers
    );

    /**
     * @notice Registers the `operator`'s pubkey for the specified `quorumNumbers`.
     * @param operator The address of the operator to register.
     * @param quorumNumbers The quorum numbers the operator is registering for, where each byte is an 8 bit integer quorumNumber.
     * @dev access restricted to the RegistryCoordinator
     * @dev Preconditions (these are assumed, not validated in this contract):
     *         1) `quorumNumbers` has no duplicates
     *         2) `quorumNumbers.length` != 0
     *         3) `quorumNumbers` is ordered in ascending order
     *         4) the operator is not already registered
     */
    function registerOperator(address operator, bytes calldata quorumNumbers) external;

    /**
     * @notice Deregisters the `operator`'s pubkey for the specified `quorumNumbers`.
     * @param operator The address of the operator to deregister.
     * @param quorumNumbers The quorum numbers the operator is deregistering from, where each byte is an 8 bit integer quorumNumber.
     * @dev access restricted to the RegistryCoordinator
     * @dev Preconditions (these are assumed, not validated in this contract):
     *         1) `quorumNumbers` has no duplicates
     *         2) `quorumNumbers.length` != 0
     *         3) `quorumNumbers` is ordered in ascending order
     *         4) the operator is not already deregistered
     *         5) `quorumNumbers` is a subset of the quorumNumbers that the operator is registered for
     */ 
    function deregisterOperator(address operator, bytes calldata quorumNumbers) external;
    
    /**
     * @notice Initializes a new quorum by pushing its first apk update
     * @param quorumNumber The number of the new quorum
     */
    function initializeQuorum(uint8 quorumNumber) external;

    /**
     * @notice mapping from operator address to pubkey hash.
     * Returns *zero* if the `operator` has never registered, and otherwise returns the hash of the public key of the operator.
     */
    function operatorToPubkeyHash(address operator) external view returns (bytes32);

    /**
     * @notice mapping from pubkey hash to operator address.
     * Returns *zero* if no operator has ever registered the public key corresponding to `pubkeyHash`,
     * and otherwise returns the (unique) registered operator who owns the BLS public key that is the preimage of `pubkeyHash`.
     */
    function pubkeyHashToOperator(bytes32 pubkeyHash) external view returns (address);

    /**
     * @notice Called by the RegistryCoordinator register an operator as the owner of a BLS public key.
     * @param operator is the operator for whom the key is being registered
     * @param params contains the G1 & G2 public keys of the operator, and a signature proving their ownership
     * @param pubkeyRegistrationMessageHash is a hash that the operator must sign to prove key ownership
     */
    function registerBLSPublicKey(
        address operator,
        PubkeyRegistrationParams calldata params,
        BN254.G1Point calldata pubkeyRegistrationMessageHash
    ) external returns (bytes32 operatorId);

    /**
     * @notice Returns the pubkey and pubkey hash of an operator
     * @dev Reverts if the operator has not registered a valid pubkey
     */
    function getRegisteredPubkey(address operator) external view returns (BN254.G1Point memory, bytes32);

    /// @notice Returns the current APK for the provided `quorumNumber `
    function getApk(uint8 quorumNumber) external view returns (BN254.G1Point memory);

    /// @notice Returns the index of the quorumApk index at `blockNumber` for the provided `quorumNumber`
    function getApkIndicesAtBlockNumber(bytes calldata quorumNumbers, uint256 blockNumber) external view returns(uint32[] memory);

    /// @notice Returns the `ApkUpdate` struct at `index` in the list of APK updates for the `quorumNumber`
    function getApkUpdateAtIndex(uint8 quorumNumber, uint256 index) external view returns (ApkUpdate memory);

    /// @notice Returns the operator address for the given `pubkeyHash`
    function getOperatorFromPubkeyHash(bytes32 pubkeyHash) external view returns (address);

    /**
     * @notice get 24 byte hash of the apk of `quorumNumber` at `blockNumber` using the provided `index`;
     * called by checkSignatures in BLSSignatureChecker.sol.
     * @param quorumNumber is the quorum whose ApkHash is being retrieved
     * @param blockNumber is the number of the block for which the latest ApkHash will be retrieved
     * @param index is the index of the apkUpdate being retrieved from the list of quorum apkUpdates in storage
     */
    function getApkHashAtBlockNumberAndIndex(uint8 quorumNumber, uint32 blockNumber, uint256 index) external view returns (bytes24);

    /// @notice returns the ID used to identify the `operator` within this AVS.
    /// @dev Returns zero in the event that the `operator` has never registered for the AVS
    function getOperatorId(address operator) external view returns (bytes32);
}
