
# Encoding

## Overall Requirements

Within EigenDA, blobs are encoded so that they can be scalably distributed among the DA nodes. The EigenDA encoding module is designed to meet the following security requirements:
1. Adversarial tolerance for DA nodes: We need to have tolerance to arbitrary adversarial behavior by DA nodes up to some threshold, which is discussed in other sections. Note that simple sharding approaches such as duplicating slices of the blob data have good tolerance to random node dropout, but poor tolerance to worst-case adversarial behavior.
2. Adversarial tolerance for disperser: We do not want to put trust assumptions on the encoder or rely on fraud proofs to detect if an encoding is done incorrectly.

## Interfaces

From a system standpoint, the encoder module must satisfy the following interface, which codifies requirements 1 and 2 above.

```go
type EncodingParams struct {
	ChunkLength uint // ChunkSize is the length of the chunk in symbols
	NumChunks   uint
}

// Encoder is responsible for encoding, decoding, and chunk verification
type Encoder interface {

	// GetBlobLength converts from blob size in bytes to blob size in symbols. This is necessary because different encoder backends
	// may use different symbol sizes
	GetBlobLength(blobSize uint) uint

	// GetEncodingParams takes in the minimum chunk length and the minimum number of chunks and returns the encoding parameters given any
	// additional constraints from the encoder backend. For instance, both the ChunkLength and NumChunks must typically be powers of 2.
	// The ChunkLength returned here should be used in constructing the BlobHeader.
	GetEncodingParams(minChunkLength, minNumChunks uint) (EncodingParams, error)

	// Encode takes in a blob and returns the commitments and encoded chunks. The encoding will satisfy the property that
	// for any number M such that M*params.ChunkLength > BlobCommitments.Length, then any set of M chunks will be sufficient to
	// reconstruct the blob.
	Encode(data [][]byte, params EncodingParams) (BlobCommitments, []*Chunk, error)

	// VerifyChunks takes in the chunks, indices, commitments, and encoding parameters and returns an error if the chunks are invalid
	// VerifyChunks also verifies that the ChunkLength contained in the BlobCommitments is consistent with the LengthProof.
	VerifyChunks(chunks []*Chunk, indices []ChunkNumber, commitments BlobCommitments, params EncodingParams) error

	// Decode takes in the chunks, indices, and encoding parameters and returns the decoded blob
	Decode(chunks []*Chunk, indices []ChunkNumber, params EncodingParams, inputSize uint64) ([]byte, error)
}
```

Notice that these interfaces only support a global chunk size across all the encoded chunks for a given encoded blob. This constraint derives mostly from the design of the [KZG FFT Encoder Backend](#the-kzg-fft-encoder-backend) which generates multireveal proofs in an amortized fashion.


## Trustless Encoding via KZG and Reed-Solomon

EigenDA uses a combination of Reed-Solomon (RS) erasure coding and KZG polynomial commitments to perform trustless  encoding. In this section, we provide a high level overview of how the EigenDA encoding module works and how it achieves these properties.

### Basic Reed Solomon Encoding

Basic RS encoding is used to achieve the first requirement of tolerance to adversarial node behavior. This looks like the following:

1. The blob data is represented as a string of symbols, where each symbol is elements in a certain finite field. The number of symbols is called the `BlobLength`
2. These symbols are interpreted as the coefficients of a `BlobLength`-1 degree polynomial.
3. This polynomial is evaluated at `NumChunks`*`ChunkLength` distinct indices.
4. Chunks are constructed, where each chunk consists of the polynomial evaluations at `ChunkLength` distinct indices.

Notice that given any number of chunks $M$ such that $M$*`ChunkLength` > `BlobLength`, via [polynomial interpolation](https://en.wikipedia.org/wiki/Polynomial_interpolation) it is possible to reconstruct the original polynomial, and therefore its coefficients which represent the original blob. Thus, this basic RS encoding scheme satisfies the requirement of the `Encoder.Encode` interface.

### Validation via KZG

Without modification, RS encoding has the following important problem: Suppose that a user asks an untrusted disperser to encode data and send it to the nodes. Even if the user is satisfied that each node has received some chunk of data, there is no way to know how the disperser went about constructing those chunks. [KZG polynomial commitments](https://dankradfeist.de/ethereum/2020/06/16/kate-polynomial-commitments.html) provide a solution to this problem.

#### Encoded Chunk Verification
KZG commitments provide three important primitives, for a polynomial $p(X) = \sum_{i}c_iX^i$:
- `commit(p(X))` returns a `Commitment` which is used to identify the polynomial.
- `prove(p(X),indices)` returns a `Proof` which can be used to verify that a set of evaluations lies on the polynomial.
- `verify(Commitment,Proof,evals,indices)` returns a `bool` indicating whether the committed polynomial evaluates to `evals` and the provided `indices`.

#### Blob Size Verification
KZG commitments also can be used to verify the degree of the original polynomial, which in turn corresponds to the size of the encoded blob.

The KZG commitment relies on a structured random string (SRS) containing a generator point $G$ multiplied by all of the powers of some secret field element $\tau$, up to some maximum power $n$. This means that it is not possible to use this SRS to commit to a polynomial of degree greater than $n$. A consequence of this is that if $p(x)$ is a polynomial of degree greater than $m$, it will not be possible to commit to the polynomial $x^{n-m}p(x)$.

The scheme thus works as follows: If the disperser wishes to claim that the polynomial $p(x)$ is of degree less than or equal to $m$, they must provide along with the commitment $C_1$ to $p$, a commitment $C_2$ to $q(x) = x^{n-m}p(x)$. A verifier can request the disperser to open both polynomials at a random point $y$, verify the values of $p(y)$ and $q(y)$, and then check that $q(y) = y^{n-m}p(y)$. If these checks pass, the verifier knows that 1) $deg(q) = deg(p) + n - m$, 2) the disperser was able to make a commitment to $q$, and so $deg(q) \le n$, and therefore 3), $deg(p) \le m$. In practice, this protocol can be made non-interactive using the Fiat-Shamir heuristic.

Note: The blob length verification here allows for the blob length to be upper-bounded; it cannot be used to prove the exact blob length.

## The KZG-FFT Encoder Backend

A design for an efficient encoding backend that makes use of amortized kzg multi-reveal generation is described [here](../../../design/encoding.md).

## Validation Actions

When a DA node receives a `StoreChunks` request, it performs the following validation actions relative to each blob header:
- It uses `GetOperatorAssignment` of the `AssignmentCoordinator` interface to calculate the chunk indices for which it is responsible and the total number of chunks, `TotalChunks`.
- It instantiates an encoder using the `ChunkLength` from the `BlobHeader` and the `TotalChunks`, and uses `VerifyChunks` to verify that the data contained within the chunks lies on the committed polynomial at the correct indices.
- The `VerifyChunks` method also verifies that the `Length` contained in the `BlobCommitments` struct is valid based on the `LengthProof`.
