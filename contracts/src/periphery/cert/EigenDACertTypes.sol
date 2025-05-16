// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {EigenDATypesV2 as DATypesV2} from "src/core/libraries/v2/EigenDATypesV2.sol";

/// @title EigenDACertTypes
/// @notice This library defines the types for each EigenDA certificate version.
/// @dev It is required that RBN be located in positions 32:64 (padded) in the ABI encoded certificate.
library EigenDACertTypes {
    struct EigenDACertV3 {
        DATypesV2.BatchHeaderV2 batchHeader;
        DATypesV2.BlobInclusionInfo blobInclusionInfo;
        DATypesV1.NonSignerStakesAndSignature nonSignerStakesAndSignature;
        bytes signedQuorumNumbers;
    }
}
