package workers

import (
	"testing"
)

var _ KeyHandler = (*MockKeyHandler)(nil)

// MockKeyHandler is a stand-in for the blob verifier's UnconfirmedKeyHandler.
type MockKeyHandler struct {
	t *testing.T

	ProvidedKey      []byte
	ProvidedChecksum [16]byte
	ProvidedSize     uint

	// Incremented each time AddUnconfirmedKey is called.
	Count uint
}

func NewMockKeyHandler(t *testing.T) *MockKeyHandler {
	return &MockKeyHandler{
		t: t,
	}
}

func (m *MockKeyHandler) AddUnconfirmedKey(key []byte, checksum [16]byte, size uint) {
	m.ProvidedKey = key
	m.ProvidedChecksum = checksum
	m.ProvidedSize = size

	m.Count++
}
