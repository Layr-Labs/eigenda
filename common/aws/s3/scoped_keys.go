package s3

import (
	"fmt"

	corev2 "github.com/Layr-Labs/eigenda/core/v2"
)

const (
	// prefixLength is the number of characters to use from the base key to form the prefix.
	// Assuming keys take the form of a random hash in hex, 3 will yield 16^3 = 4096 possible prefixes.
	// This is currently hard coded because it is not expected to change, and it would require migration
	// to change it that we have not yet implemented.
	prefixLength = 3

	// blobNamespace is the namespace for a blob key.
	blobNamespace = "blob"

	// chunkNamespace is the namespace for a chunk key.
	chunkNamespace = "chunk"

	// proofNamespace is the namespace for a proof key.
	proofNamespace = "proof"
)

// ScopedKey returns a key that is scoped to a "namespace". Keys take the form of "prefix/namespace/baseKey".
// Although there is no runtime enforcement, neither the base key nor the namespace should contain any
// non-alphanumeric characters.
func ScopedKey(namespace string, baseKey string, prefixLength int) string {
	var prefix string
	if prefixLength > len(baseKey) {
		prefix = baseKey
	} else {
		prefix = baseKey[:prefixLength]
	}

	return fmt.Sprintf("%s/%s/%s", prefix, namespace, baseKey)
}

// ScopedBlobKey returns a key scoped to the blob namespace. Used to name files containing blobs in S3.
// A key scoped for blobs will never collide with a key scoped for chunks or proofs.
func ScopedBlobKey(blobKey corev2.BlobKey) string {
	return ScopedKey(blobNamespace, blobKey.Hex(), prefixLength)
}

// ScopedChunkKey returns a key scoped to the chunk namespace. Used to name files containing chunks in S3.
// A key scoped for chunks will never collide with a key scoped for blobs or proofs.
func ScopedChunkKey(blobKey corev2.BlobKey) string {
	return ScopedKey(chunkNamespace, blobKey.Hex(), prefixLength)
}

// ScopedProofKey returns a key scoped to the proof namespace. Used to name files containing proofs in S3.
// A key scoped for proofs will never collide with a key scoped for blobs or chunks.
func ScopedProofKey(blobKey corev2.BlobKey) string {
	return ScopedKey(proofNamespace, blobKey.Hex(), prefixLength)
}
