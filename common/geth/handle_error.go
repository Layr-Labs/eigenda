package geth

import (
	"errors"

	"github.com/ethereum/go-ethereum/rpc"
)

// handleHttpError returns a boolean indicating if error atrributes to remote RPC
func (f *FailoverController) handleHttpError(httpRespError rpc.HTTPError) bool {
	sc := httpRespError.StatusCode
	f.Logger.Info("[RPC Error]", "Status Code", sc)
	if sc >= 200 && sc < 300 {
		// 2xx error
		return false
	} else if sc >= 300 && sc < 400 {
		// 4xx error, Client Error. Alchemy documents 400,401,403,429
		// https://docs.alchemy.com/reference/error-reference
		return false
	} else if sc >= 500 {
		// 5xx codes, Server Error, Alchemy documents 500, 503
		return true
	}
	// by default, attribute to rpc
	return true
}

// handleError returns a boolean indicating if error atrributes to remote RPC
func (f *FailoverController) handleError(err error) bool {
	var httpRespError rpc.HTTPError
	if errors.As(err, &httpRespError) {
		// if error is http error
		return f.handleHttpError(httpRespError)
	} else {
		// If no http response is returned, it is a connection issue,
		// since we can't accurately attribute the network issue to neither sender nor receiver
		// side. Optimistically, switch rpc client
		return true
	}

}
