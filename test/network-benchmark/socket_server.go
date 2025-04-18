package network_benchmark

import (
	"encoding/binary"
	"io"
	"net"
	"sync"
)

// socketServer implements TestServer interface for direct socket communication
type socketServer struct {
	listener   net.Listener
	randomData *reusableRandomness
	mu         sync.RWMutex
	clients    []net.Conn
	done       chan struct{}
}

// NewSocketServer creates a new socketServer instance
func NewSocketServer(address string) (TestServer, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	server := &socketServer{
		listener: listener,
		clients:  make([]net.Conn, 0),
		done:     make(chan struct{}),
	}

	// Start accepting connections in a goroutine
	go server.acceptConnections()

	return server, nil
}

// acceptConnections handles incoming client connections
func (s *socketServer) acceptConnections() {
	for {
		select {
		case <-s.done:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				// If the listener was closed, just exit
				select {
				case <-s.done:
					return
				default:
					continue
				}
			}

			// Add client to the list
			s.mu.Lock()
			s.clients = append(s.clients, conn)
			s.mu.Unlock()

			// Handle client in a goroutine
			go s.handleClient(conn)
		}
	}
}

// handleClient processes requests from a connected client
func (s *socketServer) handleClient(conn net.Conn) {
	defer func() {
		conn.Close()
		s.removeClient(conn)
	}()

	for {
		select {
		case <-s.done:
			return
		default:
			// Request format:
			// - 8 bytes (int64): size
			// - 8 bytes (int64): seed

			// Read request size (8 bytes)
			var size int64
			err := binary.Read(conn, binary.BigEndian, &size)
			if err != nil {
				if err == io.EOF || isConnectionClosed(err) {
					return // Client disconnected
				}
				continue // Other error, try again
			}

			// Read seed (8 bytes)
			var seed int64
			err = binary.Read(conn, binary.BigEndian, &seed)
			if err != nil {
				return // Error reading seed
			}

			// Generate data
			s.mu.RLock()
			data := s.randomData.getData(size, seed)
			s.mu.RUnlock()

			// Response format:
			// - 8 bytes (int64): data length
			// - N bytes: data

			// Write data length
			dataLen := int64(len(data))
			err = binary.Write(conn, binary.BigEndian, dataLen)
			if err != nil {
				return // Error writing data length
			}

			// Write actual data
			_, err = conn.Write(data)
			if err != nil {
				return // Error writing data
			}
		}
	}
}

// removeClient removes a client from the clients list
func (s *socketServer) removeClient(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, client := range s.clients {
		if client == conn {
			// Remove this client
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
			break
		}
	}
}

// SetRandomData sets the random data source for the server
func (s *socketServer) SetRandomData(randomData *reusableRandomness) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.randomData = randomData
}

// Close shuts down the server
func (s *socketServer) Close() error {
	close(s.done)

	// Close listener
	err := s.listener.Close()

	// Close all client connections
	s.mu.Lock()
	for _, client := range s.clients {
		client.Close()
	}
	s.clients = nil
	s.mu.Unlock()

	return err
}

// isConnectionClosed checks if an error indicates a closed connection
func isConnectionClosed(err error) bool {
	if err == nil {
		return false
	}

	// Check for common connection closed error strings
	errStr := err.Error()
	return errStr == "EOF" ||
		errStr == "use of closed network connection" ||
		errStr == "connection reset by peer"
}
