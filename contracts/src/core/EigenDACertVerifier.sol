// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDACertVerifierV1} from "src/periphery/EigenDACertVerifierV1.sol";
import {EigenDACertVerifierV2} from "src/periphery/EigenDACertVerifierV2.sol";
import {IEigenDAThresholdRegistry} from "../interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "../interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "../interfaces/IEigenDASignatureVerifier.sol";
import {EigenDACertVerificationV1Lib as CertV1Lib} from "src/libraries/EigenDACertVerificationV1Lib.sol";
import {EigenDACertVerificationV2Lib as CertV2Lib} from "src/libraries/EigenDACertVerificationV2Lib.sol";
import {OperatorStateRetriever} from "../../lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "../../lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {IEigenDARelayRegistry} from "../interfaces/IEigenDARelayRegistry.sol";
import "../interfaces/IEigenDAStructs.sol";

/**
 * @title A CertVerifier is an immutable contract that is used by a consumer to verify EigenDA blob certificates
 * @notice For V2 verification this contract is deployed with immutable security thresholds and required quorum numbers,
 *         to change these values or verification behavior a new CertVerifier must be deployed
 */
contract EigenDACertVerifier is EigenDACertVerifierV1, EigenDACertVerifierV2 {
    constructor(
        IEigenDAThresholdRegistry _eigenDAThresholdRegistry,
        IEigenDABatchMetadataStorage _eigenDABatchMetadataStorage,
        IEigenDASignatureVerifier _eigenDASignatureVerifier,
        IEigenDARelayRegistry _eigenDARelayRegistry,
        OperatorStateRetriever _operatorStateRetriever,
        IRegistryCoordinator _registryCoordinator,
        SecurityThresholds memory _securityThresholdsV2,
        bytes memory _quorumNumbersRequiredV2
    )
        EigenDACertVerifierV1(_eigenDAThresholdRegistry, _eigenDABatchMetadataStorage)
        EigenDACertVerifierV2(
            _eigenDAThresholdRegistry,
            _eigenDASignatureVerifier,
            _eigenDARelayRegistry,
            _operatorStateRetriever,
            _registryCoordinator,
            _securityThresholdsV2,
            _quorumNumbersRequiredV2
        )
    {}

    function verifyDACertSecurityParams(
        VersionedBlobParams memory blobParams,
        SecurityThresholds memory securityThresholds
    ) external pure {
        (CertV2Lib.StatusCode status, bytes memory statusParams) =
            CertV2Lib.checkSecurityParams(blobParams, securityThresholds);
        CertV2Lib.revertOnError(status, statusParams);
    }

    function verifyDACertSecurityParams(uint16 version, SecurityThresholds memory securityThresholds) external view {
        (CertV2Lib.StatusCode status, bytes memory statusParams) =
            CertV2Lib.checkSecurityParams(getBlobParams(version), securityThresholds);
        CertV2Lib.revertOnError(status, statusParams);
    }

    /**
     * @notice Returns the threshold registry contract
     * @return The IEigenDAThresholdRegistry contract
     * @dev Overrides both V1 and V2 implementations to use a single registry for consistency
     */
    function _thresholdRegistry()
        internal
        view
        override(EigenDACertVerifierV1, EigenDACertVerifierV2)
        returns (IEigenDAThresholdRegistry)
    {
        // This contract enforces that V1 and V2 use the same registry on construction. So we just choose V2.
        return eigenDAThresholdRegistryV2;
    }
}
