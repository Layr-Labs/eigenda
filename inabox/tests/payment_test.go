package integration_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	disperserpb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/meterer"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Inabox Integration", func() {
	It("test payment metering", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		gasTipCap, gasFeeCap, err := ethClient.GetLatestGasCaps(ctx)
		Expect(err).To(BeNil())

		privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded"
		signer := auth.NewLocalBlobRequestSigner(privateKeyHex)
		// Disperser configs: rsv window 60s, min chargeable size 100 bytes, price per chargeable 100, global limit 500
		// -> need to check the mock, can't just use any account for the disperser client, consider using static wallets...

		// say with dataLength of 150 bytes, within a window, we can send 7 blobs with overflow of 50 bytes
		// the later requests is then 250 bytes, try send 4 blobs within a second, 2 of them would fail but not charged for
		// wait for a second, retry, and that should allow ondemand to work
		privateKey, err := crypto.HexToECDSA(privateKeyHex[2:]) // Remove "0x" prefix
		Expect(err).To(BeNil())
		disp := clients.NewDisperserClient(&clients.Config{
			Hostname: "localhost",
			Port:     "32003",
			Timeout:  10 * time.Second,
		}, signer, clients.NewAccountant(meterer.DummyReservation, meterer.DummyOnDemandPayment, 60, meterer.DummyMinimumChargeableSize, meterer.DummyMinimumChargeablePayment, privateKey))

		Expect(disp).To(Not(BeNil()))

		singleBlobSize := meterer.DummyMinimumChargeableSize
		data := make([]byte, singleBlobSize)
		_, err = rand.Read(data)
		Expect(err).To(BeNil())

		paddedData := codec.ConvertByPaddingEmptyByte(data)

		// requests that count towards either reservation or payments
		paidBlobStatus := []disperser.BlobStatus{}
		paidKeys := [][]byte{}
		// TODO: payment calculation unit consistency
		for i := 0; i < (int(meterer.DummyReservationBytesLimit)+int(meterer.DummyPaymentLimit))/int(singleBlobSize); i++ {
			blobStatus, key, err := disp.PaidDisperseBlob(ctx, paddedData, []uint8{0})
			Expect(err).To(BeNil())
			Expect(key).To(Not(BeNil()))
			Expect(blobStatus).To(Not(BeNil()))
			Expect(*blobStatus).To(Equal(disperser.Processing))
			paidBlobStatus = append(paidBlobStatus, *blobStatus)
			paidKeys = append(paidKeys, key)
		}

		// requests that aren't covered by reservation or on-demand payment
		blobStatus, key, err := disp.PaidDisperseBlob(ctx, paddedData, []uint8{0})
		Expect(err).To(Not(BeNil()))
		Expect(key).To(BeNil())
		Expect(blobStatus).To(BeNil())

		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()

		var replies = make([]*disperserpb.BlobStatusReply, len(paidBlobStatus))
		// now make sure all the paid blobs get confirmed
	loop:
		for {
			select {
			case <-ctx.Done():
				Fail("timed out")
			case <-ticker.C:
				notConfirmed := false
				for i, key := range paidKeys {
					reply, err := disp.GetBlobStatus(context.Background(), key)
					Expect(err).To(BeNil())
					Expect(reply).To(Not(BeNil()))
					status, err := disperser.FromBlobStatusProto(reply.GetStatus())
					Expect(err).To(BeNil())
					if *status != disperser.Confirmed {
						notConfirmed = true
					}
					replies[i] = reply
					paidBlobStatus[i] = *status
				}

				if notConfirmed {
					mineAnvilBlocks(numConfirmations + 1)
					continue
				}

				for _, reply := range replies {
					blobHeader := blobHeaderFromProto(reply.GetInfo().GetBlobHeader())
					verificationProof := blobVerificationProofFromProto(reply.GetInfo().GetBlobVerificationProof())
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
				}

				break loop
			}
		}
		for _, status := range paidBlobStatus {
			Expect(status).To(Equal(disperser.Confirmed))
		}

		ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		for _, reply := range replies {
			retrieved, err := retrievalClient.RetrieveBlob(ctx,
				[32]byte(reply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
				reply.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
				uint(reply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber()),
				[32]byte(reply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
				0, // retrieve blob 1 from quorum 0
			)
			Expect(err).To(BeNil())
			restored := codec.RemoveEmptyByteFromPaddedBytes(retrieved)
			Expect(bytes.TrimRight(restored, "\x00")).To(Equal(bytes.TrimRight(data, "\x00")))

			_, err = retrievalClient.RetrieveBlob(ctx,
				[32]byte(reply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
				reply.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
				uint(reply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber()),
				[32]byte(reply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
				1, // retrieve blob 1 from quorum 1
			)
			Expect(err).NotTo(BeNil())
		}
	})
})
