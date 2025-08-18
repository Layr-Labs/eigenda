package mock

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/common"
	v2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"google.golang.org/grpc"
)

// DisperserRPC is a mock implementation of disperser_rpc.DisperserClient
type DisperserRPC struct {
	DisperseCount     int
	DisperseMutex     sync.Mutex
	DisperseCallTimes []time.Time
	DisperseDelay     time.Duration
}

// NewDisperserRPC creates a new mock DisperserRPC with default values
func NewDisperserRPC() *DisperserRPC {
	return &DisperserRPC{
		DisperseCount:     0,
		DisperseCallTimes: []time.Time{},
		DisperseDelay:     0,
	}
}

// DisperseBlob is a mock implementation that simulates a delay in processing
func (m *DisperserRPC) DisperseBlob(ctx context.Context, in *v2.DisperseBlobRequest, opts ...grpc.CallOption) (*v2.DisperseBlobReply, error) {
	m.DisperseMutex.Lock()
	callTime := time.Now()
	m.DisperseCallTimes = append(m.DisperseCallTimes, callTime)
	m.DisperseCount++
	m.DisperseMutex.Unlock()

	// Simulate processing time
	time.Sleep(m.DisperseDelay)

	blobKey := [32]byte{1, 2, 3}
	return &v2.DisperseBlobReply{
		BlobKey: blobKey[:],
		Result:  v2.BlobStatus_QUEUED,
	}, nil
}

// GetBlobStatus is a mock implementation
func (m *DisperserRPC) GetBlobStatus(ctx context.Context, in *v2.BlobStatusRequest, opts ...grpc.CallOption) (*v2.BlobStatusReply, error) {
	return &v2.BlobStatusReply{}, nil
}

// GetBlobCommitment is a mock implementation
func (m *DisperserRPC) GetBlobCommitment(ctx context.Context, in *v2.BlobCommitmentRequest, opts ...grpc.CallOption) (*v2.BlobCommitmentReply, error) {
	return &v2.BlobCommitmentReply{
		BlobCommitment: &common.BlobCommitment{
			Length: 32,
		},
	}, nil
}

// GetPaymentState is a mock implementation
func (m *DisperserRPC) GetPaymentState(ctx context.Context, in *v2.GetPaymentStateRequest, opts ...grpc.CallOption) (*v2.GetPaymentStateReply, error) {
	// Create a mock payment state response with valid global parameters
	return &v2.GetPaymentStateReply{
		PaymentGlobalParams: &v2.PaymentGlobalParams{
			MinNumSymbols:          32,   // Ensure non-zero value to avoid division by zero
			PricePerSymbol:         100,  // Mock price
			ReservationWindow:      3600, // Mock window
			GlobalSymbolsPerSecond: 1000, // Mock rate limit
		},
		Reservation: &v2.Reservation{
			SymbolsPerSecond: 100,
			StartTimestamp:   uint32(time.Now().Unix() - 3600), // Start an hour ago
			EndTimestamp:     uint32(time.Now().Unix() + 3600), // End an hour from now
			QuorumNumbers:    []uint32{1},                      // Allow quorum 1
		},
		CumulativePayment:        big.NewInt(0).Bytes(),
		OnchainCumulativePayment: big.NewInt(1000000).Bytes(), // Allow some payment
	}, nil
}
