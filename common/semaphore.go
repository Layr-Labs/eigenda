package common

import (
	"context"
	"fmt"
	"sync/atomic"
)

// semaphoreChannelSize is the size of the channels used to communicate with the semaphore control loop.
const semaphoreChannelSize = 64

// Implements a semaphore. It baffles me why this is not in the standard libraries.
//
// This semaphore is "fair". That is, if N sequential calls to Acquire() are made, then each call will acquire their
// tokens in the order that they were requested. If call N requests many tokens and is temporarily blocked, and if
// call N+1 requests few tokens and those tokens are technically available, then call N+1 will not be able to acquire
// those tokens until call N has acquired its tokens (or timed out).
type Semaphore struct {
	// The total number of tokens available in the semaphore.
	totalTokens uint64

	// The number of tokens currently taken by calls to Acquire().
	acquiredTokens uint64

	// Requests to acquire tokens are sent to this channel.
	acquireChan chan *acquireReqeust

	// Requests to release tokens are sent to this channel.
	releaseChan chan *releaseRequest

	// This channel is closed when the semaphore is closed.
	shutdownChan chan struct{}

	// Used to make calling Close() idempotent.
	closeCalled atomic.Bool

	// The next acquireRequest that will be processed. May be nil. Only accessed from the control loop goroutine.
	nextAcquireRequest *acquireReqeust
}

// NewSemaphore creates a new semaphore with the specified number of available tokens.
func NewSemaphore(tokens uint64) (*Semaphore, error) {
	if tokens == 0 {
		return nil, fmt.Errorf("tokens must be greater than zero")
	}

	sem := &Semaphore{
		totalTokens:  tokens,
		acquireChan:  make(chan *acquireReqeust, semaphoreChannelSize),
		releaseChan:  make(chan *releaseRequest, semaphoreChannelSize),
		shutdownChan: make(chan struct{}),
	}
	go sem.controlLoop()

	return sem, nil
}

// A message sent to the control loop that requests a certain number of tokens to be acquired.
type acquireReqeust struct {
	// The context used by the caller to await for the acquire request to complete. If this context is canceled,
	// then we can abort without allocating any tokens.
	ctx context.Context

	// The number of tokens requested.
	tokens uint64

	// this channel produces a value when either the tokens are acquired (true), or the request failed (false).
	acquiredChan chan bool
}

// A message sent to the control loop that requests a certain number of tokens to be released.
type releaseRequest struct {
	// The number of tokens to release.
	tokens uint64

	// this channel produces a value when the tokens are released (true), or the request failed (false). The request
	// will only fail if more tokens are released than were acquired.
	releaseChan chan bool
}

// Acquire acquires the specified number of tokens from the semaphore, blocking until they are available.
// If the context is canceled, it returns immediately with an error. If this method does not return an error,
// Release MUST be called with the same number of tokens to release them back to the semaphore. If this method does
// return an error, Release MUST NOT be called, and the tokens should be treated as if they were never acquired.
func (s *Semaphore) Acquire(ctx context.Context, tokens uint64) error {

	if tokens > s.totalTokens {
		return fmt.Errorf("cannot acquire %d tokens, only %d available", tokens, s.totalTokens)
	}

	// Send the request to the control loop.

	request := &acquireReqeust{
		ctx:          ctx,
		tokens:       tokens,
		acquiredChan: make(chan bool, 1),
	}

	select {
	case <-s.shutdownChan:
		// Close() has been called.
		return fmt.Errorf("semaphore closed")
	case <-ctx.Done():
		// Context was canceled before we could send the request.
		// The tokens were never acquired, no need to release them.
		return fmt.Errorf("context canceled")
	case s.acquireChan <- request:
		// Request was successfully sent to the control loop.
	}

	// Await a response from the control loop.

	select {
	case <-s.shutdownChan:
		// Close() has been called. It doesn't matter if we release the tokens or not.
		return fmt.Errorf("semaphore closed")
	case <-ctx.Done():
		// The context was canceled after we sent the request, but before we got a response.
		// Depending on the response, we may need to release these tokens.
		go func() {
			select {
			case <-s.shutdownChan:
				// Close() has been called, no need to release the tokens.
				return
			case success := <-request.acquiredChan:
				if success {
					// We acquired the tokens, but the caller will not release them.
					_ = s.Release(tokens)
				}
			}
		}()
		return fmt.Errorf("context canceled")
	case success := <-request.acquiredChan:
		// We got a response from the control loop.
		if success {
			// Tokens were successfully acquired.
			return nil
		}
		return fmt.Errorf("failed to acquire %d tokens", tokens)
	}
}

// Release releases the specified number of tokens back to the semaphore. Must be called with the same number of
// tokens that were acquired. It is legal to release tokens in smaller amounts than were acquired. If more tokens
// are released than were acquired, an error will be returned. Do not call this method if the Acquire() call
// returned an error, as the tokens were never acquired in that case.
func (s *Semaphore) Release(tokens uint64) error {
	if tokens > s.totalTokens {
		return fmt.Errorf("cannot release %d tokens, only %d available", tokens, s.totalTokens)
	}

	// Unlike Acquire(), this method does not accept a context from the caller. Although this method may block
	// for a very short amount of time, the control loop will immediately process release requests as fast
	// as it can pop them out of the channel. In practice, this means that Release() will always return
	// very quickly, and will never block for more than a few microseconds.

	// Send the request to the control loop.

	request := &releaseRequest{
		tokens:      tokens,
		releaseChan: make(chan bool, 1),
	}

	select {
	case <-s.shutdownChan:
		// Close() has been called. It's ok to abort and return immediately, since releasing tokens no longer matters.
		return fmt.Errorf("semaphore closed")
	case s.releaseChan <- request:
		// Request was successfully sent to the control loop.
	}

	select {
	case <-s.shutdownChan:
		// Close() has been called. It's ok to abort and return immediately, since releasing tokens no longer matters.
		return fmt.Errorf("semaphore closed")
	case success := <-request.releaseChan:
		if success {
			return nil
		}
		return fmt.Errorf("failed to release %d tokens", tokens)
	}
}

// Close releases resources associated with the semaphore. If there are pending Acquire() or Release() calls, those
// methods may return an error as a result to this call to Close().
func (s *Semaphore) Close() {
	if s.closeCalled.CompareAndSwap(false, true) {
		close(s.shutdownChan)
	}
}

// controlLoop is the main loop that processes acquire and release requests. Operations are run in serialized order
// here to simplify threading logic/safety.
func (s *Semaphore) controlLoop() {
	for {
		if s.nextAcquireRequest == nil {
			// We don't currently have any pending acquire requests.
			select {
			case <-s.shutdownChan:
				// Close() has been called, exit the loop.
				return
			case req := <-s.releaseChan:
				// Important: handle release requests before acquire requests
				s.handleReleaseRequest(req)
			case req := <-s.acquireChan:
				// We have a new acquire request.
				s.nextAcquireRequest = req
				s.handleAcquireRequest()
			}
		} else {
			// There is currently a pending acquire request. We will only hit this block if we previously had a request
			// that we could not immediately fulfill. In that case, we need to wait for the next release request to
			// come in the door before we can process the pending request, or for the pending request to time out.

			select {
			case <-s.shutdownChan:
				// Close() has been called, exit the loop.
				return
			case req := <-s.releaseChan:
				// Important: handle release requests before acquire requests
				s.handleReleaseRequest(req)

				// make another attempt at processing the pending acquire request now that we've freed up some tokens
				s.handleAcquireRequest()
			case <-s.nextAcquireRequest.ctx.Done():
				// The context for the pending acquire request was canceled.
				s.nextAcquireRequest.acquiredChan <- false
				s.nextAcquireRequest = nil
			}
		}
	}
}

// handle a request to release tokens
func (s *Semaphore) handleReleaseRequest(req *releaseRequest) {
	if req.tokens > s.acquiredTokens {
		// error case: more tokens released than previously acquired
		req.releaseChan <- false
	} else {
		// success case: release the tokens
		s.acquiredTokens -= req.tokens
		req.releaseChan <- true
	}
}

// handle the pending acquire request stored in s.nextAcquireRequest.
func (s *Semaphore) handleAcquireRequest() {
	request := s.nextAcquireRequest

	remainingTokens := s.totalTokens - s.acquiredTokens
	if request.tokens <= remainingTokens {
		s.acquiredTokens += request.tokens
		request.acquiredChan <- true
		s.nextAcquireRequest = nil
	}

	// If we get here without acquiring the tokens, it means that we could not fulfill the request immediately.
	// We will make another attempt to process this request either when the request's context is canceled, or when
	// a release request comes in that frees up some tokens.
}
