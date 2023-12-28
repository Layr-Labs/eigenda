# Data Model

### Quorum Information

```go
// QuorumID is a unique identifier for a quorum; initially EigenDA will support up to 256 quorums
type QuorumID = uint8

// SecurityParam contains the quorum ID and the adversary threshold for the quorum;
type SecurityParam struct {
	QuorumID QuorumID
	// AdversaryThreshold is the maximum amount of stake that can be controlled by an adversary in the quorum as a percentage of the total stake in the quorum
	AdversaryThreshold uint8
	// QuorumThreshold is the amount of stake that must sign a message for it to be considered valid as a percentage of the total stake in the quorum
	QuorumThreshold uint8 `json:"quorum_threshold"`
}

// QuorumResult contains the quorum ID and the amount signed for the quorum
type QuorumResult struct {
	QuorumID QuorumID
	// PercentSigned is percentage of the total stake for the quorum that signed for a particular batch.
	PercentSigned uint8
}
```

### Blob Requests

```go
// BlobHeader contains the original data size of a blob and the security required
type BlobRequestHeader struct {
	// BlobSize is the size of the original data in bytes
	BlobSize uint32
	// For a blob to be accepted by EigenDA, it satisfies the AdversaryThreshold of each quorum contained in SecurityParams
	SecurityParams []SecurityParam
}
```

### Data Headers

```go
type BlobHeader struct {
	BlobCommitments
	// QuorumInfos contains the quorum specific parameters for the blob
	QuorumInfos []*BlobQuorumInfo
}

// BlobQuorumInfo contains the quorum IDs and parameters for a blob specific to a given quorum
type BlobQuorumInfo struct {
	SecurityParam
	// ChunkLength is the number of symbols in a chunk
	ChunkLength uint
}

// BlomCommitments contains the blob's commitment, degree proof, and the actual degree.
type BlobCommitments struct {
	Commitment  *Commitment
	LengthProof *Commitment
	Length      uint
}

// BatchHeader contains the metadata associated with a Batch for which DA nodes must attest; DA nodes sign on the hash of the batch header
type BatchHeader struct {
	// BlobHeaders contains the headers of the blobs in the batch
	BlobHeaders []*BlobHeader
	// QuorumResults contains the quorum parameters for each quorum that must sign the batch; all quorum parameters must be satisfied
	// for the batch to be considered valid
	QuorumResults []QuorumResult
	// ReferenceBlockNumber is the block number at which all operator information (stakes, indexes, etc.) is taken from
	ReferenceBlockNumber uint
	// BatchRoot is the root of a Merkle tree whose leaves are the hashes of the blobs in the batch
	BatchRoot [32]byte
}
```

### Encoded Data Products

```go
// EncodedBatch is a container for a batch of blobs. DA nodes receive and attest to the blobs in a batch together to amortize signature verification costs
type EncodedBatch struct {
	ChunkBatches map[OperatorID]ChunkBatch
}

// Chunks

// Chunk is the smallest unit that is distributed to DA nodes, including both data and the associated polynomial opening proofs.
// A chunk corresponds to a set of evaluations of the global polynomial whose coefficients are used to construct the blob Commitment.
type Chunk struct {
	// The Coeffs field contains the coefficients of the polynomial which interpolates these evaluations. This is the same as the
	// interpolating polynomial, I(X), used in the KZG multi-reveal (https://dankradfeist.de/ethereum/2020/06/16/kate-polynomial-commitments.html#multiproofs)
	Coeffs []Symbol
	Proof  Proof
}

func (c *Chunk) Length() int {
	return len(c.Coeffs)
}

// ChunkBatch is the collection of chunks associated with a single operator and a single batch.
type ChunkBatch struct {
	// Bundles contains the chunks associated with each blob in the batch; each bundle contains the chunks associated with a single blob
	// The number of bundles should be equal to the total number of blobs in the batch. The number of chunks per bundle will vary
	Bundles [][]*Chunk
}
```

### DA Node State

```go
type StakeAmount *big.Int

// OperatorInfo contains information about an operator which is stored on the blockchain state,
// corresponding to a particular quorum
type OperatorInfo struct {
	// Stake is the amount of stake held by the operator in the quorum
	Stake StakeAmount
	// Index is the index of the operator within the quorum
	Index OperatorIndex
}

// OperatorState contains information about the current state of operators which is stored in the blockchain state
type OperatorState struct {
	// Operators is a map from quorum ID to a map from the operators in that quorum to their StoredOperatorInfo. Membership
	// in the map implies membership in the quorum.
	Operators map[QuorumID]map[OperatorID]*OperatorInfo
	// Totals is a map from quorum ID to the total stake (Stake) and total count (Index) of all operators in that quorum
	Totals map[QuorumID]*OperatorInfo
	// BlockNumber is the block number at which this state was retrieved
	BlockNumber uint
}
```
