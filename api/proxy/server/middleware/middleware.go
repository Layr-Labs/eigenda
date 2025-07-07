package middleware

import (
	"net/http"

	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// Helper function to chain middlewares in the correct order
// Context -> Logging -> Metrics -> Error Handling -> Handler
//
// This should only be used for cert POST and GET routes,
// as the middlewares are currently not compatible with
// other generic routes (e.g. /health, /version, etc.)
//
// TODO: make our middlewares compatible with all routes, if possible.
func WithCertMiddlewares(
	handler func(http.ResponseWriter, *http.Request) error,
	log logging.Logger,
	m metrics.Metricer,
	mode commitments.CommitmentMode,
) http.HandlerFunc {
	return withRequestContext(
		withLogging(
			withMetrics(
				withErrorHandling(handler),
				m,
				mode,
			),
			log,
			mode,
		),
	)
}
