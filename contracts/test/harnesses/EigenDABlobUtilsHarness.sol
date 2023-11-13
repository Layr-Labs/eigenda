// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.9;

import "../../src/libraries/EigenDABlobUtils.sol";
import "forge-std/Test.sol";

contract EigenDABlobUtilsHarness is Test {    

    function verifyBlob(
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        IEigenDAServiceManager eigenDAServiceManager,
        EigenDABlobUtils.BlobVerificationProof calldata blobVerificationProof
    ) external view {
        EigenDABlobUtils.verifyBlob(blobHeader, eigenDAServiceManager, blobVerificationProof);
    }
}
