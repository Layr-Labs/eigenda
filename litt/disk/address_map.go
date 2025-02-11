package disk

import (
	"fmt"
)

// addressMap manages a mapping between keys and the location of their data on disk (i.e. "addresses").
type addressMap struct {
	// addresses maps keys to their locations on disk.
	addresses map[string]address
}

// TODO thread safety

// setAddress sets the address of a key in a table.
func (k *addressMap) setAddress(key []byte, addr address) {
	k.addresses[string(key)] = addr
}

// getAddress gets the address of a key in a table.
func (k *addressMap) getAddress(key []byte) (address, error) {
	a, ok := k.addresses[string(key)]
	if !ok {
		return 0, fmt.Errorf("key not found in address map: %s", key)
	}
	return a, nil
}

// deleteAddress deletes a number of addresses from a table.
func (k *addressMap) deleteAddresses(keys [][]byte) {
	for _, key := range keys {
		delete(k.addresses, string(key))
	}
}
