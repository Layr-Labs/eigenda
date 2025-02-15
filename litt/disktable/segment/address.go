package segment

// Address describes the location of data on disk.
// The first 4 bytes are the file ID, and the second 4 bytes are the offset of the data within the file.
type Address uint64

// NewAddress creates a new address
func NewAddress(fileID uint32, offset uint32) Address {
	return Address(uint64(fileID)<<32 | uint64(offset))
}

// FileID returns the file ID of the value address.
func (a Address) FileID() uint32 {
	return uint32(a >> 32)
}

// Offset returns the offset of the value address.
func (a Address) Offset() uint32 {
	return uint32(a)
}
