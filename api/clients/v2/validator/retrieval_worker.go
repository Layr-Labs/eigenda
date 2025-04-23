package validator

import (
	"context"
	"errors"
	"fmt"
	"math"
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

/*
	                        ...........................................
	                        .    chunk downloading/verification       .
	                        .          happens in parallel   	      .
	                        .         (for different chunks)  	      .
	                        ...........................................
	                                 |                      |
	                                 v                      v
	+-----------------+     +-----------------+     +-----------------+     +-----------------+
	|     request     | --> | download chunks | --> |  verify chunks  | --> |  decode blob    |
	+-----------------+     +-----------------+     +-----------------+     +-----------------+

                          Each chunk goes through the following states:

               +-----------------------------------+
               |                                   |
               |                                   v
               |   available -> downloading -> downloaded -> verifying -> verified
               |                 |       |                       |
               |                 v       |                       |
               |   pessimistic timeout   |                       |
               |        |        |       |                       |
               |        |        |       v                       |
               +--------+        +---> failed <------------------+

== DownloadPessimism ==

The setting DownloadPessimism determines the number of chunks that should be downloaded. It wants to make sure
that there are at least a certain number of chunks with one of the following states:

  - downloading
  - downloaded
  - verifying
  - verified

If there are an insufficient number of chunks with one of these states, then it will schedule more chunks in
to be downloaded, if possible.

== VerificationPessimism ==

The setting VerificationPessimism determines the number of chunks that should be verified. It wants to make sure
that there are at least a certain number of chunks with one of the following states:

  - verifying
  - verified

If there are an insufficient number of chunks with one of these states, then it will cause chunks with the state
"downloaded" to be verified. If there are no remaining chunks with the state "downloaded", then it will wait for
chunks to enter this state.

*/

const (
	// chunks that can be downloaded
	available = iota
	// chunks that are currently being downloaded
	downloading
	// chunks that have been downloaded
	downloaded
	// chunks that are currently being verified
	verifying
	// chunks that have been verified
	verified
	// chunks that have failed to be downloaded or verified
	failed
	// chunks that have timed of the pessimistic download timeout, but still may be downloaded successfully
	pessimisticTimeout
)

// retrievalWorker implements the distributed retrieval process for a specified blob (i.e. reading the blobs from
// validators). It is responsible for coordinating the lifecycle of this retrieval workflow.
type retrievalWorker struct {
	ctx                      context.Context
	downloadAndVerifyContext context.Context
	cancel                   context.CancelFunc
	logger                   logging.Logger
	config                   *ValidatorClientConfig

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

	// If true, then the validator retrieval client will log detailed information about the download process.
	detailedLogging bool

	// Used to collect metrics.
	probe *common.SequenceProbe

	///////////////////////////////////////////////////////////////////////////////////////////////////
	// All variables below this line are for use only on the downloadBlobFromValidators() goroutine  //
	///////////////////////////////////////////////////////////////////////////////////////////////////

	// The order in which to download chunks from the operators.
	downloadOrder []core.OperatorID

	// The index of the next operator to download from.
	nextDownloadIndex int

	// The total number of chunks.
	totalChunkCount uint32

	// TODO verify technical accuracy of this
	// The minimum number of chunks needed to reconstruct the blob.
	minimumChunkCount uint32

	// The number of chunks we'd like to download.
	targetDownloadCount uint32

	// The number of chunks we'd like to verify.
	targetVerifiedCount uint32

	// This queue is used to determine when the pessimistic timeout for a download has been reached.
	downloadsInProgressQueue queues.Queue

	// Contains chunks that have been downloaded but not yet scheduled for verification.
	downloadedChunksQueue queues.Queue

	// Contains chunks that have been verified.
	verifiedChunksQueue queues.Queue

	// The status of the chunks from each validator. The key is the operator ID, and the value is the
	// status of the chunk.
	chunkStatusMap map[core.OperatorID]int

	// Counts the number of chunks in each status.
	chunkStatusCounts map[int]int
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

// newRetrievalWorker creates a new retrieval worker.
func newRetrievalWorker(
	ctx context.Context,
	logger logging.Logger,
	config *ValidatorClientConfig,
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
	totalChunkCount := uint32(encodingParams.NumChunks)
	availableChunkCount := uint32(encodingParams.NumChunks)
	minimumChunkCount := availableChunkCount / blobParams.CodingRate

	downloadOrder := make([]core.OperatorID, 0, len(operators))
	for opID := range operators {
		downloadOrder = append(downloadOrder, opID)
	}

	downloadAndVerifyCtx, cancel := context.WithCancel(ctx)

	chunkStatusMap := make(map[core.OperatorID]int)
	chunkStatusCounts := make(map[int]int)
	for _, opID := range downloadOrder {
		chunkStatusMap[opID] = available
		chunkStatusCounts[available] += int(assignments[opID].NumChunks)
	}

	targetDownloadCount := uint32(math.Ceil(float64(minimumChunkCount) * config.DownloadPessimism))
	targetVerifiedCount := uint32(math.Ceil(float64(minimumChunkCount) * config.VerificationPessimism))

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
		downloadsInProgressQueue: linkedlistqueue.New(),
		downloadedChunksQueue:    linkedlistqueue.New(),
		verifiedChunksQueue:      linkedlistqueue.New(),
		downloadOrder:            downloadOrder,
		totalChunkCount:          totalChunkCount,
		chunkStatusMap:           chunkStatusMap,
		chunkStatusCounts:        chunkStatusCounts,
		targetDownloadCount:      targetDownloadCount,
		targetVerifiedCount:      targetVerifiedCount,
		detailedLogging:          config.DetailedLogging,
	}, nil
}

// TODO unit test this workflow

// updateChunkStatus updates the status of a chunk from a given operator. It also updates the
// counts of the various chunk statuses.
func (w *retrievalWorker) updateChunkStatus(operatorID core.OperatorID, status int) {
	operatorChunks := w.assignments[operatorID].NumChunks

	oldStatus, ok := w.chunkStatusMap[operatorID]
	if ok {
		w.chunkStatusCounts[oldStatus] -= int(operatorChunks)
	}
	w.chunkStatusMap[operatorID] = status
	w.chunkStatusCounts[status] += int(operatorChunks)
}

// getChunkStatus returns the status of a chunk from a given operator.
func (w *retrievalWorker) getChunkStatus(operatorID core.OperatorID) int {
	return w.chunkStatusMap[operatorID]
}

// getStatusCount returns the number of chunks with one of the given statuses.
func (w *retrievalWorker) getStatusCount(statuses ...int) uint32 {
	total := 0
	for _, status := range statuses {
		if count, ok := w.chunkStatusCounts[status]; ok {
			total += count
		}
	}
	return uint32(total)
}

// downloadBlobFromValidators downloads the blob from the validators.
func (w *retrievalWorker) downloadBlobFromValidators() ([]byte, error) {
	// Once everything is completed, we can abort all ongoing work.
	defer w.cancel()

	w.probe.SetStage("download_and_verify")

	controlLoopTicker := time.NewTicker(w.config.ControlLoopPeriod)
	for {
		if w.getStatusCount(verified) >= w.minimumChunkCount {
			// We've verified enough chunks to reconstruct the blob
			break
		}
		if w.getStatusCount(failed) >= w.totalChunkCount-w.minimumChunkCount {
			// We've failed too many chunks, reconstruction is no longer possible
			break
		}

		w.scheduleDownloads()
		w.scheduleVerifications()

		select {
		case <-w.ctx.Done():
			return nil, fmt.Errorf(
				"retrieval worker context cancelled, blobKey: %s: %w",
				w.blobKey.Hex(), w.ctx.Err())
		case message := <-w.downloadStartedChan:
			w.downloadsInProgressQueue.Enqueue(message)
		case <-controlLoopTicker.C:
			w.checkPessimisticTimeout()
		case message := <-w.downloadCompletedChan:
			w.handleCompletedDownload(message)
		case message := <-w.verificationCompleteChan:
			w.handleVerificationCompleted(message)
		}
	}

	// This aborts all unfinished download/verification work.
	w.cancel()

	verifiedCount := w.getStatusCount(verified)
	if verifiedCount < w.minimumChunkCount {
		return nil, fmt.Errorf("not enough chunks verified: %d < %d", verifiedCount, w.minimumChunkCount)
	}

	w.probe.SetStage("decode")
	return w.decodeChunks()
}

// checkPessimisticTimeout checks to see if any downloads in progress have exceeded the pessimistic timeout.
// These downloads are not cancelled, but this timeout may result in other chunks being downloaded.
func (w *retrievalWorker) checkPessimisticTimeout() {
	for !w.downloadsInProgressQueue.Empty() {
		next, _ := w.downloadsInProgressQueue.Peek()
		operatorID := next.(*downloadStarted).operatorID
		downloadStart := next.(*downloadStarted).downloadStart

		if w.getChunkStatus(operatorID) != downloading {
			// The operator has finished downloading, we can remove it from the queue.
			w.downloadsInProgressQueue.Dequeue()
		} else if time.Since(downloadStart) > w.config.PessimisticTimeout {
			// Too much time has passed. Assume that the operator is not responding.
			if w.detailedLogging {
				w.logger.Debug("soft timeout exceeded for chunk download",
					"operator", operatorID.Hex())
			}
			w.downloadsInProgressQueue.Dequeue()

			w.updateChunkStatus(operatorID, pessimisticTimeout)
		} else {
			// The next download has not yet timed out.
			break
		}
	}
}

// handleCompletedDownload handles the completion of a download.
func (w *retrievalWorker) handleCompletedDownload(message *validatorDownloadCompleted) {
	if message.Err == nil {
		if w.detailedLogging {
			w.logger.Debug("downloaded chunks from operator",
				"operator", message.OperatorID.Hex(),
				"blobKey", w.blobKey.Hex())
		}
		w.downloadedChunksQueue.Enqueue(message)
		w.updateChunkStatus(message.OperatorID, downloaded)
	} else {
		w.logger.Warn("failed to download chunk data",
			"operator", message.OperatorID.Hex(),
			"blobKey", w.blobKey.Hex(),
			"err", message.Err)

		w.updateChunkStatus(message.OperatorID, failed)
	}
}

// handleVerificationCompleted handles the completion of a verification.
func (w *retrievalWorker) handleVerificationCompleted(message *verificationCompleted) {
	if message.Err == nil {
		if w.detailedLogging {
			w.logger.Debug("verified chunks from operator",
				"operator", message.OperatorID.Hex(),
				"blobKey", w.blobKey.Hex())
		}
		w.verifiedChunksQueue.Enqueue(message)
		w.updateChunkStatus(message.OperatorID, verified)
	} else {
		w.logger.Warn("failed to verify chunk data",
			"operator", message.OperatorID.Hex(),
			"blobKey", w.blobKey.Hex(),
			"err", message.Err)

		w.updateChunkStatus(message.OperatorID, failed)
	}
}

// scheduleDownloads schedules downloads as needed.
func (w *retrievalWorker) scheduleDownloads() {
	for w.nextDownloadIndex < len(w.downloadOrder) {
		if w.getStatusCount(downloading, downloaded, verifying, verified) >= w.targetDownloadCount {
			// We've requested enough downloads
			break
		}
		operatorID := w.downloadOrder[w.nextDownloadIndex]
		w.updateChunkStatus(operatorID, downloading)

		w.connectionPool.Submit(func() {
			w.downloadChunks(operatorID)
		})

		w.nextDownloadIndex++
	}
}

// scheduleVerifications schedules verifications as needed.
func (w *retrievalWorker) scheduleVerifications() {
	for !w.downloadedChunksQueue.Empty() {
		if w.getStatusCount(verifying, verified) > w.targetVerifiedCount {
			// We've requested enough verifications
			break
		}

		next, _ := w.downloadedChunksQueue.Dequeue()
		reply := next.(*validatorDownloadCompleted).Reply
		operatorID := next.(*validatorDownloadCompleted).OperatorID

		w.updateChunkStatus(operatorID, verifying)

		w.computePool.Submit(func() {
			w.deserializeAndVerifyChunks(operatorID, reply)
		})
	}
}

// decodeBlob decodes the blob from the chunks.
func (w *retrievalWorker) decodeChunks() ([]byte, error) {
	w.logger.Info("decoding blob", "blobKey", w.blobKey.Hex())

	chunks := make([]*encoding.Frame, 0)
	indices := make([]uint, 0)

	for !w.verifiedChunksQueue.Empty() {
		next, _ := w.verifiedChunksQueue.Dequeue()
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
}

// downloadChunks downloads the chunk data from the specified operator.
func (w *retrievalWorker) downloadChunks(operatorID core.OperatorID) {
	if w.detailedLogging {
		w.logger.Debug("downloading chunks",
			"operator", operatorID.Hex(),
			"blobKey", w.blobKey.Hex())
	}

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

	if w.detailedLogging {
		w.logger.Debug("verifying chunks",
			"operator", operatorID.Hex(),
			"blobKey", w.blobKey.Hex())
	}

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

	if w.detailedLogging {
		w.logger.Debug("decoding blob", "blobKey", w.blobKey.Hex())
	}

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
