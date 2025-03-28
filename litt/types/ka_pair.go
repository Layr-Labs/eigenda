package types

// KAPair represents a key-address pair.
type KAPair struct {
	// Key is the key.
	Key []byte
	// Address is the address that describes where the value associated with the key is stored.
	Address Address
}
