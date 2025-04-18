// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDACertVerifierV1} from "src/periphery/EigenDACertVerifierV1.sol";
import {EigenDACertVerifierV2} from "src/periphery/EigenDACertVerifierV2.sol";
import {IEigenDAThresholdRegistry} from "src/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDABatchMetadataStorage} from "src/interfaces/IEigenDABatchMetadataStorage.sol";
import {IEigenDASignatureVerifier} from "src/interfaces/IEigenDASignatureVerifier.sol";
import {IEigenDARelayRegistry} from "src/interfaces/IEigenDARelayRegistry.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import "src/interfaces/IEigenDAStructs.sol";

/**
 * @title EigenDACertVerifierV1_V2 - Combined V1 and V2 certificate verification
 * @author Layr Labs, Inc.
 * @notice A unified certificate verifier supporting both V1 and V2 verification methods
 */
contract EigenDACertVerifierV1_V2 is EigenDACertVerifierV1, EigenDACertVerifierV2 {
    /**
     * @notice Constructor for combined V1 and V2 certificate verifier
     * @param _eigenDAThresholdRegistry The EigenDAThresholdRegistry contract
     * @param _eigenDABatchMetadataStorage The EigenDABatchMetadataStorage contract
     * @param _eigenDASignatureVerifier The EigenDASignatureVerifier contract
     * @param _eigenDARelayRegistry The EigenDARelayRegistry contract
     * @param _operatorStateRetriever The OperatorStateRetriever contract
     * @param _registryCoordinator The RegistryCoordinator contract
     * @param _securityThresholdsV2 Security thresholds for V2 cert verification
     * @param _quorumNumbersRequiredV2 Required quorum numbers for V2 cert verification
     */
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
            _securityThresholdsV2
        )
    {}

    /**
     * @notice Gets the quorum numbers required for verification
     * @return bytes The required quorum numbers
     */
    function getQuorumNumbersRequired() external view returns (bytes memory) {
        return _quorumNumbersRequired();
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

    /**
     * @notice Returns the quorum numbers required for verification
     * @return bytes The required quorum numbers
     * @dev Overrides both V1 and V2 implementations to ensure consistency
     */
    function _quorumNumbersRequired()
        internal
        view
        override(EigenDACertVerifierV1, EigenDACertVerifierV2)
        returns (bytes memory)
    {
        return _thresholdRegistry().quorumNumbersRequired();
    }
}
