package integration_test

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Inabox v2 Integration", func() {
	It("reservation runs out", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdee"
		signer := auth.NewLocalBlobRequestSigner(privateKeyHex)

		disp, err := clients.NewDisperserClient(&clients.DisperserClientConfig{
			Hostname: "localhost",
			Port:     "32005",
		}, signer, nil, nil)
		Expect(err).To(BeNil())
		Expect(disp).To(Not(BeNil()))

		data1 := make([]byte, 992)
		_, err = rand.Read(data1)
		Expect(err).To(BeNil())
		data2 := make([]byte, 123)
		_, err = rand.Read(data2)
		Expect(err).To(BeNil())

		paddedData1 := codec.ConvertByPaddingEmptyByte(data1)
		paddedData2 := codec.ConvertByPaddingEmptyByte(data2)

		blobStatus1, key1, err := disp.DisperseBlob(ctx, paddedData1, 0, []uint8{0, 1}, 0)
		Expect(err).To(BeNil())
		Expect(key1).To(Not(BeNil()))
		Expect(blobStatus1).To(Not(BeNil()))
		Expect(*blobStatus1).To(Equal(dispv2.Queued))

		blobStatus2, key2, err := disp.DisperseBlob(ctx, paddedData2, 0, []uint8{0}, 0)
		Expect(err).To(BeNil())
		Expect(key2).To(Not(BeNil()))
		Expect(blobStatus2).To(Not(BeNil()))
		Expect(*blobStatus2).To(Equal(dispv2.Queued))

		blobStatus3, key3, err := disp.DisperseBlob(ctx, paddedData2, 0, []uint8{0}, 0)
		Expect(err.Error()).To(ContainSubstring("neither reservation nor on-demand payment is available"))
		Expect(key3).To(Not(BeNil()))
		Expect(blobStatus3).To(BeNil())
	})

	It("ondemand runs out", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
		signer := auth.NewLocalBlobRequestSigner(privateKeyHex)

		disp, err := clients.NewDisperserClient(&clients.DisperserClientConfig{
			Hostname: "localhost",
			Port:     "32005",
		}, signer, nil, nil)
		Expect(err).To(BeNil())
		Expect(disp).To(Not(BeNil()))

		data1 := make([]byte, 992)
		_, err = rand.Read(data1)
		Expect(err).To(BeNil())
		data2 := make([]byte, 123)
		_, err = rand.Read(data2)
		Expect(err).To(BeNil())

		paddedData1 := codec.ConvertByPaddingEmptyByte(data1)
		paddedData2 := codec.ConvertByPaddingEmptyByte(data2)

		blobStatus1, key1, err := disp.DisperseBlob(ctx, paddedData1, 0, []uint8{0, 1}, 0)
		Expect(err).To(BeNil())
		Expect(key1).To(Not(BeNil()))
		Expect(blobStatus1).To(Not(BeNil()))
		Expect(*blobStatus1).To(Equal(dispv2.Queued))

		blobStatus2, key2, err := disp.DisperseBlob(ctx, paddedData2, 0, []uint8{0}, 0)
		Expect(err.Error()).To(ContainSubstring("neither reservation nor on-demand payment is available"))
		Expect(blobStatus2).To(BeNil())
		Expect(key2).To(BeNil())
	})
})
