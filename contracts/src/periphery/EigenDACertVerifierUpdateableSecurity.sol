// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifier} from "src/interfaces/IEigenDACertVerifier.sol";
import {IEigenDAThresholdRegistry} from "src/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "src/interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "src/interfaces/IEigenDASignatureVerifier.sol";
import {EigenDACertVerificationUtils} from "src/libraries/EigenDACertVerificationUtils.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {IEigenDARelayRegistry} from "src/interfaces/IEigenDARelayRegistry.sol";
import {IEigenDAThresholdRegistry} from "src/interfaces/IEigenDAThresholdRegistry.sol";
import "src/interfaces/IEigenDAStructs.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

/**
 * @title A CertVerifier is an immutable contract that is used by a consumer to verify EigenDA blob certificates
 * @notice For V2 verification this contract is deployed with immutable security thresholds and required quorum numbers,
 *         to change these values or verification behavior a new CertVerifier must be deployed
 */
contract EigenDACertVerifierUpdateableSecurity is IEigenDACertVerifier, OwnableUpgradeable {
    /// @notice The EigenDAThresholdRegistry contract address
    IEigenDAThresholdRegistry public immutable eigenDAThresholdRegistry;

    /// @notice The EigenDABatchMetadataStorage contract address
    /// @dev On L1 this contract is the EigenDA Service Manager contract
    IEigenDABatchMetadataStorage public immutable eigenDABatchMetadataStorage;

    /// @notice The EigenDASignatureVerifier contract address
    /// @dev On L1 this contract is the EigenDA Service Manager contract
    IEigenDASignatureVerifier public immutable eigenDASignatureVerifier;

    /// @notice The EigenDARelayRegistry contract address
    IEigenDARelayRegistry public immutable eigenDARelayRegistry;

    /// @notice The EigenDA middleware OperatorStateRetriever contract address
    OperatorStateRetriever public immutable operatorStateRetriever;

    /// @notice The EigenDA middleware RegistryCoordinator contract address
    IRegistryCoordinator public immutable registryCoordinator;

    mapping(uint32 => SecurityThresholds) public securityThresholdsV2;
    mapping(uint32 => bytes) public quorumNumbersRequiredV2;
    uint32[] public rbns;

    event SecurityThresholdsAndQuorumNumbersAdded(
        uint32 indexed rbn, uint8 indexed confirmationThreshold, uint8 indexed adversaryThreshold, bytes quorumNumbers
    );

    constructor(
        IEigenDAThresholdRegistry _eigenDAThresholdRegistry,
        IEigenDABatchMetadataStorage _eigenDABatchMetadataStorage,
        IEigenDASignatureVerifier _eigenDASignatureVerifier,
        IEigenDARelayRegistry _eigenDARelayRegistry,
        OperatorStateRetriever _operatorStateRetriever,
        IRegistryCoordinator _registryCoordinator
    ) {
        _disableInitializers;
        eigenDAThresholdRegistry = _eigenDAThresholdRegistry;
        eigenDABatchMetadataStorage = _eigenDABatchMetadataStorage;
        eigenDASignatureVerifier = _eigenDASignatureVerifier;
        eigenDARelayRegistry = _eigenDARelayRegistry;
        operatorStateRetriever = _operatorStateRetriever;
        registryCoordinator = _registryCoordinator;
    }

    function initialize(address _initialOwner) external initializer {
        _transferOwnership(_initialOwner);
    }

    function addSecurityThresholdsAndQuorum(
        uint32 referenceBlockNumber,
        SecurityThresholds memory securityThresholds,
        bytes memory quorumNumbers
    ) external onlyOwner {
        require(referenceBlockNumber > block.number, "Reference block number must be in the future");
        require(
            rbns.length == 0 || referenceBlockNumber > rbns[rbns.length - 1],
            "Reference block number must be greater than the last registered RBN"
        );
        securityThresholdsV2[referenceBlockNumber] = securityThresholds;
        quorumNumbersRequiredV2[referenceBlockNumber] = quorumNumbers;
        rbns.push(referenceBlockNumber);
        emit SecurityThresholdsAndQuorumNumbersAdded(
            referenceBlockNumber,
            securityThresholds.confirmationThreshold,
            securityThresholds.adversaryThreshold,
            quorumNumbers
        );
    }

    function getSecurityParamsAt(uint32 rbn) external view returns (SecurityThresholds memory) {
        return securityThresholdsV2[rbn];
    }

    function getQuorumNumbersAt(uint32 rbn) external view returns (bytes memory) {
        return quorumNumbersRequiredV2[rbn];
    }

    ///////////////////////// V1 ///////////////////////////////

    /**
     * @notice Verifies a the blob cert is valid for the required quorums
     * @param blobHeader The blob header to verify
     * @param blobVerificationProof The blob cert verification proof to verify
     */
    function verifyDACertV1(BlobHeader calldata blobHeader, BlobVerificationProof calldata blobVerificationProof)
        public
        view
    {
        uint32 closestRbn =
            _findClosestRegisteredRBN(blobVerificationProof.batchMetadata.batchHeader.referenceBlockNumber);
        EigenDACertVerificationUtils._verifyDACertV1ForQuorums(
            eigenDAThresholdRegistry,
            eigenDABatchMetadataStorage,
            blobHeader,
            blobVerificationProof,
            quorumNumbersRequiredV2[closestRbn]
        );
    }

    /**
     * @notice Verifies a batch of blob certs for the required quorums
     * @param blobHeaders The blob headers to verify
     * @param blobVerificationProofs The blob cert verification proofs to verify against
     */
    function verifyDACertsV1(BlobHeader[] calldata blobHeaders, BlobVerificationProof[] calldata blobVerificationProofs)
        external
        view
    {
        for (uint256 i; i < blobHeaders.length; i++) {
            verifyDACertV1(blobHeaders[i], blobVerificationProofs[i]);
        }
    }

    ///////////////////////// V2 ///////////////////////////////

    /**
     * @notice Verifies a blob cert using the immutable required quorums and security thresholds set in the constructor
     * @param batchHeader The batch header of the blob
     * @param blobInclusionInfo The inclusion proof for the blob cert
     * @param nonSignerStakesAndSignature The nonSignerStakesAndSignature to verify the blob cert against
     * @param signedQuorumNumbers The signed quorum numbers corresponding to the nonSignerStakesAndSignature
     */
    function verifyDACertV2(
        BatchHeaderV2 calldata batchHeader,
        BlobInclusionInfo calldata blobInclusionInfo,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature,
        bytes memory signedQuorumNumbers
    ) external view {
        uint32 closestRbn = _findClosestRegisteredRBN(batchHeader.referenceBlockNumber);
        EigenDACertVerificationUtils._verifyDACertV2ForQuorums(
            eigenDAThresholdRegistry,
            eigenDASignatureVerifier,
            eigenDARelayRegistry,
            batchHeader,
            blobInclusionInfo,
            nonSignerStakesAndSignature,
            securityThresholdsV2[closestRbn],
            quorumNumbersRequiredV2[closestRbn],
            signedQuorumNumbers
        );
    }

    /**
     * @notice Verifies a blob cert using the immutable required quorums and security thresholds set in the constructor
     * @param signedBatch The signed batch to verify the blob cert against
     * @param blobInclusionInfo The inclusion proof for the blob cert
     */
    function verifyDACertV2FromSignedBatch(
        SignedBatch calldata signedBatch,
        BlobInclusionInfo calldata blobInclusionInfo
    ) external view {
        uint32 closestRbn = _findClosestRegisteredRBN(signedBatch.batchHeader.referenceBlockNumber);
        EigenDACertVerificationUtils._verifyDACertV2ForQuorumsFromSignedBatch(
            eigenDAThresholdRegistry,
            eigenDASignatureVerifier,
            eigenDARelayRegistry,
            operatorStateRetriever,
            registryCoordinator,
            signedBatch,
            blobInclusionInfo,
            securityThresholdsV2[closestRbn],
            quorumNumbersRequiredV2[closestRbn]
        );
    }

    /**
     * @notice Thin try/catch wrapper around verifyDACertV2 that returns false instead of panicing
     * @dev The Steel library (https://github.com/risc0/risc0-ethereum/tree/main/crates/steel)
     *      currently has a limitation that it can only create zk proofs for functions that return a value
     * @param batchHeader The batch header of the blob
     * @param blobInclusionInfo The inclusion proof for the blob cert
     * @param nonSignerStakesAndSignature The nonSignerStakesAndSignature to verify the blob cert against
     * @param signedQuorumNumbers The signed quorum numbers corresponding to the nonSignerStakesAndSignature
     */
    function verifyDACertV2ForZKProof(
        BatchHeaderV2 calldata batchHeader,
        BlobInclusionInfo calldata blobInclusionInfo,
        NonSignerStakesAndSignature calldata nonSignerStakesAndSignature,
        bytes memory signedQuorumNumbers
    ) external view returns (bool) {
        uint32 closestRbn = _findClosestRegisteredRBN(batchHeader.referenceBlockNumber);
        try EigenDACertVerificationUtils.verifyDACertV2ForQuorumsExternal(
            eigenDAThresholdRegistry,
            eigenDASignatureVerifier,
            eigenDARelayRegistry,
            batchHeader,
            blobInclusionInfo,
            nonSignerStakesAndSignature,
            securityThresholdsV2[closestRbn],
            quorumNumbersRequiredV2[closestRbn],
            signedQuorumNumbers
        ) {
            return true;
        } catch {
            return false;
        }
    }

    ///////////////////////// HELPER FUNCTIONS ///////////////////////////////

    /**
     * @notice Returns the nonSignerStakesAndSignature for a given blob cert and signed batch
     * @param signedBatch The signed batch to get the nonSignerStakesAndSignature for
     * @return nonSignerStakesAndSignature The nonSignerStakesAndSignature for the given signed batch attestation
     */
    function getNonSignerStakesAndSignature(SignedBatch calldata signedBatch)
        external
        view
        returns (NonSignerStakesAndSignature memory)
    {
        (NonSignerStakesAndSignature memory nonSignerStakesAndSignature,) = EigenDACertVerificationUtils
            ._getNonSignerStakesAndSignature(operatorStateRetriever, registryCoordinator, signedBatch);
        return nonSignerStakesAndSignature;
    }

    /**
     * @notice Verifies the security parameters for a blob cert
     * @param blobParams The blob params to verify
     * @param securityThresholds The security thresholds to verify against
     */
    function verifyDACertSecurityParams(
        VersionedBlobParams memory blobParams,
        SecurityThresholds memory securityThresholds
    ) external pure {
        EigenDACertVerificationUtils._verifyDACertSecurityParams(blobParams, securityThresholds);
    }

    /**
     * @notice Verifies the security parameters for a blob cert
     * @param version The version of the blob to verify
     * @param securityThresholds The security thresholds to verify against
     */
    function verifyDACertSecurityParams(uint16 version, SecurityThresholds memory securityThresholds) external view {
        EigenDACertVerificationUtils._verifyDACertSecurityParams(
            eigenDAThresholdRegistry.getBlobParams(version), securityThresholds
        );
    }

    /// @notice Given an RBN, find the closest RBN registered in this contract that is less than or equal to the given RBN.
    /// @param referenceBlockNumber The reference block number to find the closest RBN for
    /// @return closestRBN The closest RBN registered in this contract that is less than or equal to the given RBN.
    function _findClosestRegisteredRBN(uint32 referenceBlockNumber) internal view returns (uint32) {
        // It is assumed that the latest RBNs are the most likely to be used.
        require(rbns.length > 0, "No rbn available");

        uint256 rbnMaxIndex = rbns.length - 1; // cache to memory
        for (uint256 i; i < rbns.length; i++) {
            uint32 rbnMem = rbns[rbnMaxIndex - i];
            if (rbnMem <= referenceBlockNumber) {
                return rbnMem;
            }
        }
        revert("No rbn found");
    }
}
