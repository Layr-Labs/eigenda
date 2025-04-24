// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifierV2} from "src/interfaces/IEigenDACertVerifierV2.sol";
import {IEigenDAThresholdRegistry} from "src/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "src/interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "src/interfaces/IEigenDASignatureVerifier.sol";
import {EigenDACertVerificationV1Lib as CertV1Lib} from "src/libraries/EigenDACertVerificationV1Lib.sol";
import {EigenDACertVerificationV2Lib as CertV2Lib} from "src/libraries/EigenDACertVerificationV2Lib.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {IEigenDARelayRegistry} from "src/interfaces/IEigenDARelayRegistry.sol";
import {
    BatchHeaderV2,
    BlobInclusionInfo,
    NonSignerStakesAndSignature,
    SignedBatch,
    SecurityThresholds,
    EigenDACertV2
} from "src/interfaces/IEigenDAStructs.sol";

/**
 * @title A CertVerifier is an immutable contract that is used by a consumer to verify EigenDA blob certificates
 * @notice For V2 verification this contract is deployed with immutable security thresholds and required quorum numbers,
 *         to change these values or verification behavior a new CertVerifier must be deployed
 */
contract EigenDACertVerifierV2 is IEigenDACertVerifierV2 {
    error InvalidSecurityThresholds();

    /// @notice The EigenDAThresholdRegistry contract address
    IEigenDAThresholdRegistry public immutable eigenDAThresholdRegistryV2;

    /// @notice The EigenDASignatureVerifier contract address
    IEigenDASignatureVerifier public immutable eigenDASignatureVerifierV2;

    /// @notice The EigenDA middleware RegistryCoordinator contract address
    IRegistryCoordinator public immutable registryCoordinatorV2;

    SecurityThresholds public securityThresholdsV2;

    bytes public quorumNumbersRequiredV2;

    /**
     * @notice Constructor for the EigenDA V2 certificate verifier
     * @param _eigenDAThresholdRegistryV2 The address of the EigenDAThresholdRegistry contract
     * @param _eigenDASignatureVerifierV2 The address of the EigenDASignatureVerifier contract
     * @param _registryCoordinatorV2 The address of the RegistryCoordinator contract
     * @param _securityThresholdsV2 The security thresholds for verification
     */
    constructor(
        IEigenDAThresholdRegistry _eigenDAThresholdRegistryV2,
        IEigenDASignatureVerifier _eigenDASignatureVerifierV2,
        IRegistryCoordinator _registryCoordinatorV2,
        SecurityThresholds memory _securityThresholdsV2,
        bytes memory _quorumNumbersRequiredV2
    ) {
        if (_securityThresholdsV2.confirmationThreshold <= _securityThresholdsV2.adversaryThreshold) {
            revert InvalidSecurityThresholds();
        }
        eigenDAThresholdRegistryV2 = _eigenDAThresholdRegistryV2;
        eigenDASignatureVerifierV2 = _eigenDASignatureVerifierV2;
        registryCoordinatorV2 = _registryCoordinatorV2;
        securityThresholdsV2 = _securityThresholdsV2;
        quorumNumbersRequiredV2 = _quorumNumbersRequiredV2;
    }

    /**
     * @notice Verifies a blob cert using the immutable required quorums and security thresholds set in the constructor
     * @param cert The EigenDACertV2 to verify
     */
    function verifyDACertV2(EigenDACertV2 calldata cert) external view {
        CertV2Lib.verifyDACertV2(
            _thresholdRegistry(), _signatureVerifier(), cert, _securityThresholds(), _quorumNumbersRequired()
        );
    }

    /**
     * @notice Thin try/catch wrapper around verifyDACertV2 that returns false instead of panicing
     * @dev The Steel library (https://github.com/risc0/risc0-ethereum/tree/main/crates/steel)
     *      currently has a limitation that it can only create zk proofs for functions that return a value
     * @param cert The EigenDACertV2 to verify
     */
    function checkDACert(EigenDACertV2 calldata cert) external view returns (bool) {
        (CertV2Lib.StatusCode err,) = CertV2Lib.checkDACertV2(
            _thresholdRegistry(), _signatureVerifier(), cert, _securityThresholds(), _quorumNumbersRequired()
        );
        if (err == CertV2Lib.StatusCode.SUCCESS) {
            return true;
        } else {
            return false;
        }
    }

    /**
     * @notice Returns the threshold registry contract
     * @return The IEigenDAThresholdRegistry contract
     * @dev Can be overridden by derived contracts
     */
    function _thresholdRegistry() internal view virtual returns (IEigenDAThresholdRegistry) {
        return eigenDAThresholdRegistryV2;
    }

    /**
     * @notice Returns the signature verifier contract
     * @return The IEigenDASignatureVerifier contract
     * @dev Can be overridden by derived contracts
     */
    function _signatureVerifier() internal view virtual returns (IEigenDASignatureVerifier) {
        return eigenDASignatureVerifierV2;
    }

    /**
     * @notice Returns the registry coordinator contract
     * @return The IRegistryCoordinator contract
     * @dev Can be overridden by derived contracts
     */
    function _registryCoordinator() internal view virtual returns (IRegistryCoordinator) {
        return registryCoordinatorV2;
    }

    /**
     * @notice Returns the security thresholds used for verification
     * @return The SecurityThresholds struct with confirmation and adversary thresholds
     * @dev Can be overridden by derived contracts
     */
    function _securityThresholds() internal view virtual returns (SecurityThresholds memory) {
        return securityThresholdsV2;
    }

    /**
     * @notice Returns the quorum numbers required for verification
     * @return bytes The required quorum numbers
     * @dev Can be overridden by derived contracts
     */
    function _quorumNumbersRequired() internal view virtual returns (bytes memory) {
        return quorumNumbersRequiredV2;
    }
}
