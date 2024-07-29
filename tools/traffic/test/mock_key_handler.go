package test

import (
	"sync"
	"testing"
)

// mockKeyHandler is a stand-in for the blob verifier's UnconfirmedKeyHandler.
type mockKeyHandler struct {
	t *testing.T

	ProvidedKey      []byte
	ProvidedChecksum [16]byte
	ProvidedSize     uint

	// Incremented each time AddUnconfirmedKey is called.
	Count uint

	lock *sync.Mutex
}

func newMockKeyHandler(t *testing.T, lock *sync.Mutex) *mockKeyHandler {
	return &mockKeyHandler{
		t:    t,
		lock: lock,
	}
}

func (m *mockKeyHandler) AddUnconfirmedKey(key *[]byte, checksum *[16]byte, size uint) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.ProvidedKey = *key
	m.ProvidedChecksum = *checksum
	m.ProvidedSize = size

	m.Count++
}
