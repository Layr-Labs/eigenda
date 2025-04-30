// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifier} from "src/interfaces/IEigenDACertVerifier.sol";
import {IEigenDAThresholdRegistry} from "src/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "src/interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "src/interfaces/IEigenDASignatureVerifier.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {EigenDACertVerificationV3Lib as CertV3Lib} from "src/libraries/V3/EigenDACertVerificationV3Lib.sol";
import {EigenDAV3Cert} from "src/libraries/V3/EigenDATypesV3.sol";
import {EigenDATypesV2 as DATypesV2} from "src/libraries/V2/EigenDATypesV2.sol";
import {EigenDATypesV1 as DATypesV1} from "src/libraries/V1/EigenDATypesV1.sol";
import {EigenDACertVerificationV2Lib as CertV2Lib} from "src/libraries/V2/EigenDACertVerificationV2Lib.sol";

contract EigenDACertVerifierV3 is IEigenDACertVerifier {
    error InvalidSecurityThresholds();

    IEigenDAThresholdRegistry public immutable eigenDAThresholdRegistry;

    IEigenDASignatureVerifier public immutable eigenDASignatureVerifier;

    IRegistryCoordinator public immutable registryCoordinator;

    DATypesV1.SecurityThresholds public securityThresholds;

    bytes public quorumNumbersRequired;

    uint8 public constant CERT_VERSION = 3;

    constructor(
        IEigenDAThresholdRegistry _eigenDAThresholdRegistry,
        IEigenDASignatureVerifier _eigenDASignatureVerifier,
        IRegistryCoordinator _registryCoordinator,
        DATypesV1.SecurityThresholds memory _securityThresholds,
        bytes memory _quorumNumbersRequired
    ) {
        if (_securityThresholds.confirmationThreshold <= _securityThresholds.adversaryThreshold) {
            revert InvalidSecurityThresholds();
        }
        eigenDAThresholdRegistry = _eigenDAThresholdRegistry;
        eigenDASignatureVerifier = _eigenDASignatureVerifier;
        registryCoordinator = _registryCoordinator;
        securityThresholds = _securityThresholds;
        quorumNumbersRequired = _quorumNumbersRequired;
    }

    function verifyDACert(bytes calldata certBytes) external view {
        CertV3Lib.verifyDACert(
            eigenDAThresholdRegistry, eigenDASignatureVerifier, certBytes, securityThresholds, quorumNumbersRequired
        );
    }

    function checkDACert(bytes calldata certBytes) external view returns (uint8) {
        (CertV2Lib.StatusCode status,) = CertV3Lib.checkDACert(
            eigenDAThresholdRegistry, eigenDASignatureVerifier, certBytes, securityThresholds, quorumNumbersRequired
        );
        return uint8(status);
    }

}