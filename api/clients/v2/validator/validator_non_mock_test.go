package validator

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/validator/internal"
	grpcnode "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	testrandom "github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gammazero/workerpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var duplicatedIndex uint32 = 1

// TestNonMockedValidatorClientWorkflow tests validator client retrieval using real KZG prover and verifier
// This creates actual encoded blobs with proper KZG commitments, rather than using mocked components
func TestNonMockedValidatorClientWorkflow(t *testing.T) {
	// Set up KZG components (prover and verifier)
	p, v, err := makeTestEncodingComponents()
	require.NoError(t, err)

	// Set up test environment
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
	config.DownloadPessimism = rand.Float64Range(1.0, 3.0)
	config.VerificationPessimism = rand.Float64Range(1.0, 3.0)

	// Set up workerpools
	connectionPool := workerpool.New(8)
	computePool := workerpool.New(8)

	// Create test data
	quorumID := (core.QuorumID)(rand.Uint32Range(0, 10))
	quorumNumbers := []core.QuorumID{quorumID}
	validatorCount := rand.Int32Range(20, 40)

	// Create blob and get commitments
	blobVersion := v2.BlobVersion(0)
	blobHeader, data := makeTestBlob(t, p, blobVersion, 2, quorumNumbers)
	blobKey, err := blobHeader.BlobKey()
	require.NoError(t, err)

	// Create stakes and chain state
	stakes := map[core.QuorumID]map[core.OperatorID]int{
		quorumID: {},
	}
	for i := 0; i < int(validatorCount); i++ {
		operatorID := (core.OperatorID)(rand.PrintableBytes(32))
		stakes[quorumID][operatorID] = rand.Intn(100) + 1
	}

	dat, err := coremock.NewChainDataMock(stakes)
	require.NoError(t, err)

	// Prepare blobs with real encoding
	// This creates actual sharded blobs for each operator with valid KZG proofs
	operatorState, err := dat.GetOperatorState(context.Background(), 0, quorumNumbers)
	require.NoError(t, err)

	// Get encoding parameters
	encodingParams, err := v2.GetEncodingParams(blobHeader.BlobCommitments.Length, blobParams)
	require.NoError(t, err)

	// Get assignments
	assignments, err := v2.GetAssignments(operatorState, blobParams, quorumNumbers, blobKey[:])
	require.NoError(t, err)

	// Put a redundant index in each assignment
	for opID, assignment := range assignments {
		assignment.Indices = append(assignment.Indices, duplicatedIndex)
		assignments[opID] = assignment
	}

	// Create the actual blob chunks using the prover
	chunks, err := p.GetFrames(data, encodingParams)
	require.NoError(t, err)

	// Store chunks by operator
	operatorChunks := make(map[core.OperatorID][]*encoding.Frame)
	for opID, assignment := range assignments {
		operatorChunks[opID] = make([]*encoding.Frame, assignment.NumChunks())
		for i := uint32(0); i < assignment.NumChunks(); i++ {
			operatorChunks[opID][i] = chunks[assignment.Indices[i]]
		}
	}

	// Calculate max chunks per validator
	maximumChunksPerValidator := uint32(0)
	for _, assgn := range assignments {
		numChunks := assgn.NumChunks()
		if numChunks > maximumChunksPerValidator {
			maximumChunksPerValidator = numChunks
		}
	}

	// The minimum number of chunks needed for reconstruction
	minimumChunkCount := blobParams.NumChunks / blobParams.CodingRate

	// Create necessary tracking for downloads
	chunksDownloaded := sync.Map{}
	chunksDownloadedCount := atomic.Uint32{}
	downloadSet := sync.Map{}
	framesSentToDecoding := sync.Map{}
	framesSentToDecodingCount := atomic.Uint32{}

	// Create a custom GRPC manager that simulates getting frames from nodes
	grpcManager := &customValidatorGRPCManager{
		operatorChunks:        operatorChunks,
		downloadSet:           &downloadSet,
		chunksDownloaded:      &chunksDownloaded,
		chunksDownloadedCount: &chunksDownloadedCount,
		assignments:           assignments,
	}

	// Configure the client to use our custom GRPC manager but real deserializer and decoder
	config.UnsafeValidatorGRPCManagerFactory = func(logger logging.Logger, operators map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo) internal.ValidatorGRPCManager {
		return grpcManager
	}

	// We're using the real deserializer and decoder, but tracking frame counts
	originalChunkDeserializerFactory := config.UnsafeChunkDeserializerFactory
	config.UnsafeChunkDeserializerFactory = func(assignments map[core.OperatorID]v2.Assignment, verifier encoding.Verifier) internal.ChunkDeserializer {
		realDeserializer := originalChunkDeserializerFactory(assignments, verifier)
		return &instrumentedChunkDeserializer{
			ChunkDeserializer: realDeserializer,
		}
	}

	originalBlobDecoderFactory := config.UnsafeBlobDecoderFactory
	config.UnsafeBlobDecoderFactory = func(verifier encoding.Verifier) internal.BlobDecoder {
		realDecoder := originalBlobDecoderFactory(verifier)
		return &instrumentedBlobDecoder{
			t:                         t,
			BlobDecoder:               realDecoder,
			framesSentToDecoding:      &framesSentToDecoding,
			framesSentToDecodingCount: &framesSentToDecodingCount,
		}
	}

	// Create a worker with all the real components
	worker, err := newRetrievalWorker(
		context.Background(),
		logger,
		config,
		connectionPool,
		computePool,
		grpcManager,
		config.UnsafeChunkDeserializerFactory(assignments, v),
		config.UnsafeBlobDecoderFactory(v),
		assignments,
		minimumChunkCount,
		&encodingParams,
		blobHeader,
		blobKey,
		nil)
	require.NoError(t, err)

	// Execute the retrieval
	retrievedData, err := worker.retrieveBlobFromValidators()
	require.NoError(t, err)

	// Verify results
	require.Equal(t, data, retrievedData)

	// The number of downloads should exceed the pessimistic threshold by no more than the
	// maximum chunk count of any individual operator
	pessimisticDownloadThreshold := uint32(math.Ceil(float64(minimumChunkCount) * config.DownloadPessimism))
	maxToDownload := uint32(math.Ceil(float64(pessimisticDownloadThreshold)*config.VerificationPessimism)) + maximumChunksPerValidator
	fmt.Println("minimumChunkCount", minimumChunkCount)
	fmt.Println("config.DownloadPessimism", config.DownloadPessimism)
	fmt.Println("pessimisticDownloadThreshold", pessimisticDownloadThreshold)
	fmt.Println("config.VerificationPessimism", config.VerificationPessimism)
	fmt.Println("maximumChunksPerValidator", maximumChunksPerValidator)
	fmt.Println("maxToDownload", maxToDownload)
	fmt.Println("chunksDownloadedCount", chunksDownloadedCount.Load())
	require.GreaterOrEqual(t, maxToDownload, chunksDownloadedCount.Load())

	// The number of chunks verified should exceed the pessimistic threshold by no more than the
	// maximum chunk count of any individual operator
	pessimisticVerificationThreshold := uint32(math.Ceil(float64(minimumChunkCount) * config.VerificationPessimism))
	maxToVerify := pessimisticVerificationThreshold + maximumChunksPerValidator
	require.GreaterOrEqual(t, maxToVerify, framesSentToDecodingCount.Load())
}

// makeTestEncodingComponents makes a prover and verifier for KZG
func makeTestEncodingComponents() (encoding.Prover, encoding.Verifier, error) {
	config := &kzg.KzgConfig{
		G1Path:          "../../../../inabox/resources/kzg/g1.point.300000",
		G2Path:          "../../../../inabox/resources/kzg/g2.point.300000",
		CacheDir:        "../../../../inabox/resources/kzg/SRSTables",
		SRSOrder:        8192,
		SRSNumberToLoad: 8192,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    true,
	}

	p, err := prover.NewProver(config, nil)
	if err != nil {
		return nil, nil, err
	}

	v, err := verifier.NewVerifier(config, nil)
	if err != nil {
		return nil, nil, err
	}

	return p, v, nil
}

// makeTestBlob creates a test blob with valid commitments
func makeTestBlob(t *testing.T, p encoding.Prover, version v2.BlobVersion, length int, quorums []core.QuorumID) (*v2.BlobHeaderWithHashedPayment, []byte) {
	data := make([]byte, length*31)
	_, err := rand.Read(data)
	if err != nil {
		t.Fatal(err)
	}

	data = codec.ConvertByPaddingEmptyByte(data)

	commitments, err := p.GetCommitmentsForPaddedLength(data)
	if err != nil {
		t.Fatal(err)
	}

	header := &v2.BlobHeaderWithHashedPayment{
		BlobVersion:         version,
		QuorumNumbers:       quorums,
		BlobCommitments:     commitments,
		PaymentMetadataHash: [32]byte{},
	}

	return header, data
}

// customValidatorGRPCManager simulates a network of validator nodes with pre-populated chunks
type customValidatorGRPCManager struct {
	operatorChunks        map[core.OperatorID][]*encoding.Frame
	downloadSet           *sync.Map
	chunksDownloaded      *sync.Map
	chunksDownloadedCount *atomic.Uint32
	assignments           map[core.OperatorID]v2.Assignment
}

func (m *customValidatorGRPCManager) DownloadChunks(
	ctx context.Context,
	key v2.BlobKey,
	operatorID core.OperatorID,
	quorumID core.QuorumID,
) (*grpcnode.GetChunksReply, error) {
	// Make sure this is for a valid operator ID
	frames, ok := m.operatorChunks[operatorID]
	if !ok {
		return nil, nil
	}

	// Only permit downloads to happen once per operator
	_, ok = m.downloadSet.Load(operatorID)
	if ok {
		return nil, nil
	}
	m.downloadSet.Store(operatorID, struct{}{})

	// De-duplicate chunks when counting
	assgn := m.assignments[operatorID]
	numChunks := uint32(0)
	for _, i := range assgn.Indices {
		_, ok = m.chunksDownloaded.Load(i)
		if !ok {
			m.chunksDownloaded.Store(i, struct{}{})
			numChunks++
		}
	}
	m.chunksDownloadedCount.Add(numChunks)

	// Convert frames to bytes for gRPC response
	chunks := make([][]byte, len(frames))
	for i, frame := range frames {
		serialized, err := frame.SerializeGnark()
		if err != nil {
			return nil, err
		}
		chunks[i] = serialized
	}

	return &grpcnode.GetChunksReply{
		Chunks: chunks,
	}, nil
}

// instrumentedChunkDeserializer wraps a real deserializer but adds instrumentation
type instrumentedChunkDeserializer struct {
	internal.ChunkDeserializer
}

// instrumentedBlobDecoder wraps a real decoder but adds instrumentation for counting
type instrumentedBlobDecoder struct {
	t *testing.T
	internal.BlobDecoder
	framesSentToDecoding      *sync.Map
	framesSentToDecodingCount *atomic.Uint32
}

func (d *instrumentedBlobDecoder) DecodeBlob(
	blobKey v2.BlobKey,
	chunks []*encoding.Frame,
	indices []uint,
	encodingParams *encoding.EncodingParams,
	blobCommitments *encoding.BlobCommitments,
) ([]byte, error) {
	// Count de-duplicated frames
	frameCount := uint32(0)
	for _, i := range indices {
		_, ok := d.framesSentToDecoding.Load(i)
		if !ok {
			d.framesSentToDecoding.Store(i, struct{}{})
			frameCount++
		}
	}
	d.framesSentToDecodingCount.Add(frameCount)

	// Count the number of times the duplicated index was sent to decoding
	duplicatedIndexCount := 0
	for _, i := range indices {
		if i == uint(duplicatedIndex) {
			duplicatedIndexCount++
		}
	}
	assert.GreaterOrEqual(d.t, duplicatedIndexCount, 2)

	// Call the real decoder
	return d.BlobDecoder.DecodeBlob(blobKey, chunks, indices, encodingParams, blobCommitments)
}
