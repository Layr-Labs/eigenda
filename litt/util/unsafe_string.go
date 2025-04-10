package util

import "unsafe"

// UnsafeBytesToString converts a byte slice to a string without copying the data.
// Note that once converted in this way, it is not safe to modify the byte slice for any reason.
func UnsafeBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
