package middleware

import "net/http"

// Used to capture the status code of the response, so that we can use it in metrics
// and logging middlewares. See https://github.com/golang/go/issues/18997
// For most routes, the status is written by the error middleware.
// We could potentially instead just return the status code from the error middleware
// to the outer layer middlewares. Not sure which way is better.
//
// TODO: right now instantiating a separate scw for logging and metrics... is there a better way?
// TODO: should we capture more information about the response, like GET vs POST, etc?
type statusCaptureWriter struct {
	http.ResponseWriter
	status int
}

func (scw *statusCaptureWriter) WriteHeader(status int) {
	scw.status = status
	scw.ResponseWriter.WriteHeader(status)
}

func newStatusCaptureWriter(w http.ResponseWriter) *statusCaptureWriter {
	return &statusCaptureWriter{
		ResponseWriter: w,
		// 200 status code is only added to response by outer layer http framework,
		// since WriteHeader(200) is typically not called by handlers.
		// So we initialize status as 200, and assume that any other status code
		// will be set by the handler.
		status: http.StatusOK,
	}
}
