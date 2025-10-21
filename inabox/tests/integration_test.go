package integration_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	disperserpb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	certTypes "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV1"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding/codec"
	integration "github.com/Layr-Labs/eigenda/inabox/tests"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
)

func TestEndToEndScenario(t *testing.T) {
	// Create a fresh test harness for this test
	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err, "Failed to create test context")
	defer testHarness.Cleanup()

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*15)
	defer cancel()

	privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded"
	signer := auth.NewLocalBlobRequestSigner(privateKeyHex)

	disp, err := clients.NewDisperserClient(&clients.Config{
		Hostname: "localhost",
		Port:     "32001",
		Timeout:  10 * time.Second,
	}, signer)
	require.NoError(t, err)
	require.NotNil(t, disp)

	data := make([]byte, 1024)
	_, err = rand.Read(data)
	require.NoError(t, err)

	paddedData := codec.ConvertByPaddingEmptyByte(data)

	blobStatus1, key1, err := disp.DisperseBlob(ctx, paddedData, []uint8{})
	require.NoError(t, err)
	require.NotNil(t, key1)
	require.NotNil(t, blobStatus1)
	require.Equal(t, disperser.Processing, *blobStatus1)

	blobStatus2, key2, err := disp.DisperseBlobAuthenticated(ctx, paddedData, []uint8{})
	require.NoError(t, err)
	require.NotNil(t, key2)
	require.NotNil(t, blobStatus2)
	require.Equal(t, disperser.Processing, *blobStatus2)

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	var reply1 *disperserpb.BlobStatusReply
	var reply2 *disperserpb.BlobStatusReply

	for loop := true; loop; {
		select {
		case <-ctx.Done():
			t.Fatal("timed out")
		case <-ticker.C:
			reply1, err = disp.GetBlobStatus(ctx, key1)
			require.NoError(t, err)
			require.NotNil(t, reply1)
			blobStatus1, err = disperser.FromBlobStatusProto(reply1.GetStatus())
			require.NoError(t, err)

			reply2, err = disp.GetBlobStatus(ctx, key2)
			require.NoError(t, err)
			require.NotNil(t, reply2)
			blobStatus2, err = disperser.FromBlobStatusProto(reply2.GetStatus())
			require.NoError(t, err)

			if *blobStatus1 != disperser.Confirmed || *blobStatus2 != disperser.Confirmed {
				integration.MineAnvilBlocks(t, testHarness.RPCClient, testHarness.NumConfirmations+1)
				continue
			}
			blobHeader := blobHeaderFromProto(reply1.GetInfo().GetBlobHeader())
			verificationProof := blobVerificationProofFromProto(reply1.GetInfo().GetBlobVerificationProof())
			err = testHarness.EigenDACertVerifierV1.VerifyDACertV1(&bind.CallOpts{}, blobHeader, verificationProof)
			require.NoError(t, err)
			integration.MineAnvilBlocks(t, testHarness.RPCClient, testHarness.NumConfirmations+1)

			blobHeader = blobHeaderFromProto(reply2.GetInfo().GetBlobHeader())
			verificationProof = blobVerificationProofFromProto(reply2.GetInfo().GetBlobVerificationProof())
			err = testHarness.EigenDACertVerifierV1.VerifyDACertV1(&bind.CallOpts{}, blobHeader, verificationProof)
			require.NoError(t, err)
			loop = false
		}
	}
	require.Equal(t, disperser.Confirmed, *blobStatus1)
	require.Equal(t, disperser.Confirmed, *blobStatus2)

	ctx, cancel = context.WithTimeout(t.Context(), time.Second*5)
	defer cancel()
	retrieved, err := testHarness.RetrievalClient.RetrieveBlob(ctx,
		[32]byte(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
		reply1.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
		uint(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber()),
		[32]byte(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
		0, // retrieve blob 1 from quorum 0
	)
	require.NoError(t, err)

	restored := codec.RemoveEmptyByteFromPaddedBytes(retrieved)
	require.Equal(t, bytes.TrimRight(data, "\x00"), bytes.TrimRight(restored, "\x00"))

	_, err = testHarness.RetrievalClient.RetrieveBlob(ctx,
		[32]byte(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
		reply1.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
		uint(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber()),
		[32]byte(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
		1, // retrieve blob 1 from quorum 1
	)
	require.NoError(t, err)

	_, err = testHarness.RetrievalClient.RetrieveBlob(ctx,
		[32]byte(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
		reply1.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
		uint(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber()),
		[32]byte(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
		2, // retrieve blob 1 from quorum 2
	)
	require.Error(t, err)

	retrieved, err = testHarness.RetrievalClient.RetrieveBlob(ctx,
		[32]byte(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
		reply2.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
		uint(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber()),
		[32]byte(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
		0, // retrieve from quorum 0
	)
	require.NoError(t, err)
	restored = codec.RemoveEmptyByteFromPaddedBytes(retrieved)
	require.Equal(t, bytes.TrimRight(data, "\x00"), bytes.TrimRight(restored, "\x00"))
	_, err = testHarness.RetrievalClient.RetrieveBlob(ctx,
		[32]byte(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
		reply2.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
		uint(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber()),
		[32]byte(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
		1, // retrieve from quorum 1
	)
	require.NoError(t, err)
	_, err = testHarness.RetrievalClient.RetrieveBlob(ctx,
		[32]byte(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
		reply2.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
		uint(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber()),
		[32]byte(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
		2, // retrieve from quorum 2
	)
	require.Error(t, err)
}

func blobHeaderFromProto(blobHeader *disperserpb.BlobHeader) certTypes.EigenDATypesV1BlobHeader {
	quorums := make([]certTypes.EigenDATypesV1QuorumBlobParam, len(blobHeader.GetBlobQuorumParams()))
	for i, quorum := range blobHeader.GetBlobQuorumParams() {
		quorums[i] = certTypes.EigenDATypesV1QuorumBlobParam{
			QuorumNumber:                    uint8(quorum.GetQuorumNumber()),
			AdversaryThresholdPercentage:    uint8(quorum.GetAdversaryThresholdPercentage()),
			ConfirmationThresholdPercentage: uint8(quorum.GetConfirmationThresholdPercentage()),
			ChunkLength:                     quorum.GetChunkLength(),
		}
	}
	return certTypes.EigenDATypesV1BlobHeader{
		Commitment: certTypes.BN254G1Point{
			X: new(big.Int).SetBytes(blobHeader.GetCommitment().GetX()),
			Y: new(big.Int).SetBytes(blobHeader.GetCommitment().GetY()),
		},
		DataLength:       blobHeader.GetDataLength(),
		QuorumBlobParams: quorums,
	}
}

func blobVerificationProofFromProto(verificationProof *disperserpb.BlobVerificationProof) certTypes.EigenDATypesV1BlobVerificationProof {
	batchMetadataProto := verificationProof.GetBatchMetadata()
	batchHeaderProto := verificationProof.GetBatchMetadata().GetBatchHeader()
	var batchRoot [32]byte
	copy(batchRoot[:], batchHeaderProto.GetBatchRoot())
	batchHeader := certTypes.EigenDATypesV1BatchHeader{
		BlobHeadersRoot:       batchRoot,
		QuorumNumbers:         batchHeaderProto.GetQuorumNumbers(),
		SignedStakeForQuorums: batchHeaderProto.GetQuorumSignedPercentages(),
		ReferenceBlockNumber:  batchHeaderProto.GetReferenceBlockNumber(),
	}
	var sig [32]byte
	copy(sig[:], batchMetadataProto.GetSignatoryRecordHash())
	batchMetadata := certTypes.EigenDATypesV1BatchMetadata{
		BatchHeader:             batchHeader,
		SignatoryRecordHash:     sig,
		ConfirmationBlockNumber: batchMetadataProto.GetConfirmationBlockNumber(),
	}
	return certTypes.EigenDATypesV1BlobVerificationProof{
		BatchId:        verificationProof.GetBatchId(),
		BlobIndex:      verificationProof.GetBlobIndex(),
		BatchMetadata:  batchMetadata,
		InclusionProof: verificationProof.GetInclusionProof(),
		QuorumIndices:  verificationProof.GetQuorumIndexes(),
	}
}
