package proxyerrors

import (
	"errors"

	"github.com/Layr-Labs/eigenda/api"
)

// 503 is returned to tell the caller (batcher) to failover to ethda b/c eigenda is temporarily down
func Is503(err error) bool {
	// TODO: would be cleaner to define a sentinel error in eigenda-core and use that instead
	return errors.Is(err, &api.ErrorFailover{})
}
