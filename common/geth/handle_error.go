package geth

import (
	"errors"

	"github.com/ethereum/go-ethereum/rpc"
)

// handleHttpError returns a boolean indicating if error atrributes to remote RPC
func (f *RPCStatistics) handleHttpError(httpRespError rpc.HTTPError) bool {
	sc := httpRespError.StatusCode
	f.Logger.Info("[HTTP Response Error]", "Status Code", sc)
	if sc >= 200 && sc < 300 {
		// 2xx error
		return false
	}

	if sc >= 400 && sc < 500 {
		// 4xx error, Client Error. Alchemy documents 400,401,403,429
		// https://docs.alchemy.com/reference/error-reference
		return false
	}

	if sc >= 500 {
		// 5xx codes, Server Error, Alchemy documents 500, 503
		return true
	}

	// by default, attribute to rpc
	return true
}

// handleJsonRPCError returns a boolean indicating if error atrributes to remote Server
func (f *RPCStatistics) handleJsonRPCError(rpcError rpc.Error) bool {
	ec := rpcError.ErrorCode()

	// Based on JSON-RPC 2.0, https://www.jsonrpc.org/specification#error_object
	// Parse Error, Invalid Request,Method not found,Invalid params,Internal error
	if ec == -32700 || ec == -32600 || ec == -32601 || ec == -32602 || ec == -32603 {
		return false
	}

	// server error
	if ec >= -32099 && ec <= -32000 {
		return true
	}

	// execution revert, see https://docs.alchemy.com/reference/error-reference
	if ec == 3 {
		return false
	}

	return true
}

// handleError returns a boolean indicating if error atrributes to remote Server
func (f *RPCStatistics) handleError(err error) bool {

	var httpRespError rpc.HTTPError
	if errors.As(err, &httpRespError) {
		// if error is http error
		return f.handleHttpError(httpRespError)
	} else {
		// it might be websocket error or ipc error. Parse json error code
		var nonHttpRespError rpc.Error
		if errors.As(err, &nonHttpRespError) {
			return f.handleJsonRPCError(nonHttpRespError)
		}

		// If no http response is returned, it is a connection issue,
		// since we can't accurately attribute the network issue to neither sender nor receiver
		// side. Optimistically, switch rpc client
		return true
	}

}
