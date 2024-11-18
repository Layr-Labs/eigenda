package authentication

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/core"
	"sync"
	"time"
)

// RequestAuthenticator authenticates requests to the relay service. This object is thread safe.
type RequestAuthenticator interface {
	// AuthenticateGetChunksRequest authenticates a GetChunksRequest, returning an error if the request is invalid.
	// The address is the address of the peer that sent the request. This may be used to cache authentication results
	// in order to save server resources.
	AuthenticateGetChunksRequest(
		address string,
		request *pb.GetChunksRequest,
		now time.Time) error
}

// authenticationTimeout is used to track the expiration of an authentication.
type authenticationTimeout struct {
	clientID   string
	expiration time.Time
}

var _ RequestAuthenticator = &requestAuthenticator{}

type requestAuthenticator struct {
	ics core.IndexedChainState

	// authenticatedClients is a set of client IDs that have been recently authenticated.
	authenticatedClients map[string]struct{}

	// authenticationTimeouts is a list of authentications that have been performed, along with their expiration times.
	authenticationTimeouts []*authenticationTimeout

	// authenticationTimeoutDuration is the duration for which an authentication is valid.
	// If this is zero, then authentication saving is disabled, and each request will be authenticated independently.
	authenticationTimeoutDuration time.Duration

	// savedAuthLock is used for thread safe atomic modification of the authenticatedClients map and the
	// authenticationTimeouts queue.
	savedAuthLock sync.Mutex
}

// NewRequestAuthenticator creates a new RequestAuthenticator.
func NewRequestAuthenticator(
	ics core.IndexedChainState,
	authenticationTimeoutDuration time.Duration) RequestAuthenticator {

	return &requestAuthenticator{
		ics:                           ics,
		authenticatedClients:          make(map[string]struct{}),
		authenticationTimeouts:        make([]*authenticationTimeout, 0),
		authenticationTimeoutDuration: authenticationTimeoutDuration,
	}
}

func (a *requestAuthenticator) AuthenticateGetChunksRequest(
	address string,
	request *pb.GetChunksRequest,
	now time.Time) error {

	if a == nil {
		// do not enforce authentication if the authenticator is nil
		return nil
	}

	if a.isAuthenticationStillValid(now, address) {
		// We've recently authenticated this client. Do not authenticate again for a while.
		return nil
	}

	blockNumber, err := a.ics.GetCurrentBlockNumber()
	if err != nil {
		return fmt.Errorf("failed to get current block number: %w", err)
	}
	operators, err := a.ics.GetIndexedOperators(context.Background(), blockNumber)
	if err != nil {
		return fmt.Errorf("failed to get operators: %w", err)
	}

	operatorID := core.OperatorID(request.RequesterId)
	operator, ok := operators[operatorID]
	if !ok {
		return errors.New("operator not found")
	}
	key := operator.PubkeyG2

	g1Point, err := (&core.G1Point{}).Deserialize(request.RequesterSignature)
	if err != nil {
		return fmt.Errorf("failed to deserialize signature: %w", err)
	}

	signature := core.Signature{
		G1Point: g1Point,
	}

	hash := HashGetChunksRequest(request)
	isValid := signature.Verify(key, ([32]byte)(hash))

	if !isValid {
		return errors.New("signature verification failed")
	}

	a.saveAuthenticationResult(now, address)
	return nil
}

// saveAuthenticationResult saves the result of an authentication.
func (a *requestAuthenticator) saveAuthenticationResult(now time.Time, address string) {
	if a.authenticationTimeoutDuration == 0 {
		// Authentication saving is disabled.
		return
	}

	a.savedAuthLock.Lock()
	defer a.savedAuthLock.Unlock()

	a.authenticatedClients[address] = struct{}{}
	a.authenticationTimeouts = append(a.authenticationTimeouts,
		&authenticationTimeout{
			clientID:   address,
			expiration: now.Add(a.authenticationTimeoutDuration),
		})
}

// isAuthenticationStillValid returns true if the client at the given address has been authenticated recently.
func (a *requestAuthenticator) isAuthenticationStillValid(now time.Time, address string) bool {
	if a.authenticationTimeoutDuration == 0 {
		// Authentication saving is disabled.
		return false
	}

	a.savedAuthLock.Lock()
	defer a.savedAuthLock.Unlock()

	a.removeOldAuthentications(now)
	_, ok := a.authenticatedClients[address]
	return ok
}

// removeOldAuthentications removes any authentications that have expired.
func (a *requestAuthenticator) removeOldAuthentications(now time.Time) {
	index := 0
	for ; index < len(a.authenticationTimeouts); index++ {
		if a.authenticationTimeouts[index].expiration.After(now) {
			break
		}
		delete(a.authenticatedClients, a.authenticationTimeouts[index].clientID)
	}
	if index > 0 {
		a.authenticationTimeouts = a.authenticationTimeouts[index:]
	}
}
