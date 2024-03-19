package geth

import (
	"errors"

	"github.com/ethereum/go-ethereum/rpc"
)

// handleHttpError returns a boolean indicating if the current RPC should be rotated
func (f *RPCStatistics) handleHttpError(httpRespError rpc.HTTPError) bool {
	sc := httpRespError.StatusCode
	if sc >= 200 && sc < 300 {
		// 2xx error
		return false
	}
	// Default to rotation the current leader, because it allows a higher chance to get the query completed.
	// When a http query failed for non2xx, either sender or remote server is at fault or both.
	// Since the software cannot be patched at runtime, there is no immediate remedy when sender is at fault.
	// If the fault is at receiver for any reason, defaulting to retry increase the chance of
	// having a successful return.
	f.Logger.Info("[HTTP Response Error]", "Status Code", sc, "Error", httpRespError)
	return true
}

// handleJsonRPCError returns a boolean indicating if the current RPC should be rotated
// It could be a http2xx error, websocket error or ipc error
func (f *RPCStatistics) handleJsonRPCError(rpcError rpc.Error) bool {
	ec := rpcError.ErrorCode()
	f.Logger.Info("[JSON RPC Response Error]", "Error Code", ec, "Error", rpcError)

	// Based on JSON-RPC 2.0, https://www.jsonrpc.org/specification#error_object.
	// Rotating in those case, because the software wants to operate in optimistic case.
	// The same reason as above
	if ec >= -32768 && ec <= -32000 {
		return true
	}

	// default to false, because http 2xx error can return any error code, https://docs.alchemy.com/reference/error-reference
	//
	// Todo, it is useful to have local information to distinguish if the current connection is a http connection.
	// It is useful
	return false
}

// handleError returns a boolean indicating if the current connection should be rotated.
// Because library of the sender uses geth, which supports only 3 types of connections,
// we can categorize the error as HTTP error, Websocket error and IPC error.
//
// If the error is http, non2xx error would generate HTTP error, https://github.com/ethereum/go-ethereum/blob/master/rpc/http.go#L233
// but a 2xx http response could contain JSON RPC error, https://github.com/ethereum/go-ethereum/blob/master/rpc/http.go#L181
// If the error is Websocket or IPC, we only look for JSON error, https://github.com/ethereum/go-ethereum/blob/master/rpc/json.go#L67

func (f *RPCStatistics) handleError(err error) bool {

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
			return f.handleJsonRPCError(rpcError)
		}

		// If no http response or no rpc response is returned, it is a connection issue,
		// since we can't accurately attribute the network issue to neither sender nor receiver
		// side. Optimistically, switch rpc client
		f.Logger.Info("[Default Response Error]", err)
		return true
	}

}
