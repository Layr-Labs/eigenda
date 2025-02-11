package disk

// address describes the location of data on disk.
// The first 4 bytes are the file ID, and the second 4 bytes are the offset of the data within the file.
type address uint64

// newKeyMap creates a new address
func newAddress(fileID uint32, offset uint32) address {
	return address(uint64(fileID)<<32 | uint64(offset))
}

// fileID returns the file ID of the value address.
func (a address) fileID() uint32 {
	return uint32(a >> 32)
}

// offset returns the offset of the value address.
func (a address) offset() uint32 {
	return uint32(a)
}
