package client

import "encoding/binary"

func ConvertIntToVarUInt(v int) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, uint64(v))
	return buf[:n]
}
