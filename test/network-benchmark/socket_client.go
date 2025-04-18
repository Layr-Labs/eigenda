package network_benchmark

import (
	"encoding/binary"
	"io"
	"net"
	"time"
)

// socketClient implements TestClient interface for direct socket communication
type socketClient struct {
	conn    net.Conn
	timeout time.Duration
}

// NewSocketClient creates a new socketClient instance
func NewSocketClient(serverAddress string) (TestClient, error) {
	// Connect to the server
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}

	return &socketClient{
		conn:    conn,
		timeout: 30 * time.Second, // Default timeout
	}, nil
}

// GetData retrieves data of the specified size from the server
func (c *socketClient) GetData(size int64, seed int64) ([]byte, error) {
	// Set a timeout for the operation
	err := c.conn.SetDeadline(time.Now().Add(c.timeout))
	if err != nil {
		return nil, err
	}

	// Request format:
	// - 8 bytes (int64): size
	// - 8 bytes (int64): seed

	// Write size
	err = binary.Write(c.conn, binary.BigEndian, size)
	if err != nil {
		return nil, err
	}

	// Write seed
	err = binary.Write(c.conn, binary.BigEndian, seed)
	if err != nil {
		return nil, err
	}

	// Response format:
	// - 8 bytes (int64): data length
	// - N bytes: data

	// Read data length
	var dataLen int64
	err = binary.Read(c.conn, binary.BigEndian, &dataLen)
	if err != nil {
		return nil, err
	}

	// Read the actual data
	data := make([]byte, dataLen)
	_, err = io.ReadFull(c.conn, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// SetTimeout sets the timeout for client operations
func (c *socketClient) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// Close closes the connection to the server
func (c *socketClient) Close() error {
	return c.conn.Close()
}