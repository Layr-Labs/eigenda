package server

import (
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/commitments"
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
