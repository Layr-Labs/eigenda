package integration_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	disperserpb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	rollupbindings "github.com/Layr-Labs/eigenda/contracts/bindings/MockRollup"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/disperser"

	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func mineAnvilBlocks(numBlocks int) {
	for i := 0; i < numBlocks; i++ {
		err := rpcClient.CallContext(context.Background(), nil, "evm_mine")
		Expect(err).To(BeNil())
	}
}

var _ = Describe("Inabox Integration", func() {
	It("test end to end scenario", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		gasTipCap, gasFeeCap, err := ethClient.GetLatestGasCaps(ctx)
		Expect(err).To(BeNil())

		privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded"
		signer := auth.NewLocalBlobRequestSigner(privateKeyHex)

		disp, err := clients.NewDisperserClient(&clients.Config{
			Hostname: "localhost",
			Port:     "32003",
			Timeout:  10 * time.Second,
		}, signer)
		Expect(err).To(BeNil())
		Expect(disp).To(Not(BeNil()))

		data := make([]byte, 1024)
		_, err = rand.Read(data)
		Expect(err).To(BeNil())

		paddedData := codec.ConvertByPaddingEmptyByte(data)

		blobStatus1, key1, err := disp.DisperseBlob(ctx, paddedData, []uint8{})
		Expect(err).To(BeNil())
		Expect(key1).To(Not(BeNil()))
		Expect(blobStatus1).To(Not(BeNil()))
		Expect(*blobStatus1).To(Equal(disperser.Processing))

		blobStatus2, key2, err := disp.DisperseBlobAuthenticated(ctx, paddedData, []uint8{})
		Expect(err).To(BeNil())
		Expect(key2).To(Not(BeNil()))
		Expect(blobStatus2).To(Not(BeNil()))
		Expect(*blobStatus2).To(Equal(disperser.Processing))

		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()

		var reply1 *disperserpb.BlobStatusReply
		var reply2 *disperserpb.BlobStatusReply

		for loop := true; loop; {
			select {
			case <-ctx.Done():
				Fail("timed out")
			case <-ticker.C:
				reply1, err = disp.GetBlobStatus(context.Background(), key1)
				Expect(err).To(BeNil())
				Expect(reply1).To(Not(BeNil()))
				blobStatus1, err = disperser.FromBlobStatusProto(reply1.GetStatus())
				Expect(err).To(BeNil())

				reply2, err = disp.GetBlobStatus(context.Background(), key2)
				Expect(err).To(BeNil())
				Expect(reply2).To(Not(BeNil()))
				blobStatus2, err = disperser.FromBlobStatusProto(reply2.GetStatus())
				Expect(err).To(BeNil())

				if *blobStatus1 != disperser.Confirmed || *blobStatus2 != disperser.Confirmed {
					mineAnvilBlocks(numConfirmations + 1)
					continue
				}
				blobHeader := blobHeaderFromProto(reply1.GetInfo().GetBlobHeader())
				verificationProof := blobVerificationProofFromProto(reply1.GetInfo().GetBlobVerificationProof())
				opts, err := ethClient.GetNoSendTransactOpts()
				Expect(err).To(BeNil())
				tx, err := mockRollup.PostCommitment(opts, blobHeader, verificationProof)
				Expect(err).To(BeNil())
				tx, err = ethClient.UpdateGas(ctx, tx, nil, gasTipCap, gasFeeCap)
				Expect(err).To(BeNil())
				err = ethClient.SendTransaction(ctx, tx)
				Expect(err).To(BeNil())
				mineAnvilBlocks(numConfirmations + 1)
				_, err = ethClient.EnsureTransactionEvaled(ctx, tx, "PostCommitment")
				Expect(err).To(BeNil())

				blobHeader = blobHeaderFromProto(reply2.GetInfo().GetBlobHeader())
				verificationProof = blobVerificationProofFromProto(reply2.GetInfo().GetBlobVerificationProof())
				tx, err = mockRollup.PostCommitment(opts, blobHeader, verificationProof)
				Expect(err).To(BeNil())
				tx, err = ethClient.UpdateGas(ctx, tx, nil, gasTipCap, gasFeeCap)
				Expect(err).To(BeNil())
				err = ethClient.SendTransaction(ctx, tx)
				Expect(err).To(BeNil())
				mineAnvilBlocks(numConfirmations + 1)
				_, err = ethClient.EnsureTransactionEvaled(ctx, tx, "PostCommitment")
				Expect(err).To(BeNil())
				loop = false
			}
		}
		Expect(*blobStatus1).To(Equal(disperser.Confirmed))
		Expect(*blobStatus2).To(Equal(disperser.Confirmed))

		ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		retrieved, err := retrievalClient.RetrieveBlob(ctx,
			[32]byte(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
			reply1.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
			uint(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber()),
			[32]byte(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
			0, // retrieve blob 1 from quorum 0
		)
		Expect(err).To(BeNil())

		restored := codec.RemoveEmptyByteFromPaddedBytes(retrieved)
		Expect(bytes.TrimRight(restored, "\x00")).To(Equal(bytes.TrimRight(data, "\x00")))

		_, err = retrievalClient.RetrieveBlob(ctx,
			[32]byte(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
			reply1.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
			uint(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber()),
			[32]byte(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
			1, // retrieve blob 1 from quorum 1
		)
		Expect(err).To(BeNil())

		_, err = retrievalClient.RetrieveBlob(ctx,
			[32]byte(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
			reply1.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
			uint(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber()),
			[32]byte(reply1.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
			2, // retrieve blob 1 from quorum 2
		)
		Expect(err).NotTo(BeNil())

		retrieved, err = retrievalClient.RetrieveBlob(ctx,
			[32]byte(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
			reply2.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
			uint(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber()),
			[32]byte(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
			0, // retrieve from quorum 0
		)
		Expect(err).To(BeNil())
		restored = codec.RemoveEmptyByteFromPaddedBytes(retrieved)
		Expect(bytes.TrimRight(restored, "\x00")).To(Equal(bytes.TrimRight(data, "\x00")))
		_, err = retrievalClient.RetrieveBlob(ctx,
			[32]byte(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
			reply2.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
			uint(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber()),
			[32]byte(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
			1, // retrieve from quorum 1
		)
		Expect(err).To(BeNil())
		_, err = retrievalClient.RetrieveBlob(ctx,
			[32]byte(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
			reply2.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
			uint(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber()),
			[32]byte(reply2.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
			2, // retrieve from quorum 2
		)
		Expect(err).NotTo(BeNil())
	})
})

func blobHeaderFromProto(blobHeader *disperserpb.BlobHeader) rollupbindings.BlobHeader {
	quorums := make([]rollupbindings.QuorumBlobParam, len(blobHeader.GetBlobQuorumParams()))
	for i, quorum := range blobHeader.GetBlobQuorumParams() {
		quorums[i] = rollupbindings.QuorumBlobParam{
			QuorumNumber:                    uint8(quorum.GetQuorumNumber()),
			AdversaryThresholdPercentage:    uint8(quorum.GetAdversaryThresholdPercentage()),
			ConfirmationThresholdPercentage: uint8(quorum.GetConfirmationThresholdPercentage()),
			ChunkLength:                     quorum.ChunkLength,
		}
	}
	return rollupbindings.BlobHeader{
		Commitment: rollupbindings.BN254G1Point{
			X: new(big.Int).SetBytes(blobHeader.GetCommitment().X),
			Y: new(big.Int).SetBytes(blobHeader.GetCommitment().Y),
		},
		DataLength:       blobHeader.GetDataLength(),
		QuorumBlobParams: quorums,
	}
}

func blobVerificationProofFromProto(verificationProof *disperserpb.BlobVerificationProof) rollupbindings.BlobVerificationProof {
	batchMetadataProto := verificationProof.GetBatchMetadata()
	batchHeaderProto := verificationProof.GetBatchMetadata().GetBatchHeader()
	var batchRoot [32]byte
	copy(batchRoot[:], batchHeaderProto.GetBatchRoot())
	batchHeader := rollupbindings.BatchHeader{
		BlobHeadersRoot:       batchRoot,
		QuorumNumbers:         batchHeaderProto.GetQuorumNumbers(),
		SignedStakeForQuorums: batchHeaderProto.GetQuorumSignedPercentages(),
		ReferenceBlockNumber:  batchHeaderProto.GetReferenceBlockNumber(),
	}
	var sig [32]byte
	copy(sig[:], batchMetadataProto.GetSignatoryRecordHash())
	batchMetadata := rollupbindings.BatchMetadata{
		BatchHeader:             batchHeader,
		SignatoryRecordHash:     sig,
		ConfirmationBlockNumber: batchMetadataProto.GetConfirmationBlockNumber(),
	}
	return rollupbindings.BlobVerificationProof{
		BatchId:        verificationProof.GetBatchId(),
		BlobIndex:      verificationProof.GetBlobIndex(),
		BatchMetadata:  batchMetadata,
		InclusionProof: verificationProof.GetInclusionProof(),
		QuorumIndices:  verificationProof.GetQuorumIndexes(),
	}
}
