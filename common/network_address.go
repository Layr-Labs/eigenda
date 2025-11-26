package common

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

// Represents a network address composed of a hostname and port.
type NetworkAddress struct {
	hostname string
	port     uint16
}

// Creates a new NetworkAddress with the given hostname and port.
// Returns an error if the hostname is empty, or the port is outside the valid range (1-65535).
func NewNetworkAddress[T constraints.Integer](hostname string, genericPort T) (*NetworkAddress, error) {
	if strings.TrimSpace(hostname) == "" {
		return nil, fmt.Errorf("hostname must not be empty")
	}
	port := uint64(genericPort)
	if port == 0 || port > 65535 {
		return nil, fmt.Errorf("port must be in range 1-65535, got %d", genericPort)
	}
	return &NetworkAddress{
		hostname: hostname,
		port:     uint16(port),
	}, nil
}

// Creates a new NetworkAddress by parsing a string in "hostname:port" format.
// Returns an error if the format is invalid or the port is outside the valid range (1-65535).
func NewNetworkAddressFromString(address string) (*NetworkAddress, error) {
	hostname, portStr, err := net.SplitHostPort(address)
	if err != nil {
		return nil, fmt.Errorf("invalid address format %q: %w", address, err)
	}

	port, err := strconv.ParseUint(portStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid port %q: %w", portStr, err)
	}

	return NewNetworkAddress(hostname, port)
}

// Returns the hostname.
func (n *NetworkAddress) Hostname() string {
	return n.hostname
}

// Returns the port.
func (n *NetworkAddress) Port() uint16 {
	return n.port
}

// Returns true if this network address is equal to the other network address.
func (n *NetworkAddress) Equals(other *NetworkAddress) bool {
	if n == nil || other == nil {
		return false
	}
	return n.hostname == other.hostname && n.port == other.port
}

// Validates the network address.
func (n *NetworkAddress) Check() error {
	if n == nil {
		return fmt.Errorf("network address is nil")
	}
	if n.hostname == "" {
		return fmt.Errorf("hostname is empty")
	}
	if n.port == 0 {
		return fmt.Errorf("port is zero")
	}
	return nil
}

// Returns the network address in the format "hostname:port".
func (n *NetworkAddress) String() string {
	return net.JoinHostPort(n.hostname, fmt.Sprintf("%d", n.port))
}
