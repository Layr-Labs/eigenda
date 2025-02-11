package littdb

// address describes the location of data on disk.
// The first 4 bytes are the file ID, and the second 4 bytes are the offset of the data within the file.
type address uint64

// table is a map of keys within a table to their location on disk. The key in the map is a byte array stored in
// string format (golang maps don't support byte arrays as keys).
type table map[string]address

// keyMap manages a mapping between keys and the location of their data on disk.
type keyMap struct {
	// tables is a map of table names to their key maps.
	tables map[string]table
}

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

// setAddress sets the address of a key in a table.
func (k keyMap) setAddress(table string, key []byte, addr address) {
	// TODO
}

// getAddress gets the address of a key in a table.
func (k keyMap) getAddress(table string, key []byte) (address, error) {
	// TODO
	return 0, nil
}

// deleteAddress deletes a number of addresses from a table.
func (k keyMap) deleteAddresses(table string, keys [][]byte) {
	// TODO
}
