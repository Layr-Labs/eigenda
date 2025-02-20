package segment

import "fmt"

// Address describes the location of data on disk.
// The first 4 bytes are the file ID, and the second 4 bytes are the offset of the data within the file.
type Address uint64

// NewAddress creates a new address
func NewAddress(index uint32, offset uint32) Address {
	return Address(uint64(index)<<32 | uint64(offset))
}

// Index returns the file index of the value address.
func (a Address) Index() uint32 {
	return uint32(a >> 32)
}

// Offset returns the offset of the value address.
func (a Address) Offset() uint32 {
	return uint32(a)
}

// String returns a string representation of the address.
func (a Address) String() string {
	return fmt.Sprintf("(%d:%d)", a.Index(), a.Offset())
}
