package geth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeRpcUrl(t *testing.T) {
	require.Equal(t, "https://rpc.example.com", SanitizeRpcUrl("https://user:password@rpc.example.com"))
	require.Equal(t, "https://rpc.example.com", SanitizeRpcUrl("https://rpc.example.com/v2/SECRET_API_KEY"))
	require.Equal(t, "https://rpc.example.com", SanitizeRpcUrl("https://rpc.example.com?apikey=SECRET"))
	require.Equal(t, "https://rpc.example.com", SanitizeRpcUrl("https://SECRET_KEY@rpc.example.com"))
	require.Equal(t, "wss://rpc.example.com", SanitizeRpcUrl("wss://SECRET@rpc.example.com/ws"))
	require.Equal(t, "[malformed-url]", SanitizeRpcUrl("user:pass@example.com"))
	require.Equal(t, "[invalid-url]", SanitizeRpcUrl("://invalid"))
}
