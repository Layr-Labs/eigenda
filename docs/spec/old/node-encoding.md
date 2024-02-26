
# Blob Encoding Constraints

The DA nodes expect for blobs to correspond to evaluations of a polynomial of a certain degree. The blob payload delivered to the node contains a KZG polynomial commitment identifying the polynomial, as well as a separate commitment allowing the node to verify its degree. The payload also contains KZG reveal proofs allowing the node to verify that its received data corresponds to evaluations of the polynomial at a specific evaluation index. This document describes in detail how the node performs all verifications, including calculating the evaluation indices.

Based on the `referenceBlockNumber` contained in the [`DataStoreHeader`](./types/node-types.md#datastoreheader) structure, the DA node construct a [`StateView`](./types/node-types.md#stateview) object using the [Indexer](./node.md#indexer) service.

The operator node will then perform the following checks for each quorum they are a part of to ensure that the `StoreChunksRequest` is valid: 
1. Confirm Header Consistency: 
    - Generate the encoding params `numPar` and `numSys` (from `AdversaryThresholdBPs`, `QuorumThresholdBPs` and `stateView`)
    - Check that `degree` properly derives from `origDataSize` and `numSys`
2. Confirm Degree of KZG Commitment: Use the `lowDegreeProof` to verify that `kzgCommitment` commits to a polynomial of degree less than or equal to `numSysE`*`degreeE`.
3. Verify frame for the quorum against KZG Commitment: Use the `headerHash` to determine the proper indices for the chunks held by the operator, and then use the multireveal proofs contained in the chunks to verify each chunk against `kzgCommitment` commit.

Indices are determined using the following formula:

