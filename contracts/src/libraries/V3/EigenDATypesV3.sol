// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV2 as DATypesV2} from "src/libraries/V2/EigenDATypesV2.sol";
import {NonSignerStakesAndSignature} from "src/interfaces/IEigenDAStructs.sol";

struct EigenDAV3Cert {
    DATypesV2.BatchHeaderV2 batchHeader;
    DATypesV2.BlobInclusionInfo blobInclusionInfo;
    NonSignerStakesAndSignature nonSignerStakesAndSignature;
    bytes signedQuorumNumbers;
}
