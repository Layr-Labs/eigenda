package validator

import (
	"context"
	"errors"
	"math"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/validator/mock"
	grpcnode "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testutils"
	testrandom "github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/gammazero/workerpool"
	"github.com/stretchr/testify/require"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"

	coremock "github.com/Layr-Labs/eigenda/core/mock"
)

var (
	blobParams = &core.BlobVersionParameters{
		NumChunks:                   8192,
		CodingRate:                  8,
		ReconstructionThresholdBips: 1666,
		SamplesPerUnit:              20,
		NumUnits:                    393,
	}

	blobKey1 = []byte("hello0")
)

func MakeRandomAssignment(t *testing.T, rand *testrandom.TestRandom, validatorCount int32, quorumID core.QuorumID) map[core.OperatorID]v2.Assignment {

	stakes := map[core.QuorumID]map[core.OperatorID]int{
		quorumID: {},
	}
	for i := 0; i < int(validatorCount); i++ {
		operatorID := (core.OperatorID)(rand.PrintableBytes(32))
		stakes[quorumID][operatorID] = rand.Intn(100) + 1
	}

	dat, err := coremock.NewChainDataMock(stakes)
	if err != nil {
		t.Fatal(err)
	}

	state := dat.GetTotalOperatorState(context.Background(), 0)

	assignments, err := v2.GetAssignments(state.OperatorState, blobParams, []core.QuorumID{quorumID}, blobKey1[:])
	require.NoError(t, err)

	return assignments
}

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
	maximumChunksPerValidator := uint32(0)
	quorumID := (core.QuorumID)(rand.Uint32Range(0, 10))

	// Simulated chunks for each operator
	operatorChunks := make(map[core.OperatorID][][]byte, validatorCount)

	// Create assignment
	assignments := MakeRandomAssignment(t, rand, validatorCount, quorumID)

	for operatorID, assgn := range assignments {

		numChunks := assgn.NumChunks()
		if numChunks > uint32(maximumChunksPerValidator) {
			maximumChunksPerValidator = numChunks
		}

		operatorChunks[operatorID] = make([][]byte, numChunks)
		for j := 0; j < int(numChunks); j++ {
			operatorChunks[operatorID][j] = rand.PrintableBytes(8)
		}
	}

	// The number of chunks needed to reconstruct the blob
	minimumChunkCount := blobParams.NumChunks / blobParams.CodingRate

	// the number of chunks downloaded
	chunksDownloaded := sync.Map{}
	chunksDownloadedCount := atomic.Uint32{}
	// a set of operators that have provided chunks
	downloadSet := sync.Map{}
	mockGRPCManager := &mock.MockValidatorGRPCManager{}
	mockGRPCManager.DownloadChunksFunction = func(
		ctx context.Context,
		key v2.BlobKey,
		operatorID core.OperatorID,
		quorumID core.QuorumID,
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

		// TODO(Claude): Please apply this de-duplication logic to the other tests as well. It also needs to be performed for the
		// framesSentToDecoding variable.
		assgn := assignments[operatorID]
		numChunks := uint32(0)
		for _, i := range assgn.Indices {
			_, ok = chunksDownloaded.Load(i)
			if !ok {
				chunksDownloaded.Store(i, struct{}{})
				numChunks++
			}
		}
		chunksDownloadedCount.Add(numChunks)

		return &grpcnode.GetChunksReply{
			Chunks: chunks,
		}, nil
	}

	// the set of operators we have verified the chunks of
	verificationSet := sync.Map{}
	mockDeserializer := &mock.MockChunkDeserializer{}
	mockDeserializer.DeserializeAndVerifyFunction = func(
		blobKey v2.BlobKey,
		operatorID core.OperatorID,
		getChunksReply *grpcnode.GetChunksReply,
		blobCommitments *encoding.BlobCommitments,
		encodingParams *encoding.EncodingParams,
	) ([]*encoding.Frame, error) {

		// verify we have the expected blob key
		require.Equal(t, blobKey, blobKey)

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
	framesSentToDecoding := sync.Map{}
	framesSentToDecodingCount := atomic.Uint32{}

	mockDecoder := &mock.MockBlobDecoder{}
	mockDecoder.DecodeBlobFunction = func(
		key v2.BlobKey,
		chunks []*encoding.Frame,
		indices []uint,
		encodingParams *encoding.EncodingParams,
		blobCommitments *encoding.BlobCommitments,
	) ([]byte, error) {

		// verify we have the expected blob key
		require.Equal(t, blobKey, key)

		// we shouldn't have called decode before
		require.False(t, decodeCalled.Load())
		decodeCalled.Store(true)

		// De-duplicate frames when counting
		frameCount := uint32(0)
		for _, i := range indices {
			_, ok := framesSentToDecoding.Load(i)
			if !ok {
				framesSentToDecoding.Store(i, struct{}{})
				frameCount++
			}
		}
		framesSentToDecodingCount.Add(frameCount)

		return decodedBytes, nil
	}

	blobHeader := &v2.BlobHeaderWithoutPayment{
		BlobVersion:         0,
		QuorumNumbers:       []core.QuorumID{quorumID},
		BlobCommitments:     MockCommitment(t),
		PaymentMetadataHash: [32]byte{},
	}

	worker, err := newRetrievalWorker(
		context.Background(),
		logger,
		config,
		connectionPool,
		computePool,
		mockGRPCManager,
		mockDeserializer,
		mockDecoder,
		assignments,
		minimumChunkCount,
		nil,
		blobHeader,
		blobKey,
		nil)
	require.NoError(t, err)

	blob, err := worker.retrieveBlobFromValidators()
	require.NoError(t, err)
	require.Equal(t, decodedBytes, blob)

	// The number of downloads should exceed the pessimistic threshold by no more than the
	// maximum chunk count of any individual operator
	pessimisticDownloadThreshold := uint32(math.Ceil(float64(minimumChunkCount) * config.DownloadPessimism))
	maxToDownload := uint32(math.Ceil(float64(pessimisticDownloadThreshold)*config.VerificationPessimism)) +
		maximumChunksPerValidator
	require.GreaterOrEqual(t, maxToDownload, chunksDownloadedCount.Load())

	// The number of chunks verified should exceed the pessimistic threshold by no more than the
	// maximum chunk count of any individual operator
	pessimisticVerificationThreshold := uint32(math.Ceil(float64(minimumChunkCount) * config.VerificationPessimism))
	maxToVerify := pessimisticVerificationThreshold + maximumChunksPerValidator
	require.GreaterOrEqual(t, maxToVerify, framesSentToDecodingCount.Load())
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
	maximumChunksPerValidator := uint32(0)
	quorumID := (core.QuorumID)(rand.Uint32Range(0, 10))

	connectionPool := workerpool.New(int(validatorCount))
	computePool := workerpool.New(int(validatorCount))

	// The assignments for each operator (i.e. how many chunks it is responsible for)

	assignments := MakeRandomAssignment(t, rand, validatorCount, quorumID)

	// Simulated chunks for each operator
	operatorChunks := make(map[core.OperatorID][][]byte, validatorCount)
	// Allows the test to block a download, download does not complete until element is inserted into chan
	downloadLocks := make(map[core.OperatorID]chan struct{}, validatorCount)

	for operatorID, assgn := range assignments {
		numChunks := assgn.NumChunks()
		if numChunks > uint32(maximumChunksPerValidator) {
			maximumChunksPerValidator = numChunks
		}
		operatorChunks[operatorID] = make([][]byte, numChunks)
		for j := 0; j < int(numChunks); j++ {
			operatorChunks[operatorID][j] = rand.PrintableBytes(8)
		}
		downloadLocks[operatorID] = make(chan struct{}, 1)
	}

	// The number of chunks needed to reconstruct the blob
	minimumChunkCount := blobParams.NumChunks / blobParams.CodingRate

	// the number of chunks downloaded (de-duplicated)
	chunksDownloaded := sync.Map{}
	chunksDownloadedCount := atomic.Uint32{}
	timedOutDownloads := atomic.Uint32{}
	// a set of operators that have provided chunks
	downloadSet := sync.Map{}
	mockGRPCManager := &mock.MockValidatorGRPCManager{}
	mockGRPCManager.DownloadChunksFunction = func(
		ctx context.Context,
		key v2.BlobKey,
		operatorID core.OperatorID,
		quorumID core.QuorumID,
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

		// De-duplicate chunks when counting
		assgn := assignments[operatorID]
		numChunks := uint32(0)
		for _, i := range assgn.Indices {
			_, ok = chunksDownloaded.Load(i)
			if !ok {
				chunksDownloaded.Store(i, struct{}{})
				numChunks++
			}
		}
		chunksDownloadedCount.Add(numChunks)

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
	mockDeserializer := &mock.MockChunkDeserializer{}
	mockDeserializer.DeserializeAndVerifyFunction = func(
		blobKey v2.BlobKey,
		operatorID core.OperatorID,
		getChunksReply *grpcnode.GetChunksReply,
		blobCommitments *encoding.BlobCommitments,
		encodingParams *encoding.EncodingParams,
	) ([]*encoding.Frame, error) {

		// verify we have the expected blob key
		require.Equal(t, blobKey, blobKey)

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
	framesSentToDecoding := sync.Map{}
	framesSentToDecodingCount := atomic.Uint32{}
	mockDecoder := &mock.MockBlobDecoder{}
	mockDecoder.DecodeBlobFunction = func(
		key v2.BlobKey,
		chunks []*encoding.Frame,
		indices []uint,
		encodingParams *encoding.EncodingParams,
		blobCommitments *encoding.BlobCommitments,
	) ([]byte, error) {

		// verify we have the expected blob key
		require.Equal(t, blobKey, key)

		// we shouldn't have called decode before
		require.False(t, decodeCalled.Load())
		decodeCalled.Store(true)

		// De-duplicate frames when counting
		frameCount := uint32(0)
		for _, i := range indices {
			_, ok := framesSentToDecoding.Load(i)
			if !ok {
				framesSentToDecoding.Store(i, struct{}{})
				frameCount++
			}
		}
		framesSentToDecodingCount.Add(frameCount)

		return decodedBytes, nil
	}

	blobHeader := &v2.BlobHeaderWithoutPayment{
		BlobVersion:         0,
		QuorumNumbers:       []core.QuorumID{quorumID},
		BlobCommitments:     MockCommitment(t),
		PaymentMetadataHash: [32]byte{},
	}

	worker, err := newRetrievalWorker(
		context.Background(),
		logger,
		config,
		connectionPool,
		computePool,
		mockGRPCManager,
		mockDeserializer,
		mockDecoder,
		assignments,
		minimumChunkCount,
		nil,
		blobHeader,
		blobKey,
		nil)
	require.NoError(t, err)

	downloadFinishedChan := make(chan struct{}, 1)
	var downloadFinished bool
	var blob []byte
	go func() {
		blob, err = worker.retrieveBlobFromValidators()
		require.Equal(t, decodedBytes, blob)
		downloadFinished = true
		downloadFinishedChan <- struct{}{}
	}()

	pessimisticDownloadThreshold := uint32(math.Ceil(float64(minimumChunkCount) * config.DownloadPessimism))

	// Wait until we've scheduled all the downloads.
	testutils.AssertEventuallyTrue(
		t,
		func() bool {
			return chunksDownloadedCount.Load() >= pessimisticDownloadThreshold
		},
		time.Second)

	require.False(t, downloadFinished)
	initialDownloadsScheduled := chunksDownloadedCount.Load()

	// Advance the clock past the point when pessimistic thresholds trigger for the download.
	newTime := start.Add(config.PessimisticTimeout + 1*time.Second)
	fakeClock.Store(&newTime)

	// Wait until we've scheduled the additional downloads.
	testutils.AssertEventuallyTrue(
		t,
		func() bool {
			return chunksDownloadedCount.Load()-initialDownloadsScheduled >= pessimisticDownloadThreshold
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
	require.GreaterOrEqual(t, maxToVerify, framesSentToDecodingCount.Load())
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
	maximumChunksPerValidator := uint32(0)
	quorumID := (core.QuorumID)(rand.Uint32Range(0, 10))

	// The assignments for each operator (i.e. how many chunks it is responsible for)
	assignments := MakeRandomAssignment(t, rand, validatorCount, quorumID)
	// Simulated chunks for each operator
	operatorChunks := make(map[core.OperatorID][][]byte, validatorCount)

	for operatorID, assgn := range assignments {
		numChunks := assgn.NumChunks()
		if numChunks > uint32(maximumChunksPerValidator) {
			maximumChunksPerValidator = numChunks
		}
		operatorChunks[operatorID] = make([][]byte, numChunks)
		for j := 0; j < int(numChunks); j++ {
			operatorChunks[operatorID][j] = rand.PrintableBytes(8)
		}
	}

	// The number of chunks needed to reconstruct the blob
	minimumChunkCount := blobParams.NumChunks / blobParams.CodingRate

	// the number of chunks downloaded (de-duplicated)
	chunksDownloaded := sync.Map{}
	chunksDownloadedCount := atomic.Uint32{}
	// a set of operators that have provided chunks
	downloadSet := sync.Map{}
	mockGRPCManager := &mock.MockValidatorGRPCManager{}
	mockGRPCManager.DownloadChunksFunction = func(
		ctx context.Context,
		key v2.BlobKey,
		operatorID core.OperatorID,
		quorumID core.QuorumID,
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

		// De-duplicate chunks when counting
		assgn := assignments[operatorID]
		numChunks := uint32(0)
		for _, i := range assgn.Indices {
			_, ok = chunksDownloaded.Load(i)
			if !ok {
				chunksDownloaded.Store(i, struct{}{})
				numChunks++
			}
		}
		chunksDownloadedCount.Add(numChunks)

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
	mockDeserializer := &mock.MockChunkDeserializer{}
	mockDeserializer.DeserializeAndVerifyFunction = func(
		blobKey v2.BlobKey,
		operatorID core.OperatorID,
		getChunksReply *grpcnode.GetChunksReply,
		blobCommitments *encoding.BlobCommitments,
		encodingParams *encoding.EncodingParams,
	) ([]*encoding.Frame, error) {

		// verify we have the expected blob key
		require.Equal(t, blobKey, blobKey)

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
	framesSentToDecoding := sync.Map{}
	framesSentToDecodingCount := atomic.Uint32{}
	mockDecoder := &mock.MockBlobDecoder{}
	mockDecoder.DecodeBlobFunction = func(
		key v2.BlobKey,
		chunks []*encoding.Frame,
		indices []uint,
		encodingParams *encoding.EncodingParams,
		blobCommitments *encoding.BlobCommitments,
	) ([]byte, error) {

		// verify we have the expected blob key
		require.Equal(t, blobKey, key)

		// we shouldn't have called decode before
		require.False(t, decodeCalled.Load())
		decodeCalled.Store(true)

		// De-duplicate frames when counting
		frameCount := uint32(0)
		for _, i := range indices {
			_, ok := framesSentToDecoding.Load(i)
			if !ok {
				framesSentToDecoding.Store(i, struct{}{})
				frameCount++
			}
		}
		framesSentToDecodingCount.Add(frameCount)

		return decodedBytes, nil
	}

	blobHeader := &v2.BlobHeaderWithoutPayment{
		BlobVersion:         0,
		QuorumNumbers:       []core.QuorumID{quorumID},
		BlobCommitments:     MockCommitment(t),
		PaymentMetadataHash: [32]byte{},
	}

	worker, err := newRetrievalWorker(
		context.Background(),
		logger,
		config,
		connectionPool,
		computePool,
		mockGRPCManager,
		mockDeserializer,
		mockDecoder,
		assignments,
		minimumChunkCount,
		nil,
		blobHeader,
		blobKey,
		nil)
	require.NoError(t, err)

	blob, err := worker.retrieveBlobFromValidators()
	require.NoError(t, err)
	require.Equal(t, decodedBytes, blob)

	// The number of downloads should exceed the pessimistic threshold by no more than the
	// maximum chunk count of any individual operator, plus the number of failed verifications.
	pessimisticDownloadThreshold := uint32(math.Ceil(float64(minimumChunkCount) * config.DownloadPessimism))
	maxToDownload := uint32(math.Ceil(float64(pessimisticDownloadThreshold)*config.VerificationPessimism)) +
		maximumChunksPerValidator + failedChunkCount.NumChunks()
	require.GreaterOrEqual(t, maxToDownload, chunksDownloadedCount.Load())

	// The number of chunks verified should exceed the pessimistic threshold by no more than the
	// maximum chunk count of any individual operator, plus the number of failed verifications.
	pessimisticVerificationThreshold := uint32(math.Ceil(float64(minimumChunkCount) * config.VerificationPessimism))
	maxToVerify := pessimisticVerificationThreshold + maximumChunksPerValidator + failedChunkCount.NumChunks()
	require.GreaterOrEqual(t, maxToVerify, framesSentToDecodingCount.Load())
}

func MockCommitment(t *testing.T) encoding.BlobCommitments {
	var X1, Y1 fp.Element
	X1 = *X1.SetBigInt(big.NewInt(1))
	Y1 = *Y1.SetBigInt(big.NewInt(2))

	var lengthXA0, lengthXA1, lengthYA0, lengthYA1 fp.Element
	_, err := lengthXA0.SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781")
	require.NoError(t, err)
	_, err = lengthXA1.SetString("11559732032986387107991004021392285783925812861821192530917403151452391805634")
	require.NoError(t, err)
	_, err = lengthYA0.SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930")
	require.NoError(t, err)
	_, err = lengthYA1.SetString("4082367875863433681332203403145435568316851327593401208105741076214120093531")
	require.NoError(t, err)

	var lengthProof, lengthCommitment bn254.G2Affine
	lengthProof.X.A0 = lengthXA0
	lengthProof.X.A1 = lengthXA1
	lengthProof.Y.A0 = lengthYA0
	lengthProof.Y.A1 = lengthYA1

	lengthCommitment = lengthProof

	return encoding.BlobCommitments{
		Commitment: &encoding.G1Commitment{
			X: X1,
			Y: Y1,
		},
		LengthCommitment: (*encoding.G2Commitment)(&lengthCommitment),
		LengthProof:      (*encoding.G2Commitment)(&lengthProof),
		Length:           10,
	}
}
