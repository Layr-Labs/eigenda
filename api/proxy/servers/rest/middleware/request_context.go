package middleware

import (
	"context"
	"net/http"
)

// withRequestContext initializes the request context (outermost middleware)
func withRequestContext(
	handleFn func(http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestContext := &RequestContext{
			// CertVersion is only known and set after parsing the request,
			// so we initialize it to a default value.
			// TODO: should this flow via some other means..?
			CertVersion: "unknown",
		}

		// Add context to request
		rWithRequestContext := r.WithContext(context.WithValue(r.Context(), RequestContextKey, requestContext))

		handleFn(w, rWithRequestContext)

		// RequestContext middleware is the outermost middleware,
		// so there is nothing to do after the handler is called.
	}
}

// RequestContext holds request-specific data that middlewares need to share
type RequestContext struct {
	CertVersion string
}

// ContextKey is used to store CertVersion in the request context
// A custom type is used to avoid collisions with other context keys.
// See https://pkg.go.dev/context#WithValue
type ContextKey string

const RequestContextKey ContextKey = "RequestContext"

// getRequestContext retrieves the RequestContext from the request
func getRequestContext(r *http.Request) *RequestContext {
	if ctx, ok := r.Context().Value(RequestContextKey).(*RequestContext); ok {
		return ctx
	}
	return nil
}

// SetCertVersion is public because it allows handlers to set the certificate version.
func SetCertVersion(r *http.Request, certVersion string) {
	if ctx := getRequestContext(r); ctx != nil {
		ctx.CertVersion = certVersion
	}
}

// getCertVersion is private because it is only used by the middlewares.
func getCertVersion(r *http.Request) string {
	if ctx := getRequestContext(r); ctx != nil {
		return ctx.CertVersion
	}
	return "unknown"
}
