package geth

import (
	"errors"

	"github.com/ethereum/go-ethereum/rpc"
)

type ImmediateAction int

const (
	Return ImmediateAction = iota
	Retry
)

type NextEndpoint int

const (
	NewRPC = iota
	CurrentRPC
)

// handleHttpError returns a boolean indicating if the current RPC should be rotated
// the second boolean indicating if should giveup immediately
func (f *RPCStatistics) handleHttpError(httpRespError rpc.HTTPError) (NextEndpoint, ImmediateAction) {
	sc := httpRespError.StatusCode
	// Default to rotation the current RPC, because it allows a higher chance to get the query completed.
	f.Logger.Info("[HTTP Response Error]", "Status Code", sc, "Error", httpRespError)

	if sc >= 200 && sc < 300 {
		// 2xx error, however it should not be reachable
		return CurrentRPC, Return
	}

	if sc >= 400 && sc < 500 {
		// 403 Forbidden, 429 Too many Requests. We should rotate
		if sc == 403 || sc == 429 {
			return NewRPC, Retry
		}
		return CurrentRPC, Retry
	}

	// 500
	return NewRPC, Retry
}

// handleError returns a boolean indicating if the current connection should be rotated.
// Because library of the sender uses geth, which supports only 3 types of connections,
// we can categorize the error as HTTP error, Websocket error and IPC error.
//
// If the error is http, non2xx error would generate HTTP error, https://github.com/ethereum/go-ethereum/blob/master/rpc/http.go#L233
// but a 2xx http response could contain JSON RPC error, https://github.com/ethereum/go-ethereum/blob/master/rpc/http.go#L181
// If the error is Websocket or IPC, we only look for JSON error, https://github.com/ethereum/go-ethereum/blob/master/rpc/json.go#L67

func (f *RPCStatistics) handleError(err error) (NextEndpoint, ImmediateAction) {

	var httpRespError rpc.HTTPError
	if errors.As(err, &httpRespError) {
		// if error is http error, i.e. non 2xx error, it is handled here
		// if it is 2xx error, the error message is nil, https://github.com/ethereum/go-ethereum/blob/master/rpc/http.go,
		// execution does not entere here.
		return f.handleHttpError(httpRespError)
	} else {
		// it might be http2xx error, websocket error or ipc error. Parse json error code
		var rpcError rpc.Error
		if errors.As(err, &rpcError) {
			ec := rpcError.ErrorCode()
			f.Logger.Info("[JSON RPC Response Error]", "Error Code", ec, "Error", rpcError)
			// we always attribute JSON RPC error as sender's fault, i.e no connection rotation
			return CurrentRPC, Return
		}

		// If no http response or no rpc response is returned, it is a connection issue,
		// since we can't accurately attribute the network issue to neither sender nor receiver
		// side. Optimistically, switch rpc client
		f.Logger.Info("[Default Response Error]", err)
		return NewRPC, Retry
	}

}
