package integration_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Inabox v2 Integration - Payment", func() {
	It("test reserved payment only scenario", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
		defer cancel()

		privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdee"
		signer := auth.NewLocalBlobRequestSigner(privateKeyHex)

		accountID, err := signer.GetAccountID()
		fmt.Println("accountID", accountID)
		Expect(err).To(BeNil())

		disp, err := clients.NewDisperserClient(&clients.DisperserClientConfig{
			Hostname: "localhost",
			Port:     "32005",
		}, signer, nil, nil)
		Expect(err).To(BeNil())
		Expect(disp).To(Not(BeNil()))
		err = disp.PopulateAccountant(ctx)
		Expect(err).To(BeNil())

		data1 := make([]byte, 992)
		_, err = rand.Read(data1)
		Expect(err).To(BeNil())
		data2 := make([]byte, 123)
		_, err = rand.Read(data2)
		Expect(err).To(BeNil())

		paddedData1 := codec.ConvertByPaddingEmptyByte(data1)

		blobStatus1, key1, err := disp.DisperseBlob(ctx, paddedData1, 0, []uint8{0, 1}, 0)
		Expect(err).To(Not(BeNil()))
		Expect(key1).To(BeNil())
		Expect(blobStatus1).To(BeNil())
	})
	It("test ondemand payment only scenario", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
		defer cancel()

		privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdee"
		signer := auth.NewLocalBlobRequestSigner(privateKeyHex)

		disp, err := clients.NewDisperserClient(&clients.DisperserClientConfig{
			Hostname: "localhost",
			Port:     "32005",
		}, signer, nil, nil)
		Expect(err).To(BeNil())
		Expect(disp).To(Not(BeNil()))
		err = disp.PopulateAccountant(ctx)
		Expect(err).To(BeNil())

		data1 := make([]byte, 992)
		_, err = rand.Read(data1)
		Expect(err).To(BeNil())
		data2 := make([]byte, 123)
		_, err = rand.Read(data2)
		Expect(err).To(BeNil())

		paddedData1 := codec.ConvertByPaddingEmptyByte(data1)

		blobStatus1, key1, err := disp.DisperseBlob(ctx, paddedData1, 0, []uint8{0, 1}, 0)
		Expect(err).To(Not(BeNil()))
		Expect(key1).To(BeNil())
		Expect(blobStatus1).To(BeNil())
	})
	It("test failing payment scenario", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
		defer cancel()

		privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdee"
		signer := auth.NewLocalBlobRequestSigner(privateKeyHex)

		disp, err := clients.NewDisperserClient(&clients.DisperserClientConfig{
			Hostname: "localhost",
			Port:     "32005",
		}, signer, nil, nil)
		Expect(err).To(BeNil())
		Expect(disp).To(Not(BeNil()))
		err = disp.PopulateAccountant(ctx)
		Expect(err).To(BeNil())

		data1 := make([]byte, 992)
		_, err = rand.Read(data1)
		Expect(err).To(BeNil())
		data2 := make([]byte, 123)
		_, err = rand.Read(data2)
		Expect(err).To(BeNil())

		paddedData1 := codec.ConvertByPaddingEmptyByte(data1)

		blobStatus1, key1, err := disp.DisperseBlob(ctx, paddedData1, 0, []uint8{0, 1}, 0)
		Expect(err).To(Not(BeNil()))
		Expect(key1).To(BeNil())
		Expect(blobStatus1).To(BeNil())
	})
})
