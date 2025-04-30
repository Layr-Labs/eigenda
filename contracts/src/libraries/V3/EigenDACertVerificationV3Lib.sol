// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDAV3Cert} from "src/libraries/V3/EigenDATypesV3.sol";
import {EigenDATypesV2 as DATypesV2} from "src/libraries/V2/EigenDATypesV2.sol";
import {EigenDATypesV1 as DATypesV1} from "src/libraries/V1/EigenDATypesV1.sol";
import {EigenDACertVerificationV2Lib as V2Lib} from "src/libraries/V2/EigenDACertVerificationV2Lib.sol";
import {IEigenDAThresholdRegistry} from "src/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/interfaces/IEigenDASignatureVerifier.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";

library EigenDACertVerificationV3Lib {
    function decodeCert(bytes calldata data) internal pure returns (EigenDAV3Cert calldata cert) {
        assembly {
            cert := data.offset
        }
    }

    function verifyDACert(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDASignatureVerifier eigenDASignatureVerifier,
        bytes calldata certBytes,
        DATypesV1.SecurityThresholds memory securityThresholds,
        bytes memory requiredQuorumNumbers
    ) internal view {
        EigenDAV3Cert calldata cert = decodeCert(certBytes);
        V2Lib.verifyDACertV2(
            eigenDAThresholdRegistry,
            eigenDASignatureVerifier,
            cert.batchHeader,
            cert.blobInclusionInfo,
            cert.nonSignerStakesAndSignature,
            securityThresholds,
            requiredQuorumNumbers,
            cert.signedQuorumNumbers
        );
    }

    function checkDACert(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDASignatureVerifier eigenDASignatureVerifier,
        bytes calldata certBytes,
        DATypesV1.SecurityThresholds memory securityThresholds,
        bytes memory requiredQuorumNumbers
    ) internal view returns (V2Lib.StatusCode, bytes memory) {
        EigenDAV3Cert calldata cert = decodeCert(certBytes);
        return V2Lib.checkDACertV2(
            eigenDAThresholdRegistry,
            eigenDASignatureVerifier,
            cert.batchHeader,
            cert.blobInclusionInfo,
            cert.nonSignerStakesAndSignature,
            securityThresholds,
            requiredQuorumNumbers,
            cert.signedQuorumNumbers
        );
    }
}
