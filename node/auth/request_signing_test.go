package auth

import (
	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHashing(t *testing.T) {
	rand := random.NewTestRandom(t)

	request := RandomStoreChunksRequest(rand)
	originalRequestHash := hashing.HashStoreChunksRequest(request)

	// modifying the signature should not change the hash
	request.Signature = rand.Bytes(32)
	hash := hashing.HashStoreChunksRequest(request)
	require.Equal(t, originalRequestHash, hash)

	// modify the disperser id
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.DisperserID = request.DisperserID + 1
	hash = hashing.HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// remove a blob cert
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates = request.Batch.BlobCertificates[:len(request.Batch.BlobCertificates)-1]
	hash = hashing.HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify a relay
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].Relays[0] = request.Batch.BlobCertificates[0].Relays[0] + 1
	hash = hashing.HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, remove a relay
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].Relays =
		request.Batch.BlobCertificates[0].Relays[:len(request.Batch.BlobCertificates[0].Relays)-1]
	hash = hashing.HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, add a relay
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].Relays = append(request.Batch.BlobCertificates[0].Relays, rand.Uint32())
	hash = hashing.HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify a quorum number
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers[0] =
		request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers[0] + 1
	hash = hashing.HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, remove a quorum number
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers =
		request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers[:len(
			request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers)-1]
	hash = hashing.HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, add a quorum number
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers = append(
		request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers, rand.Uint32())
	hash = hashing.HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the Commitment.Commitment
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.Commitment.Commitment = rand.Bytes(32)
	hash = hashing.HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the Commitment.LengthCommitment
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.Commitment.LengthCommitment = rand.Bytes(32)
	hash = hashing.HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the Commitment.LengthProof
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.Commitment.LengthProof = rand.Bytes(32)
	hash = hashing.HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the Commitment.Length
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.Commitment.Length = rand.Uint32()
	hash = hashing.HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the PaymentHeader.AccountId
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.PaymentHeader.AccountId = rand.String(32)
	hash = hashing.HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the PaymentHeader.ReservationPeriod
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.PaymentHeader.ReservationPeriod = rand.Uint32()
	hash = hashing.HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the PaymentHeader.CumulativePayment
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.PaymentHeader.CumulativePayment = rand.Bytes(32)
	hash = hashing.HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the PaymentHeader.Salt
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.PaymentHeader.Salt = rand.Uint32()
	hash = hashing.HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the Signature
	rand.Reset()
	request = RandomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.Signature = rand.Bytes(32)
	hash = hashing.HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)
}

func TestRequestSigning(t *testing.T) {
	rand := random.NewTestRandom(t)

	public, private := rand.ECDSA()
	publicAddress := crypto.PubkeyToAddress(*public)

	request := RandomStoreChunksRequest(rand)

	signature, err := SignStoreChunksRequest(private, request)
	require.NoError(t, err)
	request.Signature = signature

	err = VerifyStoreChunksRequest(publicAddress, request)
	require.NoError(t, err)

	// Using a different public key should make the signature invalid
	otherPublic, _ := rand.ECDSA()
	otherPublicAddress := crypto.PubkeyToAddress(*otherPublic)
	err = VerifyStoreChunksRequest(otherPublicAddress, request)
	require.Error(t, err)

	// Changing a byte in the signature should make it invalid
	alteredSignature := make([]byte, len(signature))
	copy(alteredSignature, signature)
	alteredSignature[0] = alteredSignature[0] + 1
	request.Signature = alteredSignature
	err = VerifyStoreChunksRequest(publicAddress, request)
	require.Error(t, err)

	// Changing a field in the request should make it invalid
	request.DisperserID = request.DisperserID + 1
	request.Signature = signature
	err = VerifyStoreChunksRequest(publicAddress, request)
	require.Error(t, err)
}
