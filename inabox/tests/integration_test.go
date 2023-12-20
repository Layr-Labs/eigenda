package integration_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"math/big"
	"time"

	disperserpb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	rollupbindings "github.com/Layr-Labs/eigenda/contracts/bindings/MockRollup"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/tools/traffic"

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

		optsWithValue, err := ethClient.GetNoSendTransactOpts()
		Expect(err).To(BeNil())
		optsWithValue.Value = big.NewInt(1e18)
		tx, err := mockRollup.RegisterValidator(optsWithValue)
		Expect(err).To(BeNil())
		tx, err = ethClient.UpdateGas(ctx, tx, optsWithValue.Value)
		Expect(err).To(BeNil())
		err = ethClient.SendTransaction(ctx, tx)
		Expect(err).To(BeNil())
		mineAnvilBlocks(numConfirmations + 1)
		_, err = ethClient.EnsureTransactionEvaled(ctx, tx, "RegisterValidator")
		Expect(err).To(BeNil())

		disp := traffic.NewDisperserClient(&traffic.Config{
			Hostname:        "localhost",
			GrpcPort:        "32003",
			NumInstances:    1,
			DataSize:        1000_000,
			RequestInterval: 1 * time.Second,
			Timeout:         10 * time.Second,
		})
		Expect(disp).To(Not(BeNil()))

		data := make([]byte, 1024)
		_, err = rand.Read(data)
		Expect(err).To(BeNil())

		blobStatus1, key1, err := disp.DisperseBlob(ctx, data, 0, 100, 80)
		Expect(err).To(BeNil())
		Expect(key1).To(Not(BeNil()))
		Expect(blobStatus1).To(Not(BeNil()))
		Expect(*blobStatus1).To(Equal(disperser.Processing))

		blobStatus2, key2, err := disp.DisperseBlob(ctx, data, 0, 100, 80)
		Expect(err).To(BeNil())
		Expect(key2).To(Not(BeNil()))
		Expect(blobStatus2).To(Not(BeNil()))
		Expect(*blobStatus2).To(Equal(disperser.Processing))

		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()

		var blobStatus *disperser.BlobStatus
		var reply *disperserpb.BlobStatusReply
	loop:
		for {
			select {
			case <-ctx.Done():
				Fail("timed out")
			case <-ticker.C:
				reply, err = disp.GetBlobStatus(context.Background(), key1)
				Expect(err).To(BeNil())
				Expect(reply).To(Not(BeNil()))
				blobStatus, err = disperser.FromBlobStatusProto(reply.GetStatus())
				Expect(err).To(BeNil())
				if *blobStatus == disperser.Confirmed {
					blobHeader := blobHeaderFromProto(reply.GetInfo().GetBlobHeader())
					verificationProof := blobVerificationProofFromProto(reply.GetInfo().GetBlobVerificationProof())
					opts, err := ethClient.GetNoSendTransactOpts()
					Expect(err).To(BeNil())
					tx, err := mockRollup.PostCommitment(opts, blobHeader, verificationProof)
					Expect(err).To(BeNil())
					tx, err = ethClient.UpdateGas(ctx, tx, nil)
					Expect(err).To(BeNil())
					err = ethClient.SendTransaction(ctx, tx)
					Expect(err).To(BeNil())
					mineAnvilBlocks(numConfirmations + 1)
					_, err = ethClient.EnsureTransactionEvaled(ctx, tx, "PostCommitment")
					Expect(err).To(BeNil())
					break loop
				} else {
					mineAnvilBlocks(numConfirmations + 1)
				}
			}
		}
		Expect(*blobStatus).To(Equal(disperser.Confirmed))

		ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		retrieved, err := retrievalClient.RetrieveBlob(ctx,
			[32]byte(reply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
			reply.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
			uint(reply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber()),
			[32]byte(reply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
			0,
		)
		Expect(err).To(BeNil())
		Expect(bytes.TrimRight(retrieved, "\x00")).To(Equal(bytes.TrimRight(data, "\x00")))
	})
})

func blobHeaderFromProto(blobHeader *disperserpb.BlobHeader) rollupbindings.IEigenDAServiceManagerBlobHeader {
	commitmentBytes := blobHeader.GetCommitment()
	commitment, err := new(core.Commitment).Deserialize(commitmentBytes)
	Expect(err).To(BeNil())
	quorums := make([]rollupbindings.IEigenDAServiceManagerQuorumBlobParam, len(blobHeader.GetBlobQuorumParams()))
	for i, quorum := range blobHeader.GetBlobQuorumParams() {
		quorums[i] = rollupbindings.IEigenDAServiceManagerQuorumBlobParam{
			QuorumNumber:                 uint8(quorum.GetQuorumNumber()),
			AdversaryThresholdPercentage: uint8(quorum.GetAdversaryThresholdPercentage()),
			QuorumThresholdPercentage:    uint8(quorum.GetQuorumThresholdPercentage()),
			ChunkLength:                  quorum.ChunkLength,
		}
	}

	return rollupbindings.IEigenDAServiceManagerBlobHeader{
		Commitment: rollupbindings.BN254G1Point{
			X: commitment.X.BigInt(new(big.Int)),
			Y: commitment.Y.BigInt(new(big.Int)),
		},
		DataLength:       blobHeader.GetDataLength(),
		QuorumBlobParams: quorums,
	}
}

func blobVerificationProofFromProto(verificationProof *disperserpb.BlobVerificationProof) rollupbindings.EigenDABlobUtilsBlobVerificationProof {
	batchMetadataProto := verificationProof.GetBatchMetadata()
	batchHeaderProto := verificationProof.GetBatchMetadata().GetBatchHeader()
	var batchRoot [32]byte
	copy(batchRoot[:], batchHeaderProto.GetBatchRoot())
	batchHeader := rollupbindings.IEigenDAServiceManagerBatchHeader{
		BlobHeadersRoot:            batchRoot,
		QuorumNumbers:              batchHeaderProto.GetQuorumNumbers(),
		QuorumThresholdPercentages: batchHeaderProto.GetQuorumSignedPercentages(),
		ReferenceBlockNumber:       batchHeaderProto.GetReferenceBlockNumber(),
	}
	var sig [32]byte
	copy(sig[:], batchMetadataProto.GetSignatoryRecordHash())
	fee := new(big.Int).SetBytes(batchMetadataProto.GetFee())
	batchMetadata := rollupbindings.IEigenDAServiceManagerBatchMetadata{
		BatchHeader:             batchHeader,
		SignatoryRecordHash:     sig,
		Fee:                     fee,
		ConfirmationBlockNumber: batchMetadataProto.GetConfirmationBlockNumber(),
	}
	return rollupbindings.EigenDABlobUtilsBlobVerificationProof{
		BatchId:                verificationProof.GetBatchId(),
		BlobIndex:              uint8(verificationProof.GetBlobIndex()),
		BatchMetadata:          batchMetadata,
		InclusionProof:         verificationProof.GetInclusionProof(),
		QuorumThresholdIndexes: verificationProof.GetQuorumIndexes(),
	}
}
