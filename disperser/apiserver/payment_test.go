package apiserver_test

import (
	"context"
	"crypto/rand"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/assert"
	tmock "github.com/stretchr/testify/mock"
	"google.golang.org/grpc/peer"
)

func TestDispersePaidBlob(t *testing.T) {

	transactor := &mock.MockWriter{}
	transactor.On("GetCurrentBlockNumber").Return(uint32(100), nil)
	transactor.On("GetQuorumCount").Return(uint8(2), nil)
	quorumParams := []core.SecurityParam{
		{QuorumID: 0, AdversaryThreshold: 80, ConfirmationThreshold: 100},
		{QuorumID: 1, AdversaryThreshold: 80, ConfirmationThreshold: 100},
	}
	transactor.On("GetQuorumSecurityParams", tmock.Anything).Return(quorumParams, nil)
	transactor.On("GetRequiredQuorumNumbers", tmock.Anything).Return([]uint8{0, 1}, nil)

	quorums := []uint32{0, 1}

	dispersalServer := newTestServer(transactor, t.Name())

	data := make([]byte, 1024)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	data = codec.ConvertByPaddingEmptyByte(data)

	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 51001,
		},
	}
	ctx := peer.NewContext(context.Background(), p)

	transactor.On("GetRequiredQuorumNumbers", tmock.Anything).Return([]uint8{0, 1}, nil).Twice()

	pk := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdeb"
	signer := auth.NewPaymentSigner(pk)

	symbolLength := encoding.GetBlobLength(uint(len(data)))
	// disperse on-demand payment
	for i := 1; i < 3; i++ {
		pm := &core.PaymentMetadata{
			AccountID:         signer.GetAccountID(),
			BinIndex:          0,
			CumulativePayment: big.NewInt(int64(int(symbolLength) * i * encoding.BYTES_PER_SYMBOL)),
		}
		sig, err := signer.SignBlobPayment(pm)
		assert.NoError(t, err)
		reply, err := dispersalServer.DispersePaidBlob(ctx, &pb.DispersePaidBlobRequest{
			Data:             data,
			QuorumNumbers:    quorums,
			PaymentHeader:    pm.ConvertToProtoPaymentHeader(),
			PaymentSignature: sig,
		})
		assert.NoError(t, err)
		assert.Equal(t, reply.GetResult(), pb.BlobStatus_PROCESSING)
		assert.NotNil(t, reply.GetRequestId())
	}

	// exceeded payment limit
	pm := &core.PaymentMetadata{
		AccountID:         signer.GetAccountID(),
		BinIndex:          0,
		CumulativePayment: big.NewInt(int64(symbolLength*3)*encoding.BYTES_PER_SYMBOL - 1),
	}
	sig, err := signer.SignBlobPayment(pm)
	assert.NoError(t, err)
	_, err = dispersalServer.DispersePaidBlob(ctx, &pb.DispersePaidBlobRequest{
		Data:             data,
		QuorumNumbers:    quorums,
		PaymentHeader:    pm.ConvertToProtoPaymentHeader(),
		PaymentSignature: sig,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request claims a cumulative payment greater than the on-chain deposit")

	// disperse paid reservation (any quorum number)
	// TODO: somehow meterer is not defined as a method or field in dispersalServer; reservationWindow we set was 1
	for i := 0; i < 2; i++ {
		binIndex := meterer.GetBinIndex(uint64(time.Now().Unix()), 1)
		pm := &core.PaymentMetadata{
			AccountID:         signer.GetAccountID(),
			BinIndex:          binIndex,
			CumulativePayment: big.NewInt(0),
		}
		sig, err = signer.SignBlobPayment(pm)
		assert.NoError(t, err)
		reply, err := dispersalServer.DispersePaidBlob(ctx, &pb.DispersePaidBlobRequest{
			Data:             data,
			QuorumNumbers:    []uint32{1},
			PaymentHeader:    pm.ConvertToProtoPaymentHeader(),
			PaymentSignature: sig,
		})
		assert.NoError(t, err)
		assert.Equal(t, reply.GetResult(), pb.BlobStatus_PROCESSING)
		assert.NotNil(t, reply.GetRequestId())

	}
	binIndex := meterer.GetBinIndex(uint64(time.Now().Unix()), 1)
	pm = &core.PaymentMetadata{
		AccountID:         signer.GetAccountID(),
		BinIndex:          binIndex,
		CumulativePayment: big.NewInt(0),
	}
	sig, err = signer.SignBlobPayment(pm)
	assert.NoError(t, err)
	_, err = dispersalServer.DispersePaidBlob(ctx, &pb.DispersePaidBlobRequest{
		Data:             data,
		QuorumNumbers:    []uint32{1},
		PaymentHeader:    pm.ConvertToProtoPaymentHeader(),
		PaymentSignature: sig,
	})
	assert.Contains(t, err.Error(), "bin has already been filled")

	// invalid bin index
	binIndex = meterer.GetBinIndex(uint64(time.Now().Unix())/2, 1)
	pm = &core.PaymentMetadata{
		AccountID:         signer.GetAccountID(),
		BinIndex:          binIndex,
		CumulativePayment: big.NewInt(0),
	}
	sig, err = signer.SignBlobPayment(pm)
	_, err = dispersalServer.DispersePaidBlob(ctx, &pb.DispersePaidBlobRequest{
		Data:             data,
		QuorumNumbers:    []uint32{1},
		PaymentHeader:    pm.ConvertToProtoPaymentHeader(),
		PaymentSignature: sig,
	})
	assert.Contains(t, err.Error(), "invalid bin index for reservation")
}
