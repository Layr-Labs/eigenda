package signingrate

import (
	"context"
	"errors"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
)

var _ SigningRateTracker = (*threadsafeSigningRateTracker)(nil)

// Size of the channel buffer for requests to the internal goroutine. This is intentionally sized large in order
// to absorb bursts of updates. In practice, this won't be necessary unless the number of validators is very large
// or the batch rate is very high.
const channelSize = 4096

// A thread-safe wrapper around a SigningRateTracker. Although we could implement this using a mutex,
// we instead use a goroutine and a channel to serialize access to the underlying SigningRateTracker.
// This allows operations such as ReportSuccess and ReportFailure to be non-blocking, which is important
// for performance. These methods are called many times for each batch processed, and we don't want
// to block the main processing loop on mutex contention.
type threadsafeSigningRateTracker struct {
	ctx    context.Context
	cancel context.CancelFunc

	// The base signing rate tracker that does the actual work.
	base SigningRateTracker

	// Channel for sending requests to the internal goroutine.
	requests chan any
}

// a request to invoke GetValidatorSigningRate
type getValidatorSigningRateRequest struct {
	operatorID []byte
	startTime  time.Time
	endTime    time.Time
	responseCh chan *getValidatorSigningRateResponse
}

// holds a response to GetValidatorSigningRate
type getValidatorSigningRateResponse struct {
	result *validator.ValidatorSigningRate
	err    error
}

// a request to invoke GetSigningRateDump
type getSigningRateDumpRequest struct {
	startTime  time.Time
	responseCh chan *getSigningRateDumpResponse
}

// holds a response to GetSigningRateDump
type getSigningRateDumpResponse struct {
	result []*validator.SigningRateBucket
	err    error
}

// a request to invoke GetUnflushedBuckets
type getUnflushedBucketsRequest struct {
	responseCh chan *getUnflushedBucketsResponse
}

// holds a response to GetUnflushedBuckets
type getUnflushedBucketsResponse struct {
	result []*validator.SigningRateBucket
	err    error
}

// a request to invoke UpdateLastBucket
type updateLastBucketRequest struct {
	bucket *validator.SigningRateBucket
	now    time.Time
}

// a request to invoke ReportSuccess
type reportSuccessRequest struct {
	now            time.Time
	id             core.OperatorID
	batchSize      uint64
	signingLatency time.Duration
}

// a request to invoke ReportFailure
type reportFailureRequest struct {
	now       time.Time
	id        core.OperatorID
	batchSize uint64
}

func NewThreadsafeSigningRateTracker(base SigningRateTracker) SigningRateTracker {
	ctx, cancel := context.WithCancel(context.Background())

	tracker := &threadsafeSigningRateTracker{
		ctx:      ctx,
		cancel:   cancel,
		base:     base,
		requests: make(chan any, channelSize),
	}

	go tracker.controlLoop()

	return tracker
}

func (t *threadsafeSigningRateTracker) GetValidatorSigningRate(
	operatorID []byte,
	startTime time.Time,
	endTime time.Time,
) (*validator.ValidatorSigningRate, error) {

	request := &getValidatorSigningRateRequest{
		operatorID: operatorID,
		startTime:  startTime,
		endTime:    endTime,
		responseCh: make(chan *getValidatorSigningRateResponse, 1),
	}

	// Send the request
	select {
	case <-t.ctx.Done():
		return nil, errors.New("signing rate tracker is shutting down")
	case t.requests <- request:
	}

	// await the response
	select {
	case <-t.ctx.Done():
		return nil, errors.New("signing rate tracker is shutting down")
	case response := <-request.responseCh:
		return response.result, response.err
	}
}

func (t *threadsafeSigningRateTracker) GetSigningRateDump(
	startTime time.Time,
) ([]*validator.SigningRateBucket, error) {

	request := &getSigningRateDumpRequest{
		startTime:  startTime,
		responseCh: make(chan *getSigningRateDumpResponse, 1),
	}

	// Send the request
	select {
	case <-t.ctx.Done():
		return nil, errors.New("signing rate tracker is shutting down")
	case t.requests <- request:
	}

	// await the response
	select {
	case <-t.ctx.Done():
		return nil, errors.New("signing rate tracker is shutting down")
	case response := <-request.responseCh:
		return response.result, response.err
	}
}

func (t *threadsafeSigningRateTracker) GetUnflushedBuckets() ([]*validator.SigningRateBucket, error) {

	request := &getUnflushedBucketsRequest{
		responseCh: make(chan *getUnflushedBucketsResponse, 1),
	}

	// Send the request
	select {
	case <-t.ctx.Done():
		return nil, errors.New("signing rate tracker is shutting down")
	case t.requests <- request:
	}

	// await the response
	select {
	case <-t.ctx.Done():
		return nil, errors.New("signing rate tracker is shutting down")
	case response := <-request.responseCh:
		return response.result, nil
	}
}

func (t *threadsafeSigningRateTracker) ReportSuccess(
	now time.Time,
	id core.OperatorID,
	batchSize uint64,
	signingLatency time.Duration,
) {

	request := &reportSuccessRequest{
		now:            now,
		id:             id,
		batchSize:      batchSize,
		signingLatency: signingLatency,
	}

	select {
	case <-t.ctx.Done():
		// things are being torn down, just drop the request
	case t.requests <- request:
	}

}

func (t *threadsafeSigningRateTracker) ReportFailure(
	now time.Time,
	id core.OperatorID,
	batchSize uint64,
) {
	request := &reportFailureRequest{
		now:       now,
		id:        id,
		batchSize: batchSize,
	}

	select {
	case <-t.ctx.Done():
		// things are being torn down, just drop the request
	case t.requests <- request:
	}
}

func (t *threadsafeSigningRateTracker) UpdateLastBucket(now time.Time, bucket *validator.SigningRateBucket) {
	request := &updateLastBucketRequest{
		bucket: bucket,
		now:    now,
	}

	select {
	case <-t.ctx.Done():
		// things are being torn down, just drop the request
	case t.requests <- request:
	}
}

func (t *threadsafeSigningRateTracker) Close() {
	t.cancel()
}

// Serialize access to the underlying SigningRateTracker.
func (t *threadsafeSigningRateTracker) controlLoop() {
	for {
		select {
		case <-t.ctx.Done():
			return
		case req := <-t.requests:
			switch typedRequest := req.(type) {

			case *getValidatorSigningRateRequest:
				result, err := t.base.GetValidatorSigningRate(
					typedRequest.operatorID,
					typedRequest.startTime,
					typedRequest.endTime)
				response := &getValidatorSigningRateResponse{
					result: result,
					err:    err,
				}
				select {
				case <-t.ctx.Done():
					return
				case typedRequest.responseCh <- response:
				}

			case *getSigningRateDumpRequest:
				result, err := t.base.GetSigningRateDump(typedRequest.startTime)
				response := &getSigningRateDumpResponse{
					result: result,
					err:    err,
				}
				select {
				case <-t.ctx.Done():
					return
				case typedRequest.responseCh <- response:
				}

			case *updateLastBucketRequest:
				t.base.UpdateLastBucket(typedRequest.now, typedRequest.bucket)

			case *getUnflushedBucketsRequest:
				result, err := t.base.GetUnflushedBuckets()
				response := &getUnflushedBucketsResponse{
					result: result,
					err:    err,
				}
				select {
				case <-t.ctx.Done():
					return
				case typedRequest.responseCh <- response:
				}

			case *reportSuccessRequest:
				t.base.ReportSuccess(
					typedRequest.now,
					typedRequest.id,
					typedRequest.batchSize,
					typedRequest.signingLatency)

			case *reportFailureRequest:
				t.base.ReportFailure(typedRequest.now, typedRequest.id, typedRequest.batchSize)
			}
		}
	}
}
