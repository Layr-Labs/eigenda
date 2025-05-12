// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {EigenDATypesV2 as DATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";
import {EigenDACertVerificationLib as CertLib} from "src/periphery/cert/libraries/EigenDACertVerificationLib.sol";

contract EigenDAStateRetriever is OperatorStateRetriever {
    /**
     * @notice Gets nonSignerStakesAndSignature for a given signed batch
     * @param operatorStateRetriever The operator state retriever contract
     * @param registryCoordinator The registry coordinator contract
     * @param signedBatch The signed batch
     * @return nonSignerStakesAndSignature The non-signer stakes and signature
     * @return signedQuorumNumbers The signed quorum numbers
     */
    function getNonSignerStakesAndSignature(
        OperatorStateRetriever operatorStateRetriever,
        IRegistryCoordinator registryCoordinator,
        DATypesV2.SignedBatch memory signedBatch
    )
        external
        view
        returns (
            DATypesV1.NonSignerStakesAndSignature memory nonSignerStakesAndSignature,
            bytes memory signedQuorumNumbers
        )
    {
        return CertLib.getNonSignerStakesAndSignature(operatorStateRetriever, registryCoordinator, signedBatch);
    }
}
