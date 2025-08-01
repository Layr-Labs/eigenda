package auth

import (
	"testing"

	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestHashing(t *testing.T) {
	rand := random.NewTestRandom()

	request := RandomStoreChunksRequest(rand)
	originalRequestHash, err := hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)

	// modifying the signature should not change the hash
	request.Signature = rand.Bytes(32)
	hash, err := hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.Equal(t, originalRequestHash, hash)

	// modify the disperser id
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.DisperserID = request.GetDisperserID() + 1
	hash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.NotEqual(t, originalRequestHash, hash)

	// remove a blob cert
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates = request.GetBatch().GetBlobCertificates()[:len(request.GetBatch().GetBlobCertificates())-1]
	hash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify a relay
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].RelayKeys[0] = request.GetBatch().GetBlobCertificates()[0].GetRelayKeys()[0] + 1
	hash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, remove a relay
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].RelayKeys =
		request.GetBatch().GetBlobCertificates()[0].GetRelayKeys()[:len(request.GetBatch().GetBlobCertificates()[0].GetRelayKeys())-1]
	hash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, add a relay
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].RelayKeys = append(request.Batch.BlobCertificates[0].RelayKeys, rand.Uint32())
	hash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify a quorum number
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers[0] =
		request.GetBatch().GetBlobCertificates()[0].GetBlobHeader().GetQuorumNumbers()[0] + 1
	hash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, remove a quorum number
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers =
		request.GetBatch().GetBlobCertificates()[0].GetBlobHeader().GetQuorumNumbers()[:len(
			request.GetBatch().GetBlobCertificates()[0].GetBlobHeader().GetQuorumNumbers())-1]
	hash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, add a quorum number
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers = append(
		request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers, rand.Uint32())
	hash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the Commitment.Commitment
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.Commitment.Commitment = rand.Bytes(32)
	hash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the Commitment.LengthCommitment
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.Commitment.LengthCommitment = rand.Bytes(32)
	hash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the Commitment.LengthProof
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.Commitment.LengthProof = rand.Bytes(32)
	hash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the Commitment.Length
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.Commitment.Length = rand.Uint32()
	hash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the PaymentHeader.AccountId
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.PaymentHeader.AccountId = rand.String(32)
	hash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the PaymentHeader.Timestamp
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.PaymentHeader.Timestamp = rand.Time().UnixMicro()
	hash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the PaymentHeader.CumulativePayment
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.PaymentHeader.CumulativePayment = rand.Bytes(32)
	hash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the Signature
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].Signature = rand.Bytes(32)
	hash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.NotEqual(t, originalRequestHash, hash)

	// nil header
	request = RandomStoreChunksRequest(rand)
	request.Batch.Header = nil
	hash, err = hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.NotEqual(t, originalRequestHash, hash)
}

func TestRequestSigning(t *testing.T) {
	rand := random.NewTestRandom()

	public, private, err := rand.ECDSA()
	require.NoError(t, err)
	publicAddress := crypto.PubkeyToAddress(*public)

	request := RandomStoreChunksRequest(rand)

	signature, err := SignStoreChunksRequest(private, request)
	require.NoError(t, err)
	request.Signature = signature

	hash, err := VerifyStoreChunksRequest(publicAddress, request)
	require.NoError(t, err)
	expectedHash, err := hashing.HashStoreChunksRequest(request)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)

	// Using a different public key should make the signature invalid
	otherPublic, _, err := rand.ECDSA()
	require.NoError(t, err)
	otherPublicAddress := crypto.PubkeyToAddress(*otherPublic)
	_, err = VerifyStoreChunksRequest(otherPublicAddress, request)
	require.Error(t, err)

	// Changing a byte in the signature should make it invalid
	alteredSignature := make([]byte, len(signature))
	copy(alteredSignature, signature)
	alteredSignature[0] = alteredSignature[0] + 1
	request.Signature = alteredSignature
	_, err = VerifyStoreChunksRequest(publicAddress, request)
	require.Error(t, err)

	// Changing a field in the request should make it invalid
	request.DisperserID = request.GetDisperserID() + 1
	request.Signature = signature
	_, err = VerifyStoreChunksRequest(publicAddress, request)
	require.Error(t, err)
}
