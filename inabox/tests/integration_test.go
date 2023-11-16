package integration_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"math/big"
	"time"

	disperserpb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	retrieverpb "github.com/Layr-Labs/eigenda/api/grpc/retriever"
	rollupbindings "github.com/Layr-Labs/eigenda/contracts/bindings/MockRollup"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/tools/traffic"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Inabox Integration", func() {
	It("test end to end scenario", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		optsWithValue := new(bind.TransactOpts)
		*optsWithValue = *ethClient.GetNoSendTransactOpts()
		optsWithValue.Value = big.NewInt(1e18)
		tx, err := mockRollup.RegisterValidator(optsWithValue)
		Expect(err).To(BeNil())
		_, err = ethClient.EstimateGasPriceAndLimitAndSendTx(ctx, tx, "RegisterValidator", big.NewInt(1e18))
		Expect(err).To(BeNil())

		disp := traffic.NewDisperserClient(&traffic.Config{
			Hostname:        "localhost",
			GrpcPort:        "32001",
			NumInstances:    1,
			DataSize:        1000_000,
			RequestInterval: 1 * time.Second,
			Timeout:         10 * time.Second,
		})
		Expect(disp).To(Not(BeNil()))

		// Must not end in 0x00 or else this test will fail in a flakey manner because padding bytes are greedily trimmed.
		data := make([]byte, 1023)
		_, err = rand.Read(data)
		data = append(data, 0x01)
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
					tx, err := mockRollup.PostCommitment(ethClient.GetNoSendTransactOpts(), blobHeader, verificationProof)
					Expect(err).To(BeNil())
					_, err = ethClient.EstimateGasPriceAndLimitAndSendTx(ctx, tx, "PostCommitment", nil)
					Expect(err).To(BeNil())
					break loop
				}
			}
		}
		Expect(*blobStatus).To(Equal(disperser.Confirmed))

		ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		conn, err := grpc.Dial(
			"localhost:32014",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).To(BeNil())
		defer conn.Close()

		retrieverClient := retrieverpb.NewRetrieverClient(conn)

		retrieveReply, err := retrieverClient.RetrieveBlob(ctx,
			&retrieverpb.BlobRequest{
				BatchHeaderHash:      []byte(reply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
				BlobIndex:            reply.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
				ReferenceBlockNumber: reply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber(),
				QuorumId:             0,
			},
		)
		Expect(err).To(BeNil())
		retrieved := retrieveReply.Data
		Expect(bytes.TrimRight(retrieved, "\x00")).To(Equal(data))
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
			QuantizationParameter:        uint8(quorum.GetQuantizationParam()),
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
