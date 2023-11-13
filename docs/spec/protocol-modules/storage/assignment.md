
# Assignment

The assignment functionality within EigenDA is carried out by the `AssignmentCoordinator` which is responsible for taking the current OperatorState and the security requirements represented by a given QuorumParams and determining or validating system parameters that will satisfy these security requirements given the OperatorStates. There are two classes of parameters that must be determined or validated:

1) the chunk indices that will be assigned to each DA node.
2) the length of each chunk (measured in number of symbols). In keeping with the constraint imposed by the Encoding module, all chunks must have the same length, so this is paramter is a scalar.

As illustrated in the interface that follows, the assignment of indices does not depends on the security parameters such as quorum threshold and adversary threshold. As these parameters change, the only effect on the resulting assignments will be that the chunk length changes.

The `AssignmentCoordinator` is used by the disperser to determine or validate the `EncodingParams` struct used to encode a data blob, consisting of the total number of chunks (i.e., the total number of indices) and the length of the chunk. We illustrate this in the next section.

## Interface

The AssignmentCoordinator must implement the following interface, which facilitates with the above tasks:

```go
type AssignmentCoordinator interface {

	// GetAssignments calculates the full set of node assignments. The assignment of indices to nodes depends only on the OperatorState
	// for a given quorum and the quantizationFactor. In particular, it does not depend on the security parameters.
	GetAssignments(state *OperatorState, quorumID QuorumID, quantizationFactor uint) (map[OperatorID]Assignment, AssignmentInfo, error)

	// GetOperatorAssignment calculates the assignment for a specific DA node
	GetOperatorAssignment(state *OperatorState, quorum QuorumID, quantizationFactor uint, id OperatorID) (Assignment, AssignmentInfo, error)

	// GetMinimumChunkLength calculates the minimum chunkLength that is sufficient for a given blob for each quorum
	GetMinimumChunkLength(numOperators, blobLength, quantizationFactor uint, quorumThreshold, adversaryThreshold uint8) uint

	// GetChunkLengthFromHeader calculates the chunk length from the blob header
	GetChunkLengthFromHeader(state *OperatorState, header *BlobQuorumInfo) (uint, error)
}
```

The `AssignmentCoordinator` can be used to get the `EncodingParams` struct in the following manner:

```go
    //  quorumThreshold, adversaryThreshold, blobSize, quorumID and quantizationFactor are given

    // Get assignments
	assignments, info, _ := asn.GetAssignments(state, quorumID, quantizationFactor)

    // Get minimum chunk length
	blobLength := enc.GetBlobLength(blobSize)
	numOperators := uint(len(state.Operators[quorumID]))
	chunkLength := asn.GetMinimumChunkLength(numOperators, blobLength, quantizationFactor, quorumThreshold, adversaryThreshold)

    // Get encoding params
	params, _ := enc.GetEncodingParams(chunkLength, info.TotalChunks)
```

## Standard Assignment Security Logic

The standard assignment coordinator implements a very simple logic for determining the number of chunks per node and the chunk length, which we describe here. More background concerning this design can be found in the [Design Document](../../../design/assignment.md)

**Index Assignment**.
For each operator $i$, let $S_i$ signify the amount of stake held by that operator. The number of chunks assigned to an operator $i$ is given by the equation

$$m_i = \text{ceil}\left(\frac{\rho nS_i}{\sum_j S_j}\right) \tag{1}$$

where $n$ is the total number of operators and $\rho$ is called the quantization factor.

**Chunk Length**.
To determine the chunk length, we first set the reconstruction threshold at the level of chunks (i.e., the number of chunks from which we wish to be able to reconstruct the original blob) as

$$m = ceil(n\rho\gamma)$$

where  $\gamma = \beta-\alpha$, with $\alpha$ and $\beta$ as defined in the [Storage Overview](./overview.md).

We can then derive the chunk length as $C \ge B/m$, where $B$ is the length of the blob.

**Correctness**.
Let's show that any set of operators $U_q \setminus U_a$ will have a complete blob. The amount of data held by these operators is given by

$$
\sum_{i \in U_q \setminus U_a} m_i C
$$

We first notice from (1) and from the definitions of $U_q$ and $U_a$ that

$$
\sum_{i \in U_q \setminus U_a} m_i \ge  n\rho\frac{\sum_{i \in U_q \setminus U_a} S_i}{\sum_jS_j}  = n\rho\frac{\sum_{i \in U_q} S_i - \sum_{i \in U_a} S_i}{\sum_jS_j} \ge n\rho(\beta - \alpha) = n\rho\gamma = ceil(n\rho\gamma) = m.
$$

The second from last equality follows because each $m_i$ is an integer, so the summation must be an integer. Substituting for $C$, we see that

$$
\sum_{i \in U_q \setminus U_a} m_i C \ge mC \ge B. \tag{2}
$$

Thus, the reconstruction requirement from the [Encoding](./encoding.md) module is satisfied. Notice that the final inequality of this equation can be written as

`ceil(EncodedBlobLength*(QuorumThreshold-AdversaryThreshold)) >= BlobLength`

with the following mappings:
- `EncodedBlobLength` = $n\rho C$
- `QuorumThreshold` = $\beta$
- `AdversaryThreshold` = $\alpha$
- `BlobLength` = $B$.

## Validation Actions

Validation with respect to assignments is performed at different layers of the protocol:

### DA Nodes

When the a DA node receives a `StoreChunks` request, it performs the following validation actions relative to each blob header:
- It uses `GetOperatorAssignment` to calculate the chunk indices for which it is responsible, and verifies that each of the chunks that it has recieved lies on the polynomial at these indices (see [Encoding validation actions](./encoding.md#validation-actions))
- It validates that the `Length` contained in the `BlobHeader` is valid. (see [Encoding validation actions](./encoding.md#validation-actions))
- It determines the `ChunkLength` of the received chunks.
- It validates that the `EncodedBlobLength` of the `BlobHeader` satisfies `BlobHeader.EncodedBlobLength = ChunkLength*BlobHeader.QuantizationFactor*State.NumOperators`

This step ensures that each honest node has received the blobs for which it is accountable under the [Standard Assignment Coordinator](#standard-assignment-security-logic), and that the chunk Length and quantization parameters are consistent across all of the honest DA nodes.

### Rollup Smart Contract

When the rollup confirms its blob against the EigenDA batch, it performs the following checks for each quorum

- Check that `BlobHeader.EncodedBlobLength*(BatchHeader.QuorumThreshold[quorumId] - BlobHeader.AdversaryThreshold) > BlobHeader.Length`

Together, these checks ensure that Equation (2) is satisfied.

The check by the rollup smart contract also serves to ensure that the `QuorumThreshold` for the blob is greater than the `AdversaryThreshold`. This means that if the `EncodedBlobLength` was set incorrectly by the disperser and the adversarial contingent of the DA nodes is within the specified threshold, the batch cannot be confirmed as a sufficient number of nodes will not sign.
