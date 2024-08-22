package workers

import (
	"github.com/stretchr/testify/mock"
)

var _ KeyHandler = (*MockKeyHandler)(nil)

// MockKeyHandler is a stand-in for the blob verifier's UnconfirmedKeyHandler.
type MockKeyHandler struct {
	mock mock.Mock

	ProvidedKey      []byte
	ProvidedChecksum [16]byte
	ProvidedSize     uint
}

func (m *MockKeyHandler) AddUnconfirmedKey(key []byte, checksum [16]byte, size uint) {
	m.mock.Called(key, checksum, size)

	m.ProvidedKey = key
	m.ProvidedChecksum = checksum
	m.ProvidedSize = size
}
