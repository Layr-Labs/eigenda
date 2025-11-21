package clients

import "context"

// DisperserConnectionInfo contains the information needed to connect to a disperser
type DisperserConnectionInfo struct {
	Hostname string
	Port     uint16
}

// DisperserRegistry provides access to disperser information from the DisperserRegistry contract.
type DisperserRegistry interface {
	// Returns the set of dispersers that network participants should interact with by default.
	GetDefaultDispersers(ctx context.Context) ([]uint32, error)
	// Returns the set of dispersers that support on-demand payments
	GetOnDemandDispersers(ctx context.Context) ([]uint32, error)
	// Returns the connection information for a specific disperser
	GetDisperserConnectionInfo(ctx context.Context, disperserID uint32) (*DisperserConnectionInfo, error)
}
