package validator

import (
	"context"
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	grpcnode "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testutils"
	testrandom "github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/gammazero/workerpool"
	"github.com/stretchr/testify/require"
)

func TestBasicWorkflow(t *testing.T) {
	rand := testrandom.NewTestRandom()
	start := rand.Time()

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	fakeClock := atomic.Pointer[time.Time]{}
	fakeClock.Store(&start)

	config := DefaultClientConfig()
	config.ControlLoopPeriod = 50 * time.Microsecond
	config.TimeSource = func() time.Time {
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
	minimumChunkCount := uint32(rand.Float64Range(1.0/8.0, 1.0/8.0) * float64(totalChunkCount))

	// the number of chunks downloaded
	chunksDownloaded := atomic.Uint32{}
	// a set of operators that have provided chunks
	downloadSet := sync.Map{}
	config.UnsafeDownloadChunksFunction = func(
		ctx context.Context,
		key v2.BlobKey,
		operatorID core.OperatorID,
	) (*grpcnode.GetChunksReply, error) {

		// verify we have the expected blob key
		require.Equal(t, blobKey, key)

		// make sure this is for a valid operator ID
		chunks, ok := operatorChunks[operatorID]
		require.True(t, ok)

		// only permit downloads to happen once per operator
		_, ok = downloadSet.Load(operatorID)
		require.False(t, ok)
		downloadSet.Store(operatorID, struct{}{})

		chunksDownloaded.Add(uint32(len(chunks)))

		return &grpcnode.GetChunksReply{
			Chunks: chunks,
		}, nil
	}

	// the set of operators we have verified the chunks of
	verificationSet := sync.Map{}
	config.UnsafeDeserializeAndVerifyFunction = func(
		ctx context.Context,
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
		_, ok = downloadSet.Load(operatorID)
		require.True(t, ok)

		// make sure we have not previously verified data from this operator
		_, ok = verificationSet.Load(operatorID)
		require.False(t, ok)
		verificationSet.Store(operatorID, struct{}{})

		// make sure the chunks are the ones we expect for this operator
		require.Equal(t, len(chunks), len(getChunksReply.Chunks))
		for i, chunk := range getChunksReply.Chunks {
			require.Equal(t, chunks[i], chunk)
		}

		frames := make([]*encoding.Frame, len(chunks))
		for i := range chunks {
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
		ctx context.Context,
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
	require.NoError(t, err)
	require.Equal(t, decodedBytes, blob)

	// The number of downloads should exceed the pessimistic threshold by no more than the
	// maximum chunk count of any individual operator
	pessimisticDownloadThreshold := uint32(math.Ceil(float64(minimumChunkCount) * config.DownloadPessimism))
	maxToDownload := uint32(math.Ceil(float64(pessimisticDownloadThreshold)*config.VerificationPessimism)) +
		maximumChunksPerValidator
	require.GreaterOrEqual(t, maxToDownload, chunksDownloaded.Load())

	// The number of chunks verified should exceed the pessimistic threshold by no more than the
	// maximum chunk count of any individual operator
	pessimisticVerificationThreshold := uint32(math.Ceil(float64(minimumChunkCount) * config.VerificationPessimism))
	maxToVerify := pessimisticVerificationThreshold + maximumChunksPerValidator
	require.GreaterOrEqual(t, maxToVerify, framesSentToDecoding.Load())
}

func TestDownloadTimeout(t *testing.T) {
	rand := testrandom.NewTestRandom()
	start := rand.Time()

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	fakeClock := atomic.Pointer[time.Time]{}
	fakeClock.Store(&start)

	config := DefaultClientConfig()
	config.ControlLoopPeriod = 50 * time.Microsecond
	config.TimeSource = func() time.Time {
		return *fakeClock.Load()
	}
	config.DownloadPessimism = rand.Float64Range(1.0, 2.0)
	config.VerificationPessimism = rand.Float64Range(1.0, 2.0)
	config.PessimisticTimeout = time.Second
	config.DownloadTimeout = 10 * time.Second

	blobKey := (v2.BlobKey)(rand.Bytes(32))

	validatorCount := rand.Int32Range(50, 100)
	minimumChunksPerValidator := uint32(1)
	maximumChunksPerValidator := rand.Uint32Range(50, 100)
	quorumID := (core.QuorumID)(rand.Uint32Range(0, 10))

	connectionPool := workerpool.New(int(validatorCount))
	computePool := workerpool.New(int(validatorCount))

	// The assignments for each operator (i.e. how many chunks it is responsible for)
	assignments := make(map[core.OperatorID]v2.Assignment, validatorCount)
	// The total number of chunks
	totalChunkCount := uint32(0)
	// Simulated chunks for each operator
	operatorChunks := make(map[core.OperatorID][][]byte, validatorCount)
	// Allows the test to block a download, download does not complete until element is inserted into chan
	downloadLocks := make(map[core.OperatorID]chan struct{}, validatorCount)

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
		downloadLocks[operatorID] = make(chan struct{}, 1)
	}

	// The number of chunks needed to reconstruct the blob
	minimumChunkCount := uint32(rand.Float64Range(1.0/8.0, 1.0/8.0) * float64(totalChunkCount))

	// the number of chunks downloaded
	chunksDownloaded := atomic.Uint32{}
	timedOutDownloads := atomic.Uint32{}
	// a set of operators that have provided chunks
	downloadSet := sync.Map{}
	config.UnsafeDownloadChunksFunction = func(
		ctx context.Context,
		key v2.BlobKey,
		operatorID core.OperatorID,
	) (*grpcnode.GetChunksReply, error) {
		// verify we have the expected blob key
		require.Equal(t, blobKey, key)

		// make sure this is for a valid operator ID
		chunks, ok := operatorChunks[operatorID]
		require.True(t, ok)

		// only permit downloads to happen once per operator
		_, ok = downloadSet.Load(operatorID)
		require.False(t, ok)
		downloadSet.Store(operatorID, struct{}{})

		chunksDownloaded.Add(uint32(len(chunks)))

		// wait until the download is unlocked
		select {
		case <-downloadLocks[operatorID]:
		case <-ctx.Done():
			timedOutDownloads.Add(uint32(len(chunks)))
		}

		return &grpcnode.GetChunksReply{
			Chunks: chunks,
		}, nil
	}

	// the set of operators we have verified the chunks of
	verificationSet := sync.Map{}
	config.UnsafeDeserializeAndVerifyFunction = func(
		ctx context.Context,
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
		_, ok = downloadSet.Load(operatorID)
		require.True(t, ok)

		// make sure we have not previously verified data from this operator
		_, ok = verificationSet.Load(operatorID)
		require.False(t, ok)
		verificationSet.Store(operatorID, struct{}{})

		// make sure the chunks are the ones we expect for this operator
		require.Equal(t, len(chunks), len(getChunksReply.Chunks))
		for i, chunk := range getChunksReply.Chunks {
			require.Equal(t, chunks[i], chunk)
		}

		frames := make([]*encoding.Frame, len(chunks))
		for i := range chunks {
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
		ctx context.Context,
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

	downloadFinishedChan := make(chan struct{}, 1)
	var downloadFinished bool
	var blob []byte
	go func() {
		blob, err = worker.downloadBlobFromValidators()
		require.Equal(t, decodedBytes, blob)
		downloadFinished = true
		downloadFinishedChan <- struct{}{}
	}()

	pessimisticDownloadThreshold := uint32(math.Ceil(float64(minimumChunkCount) * config.DownloadPessimism))

	// Wait until we've scheduled all the downloads.
	testutils.AssertEventuallyTrue(
		t,
		func() bool {
			return chunksDownloaded.Load() >= pessimisticDownloadThreshold
		},
		time.Second)

	require.False(t, downloadFinished)
	initialDownloadsScheduled := chunksDownloaded.Load()

	// Advance the clock past the point when pessimistic thresholds trigger for the download.
	newTime := start.Add(config.PessimisticTimeout + 1*time.Second)
	fakeClock.Store(&newTime)

	// Wait until we've scheduled the additional downloads.
	testutils.AssertEventuallyTrue(
		t,
		func() bool {
			return chunksDownloaded.Load()-initialDownloadsScheduled >= pessimisticDownloadThreshold
		},
		time.Second)

	require.False(t, downloadFinished)

	// None of the downloads should have hit the full timeout.
	require.Equal(t, uint32(0), timedOutDownloads.Load())

	// Now, unblock all the downloads except for one. We will let that one download timeout.
	unblockCount := 0
	for operatorID := range downloadLocks {
		// Don't unblock the last one
		if unblockCount == int(validatorCount-1) {
			break
		}
		unblockCount++

		downloadLocks[operatorID] <- struct{}{}
	}

	// Advance the clock to trigger the timeout.
	newTime = start.Add(config.DownloadTimeout + 100*time.Second)

	// At least one download should time out. In theory, multiple downloads could have timed out (since
	// the context is cancelled after the download is completed).
	testutils.AssertEventuallyTrue(t, func() bool {
		return timedOutDownloads.Load() > uint32(0)
	}, time.Second)

	// Wait for the blob to be downloaded.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	select {
	case <-downloadFinishedChan:
		// continue with test
	case <-ctx.Done():
		require.Fail(t, "download did not finish in time")
	}

	// The number of chunks verified should exceed the pessimistic threshold by no more than the
	// maximum chunk count of any individual operator
	pessimisticVerificationThreshold := uint32(math.Ceil(float64(minimumChunkCount) * config.VerificationPessimism))
	maxToVerify := pessimisticVerificationThreshold + maximumChunksPerValidator
	require.GreaterOrEqual(t, maxToVerify, framesSentToDecoding.Load())
}

func TestFailedVerification(t *testing.T) {
	rand := testrandom.NewTestRandom()
	start := rand.Time()

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	fakeClock := atomic.Pointer[time.Time]{}
	fakeClock.Store(&start)

	config := DefaultClientConfig()
	config.ControlLoopPeriod = 50 * time.Microsecond
	config.TimeSource = func() time.Time {
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
	minimumChunkCount := uint32(rand.Float64Range(1.0/8.0, 1.0/8.0) * float64(totalChunkCount))

	// the number of chunks downloaded
	chunksDownloaded := atomic.Uint32{}
	// a set of operators that have provided chunks
	downloadSet := sync.Map{}
	config.UnsafeDownloadChunksFunction = func(
		ctx context.Context,
		key v2.BlobKey,
		operatorID core.OperatorID,
	) (*grpcnode.GetChunksReply, error) {

		// verify we have the expected blob key
		require.Equal(t, blobKey, key)

		// make sure this is for a valid operator ID
		chunks, ok := operatorChunks[operatorID]
		require.True(t, ok)

		// only permit downloads to happen once per operator
		_, ok = downloadSet.Load(operatorID)
		require.False(t, ok)
		downloadSet.Store(operatorID, struct{}{})

		chunksDownloaded.Add(uint32(len(chunks)))

		return &grpcnode.GetChunksReply{
			Chunks: chunks,
		}, nil
	}

	// Intentionally cause this operator to fail verification
	var operatorWithInvalidChunks core.OperatorID
	for operatorID := range operatorChunks {
		operatorWithInvalidChunks = operatorID
		break
	}
	failedChunkCount := assignments[operatorWithInvalidChunks]

	// the set of operators we have verified the chunks of
	verificationSet := sync.Map{}
	config.UnsafeDeserializeAndVerifyFunction = func(
		ctx context.Context,
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
		_, ok = downloadSet.Load(operatorID)
		require.True(t, ok)

		// make sure we have not previously verified data from this operator
		_, ok = verificationSet.Load(operatorID)
		require.False(t, ok)
		verificationSet.Store(operatorID, struct{}{})

		// make sure the chunks are the ones we expect for this operator
		require.Equal(t, len(chunks), len(getChunksReply.Chunks))
		for i, chunk := range getChunksReply.Chunks {
			require.Equal(t, chunks[i], chunk)
		}

		if operatorID == operatorWithInvalidChunks {
			return nil, errors.New("this is an intentional failure")
		}

		frames := make([]*encoding.Frame, len(chunks))
		for i := range chunks {
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
		ctx context.Context,
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
	require.NoError(t, err)
	require.Equal(t, decodedBytes, blob)

	// The number of downloads should exceed the pessimistic threshold by no more than the
	// maximum chunk count of any individual operator, plus the number of failed verifications.
	pessimisticDownloadThreshold := uint32(math.Ceil(float64(minimumChunkCount) * config.DownloadPessimism))
	maxToDownload := uint32(math.Ceil(float64(pessimisticDownloadThreshold)*config.VerificationPessimism)) +
		maximumChunksPerValidator + failedChunkCount.NumChunks
	require.GreaterOrEqual(t, maxToDownload, chunksDownloaded.Load())

	// The number of chunks verified should exceed the pessimistic threshold by no more than the
	// maximum chunk count of any individual operator, plus the number of failed verifications.
	pessimisticVerificationThreshold := uint32(math.Ceil(float64(minimumChunkCount) * config.VerificationPessimism))
	maxToVerify := pessimisticVerificationThreshold + maximumChunksPerValidator + failedChunkCount.NumChunks
	require.GreaterOrEqual(t, maxToVerify, framesSentToDecoding.Load())
}
