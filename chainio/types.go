package chainio

import "encoding/hex"

type QuorumID = uint8

type OperatorID = [32]byte

func GetOperatorHex(id OperatorID) string {
	return hex.EncodeToString(id[:])
}
