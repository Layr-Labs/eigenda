// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDACertTypes as CT} from "src/periphery/cert/EigenDACertTypes.sol";
import {EigenDATypesV2 as DATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {EigenDACertVerificationV2Lib as V2Lib} from "src/periphery/cert/v2/EigenDACertVerificationV2Lib.sol";
import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";
import {OperatorStateRetriever} from "lib/eigenlayer-middleware/src/OperatorStateRetriever.sol";

library EigenDACertVerificationV3Lib {
    function decodeCert(bytes calldata data) internal pure returns (CT.EigenDACertV3 memory cert) {
        return abi.decode(data, (CT.EigenDACertV3));
    }

    function verifyDACert(
        IEigenDAThresholdRegistry eigenDAThresholdRegistry,
        IEigenDASignatureVerifier eigenDASignatureVerifier,
        bytes calldata certBytes,
        DATypesV1.SecurityThresholds memory securityThresholds,
        bytes memory requiredQuorumNumbers
    ) internal view {
        CT.EigenDACertV3 memory cert = decodeCert(certBytes);
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
        CT.EigenDACertV3 memory cert = decodeCert(certBytes);
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
