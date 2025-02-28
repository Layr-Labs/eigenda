package tablestore

import (
	"encoding/binary"
	"time"
)

// prependTimestamp prepends the given timestamp to the given base byte slice. The timestamp is
// stored as an 8-byte big-endian integer.
func prependTimestamp(
	timestamp time.Time,
	baseValue []byte) []byte {

	result := make([]byte, len(baseValue)+8)
	unixNano := timestamp.UnixNano()
	binary.BigEndian.PutUint64(result, uint64(unixNano))

	copy(result[8:], baseValue)

	return result
}

// parsePrependedTimestamp extracts the timestamp and base key from the given byte slice. This method
// is the inverse of prependTimestamp.
func parsePrependedTimestamp(data []byte) (timestamp time.Time, baseValue []byte) {
	expiryUnixNano := int64(binary.BigEndian.Uint64(data))
	timestamp = time.Unix(0, expiryUnixNano)
	baseValue = data[8:]
	return timestamp, baseValue
}
