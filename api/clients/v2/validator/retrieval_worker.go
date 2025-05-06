package validator

import (
	"context"
	"fmt"
	"math"
	"time"

	grpcnode "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/emirpasic/gods/queues"
	"github.com/emirpasic/gods/queues/linkedlistqueue"
	"github.com/gammazero/workerpool"
)

/*
	                        ...........................................
	                        .    chunk downloading/verification       .
	                        .          happens in parallel            .
	                        .         (for different chunks)          .
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

If there are an insufficient number of chunks with one of these states, then it will schedule more chunks
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

// chunkStatus is the status of a chunk download/verification from a particular validator.
type chunkStatus int

const (
	// chunks that can be downloaded
	available chunkStatus = iota
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

// TODO(cody.littley): the following improvements can be made in the future
//  - the ability to download a blob from multiple quorums at once (useful if many validators are offline)
//  - check to see if it's faster to send the bare minimum number of chunks to decoding, modify code accordingly
//  - more granular metrics via sequence probe (requires sequence probe enhancements)

// retrievalWorker implements the distributed retrieval process for a specified blob (i.e. reading the blobs from
// validators). It is responsible for coordinating the lifecycle of this retrieval workflow.
type retrievalWorker struct {
	ctx                  context.Context
	downloadAndVerifyCtx context.Context
	cancel               context.CancelFunc
	logger               logging.Logger
	config               *ValidatorClientConfig

	// Responsible for talking to the validator nodes via gRPCs.
	validatorGRPCManager ValidatorGRPCManager

	// Responsible for deserializing and verifying chunk data.
	chunkDeserializer ChunkDeserializer

	// A pool of workers for network intensive operations (e.g. downloading blob data).
	connectionPool *workerpool.WorkerPool

	// A pool of workers for CPU intensive operations (e.g. deserializing and verifying blob data).
	computePool *workerpool.WorkerPool

	// The encoding parameters for the blob.
	encodingParams *encoding.EncodingParams

	// The assignments for the operators, i.e. which operators are responsible for which chunks.
	assignments map[core.OperatorID]v2.Assignment

	// The quorum to download from.
	quorumID core.QuorumID

	// The blob key to download.
	blobKey v2.BlobKey

	// Verifies chunk data.
	verifier encoding.Verifier

	// The commitments for the blob.
	blobCommitments *encoding.BlobCommitments

	// When a thread begins downloading chunk data, it will send a message to the downloadStartedChan.
	downloadStartedChan chan *downloadStarted

	// When a thread completes downloading chunk data, it will send a message to the downloadCompletedChan.
	downloadCompletedChan chan *downloadCompleted

	// When a thread completes verifying chunk data, it will send a message to the verificationCompletedChan.
	verificationCompletedChan chan *verificationCompleted

	// When a thread completes decoding chunk data, it will send a message to the decodeResponseChan.
	decodeResponseChan chan *decodeCompleted

	// The function used to decode the blob from the chunks. This can be set to a custom function
	// for testing purposes.
	decodeBlobFunction DecodeBlobFunction

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
	// status of the chunks.
	chunkStatusMap map[core.OperatorID]chunkStatus

	// Counts the number of chunks in each status.
	chunkStatusCounts map[chunkStatus]int
}

// downloadStarted is used to signal that a download of chunk data has been initiated.
type downloadStarted struct {
	operatorID    core.OperatorID
	downloadStart time.Time
}

// downloadCompleted is used to signal that a download of chunk data has completed.
type downloadCompleted struct {
	OperatorID core.OperatorID
	Reply      *grpcnode.GetChunksReply
	Err        error
}

// verificationCompleted is used to signal that verification of chunk data has completed.
type verificationCompleted struct {
	OperatorID core.OperatorID
	chunks     []*encoding.Frame
	Err        error
}

// decodeCompleted is used to signal that decoding of chunk data has completed.
type decodeCompleted struct {
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
	validatorGRPCManager ValidatorGRPCManager,
	chunkDeserializer ChunkDeserializer,
	assignments map[core.OperatorID]v2.Assignment,
	totalChunkCount uint32,
	minimumChunkCount uint32,
	encodingParams *encoding.EncodingParams,
	quorumID core.QuorumID,
	blobKey v2.BlobKey,
	verifier encoding.Verifier,
	blobCommitments *encoding.BlobCommitments,
	probe *common.SequenceProbe,
) (*retrievalWorker, error) {
	if config.DownloadPessimism < 1.0 {
		return nil, fmt.Errorf("downloadPessimism must be greater than or equal to 1.0, got %f",
			config.DownloadPessimism)
	}
	if config.VerificationPessimism < 1.0 {
		return nil, fmt.Errorf(
			"verificationPessimism must be greater than or equal to 1.0, got %f", config.VerificationPessimism)
	}
	if minimumChunkCount == 0 {
		return nil, fmt.Errorf("minimumChunkCount must be greater than 0")
	}

	downloadOrder := make([]core.OperatorID, 0, len(assignments))
	for opID := range assignments {
		downloadOrder = append(downloadOrder, opID)
	}

	// The retrieval worker uses two contexts. The downloadAndVerifyCtx is cancelled once a sufficient number
	// of chunks have been downloaded and verified. This causes all ongoing downloads and verifications to be
	// aborted (if possible), since they are not needed. There are other operations that require a context after
	// the downloadAndVerifyCtx is cancelled, so we need to keep a reference the original context as well.
	downloadAndVerifyCtx, cancel := context.WithCancel(ctx)

	chunkStatusMap := make(map[core.OperatorID]chunkStatus)
	chunkStatusCounts := make(map[chunkStatus]int)
	for _, opID := range downloadOrder {
		chunkStatusMap[opID] = available
		chunkStatusCounts[available] += int(assignments[opID].NumChunks)
	}

	targetDownloadCount := uint32(math.Ceil(float64(minimumChunkCount) * config.DownloadPessimism))
	targetVerifiedCount := uint32(math.Ceil(float64(minimumChunkCount) * config.VerificationPessimism))

	worker := &retrievalWorker{
		ctx:                       ctx,
		downloadAndVerifyCtx:      downloadAndVerifyCtx,
		cancel:                    cancel,
		logger:                    logger,
		config:                    config,
		connectionPool:            connectionPool,
		computePool:               computePool,
		encodingParams:            encodingParams,
		assignments:               assignments,
		quorumID:                  quorumID,
		blobKey:                   blobKey,
		verifier:                  verifier,
		blobCommitments:           blobCommitments,
		downloadStartedChan:       make(chan *downloadStarted, len(assignments)),
		downloadCompletedChan:     make(chan *downloadCompleted, len(assignments)),
		verificationCompletedChan: make(chan *verificationCompleted, len(assignments)),
		decodeResponseChan:        make(chan *decodeCompleted, 1),
		probe:                     probe,
		minimumChunkCount:         minimumChunkCount,
		downloadsInProgressQueue:  linkedlistqueue.New(),
		downloadedChunksQueue:     linkedlistqueue.New(),
		verifiedChunksQueue:       linkedlistqueue.New(),
		downloadOrder:             downloadOrder,
		totalChunkCount:           totalChunkCount,
		chunkStatusMap:            chunkStatusMap,
		chunkStatusCounts:         chunkStatusCounts,
		targetDownloadCount:       targetDownloadCount,
		targetVerifiedCount:       targetVerifiedCount,
		validatorGRPCManager:      validatorGRPCManager,
		chunkDeserializer:         chunkDeserializer,
		decodeBlobFunction:        config.UnsafeDecodeBlobFunction,
	}

	if worker.decodeBlobFunction == nil {
		worker.decodeBlobFunction = worker.standardDecodeBlob
	}

	return worker, nil
}

// updateChunkStatus updates the status of a chunk from a given operator. It also updates the
// counts of the various chunk statuses.
func (w *retrievalWorker) updateChunkStatus(operatorID core.OperatorID, status chunkStatus) {
	operatorChunks := w.assignments[operatorID].NumChunks

	oldStatus, ok := w.chunkStatusMap[operatorID]
	if ok {
		w.chunkStatusCounts[oldStatus] -= int(operatorChunks)
	}
	w.chunkStatusMap[operatorID] = status
	w.chunkStatusCounts[status] += int(operatorChunks)
}

// getChunkStatus returns the status of a chunk from a given operator.
func (w *retrievalWorker) getChunkStatus(operatorID core.OperatorID) chunkStatus {
	return w.chunkStatusMap[operatorID]
}

// getStatusCount returns the number of chunks with one of the given statuses.
func (w *retrievalWorker) getStatusCount(statuses ...chunkStatus) uint32 {
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
	// Defer a cancellation just in case we return early. There are no negative side effects if the context
	// is cancelled more than once.
	defer w.cancel()

	w.probe.SetStage("download_and_verify")

	controlLoopTicker := time.NewTicker(w.config.ControlLoopPeriod)
	for {
		if w.getStatusCount(verified) >= w.minimumChunkCount {
			// We've verified enough chunks to reconstruct the blob
			break
		}
		if w.getStatusCount(failed) > w.totalChunkCount-w.minimumChunkCount {
			// We've failed too many chunks, reconstruction is no longer possible
			break
		}

		w.scheduleDownloads()
		w.scheduleVerifications()

		select {
		case <-w.downloadAndVerifyCtx.Done():
			return nil, fmt.Errorf(
				"retrieval worker context cancelled, blobKey: %s: %w",
				w.blobKey.Hex(), w.downloadAndVerifyCtx.Err())
		case message := <-w.downloadStartedChan:
			w.downloadsInProgressQueue.Enqueue(message)
		case <-controlLoopTicker.C:
			w.checkPessimisticTimeout()
		case message := <-w.downloadCompletedChan:
			w.handleCompletedDownload(message)
		case message := <-w.verificationCompletedChan:
			w.handleVerificationCompleted(message)
		}
	}

	// This aborts all unfinished download/verification work.
	w.cancel()

	verifiedCount := w.getStatusCount(verified)
	if verifiedCount < w.minimumChunkCount {
		return nil, fmt.Errorf("not enough chunks verified: %d < %d", verifiedCount, w.minimumChunkCount)
	}
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
			if w.config.DetailedLogging {
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
func (w *retrievalWorker) handleCompletedDownload(message *downloadCompleted) {
	if message.Err == nil {
		if w.config.DetailedLogging {
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
		if w.config.DetailedLogging {
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
		if w.getStatusCount(verifying, verified) >= w.targetVerifiedCount {
			// We've requested enough verifications
			break
		}

		next, _ := w.downloadedChunksQueue.Dequeue()
		reply := next.(*downloadCompleted).Reply
		operatorID := next.(*downloadCompleted).OperatorID

		w.updateChunkStatus(operatorID, verifying)

		w.computePool.Submit(func() {
			w.deserializeAndVerifyChunks(w.ctx, operatorID, reply)
		})
	}
}

// decodeBlob decodes the blob from the chunks.
func (w *retrievalWorker) decodeChunks() ([]byte, error) {
	w.probe.SetStage("decode")

	if w.config.DetailedLogging {
		w.logger.Info("decoding blob", "blobKey", w.blobKey.Hex())
	}

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
	if w.config.DetailedLogging {
		w.logger.Debug("downloading chunks",
			"operator", operatorID.Hex(),
			"blobKey", w.blobKey.Hex())
	}

	// Report back to the control loop when the download started. This may be later than when
	// the download was scheduled if there is contention for the connection pool.
	w.downloadStartedChan <- &downloadStarted{
		operatorID:    operatorID,
		downloadStart: w.config.TimeSource(),
	}

	ctx, cancel := context.WithTimeout(w.downloadAndVerifyCtx, w.config.DownloadTimeout)
	defer cancel()

	reply, err := w.validatorGRPCManager.DownloadChunks(ctx, w.blobKey, operatorID, w.quorumID)

	w.downloadCompletedChan <- &downloadCompleted{
		OperatorID: operatorID,
		Reply:      reply,
		Err:        err,
	}
}

// deserializeAndVerifyChunks deserializes the chunks from the GetChunksReply and sends them to the chunksChan.
func (w *retrievalWorker) deserializeAndVerifyChunks(
	_ context.Context,
	operatorID core.OperatorID,
	getChunksReply *grpcnode.GetChunksReply,
) {

	if w.downloadAndVerifyCtx.Err() != nil {
		// blob is already finished
		return
	}

	if w.config.DetailedLogging {
		w.logger.Debug("verifying chunks",
			"operator", operatorID.Hex(),
			"blobKey", w.blobKey.Hex())
	}

	chunks, err := w.chunkDeserializer.DeserializeAndVerify(
		w.blobKey,
		operatorID,
		getChunksReply,
		w.blobCommitments,
		w.encodingParams)

	w.verificationCompletedChan <- &verificationCompleted{
		OperatorID: operatorID,
		chunks:     chunks,
		Err:        err,
	}
}

// decodeBlob decodes the blob from the chunks and indices.
func (w *retrievalWorker) decodeBlob(chunks []*encoding.Frame, indices []uint) {
	if w.config.DetailedLogging {
		w.logger.Debug("decoding blob", "blobKey", w.blobKey.Hex())
	}

	// Unless this has been overridden for testing, this is just a call to standardDecodeBlob().
	blob, err := w.decodeBlobFunction(w.ctx, w.blobKey, chunks, indices)

	w.decodeResponseChan <- &decodeCompleted{
		Blob: blob,
		Err:  err,
	}
}

// decodeBlobFunction is the default function used to decode the blob from the chunks.
func (w *retrievalWorker) standardDecodeBlob(
	_ context.Context,
	_ v2.BlobKey,
	chunks []*encoding.Frame,
	indices []uint,
) ([]byte, error) {

	return w.verifier.Decode(
		chunks,
		indices,
		*w.encodingParams,
		uint64(w.blobCommitments.Length)*encoding.BYTES_PER_SYMBOL)
}
