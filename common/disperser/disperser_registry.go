package disperser

import "context"

// DisperserRegistry provides access to disperser information from the DisperserRegistry contract.
type DisperserRegistry interface {
	// Returns the list of dispersers that network participants should interact with by default.
	GetDefaultDispersers(ctx context.Context) ([]uint32, error)
	// Returns the list of dispersers that support on-demand payments
	GetOnDemandDispersers(ctx context.Context) ([]uint32, error)
	// Returns the gRPC URI for a specific disperser in "hostname:port" format
	GetDisperserGrpcUri(ctx context.Context, disperserID uint32) (string, error)
}
