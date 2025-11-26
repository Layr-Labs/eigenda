package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewNetworkAddress(t *testing.T) {
	addr, err := NewNetworkAddress("example.com", 8080)
	require.NoError(t, err)
	require.Equal(t, "example.com", addr.Hostname())
	require.Equal(t, uint16(8080), addr.Port())
	require.Equal(t, "example.com:8080", addr.String())
}

func TestNewNetworkAddressWithDifferentIntegerTypes(t *testing.T) {
	var portInt int = 443
	var portUint16 uint16 = 443
	var portUint64 uint64 = 443

	addr1, err := NewNetworkAddress("host1.com", portInt)
	require.NoError(t, err)
	require.Equal(t, uint16(443), addr1.Port())

	addr2, err := NewNetworkAddress("host2.com", portUint16)
	require.NoError(t, err)
	require.Equal(t, uint16(443), addr2.Port())

	addr3, err := NewNetworkAddress("host3.com", portUint64)
	require.NoError(t, err)
	require.Equal(t, uint16(443), addr3.Port())
}

func TestNewNetworkAddressValidation(t *testing.T) {
	_, err := NewNetworkAddress("", 8080)
	require.Error(t, err)
	_, err = NewNetworkAddress("   ", 8080)
	require.Error(t, err)
	_, err = NewNetworkAddress("example.com", 0)
	require.Error(t, err)
	_, err = NewNetworkAddress("example.com", -1)
	require.Error(t, err)
	_, err = NewNetworkAddress("example.com", 65536)
	require.Error(t, err)
}

func TestNewNetworkAddressFromString(t *testing.T) {
	addr, err := NewNetworkAddressFromString("example.com:8080")
	require.NoError(t, err)
	require.Equal(t, "example.com", addr.Hostname())
	require.Equal(t, uint16(8080), addr.Port())

	_, err = NewNetworkAddressFromString("example.com")
	require.Error(t, err)
	_, err = NewNetworkAddressFromString("example.com:notaport")
	require.Error(t, err)
	_, err = NewNetworkAddressFromString("example.com:99999")
	require.Error(t, err)
}

func TestNetworkAddressEquals(t *testing.T) {
	addr1, _ := NewNetworkAddress("example.com", 8080)
	addr2, _ := NewNetworkAddress("example.com", 8080)
	addr3, _ := NewNetworkAddress("other.com", 8080)
	addr4, _ := NewNetworkAddress("example.com", 9090)

	require.True(t, addr1.Equals(addr2))
	require.False(t, addr1.Equals(addr3))
	require.False(t, addr1.Equals(addr4))
	require.False(t, addr1.Equals(nil))

	var nilAddr *NetworkAddress
	require.False(t, nilAddr.Equals(addr1))
	require.False(t, nilAddr.Equals(nil))
}
