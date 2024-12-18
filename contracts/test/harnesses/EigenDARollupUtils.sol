
pragma solidity ^0.8.9;


import {IEigenDAServiceManager} from "../../src/interfaces/IEigenDAServiceManager.sol";
import {IDummyServiceManager} from "./EigenDAServiceManager.sol";

/**
 * @title Library of functions to be used by smart contracts wanting to prove blobs on EigenDA and open KZG commitments.
 * @author Layr Labs, Inc.
 */
library DummyRollupUtils {

    // STRUCTS
    struct BlobVerificationProof {
        uint32 batchId;
        uint32 blobIndex;
        IEigenDAServiceManager.BatchMetadata batchMetadata;
        bytes inclusionProof;
        bytes quorumIndices;
    }
    
    /**
     * @notice Verifies the inclusion of a blob within a batch confirmed in `eigenDAServiceManager` and its trust assumptions
     * @param blobHeader the header of the blob containing relevant attributes of the blob
     * @param eigenDAServiceManager the contract in which the batch was confirmed 
     * @param blobVerificationProof the relevant data needed to prove inclusion of the blob and that the trust assumptions were as expected
     */
    function verifyBlob(
        IEigenDAServiceManager.BlobHeader calldata blobHeader,
        IDummyServiceManager eigenDAServiceManager,
        BlobVerificationProof calldata blobVerificationProof
    ) public {
        if (!eigenDAServiceManager.verifyReturn()) {
            revert("Rollup blob verification failed");
        }
    }

}
