package auth

import (
	"github.com/Layr-Labs/eigenda/api/grpc/common"
	v2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	grpc "github.com/Layr-Labs/eigenda/api/grpc/node/v2"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"testing"
)

func randomStoreChunksRequest(rand *random.TestRandom) *grpc.StoreChunksRequest {
	certificateCount := rand.Intn(10) + 1
	blobCertificates := make([]*v2.BlobCertificate, certificateCount)
	for i := 0; i < certificateCount; i++ {

		relayCount := rand.Intn(10) + 1
		relays := make([]uint32, relayCount)
		for j := 0; j < relayCount; j++ {
			relays[j] = rand.Uint32()
		}

		quorumCount := rand.Intn(10) + 1
		quorumNumbers := make([]uint32, quorumCount)
		for j := 0; j < quorumCount; j++ {
			quorumNumbers[j] = rand.Uint32()
		}

		blobCertificates[i] = &v2.BlobCertificate{
			BlobHeader: &v2.BlobHeader{
				Version:       rand.Uint32(),
				QuorumNumbers: quorumNumbers,
				Commitment: &common.BlobCommitment{
					Commitment:       rand.Bytes(32),
					LengthCommitment: rand.Bytes(32),
					LengthProof:      rand.Bytes(32),
					Length:           rand.Uint32(),
				},
				PaymentHeader: &common.PaymentHeader{
					AccountId:         rand.String(32),
					ReservationPeriod: rand.Uint32(),
					CumulativePayment: rand.Bytes(32),
					Salt:              rand.Uint32(),
				},
				Signature: rand.Bytes(32),
			},
			Relays: relays,
		}
	}

	return &grpc.StoreChunksRequest{
		Batch: &v2.Batch{
			Header: &v2.BatchHeader{
				BatchRoot:            rand.Bytes(32),
				ReferenceBlockNumber: rand.Uint64(),
			},
			BlobCertificates: blobCertificates,
		},
		DisperserID: rand.Uint32(),
		Signature:   rand.Bytes(32),
	}
}

func TestHashing(t *testing.T) {
	rand := random.NewTestRandom(t)

	request := randomStoreChunksRequest(rand)
	originalRequestHash := HashStoreChunksRequest(request)

	// modifying the signature should not change the hash
	request.Signature = rand.Bytes(32)
	hash := HashStoreChunksRequest(request)
	require.Equal(t, originalRequestHash, hash)

	// modify the disperser id
	rand.Reset()
	request = randomStoreChunksRequest(rand)
	request.DisperserID = request.DisperserID + 1
	hash = HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// remove a blob cert
	rand.Reset()
	request = randomStoreChunksRequest(rand)
	request.Batch.BlobCertificates = request.Batch.BlobCertificates[:len(request.Batch.BlobCertificates)-1]
	hash = HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify a relay
	rand.Reset()
	request = randomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].Relays[0] = request.Batch.BlobCertificates[0].Relays[0] + 1
	hash = HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, remove a relay
	rand.Reset()
	request = randomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].Relays =
		request.Batch.BlobCertificates[0].Relays[:len(request.Batch.BlobCertificates[0].Relays)-1]
	hash = HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, add a relay
	rand.Reset()
	request = randomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].Relays = append(request.Batch.BlobCertificates[0].Relays, rand.Uint32())
	hash = HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify a quorum number
	rand.Reset()
	request = randomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers[0] =
		request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers[0] + 1
	hash = HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, remove a quorum number
	rand.Reset()
	request = randomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers =
		request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers[:len(
			request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers)-1]
	hash = HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, add a quorum number
	rand.Reset()
	request = randomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers = append(
		request.Batch.BlobCertificates[0].BlobHeader.QuorumNumbers, rand.Uint32())
	hash = HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the Commitment.Commitment
	rand.Reset()
	request = randomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.Commitment.Commitment = rand.Bytes(32)
	hash = HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the Commitment.LengthCommitment
	rand.Reset()
	request = randomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.Commitment.LengthCommitment = rand.Bytes(32)
	hash = HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the Commitment.LengthProof
	rand.Reset()
	request = randomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.Commitment.LengthProof = rand.Bytes(32)
	hash = HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the Commitment.Length
	rand.Reset()
	request = randomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.Commitment.Length = rand.Uint32()
	hash = HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the PaymentHeader.AccountId
	rand.Reset()
	request = randomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.PaymentHeader.AccountId = rand.String(32)
	hash = HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the PaymentHeader.ReservationPeriod
	rand.Reset()
	request = randomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.PaymentHeader.ReservationPeriod = rand.Uint32()
	hash = HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the PaymentHeader.CumulativePayment
	rand.Reset()
	request = randomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.PaymentHeader.CumulativePayment = rand.Bytes(32)
	hash = HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the PaymentHeader.Salt
	rand.Reset()
	request = randomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.PaymentHeader.Salt = rand.Uint32()
	hash = HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)

	// within a blob cert, modify the Signature
	rand.Reset()
	request = randomStoreChunksRequest(rand)
	request.Batch.BlobCertificates[0].BlobHeader.Signature = rand.Bytes(32)
	hash = HashStoreChunksRequest(request)
	require.NotEqual(t, originalRequestHash, hash)
}

func TestRequestSigning(t *testing.T) {
	rand := random.NewTestRandom(t)

	public, private := rand.ECDSA()
	publicAddress := crypto.PubkeyToAddress(*public)

	request := randomStoreChunksRequest(rand)

	signature, err := SignStoreChunksRequest(private, request)
	require.NoError(t, err)

	err = VerifyStoreChunksRequest(publicAddress, request, signature)
	require.NoError(t, err)

	// Adding the signature to the request should not change the hash, so it should still be valid
	request.Signature = signature
	err = VerifyStoreChunksRequest(publicAddress, request, signature)
	require.NoError(t, err)

	// Using a different public key should make the signature invalid
	otherPublic, _ := rand.ECDSA()
	otherPublicAddress := crypto.PubkeyToAddress(*otherPublic)
	err = VerifyStoreChunksRequest(otherPublicAddress, request, signature)
	require.Error(t, err)

	// Changing a byte in the signature should make it invalid
	alteredSignature := make([]byte, len(signature))
	copy(alteredSignature, signature)
	alteredSignature[0] = alteredSignature[0] + 1
	err = VerifyStoreChunksRequest(publicAddress, request, alteredSignature)
	require.Error(t, err)

	// Changing a field in the request should make it invalid
	request.DisperserID = request.DisperserID + 1
	err = VerifyStoreChunksRequest(publicAddress, request, signature)
	require.Error(t, err)
}
