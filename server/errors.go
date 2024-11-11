package server

import (
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/commitments"
	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MetaError includes both an error and commitment metadata
type MetaError struct {
	Err  error
	Meta commitments.CommitmentMeta
}

func (me MetaError) Error() string {
	return fmt.Sprintf("Error: %s (Mode: %s, CertVersion: %b)",
		me.Err.Error(),
		me.Meta.Mode,
		me.Meta.CertVersion)
}

func (me MetaError) Unwrap() error {
	return me.Err
}

func is400(err error) bool {
	// proxy requests are super simple (clients basically only pass bytes), so the only 400 possible
	// is passing a blob that's too big.
	//
	// Any 400s returned by the disperser are due to formatting bugs in proxy code, for eg. badly
	// IFFT'ing or encoding the blob, so we shouldn't return a 400 to the client.
	// See https://github.com/Layr-Labs/eigenda/blob/bee55ed9207f16153c3fd8ebf73c219e68685def/api/errors.go#L22
	// for the 400s returned by the disperser server (currently only INVALID_ARGUMENT).
	return errors.Is(err, common.ErrProxyOversizedBlob)
}

func is429(err error) bool {
	// grpc RESOURCE_EXHAUSTED is returned by the disperser server when the client has sent too many requests
	// in a short period of time. This is a client-side issue, so we should return the 429 to the client.
	st, isGRPCError := status.FromError(err)
	return isGRPCError && st.Code() == codes.ResourceExhausted
}

// 503 is returned to tell the caller (batcher) to failover to ethda b/c eigenda is temporarily down
func is503(err error) bool {
	// TODO: would be cleaner to define a sentinel error in eigenda-core and use that instead
	return errors.Is(err, &api.ErrorFailover{})
}
