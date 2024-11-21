// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.9;

import "../../src/libraries/EigenDARollupUtils.sol";
import "forge-std/Test.sol";
import "../../src/interfaces/IEigenDAStructs.sol";

contract EigenDABlobUtilsHarness is Test {    

    function verifyBlob(
        BlobHeader calldata blobHeader,
        IEigenDAServiceManager eigenDAServiceManager,
        BlobVerificationProof calldata blobVerificationProof
    ) external view {
        EigenDARollupUtils.verifyBlob(blobHeader, eigenDAServiceManager, blobVerificationProof);
    }

    function verifyBlobs(
        BlobHeader[] calldata blobHeaders,
        IEigenDAServiceManager eigenDAServiceManager,
        BlobVerificationProof[] calldata blobVerificationProofs
    ) external view {
        EigenDARollupUtils.verifyBlobs(blobHeaders, eigenDAServiceManager, blobVerificationProofs);
    }
}
