package validator

import (
	"context"
	"errors"
	"fmt"
	"time"

	grpcnode "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/docker/go-units"
	"github.com/emirpasic/gods/queues"
	"github.com/emirpasic/gods/queues/linkedlistqueue"
	"github.com/gammazero/workerpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ValidatorRetrievalConfig contains the configuration for the validator retrieval client.
type ValidatorRetrievalConfig struct {
	// If 1.0, then the validator retrieval client will attempt to download the exact number of chunks
	// needed to reconstruct the blob. If higher than 1.0, then the validator retrieval client will
	// pessimistically assume that some operators will not respond in time, and will download
	// additional chunks from other operators. For example, at 2.0, the validator retrieval client
	// will download twice the number of chunks needed to reconstruct the blob. Setting this to below
	// 1.0 is not supported.
	DownloadPessimism float64
	// If 1.0, then the validator retrieval client will attempt to verify the exact number of chunks
	// needed to reconstruct the blob. If higher than 1.0, then the validator retrieval client will
	// pessimistically assume that some operators sent invalid chunks, and will verify additional chunks
	// from other operators. For example, at 2.0, the validator retrieval client will verify twice the number of
	// chunks needed to reconstruct the blob. Setting this to below 1.0 is not supported.
	VerificationPessimism float64
	// After this amount of time passes, the validator retrieval client will assume that the operator is not
	// responding, and will start downloading from a different operator. The download is not terminated when
	// this timeout is reached.
	PessimisticTimeout time.Duration
	// The absolute limit on the time to wait for a download to complete. If this timeout is reached, the
	// download will be terminated.
	DownloadTimeout time.Duration
	// The control loop periodically wakes up to do work. This is the period of that control loop.
	ControlLoopPeriod time.Duration
}

// retrievalWorker implements the distributed retrieval process for a specified blob (i.e. reading the blobs from
// validators). It is responsible for coordinating the lifecycle of this retrieval workflow.
type retrievalWorker struct {
	ctx                      context.Context
	downloadAndVerifyContext context.Context
	cancel                   context.CancelFunc
	logger                   logging.Logger
	config                   *ValidatorRetrievalConfig

	// A pool of workers for network intensive operations (e.g. downloading blob data).
	connectionPool *workerpool.WorkerPool

	// A pool of workers for CPU intensive operations (e.g. deserializing and verifying blob data).
	computePool *workerpool.WorkerPool

	// The operators that are expected to have chunk data for the blob.
	operators map[core.OperatorID]*core.OperatorInfo

	// Information about operators.
	operatorState *core.OperatorState

	// The encoding parameters for the blob.
	encodingParams encoding.EncodingParams

	// The assignments for the operators, i.e. which operators are responsible for which chunks.
	assignments map[core.OperatorID]v2.Assignment

	// The quorum to download from.
	quorumID core.QuorumID

	// The blob key to download.
	blobKey v2.BlobKey

	// Verifies chunk data.
	verifier encoding.Verifier

	// The commitments for the blob.
	blobCommitments encoding.BlobCommitments

	// When a thread begins downloading chunk data, it will send a message to the downloadStartedChan.
	downloadStartedChan chan *downloadStarted

	// When a thread completes downloading chunk data, it will send a message to the downloadCompletedChan.
	downloadCompletedChan chan *validatorDownloadCompleted

	// When a thread completes verifying chunk data, it will send a message to the verificationCompleteChan.
	verificationCompleteChan chan *verificationCompleted

	// When a thread completes decoding chunk data, it will send a message to the decodeResponseChan.
	decodeResponseChan chan *decodeResult

	// Used to collect metrics.
	probe *common.SequenceProbe

	///////////////////////////////////////////////////////////////////////////////////////////////////
	// All variables below this line are for use only on the downloadBlobFromValidators() goroutine  //
	///////////////////////////////////////////////////////////////////////////////////////////////////

	// The order in which to download chunks from the operators.
	downloadOrder []core.OperatorID

	// The index of the next operator to download from.
	nextDownloadIndex int

	// TODO verify technical accuracy of this
	// The minimum number of chunks needed to reconstruct the blob.
	minimumChunkCount uint32

	// The number of chunks we'd like to download.
	targetDownloadCount uint32

	// The number of chunks we'd like to verify.
	targetVerifiedCount uint32

	// The number of chunks currently in the process of being downloaded, or which have already been downloaded.
	// If chunk data is found to be invalid, then those chunks will be removed from this count.
	chunksBeingDownloaded uint32

	// The number of chunks currently in the process of being validated, or which have already been validated.
	// If chunk data is found to be invalid, then those chunks will be removed from this count.
	chunksBeingVerified uint32

	// This queue is used to determine when the pessimistic timeout for a download has been reached.
	downloadsInProgress queues.Queue

	// The set of validators that have completed their download (successfully or unsuccessfully).
	completedDownloadSet map[core.OperatorID]struct{}

	// Contains chunks that have been downloaded but not yet verified.
	unverifiedChunks queues.Queue

	// Contains chunks that have been verified.
	verifiedChunks queues.Queue

	// The number of chunks that have been verified.
	verifiedChunkCount uint32
}

// downloadStarted is used to signal that a download of chunk data has been initiated.
type downloadStarted struct {
	operatorID    core.OperatorID
	downloadStart time.Time
}

// downloadCompleted is used to signal that a download of chunk data has completed.
type validatorDownloadCompleted struct {
	OperatorID core.OperatorID
	Reply      *grpcnode.GetChunksReply
	Err        error
}

// verificationStarted is used to signal that verification of chunk data has completed.
type verificationCompleted struct {
	OperatorID core.OperatorID
	chunks     []*encoding.Frame
	Err        error
}

// decodeResult is used to signal that decoding of chunk data has completed.
type decodeResult struct {
	Blob []byte
	Err  error
}

// DefaultValidatorRetrievalConfig returns the default configuration for the validator retrieval client.
func DefaultValidatorRetrievalConfig() *ValidatorRetrievalConfig {
	return &ValidatorRetrievalConfig{
		DownloadPessimism:     2.0,
		VerificationPessimism: 1.0,
		PessimisticTimeout:    10 * time.Second,
		DownloadTimeout:       30 * time.Second,
		ControlLoopPeriod:     1 * time.Second,
	}
}

// newRetrievalWorker creates a new retrieval worker.
func newRetrievalWorker(
	ctx context.Context,
	logger logging.Logger,
	config *ValidatorRetrievalConfig,
	connectionPool *workerpool.WorkerPool,
	computePool *workerpool.WorkerPool,
	chainState core.ChainState,
	referenceBlockNumber uint64,
	blobVersion v2.BlobVersion,
	ethClient core.Reader,
	quorumID core.QuorumID,
	blobKey v2.BlobKey,
	verifier encoding.Verifier,
	blobCommitments encoding.BlobCommitments,
	probe *common.SequenceProbe,
) (*retrievalWorker, error) {
	if config.DownloadPessimism < 1.0 {
		return nil, fmt.Errorf("downloadPessimism must be greater than 1.0, got %f", config.DownloadPessimism)
	}
	if config.VerificationPessimism < 1.0 {
		return nil, fmt.Errorf(
			"verificationPessimism must be greater than 1.0, got %f", config.VerificationPessimism)
	}

	probe.SetStage("verify_commitment")
	commitmentBatch := []encoding.BlobCommitments{blobCommitments}
	err := verifier.VerifyCommitEquivalenceBatch(commitmentBatch)
	if err != nil {
		return nil, err
	}

	probe.SetStage("get_operator_state")
	operatorState, err := chainState.GetOperatorStateWithSocket(ctx, uint(referenceBlockNumber), []core.QuorumID{quorumID})
	if err != nil {
		return nil, err
	}
	operators, ok := operatorState.Operators[quorumID]
	if !ok {
		return nil, fmt.Errorf("no quorum with ID: %d", quorumID)
	}

	probe.SetStage("get_blob_versions")
	blobVersions, err := ethClient.GetAllVersionedBlobParams(ctx)
	if err != nil {
		return nil, err
	}

	blobParams, ok := blobVersions[blobVersion]
	if !ok {
		return nil, fmt.Errorf("invalid blob version %d", blobVersion)
	}

	probe.SetStage("get_encoding_params")
	encodingParams, err := v2.GetEncodingParams(blobCommitments.Length, blobParams)
	if err != nil {
		return nil, err
	}

	probe.SetStage("get_assignments")
	assignments, err := v2.GetAssignments(operatorState, blobParams, quorumID)
	if err != nil {
		return nil, errors.New("failed to get assignments")
	}

	// TODO verify technical accuracy of this
	minimumChunkCount := uint32(encodingParams.NumChunks) / blobParams.CodingRate

	downloadOrder := make([]core.OperatorID, 0, len(operators))
	for opID := range operators {
		downloadOrder = append(downloadOrder, opID)
	}

	downloadAndVerifyCtx, cancel := context.WithCancel(ctx)

	return &retrievalWorker{
		ctx:                      ctx,
		downloadAndVerifyContext: downloadAndVerifyCtx,
		cancel:                   cancel,
		logger:                   logger,
		config:                   config,
		connectionPool:           connectionPool,
		computePool:              computePool,
		operators:                operators,
		operatorState:            operatorState,
		encodingParams:           encodingParams,
		assignments:              assignments,
		quorumID:                 quorumID,
		blobKey:                  blobKey,
		verifier:                 verifier,
		blobCommitments:          blobCommitments,
		downloadStartedChan:      make(chan *downloadStarted, len(operators)),
		downloadCompletedChan:    make(chan *validatorDownloadCompleted, len(operators)),
		verificationCompleteChan: make(chan *verificationCompleted, len(operators)),
		decodeResponseChan:       make(chan *decodeResult, len(operators)),
		probe:                    probe,
		minimumChunkCount:        minimumChunkCount,
		completedDownloadSet:     make(map[core.OperatorID]struct{}),
		downloadsInProgress:      linkedlistqueue.New(),
		unverifiedChunks:         linkedlistqueue.New(),
		verifiedChunks:           linkedlistqueue.New(),
		downloadOrder:            downloadOrder,
	}, nil
}

// TODO document this workflow
// TODO unit test this workflow

// downloadBlobFromValidators downloads the blob from the validators.
func (w *retrievalWorker) downloadBlobFromValidators() ([]byte, error) {
	// Once everything is completed, we can abort all ongoing work.
	defer w.cancel()

	// To start things off, request an initial number of downloads.
	for w.nextDownloadIndex < len(w.downloadOrder) {
		if w.chunksBeingDownloaded > w.targetDownloadCount {
			// We've requested enough downloads
			break
		}
		operatorID := w.downloadOrder[w.nextDownloadIndex]
		assignedChunks := w.assignments[operatorID].NumChunks
		w.chunksBeingDownloaded += assignedChunks
		w.connectionPool.Submit(func() {
			w.downloadChunks(operatorID)
		})

		w.nextDownloadIndex++
	}

	// TODO encapsulate parts of this loop to make it easier to read
	w.probe.SetStage("download_and_verify")

	controlLoopTicker := time.NewTicker(w.config.ControlLoopPeriod)
	for w.verifiedChunkCount < w.minimumChunkCount {
		select {
		case <-w.ctx.Done():
			return nil, fmt.Errorf(
				"retrieval worker context cancelled, blobKey: %s: %w",
				w.blobKey.Hex(), w.ctx.Err())
		case message := <-w.downloadStartedChan:
			w.downloadsInProgress.Enqueue(message)
		case <-controlLoopTicker.C:
			// Check to see if any downloads in progress have exceeded the pessimistic timeout.
			for !w.downloadsInProgress.Empty() {
				next, _ := w.downloadsInProgress.Peek()
				operatorID := next.(*downloadStarted).operatorID
				downloadStart := next.(*downloadStarted).downloadStart

				if _, ok := w.completedDownloadSet[operatorID]; ok {
					// The operator has finished downloading, we can remove it from the queue.
					w.downloadsInProgress.Dequeue()
				} else if time.Since(downloadStart) > w.config.PessimisticTimeout {
					// Too much time has passed. Assume that the operator is not responding.
					w.downloadsInProgress.Dequeue()
					operatorChunks := w.assignments[operatorID].NumChunks
					w.chunksBeingDownloaded -= operatorChunks
				} else {
					// The next download has not yet timed out.
					break
				}
			}
		case message := <-w.downloadCompletedChan:
			if message.Err == nil {
				w.logger.Debug("downloaded chunks from operator",
					"operator", message.OperatorID.Hex(),
					"blobKey", w.blobKey.Hex())
				w.unverifiedChunks.Enqueue(message)
			} else {
				w.logger.Warn("failed to download chunk data",
					"operator", message.OperatorID.Hex(),
					"blobKey", w.blobKey.Hex(),
					"err", message.Err)

				chunkCount := w.assignments[message.OperatorID].NumChunks
				w.chunksBeingDownloaded -= chunkCount
			}
		case message := <-w.verificationCompleteChan:
			if message.Err == nil {
				w.logger.Debug("verified chunks from operator",
					"operator", message.OperatorID.Hex(),
					"blobKey", w.blobKey.Hex())
				w.verifiedChunks.Enqueue(message)
				w.verifiedChunkCount += uint32(len(message.chunks))
			} else {
				w.logger.Warn("failed to verify chunk data",
					"operator", message.OperatorID.Hex(),
					"blobKey", w.blobKey.Hex(),
					"err", message.Err)

				chunkCount := w.assignments[message.OperatorID].NumChunks
				w.chunksBeingVerified -= chunkCount
				w.chunksBeingDownloaded -= chunkCount
			}
		}

		// Check to see if we need to schedule more downloads.
		for w.nextDownloadIndex < len(w.downloadOrder) {
			if w.chunksBeingDownloaded > w.targetDownloadCount {
				// We've requested enough downloads
				break
			}
			operatorID := w.downloadOrder[w.nextDownloadIndex]
			assignedChunks := w.assignments[operatorID].NumChunks
			w.chunksBeingDownloaded += assignedChunks
			w.connectionPool.Submit(func() {
				w.downloadChunks(operatorID)
			})

			w.nextDownloadIndex++
		}

		// Check to see if we need to schedule more verifications.
		for !w.unverifiedChunks.Empty() {
			if w.chunksBeingVerified > w.targetVerifiedCount {
				// We've requested enough verifications
				break
			}

			next, _ := w.unverifiedChunks.Dequeue()
			reply := next.(*validatorDownloadCompleted).Reply
			operatorID := next.(*validatorDownloadCompleted).OperatorID

			w.chunksBeingVerified += uint32(len(reply.Chunks))

			w.computePool.Submit(func() {
				w.deserializeAndVerifyChunks(operatorID, reply)
			})
		}
	}

	// This aborts all unfinished download/verification work.
	w.cancel()

	if w.verifiedChunkCount < w.minimumChunkCount {
		return nil, fmt.Errorf("not enough chunks verified: %d < %d", w.verifiedChunkCount, w.minimumChunkCount)
	}

	// TODO can this last part be encapsulated in to a helper function?

	w.probe.SetStage("decode")
	w.logger.Info("decoding blob", "blobKey", w.blobKey.Hex())

	chunks := make([]*encoding.Frame, 0)
	indices := make([]uint, 0)

	for !w.verifiedChunks.Empty() {
		next, _ := w.verifiedChunks.Dequeue()
		operatorID := next.(*verificationCompleted).OperatorID
		operatorChunks := next.(*verificationCompleted).chunks

		assignment := w.assignments[operatorID]
		assignmentIndices := make([]uint, len(assignment.GetIndices()))

		for i, index := range assignment.GetIndices() {
			assignmentIndices[i] = uint(index)
		}

		chunks = append(chunks, operatorChunks...)
		indices = append(indices, assignmentIndices...)
	}

	// TODO: faster to truncate the number of chunks if we have more than we need?

	w.computePool.Submit(func() {
		w.decodeBlob(chunks, indices)
	})

	select {
	case <-w.ctx.Done():
		return nil, fmt.Errorf("retrieval worker context cancelled: %w", w.ctx.Err())
	case decodeResponse := <-w.decodeResponseChan:
		if decodeResponse.Err == nil {
			return decodeResponse.Blob, nil
		} else {
			return nil, fmt.Errorf("failed to decode blob: %w", decodeResponse.Err)
		}
	}

	// TODO abort all ongoing work when we terminate
}

// downloadChunks downloads the chunk data from the specified operator.
func (w *retrievalWorker) downloadChunks(operatorID core.OperatorID) {
	w.logger.Debug("downloading chunks",
		"operator", operatorID.Hex(),
		"blobKey", w.blobKey.Hex())

	ctx, cancel := context.WithTimeout(w.downloadAndVerifyContext, w.config.DownloadTimeout)
	defer cancel()

	// Report back to the control loop when the download started. This may be later than when
	// the download was scheduled if there is contention for the connection pool.
	w.downloadStartedChan <- &downloadStarted{
		operatorID:    operatorID,
		downloadStart: time.Now(),
	}

	// TODO we can get a tighter bound?
	maxBlobSize := 16 * units.MiB // maximum size of the original blob
	encodingRate := 8             // worst case scenario if one validator has 100% stake
	fudgeFactor := units.MiB      // to allow for some overhead from things like protobuf encoding
	maxMessageSize := maxBlobSize*encodingRate + fudgeFactor

	opInfo := w.operatorState.Operators[w.quorumID][operatorID]
	conn, err := grpc.NewClient(
		opInfo.Socket.GetV2RetrievalSocket(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMessageSize)),
	)
	defer func() {
		err := conn.Close()
		if err != nil {
			w.logger.Error("validator retriever failed to close connection", "err", err)
		}
	}()
	if err != nil {
		w.downloadCompletedChan <- &validatorDownloadCompleted{
			OperatorID: operatorID,
			Err:        fmt.Errorf("failed to create connection to operator %s: %w", operatorID.Hex(), err),
		}
		return
	}

	n := grpcnode.NewRetrievalClient(conn)
	request := &grpcnode.GetChunksRequest{
		BlobKey:  w.blobKey[:],
		QuorumId: uint32(w.quorumID),
	}

	reply, err := n.GetChunks(ctx, request)
	if err != nil {
		w.downloadCompletedChan <- &validatorDownloadCompleted{
			OperatorID: operatorID,
			Err:        fmt.Errorf("failed to get chunks from operator %s: %w", operatorID.Hex(), err),
		}
		return
	}

	w.downloadCompletedChan <- &validatorDownloadCompleted{
		OperatorID: operatorID,
		Reply:      reply,
	}
}

// deserializeAndVerifyChunks deserializes the chunks from the GetChunksReply and sends them to the chunksChan.
func (w *retrievalWorker) deserializeAndVerifyChunks(
	operatorID core.OperatorID,
	getChunksReply *grpcnode.GetChunksReply,
) {

	if w.downloadAndVerifyContext.Err() != nil {
		// blob is already finished
		return
	}

	w.logger.Debug("verifying chunks",
		"operator", operatorID.Hex(),
		"blobKey", w.blobKey.Hex())

	chunks := make([]*encoding.Frame, len(getChunksReply.GetChunks()))
	for i, data := range getChunksReply.GetChunks() {
		chunk, err := new(encoding.Frame).DeserializeGnark(data)
		if err != nil {
			w.verificationCompleteChan <- &verificationCompleted{
				OperatorID: operatorID,
				Err:        fmt.Errorf("failed to deserialize chunk from operator %s: %w", operatorID.Hex(), err),
			}
			return
		}

		chunks[i] = chunk
	}

	assignment := w.assignments[operatorID]

	assignmentIndices := make([]uint, len(assignment.GetIndices()))
	for i, index := range assignment.GetIndices() {
		assignmentIndices[i] = uint(index)
	}

	// TODO this can be parallelized!
	err := w.verifier.VerifyFrames(chunks, assignmentIndices, w.blobCommitments, w.encodingParams)
	if err != nil {
		w.verificationCompleteChan <- &verificationCompleted{
			OperatorID: operatorID,
			Err:        fmt.Errorf("failed to verify chunks from operator %s: %w", operatorID.Hex(), err),
		}
		return
	}

	w.verificationCompleteChan <- &verificationCompleted{
		OperatorID: operatorID,
		chunks:     chunks,
	}
}

// decodeBlob decodes the blob from the chunks and indices.
func (w *retrievalWorker) decodeBlob(chunks []*encoding.Frame, indices []uint) {

	w.logger.Debug("decoding blob", "blobKey", w.blobKey.Hex())

	blob, err := w.verifier.Decode(
		chunks,
		indices,
		w.encodingParams,
		uint64(w.blobCommitments.Length)*encoding.BYTES_PER_SYMBOL,
	)
	w.decodeResponseChan <- &decodeResult{
		Blob: blob,
		Err:  err,
	}
}
