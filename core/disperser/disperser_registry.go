package disperser

import "context"

// DisperserRegistry provides access to disperser information from the DisperserRegistry contract.
type DisperserRegistry interface {
	// Returns the set of dispersers that network participants should interact with by default.
	GetDefaultDispersers(ctx context.Context) (map[uint32]struct{}, error)
	// Returns the set of dispersers that support on-demand payments
	GetOnDemandDispersers(ctx context.Context) (map[uint32]struct{}, error)
	// Returns the gRPC URI for a specific disperser in "hostname:port" format
	GetDisperserGrpcUri(ctx context.Context, disperserID uint32) (string, error)
}
