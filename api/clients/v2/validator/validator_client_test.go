package validator

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	grpcnode "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	testrandom "github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/gammazero/workerpool"
	"github.com/stretchr/testify/require"
)

// TODO
//  - test what happens when a pessimistic download timeout triggers
//  - test what happens when a long download timeout triggers
//  - test what happens when chunks are invalid
//  - pessimistic timeout -> slow download eventually finishes -> validation on a chunk fails
//  - respecting thread pool limits

func TestBasicWorkflow(t *testing.T) {
	rand := testrandom.NewTestRandom()
	start := rand.Time()

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	fakeClock := atomic.Pointer[time.Time]{}
	fakeClock.Store(&start)

	config := DefaultClientConfig()
	config.ControlLoopPeriod = 50 * time.Microsecond
	config.timeSource = func() time.Time {
		return *fakeClock.Load()
	}
	config.DownloadPessimism = rand.Float64Range(1.0, 8.0)
	config.VerificationPessimism = rand.Float64Range(1.0, 8.0)
	connectionPool := workerpool.New(8)
	computePool := workerpool.New(8)

	blobKey := (v2.BlobKey)(rand.Bytes(32))

	validatorCount := rand.Int32Range(50, 100)
	minimumChunksPerValidator := uint32(1)
	maximumChunksPerValidator := rand.Uint32Range(50, 100)
	quorumID := (core.QuorumID)(rand.Uint32Range(0, 10))

	// The assignments for each operator (i.e. how many chunks it is responsible for)
	assignments := make(map[core.OperatorID]v2.Assignment, validatorCount)
	// The total number of chunks
	totalChunkCount := uint32(0)
	// Simulated chunks for each operator
	operatorChunks := make(map[core.OperatorID][][]byte, validatorCount)

	for i := 0; i < int(validatorCount); i++ {
		operatorID := (core.OperatorID)(rand.PrintableBytes(32))
		numChunks := rand.Uint32Range(minimumChunksPerValidator, maximumChunksPerValidator)
		totalChunkCount += numChunks
		assignments[operatorID] = v2.Assignment{
			NumChunks: numChunks,
		}
		operatorChunks[operatorID] = make([][]byte, numChunks)
		for j := 0; j < int(numChunks); j++ {
			operatorChunks[operatorID][j] = rand.PrintableBytes(8)
		}
	}

	// The number of chunks needed to reconstruct the blob
	minimumChunkCount := uint32(rand.Float64Range(1.0/8.0, 1.0/8.0))

	// the number of chunks downloaded
	chunksDownloaded := atomic.Uint32{}
	// a set of operators that have provided chunks
	downloadSet := make(map[core.OperatorID]struct{}, validatorCount)
	config.UnsafeDownloadChunksFunction = func(
		key v2.BlobKey,
		operatorID core.OperatorID,
	) (*grpcnode.GetChunksReply, error) {

		// verify we have the expected blob key
		require.Equal(t, blobKey, key)

		// make sure this is for a valid operator ID
		chunks, ok := operatorChunks[operatorID]
		require.True(t, ok)

		// only permit downloads to happen once per operator
		_, ok = downloadSet[operatorID]
		require.False(t, ok)
		downloadSet[operatorID] = struct{}{}

		chunksDownloaded.Add(uint32(len(chunks)))

		return &grpcnode.GetChunksReply{
			Chunks: chunks,
		}, nil
	}

	// the set of operators we have verified the chunks of
	verificationSet := make(map[core.OperatorID]struct{}, validatorCount)
	config.UnsafeDeserializeAndVerifyFunction = func(
		key v2.BlobKey,
		operatorID core.OperatorID,
		getChunksReply *grpcnode.GetChunksReply,
	) ([]*encoding.Frame, error) {

		// verify we have the expected blob key
		require.Equal(t, blobKey, key)

		// make sure this is for a valid operator ID
		chunks, ok := operatorChunks[operatorID]
		require.True(t, ok)

		// make sure we previously downloaded from this operator
		_, ok = downloadSet[operatorID]
		require.True(t, ok)

		// make sure we have not previously verified data from this operator
		_, ok = verificationSet[operatorID]
		require.False(t, ok)

		// make sure the chunks are the ones we expect for this operator
		require.Equal(t, len(chunks), len(getChunksReply.Chunks))
		for i, chunk := range getChunksReply.Chunks {
			require.Equal(t, chunks[i], chunk)
		}

		frames := make([]*encoding.Frame, len(chunks))
		for i, _ := range chunks {
			// Unfortunately, it's complicated to generate random frame data.
			// So just use placeholders.
			frames[i] = &encoding.Frame{}
		}

		return frames, nil
	}

	decodeCalled := atomic.Bool{}
	decodedBytes := rand.PrintableBytes(32)
	framesSentToDecoding := atomic.Uint32{}
	config.UnsafeDecodeBlobFunction = func(
		key v2.BlobKey,
		chunks []*encoding.Frame,
		indices []uint,
	) ([]byte, error) {

		// verify we have the expected blob key
		require.Equal(t, blobKey, key)

		// we shouldn't have called decode before
		require.False(t, decodeCalled.Load())
		decodeCalled.Store(true)

		framesSentToDecoding.Store(uint32(len(chunks)))

		return decodedBytes, nil
	}

	worker, err := newRetrievalWorker(
		context.Background(),
		logger,
		config,
		connectionPool,
		computePool,
		nil,
		assignments,
		totalChunkCount,
		minimumChunkCount,
		nil,
		quorumID,
		blobKey,
		nil,
		nil,
		nil)
	require.NoError(t, err)

	blob, err := worker.downloadBlobFromValidators()
	require.Equal(t, decodedBytes, blob)

	// The number of downloads should exceed the pessimistic threshold by no more than the
	// maximum chunk count of any individual operator
	pessimisticThreshold := uint32(float64(minimumChunkCount) * config.DownloadPessimism)
	maxToDownload := pessimisticThreshold + maximumChunksPerValidator
	require.GreaterOrEqual(t, maxToDownload, chunksDownloaded.Load())

	// The number of chunks verified should exceed the pessimistic threshold by no more than the
	// maximum chunk count of any individual operator
	pessimisticThreshold = uint32(float64(minimumChunkCount) * config.VerificationPessimism)
	maxToVerify := pessimisticThreshold + maximumChunksPerValidator
	require.GreaterOrEqual(t, maxToVerify, framesSentToDecoding.Load())
}
