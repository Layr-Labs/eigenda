package workers

// KeyHandler is an interface describing an object that can accept unconfirmed keys.
type KeyHandler interface {
	// AddUnconfirmedKey accepts an unconfirmed blob key, the checksum of the blob, and the size of the blob in bytes.
	AddUnconfirmedKey(key []byte, checksum [16]byte, size uint)
}
